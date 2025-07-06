<script lang="ts">
	import { page } from '$app/stores';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { goto } from '$app/navigation';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { getContext } from 'svelte';
	import { getAllShops, deleteShop } from '$lib/api';
	import { toast } from 'svelte-sonner';
	import type { Shop } from '$lib/types';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { Badge } from '$lib/components/ui/badge';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { 
		Plus, 
		Store, 
		TrendingUp, 
		Users, 
		ShoppingCart, 
		MoreVertical, 
		Edit, 
		Trash2,
		ExternalLink,
		BarChart3,
		Calendar,
		DollarSign,
		Activity
	} from 'lucide-svelte';

	const authFetch = getContext<typeof fetch>('authFetch');
	const queryClient = useQueryClient();

	const shops = createQuery<Shop[], Error>({
		queryKey: ['shops'],
		queryFn: () => getAllShops(authFetch)
	});

	$: if ($shops.status === 'error' && $shops.error) {
		toast.error($shops.error.message || 'Failed to load shops');
	}

	function goToStore(storeId: string) {
		goto(`/${storeId}`);
	}

	function createNewStore() {
		goto('/account/create-store');
	}

	function handleKeydown(event: any) {
		if (event.key === 'Enter') {
			createNewStore();
		}
	}

	let showDeleteConfirmation = false;
	let deleteConfirmText = '';
	let expectedDeleteConfirmText = '';
	let storeToDelete: Shop | null = null;
	$: isDeleteConfirmValid = deleteConfirmText === expectedDeleteConfirmText;

	function handleOpenDeleteDialog(store: Shop) {
		storeToDelete = store;
		showDeleteConfirmation = true;
		deleteConfirmText = '';
		expectedDeleteConfirmText = `DELETE ${store.title}`;
	}

	function handleDeleteDialogChange(open: boolean) {
		showDeleteConfirmation = open;
		if (!open) {
			deleteConfirmText = '';
			storeToDelete = null;
		}
	}

	async function handleDeleteStoreConfirmed() {
		if (!isDeleteConfirmValid || !storeToDelete) return;
		try {
			await deleteShop(storeToDelete.shop_id, authFetch);
			await queryClient.invalidateQueries({ queryKey: ['shops'] });
			toast.success('Store deleted');
		} catch (error) {
			console.error(error);
			toast.error('Failed to delete store');
		} finally {
			showDeleteConfirmation = false;
			storeToDelete = null;
			deleteConfirmText = '';
		}
	}
</script>

