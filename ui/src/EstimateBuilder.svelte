<script>
	import SectionBlock from './lib/SectionBlock.svelte';
	import FooterSummary from './lib/FooterSummary.svelte';
	import SaveStatus from './lib/SaveStatus.svelte';
	import { formatPercent } from './lib/markup.js';
	import { createAutoSave } from './lib/autosave.svelte.js';
	import { createUndoStack } from './lib/undo.svelte.js';
	import { nanoid } from 'nanoid';

	let { projectId } = $props();

	let estimate = $state(null);
	let error = $state(null);
	let loading = $state(true);
	let showAddSection = $state(false);
	let newSectionName = $state('');

	// Auto-save controller
	const autoSave = createAutoSave(projectId);
	const undoStack = createUndoStack();

	async function fetchEstimate() {
		try {
			const res = await fetch(`/api/estimate/${projectId}`);
			if (!res.ok) {
				throw new Error(`Failed to load estimate: ${res.status}`);
			}
			estimate = await res.json();
			// Register the getter for auto-save
			autoSave.register(() => estimate);
		} catch (e) {
			error = e.message;
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (projectId) {
			fetchEstimate();
		}
	});

	// Expose dirty state to the DOM for external scripts
	$effect(() => {
		const root = document.getElementById('estimate-root');
		if (root) {
			root.dataset.dirty = (autoSave.status === 'dirty' || autoSave.status === 'saving') ? 'true' : 'false';
		}
	});

	// Navigation guards + Ctrl+Z undo
	$effect(() => {
		function handleBeforeUnload(e) {
			if (autoSave.status === 'dirty' || autoSave.status === 'saving') {
				e.preventDefault();
			}
		}

		function handleLinkClick(e) {
			if (autoSave.status !== 'dirty' && autoSave.status !== 'saving') return;
			const link = e.target.closest('a[href]');
			if (!link) return;
			// Only intercept same-origin navigation links
			if (link.origin !== window.location.origin) return;
			// Don't intercept links inside the estimate builder itself
			if (link.closest('#estimate-root')) return;
			e.preventDefault();
			if (confirm('You have unsaved changes. Leave without saving?')) {
				window.location.href = link.href;
			}
		}

		function handleKeydown(e) {
			if ((e.ctrlKey || e.metaKey) && e.key === 'z' && !e.shiftKey) {
				e.preventDefault();
				if (estimate && undoStack.undo(estimate)) {
					handleChange();
				}
			}
			// Ctrl+S to force save
			if ((e.ctrlKey || e.metaKey) && e.key === 's') {
				e.preventDefault();
				if (autoSave.status === 'dirty') {
					autoSave.save();
				}
			}
		}

		window.addEventListener('beforeunload', handleBeforeUnload);
		document.addEventListener('click', handleLinkClick, true);
		window.addEventListener('keydown', handleKeydown);
		return () => {
			window.removeEventListener('beforeunload', handleBeforeUnload);
			document.removeEventListener('click', handleLinkClick, true);
			window.removeEventListener('keydown', handleKeydown);
			autoSave.destroy();
		};
	});

	function handleChange() {
		autoSave.markDirty();
	}

	function takeSnapshot() {
		if (estimate) undoStack.snapshot(estimate);
	}

	function addSection() {
		takeSnapshot();
		const name = newSectionName.trim() || 'New Section';
		estimate.sections.push({
			id: nanoid(),
			name,
			sort_order: estimate.sections.length,
			subcategories: [],
		});
		newSectionName = '';
		showAddSection = false;
		handleChange();
	}

	function deleteSection(sectionId) {
		takeSnapshot();
		const idx = estimate.sections.findIndex(s => s.id === sectionId);
		if (idx !== -1) {
			estimate.sections.splice(idx, 1);
			handleChange();
		}
	}

	function handleGlobalMarkup(type, e) {
		const val = parseFloat(e.target.value);
		if (!isNaN(val) && val >= 0) {
			estimate.globals[`${type}_markup`] = val;
			handleChange();
		}
	}

	const MARKUP_TYPES = [
		{ key: 'materials', label: 'Materials', bg: 'bg-blue-900/30', text: 'text-blue-400', border: 'border-blue-800' },
		{ key: 'labor', label: 'Labor', bg: 'bg-amber-900/30', text: 'text-amber-400', border: 'border-amber-800' },
		{ key: 'equipment', label: 'Equipment', bg: 'bg-purple-900/30', text: 'text-purple-400', border: 'border-purple-800' },
		{ key: 'subs', label: 'Subs', bg: 'bg-green-900/30', text: 'text-green-400', border: 'border-green-800' },
		{ key: 'other', label: 'Other', bg: 'bg-white/[0.06]', text: 'text-[var(--color-concrete)]', border: 'border-white/[0.08]' },
	];
</script>

