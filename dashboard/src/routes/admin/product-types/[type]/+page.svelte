<script lang="ts">
	import File from 'lucide-svelte/icons/file';
	import ListFilter from 'lucide-svelte/icons/list-filter';
	import Ellipsis from 'lucide-svelte/icons/ellipsis';
	import CirclePlus from 'lucide-svelte/icons/circle-plus';
	import * as Table from '$lib/components/ui/table';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Tabs from '$lib/components/ui/tabs';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { getContext } from 'svelte';
	import { api } from '$lib/api';
	import type { Product } from '$lib/types';
	const client = useQueryClient()
	const authFetch = getContext('authFetch')
	import type { PageData } from './$types';

	export let data: PageData;

	let limit = 10
	let after = 0

	const products = createQuery<
	Product[],
	Error
	>({
	queryKey: ['products', 'type', data.typeId, limit, after],
	queryFn: () => api(authFetch as any).getProductsByType(data.typeId, limit, after),
	})
</script>

<div>
	<main class="grid flex-1 items-start gap-4 p-4 sm:px-6 sm:py-0 md:gap-8">
		<Tabs.Root value="all">
			<div class="flex items-center">
				<Tabs.List>
					<Tabs.Trigger value="all">All</Tabs.Trigger>
					<Tabs.Trigger value="active">Active</Tabs.Trigger>
					<Tabs.Trigger value="draft">Draft</Tabs.Trigger>
					<Tabs.Trigger value="archived" class="hidden sm:flex">Archived</Tabs.Trigger>
				</Tabs.List>
				<div class="ml-auto flex items-center gap-2">
					<DropdownMenu.Root>
						<DropdownMenu.Trigger asChild let:builder>
							<Button builders={[builder]} variant="outline" size="sm" class="h-7 gap-1">
								<ListFilter class="h-3.5 w-3.5" />
								<span class="sr-only sm:not-sr-only sm:whitespace-nowrap"> Filter </span>
							</Button>
						</DropdownMenu.Trigger>
						<DropdownMenu.Content align="end">
							<DropdownMenu.Label>Filter by</DropdownMenu.Label>
							<DropdownMenu.Separator />
							<DropdownMenu.CheckboxItem checked>Active</DropdownMenu.CheckboxItem>
							<DropdownMenu.CheckboxItem>Draft</DropdownMenu.CheckboxItem>
							<DropdownMenu.CheckboxItem>Archived</DropdownMenu.CheckboxItem>
						</DropdownMenu.Content>
					</DropdownMenu.Root>
					<Button size="sm" variant="outline" class="h-7 gap-1">
						<File class="h-3.5 w-3.5" />
						<span class="sr-only sm:not-sr-only sm:whitespace-nowrap"> Export </span>
					</Button>
					<Button size="sm" class="h-7 gap-1" href="/admin/product-types/{data?.typeId|| ''}/products/create">
						<CirclePlus class="h-3.5 w-3.5" />
						<span class="sr-only sm:not-sr-only sm:whitespace-nowrap"> Add Product </span>
					</Button>
				</div>
			</div>
			<Tabs.Content value="all">
				<Card.Root>
					<Card.Header>
						<Card.Title>Products</Card.Title>
						<Card.Description>
							Manage your products and view their sales performance.
						</Card.Description>
					</Card.Header>
					<Card.Content>
						<Table.Root>
							<Table.Header>
								<Table.Row>
									<Table.Head class="hidden w-[100px] sm:table-cell">
										<span class="sr-only">Image</span>
									</Table.Head>
									<Table.Head>Name</Table.Head>
									<Table.Head>Status</Table.Head>
									<Table.Head>Price</Table.Head>
									<Table.Head class="hidden md:table-cell">Total Sales</Table.Head>
									<Table.Head class="hidden md:table-cell">Created at</Table.Head>
									<Table.Head>
										<span class="sr-only">Actions</span>
									</Table.Head>
								</Table.Row>
							</Table.Header>
							<Table.Body>
								{#if $products.status === 'pending'}
								<p>Loading...</p>
								{:else if $products.status === 'error'}
								<span>Error: {$products.error.message}</span>
								{:else}
									{#each $products.data as product}
										<Table.Row>
											<Table.Cell class="hidden sm:table-cell">
												<img
													alt={product?.images?.[0]?.alt || "Product image"}
													class="aspect-square rounded-md object-cover"
													height="64"
													width="64"
													src={product?.images?.[0]?.url || "/images/placeholder.png"}
												/>
											</Table.Cell>
											<Table.Cell class="font-medium">{product?.title || 'Untitled'}</Table.Cell>
											<Table.Cell>
												<Badge variant="outline">{product?.status || 'Unknown'}</Badge>
											</Table.Cell>
											<Table.Cell>
												{product?.variants[0]?.price || 'N/A'}
											</Table.Cell>
											<Table.Cell class="hidden md:table-cell">25</Table.Cell>
											<Table.Cell class="hidden md:table-cell">
												{product?.variants?.[0]?.created_at ? new Date(product.variants[0].created_at).toLocaleDateString() : 'N/A'}
											</Table.Cell>
											<Table.Cell>
												<DropdownMenu.Root>
													<DropdownMenu.Trigger asChild let:builder>
														<Button
															aria-haspopup="true"
															size="icon"
															variant="ghost"
															builders={[builder]}
														>
															<Ellipsis class="h-4 w-4" />
															<span class="sr-only">Toggle menu</span>
														</Button>
													</DropdownMenu.Trigger>
													<DropdownMenu.Content align="end">
														<DropdownMenu.Label>Actions</DropdownMenu.Label>
														<DropdownMenu.Item href="/admin/product-types/{data?.typeId|| ''}/products/{product?.product_id || ''}/">
															Edit
														</DropdownMenu.Item>
														<DropdownMenu.Item>Delete</DropdownMenu.Item>
													</DropdownMenu.Content>
												</DropdownMenu.Root>
											</Table.Cell>
										</Table.Row>
									{/each}
								{/if}
							</Table.Body>
						</Table.Root>
					</Card.Content>
					<Card.Footer>
						<div class="text-xs text-muted-foreground">
							Showing <strong>1-10</strong> of <strong>32</strong> products
						</div>
					</Card.Footer>
				</Card.Root>
			</Tabs.Content>
		</Tabs.Root>
	</main>
</div>
