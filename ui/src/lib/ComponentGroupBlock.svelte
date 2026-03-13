<script>
	import LineItemRow from './LineItemRow.svelte';

	let { group, globals, markupOverrides, markupEnabled } = $props();
</script>

<div class="ml-4 mt-2">
	<div class="flex items-center gap-2 mb-1">
		<span class="text-xs font-semibold text-slate-500 uppercase tracking-wide">{group.name}</span>
		<span class="text-xs text-slate-400">({group.line_items.length})</span>
	</div>
	{#if group.line_items.length > 0}
		<table class="w-full">
			<tbody>
				{#each group.line_items as item (item.id)}
					<LineItemRow {item} {globals} {markupOverrides} {markupEnabled} />
				{/each}
			</tbody>
		</table>
	{/if}
</div>
