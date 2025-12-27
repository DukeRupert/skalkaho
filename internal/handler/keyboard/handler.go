package keyboard

import (
	"log/slog"

	"github.com/dukerupert/skalkaho/internal/domain"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/dukerupert/skalkaho/internal/templates/keyboard"
)

// Handler handles keyboard-centric UI HTTP requests.
type Handler struct {
	queries  *repository.Queries
	renderer *keyboard.Renderer
	logger   *slog.Logger
}

// NewHandler creates a new keyboard UI handler.
func NewHandler(queries *repository.Queries, renderer *keyboard.Renderer, logger *slog.Logger) *Handler {
	return &Handler{
		queries:  queries,
		renderer: renderer,
		logger:   logger,
	}
}

// calculateTotals computes job totals from repository types.
func (h *Handler) calculateTotals(job repository.Job, categories []repository.Category, lineItems []repository.LineItem) domain.JobTotal {
	// Convert to domain types
	domainJob := &domain.Job{
		ID:               job.ID,
		SurchargePercent: job.SurchargePercent,
		SurchargeMode:    domain.SurchargeMode(job.SurchargeMode),
	}

	domainCategories := make([]*domain.Category, len(categories))
	for i, cat := range categories {
		var parentID *string
		if cat.ParentID.Valid {
			parentID = &cat.ParentID.String
		}
		var surcharge *float64
		if cat.SurchargePercent.Valid {
			surcharge = &cat.SurchargePercent.Float64
		}
		domainCategories[i] = &domain.Category{
			ID:               cat.ID,
			JobID:            cat.JobID,
			ParentID:         parentID,
			SurchargePercent: surcharge,
		}
	}

	domainLineItems := make([]*domain.LineItem, len(lineItems))
	for i, item := range lineItems {
		var surcharge *float64
		if item.SurchargePercent.Valid {
			surcharge = &item.SurchargePercent.Float64
		}
		domainLineItems[i] = &domain.LineItem{
			ID:               item.ID,
			CategoryID:       item.CategoryID,
			Type:             domain.LineItemType(item.Type),
			Quantity:         item.Quantity,
			UnitPrice:        item.UnitPrice,
			SurchargePercent: surcharge,
		}
	}

	return domain.CalculateJobTotal(domainJob, domainCategories, domainLineItems)
}

// calculateCategoryTotal computes totals for a single category.
func (h *Handler) calculateCategoryTotal(categoryID string, job repository.Job, categories []repository.Category, lineItems []repository.LineItem) domain.CategoryTotal {
	domainJob := &domain.Job{
		ID:               job.ID,
		SurchargePercent: job.SurchargePercent,
		SurchargeMode:    domain.SurchargeMode(job.SurchargeMode),
	}

	domainCategories := make([]*domain.Category, len(categories))
	for i, cat := range categories {
		var parentID *string
		if cat.ParentID.Valid {
			parentID = &cat.ParentID.String
		}
		var surcharge *float64
		if cat.SurchargePercent.Valid {
			surcharge = &cat.SurchargePercent.Float64
		}
		domainCategories[i] = &domain.Category{
			ID:               cat.ID,
			JobID:            cat.JobID,
			ParentID:         parentID,
			SurchargePercent: surcharge,
		}
	}

	domainLineItems := make([]*domain.LineItem, len(lineItems))
	for i, item := range lineItems {
		var surcharge *float64
		if item.SurchargePercent.Valid {
			surcharge = &item.SurchargePercent.Float64
		}
		domainLineItems[i] = &domain.LineItem{
			ID:               item.ID,
			CategoryID:       item.CategoryID,
			Type:             domain.LineItemType(item.Type),
			Quantity:         item.Quantity,
			UnitPrice:        item.UnitPrice,
			SurchargePercent: surcharge,
		}
	}

	return domain.CalculateCategoryTotal(categoryID, domainJob, domainCategories, domainLineItems)
}

// getCategoryDepth returns the depth of a category (1 = top level).
func (h *Handler) getCategoryDepth(categories []repository.Category, categoryID string) int {
	categoryByID := make(map[string]repository.Category)
	for _, cat := range categories {
		categoryByID[cat.ID] = cat
	}

	depth := 1
	current := categoryByID[categoryID]
	for current.ParentID.Valid {
		depth++
		current = categoryByID[current.ParentID.String]
	}
	return depth
}

// getBreadcrumbs builds the breadcrumb trail for a category.
func (h *Handler) getBreadcrumbs(categories []repository.Category, categoryID string, job repository.Job) []Breadcrumb {
	categoryByID := make(map[string]repository.Category)
	for _, cat := range categories {
		categoryByID[cat.ID] = cat
	}

	var trail []Breadcrumb
	current := categoryByID[categoryID]

	// Build trail from current to root
	for {
		trail = append([]Breadcrumb{{
			ID:   current.ID,
			Name: current.Name,
			Type: "category",
		}}, trail...)

		if !current.ParentID.Valid {
			break
		}
		current = categoryByID[current.ParentID.String]
	}

	// Prepend job
	trail = append([]Breadcrumb{{
		ID:   job.ID,
		Name: job.Name,
		Type: "job",
	}}, trail...)

	return trail
}

// Breadcrumb represents a navigation breadcrumb.
type Breadcrumb struct {
	ID   string
	Name string
	Type string // "job" or "category"
}

// helper to check if a category can have subcategories
func canAddSubcategory(depth int) bool {
	return depth < 3
}
