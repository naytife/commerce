import { sveltekit } from '@sveltejs/kit/vite'
import { defineConfig } from 'vite'

export default defineConfig({
	plugins: [sveltekit()],
	define: {
		global: 'globalThis',
		'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV || 'development'),
	},
	build: {
		target: 'es2020',
		minify: true,
		sourcemap: false,
	},
	optimizeDeps: {
		include: [
			'svelte',
			'@sveltejs/kit',
			'clsx',
			'tailwind-merge',
			'lucide-svelte'
		],
		exclude: [
			'fsevents',
			'lightningcss'
		]
	}
})