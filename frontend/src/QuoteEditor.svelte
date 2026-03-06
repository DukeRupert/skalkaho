<script lang="ts">
  import { onMount } from 'svelte';
  import { fetchJob, patchItem, createItem, deleteItem, searchTemplates } from './api';
  import { calculateTotals } from './surcharge';
  import { autofocus } from './actions';
  import type {
    QuoteResponse,
    QuoteCategory,
    QuoteLineItem,
    QuoteItemType,
    QuoteTotals,
    GridRow,
    EditableColumn,
    TemplateResult,
  } from './types';

  let { jobId }: { jobId: string } = $props();

  // --- State ---
  let data: QuoteResponse | null = $state(null);
  let error: string | null = $state(null);
  let loading = $state(true);
  let selectedRow = $state(-1);
  let selectedCol = $state(0);
  let editing = $state(false);
  let editValue = $state('');
  let gridEl: HTMLDivElement | undefined = $state();

  // --- Item search modal state ---
  let searchOpen = $state(false);
  let searchType = $state('material');
  let searchCategoryId = $state('');
  let searchQuery = $state('');
  let searchResults: TemplateResult[] = $state([]);
  let searchSelectedIdx = $state(0);
  let searchLoading = $state(false);
  let searchDebounceTimer: ReturnType<typeof setTimeout> | undefined;
  let searchResultsEl: HTMLDivElement | undefined = $state();

  const COLUMNS: EditableColumn[] = ['name', 'quantity', 'unit', 'unit_price', 'surcharge_percent'];
  const COL_HEADERS = ['Name', 'Qty', 'Unit', 'Price', 'Markup %'];
  const COL_WIDTHS = ['flex-[3]', 'w-20', 'w-20', 'w-24', 'w-20'];

  // --- Derived ---
  let gridRows: GridRow[] = $derived.by(() => {
    if (!data) return [];
    return buildFlatRows(data.categories, 0);
  });

  let liveTotals: QuoteTotals = $derived.by(() => {
    if (!data) return { subtotal: 0, surcharge_total: 0, grand_total: 0, type_subtotals: {} };
    return calculateTotals(data.job, data.categories, data.custom_types);
  });

  // --- Helpers ---
  function buildFlatRows(categories: QuoteCategory[], depth: number): GridRow[] {
    const rows: GridRow[] = [];
    for (const cat of categories) {
      rows.push({ kind: 'category', categoryId: cat.id, categoryName: cat.name, depth });
      for (const item of cat.items) {
        rows.push({ kind: 'item', categoryId: cat.id, depth: depth + 1, item });
      }
      rows.push(...buildFlatRows(cat.children, depth + 1));
    }
    return rows;
  }

  function getTypeColor(type: string): string {
    switch (type) {
      case 'material': return '#2d5a47';
      case 'labor': return '#a35a2a';
      case 'equipment': return '#3d4450';
      default: {
        const ct = data?.custom_types.find(t => t.slug === type);
        return ct?.color || '#3d4450';
      }
    }
  }

  function getTypeBgClass(type: string): string {
    switch (type) {
      case 'material': return 'bg-forest-50';
      case 'labor': return 'bg-copper-50';
      case 'equipment': return 'bg-slate-50';
      default: return 'bg-slate-50';
    }
  }

  function formatMoney(n: number): string {
    return '$' + n.toFixed(2).replace(/\B(?=(\d{3})+(?!\d))/g, ',');
  }

  function formatNumber(n: number): string {
    // Remove trailing zeros after decimal
    const s = n.toString();
    return s;
  }

  function getCellValue(item: QuoteLineItem, col: EditableColumn): string {
    switch (col) {
      case 'name': return item.name;
      case 'quantity': return formatNumber(item.quantity);
      case 'unit': return item.unit;
      case 'unit_price': return formatNumber(item.unit_price);
      case 'surcharge_percent':
        return item.surcharge_percent != null ? formatNumber(item.surcharge_percent) : '';
    }
  }

  function getCellDisplayValue(item: QuoteLineItem, col: EditableColumn): string {
    switch (col) {
      case 'name': return item.name;
      case 'quantity': return formatNumber(item.quantity);
      case 'unit': return item.unit;
      case 'unit_price': return formatMoney(item.unit_price);
      case 'surcharge_percent':
        return item.surcharge_percent != null ? item.surcharge_percent + '%' : 'inherit';
    }
  }

  function getItemTotal(item: QuoteLineItem): number {
    return item.quantity * item.unit_price;
  }

  // --- Data mutation helpers ---
  function findAndUpdateItem(categories: QuoteCategory[], itemId: string, updater: (item: QuoteLineItem) => QuoteLineItem): boolean {
    for (const cat of categories) {
      const idx = cat.items.findIndex(i => i.id === itemId);
      if (idx !== -1) {
        cat.items[idx] = updater(cat.items[idx]);
        return true;
      }
      if (findAndUpdateItem(cat.children, itemId, updater)) return true;
    }
    return false;
  }

  function findAndRemoveItem(categories: QuoteCategory[], itemId: string): QuoteLineItem | null {
    for (const cat of categories) {
      const idx = cat.items.findIndex(i => i.id === itemId);
      if (idx !== -1) {
        return cat.items.splice(idx, 1)[0];
      }
      const found = findAndRemoveItem(cat.children, itemId);
      if (found) return found;
    }
    return null;
  }

  function findCategory(categories: QuoteCategory[], catId: string): QuoteCategory | null {
    for (const cat of categories) {
      if (cat.id === catId) return cat;
      const found = findCategory(cat.children, catId);
      if (found) return found;
    }
    return null;
  }

  // --- Item search modal ---
  function openSearch(type: string, categoryId: string) {
    searchType = type;
    searchCategoryId = categoryId;
    searchQuery = '';
    searchResults = [];
    searchSelectedIdx = 0;
    searchOpen = true;
    searchLoading = false;
  }

  function closeSearch() {
    searchOpen = false;
    if (searchDebounceTimer) clearTimeout(searchDebounceTimer);
    // Refocus the grid
    tick().then(() => gridEl?.focus());
  }

  function handleSearchInput(e: Event) {
    const value = (e.target as HTMLInputElement).value;
    searchQuery = value;
    searchSelectedIdx = 0;

    if (searchDebounceTimer) clearTimeout(searchDebounceTimer);
    if (value.length === 0) {
      searchResults = [];
      searchLoading = false;
      return;
    }
    searchLoading = true;
    searchDebounceTimer = setTimeout(async () => {
      try {
        searchResults = await searchTemplates(value, searchType);
      } catch (err) {
        console.error('Search failed:', err);
        searchResults = [];
      } finally {
        searchLoading = false;
      }
    }, 200);
  }

  async function pickTemplate(template: TemplateResult | null) {
    closeSearch();
    if (!data) return;

    const name = template ? template.name : 'New Item';
    const unit = template ? template.default_unit : (searchType === 'labor' ? 'hr' : 'ea');
    const unitPrice = template ? template.default_price : 0;

    try {
      const resp = await createItem(jobId, {
        category_id: searchCategoryId,
        type: searchType,
        name,
        quantity: 1,
        unit: unit,
        unit_price: unitPrice,
      });

      const cat = findCategory(data.categories, searchCategoryId);
      if (cat) {
        cat.items = [...cat.items, resp.item];
        data.totals = resp.totals;

        await tick();
        const newIdx = gridRows.findIndex(r => r.item?.id === resp.item.id);
        if (newIdx >= 0) {
          selectedRow = newIdx;
          // Template picked: jump to quantity and start editing
          selectedCol = 1;
          startEditing();
        }
      }
    } catch (e) {
      console.error('Failed to create item:', e);
    }
  }

  async function pickSearchAsNewItem() {
    const name = searchQuery.trim() || 'New Item';
    const editName = !searchQuery.trim();
    closeSearch();
    if (!data) return;

    try {
      const resp = await createItem(jobId, {
        category_id: searchCategoryId,
        type: searchType,
        name,
        quantity: 1,
        unit: searchType === 'labor' ? 'hr' : 'ea',
        unit_price: 0,
      });

      const cat = findCategory(data.categories, searchCategoryId);
      if (cat) {
        cat.items = [...cat.items, resp.item];
        data.totals = resp.totals;

        await tick();
        const newIdx = gridRows.findIndex(r => r.item?.id === resp.item.id);
        if (newIdx >= 0) {
          selectedRow = newIdx;
          // No name: edit name. Has name: edit quantity.
          selectedCol = editName ? 0 : 1;
          startEditing();
        }
      }
    } catch (e) {
      console.error('Failed to create item:', e);
    }
  }

  function scrollSearchResult() {
    tick().then(() => {
      if (!searchResultsEl) return;
      const item = searchResultsEl.children[searchSelectedIdx] as HTMLElement | undefined;
      item?.scrollIntoView({ block: 'nearest' });
    });
  }

  function handleSearchKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      e.preventDefault();
      closeSearch();
    } else if (e.key === 'ArrowDown') {
      e.preventDefault();
      if (searchResults.length > 0) {
        searchSelectedIdx = Math.min(searchSelectedIdx + 1, searchResults.length - 1);
        scrollSearchResult();
      }
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      searchSelectedIdx = Math.max(searchSelectedIdx - 1, 0);
      scrollSearchResult();
    } else if (e.key === 'Enter') {
      e.preventDefault();
      if (searchResults.length > 0 && searchSelectedIdx < searchResults.length) {
        pickTemplate(searchResults[searchSelectedIdx]);
      } else {
        // No results - create item with typed name (or blank)
        pickSearchAsNewItem();
      }
    }
  }

  // --- Actions ---
  async function commitEdit() {
    if (!editing || !data) return;
    const row = gridRows[selectedRow];
    if (!row?.item) return;

    const col = COLUMNS[selectedCol];
    const item = row.item;
    const oldValue = getCellValue(item, col);

    if (editValue === oldValue) {
      editing = false;
      tick().then(() => gridEl?.focus());
      return;
    }

    // Build patch
    const patch: Record<string, any> = {};
    if (col === 'name') {
      patch.name = editValue;
    } else if (col === 'quantity') {
      const v = parseFloat(editValue);
      if (isNaN(v) || v <= 0) { editing = false; tick().then(() => gridEl?.focus()); return; }
      patch.quantity = v;
    } else if (col === 'unit') {
      patch.unit = editValue;
    } else if (col === 'unit_price') {
      const v = parseFloat(editValue);
      if (isNaN(v)) { editing = false; tick().then(() => gridEl?.focus()); return; }
      patch.unit_price = v;
    } else if (col === 'surcharge_percent') {
      if (editValue === '' || editValue === 'inherit') {
        patch.surcharge_percent = null;
      } else {
        const v = parseFloat(editValue);
        if (isNaN(v)) { editing = false; tick().then(() => gridEl?.focus()); return; }
        patch.surcharge_percent = v;
      }
    }

    // Optimistic update
    findAndUpdateItem(data.categories, item.id, (i) => ({ ...i, ...patch }));
    editing = false;
    tick().then(() => gridEl?.focus());

    try {
      const resp = await patchItem(jobId, item.id, patch);
      // Update with server-confirmed values
      findAndUpdateItem(data.categories, item.id, () => resp.item);
      data.totals = resp.totals;
    } catch (e) {
      // Rollback on failure
      findAndUpdateItem(data.categories, item.id, () => item);
      console.error('Failed to save:', e);
    }
  }

  async function addItem(categoryId: string, type: string) {
    if (!data) return;

    try {
      const resp = await createItem(jobId, {
        category_id: categoryId,
        type,
        name: 'New Item',
        quantity: 1,
        unit: type === 'labor' ? 'hr' : 'ea',
        unit_price: 0,
      });

      // Insert the new item into the category
      const cat = findCategory(data.categories, categoryId);
      if (cat) {
        cat.items = [...cat.items, resp.item];
        data.totals = resp.totals;

        // Select the new item's name cell
        await tick();
        const newIdx = gridRows.findIndex(r => r.item?.id === resp.item.id);
        if (newIdx >= 0) {
          selectedRow = newIdx;
          selectedCol = 0; // name column
          startEditing();
        }
      }
    } catch (e) {
      console.error('Failed to create item:', e);
    }
  }

  async function removeItem(itemId: string) {
    if (!data) return;
    if (!confirm('Delete this item?')) return;

    const removed = findAndRemoveItem(data.categories, itemId);
    if (!removed) return;

    try {
      const resp = await deleteItem(jobId, itemId);
      data.totals = resp.totals;
    } catch (e) {
      // Rollback
      const cat = findCategory(data.categories, removed.category_id);
      if (cat) cat.items = [...cat.items, removed];
      console.error('Failed to delete item:', e);
    }
  }

  function startEditing() {
    const row = gridRows[selectedRow];
    if (row?.kind !== 'item' || !row.item) return;
    editValue = getCellValue(row.item, COLUMNS[selectedCol]);
    editing = true;
  }

  function cancelEdit() {
    editing = false;
    tick().then(() => gridEl?.focus());
  }

  // Svelte tick for waiting for DOM updates
  function tick(): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, 0));
  }

  // --- Keyboard handler ---
  function handleKeydown(e: KeyboardEvent) {
    // Don't handle if search modal is open or typing in an input outside our grid
    if (searchOpen) return;
    const target = e.target as HTMLElement;
    if (target.tagName === 'SELECT' || target.tagName === 'TEXTAREA') return;

    if (editing) {
      if (e.key === 'Escape') {
        e.preventDefault();
        cancelEdit();
      } else if (e.key === 'Enter') {
        e.preventDefault();
        commitEdit();
      } else if (e.key === 'Tab') {
        e.preventDefault();
        commitEdit().then(() => {
          // Tab = next row, same column
          if (e.shiftKey) {
            moveToNextItemRow(-1);
          } else {
            moveToNextItemRow(1);
          }
          startEditing();
        });
      }
      return;
    }

    // Navigation mode
    switch (e.key) {
      case 'ArrowDown':
      case 'j':
        e.preventDefault();
        moveToNextItemRow(1);
        break;
      case 'ArrowUp':
      case 'k':
        e.preventDefault();
        moveToNextItemRow(-1);
        break;
      case 'ArrowRight':
      case 'l':
        // 'l' adds labor item when not in edit mode; only ArrowRight navigates columns
        if (e.key === 'ArrowRight') {
          e.preventDefault();
          selectedCol = Math.min(selectedCol + 1, COLUMNS.length - 1);
        }
        break;
      case 'ArrowLeft':
      case 'h':
        if (e.key === 'ArrowLeft') {
          e.preventDefault();
          selectedCol = Math.max(selectedCol - 1, 0);
        }
        break;
      case 'Tab':
        e.preventDefault();
        if (e.shiftKey) {
          moveToNextItemRow(-1);
        } else {
          moveToNextItemRow(1);
        }
        break;
      case 'Enter':
        e.preventDefault();
        startEditing();
        break;
      case 'm': {
        e.preventDefault();
        const catId = gridRows[selectedRow]?.categoryId;
        if (catId) openSearch('material', catId);
        break;
      }
      case 'l': {
        e.preventDefault();
        const catId = gridRows[selectedRow]?.categoryId;
        if (catId) openSearch('labor', catId);
        break;
      }
      case 'e': {
        e.preventDefault();
        const catId = gridRows[selectedRow]?.categoryId;
        if (catId) openSearch('equipment', catId);
        break;
      }
      case 'd':
      case 'Delete': {
        e.preventDefault();
        const row = gridRows[selectedRow];
        if (row?.item) removeItem(row.item.id);
        break;
      }
      case 'c': {
        e.preventDefault();
        // Top-level category
        const container = document.getElementById('category-form-container');
        if (container) {
          container.removeAttribute('data-parent-id');
          (window as any).showCategoryForm?.();
        }
        break;
      }
      case 'C': {
        // Subcategory under selected category
        e.preventDefault();
        const row = gridRows[selectedRow];
        const parentId = row?.categoryId;
        if (parentId) {
          const container = document.getElementById('category-form-container');
          if (container) {
            container.setAttribute('data-parent-id', parentId);
            (window as any).showCategoryForm?.();
          }
        }
        break;
      }
      default:
        // If a printable character is pressed and we're on an item, start editing
        if (e.key.length === 1 && !e.ctrlKey && !e.metaKey && !e.altKey) {
          const row = gridRows[selectedRow];
          if (row?.kind === 'item') {
            e.preventDefault();
            editValue = e.key;
            editing = true;
          }
        }
    }
  }

  function moveToNextItemRow(direction: 1 | -1) {
    if (gridRows.length === 0) return;

    // If current row is a category or out of bounds, jump to first/last item
    if (selectedRow < 0 || selectedRow >= gridRows.length || gridRows[selectedRow].kind === 'category') {
      if (direction === 1) {
        const first = gridRows.findIndex(r => r.kind === 'item');
        if (first >= 0) selectedRow = first;
      } else {
        for (let i = gridRows.length - 1; i >= 0; i--) {
          if (gridRows[i].kind === 'item') { selectedRow = i; break; }
        }
      }
      return;
    }

    let next = selectedRow + direction;
    // Skip category header rows
    while (next >= 0 && next < gridRows.length && gridRows[next].kind === 'category') {
      next += direction;
    }
    if (next >= 0 && next < gridRows.length) {
      selectedRow = next;
    }
  }

  // --- Lifecycle ---
  onMount(async () => {
    try {
      data = await fetchJob(jobId);
    } catch (e: any) {
      error = e.message || 'Failed to load job';
    } finally {
      loading = false;
    }
  });
