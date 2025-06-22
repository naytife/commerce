<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { writable, derived } from 'svelte/store';
	import { allProducts, categories } from '$lib/data/store';
	import { Search, Filter, X, SlidersHorizontal } from 'lucide-svelte';
	import type { Product } from '$lib/types';
	
	const dispatch = createEventDispatcher();
	
	export let searchQuery = '';
	export let priceRange = [0, 1000];
	export let selectedCategories: string[] = [];
	export let sortBy = 'name';
	export let showFilters = false;

	// Create filters store for reactive filtering
	const filters = writable({
		query: '',
		categories: [] as string[],
		minPrice: 0,
		maxPrice: 1000,
		sortBy: 'name'
	});

	// Derived store for filtered products
	const filteredProducts = derived(
		[allProducts, filters],
		([$allProducts, $filters]) => {
			let products = $allProducts.filter((product: Product) => {
				// Text search
				if ($filters.query && !product.title.toLowerCase().includes($filters.query.toLowerCase())) {
					return false;
				}
				
				// Category filter
				if ($filters.categories.length > 0) {
					const categoryAttr = product.attributes.find(attr => attr.title === 'Category');
					const productCategory = categoryAttr?.value;
					if (!productCategory || !$filters.categories.includes(productCategory)) {
						return false;
					}
				}
				
				// Price filter
				if (product.defaultVariant.price < $filters.minPrice || product.defaultVariant.price > $filters.maxPrice) {
					return false;
				}
				
				return true;
			});

			// Apply sorting
			return products.sort((a, b) => {
				switch ($filters.sortBy) {
					case 'name-desc':
						return b.title.localeCompare(a.title);
					case 'price':
						return a.defaultVariant.price - b.defaultVariant.price;
					case 'price-desc':
						return b.defaultVariant.price - a.defaultVariant.price;
					case 'newest':
						return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
					case 'oldest':
						return new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime();
					default: // 'name'
						return a.title.localeCompare(b.title);
				}
			});
		}
	);
	
	const sortOptions = [
		{ value: 'name', label: 'Name (A-Z)' },
		{ value: 'name-desc', label: 'Name (Z-A)' },
		{ value: 'price', label: 'Price (Low to High)' },
		{ value: 'price-desc', label: 'Price (High to Low)' },
		{ value: 'newest', label: 'Newest First' },
		{ value: 'oldest', label: 'Oldest First' }
	];

	// Update filters when props change
	$: updateFilters();

	function updateFilters() {
		filters.update(f => ({
			...f,
			query: searchQuery,
			categories: selectedCategories,
			minPrice: priceRange[0],
			maxPrice: priceRange[1],
			sortBy
		}));
	}
	
	function handleSearch() {
		updateFilters();
		dispatch('search', {
			query: searchQuery,
			priceRange,
			categories: selectedCategories,
			sortBy,
			results: $filteredProducts
		});
	}
	
	function clearFilters() {
		searchQuery = '';
		priceRange = [0, 1000];
		selectedCategories = [];
		sortBy = 'name';
		updateFilters();
		dispatch('clear');
	}
	
	function toggleCategory(category: string) {
		if (selectedCategories.includes(category)) {
			selectedCategories = selectedCategories.filter(c => c !== category);
		} else {
			selectedCategories = [...selectedCategories, category];
		}
		handleSearch();
	}
	
	$: hasActiveFilters = searchQuery.length > 0 || selectedCategories.length > 0 || priceRange[0] > 0 || priceRange[1] < 1000;
</script>

