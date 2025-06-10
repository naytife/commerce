<script lang="ts">
  import { cart } from '$lib/stores/cart';
  import { currencySymbol } from '$lib/stores/currency';
  import { derived } from 'svelte/store';
  import { goto } from '$app/navigation';
  import { Trash2 } from 'lucide-svelte';

  const total = derived(cart, items =>
    items.reduce((sum, item) => sum + item.price * item.quantity, 0)
  );
</script>

<section class="bg-white py-12 antialiased dark:bg-gray-900">
  <div class="mx-auto max-w-7xl px-6">
    <div class="border-b border-gray-200 dark:border-gray-700 pb-6 mb-8">
      <h1 class="text-2xl font-medium text-gray-900 dark:text-white tracking-tight">Shopping Cart</h1>
    </div>
    
    {#if $cart.length > 0}
      <div class="mb-8">
        <table class="w-full text-sm text-left">
          <thead class="text-xs uppercase tracking-wider border-b border-gray-200 dark:border-gray-700 text-gray-700 dark:text-gray-400">
            <tr>
              <th scope="col" class="px-6 py-4 font-medium">Product</th>
              <th scope="col" class="px-6 py-4 font-medium text-right">Price</th>
              <th scope="col" class="px-6 py-4 font-medium text-center">Quantity</th>
              <th scope="col" class="px-6 py-4 font-medium text-right">Total</th>
              <th scope="col" class="px-6 py-4 font-medium text-center">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
            {#each $cart as item (item.id)}
              <tr class="bg-white dark:bg-gray-900">
                <td class="px-6 py-6">
                  <div class="flex items-center gap-4">
                    {#if item.image}
                      <div class="w-20 h-20 bg-gray-100 dark:bg-gray-800 flex-shrink-0">
                        <img src={item.image} alt={item.title} class="w-full h-full object-cover" />
                      </div>
                    {/if}
                    <a href={`/products/${item.id}/${item.slug}`} class="font-medium text-gray-900 dark:text-white hover:text-primary-700 dark:hover:text-primary-500 transition-colors duration-200">{item.title}</a>
                  </div>
                </td>
                <td class="px-6 py-6 text-right font-medium text-gray-700 dark:text-gray-300">{$currencySymbol}{item.price.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</td>
                <td class="px-6 py-6">
                  <div class="flex justify-center">
                    <div class="flex border border-gray-200 dark:border-gray-700">
                      <button 
                        on:click={() => cart.updateQuantity(item.id, Math.max(1, item.quantity - 1))}
                        class="w-10 h-10 flex items-center justify-center text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors duration-200"
                      >-</button>
                      <input
                        type="number"
                        min="1"
                        value={item.quantity}
                        on:change={(e: Event) => cart.updateQuantity(item.id, +(e.target as HTMLInputElement).value)}
                        class="w-16 h-10 text-center border-x border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:outline-none"
                      />
                      <button 
                        on:click={() => cart.updateQuantity(item.id, item.quantity + 1)}
                        class="w-10 h-10 flex items-center justify-center text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors duration-200"
                      >+</button>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-6 text-right font-medium text-gray-900 dark:text-white">{$currencySymbol}{(item.price * item.quantity).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</td>
                <td class="px-6 py-6">
                  <div class="flex justify-center">
                    <button 
                      on:click={() => cart.remove(item.id)} 
                      class="text-gray-400 hover:text-red-600 dark:text-gray-500 dark:hover:text-red-500 transition-colors duration-200"
                    >
                      <span class="sr-only">Remove</span>
                      <Trash2 class="w-5 h-5" />
                    </button>
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
      
      <div class="border-t border-gray-200 dark:border-gray-700 pt-6">
        <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-6">
          <div class="w-full md:w-auto">
            <a href="/" class="inline-flex items-center text-sm text-gray-600 hover:text-primary-700 dark:text-gray-400 dark:hover:text-primary-500 transition-colors duration-200">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
              </svg>
              Continue Shopping
            </a>
          </div>
          
          <div class="w-full md:w-auto flex flex-col items-end">
            <div class="flex justify-between items-center gap-x-8 mb-4">
              <span class="text-base font-medium text-gray-700 dark:text-gray-300">Subtotal:</span>
              <span class="text-xl font-semibold text-gray-900 dark:text-white">{$currencySymbol}{$total.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
            </div>
            <p class="text-xs text-gray-500 dark:text-gray-400 mb-4 text-right">Shipping and taxes calculated at checkout</p>
            <button 
              on:click={() => goto('/checkout')} 
              class="bg-primary-700 hover:bg-primary-800 focus:ring-4 focus:ring-primary-300 focus:ring-offset-2 dark:focus:ring-offset-gray-900 text-white font-medium tracking-wide px-8 py-3 transition-colors duration-200"
            >
              Proceed to Checkout
            </button>
          </div>
        </div>
      </div>
    {:else}
      <div class="text-center py-12 border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800">
        <p class="text-gray-500 dark:text-gray-400 mb-6">Your cart is empty.</p>
        <a href="/" class="inline-flex items-center text-sm font-medium text-primary-700 hover:text-primary-800 dark:text-primary-500 dark:hover:text-primary-400 transition-colors duration-200">
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
          </svg>
          Continue Shopping
        </a>
      </div>
    {/if}
  </div>
</section> 