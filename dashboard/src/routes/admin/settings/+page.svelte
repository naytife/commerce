<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs';
	import { Tabs as TabsPrimitive } from 'bits-ui';
	import type { PageData } from './$types';
	import { api } from '$lib/api';
	import type { Shop } from '$lib/types';
	import { getContext } from 'svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { setContext } from 'svelte';
	
	// Import tab components
	import GeneralTab from '$lib/components/admin/settings/GeneralTab.svelte';
	import SEOTab from '$lib/components/admin/settings/SEOTab.svelte';
	import SocialTab from '$lib/components/admin/settings/SocialTab.svelte';
	import DomainTab from '$lib/components/admin/settings/DomainTab.svelte';
	import PaymentsTab from '$lib/components/admin/settings/PaymentsTab.svelte';
	import ImagesTab from '$lib/components/admin/settings/ImagesTab.svelte';

	export let data: PageData;
	const authFetch: (input: RequestInfo | URL, init?: RequestInit | undefined) => Promise<Response> = getContext('authFetch');
	
	// Make authFetch available to all child components
	setContext('authFetch', authFetch);
	
	// Get query client for refetching
	const queryClient = useQueryClient();
	
	// Function to refetch shop data
	async function refetchShopData() {
		await queryClient.invalidateQueries({
			queryKey: ['shop', 'gossip']
		});
	}
	
	// Make refetchShopData available to child components
	setContext('refetchShopData', refetchShopData);

	$: shopQuery = createQuery<Shop, Error>({
		queryKey: ['shop'],
		queryFn: () => api(authFetch).getShop()
	});

	// Initialize shop data from query
	let shop: Partial<Shop> = {};
	$: if ($shopQuery.data) {
		shop = { ...$shopQuery.data };
	}

	// Define tab values
	const tabs = [
		{ id: 'general', title: 'General' },
		{ id: 'images', title: 'Images' },
		{ id: 'seo', title: 'SEO' },
		{ id: 'social', title: 'Social Media' },
		{ id: 'domain', title: 'Domain' },
		{ id: 'payments', title: 'Payment Methods' }
	];
	let activeTab = 'general';
	
	// Currency display value
	$: currencyDisplay = shop.currency_code || "Select currency";
</script>

<div class="space-y-6">
	<h2 class="text-2xl font-bold">Store Settings</h2>
	
	<TabsPrimitive.Root value={activeTab} onValueChange={(val) => {
		if (val) activeTab = val; 
	}}>
		<Tabs.List class="grid grid-cols-6 w-full max-w-4xl">
			{#each tabs as tab}
				<Tabs.Trigger value={tab.id} class="w-full">{tab.title}</Tabs.Trigger>
			{/each}
		</Tabs.List>
		
		<!-- General Settings Tab -->
		<Tabs.Content value="general" class="mt-6">
			<GeneralTab {shop} />
		</Tabs.Content>

		<!-- Images Settings Tab -->
		<Tabs.Content value="images" class="mt-6">
			<ImagesTab {shop} />
		</Tabs.Content>
		
		<!-- SEO Settings Tab -->
		<Tabs.Content value="seo" class="mt-6">
			<SEOTab {shop} />
		</Tabs.Content>
		
		<!-- Social Media Tab -->
		<Tabs.Content value="social" class="mt-6">
			<SocialTab {shop} />
		</Tabs.Content>
		
		<!-- Domain Settings Tab -->
		<Tabs.Content value="domain" class="mt-6">
			<DomainTab {shop} />
		</Tabs.Content>

		<!-- Payment Methods Tab -->
		<Tabs.Content value="payments" class="mt-6">
			<PaymentsTab {shop} />
		</Tabs.Content>
	</TabsPrimitive.Root>
</div>
