import { sveltekit } from '@sveltejs/kit/vite'
import houdini from 'houdini/vite'
import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
	plugins: [tailwindcss(),houdini(), sveltekit()],
	server: {
		host: '0.0.0.0', // Allow external access
		allowedHosts: ['gossip.naytife.dev'], // Add your local domain here
	  },
});
