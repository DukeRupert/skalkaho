<script>
	import LineItemRow from './LineItemRow.svelte';
	import Autocomplete from './Autocomplete.svelte';
	import { nanoid } from 'nanoid';

	let { group, globals, markupOverrides, markupEnabled, onchange, onsnapshot, ondelete, materialsDb, ratesDb, subcontractorsDb } = $props();

	let showAddForm = $state(false);
	let addingToGroup = $state(null);
	let isEditing = $state(false);
	let editName = $state(group.name);
	let showAddVisualGroup = $state(false);
	let newVisualGroupName = $state('');

	// Derive visual groups from line items
	let visualGroups = $derived.by(() => {
		const groups = new Map();
		const ungrouped = [];
		for (const item of group.line_items) {
			if (item.visual_group) {
				if (!groups.has(item.visual_group)) groups.set(item.visual_group, []);
				groups.get(item.visual_group).push(item);
			} else {
				ungrouped.push(item);
			}
		}
		return { groups, ungrouped };
	});

	let collapsedGroups = $state({});

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
			subcontractor_id: itemData.subcontractor_id || null,
			price_override: false,
			description: null,
			sort_order: group.line_items.length,
			component_group_id: group.id,
			visual_group: addingToGroup,
		});
		showAddForm = false;
		addingToGroup = null;
		onchange?.();
	}

	function startAddItem(groupName = null) {
		addingToGroup = groupName;
		showAddForm = true;
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

	function addVisualGroup() {
		const name = newVisualGroupName.trim();
		if (!name) return;
		newVisualGroupName = '';
		showAddVisualGroup = false;
		startAddItem(name);
	}

	function deleteVisualGroup(groupName) {
		onsnapshot?.();
		for (const item of group.line_items) {
			if (item.visual_group === groupName) {
				item.visual_group = null;
			}
		}
		onchange?.();
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
			title="Delete section"
		>
			<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
			</svg>
		</button>
	</div>

	{#if group.line_items.length > 0}
		<table class="w-full">
			<tbody>
				<!-- Ungrouped items -->
				{#each visualGroups.ungrouped as item (item.id)}
					<LineItemRow {item} {globals} {markupOverrides} {markupEnabled} {onchange} ondelete={deleteItem} />
				{/each}
				<!-- Visual groups -->
				{#each [...visualGroups.groups.entries()] as [vgName, vgItems] (vgName)}
					{@const isVgCollapsed = collapsedGroups[vgName] ?? false}
					<tr class="border-b border-white/[0.04]">
						<td colspan="9" class="px-0 py-0">
							<div class="flex items-center gap-2 py-1.5 px-1 bg-white/[0.02] group/vg">
								<button
									onclick={() => collapsedGroups[vgName] = !isVgCollapsed}
									class="flex items-center gap-1.5 text-xs font-semibold text-[var(--color-sage)] uppercase tracking-wider font-[var(--font-ui)]"
								>
									<svg class="w-2.5 h-2.5 text-white/25 transition-transform {isVgCollapsed ? '' : 'rotate-90'}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
									</svg>
									{vgName}
								</button>
								<span class="text-xs text-white/20 font-[var(--font-body)]">({vgItems.length})</span>
								<button
									onclick={() => deleteVisualGroup(vgName)}
									class="opacity-0 group-hover/vg:opacity-100 text-white/20 hover:text-red-400 transition-opacity p-0.5 ml-auto"
									title="Ungroup items"
								>
									<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
									</svg>
								</button>
							</div>
						</td>
					</tr>
					{#if !isVgCollapsed}
						{#each vgItems as item (item.id)}
							<LineItemRow {item} {globals} {markupOverrides} {markupEnabled} {onchange} ondelete={deleteItem} />
						{/each}
						{#if showAddForm && addingToGroup === vgName}
							<tr><td colspan="9" class="px-1 py-1">
								<Autocomplete
									{materialsDb}
									{ratesDb}
									{subcontractorsDb}
									categoryType="materials"
									onselect={addItem}
									oncancel={() => { showAddForm = false; addingToGroup = null; }}
								/>
							</td></tr>
						{:else}
							<tr><td colspan="9" class="px-1 py-0.5">
								<button
									onclick={() => startAddItem(vgName)}
									class="text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] transition-colors ml-4"
								>
									<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
									</svg>
									Add Item
								</button>
							</td></tr>
						{/if}
					{/if}
				{/each}
				<!-- Pending new visual group (no items yet) -->
				{#if showAddForm && addingToGroup !== null && !visualGroups.groups.has(addingToGroup)}
					<tr class="border-b border-white/[0.04]">
						<td colspan="9" class="px-0 py-0">
							<div class="flex items-center gap-2 py-1.5 px-1 bg-white/[0.02]">
								<span class="flex items-center gap-1.5 text-xs font-semibold text-[var(--color-sage)] uppercase tracking-wider font-[var(--font-ui)]">
									<svg class="w-2.5 h-2.5 text-white/25 rotate-90" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
									</svg>
									{addingToGroup}
								</span>
							</div>
						</td>
					</tr>
					<tr><td colspan="9" class="px-1 py-1">
						<Autocomplete
							{materialsDb}
							{ratesDb}
							{subcontractorsDb}
							categoryType="materials"
							onselect={addItem}
							oncancel={() => { showAddForm = false; addingToGroup = null; }}
						/>
					</td></tr>
				{/if}
			</tbody>
		</table>
	{/if}

	<!-- Pending new visual group when no table exists yet -->
	{#if showAddForm && addingToGroup !== null && !visualGroups.groups.has(addingToGroup) && group.line_items.length === 0}
		<div class="mb-1">
			<div class="flex items-center gap-2 py-1.5 px-1 bg-white/[0.02] border-b border-white/[0.04]">
				<span class="flex items-center gap-1.5 text-xs font-semibold text-[var(--color-sage)] uppercase tracking-wider font-[var(--font-ui)]">
					<svg class="w-2.5 h-2.5 text-white/25 rotate-90" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
					</svg>
					{addingToGroup}
				</span>
			</div>
			<Autocomplete
				{materialsDb}
				{ratesDb}
				{subcontractorsDb}
				categoryType="materials"
				onselect={addItem}
				oncancel={() => { showAddForm = false; addingToGroup = null; }}
			/>
		</div>
	{/if}

	<!-- Add visual group -->
	{#if showAddVisualGroup}
		<div class="flex items-center gap-2 mt-1">
			<input
				type="text"
				bind:value={newVisualGroupName}
				placeholder="Group label"
				class="flex-1 px-0 py-1 bg-transparent border-0 border-b border-white/[0.08] text-[var(--color-white)] text-sm font-[var(--font-body)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/20"
				onkeydown={(e) => { if (e.key === 'Enter') addVisualGroup(); if (e.key === 'Escape') showAddVisualGroup = false; }}
			/>
			<button onclick={addVisualGroup} class="px-2 py-1 bg-[var(--color-sunburst)] text-[var(--color-ink)] text-xs font-[var(--font-ui)] font-bold uppercase tracking-wide hover:brightness-110">Add</button>
			<button onclick={() => showAddVisualGroup = false} class="px-2 py-1 text-[var(--color-muted-text)] text-xs hover:text-[var(--color-white)] font-[var(--font-ui)]">Cancel</button>
		</div>
	{:else}
		<button
			onclick={() => showAddVisualGroup = true}
			class="mt-1 text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] transition-colors"
		>
			<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7"/>
			</svg>
			Add Label Group
		</button>
	{/if}

	<!-- Add ungrouped item -->
	{#if showAddForm && addingToGroup === null}
		<div class="mt-1 mb-1">
			<Autocomplete
				{materialsDb}
				{ratesDb}
				{subcontractorsDb}
				categoryType="materials"
				onselect={addItem}
				oncancel={() => { showAddForm = false; addingToGroup = null; }}
			/>
		</div>
	{:else}
		<button
			onclick={() => startAddItem(null)}
			class="mt-1 text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] transition-colors"
		>
			<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
			</svg>
			Add Item
		</button>
	{/if}
</div>
