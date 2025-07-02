<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea';
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
				seo_title: shop.seo_title,
				seo_description: shop.seo_description,
				seo_keywords: shop.seo_keywords
			};
			
			// Call the API to update the shop data
			await api(authFetch).updateShop( updateData);
			
			// Refetch shop data to update the UI
			await refetchShopData();
			initialShop = deepClone(shop); // Reset initial state after update
		} catch (error) {
			console.error('Error updating SEO settings:', error);
			toast.error('Failed to update SEO settings');
		}
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>SEO Settings</Card.Title>
		<Card.Description>Optimize your store for search engines</Card.Description>
	</Card.Header>
	<Card.Content>
		<form method="POST" class="space-y-4" id="seo-settings-form" on:submit={handleFormSubmit}>
			<input type="hidden" name="form-id" value="seo-settings-form" />
			<div class="flex w-full flex-col gap-3">
				<Label for="seo-title">SEO Title</Label>
				<Input 
					id="seo-title" 
					name="seo-title"
					bind:value={shop.seo_title}
					placeholder="Your store SEO title" 
				/>
				<p class="text-muted-foreground text-sm">Title that appears in search engine results.</p>
			</div>
			
			<div class="flex w-full flex-col gap-3">
				<Label for="seo-description">SEO Description</Label>
				<Textarea 
					id="seo-description" 
					name="seo-description"
					bind:value={shop.seo_description}
					placeholder="Brief description of your store for search engines" 
				/>
				<p class="text-muted-foreground text-sm">Description that appears in search engine results.</p>
			</div>
			
			<div class="flex w-full flex-col gap-3">
				<Label for="seo-keywords">SEO Keywords</Label>
				<Input 
					id="seo-keywords" 
					name="seo-keywords"
					placeholder="Comma-separated keywords" 
					value={shop.seo_keywords ? shop.seo_keywords.join(', ') : ''}
					on:input={(e) => {
						shop.seo_keywords = e.currentTarget.value.split(',').map(k => k.trim()).filter(Boolean);
					}}
				/>
				<p class="text-muted-foreground text-sm">Keywords to help with search engine indexing.</p>
			</div>
			
			<Button type="submit" disabled={!hasChanges}>Update SEO Settings</Button>
		</form>
	</Card.Content>
</Card.Root>