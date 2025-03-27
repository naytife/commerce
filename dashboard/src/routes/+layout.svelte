<script lang="ts">
	import { writable } from 'svelte/store'
	import { setContext } from 'svelte'
	import { ModeWatcher } from 'mode-watcher'
	import '../app.css'
	import { QueryClientProvider } from '@tanstack/svelte-query'
	import type { PageData } from './$types'

	export let data: PageData

	// Create a writable store for the session to make it reactive
	const sessionStore = writable(data.session)

	// Update the store whenever data.session changes
	$: sessionStore.set(data.session)

	// Define an authenticated fetch function that uses the latest session
	const authFetch = async (url: string, options: RequestInit = {}) => {
		const session = $sessionStore
		const accessToken = session?.access_token
		if (accessToken) {
		options.headers = {
			...options.headers,
			Authorization: `Bearer ${accessToken}`,
		}
		}
		return fetch(url, options)
	}

	// Share authFetch via context
	setContext('authFetch', authFetch)
  </script>
  
  <QueryClientProvider client={data.queryClient}>
	<ModeWatcher defaultMode={'light'} />
	<slot></slot>
  </QueryClientProvider>