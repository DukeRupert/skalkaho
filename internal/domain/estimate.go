package domain

// EstimatePayload is the JSON contract between the Go API and the Svelte estimate builder.
type EstimatePayload struct {
	Project          EstimateProject          `json:"project"`
	Globals          MarkupGlobals            `json:"globals"`
	Sections         []EstimateSection        `json:"sections"`
	MaterialsDB      []MaterialDBEntry        `json:"materials_db"`
	RatesDB          []RateDBEntry            `json:"rates_db"`
	SubcontractorsDB []SubcontractorDBEntry   `json:"subcontractors_db"`
}

// EstimateProject is the minimal project info included in the estimate payload.
type EstimateProject struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// MarkupGlobals holds the 5 global markup percentages for a project.
type MarkupGlobals struct {
	MaterialsMarkup float64 `json:"materials_markup"`
	LaborMarkup     float64 `json:"labor_markup"`
	EquipmentMarkup float64 `json:"equipment_markup"`
	SubsMarkup      float64 `json:"subs_markup"`
	OtherMarkup     float64 `json:"other_markup"`
}

// EstimateSection is the top level of the 4-level estimate hierarchy.
type EstimateSection struct {
	ID            string                `json:"id"`
	Name          string                `json:"name"`
	SortOrder     int                   `json:"sort_order"`
	Subcategories []EstimateSubcategory `json:"subcategories"`
}

// MarkupOverrides holds per-subcategory markup overrides. Nil means inherit global.
type MarkupOverrides struct {
	Materials *float64 `json:"materials"`
	Labor     *float64 `json:"labor"`
	Equipment *float64 `json:"equipment"`
	Subs      *float64 `json:"subs"`
	Other     *float64 `json:"other"`
}

// MarkupEnabled holds per-subcategory enable/disable toggles for each markup type.
type MarkupEnabled struct {
	Materials bool `json:"materials"`
	Labor     bool `json:"labor"`
	Equipment bool `json:"equipment"`
	Subs      bool `json:"subs"`
	Other     bool `json:"other"`
}

// EstimateSubcategory is the second level, with markup controls and lump sum.
type EstimateSubcategory struct {
	ID              string                   `json:"id"`
	Name            string                   `json:"name"`
	SortOrder       int                      `json:"sort_order"`
	LumpSum         float64                  `json:"lump_sum"`
	MarkupOverrides MarkupOverrides          `json:"markup_overrides"`
	MarkupEnabled   MarkupEnabled            `json:"markup_enabled"`
	ComponentGroups []EstimateComponentGroup `json:"component_groups"`
	LineItems       []EstimateLineItem       `json:"line_items"`
}

// EstimateComponentGroup is the third level, an optional grouping.
type EstimateComponentGroup struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	SortOrder int                `json:"sort_order"`
	LineItems []EstimateLineItem `json:"line_items"`
}

// CategoryType enumerates the 5 markup categories for line items.
type CategoryType string

const (
	CategoryTypeMaterials CategoryType = "materials"
	CategoryTypeLabor     CategoryType = "labor"
	CategoryTypeEquipment CategoryType = "equipment"
	CategoryTypeSubs      CategoryType = "subs"
	CategoryTypeOther     CategoryType = "other"
)

// EstimateLineItem is the fourth level, an individual cost entry.
type EstimateLineItem struct {
	ID               string       `json:"id"`
	CategoryType     CategoryType `json:"category_type"`
	ItemName         string       `json:"item_name"`
	Quantity         float64      `json:"quantity"`
	Unit             string       `json:"unit"`
	UnitPrice        float64      `json:"unit_price"`
	IsCustom         bool         `json:"is_custom"`
	MaterialID       *string      `json:"material_id,omitempty"`
	SubcontractorID  *string      `json:"subcontractor_id,omitempty"`
	PriceOverride    bool         `json:"price_override"`
	Description      *string      `json:"description,omitempty"`
	SortOrder        int          `json:"sort_order"`
	ComponentGroupID *string      `json:"component_group_id,omitempty"`
	VisualGroup      *string      `json:"visual_group,omitempty"`
}

// MaterialDBEntry represents a material from the database for autocomplete.
type MaterialDBEntry struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Supplier     string  `json:"supplier"`
	UnitPrice    float64 `json:"unit_price"`
	Unit         string  `json:"unit"`
	SupplierCode string  `json:"supplier_code,omitempty"`
}

// SubcontractorDBEntry represents a subcontractor from the directory for the picker.
type SubcontractorDBEntry struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Company      string `json:"company"`
	PrimaryTrade string `json:"primary_trade"`
}

// RateDBEntry represents a labor/equipment/subs rate for autocomplete.
type RateDBEntry struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Rate     float64 `json:"rate"`
	Unit     string  `json:"unit"`
}

// ResolveMarkup resolves the effective markup percentage for a line item
// within a subcategory, following the spec's resolution order:
// 1. Check if the relevant type's markup_enabled toggle is false -> use 0%
// 2. Check if the subcategory has a non-null override for that type -> use it
// 3. Fall back to the global default for that type
func ResolveMarkup(categoryType CategoryType, globals MarkupGlobals, overrides MarkupOverrides, enabled MarkupEnabled) float64 {
	switch categoryType {
	case CategoryTypeMaterials:
		if !enabled.Materials {
			return 0
		}
		if overrides.Materials != nil {
			return *overrides.Materials
		}
		return globals.MaterialsMarkup
	case CategoryTypeLabor:
		if !enabled.Labor {
			return 0
		}
		if overrides.Labor != nil {
			return *overrides.Labor
		}
		return globals.LaborMarkup
	case CategoryTypeEquipment:
		if !enabled.Equipment {
			return 0
		}
		if overrides.Equipment != nil {
			return *overrides.Equipment
		}
		return globals.EquipmentMarkup
	case CategoryTypeSubs:
		if !enabled.Subs {
			return 0
		}
		if overrides.Subs != nil {
			return *overrides.Subs
		}
		return globals.SubsMarkup
	case CategoryTypeOther:
		if !enabled.Other {
			return 0
		}
		if overrides.Other != nil {
			return *overrides.Other
		}
		return globals.OtherMarkup
	default:
		return 0
	}
}
