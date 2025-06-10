<script lang="ts">
	import { onMount } from 'svelte';
	import { initFlowbite } from 'flowbite';
	import { cart } from '$lib/stores/cart';
	import { currencyCode, currencySymbol } from '$lib/stores/currency';
	import type { PageData } from './$houdini';
	/** @type { import('./$houdini').PageData } */

    export let data: PageData

    $: ({ ProductsQuery, ShopQuery } = data)
	
	onMount(() => {
		initFlowbite();
	});

	// Define interface for product structure based on GraphQL schema
	interface ProductAttribute {
		title: string;
		value: string;
	}

	interface ProductVariant {
		variationId: number;
		id: string;
		price: number;
		attributes: ProductAttribute[];
		availableQuantity: number;
		description: string;
		stockStatus: 'IN_STOCK' | 'OUT_OF_STOCK' | 'PREORDER';
	}

	interface ProductNode {
		id: string;
		productId: number;
		slug: string;
		title: string;
		description: string;
		images: { url: string; altText?: string }[];
		defaultVariant: ProductVariant;
	}

	interface ProductEdge {
		cursor: string;
		node: ProductNode;
	}

	function constructProductUrl(product: ProductEdge): string {
		const productId = product.node.productId;
		const attributesSlug = product.node.defaultVariant.attributes
			.sort((a: ProductAttribute, b: ProductAttribute) => a.title.localeCompare(b.title))
			.map((attr: ProductAttribute) => encodeURIComponent(attr.value.toLowerCase()))
			.join('-');
		const variationId = product.node.defaultVariant.variationId;
		return `/products/${productId}/${product.node.slug}-${attributesSlug}-${variationId}`;
	}

    $: if ($ShopQuery.data?.shop?.currencyCode) currencyCode.set($ShopQuery.data.shop.currencyCode);
</script>

