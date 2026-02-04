package keyboard

import (
	"bytes"
	"database/sql"
	"net/http"
	"time"

	"github.com/dukerupert/skalkaho/internal/domain"
	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/google/uuid"
)

// EstimateCategoryWithChildren represents a tier 1 category with its tier 2 children.
type EstimateCategoryWithChildren struct {
	repository.EstimateCategory
	Children []repository.EstimateCategory
}

// ListEstimates shows all estimates for a job.
func (h *Handler) ListEstimates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	jobID := r.PathValue("jobID")
	orgID := GetOrgID(ctx)

	job, err := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    jobID,
		OrgID: orgID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Failed to load job", http.StatusInternalServerError)
		return
	}

	estimates, err := h.queries.ListEstimatesByJob(ctx, repository.ListEstimatesByJobParams{
		JobID: jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to list estimates", "error", err)
		http.Error(w, "Failed to load estimates", http.StatusInternalServerError)
		return
	}

	// Get client if available
	var client *repository.Client
	if job.ClientID.Valid {
		c, err := h.queries.GetClient(ctx, repository.GetClientParams{
			ID:    job.ClientID.String,
			OrgID: orgID,
		})
		if err == nil {
			client = &c
		}
	}

	data := map[string]interface{}{
		"Job":       job,
		"Estimates": estimates,
		"Client":    client,
	}

	if err := h.renderer.Render(w, "estimates_list", data); err != nil {
		logger.Error("failed to render estimates list", "error", err)
	}
}

// GetNewEstimateForm returns the form for creating a new estimate.
// Shows a client selector for jobs without a client.
func (h *Handler) GetNewEstimateForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	jobID := r.PathValue("jobID")
	orgID := GetOrgID(ctx)

	job, err := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    jobID,
		OrgID: orgID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Failed to load job", http.StatusInternalServerError)
		return
	}

	// Show client selection form
	clients, err := h.queries.ListClients(ctx, orgID)
	if err != nil {
		logger.Error("failed to list clients", "error", err)
		clients = nil
	}

	data := map[string]interface{}{
		"Job":     job,
		"Clients": clients,
	}

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "new_estimate_form", data); err != nil {
		logger.Error("failed to render new estimate form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// CreateEstimate creates a new estimate from a job's current state.
func (h *Handler) CreateEstimate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	jobID := r.PathValue("jobID")
	orgID := GetOrgID(ctx)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get job
	job, err := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    jobID,
		OrgID: orgID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Failed to load job", http.StatusInternalServerError)
		return
	}

	// If job has no client and form provided one, update the job first
	if !job.ClientID.Valid {
		clientID := r.FormValue("client_id")
		if clientID == "" {
			http.Error(w, "Client is required for estimates", http.StatusBadRequest)
			return
		}

		// Update job with client
		_, err = h.queries.UpdateJob(ctx, repository.UpdateJobParams{
			ID:                        jobID,
			Name:                      job.Name,
			CustomerName:              job.CustomerName,
			SurchargePercent:          job.SurchargePercent,
			MaterialSurchargePercent:  job.MaterialSurchargePercent,
			LaborSurchargePercent:     job.LaborSurchargePercent,
			EquipmentSurchargePercent: job.EquipmentSurchargePercent,
			SurchargeMode:             job.SurchargeMode,
			Status:                    job.Status,
			ExpiresAt:                 job.ExpiresAt,
			ClientID:                  sql.NullString{String: clientID, Valid: true},
		})
		if err != nil {
			logger.Error("failed to update job with client", "error", err)
			http.Error(w, "Failed to update job", http.StatusInternalServerError)
			return
		}
	}

	h.doCreateEstimate(w, r, jobID, orgID)
}

