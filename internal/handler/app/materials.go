package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/dukerupert/skalkaho/internal/repository"
)

// MaterialsPageData extends PageData with materials-specific data.
type MaterialsPageData struct {
	PageData
	Materials    []repository.Material
	Suppliers    []repository.Supplier
	SupplierID   string
	PriceSource  string
	Search       string
	EditMaterial *repository.Material
}

// ListMaterials renders the materials list with supplier tabs, search, and filters.
func (h *Handler) ListMaterials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	supplierID := r.URL.Query().Get("supplier")
	priceSource := r.URL.Query().Get("source")
	search := r.URL.Query().Get("search")

	suppliers, err := h.queries.ListSuppliers(ctx)
	if err != nil {
		h.logger.Error("listing suppliers", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var materials []repository.Material

	switch {
	case search != "" && supplierID != "":
		materials, err = h.queries.SearchMaterialsBySupplier(ctx, repository.SearchMaterialsBySupplierParams{
			SupplierID: supplierID,
			SearchTerm: search,
		})
	case search != "":
		materials, err = h.queries.SearchMaterials(ctx, search)
	case supplierID != "" && priceSource != "":
		materials, err = h.queries.ListMaterialsBySupplierAndSource(ctx, repository.ListMaterialsBySupplierAndSourceParams{
			SupplierID:  supplierID,
			PriceSource: priceSource,
		})
	case supplierID != "":
		materials, err = h.queries.ListMaterialsBySupplier(ctx, supplierID)
	case priceSource != "":
		materials, err = h.queries.ListMaterialsByPriceSource(ctx, priceSource)
	default:
		materials, err = h.queries.ListMaterials(ctx)
	}
	if err != nil {
		h.logger.Error("listing materials", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := MaterialsPageData{
		PageData:    h.pageData(r, "materials"),
		Materials:   materials,
		Suppliers:   suppliers,
		SupplierID:  supplierID,
		PriceSource: priceSource,
		Search:      search,
	}

	if r.Header.Get("HX-Request") == "true" {
		if err := h.renderer.RenderPartial(w, "materials.html", "material-rows", data); err != nil {
			h.logger.Error("rendering material rows", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if err := h.renderer.Render(w, "materials.html", data); err != nil {
		h.logger.Error("rendering materials", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// CreateMaterial handles POST /materials.
func (h *Handler) CreateMaterial(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	supplierID := r.FormValue("supplier_id")
	if name == "" || supplierID == "" {
		http.Error(w, "Name and supplier are required", http.StatusBadRequest)
		return
	}

	unitPrice, _ := strconv.ParseFloat(r.FormValue("unit_price"), 64)
	priceSource := r.FormValue("price_source")
	if priceSource == "" {
		priceSource = "Manual"
	}

	id := uuid.New().String()[:20]
	_, err := h.queries.CreateMaterial(r.Context(), repository.CreateMaterialParams{
		ID:           id,
		Name:         name,
		SupplierID:   supplierID,
		UnitPrice:    unitPrice,
		Unit:         r.FormValue("unit"),
		SupplierCode: toNullString(r.FormValue("supplier_code")),
		PriceSource:  priceSource,
	})
	if err != nil {
		h.logger.Error("creating material", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/materials?supplier=%s", supplierID), http.StatusSeeOther)
}

// GetMaterialEditForm returns the edit form partial for a material.
func (h *Handler) GetMaterialEditForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	material, err := h.queries.GetMaterial(r.Context(), id)
	if err != nil {
		h.logger.Error("getting material", "error", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	suppliers, err := h.queries.ListSuppliers(r.Context())
	if err != nil {
		h.logger.Error("listing suppliers", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := MaterialsPageData{
		PageData:     h.pageData(r, "materials"),
		Suppliers:    suppliers,
		EditMaterial: &material,
	}

	if err := h.renderer.RenderPartial(w, "materials.html", "edit-modal", data); err != nil {
		h.logger.Error("rendering edit form", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// UpdateMaterial handles POST /materials/{id}.
func (h *Handler) UpdateMaterial(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	supplierID := r.FormValue("supplier_id")
	if name == "" || supplierID == "" {
		http.Error(w, "Name and supplier are required", http.StatusBadRequest)
		return
	}

	unitPrice, _ := strconv.ParseFloat(r.FormValue("unit_price"), 64)
	priceSource := r.FormValue("price_source")
	if priceSource == "" {
		priceSource = "Manual"
	}

	_, err := h.queries.UpdateMaterial(r.Context(), repository.UpdateMaterialParams{
		ID:           id,
		Name:         name,
		SupplierID:   supplierID,
		UnitPrice:    unitPrice,
		Unit:         r.FormValue("unit"),
		SupplierCode: toNullString(r.FormValue("supplier_code")),
		PriceSource:  priceSource,
	})
	if err != nil {
		h.logger.Error("updating material", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/materials?supplier=%s", supplierID), http.StatusSeeOther)
}

// DeleteMaterial handles DELETE /materials/{id}.
func (h *Handler) DeleteMaterial(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.queries.DeleteMaterial(r.Context(), id); err != nil {
		h.logger.Error("deleting material", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/materials")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/materials", http.StatusSeeOther)
}

// CreateSupplier handles POST /suppliers.
func (h *Handler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Supplier name is required", http.StatusBadRequest)
		return
	}

	// Get current max sort order
	suppliers, _ := h.queries.ListSuppliers(r.Context())
	sortOrder := int64(len(suppliers))

	id := uuid.New().String()[:20]
	_, err := h.queries.CreateSupplier(r.Context(), repository.CreateSupplierParams{
		ID:        id,
		Name:      name,
		SortOrder: sortOrder,
	})
	if err != nil {
		h.logger.Error("creating supplier", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/materials?supplier=%s", id), http.StatusSeeOther)
}

// DeleteSupplier handles DELETE /suppliers/{id}.
func (h *Handler) DeleteSupplier(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.queries.DeleteSupplier(r.Context(), id); err != nil {
		h.logger.Error("deleting supplier", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/materials")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/materials", http.StatusSeeOther)
}
