<script>
	let { status, savedAt } = $props();

	let statusText = $derived.by(() => {
		switch (status) {
			case 'clean':
				if (savedAt) {
					const ago = Math.round((Date.now() - savedAt.getTime()) / 1000);
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
				return 'Save failed — retrying';
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
			case 'error': return 'bg-red-400';
			default: return 'bg-slate-400';
		}
	});
</script>

<div class="flex items-center gap-1.5 text-xs {statusColor}">
	<span class="w-2 h-2 rounded-full {dotColor}"></span>
	{statusText}
</div>
