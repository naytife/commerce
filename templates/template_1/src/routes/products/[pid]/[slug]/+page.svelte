<script lang="ts">

	import * as Card from '$lib/components/ui/card';
	import * as Accordion from '$lib/components/ui/accordion';
	import Minus from 'lucide-svelte/icons/minus';
	import Plus from 'lucide-svelte/icons/plus';
	import ShoppingBag from 'lucide-svelte/icons/shopping-bag';
	import * as Select from '$lib/components/ui/select';
	import type { PageData } from './$types';
	import { onMount } from 'svelte';
	import { afterNavigate } from '$app/navigation';
	import { initFlowbite } from 'flowbite';
	import { page } from '$app/stores';
	import { get } from 'svelte/store';
	import { goto } from '$app/navigation';
	import { cart } from '$lib/stores/cart';
	import { currencySymbol } from '$lib/stores/currency';

	export let data: PageData;
	let selectedImage: number = 0;
	let ProductQuery: any;
	$: ProductQuery = data.ProductQuery;

	onMount(() => {
		initFlowbite();
		syncActiveVariant();
	});

	let attributeOptions: { title: string; values: string[] }[] = [];

	// store user-selected values per attribute (used for swatches)
	let selectedAttributes: Record<string, string> = {};
	// will hold the currently active variation
	let activeVariant: any;
	let quantity: number = 1;

	// sync activeVariant and selectedAttributes for current slug
	function syncActiveVariant() {
		const pageData = get(page);
		const slug = pageData.params.slug;
		const variants = ProductQuery.data?.product?.variants;
		if (!slug || !variants) return;
		// extract variationId from slug or hash
		const parts = slug.split('-');
		let variationIdStr = parts.pop() || '';
		let variationId = Number(variationIdStr);
		if (isNaN(variationId) || variationId <= 0) {
			const hash = pageData.url.hash.substring(1);
			if (hash) {
				const hashParts = hash.split('-');
				variationIdStr = hashParts.pop() || '';
				variationId = Number(variationIdStr);
			}
		}
		const variant = variants.find((v: any) => v.variationId === variationId);
		if (variant) {
			activeVariant = variant;
			selectedAttributes = {};
			variant.attributes.forEach((a: any) => {
				selectedAttributes[a.title] = a.value;
			});
		}
	}

	// re-sync when navigating in-app
	afterNavigate(() => {
		syncActiveVariant();
	});

	onMount(() => {
	if (ProductQuery.data?.product?.variants) {
		const optionsMap = new Map<string, Set<string>>();
		for (const variant of ProductQuery.data.product.variants) {
		for (const attr of variant.attributes) {
			if (!optionsMap.has(attr.title)) {
			optionsMap.set(attr.title, new Set());
			}
			optionsMap.get(attr.title)?.add(attr.value);
		}
		}
		attributeOptions = Array.from(optionsMap.entries()).map(([title, valuesSet]) => ({
		title,
		values: Array.from(valuesSet)
		}));
	}
	});

	$: if (activeVariant) {
		if (quantity < 1) quantity = 1;
		if (quantity > activeVariant.availableQuantity) quantity = activeVariant.availableQuantity;
	}

	function constructVariantUrl(variant: any): string {
		const product = ProductQuery.data?.product;
		if (!product) return '';
		const productId = product.productId;
		const slug = product.slug;
		const attrsSlug = [...variant.attributes]
			.sort((a, b) => a.title.localeCompare(b.title))
			.map(a => encodeURIComponent(a.value.toLowerCase()))
			.join('-');
		return `/products/${productId}/${slug}-${attrsSlug}-${variant.variationId}`;
	}

	// find a variant matching a given set of attribute selections
	function findVariantByAttributes(attrs: Record<string, string>): any {
		return ProductQuery.data?.product?.variants.find((variant: any) =>
			variant.attributes.every((a: { title: string; value: string }) => attrs[a.title] === a.value)
		);
	}

	// handle selection changes in non-color attributes
	function handleAttributeSelect(title: string, value: string) {
		const attrs = { ...selectedAttributes, [title]: value };
		const variant = findVariantByAttributes(attrs);
		if (variant) {
			goto(constructVariantUrl(variant));
		} else {
			selectedAttributes[title] = value;
		}
	}
