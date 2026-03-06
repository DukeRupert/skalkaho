import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
  plugins: [svelte({ compilerOptions: { css: 'injected' } })],
  build: {
    lib: {
      entry: 'src/main.ts',
      name: 'QuoteEditor',
      formats: ['iife'],
      fileName: () => 'quote-editor.iife.js',
    },
    outDir: '../static/js',
    emptyOutDir: false,
    sourcemap: true,
  },
});
