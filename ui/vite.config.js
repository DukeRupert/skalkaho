import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import path from 'path';

export default defineConfig({
	plugins: [svelte()],
	build: {
		outDir: path.resolve(__dirname, '../static/estimate-builder'),
		emptyOutDir: true,
		// Library mode: single entry point, no HTML
		lib: {
			entry: path.resolve(__dirname, 'src/main.js'),
			formats: ['es'],
			fileName: 'estimate-builder',
		},
		rollupOptions: {
			output: {
				// Predictable filenames (no hash) for simple <script> include
				entryFileNames: 'estimate-builder.js',
				assetFileNames: 'estimate-builder[extname]',
			},
		},
	},
});
