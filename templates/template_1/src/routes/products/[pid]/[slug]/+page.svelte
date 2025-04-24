<script lang="ts">

	import * as Card from '$lib/components/ui/card';
	import * as Accordion from '$lib/components/ui/accordion';
	import { Info, MapPin } from 'lucide-svelte';
	import * as Select from '$lib/components/ui/select';
	import type { PageData } from './$houdini'
	import { onMount } from 'svelte';
	import { afterNavigate } from '$app/navigation';
	import { initFlowbite } from 'flowbite';
	import { page } from '$app/stores';
	import { get } from 'svelte/store';
	import { goto } from '$app/navigation';

    export let data: PageData
    let selectedImage: number = 0;
    $: ({ ProductQuery } = data)
	const phoneNumber = '09150610759';
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
		const variants = get(ProductQuery).data?.product?.variants;
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
	if ($ProductQuery.data?.product?.variants) {
		const optionsMap = new Map<string, Set<string>>();
		for (const variant of $ProductQuery.data.product.variants) {
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
		const product = $ProductQuery.data?.product;
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
		return $ProductQuery.data?.product?.variants.find((variant: any) =>
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

<section class="bg-white py-8 antialiased dark:bg-gray-900 md:py-16">
	<div class="mx-auto max-w-(--breakpoint-xl) px-4 2xl:px-0">
		<div class="lg: grid gap-4 md:grid-cols-[60%_40%] lg:gap-10">
			<div class="flex flex-col gap-4">
				<div class="gap-2 lg:grid lg:grid-cols-4">
					<div class="col-span-1 flex h-28 w-28 flex-col gap-4">
						{#each $ProductQuery.data?.product?.images ?? [] as img, index}
							<button class="rounded-md p-2 ring-[1.5px] ring-gray-400" on:click={() => selectedImage = index}>
								<img src={img.url} alt={img.altText} class="aspect-square w-full rounded-md object-cover" />
							</button>
						{/each}
					</div>
					<img src={$ProductQuery.data?.product?.images?.[selectedImage]?.url} alt={$ProductQuery.data?.product?.images?.[selectedImage]?.altText ?? 'Product'} class="col-span-4 col-start-2 aspect-square w-full rounded-md object-cover" height="300" width="300" />
				</div>
				<div class="">
					<Accordion.Root value="item-2" class="max-w-[70%] md:w-full md:max-w-full">
						<Accordion.Item value="item-1">
							<Accordion.Trigger>Product Details</Accordion.Trigger>
							<Accordion.Content>
								{$ProductQuery.data?.product?.description}
							</Accordion.Content>
						</Accordion.Item>
						<Accordion.Item value="item-2">
							<Accordion.Trigger>Specifications</Accordion.Trigger>
							<Accordion.Content>
								{#if $ProductQuery.data?.product?.attributes && $ProductQuery.data.product.attributes.length > 0}
									<ul class="list-disc list-inside space-y-1">
										{#each $ProductQuery.data.product.attributes as attr}
											<li><span class="font-semibold">{attr.title}:</span> {attr.value}</li>
										{/each}
									</ul>
								{:else}
									<p>No specifications available.</p>
								{/if}
							</Accordion.Content>
						</Accordion.Item>
					</Accordion.Root>
				</div>
			</div>
			<div class="mt-6 sm:mt-8 lg:mt-0">
				<Card.Root>
					<Card.Content>
						<h1 class="text-xl font-semibold text-gray-900 dark:text-white sm:text-2xl">
							{$ProductQuery.data?.product?.title}
						</h1>
						<div class="mt-4 sm:flex sm:items-center sm:gap-4">
							<p
								class="flex items-center rounded-sm bg-red-600 px-2 py-1 text-xs font-medium leading-none text-white"
							>
								the last 2 products
							</p>
							<div class="mt-2 flex items-center gap-2 sm:mt-0">
								<div class="flex items-center gap-1">
									<svg
										class="h-4 w-4 text-yellow-300"
										aria-hidden="true"
										xmlns="http://www.w3.org/2000/svg"
										width="24"
										height="24"
										fill="currentColor"
										viewBox="0 0 24 24"
									>
										<path
											d="M13.849 4.22c-.684-1.626-3.014-1.626-3.698 0L8.397 8.387l-4.552.361c-1.775.14-2.495 2.331-1.142 3.477l3.468 2.937-1.06 4.392c-.413 1.713 1.472 3.067 2.992 2.149L12 19.35l3.897 2.354c1.52.918 3.405-.436 2.992-2.15l-1.06-4.39 3.468-2.938c1.353-1.146.633-3.336-1.142-3.477l-4.552-.36-1.754-4.17Z"
										/>
									</svg>
									<svg
										class="h-4 w-4 text-yellow-300"
										aria-hidden="true"
										xmlns="http://www.w3.org/2000/svg"
										width="24"
										height="24"
										fill="currentColor"
										viewBox="0 0 24 24"
									>
										<path
											d="M13.849 4.22c-.684-1.626-3.014-1.626-3.698 0L8.397 8.387l-4.552.361c-1.775.14-2.495 2.331-1.142 3.477l3.468 2.937-1.06 4.392c-.413 1.713 1.472 3.067 2.992 2.149L12 19.35l3.897 2.354c1.52.918 3.405-.436 2.992-2.15l-1.06-4.39 3.468-2.938c1.353-1.146.633-3.336-1.142-3.477l-4.552-.36-1.754-4.17Z"
										/>
									</svg>
									<svg
										class="h-4 w-4 text-yellow-300"
										aria-hidden="true"
										xmlns="http://www.w3.org/2000/svg"
										width="24"
										height="24"
										fill="currentColor"
										viewBox="0 0 24 24"
									>
										<path
											d="M13.849 4.22c-.684-1.626-3.014-1.626-3.698 0L8.397 8.387l-4.552.361c-1.775.14-2.495 2.331-1.142 3.477l3.468 2.937-1.06 4.392c-.413 1.713 1.472 3.067 2.992 2.149L12 19.35l3.897 2.354c1.52.918 3.405-.436 2.992-2.15l-1.06-4.39 3.468-2.938c1.353-1.146.633-3.336-1.142-3.477l-4.552-.36-1.754-4.17Z"
										/>
									</svg>
									<svg
										class="h-4 w-4 text-yellow-300"
										aria-hidden="true"
										xmlns="http://www.w3.org/2000/svg"
										width="24"
										height="24"
										fill="currentColor"
										viewBox="0 0 24 24"
									>
										<path
											d="M13.849 4.22c-.684-1.626-3.014-1.626-3.698 0L8.397 8.387l-4.552.361c-1.775.14-2.495 2.331-1.142 3.477l3.468 2.937-1.06 4.392c-.413 1.713 1.472 3.067 2.992 2.149L12 19.35l3.897 2.354c1.52.918 3.405-.436 2.992-2.15l-1.06-4.39 3.468-2.938c1.353-1.146.633-3.336-1.142-3.477l-4.552-.36-1.754-4.17Z"
										/>
									</svg>
									<svg
										class="h-4 w-4 text-yellow-300"
										aria-hidden="true"
										xmlns="http://www.w3.org/2000/svg"
										width="24"
										height="24"
										fill="currentColor"
										viewBox="0 0 24 24"
									>
										<path
											d="M13.849 4.22c-.684-1.626-3.014-1.626-3.698 0L8.397 8.387l-4.552.361c-1.775.14-2.495 2.331-1.142 3.477l3.468 2.937-1.06 4.392c-.413 1.713 1.472 3.067 2.992 2.149L12 19.35l3.897 2.354c1.52.918 3.405-.436 2.992-2.15l-1.06-4.39 3.468-2.938c1.353-1.146.633-3.336-1.142-3.477l-4.552-.36-1.754-4.17Z"
										/>
									</svg>
								</div>

								<a
									href="/"
									class="text-sm font-medium leading-none text-gray-900 underline hover:no-underline dark:text-white"
								>
									345 Reviews
								</a>
							</div>
						</div>
						<p class="mt-4 flex items-center gap-2 text-sm font-medium leading-none text-blue-600">
							<MapPin class="blue-600" /> Deliver to Bonnie Green- Sacramento 23647
						</p>
						<div class="mt-6 flex items-center justify-between">
							<p class="text-2xl font-extrabold text-gray-900 dark:text-white sm:text-3xl">
								{#if activeVariant}
									{'$' + activeVariant.price.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
								{/if}
							</p>
							<div class="flex items-center gap-4">
								<p class="flex items-center gap-1 text-sm font-semibold">
									Quantity <Info class="h-5 w-4" />
								</p>
								<input
									type="number"
									min="1"
									max={activeVariant?.availableQuantity ?? 1}
									step="1"
									bind:value={quantity}
									class="w-20 p-2 border border-gray-300 rounded-md text-center shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
								/>
							</div>
						</div>

						<div class="mt-6 flex flex-col sm:mt-8 sm:flex sm:items-center sm:gap-4">
							<a
								href="/"
								title=""
								class="hover:text-primary-700 flex w-full items-center justify-center rounded-lg border border-gray-200 bg-white px-5 py-2.5 text-sm font-medium text-gray-900 hover:bg-gray-100 focus:z-10 focus:outline-hidden focus:ring-4 focus:ring-gray-100 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white dark:focus:ring-gray-700"
								role="button"
							>
								<svg
									class="-ms-2 me-2 h-5 w-5"
									aria-hidden="true"
									xmlns="http://www.w3.org/2000/svg"
									width="24"
									height="24"
									fill="none"
									viewBox="0 0 24 24"
								>
									<path
										stroke="currentColor"
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M12.01 6.001C6.5 1 1 8 5.782 13.001L12.011 20l6.23-7C23 8 17.5 1 12.01 6.002Z"
									/>
								</svg>
								Add to favorites
							</a>

							<a
								href={`https://wa.me/${phoneNumber}`}
								target="_blank"
								rel="noopener noreferrer"
								title=""
								class="hover:bg-primary-800 focus:ring-primary-300 dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800 mt-4 flex w-full items-center justify-center rounded-lg bg-green-600 px-5 py-2.5 text-sm font-medium text-white focus:outline-hidden focus:ring-4 sm:mt-0"
								role="button"
							>
								Add to Cart
							</a>
						</div>
						<p class="mt-6 text-sm">
							Also available at competitive prices from <a
								class="text-blue-600"
								href="http://retailauhtority.com">authorized retailers</a
							>, with optional Premium delivery for expedited shipping
						</p>

						<hr class="my-6 border-gray-200 dark:border-gray-800 md:my-8" />
						{#each attributeOptions as { title, values }}
							<div class="mt-6 flex w-full flex-col gap-2">
								<p class="text-lg font-bold">{title}</p>
								{#if title === 'Color'}
									<div class="flex gap-2">
										{#each values as value}
											{#if findVariantByAttributes({...selectedAttributes, Color: value})}
												<a
													href={constructVariantUrl(findVariantByAttributes({...selectedAttributes, Color: value}))}
													class="w-8 h-8 rounded-full ring-2 block"
													class:ring-primary-600={selectedAttributes[title] === value}
													class:ring-gray-300={selectedAttributes[title] !== value}
													style="background-color: {value}"
													aria-label={`View ${title} ${value}`}
												></a>
											{:else}
												<span
													class="w-8 h-8 rounded-full ring-2 opacity-50 block"
													style="background-color: {value}"
													aria-label={`Unavailable ${title} ${value}`}
												></span>
											{/if}
										{/each}
									</div>
								{:else}
									<Select.Root>
										<Select.Trigger class="w-full p-4">
											<Select.Value placeholder={selectedAttributes[title] || `Select ${title}`} />
										</Select.Trigger>
										<Select.Content>
											{#each values as value}
												<Select.Item value={value} on:click={() => handleAttributeSelect(title, value)}>{value}</Select.Item>
											{/each}
										</Select.Content>
									</Select.Root>
								{/if}
							</div>
						{/each}
					</Card.Content>
				</Card.Root>
			</div>
		</div>
	</div>
</section>
