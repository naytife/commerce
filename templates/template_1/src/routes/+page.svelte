<script lang="ts">
	import { onMount } from 'svelte';
	import { initFlowbite } from 'flowbite';
	import { cart } from '$lib/stores/cart';
	import type { PageData } from './$houdini';
	/** @type { import('./$houdini').PageData } */

    export let data: PageData

    $: ({ ProductsQuery } = data)
	
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
</script>

<section class="bg-gray-50 py-8 antialiased dark:bg-gray-900 md:py-12">
	<div class="mx-auto max-w-(--breakpoint-xl) px-4 2xl:px-0">
		<!-- Heading & Filters -->
		<div class="mb-4 items-end justify-between space-y-4 sm:flex sm:space-y-0 md:mb-8">
			<div>
				<nav class="flex" aria-label="Breadcrumb">
					<ol class="inline-flex items-center space-x-1 md:space-x-2 rtl:space-x-reverse">
						<li class="inline-flex items-center">
							<a
							aria-current="page"
								href="/"
								class="hover:text-primary-600 inline-flex items-center text-sm font-medium text-gray-700 dark:text-gray-400 dark:hover:text-white"
							>
								<svg
									class="me-2.5 h-3 w-3"
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
				<h2 class="mt-3 text-xl font-semibold text-gray-900 dark:text-white sm:text-2xl">
					Products
				</h2>
			</div>
			<div class="flex items-center space-x-4">
				<button
					data-modal-toggle="filterModal"
					data-modal-target="filterModal"
					type="button"
					class="hover:text-primary-700 flex w-full items-center justify-center rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm font-medium text-gray-900 hover:bg-gray-100 focus:z-10 focus:outline-hidden focus:ring-4 focus:ring-gray-100 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white dark:focus:ring-gray-700 sm:w-auto"
				>
					<svg
						class="-ms-0.5 me-2 h-4 w-4"
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
						class="-me-0.5 ms-2 h-4 w-4"
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
					class="hover:text-primary-700 flex w-full items-center justify-center rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm font-medium text-gray-900 hover:bg-gray-100 focus:z-10 focus:outline-hidden focus:ring-4 focus:ring-gray-100 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white dark:focus:ring-gray-700 sm:w-auto"
				>
					<svg
						class="-ms-0.5 me-2 h-4 w-4"
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
						class="-me-0.5 ms-2 h-4 w-4"
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
					class="z-50 hidden w-40 divide-y divide-gray-100 rounded-lg bg-white shadow-sm dark:bg-gray-700"
					data-popper-placement="bottom"
				>
					<ul
						class="p-2 text-left text-sm font-medium text-gray-500 dark:text-gray-400"
						aria-labelledby="sortDropdownButton"
					>
						<li>
							<a
								href="/"
								class="group inline-flex w-full items-center rounded-md px-3 py-2 text-sm text-gray-500 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-600 dark:hover:text-white"
							>
								The most popular
							</a>
						</li>
						<li>
							<a
								href="/"
								class="group inline-flex w-full items-center rounded-md px-3 py-2 text-sm text-gray-500 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-600 dark:hover:text-white"
							>
								Newest
							</a>
						</li>
						<li>
							<a
								href="/"
								class="group inline-flex w-full items-center rounded-md px-3 py-2 text-sm text-gray-500 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-600 dark:hover:text-white"
							>
								Increasing price
							</a>
						</li>
						<li>
							<a
								href="/"
								class="group inline-flex w-full items-center rounded-md px-3 py-2 text-sm text-gray-500 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-600 dark:hover:text-white"
							>
								Decreasing price
							</a>
						</li>
						<li>
							<a
								href="/"
								class="group inline-flex w-full items-center rounded-md px-3 py-2 text-sm text-gray-500 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-600 dark:hover:text-white"
							>
								No. reviews
							</a>
						</li>
						<li>
							<a
								href="/"
								class="group inline-flex w-full items-center rounded-md px-3 py-2 text-sm text-gray-500 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-600 dark:hover:text-white"
							>
								Discount %
							</a>
						</li>
					</ul>
				</div>
			</div>
		</div>
		<div class="mb-4 grid gap-4 sm:grid-cols-2 md:mb-8 lg:grid-cols-3 xl:grid-cols-4">
			{#each ($ProductsQuery.data?.products?.edges || []) as product, i (product?.node?.id)}
			{@const typedProduct = product as ProductEdge}
			<article class="overflow-hidden rounded-lg bg-white dark:bg-gray-800 shadow-lg hover:shadow-2xl transition flex flex-col">
				<div class="relative h-56 w-full">
					<img
						src={typedProduct.node.images[0]?.url || 'https://flowbite.s3.amazonaws.com/blocks/e-commerce/imac-front.svg'}
						alt={typedProduct.node.images[0]?.altText || typedProduct.node.title}
						class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
					/>
					<a href={constructProductUrl(typedProduct)} class="absolute inset-0"></a>
					{#if typedProduct.node.defaultVariant.stockStatus === 'OUT_OF_STOCK'}
						<span class="absolute top-2 left-2 bg-red-600 text-white text-xs font-bold px-2 py-1 rounded">Sold out</span>
					{/if}
				</div>
				<div class="flex-1 p-4 space-y-1">
					<h3 class="text-lg font-semibold text-gray-900 dark:text-white hover:text-primary-600 transition-colors">
						<a href={constructProductUrl(typedProduct)}>{typedProduct.node.title}</a>
					</h3>
					<p class="text-sm text-gray-500 dark:text-gray-400 line-clamp-2">{typedProduct.node.description}</p>
				</div>
				<footer class="px-4 pb-4 flex items-center justify-between">
					<p class="text-xl font-bold text-gray-900 dark:text-white">NGN {typedProduct.node.defaultVariant.price.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</p>
					<button
						type="button"
						on:click={() =>
							cart.add(
								{
									id: typedProduct.node.defaultVariant.variationId.toString(),
									title: typedProduct.node.title,
									price: typedProduct.node.defaultVariant.price,
									image: typedProduct.node.images[0]?.url,
									slug: typedProduct.node.slug
								},
								1
							)
						}
						class="inline-flex items-center gap-2 bg-gradient-to-r from-primary-600 to-primary-700 hover:from-primary-700 hover:to-primary-800 focus:outline-none focus:ring-4 focus:ring-primary-300 text-white px-4 py-2 rounded-lg shadow"
					>
						<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path d="M3 3h2l1 5h13l1-5h2m-6 15a2 2 0 11-4 0 2 2 0 014 0zm6 0a2 2 0 11-4 0 2 2 0 014 0z" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
						</svg>
						<span class="font-medium">Add to cart</span>
					</button>
				</footer>
			</article>
			{/each}
			
		</div>
		<div class="w-full text-center">
			<button
				type="button"
				class="hover:text-primary-700 rounded-lg border border-gray-200 bg-white px-5 py-2.5 text-sm font-medium text-gray-900 hover:bg-gray-100 focus:z-10 focus:outline-hidden focus:ring-4 focus:ring-gray-100 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white dark:focus:ring-gray-700"
				>Show more</button
			>
		</div>
	</div>
	<!-- Filter modal -->
	<form
		action="#"
		method="get"
		id="filterModal"
		tabindex="-1"
		aria-hidden="true"
		class="fixed left-0 right-0 top-0 z-50 hidden h-modal w-full overflow-y-auto overflow-x-hidden p-4 md:inset-0 md:h-full"
	>
		<div class="relative h-full w-full max-w-xl md:h-auto">
			<!-- Modal content -->
			<div class="relative rounded-lg bg-white shadow-sm dark:bg-gray-800">
				<!-- Modal header -->
				<div class="flex items-start justify-between rounded-t p-4 md:p-5">
					<h3 class="text-lg font-normal text-gray-500 dark:text-gray-400">Filters</h3>
					<button
						type="button"
						class="ml-auto inline-flex items-center rounded-lg bg-transparent p-1.5 text-sm text-gray-400 hover:bg-gray-100 hover:text-gray-900 dark:hover:bg-gray-600 dark:hover:text-white"
						data-modal-toggle="filterModal"
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
								d="M6 18 17.94 6M18 18 6.06 6"
							/>
						</svg>
						<span class="sr-only">Close modal</span>
					</button>
				</div>
				<!-- Modal body -->
				<div class="px-4 md:px-5">
					<div class="mb-4 border-b border-gray-200 dark:border-gray-700">
						<ul
							class="-mb-px flex flex-wrap text-center text-sm font-medium"
							id="myTab"
							data-tabs-toggle="#myTabContent"
							role="tablist"
						>
							<li class="mr-1" role="presentation">
								<button
									class="inline-block pb-2 pr-1"
									id="brand-tab"
									data-tabs-target="#brand"
									type="button"
									role="tab"
									aria-controls="profile"
									aria-selected="false">Brand</button
								>
							</li>
							<li class="mr-1" role="presentation">
								<button
									class="inline-block px-2 pb-2 hover:border-gray-300 hover:text-gray-600 dark:hover:text-gray-300"
									id="advanced-filers-tab"
									data-tabs-target="#advanced-filters"
									type="button"
									role="tab"
									aria-controls="advanced-filters"
									aria-selected="false">Advanced Filters</button
								>
							</li>
						</ul>
					</div>
					<div id="myTabContent">
						<div
							class="grid grid-cols-2 gap-4 md:grid-cols-3"
							id="brand"
							role="tabpanel"
							aria-labelledby="brand-tab"
						>
							<div class="space-y-2">
								<h5 class="text-lg font-medium uppercase text-black dark:text-white">A</h5>

								<div class="flex items-center">
									<input
										id="apple"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="apple"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Apple (56)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="asus"
										type="checkbox"
										value=""
										checked
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="asus"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Asus (97)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="acer"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="acer"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Acer (234)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="allview"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="allview"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Allview (45)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="atari"
										type="checkbox"
										value=""
										checked
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="asus"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Atari (176)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="amd"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="amd"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										AMD (49)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="aruba"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="aruba"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Aruba (16)
									</label>
								</div>
							</div>

							<div class="space-y-2">
								<h5 class="text-lg font-medium uppercase text-black dark:text-white">B</h5>

								<div class="flex items-center">
									<input
										id="beats"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="beats"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Beats (56)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="bose"
										type="checkbox"
										value=""
										checked
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="bose"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Bose (97)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="benq"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="benq"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										BenQ (45)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="bosch"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="bosch"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Bosch (176)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="brother"
										type="checkbox"
										value=""
										checked
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="brother"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Brother (176)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="biostar"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="biostar"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Biostar (49)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="braun"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="braun"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Braun (16)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="blaupunkt"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="blaupunkt"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Blaupunkt (45)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="benq2"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="benq2"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										BenQ (23)
									</label>
								</div>
							</div>

							<div class="space-y-2">
								<h5 class="text-lg font-medium uppercase text-black dark:text-white">C</h5>

								<div class="flex items-center">
									<input
										id="canon"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="canon"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Canon (49)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="cisco"
										type="checkbox"
										value=""
										checked
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="cisco"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Cisco (97)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="cowon"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="cowon"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Cowon (234)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="clevo"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="clevo"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Clevo (45)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="corsair"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="corsair"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Corsair (15)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="csl"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="csl"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Canon (49)
									</label>
								</div>
							</div>

							<div class="space-y-2">
								<h5 class="text-lg font-medium uppercase text-black dark:text-white">D</h5>

								<div class="flex items-center">
									<input
										id="dell"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="dell"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Dell (56)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="dogfish"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="dogfish"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Dogfish (24)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="dyson"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="dyson"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Dyson (234)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="dobe"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="dobe"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Dobe (5)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="digitus"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="digitus"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Digitus (1)
									</label>
								</div>
							</div>

							<div class="space-y-2">
								<h5 class="text-lg font-medium uppercase text-black dark:text-white">E</h5>

								<div class="flex items-center">
									<input
										id="emetec"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="emetec"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Emetec (56)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="extreme"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="extreme"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Extreme (10)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="elgato"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="elgato"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Elgato (234)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="emerson"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="emerson"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Emerson (45)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="emi"
										type="checkbox"
										value=""
										checked
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="emi"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										EMI (176)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="fugoo"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="fugoo"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Fugoo (49)
									</label>
								</div>
							</div>

							<div class="space-y-2">
								<h5 class="text-lg font-medium uppercase text-black dark:text-white">F</h5>

								<div class="flex items-center">
									<input
										id="fujitsu"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="fujitsu"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Fujitsu (97)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="fitbit"
										type="checkbox"
										value=""
										checked
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="fitbit"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Fitbit (56)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="foxconn"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="foxconn"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Foxconn (234)
									</label>
								</div>

								<div class="flex items-center">
									<input
										id="floston"
										type="checkbox"
										value=""
										class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
									/>

									<label
										for="floston"
										class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
									>
										Floston (45)
									</label>
								</div>
							</div>
						</div>
					</div>

					<div
						class="space-y-4"
						id="advanced-filters"
						role="tabpanel"
						aria-labelledby="advanced-filters-tab"
					>
						<div class="grid grid-cols-1 gap-8 md:grid-cols-2">
							<div class="grid grid-cols-2 gap-3">
								<div>
									<label
										for="min-price"
										class="block text-sm font-medium text-gray-900 dark:text-white"
									>
										Min Price
									</label>
									<input
										id="min-price"
										type="range"
										min="0"
										max="7000"
										value="300"
										step="1"
										class="h-2 w-full cursor-pointer appearance-none rounded-lg bg-gray-200 dark:bg-gray-700"
									/>
								</div>

								<div>
									<label
										for="max-price"
										class="block text-sm font-medium text-gray-900 dark:text-white"
									>
										Max Price
									</label>
									<input
										id="max-price"
										type="range"
										min="0"
										max="7000"
										value="3500"
										step="1"
										class="h-2 w-full cursor-pointer appearance-none rounded-lg bg-gray-200 dark:bg-gray-700"
									/>
								</div>

								<div class="col-span-2 flex items-center justify-between space-x-2">
									<input
										type="number"
										id="min-price-input"
										value="300"
										min="0"
										max="7000"
										class="focus:border-primary-500 focus:ring-primary-500 dark:focus:border-primary-500 dark:focus:ring-primary-500 block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder:text-gray-400"
										placeholder=""
										required
									/>

									<div class="shrink-0 text-sm font-medium dark:text-gray-300">to</div>

									<input
										type="number"
										id="max-price-input"
										value="3500"
										min="0"
										max="7000"
										class="focus:border-primary-500 focus:ring-primary-500 dark:focus:border-primary-500 dark:focus:ring-primary-500 block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder:text-gray-400"
										placeholder=""
										required
									/>
								</div>
							</div>

							<div class="space-y-3">
								<div>
									<label
										for="min-delivery-time"
										class="block text-sm font-medium text-gray-900 dark:text-white"
									>
										Min Delivery Time (Days)
									</label>

									<input
										id="min-delivery-time"
										type="range"
										min="3"
										max="50"
										value="30"
										step="1"
										class="h-2 w-full cursor-pointer appearance-none rounded-lg bg-gray-200 dark:bg-gray-700"
									/>
								</div>

								<input
									type="number"
									id="min-delivery-time-input"
									value="30"
									min="3"
									max="50"
									class="focus:border-primary-500 focus:ring-primary-500 dark:focus:border-primary-500 dark:focus:ring-primary-500 block w-full rounded-lg border border-gray-300 bg-gray-50 p-2.5 text-sm text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder:text-gray-400"
									placeholder=""
									required
								/>
							</div>
						</div>

						<div>
							<h6 class="mb-2 text-sm font-medium text-black dark:text-white">Condition</h6>

							<ul
								class="flex w-full items-center border border-gray-200 bg-white text-sm font-medium text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-white"
							>
								<li class="w-full border-r border-gray-200 dark:border-gray-600">
									<div class="flex items-center pl-3">
										<input
											id="condition-all"
											type="radio"
											value=""
											name="list-radio"
											checked
											class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-500 dark:bg-gray-600 dark:ring-offset-gray-700"
										/>
										<label
											for="condition-all"
											class="ml-2 w-full py-3 text-sm font-medium text-gray-900 dark:text-gray-300"
										>
											All
										</label>
									</div>
								</li>
								<li class="w-full border-r border-gray-200 dark:border-gray-600">
									<div class="flex items-center pl-3">
										<input
											id="condition-new"
											type="radio"
											value=""
											name="list-radio"
											class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-500 dark:bg-gray-600 dark:ring-offset-gray-700"
										/>
										<label
											for="condition-new"
											class="ml-2 w-full py-3 text-sm font-medium text-gray-900 dark:text-gray-300"
										>
											New
										</label>
									</div>
								</li>
								<li class="w-full">
									<div class="flex items-center pl-3">
										<input
											id="condition-used"
											type="radio"
											value=""
											name="list-radio"
											class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-500 dark:bg-gray-600 dark:ring-offset-gray-700"
										/>
										<label
											for="condition-used"
											class="ml-2 w-full py-3 text-sm font-medium text-gray-900 dark:text-gray-300"
										>
											Used
										</label>
									</div>
								</li>
							</ul>
						</div>

						<div class="grid grid-cols-2 gap-4 md:grid-cols-3">
							<div>
								<h6 class="mb-2 text-sm font-medium text-black dark:text-white">Colour</h6>
								<div class="space-y-2">
									<div class="flex items-center">
										<input
											id="blue"
											type="checkbox"
											value=""
											class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
										/>

										<label
											for="blue"
											class="ml-2 flex items-center text-sm font-medium text-gray-900 dark:text-gray-300"
										>
											<div class="bg-primary-600 mr-2 h-3.5 w-3.5 rounded-full"></div>
											Blue
										</label>
									</div>

									<div class="flex items-center">
										<input
											id="gray"
											type="checkbox"
											value=""
											class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
										/>

										<label
											for="gray"
											class="ml-2 flex items-center text-sm font-medium text-gray-900 dark:text-gray-300"
										>
											<div class="mr-2 h-3.5 w-3.5 rounded-full bg-gray-400"></div>
											Gray
										</label>
									</div>

									<div class="flex items-center">
										<input
											id="green"
											type="checkbox"
											value=""
											checked
											class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
										/>

										<label
											for="green"
											class="ml-2 flex items-center text-sm font-medium text-gray-900 dark:text-gray-300"
										>
											<div class="mr-2 h-3.5 w-3.5 rounded-full bg-green-400"></div>
											Green
										</label>
									</div>

									<div class="flex items-center">
										<input
											id="pink"
											type="checkbox"
											value=""
											class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
										/>

										<label
											for="pink"
											class="ml-2 flex items-center text-sm font-medium text-gray-900 dark:text-gray-300"
										>
											<div class="mr-2 h-3.5 w-3.5 rounded-full bg-pink-400"></div>
											Pink
										</label>
									</div>

									<div class="flex items-center">
										<input
											id="red"
											type="checkbox"
											value=""
											checked
											class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
										/>

										<label
											for="red"
											class="ml-2 flex items-center text-sm font-medium text-gray-900 dark:text-gray-300"
										>
											<div class="mr-2 h-3.5 w-3.5 rounded-full bg-red-500"></div>
											Red
										</label>
									</div>
								</div>
							</div>

							<div>
								<h6 class="mb-2 text-sm font-medium text-black dark:text-white">Rating</h6>
								<div class="space-y-2">
									<div class="flex items-center">
										<input
											id="five-stars"
											type="radio"
											value=""
											name="rating"
											class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
										/>
										<label for="five-stars" class="ml-2 flex items-center">
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-yellow-400"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>First star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-yellow-400"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Second star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-yellow-400"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Third star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-yellow-400"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Fourth star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-yellow-400"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Fifth star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
										</label>
									</div>

									<div class="flex items-center">
										<input
											id="four-stars"
											type="radio"
											value=""
											name="rating"
											class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
										/>
										<label for="four-stars" class="ml-2 flex items-center">
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-yellow-400"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>First star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-yellow-400"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Second star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-yellow-400"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Third star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-yellow-400"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Fourth star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-gray-300 dark:text-gray-500"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Fifth star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
										</label>
									</div>

									<div class="flex items-center">
										<input
											id="three-stars"
											type="radio"
											value=""
											name="rating"
											checked
											class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
										/>
										<label for="three-stars" class="ml-2 flex items-center">
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-yellow-400"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>First star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-yellow-400"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Second star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-yellow-400"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Third star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-gray-300 dark:text-gray-500"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Fourth star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-gray-300 dark:text-gray-500"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Fifth star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
										</label>
									</div>

									<div class="flex items-center">
										<input
											id="two-stars"
											type="radio"
											value=""
											name="rating"
											class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
										/>
										<label for="two-stars" class="ml-2 flex items-center">
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-yellow-400"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>First star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-yellow-400"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Second star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-gray-300 dark:text-gray-500"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Third star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-gray-300 dark:text-gray-500"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Fourth star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-gray-300 dark:text-gray-500"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Fifth star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
										</label>
									</div>

									<div class="flex items-center">
										<input
											id="one-star"
											type="radio"
											value=""
											name="rating"
											class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
										/>
										<label for="one-star" class="ml-2 flex items-center">
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-yellow-400"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>First star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-gray-300 dark:text-gray-500"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Second star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-gray-300 dark:text-gray-500"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Third star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-gray-300 dark:text-gray-500"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Fourth star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
											<svg
												aria-hidden="true"
												class="h-5 w-5 text-gray-300 dark:text-gray-500"
												fill="currentColor"
												viewBox="0 0 20 20"
												xmlns="http://www.w3.org/2000/svg"
											>
												<title>Fifth star</title>
												<path
													d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
												></path>
											</svg>
										</label>
									</div>
								</div>
							</div>

							<div>
								<h6 class="mb-2 text-sm font-medium text-black dark:text-white">Weight</h6>

								<div class="space-y-2">
									<div class="flex items-center">
										<input
											id="under-1-kg"
											type="checkbox"
											value=""
											class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
										/>

										<label
											for="under-1-kg"
											class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
										>
											Under 1 kg
										</label>
									</div>

									<div class="flex items-center">
										<input
											id="1-1-5-kg"
											type="checkbox"
											value=""
											checked
											class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
										/>

										<label
											for="1-1-5-kg"
											class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
										>
											1-1,5 kg
										</label>
									</div>

									<div class="flex items-center">
										<input
											id="1-5-2-kg"
											type="checkbox"
											value=""
											class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
										/>

										<label
											for="1-5-2-kg"
											class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
										>
											1,5-2 kg
										</label>
									</div>

									<div class="flex items-center">
										<input
											id="2-5-3-kg"
											type="checkbox"
											value=""
											class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
										/>

										<label
											for="2-5-3-kg"
											class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
										>
											2,5-3 kg
										</label>
									</div>

									<div class="flex items-center">
										<input
											id="over-3-kg"
											type="checkbox"
											value=""
											class="text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 h-4 w-4 rounded border-gray-300 bg-gray-100 focus:ring-2 dark:border-gray-600 dark:bg-gray-700 dark:ring-offset-gray-800"
										/>

										<label
											for="over-3-kg"
											class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
										>
											Over 3 kg
										</label>
									</div>
								</div>
							</div>
						</div>

						<div>
							<h6 class="mb-2 text-sm font-medium text-black dark:text-white">Delivery type</h6>

							<ul class="grid grid-cols-2 gap-4">
								<li>
									<input
										type="radio"
										id="delivery-usa"
										name="delivery"
										value="delivery-usa"
										class="peer hidden"
										checked
									/>
									<label
										for="delivery-usa"
										class="peer-checked:border-primary-600 peer-checked:text-primary-600 dark:peer-checked:text-primary-500 inline-flex w-full cursor-pointer items-center justify-between border border-gray-200 bg-white p-2 text-gray-500 hover:bg-gray-100 hover:text-gray-600 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-300 md:p-5"
									>
										<div class="block">
											<div class="w-full text-lg font-semibold">USA</div>
											<div class="w-full">Delivery only for USA</div>
										</div>
									</label>
								</li>
								<li>
									<input
										type="radio"
										id="delivery-europe"
										name="delivery"
										value="delivery-europe"
										class="peer hidden"
									/>
									<label
										for="delivery-europe"
										class="peer-checked:border-primary-600 peer-checked:text-primary-600 dark:peer-checked:text-primary-500 inline-flex w-full cursor-pointer items-center justify-between border border-gray-200 bg-white p-2 text-gray-500 hover:bg-gray-100 hover:text-gray-600 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-300 md:p-5"
									>
										<div class="block">
											<div class="w-full text-lg font-semibold">Europe</div>
											<div class="w-full">Delivery only for USA</div>
										</div>
									</label>
								</li>
								<li>
									<input
										type="radio"
										id="delivery-asia"
										name="delivery"
										value="delivery-asia"
										class="peer hidden"
										checked
									/>
									<label
										for="delivery-asia"
										class="peer-checked:border-primary-600 peer-checked:text-primary-600 dark:peer-checked:text-primary-500 inline-flex w-full cursor-pointer items-center justify-between border border-gray-200 bg-white p-2 text-gray-500 hover:bg-gray-100 hover:text-gray-600 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-300 md:p-5"
									>
										<div class="block">
											<div class="w-full text-lg font-semibold">Asia</div>
											<div class="w-full">Delivery only for Asia</div>
										</div>
									</label>
								</li>
								<li>
									<input
										type="radio"
										id="delivery-australia"
										name="delivery"
										value="delivery-australia"
										class="peer hidden"
									/>
									<label
										for="delivery-australia"
										class="peer-checked:border-primary-600 peer-checked:text-primary-600 dark:peer-checked:text-primary-500 inline-flex w-full cursor-pointer items-center justify-between border border-gray-200 bg-white p-2 text-gray-500 hover:bg-gray-100 hover:text-gray-600 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-gray-300 md:p-5"
									>
										<div class="block">
											<div class="w-full text-lg font-semibold">Australia</div>
											<div class="w-full">Delivery only for Australia</div>
										</div>
									</label>
								</li>
							</ul>
						</div>
					</div>
				</div>

				<!-- Modal footer -->
				<div class="flex items-center space-x-4 rounded-b p-4 dark:border-gray-600 md:p-5">
					<button
						type="submit"
						class="bg-primary-700 hover:bg-primary-800 focus:ring-primary-300 dark:bg-primary-700 dark:hover:bg-primary-800 dark:focus:ring-primary-800 rounded-lg px-5 py-2.5 text-center text-sm font-medium text-white focus:outline-hidden focus:ring-4"
						>Show 50 results</button
					>
					<button
						type="reset"
						class="hover:text-primary-700 border border-gray-200 bg-white px-5 py-2.5 text-sm font-medium text-gray-900 hover:bg-gray-100 focus:z-10 focus:outline-hidden focus:ring-4 focus:ring-gray-200 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white dark:focus:ring-gray-700"
						>Reset</button
					>
				</div>
			</div>
		</div>
	</form>
</section>
