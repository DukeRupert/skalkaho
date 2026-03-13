<script>
	import { resolveMarkup, afterMarkupPrice, lineItemTotal, formatMoney } from './markup.js';

	let { item, globals, markupOverrides, markupEnabled, onchange, ondelete } = $props();

	const CATEGORY_TYPES = ['materials', 'labor', 'equipment', 'subs', 'other'];

	const TYPE_COLORS = {
		materials: 'bg-blue-100 text-blue-700',
		labor: 'bg-amber-100 text-amber-700',
		equipment: 'bg-purple-100 text-purple-700',
		subs: 'bg-green-100 text-green-700',
		other: 'bg-slate-100 text-slate-600',
	};

	let markup = $derived(resolveMarkup(item.category_type, globals, markupOverrides, markupEnabled));
	let markedUpPrice = $derived(afterMarkupPrice(item.unit_price, markup));
	let total = $derived(lineItemTotal(item.quantity, item.unit_price, markup));
	let typeColor = $derived(TYPE_COLORS[item.category_type] || TYPE_COLORS.other);

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

	function handleDelete() {
		ondelete?.(item.id);
	}
</script>

<tr class="border-b border-slate-100 hover:bg-slate-50 text-sm group">
	<!-- Type selector -->
	<td class="px-1 py-1 w-24">
		<select
			value={item.category_type}
			onchange={handleTypeChange}
			class="w-full text-xs px-1 py-1 rounded border-0 bg-transparent font-medium cursor-pointer
				focus:ring-2 focus:ring-blue-400 focus:bg-white {typeColor}"
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
			class="w-full px-1 py-0.5 text-slate-800 bg-transparent border-0 rounded
				focus:ring-2 focus:ring-blue-400 focus:bg-white"
			placeholder="Item name"
		/>
		{#if item.description}
			<span class="block text-xs text-slate-400 mt-0.5 px-1">{item.description}</span>
		{/if}
	</td>
	<!-- Qty -->
	<td class="px-1 py-1 w-20">
		<input
			type="number"
			value={item.quantity}
			oninput={handleQtyInput}
			step="any"
			min="0"
			class="w-full text-right font-mono px-1 py-0.5 bg-transparent border-0 rounded
				focus:ring-2 focus:ring-blue-400 focus:bg-white"
		/>
	</td>
	<!-- Unit -->
	<td class="px-1 py-1 w-16">
		<input
			type="text"
			value={item.unit}
			oninput={handleUnitInput}
			class="w-full text-center text-slate-500 px-1 py-0.5 bg-transparent border-0 rounded
				focus:ring-2 focus:ring-blue-400 focus:bg-white text-sm"
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
			class="w-full text-right font-mono px-1 py-0.5 bg-transparent border-0 rounded
				focus:ring-2 focus:ring-blue-400 focus:bg-white"
		/>
	</td>
	<!-- Markup % (read-only) -->
	<td class="px-2 py-1.5 text-right font-mono text-slate-400 w-16 text-xs">{markup}%</td>
	<!-- After markup (computed) -->
	<td class="px-2 py-1.5 text-right font-mono text-slate-500 w-24 text-xs">{formatMoney(markedUpPrice)}</td>
	<!-- Total (computed) -->
	<td class="px-2 py-1.5 text-right font-mono font-medium w-24">{formatMoney(total)}</td>
	<!-- Delete -->
	<td class="px-1 py-1 w-8">
		<button
			onclick={handleDelete}
			class="opacity-0 group-hover:opacity-100 text-slate-400 hover:text-red-500 transition-opacity p-0.5"
			title="Delete item"
		>
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
			</svg>
		</button>
	</td>
</tr>
