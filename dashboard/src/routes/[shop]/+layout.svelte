<script lang="ts">
	import '../../app.css'
	import {
		Home,
		LineChart,
		Package,
		Package2,
		PanelLeft,
		Search,
		Settings,
		ShoppingCart,
		Workflow,
		UsersRound,
		Boxes,
		Sun,
		Moon,
		Sparkles,
		Bell,
		User,
		Menu,
		X,
		Palette
	} from 'lucide-svelte';

	import * as Breadcrumb from '$lib/components/ui/breadcrumb';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Input } from '$lib/components/ui/input';
	import * as Sheet from '$lib/components/ui/sheet';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { page } from '$app/stores';
	import { SignOut } from '@auth/sveltekit/components';
	import type { PageData } from './$types' 
	import { setContext } from 'svelte';
	import { writable } from 'svelte/store';
	import type { Session as AuthSession } from '@auth/sveltekit';
	import { onMount } from 'svelte';
	import { fetchShopIdFromSubdomain, currentShopId } from '$lib/api';
	import { QueryClientProvider } from '@tanstack/svelte-query'
	import { Toaster } from '$lib/components/ui/sonner'
	import { toggleMode, mode } from 'mode-watcher';
	import { createAuthenticatedFetch } from '$lib/auth-fetch';
	import DeploymentIndicator from '$lib/components/deployment/DeploymentIndicator.svelte';
	import { deploymentStore } from '$lib/stores/deployment';

	interface CustomSession extends AuthSession {
		access_token?: string;
		user?: {
			name?: string;
			email?: string;
		};
	}
	export let data: PageData & { session: CustomSession | null };

	// Create a writable store for the session to make it reactive
	const sessionStore = writable<CustomSession | null>(data.session);

	// Update the store whenever data.session changes
	$: sessionStore.set(data.session)

	// Define an authenticated fetch function that uses the latest session
	const authFetch = async (input: RequestInfo | URL, options: RequestInit = {}) => {
		const session = $sessionStore
		const accessToken = session?.access_token
		if (accessToken) {
		options.headers = {
			...options.headers,
			Authorization: `Bearer ${accessToken}`,
		}
		}
		const baseFetch = (input: RequestInfo | URL, init?: RequestInit) => fetch(input, init);
		const authenticatedFetch = createAuthenticatedFetch(baseFetch);
		return authenticatedFetch(input, options);
	}

	// Share authFetch via context
	setContext('authFetch', authFetch)

	// Enhanced state management
	let searchFocused = false;
	let sidebarExpanded = false;

	// Function to get breadcrumb items based on current route
	$: breadcrumbItems = () => {
		const path = $page.url.pathname;
		const shopId = $page.params.shop;
		const items = [];

		// Always add Dashboard as first item
		items.push({ label: 'Dashboard', href: `/${shopId}`, type: 'link' });

		// Add other items based on path
		if (path.includes('/product-types')) {
			items.push({ label: 'Product Types', href: `/${shopId}/product-types`, type: 'link' });
			if (path.includes('/create')) {
				items.push({ label: 'Create', type: 'page' });
			} else if (path.match(/\/product-types\/\d+/)) {
				items.push({ label: 'Edit', type: 'page' });
			}
		} else if (path.includes('/products')) {
			items.push({ label: 'Products', href: `/${shopId}/products`, type: 'link' });
			if (path.includes('/create')) {
				items.push({ label: 'Create', type: 'page' });
			} else if (path.match(/\/products\/\d+/)) {
				items.push({ label: 'Edit', type: 'page' });
			}
		} else if (path.includes('/orders')) {
			items.push({ label: 'Orders', href: `/${shopId}/orders`, type: 'link' });
		} else if (path.includes('/customers')) {
			items.push({ label: 'Customers', href: `/${shopId}/customers`, type: 'link' });
			if (path.includes('/create')) {
				items.push({ label: 'Create', type: 'page' });
			} else if (path.match(/\/customers\/\d+/)) {
				items.push({ label: 'View', type: 'page' });
			}
		} else if (path.includes('/inventory')) {
			items.push({ label: 'Inventory', href: `/${shopId}/inventory`, type: 'link' });
		} else if (path.includes('/categories')) {
			items.push({ label: 'Categories', href: `/${shopId}/categories`, type: 'link' });
			if (path.includes('/create')) {
				items.push({ label: 'Create', type: 'page' });
			}
		} else if (path.includes('/templates')) {
			items.push({ label: 'Templates', href: `/${shopId}/templates`, type: 'link' });
		} else if (path.includes('/settings')) {
			items.push({ label: 'Settings', href: `/${shopId}/settings`, type: 'link' });
			if (path.includes('/domain')) {
				items.push({ label: 'Domain', type: 'page' });
			} else if (path.includes('/socials')) {
				items.push({ label: 'Social', type: 'page' });
			}
		}

		return items;
	};

	// Add this to initialize shop ID on page load
	onMount(() => {
		const unsubscribe = page.subscribe(($page) => {
			if ($page.params.shop) {
				fetchShopIdFromSubdomain($page.params.shop, authFetch);
			}
		});
		
		// Set up deployment store with auth fetch
		deploymentStore.setAuthFetch(authFetch);
		
		return unsubscribe;
	})

	// reactive current path and active helper
	$: currentPath = $page.url.pathname;
	function isActive(path: string): boolean {
		return currentPath === path || currentPath.startsWith(path + '/');
	}

	// Navigation items for sidebar
	const navItems = [
		{ href: '', icon: Home, label: 'Dashboard' },
		{ href: 'orders', icon: ShoppingCart, label: 'Orders' },
		{ href: 'product-types', icon: Boxes, label: 'Product Types' },
		{ href: 'customers', icon: UsersRound, label: 'Customers' },
		{ href: 'inventory', icon: Package, label: 'Inventory' },
		// { href: 'categories', icon: Workflow, label: 'Categories' },
		// { href: 'templates', icon: Palette, label: 'Templates' }
	];
