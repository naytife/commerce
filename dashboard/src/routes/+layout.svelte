<script lang="ts">
	import '../app.css'
	import { QueryClientProvider } from '@tanstack/svelte-query'
	import { ModeWatcher } from 'mode-watcher';
	import { Toaster } from 'svelte-sonner';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	
	export let data: { queryClient: any, session?: any }
	
	// Monitor session for expiry and errors
	$: {
		if (data.session && $page.route.id && !$page.route.id.includes('login') && !$page.route.id.includes('signin')) {
			// Check for refresh token errors
			if (data.session.error === 'RefreshTokenError') {
				console.log('Session has refresh error, redirecting to login');
				goto('/login');
			}
			
			// Check for expired tokens
			if (data.session.access_token_expires && Date.now() > Number(data.session.access_token_expires) * 1000) {
				console.log('Session expired, redirecting to login');
				goto('/login');
			}
		}
	}
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