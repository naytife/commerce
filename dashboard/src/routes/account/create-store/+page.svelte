<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import { createShop, api, checkSubdomainAvailability } from '$lib/api';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import { getContext, onMount } from 'svelte';
	import { ChevronLeft, ChevronRight, CheckCircle, XCircle, Loader2 } from 'lucide-svelte';
	import type { PredefinedProductType } from '$lib/types';

	let title = '';
	let subdomain = '';
	// Slugify function for subdomain
	$: slug = subdomain
		.trim()
		.toLowerCase()
		.replace(/\s+/g, '-') // replace spaces with dashes
		.replace(/[^a-z0-9-]/g, ''); // remove invalid chars

	let loading = false;
	let currentStep = 1;
	let selectedProductTypeValue: { value: string; label: string } | undefined = undefined;
	let productTypes: PredefinedProductType[] = [];
	let loadingProductTypes = true;
	
	// Subdomain availability check state
	let checkingSubdomain = false;
	let subdomainAvailability: { available: boolean; message: string; checked: boolean } = {
		available: false,
		message: '',
		checked: false
	};

	const authFetch = getContext('authFetch') as typeof fetch;

	// Get selected product type object
	$: selectedProductType = productTypes.find(pt => pt.id === selectedProductTypeValue?.value) || null;

	// Reset subdomain availability when slug changes
	$: if (slug) {
		subdomainAvailability = { available: false, message: '', checked: false };
	}

	// Check subdomain availability
	async function checkSubdomainAvailabilityHandler() {
		if (!slug.trim()) {
			toast.error('Please enter a subdomain');
			return;
		}

		checkingSubdomain = true;
		try {
			const result = await checkSubdomainAvailability(slug, authFetch);
			subdomainAvailability = {
				available: result.available,
				message: result.message,
				checked: true
			};
		} catch (error) {
			console.error('Error checking subdomain availability:', error);
			toast.error('Failed to check subdomain availability');
			subdomainAvailability = { available: false, message: 'Failed to check availability', checked: true };
		} finally {
			checkingSubdomain = false;
		}
	}

	// Load product types on mount
	onMount(async () => {
		try {
			productTypes = await api(authFetch).getPredefinedProductTypes();
		} catch (error) {
			console.error('Failed to load product types:', error);
			toast.error('Failed to load product types');
		} finally {
			loadingProductTypes = false;
		}
	});

	function nextStep() {
		if (currentStep === 1) {
			if (!title.trim() || !slug.trim()) {
				toast.error('Please provide both a title and subdomain');
				return;
			}
			
			// Check if subdomain availability has been verified
			if (!subdomainAvailability.checked) {
				toast.error('Please check subdomain availability first');
				return;
			}
			
			// Check if subdomain is available
			if (!subdomainAvailability.available) {
				toast.error('This subdomain is not available. Please choose a different one.');
				return;
			}
		}
		currentStep++;
	}

	function prevStep() {
		currentStep--;
	}

	async function handleCreate() {
		if (!title.trim() || !slug.trim()) {
			toast.error('Please provide both a title and subdomain');
			return;
		}
		
		if (!selectedProductType) {
			toast.error('Please select a product type');
			return;
		}

		loading = true;
		try {
			const shop = await createShop({ subdomain: slug, title }, authFetch);
			
			// Create product type from template for the new shop
			try {
				const response = await authFetch(`http://127.0.0.1:8080/v1/shops/${shop.shop_id}/product-types/from-template`, {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({ template_id: selectedProductTypeValue?.value }),
				});

				if (!response.ok) {
					const errorData = await response.json();
					throw new Error(errorData.message || 'Failed to create product type');
				}
			} catch (productTypeError) {
				console.error('Failed to create product type:', productTypeError);
				toast.error('Store created, but failed to add product type. You can add it later.');
			}
			
			toast.success('Store created with your first product type!');
			goto('/account');
		} catch (error) {
			console.error(error);
			toast.error('Failed to create store');
		} finally {
			loading = false;
		}
	}
