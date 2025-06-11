<script lang="ts">
	import Package from 'lucide-svelte/icons/package';
	import Truck from 'lucide-svelte/icons/truck';
	import Download from 'lucide-svelte/icons/download';
	import Zap from 'lucide-svelte/icons/zap';
	import { Plus, Edit, Trash2, Eye } from 'lucide-svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { page } from '$app/stores';
	import type { PageData } from './$types';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { getContext } from 'svelte';
	import { api } from '$lib/api';
	import type { ProductType } from '$lib/types';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	const authFetch = getContext<typeof fetch>('authFetch');
	const queryClient = useQueryClient();

	$: shopParam = $page.params.shop;

	const productTypes = createQuery<ProductType[], Error>({
		queryKey: [`shop-${shopParam}-product-types`],
		queryFn: () => api(authFetch).getProductTypes(),
	});

	const createNewProductType = () => {
		goto(`/${$page.params.shop}/product-types/create`);
	};

	onMount(async () => {
		await queryClient.prefetchQuery({
			queryKey: [`shop-${shopParam}-product-types`],
			queryFn: () => api(authFetch).getProductTypes(),
		});
	});
</script>

<div class="space-y-8">
	<!-- Enhanced Header Section -->
	<div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
		<div>
			<h1 class="text-3xl font-bold text-foreground mb-2">Product Types</h1>
			<p class="text-muted-foreground">
				Manage your product categories and their configurations
			</p>
		</div>
		<div class="flex items-center gap-3">
			<Button on:click={createNewProductType} class="btn-gradient shadow-brand">
				<Plus class="w-4 h-4 mr-2" />
				Add Product Type
			</Button>
		</div>
	</div>

	<!-- Enhanced Stats Cards -->
	<div class="grid gap-6 md:grid-cols-3">
		<div class="card-interactive text-center">
			<div class="w-12 h-12 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-brand">
				<Package class="w-6 h-6 text-white" />
			</div>
			<div class="text-2xl font-bold text-foreground mb-1">{$productTypes.data?.length || 0}</div>
			<div class="text-sm text-muted-foreground">Total Types</div>
		</div>
		
		<div class="card-interactive text-center">
			<div class="w-12 h-12 bg-gradient-to-br from-success to-emerald-500 rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-brand">
				<Truck class="w-6 h-6 text-white" />
			</div>
			<div class="text-2xl font-bold text-foreground mb-1">{$productTypes.data?.filter(pt => pt.shippable).length || 0}</div>
			<div class="text-sm text-muted-foreground">Shippable</div>
		</div>
		
		<div class="card-interactive text-center">
			<div class="w-12 h-12 bg-gradient-to-br from-secondary to-accent rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-brand">
				<Download class="w-6 h-6 text-white" />
			</div>
			<div class="text-2xl font-bold text-foreground mb-1">{$productTypes.data?.filter(pt => pt.digital).length || 0}</div>
			<div class="text-sm text-muted-foreground">Digital</div>
		</div>
	</div>

	<!-- Enhanced Main Content -->
	<Card.Root class="border-0 shadow-2xl shadow-slate-900/10 bg-white/50 backdrop-blur-sm dark:bg-slate-900/50 dark:shadow-slate-900/20">
		<Card.Content class="p-6">
			{#if $productTypes.status === 'pending'}
				<div class="text-center py-12">
					<div class="flex flex-col items-center gap-3">
						<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
						<p class="text-muted-foreground">Loading product types...</p>
					</div>
				</div>
			{:else if $productTypes.status === 'error'}
				<div class="text-center py-12">
					<div class="flex flex-col items-center gap-3">
						<div class="p-3 bg-red-100 rounded-full dark:bg-red-900/20">
							<Zap class="h-6 w-6 text-red-600 dark:text-red-400" />
						</div>
						<p class="text-red-600 dark:text-red-400 font-medium">
							Error: {$productTypes.error.message}
						</p>
					</div>
				</div>
			{:else if ($productTypes.data?.length || 0) === 0}
				<div class="text-center py-12">
					<div class="flex flex-col items-center gap-3">
						<div class="p-3 bg-slate-100 rounded-full dark:bg-slate-800">
							<Package class="h-6 w-6 text-slate-400" />
						</div>
						<p class="text-slate-600 dark:text-slate-400">
							No product types found
						</p>
						<Button on:click={createNewProductType} class="mt-4">
							<Plus class="w-4 h-4 mr-2" />
							Create Your First Product Type
						</Button>
					</div>
				</div>
			{:else}
				<!-- Card Grid Layout -->
				<div class="grid gap-6 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
					{#each $productTypes.data as pt}
						<Card.Root class="group overflow-hidden border-0 shadow-lg hover:shadow-xl transition-all duration-300 bg-white/80 backdrop-blur-sm dark:bg-slate-900/80 hover:scale-[1.02] rounded-2xl">
							<!-- Card Header with Icon and Title -->
							<Card.Header class="pb-4">
								<div class="flex items-start justify-between">
									<div class="flex items-center gap-3">
										<div class="p-3 bg-gradient-to-br from-primary/10 to-accent/10 rounded-xl group-hover:from-primary/20 group-hover:to-accent/20 transition-all duration-300">
											<Package class="h-6 w-6 text-primary" />
										</div>
										<div>
											<h3 class="font-semibold text-foreground text-lg leading-tight">
												{pt.title}
											</h3>
											<p class="text-sm text-muted-foreground mt-1">
												SKU: <code class="px-1.5 py-0.5 bg-slate-100 dark:bg-slate-800 rounded text-xs font-mono">{pt.sku_substring}</code>
											</p>
										</div>
									</div>
								</div>
							</Card.Header>

							<!-- Card Content with Type Badges -->
							<Card.Content class="pt-0 pb-4">
								<div class="flex flex-wrap gap-2 mb-4">
									<Badge 
										variant={pt.digital ? 'secondary' : 'default'}
										class="font-medium {pt.digital 
											? 'bg-purple-100 text-purple-700 hover:bg-purple-200 dark:bg-purple-900/30 dark:text-purple-300' 
											: 'bg-blue-100 text-blue-700 hover:bg-blue-200 dark:bg-blue-900/30 dark:text-blue-300'
										} transition-colors duration-200"
									>
										{#if pt.digital}
											<Download class="h-3 w-3 mr-1" />
											Digital
										{:else}
											<Package class="h-3 w-3 mr-1" />
											Physical
										{/if}
									</Badge>
									
									<Badge 
										variant={pt.shippable ? 'secondary' : 'destructive'}
										class="font-medium {pt.shippable 
											? 'bg-green-100 text-green-700 hover:bg-green-200 dark:bg-green-900/30 dark:text-green-300' 
											: 'bg-red-100 text-red-700 hover:bg-red-200 dark:bg-red-900/30 dark:text-red-300'
										} transition-colors duration-200"
									>
										{#if pt.shippable}
											<Truck class="h-3 w-3 mr-1" />
											Shippable
										{:else}
											Not Shippable
										{/if}
									</Badge>
								</div>
							</Card.Content>

							<!-- Card Footer with Action Buttons -->
							<Card.Footer class="pt-0 pb-6">
								<div class="flex gap-2 w-full">
									<Button 
										href={`/${$page.params.shop}/product-types/${pt.id}`}
										variant="outline" 
										size="sm"
										class="flex-1 h-9 gap-2 bg-white/80 hover:bg-white border-white/20 shadow-sm hover:shadow-md transition-all duration-300 rounded-lg dark:bg-slate-800/80 dark:hover:bg-slate-700 dark:border-slate-700/50"
									>
										<Eye class="h-3.5 w-3.5" />
										View
									</Button>
									
									<Button 
										href={`/${$page.params.shop}/product-types/${pt.id}/edit`}
										variant="outline" 
										size="sm"
										class="flex-1 h-9 gap-2 bg-white/80 hover:bg-white border-white/20 shadow-sm hover:shadow-md transition-all duration-300 rounded-lg dark:bg-slate-800/80 dark:hover:bg-slate-700 dark:border-slate-700/50"
									>
										<Edit class="h-3.5 w-3.5" />
										Edit
									</Button>
									
									<Button 
										variant="outline" 
										size="sm"
										class="h-9 px-3 text-red-600 hover:text-red-700 hover:bg-red-50 border-red-200 hover:border-red-300 dark:text-red-400 dark:hover:text-red-300 dark:hover:bg-red-900/20 dark:border-red-800 dark:hover:border-red-700 transition-all duration-300 rounded-lg"
									>
										<Trash2 class="h-3.5 w-3.5" />
									</Button>
								</div>
							</Card.Footer>
						</Card.Root>
					{/each}
				</div>
			{/if}
		</Card.Content>
	</Card.Root>
</div>