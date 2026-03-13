<script>
	import LineItemRow from './LineItemRow.svelte';
	import Autocomplete from './Autocomplete.svelte';
	import { nanoid } from 'nanoid';

	let { group, globals, markupOverrides, markupEnabled, onchange, materialsDb, ratesDb } = $props();

	let showAddForm = $state(false);

	function addItem(itemData) {
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
		const idx = group.line_items.findIndex(li => li.id === itemId);
		if (idx !== -1) {
			group.line_items.splice(idx, 1);
			onchange?.();
		}
	}
</script>

<div class="ml-4 mt-2">
	<div class="flex items-center gap-2 mb-1">
		<span class="text-xs font-semibold text-slate-500 uppercase tracking-wide">{group.name}</span>
		<span class="text-xs text-slate-400">({group.line_items.length})</span>
		<button
			onclick={() => showAddForm = !showAddForm}
			class="text-xs text-blue-500 hover:text-blue-700 ml-auto"
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