</script>



<div class="container relative flex h-svh flex-col items-center justify-center lg:px-0">
	<Card.Root class="w-full max-w-md">
		<Card.Header>
			<Card.Title class="text-2xl">Create Store</Card.Title>
			<Card.Description>
				{#if currentStep === 1}
					Enter your store details to get started.
				{:else}
					Choose what type of products you'll be selling.
				{/if}
			</Card.Description>
		</Card.Header>

		<Card.Content class="space-y-4">
			{#if currentStep === 1}
				<!-- Step 1: Store Details -->
				<div class="space-y-2">
					<Label for="title">Store Name</Label>
					<Input 
						id="title" 
						type="text" 
						placeholder="My Awesome Store" 
						required 
						bind:value={title} 
					/>
				</div>
				<div class="space-y-2">
					<Label for="subdomain">Store URL</Label>
					<div class="flex gap-2">
						<Input
							id="subdomain"
							type="text"
							placeholder="mystore"
							required
							bind:value={subdomain}
							class="flex-1"
						/>
						<Button 
							variant="outline" 
							size="sm"
							on:click={checkSubdomainAvailabilityHandler}
							disabled={!slug.trim() || checkingSubdomain}
							class="shrink-0"
						>
							{#if checkingSubdomain}
								<Loader2 class="w-4 h-4 animate-spin" />
							{:else}
								Check
							{/if}
						</Button>
					</div>
					
					{#if slug}
						<p class="text-sm text-muted-foreground">
							Your store will be available at: 
							<span class="font-medium text-foreground">{slug}.naytife.com</span>
						</p>
					{/if}
					
					<!-- Subdomain availability status -->
					{#if subdomainAvailability.checked}
						<div class="flex items-center gap-2 text-sm">
							{#if subdomainAvailability.available}
								<CheckCircle class="w-4 h-4 text-green-600" />
								<span class="text-green-600">{subdomainAvailability.message}</span>
							{:else}
								<XCircle class="w-4 h-4 text-red-600" />
								<span class="text-red-600">{subdomainAvailability.message}</span>
							{/if}
						</div>
					{/if}
				</div>
			{:else if currentStep === 2}
				<!-- Step 2: Product Type Selection -->
				{#if loadingProductTypes}
					<div class="text-center py-8">
						<div class="animate-spin rounded-full h-6 w-6 border-b-2 border-primary mx-auto mb-2"></div>
						<p class="text-sm text-muted-foreground">Loading product types...</p>
					</div>
				{:else}
					<div class="space-y-2">
						<Label for="product-type">Product Type</Label>
						<Select.Root bind:selected={selectedProductTypeValue}>
							<Select.Trigger>
								<Select.Value placeholder="Select a product type" />
							</Select.Trigger>
							<Select.Content>
								{#each productTypes as productType}
									<Select.Item value={productType.id} label="{productType.icon} {productType.title}">
										{productType.icon} {productType.title}
									</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
						{#if selectedProductType}
							<p class="text-sm text-muted-foreground">
								{selectedProductType.description}
							</p>
						{/if}
					</div>
				{/if}
			{/if}
		</Card.Content>

		<Card.Footer class="flex {currentStep === 1 ? 'justify-end' : 'justify-between'}">
			{#if currentStep === 1}
				<Button 
					on:click={nextStep} 
					disabled={!title.trim() || !slug.trim() || !subdomainAvailability.checked || !subdomainAvailability.available}
				>
					Next <ChevronRight class="ml-1 h-4 w-4" />
				</Button>
			{:else if currentStep === 2}
				<Button variant="outline" on:click={prevStep}>
					<ChevronLeft class="mr-1 h-4 w-4" /> Back
				</Button>
				<Button 
					on:click={handleCreate} 
					disabled={loading || !selectedProductTypeValue}
				>
					{loading ? 'Creating...' : 'Create Store'}
				</Button>
			{/if}
		</Card.Footer>
	</Card.Root>
</div>
