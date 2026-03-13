<script>
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
	<div class="estimate-builder">
		<!-- Topbar with project info and save status -->
		<div class="sticky top-0 z-10 bg-slate-900 text-white px-4 py-3 flex items-center justify-between border-b border-slate-700">
			<div>
				<h1 class="text-lg font-semibold">{estimate.project.name}</h1>
				<span class="text-xs text-slate-400">Estimate Builder</span>
			</div>
			<div class="flex items-center gap-4">
				<span class="text-xs px-2 py-1 rounded bg-slate-700 text-slate-300 uppercase tracking-wide">
					{estimate.project.status}
				</span>
			</div>
		</div>

		<!-- Global markup toolbar -->
		<div class="bg-slate-50 border-b border-slate-200 px-4 py-2">
			<div class="flex items-center gap-4 text-sm">
				<span class="font-medium text-slate-600">Global Markup:</span>
				<span class="text-slate-500">Materials {estimate.globals.materials_markup}%</span>
				<span class="text-slate-500">Labor {estimate.globals.labor_markup}%</span>
				<span class="text-slate-500">Equipment {estimate.globals.equipment_markup}%</span>
				<span class="text-slate-500">Subs {estimate.globals.subs_markup}%</span>
				<span class="text-slate-500">Other {estimate.globals.other_markup}%</span>
			</div>
		</div>

		<!-- Sections -->
		<div class="p-4">
			{#if estimate.sections.length === 0}
				<div class="text-center py-16 text-slate-400">
					<p class="text-lg">No sections yet</p>
					<p class="text-sm mt-1">Sections, subcategories, and line items will appear here.</p>
				</div>
			{:else}
				{#each estimate.sections as section}
					<div class="mb-4 border border-slate-200 rounded-lg overflow-hidden">
						<div class="bg-slate-100 px-4 py-2 font-semibold text-slate-700">
							{section.name}
						</div>
						{#each section.subcategories as subcat}
							<div class="border-t border-slate-200 px-4 py-2">
								<div class="font-medium text-slate-600 text-sm">{subcat.name}</div>
								{#if subcat.line_items.length > 0}
									<div class="mt-1 text-xs text-slate-400">
										{subcat.line_items.length} ungrouped item{subcat.line_items.length !== 1 ? 's' : ''}
									</div>
								{/if}
								{#each subcat.component_groups as group}
									<div class="ml-4 mt-1 text-xs text-slate-500">
										{group.name} ({group.line_items.length} item{group.line_items.length !== 1 ? 's' : ''})
									</div>
								{/each}
							</div>
						{/each}
					</div>
				{/each}
			{/if}
		</div>

		<!-- Debug: payload summary -->
		<div class="px-4 py-2 bg-slate-50 border-t border-slate-200 text-xs text-slate-400">
			{estimate.sections.length} sections |
			{estimate.materials_db.length} materials |
			{estimate.rates_db.length} rates
		</div>
	</div>
{/if}
