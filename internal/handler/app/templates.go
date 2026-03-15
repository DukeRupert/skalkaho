package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/dukerupert/skalkaho/internal/repository"
)

// ── Data types for template pages ──

// TemplateTreeNode represents a section with its children for the tree editor.
type TemplateTreeNode struct {
	Section       repository.TemplateSection
	Subcategories []TemplateSubcategoryNode
}

// TemplateSubcategoryNode represents a subcategory with its component groups.
type TemplateSubcategoryNode struct {
	Subcategory     repository.TemplateSubcategory
	ComponentGroups []repository.TemplateComponentGroup
}

// TemplatesPageData holds data for the templates list page.
type TemplatesPageData struct {
	PageData
	Templates []TemplateListItem
}

// TemplateListItem extends Template with section count.
type TemplateListItem struct {
	repository.Template
	SectionCount int64
}

// TemplateDetailData holds data for the template detail/edit page.
type TemplateDetailData struct {
	PageData
	Template repository.Template
	Tree     []TemplateTreeNode
}

func newID() string {
	return uuid.New().String()[:20]
}

// ── List / CRUD ──

// ListTemplates renders the templates list page.
func (h *Handler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	templates, err := h.queries.ListTemplates(ctx)
	if err != nil {
		h.logger.Error("listing templates", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	items := make([]TemplateListItem, len(templates))
	for i, t := range templates {
		count, err := h.queries.CountTemplateSections(ctx, t.ID)
		if err != nil {
			count = 0
		}
		items[i] = TemplateListItem{Template: t, SectionCount: count}
	}

	data := TemplatesPageData{
		PageData:  h.pageData(r, "templates"),
		Templates: items,
	}

	if err := h.renderer.Render(w, "templates.html", data); err != nil {
		h.logger.Error("rendering templates", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// NewTemplateModal renders the new template form modal.
func (h *Handler) NewTemplateModal(w http.ResponseWriter, r *http.Request) {
	if err := h.renderer.RenderPartial(w, "templates.html", "new-template-modal", nil); err != nil {
		h.logger.Error("rendering new template modal", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// CreateTemplate handles POST /templates.
func (h *Handler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Template name is required", http.StatusBadRequest)
		return
	}

	id := newID()
	_, err := h.queries.CreateTemplate(r.Context(), repository.CreateTemplateParams{
		ID:          id,
		Name:        name,
		Description: toNullString(r.FormValue("description")),
	})
	if err != nil {
		h.logger.Error("creating template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/templates/%s", id), http.StatusSeeOther)
}

// GetTemplate renders the template detail/tree editor page.
func (h *Handler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	tmpl, err := h.queries.GetTemplate(ctx, id)
	if err != nil {
		h.logger.Error("getting template", "error", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	tree, err := h.loadTemplateTree(ctx, id)
	if err != nil {
		h.logger.Error("loading template tree", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := TemplateDetailData{
		PageData: h.pageData(r, "templates"),
		Template: tmpl,
		Tree:     tree,
	}

	if err := h.renderer.Render(w, "template_detail.html", data); err != nil {
		h.logger.Error("rendering template detail", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// GetTemplateEditModal renders the edit template modal.
func (h *Handler) GetTemplateEditModal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	tmpl, err := h.queries.GetTemplate(r.Context(), id)
	if err != nil {
		h.logger.Error("getting template for edit modal", "error", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	data := struct {
		Template repository.Template
	}{Template: tmpl}

	if err := h.renderer.RenderPartial(w, "template_detail.html", "edit-template-modal", data); err != nil {
		h.logger.Error("rendering edit template modal", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// UpdateTemplate handles POST /templates/{id}.
func (h *Handler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Template name is required", http.StatusBadRequest)
		return
	}

	_, err := h.queries.UpdateTemplate(r.Context(), repository.UpdateTemplateParams{
		ID:          id,
		Name:        name,
		Description: toNullString(r.FormValue("description")),
	})
	if err != nil {
		h.logger.Error("updating template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/templates/%s", id), http.StatusSeeOther)
}

// DeleteTemplate handles DELETE /templates/{id}.
func (h *Handler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.queries.DeleteTemplate(r.Context(), id); err != nil {
		h.logger.Error("deleting template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if isHTMXPartial(r) {
		w.Header().Set("HX-Redirect", "/templates")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/templates", http.StatusSeeOther)
}

// ── Sections ──

// CreateTemplateSection handles POST /templates/{id}/sections.
func (h *Handler) CreateTemplateSection(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	templateID := r.PathValue("id")
	name := r.FormValue("name")
	if name == "" {
		name = "New Section"
	}

	sortOrder := parseSortOrder(r.FormValue("sort_order"))

	_, err := h.queries.CreateTemplateSection(r.Context(), repository.CreateTemplateSectionParams{
		ID:         newID(),
		TemplateID: templateID,
		Name:       name,
		SortOrder:  sortOrder,
	})
	if err != nil {
		h.logger.Error("creating template section", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.renderTemplateTree(w, r, templateID)
}

// UpdateTemplateSection handles POST /templates/{id}/sections/{sid}.
func (h *Handler) UpdateTemplateSection(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	sid := r.PathValue("sid")
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	sortOrder := parseSortOrder(r.FormValue("sort_order"))

	_, err := h.queries.UpdateTemplateSection(r.Context(), repository.UpdateTemplateSectionParams{
		ID:        sid,
		Name:      name,
		SortOrder: sortOrder,
	})
	if err != nil {
		h.logger.Error("updating template section", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.renderTemplateTree(w, r, r.PathValue("id"))
}

// DeleteTemplateSection handles DELETE /templates/{id}/sections/{sid}.
func (h *Handler) DeleteTemplateSection(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	if err := h.queries.DeleteTemplateSection(r.Context(), sid); err != nil {
		h.logger.Error("deleting template section", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.renderTemplateTree(w, r, r.PathValue("id"))
}

// ── Subcategories ──

// CreateTemplateSubcategory handles POST /templates/{id}/sections/{sid}/subcategories.
func (h *Handler) CreateTemplateSubcategory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	sectionID := r.PathValue("sid")
	name := r.FormValue("name")
	if name == "" {
		name = "New Subcategory"
	}

	sortOrder := parseSortOrder(r.FormValue("sort_order"))

	_, err := h.queries.CreateTemplateSubcategory(r.Context(), repository.CreateTemplateSubcategoryParams{
		ID:                newID(),
		TemplateSectionID: sectionID,
		Name:              name,
		SortOrder:         sortOrder,
	})
	if err != nil {
		h.logger.Error("creating template subcategory", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.renderTemplateTree(w, r, r.PathValue("id"))
}

// UpdateTemplateSubcategory handles POST /templates/{id}/subcategories/{scid}.
func (h *Handler) UpdateTemplateSubcategory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	scid := r.PathValue("scid")
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	sortOrder := parseSortOrder(r.FormValue("sort_order"))

	_, err := h.queries.UpdateTemplateSubcategory(r.Context(), repository.UpdateTemplateSubcategoryParams{
		ID:        scid,
		Name:      name,
		SortOrder: sortOrder,
	})
	if err != nil {
		h.logger.Error("updating template subcategory", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.renderTemplateTree(w, r, r.PathValue("id"))
}

// DeleteTemplateSubcategory handles DELETE /templates/{id}/subcategories/{scid}.
func (h *Handler) DeleteTemplateSubcategory(w http.ResponseWriter, r *http.Request) {
	scid := r.PathValue("scid")
	if err := h.queries.DeleteTemplateSubcategory(r.Context(), scid); err != nil {
		h.logger.Error("deleting template subcategory", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.renderTemplateTree(w, r, r.PathValue("id"))
}

// ── Component Groups ──

// CreateTemplateComponentGroup handles POST /templates/{id}/subcategories/{scid}/groups.
func (h *Handler) CreateTemplateComponentGroup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	subcategoryID := r.PathValue("scid")
	name := r.FormValue("name")
	if name == "" {
		name = "New Group"
	}

	sortOrder := parseSortOrder(r.FormValue("sort_order"))

	_, err := h.queries.CreateTemplateComponentGroup(r.Context(), repository.CreateTemplateComponentGroupParams{
		ID:                    newID(),
		TemplateSubcategoryID: subcategoryID,
		Name:                  name,
		SortOrder:             sortOrder,
	})
	if err != nil {
		h.logger.Error("creating template component group", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.renderTemplateTree(w, r, r.PathValue("id"))
}

// UpdateTemplateComponentGroup handles POST /templates/{id}/groups/{gid}.
func (h *Handler) UpdateTemplateComponentGroup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	gid := r.PathValue("gid")
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	sortOrder := parseSortOrder(r.FormValue("sort_order"))

	_, err := h.queries.UpdateTemplateComponentGroup(r.Context(), repository.UpdateTemplateComponentGroupParams{
		ID:        gid,
		Name:      name,
		SortOrder: sortOrder,
	})
	if err != nil {
		h.logger.Error("updating template component group", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.renderTemplateTree(w, r, r.PathValue("id"))
}

// DeleteTemplateComponentGroup handles DELETE /templates/{id}/groups/{gid}.
func (h *Handler) DeleteTemplateComponentGroup(w http.ResponseWriter, r *http.Request) {
	gid := r.PathValue("gid")
	if err := h.queries.DeleteTemplateComponentGroup(r.Context(), gid); err != nil {
		h.logger.Error("deleting template component group", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.renderTemplateTree(w, r, r.PathValue("id"))
}

// ── Stamp Function ──

// stampTemplate copies a template's structure into an existing project.
// All inserts run inside the provided transaction.
func (h *Handler) stampTemplate(ctx context.Context, tx *sql.Tx, templateID, projectID string) error {
	qtx := h.queries.WithTx(tx)

	sections, err := qtx.ListTemplateSections(ctx, templateID)
	if err != nil {
		return fmt.Errorf("listing template sections: %w", err)
	}

	for _, sec := range sections {
		newSectionID := newID()
		_, err := qtx.CreateSection(ctx, repository.CreateSectionParams{
			ID:        newSectionID,
			ProjectID: projectID,
			Name:      sec.Name,
			SortOrder: sec.SortOrder,
		})
		if err != nil {
			return fmt.Errorf("creating section %q: %w", sec.Name, err)
		}

		subcats, err := qtx.ListTemplateSubcategories(ctx, sec.ID)
		if err != nil {
			return fmt.Errorf("listing template subcategories: %w", err)
		}

		for _, sub := range subcats {
			newSubID := newID()
			_, err := qtx.CreateSubcategory(ctx, repository.CreateSubcategoryParams{
				ID:                     newSubID,
				SectionID:              newSectionID,
				Name:                   sub.Name,
				SortOrder:              sub.SortOrder,
				LumpSum:                0,
				MaterialsMarkupEnabled: true,
				LaborMarkupEnabled:     true,
				EquipmentMarkupEnabled: true,
				SubsMarkupEnabled:      true,
				OtherMarkupEnabled:     true,
			})
			if err != nil {
				return fmt.Errorf("creating subcategory %q: %w", sub.Name, err)
			}

			groups, err := qtx.ListTemplateComponentGroups(ctx, sub.ID)
			if err != nil {
				return fmt.Errorf("listing template component groups: %w", err)
			}

			for _, grp := range groups {
				_, err := qtx.CreateComponentGroup(ctx, repository.CreateComponentGroupParams{
					ID:            newID(),
					SubcategoryID: newSubID,
					Name:          grp.Name,
					SortOrder:     grp.SortOrder,
				})
				if err != nil {
					return fmt.Errorf("creating component group %q: %w", grp.Name, err)
				}
			}
		}
	}

	return nil
}

// ── Helpers ──

// loadTemplateTree loads the full section → subcategory → component group tree.
func (h *Handler) loadTemplateTree(ctx context.Context, templateID string) ([]TemplateTreeNode, error) {
	sections, err := h.queries.ListTemplateSections(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("listing sections: %w", err)
	}

	tree := make([]TemplateTreeNode, len(sections))
	for i, sec := range sections {
		subcats, err := h.queries.ListTemplateSubcategories(ctx, sec.ID)
		if err != nil {
			return nil, fmt.Errorf("listing subcategories: %w", err)
		}

		subNodes := make([]TemplateSubcategoryNode, len(subcats))
		for j, sub := range subcats {
			groups, err := h.queries.ListTemplateComponentGroups(ctx, sub.ID)
			if err != nil {
				return nil, fmt.Errorf("listing component groups: %w", err)
			}
			subNodes[j] = TemplateSubcategoryNode{
				Subcategory:     sub,
				ComponentGroups: groups,
			}
		}

		tree[i] = TemplateTreeNode{
			Section:       sec,
			Subcategories: subNodes,
		}
	}

	return tree, nil
}

// renderTemplateTree re-renders the tree partial after a mutation.
func (h *Handler) renderTemplateTree(w http.ResponseWriter, r *http.Request, templateID string) {
	tmpl, err := h.queries.GetTemplate(r.Context(), templateID)
	if err != nil {
		h.logger.Error("getting template for tree render", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tree, err := h.loadTemplateTree(r.Context(), templateID)
	if err != nil {
		h.logger.Error("loading template tree", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := TemplateDetailData{
		Template: tmpl,
		Tree:     tree,
	}

	if err := h.renderer.RenderPartial(w, "template_detail.html", "template-tree", data); err != nil {
		h.logger.Error("rendering template tree", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func parseSortOrder(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}
