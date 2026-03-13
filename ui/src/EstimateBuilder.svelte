<script>
	import SectionBlock from './lib/SectionBlock.svelte';
	import FooterSummary from './lib/FooterSummary.svelte';
	import { formatPercent } from './lib/markup.js';

	let { projectId } = $props();

	let estimate = $state(null);
	let error = $state(null);
	let loading = $state(true);

	async function fetchEstimate() {
		try {
			const res = await fetch(`/api/estimate/${projectId}`);
			if (!res.ok) {
				throw new Error(`Failed to load estimate: ${res.status}`);
			}
			estimate = await res.json();
		} catch (e) {
			error = e.message;
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (projectId) {
			fetchEstimate();
		}
	});
</script>

{#if loading}
	<div class="flex items-center justify-center h-64">
		<div class="text-slate-500">Loading estimate...</div>
	</div>
{:else if error}
	<div class="flex items-center justify-center h-64">
		<div class="text-red-500">{error}</div>
	</div>
{:else if estimate}
	<div class="estimate-builder pb-16">
		<!-- Topbar -->
		<div class="sticky top-0 z-10 bg-slate-900 text-white px-4 py-3 flex items-center justify-between border-b border-slate-700">
			<div>
				<h1 class="text-lg font-semibold">{estimate.project.name}</h1>
				<span class="text-xs text-slate-400">Estimate Builder</span>
			</div>
			<div class="flex items-center gap-3">
				<span class="text-xs px-2 py-1 rounded bg-slate-700 text-slate-300 uppercase tracking-wide">
					{estimate.project.status}
				</span>
			</div>
		</div>

		<!-- Global markup toolbar -->
		<div class="bg-slate-50 border-b border-slate-200 px-4 py-2">
			<div class="flex items-center gap-4 text-sm">
				<span class="font-medium text-slate-600">Global Markup:</span>
				<span class="px-2 py-0.5 rounded bg-blue-50 text-blue-700 font-mono text-xs">
					Materials {formatPercent(estimate.globals.materials_markup)}
				</span>
				<span class="px-2 py-0.5 rounded bg-amber-50 text-amber-700 font-mono text-xs">
					Labor {formatPercent(estimate.globals.labor_markup)}
				</span>
				<span class="px-2 py-0.5 rounded bg-purple-50 text-purple-700 font-mono text-xs">
					Equipment {formatPercent(estimate.globals.equipment_markup)}
				</span>
				<span class="px-2 py-0.5 rounded bg-green-50 text-green-700 font-mono text-xs">
					Subs {formatPercent(estimate.globals.subs_markup)}
				</span>
				<span class="px-2 py-0.5 rounded bg-slate-100 text-slate-600 font-mono text-xs">
					Other {formatPercent(estimate.globals.other_markup)}
				</span>
			</div>
		</div>

		<!-- Sections -->
		<div class="p-4">
			{#if estimate.sections.length === 0}
				<div class="text-center py-16 text-slate-400">
					<svg class="w-12 h-12 mx-auto mb-3 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
					</svg>
					<p class="text-lg font-medium">No sections yet</p>
					<p class="text-sm mt-1">Add a section to start building your estimate.</p>
				</div>
			{:else}
				{#each estimate.sections as section (section.id)}
					<SectionBlock {section} globals={estimate.globals} />
				{/each}
			{/if}
		</div>

		<!-- Footer summary bar -->
		<FooterSummary {estimate} />
	</div>
{/if}