<div class="min-h-screen relative">
	<!-- Enhanced Header Section -->
	<div class="mb-12 relative">
		<!-- Welcome Section -->
		<div class="glass rounded-3xl p-8 border border-border/50 mb-8">
			<div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
				<div>
					<h1 class="text-3xl md:text-4xl font-bold mb-3">
						Welcome back!
					</h1>
					<p class="text-lg text-muted-foreground max-w-2xl leading-relaxed">
						Manage all your e-commerce stores from one powerful dashboard. 
						<span class="text-primary font-semibold">Scale your business</span> with advanced analytics and automation.
					</p>
				</div>
				<Button 
					on:click={createNewStore} 
					size="lg" 
					class="btn-gradient px-8 py-4 rounded-2xl shadow-brand group whitespace-nowrap"
				>
					<Plus class="w-5 h-5 mr-2 group-hover:rotate-90 transition-transform" />
					Create New Store
				</Button>
			</div>
		</div>
	</div>

	{#if $shops.status === 'pending'}
		<!-- Enhanced Loading State -->
		<div class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
			{#each Array(6) as _}
				<div class="card-elevated animate-pulse">
					<div class="space-y-6">
						<div class="space-y-3">
							<div class="h-6 bg-gradient-to-r from-muted to-muted/70 rounded-lg"></div>
							<div class="h-4 bg-gradient-to-r from-muted/70 to-muted/50 rounded w-3/4"></div>
						</div>
						<div class="space-y-2">
							<div class="h-3 bg-gradient-to-r from-muted/50 to-muted/30 rounded w-1/2"></div>
							<div class="h-3 bg-gradient-to-r from-muted/50 to-muted/30 rounded w-2/3"></div>
						</div>
						<div class="flex gap-2">
							<div class="h-8 bg-gradient-to-r from-muted/40 to-muted/20 rounded-xl flex-1"></div>
							<div class="h-8 w-10 bg-gradient-to-r from-muted/40 to-muted/20 rounded-xl"></div>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{:else if $shops.status === 'error'}
		<!-- Enhanced Error State -->
		<div class="card-elevated text-center py-16">
			<div class="w-16 h-16 bg-gradient-to-br from-destructive to-red-500 rounded-full flex items-center justify-center mx-auto mb-6 shadow-lg">
				<Activity class="w-8 h-8 text-white" />
			</div>
			<h3 class="text-xl font-semibold mb-2 text-foreground">Unable to load stores</h3>
			<p class="text-muted-foreground mb-6 max-w-md mx-auto">
				We encountered an issue loading your stores. Please check your connection and try again.
			</p>
			<Button 
				on:click={() => queryClient.invalidateQueries({ queryKey: ['shops'] })} 
				variant="outline" 
				class="glass border-border/50"
			>
				Try Again
			</Button>
		</div>
	{:else if $shops.data && $shops.data.length === 0}
		<!-- Enhanced Empty State -->
		<div class="card-elevated text-center py-20">
			<div class="w-20 h-20 bg-gradient-to-br from-primary to-accent rounded-full flex items-center justify-center mx-auto mb-8 shadow-brand">
				<Store class="w-10 h-10 text-white" />
			</div>
			<h3 class="text-2xl font-bold mb-4 text-foreground">Create your first store</h3>
			<p class="text-lg text-muted-foreground mb-8 max-w-2xl mx-auto leading-relaxed">
				Start your e-commerce journey today. Build a beautiful online store, 
				manage inventory, and grow your business with our powerful platform.
			</p>
			<div class="flex flex-col sm:flex-row gap-4 justify-center">
				<Button 
					on:click={createNewStore} 
					size="lg" 
					class="btn-gradient px-8 py-4 rounded-2xl shadow-brand group"
				>
					<Plus class="w-5 h-5 mr-2 group-hover:rotate-90 transition-transform" />
					Create Your First Store
				</Button>
				<Button 
					variant="outline" 
					size="lg" 
					class="glass border-border/50 px-8 py-4 rounded-2xl"
				>
					Watch Demo
				</Button>
			</div>
		</div>
	{:else if $shops.data}
		<!-- Enhanced Stores Grid -->
		<div class="space-y-6">
			<div class="flex items-center justify-between">
				<h2 class="text-2xl font-bold text-foreground">Your Stores</h2>
				<div class="text-sm text-muted-foreground">
					{$shops.data.length} {$shops.data.length === 1 ? 'store' : 'stores'}
				</div>
			</div>
			
			<div class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
				{#each $shops.data as store (store.shop_id)}
					<div class="card-interactive group">
						<!-- Store Header -->
						<div class="flex items-start justify-between mb-6">
							<div class="flex items-center gap-3">
								<div class="w-12 h-12 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center shadow-brand group-hover:scale-110 transition-transform">
									<Store class="w-6 h-6 text-white" />
								</div>
								<div>
									<h3 class="font-semibold text-lg text-foreground group-hover:text-primary transition-colors">
										{store.title}
									</h3>
									<Badge variant="secondary" class="text-xs mt-1">
										{store.status || 'Active'}
									</Badge>
								</div>
							</div>
							
							<!-- Actions Menu -->
							<DropdownMenu.Root>
								<DropdownMenu.Trigger asChild let:builder>
									<Button 
										builders={[builder]} 
										variant="ghost" 
										size="icon" 
										class="opacity-0 group-hover:opacity-100 transition-opacity rounded-xl hover:bg-surface-elevated"
									>
										<MoreVertical class="w-4 h-4" />
									</Button>
								</DropdownMenu.Trigger>
								<DropdownMenu.Content class="glass border-border/50 shadow-glass rounded-xl">
									<DropdownMenu.Item 
										on:click={() => goToStore(store.subdomain)} 
										class="flex items-center gap-2 rounded-lg"
									>
										<ExternalLink class="w-4 h-4" />
										Open Dashboard
									</DropdownMenu.Item>
									<DropdownMenu.Item class="flex items-center gap-2 rounded-lg">
										<Edit class="w-4 h-4" />
										Edit Store
									</DropdownMenu.Item>
									<DropdownMenu.Separator />
									<DropdownMenu.Item 
										on:click={() => handleOpenDeleteDialog(store)} 
										class="flex items-center gap-2 text-destructive rounded-lg"
									>
										<Trash2 class="w-4 h-4" />
										Delete Store
									</DropdownMenu.Item>
								</DropdownMenu.Content>
							</DropdownMenu.Root>
						</div>

						<!-- Store Stats -->
						<div class="grid grid-cols-2 gap-4 mb-6">
							<div class="text-center p-3 bg-surface-elevated rounded-xl">
								<div class="text-sm text-muted-foreground mb-1">Revenue</div>
								<div class="font-semibold text-foreground">₦2,458</div>
							</div>
							<div class="text-center p-3 bg-surface-elevated rounded-xl">
								<div class="text-sm text-muted-foreground mb-1">Orders</div>
								<div class="font-semibold text-foreground">127</div>
							</div>
						</div>

						<!-- Store Description -->
						<p class="text-sm text-muted-foreground mb-6 line-clamp-2">
							{store.about || 'A modern e-commerce store built with Naytife Commerce platform.'}
						</p>

						<!-- Action Buttons -->
						<div class="flex gap-2">
							<Button 
								on:click={() => goToStore(store.subdomain)} 
								class="flex-1 btn-gradient rounded-xl group"
							>
								<BarChart3 class="w-4 h-4 mr-2 group-hover:scale-110 transition-transform" />
								Dashboard
							</Button>
							<Button 
								variant="outline" 
								size="icon" 
								class="glass border-border/50 rounded-xl hover:scale-105 transition-transform"
								on:click={() => window.open(`https://${store.subdomain}.naytife.com`, '_blank')}
							>
								<ExternalLink class="w-4 h-4" />
							</Button>
						</div>

						<!-- Last Updated -->
						<div class="flex items-center gap-2 mt-4 pt-4 border-t border-border/50 text-xs text-muted-foreground">
							<Calendar class="w-3 h-3" />
							Updated {store.updated_at ? new Date(store.updated_at).toLocaleDateString() : 'recently'}
						</div>
					</div>
				{/each}
			</div>
		</div>
	{/if}

	<!-- Delete Confirmation Dialog -->
	<Dialog.Root bind:open={showDeleteConfirmation} onOpenChange={handleDeleteDialogChange}>
		<Dialog.Content class="glass border-border/50 shadow-glass rounded-2xl max-w-md">
			<Dialog.Header>
				<Dialog.Title class="text-xl font-semibold text-foreground">Delete Store</Dialog.Title>
				<Dialog.Description class="text-muted-foreground">
					This action cannot be undone. This will permanently delete your store and remove all associated data.
				</Dialog.Description>
			</Dialog.Header>
			<div class="space-y-4">
				<div class="p-4 bg-destructive/10 border border-destructive/20 rounded-xl">
					<p class="text-sm text-destructive font-medium mb-2">
						To confirm, type <span class="font-mono bg-destructive/20 px-2 py-1 rounded">{expectedDeleteConfirmText}</span>
					</p>
				</div>
				<div class="space-y-2">
					<Label for="delete-confirm" class="text-sm font-medium">Confirmation</Label>
					<Input
						id="delete-confirm"
						bind:value={deleteConfirmText}
						placeholder={expectedDeleteConfirmText}
						class="form-input"
					/>
				</div>
			</div>
			<Dialog.Footer class="flex flex-col-reverse sm:flex-row gap-2">
				<Button variant="outline" on:click={() => showDeleteConfirmation = false} class="glass border-border/50">
					Cancel
				</Button>
				<Button 
					on:click={handleDeleteStoreConfirmed} 
					disabled={!isDeleteConfirmValid}
					variant="destructive"
					class="disabled:opacity-50"
				>
					<Trash2 class="w-4 h-4 mr-2" />
					Delete Store
				</Button>
			</Dialog.Footer>
		</Dialog.Content>
	</Dialog.Root>
</div>

<style>
</style>