package app

import (
	"fmt"
	"net/http"

	"github.com/dukerupert/skalkaho/internal/domain"
	"github.com/dukerupert/skalkaho/internal/repository"
)

// OverviewPageData holds all data for the project overview page.
type OverviewPageData struct {
	PageData
	FullProject repository.Project
	Client      *repository.Client
	CostSummary domain.ProjectCostSummary
	Sections    []domain.EstimateSection
}

// GetProjectOverview renders the full project overview page with cost breakdown.
func (h *Handler) GetProjectOverview(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	ctx := r.Context()

	project, err := h.queries.GetProject(ctx, projectID)
	if err != nil {
		h.logger.Error("getting project", "error", err)
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Load client if linked
	var client *repository.Client
	if project.ClientID.Valid {
		c, err := h.queries.GetClient(ctx, project.ClientID.String)
		if err == nil {
			client = &c
		}
	}

	// Load estimate data and calculate costs using shared function
	sections, err := h.loadEstimateSections(r, projectID)
	if err != nil {
		h.logger.Error("loading estimate sections", "error", err)
		// Don't fail — just show empty costs
		sections = nil
	}

	globals := domain.MarkupGlobals{
		MaterialsMarkup: project.MaterialsMarkup,
		LaborMarkup:     project.LaborMarkup,
		EquipmentMarkup: project.EquipmentMarkup,
		SubsMarkup:      project.SubsMarkup,
		OtherMarkup:     project.OtherMarkup,
	}
	costSummary := domain.CalculateProjectCosts(sections, globals)

	data := OverviewPageData{
		PageData:    h.pageData(r, "overview"),
		FullProject: project,
		Client:      client,
		CostSummary: costSummary,
		Sections:    sections,
	}
	data.Project = &ProjectStub{ID: project.ID, Name: project.Name}

	if err := h.renderer.Render(w, "project_overview.html", data); err != nil {
		h.logger.Error("rendering project overview", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// GetStatusModal returns the status selector modal HTML.
func (h *Handler) GetStatusModal(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	project, err := h.queries.GetProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	data := struct {
		ProjectID string
		Status    string
	}{
		ProjectID: project.ID,
		Status:    project.Status,
	}

	if err := h.renderer.RenderPartial(w, "project_overview.html", "status-modal", data); err != nil {
		h.logger.Error("rendering status modal", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// loadEstimateSections loads the full section hierarchy for cost calculation.
func (h *Handler) loadEstimateSections(r *http.Request, projectID string) ([]domain.EstimateSection, error) {
	ctx := r.Context()

	dbSections, err := h.queries.ListSectionsByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("listing sections: %w", err)
	}

	sections := make([]domain.EstimateSection, 0, len(dbSections))
	for _, s := range dbSections {
		subcats, err := h.buildSubcategories(r, s.ID)
		if err != nil {
			return nil, fmt.Errorf("building subcategories for section %s: %w", s.ID, err)
		}
		sections = append(sections, domain.EstimateSection{
			ID:            s.ID,
			Name:          s.Name,
			SortOrder:     int(s.SortOrder),
			Subcategories: subcats,
		})
	}

	return sections, nil
}
