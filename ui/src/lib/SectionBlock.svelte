<script>
	import SubcategoryBlock from './SubcategoryBlock.svelte';
	import { sectionTotals, formatMoney } from './markup.js';
	import { nanoid } from 'nanoid';

	let { section, globals, collapsed = false, onchange, onsnapshot, ondelete, materialsDb = [], ratesDb = [] } = $props();

	let isCollapsed = $state(collapsed);
	let isEditing = $state(false);
	let editName = $state(section.name);
	let showAddSubcat = $state(false);
	let newSubcatName = $state('');
	let totals = $derived(sectionTotals(section, globals));

	let totalItems = $derived(
		section.subcategories.reduce((sum, sc) =>
			sum + sc.line_items.length +
			sc.component_groups.reduce((gs, g) => gs + g.line_items.length, 0),
		0)
	);

	function startRename() {
		editName = section.name;
		isEditing = true;
	}

	function commitRename() {
		const name = editName.trim();
		if (name && name !== section.name) {
			onsnapshot?.();
			section.name = name;
			onchange?.();
		}
		isEditing = false;
	}

	function addSubcategory() {
		onsnapshot?.();
		const name = newSubcatName.trim() || 'New Subcategory';
		section.subcategories.push({
			id: nanoid(),
			name,
			sort_order: section.subcategories.length,
			lump_sum: 0,
			markup_overrides: { materials: null, labor: null, equipment: null, subs: null, other: null },
			markup_enabled: { materials: true, labor: true, equipment: true, subs: true, other: true },
			line_items: [],
			component_groups: [],
		});
		newSubcatName = '';
		showAddSubcat = false;
		onchange?.();
	}

	function deleteSubcategory(subcatId) {
		onsnapshot?.();
		const idx = section.subcategories.findIndex(sc => sc.id === subcatId);
		if (idx !== -1) {
			section.subcategories.splice(idx, 1);
			onchange?.();
		}
	}
</script>

<div class="mb-6 border border-white/[0.06] overflow-hidden" style="background: var(--color-ink);">
	<!-- Section header -->
	<div class="flex items-center justify-between px-5 py-3 group/sec" style="background: var(--color-granite);">
		<button
			class="flex items-center gap-3 hover:bg-white/[0.04] -ml-2 px-2 py-1 transition-colors text-left flex-1 min-w-0"
			onclick={() => isCollapsed = !isCollapsed}
		>
			<svg class="w-4 h-4 text-white/30 transition-transform shrink-0 {isCollapsed ? '' : 'rotate-90'}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
			</svg>
			{#if isEditing}
				<!-- svelte-ignore a11y_autofocus -->
				<input
					type="text"
					bind:value={editName}
					autofocus
					class="bg-white/[0.06] text-[var(--color-white)] px-2 py-0.5 text-sm font-semibold border border-white/[0.08] focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] font-[var(--font-ui)] uppercase tracking-wide"
					onclick={(e) => e.stopPropagation()}
					onkeydown={(e) => { if (e.key === 'Enter') commitRename(); if (e.key === 'Escape') isEditing = false; }}
					onblur={commitRename}
				/>
			{:else}
				<span class="font-semibold font-[var(--font-ui)] uppercase tracking-wide text-[var(--color-white)]">{section.name}</span>
				<button
					onclick={(e) => { e.stopPropagation(); startRename(); }}
					class="opacity-0 group-hover/sec:opacity-100 text-white/20 hover:text-[var(--color-sunburst)] transition-opacity p-0.5 shrink-0"
					title="Rename section"
				>
					<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
					</svg>
				</button>
			{/if}
			<span class="text-xs text-white/30 font-[var(--font-body)]">
				{section.subcategories.length} subcategor{section.subcategories.length !== 1 ? 'ies' : 'y'}
				&middot; {totalItems} item{totalItems !== 1 ? 's' : ''}
			</span>
		</button>
		<div class="flex items-center gap-3">
			<span class="font-mono font-semibold text-[var(--color-white)]">{formatMoney(totals.withMarkup)}</span>
			<button
				onclick={() => ondelete?.(section.id)}
				class="text-white/20 hover:text-red-400 transition-colors p-1"
				title="Delete section"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
				</svg>
			</button>
		</div>
	</div>

	{#if !isCollapsed}
		{#if section.subcategories.length === 0 && !showAddSubcat}
			<div class="px-5 py-10 text-center text-[var(--color-muted-text)] text-sm font-[var(--font-body)]">
				No subcategories yet
			</div>
		{:else}
			{#each section.subcategories as subcat (subcat.id)}
				<SubcategoryBlock {subcat} {globals} {onchange} {onsnapshot} ondelete={deleteSubcategory} {materialsDb} {ratesDb} />
			{/each}
		{/if}

		<div class="px-5 pb-4">
			{#if showAddSubcat}
				<div class="flex items-center gap-2 mt-2">
					<input
						type="text"
						bind:value={newSubcatName}
						placeholder="Subcategory name"
						class="flex-1 px-0 py-1.5 bg-transparent border-0 border-b border-white/[0.08] text-[var(--color-white)] text-sm font-[var(--font-body)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/20"
						onkeydown={(e) => { if (e.key === 'Enter') addSubcategory(); if (e.key === 'Escape') showAddSubcat = false; }}
					/>
					<button onclick={addSubcategory} class="px-3 py-1.5 bg-[var(--color-sunburst)] text-[var(--color-ink)] text-xs font-[var(--font-ui)] font-bold uppercase tracking-wide hover:brightness-110">Add</button>
					<button onclick={() => showAddSubcat = false} class="px-2 py-1.5 text-[var(--color-muted-text)] text-xs hover:text-[var(--color-white)] font-[var(--font-ui)]">Cancel</button>
				</div>
			{:else}
				<button
					onclick={() => showAddSubcat = true}
					class="mt-2 text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] uppercase tracking-wide transition-colors"
				>
					<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
					</svg>
					Add Subcategory
				</button>
			{/if}
		</div>
	{/if}
</div>
