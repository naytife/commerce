<script lang="ts">
  import { cart } from '$lib/stores/cart';
  import { derived } from 'svelte/store';
  import { goto } from '$app/navigation';

  const total = derived(cart, items =>
    items.reduce((sum, item) => sum + item.price * item.quantity, 0)
  );
</script>

<section class="bg-white py-8 antialiased dark:bg-gray-900">
  <div class="mx-auto max-w-screen-lg px-4 2xl:px-0">
    <h1 class="text-2xl font-semibold text-gray-900 dark:text-white mb-6">Your Shopping Cart</h1>
    {#if $cart.length > 0}
      <div class="overflow-x-auto">
        <table class="w-full text-sm text-left text-gray-500 dark:text-gray-400">
          <thead class="text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
            <tr>
              <th scope="col" class="px-6 py-3">Product</th>
              <th scope="col" class="px-6 py-3">Price</th>
              <th scope="col" class="px-6 py-3">Quantity</th>
              <th scope="col" class="px-6 py-3">Total</th>
              <th scope="col" class="px-6 py-3">Actions</th>
            </tr>
          </thead>
          <tbody>
            {#each $cart as item (item.id)}
              <tr class="bg-white border-b dark:bg-gray-800 dark:border-gray-700">
                <td class="px-6 py-4 flex items-center gap-4">
                  {#if item.image}
                    <img src={item.image} alt={item.title} class="w-16 h-16 object-cover rounded" />
                  {/if}
                  <a href={`/products/${item.id}/${item.slug}`} class="font-medium text-gray-900 dark:text-white hover:underline">{item.title}</a>
                </td>
                <td class="px-6 py-4">${item.price.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</td>
                <td class="px-6 py-4">
                  <input
                    type="number"
                    min="1"
                    value={item.quantity}
                    on:change={(e: Event) => cart.updateQuantity(item.id, +(e.target as HTMLInputElement).value)}
                    class="w-20 p-1 border border-gray-300 rounded text-center dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                  />
                </td>
                <td class="px-6 py-4">${(item.price * item.quantity).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</td>
                <td class="px-6 py-4">
                  <button on:click={() => cart.remove(item.id)} class="text-red-600 hover:text-red-800 dark:text-red-500 dark:hover:text-red-700">Remove</button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
      <div class="mt-6 flex justify-end items-center gap-4">
        <span class="text-xl font-semibold text-gray-900 dark:text-white">Subtotal:</span>
        <span class="text-2xl font-bold text-gray-900 dark:text-white">${$total.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
      </div>
      <div class="mt-6 flex justify-end">
        <button on:click={() => goto('/checkout')} class="bg-primary-600 hover:bg-primary-700 text-white font-medium rounded px-6 py-2 focus:outline-none focus:ring-2 focus:ring-primary-300 dark:focus:ring-primary-800">Proceed to Checkout</button>
      </div>
    {:else}
      <p class="text-center text-gray-500 dark:text-gray-400">Your cart is empty.</p>
      <div class="mt-4 text-center">
        <a href="/products" class="text-primary-600 hover:underline dark:text-primary-500">Continue Shopping</a>
      </div>
    {/if}
  </div>
</section> 