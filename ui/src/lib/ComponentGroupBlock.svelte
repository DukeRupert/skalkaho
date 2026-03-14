<script>
	import LineItemRow from './LineItemRow.svelte';
	import Autocomplete from './Autocomplete.svelte';
	import { nanoid } from 'nanoid';

	let { group, globals, markupOverrides, markupEnabled, onchange, onsnapshot, ondelete, materialsDb, ratesDb } = $props();

	let showAddForm = $state(false);
	let isEditing = $state(false);
	let editName = $state(group.name);

	function addItem(itemData) {
		onsnapshot?.();
		group.line_items.push({
			id: nanoid(),
			category_type: itemData.category_type,
			item_name: itemData.item_name,
			quantity: 1,
			unit: itemData.unit,
			unit_price: itemData.unit_price,
			is_custom: itemData.is_custom,
			material_id: itemData.material_id,
			price_override: false,
			description: null,
			sort_order: group.line_items.length,
			component_group_id: group.id,
		});
		showAddForm = false;
		onchange?.();
	}

	function deleteItem(itemId) {
		onsnapshot?.();
		const idx = group.line_items.findIndex(li => li.id === itemId);
		if (idx !== -1) {
			group.line_items.splice(idx, 1);
			onchange?.();
		}
	}

	function startRename() {
		editName = group.name;
		isEditing = true;
	}

	function commitRename() {
		const name = editName.trim();
		if (name && name !== group.name) {
			onsnapshot?.();
			group.name = name;
			onchange?.();
		}
		isEditing = false;
	}
</script>

<div class="ml-4 mt-2">
	<div class="flex items-center gap-2 mb-1 group/cg">
		{#if isEditing}
			<!-- svelte-ignore a11y_autofocus -->
			<input
				type="text"
				bind:value={editName}
				autofocus
				class="text-xs font-semibold text-[var(--color-concrete)] uppercase tracking-wide px-1 py-0.5 border border-white/[0.08] rounded bg-white/[0.04] focus:ring-2 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] font-[var(--font-ui)]"
				onkeydown={(e) => { if (e.key === 'Enter') commitRename(); if (e.key === 'Escape') isEditing = false; }}
				onblur={commitRename}
			/>
		{:else}
			<span role="button" tabindex="0" class="text-xs font-semibold text-[var(--color-muted-text)] uppercase tracking-wide cursor-pointer font-[var(--font-ui)]" ondblclick={startRename}>{group.name}</span>
		{/if}
		<span class="text-xs text-white/40 font-[var(--font-body)]">({group.line_items.length})</span>
		<button
			onclick={() => ondelete?.(group.id)}
			class="opacity-0 group-hover/cg:opacity-100 text-white/30 hover:text-red-400 transition-opacity p-0.5"
			title="Delete group"
		>
			<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
			</svg>
		</button>
		<button
			onclick={() => showAddForm = !showAddForm}
			class="text-xs text-[var(--color-sunburst)] hover:brightness-110 ml-auto font-[var(--font-ui)]"
		>
			+ Add Item
		</button>
	</div>

	{#if showAddForm}
		<div class="mb-2">
			<Autocomplete
				{materialsDb}
				{ratesDb}
				categoryType="materials"
				onselect={addItem}
				oncancel={() => showAddForm = false}
			/>
		</div>
	{/if}

	{#if group.line_items.length > 0}
		<table class="w-full">
			<tbody>
				{#each group.line_items as item (item.id)}
					<LineItemRow {item} {globals} {markupOverrides} {markupEnabled} {onchange} ondelete={deleteItem} />
				{/each}
			</tbody>
		</table>
	{/if}
</div>
