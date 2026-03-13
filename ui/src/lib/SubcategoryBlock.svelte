<script>
	import LineItemRow from './LineItemRow.svelte';
	import ComponentGroupBlock from './ComponentGroupBlock.svelte';
	import Autocomplete from './Autocomplete.svelte';
	import { subcategoryTotals, formatMoney, resolveMarkup, formatPercent } from './markup.js';
	import { nanoid } from 'nanoid';

	let { subcat, globals, collapsed = false, onchange, onsnapshot, ondelete, materialsDb = [], ratesDb = [] } = $props();

	let isCollapsed = $state(collapsed);
	let showAddForm = $state(false);
	let isEditing = $state(false);
	let editName = $state(subcat.name);
	let showAddGroup = $state(false);
	let newGroupName = $state('');
	let totals = $derived(subcategoryTotals(subcat, globals));

	let markupSummary = $derived.by(() => {
		const types = ['materials', 'labor', 'equipment', 'subs', 'other'];
		return types.map(t => ({
			type: t,
			value: resolveMarkup(t, globals, subcat.markup_overrides, subcat.markup_enabled),
			isOverride: subcat.markup_overrides[t] != null,
			isDisabled: !subcat.markup_enabled[t],
		}));
	});

	let showMarkupPanel = $state(false);
	let hasOverrides = $derived(markupSummary.some(m => m.isOverride || m.isDisabled));
	let totalItems = $derived(
		subcat.line_items.length +
		subcat.component_groups.reduce((sum, g) => sum + g.line_items.length, 0)
	);

	function addItem(itemData) {
		onsnapshot?.();
		subcat.line_items.push({
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
			sort_order: subcat.line_items.length,
			component_group_id: null,
		});
		showAddForm = false;
		onchange?.();
	}

	function deleteItem(itemId) {
		onsnapshot?.();
		const idx = subcat.line_items.findIndex(li => li.id === itemId);
		if (idx !== -1) {
			subcat.line_items.splice(idx, 1);
			onchange?.();
		}
	}

	function startRename() {
		editName = subcat.name;
		isEditing = true;
	}

	function commitRename() {
		const name = editName.trim();
		if (name && name !== subcat.name) {
			onsnapshot?.();
			subcat.name = name;
			onchange?.();
		}
		isEditing = false;
	}

	function addComponentGroup() {
		onsnapshot?.();
		const name = newGroupName.trim() || 'New Group';
		subcat.component_groups.push({
			id: nanoid(),
			name,
			sort_order: subcat.component_groups.length,
			line_items: [],
		});
		newGroupName = '';
		showAddGroup = false;
		onchange?.();
	}

	function deleteComponentGroup(groupId) {
		onsnapshot?.();
		const idx = subcat.component_groups.findIndex(g => g.id === groupId);
		if (idx !== -1) {
			subcat.component_groups.splice(idx, 1);
			onchange?.();
		}
	}

	function handleOverrideInput(type, e) {
		const raw = e.target.value.trim();
		if (raw === '') {
			subcat.markup_overrides[type] = null;
		} else {
			const val = parseFloat(raw);
			if (!isNaN(val) && val >= 0) {
				subcat.markup_overrides[type] = val;
			}
		}
		onchange?.();
	}

	function handleToggle(type) {
		subcat.markup_enabled[type] = !subcat.markup_enabled[type];
		onchange?.();
	}

	function handleLumpSum(e) {
		const val = parseFloat(e.target.value);
		if (!isNaN(val) && val >= 0) {
			subcat.lump_sum = val;
			onchange?.();
		}
	}

	const MARKUP_TYPES = [
		{ key: 'materials', label: 'Mat', color: 'text-blue-700' },
		{ key: 'labor', label: 'Lab', color: 'text-amber-700' },
		{ key: 'equipment', label: 'Equip', color: 'text-purple-700' },
		{ key: 'subs', label: 'Subs', color: 'text-green-700' },
		{ key: 'other', label: 'Other', color: 'text-slate-600' },
	];
</script>

