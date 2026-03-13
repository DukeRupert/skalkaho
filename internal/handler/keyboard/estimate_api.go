package keyboard

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dukerupert/skalkaho/internal/domain"
	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/google/uuid"
)

// --- GET /api/estimate/{jobID} ---

// GetEstimateBuilder renders the estimate builder page with the Svelte mount point.
func (h *Handler) GetEstimateBuilder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	jobID := r.PathValue("jobID")

	job, err := h.queries.GetJob(ctx, repository.GetJobParams{ID: jobID, OrgID: orgID})
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Failed to load job", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Job": job,
	}

	if err := h.renderer.Render(w, "estimate_builder", data); err != nil {
		logger.Error("failed to render estimate builder", "error", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// GetEstimateJSON returns the full estimate payload for the Svelte estimate builder.
func (h *Handler) GetEstimateJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	jobID := r.PathValue("jobID")

	// Load job
	job, err := h.queries.GetJob(ctx, repository.GetJobParams{ID: jobID, OrgID: orgID})
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "Job not found")
			return
		}
		logger.Error("failed to get job", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to load job")
		return
	}

	// Load all estimate hierarchy data in parallel-safe order
	sections, err := h.queries.ListSectionsByJob(ctx, repository.ListSectionsByJobParams{
		JobID: jobID, OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to list sections", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to load sections")
		return
	}

	subcategories, err := h.queries.ListSubcategoriesByJob(ctx, repository.ListSubcategoriesByJobParams{
		JobID: jobID, OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to list subcategories", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to load subcategories")
		return
	}

	groups, err := h.queries.ListComponentGroupsByJob(ctx, repository.ListComponentGroupsByJobParams{
		JobID: jobID, OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to list component groups", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to load component groups")
		return
	}

	lineItems, err := h.queries.ListEstimateLineItemsByJob(ctx, repository.ListEstimateLineItemsByJobParams{
		JobID: jobID, OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to list estimate line items", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to load line items")
		return
	}

	// Load item templates for autocomplete (materials + rates)
	templates, err := h.queries.ListItemTemplates(ctx, orgID)
	if err != nil {
		logger.Error("failed to list item templates", "error", err)
		templates = []repository.ItemTemplate{}
	}

	// Build the nested payload
	payload := buildEstimatePayload(job, sections, subcategories, groups, lineItems, templates)
	writeJSON(w, http.StatusOK, payload)
}

