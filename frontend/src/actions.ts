/** Svelte action: auto-focus and select input when mounted. */
export function autofocus(node: HTMLInputElement) {
  // Use microtask to ensure the DOM is settled
  queueMicrotask(() => {
    node.focus();
    node.select();
  });
}
