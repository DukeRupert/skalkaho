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

// Test helper functions for cleaner test setup
func stringPtr(s string) *string {
	return &s
}

func makeJob(id string, surcharge float64, mode domain.SurchargeMode) *domain.Job {
	return &domain.Job{
		ID:               id,
		SurchargePercent: surcharge,
		SurchargeMode:    mode,
	}
}

func makeCategory(id, jobID string, parentID *string, surcharge *float64) *domain.Category {
	return &domain.Category{
		ID:               id,
		JobID:            jobID,
		ParentID:         parentID,
		SurchargePercent: surcharge,
	}
}

func makeLineItem(id, categoryID string, itemType domain.LineItemType, qty, price float64) *domain.LineItem {
	return &domain.LineItem{
		ID:         id,
		CategoryID: categoryID,
		Type:       itemType,
		Quantity:   qty,
		UnitPrice:  price,
	}
}

func TestCalculateJobTotal_ThreeLevelNestedCategories(t *testing.T) {
	// Job 10% → Cat-L1 5% → Cat-L2 3% → Cat-L3 2%
	// Stacking mode: total surcharge = 10 + 5 + 3 + 2 = 20% for L3 items
	job := makeJob("job-1", 10, domain.SurchargeModeStacking)

	categories := []*domain.Category{
		makeCategory("cat-l1", "job-1", nil, floatPtr(5)),
		makeCategory("cat-l2", "job-1", stringPtr("cat-l1"), floatPtr(3)),
		makeCategory("cat-l3", "job-1", stringPtr("cat-l2"), floatPtr(2)),
	}

	lineItems := []*domain.LineItem{
		// L1 item: Base 100, Surcharge 15% (10+5), Final 115
		makeLineItem("item-l1", "cat-l1", domain.LineItemTypeMaterial, 1, 100),
		// L2 item: Base 200, Surcharge 18% (10+5+3), Final 236
		makeLineItem("item-l2", "cat-l2", domain.LineItemTypeLabor, 2, 100),
		// L3 item: Base 300, Surcharge 20% (10+5+3+2), Final 360
		makeLineItem("item-l3", "cat-l3", domain.LineItemTypeMaterial, 3, 100),
	}

	result := domain.CalculateJobTotal(job, categories, lineItems)

	// Subtotal: 100 + 200 + 300 = 600
	if !floatEquals(result.Subtotal, 600) {
		t.Errorf("Subtotal = %v, want 600", result.Subtotal)
	}

	// GrandTotal: 115 + 236 + 360 = 711
	if !floatEquals(result.GrandTotal, 711) {
		t.Errorf("GrandTotal = %v, want 711", result.GrandTotal)
	}

	// MaterialSubtotal: 115 + 360 = 475
	if !floatEquals(result.MaterialSubtotal, 475) {
		t.Errorf("MaterialSubtotal = %v, want 475", result.MaterialSubtotal)
	}

	// LaborSubtotal: 236
	if !floatEquals(result.LaborSubtotal, 236) {
		t.Errorf("LaborSubtotal = %v, want 236", result.LaborSubtotal)
	}

	// SurchargeTotal: 711 - 600 = 111
	if !floatEquals(result.SurchargeTotal, 111) {
		t.Errorf("SurchargeTotal = %v, want 111", result.SurchargeTotal)
	}
}

func TestCalculateJobTotal_ThreeLevelNestedCategories_OverrideMode(t *testing.T) {
	// Override mode: deepest surcharge wins
	job := makeJob("job-1", 10, domain.SurchargeModeOverride)

	categories := []*domain.Category{
		makeCategory("cat-l1", "job-1", nil, floatPtr(5)),
		makeCategory("cat-l2", "job-1", stringPtr("cat-l1"), floatPtr(3)),
		makeCategory("cat-l3", "job-1", stringPtr("cat-l2"), floatPtr(2)),
	}

	lineItems := []*domain.LineItem{
		// L1 item: Base 100, Surcharge 5% (cat-l1 overrides job), Final 105
		makeLineItem("item-l1", "cat-l1", domain.LineItemTypeMaterial, 1, 100),
		// L2 item: Base 200, Surcharge 3% (cat-l2 is deepest), Final 206
		makeLineItem("item-l2", "cat-l2", domain.LineItemTypeLabor, 2, 100),
		// L3 item: Base 300, Surcharge 2% (cat-l3 is deepest), Final 306
		makeLineItem("item-l3", "cat-l3", domain.LineItemTypeMaterial, 3, 100),
	}

	result := domain.CalculateJobTotal(job, categories, lineItems)

	// GrandTotal: 105 + 206 + 306 = 617
	if !floatEquals(result.GrandTotal, 617) {
		t.Errorf("GrandTotal = %v, want 617", result.GrandTotal)
	}
}

