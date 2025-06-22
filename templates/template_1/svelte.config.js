import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://kit.svelte.dev/docs/integrations#preprocessors
	// for more information about preprocessors
	preprocess: vitePreprocess(),

	kit: {
		adapter: adapter({
			// default options are shown. On some platforms
			// these options are set automatically — see below
			pages: 'build',
			assets: 'build',
			fallback: 'app.html',
			precompress: false,
			strict: false
		}),
		prerender: {
			crawl: true,
			entries: ['/', '/privacy', '/terms'],
			handleHttpError: ({ path, referrer, message, status }) => {
				// Don't fail the build for backend connection errors during prerendering
				if (message.includes('ECONNREFUSED') || message.includes('500') || message.includes('fetch failed')) {
					console.warn(`Prerender warning for ${path}: ${message}`);
					return;
				}
				// Don't fail the build for 404s on placeholder links
				if (status === 404) {
					console.warn(`Prerender warning: 404 for ${path} (linked from ${referrer || 'unknown'})`);
					return;
				}
				throw new Error(`${status} ${path}`);
			},
			handleMissingId: ({ path, id, message }) => {
				// Handle missing route IDs gracefully
				console.warn(`Prerender warning: Missing ID ${id} for ${path}: ${message}`);
			}
		}
	}
};

export default config;
