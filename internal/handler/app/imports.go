package app

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/dukerupert/skalkaho/internal/repository"
)

// ImportsPageData holds data for both the index and show pages.
type ImportsPageData struct {
	PageData
	// Index fields
	Jobs      []repository.ListImportJobsRow
	Suppliers []repository.Supplier
	// Show fields
	Job          *repository.MaterialImportJob
	SupplierName string
	Flags        []repository.ListImportFlagsRow
	PendingFlags []repository.ListPendingFlagsRow
}

// ListImports renders the import jobs list page.
func (h *Handler) ListImports(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jobs, err := h.queries.ListImportJobs(ctx)
	if err != nil {
		h.logger.Error("listing import jobs", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	suppliers, err := h.queries.ListSuppliers(ctx)
	if err != nil {
		h.logger.Error("listing suppliers", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := ImportsPageData{
		PageData:  h.pageData(r, "imports"),
		Jobs:      jobs,
		Suppliers: suppliers,
	}

	if err := h.renderer.Render(w, "imports.html", data); err != nil {
		h.logger.Error("rendering imports", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// CreateImport handles POST /imports — file upload and job creation.
func (h *Handler) CreateImport(w http.ResponseWriter, r *http.Request) {
	// 10MB max upload
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too large (max 10MB)", http.StatusBadRequest)
		return
	}

	supplierID := r.FormValue("supplier_id")
	if supplierID == "" {
		http.Error(w, "Supplier is required", http.StatusBadRequest)
		return
	}

	// Validate supplier exists
	_, err := h.queries.GetSupplier(r.Context(), supplierID)
	if err != nil {
		http.Error(w, "Invalid supplier", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer func() { _ = file.Close() }()

	// Validate file extension
	ext := filepath.Ext(header.Filename)
	if ext != ".xlsx" && ext != ".xls" {
		http.Error(w, "Only .xlsx and .xls files are supported", http.StatusBadRequest)
		return
	}

	// Write to temp file
	tempID := uuid.New().String()[:12]
	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("skalkaho-import-%s%s", tempID, ext))
	tempFile, err := os.Create(tempPath)
	if err != nil {
		h.logger.Error("creating temp file", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(tempFile, file); err != nil {
		_ = tempFile.Close()
		_ = os.Remove(tempPath)
		h.logger.Error("writing temp file", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_ = tempFile.Close()

	// Create job row
	jobID := uuid.New().String()[:20]
	_, err = h.queries.CreateImportJob(r.Context(), repository.CreateImportJobParams{
		ID:           jobID,
		SupplierID:   supplierID,
		Filename:     header.Filename,
		TempFilePath: tempPath,
	})
	if err != nil {
		_ = os.Remove(tempPath)
		h.logger.Error("creating import job", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/imports/%s", jobID), http.StatusSeeOther)
}

// ShowImport renders the import job detail/status page.
func (h *Handler) ShowImport(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, err := h.queries.GetImportJob(r.Context(), id)
	if err != nil {
		h.logger.Error("getting import job", "error", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	supplier, err := h.queries.GetSupplier(r.Context(), job.SupplierID)
	if err != nil {
		h.logger.Error("getting supplier", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := ImportsPageData{
		PageData:     h.pageData(r, "imports"),
		Job:          &job,
		SupplierName: supplier.Name,
	}

	// Load flags if job is complete
	if job.Status == "complete" {
		flags, err := h.queries.ListImportFlags(r.Context(), id)
		if err != nil {
			h.logger.Error("listing import flags", "error", err)
		} else {
			data.Flags = flags
		}
		pending, err := h.queries.ListPendingFlags(r.Context(), id)
		if err != nil {
			h.logger.Error("listing pending flags", "error", err)
		} else {
			data.PendingFlags = pending
		}
	}

	if err := h.renderer.Render(w, "imports.html", data); err != nil {
		h.logger.Error("rendering import show", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ImportStatus returns the polling fragment for HTMX.
func (h *Handler) ImportStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, err := h.queries.GetImportJob(r.Context(), id)
	if err != nil {
		h.logger.Error("getting import job", "error", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	data := ImportsPageData{
		PageData: h.pageData(r, "imports"),
		Job:      &job,
	}

	supplier, err := h.queries.GetSupplier(r.Context(), job.SupplierID)
	if err == nil {
		data.SupplierName = supplier.Name
	}

	var block string
	switch job.Status {
	case "pending", "processing":
		block = "status-pending"
	case "complete":
		flags, err := h.queries.ListPendingFlags(r.Context(), id)
		if err == nil {
			data.PendingFlags = flags
		}
		allFlags, err := h.queries.ListImportFlags(r.Context(), id)
		if err == nil {
			data.Flags = allFlags
		}
		block = "status-complete"
	case "failed":
		block = "status-failed"
	}

	if err := h.renderer.RenderPartial(w, "imports.html", block, data); err != nil {
		h.logger.Error("rendering import status", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ResolveImportFlag handles POST /imports/{id}/flags/{fid}/resolve.
func (h *Handler) ResolveImportFlag(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("id")
	flagID := r.PathValue("fid")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	resolution := r.FormValue("resolution")
	if resolution != "accepted" && resolution != "rejected" && resolution != "created" {
		http.Error(w, "Invalid resolution", http.StatusBadRequest)
		return
	}

	flag, err := h.queries.GetImportFlag(r.Context(), flagID)
	if err != nil {
		h.logger.Error("getting import flag", "error", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	ctx := r.Context()

	switch resolution {
	case "accepted":
		// Update the candidate material's price with the raw price from the spreadsheet
		if flag.CandidateMaterialID.Valid && flag.RawPrice.Valid {
			candidate, err := h.queries.GetMaterial(ctx, flag.CandidateMaterialID.String)
			if err == nil {
				_, _ = h.queries.UpdateMaterial(ctx, repository.UpdateMaterialParams{
					ID:           candidate.ID,
					Name:         candidate.Name,
					SupplierID:   candidate.SupplierID,
					UnitPrice:    flag.RawPrice.Float64,
					Unit:         candidate.Unit,
					SupplierCode: candidate.SupplierCode,
					PriceSource:  "Supplier",
				})
			}
		}

	case "created":
		// Create a new material from the raw row values
		unit := "ea"
		if flag.RawUnit.Valid && flag.RawUnit.String != "" {
			unit = flag.RawUnit.String
		}
		price := 0.0
		if flag.RawPrice.Valid {
			price = flag.RawPrice.Float64
		}
		id := uuid.New().String()[:20]
		_, _ = h.queries.CreateMaterial(ctx, repository.CreateMaterialParams{
			ID:          id,
			Name:        flag.RawName,
			SupplierID:  flag.SupplierID,
			UnitPrice:   price,
			Unit:        unit,
			SupplierCode: sql.NullString{
				String: flag.RawSupplierCode.String,
				Valid:  flag.RawSupplierCode.Valid && flag.RawSupplierCode.String != "",
			},
			PriceSource: "Supplier",
		})
	}
	// "rejected" — no material changes needed

	// Mark the flag resolved
	_, err = h.queries.ResolveFlag(ctx, repository.ResolveFlagParams{
		ID:         flagID,
		Resolution: sql.NullString{String: resolution, Valid: true},
	})
	if err != nil {
		h.logger.Error("resolving flag", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return updated flag row partial for HTMX swap
	if isHTMXPartial(r) {
		// Re-fetch pending flags to render the updated row
		pending, err := h.queries.ListPendingFlags(ctx, jobID)
		if err != nil {
			h.logger.Error("listing pending flags", "error", err)
		}
		job, err := h.queries.GetImportJob(ctx, jobID)
		if err != nil {
			h.logger.Error("getting import job after resolve", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		data := ImportsPageData{
			PageData:     h.pageData(r, "imports"),
			Job:          &job,
			PendingFlags: pending,
		}
		if err := h.renderer.RenderPartial(w, "imports.html", "flag-table-body", data); err != nil {
			h.logger.Error("rendering flag table body", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/imports/%s", jobID), http.StatusSeeOther)
}
