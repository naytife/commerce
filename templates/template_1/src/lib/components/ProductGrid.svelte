<script lang="ts">
  import { featuredProducts, productsByCategory, allProducts } from '$lib/data/store'
  import type { Product } from '$lib/types'
  
  export let products: Product[] = []
  export let categoryId: string | null = null
  export let showFeatured = false
  export let showAllIfEmpty = false
  export let limit: number | null = null
  
  // Reactive product selection based on props
  $: displayProducts = products.length > 0 
    ? products 
    : categoryId 
    ? $productsByCategory.get(categoryId) || []
    : showFeatured 
    ? $featuredProducts 
    : showAllIfEmpty
    ? $allProducts
    : []
  
  // Apply limit if specified
  $: limitedProducts = limit ? displayProducts.slice(0, limit) : displayProducts
  
  function getOptimizedImageUrl(url: string, width = 400) {
    // Add image optimization parameters
    return `${url}&w=${width}&q=80&fm=webp`
  }
</script>

<div class="product-grid">
  {#if limitedProducts.length === 0}
    <div class="col-span-full text-center py-12">
      <div class="text-gray-500 dark:text-gray-400">
        <svg class="mx-auto h-12 w-12 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2 2v-5m16 0H4"/>
        </svg>
        <p class="text-lg font-medium">No products found</p>
        <p class="text-sm">Try adjusting your search or filters</p>
      </div>
    </div>
  {:else}
    {#each limitedProducts as product (product.id)}
      <div class="product-card group">
        <!-- Product Image -->
        <div class="product-image-container">
          {#if product.images?.[0]}
            <img 
              src={getOptimizedImageUrl(product.images[0].url)} 
              alt={product.images[0].altText || product.title}
              loading="lazy"
              decoding="async"
              class="product-image"
            />
          {:else}
            <div class="product-placeholder">
              <svg class="w-12 h-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
              </svg>
            </div>
          {/if}
          
          <!-- Product badges -->
          <div class="product-badges">
            {#if product.attributes.some(attr => attr.title === 'Featured' && attr.value === 'true')}
              <span class="badge badge-featured">Featured</span>
            {/if}
          </div>
          
          <!-- Quick view overlay -->
          <div class="product-overlay">
            <button class="btn btn-primary btn-sm">Quick View</button>
          </div>
        </div>
        
        <!-- Product Info -->
        <div class="product-info">
          <h3 class="product-title">
            <a href="/products/{product.productId}/{product.slug}" class="product-link">
              {product.title}
            </a>
          </h3>
          
          <p class="product-description">
            {product.description}
          </p>
          
          <!-- Price -->
          <div class="product-pricing">
            <span class="price-current">${product.defaultVariant.price.toFixed(2)}</span>
          </div>
          
          <!-- Stock Status -->
          <div class="product-stock">
            {#if product.defaultVariant.availableQuantity <= 0}
              <span class="stock-status stock-out">Out of Stock</span>
            {:else if product.defaultVariant.availableQuantity <= 5}
              <span class="stock-status stock-low">Only {product.defaultVariant.availableQuantity} left</span>
            {:else}
              <span class="stock-status stock-in">In Stock</span>
            {/if}
          </div>
          
          <!-- Actions -->
          <div class="product-actions">
            {#if product.defaultVariant.availableQuantity > 0}
              <button class="btn btn-primary btn-full">
                Add to Cart
              </button>
            {:else}
              <button class="btn btn-secondary btn-full" disabled>
                Notify When Available
              </button>
            {/if}
          </div>
        </div>
      </div>
    {/each}
  {/if}
</div>

<style>
  .product-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 1.5rem;
    padding: 1rem 0;
  }
  
  .product-card {
    @apply bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden shadow-sm hover:shadow-lg transition-all duration-300;
  }
  
  .product-image-container {
    @apply relative aspect-square overflow-hidden bg-gray-100 dark:bg-gray-700;
  }
  
  .product-image {
    @apply w-full h-full object-cover transition-transform duration-300 group-hover:scale-105;
  }
  
  .product-placeholder {
    @apply w-full h-full flex items-center justify-center;
  }
  
  .product-badges {
    @apply absolute top-2 left-2 flex flex-col gap-1;
  }
  
  .badge {
    @apply px-2 py-1 text-xs font-semibold rounded;
  }
  
  .badge-featured {
    @apply bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200;
  }
  
  .badge-sale {
    @apply bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200;
  }
  
  .product-overlay {
    @apply absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300;
  }
  
  .product-info {
    @apply p-4 space-y-3;
  }
  
  .product-title {
    @apply text-lg font-semibold;
  }
  
  .product-link {
    @apply text-gray-900 dark:text-white hover:text-blue-600 dark:hover:text-blue-400 transition-colors;
  }
  
  .product-description {
    @apply text-sm text-gray-600 dark:text-gray-300 line-clamp-2;
  }
  
  .product-pricing {
    @apply flex items-center gap-2;
  }
  
  .price-current {
    @apply text-xl font-bold text-blue-600 dark:text-blue-400;
  }
  
  .price-original {
    @apply text-sm text-gray-500 dark:text-gray-400 line-through;
  }
  
  .product-stock {
    @apply text-sm;
  }
  
  .stock-status {
    @apply font-medium;
  }
  
  .stock-in {
    @apply text-green-600 dark:text-green-400;
  }
  
  .stock-low {
    @apply text-orange-600 dark:text-orange-400;
  }
  
  .stock-out {
    @apply text-red-600 dark:text-red-400;
  }
  
  .product-actions {
    @apply mt-4;
  }
  
  .btn {
    @apply px-4 py-2 rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2;
  }
  
  .btn-primary {
    @apply bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500;
  }
  
  .btn-secondary {
    @apply bg-gray-200 text-gray-800 hover:bg-gray-300 focus:ring-gray-500 dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600;
  }
  
  .btn-sm {
    @apply px-3 py-1 text-sm;
  }
  
  .btn-full {
    @apply w-full;
  }
  
  .btn:disabled {
    @apply opacity-50 cursor-not-allowed;
  }
  
  @media (max-width: 640px) {
    .product-grid {
      grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
      gap: 1rem;
    }
  }
</style>
