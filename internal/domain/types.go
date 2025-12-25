package domain

import "time"

// SurchargeMode defines how surcharges are calculated.
type SurchargeMode string

const (
	SurchargeModeStacking SurchargeMode = "stacking"
	SurchargeModeOverride SurchargeMode = "override"
)

// LineItemType distinguishes materials, labor, and equipment.
type LineItemType string

const (
	LineItemTypeMaterial  LineItemType = "material"
	LineItemTypeLabor     LineItemType = "labor"
	LineItemTypeEquipment LineItemType = "equipment"
)

// Settings holds application-wide defaults.
type Settings struct {
	ID                      string        `json:"id"`
	DefaultSurchargeMode    SurchargeMode `json:"default_surcharge_mode"`
	DefaultSurchargePercent float64       `json:"default_surcharge_percent"`
}

// Job is the top-level container for a quote.
type Job struct {
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	CustomerName     *string       `json:"customer_name,omitempty"`
	SurchargePercent float64       `json:"surcharge_percent"`
	SurchargeMode    SurchargeMode `json:"surcharge_mode"`
	CreatedAt        time.Time     `json:"created_at"`
}

// Category represents an organizational grouping within a job.
type Category struct {
	ID               string   `json:"id"`
	JobID            string   `json:"job_id"`
	ParentID         *string  `json:"parent_id,omitempty"`
	Name             string   `json:"name"`
	SurchargePercent *float64 `json:"surcharge_percent,omitempty"`
	SortOrder        int      `json:"sort_order"`
}

// LineItem represents an individual material or labor entry.
type LineItem struct {
	ID               string       `json:"id"`
	CategoryID       string       `json:"category_id"`
	Type             LineItemType `json:"type"`
	Name             string       `json:"name"`
	Description      *string      `json:"description,omitempty"`
	Quantity         float64      `json:"quantity"`
	Unit             string       `json:"unit"`
	UnitPrice        float64      `json:"unit_price"`
	SurchargePercent *float64     `json:"surcharge_percent,omitempty"`
	SortOrder        int          `json:"sort_order"`
}

// BasePrice calculates quantity * unit_price.
func (li *LineItem) BasePrice() float64 {
	return li.Quantity * li.UnitPrice
}

// CommonUnits returns suggested units for the UI.
var CommonUnits = struct {
	Material []string
	Labor    []string
}{
	Material: []string{"ea", "sqft", "lnft", "bundle", "box", "bag", "gal", "sheet"},
	Labor:    []string{"hr", "day", "job", "sqft"},
}
