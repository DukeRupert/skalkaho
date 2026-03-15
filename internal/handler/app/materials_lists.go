package app

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/dukerupert/skalkaho/internal/repository"
)

// ── View models ────────────────────────────────────────────────

// MaterialsListIndexData holds data for the materials list index page.
type MaterialsListIndexData struct {
	PageData
	Quote     repository.Quote
	Lists     []MaterialsListView
	Sections  []repository.ListSectionsForProjectRow
	ProjectID string
}

// MaterialsListView is a template-friendly view of a materials list.
type MaterialsListView struct {
	ID          string
	Label       string
	GeneratedAt string
	Notes       string
}

// MaterialsListShowData holds data for the materials list detail/print page.
type MaterialsListShowData struct {
	PageData
	List           repository.GetMaterialsListRow
	Quote          repository.Quote
	QuoteID        string
	ProjectName    string
	SectionGroups  []SectionGroup
	GeneratedAtFmt string
	IsAllSections  bool
}

// SectionGroup groups items by section for the full-quote view.
type SectionGroup struct {
	SectionName      string
	SubcategoryGroups []SubcategoryGroup
}

// SubcategoryGroup groups items by subcategory.
type SubcategoryGroup struct {
	SubcategoryName string
	Items           []MaterialsListItem
}

// MaterialsListItem is a single line item for display.
type MaterialsListItem struct {
	ItemName string
	Quantity float64
	Unit     string
}

// ── Helpers ────────────────────────────────────────────────────

func toMaterialsListViews(rows []repository.ListMaterialsListsByQuoteRow) []MaterialsListView {
	views := make([]MaterialsListView, len(rows))
	for i, r := range rows {
		views[i] = MaterialsListView{
			ID:          r.ID,
			Label:       r.Label,
			GeneratedAt: r.GeneratedAt.Format("Jan 2, 2006"),
			Notes:       r.Notes.String,
		}
	}
	return views
}

// groupAllItems groups full-quote items by section → subcategory.
func groupAllItems(items []repository.GetMaterialsListItemsRow) []SectionGroup {
	var groups []SectionGroup
	var curSection *SectionGroup
	var curSubcat *SubcategoryGroup

	for _, item := range items {
		// New section?
		if curSection == nil || curSection.SectionName != item.SectionName {
			if curSubcat != nil && curSection != nil {
				curSection.SubcategoryGroups = append(curSection.SubcategoryGroups, *curSubcat)
			}
			if curSection != nil {
				groups = append(groups, *curSection)
			}
			curSection = &SectionGroup{SectionName: item.SectionName}
			curSubcat = nil
		}
		// New subcategory?
		if curSubcat == nil || curSubcat.SubcategoryName != item.SubcategoryName {
			if curSubcat != nil {
				curSection.SubcategoryGroups = append(curSection.SubcategoryGroups, *curSubcat)
			}
			curSubcat = &SubcategoryGroup{SubcategoryName: item.SubcategoryName}
		}
		curSubcat.Items = append(curSubcat.Items, MaterialsListItem{
			ItemName: item.ItemName,
			Quantity: item.Quantity,
			Unit:     item.Unit,
		})
	}
	if curSubcat != nil && curSection != nil {
		curSection.SubcategoryGroups = append(curSection.SubcategoryGroups, *curSubcat)
	}
	if curSection != nil {
		groups = append(groups, *curSection)
	}
	return groups
}

// groupSectionItems groups single-section items by subcategory only.
func groupSectionItems(items []repository.GetMaterialsListItemsBySectionRow) []SubcategoryGroup {
	var groups []SubcategoryGroup
	var cur *SubcategoryGroup

	for _, item := range items {
		if cur == nil || cur.SubcategoryName != item.SubcategoryName {
			if cur != nil {
				groups = append(groups, *cur)
			}
			cur = &SubcategoryGroup{SubcategoryName: item.SubcategoryName}
		}
		cur.Items = append(cur.Items, MaterialsListItem{
			ItemName: item.ItemName,
			Quantity: item.Quantity,
			Unit:     item.Unit,
		})
	}
	if cur != nil {
		groups = append(groups, *cur)
	}
	return groups
}