func TestCalculateJobTotal_MultipleCategories(t *testing.T) {
	// Multiple root-level categories with items in each
	job := makeJob("job-1", 10, domain.SurchargeModeStacking)

	categories := []*domain.Category{
		makeCategory("cat-a", "job-1", nil, floatPtr(5)),
		makeCategory("cat-b", "job-1", nil, floatPtr(8)),
		makeCategory("cat-c", "job-1", nil, nil), // No surcharge, inherits job only
	}

	lineItems := []*domain.LineItem{
		// Cat-A: 15% surcharge
		makeLineItem("item-a1", "cat-a", domain.LineItemTypeMaterial, 10, 10), // Base 100, Final 115
		makeLineItem("item-a2", "cat-a", domain.LineItemTypeLabor, 5, 20),     // Base 100, Final 115
		// Cat-B: 18% surcharge
		makeLineItem("item-b1", "cat-b", domain.LineItemTypeMaterial, 2, 50), // Base 100, Final 118
		// Cat-C: 10% surcharge (job only)
		makeLineItem("item-c1", "cat-c", domain.LineItemTypeLabor, 4, 25), // Base 100, Final 110
	}

	result := domain.CalculateJobTotal(job, categories, lineItems)

	// Subtotal: 100 + 100 + 100 + 100 = 400
	if !floatEquals(result.Subtotal, 400) {
		t.Errorf("Subtotal = %v, want 400", result.Subtotal)
	}

	// GrandTotal: 115 + 115 + 118 + 110 = 458
	if !floatEquals(result.GrandTotal, 458) {
		t.Errorf("GrandTotal = %v, want 458", result.GrandTotal)
	}

	// MaterialSubtotal: 115 + 118 = 233
	if !floatEquals(result.MaterialSubtotal, 233) {
		t.Errorf("MaterialSubtotal = %v, want 233", result.MaterialSubtotal)
	}

	// LaborSubtotal: 115 + 110 = 225
	if !floatEquals(result.LaborSubtotal, 225) {
		t.Errorf("LaborSubtotal = %v, want 225", result.LaborSubtotal)
	}
}

func TestCalculateJobTotal_TypeBreakdown(t *testing.T) {
	tests := []struct {
		name             string
		lineItems        []*domain.LineItem
		wantMaterial     float64
		wantLabor        float64
		wantGrandTotal   float64
	}{
		{
			name: "all materials",
			lineItems: []*domain.LineItem{
				makeLineItem("m1", "cat-1", domain.LineItemTypeMaterial, 10, 10),
				makeLineItem("m2", "cat-1", domain.LineItemTypeMaterial, 5, 20),
				makeLineItem("m3", "cat-1", domain.LineItemTypeMaterial, 2, 50),
			},
			wantMaterial:   345, // (100+100+100) * 1.15
			wantLabor:      0,
			wantGrandTotal: 345,
		},
		{
			name: "all labor",
			lineItems: []*domain.LineItem{
				makeLineItem("l1", "cat-1", domain.LineItemTypeLabor, 8, 25),
				makeLineItem("l2", "cat-1", domain.LineItemTypeLabor, 4, 50),
			},
			wantMaterial:   0,
			wantLabor:      460, // (200+200) * 1.15
			wantGrandTotal: 460,
		},
		{
			name: "mixed types",
			lineItems: []*domain.LineItem{
				makeLineItem("m1", "cat-1", domain.LineItemTypeMaterial, 10, 10), // 100 * 1.15 = 115
				makeLineItem("l1", "cat-1", domain.LineItemTypeLabor, 5, 30),     // 150 * 1.15 = 172.5
				makeLineItem("m2", "cat-1", domain.LineItemTypeMaterial, 2, 25),  // 50 * 1.15 = 57.5
			},
			wantMaterial:   172.5,
			wantLabor:      172.5,
			wantGrandTotal: 345,
		},
	}

	job := makeJob("job-1", 10, domain.SurchargeModeStacking)
	categories := []*domain.Category{
		makeCategory("cat-1", "job-1", nil, floatPtr(5)),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := domain.CalculateJobTotal(job, categories, tt.lineItems)

			if !floatEquals(result.MaterialSubtotal, tt.wantMaterial) {
				t.Errorf("MaterialSubtotal = %v, want %v", result.MaterialSubtotal, tt.wantMaterial)
			}
			if !floatEquals(result.LaborSubtotal, tt.wantLabor) {
				t.Errorf("LaborSubtotal = %v, want %v", result.LaborSubtotal, tt.wantLabor)
			}
			if !floatEquals(result.GrandTotal, tt.wantGrandTotal) {
				t.Errorf("GrandTotal = %v, want %v", result.GrandTotal, tt.wantGrandTotal)
			}
		})
	}
}

