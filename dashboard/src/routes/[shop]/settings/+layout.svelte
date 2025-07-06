<script lang="ts">
	import { page } from '$app/stores';
	import SettingsSidebar from '$lib/components/settings/SettingsSidebar.svelte';
	import { api } from '$lib/api';
	import { getContext, setContext } from 'svelte';
	import { useQueryClient, createQuery } from '@tanstack/svelte-query';
	import type { Shop } from '$lib/types';
	import { writable } from 'svelte/store';

	// Get authFetch from parent (set in root layout or +page)
	const authFetch = getContext<(input: RequestInfo | URL, init?: RequestInit) => Promise<Response>>('authFetch');
	setContext('authFetch', authFetch);

	// Query client for refetching
	const queryClient = useQueryClient();
	async function refetchShopData() {
		await queryClient.invalidateQueries({ queryKey: [`shop-${$page.params.shop}`] });
	}
	setContext('refetchShopData', refetchShopData);

	// Use a store for shop
	const shopStore = writable<Partial<Shop>>({});
	setContext('shop', shopStore);

	// Load shop data
	$: shopQuery = createQuery<Shop, Error>({
		queryKey: [`shop-${$page.params.shop}`],
		queryFn: () => api(authFetch).getShop(),
		enabled: !!$page.params.shop
	});
	$: if ($shopQuery.data) {
		shopStore.set({ ...$shopQuery.data });
	}
</script>

<main class="fixed top-16 left-0 right-0 bottom-0 overflow-hidden">
	<div class="h-full flex flex-col">
		<!-- Fixed header -->
		<div class="flex-shrink-0 p-4 md:p-10 pb-4 md:pb-4">
			<div class="mx-auto w-full max-w-6xl">
				<h1 class="text-3xl font-semibold">Settings</h1>
			</div>
		</div>
		
		<!-- Main content area -->
		<div class="flex-1 overflow-hidden px-4 md:px-10 pb-4 md:pb-10">
			<div class="mx-auto w-full max-w-6xl h-full flex gap-8">
				<!-- Fixed sidebar -->
				<div class="flex-shrink-0 w-[240px] lg:w-[250px]">
					<SettingsSidebar />
				</div>
				
				<!-- Scrollable content -->
				<div class="flex-1 overflow-y-auto">
					<div class="settings-content-area rounded-lg p-6 shadow-md flex flex-col gap-4">
						<slot></slot>
					</div>
				</div>
			</div>
		</div>
	</div>
</main>


