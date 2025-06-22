<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { cart } from '$lib/stores/cart';
	import { currencyCode, currencySymbol } from '$lib/stores/currency';
	import { theme, toggleTheme } from '$lib/stores/theme';
	import { initializeData, shopData, isDataLoading } from '$lib/data/store';
	import { Sun, Moon, ShoppingBag, User, Menu } from 'lucide-svelte';
	
	// Initialize optimized data loading
	onMount(async () => {
		await initializeData();
	});
	
	// Set currency from shop data when loaded
	$: if ($shopData?.currencyCode) {
		currencyCode.set($shopData.currencyCode);
	}
</script>

<nav class="border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 antialiased">
	<div class="max-w-7xl mx-auto px-6 py-5">
		<div class="flex items-center justify-between">
			<div class="flex items-center space-x-12">
				<div class="shrink-0">
					<a href="/" class="inline-flex items-center">
						{#if $shopData}
							<span class="text-xl font-bold text-gray-900 dark:text-white">
								{$shopData.title}
							</span>
						{:else if $isDataLoading}
							<div class="animate-pulse bg-gray-200 dark:bg-gray-700 h-6 w-32 rounded"></div>
						{:else}
							<span class="text-xl font-bold text-gray-900 dark:text-white">Store</span>
						{/if}
					</a>
				</div>

				<ul class="hidden lg:flex items-center space-x-8">
					<li>
						<a href="/" class="text-sm uppercase tracking-wider font-medium text-gray-900 hover:text-primary-700 dark:text-white dark:hover:text-primary-500 border-b-2 border-transparent hover:border-primary-700 dark:hover:border-primary-500 py-1 transition-all duration-200">
							Home
						</a>
					</li>
				</ul>
			</div>

			<div class="flex items-center space-x-6">
				<button id="myCartDropdownButton1" data-dropdown-toggle="myCartDropdown1" class="group relative inline-flex items-center justify-center p-2 hover:bg-gray-50 dark:hover:bg-gray-800 text-sm font-medium text-gray-900 dark:text-white transition-colors duration-200">
					<span class="sr-only">Cart</span>
					<ShoppingBag class="w-5 h-5 lg:mr-2" />
					<span class="hidden sm:inline-flex uppercase tracking-wider text-xs font-semibold">My Cart</span>
					{#if $cart.length > 0}
						<span class="absolute -top-1 -right-1 w-4 h-4 flex items-center justify-center bg-primary-700 text-white text-[10px] font-bold">{$cart.length}</span>
					{/if}
				</button>

				<div id="myCartDropdown1" class="hidden z-10 w-80 p-0 overflow-hidden bg-white shadow-xl dark:bg-gray-900 border border-gray-200 dark:border-gray-700">
					<div class="p-4 border-b border-gray-100 dark:border-gray-800">
						<h3 class="text-sm font-semibold uppercase tracking-wider text-gray-900 dark:text-white">Your Cart</h3>
					</div>
					{#if $cart.length > 0}
						<div class="max-h-96 overflow-y-auto">
							{#each $cart as item (item.id)}
								<div class="grid grid-cols-[auto_1fr_auto] gap-4 p-4 border-b border-gray-100 dark:border-gray-800">
									<div class="w-16 h-16 bg-gray-100 dark:bg-gray-800"></div>
									<div>
										<a href={`/products/${item.id}/${item.slug}`} class="text-sm font-medium text-gray-900 dark:text-white hover:text-primary-700 dark:hover:text-primary-500">{item.title}</a>
										<p class="mt-1 text-sm font-normal text-gray-500 dark:text-gray-400">Qty: {item.quantity}</p>
										<p class="mt-1 text-sm font-medium text-gray-900 dark:text-white">{$currencySymbol}{item.price.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</p>
									</div>
									<button type="button" on:click={() => cart.remove(item.id)} class="self-start text-gray-400 hover:text-red-600 dark:text-gray-500 dark:hover:text-red-500 transition-colors duration-200">
										<span class="sr-only">Remove</span>
										<svg class="h-4 w-4" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
										</svg>
									</button>
								</div>
							{/each}
						</div>
						<div class="p-4 bg-gray-50 dark:bg-gray-800">
							<a href="/cart" class="block w-full px-5 py-3 bg-primary-700 hover:bg-primary-800 focus:ring-2 focus:ring-primary-700 focus:ring-offset-2 dark:focus:ring-offset-gray-900 text-center text-sm font-medium uppercase tracking-wider text-white transition-all duration-200">
								Proceed to Cart
							</a>
						</div>
					{:else}
						<div class="p-6 text-center">
							<p class="text-sm text-gray-500 dark:text-gray-400">Your cart is empty</p>
						</div>
					{/if}
				</div>

				<button id="userDropdownButton1" data-dropdown-toggle="userDropdown1" class="inline-flex items-center justify-center p-2 hover:bg-gray-50 dark:hover:bg-gray-800 text-sm font-medium text-gray-900 dark:text-white transition-colors duration-200">
					<User class="w-5 h-5 lg:mr-2" />
					<span class="hidden lg:inline-flex uppercase tracking-wider text-xs font-semibold">Account</span>
				</button>

				<div id="userDropdown1" class="hidden z-10 w-56 bg-white shadow-xl dark:bg-gray-900 border border-gray-200 dark:border-gray-700">
					<div class="p-4 border-b border-gray-100 dark:border-gray-800">
						<h3 class="text-sm font-semibold uppercase tracking-wider text-gray-900 dark:text-white">Account</h3>
					</div>
					<ul class="p-2">
						<li><a href="/account" class="flex w-full items-center px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors duration-200"> My Account </a></li>
						<li><a href="/orders" class="flex w-full items-center px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors duration-200"> My Orders </a></li>
						<li><a href="/settings" class="flex w-full items-center px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors duration-200"> Settings </a></li>
						<li><a href="/favourites" class="flex w-full items-center px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors duration-200"> Favourites </a></li>
						<li><a href="/addresses" class="flex w-full items-center px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors duration-200"> Delivery Addresses </a></li>
						<li><a href="/billing" class="flex w-full items-center px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors duration-200"> Billing Data </a></li>
					</ul>
					<div class="p-2 border-t border-gray-100 dark:border-gray-800">
						<a href="/signout" class="flex w-full items-center px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors duration-200"> Sign Out </a>
					</div>
				</div>

				<button type="button" data-collapse-toggle="ecommerce-navbar-menu-1" aria-controls="ecommerce-navbar-menu-1" aria-expanded="false" class="inline-flex lg:hidden items-center justify-center hover:bg-gray-50 dark:hover:bg-gray-800 p-2 text-gray-900 dark:text-white transition-colors duration-200">
					<span class="sr-only">Open Menu</span>
					<Menu class="w-5 h-5" />
				</button>

				<button 
					on:click={toggleTheme} 
					type="button" 
					class="inline-flex items-center justify-center p-2 hover:bg-gray-50 dark:hover:bg-gray-800 text-sm font-medium text-gray-900 dark:text-white transition-colors duration-200"
				>
					<span class="sr-only">Toggle theme</span>
					{#if $theme === 'dark'}
						<Sun class="w-5 h-5" />
					{:else}
						<Moon class="w-5 h-5" />
					{/if}
				</button>
			</div>
		</div>

		<div id="ecommerce-navbar-menu-1" class="border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 py-4 px-4 mt-4 hidden">
			<ul class="space-y-4">
				<li>
					<a href="/" class="block text-sm uppercase tracking-wider font-medium text-gray-900 hover:text-primary-700 dark:text-white dark:hover:text-primary-500 transition-colors duration-200">Home</a>
				</li>
			</ul>
		</div>
	</div>
</nav>

<slot />

<footer class="bg-white dark:bg-gray-900 border-t border-gray-200 dark:border-gray-700">
	<div class="max-w-7xl mx-auto px-6 py-12">
		<div class="grid grid-cols-1 md:grid-cols-4 gap-12">
			<div class="col-span-1 md:col-span-2">
				<a href="/" class="flex items-center space-x-4">
					{#if $shopData}
						<span class="text-xl font-serif font-medium text-gray-900 dark:text-white">
							{$shopData.title}
						</span>
					{:else}
						<div class="animate-pulse bg-gray-200 dark:bg-gray-700 h-6 w-32 rounded"></div>
					{/if}
				</a>
				<p class="mt-6 text-sm text-gray-600 dark:text-gray-400 max-w-sm">
					{#if $shopData?.about}
						{$shopData.about}
					{:else}
						Premium quality products curated for the discerning customer. Discover the difference that quality makes.
					{/if}
				</p>
				<div class="mt-8 flex space-x-6">
					<a href="{$shopData?.facebookLink || '#'}" class="text-gray-400 hover:text-gray-500 dark:hover:text-gray-300">
						<span class="sr-only">Facebook</span>
						<svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
							<path fill-rule="evenodd" d="M22 12c0-5.523-4.477-10-10-10S2 6.477 2 12c0 4.991 3.657 9.128 8.438 9.878v-6.987h-2.54V12h2.54V9.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V12h2.773l-.443 2.89h-2.33v6.988C18.343 21.128 22 16.991 22 12z" clip-rule="evenodd" />
						</svg>
					</a>
					<a href="{$shopData?.instagramLink || '#'}" class="text-gray-400 hover:text-gray-500 dark:hover:text-gray-300">
						<span class="sr-only">Instagram</span>
						<svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
							<path fill-rule="evenodd" d="M12.315 2c2.43 0 2.784.013 3.808.06 1.064.049 1.791.218 2.427.465a4.902 4.902 0 011.772 1.153 4.902 4.902 0 011.153 1.772c.247.636.416 1.363.465 2.427.048 1.067.06 1.407.06 4.123v.08c0 2.643-.012 2.987-.06 4.043-.049 1.064-.218 1.791-.465 2.427a4.902 4.902 0 01-1.153 1.772 4.902 4.902 0 01-1.772 1.153c-.636.247-1.363.416-2.427.465-1.067.048-1.407.06-4.123.06h-.08c-2.643 0-2.987-.012-4.043-.06-1.064-.049-1.791-.218-2.427-.465a4.902 4.902 0 01-1.772-1.153 4.902 4.902 0 01-1.153-1.772c-.247-.636-.416-1.363-.465-2.427-.047-1.024-.06-1.379-.06-3.808v-.63c0-2.43.013-2.784.06-3.808.049-1.064.218-1.791.465-2.427a4.902 4.902 0 011.153-1.772A4.902 4.902 0 015.45 2.525c.636-.247 1.363-.416 2.427-.465C8.901 2.013 9.256 2 11.685 2h.63zm-.081 1.802h-.468c-2.456 0-2.784.011-3.807.058-.975.045-1.504.207-1.857.344-.467.182-.8.398-1.15.748-.35.35-.566.683-.748 1.15-.137.353-.3.882-.344 1.857-.047 1.023-.058 1.351-.058 3.807v.468c0 2.456.011 2.784.058 3.807.045.975.207 1.504.344 1.857.182.466.399.8.748 1.15.35.35.683.566 1.15.748.353.137.882.3 1.857.344 1.054.048 1.37.058 4.041.058h.08c2.597 0 2.917-.01 3.96-.058.976-.045 1.505-.207 1.858-.344.466-.182.8-.398 1.15-.748.35-.35.566-.683.748-1.15.137-.353.3-.882.344-1.857.048-1.055.058-1.37.058-4.041v-.08c0-2.597-.01-2.917-.058-3.96-.045-.976-.207-1.505-.344-1.858a3.097 3.097 0 00-.748-1.15 3.098 3.098 0 00-1.15-.748c-.353-.137-.882-.3-1.857-.344-1.023-.047-1.351-.058-3.807-.058zM12 6.865a5.135 5.135 0 110 10.27 5.135 5.135 0 010-10.27zm0 1.802a3.333 3.333 0 100 6.666 3.333 3.333 0 000-6.666zm5.338-3.205a1.2 1.2 0 110 2.4 1.2 1.2 0 010-2.4z" clip-rule="evenodd" />
						</svg>
					</a>
					<a href="{$shopData?.whatsAppLink || '#'}" class="text-gray-400 hover:text-gray-500 dark:hover:text-gray-300">
						<span class="sr-only">WhatsApp</span>
						<svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
							<path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413Z" />
						</svg>
					</a>
				</div>
			</div>
			<div>
				<h3 class="text-sm uppercase tracking-wider font-semibold text-gray-900 dark:text-white">Shop</h3>
				<ul class="mt-4 space-y-4">
					<li>
						<a href="/" class="text-sm text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white">
							New Arrivals
						</a>
					</li>
					<li>
						<a href="/" class="text-sm text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white">
							Best Sellers
						</a>
					</li>
					<li>
						<a href="/" class="text-sm text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white">
							Collections
						</a>
					</li>
					<li>
						<a href="/" class="text-sm text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white">
							Gift Cards
						</a>
					</li>
				</ul>
			</div>
			<div>
				<h3 class="text-sm uppercase tracking-wider font-semibold text-gray-900 dark:text-white">Legal</h3>
				<ul class="mt-4 space-y-4">
					<li>
						<a href="/privacy" class="text-sm text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white">
							Privacy Policy
						</a>
					</li>
					<li>
						<a href="/terms" class="text-sm text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white">
							Terms & Conditions
						</a>
					</li>
					<li>
						<a href="/shipping" class="text-sm text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white">
							Shipping Policy
						</a>
					</li>
					<li>
						<a href="/returns" class="text-sm text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white">
							Returns & Exchanges
						</a>
					</li>
				</ul>
			</div>
		</div>
		<div class="mt-12 pt-8 border-t border-gray-200 dark:border-gray-700">
			<p class="text-sm text-gray-500 dark:text-gray-400">
				© {new Date().getFullYear()} 
				{#if $shopData}
					<a href="https://{$shopData.defaultDomain}" class="hover:text-gray-900 dark:hover:text-white">{$shopData.title}</a>
				{:else if $isDataLoading}
					<span class="animate-pulse bg-gray-200 dark:bg-gray-700 h-4 w-24 rounded inline-block"></span>
				{:else}
					<span>Store</span>
				{/if}
				. All Rights Reserved.
			</p>
		</div>
	</div>
</footer>
