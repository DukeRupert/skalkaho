<script>
	import { resolveMarkup, afterMarkupPrice, lineItemTotal, formatMoney } from './markup.js';

	let { item, globals, markupOverrides, markupEnabled } = $props();

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
</script>

<tr class="border-b border-slate-100 hover:bg-slate-50 text-sm">
	<!-- Type badge -->
	<td class="px-2 py-1.5 w-20">
		<span class="text-xs px-1.5 py-0.5 rounded font-medium {typeColor}">
			{item.category_type}
		</span>
	</td>
	<!-- Item name -->
	<td class="px-2 py-1.5">
		<span class="text-slate-800">{item.item_name}</span>
		{#if item.description}
			<span class="block text-xs text-slate-400 mt-0.5">{item.description}</span>
		{/if}
	</td>
	<!-- Qty -->
	<td class="px-2 py-1.5 text-right font-mono w-20">{item.quantity}</td>
	<!-- Unit -->
	<td class="px-2 py-1.5 text-center text-slate-500 w-16">{item.unit}</td>
	<!-- Unit price -->
	<td class="px-2 py-1.5 text-right font-mono w-24">{formatMoney(item.unit_price)}</td>
	<!-- Markup % -->
	<td class="px-2 py-1.5 text-right font-mono text-slate-500 w-16">{markup}%</td>
	<!-- After markup -->
	<td class="px-2 py-1.5 text-right font-mono w-24">{formatMoney(markedUpPrice)}</td>
	<!-- Total -->
	<td class="px-2 py-1.5 text-right font-mono font-medium w-28">{formatMoney(total)}</td>
</tr>