<div class="border-t border-slate-200">
	<!-- Subcategory header -->
	<div class="flex items-center justify-between px-4 py-2 hover:bg-slate-50 transition-colors group/subcat">
		<button
			class="flex items-center gap-2 text-left flex-1 min-w-0"
			onclick={() => isCollapsed = !isCollapsed}
		>
			<svg class="w-4 h-4 text-slate-400 transition-transform {isCollapsed ? '' : 'rotate-90'}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
			</svg>
			{#if isEditing}
				<!-- svelte-ignore a11y_autofocus -->
				<input
					type="text"
					bind:value={editName}
					autofocus
					class="px-2 py-0.5 border border-slate-300 rounded text-sm font-medium text-slate-600 focus:ring-2 focus:ring-blue-400"
					onclick={(e) => e.stopPropagation()}
					onkeydown={(e) => { if (e.key === 'Enter') commitRename(); if (e.key === 'Escape') isEditing = false; }}
					onblur={commitRename}
				/>
			{:else}
				<span role="button" tabindex="0" class="font-medium text-slate-600 text-sm" ondblclick={(e) => { e.stopPropagation(); startRename(); }}>{subcat.name}</span>
			{/if}
			<span class="text-xs text-slate-400">{totalItems} item{totalItems !== 1 ? 's' : ''}</span>
			{#if hasOverrides}
				<span class="text-xs px-1.5 py-0.5 rounded bg-amber-50 text-amber-600 font-medium">overrides</span>
			{/if}
			{#if subcat.lump_sum > 0}
				<span class="text-xs px-1.5 py-0.5 rounded bg-green-50 text-green-600 font-medium">
					+{formatMoney(subcat.lump_sum)} lump sum
				</span>
			{/if}
		</button>
		<div class="flex items-center gap-2">
			<button
				onclick={() => showMarkupPanel = !showMarkupPanel}
				class="text-xs px-1.5 py-0.5 rounded bg-slate-100 text-slate-500 hover:bg-slate-200 transition-colors"
				title="Configure markup"
			>
				<svg class="w-3 h-3 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"/>
				</svg>
			</button>
			<span class="font-mono text-sm font-semibold text-slate-700">{formatMoney(totals.withMarkup)}</span>
			<button
				onclick={() => ondelete?.(subcat.id)}
				class="opacity-0 group-hover/subcat:opacity-100 text-slate-400 hover:text-red-500 transition-opacity p-0.5"
				title="Delete subcategory"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
				</svg>
			</button>
		</div>
	</div>

	{#if !isCollapsed}
		<div class="px-4 pb-3">
			<!-- Markup controls panel -->
			{#if showMarkupPanel}
				<div class="bg-slate-50 rounded-lg border border-slate-200 p-3 mb-3">
					<div class="flex items-center justify-between mb-2">
						<span class="text-xs font-semibold text-slate-600 uppercase tracking-wide">Markup Overrides</span>
						<button onclick={() => showMarkupPanel = false} class="text-xs text-slate-400 hover:text-slate-600">Close</button>
					</div>
					<div class="grid grid-cols-5 gap-2">
						{#each MARKUP_TYPES as mt}
							{@const m = markupSummary.find(x => x.type === mt.key)}
							<div class="text-center">
								<span class="block text-xs font-medium {mt.color} mb-1">{mt.label}</span>
								<div class="flex items-center justify-center gap-1 mb-1">
									<input
										type="checkbox"
										checked={subcat.markup_enabled[mt.key]}
										onchange={() => handleToggle(mt.key)}
										class="w-3 h-3 rounded border-slate-300 text-blue-600 focus:ring-blue-400"
									/>
									<span class="text-xs text-slate-400">{subcat.markup_enabled[mt.key] ? 'On' : 'Off'}</span>
								</div>
								<input
									type="number"
									value={subcat.markup_overrides[mt.key] ?? ''}
									oninput={(e) => handleOverrideInput(mt.key, e)}
									placeholder={`${globals[`${mt.key}_markup`]}%`}
									step="1"
									min="0"
									disabled={!subcat.markup_enabled[mt.key]}
									class="w-full text-center text-xs font-mono px-1 py-1 border border-slate-200 rounded
										focus:ring-1 focus:ring-blue-400 focus:border-blue-400
										{!subcat.markup_enabled[mt.key] ? 'bg-slate-100 text-slate-400' : 'bg-white'}
										{m?.isOverride ? 'border-amber-300 bg-amber-50' : ''}"
								/>
								<div class="text-xs text-slate-400 mt-0.5">
									eff: {formatPercent(m?.value ?? 0)}
								</div>
							</div>
						{/each}
					</div>
					<!-- Lump sum -->
					<div class="mt-3 pt-2 border-t border-slate-200 flex items-center gap-2">
						<span class="text-xs font-medium text-slate-600">Lump Sum:</span>
						<span class="text-xs text-slate-400">$</span>
						<input
							type="number"
							value={subcat.lump_sum}
							oninput={handleLumpSum}
							step="0.01"
							min="0"
							class="w-28 text-right text-xs font-mono px-2 py-1 border border-slate-200 rounded
								focus:ring-1 focus:ring-blue-400 focus:border-blue-400"
						/>
						<span class="text-xs text-slate-400">added post-markup</span>
					</div>
				</div>
			{:else if hasOverrides}
				<div class="flex items-center gap-3 py-1.5 px-2 bg-amber-50 rounded text-xs mb-2 cursor-pointer hover:bg-amber-100 transition-colors" onclick={() => showMarkupPanel = true} onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') showMarkupPanel = true; }} role="button" tabindex="0">
					<span class="font-medium text-amber-700">Markup:</span>
					{#each markupSummary as m}
						<span class={m.isDisabled ? 'text-slate-400 line-through' : m.isOverride ? 'text-amber-700 font-medium' : 'text-slate-500'}>
							{m.type} {formatPercent(m.value)}
						</span>
					{/each}
				</div>
			{/if}

			{#if totalItems > 0 || subcat.line_items.length > 0}
				<table class="w-full">
					<thead>
						<tr class="text-xs text-slate-400 uppercase tracking-wide border-b border-slate-100">
							<th class="px-1 py-1 text-left w-24">Type</th>
							<th class="px-1 py-1 text-left">Name</th>
							<th class="px-1 py-1 text-right w-20">Qty</th>
							<th class="px-1 py-1 text-center w-16">Unit</th>
							<th class="px-1 py-1 text-right w-24">Price</th>
							<th class="px-2 py-1 text-right w-16">Markup</th>
							<th class="px-2 py-1 text-right w-24">w/ Markup</th>
							<th class="px-2 py-1 text-right w-24">Total</th>
							<th class="w-8"></th>
						</tr>
					</thead>
					<tbody>
						{#each subcat.line_items as item (item.id)}
							<LineItemRow
								{item}
								{globals}
								markupOverrides={subcat.markup_overrides}
								markupEnabled={subcat.markup_enabled}
								{onchange}
								ondelete={deleteItem}
							/>
						{/each}
					</tbody>
				</table>
			{/if}

			{#each subcat.component_groups as group (group.id)}
				<ComponentGroupBlock
					{group}
					{globals}
					markupOverrides={subcat.markup_overrides}
					markupEnabled={subcat.markup_enabled}
					{onchange}
					{onsnapshot}
					ondelete={deleteComponentGroup}
					{materialsDb}
					{ratesDb}
				/>
			{/each}

			<!-- Add component group -->
			{#if showAddGroup}
				<div class="flex items-center gap-2 ml-4 mt-2">
					<input
						type="text"
						bind:value={newGroupName}
						placeholder="Group name"
						class="flex-1 px-2 py-1 border border-slate-300 rounded text-sm focus:ring-2 focus:ring-blue-400 focus:border-blue-400"
						onkeydown={(e) => { if (e.key === 'Enter') addComponentGroup(); if (e.key === 'Escape') showAddGroup = false; }}
					/>
					<button onclick={addComponentGroup} class="px-2 py-1 bg-slate-800 text-white text-xs rounded hover:bg-slate-700">Add</button>
					<button onclick={() => showAddGroup = false} class="px-2 py-1 text-slate-500 text-xs hover:text-slate-700">Cancel</button>
				</div>
			{:else}
				<button
					onclick={() => showAddGroup = true}
					class="ml-4 mt-2 text-xs text-blue-500 hover:text-blue-700 flex items-center gap-1"
				>
					<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
					</svg>
					Add Group
				</button>
			{/if}

			<!-- Add item button + autocomplete -->
			{#if showAddForm}
				<div class="mt-2">
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
					class="mt-2 text-sm text-blue-500 hover:text-blue-700 flex items-center gap-1"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
					</svg>
					Add Item
				</button>
			{/if}

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