{#if loading}
	<div class="flex items-center justify-center h-64">
		<div class="text-[var(--color-muted-text)] font-[var(--font-body)]">Loading estimate...</div>
	</div>
{:else if error}
	<div class="flex items-center justify-center h-64">
		<div class="text-red-400 font-[var(--font-body)]">{error}</div>
	</div>
{:else if estimate}
	<div class="estimate-builder pb-16">
		<!-- Topbar -->
		<div class="sticky top-0 z-10 bg-[var(--color-ink)] text-[var(--color-white)] px-4 py-3 flex items-center justify-between border-b border-white/[0.08]">
			<div>
				<h1 class="text-lg font-semibold font-[var(--font-ui)] uppercase tracking-wide">{estimate.project.name}</h1>
				<span class="text-xs text-[var(--color-muted-text)] font-[var(--font-ui)]">Estimate Builder</span>
			</div>
			<div class="flex items-center gap-3">
				<button
					onclick={() => { if (undoStack.undo(estimate)) handleChange(); }}
					disabled={!undoStack.canUndo}
					class="text-xs px-2 py-1 rounded bg-white/[0.06] hover:bg-white/[0.1] disabled:opacity-30 disabled:cursor-not-allowed transition-colors flex items-center gap-1 font-[var(--font-ui)]"
					title="Undo (Ctrl+Z)"
				>
					<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a5 5 0 015 5v2M3 10l4-4M3 10l4 4"/>
					</svg>
					Undo
				</button>
				<SaveStatus status={autoSave.status} savedAt={autoSave.savedAt} onsave={() => autoSave.save()} />
				<span class="text-xs px-2 py-1 rounded bg-white/[0.06] text-[var(--color-concrete)] uppercase tracking-wide font-[var(--font-ui)]">
					{estimate.project.status}
				</span>
			</div>
		</div>

		<!-- Global markup toolbar -->
		<div class="bg-[var(--color-granite)] border-b border-white/[0.06] px-4 py-2">
			<div class="flex items-center gap-4 text-sm">
				<span class="font-medium text-[var(--color-concrete)] font-[var(--font-ui)] uppercase tracking-wide text-xs">Global Markup:</span>
				{#each MARKUP_TYPES as mt}
					<label class="flex items-center gap-1 px-2 py-0.5 rounded {mt.bg} {mt.text} font-mono text-xs">
						<span class="font-medium font-[var(--font-ui)]">{mt.label}</span>
						<input
							type="number"
							value={estimate.globals[`${mt.key}_markup`]}
							oninput={(e) => handleGlobalMarkup(mt.key, e)}
							step="1"
							min="0"
							class="w-12 text-right bg-transparent border-0 p-0 font-mono text-xs focus:ring-1 focus:ring-[var(--color-sunburst)] rounded {mt.text}"
						/>
						<span>%</span>
					</label>
				{/each}
			</div>
		</div>

		<!-- Sections -->
		<div class="p-4">
			{#if estimate.sections.length === 0 && !showAddSection}
				<div class="text-center py-16 text-[var(--color-muted-text)]">
					<svg class="w-12 h-12 mx-auto mb-3 text-white/20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
					</svg>
					<p class="text-lg font-medium font-[var(--font-ui)]">No sections yet</p>
					<p class="text-sm mt-1 font-[var(--font-body)]">Add a section to start building your estimate.</p>
				</div>
			{:else}
				{#each estimate.sections as section (section.id)}
					<SectionBlock
						{section}
						globals={estimate.globals}
						onchange={handleChange}
						onsnapshot={takeSnapshot}
						ondelete={deleteSection}
						materialsDb={estimate.materials_db}
						ratesDb={estimate.rates_db}
					/>
				{/each}
			{/if}

			<!-- Add Section -->
			{#if showAddSection}
				<div class="flex items-center gap-2 mt-3">
					<input
						type="text"
						bind:value={newSectionName}
						placeholder="Section name"
						class="flex-1 px-3 py-2 bg-transparent border-0 border-b-2 border-white/[0.08] text-[var(--color-white)] text-sm font-[var(--font-body)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/30"
						onkeydown={(e) => { if (e.key === 'Enter') addSection(); if (e.key === 'Escape') showAddSection = false; }}
					/>
					<button onclick={addSection} class="px-3 py-2 bg-[var(--color-sunburst)] text-[var(--color-ink)] text-sm rounded font-[var(--font-ui)] font-semibold hover:brightness-110">Add</button>
					<button onclick={() => showAddSection = false} class="px-3 py-2 text-[var(--color-muted-text)] text-sm hover:text-[var(--color-white)] font-[var(--font-ui)]">Cancel</button>
				</div>
			{:else}
				<button
					onclick={() => showAddSection = true}
					class="mt-3 text-sm text-[var(--color-sunburst)] hover:brightness-110 flex items-center gap-1 font-[var(--font-ui)]"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
					</svg>
					Add Section
				</button>
			{/if}
		</div>

		<!-- Footer summary bar -->
		<FooterSummary {estimate} />
	</div>
{/if}
