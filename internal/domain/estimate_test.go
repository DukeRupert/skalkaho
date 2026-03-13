package domain

import "testing"

func TestResolveMarkup(t *testing.T) {
	globals := MarkupGlobals{
		MaterialsMarkup: 20,
		LaborMarkup:     25,
		EquipmentMarkup: 15,
		SubsMarkup:      10,
		OtherMarkup:     10,
	}

	allEnabled := MarkupEnabled{
		Materials: true, Labor: true, Equipment: true, Subs: true, Other: true,
	}

	noOverrides := MarkupOverrides{}

	t.Run("falls back to global defaults", func(t *testing.T) {
		tests := []struct {
			catType  CategoryType
			expected float64
		}{
			{CategoryTypeMaterials, 20},
			{CategoryTypeLabor, 25},
			{CategoryTypeEquipment, 15},
			{CategoryTypeSubs, 10},
			{CategoryTypeOther, 10},
		}
		for _, tt := range tests {
			got := ResolveMarkup(tt.catType, globals, noOverrides, allEnabled)
			if got != tt.expected {
				t.Errorf("ResolveMarkup(%s) = %v, want %v", tt.catType, got, tt.expected)
			}
		}
	})

	t.Run("uses subcategory override when present", func(t *testing.T) {
		matOverride := 35.0
		overrides := MarkupOverrides{Materials: &matOverride}
		got := ResolveMarkup(CategoryTypeMaterials, globals, overrides, allEnabled)
		if got != 35 {
			t.Errorf("ResolveMarkup(materials with override) = %v, want 35", got)
		}

		// Non-overridden types still use global
		got = ResolveMarkup(CategoryTypeLabor, globals, overrides, allEnabled)
		if got != 25 {
			t.Errorf("ResolveMarkup(labor no override) = %v, want 25", got)
		}
	})

	t.Run("returns 0 when markup disabled", func(t *testing.T) {
		disabled := MarkupEnabled{
			Materials: false, Labor: true, Equipment: true, Subs: true, Other: true,
		}
		got := ResolveMarkup(CategoryTypeMaterials, globals, noOverrides, disabled)
		if got != 0 {
			t.Errorf("ResolveMarkup(disabled materials) = %v, want 0", got)
		}

		// Disabled overrides override — even with a subcategory override, disabled = 0%
		matOverride := 35.0
		overrides := MarkupOverrides{Materials: &matOverride}
		got = ResolveMarkup(CategoryTypeMaterials, globals, overrides, disabled)
		if got != 0 {
			t.Errorf("ResolveMarkup(disabled with override) = %v, want 0", got)
		}
	})

	t.Run("unknown category type returns 0", func(t *testing.T) {
		got := ResolveMarkup("invalid", globals, noOverrides, allEnabled)
		if got != 0 {
			t.Errorf("ResolveMarkup(invalid) = %v, want 0", got)
		}
	})
}
