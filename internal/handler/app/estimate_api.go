package app

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/dukerupert/skalkaho/internal/domain"
	"github.com/dukerupert/skalkaho/internal/repository"
)

// GetEstimate handles GET /api/estimate/{id} — returns the full estimate payload as JSON.
func (h *Handler) GetEstimate(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	ctx := r.Context()

	project, err := h.queries.GetProject(ctx, projectID)
	if err != nil {
		h.logger.Error("getting project", "error", err)
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	payload, err := h.buildEstimatePayload(r, project)
	if err != nil {
		h.logger.Error("building estimate payload", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		h.logger.Error("encoding estimate", "error", err)
	}
}

// SaveEstimate handles POST /api/estimate/{id} — persists the estimate and returns the saved state.
func (h *Handler) SaveEstimate(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	ctx := r.Context()

	var payload domain.EstimatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.logger.Error("decoding estimate", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Verify project exists
	project, err := h.queries.GetProject(ctx, projectID)
	if err != nil {
		h.logger.Error("getting project", "error", err)
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Update global markups
	if err := h.queries.UpdateProjectMarkups(ctx, repository.UpdateProjectMarkupsParams{
		ID:              projectID,
		MaterialsMarkup: payload.Globals.MaterialsMarkup,
		LaborMarkup:     payload.Globals.LaborMarkup,
		EquipmentMarkup: payload.Globals.EquipmentMarkup,
		SubsMarkup:      payload.Globals.SubsMarkup,
		OtherMarkup:     payload.Globals.OtherMarkup,
	}); err != nil {
		h.logger.Error("updating markups", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Delete existing estimate data (cascade from sections)
	if err := h.queries.DeleteSectionsByProject(ctx, projectID); err != nil {
		h.logger.Error("deleting sections", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Re-create all sections, subcategories, component groups, and line items
	for _, section := range payload.Sections {
		if _, err := h.queries.CreateSection(ctx, repository.CreateSectionParams{
			ID:        section.ID,
			ProjectID: projectID,
			Name:      section.Name,
			SortOrder: int64(section.SortOrder),
		}); err != nil {
			h.logger.Error("creating section", "error", err, "section", section.ID)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		for _, subcat := range section.Subcategories {
			if _, err := h.queries.CreateSubcategory(ctx, repository.CreateSubcategoryParams{
				ID:                     subcat.ID,
				SectionID:              section.ID,
				Name:                   subcat.Name,
				SortOrder:              int64(subcat.SortOrder),
				LumpSum:                subcat.LumpSum,
				MaterialsMarkup:        toNullFloat64(subcat.MarkupOverrides.Materials),
				LaborMarkup:            toNullFloat64(subcat.MarkupOverrides.Labor),
				EquipmentMarkup:        toNullFloat64(subcat.MarkupOverrides.Equipment),
				SubsMarkup:             toNullFloat64(subcat.MarkupOverrides.Subs),
				OtherMarkup:            toNullFloat64(subcat.MarkupOverrides.Other),
				MaterialsMarkupEnabled: subcat.MarkupEnabled.Materials,
				LaborMarkupEnabled:     subcat.MarkupEnabled.Labor,
				EquipmentMarkupEnabled: subcat.MarkupEnabled.Equipment,
				SubsMarkupEnabled:      subcat.MarkupEnabled.Subs,
				OtherMarkupEnabled:     subcat.MarkupEnabled.Other,
			}); err != nil {
				h.logger.Error("creating subcategory", "error", err, "subcategory", subcat.ID)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Create component groups
			for _, cg := range subcat.ComponentGroups {
				if _, err := h.queries.CreateComponentGroup(ctx, repository.CreateComponentGroupParams{
					ID:            cg.ID,
					SubcategoryID: subcat.ID,
					Name:          cg.Name,
					SortOrder:     int64(cg.SortOrder),
				}); err != nil {
					h.logger.Error("creating component group", "error", err, "group", cg.ID)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				// Create line items in this group
				for _, li := range cg.LineItems {
					if err := h.createLineItem(r, subcat.ID, &cg.ID, li); err != nil {
						h.logger.Error("creating line item", "error", err, "item", li.ID)
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						return
					}
				}
			}

			// Create ungrouped line items
			for _, li := range subcat.LineItems {
				if err := h.createLineItem(r, subcat.ID, nil, li); err != nil {
					h.logger.Error("creating line item", "error", err, "item", li.ID)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	// Calculate and update project total
	total := domain.CalculateProjectTotal(payload.Sections, payload.Globals)
	if err := h.queries.UpdateProjectTotal(ctx, repository.UpdateProjectTotalParams{
		ID:    projectID,
		Total: total,
	}); err != nil {
		h.logger.Error("updating project total", "error", err)
	}

	// Return the saved state
	project.Total = total
	project.MaterialsMarkup = payload.Globals.MaterialsMarkup
	project.LaborMarkup = payload.Globals.LaborMarkup
	project.EquipmentMarkup = payload.Globals.EquipmentMarkup
	project.SubsMarkup = payload.Globals.SubsMarkup
	project.OtherMarkup = payload.Globals.OtherMarkup

	saved, err := h.buildEstimatePayload(r, project)
	if err != nil {
		h.logger.Error("building saved payload", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(saved); err != nil {
		h.logger.Error("encoding saved estimate", "error", err)
	}
}

func (h *Handler) createLineItem(r *http.Request, subcategoryID string, componentGroupID *string, li domain.EstimateLineItem) error {
	cgID := sql.NullString{}
	if componentGroupID != nil {
		cgID = sql.NullString{String: *componentGroupID, Valid: true}
	} else if li.ComponentGroupID != nil {
		cgID = sql.NullString{String: *li.ComponentGroupID, Valid: true}
	}

	_, err := h.queries.CreateLineItem(r.Context(), repository.CreateLineItemParams{
		ID:               li.ID,
		SubcategoryID:    subcategoryID,
		ComponentGroupID: cgID,
		CategoryType:     string(li.CategoryType),
		ItemName:         li.ItemName,
		Quantity:         li.Quantity,
		Unit:             li.Unit,
		UnitPrice:        li.UnitPrice,
		IsCustom:         li.IsCustom,
		MaterialID:       toNullString(ptrToStr(li.MaterialID)),
		SubcontractorID:  toNullString(ptrToStr(li.SubcontractorID)),
		PriceOverride:    li.PriceOverride,
		Description:      toNullString(ptrToStr(li.Description)),
		SortOrder:        int64(li.SortOrder),
		VisualGroup:      toNullString(ptrToStr(li.VisualGroup)),
	})
	return err
}

func (h *Handler) buildEstimatePayload(r *http.Request, project repository.Project) (domain.EstimatePayload, error) {
	ctx := r.Context()

	// Load sections
	dbSections, err := h.queries.ListSectionsByProject(ctx, project.ID)
	if err != nil {
		return domain.EstimatePayload{}, err
	}

	sections := make([]domain.EstimateSection, 0, len(dbSections))
	for _, s := range dbSections {
		subcats, err := h.buildSubcategories(r, s.ID)
		if err != nil {
			return domain.EstimatePayload{}, err
		}
		sections = append(sections, domain.EstimateSection{
			ID:            s.ID,
			Name:          s.Name,
			SortOrder:     int(s.SortOrder),
			Subcategories: subcats,
		})
	}

	// Load materials for autocomplete
	dbMaterials, err := h.queries.ListMaterials(ctx)
	if err != nil {
		return domain.EstimatePayload{}, err
	}

	// Build supplier name map for materials
	suppliers, err := h.queries.ListSuppliers(ctx)
	if err != nil {
		return domain.EstimatePayload{}, err
	}
	supplierNames := make(map[string]string, len(suppliers))
	for _, s := range suppliers {
		supplierNames[s.ID] = s.Name
	}

	materialsDB := make([]domain.MaterialDBEntry, 0, len(dbMaterials))
	for _, m := range dbMaterials {
		materialsDB = append(materialsDB, domain.MaterialDBEntry{
			ID:           m.ID,
			Name:         m.Name,
			Supplier:     supplierNames[m.SupplierID],
			UnitPrice:    m.UnitPrice,
			Unit:         m.Unit,
			SupplierCode: nullStrToStr(m.SupplierCode),
		})
	}

	// Load rates for autocomplete
	dbRates, err := h.queries.ListRates(ctx)
	if err != nil {
		return domain.EstimatePayload{}, err
	}

	// Build category name map for rates
	categories, err := h.queries.ListRateCategories(ctx)
	if err != nil {
		return domain.EstimatePayload{}, err
	}
	categoryNames := make(map[string]string, len(categories))
	for _, c := range categories {
		categoryNames[c.ID] = c.Name
	}

	ratesDB := make([]domain.RateDBEntry, 0, len(dbRates))
	for _, rate := range dbRates {
		ratesDB = append(ratesDB, domain.RateDBEntry{
			ID:       rate.ID,
			Name:     rate.Name,
			Category: categoryNames[rate.CategoryID],
			Rate:     rate.Rate,
			Unit:     rate.Unit,
		})
	}

	// Load subcontractors for picker
	dbSubs, err := h.queries.ListSubcontractors(ctx)
	if err != nil {
		return domain.EstimatePayload{}, err
	}
	subcontractorsDB := make([]domain.SubcontractorDBEntry, 0, len(dbSubs))
	for _, s := range dbSubs {
		subcontractorsDB = append(subcontractorsDB, domain.SubcontractorDBEntry{
			ID:           s.ID,
			Name:         s.Name,
			Company:      nullStrToStr(s.Company),
			PrimaryTrade: nullStrToStr(s.PrimaryTrade),
		})
	}

	return domain.EstimatePayload{
		Project: domain.EstimateProject{
			ID:     project.ID,
			Name:   project.Name,
			Status: project.Status,
		},
		Globals: domain.MarkupGlobals{
			MaterialsMarkup: project.MaterialsMarkup,
			LaborMarkup:     project.LaborMarkup,
			EquipmentMarkup: project.EquipmentMarkup,
			SubsMarkup:      project.SubsMarkup,
			OtherMarkup:     project.OtherMarkup,
		},
		Sections:         sections,
		MaterialsDB:      materialsDB,
		RatesDB:          ratesDB,
		SubcontractorsDB: subcontractorsDB,
	}, nil
}

func (h *Handler) buildSubcategories(r *http.Request, sectionID string) ([]domain.EstimateSubcategory, error) {
	ctx := r.Context()
	dbSubcats, err := h.queries.ListSubcategoriesBySection(ctx, sectionID)
	if err != nil {
		return nil, err
	}

	subcats := make([]domain.EstimateSubcategory, 0, len(dbSubcats))
	for _, sc := range dbSubcats {
		// Load component groups
		dbGroups, err := h.queries.ListComponentGroupsBySubcategory(ctx, sc.ID)
		if err != nil {
			return nil, err
		}

		groups := make([]domain.EstimateComponentGroup, 0, len(dbGroups))
		for _, g := range dbGroups {
			groupItems, err := h.loadLineItemsForGroup(r, sc.ID, g.ID)
			if err != nil {
				return nil, err
			}
			groups = append(groups, domain.EstimateComponentGroup{
				ID:        g.ID,
				Name:      g.Name,
				SortOrder: int(g.SortOrder),
				LineItems: groupItems,
			})
		}

		// Load ungrouped line items
		allItems, err := h.queries.ListLineItemsBySubcategory(ctx, sc.ID)
		if err != nil {
			return nil, err
		}

		ungrouped := make([]domain.EstimateLineItem, 0)
		for _, li := range allItems {
			if !li.ComponentGroupID.Valid {
				ungrouped = append(ungrouped, toDomainLineItem(li))
			}
		}

		subcats = append(subcats, domain.EstimateSubcategory{
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
			ComponentGroups: groups,
			LineItems:       ungrouped,
		})
	}

	return subcats, nil
}

func (h *Handler) loadLineItemsForGroup(r *http.Request, subcategoryID, groupID string) ([]domain.EstimateLineItem, error) {
	allItems, err := h.queries.ListLineItemsBySubcategory(r.Context(), subcategoryID)
	if err != nil {
		return nil, err
	}

	var items []domain.EstimateLineItem
	for _, li := range allItems {
		if li.ComponentGroupID.Valid && li.ComponentGroupID.String == groupID {
			items = append(items, toDomainLineItem(li))
		}
	}
	return items, nil
}

func toDomainLineItem(li repository.LineItem) domain.EstimateLineItem {
	item := domain.EstimateLineItem{
		ID:            li.ID,
		CategoryType:  domain.CategoryType(li.CategoryType),
		ItemName:      li.ItemName,
		Quantity:      li.Quantity,
		Unit:          li.Unit,
		UnitPrice:     li.UnitPrice,
		IsCustom:      li.IsCustom,
		PriceOverride: li.PriceOverride,
		SortOrder:     int(li.SortOrder),
	}
	if li.MaterialID.Valid {
		item.MaterialID = &li.MaterialID.String
	}
	if li.SubcontractorID.Valid {
		item.SubcontractorID = &li.SubcontractorID.String
	}
	if li.Description.Valid {
		item.Description = &li.Description.String
	}
	if li.ComponentGroupID.Valid {
		item.ComponentGroupID = &li.ComponentGroupID.String
	}
	if li.VisualGroup.Valid {
		item.VisualGroup = &li.VisualGroup.String
	}
	return item
}

func toNullFloat64(p *float64) sql.NullFloat64 {
	if p == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *p, Valid: true}
}

func nullFloat64ToPtr(nf sql.NullFloat64) *float64 {
	if !nf.Valid {
		return nil
	}
	return &nf.Float64
}

func nullStrToStr(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func ptrToStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