// buildEstimatePayload assembles the nested JSON payload from flat DB rows.
func buildEstimatePayload(
	job repository.Job,
	sections []repository.Section,
	subcategories []repository.Subcategory,
	groups []repository.ComponentGroup,
	lineItems []repository.EstimateLineItem,
	templates []repository.ItemTemplate,
) domain.EstimatePayload {
	// Index subcategories by section_id
	subcatsBySectionID := make(map[string][]repository.Subcategory)
	for _, sc := range subcategories {
		subcatsBySectionID[sc.SectionID] = append(subcatsBySectionID[sc.SectionID], sc)
	}

	// Index groups by subcategory_id
	groupsBySubcatID := make(map[string][]repository.ComponentGroup)
	for _, g := range groups {
		groupsBySubcatID[g.SubcategoryID] = append(groupsBySubcatID[g.SubcategoryID], g)
	}

	// Index line items by subcategory_id and component_group_id
	itemsBySubcatID := make(map[string][]repository.EstimateLineItem)
	itemsByGroupID := make(map[string][]repository.EstimateLineItem)
	for _, li := range lineItems {
		if li.ComponentGroupID.Valid {
			itemsByGroupID[li.ComponentGroupID.String] = append(itemsByGroupID[li.ComponentGroupID.String], li)
		} else {
			itemsBySubcatID[li.SubcategoryID] = append(itemsBySubcatID[li.SubcategoryID], li)
		}
	}

	// Build sections
	domainSections := make([]domain.EstimateSection, len(sections))
	for i, s := range sections {
		// Build subcategories for this section
		subcats := subcatsBySectionID[s.ID]
		domainSubcats := make([]domain.EstimateSubcategory, len(subcats))
		for j, sc := range subcats {
			// Build component groups for this subcategory
			scGroups := groupsBySubcatID[sc.ID]
			domainGroups := make([]domain.EstimateComponentGroup, len(scGroups))
			for k, g := range scGroups {
				domainGroups[k] = domain.EstimateComponentGroup{
					ID:        g.ID,
					Name:      g.Name,
					SortOrder: int(g.SortOrder),
					LineItems: convertEstimateLineItems(itemsByGroupID[g.ID]),
				}
			}

			domainSubcats[j] = domain.EstimateSubcategory{
				ID:        sc.ID,
				Name:      sc.Name,
				SortOrder: int(sc.SortOrder),
				LumpSum:   sc.LumpSum,
				MarkupOverrides: domain.MarkupOverrides{
					Materials: nullFloat64ToPtr(sc.MaterialsMarkup),
					Labor:     nullFloat64ToPtr(sc.LaborMarkup),
					Equipment: nullFloat64ToPtr(sc.EquipmentMarkup),
					Subs:      nullFloat64ToPtr(sc.SubsMarkup),
					Other:     nullFloat64ToPtr(sc.OtherMarkup),
				},
				MarkupEnabled: domain.MarkupEnabled{
					Materials: sc.MaterialsMarkupEnabled,
					Labor:     sc.LaborMarkupEnabled,
					Equipment: sc.EquipmentMarkupEnabled,
					Subs:      sc.SubsMarkupEnabled,
					Other:     sc.OtherMarkupEnabled,
				},
				ComponentGroups: domainGroups,
				LineItems:       convertEstimateLineItems(itemsBySubcatID[sc.ID]),
			}
		}

		domainSections[i] = domain.EstimateSection{
			ID:            s.ID,
			Name:          s.Name,
			SortOrder:     int(s.SortOrder),
			Subcategories: domainSubcats,
		}
	}

	// Split templates into materials_db and rates_db
	var materialsDB []domain.MaterialDBEntry
	var ratesDB []domain.RateDBEntry
	for _, t := range templates {
		switch t.Type {
		case "material":
			materialsDB = append(materialsDB, domain.MaterialDBEntry{
				ID:           t.ID,
				Name:         t.Name,
				Category:     t.Category,
				DefaultUnit:  t.DefaultUnit,
				DefaultPrice: t.DefaultPrice,
			})
		default: // labor, equipment, and any other types become rates
			ratesDB = append(ratesDB, domain.RateDBEntry{
				ID:           t.ID,
				Type:         t.Type,
				Name:         t.Name,
				Category:     t.Category,
				DefaultUnit:  t.DefaultUnit,
				DefaultPrice: t.DefaultPrice,
			})
		}
	}
	if materialsDB == nil {
		materialsDB = []domain.MaterialDBEntry{}
	}
	if ratesDB == nil {
		ratesDB = []domain.RateDBEntry{}
	}

	return domain.EstimatePayload{
		Project: domain.EstimateProject{
			ID:     job.ID,
			Name:   job.Name,
			Status: job.Status,
		},
		Globals: domain.MarkupGlobals{
			MaterialsMarkup: job.MaterialsMarkup,
			LaborMarkup:     job.LaborMarkup,
			EquipmentMarkup: job.EquipmentMarkup,
			SubsMarkup:      job.SubsMarkup,
			OtherMarkup:     job.OtherMarkup,
		},
		Sections:    domainSections,
		MaterialsDB: materialsDB,
		RatesDB:     ratesDB,
	}
}

// convertEstimateLineItems converts repository line items to domain line items.
func convertEstimateLineItems(items []repository.EstimateLineItem) []domain.EstimateLineItem {
	if len(items) == 0 {
		return []domain.EstimateLineItem{}
	}
	result := make([]domain.EstimateLineItem, len(items))
	for i, li := range items {
		result[i] = domain.EstimateLineItem{
			ID:           li.ID,
			CategoryType: domain.CategoryType(li.CategoryType),
			ItemName:     li.ItemName,
			Quantity:     li.Quantity,
			Unit:         li.Unit,
			UnitPrice:    li.UnitPrice,
			IsCustom:     li.IsCustom,
			PriceOverride: li.PriceOverride,
			Description:  nullStringToPtr(li.Description),
			SortOrder:    int(li.SortOrder),
		}
		if li.MaterialID.Valid {
			mid := int64(li.MaterialID.Int32)
			result[i].MaterialID = &mid
		}
		if li.ComponentGroupID.Valid {
			result[i].ComponentGroupID = &li.ComponentGroupID.String
		}
	}
	return result
}

// --- POST /api/estimate/{jobID} ---

