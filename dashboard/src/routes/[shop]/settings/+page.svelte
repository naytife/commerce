<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs';
	import { Tabs as TabsPrimitive } from 'bits-ui';
	import type { PageData } from './$types';
	import { api } from '$lib/api';
	import type { Shop } from '$lib/types';
	import { getContext } from 'svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { setContext } from 'svelte';
	import { page } from '$app/stores';
	// Import tab components
	import GeneralTab from '$lib/components/admin/settings/GeneralTab.svelte';
	import SEOTab from '$lib/components/admin/settings/SEOTab.svelte';
	import SocialTab from '$lib/components/admin/settings/SocialTab.svelte';
	import DomainTab from '$lib/components/admin/settings/DomainTab.svelte';
	import PaymentsTab from '$lib/components/admin/settings/PaymentsTab.svelte';
	import ImagesTab from '$lib/components/admin/settings/ImagesTab.svelte';
	import PublishButton from '$lib/components/PublishButton.svelte';
	import { 
		Settings, 
		Image, 
		Search, 
		Share2, 
		Globe, 
		CreditCard,
		Store,
		Palette,
		Shield
	} from 'lucide-svelte';

	// export let data: PageData; // Not currently used but required by SvelteKit
	const authFetch: (input: RequestInfo | URL, init?: RequestInit | undefined) => Promise<Response> = getContext('authFetch');
	
	// Make authFetch available to all child components
	setContext('authFetch', authFetch);
	
	// Get query client for refetching
	const queryClient = useQueryClient();
	
	// Function to refetch shop data
	async function refetchShopData() {
		await queryClient.invalidateQueries({
			queryKey: [`shop-${$page.params.shop}`]
		});
	}
	
	// Make refetchShopData available to child components
	setContext('refetchShopData', refetchShopData);

	$: shopQuery = createQuery<Shop, Error>({
		queryKey: [`shop-${$page.params.shop}`],
		queryFn: () => api(authFetch).getShop(),
		enabled: !!$page.params.shop
	});

	// Initialize shop data from query
	let shop: Partial<Shop> = {};
	$: if ($shopQuery.data) {
		shop = { ...$shopQuery.data };
	}

	// Define tab values with icons
	const tabs = [
		{ id: 'general', title: 'General', icon: Settings },
		{ id: 'images', title: 'Images', icon: Image },
		{ id: 'seo', title: 'SEO', icon: Search },
		{ id: 'social', title: 'Social Media', icon: Share2 },
		{ id: 'domain', title: 'Domain', icon: Globe },
		{ id: 'payments', title: 'Payment Methods', icon: CreditCard }
	];
	let activeTab = 'general';
	
	// Currency display value
	$: currencyDisplay = shop.currency_code || "Select currency";
</script>

<div class="space-y-8">
	<!-- Enhanced Header Section -->
	<div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
		<div>
			<h1 class="text-3xl font-bold text-foreground mb-2">Store Settings</h1>
			<p class="text-muted-foreground">
				Configure your store preferences, branding, and integrations
			</p>
		</div>
		<div class="flex items-center gap-3">
			<PublishButton />
			<div class="w-12 h-12 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center shadow-brand">
				<Store class="w-6 h-6 text-white" />
			</div>
		</div>
	</div>
	
	<!-- Enhanced Settings Tabs -->
	<div class="card-elevated">
		<TabsPrimitive.Root value={activeTab} onValueChange={(val) => {
			if (val) activeTab = val; 
		}}>
			<!-- Responsive tabs with horizontal scroll on mobile -->
			<Tabs.List class="flex w-full overflow-x-auto glass border-border/50 scrollbar-hide">
				{#each tabs as tab}
					<Tabs.Trigger 
						value={tab.id} 
						class="flex items-center gap-2 whitespace-nowrap data-[state=active]:bg-primary/10 data-[state=active]:text-primary py-3 px-4 min-w-fit flex-shrink-0"
					>
						<svelte:component this={tab.icon} class="w-4 h-4" />
						<span class="text-sm font-medium">{tab.title}</span>
					</Tabs.Trigger>
				{/each}
			</Tabs.List>
			
			<!-- Tab Content with consistent spacing -->
			<div class="p-6">
				<!-- General Settings Tab -->
				<Tabs.Content value="general" class="space-y-6">
					<GeneralTab {shop} />
				</Tabs.Content>

				<!-- Images Settings Tab -->
				<Tabs.Content value="images" class="space-y-6">
					<ImagesTab {shop} />
				</Tabs.Content>
				
				<!-- SEO Settings Tab -->
				<Tabs.Content value="seo" class="space-y-6">
					<SEOTab {shop} />
				</Tabs.Content>
				
				<!-- Social Media Tab -->
				<Tabs.Content value="social" class="space-y-6">
					<SocialTab {shop} />
				</Tabs.Content>
				
				<!-- Domain Settings Tab -->
				<Tabs.Content value="domain" class="space-y-6">
					<DomainTab {shop} />
				</Tabs.Content>

				<!-- Payment Methods Tab -->
				<Tabs.Content value="payments" class="space-y-6">
					<PaymentsTab />
				</Tabs.Content>
			</div>
		</TabsPrimitive.Root>
	</div>
</div>