</script>

<section class="bg-white text-gray-900 dark:bg-gray-900 dark:text-white py-12 antialiased">
	<div class="mx-auto max-w-7xl px-6">
		<div class="flex flex-col lg:flex-row gap-12">
			<div class="w-full lg:w-3/5">
				<div class="flex flex-col lg:flex-row gap-6">
					<div class="order-2 lg:order-1 flex lg:flex-col gap-4 overflow-x-auto lg:overflow-x-visible pb-4 lg:pb-0">
						{#each ProductQuery.data?.product?.images ?? [] as img, index}
							<button 
								class="border border-gray-200 dark:border-gray-700 p-1 min-w-20 h-20 flex-shrink-0 transition-colors duration-200 {selectedImage === index ? 'border-primary-700 dark:border-primary-500' : 'hover:border-gray-300 dark:hover:border-gray-600'}" 
								on:click={() => selectedImage = index}
							>
								<img src={img.url} alt={img.altText} class="w-full h-full object-cover" />
							</button>
						{/each}
					</div>
					<div class="order-1 lg:order-2 bg-gray-50 dark:bg-gray-800 flex-grow">
						<img 
							src={ProductQuery.data?.product?.images?.[selectedImage]?.url} 
							alt={ProductQuery.data?.product?.images?.[selectedImage]?.altText ?? 'Product'} 
							class="w-full h-auto lg:h-[500px] object-contain" 
						/>
					</div>
				</div>
			</div>
			
			<div class="w-full lg:w-2/5">
				<div class="border-b border-gray-200 dark:border-gray-700 pb-6 mb-6">
					<h1 class="text-2xl font-medium text-gray-900 dark:text-white mb-3">
						{ProductQuery.data?.product?.title}
					</h1>
					<p class="text-sm text-gray-500 dark:text-gray-400 mb-4">{ProductQuery.data?.product?.description}</p>
					
					<div class="flex items-center justify-between">
						{#if activeVariant}
							<p class="text-2xl font-semibold text-gray-900 dark:text-white">
								{$currencySymbol}{activeVariant.price.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
							</p>
							<div class="flex items-center">
								{#if activeVariant.stockStatus === 'IN_STOCK'}
									<span class="inline-flex items-center text-xs uppercase tracking-wider font-medium text-green-700 dark:text-green-500">
										In Stock
									</span>
								{:else if activeVariant.stockStatus === 'OUT_OF_STOCK'}
									<span class="inline-flex items-center text-xs uppercase tracking-wider font-medium text-red-700 dark:text-red-500">
										Out of Stock
									</span>
								{:else if activeVariant.stockStatus === 'PREORDER'}
									<span class="inline-flex items-center text-xs uppercase tracking-wider font-medium text-blue-700 dark:text-blue-500">
										Pre-order
									</span>
								{/if}
							</div>
						{/if}
					</div>
				</div>
				
				{#if attributeOptions.length > 0}
					<div class="mb-8">
						{#each attributeOptions as attrOption}
							<div class="mb-6">
								<h3 class="text-sm font-medium uppercase tracking-wider text-gray-700 dark:text-gray-300 mb-3">
									{attrOption.title}
								</h3>
								<div class="flex flex-wrap gap-2">
									{#each attrOption.values as value}
										{#if attrOption.title.toLowerCase() === 'color'}
											<button 
												class="w-10 h-10 border {selectedAttributes[attrOption.title] === value ? 'border-primary-700 dark:border-primary-500' : 'border-gray-200 dark:border-gray-700'} p-0.5"
												on:click={() => handleAttributeSelect(attrOption.title, value)}
												aria-label="Select {value} color"
											>
												<span class="block w-full h-full" style="background-color: {value.toLowerCase()}"></span>
											</button>
										{:else}
											<button
												class="px-4 py-2 border {selectedAttributes[attrOption.title] === value ? 'border-primary-700 dark:border-primary-500 bg-primary-50 dark:bg-primary-900/20 text-primary-700 dark:text-primary-500' : 'border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800'} transition-colors duration-200"
												on:click={() => handleAttributeSelect(attrOption.title, value)}
											>
												{value}
											</button>
										{/if}
									{/each}
								</div>
							</div>
						{/each}
					</div>
				{/if}
				
				<div class="mb-8">
					<div class="flex items-center justify-between mb-4">
						<h3 class="text-sm font-medium uppercase tracking-wider text-gray-700 dark:text-gray-300">Quantity</h3>
						<p class="text-sm text-gray-500 dark:text-gray-400">{activeVariant?.availableQuantity ?? 0} available</p>
					</div>
					<div class="flex border border-gray-200 dark:border-gray-700 w-full max-w-[180px]">
						<button 
							type="button"
							on:click={() => quantity = Math.max(1, quantity - 1)}
							class="w-12 h-12 flex items-center justify-center text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors duration-200"
							disabled={quantity <= 1}
						>
							<Minus class="w-4 h-4" />
						</button>
						<input
							type="number"
							min="1"
							max={activeVariant?.availableQuantity ?? 1}
							bind:value={quantity}
							class="w-16 h-12 text-center border-x border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:outline-none"
						/>
						<button 
							type="button"
							on:click={() => quantity = Math.min(activeVariant?.availableQuantity ?? 1, quantity + 1)}
							class="w-12 h-12 flex items-center justify-center text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors duration-200"
							disabled={quantity >= (activeVariant?.availableQuantity ?? 1)}
						>
							<Plus class="w-4 h-4" />
						</button>
					</div>
				</div>
				
				<div class="mb-8">
					<button
						type="button"
						on:click={() => cart.add({
							id: (activeVariant?.variationId ?? ProductQuery.data!.product!.defaultVariant.variationId).toString(),
							title: ProductQuery.data!.product!.title,
							price: activeVariant?.price ?? ProductQuery.data!.product!.defaultVariant.price,
							image: ProductQuery.data!.product!.images[selectedImage]?.url,
							slug: $page.params.slug
						}, quantity)}
						class="w-full bg-primary-700 hover:bg-primary-800 text-white font-medium uppercase tracking-wider py-4 flex items-center justify-center gap-2 transition-colors duration-200 focus:ring-4 focus:ring-primary-300 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
					>
						<ShoppingBag class="w-5 h-5" />
						Add to Cart
					</button>
				</div>
				
				<div class="border-t border-gray-200 dark:border-gray-700 pt-8">
					<Accordion.Root value="item-1" class="border-b border-gray-200 dark:border-gray-700">
						<Accordion.Item value="item-1" class="border-t-0">
							<Accordion.Trigger class="text-sm font-medium uppercase tracking-wider py-4">Product Details</Accordion.Trigger>
							<Accordion.Content class="pb-4 text-sm text-gray-600 dark:text-gray-400">
								{ProductQuery.data?.product?.description}
							</Accordion.Content>
						</Accordion.Item>
						<Accordion.Item value="item-2" class="border-t border-gray-200 dark:border-gray-700">
							<Accordion.Trigger class="text-sm font-medium uppercase tracking-wider py-4">Specifications</Accordion.Trigger>
							<Accordion.Content class="pb-4">
								{#if ProductQuery.data?.product?.attributes && ProductQuery.data.product.attributes.length > 0}
									<ul class="space-y-2 text-sm text-gray-600 dark:text-gray-400">
										{#each ProductQuery.data.product.attributes as attr}
											<li class="flex justify-between">
												<span class="font-medium text-gray-700 dark:text-gray-300">{attr.title}:</span> 
												<span>{attr.value}</span>
											</li>
										{/each}
									</ul>
								{:else}
									<p class="text-sm text-gray-600 dark:text-gray-400">No specifications available.</p>
								{/if}
							</Accordion.Content>
						</Accordion.Item>
						<Accordion.Item value="item-3" class="border-t border-gray-200 dark:border-gray-700">
							<Accordion.Trigger class="text-sm font-medium uppercase tracking-wider py-4">Shipping</Accordion.Trigger>
							<Accordion.Content class="pb-4 text-sm text-gray-600 dark:text-gray-400">
								Free shipping on orders over {$currencySymbol}50. Standard shipping takes 3-5 business days.
							</Accordion.Content>
						</Accordion.Item>
					</Accordion.Root>
				</div>
			</div>
		</div>
	</div>
</section>