<div class="bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-700 mb-8">
	<div class="max-w-7xl mx-auto px-6 py-6">
		<!-- Search Bar -->
		<div class="flex flex-col lg:flex-row gap-4 items-center">
			<div class="flex-1 relative">
				<Search class="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Search for products..."
					class="w-full pl-10 pr-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-800 dark:text-white transition-colors"
					on:input={handleSearch}
					on:keydown={(e) => e.key === 'Enter' && handleSearch()}
				/>
			</div>
			
			<!-- Sort & Filter Controls -->
			<div class="flex items-center gap-3">
				<select
					bind:value={sortBy}
					on:change={handleSearch}
					class="px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-800 dark:text-white"
				>
					{#each sortOptions as option}
						<option value={option.value}>{option.label}</option>
					{/each}
				</select>
				
				<button
					type="button"
					on:click={() => showFilters = !showFilters}
					class="flex items-center gap-2 px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors {hasActiveFilters ? 'bg-primary-50 border-primary-300 text-primary-700' : ''}"
				>
					<SlidersHorizontal class="w-4 h-4" />
					Filters
					{#if hasActiveFilters}
						<span class="bg-primary-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
							{selectedCategories.length + (priceRange[0] > 0 || priceRange[1] < 1000 ? 1 : 0)}
						</span>
					{/if}
				</button>
				
				{#if hasActiveFilters}
					<button
						type="button"
						on:click={clearFilters}
						class="flex items-center gap-1 px-3 py-2 text-sm text-gray-600 dark:text-gray-400 hover:text-red-600 dark:hover:text-red-400 transition-colors"
					>
						<X class="w-4 h-4" />
						Clear
					</button>
				{/if}
			</div>
		</div>
		
		<!-- Expandable Filters -->
		{#if showFilters}
			<div class="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700">
				<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
					<!-- Categories -->
					<div>
						<h3 class="text-sm font-medium text-gray-900 dark:text-white mb-3">Categories</h3>
						<div class="space-y-2">
							{#each $categories as category}
								<label class="flex items-center">
									<input
										type="checkbox"
										checked={selectedCategories.includes(category)}
										on:change={() => toggleCategory(category)}
										class="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
									/>
									<span class="ml-2 text-sm text-gray-700 dark:text-gray-300">{category}</span>
								</label>
							{/each}
						</div>
					</div>
					
					<!-- Price Range -->
					<div>
						<h3 class="text-sm font-medium text-gray-900 dark:text-white mb-3">Price Range</h3>
						<div class="space-y-3">
							<div class="flex items-center gap-2">
								<input
									type="number"
									bind:value={priceRange[0]}
									placeholder="Min"
									class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-800 dark:text-white"
									on:change={handleSearch}
								/>
								<span class="text-gray-500">-</span>
								<input
									type="number"
									bind:value={priceRange[1]}
									placeholder="Max"
									class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-800 dark:text-white"
									on:change={handleSearch}
								/>
							</div>
							<input
								type="range"
								min="0"
								max="1000"
								step="10"
								bind:value={priceRange[1]}
								class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer dark:bg-gray-700"
								on:input={handleSearch}
							/>
							<div class="flex justify-between text-xs text-gray-500">
								<span>${priceRange[0]}</span>
								<span>${priceRange[1]}</span>
							</div>
						</div>
					</div>
					
					<!-- Additional Filters -->
					<div>
						<h3 class="text-sm font-medium text-gray-900 dark:text-white mb-3">Availability</h3>
						<div class="space-y-2">
							<label class="flex items-center">
								<input
									type="checkbox"
									class="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
								/>
								<span class="ml-2 text-sm text-gray-700 dark:text-gray-300">In Stock</span>
							</label>
							<label class="flex items-center">
								<input
									type="checkbox"
									class="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
								/>
								<span class="ml-2 text-sm text-gray-700 dark:text-gray-300">On Sale</span>
							</label>
							<label class="flex items-center">
								<input
									type="checkbox"
									class="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
								/>
								<span class="ml-2 text-sm text-gray-700 dark:text-gray-300">Free Shipping</span>
							</label>
						</div>
					</div>
				</div>
			</div>
		{/if}
	</div>

	<!-- Results Summary -->
	<div class="mt-4 flex items-center justify-between text-sm text-gray-600 dark:text-gray-400">
		<div>
			Found <span class="font-semibold text-gray-900 dark:text-white">{$filteredProducts.length}</span> products
		</div>
		{#if $filteredProducts.length > 0}
			<div class="text-xs">
				Last updated: {new Date().toLocaleTimeString()}
			</div>
		{/if}
	</div>
</div>

<!-- Export filtered products for parent component to use -->
<slot filteredProducts={$filteredProducts} />