func TestCalculateJobTotal_CategoryRemoval(t *testing.T) {
	// Simulate removing a category by recalculating without its items
	job := makeJob("job-1", 10, domain.SurchargeModeStacking)

	categories := []*domain.Category{
		makeCategory("cat-a", "job-1", nil, floatPtr(5)),
		makeCategory("cat-b", "job-1", nil, floatPtr(5)),
	}

	allItems := []*domain.LineItem{
		makeLineItem("item-a1", "cat-a", domain.LineItemTypeMaterial, 10, 10), // 115
		makeLineItem("item-a2", "cat-a", domain.LineItemTypeLabor, 5, 20),     // 115
		makeLineItem("item-b1", "cat-b", domain.LineItemTypeMaterial, 2, 50),  // 115
		makeLineItem("item-b2", "cat-b", domain.LineItemTypeLabor, 4, 25),     // 115
	}

	// Calculate with all items
	fullResult := domain.CalculateJobTotal(job, categories, allItems)

	// Remove cat-b items (simulate category deletion)
	catAItemsOnly := []*domain.LineItem{
		makeLineItem("item-a1", "cat-a", domain.LineItemTypeMaterial, 10, 10),
		makeLineItem("item-a2", "cat-a", domain.LineItemTypeLabor, 5, 20),
	}

	reducedResult := domain.CalculateJobTotal(job, categories, catAItemsOnly)

	// Full: 4 items * 115 = 460
	if !floatEquals(fullResult.GrandTotal, 460) {
		t.Errorf("Full GrandTotal = %v, want 460", fullResult.GrandTotal)
	}

	// Reduced: 2 items * 115 = 230
	if !floatEquals(reducedResult.GrandTotal, 230) {
		t.Errorf("Reduced GrandTotal = %v, want 230", reducedResult.GrandTotal)
	}

	// Difference should equal the removed items' contribution
	removedContribution := fullResult.GrandTotal - reducedResult.GrandTotal
	if !floatEquals(removedContribution, 230) {
		t.Errorf("Removed contribution = %v, want 230", removedContribution)
	}
}

func TestCalculateCategoryTotal(t *testing.T) {
	job := makeJob("job-1", 10, domain.SurchargeModeStacking)

	categories := []*domain.Category{
		makeCategory("cat-root", "job-1", nil, floatPtr(5)),
		makeCategory("cat-l2", "job-1", stringPtr("cat-root"), floatPtr(3)),
		makeCategory("cat-l3", "job-1", stringPtr("cat-l2"), floatPtr(2)),
	}

	lineItems := []*domain.LineItem{
		// Root: Base 100, 15% surcharge, Final 115
		makeLineItem("item-root", "cat-root", domain.LineItemTypeMaterial, 1, 100),
		// L2: Base 200, 18% surcharge, Final 236
		makeLineItem("item-l2", "cat-l2", domain.LineItemTypeLabor, 2, 100),
		// L3: Base 300, 20% surcharge, Final 360
		makeLineItem("item-l3", "cat-l3", domain.LineItemTypeMaterial, 3, 100),
	}

	t.Run("root category includes all descendants", func(t *testing.T) {
		result := domain.CalculateCategoryTotal("cat-root", job, categories, lineItems)

		// Root total should include all items: 115 + 236 + 360 = 711
		if !floatEquals(result.Total, 711) {
			t.Errorf("Root Total = %v, want 711", result.Total)
		}

		// Subtotal: 100 + 200 + 300 = 600
		if !floatEquals(result.Subtotal, 600) {
			t.Errorf("Root Subtotal = %v, want 600", result.Subtotal)
		}

		if result.CategoryID != "cat-root" {
			t.Errorf("CategoryID = %v, want cat-root", result.CategoryID)
		}
	})

	t.Run("L2 category includes L3 items", func(t *testing.T) {
		result := domain.CalculateCategoryTotal("cat-l2", job, categories, lineItems)

		// L2 total: 236 (L2 item) + 360 (L3 item) = 596
		if !floatEquals(result.Total, 596) {
			t.Errorf("L2 Total = %v, want 596", result.Total)
		}
	})

	t.Run("L3 category only includes its own items", func(t *testing.T) {
		result := domain.CalculateCategoryTotal("cat-l3", job, categories, lineItems)

		// L3 total: 360 (only L3 item)
		if !floatEquals(result.Total, 360) {
			t.Errorf("L3 Total = %v, want 360", result.Total)
		}
	})

	t.Run("empty category returns zero", func(t *testing.T) {
		emptyCats := []*domain.Category{
			makeCategory("empty-cat", "job-1", nil, floatPtr(5)),
		}
		result := domain.CalculateCategoryTotal("empty-cat", job, emptyCats, []*domain.LineItem{})

		if result.Total != 0 {
			t.Errorf("Empty category Total = %v, want 0", result.Total)
		}
		if result.Subtotal != 0 {
			t.Errorf("Empty category Subtotal = %v, want 0", result.Subtotal)
		}
	})

	t.Run("category with only subcategory items", func(t *testing.T) {
		// Parent has no direct items, only child has items
		parentChild := []*domain.Category{
			makeCategory("parent", "job-1", nil, floatPtr(5)),
			makeCategory("child", "job-1", stringPtr("parent"), floatPtr(3)),
		}
		childItems := []*domain.LineItem{
			makeLineItem("child-item", "child", domain.LineItemTypeMaterial, 10, 10), // Base 100, 18%, Final 118
		}

		result := domain.CalculateCategoryTotal("parent", job, parentChild, childItems)

		// Parent total should include child item: 118
		if !floatEquals(result.Total, 118) {
			t.Errorf("Parent Total = %v, want 118", result.Total)
		}
	})
}

