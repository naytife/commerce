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
				title: shop.title,
				email: shop.email,
				address: shop.address,
				about: shop.about,
				phone_number: shop.phone_number,
				currency_code: shop.currency_code
			};
			
			// Call the API to update the shop data
			await api(authFetch).updateShop(updateData);
			
			// Refetch shop data to update the UI
			await refetchShopData();
			initialShop = deepClone(shop); // Reset initial state after update
		} catch (error) {
			console.error('Error updating shop settings:', error);
			toast.error('Failed to update settings');
		}
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>General Store Information</Card.Title>
		<Card.Description>Manage your store details and contact information</Card.Description>
	</Card.Header>
	<Card.Content>
		<form method="POST" class="space-y-4" id="general-settings-form" on:submit={handleFormSubmit}>
			<input type="hidden" name="form-id" value="general-settings-form" />
			<div class="flex w-full flex-col gap-3">
				<Label for="store-name">Store Name</Label>
				<Input
					id="store-name"
					name="store-name"
					bind:value={shop.title}
					placeholder="Enter your shop name"
				/>
			</div>
			<div class="flex w-full flex-col gap-3">
				<Label for="email">Email</Label>
				<Input 
					type="email" 
					id="email" 
					name="email"
					bind:value={shop.email}
					placeholder="contact@yourstore.com" 
				/>
				<p class="text-muted-foreground text-sm">Store contact email address.</p>
			</div>

			<div class="flex w-full flex-col gap-3">
				<Label for="address">Address</Label>
				<Textarea 
					id="address" 
					name="address"
					bind:value={shop.address}
					placeholder="Your store's physical address" 
				/>
				<p class="text-muted-foreground text-sm">Store physical address.</p>
			</div>

			<div class="flex w-full flex-col gap-3">
				<Label for="about">About</Label>
				<Textarea 
					id="about" 
					name="about"
					bind:value={shop.about}
					placeholder="Tell customers about your store" 
				/>
				<p class="text-muted-foreground text-sm">Tell customers about your store.</p>
			</div>

			<div class="flex w-full flex-col gap-3">
				<Label for="phone">Contact Phone Number</Label>
				<Input 
					id="phone" 
					name="phone"
					type="tel" 
					bind:value={shop.phone_number}
					placeholder="+1234567890" 
				/>
				<p class="text-muted-foreground text-sm">Store contact number.</p>
			</div>
			
			<div class="flex w-full flex-col gap-3">
				<Label for="currency">Currency</Label>
				<div class="relative">
					<select 
						id="currency" 
						name="currency" 
						class="w-full h-10 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
						bind:value={shop.currency_code}
					>
						<option value="" disabled>Select currency</option>
						<option value="USD">USD - US Dollar ($)</option>
						<option value="NGN">NGN - Nigerian Naira (₦)</option>
						<option value="EUR">EUR - Euro (€)</option>
						<option value="GBP">GBP - British Pound (£)</option>
						<option value="JPY">JPY - Japanese Yen (¥)</option>
						<option value="CAD">CAD - Canadian Dollar (C$)</option>
						<option value="AUD">AUD - Australian Dollar (A$)</option>
						<option value="CHF">CHF - Swiss Franc (CHF)</option>
						<option value="CNY">CNY - Chinese Yuan (¥)</option>
						<option value="INR">INR - Indian Rupee (₹)</option>
						<option value="ZAR">ZAR - South African Rand (R)</option>
						<option value="BRL">BRL - Brazilian Real (R$)</option>
						<option value="MXN">MXN - Mexican Peso ($)</option>
						<option value="SGD">SGD - Singapore Dollar (S$)</option>
						<option value="HKD">HKD - Hong Kong Dollar (HK$)</option>
						<option value="SEK">SEK - Swedish Krona (kr)</option>
						<option value="NOK">NOK - Norwegian Krone (kr)</option>
						<option value="DKK">DKK - Danish Krone (kr)</option>
						<option value="PLN">PLN - Polish Zloty (zł)</option>
						<option value="CZK">CZK - Czech Koruna (Kč)</option>
						<option value="HUF">HUF - Hungarian Forint (Ft)</option>
						<option value="TRY">TRY - Turkish Lira (₺)</option>
						<option value="RUB">RUB - Russian Ruble (₽)</option>
						<option value="KRW">KRW - South Korean Won (₩)</option>
						<option value="THB">THB - Thai Baht (฿)</option>
						<option value="MYR">MYR - Malaysian Ringgit (RM)</option>
						<option value="IDR">IDR - Indonesian Rupiah (Rp)</option>
						<option value="PHP">PHP - Philippine Peso (₱)</option>
						<option value="VND">VND - Vietnamese Dong (₫)</option>
						<option value="EGP">EGP - Egyptian Pound (£)</option>
						<option value="MAD">MAD - Moroccan Dirham (DH)</option>
						<option value="KES">KES - Kenyan Shilling (KSh)</option>
						<option value="GHS">GHS - Ghanaian Cedi (₵)</option>
					</select>
				</div>
				<p class="text-muted-foreground text-sm">Store's primary currency.</p>
			</div>

			<Button type="submit" disabled={!hasChanges}>Update General Settings</Button>
		</form>
	</Card.Content>
</Card.Root>