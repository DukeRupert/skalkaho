<script>
	import { resolveMarkup, afterMarkupPrice, lineItemTotal, formatMoney } from './markup.js';

	let { item, globals, markupOverrides, markupEnabled, onchange, ondelete } = $props();

	const CATEGORY_TYPES = ['materials', 'labor', 'equipment', 'subs', 'other'];

	let markup = $derived(resolveMarkup(item.category_type, globals, markupOverrides, markupEnabled));
	let markedUpPrice = $derived(afterMarkupPrice(item.unit_price, markup));
	let total = $derived(lineItemTotal(item.quantity, item.unit_price, markup));

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
	}

	function handleDelete() {
		ondelete?.(item.id);
	}
</script>

<tr class="border-b border-white/[0.04] hover:bg-white/[0.02] text-sm group">
	<!-- Type selector -->
	<td class="px-1 py-1 w-24">
		<select
			value={item.category_type}
			onchange={handleTypeChange}
			class="w-full text-xs px-1 py-1 border-0 bg-transparent font-medium cursor-pointer
				focus:ring-1 focus:ring-[var(--color-sunburst)] text-[var(--color-concrete)] font-[var(--font-ui)]"
			style="color-scheme: dark;"
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
			class="w-full px-1 py-0.5 text-[var(--color-white)] bg-transparent border-0
				focus:ring-1 focus:ring-[var(--color-sunburst)] font-[var(--font-body)]"
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
			class="w-full text-right font-mono px-1 py-0.5 bg-transparent border-0 text-[var(--color-white)]
				focus:ring-1 focus:ring-[var(--color-sunburst)]"
		/>
	</td>
	<!-- Unit -->
	<td class="px-1 py-1 w-16">
		<input
			type="text"
			value={item.unit}
			oninput={handleUnitInput}
			class="w-full text-center text-white/30 px-1 py-0.5 bg-transparent border-0
				focus:ring-1 focus:ring-[var(--color-sunburst)] text-sm font-[var(--font-body)]"
			placeholder="ea"
		/>
	</td>
	<!-- Unit price -->
	<td class="px-1 py-1 w-24">
		<input
			type="text"
			inputmode="decimal"
			value={item.unit_price.toFixed(2)}
			oninput={handlePriceInput}
			class="w-full text-right font-mono px-1 py-0.5 bg-transparent border-0 text-[var(--color-white)]
				focus:ring-1 focus:ring-[var(--color-sunburst)]"
		/>
	</td>
	<!-- Markup % (read-only) -->
	<td class="px-2 py-1.5 text-right font-mono text-white/25 w-16 text-xs">{markup}%</td>
	<!-- After markup (computed) -->
	<td class="px-2 py-1.5 text-right font-mono text-white/40 w-24 text-xs">{formatMoney(markedUpPrice)}</td>
	<!-- Total (computed) -->
	<td class="px-2 py-1.5 text-right font-mono font-medium text-[var(--color-white)] w-24">{formatMoney(total)}</td>
	<!-- Actions: description toggle + delete -->
	<td class="px-1 py-1 w-16">
		<div class="flex items-center gap-0.5">
			<button
				onclick={toggleDescription}
				class="p-0.5 transition-opacity
					{showDescription || item.description ? 'opacity-100 text-[var(--color-sunburst)]' : 'opacity-0 group-hover:opacity-100 text-white/20 hover:text-[var(--color-concrete)]'}"
				title="Toggle note"
			>
				<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z"/>
				</svg>
			</button>
			<button
				onclick={handleDelete}
				class="opacity-0 group-hover:opacity-100 text-white/20 hover:text-red-400 transition-opacity p-0.5"
				title="Delete item"
			>
				<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
				</svg>
			</button>
		</div>
	</td>
</tr>
{#if showDescription}
	<tr class="border-b border-white/[0.04]">
		<td colspan="9" class="px-2 py-1">
			<div class="flex items-center gap-2 pl-6">
				<svg class="w-3 h-3 text-white/15 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a3 3 0 013 3v1"/>
				</svg>
				<input
					type="text"
					value={item.description ?? ''}
					oninput={handleDescriptionInput}
					placeholder="Add a note..."
					class="w-full px-1 py-0.5 text-xs bg-transparent border-0 border-b border-white/[0.04] text-[var(--color-concrete)]
						focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/15 font-[var(--font-body)]"
				/>
			</div>
		</td>
	</tr>
{/if}