</script>

<QueryClientProvider client={data.queryClient}>
<div class="flex min-h-screen w-full bg-gradient-to-br from-background via-surface-elevated to-surface-muted relative overflow-hidden">
	
	<!-- Ambient background effects -->
	<div class="fixed inset-0 overflow-hidden pointer-events-none">
		<div class="absolute -top-40 -right-40 w-80 h-80 bg-gradient-to-br from-primary/10 to-accent/10 rounded-full blur-3xl animate-pulse"></div>
		<div class="absolute -bottom-40 -left-40 w-80 h-80 bg-gradient-to-tr from-accent/10 to-secondary/10 rounded-full blur-3xl animate-pulse delay-1000"></div>
	</div>

	<!-- Grid pattern overlay -->
	<div class="fixed inset-0 opacity-[0.02] dark:opacity-[0.05] pointer-events-none" style="background-image: radial-gradient(circle at 1px 1px, currentColor 1px, transparent 0); background-size: 40px 40px;"></div>
	
	<!-- Enhanced Sidebar -->
	<aside class="fixed inset-y-0 left-0 z-20 hidden w-16 hover:w-64 flex-col border-r border-border/50 glass shadow-glass sm:flex transition-all duration-300 group">
		<!-- Decorative gradient overlay -->
		<div class="absolute inset-0 bg-gradient-to-b from-primary/5 via-transparent to-accent/5 pointer-events-none rounded-r-3xl"></div>
		
		<nav class="relative flex flex-col gap-2 px-3 py-6">
			<!-- Logo with enhanced styling -->
			<a
				href="/account"
				class="group/logo relative flex h-12 w-12 shrink-0 items-center justify-center gap-3 rounded-2xl bg-surface-elevated hover:bg-surface-elevated/80 border border-border/50 shadow-glass transition-all duration-300 hover:shadow-xl hover:scale-105 active:scale-95 mb-4"
			>
				<img src="/images/logo.svg" alt="Naytife" class="h-6 w-6 transition-all group-hover/logo:scale-110" />
				
				<!-- Shine effect -->
				<div class="absolute inset-0 rounded-2xl bg-gradient-to-tr from-white/20 to-transparent opacity-0 group-hover/logo:opacity-100 transition-opacity duration-300"></div>
			</a>

			<!-- Navigation Items -->
			{#each navItems as item}
				<Tooltip.Root>
					<Tooltip.Trigger asChild let:builder>
						<a
							href="/{$page.params.shop}/{item.href}"
							class="nav-item group/nav overflow-hidden"
							class:nav-item-active={isActive(`/${$page.params.shop}/${item.href}`)}
							class:nav-item-inactive={!isActive(`/${$page.params.shop}/${item.href}`)}
							use:builder.action
							{...builder}
						>
							<div class="flex items-center gap-3 min-w-0">
								<svelte:component this={item.icon} class="h-5 w-5 flex-shrink-0 transition-all group-hover/nav:scale-110" />
								<span class="opacity-0 group-hover:opacity-100 transition-opacity duration-300 whitespace-nowrap text-sm font-medium truncate">
									{item.label}
								</span>
							</div>
							
							<!-- Active indicator -->
							{#if isActive(`/${$page.params.shop}/${item.href}`)}
								<div class="absolute -right-1 top-1/2 h-6 w-1 -translate-y-1/2 rounded-full bg-primary-foreground/30"></div>
							{/if}

							<!-- Hover effect -->
							<div class="absolute inset-0 bg-gradient-to-r from-transparent via-primary/5 to-transparent opacity-0 group-hover/nav:opacity-100 transition-opacity duration-300 rounded-xl"></div>
						</a>
					</Tooltip.Trigger>
					<Tooltip.Content side="right" class="glass border-border/50 shadow-glass ml-2 group-hover:hidden">
						{item.label}
					</Tooltip.Content>
				</Tooltip.Root>
			{/each}
		</nav>

		<!-- Settings at bottom -->
		<nav class="relative mt-auto flex flex-col gap-2 px-3 py-6">
			<Tooltip.Root>
				<Tooltip.Trigger asChild let:builder>
					<a
						href="/{$page.params.shop}/settings/store/general"
						class="nav-item group/nav overflow-hidden"
						class:nav-item-active={isActive(`/${$page.params.shop}/settings`)}
						class:nav-item-inactive={!isActive(`/${$page.params.shop}/settings`)}
						use:builder.action
						{...builder}
					>
						<div class="flex items-center gap-3 min-w-0">
							<Settings class="h-5 w-5 flex-shrink-0 transition-all group-hover/nav:scale-110 group-hover/nav:rotate-90" />
							<span class="opacity-0 group-hover:opacity-100 transition-opacity duration-300 whitespace-nowrap text-sm font-medium">
								Settings
							</span>
						</div>
						
						{#if isActive(`/${$page.params.shop}/settings`)}
							<div class="absolute -right-1 top-1/2 h-6 w-1 -translate-y-1/2 rounded-full bg-primary-foreground/30"></div>
						{/if}

						<!-- Hover effect -->
						<div class="absolute inset-0 bg-gradient-to-r from-transparent via-primary/5 to-transparent opacity-0 group-hover/nav:opacity-100 transition-opacity duration-300 rounded-xl"></div>
					</a>
				</Tooltip.Trigger>
				<Tooltip.Content side="right" class="glass border-border/50 shadow-glass ml-2 group-hover:hidden">
					Settings
				</Tooltip.Content>
			</Tooltip.Root>
		</nav>
	</aside>

	<!-- Main Content Area -->
	<div class="flex flex-col sm:gap-6 sm:py-6 sm:pl-16 w-full">
		<!-- Enhanced Header -->
		<header class="sticky top-0 z-10 flex h-16 items-center gap-4 border-b border-border/50 glass px-6 shadow-glass sm:static sm:h-auto sm:border-0 sm:glass-0 sm:shadow-none">
			
			<!-- Mobile Menu -->
			<Sheet.Root>
				<Sheet.Trigger asChild let:builder>
					<Button builders={[builder]} size="icon" variant="outline" class="sm:hidden glass border-border/50 shadow-glass backdrop-blur-strong">
						<Menu class="h-5 w-5" />
						<span class="sr-only">Toggle Menu</span>
					</Button>
				</Sheet.Trigger>
				<Sheet.Content side="left" class="sm:max-w-xs glass border-border/50 backdrop-blur-strong">
					<nav class="grid gap-4 text-lg font-medium pt-8">
						<a
							href="/account"
							class="group flex h-12 w-12 shrink-0 items-center justify-center gap-2 rounded-2xl bg-surface-elevated hover:bg-surface-elevated/80 border border-border/50 shadow-glass mb-4"
						>
							<img src="/images/logo.svg" alt="Naytife" class="h-6 w-6 transition-all group-hover:scale-110" />
							<span class="sr-only">Naytife</span>
						</a>
						
						<!-- Mobile nav items with enhanced styling -->
						{#each [{ href: '', icon: Home, label: 'Dashboard' }, ...navItems.slice(1)] as item}
							<a
								href="/{$page.params.shop}/{item.href}"
								class="nav-item"
								class:nav-item-active={isActive(`/${$page.params.shop}/${item.href}`)}
								class:nav-item-inactive={!isActive(`/${$page.params.shop}/${item.href}`)}
							>
								<svelte:component this={item.icon} class="h-5 w-5" />
								{item.label}
							</a>
						{/each}

						<a
							href="/{$page.params.shop}/settings"
							class="nav-item"
							class:nav-item-active={isActive(`/${$page.params.shop}/settings`)}
							class:nav-item-inactive={!isActive(`/${$page.params.shop}/settings`)}
						>
							<Settings class="h-5 w-5" />
							Settings
						</a>
					</nav>
				</Sheet.Content>
			</Sheet.Root>

			<!-- Breadcrumb Navigation -->
			<Breadcrumb.Root class="hidden md:flex">
				<Breadcrumb.List>
					{#each breadcrumbItems() as item, index}
						<Breadcrumb.Item>
							{#if item.type === 'link'}
								<Breadcrumb.Link href={item.href} class="text-muted-foreground hover:text-foreground transition-colors">
									{item.label}
								</Breadcrumb.Link>
							{:else}
								<Breadcrumb.Page class="text-foreground font-medium">
									{item.label}
								</Breadcrumb.Page>
							{/if}
						</Breadcrumb.Item>
						{#if index < breadcrumbItems().length - 1}
							<Breadcrumb.Separator />
						{/if}
					{/each}
				</Breadcrumb.List>
			</Breadcrumb.Root>

			<!-- Action buttons with enhanced styling -->
			<div class="flex items-center gap-2 ml-auto">
				<!-- Theme toggle -->
				<Tooltip.Root>
					<Tooltip.Trigger asChild let:builder>
						<Button 
							builders={[builder]} 
							variant="outline" 
							size="icon" 
							on:click={toggleMode}
							class="glass border-border/50 hover:bg-surface-elevated hover:scale-105 transition-all duration-300 shadow-glass rounded-2xl"
						>
							<div class="relative overflow-hidden">
								{#if $mode === 'dark'}
									<Sun class="h-5 w-5 transition-all duration-300 rotate-0 scale-100" />
								{:else}
									<Moon class="h-5 w-5 transition-all duration-300 rotate-0 scale-100" />
								{/if}
							</div>
							<span class="sr-only">Toggle Theme</span>
						</Button>
					</Tooltip.Trigger>
					<Tooltip.Content side="bottom" class="glass border-border/50 shadow-glass">
						Toggle Theme
					</Tooltip.Content>
				</Tooltip.Root>

				<!-- User menu -->
				<DropdownMenu.Root>
					<DropdownMenu.Trigger asChild let:builder>
						<Button
							builders={[builder]}
							variant="outline"
							size="icon"
							class="overflow-hidden rounded-2xl glass border-border/50 hover:bg-surface-elevated hover:scale-105 transition-all duration-300 shadow-glass p-0"
						>
							<div class="relative">
								<div class="w-8 h-8 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center">
									<User class="w-4 h-4 text-primary-foreground" />
								</div>
							</div>
						</Button>
					</DropdownMenu.Trigger>
					<DropdownMenu.Content 
						align="end" 
						class="glass border-border/50 shadow-glass rounded-2xl min-w-[200px] p-2"
					>
						<DropdownMenu.Label class="text-sm font-semibold text-foreground px-3 py-2">
							Store Dashboard
						</DropdownMenu.Label>
						<DropdownMenu.Separator class="bg-border/50" />
						<DropdownMenu.Item 
							href="/account" 
							class="nav-item rounded-xl hover:bg-surface-elevated transition-all duration-200 px-3 py-2 cursor-pointer"
						>
							<div class="flex items-center gap-3">
								<Home class="w-4 h-4 text-muted-foreground" />
								All Stores
							</div>
						</DropdownMenu.Item>
						<DropdownMenu.Item class="nav-item rounded-xl hover:bg-surface-elevated transition-all duration-200 px-3 py-2 cursor-pointer">
							<div class="flex items-center gap-3">
								<Sparkles class="w-4 h-4 text-muted-foreground" />
								Support
							</div>
						</DropdownMenu.Item>
						<DropdownMenu.Separator class="bg-border/50" />
						<DropdownMenu.Item class="nav-item rounded-xl hover:bg-destructive/10 transition-all duration-200 px-3 py-2 cursor-pointer text-destructive">
							<SignOut></SignOut>
						</DropdownMenu.Item>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</div>
		</header>

		<!-- Content area with enhanced spacing -->
		<main class="flex-1 p-4 sm:p-8">
			<div class="max-w-7xl mx-auto">
				<slot></slot>
			</div>
		</main>
	</div>

	<!-- Enhanced Toaster -->
	<Toaster 
		theme="dark"
		position="top-right"
		closeButton={true}
		richColors={true}
		toastOptions={{
			style: 'background: var(--glass-bg); border: 1px solid var(--glass-border); backdrop-filter: blur(20px);'
		}}
	/>

	<!-- Deployment Indicator -->
	<DeploymentIndicator />
</div>
</QueryClientProvider>