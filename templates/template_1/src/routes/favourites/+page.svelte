<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Heart, ArrowLeft, Star, ShoppingCart } from 'lucide-svelte';
	
	// Mock favourite products
	const favourites = [
		{
			id: 1,
			name: 'Wireless Bluetooth Headphones',
			price: 79.99,
			originalPrice: 99.99,
			image: '/images/placeholder-product.jpg',
			rating: 4.5,
			inStock: true
		},
		{
			id: 2,
			name: 'Smart Watch Series 5',
			price: 299.99,
			originalPrice: null,
			image: '/images/placeholder-product.jpg',
			rating: 4.8,
			inStock: true
		},
		{
			id: 3,
			name: 'Organic Coffee Beans',
			price: 24.99,
			originalPrice: 29.99,
			image: '/images/placeholder-product.jpg',
			rating: 4.3,
			inStock: false
		}
	];

	function removeFavourite(id: number) {
		// In a real app, this would call an API
		console.log('Removing favourite:', id);
	}

	function addToCart(id: number) {
		// In a real app, this would add to cart
		console.log('Adding to cart:', id);
	}
</script>

<svelte:head>
	<title>My Favourites</title>
	<meta name="description" content="View and manage your favourite products" />
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-6xl">
	<div class="mb-8">
		<Button variant="ghost" class="mb-4">
			<a href="/account">
				<ArrowLeft class="h-4 w-4 mr-2" />
				Back to Account
			</a>
		</Button>
		<h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">My Favourites</h1>
		<p class="text-gray-600 dark:text-gray-400">Products you've saved for later</p>
	</div>

	{#if favourites.length > 0}
		<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
			{#each favourites as product}
				<Card.Root class="relative group">
					<div class="absolute top-4 right-4 z-10">
						<Button
							variant="ghost"
							size="icon"
							class="h-8 w-8 rounded-full bg-white/80 hover:bg-white dark:bg-gray-800/80 dark:hover:bg-gray-800"
							on:click={() => removeFavourite(product.id)}
						>
							<Heart class="h-4 w-4 text-red-500 fill-current" />
						</Button>
					</div>
					
					<div class="aspect-square bg-gray-100 dark:bg-gray-800 rounded-t-lg overflow-hidden">
						<img
							src={product.image}
							alt={product.name}
							class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
						/>
					</div>
					
					<Card.Content class="p-4">
						<h3 class="font-semibold text-gray-900 dark:text-white mb-2 line-clamp-2">
							{product.name}
						</h3>
						
						<div class="flex items-center gap-1 mb-2">
							<Star class="h-4 w-4 fill-yellow-400 text-yellow-400" />
							<span class="text-sm text-gray-600 dark:text-gray-400">{product.rating}</span>
						</div>
						
						<div class="flex items-center gap-2 mb-4">
							<span class="text-lg font-bold text-gray-900 dark:text-white">
								${product.price.toFixed(2)}
							</span>
							{#if product.originalPrice}
								<span class="text-sm text-gray-500 line-through">
									${product.originalPrice.toFixed(2)}
								</span>
							{/if}
						</div>
						
						{#if product.inStock}
							<Button 
								class="w-full"
								on:click={() => addToCart(product.id)}
							>
								<ShoppingCart class="h-4 w-4 mr-2" />
								Add to Cart
							</Button>
						{:else}
							<Button variant="outline" disabled class="w-full">
								Out of Stock
							</Button>
						{/if}
					</Card.Content>
				</Card.Root>
			{/each}
		</div>
	{:else}
		<Card.Root>
			<Card.Content class="text-center py-12">
				<Heart class="h-12 w-12 mx-auto text-gray-400 mb-4" />
				<h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">No favourites yet</h3>
				<p class="text-gray-600 dark:text-gray-400 mb-4">
					Start browsing and save products you love
				</p>
				<Button>
					<a href="/">Start Shopping</a>
				</Button>
			</Card.Content>
		</Card.Root>
	{/if}
</div>

<style>
	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
