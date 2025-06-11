<script lang="ts">
	import ChevronLeft from 'lucide-svelte/icons/chevron-left';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import * as Tabs from '$lib/components/ui/tabs';
	import { Badge } from '$lib/components/ui/badge';
	import TemplatePreview from '$lib/components/product-types/TemplatePreview.svelte';
	import ErrorMessage from '$lib/components/ui/error-message.svelte';
	import { api } from '$lib/api';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import { onMount, getContext } from 'svelte';
	import { page } from '$app/stores';
	import type { PredefinedProductType } from '$lib/types';

	const authFetch = getContext<typeof fetch>('authFetch');
	$: shopParam = $page.params.shop;

	// Custom creation form state
	let title = '';
	let skuSubstring = '';
	let shippable = false;
	let digital = false;
	let skuError = '';
	let customLoading = false;

	// Predefined templates state
	let predefinedTypes: PredefinedProductType[] = [];
	let loading = false;
	let templatesLoading = true;
	let templatesError = false;
	let selectedTemplate: PredefinedProductType | null = null;

	// Generate SKU substring from title
	function generateSkuSubstring(input: string): string {
		if (!input) return '';
		
		const words = input.trim().split(/\s+/);
		let result = '';
		
		if (words.length === 1) {
			// Single word: use first and last character
			result = words[0].charAt(0) + (words[0].length > 1 ? words[0].charAt(words[0].length - 1) : '');
		} else {
			// Multiple words: use first character of each word (max 4)
			result = words.slice(0, 4).map(word => word.charAt(0)).join('');
		}
		
		return result.toUpperCase();
	}

	// Auto-generate SKU substring when title changes
	$: if (title && !skuSubstring) {
		skuSubstring = generateSkuSubstring(title);
	}

	// Force uppercase for SKU input
	function handleSkuInput(event: Event) {
		const input = event.target as HTMLInputElement;
		skuSubstring = input.value.toUpperCase();
	}

	function handleSkuFocus() {
		if (!skuSubstring) {
			skuSubstring = generateSkuSubstring(title);
		}
	}

	// Validate SKU substring
	function validateSkuSubstring(sku: string): boolean {
		if (!sku) {
			skuError = 'SKU substring is required';
			return false;
		}
		if (sku.length > 4) {
			skuError = 'SKU substring must be 4 characters or less';
			return false;
		}
		if (!/^[A-Z0-9]+$/.test(sku)) {
			skuError = 'SKU substring must contain only uppercase letters and numbers';
			return false;
		}
		skuError = '';
		return true;
	}

	// Custom product type creation
	async function handleCreateCustom() {
		if (!title.trim()) {
			toast.error('Please enter a product type name');
			return;
		}

		if (skuSubstring && !validateSkuSubstring(skuSubstring)) {
			toast.error('Please fix the SKU substring before saving');
			return;
		}

		customLoading = true;
		try {
			const productType = await api(authFetch).createProductType({
				title: title.trim(),
				sku_substring: skuSubstring,
				shippable,
				digital
			});
			toast.success('Product type created successfully');
			goto(`/${shopParam}/product-types/${productType.id}/edit`);
		} catch (error) {
			console.error('Failed to create product type:', error);
			toast.error('Failed to create product type');
		} finally {
			customLoading = false;
		}
	}

	// Template-based product type creation
	async function handleCreateFromTemplate(template: PredefinedProductType) {
		if (loading) return;
		loading = true;
		
		try {
			const result = await api(authFetch).createProductTypeFromTemplate(template.id);
			toast.success(`${template.title} product type created successfully`);
			goto(`/${shopParam}/product-types/${result.product_type.id}/edit`);
		} catch (error) {
			console.error('Failed to create product type from template:', error);
			toast.error('Failed to create product type from template');
		} finally {
			loading = false;
		}
	}

	// Load predefined templates
	onMount(async () => {
		await loadTemplates();
	});

	async function loadTemplates() {
		templatesLoading = true;
		templatesError = false;
		try {
			predefinedTypes = await api(authFetch).getPredefinedProductTypes();
		} catch (error) {
			console.error('Failed to load predefined product types:', error);
			templatesError = true;
			toast.error('Failed to load predefined product types');
		} finally {
			templatesLoading = false;
		}
	}

	// Group templates by category
	$: groupedTemplates = predefinedTypes.reduce((acc, template) => {
		if (!acc[template.category]) {
			acc[template.category] = [];
		}
		acc[template.category].push(template);
		return acc;
	}, {} as Record<string, PredefinedProductType[]>);
</script>

