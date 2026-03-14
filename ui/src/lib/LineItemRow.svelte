<script>
	import { resolveMarkup, afterMarkupPrice, lineItemTotal, formatMoney } from './markup.js';

	let { item, globals, markupOverrides, markupEnabled, onchange, ondelete } = $props();

	const CATEGORY_TYPES = ['materials', 'labor', 'equipment', 'subs', 'other'];

	const TYPE_COLORS = {
		materials: 'bg-blue-900/30 text-blue-400',
		labor: 'bg-amber-900/30 text-amber-400',
		equipment: 'bg-purple-900/30 text-purple-400',
		subs: 'bg-green-900/30 text-green-400',
		other: 'bg-white/[0.06] text-[var(--color-concrete)]',
	};

	let markup = $derived(resolveMarkup(item.category_type, globals, markupOverrides, markupEnabled));
	let markedUpPrice = $derived(afterMarkupPrice(item.unit_price, markup));
	let total = $derived(lineItemTotal(item.quantity, item.unit_price, markup));
	let typeColor = $derived(TYPE_COLORS[item.category_type] || TYPE_COLORS.other);

	let showDescription = $state(!!item.description);

	function handleInput() {
		onchange?.();
	}

	function handleQtyInput(e) {
		const val = parseFloat(e.target.value);
		if (!isNaN(val)) {
			item.quantity = val;
			handleInput();
		}
	}

	function handlePriceInput(e) {
		const val = parseFloat(e.target.value);
		if (!isNaN(val)) {
			item.unit_price = val;
			item.price_override = true;
			handleInput();
		}
	}

	function handleNameInput(e) {
		item.item_name = e.target.value;
		handleInput();
	}

	function handleUnitInput(e) {
		item.unit = e.target.value;
		handleInput();
	}

	function handleTypeChange(e) {
		item.category_type = e.target.value;
		handleInput();
	}

	function handleDescriptionInput(e) {
		item.description = e.target.value || null;
		handleInput();
	}

	function toggleDescription() {
		showDescription = !showDescription;
		if (!showDescription) {
			// Keep existing description data, just hide the field
		}
	}

	function handleDelete() {
		ondelete?.(item.id);
	}
</script>

<tr class="border-b border-white/[0.06] hover:bg-white/[0.03] text-sm group">
	<!-- Type selector -->
	<td class="px-1 py-1 w-24">
		<select
			value={item.category_type}
			onchange={handleTypeChange}
			class="w-full text-xs px-1 py-1 rounded border-0 bg-transparent font-medium cursor-pointer
				focus:ring-2 focus:ring-[var(--color-sunburst)] {typeColor} font-[var(--font-ui)]"
		>
			{#each CATEGORY_TYPES as t}
				<option value={t}>{t}</option>
			{/each}
		</select>
	</td>
	<!-- Item name -->
	<td class="px-1 py-1">
		<input
			type="text"
			value={item.item_name}
			oninput={handleNameInput}
			class="w-full px-1 py-0.5 text-[var(--color-white)] bg-transparent border-0 rounded
				focus:ring-2 focus:ring-[var(--color-sunburst)] font-[var(--font-body)]"
			placeholder="Item name"
		/>
	</td>
	<!-- Qty -->
	<td class="px-1 py-1 w-20">
		<input
			type="number"
			value={item.quantity}
			oninput={handleQtyInput}
			step="any"
			min="0"
			class="w-full text-right font-mono px-1 py-0.5 bg-transparent border-0 rounded text-[var(--color-white)]
				focus:ring-2 focus:ring-[var(--color-sunburst)]"
		/>
	</td>
	<!-- Unit -->
	<td class="px-1 py-1 w-16">
		<input
			type="text"
			value={item.unit}
			oninput={handleUnitInput}
			class="w-full text-center text-[var(--color-muted-text)] px-1 py-0.5 bg-transparent border-0 rounded
				focus:ring-2 focus:ring-[var(--color-sunburst)] text-sm font-[var(--font-body)]"
			placeholder="ea"
		/>
	</td>
	<!-- Unit price -->
	<td class="px-1 py-1 w-24">
		<input
			type="number"
			value={item.unit_price}
			oninput={handlePriceInput}
			step="0.01"
			min="0"
			class="w-full text-right font-mono px-1 py-0.5 bg-transparent border-0 rounded text-[var(--color-white)]
				focus:ring-2 focus:ring-[var(--color-sunburst)]"
		/>
	</td>
	<!-- Markup % (read-only) -->
	<td class="px-2 py-1.5 text-right font-mono text-white/40 w-16 text-xs">{markup}%</td>
	<!-- After markup (computed) -->
	<td class="px-2 py-1.5 text-right font-mono text-[var(--color-muted-text)] w-24 text-xs">{formatMoney(markedUpPrice)}</td>
	<!-- Total (computed) -->
	<td class="px-2 py-1.5 text-right font-mono font-medium text-[var(--color-white)] w-24">{formatMoney(total)}</td>
	<!-- Actions: description toggle + delete -->
	<td class="px-1 py-1 w-16">
		<div class="flex items-center gap-0.5">
			<button
				onclick={toggleDescription}
				class="opacity-0 group-hover:opacity-100 p-0.5 transition-opacity rounded
					{showDescription || item.description ? 'opacity-100 text-[var(--color-sunburst)]' : 'text-white/30 hover:text-[var(--color-concrete)]'}"
				title="Toggle description"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z"/>
				</svg>
			</button>
			<button
				onclick={handleDelete}
				class="opacity-0 group-hover:opacity-100 text-white/30 hover:text-red-400 transition-opacity p-0.5"
				title="Delete item"
			>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
				</svg>
			</button>
		</div>
	</td>
</tr>
{#if showDescription}
	<tr class="border-b border-white/[0.06]">
		<td colspan="9" class="px-2 py-1.5">
			<div class="flex items-center gap-2 pl-6">
				<svg class="w-3 h-3 text-white/20 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a3 3 0 013 3v1"/>
				</svg>
				<input
					type="text"
					value={item.description ?? ''}
					oninput={handleDescriptionInput}
					placeholder="Add a description or note..."
					class="w-full px-2 py-1 text-xs bg-transparent border-0 border-b border-white/[0.06] text-[var(--color-concrete)]
						focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/20 font-[var(--font-body)]"
				/>
			</div>
		</td>
	</tr>
{/if}
