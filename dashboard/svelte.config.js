import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://kit.svelte.dev/docs/integrations#preprocessors
	// for more information about preprocessors
	preprocess: vitePreprocess(),

	kit: {
		adapter: adapter({ out: 'build' }),
		alias: {
		},
		// Add CSRF configuration to allow cross-site form submissions for auth
		csrf: {
			checkOrigin: false
		}
	}
};

export default config;
