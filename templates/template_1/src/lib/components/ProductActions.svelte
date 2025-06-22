<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { Heart, GitCompare, X, Check } from 'lucide-svelte';
	import { writable } from 'svelte/store';
	
	const dispatch = createEventDispatcher();
	
	// Simple stores for wishlist and comparison (in real app these would be persistent)
	export const wishlistStore = writable<string[]>([]);
	export const comparisonStore = writable<string[]>([]);
	
	export let productId: string;
	export let productTitle: string;
	export const productImage: string = '';
	export const productPrice: number = 0;
	export let size: 'small' | 'medium' = 'medium';
	
	let isInWishlist = false;
	let isInComparison = false;
	let showFeedback = false;
	
	// Subscribe to stores to check if product is already added
	wishlistStore.subscribe((items) => {
		isInWishlist = items.includes(productId);
	});
	
	comparisonStore.subscribe((items) => {
		isInComparison = items.includes(productId);
	});
	
	function toggleWishlist() {
		wishlistStore.update((items) => {
			if (items.includes(productId)) {
				return items.filter(id => id !== productId);
			} else {
				return [...items, productId];
			}
		});
		
		showFeedback = true;
		setTimeout(() => showFeedback = false, 2000);
		
		dispatch('wishlistToggle', {
			productId,
			isInWishlist: !isInWishlist,
			productTitle
		});
	}
	
	function toggleComparison() {
		comparisonStore.update((items) => {
			if (items.includes(productId)) {
				return items.filter(id => id !== productId);
			} else {
				// Limit comparison to 4 items
				if (items.length >= 4) {
					return [productId, ...items.slice(0, 3)];
				}
				return [...items, productId];
			}
		});
		
		dispatch('comparisonToggle', {
			productId,
			isInComparison: !isInComparison,
			productTitle
		});
	}
	
	$: iconSize = size === 'small' ? 'w-4 h-4' : 'w-5 h-5';
	$: buttonSize = size === 'small' ? 'p-2' : 'p-2.5';
</script>

<div class="flex items-center gap-2 relative">
	<!-- Wishlist Button -->
	<button
		type="button"
		on:click={toggleWishlist}
		class="relative {buttonSize} rounded-full border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 transition-all duration-200 {isInWishlist ? 'text-red-500 border-red-300 bg-red-50 dark:bg-red-900/20' : 'text-gray-600 dark:text-gray-400'}"
		title={isInWishlist ? 'Remove from Wishlist' : 'Add to Wishlist'}
	>
		<Heart class="{iconSize} {isInWishlist ? 'fill-current' : ''}" />
		
		{#if showFeedback}
			<div class="absolute -top-8 left-1/2 transform -translate-x-1/2 bg-gray-900 text-white text-xs px-2 py-1 rounded shadow-lg whitespace-nowrap">
				{isInWishlist ? 'Added to' : 'Removed from'} wishlist
				<div class="absolute top-full left-1/2 transform -translate-x-1/2 w-0 h-0 border-l-2 border-r-2 border-t-2 border-transparent border-t-gray-900"></div>
			</div>
		{/if}
	</button>
	
	<!-- Comparison Button -->
	<button
		type="button"
		on:click={toggleComparison}
		class="relative {buttonSize} rounded-full border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 transition-all duration-200 {isInComparison ? 'text-blue-500 border-blue-300 bg-blue-50 dark:bg-blue-900/20' : 'text-gray-600 dark:text-gray-400'}"
		title={isInComparison ? 'Remove from Comparison' : 'Add to Comparison'}
	>
		<GitCompare class="{iconSize}" />
		
		{#if isInComparison}
			<div class="absolute -top-1 -right-1 w-3 h-3 bg-blue-500 text-white text-xs rounded-full flex items-center justify-center">
				<Check class="w-2 h-2" />
			</div>
		{/if}
	</button>
</div>

<!-- Comparison Count Indicator (Global) -->
{#if $comparisonStore.length > 0}
	<div class="fixed bottom-4 right-4 z-50">
		<button
			type="button"
			on:click={() => dispatch('viewComparison')}
			class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg shadow-lg flex items-center gap-2 transition-colors"
		>
			<GitCompare class="w-4 h-4" />
			Compare ({$comparisonStore.length})
		</button>
	</div>
{/if}
