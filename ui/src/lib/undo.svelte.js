/**
 * Undo system for the estimate builder.
 * Snapshots full estimate state before structural mutations.
 * Capped at 20 entries. Ctrl+Z triggers undo.
 */

const MAX_UNDO = 20;

export function createUndoStack() {
	let stack = $state([]);
	let canUndo = $derived(stack.length > 0);

	function snapshot(estimate) {
		const snap = JSON.parse(JSON.stringify({
			globals: estimate.globals,
			sections: estimate.sections,
		}));
		stack.push(snap);
		if (stack.length > MAX_UNDO) {
			stack.shift();
		}
	}

	function undo(estimate) {
		if (stack.length === 0) return false;
		const snap = stack.pop();
		estimate.globals = snap.globals;
		estimate.sections = snap.sections;
		return true;
	}

	function clear() {
		stack.length = 0;
	}

	return {
		get canUndo() { return canUndo; },
		get depth() { return stack.length; },
		snapshot,
		undo,
		clear,
	};
}
