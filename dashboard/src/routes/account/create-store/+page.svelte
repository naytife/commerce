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
	import { ChevronLeft, ChevronRight, CheckCircle, XCircle, Loader2, Eye, Palette } from 'lucide-svelte';
	import type { PredefinedProductType, StoreTemplate } from '$lib/types';

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
	let selectedTemplateValue: { value: string; label: string } | undefined = undefined;
	let productTypes: PredefinedProductType[] = [];
	let storeTemplates: StoreTemplate[] = [];
	let loadingProductTypes = true;
	let loadingTemplates = true;
	
	// Subdomain availability check state
	let checkingSubdomain = false;
	let subdomainAvailability: { available: boolean; message: string; checked: boolean } = {
		available: false,
		message: '',
		checked: false
	};

	const authFetch = getContext('authFetch') as typeof fetch;

	// Get selected template object
	$: selectedTemplate = storeTemplates.find(template => template.name === selectedTemplateValue?.value) || null;

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

	// Load data on mount
	onMount(async () => {
		// Load store templates
		try {
			storeTemplates = await api(authFetch).getTemplates();
		} catch (error) {
			console.error('Failed to load store templates:', error);
			toast.error('Failed to load store templates');
		} finally {
			loadingTemplates = false;
		}

		// Load product types
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
		} else if (currentStep === 2) {
			if (!selectedTemplateValue) {
				toast.error('Please select a store template');
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
		
		if (!selectedTemplate) {
			toast.error('Please select a store template');
			return;
		}
		
		if (!selectedProductType) {
			toast.error('Please select a product type');
			return;
		}

		loading = true;
		try {
			// Create shop with selected template
			const shop = await createShop({ 
				subdomain: slug, 
				title,
				template: selectedTemplate.name
			}, authFetch);
			
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
			
			toast.success('Store created with your selected template and first product type!');
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
				{:else if currentStep === 2}
					Choose a template for your store's design and layout.
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
				<!-- Step 2: Template Selection -->
				{#if loadingTemplates}
					<div class="text-center py-8">
						<div class="animate-spin rounded-full h-6 w-6 border-b-2 border-primary mx-auto mb-2"></div>
						<p class="text-sm text-muted-foreground">Loading store templates...</p>
					</div>
				{:else}
					<div class="space-y-4">
						<div class="space-y-2">
							<Label>Choose a Store Template</Label>
							<p class="text-sm text-muted-foreground">Select a template that best fits your store's style and needs.</p>
						</div>
						
						<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
							{#each storeTemplates as template}
								<div 
									class="relative border rounded-lg cursor-pointer transition-all duration-200 hover:border-primary hover:shadow-md {selectedTemplateValue?.value === template.name ? 'border-primary bg-primary/5 ring-2 ring-primary/20' : 'border-border'} group"
									on:click={() => selectedTemplateValue = { value: template.name, label: template.title }}
									role="button"
									tabindex="0"
									on:keydown={(e) => {
										if (e.key === 'Enter' || e.key === ' ') {
											selectedTemplateValue = { value: template.name, label: template.title };
										}
									}}
								>
									{#if selectedTemplateValue?.value === template.name}
										<div class="absolute top-2 right-2 z-10">
											<div class="bg-primary text-primary-foreground rounded-full p-1">
												<CheckCircle class="w-4 h-4" />
											</div>
										</div>
									{/if}
									
									<div class="p-4">
										{#if template.thumbnail_url}
											<div class="aspect-video w-full mb-3 rounded-md overflow-hidden bg-muted group-hover:shadow-sm transition-shadow">
												<img 
													src={template.thumbnail_url} 
													alt="{template.title} preview"
													class="w-full h-full object-cover transition-transform duration-200 group-hover:scale-105"
													loading="lazy"
												/>
											</div>
										{:else}
											<div class="aspect-video w-full mb-3 rounded-md bg-gradient-to-br from-primary/10 to-primary/20 flex items-center justify-center group-hover:from-primary/15 group-hover:to-primary/25 transition-all duration-200">
												<Palette class="w-8 h-8 text-primary" />
											</div>
										{/if}
										
										<div class="space-y-2">
											<div class="flex items-start justify-between">
												<div>
													<h4 class="font-semibold text-sm">{template.title}</h4>
													<p class="text-xs text-muted-foreground">{template.category}</p>
												</div>
												{#if template.version}
													<span class="text-xs bg-muted px-2 py-1 rounded">
														v{template.version}
													</span>
												{/if}
											</div>
											
											<p class="text-xs text-muted-foreground line-clamp-2">{template.description}</p>
											
											{#if template.features && template.features.length > 0}
												<div class="flex flex-wrap gap-1">
													{#each template.features.slice(0, 3) as feature}
														<span class="inline-flex items-center px-1.5 py-0.5 rounded text-xs bg-primary/10 text-primary">
															{feature}
														</span>
													{/each}
													{#if template.features.length > 3}
														<span class="text-xs text-muted-foreground">
															+{template.features.length - 3} more
														</span>
													{/if}
												</div>
											{/if}
											
											{#if template.preview_url}
												<div class="pt-1">
													<Button 
														variant="ghost" 
														size="sm" 
														class="h-6 px-2 text-xs gap-1" 
														on:click={(e) => {
															e.stopPropagation();
															window.open(template.preview_url, '_blank');
														}}
													>
														<Eye class="w-3 h-3" />
														Preview
													</Button>
												</div>
											{/if}
										</div>
									</div>
								</div>
							{/each}
						</div>
						
						{#if storeTemplates.length === 0}
							<div class="text-center py-8">
								<Palette class="w-12 h-12 text-muted-foreground mx-auto mb-2" />
								<p class="text-sm text-muted-foreground">No templates available</p>
							</div>
						{/if}
					</div>
				{/if}
			{:else if currentStep === 3}
				<!-- Step 3: Product Type Selection -->
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
					on:click={nextStep} 
					disabled={!selectedTemplateValue}
				>
					Next <ChevronRight class="ml-1 h-4 w-4" />
				</Button>
			{:else if currentStep === 3}
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

<style>
	.line-clamp-2 {
		display: -webkit-box;
		-webkit-box-orient: vertical;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		overflow: hidden;
	}
</style>