<section class="bg-white py-12 antialiased dark:bg-gray-900">
	<div class="mx-auto max-w-7xl px-6">
		<!-- Heading & Filters -->
		<div class="mb-8 flex flex-col space-y-4 sm:flex-row sm:items-center sm:justify-between sm:space-y-0">
			<div>
				<nav class="flex" aria-label="Breadcrumb">
					<ol class="inline-flex items-center space-x-2">
						<li class="inline-flex items-center">
							<a
								href="/"
								class="text-sm text-gray-600 hover:text-primary-700 dark:text-gray-400 dark:hover:text-primary-500 transition-colors duration-200"
							>
								<svg
									class="mr-2 h-3.5 w-3.5"
									aria-hidden="true"
									xmlns="http://www.w3.org/2000/svg"
									fill="currentColor"
									viewBox="0 0 20 20"
								>
									<path
										d="m19.707 9.293-2-2-7-7a1 1 0 0 0-1.414 0l-7 7-2 2a1 1 0 0 0 1.414 1.414L2 10.414V18a2 2 0 0 0 2 2h3a1 1 0 0 0 1-1v-4a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1v4a1 1 0 0 0 1 1h3a2 2 0 0 0 2-2v-7.586l.293.293a1 1 0 0 0 1.414-1.414Z"
									/>
								</svg>
								Home
							</a>
						</li>
					</ol>
				</nav>
				<h2 class="mt-3 text-2xl font-medium text-gray-900 dark:text-white tracking-tight">
					Products
				</h2>
			</div>
			<div class="flex items-center space-x-4">
				<button
					data-modal-toggle="filterModal"
					data-modal-target="filterModal"
					type="button"
					class="inline-flex items-center border border-gray-200 bg-white px-4 py-2.5 text-sm font-medium text-gray-900 hover:bg-gray-50 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300 dark:hover:bg-gray-800 transition-colors duration-200"
				>
					<svg
						class="mr-2 h-4 w-4"
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
							stroke-width="2"
							d="M18.796 4H5.204a1 1 0 0 0-.753 1.659l5.302 6.058a1 1 0 0 1 .247.659v4.874a.5.5 0 0 0 .2.4l3 2.25a.5.5 0 0 0 .8-.4v-7.124a1 1 0 0 1 .247-.659l5.302-6.059c.566-.646.106-1.658-.753-1.658Z"
						/>
					</svg>
					Filters
					<svg
						class="ml-2 h-4 w-4"
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
							d="m19 9-7 7-7-7"
						/>
					</svg>
				</button>
				<button
					id="sortDropdownButton1"
					data-dropdown-toggle="dropdownSort1"
					type="button"
					class="inline-flex items-center border border-gray-200 bg-white px-4 py-2.5 text-sm font-medium text-gray-900 hover:bg-gray-50 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-300 dark:hover:bg-gray-800 transition-colors duration-200"
				>
					<svg
						class="mr-2 h-4 w-4"
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
							d="M7 4v16M7 4l3 3M7 4 4 7m9-3h6l-6 6h6m-6.5 10 3.5-7 3.5 7M14 18h4"
						/>
					</svg>
					Sort
					<svg
						class="ml-2 h-4 w-4"
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
							d="m19 9-7 7-7-7"
						/>
					</svg>
				</button>
				<div
					id="dropdownSort1"
					class="z-50 hidden w-40 divide-y divide-gray-100 bg-white shadow dark:divide-gray-700 dark:bg-gray-800"
				>
					<ul
						class="p-2 text-left text-sm font-medium text-gray-500 dark:text-gray-400"
						aria-labelledby="sortDropdownButton"
					>
						<li>
							<a
								href="/"
								class="block w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-700"
							>
								The most popular
							</a>
						</li>
						<li>
							<a
								href="/"
								class="block w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-700"
							>
								Newest
							</a>
						</li>
						<li>
							<a
								href="/"
								class="block w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-700"
							>
								Price: Low to High
							</a>
						</li>
						<li>
							<a
								href="/"
								class="block w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-gray-700"
							>
								Price: High to Low
							</a>
						</li>
					</ul>
				</div>
			</div>
		</div>

		<!-- Products Grid -->
		<div class="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
			{#each $ProductsQuery.data?.products?.edges ?? [] as product}
				{#if product.node}
					<div class="group relative border border-transparent hover:border-gray-200 dark:hover:border-gray-700 transition-colors duration-200">
						<a href={constructProductUrl(product)}>
							<div class="aspect-h-1 aspect-w-1 overflow-hidden bg-gray-50 dark:bg-gray-800">
								<img
									src={product.node.images?.[0]?.url ?? ''}
									alt={product.node.images?.[0]?.altText ?? product.node.title}
									class="h-full w-full object-cover object-center"
								/>
							</div>
							<div class="mt-4 px-2 pb-5">
								<h3 class="text-sm font-medium text-gray-900 dark:text-white">
									{product.node.title}
								</h3>
								<div class="mt-2 flex items-center justify-between">
									<p class="text-sm font-medium text-gray-900 dark:text-white">
										{$currencySymbol}{product.node.defaultVariant?.price?.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) ?? ''}
									</p>
									<div>
										{#if product.node.defaultVariant?.stockStatus === 'IN_STOCK'}
											<span class="inline-flex items-center text-xs font-medium text-green-700 dark:text-green-500">
												In Stock
											</span>
										{:else if product.node.defaultVariant?.stockStatus === 'OUT_OF_STOCK'}
											<span class="inline-flex items-center text-xs font-medium text-red-700 dark:text-red-500">
												Out of Stock
											</span>
										{:else if product.node.defaultVariant?.stockStatus === 'PREORDER'}
											<span class="inline-flex items-center text-xs font-medium text-blue-700 dark:text-blue-500">
												Pre-order
											</span>
										{/if}
									</div>
								</div>
							</div>
						</a>

						<div class="absolute inset-x-0 bottom-0 h-12 translate-y-full opacity-0 group-hover:translate-y-0 group-hover:opacity-100 transition-all duration-200 bg-white dark:bg-gray-900 border-t border-gray-200 dark:border-gray-700 flex items-center justify-center">
							<button
								on:click={() => cart.add({
									id: product.node.defaultVariant.variationId.toString(),
									title: product.node.title,
									price: product.node.defaultVariant.price,
									image: product.node.images?.[0]?.url ?? '',
									slug: `${product.node.slug}-${product.node.defaultVariant.attributes
										.sort((a: ProductAttribute, b: ProductAttribute) => a.title.localeCompare(b.title))
										.map((attr: ProductAttribute) => encodeURIComponent(attr.value.toLowerCase()))
										.join('-')}-${product.node.defaultVariant.variationId}`
								}, 1)}
								class="text-sm text-primary-700 font-medium hover:text-primary-800 dark:text-primary-500 dark:hover:text-primary-400 transition-colors duration-200"
							>
								Add to Cart
							</button>
						</div>
					</div>
				{/if}
			{/each}
		</div>

		<!-- Filters Modal -->
		<div
			id="filterModal"
			tabindex="-1"
			aria-hidden="true"
			class="fixed left-0 right-0 top-0 z-50 hidden h-modal h-full w-full overflow-y-auto overflow-x-hidden p-4 md:inset-0 md:h-full"
		>
			<div class="relative h-full w-full max-w-lg md:h-auto">
				<div class="bg-white dark:bg-gray-800 shadow">
					<div class="flex items-center justify-between p-4 md:p-5 border-b border-gray-200 dark:border-gray-700">
						<h3 class="text-lg font-semibold text-gray-900 dark:text-white">
							Filter Products
						</h3>
						<button
							type="button"
							class="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 text-sm dark:hover:bg-gray-700 dark:hover:text-white"
							data-modal-hide="filterModal"
						>
							<svg
								class="h-5 w-5"
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
									d="M6 18L18 6M6 6l12 12"
								/>
							</svg>
							<span class="sr-only">Close modal</span>
						</button>
					</div>
					<div class="p-4 md:p-5">
						<div class="mb-4 space-y-4">
							<h4 class="text-sm font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">
								Categories
							</h4>
							<ul class="space-y-2 text-sm text-gray-900 dark:text-white">
								<li>
									<label class="relative flex cursor-pointer items-center ">
										<input
											type="checkbox"
											value=""
											class="peer sr-only"
										/>
										<div
											class="peer relative h-5 w-5 shrink-0 border border-gray-300 bg-white focus:outline-none focus:ring-0 peer-checked:bg-primary-600 peer-checked:border-primary-600 dark:border-gray-600 dark:bg-gray-700 peer-checked:after:absolute peer-checked:after:left-1/2 peer-checked:after:top-1/2 peer-checked:after:h-2 peer-checked:after:w-[5px] peer-checked:after:-translate-x-1/2 peer-checked:after:-translate-y-1/2 peer-checked:after:rotate-45 peer-checked:after:border-r-2 peer-checked:after:border-b-2 peer-checked:after:border-white"
										/>
										<span class="ml-3 text-sm font-medium text-gray-900 dark:text-gray-300"
											>Electronics</span
										>
									</label>
								</li>
								<li>
									<label class="relative flex cursor-pointer items-center ">
										<input
											type="checkbox"
											value=""
											class="peer sr-only"
										/>
										<div
											class="peer relative h-5 w-5 shrink-0 border border-gray-300 bg-white focus:outline-none focus:ring-0 peer-checked:bg-primary-600 peer-checked:border-primary-600 dark:border-gray-600 dark:bg-gray-700 peer-checked:after:absolute peer-checked:after:left-1/2 peer-checked:after:top-1/2 peer-checked:after:h-2 peer-checked:after:w-[5px] peer-checked:after:-translate-x-1/2 peer-checked:after:-translate-y-1/2 peer-checked:after:rotate-45 peer-checked:after:border-r-2 peer-checked:after:border-b-2 peer-checked:after:border-white"
										/>
										<span class="ml-3 text-sm font-medium text-gray-900 dark:text-gray-300"
											>Fashion</span
										>
									</label>
								</li>
								<li>
									<label class="relative flex cursor-pointer items-center ">
										<input
											type="checkbox"
											value=""
											class="peer sr-only"
										/>
										<div
											class="peer relative h-5 w-5 shrink-0 border border-gray-300 bg-white focus:outline-none focus:ring-0 peer-checked:bg-primary-600 peer-checked:border-primary-600 dark:border-gray-600 dark:bg-gray-700 peer-checked:after:absolute peer-checked:after:left-1/2 peer-checked:after:top-1/2 peer-checked:after:h-2 peer-checked:after:w-[5px] peer-checked:after:-translate-x-1/2 peer-checked:after:-translate-y-1/2 peer-checked:after:rotate-45 peer-checked:after:border-r-2 peer-checked:after:border-b-2 peer-checked:after:border-white"
										/>
										<span class="ml-3 text-sm font-medium text-gray-900 dark:text-gray-300"
											>Home & Garden</span
										>
									</label>
								</li>
							</ul>
						</div>
						<div class="mb-4 space-y-4">
							<h4 class="text-sm font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">
								Price Range
							</h4>
							<div>
								<label for="price-range" class="sr-only">Price Range</label>
								<input id="price-range" type="range" min="0" max="1000" value="500" class="w-full h-2 bg-gray-200 appearance-none dark:bg-gray-700" />
								<div class="flex justify-between mt-2 text-xs text-gray-700 dark:text-gray-300">
									<span>$0</span>
									<span>$500</span>
									<span>$1000</span>
								</div>
							</div>
						</div>
						<div class="mt-6 flex items-center justify-between space-x-4">
							<button
								data-modal-hide="filterModal"
								type="button"
								class="w-1/2 border border-gray-200 bg-white px-5 py-2.5 text-sm font-medium text-gray-500 hover:text-gray-700 dark:border-gray-700 dark:bg-transparent dark:text-gray-300 dark:hover:text-white"
							>
								Clear All
							</button>
							<button
								type="button"
								class="w-1/2 bg-primary-700 px-5 py-2.5 text-center text-sm font-medium text-white hover:bg-primary-800 dark:bg-primary-600 dark:hover:bg-primary-700"
							>
								Apply Filters
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>
