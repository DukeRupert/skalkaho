import { mount } from 'svelte';
import EstimateBuilder from './EstimateBuilder.svelte';

const target = document.getElementById('estimate-root');
if (target) {
	const projectId = target.dataset.projectId;
	mount(EstimateBuilder, {
		target,
		props: { projectId },
	});
}
