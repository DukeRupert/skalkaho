package domain

// EffectiveSurcharge calculates the applicable surcharge for a line item
// based on the job's surcharge mode and the category hierarchy.
func EffectiveSurcharge(li *LineItem, job *Job, categoryChain []*Category) float64 {
	if job.SurchargeMode == SurchargeModeOverride {
		return effectiveSurchargeOverride(li, job, categoryChain)
	}
	return effectiveSurchargeStacking(li, job, categoryChain)
}

// effectiveSurchargeOverride returns the most specific (lowest-level) surcharge.
// Priority: LineItem > deepest Category > ... > shallowest Category > Job
func effectiveSurchargeOverride(li *LineItem, job *Job, categoryChain []*Category) float64 {
	// Check line item first
	if li.SurchargePercent != nil {
		return *li.SurchargePercent
	}

	// Walk category chain from deepest to shallowest
	for i := len(categoryChain) - 1; i >= 0; i-- {
		if categoryChain[i].SurchargePercent != nil {
			return *categoryChain[i].SurchargePercent
		}
	}

	// Fall back to job surcharge
	return job.SurchargePercent
}

// effectiveSurchargeStacking sums all surcharges in the hierarchy.
// Total = Job% + Category%s + LineItem%
func effectiveSurchargeStacking(li *LineItem, job *Job, categoryChain []*Category) float64 {
	total := job.SurchargePercent

	// Add all category surcharges
	for _, cat := range categoryChain {
		if cat.SurchargePercent != nil {
			total += *cat.SurchargePercent
		}
	}

	// Add line item surcharge
	if li.SurchargePercent != nil {
		total += *li.SurchargePercent
	}

	return total
}

// FinalPrice calculates the line item total with surcharge applied.
func FinalPrice(li *LineItem, effectiveSurcharge float64) float64 {
	base := li.BasePrice()
	return base * (1 + effectiveSurcharge/100)
}

// CategoryTotal calculates the total for a category including all line items and child categories.
type CategoryTotal struct {
	CategoryID     string  `json:"category_id"`
	Subtotal       float64 `json:"subtotal"`        // Sum of base prices
	SurchargeTotal float64 `json:"surcharge_total"` // Sum of surcharges
	Total          float64 `json:"total"`           // Final total
}

// JobTotal calculates the complete job totals.
type JobTotal struct {
	Subtotal           float64 `json:"subtotal"`            // Sum of all base prices
	SurchargeTotal     float64 `json:"surcharge_total"`     // Total surcharges applied
	GrandTotal         float64 `json:"grand_total"`         // Final total
	MaterialSubtotal   float64 `json:"material_subtotal"`   // Materials only
	LaborSubtotal      float64 `json:"labor_subtotal"`      // Labor only
	EquipmentSubtotal  float64 `json:"equipment_subtotal"`  // Equipment only
}

// CalculateJobTotal computes all totals for a job.
func CalculateJobTotal(job *Job, categories []*Category, lineItems []*LineItem) JobTotal {
	var result JobTotal

	// Build category lookup for chain resolution
	categoryByID := make(map[string]*Category)
	for _, cat := range categories {
		categoryByID[cat.ID] = cat
	}

	// Build category chain for each line item's category
	categoryChains := make(map[string][]*Category)

	for _, li := range lineItems {
		// Get or build category chain
		chain, exists := categoryChains[li.CategoryID]
		if !exists {
			chain = buildCategoryChain(li.CategoryID, categoryByID)
			categoryChains[li.CategoryID] = chain
		}

		// Calculate effective surcharge and prices
		basePrice := li.BasePrice()
		effSurcharge := EffectiveSurcharge(li, job, chain)
		finalPrice := FinalPrice(li, effSurcharge)

		result.Subtotal += basePrice
		result.GrandTotal += finalPrice

		// Track by type
		switch li.Type {
		case LineItemTypeMaterial:
			result.MaterialSubtotal += finalPrice
		case LineItemTypeLabor:
			result.LaborSubtotal += finalPrice
		case LineItemTypeEquipment:
			result.EquipmentSubtotal += finalPrice
		}
	}

	result.SurchargeTotal = result.GrandTotal - result.Subtotal

	return result
}

// buildCategoryChain builds the chain from root to the specified category.
func buildCategoryChain(categoryID string, categoryByID map[string]*Category) []*Category {
	var chain []*Category
	current := categoryByID[categoryID]

	for current != nil {
		chain = append([]*Category{current}, chain...) // Prepend to build root-first
		if current.ParentID == nil {
			break
		}
		current = categoryByID[*current.ParentID]
	}

	return chain
}

// CalculateCategoryTotal computes totals for a category including all nested subcategories.
func CalculateCategoryTotal(categoryID string, job *Job, categories []*Category, lineItems []*LineItem) CategoryTotal {
	var result CategoryTotal
	result.CategoryID = categoryID

	// Build category lookup
	categoryByID := make(map[string]*Category)
	for _, cat := range categories {
		categoryByID[cat.ID] = cat
	}

	// Find all descendant category IDs (including the target category itself)
	descendantIDs := findDescendantCategories(categoryID, categories)
	descendantIDs[categoryID] = true

	// Build category chains cache
	categoryChains := make(map[string][]*Category)

	for _, li := range lineItems {
		// Only include items from this category or its descendants
		if !descendantIDs[li.CategoryID] {
			continue
		}

		// Get or build category chain
		chain, exists := categoryChains[li.CategoryID]
		if !exists {
			chain = buildCategoryChain(li.CategoryID, categoryByID)
			categoryChains[li.CategoryID] = chain
		}

		// Calculate prices
		basePrice := li.BasePrice()
		effSurcharge := EffectiveSurcharge(li, job, chain)
		finalPrice := FinalPrice(li, effSurcharge)

		result.Subtotal += basePrice
		result.Total += finalPrice
	}

	result.SurchargeTotal = result.Total - result.Subtotal

	return result
}

// findDescendantCategories returns a set of all category IDs that are descendants of the given category.
func findDescendantCategories(parentID string, categories []*Category) map[string]bool {
	result := make(map[string]bool)

	// Build children lookup
	childrenOf := make(map[string][]string)
	for _, cat := range categories {
		if cat.ParentID != nil {
			childrenOf[*cat.ParentID] = append(childrenOf[*cat.ParentID], cat.ID)
		}
	}

	// BFS to find all descendants
	queue := childrenOf[parentID]
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		result[current] = true
		queue = append(queue, childrenOf[current]...)
	}

	return result
}
