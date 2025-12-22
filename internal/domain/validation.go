package domain

import (
	"strings"
)

// ValidationError represents a single field validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// JobInput represents input for creating or updating a job.
type JobInput struct {
	Name             string        `json:"name"`
	CustomerName     *string       `json:"customer_name"`
	SurchargePercent float64       `json:"surcharge_percent"`
	SurchargeMode    SurchargeMode `json:"surcharge_mode"`
}

// Validate checks the job input for errors.
func (i *JobInput) Validate() []ValidationError {
	var errors []ValidationError

	if strings.TrimSpace(i.Name) == "" {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "Name is required",
		})
	} else if len(i.Name) > 255 {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "Name must be less than 255 characters",
		})
	}

	if i.SurchargeMode != "" && i.SurchargeMode != SurchargeModeStacking && i.SurchargeMode != SurchargeModeOverride {
		errors = append(errors, ValidationError{
			Field:   "surcharge_mode",
			Message: "Surcharge mode must be 'stacking' or 'override'",
		})
	}

	return errors
}

// CategoryInput represents input for creating or updating a category.
type CategoryInput struct {
	JobID            string   `json:"job_id"`
	ParentID         *string  `json:"parent_id"`
	Name             string   `json:"name"`
	SurchargePercent *float64 `json:"surcharge_percent"`
	SortOrder        int      `json:"sort_order"`
}

// Validate checks the category input for errors.
func (i *CategoryInput) Validate() []ValidationError {
	var errors []ValidationError

	if strings.TrimSpace(i.Name) == "" {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "Name is required",
		})
	} else if len(i.Name) > 255 {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "Name must be less than 255 characters",
		})
	}

	return errors
}

// ValidateCategoryDepth checks if adding a category at this level would exceed max depth.
// Returns an error if the resulting depth would be > 3.
func ValidateCategoryDepth(parentDepth int) *ValidationError {
	if parentDepth >= 3 {
		return &ValidationError{
			Field:   "parent_id",
			Message: "Maximum category nesting depth is 3 levels",
		}
	}
	return nil
}

// LineItemInput represents input for creating or updating a line item.
type LineItemInput struct {
	CategoryID       string       `json:"category_id"`
	Type             LineItemType `json:"type"`
	Name             string       `json:"name"`
	Description      *string      `json:"description"`
	Quantity         float64      `json:"quantity"`
	Unit             string       `json:"unit"`
	UnitPrice        float64      `json:"unit_price"`
	SurchargePercent *float64     `json:"surcharge_percent"`
	SortOrder        int          `json:"sort_order"`
}

// Validate checks the line item input for errors.
func (i *LineItemInput) Validate() []ValidationError {
	var errors []ValidationError

	if strings.TrimSpace(i.Name) == "" {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "Name is required",
		})
	} else if len(i.Name) > 255 {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "Name must be less than 255 characters",
		})
	}

	if i.Type != LineItemTypeMaterial && i.Type != LineItemTypeLabor {
		errors = append(errors, ValidationError{
			Field:   "type",
			Message: "Type must be 'material' or 'labor'",
		})
	}

	if i.Quantity <= 0 {
		errors = append(errors, ValidationError{
			Field:   "quantity",
			Message: "Quantity must be greater than 0",
		})
	}

	if strings.TrimSpace(i.Unit) == "" {
		errors = append(errors, ValidationError{
			Field:   "unit",
			Message: "Unit is required",
		})
	}

	if i.UnitPrice < 0 {
		errors = append(errors, ValidationError{
			Field:   "unit_price",
			Message: "Unit price cannot be negative",
		})
	}

	return errors
}

// SettingsInput represents input for updating settings.
type SettingsInput struct {
	DefaultSurchargeMode    SurchargeMode `json:"default_surcharge_mode"`
	DefaultSurchargePercent float64       `json:"default_surcharge_percent"`
}

// Validate checks the settings input for errors.
func (i *SettingsInput) Validate() []ValidationError {
	var errors []ValidationError

	if i.DefaultSurchargeMode != SurchargeModeStacking && i.DefaultSurchargeMode != SurchargeModeOverride {
		errors = append(errors, ValidationError{
			Field:   "default_surcharge_mode",
			Message: "Surcharge mode must be 'stacking' or 'override'",
		})
	}

	return errors
}