func TestCalculateJobTotal_EdgeCases(t *testing.T) {
	job := makeJob("job-1", 10, domain.SurchargeModeStacking)
	categories := []*domain.Category{
		makeCategory("cat-1", "job-1", nil, floatPtr(5)),
	}

	t.Run("empty job no items", func(t *testing.T) {
		result := domain.CalculateJobTotal(job, categories, []*domain.LineItem{})

		if result.Subtotal != 0 {
			t.Errorf("Subtotal = %v, want 0", result.Subtotal)
		}
		if result.GrandTotal != 0 {
			t.Errorf("GrandTotal = %v, want 0", result.GrandTotal)
		}
		if result.MaterialSubtotal != 0 {
			t.Errorf("MaterialSubtotal = %v, want 0", result.MaterialSubtotal)
		}
		if result.LaborSubtotal != 0 {
			t.Errorf("LaborSubtotal = %v, want 0", result.LaborSubtotal)
		}
	})

	t.Run("zero quantity item", func(t *testing.T) {
		items := []*domain.LineItem{
			makeLineItem("zero-qty", "cat-1", domain.LineItemTypeMaterial, 0, 100),
		}
		result := domain.CalculateJobTotal(job, categories, items)

		if result.GrandTotal != 0 {
			t.Errorf("GrandTotal = %v, want 0", result.GrandTotal)
		}
	})

	t.Run("zero price item", func(t *testing.T) {
		items := []*domain.LineItem{
			makeLineItem("zero-price", "cat-1", domain.LineItemTypeMaterial, 10, 0),
		}
		result := domain.CalculateJobTotal(job, categories, items)

		if result.GrandTotal != 0 {
			t.Errorf("GrandTotal = %v, want 0", result.GrandTotal)
		}
	})

	t.Run("large numbers precision", func(t *testing.T) {
		items := []*domain.LineItem{
			makeLineItem("large", "cat-1", domain.LineItemTypeMaterial, 1000000, 999.99),
		}
		result := domain.CalculateJobTotal(job, categories, items)

		// Base: 1,000,000 * 999.99 = 999,990,000
		// With 15% surcharge: 1,149,988,500
		expectedGrandTotal := 999990000.0 * 1.15
		if !floatEquals(result.GrandTotal, expectedGrandTotal) {
			t.Errorf("GrandTotal = %v, want %v", result.GrandTotal, expectedGrandTotal)
		}
	})

	t.Run("fractional quantities and prices", func(t *testing.T) {
		items := []*domain.LineItem{
			makeLineItem("fractional", "cat-1", domain.LineItemTypeMaterial, 2.5, 33.33),
		}
		result := domain.CalculateJobTotal(job, categories, items)

		// Base: 2.5 * 33.33 = 83.325
		// With 15% surcharge: 95.82375
		expectedGrandTotal := 83.325 * 1.15
		if !floatEquals(result.GrandTotal, expectedGrandTotal) {
			t.Errorf("GrandTotal = %v, want %v", result.GrandTotal, expectedGrandTotal)
		}
	})

	t.Run("zero job surcharge", func(t *testing.T) {
		zeroJob := makeJob("job-zero", 0, domain.SurchargeModeStacking)
		zeroCats := []*domain.Category{
			makeCategory("cat-zero", "job-zero", nil, nil), // No category surcharge either
		}
		items := []*domain.LineItem{
			makeLineItem("item", "cat-zero", domain.LineItemTypeMaterial, 10, 10),
		}

		result := domain.CalculateJobTotal(zeroJob, zeroCats, items)

		// No surcharge: base = final = 100
		if !floatEquals(result.GrandTotal, 100) {
			t.Errorf("GrandTotal = %v, want 100", result.GrandTotal)
		}
		if !floatEquals(result.SurchargeTotal, 0) {
			t.Errorf("SurchargeTotal = %v, want 0", result.SurchargeTotal)
		}
	})
}
