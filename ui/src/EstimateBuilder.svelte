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
		{ key: 'materials', label: 'Materials' },
		{ key: 'labor', label: 'Labor' },
		{ key: 'equipment', label: 'Equipment' },
		{ key: 'subs', label: 'Subs' },
		{ key: 'other', label: 'Other' },
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
	<div class="estimate-builder pb-20">
		<!-- Toolbar: undo, save status, global markup -->
		<div class="sticky top-0 z-10 border-b border-white/[0.06]" style="background: var(--color-granite);">
			<div class="flex items-center justify-between px-5 py-2">
				<div class="flex items-center gap-4">
					<button
						onclick={() => { if (undoStack.undo(estimate)) handleChange(); }}
						disabled={!undoStack.canUndo}
						class="text-xs px-2 py-1 bg-white/[0.06] hover:bg-white/[0.1] disabled:opacity-30 disabled:cursor-not-allowed transition-colors flex items-center gap-1 font-[var(--font-ui)] text-[var(--color-concrete)]"
						title="Undo (Ctrl+Z)"
					>
						<svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a5 5 0 015 5v2M3 10l4-4M3 10l4 4"/>
						</svg>
						Undo
					</button>
					<SaveStatus status={autoSave.status} savedAt={autoSave.savedAt} onsave={() => autoSave.save()} />
				</div>
				<div class="flex items-center gap-3 text-sm">
					<span class="text-[var(--color-sage)] font-[var(--font-ui)] uppercase tracking-wider text-xs font-semibold">Markup</span>
					{#each MARKUP_TYPES as mt}
						<label class="flex items-center gap-1 text-xs text-[var(--color-concrete)] font-[var(--font-ui)]">
							<span>{mt.label}</span>
							<input
								type="number"
								value={estimate.globals[`${mt.key}_markup`]}
								oninput={(e) => handleGlobalMarkup(mt.key, e)}
								step="1"
								min="0"
								class="w-10 text-right bg-transparent border-0 border-b border-white/[0.08] p-0 font-mono text-xs text-[var(--color-sunburst)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none"
							/>
							<span class="text-white/30">%</span>
						</label>
					{/each}
				</div>
			</div>
		</div>

		<!-- Sections -->
		<div class="px-5 pt-6">
			{#if estimate.sections.length === 0 && !showAddSection}
				<div class="text-center py-16 text-[var(--color-muted-text)]">
					<p class="text-lg font-[var(--font-ui)]">No sections yet</p>
					<p class="text-sm mt-1 font-[var(--font-body)]">Add a section to start building your estimate.</p>
					<button
						onclick={() => showAddSection = true}
						class="mt-4 px-4 py-2 bg-[var(--color-granite)] border border-white/[0.06] text-[var(--color-sunburst)] text-sm font-[var(--font-ui)] font-semibold uppercase tracking-wide hover:border-[var(--color-sunburst)] transition-colors"
					>
						+ Add First Section
					</button>
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
						subcontractorsDb={estimate.subcontractors_db}
					/>
				{/each}
			{/if}

			<!-- Add Section -->
			{#if showAddSection}
				<div class="flex items-center gap-2 mt-4">
					<input
						type="text"
						bind:value={newSectionName}
						placeholder="Section name (e.g. Framing, Roofing)"
						class="flex-1 px-0 py-2 bg-transparent border-0 border-b-2 border-white/[0.08] text-[var(--color-white)] text-sm font-[var(--font-body)] focus:border-[var(--color-sunburst)] focus:ring-0 focus:outline-none placeholder-white/20"
						onkeydown={(e) => { if (e.key === 'Enter') addSection(); if (e.key === 'Escape') showAddSection = false; }}
					/>
					<button onclick={addSection} class="px-4 py-2 bg-[var(--color-sunburst)] text-[var(--color-ink)] text-xs font-[var(--font-ui)] font-bold uppercase tracking-wide hover:brightness-110">Add</button>
					<button onclick={() => showAddSection = false} class="px-3 py-2 text-[var(--color-muted-text)] text-xs hover:text-[var(--color-white)] font-[var(--font-ui)]">Cancel</button>
				</div>
			{:else if estimate.sections.length > 0}
				<button
					onclick={() => showAddSection = true}
					class="mt-6 w-full py-3 border border-dashed border-white/[0.08] text-sm text-[var(--color-muted-text)] hover:text-[var(--color-sunburst)] hover:border-[var(--color-sunburst)]/30 flex items-center justify-center gap-2 font-[var(--font-ui)] uppercase tracking-wide transition-colors"
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
