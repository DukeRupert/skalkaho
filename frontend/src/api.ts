import type { QuoteResponse, QuoteLineItem, QuoteTotals, TemplateResult } from './types';

const FETCH_OPTS: RequestInit = {
  credentials: 'same-origin',
  headers: { 'Content-Type': 'application/json' },
};

async function apiFetch<T>(url: string, opts: RequestInit = {}): Promise<T> {
  const res = await fetch(url, { ...FETCH_OPTS, ...opts });
  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(body.error || `HTTP ${res.status}`);
  }
  return res.json();
}

export function fetchJob(jobId: string): Promise<QuoteResponse> {
  return apiFetch(`/api/jobs/${jobId}`);
}

export function patchItem(
  jobId: string,
  itemId: string,
  fields: Partial<QuoteLineItem>,
): Promise<{ item: QuoteLineItem; totals: QuoteTotals }> {
  return apiFetch(`/api/jobs/${jobId}/items/${itemId}`, {
    method: 'PATCH',
    body: JSON.stringify(fields),
  });
}

export function createItem(
  jobId: string,
  data: {
    category_id: string;
    type: string;
    name: string;
    quantity?: number;
    unit?: string;
    unit_price?: number;
    sort_order?: number;
  },
): Promise<{ item: QuoteLineItem; totals: QuoteTotals }> {
  return apiFetch(`/api/jobs/${jobId}/items`, {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export function deleteItem(
  jobId: string,
  itemId: string,
): Promise<{ totals: QuoteTotals }> {
  return apiFetch(`/api/jobs/${jobId}/items/${itemId}`, {
    method: 'DELETE',
  });
}

export function reorderItems(
  jobId: string,
  items: { id: string; sort_order: number }[],
): Promise<{ ok: boolean }> {
  return apiFetch(`/api/jobs/${jobId}/items/reorder`, {
    method: 'PATCH',
    body: JSON.stringify({ items }),
  });
}

export function searchTemplates(
  query: string,
  type?: string,
): Promise<TemplateResult[]> {
  const params = new URLSearchParams({ q: query });
  if (type) params.set('type', type);
  return apiFetch(`/api/items/search?${params}`);
}
