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
	import { ChevronLeft, ChevronRight, CheckCircle, XCircle, Loader2, Eye, Palette, Package } from 'lucide-svelte';
	import type { PredefinedProductType, StoreTemplate } from '$lib/types';
	import { deploymentStore } from '$lib/stores/deployment';

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
		// Set auth fetch for deployment store
		deploymentStore.setAuthFetch(authFetch);
		
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
			
			// Start deployment tracking
			deploymentStore.startDeployment(shop);
			
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
			// Redirect to the new shop instead of account page
			goto(`/${shop.subdomain}`);
		} catch (error) {
			console.error(error);
			toast.error('Failed to create store');
		} finally {
			loading = false;
		}
	}
</script>



<div class="min-h-screen bg-gradient-to-br from-background via-surface-elevated to-surface-muted relative overflow-hidden">
	<!-- Ambient background effects -->
	<div class="fixed inset-0 overflow-hidden pointer-events-none">
		<div class="absolute -top-4 -right-4 w-72 h-72 bg-gradient-to-br from-primary/20 to-accent/20 rounded-full blur-3xl animate-pulse" style="animation-delay: 0s;"></div>
		<div class="absolute top-1/3 -left-8 w-96 h-96 bg-gradient-to-br from-accent/15 to-primary/15 rounded-full blur-3xl animate-pulse" style="animation-delay: 2s;"></div>
		<div class="absolute bottom-1/4 right-1/3 w-64 h-64 bg-gradient-to-br from-secondary/10 to-primary/10 rounded-full blur-2xl animate-pulse" style="animation-delay: 4s;"></div>
	</div>
	
	<!-- Dot grid pattern -->
	<div class="fixed inset-0 opacity-[0.02] dark:opacity-[0.05] pointer-events-none" style="background-image: radial-gradient(circle at 1px 1px, currentColor 1px, transparent 0); background-size: 40px 40px;"></div>

	<div class="container relative flex min-h-screen flex-col items-center justify-start lg:px-0 py-6">
		<!-- Progress indicator -->
		<div class="mb-6 w-full max-w-2xl">
			<div class="flex items-center justify-between mb-4">
				{#each [1, 2, 3] as step}
					<div class="flex items-center {step !== 3 ? 'flex-1' : ''}">
						<div class="relative">
							<div class="w-10 h-10 rounded-full border-2 flex items-center justify-center transition-all duration-300 {
								currentStep >= step 
									? 'bg-gradient-to-br from-primary to-accent border-primary text-white shadow-brand' 
									: currentStep === step 
										? 'border-primary bg-primary/10 text-primary animate-pulse'
										: 'border-border bg-background text-muted-foreground'
							}">
								{#if currentStep > step}
									<CheckCircle class="w-5 h-5" />
								{:else}
									<span class="text-sm font-semibold">{step}</span>
								{/if}
							</div>
						</div>
						{#if step !== 3}
							<div class="flex-1 h-px mx-4 transition-all duration-300 {
								currentStep > step ? 'bg-gradient-to-r from-primary to-accent' : 'bg-border'
							}"></div>
						{/if}
					</div>
				{/each}
			</div>
			<div class="flex justify-between text-sm">
				<span class="text-center {currentStep >= 1 ? 'text-primary font-medium' : 'text-muted-foreground'}">Store Details</span>
				<span class="text-center {currentStep >= 2 ? 'text-primary font-medium' : 'text-muted-foreground'}">Choose Template</span>
				<span class="text-center {currentStep >= 3 ? 'text-primary font-medium' : 'text-muted-foreground'}">Product Type</span>
			</div>
		</div>

		<Card.Root class="w-full max-w-2xl glass backdrop-blur-xl border-border/50 shadow-glass overflow-hidden animate-scale-in">
			<Card.Header class="text-center space-y-2 pb-4">
				<div class="w-14 h-14 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center mx-auto shadow-brand animate-bounce-in">
					{#if currentStep === 1}
						<div class="w-8 h-8 bg-white rounded-xl flex items-center justify-center">
							<span class="text-primary font-bold text-lg">1</span>
						</div>
					{:else if currentStep === 2}
						<Palette class="w-8 h-8 text-white" />
					{:else}
						<Package class="w-8 h-8 text-white" />
					{/if}
				</div>
				<Card.Title class="text-2xl font-bold bg-gradient-to-r from-foreground to-foreground/80 bg-clip-text text-transparent">
					{#if currentStep === 1}
						Create Your Store
					{:else if currentStep === 2}
						Choose Your Template
					{:else}
						Select Product Type
					{/if}
				</Card.Title>
				<Card.Description class="text-base text-muted-foreground max-w-lg mx-auto">
					{#if currentStep === 1}
						Enter your store details and create a unique online presence.
					{:else if currentStep === 2}
						Pick a beautiful template that matches your brand and style.
					{:else}
						Choose what type of products you'll be selling to get started.
					{/if}
				</Card.Description>
			</Card.Header>

		<Card.Content class="space-y-6 px-8 pb-6">
			{#if currentStep === 1}
				<!-- Step 1: Store Details with enhanced design -->
				<div class="space-y-6">
					<div class="space-y-3">
						<Label for="title" class="text-base font-medium">Store Name</Label>
						<div class="relative">
							<Input 
								id="title" 
								type="text" 
								placeholder="My Awesome Store" 
								required 
								bind:value={title}
								class="form-input h-12 text-lg pl-4 pr-12"
							/>
							<div class="absolute right-4 top-1/2 -translate-y-1/2">
								{#if title.trim()}
									<CheckCircle class="w-5 h-5 text-success" />
								{/if}
							</div>
						</div>
					</div>

					<div class="space-y-3">
						<Label for="subdomain" class="text-base font-medium">Store URL</Label>
						<div class="space-y-3">
							<div class="flex gap-3">
								<div class="relative flex-1">
									<Input
										id="subdomain"
										type="text"
										placeholder="mystore"
										required
										bind:value={subdomain}
										class="form-input h-12 text-lg pl-4 pr-12"
									/>
									<div class="absolute right-4 top-1/2 -translate-y-1/2">
										{#if subdomainAvailability.checked}
											{#if subdomainAvailability.available}
												<CheckCircle class="w-5 h-5 text-success" />
											{:else}
												<XCircle class="w-5 h-5 text-destructive" />
											{/if}
										{/if}
									</div>
								</div>
								<Button 
									variant="outline" 
									size="lg"
									on:click={checkSubdomainAvailabilityHandler}
									disabled={!slug.trim() || checkingSubdomain}
									class="glass border-border/50 h-12 px-6 min-w-[120px]"
								>
									{#if checkingSubdomain}
										<Loader2 class="w-4 h-4 animate-spin mr-2" />
										Checking
									{:else}
										Check
									{/if}
								</Button>
							</div>
							
							{#if slug}
								<div class="bg-surface-elevated rounded-xl p-4 border border-border/50">
									<p class="text-sm text-muted-foreground mb-1">Your store will be available at:</p>
									<p class="font-medium text-lg text-foreground">
										<span class="text-primary">{slug}</span><span class="text-muted-foreground">.naytife.com</span>
									</p>
								</div>
							{/if}
							
							<!-- Subdomain availability status with enhanced styling -->
							{#if subdomainAvailability.checked}
								<div class="flex items-center gap-3 p-4 rounded-xl transition-all duration-300 {
									subdomainAvailability.available 
										? 'bg-success/5 border border-success/20' 
										: 'bg-destructive/5 border border-destructive/20'
								}">
									{#if subdomainAvailability.available}
										<div class="w-8 h-8 bg-success/10 rounded-full flex items-center justify-center">
											<CheckCircle class="w-4 h-4 text-success" />
										</div>
										<div>
											<p class="font-medium text-success">Available!</p>
											<p class="text-sm text-success/80">{subdomainAvailability.message}</p>
										</div>
									{:else}
										<div class="w-8 h-8 bg-destructive/10 rounded-full flex items-center justify-center">
											<XCircle class="w-4 h-4 text-destructive" />
										</div>
										<div>
											<p class="font-medium text-destructive">Not available</p>
											<p class="text-sm text-destructive/80">{subdomainAvailability.message}</p>
										</div>
									{/if}
								</div>
							{/if}
						</div>
					</div>
				</div>
			{:else if currentStep === 2}
				<!-- Step 2: Enhanced Template Selection -->
				{#if loadingTemplates}
					<div class="text-center py-8">
						<div class="w-12 h-12 bg-gradient-to-br from-primary/10 to-accent/10 rounded-2xl flex items-center justify-center mx-auto mb-3">
							<div class="animate-spin rounded-full h-6 w-6 border-b-2 border-primary"></div>
						</div>
						<p class="text-lg font-medium text-foreground mb-1">Loading Templates</p>
						<p class="text-sm text-muted-foreground">Finding the perfect designs for your store...</p>
					</div>
				{:else}
					<div class="space-y-4">
						<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
							{#each storeTemplates as template}
								<div 
									class="group relative border rounded-xl cursor-pointer transition-all duration-300 hover:shadow-lg hover:-translate-y-0.5 overflow-hidden {
										selectedTemplateValue?.value === template.name 
											? 'border-primary bg-primary/5 ring-2 ring-primary/20 shadow-brand' 
											: 'border-border glass hover:border-primary/50'
									}"
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
											<div class="w-6 h-6 bg-gradient-to-br from-primary to-accent rounded-full flex items-center justify-center shadow-brand">
												<CheckCircle class="w-4 h-4 text-white" />
											</div>
										</div>
									{/if}
									
									<!-- Large template preview image -->
									{#if template.thumbnail_url}
										<div class="aspect-[4/3] w-full bg-muted relative overflow-hidden">
											<img 
												src={template.thumbnail_url} 
												alt="{template.title} preview"
												class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
												loading="lazy"
											/>
											<!-- Subtle overlay on hover for better text readability -->
											<div class="absolute inset-0 bg-gradient-to-t from-black/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
											
											<!-- Preview button overlay -->
											{#if template.preview_url}
												<div class="absolute top-2 left-2 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
													<Button 
														variant="secondary" 
														size="sm" 
														class="h-7 px-2 text-xs gap-1.5 bg-white/90 hover:bg-white text-gray-700 border-0 shadow-sm" 
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
									{:else}
										<div class="aspect-[4/3] w-full bg-gradient-to-br from-primary/10 to-accent/10 flex items-center justify-center group-hover:from-primary/15 group-hover:to-accent/15 transition-all duration-300">
											<div class="w-12 h-12 bg-gradient-to-br from-primary to-accent rounded-xl flex items-center justify-center shadow-brand">
												<Palette class="w-6 h-6 text-white" />
											</div>
										</div>
									{/if}
									
									<!-- Compact template info -->
									<div class="p-4">
										<div class="flex items-start justify-between gap-2">
											<div class="flex-1 min-w-0">
												<h4 class="font-semibold text-base text-foreground truncate">{template.title}</h4>
												<p class="text-xs text-primary font-medium">{template.category}</p>
											</div>
											{#if template.version}
												<span class="text-xs bg-surface-elevated border border-border px-2 py-0.5 rounded text-muted-foreground flex-shrink-0">
													v{template.version}
												</span>
											{/if}
										</div>
										
										<!-- Compact feature tags -->
										{#if template.features && template.features.length > 0}
											<div class="flex flex-wrap gap-1 mt-2">
												{#each template.features.slice(0, 2) as feature}
													<span class="inline-flex items-center px-1.5 py-0.5 rounded text-xs bg-primary/10 text-primary border border-primary/20">
														{feature}
													</span>
												{/each}
												{#if template.features.length > 2}
													<span class="text-xs text-muted-foreground bg-surface-elevated px-1.5 py-0.5 rounded border border-border">
														+{template.features.length - 2}
													</span>
												{/if}
											</div>
										{/if}
									</div>
								</div>
							{/each}
						</div>
						
						{#if storeTemplates.length === 0}
							<div class="text-center py-8">
								<div class="w-12 h-12 bg-gradient-to-br from-muted to-surface-elevated rounded-2xl flex items-center justify-center mx-auto mb-3">
									<Palette class="w-6 h-6 text-muted-foreground" />
								</div>
								<h3 class="text-base font-medium text-foreground mb-1">No Templates Available</h3>
								<p class="text-sm text-muted-foreground">We're working on adding more templates for you.</p>
							</div>
						{/if}
					</div>
				{/if}
			{:else if currentStep === 3}
				<!-- Step 3: Enhanced Product Type Selection with Cards -->
				{#if loadingProductTypes}
					<div class="text-center py-8">
						<div class="w-12 h-12 bg-gradient-to-br from-primary/10 to-accent/10 rounded-2xl flex items-center justify-center mx-auto mb-3">
							<div class="animate-spin rounded-full h-6 w-6 border-b-2 border-primary"></div>
						</div>
						<p class="text-lg font-medium text-foreground mb-1">Loading Product Types</p>
						<p class="text-sm text-muted-foreground">Discovering the perfect fit for your business...</p>
					</div>
				{:else}
					<div class="space-y-4">
						<div class="grid grid-cols-1 gap-3 max-h-80 overflow-y-auto scrollbar-hide pr-2">
							{#each productTypes as productType}
								<div 
									class="group relative border rounded-xl cursor-pointer transition-all duration-300 hover:shadow-md hover:-translate-y-0.5 {
										selectedProductTypeValue?.value === productType.id 
											? 'border-primary bg-gradient-to-r from-primary/5 to-accent/5 ring-2 ring-primary/20 shadow-brand' 
											: 'border-border glass hover:border-primary/50'
									}"
									on:click={() => selectedProductTypeValue = { value: productType.id, label: productType.title }}
									role="button"
									tabindex="0"
									on:keydown={(e) => {
										if (e.key === 'Enter' || e.key === ' ') {
											selectedProductTypeValue = { value: productType.id, label: productType.title };
										}
									}}
								>
									{#if selectedProductTypeValue?.value === productType.id}
										<div class="absolute -top-1 -right-1 z-10">
											<div class="w-6 h-6 bg-gradient-to-br from-primary to-accent rounded-full flex items-center justify-center shadow-brand">
												<CheckCircle class="w-4 h-4 text-white" />
											</div>
										</div>
									{/if}
									
									<div class="p-4">
										<div class="flex items-start gap-3">
											<!-- Icon Section -->
											<div class="flex-shrink-0">
												<div class="w-12 h-12 rounded-xl bg-gradient-to-br from-primary/10 to-accent/10 flex items-center justify-center text-2xl group-hover:from-primary/15 group-hover:to-accent/15 transition-all duration-300 {
													selectedProductTypeValue?.value === productType.id ? 'bg-gradient-to-br from-primary/20 to-accent/20' : ''
												}">
													{productType.icon}
												</div>
											</div>
											
											<!-- Content Section -->
											<div class="flex-1 min-w-0">
												<div class="flex items-center justify-between mb-1">
													<h4 class="text-base font-semibold text-foreground truncate">{productType.title}</h4>
													<span class="text-xs bg-surface-elevated border border-border px-2 py-0.5 rounded text-muted-foreground whitespace-nowrap">
														{productType.category}
													</span>
												</div>
												
												<p class="text-sm text-muted-foreground mb-2 leading-relaxed line-clamp-2">{productType.description}</p>
												
												<!-- Product Features -->
												<div class="flex flex-wrap gap-1">
													{#if productType.shippable}
														<span class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-xs bg-blue-50 text-blue-700 border border-blue-200 dark:bg-blue-950 dark:text-blue-300 dark:border-blue-800">
															<svg class="w-2.5 h-2.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
																<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
															</svg>
															Shipping
														</span>
													{/if}
													{#if productType.digital}
														<span class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-xs bg-purple-50 text-purple-700 border border-purple-200 dark:bg-purple-950 dark:text-purple-300 dark:border-purple-800">
															<svg class="w-2.5 h-2.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
																<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"/>
															</svg>
															Digital
														</span>
													{/if}
													<span class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-xs bg-green-50 text-green-700 border border-green-200 dark:bg-green-950 dark:text-green-300 dark:border-green-800">
														<svg class="w-2.5 h-2.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
															<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
														</svg>
														{productType.attributes.length} Ready
													</span>
												</div>
											</div>
										</div>
									</div>
								</div>
							{/each}
						</div>
						
						{#if productTypes.length === 0}
							<div class="text-center py-12">
								<div class="w-16 h-16 bg-gradient-to-br from-muted to-surface-elevated rounded-3xl flex items-center justify-center mx-auto mb-4">
									<Package class="w-8 h-8 text-muted-foreground" />
								</div>
								<h3 class="text-lg font-medium text-foreground mb-2">No Product Types Available</h3>
								<p class="text-sm text-muted-foreground">We're working on adding more product types for you.</p>
							</div>
						{/if}

						{#if selectedProductType}
							<div class="mt-4 p-4 bg-gradient-to-r from-primary/5 to-accent/5 rounded-xl border border-primary/20">
								<div class="flex items-center gap-3">
									<div class="w-10 h-10 rounded-lg bg-gradient-to-br from-primary to-accent flex items-center justify-center text-white text-lg shadow-brand">
										{selectedProductType.icon}
									</div>
									<div class="flex-1">
										<h4 class="font-semibold text-foreground text-sm">Great Choice!</h4>
										<p class="text-xs text-muted-foreground">{selectedProductType.description}</p>
									</div>
									<div class="text-xs text-muted-foreground">
										{selectedProductType.attributes.length} attributes
									</div>
								</div>
							</div>
						{/if}
					</div>
				{/if}
			{/if}
		</Card.Content>

		<Card.Footer class="flex {currentStep === 1 ? 'justify-end' : 'justify-between'} p-6 pt-4 bg-surface-elevated/50 border-t border-border/50">
			{#if currentStep === 1}
				<Button 
					on:click={nextStep} 
					disabled={!title.trim() || !slug.trim() || !subdomainAvailability.checked || !subdomainAvailability.available}
					size="lg"
					class="btn-gradient px-8 py-3 shadow-brand h-12 min-w-[140px]"
				>
					Continue
					<ChevronRight class="ml-2 h-4 w-4" />
				</Button>
			{:else if currentStep === 2}
				<Button 
					variant="outline" 
					on:click={prevStep}
					size="lg"
					class="glass border-border/50 px-6 py-3 h-12"
				>
					<ChevronLeft class="mr-2 h-4 w-4" /> 
					Back
				</Button>
				<Button 
					on:click={nextStep} 
					disabled={!selectedTemplateValue}
					size="lg"
					class="btn-gradient px-8 py-3 shadow-brand h-12 min-w-[140px]"
				>
					Continue
					<ChevronRight class="ml-2 h-4 w-4" />
				</Button>
			{:else if currentStep === 3}
				<Button 
					variant="outline" 
					on:click={prevStep}
					size="lg"
					class="glass border-border/50 px-6 py-3 h-12"
				>
					<ChevronLeft class="mr-2 h-4 w-4" /> 
					Back
				</Button>
				<Button 
					on:click={handleCreate} 
					disabled={loading || !selectedProductTypeValue}
					size="lg"
					class="btn-gradient px-8 py-3 shadow-brand h-12 min-w-[140px] relative overflow-hidden group"
				>
					{#if loading}
						<Loader2 class="w-4 h-4 animate-spin mr-2" />
						Creating Store...
					{:else}
						<span class="relative z-10">Create Store</span>
						<div class="absolute inset-0 bg-gradient-to-r from-accent to-primary opacity-0 group-hover:opacity-20 transition-opacity"></div>
					{/if}
				</Button>
			{/if}
		</Card.Footer>
	</Card.Root>
	</div>
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
