<script>
	let { status, savedAt, onsave } = $props();

	let now = $state(Date.now());

	// Update relative time every 10 seconds
	$effect(() => {
		const interval = setInterval(() => { now = Date.now(); }, 10000);
		return () => clearInterval(interval);
	});

	let statusText = $derived.by(() => {
		switch (status) {
			case 'clean':
				if (savedAt) {
					const ago = Math.round((now - savedAt.getTime()) / 1000);
					if (ago < 5) return 'Saved just now';
					if (ago < 60) return `Saved ${ago}s ago`;
					return `Saved ${Math.round(ago / 60)}m ago`;
				}
				return 'Up to date';
			case 'dirty':
				return 'Unsaved changes';
			case 'saving':
				return 'Saving...';
			case 'error':
				return 'Save failed';
			default:
				return '';
		}
	});

	let statusColor = $derived.by(() => {
		switch (status) {
			case 'clean': return 'text-green-400';
			case 'dirty': return 'text-amber-400';
			case 'saving': return 'text-blue-400';
			case 'error': return 'text-red-400';
			default: return 'text-slate-400';
		}
	});

	let dotColor = $derived.by(() => {
		switch (status) {
			case 'clean': return 'bg-green-400';
			case 'dirty': return 'bg-amber-400';
			case 'saving': return 'bg-blue-400 animate-pulse';
			case 'error': return 'bg-red-400 animate-pulse';
			default: return 'bg-slate-400';
		}
	});
</script>

<div class="flex items-center gap-1.5 text-xs {statusColor}">
	<span class="w-2 h-2 rounded-full {dotColor}"></span>
	<span>{statusText}</span>
	{#if status === 'dirty' && onsave}
		<button
			onclick={onsave}
			class="ml-1 px-1.5 py-0.5 rounded bg-amber-500/20 hover:bg-amber-500/30 text-amber-300 transition-colors"
			title="Save now (Ctrl+S)"
		>
			Save
		</button>
	{/if}
	{#if status === 'error' && onsave}
		<button
			onclick={onsave}
			class="ml-1 px-1.5 py-0.5 rounded bg-red-500/20 hover:bg-red-500/30 text-red-300 transition-colors"
			title="Retry save"
		>
			Retry
		</button>
	{/if}
</div>
