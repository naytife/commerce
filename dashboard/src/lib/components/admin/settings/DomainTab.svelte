<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import { toast } from 'svelte-sonner';
	import type { Shop } from '$lib/types';
	import { api } from '$lib/api';
	import { getContext } from 'svelte';
	import { deepEqual, deepClone } from '$lib/utils/deepEqual';
	
	export let shop: Partial<Shop>;
	
	const authFetch: (input: RequestInfo | URL, init?: RequestInit | undefined) => Promise<Response> = getContext('authFetch');
	const refetchShopData: () => Promise<void> = getContext('refetchShopData');
	
	let initialShop = deepClone(shop);
	let hasChanges = false;
	
	$: hasChanges = !deepEqual(shop, initialShop);
	
	// Handle form submission
	async function handleFormSubmit(event: Event) {
		event.preventDefault();
		
		try {
			// Get only the fields that we want to update
			const updateData: Partial<Shop> = {
				custom_domain: shop.custom_domain
			};
			
			// Call the API to update the shop data
			await api(authFetch).updateShop(updateData);
			
			// Refetch shop data to update the UI
			await refetchShopData();
			initialShop = deepClone(shop); // Reset initial state after update
		} catch (error) {
			console.error('Error updating domain settings:', error);
			toast.error('Failed to update domain settings');
		}
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Domain Settings</Card.Title>
		<Card.Description>Configure your store's domain</Card.Description>
	</Card.Header>
	<Card.Content>
		<form method="POST" class="space-y-4" id="domain-settings-form" on:submit={handleFormSubmit}>
			<input type="hidden" name="form-id" value="domain-settings-form" />
			<div class="flex w-full flex-col gap-3">
				<Label for="domain">Domain</Label>
				<div class="flex items-center gap-2 border border-input rounded-md px-3 py-2 bg-muted">
					<span>{shop.subdomain || 'yourstore'}.naytife.com</span>
				</div>
				<p class="text-muted-foreground text-sm">Your store's domain name (not editable).</p>
			</div>
			
			<div class="flex w-full flex-col gap-3">
				<Label for="custom-domain">Custom Domain</Label>
				<Input 
					id="custom-domain" 
					name="custom-domain"
					bind:value={shop.custom_domain}
					placeholder="www.yourstore.com" 
				/>
				<p class="text-muted-foreground text-sm">Your store's custom domain (requires verification).</p>
			</div>
			
			<Button type="submit" disabled={!hasChanges}>Update Domain Settings</Button>
		</form>
	</Card.Content>
</Card.Root>