// Wire types mirroring Go QuoteResponse types exactly

export interface QuoteResponse {
  job: QuoteJob;
  categories: QuoteCategory[];
  custom_types: QuoteItemType[];
  totals: QuoteTotals;
}

export interface QuoteJob {
  id: string;
  name: string;
  status: string;
  surcharge_percent: number;
  surcharge_mode: 'stacking' | 'override';
  material_surcharge_percent: number | null;
  labor_surcharge_percent: number | null;
  equipment_surcharge_percent: number | null;
}

export interface QuoteCategory {
  id: string;
  name: string;
  parent_id: string | null;
  sort_order: number;
  surcharge_percent: number | null;
  items: QuoteLineItem[];
  children: QuoteCategory[];
}

export interface QuoteLineItem {
  id: string;
  category_id: string;
  type: string;
  name: string;
  description: string | null;
  quantity: number;
  unit: string;
  unit_price: number;
  surcharge_percent: number | null;
  sort_order: number;
  tag: string | null;
}

export interface QuoteItemType {
  slug: string;
  name: string;
  color: string;
  surcharge_percent: number | null;
}

export interface QuoteTotals {
  subtotal: number;
  surcharge_total: number;
  grand_total: number;
  type_subtotals: Record<string, number>;
}

// Grid-specific types

export type EditableColumn = 'name' | 'quantity' | 'unit' | 'unit_price' | 'surcharge_percent';

export interface GridRow {
  kind: 'category' | 'item';
  categoryId: string;
  categoryName?: string;
  depth: number;
  item?: QuoteLineItem;
}

export interface CellPosition {
  row: number;
  col: number;
}

export interface TemplateResult {
  id: number;
  type: string;
  category: string;
  name: string;
  default_unit: string;
  default_price: number;
}
