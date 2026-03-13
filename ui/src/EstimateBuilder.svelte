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
		{ key: 'materials', label: 'Materials', bg: 'bg-blue-50', text: 'text-blue-700', border: 'border-blue-200' },
		{ key: 'labor', label: 'Labor', bg: 'bg-amber-50', text: 'text-amber-700', border: 'border-amber-200' },
		{ key: 'equipment', label: 'Equipment', bg: 'bg-purple-50', text: 'text-purple-700', border: 'border-purple-200' },
		{ key: 'subs', label: 'Subs', bg: 'bg-green-50', text: 'text-green-700', border: 'border-green-200' },
		{ key: 'other', label: 'Other', bg: 'bg-slate-100', text: 'text-slate-600', border: 'border-slate-300' },
	];
</script>

{#if loading}
	<div class="flex items-center justify-center h-64">
		<div class="text-slate-500">Loading estimate...</div>
	</div>
{:else if error}
	<div class="flex items-center justify-center h-64">
		<div class="text-red-500">{error}</div>
	</div>
{:else if estimate}
	<div class="estimate-builder pb-16">
		<!-- Topbar -->
		<div class="sticky top-0 z-10 bg-slate-900 text-white px-4 py-3 flex items-center justify-between border-b border-slate-700">
			<div>
				<h1 class="text-lg font-semibold">{estimate.project.name}</h1>
				<span class="text-xs text-slate-400">Estimate Builder</span>
			</div>
			<div class="flex items-center gap-3">
				<button
					onclick={() => { if (undoStack.undo(estimate)) handleChange(); }}
					disabled={!undoStack.canUndo}
					class="text-xs px-2 py-1 rounded bg-slate-700 hover:bg-slate-600 disabled:opacity-30 disabled:cursor-not-allowed transition-colors flex items-center gap-1"
					title="Undo (Ctrl+Z)"
				>
					<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a5 5 0 015 5v2M3 10l4-4M3 10l4 4"/>
					</svg>
					Undo
				</button>
				<SaveStatus status={autoSave.status} savedAt={autoSave.savedAt} onsave={() => autoSave.save()} />
				<span class="text-xs px-2 py-1 rounded bg-slate-700 text-slate-300 uppercase tracking-wide">
					{estimate.project.status}
				</span>
			</div>
		</div>

		<!-- Global markup toolbar -->
		<div class="bg-slate-50 border-b border-slate-200 px-4 py-2">
			<div class="flex items-center gap-4 text-sm">
				<span class="font-medium text-slate-600">Global Markup:</span>
				{#each MARKUP_TYPES as mt}
					<label class="flex items-center gap-1 px-2 py-0.5 rounded {mt.bg} {mt.text} font-mono text-xs">
						<span class="font-medium">{mt.label}</span>
						<input
							type="number"
							value={estimate.globals[`${mt.key}_markup`]}
							oninput={(e) => handleGlobalMarkup(mt.key, e)}
							step="1"
							min="0"
							class="w-12 text-right bg-transparent border-0 p-0 font-mono text-xs focus:ring-1 focus:ring-blue-400 rounded {mt.text}"
						/>
						<span>%</span>
					</label>
				{/each}
			</div>
		</div>

		<!-- Sections -->
		<div class="p-4">
			{#if estimate.sections.length === 0 && !showAddSection}
				<div class="text-center py-16 text-slate-400">
					<svg class="w-12 h-12 mx-auto mb-3 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
					</svg>
					<p class="text-lg font-medium">No sections yet</p>
					<p class="text-sm mt-1">Add a section to start building your estimate.</p>
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
						class="flex-1 px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-400 focus:border-blue-400"
						onkeydown={(e) => { if (e.key === 'Enter') addSection(); if (e.key === 'Escape') showAddSection = false; }}
					/>
					<button onclick={addSection} class="px-3 py-2 bg-slate-800 text-white text-sm rounded-lg hover:bg-slate-700">Add</button>
					<button onclick={() => showAddSection = false} class="px-3 py-2 text-slate-500 text-sm hover:text-slate-700">Cancel</button>
				</div>
			{:else}
				<button
					onclick={() => showAddSection = true}
					class="mt-3 text-sm text-blue-500 hover:text-blue-700 flex items-center gap-1"
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
