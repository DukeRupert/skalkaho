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

// JobWithTotal wraps a Job with its calculated grand total.
type JobWithTotal struct {
	repository.Job
	GrandTotal float64
}

// PaginationData holds pagination state for templates.
type PaginationData struct {
	CurrentPage int
	TotalPages  int
	TotalItems  int64
	HasPrev     bool
	HasNext     bool
}

// ListJobs shows the keyboard-centric jobs list with pagination and filtering.
func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

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
	totalItems, err := h.queries.CountJobs(ctx, status)
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
	params := repository.ListJobsPaginatedParams{
		Status: status,
		Offset: offset,
		Limit:  pageSize,
	}

	switch sortBy {
	case "oldest":
		jobs, err = h.queries.ListJobsPaginatedOldest(ctx, repository.ListJobsPaginatedOldestParams{
			Status: status,
			Offset: offset,
			Limit:  pageSize,
		})
	case "name_asc":
		jobs, err = h.queries.ListJobsPaginatedByName(ctx, repository.ListJobsPaginatedByNameParams{
			Status: status,
			Offset: offset,
			Limit:  pageSize,
		})
	case "name_desc":
		jobs, err = h.queries.ListJobsPaginatedByNameDesc(ctx, repository.ListJobsPaginatedByNameDescParams{
			Status: status,
			Offset: offset,
			Limit:  pageSize,
		})
	default: // newest
		jobs, err = h.queries.ListJobsPaginated(ctx, params)
	}

	if err != nil {
		logger.Error("failed to list jobs", "error", err)
		http.Error(w, "Failed to load jobs", http.StatusInternalServerError)
		return
	}

	// Calculate totals for each job
	jobsWithTotals := make([]JobWithTotal, len(jobs))
	for i, job := range jobs {
		categories, _ := h.queries.ListCategoriesByJob(ctx, job.ID)
		lineItems, _ := h.queries.ListLineItemsByJob(ctx, job.ID)
		totals := h.calculateTotals(job, categories, lineItems)
		jobsWithTotals[i] = JobWithTotal{
			Job:        job,
			GrandTotal: totals.GrandTotal,
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
	jobID := r.PathValue("id")

	job, err := h.queries.GetJob(ctx, jobID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Failed to load job", http.StatusInternalServerError)
		return
	}

	categories, err := h.queries.ListCategoriesByJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to list categories", "error", err)
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	lineItems, err := h.queries.ListLineItemsByJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to list line items", "error", err)
		http.Error(w, "Failed to load line items", http.StatusInternalServerError)
		return
	}

	// Get only top-level categories
	topLevelCategories := make([]repository.Category, 0)
	for _, cat := range categories {
		if !cat.ParentID.Valid {
			topLevelCategories = append(topLevelCategories, cat)
		}
	}

	// Calculate totals for each category
	type CategoryWithTotal struct {
		repository.Category
		Total float64
	}

	categoriesWithTotals := make([]CategoryWithTotal, len(topLevelCategories))
	for i, cat := range topLevelCategories {
		catTotal := h.calculateCategoryTotal(cat.ID, job, categories, lineItems)
		categoriesWithTotals[i] = CategoryWithTotal{
			Category: cat,
			Total:    catTotal.Total,
		}
	}

	totals := h.calculateTotals(job, categories, lineItems)

	// Build category tree for sidebar navigation
	categoryTree := buildCategoryTree(categories)

	data := map[string]interface{}{
		"Job":               job,
		"Categories":        categoriesWithTotals,
		"Totals":            totals,
		"SelectedIndex":     0,
		"CategoryTree":      categoryTree,
		"CurrentCategoryID": "",
	}

	if err := h.renderer.Render(w, "job", data); err != nil {
		logger.Error("failed to render job page", "error", err)
	}
}

// CreateJob creates a new job and redirects to it.
func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		name = "New Quote"
	}

	settings, err := h.queries.GetSettings(ctx)
	if err != nil {
		logger.Error("failed to get settings", "error", err)
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	job, err := h.queries.CreateJob(ctx, repository.CreateJobParams{
		ID:               uuid.New().String(),
		Name:             name,
		CustomerName:     sql.NullString{},
		SurchargePercent: settings.DefaultSurchargePercent,
		SurchargeMode:    settings.DefaultSurchargeMode,
		Status:           "draft",
		ExpiresAt:        sql.NullString{},
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
	existingJob, err := h.queries.GetJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	status := r.FormValue("status")
	if status == "" {
		status = existingJob.Status
	}

	expiresAt := sql.NullString{}
	if ea := r.FormValue("expires_at"); ea != "" {
		expiresAt = sql.NullString{String: ea, Valid: true}
	} else {
		expiresAt = existingJob.ExpiresAt
	}

	_, err = h.queries.UpdateJob(ctx, repository.UpdateJobParams{
		ID:               jobID,
		Name:             r.FormValue("name"),
		CustomerName:     customerName,
		SurchargePercent: surchargePercent,
		SurchargeMode:    r.FormValue("surcharge_mode"),
		Status:           status,
		ExpiresAt:        expiresAt,
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
	jobID := r.PathValue("id")

	if err := h.queries.DeleteJob(ctx, jobID); err != nil {
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

	var buf bytes.Buffer
	if err := h.renderer.RenderPartial(&buf, "job_form", nil); err != nil {
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
	jobID := r.PathValue("id")

	job, err := h.queries.GetJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Job": job,
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
	jobID := r.PathValue("id")

	job, err := h.queries.GetJob(ctx, jobID)
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
	jobID := r.PathValue("id")

	job, err := h.queries.GetJob(ctx, jobID)
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
		ID:               jobID,
		Name:             name,
		CustomerName:     job.CustomerName,
		SurchargePercent: job.SurchargePercent,
		SurchargeMode:    job.SurchargeMode,
		Status:           job.Status,
		ExpiresAt:        job.ExpiresAt,
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

// UpdateMarkup updates a job's markup percentage.
func (h *Handler) UpdateMarkup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	jobID := r.PathValue("id")

	job, err := h.queries.GetJob(ctx, jobID)
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

	_, err = h.queries.UpdateJob(ctx, repository.UpdateJobParams{
		ID:               jobID,
		Name:             job.Name,
		CustomerName:     job.CustomerName,
		SurchargePercent: surchargePercent,
		SurchargeMode:    job.SurchargeMode,
		Status:           job.Status,
		ExpiresAt:        job.ExpiresAt,
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
	jobID := r.PathValue("id")

	job, err := h.queries.GetJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	lineItems, err := h.queries.ListLineItemsByJob(ctx, jobID)
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

// GetSiteMaterials shows materials and equipment broken down by category.
func (h *Handler) GetSiteMaterials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := middleware.LoggerFromContext(ctx)
	jobID := r.PathValue("id")

	job, err := h.queries.GetJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to get job", "error", err)
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	categories, err := h.queries.ListCategoriesByJob(ctx, jobID)
	if err != nil {
		logger.Error("failed to list categories", "error", err)
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	lineItems, err := h.queries.ListLineItemsByJob(ctx, jobID)
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
