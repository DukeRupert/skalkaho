<script>
	import LineItemRow from './LineItemRow.svelte';
	import ComponentGroupBlock from './ComponentGroupBlock.svelte';
	import { subcategoryTotals, formatMoney, resolveMarkup, formatPercent } from './markup.js';

	let { subcat, globals, collapsed = false } = $props();

	let isCollapsed = $state(collapsed);
	let totals = $derived(subcategoryTotals(subcat, globals));

	// Build markup summary for display
	let markupSummary = $derived.by(() => {
		const types = ['materials', 'labor', 'equipment', 'subs', 'other'];
		return types.map(t => ({
			type: t,
			value: resolveMarkup(t, globals, subcat.markup_overrides, subcat.markup_enabled),
			isOverride: subcat.markup_overrides[t] != null,
			isDisabled: !subcat.markup_enabled[t],
		}));
	});

	let hasOverrides = $derived(markupSummary.some(m => m.isOverride || m.isDisabled));
	let totalItems = $derived(
		subcat.line_items.length +
		subcat.component_groups.reduce((sum, g) => sum + g.line_items.length, 0)
	);
</script>

<div class="border-t border-slate-200">
	<!-- Subcategory header -->
	<button
		class="w-full flex items-center justify-between px-4 py-2 hover:bg-slate-50 transition-colors text-left"
		onclick={() => isCollapsed = !isCollapsed}
	>
		<div class="flex items-center gap-2">
			<svg class="w-4 h-4 text-slate-400 transition-transform {isCollapsed ? '' : 'rotate-90'}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
			</svg>
			<span class="font-medium text-slate-600 text-sm">{subcat.name}</span>
			<span class="text-xs text-slate-400">{totalItems} item{totalItems !== 1 ? 's' : ''}</span>
			{#if hasOverrides}
				<span class="text-xs px-1.5 py-0.5 rounded bg-amber-50 text-amber-600 font-medium">overrides</span>
			{/if}
			{#if subcat.lump_sum > 0}
				<span class="text-xs px-1.5 py-0.5 rounded bg-green-50 text-green-600 font-medium">
					+{formatMoney(subcat.lump_sum)} lump sum
				</span>
			{/if}
		</div>
		<span class="font-mono text-sm font-semibold text-slate-700">{formatMoney(totals.withMarkup)}</span>
	</button>

	<!-- Subcategory content -->
	{#if !isCollapsed}
		<div class="px-4 pb-3">
			<!-- Markup row (if overrides) -->
			{#if hasOverrides}
				<div class="flex items-center gap-3 py-1.5 px-2 bg-amber-50 rounded text-xs mb-2">
					<span class="font-medium text-amber-700">Markup:</span>
					{#each markupSummary as m}
						<span class={m.isDisabled ? 'text-slate-400 line-through' : m.isOverride ? 'text-amber-700 font-medium' : 'text-slate-500'}>
							{m.type} {formatPercent(m.value)}
						</span>
					{/each}
				</div>
			{/if}

			<!-- Line items table header -->
			{#if totalItems > 0}
				<table class="w-full">
					<thead>
						<tr class="text-xs text-slate-400 uppercase tracking-wide border-b border-slate-100">
							<th class="px-2 py-1 text-left w-20">Type</th>
							<th class="px-2 py-1 text-left">Name</th>
							<th class="px-2 py-1 text-right w-20">Qty</th>
							<th class="px-2 py-1 text-center w-16">Unit</th>
							<th class="px-2 py-1 text-right w-24">Price</th>
							<th class="px-2 py-1 text-right w-16">Markup</th>
							<th class="px-2 py-1 text-right w-24">w/ Markup</th>
							<th class="px-2 py-1 text-right w-28">Total</th>
						</tr>
					</thead>
					<tbody>
						<!-- Ungrouped items -->
						{#each subcat.line_items as item (item.id)}
							<LineItemRow
								{item}
								{globals}
								markupOverrides={subcat.markup_overrides}
								markupEnabled={subcat.markup_enabled}
							/>
						{/each}
					</tbody>
				</table>
			{/if}

			<!-- Component groups -->
			{#each subcat.component_groups as group (group.id)}
				<ComponentGroupBlock
					{group}
					{globals}
					markupOverrides={subcat.markup_overrides}
					markupEnabled={subcat.markup_enabled}
				/>
			{/each}

			<!-- Subcategory total row -->
			<div class="flex justify-between items-center mt-2 pt-2 border-t border-slate-100 text-sm">
				<span class="text-slate-500">
					Subtotal: {formatMoney(totals.base)}
					{#if subcat.lump_sum > 0}
						<span class="text-green-600"> + {formatMoney(subcat.lump_sum)} lump sum</span>
					{/if}
				</span>
				<span class="font-mono font-semibold text-slate-700">{formatMoney(totals.withMarkup)}</span>
			</div>
		</div>
	{/if}
</div>
