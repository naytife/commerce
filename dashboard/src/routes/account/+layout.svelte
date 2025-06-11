<script lang="ts">
    import {
		Home,
		LineChart,
		Package,
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
		User
	} from 'lucide-svelte';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Input } from '$lib/components/ui/input';
	import * as Sheet from '$lib/components/ui/sheet';
    import { Button } from '$lib/components/ui/button';
    import * as Tooltip from '$lib/components/ui/tooltip';
    import { SignOut } from '@auth/sveltekit/components';
	import { setContext } from 'svelte';
	import { writable } from 'svelte/store';
	import { toggleMode, mode } from 'mode-watcher';
	export let data: { session: any };
    
	const sessionStore = writable(data.session);
	$: sessionStore.set(data.session);
	const authFetch = async (url: string, options: RequestInit = {}) => {
		const session = $sessionStore;
		const accessToken = session?.access_token;
		if (accessToken) {
			options.headers = {
				...options.headers,
				Authorization: `Bearer ${accessToken}`,
			};
		}
		return fetch(url, options);
	};
	setContext('authFetch', authFetch);

	// Enhanced state for animations
	let searchFocused = false;
	let headerScrolled = false;

	// Handle scroll for dynamic header
	function handleScroll() {
		headerScrolled = window.scrollY > 10;
	}
</script>

<svelte:window on:scroll={handleScroll} />

