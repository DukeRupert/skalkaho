/**
 * Auto-save with debounced POST.
 *
 * Save states:
 * - clean: no unsaved changes
 * - dirty: changes pending, timer running
 * - saving: POST in flight
 * - error: last save failed
 */

/**
 * Create an auto-save controller.
 * @param {string} projectId - job/project ID for the API URL
 * @param {number} debounceMs - debounce delay in milliseconds (default 2000)
 */
export function createAutoSave(projectId, debounceMs = 2000) {
	let status = $state('clean');
	let savedAt = $state(null);
	let timer = null;
	let getEstimateFn = null;

	function register(fn) {
		getEstimateFn = fn;
	}

	function markDirty() {
		status = 'dirty';
		if (timer) clearTimeout(timer);
		timer = setTimeout(() => save(), debounceMs);
	}

	async function save() {
		if (timer) {
			clearTimeout(timer);
			timer = null;
		}

		if (!getEstimateFn) return null;
		const estimate = getEstimateFn();
		if (!estimate) return null;

		status = 'saving';

		try {
			const res = await fetch(`/api/estimate/${projectId}`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(estimate),
			});

			if (!res.ok) {
				const err = await res.json().catch(() => ({ error: `HTTP ${res.status}` }));
				throw new Error(err.error || `Save failed: ${res.status}`);
			}

			const saved = await res.json();
			status = 'clean';
			savedAt = new Date();
			return saved;
		} catch (e) {
			status = 'error';
			console.error('Auto-save failed:', e.message);
			timer = setTimeout(() => save(), debounceMs * 2);
			return null;
		}
	}

	function destroy() {
		if (timer) clearTimeout(timer);
	}

	return {
		register,
		markDirty,
		save,
		destroy,
		get status() { return status; },
		get savedAt() { return savedAt; },
	};
}