// ── Materials List Index ──────────────────────────────────────

// MaterialsListIndex shows all materials lists for a quote.
func (h *Handler) MaterialsListIndex(w http.ResponseWriter, r *http.Request) {
	quoteID := r.PathValue("id")
	ctx := r.Context()

	quote, err := h.queries.GetQuote(ctx, quoteID)
	if err != nil {
		h.logger.Error("getting quote", "error", err)
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	project, err := h.queries.GetProject(ctx, quote.ProjectID)
	if err != nil {
		h.logger.Error("getting project", "error", err)
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	lists, err := h.queries.ListMaterialsListsByQuote(ctx, quoteID)
	if err != nil {
		h.logger.Error("listing materials lists", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	sections, err := h.queries.ListSectionsForProject(ctx, quote.ProjectID)
	if err != nil {
		h.logger.Error("listing sections for project", "error", err)
		sections = nil
	}

	data := MaterialsListIndexData{
		PageData:  h.pageData(r, "projects"),
		Quote:     quote,
		Lists:     toMaterialsListViews(lists),
		Sections:  sections,
		ProjectID: quote.ProjectID,
	}
	data.Project = &ProjectStub{ID: project.ID, Name: project.Name}

	if err := h.renderer.Render(w, "materials_lists.html", data); err != nil {
		h.logger.Error("rendering materials list index", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ── Generate handlers ──────────────────────────────────────────

// GenerateAllMaterialsList creates one materials list covering the entire quote.
func (h *Handler) GenerateAllMaterialsList(w http.ResponseWriter, r *http.Request) {
	quoteID := r.PathValue("id")
	ctx := r.Context()

	quote, err := h.queries.GetQuote(ctx, quoteID)
	if err != nil {
		h.logger.Error("getting quote", "error", err)
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	if quote.Status != "signed" {
		http.Error(w, "Materials lists can only be generated for signed quotes", http.StatusBadRequest)
		return
	}

	_, err = h.queries.CreateMaterialsList(ctx, repository.CreateMaterialsListParams{
		ID:      newID(),
		QuoteID: quoteID,
		Label:   "All Categories",
	})
	if err != nil {
		h.logger.Error("creating materials list", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/quotes/%s/materials-lists", quoteID), http.StatusSeeOther)
}

// GeneratePerSectionMaterialsList creates one materials list per section on the quote.
func (h *Handler) GeneratePerSectionMaterialsList(w http.ResponseWriter, r *http.Request) {
	quoteID := r.PathValue("id")
	ctx := r.Context()

	quote, err := h.queries.GetQuote(ctx, quoteID)
	if err != nil {
		h.logger.Error("getting quote", "error", err)
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	if quote.Status != "signed" {
		http.Error(w, "Materials lists can only be generated for signed quotes", http.StatusBadRequest)
		return
	}

	sections, err := h.queries.ListSectionsForProject(ctx, quote.ProjectID)
	if err != nil {
		h.logger.Error("listing sections", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, section := range sections {
		_, err := h.queries.CreateMaterialsList(ctx, repository.CreateMaterialsListParams{
			ID:        newID(),
			QuoteID:   quoteID,
			SectionID: sql.NullString{String: section.ID, Valid: true},
			Label:     section.Name,
		})
		if err != nil {
			h.logger.Error("creating per-section materials list", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/quotes/%s/materials-lists", quoteID), http.StatusSeeOther)
}

// GenerateForSectionMaterialsList creates a materials list for a single section.
func (h *Handler) GenerateForSectionMaterialsList(w http.ResponseWriter, r *http.Request) {
	quoteID := r.PathValue("id")
	sectionID := r.PathValue("sid")
	ctx := r.Context()

	quote, err := h.queries.GetQuote(ctx, quoteID)
	if err != nil {
		h.logger.Error("getting quote", "error", err)
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	if quote.Status != "signed" {
		http.Error(w, "Materials lists can only be generated for signed quotes", http.StatusBadRequest)
		return
	}

	// Verify the section belongs to this project
	section, err := h.queries.GetSection(ctx, sectionID)
	if err != nil {
		h.logger.Error("getting section", "error", err)
		http.Error(w, "Section not found", http.StatusNotFound)
		return
	}

	if section.ProjectID != quote.ProjectID {
		http.Error(w, "Section does not belong to this project", http.StatusBadRequest)
		return
	}

	list, err := h.queries.CreateMaterialsList(ctx, repository.CreateMaterialsListParams{
		ID:        newID(),
		QuoteID:   quoteID,
		SectionID: sql.NullString{String: sectionID, Valid: true},
		Label:     section.Name,
	})
	if err != nil {
		h.logger.Error("creating materials list", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/quotes/%s/materials-lists/%s", quoteID, list.ID), http.StatusSeeOther)
}

// ── Show / Print ───────────────────────────────────────────────

// MaterialsListShow renders a single materials list (view + print).
func (h *Handler) MaterialsListShow(w http.ResponseWriter, r *http.Request) {
	quoteID := r.PathValue("id")
	listID := r.PathValue("lid")
	ctx := r.Context()

	list, err := h.queries.GetMaterialsList(ctx, listID)
	if err != nil {
		h.logger.Error("getting materials list", "error", err)
		http.Error(w, "Materials list not found", http.StatusNotFound)
		return
	}

	if list.QuoteID != quoteID {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	quote, err := h.queries.GetQuote(ctx, quoteID)
	if err != nil {
		h.logger.Error("getting quote", "error", err)
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	project, err := h.queries.GetProject(ctx, quote.ProjectID)
	if err != nil {
		h.logger.Error("getting project", "error", err)
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	isAllSections := !list.SectionID.Valid
	var sectionGroups []SectionGroup

	if isAllSections {
		// Full quote — group by section → subcategory
		items, err := h.queries.GetMaterialsListItems(ctx, quote.ProjectID)
		if err != nil {
			h.logger.Error("getting materials list items", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		sectionGroups = groupAllItems(items)
	} else {
		// Single section — group by subcategory only
		items, err := h.queries.GetMaterialsListItemsBySection(ctx, list.SectionID.String)
		if err != nil {
			h.logger.Error("getting materials list items by section", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		subcatGroups := groupSectionItems(items)
		sectionGroups = []SectionGroup{
			{SectionName: list.Label, SubcategoryGroups: subcatGroups},
		}
	}

	data := MaterialsListShowData{
		PageData:       h.pageData(r, "projects"),
		List:           list,
		Quote:          quote,
		QuoteID:        quoteID,
		ProjectName:    project.Name,
		SectionGroups:  sectionGroups,
		GeneratedAtFmt: list.GeneratedAt.Format("January 2, 2006"),
		IsAllSections:  isAllSections,
	}
	data.Project = &ProjectStub{ID: project.ID, Name: project.Name}

	if err := h.renderer.Render(w, "materials_list_show.html", data); err != nil {
		h.logger.Error("rendering materials list show", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ── Update notes ───────────────────────────────────────────────

// UpdateMaterialsListNotes updates notes on a materials list.
func (h *Handler) UpdateMaterialsListNotes(w http.ResponseWriter, r *http.Request) {
	listID := r.PathValue("lid")
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err := h.queries.UpdateMaterialsListNotes(ctx, repository.UpdateMaterialsListNotesParams{
		ID:    listID,
		Notes: toNullString(r.FormValue("notes")),
	})
	if err != nil {
		h.logger.Error("updating materials list notes", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Saved")
}

// ── Delete ─────────────────────────────────────────────────────

// DeleteMaterialsList deletes a materials list.
func (h *Handler) DeleteMaterialsList(w http.ResponseWriter, r *http.Request) {
	quoteID := r.PathValue("id")
	listID := r.PathValue("lid")
	ctx := r.Context()

	list, err := h.queries.GetMaterialsList(ctx, listID)
	if err != nil {
		h.logger.Error("getting materials list", "error", err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if list.QuoteID != quoteID {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if err := h.queries.DeleteMaterialsList(ctx, listID); err != nil {
		h.logger.Error("deleting materials list", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", fmt.Sprintf("/quotes/%s/materials-lists", quoteID))
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/quotes/%s/materials-lists", quoteID), http.StatusSeeOther)
}
