<script>
	import SubcategoryBlock from './SubcategoryBlock.svelte';
	import { sectionTotals, formatMoney } from './markup.js';

	let { section, globals, collapsed = false, onchange, materialsDb = [], ratesDb = [] } = $props();

	let isCollapsed = $state(collapsed);
	let totals = $derived(sectionTotals(section, globals));

	let totalItems = $derived(
		section.subcategories.reduce((sum, sc) =>
			sum + sc.line_items.length +
			sc.component_groups.reduce((gs, g) => gs + g.line_items.length, 0),
		0)
	);
</script>

<div class="mb-4 border border-slate-200 rounded-lg overflow-hidden bg-white">
	<button
		class="w-full flex items-center justify-between px-4 py-3 bg-slate-800 text-white hover:bg-slate-700 transition-colors text-left"
		onclick={() => isCollapsed = !isCollapsed}
	>
		<div class="flex items-center gap-3">
			<svg class="w-4 h-4 text-slate-400 transition-transform {isCollapsed ? '' : 'rotate-90'}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
			</svg>
			<span class="font-semibold">{section.name}</span>
			<span class="text-xs text-slate-400">
				{section.subcategories.length} subcategor{section.subcategories.length !== 1 ? 'ies' : 'y'}
				&middot; {totalItems} item{totalItems !== 1 ? 's' : ''}
			</span>
		</div>
		<span class="font-mono font-semibold">{formatMoney(totals.withMarkup)}</span>
	</button>

	{#if !isCollapsed}
		{#if section.subcategories.length === 0}
			<div class="px-4 py-8 text-center text-slate-400 text-sm">
				No subcategories yet
			</div>
		{:else}
			{#each section.subcategories as subcat (subcat.id)}
				<SubcategoryBlock {subcat} {globals} {onchange} {materialsDb} {ratesDb} />
			{/each}
		{/if}
	{/if}
</div>
