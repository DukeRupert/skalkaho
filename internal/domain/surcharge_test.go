package domain_test

import (
	"testing"

	"github.com/dukerupert/skalkaho/internal/domain"
)

func floatPtr(f float64) *float64 {
	return &f
}

func TestEffectiveSurcharge_StackingMode(t *testing.T) {
	tests := []struct {
		name          string
		job           *domain.Job
		categoryChain []*domain.Category
		lineItem      *domain.LineItem
		want          float64
	}{
		{
			name: "all levels with surcharges",
			job: &domain.Job{
				SurchargePercent: 15,
				SurchargeMode:    domain.SurchargeModeStacking,
			},
			categoryChain: []*domain.Category{
				{SurchargePercent: floatPtr(10)},
			},
			lineItem: &domain.LineItem{
				SurchargePercent: floatPtr(5),
			},
			want: 30, // 15 + 10 + 5
		},
		{
			name: "job only",
			job: &domain.Job{
				SurchargePercent: 15,
				SurchargeMode:    domain.SurchargeModeStacking,
			},
			categoryChain: []*domain.Category{
				{SurchargePercent: nil},
			},
			lineItem: &domain.LineItem{
				SurchargePercent: nil,
			},
			want: 15,
		},
		{
			name: "nested categories",
			job: &domain.Job{
				SurchargePercent: 10,
				SurchargeMode:    domain.SurchargeModeStacking,
			},
			categoryChain: []*domain.Category{
				{SurchargePercent: floatPtr(5)},  // top level
				{SurchargePercent: floatPtr(3)},  // level 2
				{SurchargePercent: floatPtr(2)},  // level 3
			},
			lineItem: &domain.LineItem{
				SurchargePercent: nil,
			},
			want: 20, // 10 + 5 + 3 + 2
		},
		{
			name: "zero surcharges",
			job: &domain.Job{
				SurchargePercent: 0,
				SurchargeMode:    domain.SurchargeModeStacking,
			},
			categoryChain: []*domain.Category{
				{SurchargePercent: nil},
			},
			lineItem: &domain.LineItem{
				SurchargePercent: nil,
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.EffectiveSurcharge(tt.lineItem, tt.job, tt.categoryChain)
			if got != tt.want {
				t.Errorf("EffectiveSurcharge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEffectiveSurcharge_OverrideMode(t *testing.T) {
	tests := []struct {
		name          string
		job           *domain.Job
		categoryChain []*domain.Category
		lineItem      *domain.LineItem
		want          float64
	}{
		{
			name: "line item overrides all",
			job: &domain.Job{
				SurchargePercent: 15,
				SurchargeMode:    domain.SurchargeModeOverride,
			},
			categoryChain: []*domain.Category{
				{SurchargePercent: floatPtr(10)},
			},
			lineItem: &domain.LineItem{
				SurchargePercent: floatPtr(5),
			},
			want: 5, // Line item wins
		},
		{
			name: "deepest category overrides",
			job: &domain.Job{
				SurchargePercent: 15,
				SurchargeMode:    domain.SurchargeModeOverride,
			},
			categoryChain: []*domain.Category{
				{SurchargePercent: floatPtr(10)}, // top level
				{SurchargePercent: floatPtr(8)},  // level 2 - deepest with value
			},
			lineItem: &domain.LineItem{
				SurchargePercent: nil,
			},
			want: 8, // Deepest category wins
		},
		{
			name: "falls back to job",
			job: &domain.Job{
				SurchargePercent: 15,
				SurchargeMode:    domain.SurchargeModeOverride,
			},
			categoryChain: []*domain.Category{
				{SurchargePercent: nil},
			},
			lineItem: &domain.LineItem{
				SurchargePercent: nil,
			},
			want: 15, // Job wins when nothing else set
		},
		{
			name: "skips nil categories to find value",
			job: &domain.Job{
				SurchargePercent: 15,
				SurchargeMode:    domain.SurchargeModeOverride,
			},
			categoryChain: []*domain.Category{
				{SurchargePercent: floatPtr(10)}, // top level has value
				{SurchargePercent: nil},          // level 2 - nil
				{SurchargePercent: nil},          // level 3 - nil
			},
			lineItem: &domain.LineItem{
				SurchargePercent: nil,
			},
			want: 10, // Top level category wins
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.EffectiveSurcharge(tt.lineItem, tt.job, tt.categoryChain)
			if got != tt.want {
				t.Errorf("EffectiveSurcharge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFinalPrice(t *testing.T) {
	tests := []struct {
		name              string
		lineItem          *domain.LineItem
		effectiveSurcharge float64
		want              float64
	}{
		{
			name: "basic calculation",
			lineItem: &domain.LineItem{
				Quantity:  10,
				UnitPrice: 100,
			},
			effectiveSurcharge: 15,
			want:               1150, // 1000 * 1.15
		},
		{
			name: "zero surcharge",
			lineItem: &domain.LineItem{
				Quantity:  5,
				UnitPrice: 20,
			},
			effectiveSurcharge: 0,
			want:               100, // 100 * 1.00
		},
		{
			name: "decimal quantity",
			lineItem: &domain.LineItem{
				Quantity:  2.5,
				UnitPrice: 40,
			},
			effectiveSurcharge: 10,
			want:               110, // 100 * 1.10
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.FinalPrice(tt.lineItem, tt.effectiveSurcharge)
			if !floatEquals(got, tt.want) {
				t.Errorf("FinalPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

// floatEquals compares two floats with tolerance for floating point precision.
func floatEquals(a, b float64) bool {
	const epsilon = 0.0001
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < epsilon
}

func TestCalculateJobTotal(t *testing.T) {
	job := &domain.Job{
		ID:               "job-1",
		SurchargePercent: 10,
		SurchargeMode:    domain.SurchargeModeStacking,
	}

	categories := []*domain.Category{
		{ID: "cat-1", JobID: "job-1", ParentID: nil, SurchargePercent: floatPtr(5)},
	}

	lineItems := []*domain.LineItem{
		{
			ID:         "item-1",
			CategoryID: "cat-1",
			Type:       domain.LineItemTypeMaterial,
			Quantity:   10,
			UnitPrice:  100,
			// Base: 1000, Surcharge: 15%, Final: 1150
		},
		{
			ID:         "item-2",
			CategoryID: "cat-1",
			Type:       domain.LineItemTypeLabor,
			Quantity:   5,
			UnitPrice:  50,
			// Base: 250, Surcharge: 15%, Final: 287.5
		},
	}

	result := domain.CalculateJobTotal(job, categories, lineItems)

	expectedSubtotal := 1250.0  // 1000 + 250
	expectedGrandTotal := 1437.5 // 1150 + 287.5
	expectedMaterial := 1150.0
	expectedLabor := 287.5

	if result.Subtotal != expectedSubtotal {
		t.Errorf("Subtotal = %v, want %v", result.Subtotal, expectedSubtotal)
	}

	if result.GrandTotal != expectedGrandTotal {
		t.Errorf("GrandTotal = %v, want %v", result.GrandTotal, expectedGrandTotal)
	}

	if result.MaterialSubtotal != expectedMaterial {
		t.Errorf("MaterialSubtotal = %v, want %v", result.MaterialSubtotal, expectedMaterial)
	}

	if result.LaborSubtotal != expectedLabor {
		t.Errorf("LaborSubtotal = %v, want %v", result.LaborSubtotal, expectedLabor)
	}

	expectedSurchargeTotal := expectedGrandTotal - expectedSubtotal
	if result.SurchargeTotal != expectedSurchargeTotal {
		t.Errorf("SurchargeTotal = %v, want %v", result.SurchargeTotal, expectedSurchargeTotal)
	}
}
