package keyboard

import (
	"bytes"
	"database/sql"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/google/uuid"
)

const pageSize = 20

// SpreadsheetSection represents a category section in the spreadsheet view.
type SpreadsheetSection struct {
	ID       string
	Name     string
	Depth    int
	Total    float64
	JobID    string
	Items    []SpreadsheetItem
	Children []SpreadsheetSection
}

// SpreadsheetItem wraps a LineItem with a color for the spreadsheet view.
type SpreadsheetItem struct {
	repository.LineItem
	Color string // "forest", "copper", "slate", or custom type color
}

// getItemTypeColor maps a type slug to a Tailwind color string.
func getItemTypeColor(slug string, customTypes []repository.JobItemType) string {
	switch slug {
	case "material":
		return "forest"
	case "labor":
		return "copper"
	case "equipment":
		return "slate"
	default:
		for _, ct := range customTypes {
			if ct.Slug == slug {
				return ct.Color
			}
		}
		return "slate"
	}
}

// buildSpreadsheetSections builds a hierarchical list of spreadsheet sections from categories and items.
func (h *Handler) buildSpreadsheetSections(job repository.Job, categories []repository.Category, lineItems []repository.LineItem, customTypes []repository.JobItemType) []SpreadsheetSection {
	// Group items by category ID
	itemsByCategory := make(map[string][]repository.LineItem)
	for _, item := range lineItems {
		itemsByCategory[item.CategoryID] = append(itemsByCategory[item.CategoryID], item)
	}

	// Group children by parent ID
	childrenByParent := make(map[string][]repository.Category)
	for _, cat := range categories {
		if cat.ParentID.Valid {
			childrenByParent[cat.ParentID.String] = append(childrenByParent[cat.ParentID.String], cat)
		}
	}

	// Recursive builder
	var buildSection func(cat repository.Category, depth int) SpreadsheetSection
	buildSection = func(cat repository.Category, depth int) SpreadsheetSection {
		catTotal := h.calculateCategoryTotal(cat.ID, job, categories, lineItems, customTypes)

		// Convert items with colors
		rawItems := itemsByCategory[cat.ID]
		items := make([]SpreadsheetItem, len(rawItems))
		for i, item := range rawItems {
			items[i] = SpreadsheetItem{
				LineItem: item,
				Color:    getItemTypeColor(item.Type, customTypes),
			}
		}

		// Build children recursively
		children := make([]SpreadsheetSection, 0)
		for _, child := range childrenByParent[cat.ID] {
			children = append(children, buildSection(child, depth+1))
		}

		return SpreadsheetSection{
			ID:       cat.ID,
			Name:     cat.Name,
			Depth:    depth,
			Total:    catTotal.Total,
			JobID:    cat.JobID,
			Items:    items,
			Children: children,
		}
	}

	// Start from root categories (no parent)
	var sections []SpreadsheetSection
	for _, cat := range categories {
		if !cat.ParentID.Valid {
			sections = append(sections, buildSection(cat, 1))
		}
	}

	return sections
}

// JobWithTotal wraps a Job with its calculated grand total and client info.
type JobWithTotal struct {
	repository.Job
	GrandTotal            float64
	ClientName            string
	EstimateCount         int64
	LatestEstimateStatus  string
	LatestSignatureStatus string
}

// PaginationData holds pagination state for templates.
type PaginationData struct {
	CurrentPage int
	TotalPages  int
	TotalItems  int64
	HasPrev     bool
	HasNext     bool
}

// EstimateWithStatus wraps an Estimate with its signature status.
type EstimateWithStatus struct {
	repository.Estimate
	SignatureStatus string // "", "pending", "signed", "expired", "cancelled"
}

