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
	const authFetch = getContext<typeof fetch>('authFetch');

	const productTypes = createQuery<ProductType[], Error>({
		queryKey: ['product-types'],
		queryFn: () => api(authFetch).getProductTypes()
	});
</script>

<div>
	<main class="grid flex-1 items-start gap-4 p-4 sm:px-6 sm:py-0 md:gap-8">
		<Tabs.Root value="all">
			<div class="flex items-center">
				<!-- <Tabs.List>
					<Tabs.Trigger value="all">All</Tabs.Trigger>
					<Tabs.Trigger value="active">Active</Tabs.Trigger>
					<Tabs.Trigger value="draft">Draft</Tabs.Trigger>
					<Tabs.Trigger value="archived" class="hidden sm:flex">Archived</Tabs.Trigger>
				</Tabs.List> -->
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
					<Button size="sm" class="h-7 gap-1" href="/admin/product-types/create">
						<CirclePlus class="h-3.5 w-3.5" />
						<span class="sr-only sm:not-sr-only sm:whitespace-nowrap"> Add Product Type </span>
					</Button>
				</div>
			</div>
			<Tabs.Content value="all">
				<Card.Root>
					<Card.Header>
						<Card.Title>Product Types</Card.Title>
						<Card.Description>
							view your products types and edit their different variations.
						</Card.Description>
					</Card.Header>
					<Card.Content class="grid grid-cols-1 gap-6 p-4 sm:grid-cols-2 lg:grid-cols-4">
						{#if $productTypes.status === 'pending'}
							<p>Loading...</p>
						{:else if $productTypes.status === 'error'}
							<span>Error: {$productTypes.error.message}</span>
						{:else}
							{#each $productTypes.data as productType}
								<Card.Root class="relative rounded-2xl p-4 shadow-lg">
									<CardContent class="mt-4">
										<h2 class="text-xl font-semibold">{productType.title}</h2>
										<p class="mt-2 text-sm text-gray-700">
											{productType.digital ? 'Digital Product' : 'Physical Product'}
										</p>
										<p class="mt-2 text-sm text-gray-700">
											{productType.shippable ? 'Shippable' : 'Not Shippable'}
										</p>
										<div class="mt-4">
											<Button
												size="sm"
												class="w-full"
												href="/admin/product-types/{productType.id}"
											>
												<CirclePlus class="mr-2 h-4 w-4" />
												View Products
											</Button>
										</div>
									</CardContent>
									<Card.Content>
										<DropdownMenu.Root>
											<DropdownMenu.Trigger asChild let:builder>
												<Button
													class="absolute right-4 top-4"
													aria-haspopup="true"
													size="icon"
													variant="ghost"
													builders={[builder]}
												>
													<EllipsisVertical class="h-5 w-5" />
													<span class="sr-only">Toggle menu</span>
												</Button>
											</DropdownMenu.Trigger>
											<DropdownMenu.Content align="end">
												<DropdownMenu.Label>Actions</DropdownMenu.Label>
												<DropdownMenu.Item href="/admin/product-types/{productType?.id || ''}/edit">
													Edit
												</DropdownMenu.Item>
												<DropdownMenu.Item>Delete</DropdownMenu.Item>
											</DropdownMenu.Content>
										</DropdownMenu.Root>
									</Card.Content>
								</Card.Root>
							{/each}
						{/if}
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
