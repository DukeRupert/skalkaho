package domain

// CalculateProjectTotal computes the grand total for a project's estimate.
// It sums all line item costs (quantity * unit_price) across all sections
// and subcategories, applying the resolved markup for each line item.
// This is the single source of truth for cost calculation, used by both
// the estimate save handler and the project overview.
func CalculateProjectTotal(sections []EstimateSection, globals MarkupGlobals) float64 {
	summary := CalculateProjectCosts(sections, globals)
	return summary.GrandTotal
}

// ProjectCostSummary holds the calculated cost breakdown for a project.
type ProjectCostSummary struct {
	// Per-category base costs (before markup)
	MaterialsBase float64 `json:"materials_base"`
	LaborBase     float64 `json:"labor_base"`
	EquipmentBase float64 `json:"equipment_base"`
	SubsBase      float64 `json:"subs_base"`
	OtherBase     float64 `json:"other_base"`

	// Per-category markup amounts
	MaterialsMarkup float64 `json:"materials_markup"`
	LaborMarkup     float64 `json:"labor_markup"`
	EquipmentMarkup float64 `json:"equipment_markup"`
	SubsMarkup      float64 `json:"subs_markup"`
	OtherMarkup     float64 `json:"other_markup"`

	// Lump sum total across all subcategories
	LumpSumTotal float64 `json:"lump_sum_total"`

	// Totals
	BaseTotal   float64 `json:"base_total"`
	MarkupTotal float64 `json:"markup_total"`
	GrandTotal  float64 `json:"grand_total"`

	// Per-section breakdown
	Sections []SectionCostSummary `json:"sections"`
}

// SectionCostSummary holds cost breakdown for a single section.
type SectionCostSummary struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Total float64 `json:"total"`
}

// CalculateProjectCosts computes the full cost breakdown for a project.
// It resolves markup per line item using the subcategory overrides and globals.
func CalculateProjectCosts(sections []EstimateSection, globals MarkupGlobals) ProjectCostSummary {
	var summary ProjectCostSummary
	summary.Sections = make([]SectionCostSummary, 0, len(sections))

	for _, section := range sections {
		sectionTotal := 0.0

		for _, subcat := range section.Subcategories {
			summary.LumpSumTotal += subcat.LumpSum
			sectionTotal += subcat.LumpSum

			// Process all line items (ungrouped + grouped)
			allItems := subcat.LineItems
			for _, cg := range subcat.ComponentGroups {
				allItems = append(allItems, cg.LineItems...)
			}

			for _, li := range allItems {
				baseCost := li.Quantity * li.UnitPrice
				markupPct := ResolveMarkup(li.CategoryType, globals, subcat.MarkupOverrides, subcat.MarkupEnabled)
				markupAmt := baseCost * (markupPct / 100)
				totalCost := baseCost + markupAmt

				switch li.CategoryType {
				case CategoryTypeMaterials:
					summary.MaterialsBase += baseCost
					summary.MaterialsMarkup += markupAmt
				case CategoryTypeLabor:
					summary.LaborBase += baseCost
					summary.LaborMarkup += markupAmt
				case CategoryTypeEquipment:
					summary.EquipmentBase += baseCost
					summary.EquipmentMarkup += markupAmt
				case CategoryTypeSubs:
					summary.SubsBase += baseCost
					summary.SubsMarkup += markupAmt
				case CategoryTypeOther:
					summary.OtherBase += baseCost
					summary.OtherMarkup += markupAmt
				}

				sectionTotal += totalCost
			}
		}

		summary.Sections = append(summary.Sections, SectionCostSummary{
			ID:    section.ID,
			Name:  section.Name,
			Total: sectionTotal,
		})
	}

	summary.BaseTotal = summary.MaterialsBase + summary.LaborBase + summary.EquipmentBase + summary.SubsBase + summary.OtherBase
	summary.MarkupTotal = summary.MaterialsMarkup + summary.LaborMarkup + summary.EquipmentMarkup + summary.SubsMarkup + summary.OtherMarkup
	summary.GrandTotal = summary.BaseTotal + summary.MarkupTotal + summary.LumpSumTotal

	return summary
}
