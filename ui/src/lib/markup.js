/**
 * Markup calculation engine for the estimate builder.
 * Mirrors the Go domain.ResolveMarkup logic exactly.
 */

/**
 * Resolve the effective markup percentage for a line item.
 *
 * Resolution order (per spec Section 3.4):
 * 1. Check if the relevant type's markup_enabled toggle is false → 0%
 * 2. Check if the subcategory has a non-null override for that type → use it
 * 3. Fall back to the global default for that type
 *
 * @param {string} categoryType - 'materials' | 'labor' | 'equipment' | 'subs' | 'other'
 * @param {object} globals - { materials_markup, labor_markup, equipment_markup, subs_markup, other_markup }
 * @param {object} overrides - { materials, labor, equipment, subs, other } (null = inherit)
 * @param {object} enabled - { materials, labor, equipment, subs, other } (boolean)
 * @returns {number} effective markup percentage
 */
export function resolveMarkup(categoryType, globals, overrides, enabled) {
	const typeMap = {
		materials: { global: 'materials_markup', override: 'materials', enabled: 'materials' },
		labor: { global: 'labor_markup', override: 'labor', enabled: 'labor' },
		equipment: { global: 'equipment_markup', override: 'equipment', enabled: 'equipment' },
		subs: { global: 'subs_markup', override: 'subs', enabled: 'subs' },
		other: { global: 'other_markup', override: 'other', enabled: 'other' },
	};

	const mapping = typeMap[categoryType];
	if (!mapping) return 0;

	// Step 1: disabled = 0%
	if (!enabled[mapping.enabled]) return 0;

	// Step 2: subcategory override
	const override = overrides[mapping.override];
	if (override != null) return override;

	// Step 3: global default
	return globals[mapping.global];
}

/**
 * Calculate the after-markup unit price for a line item.
 * @param {number} unitPrice - base unit price
 * @param {number} markupPercent - effective markup percentage
 * @returns {number} unit price with markup applied
 */
export function afterMarkupPrice(unitPrice, markupPercent) {
	return unitPrice * (1 + markupPercent / 100);
}

/**
 * Calculate line item total (quantity * after-markup unit price).
 * @param {number} quantity
 * @param {number} unitPrice
 * @param {number} markupPercent
 * @returns {number}
 */
export function lineItemTotal(quantity, unitPrice, markupPercent) {
	return quantity * afterMarkupPrice(unitPrice, markupPercent);
}

/**
 * Calculate totals for a subcategory.
 * @param {object} subcat - subcategory with line_items, component_groups, lump_sum
 * @param {object} globals - global markup percentages
 * @returns {{ base: number, withMarkup: number, byType: object }}
 */
export function subcategoryTotals(subcat, globals) {
	const byType = { materials: 0, labor: 0, equipment: 0, subs: 0, other: 0 };
	let base = 0;
	let withMarkup = 0;

	const processItem = (item) => {
		const markup = resolveMarkup(
			item.category_type,
			globals,
			subcat.markup_overrides,
			subcat.markup_enabled,
		);
		const itemBase = item.quantity * item.unit_price;
		const itemTotal = lineItemTotal(item.quantity, item.unit_price, markup);
		base += itemBase;
		withMarkup += itemTotal;
		if (byType[item.category_type] !== undefined) {
			byType[item.category_type] += itemTotal;
		}
	};

	// Ungrouped items
	for (const item of subcat.line_items) {
		processItem(item);
	}

	// Grouped items
	for (const group of subcat.component_groups) {
		for (const item of group.line_items) {
			processItem(item);
		}
	}

	// Add lump sum
	withMarkup += subcat.lump_sum;

	return { base, withMarkup, byType };
}

/**
 * Calculate totals for a section.
 * @param {object} section - section with subcategories
 * @param {object} globals - global markup percentages
 * @returns {{ base: number, withMarkup: number, byType: object }}
 */
export function sectionTotals(section, globals) {
	const byType = { materials: 0, labor: 0, equipment: 0, subs: 0, other: 0 };
	let base = 0;
	let withMarkup = 0;

	for (const subcat of section.subcategories) {
		const st = subcategoryTotals(subcat, globals);
		base += st.base;
		withMarkup += st.withMarkup;
		for (const type of Object.keys(byType)) {
			byType[type] += st.byType[type];
		}
	}

	return { base, withMarkup, byType };
}

/**
 * Calculate totals for the entire estimate.
 * @param {object} estimate - full estimate payload
 * @returns {{ base: number, withMarkup: number, byType: object }}
 */
export function estimateTotals(estimate) {
	const byType = { materials: 0, labor: 0, equipment: 0, subs: 0, other: 0 };
	let base = 0;
	let withMarkup = 0;

	for (const section of estimate.sections) {
		const st = sectionTotals(section, estimate.globals);
		base += st.base;
		withMarkup += st.withMarkup;
		for (const type of Object.keys(byType)) {
			byType[type] += st.byType[type];
		}
	}

	return { base, withMarkup, byType };
}

/**
 * Format a number as currency.
 * @param {number} amount
 * @returns {string}
 */
export function formatMoney(amount) {
	return new Intl.NumberFormat('en-US', {
		style: 'currency',
		currency: 'USD',
		minimumFractionDigits: 2,
	}).format(amount);
}

/**
 * Format a number as a percentage.
 * @param {number} amount
 * @returns {string}
 */
export function formatPercent(amount) {
	return `${amount}%`;
}