// ListJobs shows the keyboard-centric jobs list with pagination and filtering.
func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)

	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	status := r.URL.Query().Get("status")
	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = "newest"
	}

	offset := int64((page - 1) * pageSize)

	// Get total count for pagination
	totalItems, err := h.queries.CountJobs(ctx, repository.CountJobsParams{
		Status: status,
		OrgID:  orgID,
	})
	if err != nil {
		logger.Error("failed to count jobs", "error", err)
		http.Error(w, "Failed to load jobs", http.StatusInternalServerError)
		return
	}

	totalPages := int(totalItems+pageSize-1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	// Get jobs based on sort order
	var jobs []repository.Job

	switch sortBy {
	case "oldest":
		jobs, err = h.queries.ListJobsPaginatedOldest(ctx, repository.ListJobsPaginatedOldestParams{
			Status: status,
			OrgID:  orgID,
			Offset: int32(offset),
			Limit:  int32(pageSize),
		})
	case "name_asc":
		jobs, err = h.queries.ListJobsPaginatedByName(ctx, repository.ListJobsPaginatedByNameParams{
			Status: status,
			OrgID:  orgID,
			Offset: int32(offset),
			Limit:  int32(pageSize),
		})
	case "name_desc":
		jobs, err = h.queries.ListJobsPaginatedByNameDesc(ctx, repository.ListJobsPaginatedByNameDescParams{
			Status: status,
			OrgID:  orgID,
			Offset: int32(offset),
			Limit:  int32(pageSize),
		})
	default: // newest
		jobs, err = h.queries.ListJobsPaginated(ctx, repository.ListJobsPaginatedParams{
			Status: status,
			OrgID:  orgID,
			Offset: int32(offset),
			Limit:  int32(pageSize),
		})
	}

	if err != nil {
		logger.Error("failed to list jobs", "error", err)
		http.Error(w, "Failed to load jobs", http.StatusInternalServerError)
		return
	}

	// Calculate totals for each job and get client names + estimate status
	jobsWithTotals := make([]JobWithTotal, len(jobs))
	for i, job := range jobs {
		categories, _ := h.queries.ListCategoriesByJob(ctx, repository.ListCategoriesByJobParams{
			JobID: job.ID,
			OrgID: orgID,
		})
		lineItems, _ := h.queries.ListLineItemsByJob(ctx, repository.ListLineItemsByJobParams{
			JobID: job.ID,
			OrgID: orgID,
		})
		customTypes, _ := h.queries.ListJobItemTypes(ctx, repository.ListJobItemTypesParams{
			JobID: job.ID,
			OrgID: orgID,
		})
		totals := h.calculateTotals(job, categories, lineItems, customTypes)

		var clientName string
		if job.ClientID.Valid {
			if client, err := h.queries.GetClient(ctx, repository.GetClientParams{
				ID:    job.ClientID.String,
				OrgID: orgID,
			}); err == nil {
				clientName = client.Name
			}
		} else if job.CustomerName.Valid {
			clientName = job.CustomerName.String
		}

		// Get estimate and signature status
		var estimateCount int64
		var latestEstimateStatus, latestSignatureStatus string

		estimates, _ := h.queries.ListEstimatesByJob(ctx, repository.ListEstimatesByJobParams{
			JobID: job.ID,
			OrgID: orgID,
		})
		estimateCount = int64(len(estimates))
		if len(estimates) > 0 {
			latestEstimateStatus = estimates[0].Status // First is latest (ordered by version DESC)

			// Get signature status for latest estimate
			if sigReq, err := h.queries.GetSignatureRequestByEstimate(ctx, repository.GetSignatureRequestByEstimateParams{
				EstimateID: estimates[0].ID,
				OrgID:      orgID,
			}); err == nil {
				latestSignatureStatus = sigReq.Status
			}
		}

		jobsWithTotals[i] = JobWithTotal{
			Job:                   job,
			GrandTotal:            totals.GrandTotal,
			ClientName:            clientName,
			EstimateCount:         estimateCount,
			LatestEstimateStatus:  latestEstimateStatus,
			LatestSignatureStatus: latestSignatureStatus,
		}
	}

	pagination := PaginationData{
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
	}

	data := map[string]interface{}{
		"Jobs":          jobsWithTotals,
		"SelectedIndex": 0,
		"Pagination":    pagination,
		"Status":        status,
		"Sort":          sortBy,
	}

	if err := h.renderer.Render(w, "jobs_list", data); err != nil {
		logger.Error("failed to render jobs list", "error", err)
	}
}

// GetJob shows a single job with its categories.
func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	jobID := r.PathValue("id")

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

	// Get custom item types for this job (needed for surcharge calculations)
	customTypes, err := h.queries.ListJobItemTypes(ctx, repository.ListJobItemTypesParams{
		JobID: jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to list job item types", "error", err)
		customTypes = []repository.JobItemType{} // Continue with empty list
	}

	totals := h.calculateTotals(job, categories, lineItems, customTypes)

	// Build spreadsheet sections
	sections := h.buildSpreadsheetSections(job, categories, lineItems, customTypes)

	// Get client if associated
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

	// Get estimates with signature status
	estimates, _ := h.queries.ListEstimatesByJob(ctx, repository.ListEstimatesByJobParams{
		JobID: jobID,
		OrgID: orgID,
	})
	estimatesWithStatus := make([]EstimateWithStatus, len(estimates))
	for i, est := range estimates {
		var sigStatus string
		if sigReq, err := h.queries.GetSignatureRequestByEstimate(ctx, repository.GetSignatureRequestByEstimateParams{
			EstimateID: est.ID,
			OrgID:      orgID,
		}); err == nil {
			sigStatus = sigReq.Status
		}
		estimatesWithStatus[i] = EstimateWithStatus{
			Estimate:        est,
			SignatureStatus: sigStatus,
		}
	}

	data := map[string]interface{}{
		"Job":           job,
		"Sections":      sections,
		"Totals":        totals,
		"SelectedIndex": 0,
		"Client":        client,
		"Estimates":     estimatesWithStatus,
		"CustomTypes":   customTypes,
	}

	if err := h.renderer.Render(w, "job", data); err != nil {
		logger.Error("failed to render job page", "error", err)
	}
}