<main class="grid flex-1 items-start gap-4 p-4 sm:px-6 sm:py-0 md:gap-8">
	<div class="mx-auto grid max-w-[59rem] flex-1 auto-rows-max gap-4">
		<div class="flex items-center gap-4">
			<Button variant="outline" size="icon" class="h-7 w-7" href="/{shopParam}/product-types">
				<ChevronLeft class="h-4 w-4" />
				<span class="sr-only">Back</span>
			</Button>
			<h1 class="flex-1 shrink-0 whitespace-nowrap text-xl font-semibold tracking-tight sm:grow-0">
				Create Product Type
			</h1>
		</div>

		<Tabs.Root value="templates" class="w-full">
			<Tabs.List class="grid w-full grid-cols-2">
				<Tabs.Trigger value="templates">Choose Template</Tabs.Trigger>
				<Tabs.Trigger value="custom">Create Custom</Tabs.Trigger>
			</Tabs.List>

			<!-- Template Selection Tab -->
			<Tabs.Content value="templates" class="space-y-4">
				<Card.Root>
					<Card.Header>
						<Card.Title>Choose a Product Type Template</Card.Title>
						<Card.Description>
							Select from our pre-configured product types with optimized attributes and variations.
						</Card.Description>
					</Card.Header>
					<Card.Content>
						{#if templatesLoading}
							<div class="text-center py-8">
								<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-2"></div>
								<p class="text-muted-foreground">Loading templates...</p>
							</div>
						{:else if templatesError}
							<ErrorMessage 
								title="Failed to load templates"
								message="Unable to load predefined product type templates. Please check your connection and try again."
								onRetry={loadTemplates}
							/>
						{:else if Object.keys(groupedTemplates).length > 0}
							{#each Object.entries(groupedTemplates) as [category, templates]}
								<div class="mb-6">
									<h3 class="mb-3 text-lg font-semibold">{category}</h3>
									<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
										{#each templates as template}
											<Card.Root class="transition-all hover:shadow-md">
												<Card.Content class="p-4">
													<div class="flex items-start gap-3">
														<div class="text-2xl">{template.icon}</div>
														<div class="flex-1 space-y-2">
															<h4 class="font-semibold">{template.title}</h4>
															<p class="text-sm text-muted-foreground line-clamp-2">{template.description}</p>
															<div class="flex items-center gap-2 flex-wrap">
																<Badge variant="secondary" class="text-xs">
																	SKU: {template.sku_substring}
																</Badge>
																{#if template.shippable}
																	<Badge variant="outline" class="text-xs">Physical</Badge>
																{/if}
																{#if template.digital}
																	<Badge variant="outline" class="text-xs">Digital</Badge>
																{/if}
															</div>
															<div class="text-xs text-muted-foreground">
																{template.attributes.length} pre-configured attributes
															</div>
															<div class="flex gap-2 pt-2">
																<TemplatePreview 
																	{template} 
																	onCreateFromTemplate={handleCreateFromTemplate}
																	{loading}
																/>
																<Button 
																	size="sm" 
																	class="flex-1"
																	on:click={() => handleCreateFromTemplate(template)}
																	disabled={loading}
																>
																	{loading ? 'Creating...' : 'Create'}
																</Button>
															</div>
														</div>
													</div>
												</Card.Content>
											</Card.Root>
										{/each}
									</div>
								</div>
							{/each}
						{:else}
							<div class="text-center py-8">
								<p class="text-muted-foreground">No templates available</p>
							</div>
						{/if}
					</Card.Content>
				</Card.Root>
			</Tabs.Content>

			<!-- Custom Creation Tab -->
			<Tabs.Content value="custom" class="space-y-4">
				<Card.Root>
					<Card.Header>
						<Card.Title>Custom Product Type</Card.Title>
						<Card.Description>
							Create a custom product type from scratch. You can add attributes after creation.
						</Card.Description>
					</Card.Header>
					<Card.Content>
						<div class="grid gap-6">
							<div class="grid gap-3">
								<Label for="name">Name</Label>
								<Input id="name" type="text" class="w-full" bind:value={title} />
							</div>
							<div class="grid gap-3">
								<Label for="skuSubstring">SKU Substring</Label>
								<div class="flex flex-col gap-1">
									<div class="flex gap-2">
										<Input 
											id="skuSubstring" 
											type="text" 
											maxlength={4}
											class="w-full {skuError ? 'border-red-500' : ''}" 
											bind:value={skuSubstring}
											on:input={handleSkuInput}
											on:focus={handleSkuFocus}
											placeholder="Will auto-generate when focused"
										/>
									</div>
									{#if skuError}
										<span class="text-sm text-red-500">{skuError}</span>
									{:else}
										<span class="text-sm text-muted-foreground">
											Used as prefix for product SKUs (max 4 characters, uppercase letters and numbers only)
										</span>
									{/if}
								</div>
							</div>
							<div class="grid gap-3">
								<Label>Product Properties</Label>
								<div class="flex flex-col gap-3">
									<div class="flex items-center space-x-2">
										<Switch id="physical" bind:checked={shippable} />
										<Label for="physical">Physical product (requires shipping)</Label>
									</div>
									<div class="flex items-center space-x-2">
										<Switch id="digital" bind:checked={digital} />
										<Label for="digital">Digital product</Label>
									</div>
								</div>
							</div>
						</div>
					</Card.Content>
					<Card.Footer>
						<Button on:click={handleCreateCustom} class="w-full" disabled={customLoading}>
							{customLoading ? 'Creating...' : 'Create Custom Product Type'}
						</Button>
					</Card.Footer>
				</Card.Root>
			</Tabs.Content>
		</Tabs.Root>
	</div>
</main>
