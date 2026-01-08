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

// StandardTypes returns the list of built-in line item types.
var StandardTypes = []LineItemType{
	LineItemTypeMaterial,
	LineItemTypeLabor,
	LineItemTypeEquipment,
}

// IsStandardType returns true if the type is one of the built-in types.
func IsStandardType(t LineItemType) bool {
	return t == LineItemTypeMaterial || t == LineItemTypeLabor || t == LineItemTypeEquipment
}

// JobItemType represents a custom line item type for a specific job.
type JobItemType struct {
	ID               string    `json:"id"`
	JobID            string    `json:"job_id"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	Color            string    `json:"color"`
	SortOrder        int       `json:"sort_order"`
	SurchargePercent *float64  `json:"surcharge_percent,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

// Settings holds application-wide defaults.
type Settings struct {
	ID                      string        `json:"id"`
	DefaultSurchargeMode    SurchargeMode `json:"default_surcharge_mode"`
	DefaultSurchargePercent float64       `json:"default_surcharge_percent"`
}

// Job is the top-level container for a quote.
type Job struct {
	ID                        string        `json:"id"`
	Name                      string        `json:"name"`
	CustomerName              *string       `json:"customer_name,omitempty"`
	SurchargePercent          float64       `json:"surcharge_percent"`
	MaterialSurchargePercent  *float64      `json:"material_surcharge_percent,omitempty"`
	LaborSurchargePercent     *float64      `json:"labor_surcharge_percent,omitempty"`
	EquipmentSurchargePercent *float64      `json:"equipment_surcharge_percent,omitempty"`
	SurchargeMode             SurchargeMode `json:"surcharge_mode"`
	CreatedAt                 time.Time     `json:"created_at"`
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
	Tag              *string      `json:"tag,omitempty"`
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

// EstimateStatus defines the lifecycle states for an estimate.
type EstimateStatus string

const (
	EstimateStatusDraft    EstimateStatus = "draft"
	EstimateStatusSent     EstimateStatus = "sent"
	EstimateStatusAccepted EstimateStatus = "accepted"
	EstimateStatusRejected EstimateStatus = "rejected"
)

// Estimate represents a client-facing snapshot of a quote.
type Estimate struct {
	ID          string         `json:"id"`
	JobID       string         `json:"job_id"`
	Version     int            `json:"version"`
	Status      EstimateStatus `json:"status"`
	GrandTotal  float64        `json:"grand_total"`
	Notes       *string        `json:"notes,omitempty"`
	SentAt      *time.Time     `json:"sent_at,omitempty"`
	RespondedAt *time.Time     `json:"responded_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}

// EstimateCategory represents a category snapshot in an estimate.
type EstimateCategory struct {
	ID               string  `json:"id"`
	EstimateID       string  `json:"estimate_id"`
	CategoryID       string  `json:"category_id"`
	ParentCategoryID *string `json:"parent_category_id,omitempty"`
	Tier             int     `json:"tier"`
	Name             string  `json:"name"`
	Description      *string `json:"description,omitempty"`
	Total            float64 `json:"total"`
	SortOrder        int     `json:"sort_order"`
}
