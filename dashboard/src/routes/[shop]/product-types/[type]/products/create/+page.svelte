<!-- Product creation page under product type -->
<script lang="ts">
	import Trash2 from 'lucide-svelte/icons/trash-2';
	import ChevronLeft from 'lucide-svelte/icons/chevron-left';
	import Upload from 'lucide-svelte/icons/upload';

	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Form from '$lib/components/ui/form';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import * as Table from '$lib/components/ui/table';
	import { Textarea } from '$lib/components/ui/textarea';
	import * as ToggleGroup from '$lib/components/ui/toggle-group';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { CirclePlus } from 'lucide-svelte';
	import { getContext } from 'svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { api } from '$lib/api';
	import type { ProductTypeAttribute, AttributeOption, ProductCreatePayload, ProductAttribute, ProductVariant, ProductType, Shop } from '$lib/types';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { getCurrencySymbol, formatAsCurrency, parseCurrencyInput } from '$lib/utils/currency';
	interface PageData {
		typeId: string;
	}

	export let data: PageData;
	const authFetch = getContext<typeof fetch>('authFetch');
	const queryClient = useQueryClient();

	// Get shop currency
	const shopQuery = createQuery<Shop, Error>({
		queryKey: [`shop-${$page.params.shop}`],
		queryFn: () => api(authFetch).getShop(),
		enabled: !!$page.params.shop
	});
	$: currencyCode = $shopQuery.data?.currency_code || 'USD';
	$: currencySymbol = getCurrencySymbol(currencyCode);

	// Get product type details
	const productTypeQuery = createQuery<ProductType>({
		queryKey: ['product-type', Number(data.typeId)],
		queryFn: () => api(authFetch).getProductTypeById(Number(data.typeId))
	});

	const typeAttributesQuery = createQuery<ProductTypeAttribute[]>({
		queryKey: ['product-type-attributes', Number(data.typeId)],
		queryFn: () => api(authFetch).getProductTypeAttributes(Number(data.typeId))
	});

	// Form state
	let title = '';
	let description = '';
	let attributes: ProductAttribute[] = [];
	let variants: ProductVariant[] = [];

	// For storing removed variants data
	let removedVariants: ProductVariant[] = [];

	// Helper function to create a ProductAttribute
	function createProductAttribute(attribute_id: number, value: string, option_id?: number): ProductAttribute {
		const attr: ProductAttribute = { 
			attribute_id, 
			value 
		};
		
		if (option_id !== undefined) {
			attr.attribute_option_id = option_id;
		}
		
		return attr;
	}
	
	// Generate all combinations of variant attributes based on selected options
	function generateCombinations(attributes: ProductTypeAttribute[]): { title: string; value: string }[][] {
		if (attributes.length === 0) return [[]];
		const [first, ...rest] = attributes;
		const combinations = generateCombinations(rest);
		
		// Get selected options for this attribute, or all options if none selected
		const selectedOptionsForAttr = selectedOptions[first.attribute_id] || [];
		const options: AttributeOption[] = Array.isArray(first.options) ? first.options : [];
		
		// Filter options based on selection
		const filteredOptions = selectedOptionsForAttr.length > 0 
			? options.filter(option => selectedOptionsForAttr.includes(option.value))
			: options;
		
		if (filteredOptions.length === 0) {
			// If no options, treat as a single option with empty value to keep combinations
			return combinations.map(combination => [{ title: first.title, value: '' }, ...combination]);
		}
		return filteredOptions.flatMap(option =>
			combinations.map(combination => [{ title: first.title, value: option.value }, ...combination])
		);
	}

	$: variantAttributes = $typeAttributesQuery.data?.filter((attr) => attr.applies_to === 'ProductVariation') || [];
	
	// State for generated combinations
	let combinations: { title: string; value: string }[][] = [];
	
	$: if (variationsConfigured && variantAttributes.length > 0) {
		console.log('Generating combinations for attributes:', variantAttributes, 'with selected options:', selectedOptions);
		combinations = generateCombinations(variantAttributes);
		console.log('Generated combinations:', combinations);
	} else {
		combinations = [];
	}

	// Initialize variation setup when variant attributes are loaded
	$: {
		if (variantAttributes.length > 0 && !variationsConfigured && Object.keys(selectedOptions).length === 0) {
			// Initialize selectedOptions with empty arrays for each variant attribute
			const initialOptions: Record<number, string[]> = {};
			variantAttributes.forEach(attr => {
				initialOptions[attr.attribute_id] = [];
			});
			selectedOptions = initialOptions;
			
			// Show variation setup if there are variant attributes
			showVariationSetup = variantAttributes.length > 0;
		}
	}

	// Generate a unique key for a variant based on its attribute values
	function getVariantKey(variant: ProductVariant): string {
		if (!variant.attributes || variant.attributes.length === 0) return '';
		
		return variant.attributes
			.map(attr => `${attr.attribute_id}:${attr.value}`)
			.sort()
			.join('|');
	}

	// Generate a unique key for a combination
	function getCombinationKey(combination: { title: string; value: string }[]): string {
		if (!combination || combination.length === 0) return '';
		
		return combination
			.map(attr => {
				const variantAttr = variantAttributes.find((va: ProductTypeAttribute) => va.title === attr.title);
				if (!variantAttr) return '';
				return `${variantAttr.attribute_id}:${attr.value}`;
			})
			.filter(key => key !== '')
			.sort()
			.join('|');
	}

	// Helper function to generate a unique key for a variation
	function getVariationKey(combination: { title: string; value: string }[]): string {
		return combination.map(attr => `${attr.title}:${attr.value}`).sort().join('|');
	}

	// Track disabled state for each variation
	let disabledVariations: any[] = [];

	// Separate lists for active and disabled variations
	let activeCombinations: { title: string; value: string }[][] = [];
	let disabledCombinations: { title: string; value: string }[][] = [];

	// Track which variations we've already seen (by unique key)
	let seenVariations = new Set<string>();

	// Variation selection flow state
	let showVariationSetup = false;
	let selectedOptions: Record<number, string[]> = {}; // attribute_id -> selected option values
	let variationsConfigured = false;

	// Update active combinations and initialize variants when combinations change
	$: if (combinations && combinations.length > 0 && variationsConfigured && activeCombinations.length === 0) {
		console.log('Initializing active combinations and variants from combinations:', combinations.length);
		
		// Clear previous state
		seenVariations.clear();
		
		// Initialize active combinations with all combinations
		// Filter out any duplicates based on attribute values
		const newActiveCombinations = combinations.filter((combination: { title: string; value: string }[]) => {
			const key = getVariationKey(combination);
			if (seenVariations.has(key)) {
				return false;
			}
			seenVariations.add(key);
			return true;
		});
		
		activeCombinations = newActiveCombinations;
		
		// Initialize variants without SKUs - backend will generate them
		variants = activeCombinations.map((combination, index) => ({
			attributes: combination.map(attr => {
				const variantAttr = variantAttributes.find((va: ProductTypeAttribute) => va.title === attr.title);
				let attribute_option_id;
				let attribute_id = 0;
				if (variantAttr) {
					attribute_id = variantAttr.attribute_id;
					if (variantAttr.options) {
						const option = variantAttr.options.find((o: AttributeOption) => o.value === attr.value);
						if (option) {
							attribute_option_id = option.attribute_option_id;
						}
					}
				}
				return createProductAttribute(
					attribute_id,
					attr.value,
					attribute_option_id
				);
			}),
			available_quantity: 1,
			description: '',
			is_default: index === 0, // Make first variant the default
			price: 0.00,
			seo_description: '',
			seo_keywords: [],
			seo_title: ''
			// No SKU or ID - will be generated by the backend
		}));
		
		console.log('Created variants:', variants.length);
	}

	// Function to check if a combination exists in a variant
	function isCombinationInVariant(combination: { title: string; value: string }[], variant: ProductVariant): boolean {
		if (!combination || combination.length === 0 || !variant || !variant.attributes) {
			return false;
		}

		// Convert combination to a map for easier lookup
		const combMap = new Map();
		combination.forEach(combAttr => {
			const variantAttr = variantAttributes.find((attr: ProductTypeAttribute) => attr.title === combAttr.title);
			if (variantAttr) {
				combMap.set(variantAttr.attribute_id, combAttr.value);
			}
		});

		// Check if variant attributes match
		for (const [attrId, combValue] of combMap.entries()) {
			const matchingVariantAttr = variant.attributes.find(attr => attr.attribute_id === attrId);
			if (!matchingVariantAttr || matchingVariantAttr.value !== combValue) {
				return false;
			}
		}

		return true;
	}

	// Function to toggle variation state
	function toggleVariation(index: number, isDisabled = false) {
		if (isDisabled) {
			// Move from disabled to active
			const variation = disabledCombinations[index];
			disabledCombinations.splice(index, 1);
			
			// Get the combination key
			const combinationKey = getCombinationKey(variation);
			
			// Check if there's a previously removed variant with the same attributes
			const previousVariant = removedVariants.find(v => getVariantKey(v) === combinationKey);
			
			if (previousVariant) {
				// Use the previously removed variant
				variants = [...variants, previousVariant];
				// Remove it from the removed variants array
				removedVariants = removedVariants.filter(v => getVariantKey(v) !== combinationKey);
			} else {
				// Create a new variant for this combination
				const newVariant = {
					attributes: variation.map(attr => {
						const variantAttr = variantAttributes.find((va: ProductTypeAttribute) => va.title === attr.title);
						let attribute_option_id;
						let attribute_id = 0;
						if (variantAttr) {
							attribute_id = variantAttr.attribute_id;
							if (variantAttr.options) {
								const option = variantAttr.options.find((o: AttributeOption) => o.value === attr.value);
								if (option) {
									attribute_option_id = option.attribute_option_id;
								}
							}
						}
						return createProductAttribute(
							attribute_id,
							attr.value,
							attribute_option_id
						);
					}),
					available_quantity: 1,
					description: '',
					is_default: variants.length === 0, // Make first variant the default
					price: 0.00,
					seo_description: '',
					seo_keywords: [],
					seo_title: ''
					// No SKU or ID - will be generated by the backend
				};
				
				variants = [...variants, newVariant];
			}
			
			activeCombinations = [...activeCombinations, variation];
		} else {
			// Move from active to disabled
			const combination = activeCombinations[index];
			
			// Find the variant that matches this combination
			const variantIndex = variants.findIndex(variant => 
				isCombinationInVariant(combination, variant)
			);
			
			if (variantIndex !== -1) {
				// Store the removed variant for potential later use
				removedVariants = [...removedVariants, variants[variantIndex]];
				// Remove the variant
				variants = variants.filter((_, i) => i !== variantIndex);
			}
			
			// Remove the combination from active and add to disabled
			activeCombinations.splice(index, 1);
			disabledCombinations = [...disabledCombinations, combination];
		}
	}

	// Handle option selection for variation attributes
	function toggleOption(attributeId: number, optionValue: string) {
		const currentSelected = selectedOptions[attributeId] || [];
		const index = currentSelected.indexOf(optionValue);
		
		if (index === -1) {
			// Add option
			selectedOptions[attributeId] = [...currentSelected, optionValue];
		} else {
			// Remove option
			selectedOptions[attributeId] = currentSelected.filter(value => value !== optionValue);
		}
		selectedOptions = { ...selectedOptions };
	}

	// Check if an option is selected
	function isOptionSelected(attributeId: number, optionValue: string): boolean {
		return selectedOptions[attributeId]?.includes(optionValue) || false;
	}

	// Get total number of possible combinations
	function getTotalCombinations(): number {
		let total = 1;
		variantAttributes.forEach(attr => {
			const selectedCount = selectedOptions[attr.attribute_id]?.length || 0;
			if (selectedCount > 0) {
				total *= selectedCount;
			} else {
				// If no options selected for any attribute, total should be 0
				total = 0;
			}
		});
		return total;
	}

	// Complete variation setup and generate combinations
	function completeVariationSetup() {
		// Validate that at least one option is selected for each attribute
		const hasValidSelections = variantAttributes.every(attr => {
			const selected = selectedOptions[attr.attribute_id];
			return selected && selected.length > 0;
		});

		if (!hasValidSelections) {
			alert('Please select at least one option for each variation attribute.');
			return;
		}

		// Reset combinations state first
		activeCombinations = [];
		disabledCombinations = [];
		variants = [];
		seenVariations.clear();
		
		// Then set the flags to trigger reactive updates
		variationsConfigured = true;
		showVariationSetup = false;
	}

	// Reset variation setup
	function resetVariationSetup() {
		selectedOptions = {};
		variationsConfigured = false;
		showVariationSetup = true;
		activeCombinations = [];
		disabledCombinations = [];
		variants = [];
		seenVariations.clear();
		
		// Re-initialize selectedOptions
		if (variantAttributes.length > 0) {
			const initialOptions: Record<number, string[]> = {};
			variantAttributes.forEach(attr => {
				initialOptions[attr.attribute_id] = [];
			});
			selectedOptions = initialOptions;
		}
	}

	// Handle Save Product
	async function saveProduct() {
		// Check if variations need to be configured
		if (variantAttributes.length > 0 && !variationsConfigured) {
			alert('Please configure variations before saving the product.');
			return;
		}

		const productPayload: ProductCreatePayload = {
			title,
			description,
			attributes,
			variants: variants.map(variant => {
				// Create a clean variant object without id and sku
				const { id, sku, ...cleanVariant } = variant as any;
				return cleanVariant;
			}),
		};

		try {
			const createdProduct = await api(authFetch).createProduct(Number(data.typeId), productPayload);
			console.log('Product created:', createdProduct);
			// Use direct browser navigation instead of goto for more reliable redirection
			if (createdProduct && createdProduct.id) {
				// First try goto
				try {
					await goto(`/${$page.params.shop}/product-types/${$page.params.type}/products/${createdProduct.id}/`);
				} catch (error) {
					console.error('Error with goto navigation:', error);
					// Fallback to direct browser navigation
					window.location.href = `/${$page.params.shop}/product-types/${$page.params.type}/products/${createdProduct.id}/`;
				}
			}
		} catch (error) {
			console.error('Failed to create product:', error);
		}
	}

	onMount(async () => {
		await queryClient.prefetchQuery({
			queryKey: [`shop-${$page.params.shop}-product-type-attributes`, $page.params.type],
			queryFn: () => api(authFetch).getProductTypeAttributes(parseInt($page.params.type)),
		});
	});
