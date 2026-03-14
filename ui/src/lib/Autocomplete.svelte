<script>
	let { materialsDb = [], ratesDb = [], categoryType = 'materials', onselect, oncancel } = $props();

	let query = $state('');
	let selectedIndex = $state(0);
	let inputEl;

	// Filter results based on query and category type
	let results = $derived.by(() => {
		if (query.length < 1) return [];

		const q = query.toLowerCase();
		let pool;

		if (categoryType === 'materials') {
			pool = materialsDb.map(m => ({
				id: m.id,
				name: m.name,
				category: m.supplier || '',
				unit: m.unit,
				price: m.unit_price,
				type: 'materials',
				source: 'material',
			}));
		} else {
			pool = ratesDb
				.filter(r => {
					if (categoryType === 'labor') return r.category === 'Labor';
					if (categoryType === 'equipment') return r.category === 'Equipment Rentals';
					if (categoryType === 'subs') return r.category === 'Subcontractors';
					if (categoryType === 'other') return r.category === 'Other';
					return true;
				})
				.map(r => ({
					id: r.id,
					name: r.name,
					category: r.category,
					unit: r.unit,
					price: r.rate,
					type: categoryType,
					source: 'rate',
				}));
		}

		return pool
			.filter(item => item.name.toLowerCase().includes(q) || item.category.toLowerCase().includes(q))
			.slice(0, 15);
	});

	function handleKeydown(e) {
		if (e.key === 'ArrowDown') {
			e.preventDefault();
			selectedIndex = Math.min(selectedIndex + 1, results.length - 1);
		} else if (e.key === 'ArrowUp') {
			e.preventDefault();
			selectedIndex = Math.max(selectedIndex - 1, 0);
		} else if (e.key === 'Enter') {
			e.preventDefault();
			if (results.length > 0 && selectedIndex < results.length) {
				selectItem(results[selectedIndex]);
			} else if (query.trim()) {
				selectCustom();
			}
		} else if (e.key === 'Escape') {
			e.preventDefault();
			oncancel?.();
		}
	}

	function selectItem(item) {
		onselect?.({
			item_name: item.name,
			unit: item.unit,
			unit_price: item.price,
			category_type: item.type,
			is_custom: false,
			material_id: item.source === 'material' ? item.id : null,
		});
	}

	function selectCustom() {
		onselect?.({
			item_name: query.trim(),
			unit: 'ea',
			unit_price: 0,
			category_type: categoryType,
			is_custom: true,
			material_id: null,
		});
	}

	$effect(() => {
		selectedIndex = 0;
	});

	$effect(() => {
		inputEl?.focus();
	});
</script>

<div class="relative">
	<input
		bind:this={inputEl}
		type="text"
		bind:value={query}
		onkeydown={handleKeydown}
		placeholder="Search items or type a name..."
		class="w-full px-3 py-2 text-sm bg-white/[0.04] border border-white/[0.08] text-[var(--color-white)]
			focus:ring-1 focus:ring-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] outline-none placeholder-white/20 font-[var(--font-body)]"
	/>

	{#if results.length > 0}
		<div class="absolute top-full left-0 right-0 mt-1 border border-white/[0.08] shadow-lg z-50 max-h-64 overflow-y-auto" style="background: var(--color-granite);">
			{#each results as item, i}
				<button
					class="w-full text-left px-3 py-2 text-sm flex items-center justify-between
						hover:bg-white/[0.04] {i === selectedIndex ? 'bg-white/[0.06]' : ''} transition-colors"
					onclick={() => selectItem(item)}
					onmouseenter={() => selectedIndex = i}
				>
					<div>
						<span class="text-[var(--color-white)] font-[var(--font-body)]">{item.name}</span>
						<span class="text-xs text-white/30 ml-2 font-[var(--font-body)]">{item.category}</span>
					</div>
					<div class="flex items-center gap-2 text-xs">
						<span class="font-mono text-[var(--color-sunburst)]">${item.price.toFixed(2)}</span>
						<span class="text-white/30">/ {item.unit}</span>
					</div>
				</button>
			{/each}
		</div>
	{:else if query.length > 0}
		<div class="absolute top-full left-0 right-0 mt-1 border border-white/[0.08] shadow-lg z-50" style="background: var(--color-granite);">
			<button
				class="w-full text-left px-3 py-2 text-sm text-[var(--color-muted-text)] hover:bg-white/[0.04] font-[var(--font-body)] transition-colors"
				onclick={selectCustom}
			>
				Add custom item: <span class="font-medium text-[var(--color-white)]">"{query}"</span>
			</button>
		</div>
	{/if}
</div>
