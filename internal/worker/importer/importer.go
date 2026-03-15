package importer

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"

	"github.com/dukerupert/skalkaho/internal/repository"
)

// Importer orchestrates the full material import pipeline.
type Importer struct {
	queries *repository.Queries
	db      *sql.DB
	logger  *slog.Logger
}

// New creates a new Importer.
func New(queries *repository.Queries, db *sql.DB, logger *slog.Logger) *Importer {
	return &Importer{
		queries: queries,
		db:      db,
		logger:  logger,
	}
}

// Handle processes a single import job. This is the function registered with the worker.
func (imp *Importer) Handle(ctx context.Context, job repository.MaterialImportJob) error {
	defer func() {
		// Always try to clean up the temp file
		_ = os.Remove(job.TempFilePath)
	}()

	// 1. Parse Excel
	rows, err := ParseExcel(job.TempFilePath)
	if err != nil {
		return fmt.Errorf("parsing excel: %w", err)
	}

	if len(rows) == 0 {
		return fmt.Errorf("no data rows found in spreadsheet")
	}

	// 2. Update total count
	total := int32(len(rows))
	_, err = imp.queries.UpdateImportJobProgress(ctx, repository.UpdateImportJobProgressParams{
		ID:        job.ID,
		RowsTotal: sql.NullInt32{Int32: total, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("updating row total: %w", err)
	}

	// 3. Load existing materials for this supplier
	existing, err := imp.queries.ListMaterialsBySupplier(ctx, job.SupplierID)
	if err != nil {
		return fmt.Errorf("loading existing materials: %w", err)
	}

	// Build lookup by supplier_code
	byCode := make(map[string]repository.Material)
	for _, m := range existing {
		if m.SupplierCode.Valid && m.SupplierCode.String != "" {
			byCode[m.SupplierCode.String] = m
		}
	}

	// 4. Split rows into tiers
	var tier1, tier2, tier3 []RawRow
	for _, row := range rows {
		if row.SupplierCode == "" {
			tier3 = append(tier3, row)
			continue
		}
		if _, exists := byCode[row.SupplierCode]; exists {
			tier1 = append(tier1, row)
		} else {
			tier2 = append(tier2, row)
		}
	}

	imp.logger.Info("import tiers",
		"job_id", job.ID,
		"tier1_update", len(tier1),
		"tier2_create", len(tier2),
		"tier3_match", len(tier3),
	)

	// 5-7. Process all tiers in a transaction
	var created, updated, skipped, flagged, failed int32

	tx, err := imp.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	qtx := imp.queries.WithTx(tx)

	// Tier 1: supplier_code match — update price
	for _, row := range tier1 {
		m := byCode[row.SupplierCode]
		if m.PriceSource == "Manual" {
			skipped++
			continue
		}
		_, err := qtx.UpdateMaterial(ctx, repository.UpdateMaterialParams{
			ID:           m.ID,
			Name:         m.Name,
			SupplierID:   m.SupplierID,
			UnitPrice:    row.Price,
			Unit:         m.Unit,
			SupplierCode: m.SupplierCode,
			PriceSource:  "Supplier",
		})
		if err != nil {
			imp.logger.Error("tier1 update failed", "row", row.RowIndex, "error", err)
			failed++
			continue
		}
		updated++
	}

	// Tier 2: supplier_code present, no match — normalize names, then create new materials
	var tier2Names map[int]string
	if len(tier2) > 0 {
		names, err := NormalizeNames(ctx, tier2)
		if err != nil {
			imp.logger.Error("tier2 name normalization failed, using raw names", "error", err)
		} else {
			tier2Names = names
		}
	}

	for _, row := range tier2 {
		unit := row.Unit
		if unit == "" {
			unit = "ea"
		}
		name := row.Name
		if tier2Names != nil {
			if canonical, ok := tier2Names[row.RowIndex]; ok {
				name = canonical
			}
		}
		id := uuid.New().String()[:20]
		_, err := qtx.CreateMaterial(ctx, repository.CreateMaterialParams{
			ID:           id,
			Name:         name,
			SupplierID:   job.SupplierID,
			UnitPrice:    row.Price,
			Unit:         unit,
			SupplierCode: sql.NullString{String: row.SupplierCode, Valid: true},
			PriceSource:  "Supplier",
		})
		if err != nil {
			imp.logger.Error("tier2 create failed", "row", row.RowIndex, "error", err)
			failed++
			continue
		}
		created++
	}

	// Tier 3: no supplier_code — use Claude matching (also returns canonical names)
	if len(tier3) > 0 {
		matches, err := BatchMatch(ctx, tier3, existing)
		if err != nil {
			imp.logger.Error("claude batch match failed", "error", err)
			// Don't fail the whole job — treat all tier3 as new materials with raw names
			for _, row := range tier3 {
				unit := row.Unit
				if unit == "" {
					unit = "ea"
				}
				id := uuid.New().String()[:20]
				_, err := qtx.CreateMaterial(ctx, repository.CreateMaterialParams{
					ID:          id,
					Name:        row.Name,
					SupplierID:  job.SupplierID,
					UnitPrice:   row.Price,
					Unit:        unit,
					PriceSource: "Supplier",
				})
				if err != nil {
					failed++
				} else {
					created++
				}
			}
		} else {
			for _, match := range matches {
				// Find the original row
				var row RawRow
				for _, r := range tier3 {
					if r.RowIndex == match.RowIndex {
						row = r
						break
					}
				}

				// Use canonical name from Claude, fall back to raw name
				canonicalName := match.CanonicalName
				if canonicalName == "" {
					canonicalName = row.Name
				}

				switch {
				case match.Confidence >= 0.98 && match.CandidateID != "":
					// Auto-update price
					candidate, err := qtx.GetMaterial(ctx, match.CandidateID)
					if err != nil {
						failed++
						continue
					}
					if candidate.PriceSource == "Manual" {
						skipped++
						continue
					}
					_, err = qtx.UpdateMaterial(ctx, repository.UpdateMaterialParams{
						ID:           candidate.ID,
						Name:         candidate.Name,
						SupplierID:   candidate.SupplierID,
						UnitPrice:    row.Price,
						Unit:         candidate.Unit,
						SupplierCode: candidate.SupplierCode,
						PriceSource:  "Supplier",
					})
					if err != nil {
						failed++
					} else {
						updated++
					}

				case match.Confidence >= 0.75 && match.CandidateID != "":
					// Flag for review — store canonical name in reasoning for display
					reasoning := match.Reasoning
					if canonicalName != row.Name {
						reasoning = fmt.Sprintf("%s | Suggested name: %s", match.Reasoning, canonicalName)
					}
					flagID := uuid.New().String()[:20]
					_, err := qtx.CreateImportFlag(ctx, repository.CreateImportFlagParams{
						ID:                  flagID,
						ImportJobID:         job.ID,
						SupplierID:          job.SupplierID,
						RawName:             canonicalName,
						RawUnit:             sql.NullString{String: row.Unit, Valid: row.Unit != ""},
						RawPrice:            sql.NullFloat64{Float64: row.Price, Valid: true},
						RawSupplierCode:     sql.NullString{String: row.SupplierCode, Valid: row.SupplierCode != ""},
						CandidateMaterialID: sql.NullString{String: match.CandidateID, Valid: true},
						Confidence:          sql.NullFloat64{Float64: match.Confidence, Valid: true},
						ClaudeReasoning:     sql.NullString{String: reasoning, Valid: reasoning != ""},
					})
					if err != nil {
						failed++
					} else {
						flagged++
					}

				default:
					// Low confidence or no match — create new material with canonical name
					unit := row.Unit
					if unit == "" {
						unit = "ea"
					}
					id := uuid.New().String()[:20]
					_, err := qtx.CreateMaterial(ctx, repository.CreateMaterialParams{
						ID:          id,
						Name:        canonicalName,
						SupplierID:  job.SupplierID,
						UnitPrice:   row.Price,
						Unit:        unit,
						PriceSource: "Supplier",
					})
					if err != nil {
						failed++
					} else {
						created++
					}
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	// 8. Update progress with final counts
	_, err = imp.queries.UpdateImportJobProgress(ctx, repository.UpdateImportJobProgressParams{
		ID:          job.ID,
		RowsTotal:   sql.NullInt32{Int32: total, Valid: true},
		RowsCreated: sql.NullInt32{Int32: created, Valid: true},
		RowsUpdated: sql.NullInt32{Int32: updated, Valid: true},
		RowsSkipped: sql.NullInt32{Int32: skipped, Valid: true},
		RowsFlagged: sql.NullInt32{Int32: flagged, Valid: true},
		RowsFailed:  sql.NullInt32{Int32: failed, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("updating final progress: %w", err)
	}

	// 9. Mark job complete
	_, err = imp.queries.CompleteImportJob(ctx, job.ID)
	if err != nil {
		return fmt.Errorf("completing job: %w", err)
	}

	imp.logger.Info("import job complete",
		"job_id", job.ID,
		"total", total,
		"created", created,
		"updated", updated,
		"skipped", skipped,
		"flagged", flagged,
		"failed", failed,
	)

	return nil
}
