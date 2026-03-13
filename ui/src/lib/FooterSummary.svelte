<script>
	import { estimateTotals, formatMoney } from './markup.js';

	let { estimate } = $props();

	let totals = $derived(estimateTotals(estimate));

	// Only show type columns that have non-zero values
	let activeTypes = $derived(
		Object.entries(totals.byType)
			.filter(([, value]) => value > 0)
			.map(([type, value]) => ({ type, value }))
	);

	const TYPE_LABELS = {
		materials: 'Materials',
		labor: 'Labor',
		equipment: 'Equipment',
		subs: 'Subs',
		other: 'Other',
	};
</script>

<div class="fixed bottom-0 left-0 right-0 bg-slate-900 text-white border-t border-slate-700 z-20">
	<div class="flex items-center justify-between px-4 py-3">
		<div class="flex items-center gap-6 text-sm">
			<!-- Base cost -->
			<div class="flex flex-col">
				<span class="text-xs text-slate-400 uppercase tracking-wide">Base Cost</span>
				<span class="font-mono">{formatMoney(totals.base)}</span>
			</div>

			<!-- Divider -->
			<div class="w-px h-8 bg-slate-700"></div>

			<!-- Per-type totals (only non-zero) -->
			{#each activeTypes as { type, value }}
				<div class="flex flex-col">
					<span class="text-xs text-slate-400 uppercase tracking-wide">{TYPE_LABELS[type]}</span>
					<span class="font-mono">{formatMoney(value)}</span>
				</div>
			{/each}

			{#if activeTypes.length === 0}
				<span class="text-slate-500 text-xs">No items yet</span>
			{/if}
		</div>

		<!-- Grand total -->
		<div class="flex flex-col items-end">
			<span class="text-xs text-slate-400 uppercase tracking-wide">Total</span>
			<span class="font-mono text-lg font-bold">{formatMoney(totals.withMarkup)}</span>
		</div>
	</div>
</div>