// doCreateEstimate performs the actual estimate creation.
func (h *Handler) doCreateEstimate(w http.ResponseWriter, r *http.Request, jobID string, orgID uuid.NullUUID) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	// Get job (refresh after potential update)
	job, err := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Failed to load job", http.StatusInternalServerError)
		return
	}

	// Ensure client exists
	if !job.ClientID.Valid {
		http.Error(w, "Job must have a client to create an estimate", http.StatusBadRequest)
		return
	}

	// Get all categories and line items for calculation
	categories, err := h.queries.ListCategoriesByJob(ctx, repository.ListCategoriesByJobParams{
		JobID: jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to list categories", "error", err)
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	lineItems, err := h.queries.ListLineItemsByJob(ctx, repository.ListLineItemsByJobParams{
		JobID: jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to list line items", "error", err)
		http.Error(w, "Failed to load line items", http.StatusInternalServerError)
		return
	}

	customTypes, err := h.queries.ListJobItemTypes(ctx, repository.ListJobItemTypesParams{
		JobID: jobID,
		OrgID: orgID,
	})
	if err != nil {
		customTypes = []repository.JobItemType{}
	}

	// Calculate job total
	totals := h.calculateTotals(job, categories, lineItems, customTypes)

	// Get next version number
	maxVersionResult, err := h.queries.GetLatestEstimateVersion(ctx, repository.GetLatestEstimateVersionParams{
		JobID: jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get latest version", "error", err)
		http.Error(w, "Failed to create estimate", http.StatusInternalServerError)
		return
	}
	nextVersion := maxVersionResult + 1

	// Create estimate
	estimateID := uuid.New().String()
	estimate, err := h.queries.CreateEstimate(ctx, repository.CreateEstimateParams{
		ID:         estimateID,
		OrgID:      orgID,
		JobID:      jobID,
		Version:    nextVersion,
		Status:     string(domain.EstimateStatusDraft),
		GrandTotal: totals.GrandTotal,
		Notes:      sql.NullString{},
	})
	if err != nil {
		logger.Error("failed to create estimate", "error", err)
		http.Error(w, "Failed to create estimate", http.StatusInternalServerError)
		return
	}

	// Create estimate categories for tier 1 and tier 2 only
	sortOrder := 0
	for _, cat := range categories {
		depth := h.getCategoryDepth(categories, cat.ID)
		if depth > 2 {
			continue // Skip tier 3
		}

		// Calculate category total (including all descendants)
		catTotal := h.calculateCategoryTotal(cat.ID, job, categories, lineItems, customTypes)

		tier := depth
		var parentCatID sql.NullString
		if cat.ParentID.Valid {
			parentCatID = cat.ParentID
		}

		_, err := h.queries.CreateEstimateCategory(ctx, repository.CreateEstimateCategoryParams{
			ID:               uuid.New().String(),
			OrgID:            orgID,
			EstimateID:       estimateID,
			CategoryID:       cat.ID,
			ParentCategoryID: parentCatID,
			Tier:             int64(tier),
			Name:             cat.Name,
			Description:      sql.NullString{},
			Total:            catTotal.Total,
			SortOrder:        int64(sortOrder),
		})
		if err != nil {
			logger.Error("failed to create estimate category", "error", err, "category_id", cat.ID)
		}
		sortOrder++
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/estimates/"+estimate.ID)
		return
	}

	http.Redirect(w, r, "/estimates/"+estimate.ID, http.StatusSeeOther)
}

// GetEstimate shows a single estimate with its categories.
func (h *Handler) GetEstimate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	estimateID := r.PathValue("id")
	orgID := GetOrgID(ctx)

	estimate, err := h.queries.GetEstimate(ctx, repository.GetEstimateParams{
		ID:    estimateID,
		OrgID: orgID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Estimate not found", http.StatusNotFound)
			return
		}
		logger.Error("failed to get estimate", "error", err)
		http.Error(w, "Failed to load estimate", http.StatusInternalServerError)
		return
	}

	job, err := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    estimate.JobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Failed to load job", http.StatusInternalServerError)
		return
	}

	categories, err := h.queries.ListEstimateCategoriesByEstimate(ctx, repository.ListEstimateCategoriesByEstimateParams{
		EstimateID: estimateID,
		OrgID:      orgID,
	})
	if err != nil {
		logger.Error("failed to list estimate categories", "error", err)
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	// Build hierarchy: tier 1 with their tier 2 children
	categoryTree := buildEstimateCategoryTree(categories)

	// Get client if available
	var client *repository.Client
	if job.ClientID.Valid {
		c, err := h.queries.GetClient(ctx, repository.GetClientParams{
			ID:    job.ClientID.String,
			OrgID: orgID,
		})
		if err == nil {
			client = &c
		}
	}

	// Get latest signature request if any
	var signatureRequest *repository.SignatureRequest
	req, err := h.queries.GetSignatureRequestByEstimate(ctx, repository.GetSignatureRequestByEstimateParams{
		EstimateID: estimateID,
		OrgID:      orgID,
	})
	if err == nil {
		signatureRequest = &req
	}

	data := map[string]any{
		"Estimate":         estimate,
		"Job":              job,
		"Categories":       categoryTree,
		"Client":           client,
		"SignatureRequest": signatureRequest,
	}

	if err := h.renderer.Render(w, "estimate", data); err != nil {
		logger.Error("failed to render estimate", "error", err)
	}
}

// buildEstimateCategoryTree organizes categories into tier 1 with tier 2 children.
func buildEstimateCategoryTree(categories []repository.EstimateCategory) []EstimateCategoryWithChildren {
	var tree []EstimateCategoryWithChildren
	childrenByParent := make(map[string][]repository.EstimateCategory)

	// First pass: collect tier 2 categories by parent
	for _, cat := range categories {
		if cat.Tier == 2 && cat.ParentCategoryID.Valid {
			childrenByParent[cat.ParentCategoryID.String] = append(
				childrenByParent[cat.ParentCategoryID.String], cat)
		}
	}

	// Second pass: build tree from tier 1
	for _, cat := range categories {
		if cat.Tier == 1 {
			tree = append(tree, EstimateCategoryWithChildren{
				EstimateCategory: cat,
				Children:         childrenByParent[cat.CategoryID],
			})
		}
	}

	return tree
}

// UpdateEstimateCategoryDescription updates the description for an estimate category.
func (h *Handler) UpdateEstimateCategoryDescription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	categoryID := r.PathValue("id")
	orgID := GetOrgID(ctx)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	description := sql.NullString{}
	if desc := r.FormValue("description"); desc != "" {
		description = sql.NullString{String: desc, Valid: true}
	}

	cat, err := h.queries.UpdateEstimateCategoryDescription(ctx, repository.UpdateEstimateCategoryDescriptionParams{
		ID:          categoryID,
		Description: description,
	})
	if err != nil {
		logger.Error("failed to update category description", "error", err)
		http.Error(w, "Failed to update description", http.StatusInternalServerError)
		return
	}

	// Return to estimate page
	estCat, err := h.queries.GetEstimateCategory(ctx, repository.GetEstimateCategoryParams{
		ID:    categoryID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get estimate category", "error", err)
		http.Error(w, "Failed to get category", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		// Return updated category row partial
		var buf bytes.Buffer
		if err := h.renderer.RenderPartial(&buf, "estimate_category_row", map[string]interface{}{
			"Category": cat,
		}); err != nil {
			// If partial doesn't exist, just redirect
			w.Header().Set("HX-Redirect", "/estimates/"+estCat.EstimateID)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(buf.Bytes())
		return
	}

	http.Redirect(w, r, "/estimates/"+estCat.EstimateID, http.StatusSeeOther)
}

// SendEstimate updates status to sent and logs email (stub).
func (h *Handler) SendEstimate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	estimateID := r.PathValue("id")
	orgID := GetOrgID(ctx)

	estimate, err := h.queries.GetEstimate(ctx, repository.GetEstimateParams{
		ID:    estimateID,
		OrgID: orgID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Estimate not found", http.StatusNotFound)
			return
		}
		logger.Error("failed to get estimate", "error", err)
		http.Error(w, "Failed to load estimate", http.StatusInternalServerError)
		return
	}

	if estimate.Status != string(domain.EstimateStatusDraft) {
		http.Error(w, "Only draft estimates can be sent", http.StatusBadRequest)
		return
	}

	// Get job and client for email stub
	job, _ := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    estimate.JobID,
		OrgID: orgID,
	})
	var clientEmail string
	if job.ClientID.Valid {
		client, err := h.queries.GetClient(ctx, repository.GetClientParams{
			ID:    job.ClientID.String,
			OrgID: orgID,
		})
		if err == nil && client.Email.Valid {
			clientEmail = client.Email.String
		}
	}

	// Update status to sent
	now := time.Now().Format(time.RFC3339)
	_, err = h.queries.UpdateEstimate(ctx, repository.UpdateEstimateParams{
		ID:          estimateID,
		Status:      string(domain.EstimateStatusSent),
		Notes:       estimate.Notes,
		SentAt:      sql.NullString{String: now, Valid: true},
		RespondedAt: estimate.RespondedAt,
	})
	if err != nil {
		logger.Error("failed to update estimate status", "error", err)
		http.Error(w, "Failed to send estimate", http.StatusInternalServerError)
		return
	}

	// EMAIL STUB: Log instead of sending
	logger.Info("EMAIL STUB: Would send estimate",
		"estimate_id", estimateID,
		"job_name", job.Name,
		"client_email", clientEmail,
		"version", estimate.Version,
		"total", estimate.GrandTotal,
	)

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/estimates/"+estimateID)
		return
	}

	http.Redirect(w, r, "/estimates/"+estimateID, http.StatusSeeOther)
}

