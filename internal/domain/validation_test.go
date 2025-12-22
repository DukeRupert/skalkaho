package domain_test

import (
	"strings"
	"testing"

	"github.com/dukerupert/skalkaho/internal/domain"
)

func TestJobInput_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     domain.JobInput
		wantErr   bool
		wantField string
	}{
		{
			name: "valid input",
			input: domain.JobInput{
				Name:          "Smith Kitchen Remodel",
				SurchargeMode: domain.SurchargeModeStacking,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			input: domain.JobInput{
				Name: "",
			},
			wantErr:   true,
			wantField: "name",
		},
		{
			name: "whitespace name",
			input: domain.JobInput{
				Name: "   ",
			},
			wantErr:   true,
			wantField: "name",
		},
		{
			name: "name too long",
			input: domain.JobInput{
				Name: strings.Repeat("a", 256),
			},
			wantErr:   true,
			wantField: "name",
		},
		{
			name: "invalid surcharge mode",
			input: domain.JobInput{
				Name:          "Valid Name",
				SurchargeMode: "invalid",
			},
			wantErr:   true,
			wantField: "surcharge_mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := tt.input.Validate()

			if tt.wantErr && len(errors) == 0 {
				t.Error("expected validation error, got none")
			}

			if !tt.wantErr && len(errors) > 0 {
				t.Errorf("expected no errors, got %v", errors)
			}

			if tt.wantField != "" {
				found := false
				for _, err := range errors {
					if err.Field == tt.wantField {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error on field %q, got %v", tt.wantField, errors)
				}
			}
		})
	}
}

func TestCategoryInput_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     domain.CategoryInput
		wantErr   bool
		wantField string
	}{
		{
			name: "valid input",
			input: domain.CategoryInput{
				Name: "Electrical",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			input: domain.CategoryInput{
				Name: "",
			},
			wantErr:   true,
			wantField: "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := tt.input.Validate()

			if tt.wantErr && len(errors) == 0 {
				t.Error("expected validation error, got none")
			}

			if !tt.wantErr && len(errors) > 0 {
				t.Errorf("expected no errors, got %v", errors)
			}
		})
	}
}

func TestValidateCategoryDepth(t *testing.T) {
	tests := []struct {
		name        string
		parentDepth int
		wantErr     bool
	}{
		{
			name:        "depth 0 allowed",
			parentDepth: 0,
			wantErr:     false,
		},
		{
			name:        "depth 2 allowed",
			parentDepth: 2,
			wantErr:     false,
		},
		{
			name:        "depth 3 not allowed",
			parentDepth: 3,
			wantErr:     true,
		},
		{
			name:        "depth 4 not allowed",
			parentDepth: 4,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidateCategoryDepth(tt.parentDepth)

			if tt.wantErr && err == nil {
				t.Error("expected validation error, got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestLineItemInput_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     domain.LineItemInput
		wantErr   bool
		wantField string
	}{
		{
			name: "valid material",
			input: domain.LineItemInput{
				Type:      domain.LineItemTypeMaterial,
				Name:      "2x4 Lumber",
				Quantity:  10,
				Unit:      "ea",
				UnitPrice: 5.99,
			},
			wantErr: false,
		},
		{
			name: "valid labor",
			input: domain.LineItemInput{
				Type:      domain.LineItemTypeLabor,
				Name:      "Electrician",
				Quantity:  8,
				Unit:      "hr",
				UnitPrice: 75,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			input: domain.LineItemInput{
				Type:      domain.LineItemTypeMaterial,
				Name:      "",
				Quantity:  1,
				Unit:      "ea",
				UnitPrice: 10,
			},
			wantErr:   true,
			wantField: "name",
		},
		{
			name: "invalid type",
			input: domain.LineItemInput{
				Type:      "invalid",
				Name:      "Item",
				Quantity:  1,
				Unit:      "ea",
				UnitPrice: 10,
			},
			wantErr:   true,
			wantField: "type",
		},
		{
			name: "zero quantity",
			input: domain.LineItemInput{
				Type:      domain.LineItemTypeMaterial,
				Name:      "Item",
				Quantity:  0,
				Unit:      "ea",
				UnitPrice: 10,
			},
			wantErr:   true,
			wantField: "quantity",
		},
		{
			name: "negative quantity",
			input: domain.LineItemInput{
				Type:      domain.LineItemTypeMaterial,
				Name:      "Item",
				Quantity:  -1,
				Unit:      "ea",
				UnitPrice: 10,
			},
			wantErr:   true,
			wantField: "quantity",
		},
		{
			name: "empty unit",
			input: domain.LineItemInput{
				Type:      domain.LineItemTypeMaterial,
				Name:      "Item",
				Quantity:  1,
				Unit:      "",
				UnitPrice: 10,
			},
			wantErr:   true,
			wantField: "unit",
		},
		{
			name: "negative price",
			input: domain.LineItemInput{
				Type:      domain.LineItemTypeMaterial,
				Name:      "Item",
				Quantity:  1,
				Unit:      "ea",
				UnitPrice: -10,
			},
			wantErr:   true,
			wantField: "unit_price",
		},
		{
			name: "zero price allowed",
			input: domain.LineItemInput{
				Type:      domain.LineItemTypeMaterial,
				Name:      "Free Item",
				Quantity:  1,
				Unit:      "ea",
				UnitPrice: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := tt.input.Validate()

			if tt.wantErr && len(errors) == 0 {
				t.Error("expected validation error, got none")
			}

			if !tt.wantErr && len(errors) > 0 {
				t.Errorf("expected no errors, got %v", errors)
			}

			if tt.wantField != "" {
				found := false
				for _, err := range errors {
					if err.Field == tt.wantField {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error on field %q, got %v", tt.wantField, errors)
				}
			}
		})
	}
}
