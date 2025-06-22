<script lang="ts">
	import ChevronLeft from 'lucide-svelte/icons/chevron-left';
	import Upload from 'lucide-svelte/icons/upload';
	import Trash2 from 'lucide-svelte/icons/trash-2';
	import { X } from 'lucide-svelte';

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
	import { CirclePlus } from 'lucide-svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { api } from '$lib/api'
	import type { Product, ProductAttribute, ProductVariant, ProductTypeAttribute, AttributeOption, ProductType, ProductImage, Shop } from '$lib/types'
	import { getContext } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import type { PageData } from './$types';
	import { uploadToR2, deleteFromR2 } from '$lib/cloudflare/r2-client';
	import { toast } from 'svelte-sonner';
	import { onMount } from 'svelte';
	import { getCurrencySymbol, formatAsCurrency, parseCurrencyInput } from '$lib/utils/currency';
	const authFetch = getContext('authFetch');
	const queryClient = useQueryClient();
	// Define an extended product type to handle category_id
	interface ExtendedProduct extends Product {
		category_id?: number;
	}

	export let data: PageData;

	// Get shop currency
	const shopQuery = createQuery<Shop, Error>({
		queryKey: [`shop-${$page.params.shop}`],
		queryFn: () => api(authFetch as any).getShop(),
		enabled: !!$page.params.shop
	});
	$: currencyCode = $shopQuery.data?.currency_code || 'USD';
	$: currencySymbol = getCurrencySymbol(currencyCode);

	// Get product data
	const productQuery = createQuery<ExtendedProduct>({
		queryKey: [`shop-${$page.params.shop}-product`, $page.params.productId],
		queryFn: async () => {
			return api(authFetch as any).getProductById(parseInt($page.params.productId));
		},
	});

	// Track product images loading state
	let imagesLoading = false;
	let imagesError: Error | null = null;

	// Separate query for product images
	const productImagesQuery = createQuery<ProductImage[]>({
		queryKey: [`shop-${$page.params.shop}-product-images`, $page.params.productId],
		queryFn: async () => {
			return api(authFetch as any).getProductImages(parseInt($page.params.productId));
		},
		enabled: !!$page.params.productId
	});

	// Product type attributes query
	let productTypeId: number | undefined;
	// Try to get it from product or from the page data
	$: productTypeId = $productQuery.data?.product_type_id || $productQuery.data?.category_id || data.productTypeId;

	// Get product type data for SKU generation
	const productTypeQuery = createQuery<ProductType>({
		queryKey: [`shop-${$page.params.shop}-product-type`, productTypeId],
		queryFn: async() => {
			if (!productTypeId) throw new Error('Product type ID not available');
			return api(authFetch as any).getProductTypeById(productTypeId);
		},
		enabled: !!productTypeId
	});

	// Get product type attributes for this product
	const typeAttributesQuery = createQuery<ProductTypeAttribute[]>({
		queryKey: [`shop-${$page.params.shop}-product-type-attributes`, productTypeId],
		queryFn: async () => {
			if (!productTypeId) throw new Error('Product type ID not available');
			return api(authFetch as any).getProductTypeAttributes(productTypeId);
		},
	});

	// For storing the updated product data
	let title = '';
	let description = '';
	let attributes: ProductAttribute[] = [];
	let variants: ProductVariant[] = [];
	let removedVariants: ProductVariant[] = []; // Add array to store removed variants

	// Update local state when product data is fetched
	$: if ($productQuery.data) {
		title = $productQuery.data.title;
		description = $productQuery.data.description || '';
		attributes = $productQuery.data.attributes || [];
		variants = $productQuery.data.variants || [];
		removedVariants = []; // Reset removed variants when loading new data
	}

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

	// Generate all combinations of variant attributes
	function generateCombinations(attributes: ProductTypeAttribute[]): { title: string; value: string }[][] {
		if (attributes.length === 0) return [[]];
		const [first, ...rest] = attributes;
		const combinations = generateCombinations(rest);
		const options: AttributeOption[] = Array.isArray(first.options) ? first.options : [];
		if (options.length === 0) {
			// If no options, treat as a single option with empty value to keep combinations
			return combinations.map(combination => [{ title: first.title, value: '' }, ...combination]);
		}
		return options.flatMap(option =>
			combinations.map(combination => [{ title: first.title, value: option.value }, ...combination])
		);
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
				const varAttr = variantAttributes.find(va => va.title === attr.title);
				if (!varAttr) return '';
				return `${varAttr.attribute_id}:${attr.value}`;
			})
			.filter(key => key !== '')
			.sort()
			.join('|');
	}

	// Generate a unique key for a variation
	function getVariationKey(combination: { title: string; value: string }[]): string {
		return combination.map(attr => `${attr.title}:${attr.value}`).sort().join('|');
	}

	// Generate SKU using the product type's sku_substring and a sequential number
	function generateSku(index: number): string {
		if (!$productTypeQuery.data?.sku_substring) return '';
		
		// Get product type substring
		const prefix = $productTypeQuery.data.sku_substring;
		
		// Format the index as a 4-digit number with leading zeros
		const suffix = String(index + 1).padStart(4, '0');
		
		return `${prefix}-${suffix}`;
	}

	// Find the variant attributes
	$: variantAttributes = $typeAttributesQuery.data?.filter((attr) => attr.applies_to === 'ProductVariation') || [];
	
	// Generate all possible combinations
	$: combinations = generateCombinations(variantAttributes);

	// Function to check if a combination exists in a variant
	function isCombinationInVariant(combination: { title: string; value: string }[], variant: ProductVariant): boolean {
		if (!combination || combination.length === 0 || !variant || !variant.attributes) {
			return false;
		}

		// Convert combination to a map for easier lookup
		const combMap = new Map();
		combination.forEach(combAttr => {
			const varAttr = variantAttributes.find(attr => attr.title === combAttr.title);
			if (varAttr) {
				combMap.set(varAttr.attribute_id, combAttr.value);
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

	// Separate lists for active and disabled variations
	let activeCombinations: { title: string; value: string }[][] = [];
	let disabledCombinations: { title: string; value: string }[][] = [];
	
	// Track which variations we've already seen
	let seenVariations = new Set<string>();

	// Process combinations when product and variant attributes load
	$: if ($typeAttributesQuery.data && combinations.length > 0) {
		// Reset collections
		activeCombinations = [];
		disabledCombinations = [];
		seenVariations = new Set();

		// First identify combinations that match existing variants
		combinations.forEach(combination => {
			const key = getVariationKey(combination);
			
			// Skip if we've already seen this variation
			if (seenVariations.has(key)) return;
			seenVariations.add(key);
			
			// Check if this combination exists in any variant
			const matchingVariant = variants.some(variant => isCombinationInVariant(combination, variant));
			
			if (matchingVariant) {
				activeCombinations.push(combination);
			} else {
				disabledCombinations.push(combination);
			}
		});
	}

	// Function to create a new variant for a combination
	function createVariantFromCombination(combination: { title: string; value: string }[]): ProductVariant {
		// Calculate next index for SKU generation
		const nextIndex = variants.length;
		
		return {
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
			available_quantity: 1,
			description: '',
			is_default: variants.length === 0, // Make first variant the default
			price: 0.00,
			seo_description: null,
			seo_keywords: null,
			seo_title: null,
			sku: generateSku(nextIndex),
			slug: '',
			created_at: null,
			updated_at: null
		};
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
				const newVariant = createVariantFromCombination(variation);
				// Make sure it has a unique SKU (in case we're adding after removing others)
				newVariant.sku = generateSku(variants.length);
				variants = [...variants, newVariant];
			}
			
			activeCombinations = [...activeCombinations, variation];
			disabledCombinations = [...disabledCombinations];
		} else {
			// Move from active to disabled
			const combination = activeCombinations[index];
			activeCombinations.splice(index, 1);
			
			// Find the matching variant
			const variantIndex = variants.findIndex(variant => isCombinationInVariant(combination, variant));
			if (variantIndex !== -1) {
				// Store the removed variant for potential later use
				removedVariants = [...removedVariants, variants[variantIndex]];
				// Remove it from active variants
				variants = variants.filter((_, i) => i !== variantIndex);
			}
			
			disabledCombinations = [...disabledCombinations, combination];
			activeCombinations = [...activeCombinations];
		}
	}

	// Save updated product
	async function updateProduct() {
		try {
			const updatedProduct = await api(authFetch as any).updateProduct(parseInt($page.params.productId), {
				title,
				description,
				attributes,
				status: $productQuery.data?.status,
				variants: variants.map(variant => ({
					...variant,
					seo_keywords: variant.seo_keywords || null,
					seo_description: variant.seo_description || null,
					seo_title: variant.seo_title || null
				}))
			});
			
			// Invalidate and refetch the query
			await queryClient.invalidateQueries({ queryKey: [`shop-${$page.params.shop}-product`, $page.params.productId] });
		} catch (error) {
			console.error('Failed to update product:', error);
		}
	}

	// For file uploads
	let isUploading = false;
	let fileInput: HTMLInputElement;

	// Handle image upload
	async function handleImageUpload(event: Event) {
		const target = event.target as HTMLInputElement;
		const files = target.files;
		
		// Use the route parameter directly as the product ID
		const productId = $page.params.productId;
		
		if (!files || files.length === 0) {
			toast.error('No file selected');
			return;
		}
		
		if (!productId) {
			toast.error('Invalid product ID in the URL');
			return;
		}
		
		isUploading = true;
		
		let cloudflareData: { url: string; filename: string } | null = null;
		
		try {
			const file = files[0];
			
			// Upload to Cloudflare R2 with product ID
			cloudflareData = await uploadToR2(file, productId);
			
			// Add the image to the product in the database
			const isPrimary = !$productImagesQuery.data || $productImagesQuery.data.length === 0;
			
			try {
				await api(authFetch as any).addProductImage(productId, {
					url: cloudflareData.url,
					alt: file.name, // Use the filename as alt text by default
					filename: cloudflareData.filename,
					is_primary: isPrimary
				});
				
				// Invalidate and refetch the query
				await queryClient.invalidateQueries({ queryKey: [`shop-${$page.params.shop}-product-images`, $page.params.productId] });
			} catch (error) {
				// If adding to database fails, clean up the uploaded file from Cloudflare
				if (cloudflareData) {
					try {
						await deleteFromR2(cloudflareData.filename);
						console.log('Cleaned up Cloudflare file after database failure');
					} catch (cleanupError) {
						console.error('Failed to clean up Cloudflare file:', cleanupError);
					}
				}
				throw error; // Re-throw to be caught by the outer catch
			}
		} catch (error) {
			console.error('Image upload failed:', error);
			toast.error('Failed to upload image');
		} finally {
			isUploading = false;
			// Reset the file input
			if (fileInput) fileInput.value = '';
		}
	}

	// Function to extract filename from URL
	function extractFilenameFromUrl(url: string): string {
		// Extract the last part of the URL path after the last slash
		const urlParts = url.split('/');
		return urlParts[urlParts.length - 1];
	}

	// Extract the full path for Cloudflare R2 deletion
	function extractCloudflareObjectKey(url: string): string {
		try {
			// The URL format from Cloudflare is like:
			// https://7ff73a3850a237e871a1a720395d4fb7.r2.cloudflarestorage.com/products/88/1745184806440-unnamed.png
			
			// Parse the URL to extract the pathname
			const urlObj = new URL(url);
			// Remove the leading slash if present
			const pathParts = urlObj.pathname.split('/');
			// Remove the first empty segment if the path starts with a slash
			if (pathParts[0] === '') pathParts.shift();
			
			// Return the path that should be in the format: 'products/productId/timestamp-filename'
			return pathParts.join('/');
		} catch (e) {
			console.error('Failed to parse URL:', e);
			// Fallback: try to extract the path using regex
			const match = url.match(/([^\/]+\/[^\/]+\/[^\/]+)$/);
			return match ? match[1] : '';
		}
	}

	// Handle image deletion
	async function deleteImage(image: ProductImage) {
		// Use the route parameter directly as the product ID
		const productId = $page.params.productId;
		
		if (!productId) {
			toast.error('Invalid product ID in the URL');
			return;
		}

		if (!image.id) {
			toast.error('Cannot delete image: Image ID is not available');
			return;
		}

		if (!confirm('Are you sure you want to delete this image?')) {
			return;
		}
		
		try {
			// First delete from database to avoid orphaned references
			const deleted = await api(authFetch as any).deleteProductImage(productId, image.id);
			
			if (deleted) {
				try {
					// Try to extract the object key from the URL
					const objectKey = image.filename || extractCloudflareObjectKey(image.url);
					
					if (objectKey) {
						console.log('Deleting from Cloudflare R2:', objectKey);
						await deleteFromR2(objectKey);
						console.log('Successfully deleted from Cloudflare R2');
					} else {
						console.warn('Could not determine object key for deletion from URL:', image.url);
						toast.warning('Image deleted from database but could not determine storage path');
					}
				} catch (cloudflareError) {
					console.error('Failed to delete from Cloudflare R2:', cloudflareError);
					// Don't throw here, as the database deletion was successful
					toast.warning('Image deleted from database but failed to delete from storage');
				}
				
				// Invalidate and refetch the query
				await queryClient.invalidateQueries({ queryKey: [`shop-${$page.params.shop}-product-images`, $page.params.productId] });
			} else {
				throw new Error('Failed to delete from database');
			}
		} catch (error) {
			console.error('Failed to delete image:', error);
			toast.error('Failed to delete image');
		}
	}

	// Set image as primary by moving it to the first position in the array
	async function setPrimaryImage(image: ProductImage) {
		// Use the route parameter directly as the product ID
		const productId = $page.params.productId;
		
		if (!productId) {
			toast.error('Invalid product ID in the URL');
			return;
		}

		if (!image.id) {
			toast.error('Cannot set primary image: Image ID is not available');
			return;
		}

		if (primaryImage && primaryImage.id === image.id) return; // Already primary
		
		try {
			await api(authFetch as any).setProductPrimaryImage(productId, image.id);
			
			// Invalidate and refetch the query
			await queryClient.invalidateQueries({ queryKey: [`shop-${$page.params.shop}-product-images`, $page.params.productId] });
			
			// Also update the local UI immediately for better user experience
			if ($productImagesQuery.data) {
				// Find image index
				const imageIndex = $productImagesQuery.data.findIndex(img => img.id === image.id);
				if (imageIndex > 0) {
					// Create a new array with the selected image as the first element
					const reorderedImages = [
						$productImagesQuery.data[imageIndex],
						...$productImagesQuery.data.slice(0, imageIndex),
						...$productImagesQuery.data.slice(imageIndex + 1)
					];
					// Update the query data
					queryClient.setQueryData([`shop-${$page.params.shop}-product-images`, $page.params.productId], reorderedImages);
				}
			}
		} catch (error) {
			console.error('Failed to set primary image:', error);
			toast.error('Failed to set primary image');
		}
	}

	// Get primary image for display - no 'is_primary' flag in API response, so we'll use the first image as primary
	$: primaryImage = $productImagesQuery.data && $productImagesQuery.data.length > 0 ? 
		$productImagesQuery.data[0] : null;

	// Get other images for thumbnails - all images except the first one
	$: otherImages = $productImagesQuery.data && $productImagesQuery.data.length > 1 ? 
		$productImagesQuery.data.slice(1) : [];

	// Handle image loading state
	$: imagesLoading = $productImagesQuery.isLoading || $productImagesQuery.isFetching;
	$: imagesError = $productImagesQuery.error ? new Error($productImagesQuery.error.message) : null;

	onMount(async () => {
		await queryClient.prefetchQuery({
			queryKey: [`shop-${$page.params.shop}-product`, $page.params.productId],
			queryFn: () => api(authFetch as any).getProductById(parseInt($page.params.productId)),
		});
	});
</script>

<main class="grid flex-1 items-start gap-4 p-4 sm:px-6 sm:py-0 md:gap-8">
	{#if $productQuery.isPending}
		<span>Loading...</span>
	{/if}
	{#if $productQuery.error}
		<span>Error: {$productQuery.error.message}</span>
	{/if}
	{#if $productQuery.data}
		<div class="mx-auto grid max-w-[59rem] flex-1 auto-rows-max gap-4">
			<div class="flex items-center gap-4">
				<Button variant="outline" size="icon" class="h-7 w-7" href={`/${$page.params.shop}/product-types/${$page.params.type}/products`}>
					<ChevronLeft class="h-4 w-4" />
					<span class="sr-only">Back</span>
				</Button>
				<h1 class="flex-1 shrink-0 whitespace-nowrap text-xl font-semibold tracking-tight sm:grow-0">
					{$productQuery.data.title}
				</h1>
				<Badge variant="outline" class="ml-auto sm:ml-0">{$productQuery.data.status}</Badge>
				<div class="hidden items-center gap-2 md:ml-auto md:flex">
					<Button variant="outline" size="sm">Discard</Button>
					<Button size="sm" on:click={updateProduct}>Save Product</Button>
				</div>
			</div>
			<div class="grid gap-4 md:grid-cols-[1fr_250px] lg:grid-cols-3 lg:gap-8">
				<div class="grid auto-rows-max items-start gap-4 lg:col-span-2 lg:gap-8">
					<Card.Root>
						<Card.Header>
							<Card.Title>Product Details</Card.Title>
							<Card.Description>Edit product information</Card.Description>
						</Card.Header>
						<Card.Content>
							<div class="grid gap-6">
								<div class="grid gap-3">
									<Label for="name">Name</Label>
									<Input id="name" type="text" class="w-full" bind:value={title} />
								</div>
								<div class="grid gap-3">
									<Label for="description">Description</Label>
									<Textarea
										id="description"
										bind:value={description}
										class="min-h-32"
									/>
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
													{#if attribute.options && attribute.options.length > 0}
														{#key attribute.attribute_id}
															<Select.Root>
																<Select.Trigger id={attribute.title} aria-label={attribute.title}>
																	{#if attribute.data_type === 'Color' && attributes.find(attr => attr.attribute_id === attribute.attribute_id)?.value}
																		<Select.Value>
																			<div class="flex items-center gap-2">
																				<div class="h-4 w-4 rounded-full" style="background-color: {attributes.find(attr => attr.attribute_id === attribute.attribute_id)?.value}"></div>
																				{attributes.find(attr => attr.attribute_id === attribute.attribute_id)?.value}
																			</div>
																		</Select.Value>
																	{:else}
																		<Select.Value placeholder={attributes.find(attr => attr.attribute_id === attribute.attribute_id)?.value || `Select ${attribute.title}`} />
																	{/if}
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
																					attributes[index].value = value;
																					attributes[index].attribute_option_id = option.attribute_option_id;
																					attributes = [...attributes];
																				} else {
																					let newAttr = createProductAttribute(
																						attribute.attribute_id,
																						value,
																						option.attribute_option_id
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
																	let newAttr = createProductAttribute(
																		attribute.attribute_id,
																		value
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
							<Card.Description>Manage product variations</Card.Description>
						</Card.Header>
						<Card.Content class="max-h-60 max-w-[53rem] overflow-y-auto">
							{#if $typeAttributesQuery.isPending}
								<p>Loading attribute variations...</p>
							{:else if $typeAttributesQuery.error}
								<p class="text-red-500">Error loading variations: {$typeAttributesQuery.error.message}</p>
							{:else if !variantAttributes || variantAttributes.length === 0}
								<p class="text-yellow-400">No variant attributes found for this product type</p>
							{:else}
								<h3 class="text-lg font-semibold">Active Variations</h3>
								{#if activeCombinations.length === 0}
									<p class="text-gray-400 my-4">No active variations found. Add variations from the list below.</p>
								{:else}
									<Table.Root class="relative w-full table-fixed">
										<Table.Header class="sticky top-0 z-50 bg-gray-900">
											<Table.Row>
												<Table.Head class="w-[100px]">SKU</Table.Head>
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
												{@const variant = variants.find(v => isCombinationInVariant(combination, v))}
												{#if variant}
													<Table.Row>
														<Table.Cell class="font-semibold whitespace-nowrap">
															{variant.sku || generateSku(index)}
														</Table.Cell>
														<Table.Cell class="whitespace-nowrap">
															<Label for={`stock-active-${index}`} class="sr-only">Stock</Label>
															<Input
																id={`stock-active-${index}`}
																type="number"
																placeholder="Stock"
																value={variant.available_quantity || 0}
																on:input={(e) => {
																	const value = parseInt(e.currentTarget.value) || 0;
																	const variantIndex = variants.findIndex(v => v === variant);
																	if (variantIndex !== -1) {
																		variants[variantIndex].available_quantity = value;
																		variants = [...variants];
																	}
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
																	value={formatAsCurrency(variant.price || 0)}
																	on:input={(e) => {
																		const value = parseCurrencyInput(e.currentTarget.value);
																		const variantIndex = variants.findIndex(v => v === variant);
																		if (variantIndex !== -1) {
																			variants[variantIndex].price = value;
																			variants = [...variants];
																		}
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
												{:else}
													<Table.Row class="bg-red-900/20">
														<Table.Cell colspan={3 + variantAttributes.length} class="p-2">
															<div class="flex justify-between items-center">
																<span>Missing variant data</span>
																<Button 
																	size="sm" 
																	on:click={() => {
																		// Get the combination key
																		const combinationKey = getCombinationKey(combination);
																	
																		// Check if there's a previously removed variant with the same attributes
																		const previousVariant = removedVariants.find(v => getVariantKey(v) === combinationKey);
																	
																		if (previousVariant) {
																			// Use the previously removed variant
																			variants = [...variants, previousVariant];
																			// Remove it from the removed variants array
																			removedVariants = removedVariants.filter(v => getVariantKey(v) !== combinationKey);
																		} else {
																			// Create default variant for this combination
																			const newVariant = createVariantFromCombination(combination);
																			// Add variant
																			variants = [...variants, newVariant];
																		}
																	}}
																>
																	Create Variant
																</Button>
															</div>
														</Table.Cell>
													</Table.Row>
												{/if}
											{/each}
										</Table.Body>
									</Table.Root>
								{/if}

								{#if disabledCombinations.length > 0}
									<h3 class="mt-6 text-lg font-semibold">Add New Variations</h3>
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
								{/if}
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
											<Select.Value placeholder={$productQuery.data?.status || "Select status"} />
										</Select.Trigger>
										<Select.Content>
											<Select.Item value="draft" label="Draft" on:click={() => { if ($productQuery.data) $productQuery.data.status = 'draft'; }}>Draft</Select.Item>
											<Select.Item value="published" label="Active" on:click={() => { if ($productQuery.data) $productQuery.data.status = 'published'; }}>Active</Select.Item>
											<Select.Item value="archived" label="Archived" on:click={() => { if ($productQuery.data) $productQuery.data.status = 'archived'; }}>Archived</Select.Item>
										</Select.Content>
									</Select.Root>
								</div>
							</div>
						</Card.Content>
					</Card.Root>

					<Card.Root class="overflow-hidden">
						<Card.Header>
							<Card.Title>Product Images</Card.Title>
							<Card.Description>Manage product images</Card.Description>
						</Card.Header>
						<Card.Content>
							{#if $productImagesQuery.isPending}
								<div class="flex justify-center items-center h-40">
									<p>Loading images...</p>
								</div>
							{:else if $productImagesQuery.error}
								<div class="p-4 bg-destructive/20 rounded-md text-destructive mb-4">
									<p>Error loading images: {$productImagesQuery.error.message}</p>
								</div>
							{:else if $productImagesQuery.data && $productImagesQuery.data.length === 0}
								<div class="p-4 bg-muted rounded-md mb-4">
									<p>No images found for this product. Upload an image below.</p>
								</div>
								<div class="grid gap-2">
									<label
										class="flex aspect-square w-full cursor-pointer items-center justify-center rounded-md border border-dashed hover:bg-muted/50"
									>
										<input 
											type="file" 
											accept="image/*" 
											class="hidden" 
											on:change={handleImageUpload}
											bind:this={fileInput}
											disabled={isUploading}
										/>
										{#if isUploading}
											<span class="text-xs">Uploading...</span>
										{:else}
											<Upload class="h-6 w-6 text-muted-foreground" />
											<span class="ml-2">Upload Image</span>
										{/if}
									</label>
								</div>
							{:else}
								<div class="grid gap-2">
									<img
										alt="Product"
										class="aspect-square w-full rounded-md object-cover"
										height="300"
										src={primaryImage?.url || "/images/placeholder.png"}
										width="300"
										on:click={() => primaryImage && setPrimaryImage(primaryImage)}
									/>
									<div class="grid grid-cols-3 gap-2">
										{#if otherImages.length > 0}
											{#each otherImages.slice(0, 2) as image}
												<div class="relative">
													<button class="group" on:click={() => setPrimaryImage(image)}>
														<img
															alt="Product"
															class="aspect-square w-full rounded-md object-cover"
															height="84"
															src={image.url}
															width="84"
														/>
														<div class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50 opacity-0 group-hover:opacity-100">
															<span class="text-xs text-white">Set Primary</span>
														</div>
													</button>
													<button 
														class="absolute -top-2 -right-2 flex h-5 w-5 items-center justify-center rounded-full bg-destructive text-white hover:bg-destructive/90"
														on:click={() => deleteImage(image)}
													>
														<X class="h-3 w-3" />
													</button>
												</div>
											{/each}
										{/if}
										
										{#if !otherImages.length || otherImages.length < 2}
											<div class="col-span-{!otherImages.length ? 3 : otherImages.length === 1 ? 2 : 1}">
												<label
													class="flex aspect-square w-full cursor-pointer items-center justify-center rounded-md border border-dashed hover:bg-muted/50"
												>
													<input 
														type="file" 
														accept="image/*" 
														class="hidden" 
														on:change={handleImageUpload}
														bind:this={fileInput}
														disabled={isUploading}
													/>
													{#if isUploading}
														<span class="text-xs">Uploading...</span>
													{:else}
														<Upload class="h-4 w-4 text-muted-foreground" />
														<span class="sr-only">Upload</span>
													{/if}
												</label>
											</div>
										{/if}
									</div>
									
									{#if otherImages.length > 2}
										<div class="grid grid-cols-6 gap-2 mt-2">
											{#each otherImages.slice(2) as image, index}
												<div class="relative">
													<button class="group" on:click={() => setPrimaryImage(image)}>
														<img
															alt="Product"
															class="aspect-square w-full rounded-md object-cover"
															height="50"
															src={image.url}
															width="50"
														/>
														<div class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50 opacity-0 group-hover:opacity-100">
															<span class="text-[10px] text-white">Set Primary</span>
														</div>
													</button>
													<button 
														class="absolute -top-2 -right-2 flex h-4 w-4 items-center justify-center rounded-full bg-destructive text-white hover:bg-destructive/90"
														on:click={() => deleteImage(image)}
													>
														<X class="h-2 w-2" />
													</button>
												</div>
											{/each}
											
											{#if otherImages.length < 8}
												<label
													class="flex aspect-square w-full cursor-pointer items-center justify-center rounded-md border border-dashed hover:bg-muted/50"
												>
													<input 
														type="file" 
														accept="image/*" 
														class="hidden" 
														on:change={handleImageUpload}
														disabled={isUploading}
													/>
													{#if isUploading}
														<span class="text-[10px]">...</span>
													{:else}
														<Upload class="h-3 w-3 text-muted-foreground" />
													{/if}
												</label>
											{/if}
										</div>
									{/if}
									
									{#if primaryImage}
										<button 
											class="bg-destructive text-white p-1 text-xs rounded-sm w-full hover:bg-destructive/90"
											on:click={() => deleteImage(primaryImage)}
										>
											<div class="flex items-center justify-center gap-1">
												<Trash2 class="h-3 w-3" />
												<span>Delete Primary Image</span>
											</div>
										</button>
									{/if}
								</div>
							{/if}
						</Card.Content>
					</Card.Root>
				</div>
			</div>
			<div class="flex items-center justify-center gap-2 md:hidden">
				<Button variant="outline" size="sm">Discard</Button>
				<Button size="sm" on:click={updateProduct}>Update Product</Button>
			</div>
		</div>
	{/if}
</main>