// UpdateEstimateStatus updates the estimate status (accept/reject).
func (h *Handler) UpdateEstimateStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	estimateID := r.PathValue("id")
	orgID := GetOrgID(ctx)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	estimate, err := h.queries.GetEstimate(ctx, repository.GetEstimateParams{
		ID:    estimateID,
		OrgID: orgID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Estimate not found", http.StatusNotFound)
			return
		}
		logger.Error("failed to get estimate", "error", err)
		http.Error(w, "Failed to load estimate", http.StatusInternalServerError)
		return
	}

	status := r.FormValue("status")
	if !domain.ValidEstimateStatus(status) {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	// Set responded_at for accept/reject
	respondedAt := estimate.RespondedAt
	if status == string(domain.EstimateStatusAccepted) || status == string(domain.EstimateStatusRejected) {
		now := time.Now().Format(time.RFC3339)
		respondedAt = sql.NullString{String: now, Valid: true}
	}

	_, err = h.queries.UpdateEstimate(ctx, repository.UpdateEstimateParams{
		ID:          estimateID,
		Status:      status,
		Notes:       estimate.Notes,
		SentAt:      estimate.SentAt,
		RespondedAt: respondedAt,
	})
	if err != nil {
		logger.Error("failed to update estimate status", "error", err)
		http.Error(w, "Failed to update status", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/estimates/"+estimateID)
		return
	}

	http.Redirect(w, r, "/estimates/"+estimateID, http.StatusSeeOther)
}

// GetEstimatePreview shows a client-facing preview of the estimate.
func (h *Handler) GetEstimatePreview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	estimateID := r.PathValue("id")
	orgID := GetOrgID(ctx)

	estimate, err := h.queries.GetEstimate(ctx, repository.GetEstimateParams{
		ID:    estimateID,
		OrgID: orgID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Estimate not found", http.StatusNotFound)
			return
		}
		logger.Error("failed to get estimate", "error", err)
		http.Error(w, "Failed to load estimate", http.StatusInternalServerError)
		return
	}

	job, _ := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    estimate.JobID,
		OrgID: orgID,
	})
	categories, _ := h.queries.ListEstimateCategoriesByEstimate(ctx, repository.ListEstimateCategoriesByEstimateParams{
		EstimateID: estimateID,
		OrgID:      orgID,
	})
	categoryTree := buildEstimateCategoryTree(categories)

	var client *repository.Client
	if job.ClientID.Valid {
		c, err := h.queries.GetClient(ctx, repository.GetClientParams{
			ID:    job.ClientID.String,
			OrgID: orgID,
		})
		if err == nil {
			client = &c
		}
	}

	data := map[string]interface{}{
		"Estimate":   estimate,
		"Job":        job,
		"Categories": categoryTree,
		"Client":     client,
		"Preview":    true,
	}

	if err := h.renderer.Render(w, "estimate_preview", data); err != nil {
		logger.Error("failed to render estimate preview", "error", err)
	}
}

// DeleteEstimate deletes an estimate.
func (h *Handler) DeleteEstimate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	estimateID := r.PathValue("id")
	orgID := GetOrgID(ctx)

	estimate, err := h.queries.GetEstimate(ctx, repository.GetEstimateParams{
		ID:    estimateID,
		OrgID: orgID,
	})
	if err != nil {
		http.Error(w, "Estimate not found", http.StatusNotFound)
		return
	}

	if err := h.queries.DeleteEstimate(ctx, repository.DeleteEstimateParams{
		ID:    estimateID,
		OrgID: orgID,
	}); err != nil {
		logger.Error("failed to delete estimate", "error", err)
		http.Error(w, "Failed to delete estimate", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/jobs/"+estimate.JobID+"/estimates")
		return
	}

	http.Redirect(w, r, "/jobs/"+estimate.JobID+"/estimates", http.StatusSeeOther)
}