// CreateJob creates a new job and redirects to it.
func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		name = "New Quote"
	}

	clientID := r.FormValue("client_id")

	settings, err := h.queries.GetSettings(ctx, orgID.UUID)
	if err != nil {
		logger.Error("failed to get settings", "error", err)
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	job, err := h.queries.CreateJob(ctx, repository.CreateJobParams{
		ID:               uuid.New().String(),
		OrgID:            orgID,
		Name:             name,
		CustomerName:     sql.NullString{},
		SurchargePercent: settings.DefaultSurchargePercent,
		SurchargeMode:    settings.DefaultSurchargeMode,
		Status:           "draft",
		ExpiresAt:        sql.NullString{},
		ClientID:         toNullString(clientID),
	})
	if err != nil {
		logger.Error("failed to create job", "error", err)
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/jobs/"+job.ID)
		return
	}

	http.Redirect(w, r, "/jobs/"+job.ID, http.StatusSeeOther)
}

// UpdateJob updates a job's details.
func (h *Handler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	jobID := r.PathValue("id")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	surchargePercent, _ := strconv.ParseFloat(r.FormValue("surcharge_percent"), 64)

	customerName := sql.NullString{}
	if cn := r.FormValue("customer_name"); cn != "" {
		customerName = sql.NullString{String: cn, Valid: true}
	}

	// Get existing job to preserve status if not provided
	existingJob, err := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	status := r.FormValue("status")
	if status == "" {
		status = existingJob.Status
	}

	expiresAt := existingJob.ExpiresAt
	if ea := r.FormValue("expires_at"); ea != "" {
		expiresAt = sql.NullString{String: ea, Valid: true}
	}

	clientID := existingJob.ClientID
	if cid := r.FormValue("client_id"); cid != "" {
		clientID = sql.NullString{String: cid, Valid: true}
	} else if cid == "" && r.Form.Has("client_id") {
		// Explicitly cleared
		clientID = sql.NullString{}
	}

	_, err = h.queries.UpdateJob(ctx, repository.UpdateJobParams{
		ID:                        jobID,
		OrgID:                     orgID,
		Name:                      r.FormValue("name"),
		CustomerName:              customerName,
		SurchargePercent:          surchargePercent,
		MaterialSurchargePercent:  existingJob.MaterialSurchargePercent,
		LaborSurchargePercent:     existingJob.LaborSurchargePercent,
		EquipmentSurchargePercent: existingJob.EquipmentSurchargePercent,
		SurchargeMode:             r.FormValue("surcharge_mode"),
		Status:                    status,
		ExpiresAt:                 expiresAt,
		ClientID:                  clientID,
	})
	if err != nil {
		logger.Error("failed to update job", "error", err)
		http.Error(w, "Failed to update job", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/jobs/"+jobID)
		return
	}

	http.Redirect(w, r, "/jobs/"+jobID, http.StatusSeeOther)
}

// DeleteJob deletes a job.
func (h *Handler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	jobID := r.PathValue("id")

	if err := h.queries.DeleteJob(ctx, repository.DeleteJobParams{
		ID:    jobID,
		OrgID: orgID,
	}); err != nil {
		logger.Error("failed to delete job", "error", err)
		http.Error(w, "Failed to delete job", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GetJobForm returns an inline form for creating jobs.
func (h *Handler) GetJobForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)

	// Get clients for dropdown
	clients, err := h.queries.ListClients(ctx, orgID)
	if err != nil {
		logger.Error("failed to list clients", "error", err)
		clients = nil
	}

	data := map[string]interface{}{
		"Clients": clients,
	}

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "job_form", data); err != nil {
		logger.Error("failed to render job form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// GetMarkupForm returns an inline form for editing job markup.
func (h *Handler) GetMarkupForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	jobID := r.PathValue("id")

	job, err := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	// Get custom item types for this job
	customTypes, err := h.queries.ListJobItemTypes(ctx, repository.ListJobItemTypesParams{
		JobID: jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to list job item types", "error", err)
		customTypes = []repository.JobItemType{}
	}

	data := map[string]interface{}{
		"Job":         job,
		"CustomTypes": customTypes,
	}

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "markup_form", data); err != nil {
		logger.Error("failed to render markup form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// GetJobRenameForm returns an inline form for renaming a job.
func (h *Handler) GetJobRenameForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	jobID := r.PathValue("id")

	job, err := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Job": job,
	}

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "job_rename_form", data); err != nil {
		logger.Error("failed to render rename form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// UpdateJobName updates only a job's name.
func (h *Handler) UpdateJobName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	jobID := r.PathValue("id")

	job, err := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		name = job.Name
	}

	_, err = h.queries.UpdateJob(ctx, repository.UpdateJobParams{
		ID:                        jobID,
		OrgID:                     orgID,
		Name:                      name,
		CustomerName:              job.CustomerName,
		SurchargePercent:          job.SurchargePercent,
		MaterialSurchargePercent:  job.MaterialSurchargePercent,
		LaborSurchargePercent:     job.LaborSurchargePercent,
		EquipmentSurchargePercent: job.EquipmentSurchargePercent,
		SurchargeMode:             job.SurchargeMode,
		Status:                    job.Status,
		ExpiresAt:                 job.ExpiresAt,
		ClientID:                  job.ClientID,
	})
	if err != nil {
		logger.Error("failed to update job name", "error", err)
		http.Error(w, "Failed to update name", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/jobs/"+jobID)
		return
	}

	http.Redirect(w, r, "/jobs/"+jobID, http.StatusSeeOther)
}

// UpdateMarkup updates a job's markup percentages (per-type and default).
func (h *Handler) UpdateMarkup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	jobID := r.PathValue("id")

	job, err := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	surchargePercent, _ := strconv.ParseFloat(r.FormValue("surcharge_percent"), 64)

	// Parse per-type surcharges
	var materialSurcharge, laborSurcharge, equipmentSurcharge sql.NullFloat64
	if v := r.FormValue("material_surcharge_percent"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			materialSurcharge = sql.NullFloat64{Float64: f, Valid: true}
		}
	}
	if v := r.FormValue("labor_surcharge_percent"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			laborSurcharge = sql.NullFloat64{Float64: f, Valid: true}
		}
	}
	if v := r.FormValue("equipment_surcharge_percent"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			equipmentSurcharge = sql.NullFloat64{Float64: f, Valid: true}
		}
	}

	_, err = h.queries.UpdateJob(ctx, repository.UpdateJobParams{
		ID:                        jobID,
		OrgID:                     orgID,
		Name:                      job.Name,
		CustomerName:              job.CustomerName,
		SurchargePercent:          surchargePercent,
		MaterialSurchargePercent:  materialSurcharge,
		LaborSurchargePercent:     laborSurcharge,
		EquipmentSurchargePercent: equipmentSurcharge,
		SurchargeMode:             job.SurchargeMode,
		Status:                    job.Status,
		ExpiresAt:                 job.ExpiresAt,
		ClientID:                  job.ClientID,
	})
	if err != nil {
		logger.Error("failed to update job markup", "error", err)
		http.Error(w, "Failed to update markup", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/jobs/"+jobID)
		return
	}

	http.Redirect(w, r, "/jobs/"+jobID, http.StatusSeeOther)
}

// ReportItem represents a single item in a report (materials/equipment only).
type ReportItem struct {
	Name     string
	Quantity float64
	Unit     string
}

// CategoryReport represents a category with its items for the site materials report.
type CategoryReport struct {
	Name  string
	Items []ReportItem
}

// GetOrderList shows an aggregated list of all materials and equipment for a job.
func (h *Handler) GetOrderList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	jobID := r.PathValue("id")

	job, err := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	lineItems, err := h.queries.ListLineItemsByJob(ctx, repository.ListLineItemsByJobParams{
		JobID: jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to list line items", "error", err)
		http.Error(w, "Failed to load items", http.StatusInternalServerError)
		return
	}

	// Aggregate materials and equipment by name+unit
	itemMap := make(map[string]*ReportItem)
	for _, li := range lineItems {
		if li.Type != "material" && li.Type != "equipment" {
			continue
		}
		key := li.Name + "|" + li.Unit
		if existing, ok := itemMap[key]; ok {
			existing.Quantity += li.Quantity
		} else {
			itemMap[key] = &ReportItem{
				Name:     li.Name,
				Quantity: li.Quantity,
				Unit:     li.Unit,
			}
		}
	}

	// Convert to slice and sort alphabetically
	items := make([]ReportItem, 0, len(itemMap))
	for _, item := range itemMap {
		items = append(items, *item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	data := map[string]interface{}{
		"Job":   job,
		"Items": items,
	}

	if err := h.renderer.Render(w, "order_list", data); err != nil {
		logger.Error("failed to render order list", "error", err)
	}
}

// GetJobClientForm returns an inline form for changing the job's client.
func (h *Handler) GetJobClientForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	jobID := r.PathValue("id")

	job, err := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	// Only allow editing client in draft status
	if job.Status != "draft" {
		http.Error(w, "Client can only be changed for draft quotes", http.StatusForbidden)
		return
	}

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
	if err := h.renderer.RenderPartial(&buf, "job_client_form", data); err != nil {
		logger.Error("failed to render client form", "error", err)
		http.Error(w, "Failed to render form", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// UpdateJobClient updates only a job's client assignment.
func (h *Handler) UpdateJobClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	jobID := r.PathValue("id")

	job, err := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	// Only allow editing client in draft status
	if job.Status != "draft" {
		http.Error(w, "Client can only be changed for draft quotes", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	clientID := sql.NullString{}
	if cid := r.FormValue("client_id"); cid != "" {
		clientID = sql.NullString{String: cid, Valid: true}
	}

	_, err = h.queries.UpdateJob(ctx, repository.UpdateJobParams{
		ID:                        jobID,
		OrgID:                     orgID,
		Name:                      job.Name,
		CustomerName:              job.CustomerName,
		SurchargePercent:          job.SurchargePercent,
		MaterialSurchargePercent:  job.MaterialSurchargePercent,
		LaborSurchargePercent:     job.LaborSurchargePercent,
		EquipmentSurchargePercent: job.EquipmentSurchargePercent,
		SurchargeMode:             job.SurchargeMode,
		Status:                    job.Status,
		ExpiresAt:                 job.ExpiresAt,
		ClientID:                  clientID,
	})
	if err != nil {
		logger.Error("failed to update job client", "error", err)
		http.Error(w, "Failed to update client", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/jobs/"+jobID)
		return
	}

	http.Redirect(w, r, "/jobs/"+jobID, http.StatusSeeOther)
}

// GetSiteMaterials shows materials and equipment broken down by category.
func (h *Handler) GetSiteMaterials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	orgID := GetOrgID(ctx)
	jobID := r.PathValue("id")

	job, err := h.queries.GetJob(ctx, repository.GetJobParams{
		ID:    jobID,
		OrgID: orgID,
	})
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

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
		http.Error(w, "Failed to load items", http.StatusInternalServerError)
		return
	}

	// Build category name lookup (with full path)
	categoryNames := make(map[string]string)
	categoryParents := make(map[string]string)
	for _, cat := range categories {
		categoryNames[cat.ID] = cat.Name
		if cat.ParentID.Valid {
			categoryParents[cat.ID] = cat.ParentID.String
		}
	}

	// Build full path for each category
	getFullPath := func(catID string) string {
		parts := []string{}
		currentID := catID
		for currentID != "" {
			if name, ok := categoryNames[currentID]; ok {
				parts = append([]string{name}, parts...)
			}
			currentID = categoryParents[currentID]
		}
		return strings.Join(parts, " > ")
	}

	// Group items by category
	categoryItems := make(map[string][]ReportItem)
	for _, li := range lineItems {
		if li.Type != "material" && li.Type != "equipment" {
			continue
		}
		categoryItems[li.CategoryID] = append(categoryItems[li.CategoryID], ReportItem{
			Name:     li.Name,
			Quantity: li.Quantity,
			Unit:     li.Unit,
		})
	}

	// Build category reports (only categories with items)
	var reports []CategoryReport
	for _, cat := range categories {
		items, hasItems := categoryItems[cat.ID]
		if !hasItems {
			continue
		}
		// Sort items alphabetically
		sort.Slice(items, func(i, j int) bool {
			return items[i].Name < items[j].Name
		})
		reports = append(reports, CategoryReport{
			Name:  getFullPath(cat.ID),
			Items: items,
		})
	}

	// Sort reports by category name
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Name < reports[j].Name
	})

	data := map[string]interface{}{
		"Job":        job,
		"Categories": reports,
	}

	if err := h.renderer.Render(w, "site_materials", data); err != nil {
		logger.Error("failed to render site materials", "error", err)
	}
}
