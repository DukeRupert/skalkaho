// Client-side replication of internal/domain/surcharge.go
// Must produce identical results to the Go domain layer.

import type { QuoteJob, QuoteCategory, QuoteLineItem, QuoteItemType, QuoteTotals } from './types';

/** Get the surcharge for a specific item type (mirrors domain.GetTypeSurcharge). */
export function getTypeSurcharge(
  job: QuoteJob,
  itemType: string,
  customTypes: QuoteItemType[],
): number {
  switch (itemType) {
    case 'material':
      if (job.material_surcharge_percent != null) return job.material_surcharge_percent;
      break;
    case 'labor':
      if (job.labor_surcharge_percent != null) return job.labor_surcharge_percent;
      break;
    case 'equipment':
      if (job.equipment_surcharge_percent != null) return job.equipment_surcharge_percent;
      break;
    default:
      for (const ct of customTypes) {
        if (ct.slug === itemType && ct.surcharge_percent != null) {
          return ct.surcharge_percent;
        }
      }
  }
  return job.surcharge_percent;
}

/** Build category chain from root to target (mirrors domain.buildCategoryChain). */
function buildCategoryChain(
  categoryId: string,
  categoryById: Map<string, QuoteCategory>,
): QuoteCategory[] {
  const chain: QuoteCategory[] = [];
  let current = categoryById.get(categoryId);
  while (current) {
    chain.unshift(current);
    if (current.parent_id == null) break;
    current = categoryById.get(current.parent_id);
  }
  return chain;
}

/** Effective surcharge in override mode (most specific wins). */
function effectiveSurchargeOverride(
  item: QuoteLineItem,
  job: QuoteJob,
  chain: QuoteCategory[],
  customTypes: QuoteItemType[],
): number {
  if (item.surcharge_percent != null) return item.surcharge_percent;
  for (let i = chain.length - 1; i >= 0; i--) {
    if (chain[i].surcharge_percent != null) return chain[i].surcharge_percent!;
  }
  return getTypeSurcharge(job, item.type, customTypes);
}

/** Effective surcharge in stacking mode (all add together). */
function effectiveSurchargeStacking(
  item: QuoteLineItem,
  job: QuoteJob,
  chain: QuoteCategory[],
  customTypes: QuoteItemType[],
): number {
  let total = getTypeSurcharge(job, item.type, customTypes);
  for (const cat of chain) {
    if (cat.surcharge_percent != null) total += cat.surcharge_percent;
  }
  if (item.surcharge_percent != null) total += item.surcharge_percent;
  return total;
}

/** Calculate effective surcharge for a line item (mirrors domain.EffectiveSurcharge). */
export function effectiveSurcharge(
  item: QuoteLineItem,
  job: QuoteJob,
  chain: QuoteCategory[],
  customTypes: QuoteItemType[],
): number {
  if (job.surcharge_mode === 'override') {
    return effectiveSurchargeOverride(item, job, chain, customTypes);
  }
  return effectiveSurchargeStacking(item, job, chain, customTypes);
}

/** Flatten all categories into a map by ID. */
function flattenCategories(categories: QuoteCategory[]): Map<string, QuoteCategory> {
  const map = new Map<string, QuoteCategory>();
  function walk(cats: QuoteCategory[]) {
    for (const cat of cats) {
      map.set(cat.id, cat);
      walk(cat.children);
    }
  }
  walk(categories);
  return map;
}

/** Collect all line items from nested category tree. */
function collectAllItems(categories: QuoteCategory[]): QuoteLineItem[] {
  const items: QuoteLineItem[] = [];
  function walk(cats: QuoteCategory[]) {
    for (const cat of cats) {
      items.push(...cat.items);
      walk(cat.children);
    }
  }
  walk(categories);
  return items;
}

/** Calculate full job totals (mirrors domain.CalculateJobTotal). */
export function calculateTotals(
  job: QuoteJob,
  categories: QuoteCategory[],
  customTypes: QuoteItemType[],
): QuoteTotals {
  const categoryById = flattenCategories(categories);
  const chainCache = new Map<string, QuoteCategory[]>();
  const allItems = collectAllItems(categories);

  let subtotal = 0;
  let grandTotal = 0;
  const typeSubtotals: Record<string, number> = {};

  for (const item of allItems) {
    let chain = chainCache.get(item.category_id);
    if (!chain) {
      chain = buildCategoryChain(item.category_id, categoryById);
      chainCache.set(item.category_id, chain);
    }

    const basePrice = item.quantity * item.unit_price;
    const effSurcharge = effectiveSurcharge(item, job, chain, customTypes);
    const finalPrice = basePrice * (1 + effSurcharge / 100);

    subtotal += basePrice;
    grandTotal += finalPrice;
    typeSubtotals[item.type] = (typeSubtotals[item.type] || 0) + finalPrice;
  }

  return {
    subtotal,
    surcharge_total: grandTotal - subtotal,
    grand_total: grandTotal,
    type_subtotals: typeSubtotals,
  };
}