</script>

<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
  bind:this={gridEl}
  class="quote-editor"
  role="grid"
  tabindex="0"
  onkeydown={handleKeydown}
>
  {#if loading}
    <div class="px-4 py-8 text-center text-slate-500">Loading...</div>
  {:else if error}
    <div class="px-4 py-8 text-center text-red-600">{error}</div>
  {:else if data}
    <!-- Grid Header -->
    <div class="hidden sm:flex px-3 py-2 bg-slate-50 border-b border-slate-200 text-xs font-medium tracking-wider uppercase text-slate-500 gap-1">
      {#each COL_HEADERS as header, i}
        <div class="{COL_WIDTHS[i]} px-1 {i > 0 ? 'text-right' : ''}">
          {header}
        </div>
      {/each}
      <div class="w-20 text-right px-1">Total</div>
      <div class="w-8"></div>
    </div>

    <!-- Grid Body -->
    <div class="divide-y divide-slate-100">
      {#each gridRows as row, rowIdx}
        {#if row.kind === 'category'}
          <!-- Category Header Row -->
          <div
            class="flex items-center px-3 py-2 bg-slate-50/80 border-l-4 border-slate-300"
            style="padding-left: {0.75 + row.depth * 1}rem"
          >
            <span class="font-semibold text-sm text-slate-700">{row.categoryName}</span>
          </div>
        {:else if row.item}
          <!-- Item Row -->
          <!-- svelte-ignore a11y_click_events_have_key_events -->
          <div
            class="flex items-center px-3 py-1.5 gap-1 cursor-pointer transition-colors text-sm"
            class:bg-copper-50={rowIdx === selectedRow}
            class:hover:bg-slate-50={rowIdx !== selectedRow}
            style="padding-left: {0.75 + row.depth * 1}rem"
            onclick={() => { selectedRow = rowIdx; }}
          >
            {#each COLUMNS as col, colIdx}
              <div
                class="{COL_WIDTHS[colIdx]} px-1 {colIdx > 0 ? 'text-right' : ''} truncate"
              >
                {#if editing && rowIdx === selectedRow && colIdx === selectedCol}
                  <input
                    type="text"
                    class="w-full bg-white border border-copper-400 rounded px-1 py-0.5 text-sm outline-none {colIdx > 0 ? 'text-right' : ''}"
                    bind:value={editValue}
                    onfocus={(e) => (e.target as HTMLInputElement).select()}
                    onblur={() => commitEdit()}
                    use:autofocus
                  />
                {:else}
                  <!-- svelte-ignore a11y_click_events_have_key_events -->
                  <!-- svelte-ignore a11y_no_static_element_interactions -->
                  <span
                    class="block w-full truncate rounded px-1 py-0.5 {colIdx === selectedCol && rowIdx === selectedRow ? 'ring-2 ring-copper-400 ring-inset' : ''}"
                    class:text-slate-400={col === 'surcharge_percent' && row.item.surcharge_percent == null}
                    ondblclick={() => { selectedRow = rowIdx; selectedCol = colIdx; startEditing(); }}
                    onclick={() => { selectedCol = colIdx; }}
                  >
                    {getCellDisplayValue(row.item, col)}
                  </span>
                {/if}
              </div>
            {/each}

            <!-- Base total (qty * unit_price) -->
            <div class="w-20 text-right px-1 tabular-nums font-medium" style="color: {getTypeColor(row.item.type)}">
              {formatMoney(getItemTotal(row.item))}
            </div>

            <!-- Delete button -->
            <div class="w-8 flex justify-center">
              {#if rowIdx === selectedRow}
                <button
                  class="text-slate-400 hover:text-red-500 p-0.5"
                  title="Delete item"
                  onclick={(e) => { e.stopPropagation(); if (row.item) removeItem(row.item.id); }}
                >
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                  </svg>
                </button>
              {/if}
            </div>
          </div>
        {/if}
      {/each}
    </div>

    {#if gridRows.length === 0}
      <div class="px-4 py-8 text-center text-slate-500">
        <p>No categories yet.</p>
        <p class="text-sm mt-2">Press <kbd class="font-mono text-xs px-1.5 py-0.5 bg-slate-100 border border-slate-300 rounded text-slate-700">c</kbd> to create a category.</p>
      </div>
    {/if}

    <!-- Totals Summary -->
    <div class="mt-4 bg-white rounded-lg border border-slate-200 p-4">
      <div class="grid grid-cols-3 gap-4 text-sm">
        <div>
          <span class="text-slate-500">Materials</span>
          <p class="tabular-nums font-medium" style="color: #2d5a47">{formatMoney(liveTotals.type_subtotals['material'] || 0)}</p>
        </div>
        <div>
          <span class="text-slate-500">Labor</span>
          <p class="tabular-nums font-medium" style="color: #a35a2a">{formatMoney(liveTotals.type_subtotals['labor'] || 0)}</p>
        </div>
        <div>
          <span class="text-slate-500">Equipment</span>
          <p class="tabular-nums font-medium" style="color: #3d4450">{formatMoney(liveTotals.type_subtotals['equipment'] || 0)}</p>
        </div>
      </div>
      {#each data.custom_types as ct}
        <div class="mt-2">
          <span class="text-slate-500 text-sm">{ct.name}</span>
          <p class="tabular-nums font-medium text-sm" style="color: {ct.color}">{formatMoney(liveTotals.type_subtotals[ct.slug] || 0)}</p>
        </div>
      {/each}
      <div class="mt-3 pt-3 border-t border-slate-100 flex justify-between items-center">
        <span class="text-sm font-medium text-slate-700">Grand Total</span>
        <span class="text-xl font-semibold tabular-nums text-slate-900">{formatMoney(liveTotals.grand_total)}</span>
      </div>
    </div>

  {/if}

  <!-- Item Search Modal -->
  {#if searchOpen}
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <div class="fixed inset-0 z-50 flex items-start justify-center pt-[20vh]" onclick={() => closeSearch()}>
      <div class="absolute inset-0 bg-black/30"></div>
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <div
        class="relative bg-white rounded-lg shadow-xl border border-slate-200 w-full max-w-md mx-4"
        onclick={(e) => e.stopPropagation()}
        onkeydown={handleSearchKeydown}
      >
        <div class="px-4 pt-4 pb-2">
          <div class="flex items-center gap-2 mb-3">
            <span
              class="text-xs font-medium px-2 py-0.5 rounded-full text-white"
              style="background-color: {getTypeColor(searchType)}"
            >
              {searchType}
            </span>
            <span class="text-sm text-slate-500">Search or press Enter to create blank</span>
          </div>
          <input
            type="text"
            class="w-full border border-slate-300 rounded-md px-3 py-2 text-sm outline-none focus:border-copper-400 focus:ring-1 focus:ring-copper-400"
            placeholder="Type to search items..."
            value={searchQuery}
            oninput={handleSearchInput}
            use:autofocus
          />
        </div>
        <div class="max-h-64 overflow-y-auto" bind:this={searchResultsEl}>
          {#if searchLoading}
            <div class="px-4 py-3 text-sm text-slate-400">Searching...</div>
          {:else if searchQuery && searchResults.length === 0}
            <div class="px-4 py-3 text-sm text-slate-400">No templates found. Press Enter to create "{searchQuery}" as a new item.</div>
          {:else}
            {#each searchResults as result, idx}
              <!-- svelte-ignore a11y_click_events_have_key_events -->
              <!-- svelte-ignore a11y_no_static_element_interactions -->
              <div
                class="flex items-center justify-between px-4 py-2 cursor-pointer text-sm {idx === searchSelectedIdx ? 'bg-copper-50' : 'hover:bg-slate-50'}"
                onclick={() => pickTemplate(result)}
                onmouseenter={() => { searchSelectedIdx = idx; }}
              >
                <div class="flex-1 min-w-0">
                  <div class="font-medium text-slate-800 truncate">{result.name}</div>
                  <div class="text-xs text-slate-400">{result.category} &middot; {result.default_unit}</div>
                </div>
                <div class="text-sm tabular-nums text-slate-600 ml-3">
                  {formatMoney(result.default_price)}
                </div>
              </div>
            {/each}
          {/if}
        </div>
        <div class="px-4 py-2 border-t border-slate-100 text-xs text-slate-400 flex gap-3">
          <span><kbd class="font-mono px-1 py-0.5 bg-slate-100 border border-slate-200 rounded">↑↓</kbd> navigate</span>
          <span><kbd class="font-mono px-1 py-0.5 bg-slate-100 border border-slate-200 rounded">Enter</kbd> select</span>
          <span><kbd class="font-mono px-1 py-0.5 bg-slate-100 border border-slate-200 rounded">Esc</kbd> cancel</span>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .quote-editor:focus {
    outline: none;
  }
</style>
