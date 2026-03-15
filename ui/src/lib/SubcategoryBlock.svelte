<script>
	import LineItemRow from './LineItemRow.svelte';
	import ComponentGroupBlock from './ComponentGroupBlock.svelte';
	import Autocomplete from './Autocomplete.svelte';
	import { subcategoryTotals, formatMoney, resolveMarkup, formatPercent } from './markup.js';
	import { nanoid } from 'nanoid';

	let { subcat, globals, collapsed = false, onchange, onsnapshot, ondelete, materialsDb = [], ratesDb = [], subcontractorsDb = [] } = $props();

	let isCollapsed = $state(collapsed);
	let showAddForm = $state(false);
	let addingToGroup = $state(null); // null = ungrouped, string = visual_group name
	let isEditing = $state(false);
	let editName = $state(subcat.name);
	let showAddGroup = $state(false);
	let newGroupName = $state('');
	let showAddVisualGroup = $state(false);
	let newVisualGroupName = $state('');
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

	// Derive visual groups from line items
	let visualGroups = $derived.by(() => {
		const groups = new Map(); // name → items[]
		const ungrouped = [];
		for (const item of subcat.line_items) {
			if (item.visual_group) {
				if (!groups.has(item.visual_group)) groups.set(item.visual_group, []);
				groups.get(item.visual_group).push(item);
			} else {
				ungrouped.push(item);
			}
		}
		return { groups, ungrouped };
	});

	// Track collapsed state per visual group
	let collapsedGroups = $state({});

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
			subcontractor_id: itemData.subcontractor_id || null,
			price_override: false,
			description: null,
			sort_order: subcat.line_items.length,
			component_group_id: null,
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
		const name = newGroupName.trim() || 'New section';
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

	function addVisualGroup() {
		const name = newVisualGroupName.trim();
		if (!name) return;
		// Just create the group by adding a placeholder — the group exists when items reference it
		// For now, toggle the add-item form for this group so user can immediately add items
		newVisualGroupName = '';
		showAddVisualGroup = false;
		startAddItem(name);
	}

	function deleteVisualGroup(groupName) {
		onsnapshot?.();
		// Remove group label from all items in this group (makes them ungrouped)
		for (const item of subcat.line_items) {
			if (item.visual_group === groupName) {
				item.visual_group = null;
			}
		}
		onchange?.();
	}

	function renameVisualGroup(oldName, newName) {
		if (!newName.trim() || newName === oldName) return;
		onsnapshot?.();
		for (const item of subcat.line_items) {
			if (item.visual_group === oldName) {
				item.visual_group = newName.trim();
			}
		}
		onchange?.();
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

	const MARKUP_LABELS = { materials: 'Mat', labor: 'Lab', equipment: 'Equip', subs: 'Subs', other: 'Other' };
</script>

<div class="border-t border-white/[0.06]">
	<!-- Subcategory header -->
	<div class="flex items-center justify-between px-5 py-2 hover:bg-white/[0.02] transition-colors group/subcat">
		<button
			class="flex items-center gap-2 text-left flex-1 min-w-0"
			onclick={() => isCollapsed = !isCollapsed}
		>
			<svg class="w-3 h-3 text-white/25 transition-transform shrink-0 {isCollapsed ? '' : 'rotate-90'}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
			</svg>
			{#if isEditing}
				<!-- svelte-ignore a11y_autofocus -->
				<input
					type="text"
					bind:value={editName}
					autofocus
					class="px-2 py-0.5 border border-white/[0.08] text-sm font-medium text-[var(--color-white)] bg-white/[0.04] focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] font-[var(--font-ui)]"
					onclick={(e) => e.stopPropagation()}
					onkeydown={(e) => { if (e.key === 'Enter') commitRename(); if (e.key === 'Escape') isEditing = false; }}
					onblur={commitRename}
				/>
			{:else}
				<span class="font-medium text-[var(--color-concrete)] text-sm font-[var(--font-ui)]">{subcat.name}</span>
				<button
					onclick={(e) => { e.stopPropagation(); startRename(); }}
					class="opacity-0 group-hover/subcat:opacity-100 text-white/15 hover:text-[var(--color-sunburst)] transition-opacity p-0.5 shrink-0"
					title="Rename"
				>
					<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
					</svg>
				</button>
			{/if}
			<span class="text-xs text-white/25 font-[var(--font-body)]">{totalItems} item{totalItems !== 1 ? 's' : ''}</span>
			{#if hasOverrides}
				<span class="text-xs px-1.5 py-0.5 bg-[var(--color-sunburst)]/10 text-[var(--color-sunburst)] font-medium font-[var(--font-ui)]">overrides</span>
			{/if}
			{#if subcat.lump_sum > 0}
				<span class="text-xs px-1.5 py-0.5 bg-[var(--color-sage)]/10 text-[var(--color-sage)] font-medium font-[var(--font-ui)]">
					+{formatMoney(subcat.lump_sum)} lump sum
				</span>
			{/if}
		</button>
		<div class="flex items-center gap-2">
			<button
				onclick={() => showMarkupPanel = !showMarkupPanel}
				class="text-xs px-1.5 py-0.5 bg-white/[0.04] text-[var(--color-muted-text)] hover:bg-white/[0.08] hover:text-[var(--color-concrete)] transition-colors"
				title="Configure markup"
			>
				<svg class="w-3 h-3 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"/>
				</svg>
			</button>
			<span class="font-mono text-sm font-semibold text-[var(--color-white)]">{formatMoney(totals.withMarkup)}</span>
			<button
				onclick={() => ondelete?.(subcat.id)}
				class="opacity-0 group-hover/subcat:opacity-100 text-white/20 hover:text-red-400 transition-opacity p-0.5"
				title="Delete subcategory"
			>
				<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
				</svg>
			</button>
		</div>
	</div>

	{#if !isCollapsed}
		<div class="px-5 pb-3">
			<!-- Markup controls panel -->
			{#if showMarkupPanel}
				<div class="bg-white/[0.02] border border-white/[0.06] p-3 mb-3">
					<div class="flex items-center justify-between mb-2">
						<span class="text-xs font-semibold text-[var(--color-sage)] uppercase tracking-wider font-[var(--font-ui)]">Markup Overrides</span>
						<button onclick={() => showMarkupPanel = false} class="text-xs text-white/30 hover:text-[var(--color-white)] font-[var(--font-ui)]">Close</button>
					</div>
					<div class="grid grid-cols-5 gap-3">
						{#each ['materials', 'labor', 'equipment', 'subs', 'other'] as type}
							{@const m = markupSummary.find(x => x.type === type)}
							<div class="text-center">
								<span class="block text-xs font-medium text-[var(--color-concrete)] mb-1 font-[var(--font-ui)]">{MARKUP_LABELS[type]}</span>
								<div class="flex items-center justify-center gap-1 mb-1">
									<input
										type="checkbox"
										checked={subcat.markup_enabled[type]}
										onchange={() => handleToggle(type)}
										class="w-3 h-3 border-white/[0.1] bg-white/[0.04] text-[var(--color-sunburst)] focus:ring-[var(--color-sunburst)] rounded-sm"
									/>
									<span class="text-xs text-white/30 font-[var(--font-body)]">{subcat.markup_enabled[type] ? 'On' : 'Off'}</span>
								</div>
								<input
									type="number"
									value={subcat.markup_overrides[type] ?? ''}
									oninput={(e) => handleOverrideInput(type, e)}
									placeholder={`${globals[`${type}_markup`]}%`}
									step="1"
									min="0"
									disabled={!subcat.markup_enabled[type]}
									class="w-full text-center text-xs font-mono px-1 py-1 border border-white/[0.08]
										focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)]
										{!subcat.markup_enabled[type] ? 'bg-white/[0.01] text-white/15' : 'bg-white/[0.04] text-[var(--color-white)]'}
										{m?.isOverride ? 'border-[var(--color-sunburst)]/40 bg-[var(--color-sunburst)]/5' : ''}"
								/>
								<div class="text-xs text-white/30 mt-0.5 font-mono">
									eff: {formatPercent(m?.value ?? 0)}
								</div>
							</div>
						{/each}
					</div>
					<!-- Lump sum -->
					<div class="mt-3 pt-2 border-t border-white/[0.06] flex items-center gap-2">
						<span class="text-xs font-medium text-[var(--color-concrete)] font-[var(--font-ui)]">Lump Sum:</span>
						<span class="text-xs text-white/30">$</span>
						<input
							type="number"
							value={subcat.lump_sum}
							oninput={handleLumpSum}
							step="0.01"
							min="0"
							class="w-28 text-right text-xs font-mono px-2 py-1 border border-white/[0.08] bg-white/[0.04] text-[var(--color-white)]
								focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)]"
						/>
						<span class="text-xs text-white/30 font-[var(--font-body)]">added post-markup</span>
					</div>
				</div>
			{:else if hasOverrides}
				<div class="flex items-center gap-3 py-1.5 px-3 bg-[var(--color-sunburst)]/5 text-xs mb-2 cursor-pointer hover:bg-[var(--color-sunburst)]/8 transition-colors border-l-2 border-[var(--color-sunburst)]/30" onclick={() => showMarkupPanel = true} onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') showMarkupPanel = true; }} role="button" tabindex="0">
					<span class="font-medium text-[var(--color-sunburst)] font-[var(--font-ui)]">Markup:</span>
					{#each markupSummary as m}
						<span class="{m.isDisabled ? 'text-white/15 line-through' : m.isOverride ? 'text-[var(--color-sunburst)]' : 'text-white/30'} font-[var(--font-body)]">
							{m.type} {formatPercent(m.value)}
						</span>
					{/each}
				</div>
			{/if}

			{#if totalItems > 0 || subcat.line_items.length > 0}
				<table class="w-full">
					<thead>
						<tr class="text-xs text-white/30 uppercase tracking-wider border-b border-white/[0.06] font-[var(--font-ui)]">
							<th class="px-1 py-1.5 text-left w-24">Type</th>
							<th class="px-1 py-1.5 text-left">Name</th>
							<th class="px-1 py-1.5 text-right w-20">Qty</th>
							<th class="px-1 py-1.5 text-center w-16">Unit</th>
							<th class="px-1 py-1.5 text-right w-24">Price</th>
							<th class="px-2 py-1.5 text-right w-16">Markup</th>
							<th class="px-2 py-1.5 text-right w-24">w/ Markup</th>
							<th class="px-2 py-1.5 text-right w-24">Total</th>
							<th class="w-16"></th>
						</tr>
					</thead>
					<tbody>
						<!-- Ungrouped items first -->
						{#each visualGroups.ungrouped as item (item.id)}
							<LineItemRow
								{item}
								{globals}
								markupOverrides={subcat.markup_overrides}
								markupEnabled={subcat.markup_enabled}
								{onchange}
								ondelete={deleteItem}
							/>
						{/each}
						<!-- Visual groups -->
						{#each [...visualGroups.groups.entries()] as [groupName, groupItems] (groupName)}
							{@const isGroupCollapsed = collapsedGroups[groupName] ?? false}
							<tr class="border-b border-white/[0.04]">
								<td colspan="9" class="px-0 py-0">
									<div class="flex items-center gap-2 py-1.5 px-1 bg-white/[0.02] group/vg">
										<button
											onclick={() => collapsedGroups[groupName] = !isGroupCollapsed}
											class="flex items-center gap-1.5 text-xs font-semibold text-[var(--color-sage)] uppercase tracking-wider font-[var(--font-ui)]"
										>
											<svg class="w-2.5 h-2.5 text-white/25 transition-transform {isGroupCollapsed ? '' : 'rotate-90'}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
											</svg>
											{groupName}
										</button>
										<span class="text-xs text-white/20 font-[var(--font-body)]">({groupItems.length})</span>
										<button
											onclick={() => deleteVisualGroup(groupName)}
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
							{#if !isGroupCollapsed}
								{#each groupItems as item (item.id)}
									<LineItemRow
										{item}
										{globals}
										markupOverrides={subcat.markup_overrides}
										markupEnabled={subcat.markup_enabled}
										{onchange}
										ondelete={deleteItem}
									/>
								{/each}
								<!-- Add item within this visual group -->
								{#if showAddForm && addingToGroup === groupName}
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
											onclick={() => startAddItem(groupName)}
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
			{#if showAddForm && addingToGroup !== null && !visualGroups.groups.has(addingToGroup) && totalItems === 0}
				<div class="mb-2">
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
					{subcontractorsDb}
				/>
			{/each}

			<!-- Add section -->
			{#if showAddGroup}
				<div class="flex items-center gap-2 ml-5 mt-2">
					<input
						type="text"
						bind:value={newGroupName}
						placeholder="Section name"
						class="flex-1 px-0 py-1 bg-transparent border-0 border-b border-white/[0.08] text-[var(--color-white)] text-sm font-[var(--font-body)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/20"
						onkeydown={(e) => { if (e.key === 'Enter') addComponentGroup(); if (e.key === 'Escape') showAddGroup = false; }}
					/>
					<button onclick={addComponentGroup} class="px-2 py-1 bg-[var(--color-sunburst)] text-[var(--color-ink)] text-xs font-[var(--font-ui)] font-bold uppercase tracking-wide hover:brightness-110">Add</button>
					<button onclick={() => showAddGroup = false} class="px-2 py-1 text-[var(--color-muted-text)] text-xs hover:text-[var(--color-white)] font-[var(--font-ui)]">Cancel</button>
				</div>
			{:else}
				<button
					onclick={() => showAddGroup = true}
					class="ml-5 mt-2 text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] transition-colors"
				>
					<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
					</svg>
					Add Section
				</button>
			{/if}

			<!-- Add visual group -->
			{#if showAddVisualGroup}
				<div class="flex items-center gap-2 mt-2">
					<input
						type="text"
						bind:value={newVisualGroupName}
						placeholder="Group label (e.g. Interior, Phase 1)"
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
				<div class="mt-2">
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
					class="mt-2 text-xs text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] flex items-center gap-1 font-[var(--font-ui)] transition-colors"
				>
					<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
					</svg>
					Add Item
				</button>
			{/if}

			<div class="flex justify-between items-center mt-3 pt-2 border-t border-white/[0.06] text-sm">
				<span class="text-white/30 font-[var(--font-body)]">
					Subtotal: {formatMoney(totals.base)}
					{#if subcat.lump_sum > 0}
						<span class="text-[var(--color-sage)]"> + {formatMoney(subcat.lump_sum)} lump sum</span>
					{/if}
				</span>
				<span class="font-mono font-semibold text-[var(--color-white)]">{formatMoney(totals.withMarkup)}</span>
			</div>
		</div>
	{/if}
</div>
