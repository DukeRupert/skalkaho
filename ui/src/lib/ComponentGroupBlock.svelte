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

<div class="ml-5 mt-3 border-l-2 border-white/[0.06] pl-4">
	<div class="flex items-center gap-2 mb-1 group/cg">
		{#if isEditing}
			<!-- svelte-ignore a11y_autofocus -->
			<input
				type="text"
				bind:value={editName}
				autofocus
				class="text-xs font-semibold text-[var(--color-concrete)] uppercase tracking-wide px-1 py-0.5 border border-white/[0.08] bg-white/[0.04] focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] font-[var(--font-ui)]"
				onkeydown={(e) => { if (e.key === 'Enter') commitRename(); if (e.key === 'Escape') isEditing = false; }}
				onblur={commitRename}
			/>
		{:else}
			<span class="text-xs font-semibold text-[var(--color-sage)] uppercase tracking-wider cursor-pointer font-[var(--font-ui)]" role="button" tabindex="0" ondblclick={startRename}>{group.name}</span>
			<button
				onclick={startRename}
				class="opacity-0 group-hover/cg:opacity-100 text-white/15 hover:text-[var(--color-sunburst)] transition-opacity p-0.5"
				title="Rename"
			>
				<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
				</svg>
			</button>
		{/if}
		<span class="text-xs text-white/25 font-[var(--font-body)]">({group.line_items.length})</span>
		<button
			onclick={() => ondelete?.(group.id)}
			class="opacity-0 group-hover/cg:opacity-100 text-white/20 hover:text-red-400 transition-opacity p-0.5"
			title="Delete group"
		>
			<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
			</svg>
		</button>
	</div>

	{#if group.line_items.length > 0}
		<table class="w-full">
			<tbody>
				{#each group.line_items as item (item.id)}
					<LineItemRow {item} {globals} {markupOverrides} {markupEnabled} {onchange} ondelete={deleteItem} />
				{/each}
			</tbody>
		</table>
	{/if}

	{#if showAddForm}
		<div class="mt-1 mb-1">
			<Autocomplete
				{materialsDb}
				{ratesDb}
				categoryType="materials"
				onselect={addItem}
				oncancel={() => showAddForm = false}
			/>
		</div>
	{:else}
		<button
			onclick={() => showAddForm = true}
			class="mt-1 text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] transition-colors"
		>
			<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
			</svg>
			Add Item
		</button>
	{/if}
</div>
