import { mount } from 'svelte';
import QuoteEditor from './QuoteEditor.svelte';

const target = document.getElementById('quote-editor');
if (target) {
  const jobId = target.dataset.jobId;
  if (jobId) {
    mount(QuoteEditor, { target, props: { jobId } });
  }
}
