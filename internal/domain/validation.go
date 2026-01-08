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
	Tag              *string      `json:"tag"`
	ValidCustomTypes []string     `json:"-"` // Set by handler before validation
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

	// Type validation - accept standard types OR valid custom types
	isStandard := IsStandardType(i.Type)
	isCustomValid := false
	for _, ct := range i.ValidCustomTypes {
		if string(i.Type) == ct {
			isCustomValid = true
			break
		}
	}

	if !isStandard && !isCustomValid {
		errors = append(errors, ValidationError{
			Field:   "type",
			Message: "Type must be a valid standard or custom type",
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

// JobItemTypeInput represents input for creating or updating a custom item type.
type JobItemTypeInput struct {
	JobID            string   `json:"job_id"`
	Name             string   `json:"name"`
	Slug             string   `json:"slug"`
	Color            string   `json:"color"`
	SortOrder        int      `json:"sort_order"`
	SurchargePercent *float64 `json:"surcharge_percent,omitempty"`
}

// ValidColors are the allowed Tailwind color prefixes for custom types.
var ValidColors = []string{"forest", "copper", "slate", "amber", "rose", "violet", "cyan", "lime"}

// Validate checks the job item type input for errors.
func (i *JobItemTypeInput) Validate() []ValidationError {
	var errors []ValidationError

	if strings.TrimSpace(i.Name) == "" {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "Name is required",
		})
	} else if len(i.Name) > 100 {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "Name must be less than 100 characters",
		})
	}

	if strings.TrimSpace(i.Slug) == "" {
		errors = append(errors, ValidationError{
			Field:   "slug",
			Message: "Slug is required",
		})
	} else if !isValidSlug(i.Slug) {
		errors = append(errors, ValidationError{
			Field:   "slug",
			Message: "Slug must be lowercase alphanumeric with hyphens only",
		})
	}

	// Prevent using reserved standard type slugs
	reserved := []string{"material", "labor", "equipment"}
	for _, r := range reserved {
		if i.Slug == r {
			errors = append(errors, ValidationError{
				Field:   "slug",
				Message: "Cannot use reserved type name",
			})
			break
		}
	}

	// Color is optional - if provided, validate it
	if i.Color != "" {
		colorValid := false
		for _, c := range ValidColors {
			if i.Color == c {
				colorValid = true
				break
			}
		}
		if !colorValid {
			errors = append(errors, ValidationError{
				Field:   "color",
				Message: "Invalid color",
			})
		}
	}

	return errors
}

// isValidSlug checks if a string is a valid slug (lowercase alphanumeric with hyphens).
func isValidSlug(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		isLower := c >= 'a' && c <= 'z'
		isDigit := c >= '0' && c <= '9'
		if !isLower && !isDigit && c != '-' {
			return false
		}
	}
	return true
}

// EstimateCategoryInput represents input for updating an estimate category description.
type EstimateCategoryInput struct {
	Description *string `json:"description"`
}

// Validate checks the estimate category input for errors.
func (i *EstimateCategoryInput) Validate() []ValidationError {
	var errors []ValidationError

	if i.Description != nil && len(*i.Description) > 1000 {
		errors = append(errors, ValidationError{
			Field:   "description",
			Message: "Description must be less than 1000 characters",
		})
	}

	return errors
}

// ValidEstimateStatus checks if a status string is valid.
func ValidEstimateStatus(s string) bool {
	switch EstimateStatus(s) {
	case EstimateStatusDraft, EstimateStatusSent, EstimateStatusAccepted, EstimateStatusRejected:
		return true
	}
	return false
}
