<script lang="ts">
	import '../app.css'
	import { QueryClientProvider } from '@tanstack/svelte-query'
	import { ModeWatcher } from 'mode-watcher';
	import { Toaster } from 'svelte-sonner';
	export let data: { queryClient: any }
</script>

<QueryClientProvider client={data.queryClient}>
	<ModeWatcher defaultMode={'dark'} />
	<slot />
	<Toaster 
		theme="dark"
		position="top-right"
		closeButton={true}
		richColors={true}
		toastOptions={{
			style: 'background: var(--glass-bg); border: 1px solid var(--glass-border); backdrop-filter: blur(20px);'
		}}
	/>
</QueryClientProvider>