</script>

<main class="grid flex-1 items-start gap-4 p-4 sm:px-6 sm:py-0 md:gap-8">
	<div class="mx-auto grid max-w-[59rem] flex-1 auto-rows-max gap-4">
		<div class="flex items-center gap-4">
			<Button variant="outline" size="icon" class="h-7 w-7" href={`/${$page.params.shop}/product-types`}>
				<ChevronLeft class="h-4 w-4" />
				<span class="sr-only">Back</span>
			</Button>
			<h1 class="flex-1 shrink-0 whitespace-nowrap text-xl font-semibold tracking-tight sm:grow-0">
				Create Product
			</h1>
			<div class="hidden items-center gap-2 md:ml-auto md:flex">
				<Button variant="outline" size="sm">Discard</Button>
				<Button size="sm" on:click={saveProduct}>Save Product</Button>
			</div>
		</div>
		<div class="grid gap-4 md:grid-cols-[1fr_250px] lg:grid-cols-3 lg:gap-8">
			<div class="grid auto-rows-max items-start gap-4 lg:col-span-2 lg:gap-8">
				<Card.Root>
					<Card.Header>
						<Card.Title>Product Details</Card.Title>
						<Card.Description>Enter general product information</Card.Description>
					</Card.Header>
					<Card.Content>
						<div class="grid gap-6">
							<div class="grid gap-3">
								<Label for="name">Name</Label>
								<Input id="name" type="text" class="w-full" bind:value={title} />
							</div>
							<div class="grid gap-3">
								<Label for="description">Description</Label>
								<Textarea id="description" class="min-h-32" bind:value={description} />
							</div>
						</div>
					</Card.Content>
				</Card.Root>

				<Card.Root>
					<Card.Header>
						<Card.Title>Attributes</Card.Title>
						<Card.Description>Product specific attributes</Card.Description>
					</Card.Header>
					<Card.Content>
						<div>
							<form class="grid gap-4">
								{#if $typeAttributesQuery.data}
									{#each $typeAttributesQuery.data as attribute}
										{#if attribute.applies_to === 'Product'}
											<div class="flex w-full flex-col gap-3">
												<Label for={attribute.title}>{attribute.title}</Label>
												{#if attribute.options}
													{#key attribute.attribute_id}
														<Select.Root>
															<Select.Trigger id={attribute.title} aria-label={attribute.title}>
																<Select.Value placeholder={`Select ${attribute.title}`}>
																	{#if attribute.data_type === 'Color' && attributes.find(attr => attr.attribute_id === attribute.attribute_id)?.value}
																		<div class="flex items-center gap-2">
																			<div class="h-4 w-4 rounded-full" style="background-color: {attributes.find(attr => attr.attribute_id === attribute.attribute_id)?.value}"></div>
																			{attributes.find(attr => attr.attribute_id === attribute.attribute_id)?.value}
																		</div>
																	{/if}
																</Select.Value>
															</Select.Trigger>
															<Select.Content>
																{#each attribute.options as option}
																	<Select.Item 
																		value={option.value} 
																		label={option.value}
																		on:click={() => {
																			const value = option.value;
																			const index = attributes.findIndex(attr => attr.attribute_id === attribute.attribute_id);
																			if (index !== -1) {
																				attributes[index] = { ...attributes[index] };
																				const fullAttribute = $typeAttributesQuery.data?.find(a => a.attribute_id === attribute.attribute_id && a.applies_to === 'Product');
																				const option = fullAttribute?.options?.find(o => o.value === value);
																				if (option) {
																					attributes[index].attribute_option_id = option.attribute_option_id;
																				}
																				attributes[index].value = value;
																				attributes = [...attributes];
																			} else {
																				let attribute_option_id;
																				const fullAttribute = $typeAttributesQuery.data?.find(a => a.attribute_id === attribute.attribute_id && a.applies_to === 'Product');
																				const option = fullAttribute?.options?.find(o => o.value === value);
																				if (option) {
																					attribute_option_id = option.attribute_option_id;
																				}
																				let newAttr = createProductAttribute(
																					attribute.attribute_id,
																					value,
																					attribute_option_id
																				);
																				attributes = [...attributes, newAttr];
																			}
																		}}
																	>
																		{#if attribute.data_type === 'Color'}
																			<div class="flex items-center gap-2">
																				<div class="h-4 w-4 rounded-full" style="background-color: {option.value}"></div>
																				{option.value}
																			</div>
																		{:else}
																			{option.value}
																		{/if}
																	</Select.Item>
																{/each}
															</Select.Content>
														</Select.Root>
													{/key}
												{:else}
													<Input 
														type={attribute.data_type === 'Number' ? 'number' : 'text'}
														id={attribute.title}
														placeholder={attribute.title}
														class="w-full"
														value={attributes.find(attr => attr.attribute_id === attribute.attribute_id)?.value || ''}
														on:input={(e) => {
															const value = e.currentTarget.value;
															const index = attributes.findIndex(attr => attr.attribute_id === attribute.attribute_id);
															if (index !== -1) {
																attributes[index].value = value;
																attributes = [...attributes];
															} else {
																let attribute_option_id;
																const fullAttribute = $typeAttributesQuery.data?.find(a => a.attribute_id === attribute.attribute_id && a.applies_to === 'Product');
																const option = fullAttribute?.options?.find(o => o.value === value);
																if (option) {
																	attribute_option_id = option.attribute_option_id;
																}
																let newAttr = createProductAttribute(
																	attribute.attribute_id,
																	value,
																	attribute_option_id
																);
																attributes = [...attributes, newAttr];
															}
														}}
														required={attribute.required}
													/>
												{/if}
											</div>
										{/if}
									{/each}
								{/if}
							</form>
						</div>
					</Card.Content>
				</Card.Root>

				<Card.Root>
					<Card.Header>
						<Card.Title>Variations</Card.Title>
						<Card.Description>Define variations for all attribute combinations</Card.Description>
					</Card.Header>
					<Card.Content>
						{#if showVariationSetup && variantAttributes.length > 0}
							<!-- Variation Setup Flow -->
							<div class="space-y-6">
								<div class="text-sm text-muted-foreground">
									Select which options you want to include for each variation attribute. This will determine the available product variations.
								</div>
								
								{#each variantAttributes as attribute}
									<div class="space-y-3">
										<Label class="text-base font-medium">{attribute.title}</Label>
										<div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3">
											{#if attribute.options}
												{#each attribute.options as option}
													{@const isSelected = isOptionSelected(attribute.attribute_id, option.value)}
													<div class="flex items-center space-x-2">
														<Checkbox 
															id={`option-${attribute.attribute_id}-${option.attribute_option_id}`}
															checked={isSelected}
															on:click={() => toggleOption(attribute.attribute_id, option.value)}
														/>
														<Label 
															for={`option-${attribute.attribute_id}-${option.attribute_option_id}`}
															class="text-sm font-normal cursor-pointer flex items-center gap-2"
														>
															{#if attribute.data_type === 'Color'}
																<div class="h-4 w-4 rounded-full border border-gray-200" style="background-color: {option.value}"></div>
															{/if}
															{option.value}
														</Label>
													</div>
												{/each}
											{/if}
										</div>
									</div>
								{/each}

								<div class="flex items-center justify-between pt-4 border-t">
									<div class="text-sm text-muted-foreground">
										{getTotalCombinations()} variations will be created
									</div>
									<div class="flex gap-2">
										<Button variant="outline" on:click={resetVariationSetup}>Reset</Button>
										<Button on:click={completeVariationSetup}>Configure Variations</Button>
									</div>
								</div>
							</div>
						{:else if variationsConfigured}
							<!-- Variation Management -->
							<div class="space-y-4">
								<div class="flex items-center justify-between">
									<h3 class="text-lg font-semibold">Active Variations</h3>
									<Button variant="outline" size="sm" on:click={resetVariationSetup}>
										Reconfigure Options
									</Button>
								</div>
								
								<div class="max-h-60 max-w-[53rem] overflow-y-auto">
									<Table.Root class="relative w-full table-fixed">
										<Table.Header class="sticky top-0 z-50 bg-gray-900">
											<Table.Row>
												<Table.Head class="w-[100px]">Stock</Table.Head>
												<Table.Head class="w-[100px]">Price</Table.Head>
												{#each variantAttributes as attribute}
													<Table.Head class="w-[100px]">{attribute.title}</Table.Head>
												{/each}
												<Table.Head class="w-[50px]">Action</Table.Head>
											</Table.Row>
										</Table.Header>
										<Table.Body>
											{#each activeCombinations as combination, index}
												<Table.Row>
													<Table.Cell class="whitespace-nowrap">
														<Label for={`stock-active-${index}`} class="sr-only">Stock</Label>
														<Input
															id={`stock-active-${index}`}
															type="number"
															placeholder="Stock"
															value={variants[index]?.available_quantity || 0}
															on:input={(e) => {
																const value = parseInt(e.currentTarget.value) || 0;
																if (!variants[index]) {
																	variants[index] = {
																		attributes: combination.map(attr => {
																			const variantAttr = variantAttributes.find(va => va.title === attr.title);
																			let attribute_option_id;
																			let attribute_id = 0;
																			if (variantAttr) {
																				attribute_id = variantAttr.attribute_id;
																				if (variantAttr.options) {
																					const option = variantAttr.options.find(o => o.value === attr.value);
																					if (option) {
																						attribute_option_id = option.attribute_option_id;
																					}
																				}
																			}
																			return createProductAttribute(
																				attribute_id,
																				attr.value,
																				attribute_option_id
																			);
																		}),
																		available_quantity: 0,
																		description: '',
																		is_default: false,
																		price: 0,
																		seo_description: '',
																		seo_keywords: [],
																		seo_title: '',
																		sku: ''
																	};
																}
																variants[index].available_quantity = value;
																variants = [...variants];
															}}
														/>
													</Table.Cell>
													<Table.Cell class="whitespace-nowrap">
														<Label for={`price-active-${index}`} class="sr-only">Price</Label>
														<div class="relative">
															<span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">{currencySymbol}</span>
															<Input
																id={`price-active-${index}`}
																type="text"
																placeholder="0.00"
																class="pl-7"
																value={formatAsCurrency(variants[index]?.price || 0)}
																on:input={(e) => {
																	const value = parseCurrencyInput(e.currentTarget.value);
																	if (!variants[index]) {
																		variants[index] = {
																			attributes: combination.map(attr => {
																				const variantAttr = variantAttributes.find(va => va.title === attr.title);
																				let attribute_option_id;
																				let attribute_id = 0;
																				if (variantAttr) {
																					attribute_id = variantAttr.attribute_id;
																					if (variantAttr.options) {
																						const option = variantAttr.options.find(o => o.value === attr.value);
																						if (option) {
																							attribute_option_id = option.attribute_option_id;
																						}
																					}
																				}
																				return createProductAttribute(
																					attribute_id,
																					attr.value,
																					attribute_option_id
																				);
																			}),
																			available_quantity: 0,
																			description: '',
																			is_default: false,
																			price: 0,
																			seo_description: '',
																			seo_keywords: [],
																			seo_title: '',
																			sku: ''
																		};
																	}
																	variants[index].price = value;
																	variants = [...variants];
																}}
															/>
														</div>
													</Table.Cell>
													{#each combination as attribute}
														<Table.Cell class="whitespace-nowrap">
															{#if variantAttributes.find(attr => attr.title === attribute.title)?.data_type === 'Color'}
																<div class="flex items-center gap-2">
																	<div class="h-4 w-4 rounded-full" style="background-color: {attribute.value}"></div>
																	{attribute.value}
																</div>
															{:else}
																{attribute.value}
															{/if}
														</Table.Cell>
													{/each}
													<Table.Cell class="whitespace-nowrap">
														<Button
															size="sm"
															variant="ghost"
															on:click={() => toggleVariation(index)}
															class="gap-1"
														>
															<Trash2 class="h-3.5 w-3.5" />
														</Button>
													</Table.Cell>
												</Table.Row>
											{/each}
										</Table.Body>
									</Table.Root>
								</div>

								{#if disabledCombinations.length > 0}
									<div class="mt-6">
										<h3 class="text-lg font-semibold mb-4">Disabled Variations</h3>
										<div class="max-h-60 max-w-[53rem] overflow-y-auto">
											<Table.Root class="relative w-full table-fixed">
												<Table.Header class="sticky top-0 z-50 bg-gray-900">
													<Table.Row>
														<Table.Head class="w-[100px]">Stock</Table.Head>
														<Table.Head class="w-[100px]">Price</Table.Head>
														{#each variantAttributes as attribute}
															<Table.Head class="w-[100px]">{attribute.title}</Table.Head>
														{/each}
														<Table.Head class="w-[50px]">Action</Table.Head>
													</Table.Row>
												</Table.Header>
												<Table.Body>
													{#each disabledCombinations as combination, index}
														<Table.Row class="opacity-50">
															<Table.Cell class="whitespace-nowrap">
																<Label for={`stock-disabled-${index}`} class="sr-only">Stock</Label>
																<Input id={`stock-disabled-${index}`} type="number" placeholder="Stock" disabled />
															</Table.Cell>
															<Table.Cell class="whitespace-nowrap">
																<Label for={`price-disabled-${index}`} class="sr-only">Price</Label>
																<Input id={`price-disabled-${index}`} type="number" placeholder="Price" disabled />
															</Table.Cell>
															{#each combination as attribute}
																<Table.Cell class="whitespace-nowrap">
																	{#if variantAttributes.find(attr => attr.title === attribute.title)?.data_type === 'Color'}
																		<div class="flex items-center gap-2">
																			<div class="h-4 w-4 rounded-full" style="background-color: {attribute.value}"></div>
																			{attribute.value}
																		</div>
																	{:else}
																		{attribute.value}
																	{/if}
																</Table.Cell>
															{/each}
															<Table.Cell class="whitespace-nowrap">
																<Button
																	size="sm"
																	variant="ghost"
																	on:click={() => toggleVariation(index, true)}
																	class="gap-1"
																>
																	<CirclePlus class="h-3.5 w-3.5" />
																</Button>
															</Table.Cell>
														</Table.Row>
													{/each}
												</Table.Body>
											</Table.Root>
										</div>
									</div>
								{/if}
							</div>
						{:else}
							<!-- No variant attributes -->
							<div class="text-center py-8 text-muted-foreground">
								<p>No variation attributes are defined for this product type.</p>
								<p class="text-sm">Add variation attributes to the product type to create product variations.</p>
							</div>
						{/if}
					</Card.Content>
				</Card.Root>
			</div>

			<div class="grid auto-rows-max items-start gap-4 lg:gap-8">
				<Card.Root>
					<Card.Header>
						<Card.Title>Product Status</Card.Title>
					</Card.Header>
					<Card.Content>
						<div class="grid gap-6">
							<div class="grid gap-3">
								<Label for="status">Status</Label>
								<Select.Root>
									<Select.Trigger id="status" aria-label="Select status">
										<Select.Value placeholder="Select status" />
									</Select.Trigger>
									<Select.Content>
										<Select.Item value="draft" label="Draft">Draft</Select.Item>
										<Select.Item value="published" label="Active">Active</Select.Item>
										<Select.Item value="archived" label="Archived">Archived</Select.Item>
									</Select.Content>
								</Select.Root>
							</div>
						</div>
					</Card.Content>
				</Card.Root>

				<Card.Root class="overflow-hidden">
					<Card.Header>
						<Card.Title>Product Images</Card.Title>
						<Card.Description>Upload product images</Card.Description>
					</Card.Header>
					<Card.Content>
						<div class="grid gap-2">
							<img
								alt="Product"
								class="aspect-square w-full rounded-md object-cover"
								height="300"
								src="/images/placeholder.png"
								width="300"
							/>
							<div class="grid grid-cols-3 gap-2">
								<button>
									<img
										alt="Product"
										class="aspect-square w-full rounded-md object-cover"
										height="84"
										src="/images/placeholder.png"
										width="84"
									/>
								</button>
								<button>
									<img
										alt="Product"
										class="aspect-square w-full rounded-md object-cover"
										height="84"
										src="/images/placeholder.png"
										width="84"
									/>
								</button>
								<button
									class="flex aspect-square w-full items-center justify-center rounded-md border border-dashed"
								>
									<Upload class="h-4 w-4 text-muted-foreground" />
									<span class="sr-only">Upload</span>
								</button>
							</div>
						</div>
					</Card.Content>
				</Card.Root>
			</div>
		</div>
		<div class="flex items-center justify-center gap-2 md:hidden">
			<Button variant="outline" size="sm">Discard</Button>
			<Button size="sm" on:click={saveProduct}>Save Product</Button>
		</div>
	</div>
</main>