// SaveEstimateJSON accepts the full estimate payload and persists it in a single transaction.
func (h *Handler) SaveEstimateJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	jobID := r.PathValue("jobID")

	// Verify job exists
	job, err := h.queries.GetJob(ctx, repository.GetJobParams{ID: jobID, OrgID: orgID})
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "Job not found")
			return
		}
		logger.Error("failed to get job", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to load job")
		return
	}

	// Parse payload
	var payload domain.EstimatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validate
	if err := validateEstimatePayload(payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Begin transaction
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("failed to begin transaction", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to save estimate")
		return
	}
	defer func() { _ = tx.Rollback() }()

	qtx := h.queries.WithTx(tx)

	// Update globals on the job
	if err := qtx.UpdateJobMarkupGlobals(ctx, repository.UpdateJobMarkupGlobalsParams{
		MaterialsMarkup: payload.Globals.MaterialsMarkup,
		LaborMarkup:     payload.Globals.LaborMarkup,
		EquipmentMarkup: payload.Globals.EquipmentMarkup,
		SubsMarkup:      payload.Globals.SubsMarkup,
		OtherMarkup:     payload.Globals.OtherMarkup,
		ID:              jobID,
		OrgID:           orgID,
	}); err != nil {
		logger.Error("failed to update markup globals", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to save estimate")
		return
	}

	// Delete existing estimate data (cascades through sections -> subcategories -> groups -> items)
	if err := qtx.DeleteSectionsByJob(ctx, repository.DeleteSectionsByJobParams{
		JobID: jobID, OrgID: orgID,
	}); err != nil {
		logger.Error("failed to delete existing sections", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to save estimate")
		return
	}

	// Insert all sections, subcategories, groups, and line items
	for _, section := range payload.Sections {
		if _, err := qtx.CreateSection(ctx, repository.CreateSectionParams{
			ID:        section.ID,
			OrgID:     orgID,
			JobID:     jobID,
			Name:      section.Name,
			SortOrder: int64(section.SortOrder),
		}); err != nil {
			logger.Error("failed to create section", "error", err, "section_id", section.ID)
			writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to save section: %s", section.Name))
			return
		}

		for _, subcat := range section.Subcategories {
			if _, err := qtx.CreateSubcategory(ctx, repository.CreateSubcategoryParams{
				ID:                     subcat.ID,
				OrgID:                  orgID,
				SectionID:             section.ID,
				Name:                   subcat.Name,
				SortOrder:             int64(subcat.SortOrder),
				LumpSum:               subcat.LumpSum,
				MaterialsMarkup:       ptrToNullFloat64(subcat.MarkupOverrides.Materials),
				LaborMarkup:           ptrToNullFloat64(subcat.MarkupOverrides.Labor),
				EquipmentMarkup:       ptrToNullFloat64(subcat.MarkupOverrides.Equipment),
				SubsMarkup:            ptrToNullFloat64(subcat.MarkupOverrides.Subs),
				OtherMarkup:           ptrToNullFloat64(subcat.MarkupOverrides.Other),
				MaterialsMarkupEnabled: subcat.MarkupEnabled.Materials,
				LaborMarkupEnabled:     subcat.MarkupEnabled.Labor,
				EquipmentMarkupEnabled: subcat.MarkupEnabled.Equipment,
				SubsMarkupEnabled:      subcat.MarkupEnabled.Subs,
				OtherMarkupEnabled:     subcat.MarkupEnabled.Other,
			}); err != nil {
				logger.Error("failed to create subcategory", "error", err, "subcategory_id", subcat.ID)
				writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to save subcategory: %s", subcat.Name))
				return
			}

			// Insert component groups
			for _, group := range subcat.ComponentGroups {
				if _, err := qtx.CreateComponentGroup(ctx, repository.CreateComponentGroupParams{
					ID:            group.ID,
					OrgID:         orgID,
					SubcategoryID: subcat.ID,
					Name:          group.Name,
					SortOrder:     int64(group.SortOrder),
				}); err != nil {
					logger.Error("failed to create component group", "error", err, "group_id", group.ID)
					writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to save group: %s", group.Name))
					return
				}

				// Insert grouped line items
				for _, li := range group.LineItems {
					if err := insertEstimateLineItem(ctx, qtx, orgID, subcat.ID, &group.ID, li); err != nil {
						logger.Error("failed to create grouped line item", "error", err, "item_id", li.ID)
						writeJSONError(w, http.StatusInternalServerError, "Failed to save line item")
						return
					}
				}
			}

			// Insert ungrouped line items
			for _, li := range subcat.LineItems {
				if err := insertEstimateLineItem(ctx, qtx, orgID, subcat.ID, nil, li); err != nil {
					logger.Error("failed to create ungrouped line item", "error", err, "item_id", li.ID)
					writeJSONError(w, http.StatusInternalServerError, "Failed to save line item")
					return
				}
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		logger.Error("failed to commit transaction", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to save estimate")
		return
	}

	// Re-read and return the saved state
	_ = job // suppress unused warning
	savedJob, _ := h.queries.GetJob(ctx, repository.GetJobParams{ID: jobID, OrgID: orgID})
	sections, _ := h.queries.ListSectionsByJob(ctx, repository.ListSectionsByJobParams{JobID: jobID, OrgID: orgID})
	subcategories, _ := h.queries.ListSubcategoriesByJob(ctx, repository.ListSubcategoriesByJobParams{JobID: jobID, OrgID: orgID})
	groups, _ := h.queries.ListComponentGroupsByJob(ctx, repository.ListComponentGroupsByJobParams{JobID: jobID, OrgID: orgID})
	items, _ := h.queries.ListEstimateLineItemsByJob(ctx, repository.ListEstimateLineItemsByJobParams{JobID: jobID, OrgID: orgID})
	templates, _ := h.queries.ListItemTemplates(ctx, orgID)

	result := buildEstimatePayload(savedJob, sections, subcategories, groups, items, templates)
	writeJSON(w, http.StatusOK, result)
}

// insertEstimateLineItem inserts a single estimate line item.
func insertEstimateLineItem(
	ctx context.Context,
	qtx *repository.Queries,
	orgID uuid.NullUUID,
	subcategoryID string,
	componentGroupID *string,
	li domain.EstimateLineItem,
) error {
	var materialID sql.NullInt32
	if li.MaterialID != nil {
		materialID = sql.NullInt32{Int32: int32(*li.MaterialID), Valid: true}
	}
	var groupID sql.NullString
	if componentGroupID != nil {
		groupID = sql.NullString{String: *componentGroupID, Valid: true}
	}

	_, err := qtx.CreateEstimateLineItem(ctx, repository.CreateEstimateLineItemParams{
		ID:               li.ID,
		OrgID:            orgID,
		SubcategoryID:    subcategoryID,
		ComponentGroupID: groupID,
		CategoryType:     string(li.CategoryType),
		ItemName:         li.ItemName,
		Quantity:         li.Quantity,
		Unit:             li.Unit,
		UnitPrice:        li.UnitPrice,
		IsCustom:         li.IsCustom,
		MaterialID:       materialID,
		PriceOverride:    li.PriceOverride,
		Description:      ptrToNullString(li.Description),
		SortOrder:        int64(li.SortOrder),
	})
	return err
}

// validateEstimatePayload performs basic validation on the incoming estimate payload.
func validateEstimatePayload(p domain.EstimatePayload) error {
	for _, section := range p.Sections {
		if section.ID == "" {
			return fmt.Errorf("section missing ID")
		}
		if section.Name == "" {
			return fmt.Errorf("section missing name")
		}
		for _, subcat := range section.Subcategories {
			if subcat.ID == "" {
				return fmt.Errorf("subcategory missing ID")
			}
			if subcat.Name == "" {
				return fmt.Errorf("subcategory missing name")
			}
			for _, group := range subcat.ComponentGroups {
				if group.ID == "" {
					return fmt.Errorf("component group missing ID")
				}
				if group.Name == "" {
					return fmt.Errorf("component group missing name")
				}
				for _, li := range group.LineItems {
					if err := validateEstimateLineItem(li); err != nil {
						return err
					}
				}
			}
			for _, li := range subcat.LineItems {
				if err := validateEstimateLineItem(li); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// validateEstimateLineItem validates a single line item.
func validateEstimateLineItem(li domain.EstimateLineItem) error {
	if li.ID == "" {
		return fmt.Errorf("line item missing ID")
	}
	if li.ItemName == "" {
		return fmt.Errorf("line item missing name")
	}
	switch li.CategoryType {
	case domain.CategoryTypeMaterials, domain.CategoryTypeLabor, domain.CategoryTypeEquipment, domain.CategoryTypeSubs, domain.CategoryTypeOther:
		// valid
	default:
		return fmt.Errorf("line item has invalid category_type: %s", li.CategoryType)
	}
	return nil
}
