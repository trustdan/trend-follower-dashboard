import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://svelte.dev/docs/kit/integrations
	// for more information about preprocessors
	preprocess: vitePreprocess(),

	kit: {
		// Use static adapter for generating static files
		adapter: adapter({
			pages: 'build',          // Output directory
			assets: 'build',         // Assets directory
			fallback: 'index.html',  // SPA fallback for client-side routing
			precompress: false,      // Don't generate .br/.gz files (not needed)
			strict: true             // Fail build if any page can't be prerendered
		})
	}
};

export default config;
