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
	
	// Validation function for social links
	function validateSocialLink(link: string | undefined): boolean {
		if (!link) return true; // Empty links are valid
		const regex = /^https:\/\/[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}\/[a-zA-Z0-9._-]+/;
		return regex.test(link);
	}
	
	// Form validation state
	let errors = {
		instagram_link: false,
		facebook_link: false,
		whatsapp_link: false
	};
	
	let initialShop = deepClone(shop);
	let hasChanges = false;
	
	$: hasChanges = !deepEqual(shop, initialShop);
	
	// Handle form submission
	async function handleFormSubmit(event: Event) {
		event.preventDefault();
		
		// Reset errors
		errors = {
			instagram_link: false,
			facebook_link: false,
			whatsapp_link: false
		};
		
		// Validate social links
		let hasErrors = false;
		
		if (shop.instagram_link && !validateSocialLink(shop.instagram_link)) {
			errors.instagram_link = true;
			hasErrors = true;
		}
		
		if (shop.facebook_link && !validateSocialLink(shop.facebook_link)) {
			errors.facebook_link = true;
			hasErrors = true;
		}
		
		if (shop.whatsapp_link && !validateSocialLink(shop.whatsapp_link)) {
			errors.whatsapp_link = true;
			hasErrors = true;
		}
		
		if (hasErrors) {
			toast.error('Please check your social media links format');
			return;
		}
		
		try {
			// Get only the fields that we want to update
			const updateData: Partial<Shop> = {
				instagram_link: shop.instagram_link,
				facebook_link: shop.facebook_link,
				whatsapp_phone_number: shop.whatsapp_phone_number,
				whatsapp_link: shop.whatsapp_link
			};
			
			// Call the API to update the shop data
			await api(authFetch).updateShop( updateData);
			
			// Refetch shop data to update the UI
			await refetchShopData();
			initialShop = deepClone(shop); // Reset initial state after update
			toast.success('Social media settings updated successfully');
		} catch (error) {
			console.error('Error updating social media settings:', error);
			toast.error('Failed to update social media settings');
		}
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Social Media Settings</Card.Title>
		<Card.Description>Connect your store to social platforms</Card.Description>
	</Card.Header>
	<Card.Content>
		<form method="POST" class="space-y-4" id="social-settings-form" on:submit={handleFormSubmit}>
			<input type="hidden" name="form-id" value="social-settings-form" />
			<div class="flex w-full flex-col gap-3">
				<Label for="instagram">Instagram</Label>
				<Input 
					id="instagram" 
					name="instagram"
					bind:value={shop.instagram_link}
					placeholder="https://instagram.com/yourstorename" 
					class={errors.instagram_link ? "border-red-500" : ""}
				/>
				{#if errors.instagram_link}
					<p class="text-red-500 text-sm">Must follow format: https://site/user</p>
				{:else}
					<p class="text-muted-foreground text-sm">Your store's Instagram profile URL.</p>
				{/if}
			</div>
			
			<div class="flex w-full flex-col gap-3">
				<Label for="facebook">Facebook</Label>
				<Input 
					id="facebook" 
					name="facebook"
					bind:value={shop.facebook_link}
					placeholder="https://facebook.com/yourstorename" 
					class={errors.facebook_link ? "border-red-500" : ""}
				/>
				{#if errors.facebook_link}
					<p class="text-red-500 text-sm">Must follow format: https://site/user</p>
				{:else}
					<p class="text-muted-foreground text-sm">Your store's Facebook page URL.</p>
				{/if}
			</div>
			
			<div class="flex w-full flex-col gap-3">
				<Label for="whatsapp-number">WhatsApp Number</Label>
				<Input 
					id="whatsapp-number" 
					name="whatsapp-number"
					bind:value={shop.whatsapp_phone_number}
					placeholder="+1234567890" 
				/>
				<p class="text-muted-foreground text-sm">WhatsApp number for customer support.</p>
			</div>
			
			<div class="flex w-full flex-col gap-3">
				<Label for="whatsapp-link">WhatsApp Link</Label>
				<Input 
					id="whatsapp-link" 
					name="whatsapp-link"
					bind:value={shop.whatsapp_link}
					placeholder="https://wa.me/1234567890" 
					class={errors.whatsapp_link ? "border-red-500" : ""}
				/>
				{#if errors.whatsapp_link}
					<p class="text-red-500 text-sm">Must follow format: https://site/user</p>
				{:else}
					<p class="text-muted-foreground text-sm">Direct WhatsApp chat link.</p>
				{/if}
			</div>
			
			<Button type="submit" disabled={!hasChanges}>Update Social Media Settings</Button>
		</form>
	</Card.Content>
</Card.Root>