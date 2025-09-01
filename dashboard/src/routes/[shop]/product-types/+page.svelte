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
		queryFn: () => api(authFetch).getProductTypes()
	});

	const createNewProductType = () => {
		goto(`/${$page.params.shop}/product-types/create`);
	};

	onMount(async () => {
		await queryClient.prefetchQuery({
			queryKey: [`shop-${shopParam}-product-types`],
			queryFn: () => api(authFetch).getProductTypes()
		});
	});
</script>

<div class="space-y-8">
	<!-- Enhanced Header Section -->
	<div class="flex flex-col gap-6 lg:flex-row lg:items-center lg:justify-between">
		<div>
			<h1 class="text-foreground mb-2 text-3xl font-bold">Product Types</h1>
			<p class="text-muted-foreground">Manage your product categories and their configurations</p>
		</div>
		<div class="flex items-center gap-3">
			<Button on:click={createNewProductType} class="btn-gradient shadow-brand">
				<Plus class="mr-2 h-4 w-4" />
				Add Product Type
			</Button>
		</div>
	</div>

	<!-- Enhanced Stats Cards -->
	<!-- <div class="grid gap-6 md:grid-cols-3">
		<div class="card-interactive text-center">
			<div
				class="from-primary to-accent shadow-brand mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br"
			>
				<Package class="h-6 w-6 text-white" />
			</div>
			<div class="text-foreground mb-1 text-2xl font-bold">{$productTypes.data?.length || 0}</div>
			<div class="text-muted-foreground text-sm">Total Types</div>
		</div>

		<div class="card-interactive text-center">
			<div
				class="from-success shadow-brand mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br to-emerald-500"
			>
				<Truck class="h-6 w-6 text-white" />
			</div>
			<div class="text-foreground mb-1 text-2xl font-bold">
				{$productTypes.data?.filter((pt) => pt.shippable).length || 0}
			</div>
			<div class="text-muted-foreground text-sm">Shippable</div>
		</div>

		<div class="card-interactive text-center">
			<div
				class="from-secondary to-accent shadow-brand mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br"
			>
				<Download class="h-6 w-6 text-white" />
			</div>
			<div class="text-foreground mb-1 text-2xl font-bold">
				{$productTypes.data?.filter((pt) => pt.digital).length || 0}
			</div>
			<div class="text-muted-foreground text-sm">Digital</div>
		</div>
	</div> -->

	<!-- Enhanced Main Content -->
	<Card.Root
		class="border-0 bg-white/50 shadow-2xl shadow-slate-900/10 backdrop-blur-sm dark:bg-slate-900/50 dark:shadow-slate-900/20"
	>
		<Card.Content class="p-6">
			{#if $productTypes.status === 'pending'}
				<div class="py-12 text-center">
					<div class="flex flex-col items-center gap-3">
						<div class="border-primary h-8 w-8 animate-spin rounded-full border-b-2"></div>
						<p class="text-muted-foreground">Loading product types...</p>
					</div>
				</div>
			{:else if $productTypes.status === 'error'}
				<div class="py-12 text-center">
					<div class="flex flex-col items-center gap-3">
						<div class="rounded-full bg-red-100 p-3 dark:bg-red-900/20">
							<Zap class="h-6 w-6 text-red-600 dark:text-red-400" />
						</div>
						<p class="font-medium text-red-600 dark:text-red-400">
							Error: {$productTypes.error.message}
						</p>
					</div>
				</div>
			{:else if ($productTypes.data?.length || 0) === 0}
				<div class="py-12 text-center">
					<div class="flex flex-col items-center gap-3">
						<div class="rounded-full bg-slate-100 p-3 dark:bg-slate-800">
							<Package class="h-6 w-6 text-slate-400" />
						</div>
						<p class="text-slate-600 dark:text-slate-400">No product types found</p>
						<Button on:click={createNewProductType} class="mt-4">
							<Plus class="mr-2 h-4 w-4" />
							Create Your First Product Type
						</Button>
					</div>
				</div>
			{:else}
				<!-- Card Grid Layout -->
				<div class="grid gap-6 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
					{#each $productTypes.data as pt}
						<Card.Root
							class="group overflow-hidden rounded-2xl border-0 bg-white/80 shadow-lg backdrop-blur-sm transition-all duration-300 hover:scale-[1.02] hover:shadow-xl dark:bg-slate-900/80"
						>
							<!-- Card Header with Icon and Title -->
							<Card.Header class="pb-4">
								<div class="flex items-start justify-between">
									<div class="flex items-center gap-3">
										<div
											class="from-primary/10 to-accent/10 group-hover:from-primary/20 group-hover:to-accent/20 rounded-xl bg-gradient-to-br p-3 transition-all duration-300"
										>
											<Package class="text-primary h-6 w-6" />
										</div>
										<div>
											<h3 class="text-foreground text-lg font-semibold leading-tight">
												{pt.title}
											</h3>
											<p class="text-muted-foreground mt-1 text-sm">
												SKU: <code
													class="rounded bg-slate-100 px-1.5 py-0.5 font-mono text-xs dark:bg-slate-800"
													>{pt.sku_substring}</code
												>
											</p>
										</div>
									</div>
								</div>
							</Card.Header>

							<!-- Card Content with Type Badges -->
							<Card.Content class="pb-4 pt-0">
								<div class="mb-4 flex flex-wrap gap-2">
									<Badge
										variant={pt.digital ? 'secondary' : 'default'}
										class="font-medium {pt.digital
											? 'bg-purple-100 text-purple-700 hover:bg-purple-200 dark:bg-purple-900/30 dark:text-purple-300'
											: 'bg-blue-100 text-blue-700 hover:bg-blue-200 dark:bg-blue-900/30 dark:text-blue-300'} transition-colors duration-200"
									>
										{#if pt.digital}
											<Download class="mr-1 h-3 w-3" />
											Digital
										{:else}
											<Package class="mr-1 h-3 w-3" />
											Physical
										{/if}
									</Badge>

									<Badge
										variant={pt.shippable ? 'secondary' : 'destructive'}
										class="font-medium {pt.shippable
											? 'bg-green-100 text-green-700 hover:bg-green-200 dark:bg-green-900/30 dark:text-green-300'
											: 'bg-red-100 text-red-700 hover:bg-red-200 dark:bg-red-900/30 dark:text-red-300'} transition-colors duration-200"
									>
										{#if pt.shippable}
											<Truck class="mr-1 h-3 w-3" />
											Shippable
										{:else}
											Not Shippable
										{/if}
									</Badge>
								</div>
							</Card.Content>

							<!-- Card Footer with Action Buttons -->
							<Card.Footer class="pb-6 pt-0">
								<div class="flex w-full gap-2">
									<Button
										href={`/${$page.params.shop}/product-types/${pt.id}/products/create`}
										variant="outline"
										size="sm"
										class="h-9 flex-1 gap-2 rounded-lg border-white/20 bg-white/80 shadow-sm transition-all duration-300 hover:bg-white hover:shadow-md dark:border-slate-700/50 dark:bg-slate-800/80 dark:hover:bg-slate-700"
									>
										<Plus class="h-3.5 w-3.5" />
										Add Product
									</Button>

									<Button
										href={`/${$page.params.shop}/product-types/${pt.id}/edit`}
										variant="outline"
										size="sm"
										class="h-9 flex-1 gap-2 rounded-lg border-white/20 bg-white/80 shadow-sm transition-all duration-300 hover:bg-white hover:shadow-md dark:border-slate-700/50 dark:bg-slate-800/80 dark:hover:bg-slate-700"
									>
										<Edit class="h-3.5 w-3.5" />
										Edit
									</Button>

									<Button
										variant="outline"
										size="sm"
										class="h-9 rounded-lg border-red-200 px-3 text-red-600 transition-all duration-300 hover:border-red-300 hover:bg-red-50 hover:text-red-700 dark:border-red-800 dark:text-red-400 dark:hover:border-red-700 dark:hover:bg-red-900/20 dark:hover:text-red-300"
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
