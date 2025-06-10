<script lang="ts">
	import File from 'lucide-svelte/icons/file';
	import ListFilter from 'lucide-svelte/icons/list-filter';
	import Ellipsis from 'lucide-svelte/icons/ellipsis';
	import CirclePlus from 'lucide-svelte/icons/circle-plus';
	import Search from 'lucide-svelte/icons/search';
	import Package from 'lucide-svelte/icons/package';
	import Truck from 'lucide-svelte/icons/truck';
	import Download from 'lucide-svelte/icons/download';
	import Zap from 'lucide-svelte/icons/zap';
	import { Plus, Edit, Trash2, Eye, MoreHorizontal, Filter } from 'lucide-svelte';
	import * as Table from '$lib/components/ui/table';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Tabs from '$lib/components/ui/tabs';
	import { page } from '$app/stores';
	import type { PageData } from './$types';
	import CardContent from '$lib/components/ui/card/card-content.svelte';
	import { Popover, PopoverTrigger } from '$lib/components/ui/popover';
	import PopoverContent from '$lib/components/ui/popover/popover-content.svelte';
	import Label from '$lib/components/ui/label/label.svelte';
	import Input from '$lib/components/ui/input/input.svelte';
	import { EllipsisVertical } from 'lucide-svelte';
	import { Content, Item, Select } from '$lib/components/ui/select';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { getContext } from 'svelte';
	import { api } from '$lib/api';
	import type { ProductType } from '$lib/types';
	import type { RequestEvent } from '@sveltejs/kit';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	const authFetch = getContext<typeof fetch>('authFetch');
	const queryClient = useQueryClient();

	$: shopParam = $page.params.shop;

	const productTypes = createQuery<ProductType[], Error>({
		queryKey: [`shop-${shopParam}-product-types`],
		queryFn: () => api(authFetch).getProductTypes(),
	});

	let searchTerm: string = '';
	$: filteredProductTypes = $productTypes.data?.filter(pt =>
		pt.title.toLowerCase().includes(searchTerm.toLowerCase())
	) || [];

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
			<Button variant="outline" class="glass border-border/50">
				<Download class="w-4 h-4 mr-2" />
				Export
			</Button>
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
			<div class="w-12 h-12 bg-gradient-to-br from-secondary to-purple-500 rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-brand">
				<Download class="w-6 h-6 text-white" />
			</div>
			<div class="text-2xl font-bold text-foreground mb-1">{$productTypes.data?.filter(pt => pt.digital).length || 0}</div>
			<div class="text-sm text-muted-foreground">Digital</div>
		</div>
	</div>

	<!-- Enhanced Main Content -->
	<Card.Root class="border-0 shadow-2xl shadow-slate-900/10 bg-white/50 backdrop-blur-sm dark:bg-slate-900/50 dark:shadow-slate-900/20">
		<Card.Header class="pb-8">
			<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-6">
				<!-- Enhanced Search -->
				<div class="relative flex-1 max-w-md">
					<Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
					<Input
						type="text"
						placeholder="Search product types..."
						class="pl-10 bg-white/80 border-white/20 shadow-sm focus:shadow-lg focus:shadow-blue-500/20 transition-all duration-300 rounded-xl dark:bg-slate-800/80 dark:border-slate-700/50"
						bind:value={searchTerm}
					/>
				</div>

				<!-- Enhanced Action Buttons -->
				<div class="flex items-center gap-3">
					<DropdownMenu.Root>
						<DropdownMenu.Trigger asChild let:builder>
							<Button 
								builders={[builder]} 
								variant="outline" 
								size="sm" 
								class="h-10 gap-2 bg-white/80 hover:bg-white border-white/20 shadow-sm hover:shadow-lg transition-all duration-300 rounded-xl dark:bg-slate-800/80 dark:hover:bg-slate-700 dark:border-slate-700/50"
							>
								<ListFilter class="h-4 w-4" />
								<span class="hidden sm:inline">Filter</span>
							</Button>
						</DropdownMenu.Trigger>
						<DropdownMenu.Content align="end" class="bg-white/95 backdrop-blur-xl border-white/20 shadow-2xl dark:bg-slate-900/95 dark:border-slate-700/50">
							<DropdownMenu.Label class="font-semibold">Filter by Status</DropdownMenu.Label>
							<DropdownMenu.Separator class="bg-slate-200/50 dark:bg-slate-700/50" />
							<DropdownMenu.CheckboxItem checked class="hover:bg-blue-50 focus:bg-blue-50 dark:hover:bg-slate-800 dark:focus:bg-slate-800">
								Active
							</DropdownMenu.CheckboxItem>
							<DropdownMenu.CheckboxItem class="hover:bg-blue-50 focus:bg-blue-50 dark:hover:bg-slate-800 dark:focus:bg-slate-800">
								Draft
							</DropdownMenu.CheckboxItem>
							<DropdownMenu.CheckboxItem class="hover:bg-blue-50 focus:bg-blue-50 dark:hover:bg-slate-800 dark:focus:bg-slate-800">
								Archived
							</DropdownMenu.CheckboxItem>
						</DropdownMenu.Content>
					</DropdownMenu.Root>

					<Button 
						size="sm" 
						variant="outline" 
						class="h-10 gap-2 bg-white/80 hover:bg-white border-white/20 shadow-sm hover:shadow-lg transition-all duration-300 rounded-xl dark:bg-slate-800/80 dark:hover:bg-slate-700 dark:border-slate-700/50"
					>
						<File class="h-4 w-4" />
						<span class="hidden sm:inline">Export</span>
					</Button>
				</div>
			</div>
		</Card.Header>

		<Card.Content class="p-0">
			<!-- Enhanced Table with modern styling -->
			<div class="overflow-hidden">
				<Table.Root class="w-full">
					<Table.Header>
						<Table.Row class="border-slate-200/50 dark:border-slate-700/50 hover:bg-transparent">
							<Table.Head class="font-semibold text-slate-700 dark:text-slate-300 py-4">
								Title
							</Table.Head>
							<Table.Head class="font-semibold text-slate-700 dark:text-slate-300">
								Type
							</Table.Head>
							<Table.Head class="font-semibold text-slate-700 dark:text-slate-300">
								Shippable
							</Table.Head>
							<Table.Head class="font-semibold text-slate-700 dark:text-slate-300">
								SKU Prefix
							</Table.Head>
							<Table.Head class="text-right font-semibold text-slate-700 dark:text-slate-300">
								Actions
							</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#if $productTypes.status === 'pending'}
							<Table.Row class="hover:bg-transparent">
								<Table.Cell colspan={5} class="text-center py-12">
									<div class="flex flex-col items-center gap-3">
										<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
										<p class="text-slate-600 dark:text-slate-400">Loading product types...</p>
									</div>
								</Table.Cell>
							</Table.Row>
						{:else if $productTypes.status === 'error'}
							<Table.Row class="hover:bg-transparent">
								<Table.Cell colspan={5} class="text-center py-12">
									<div class="flex flex-col items-center gap-3">
										<div class="p-3 bg-red-100 rounded-full dark:bg-red-900/20">
											<Zap class="h-6 w-6 text-red-600 dark:text-red-400" />
										</div>
										<p class="text-red-600 dark:text-red-400 font-medium">
											Error: {$productTypes.error.message}
										</p>
									</div>
								</Table.Cell>
							</Table.Row>
						{:else if filteredProductTypes.length === 0}
							<Table.Row class="hover:bg-transparent">
								<Table.Cell colspan={5} class="text-center py-12">
									<div class="flex flex-col items-center gap-3">
										<div class="p-3 bg-slate-100 rounded-full dark:bg-slate-800">
											<Package class="h-6 w-6 text-slate-400" />
										</div>
										<p class="text-slate-600 dark:text-slate-400">
											{searchTerm ? 'No product types match your search' : 'No product types found'}
										</p>
									</div>
								</Table.Cell>
							</Table.Row>
						{:else}
							{#each filteredProductTypes as pt, index}
								<Table.Row class="border-slate-200/30 dark:border-slate-700/30 hover:bg-slate-50/50 dark:hover:bg-slate-800/30 transition-all duration-200 group">
									<Table.Cell class="py-4">
										<div class="flex items-center gap-3">
											<div class="p-2 bg-gradient-to-br from-blue-500/10 to-indigo-500/10 rounded-lg group-hover:from-blue-500/20 group-hover:to-indigo-500/20 transition-all duration-200">
												<Package class="h-4 w-4 text-blue-600 dark:text-blue-400" />
											</div>
											<div>
												<p class="font-medium text-slate-900 dark:text-slate-100">
													{pt.title}
												</p>
											</div>
										</div>
									</Table.Cell>
									<Table.Cell>
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
									</Table.Cell>
									<Table.Cell>
										<Badge 
											variant={pt.shippable ? 'secondary' : 'destructive'}
											class="font-medium {pt.shippable 
												? 'bg-green-100 text-green-700 hover:bg-green-200 dark:bg-green-900/30 dark:text-green-300' 
												: 'bg-red-100 text-red-700 hover:bg-red-200 dark:bg-red-900/30 dark:text-red-300'
											} transition-colors duration-200"
										>
											{#if pt.shippable}
												<Truck class="h-3 w-3 mr-1" />
												Yes
											{:else}
												No
											{/if}
										</Badge>
									</Table.Cell>
									<Table.Cell>
										<code class="px-2 py-1 bg-slate-100 dark:bg-slate-800 rounded-md text-sm font-mono text-slate-700 dark:text-slate-300">
											{pt.sku_substring}
										</code>
									</Table.Cell>
									<Table.Cell class="text-right">
										<DropdownMenu.Root>
											<DropdownMenu.Trigger asChild let:builder>
												<Button 
													size="icon" 
													variant="ghost" 
													builders={[builder]}
													class="h-8 w-8 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 transition-all duration-200 opacity-0 group-hover:opacity-100"
												>
													<EllipsisVertical class="h-4 w-4" />
													<span class="sr-only">Actions</span>
												</Button>
											</DropdownMenu.Trigger>
											<DropdownMenu.Content 
												align="end" 
												class="bg-white/95 backdrop-blur-xl border-white/20 shadow-2xl dark:bg-slate-900/95 dark:border-slate-700/50"
											>
												<DropdownMenu.Item 
													href={`/${$page.params.shop}/product-types/${pt.id}`}
													class="hover:bg-blue-50 focus:bg-blue-50 dark:hover:bg-slate-800 dark:focus:bg-slate-800 transition-colors duration-200"
												>
													View Details
												</DropdownMenu.Item>
												<DropdownMenu.Item 
													href={`/${$page.params.shop}/product-types/${pt.id}/edit`}
													class="hover:bg-blue-50 focus:bg-blue-50 dark:hover:bg-slate-800 dark:focus:bg-slate-800 transition-colors duration-200"
												>
													Edit
												</DropdownMenu.Item>
												<DropdownMenu.Separator class="bg-slate-200/50 dark:bg-slate-700/50" />
												<DropdownMenu.Item class="hover:bg-red-50 focus:bg-red-50 text-red-600 dark:hover:bg-red-900/20 dark:focus:bg-red-900/20 transition-colors duration-200">
													Delete
												</DropdownMenu.Item>
											</DropdownMenu.Content>
										</DropdownMenu.Root>
									</Table.Cell>
								</Table.Row>
							{/each}
						{/if}
					</Table.Body>
				</Table.Root>
			</div>
		</Card.Content>

		<!-- Enhanced Footer -->
		<Card.Footer class="bg-slate-50/50 dark:bg-slate-800/30 border-t border-slate-200/30 dark:border-slate-700/30 rounded-b-2xl">
			<div class="flex items-center justify-between w-full">
				<div class="text-sm text-slate-600 dark:text-slate-400">
					Showing <span class="font-semibold text-slate-900 dark:text-slate-100">{filteredProductTypes.length}</span> 
					of <span class="font-semibold text-slate-900 dark:text-slate-100">{$productTypes.data?.length || 0}</span> product types
				</div>
				
				{#if filteredProductTypes.length > 0}
					<div class="flex items-center gap-2">
						<Button variant="outline" size="sm" class="h-8 bg-white/80 hover:bg-white border-white/20 shadow-sm rounded-lg dark:bg-slate-800/80 dark:hover:bg-slate-700 dark:border-slate-700/50">
							Previous
						</Button>
						<Button variant="outline" size="sm" class="h-8 bg-white/80 hover:bg-white border-white/20 shadow-sm rounded-lg dark:bg-slate-800/80 dark:hover:bg-slate-700 dark:border-slate-700/50">
							Next
						</Button>
					</div>
				{/if}
			</div>
		</Card.Footer>
	</Card.Root>
</div>