<div class="flex min-h-screen w-full flex-col bg-gradient-to-br from-background via-surface-elevated to-surface-muted relative overflow-hidden">
	<!-- Ambient background effects -->
	<div class="fixed inset-0 overflow-hidden pointer-events-none">
		<div class="absolute -top-40 -right-40 w-80 h-80 bg-gradient-to-br from-primary/10 to-accent/10 rounded-full blur-3xl animate-pulse"></div>
		<div class="absolute -bottom-40 -left-40 w-80 h-80 bg-gradient-to-tr from-accent/10 to-secondary/10 rounded-full blur-3xl animate-pulse delay-1000"></div>
	</div>

	<!-- Grid pattern overlay -->
	<div class="fixed inset-0 opacity-[0.02] dark:opacity-[0.05] pointer-events-none" style="background-image: radial-gradient(circle at 1px 1px, currentColor 1px, transparent 0); background-size: 40px 40px;"></div>

	<div class="flex flex-col sm:gap-6 sm:py-6 sm:pl-16 relative z-10">
		<header
			class="sticky top-0 z-40 flex h-16 items-center gap-4 border-b border-border/50 px-4 sm:static sm:h-auto sm:border-0 sm:px-8 glass transition-all duration-300 {headerScrolled ? 'shadow-glass' : ''}"
		>
			<!-- Mobile menu with enhanced styling -->
			<Sheet.Root>
				<Sheet.Trigger asChild let:builder>
					<Button 
						builders={[builder]} 
						size="icon" 
						variant="outline" 
						class="sm:hidden glass border-border/50 hover:bg-surface-elevated hover:scale-105 transition-all duration-200 shadow-glass"
					>
						<PanelLeft class="h-5 w-5" />
						<span class="sr-only">Toggle Menu</span>
					</Button>
				</Sheet.Trigger>
				<Sheet.Content side="left" class="sm:max-w-xs glass border-r border-border/50">
					<nav class="grid gap-6 text-lg font-medium pt-8">
						<a
							href="##"
							class="group flex h-12 shrink-0 items-center justify-center gap-2 rounded-2xl bg-surface-elevated hover:bg-surface-elevated/80 border border-border/50 hover:scale-105 transition-all duration-300 shadow-glass p-2"
						>
							<img 
								src="/images/logo.svg" 
								alt="Naytife Logo" 
								class="h-8 w-8 transition-all group-hover:scale-110"
							/>
							<span class="sr-only">Naytife</span>
						</a>
					</nav>
				</Sheet.Content>
			</Sheet.Root>

			<!-- Enhanced logo for desktop -->
			<div class="hidden sm:flex items-center gap-3">
				<div class="group flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-surface-elevated hover:bg-surface-elevated/80 border border-border/50 hover:scale-105 transition-all duration-300 shadow-glass p-2 relative">
					<img 
						src="/images/logo.svg" 
						alt="Naytife Logo" 
						class="h-8 w-8 transition-all group-hover:scale-110 group-hover:rotate-12"
					/>
					<div class="absolute -top-1 -right-1 w-4 h-4 bg-gradient-to-br from-accent to-secondary rounded-full animate-pulse"></div>
				</div>
				<div class="flex flex-col">
					<span class="text-xl font-bold bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent">
						Naytife
					</span>
					<span class="text-xs text-muted-foreground">Account Dashboard</span>
				</div>
			</div>
	
			<!-- Enhanced search bar -->
			<div class="relative ml-auto flex-1 md:grow-0 max-w-md">
				<div class="relative">
					<Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground transition-colors {searchFocused ? 'text-primary' : ''}" />
					<Input
						type="search"
						placeholder="Search stores, analytics..."
						class="form-input pl-10 pr-4 w-full md:w-[300px] lg:w-[400px] transition-all duration-300 {searchFocused ? 'ring-2 ring-primary/20 border-primary' : ''}"
						on:focus={() => searchFocused = true}
						on:blur={() => searchFocused = false}
					/>
				</div>
			</div>

			<!-- Action buttons with enhanced styling -->
			<div class="flex items-center gap-2">
				<!-- Notifications -->
				<Tooltip.Root>
					<Tooltip.Trigger asChild let:builder>
						<Button 
							builders={[builder]} 
							variant="outline" 
							size="icon" 
							class="glass border-border/50 hover:bg-surface-elevated hover:scale-105 transition-all duration-300 shadow-glass rounded-2xl relative"
						>
							<Bell class="h-5 w-5" />
							<div class="absolute -top-1 -right-1 w-3 h-3 bg-gradient-to-br from-accent to-secondary rounded-full animate-pulse"></div>
							<span class="sr-only">Notifications</span>
						</Button>
					</Tooltip.Trigger>
					<Tooltip.Content side="bottom" class="glass border-border/50 shadow-glass">
						Notifications
					</Tooltip.Content>
				</Tooltip.Root>

				<!-- Theme toggle with enhanced animation -->
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

				<!-- User menu with enhanced styling -->
				<DropdownMenu.Root>
					<DropdownMenu.Trigger asChild let:builder>
						<Button
							builders={[builder]}
							variant="outline"
							size="icon"
							class="overflow-hidden rounded-2xl glass border-border/50 hover:bg-surface-elevated hover:scale-105 transition-all duration-300 shadow-glass p-0"
						>
							<div class="relative">
								<div class="w-8 h-8 bg-surface-elevated hover:bg-surface-elevated/80 border border-border/50 rounded-2xl flex items-center justify-center">
									<User class="w-4 h-4 text-foreground" />
								</div>
								<div class="absolute -top-1 -right-1 w-4 h-4 bg-gradient-to-br from-success to-emerald-500 rounded-full border-2 border-background animate-pulse"></div>
							</div>
						</Button>
					</DropdownMenu.Trigger>
					<DropdownMenu.Content 
						align="end" 
						class="glass border-border/50 shadow-glass rounded-2xl min-w-[200px] p-2"
					>
						<DropdownMenu.Label class="text-sm font-semibold text-foreground px-3 py-2">
							My Account
						</DropdownMenu.Label>
						<DropdownMenu.Separator class="bg-border/50" />
						<DropdownMenu.Item 
							href="/account" 
							class="nav-item rounded-xl hover:bg-surface-elevated transition-all duration-200 px-3 py-2 cursor-pointer"
						>
							<div class="flex items-center gap-3">
								<div class="w-2 h-2 bg-primary rounded-full"></div>
								Account
							</div>
						</DropdownMenu.Item>
						<DropdownMenu.Item class="nav-item rounded-xl hover:bg-surface-elevated transition-all duration-200 px-3 py-2 cursor-pointer">
							<div class="flex items-center gap-3">
								<Settings class="w-4 h-4 text-muted-foreground" />
								Settings
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

		<!-- Content area with enhanced spacing and backdrop -->
		<main class="flex-1 p-4 sm:p-8">
			<div class="max-w-7xl mx-auto">
				<slot></slot>
			</div>
		</main>
	</div>
</div>