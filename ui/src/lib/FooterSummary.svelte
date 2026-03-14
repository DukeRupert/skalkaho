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

<div class="fixed bottom-0 right-0 border-t-2 z-20" style="left: var(--sidebar-width, 220px); background: var(--color-granite); border-color: var(--color-sunburst);">
	<div class="flex items-center justify-between px-5 py-3">
		<div class="flex items-center gap-6 text-sm">
			<!-- Base cost -->
			<div class="flex flex-col">
				<span class="text-xs text-white/30 uppercase tracking-wider font-[var(--font-ui)]">Base Cost</span>
				<span class="font-mono text-[var(--color-concrete)]">{formatMoney(totals.base)}</span>
			</div>

			<!-- Divider -->
			<div class="w-px h-8 bg-white/[0.08]"></div>

			<!-- Per-type totals (only non-zero) -->
			{#each activeTypes as { type, value }}
				<div class="flex flex-col">
					<span class="text-xs text-white/30 uppercase tracking-wider font-[var(--font-ui)]">{TYPE_LABELS[type]}</span>
					<span class="font-mono text-[var(--color-concrete)]">{formatMoney(value)}</span>
				</div>
			{/each}

			{#if activeTypes.length === 0}
				<span class="text-[var(--color-muted-text)] text-xs font-[var(--font-body)]">No items yet</span>
			{/if}
		</div>

		<!-- Grand total -->
		<div class="flex flex-col items-end">
			<span class="text-xs text-white/30 uppercase tracking-wider font-[var(--font-ui)]">Total</span>
			<span class="font-mono text-lg font-bold text-[var(--color-sunburst)]">{formatMoney(totals.withMarkup)}</span>
		</div>
	</div>
</div>
