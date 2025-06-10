<script lang="ts">
	import ChevronLeft from 'lucide-svelte/icons/chevron-left';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Table from '$lib/components/ui/table';
	import { CirclePlus, Trash2, Pencil, AlertTriangle } from 'lucide-svelte';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import * as Select from '$lib/components/ui/select';
	import * as Dialog from '$lib/components/ui/dialog';
	import { api } from '$lib/api';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import AttributeForm from '$lib/components/product/attribute-form.svelte';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import type { Attribute, ProductType } from '$lib/types';
	import { goto } from '$app/navigation';
	import type { PageData } from './$types';
	import { toast } from 'svelte-sonner';
	import { getContext } from 'svelte';

	export let data: PageData;
	const authFetch = getContext<typeof fetch>('authFetch');
	const queryClient = useQueryClient();

	let editingAttribute: Attribute | null = null;
	let showDeleteConfirmation = false;
	let deleteConfirmText = '';
	let expectedDeleteConfirmText = '';

	const productTypeQuery = createQuery({
		queryKey: [`shop-${$page.params.shop}-product-type`, $page.params.type],
		queryFn: () => api(authFetch).getProductTypeById(parseInt($page.params.type))
	});

	// Local state for editing
	let title = '';
	let skuSubstring = '';
	let shippable = false;
	let digital = false;

	// Update local state when product type data is loaded
	$: if ($productTypeQuery.data) {
		title = $productTypeQuery.data.title;
		skuSubstring = $productTypeQuery.data.sku_substring || '';
		shippable = $productTypeQuery.data.shippable;
		digital = $productTypeQuery.data.digital;
	}

	$: expectedDeleteConfirmText = $productTypeQuery.data ? `DELETE ${$productTypeQuery.data.title}` : '';
	$: isDeleteConfirmValid = deleteConfirmText === expectedDeleteConfirmText;

	const attributesQuery = createQuery<Attribute[]>({
		queryKey: ['product-type-attributes', $page.params.type],
		queryFn: () => api(authFetch).getProductTypeAttributes(parseInt($page.params.type))
	});

	let showProductAttributeForm = false;
	let showVariationAttributeForm = false;

	$: productAttributes = $attributesQuery.data?.filter((attr: Attribute) => attr.applies_to === "Product") ?? [];
	$: variationAttributes = $attributesQuery.data?.filter((attr: Attribute) => attr.applies_to === "ProductVariation") ?? [];

	// Generate SKU substring from title
	function generateSkuSubstring(input: string): string {
		if (!input) return '';
		
		const words = input.trim().split(/\s+/);
		let result = '';
		
		if (words.length === 1) {
			// Single word: use first and last character
			const word = words[0];
			if (word.length >= 2) {
				result = word[0] + word[word.length - 1];
			} else if (word.length === 1) {
				result = word[0] + word[0];
			}
		} else {
			// Multiple words: use first character of each word
			result = words.map(word => word[0] || '').join('');
		}
		
		return result.toUpperCase();
	}

	// Generate SKU when the SKU field is focused and empty
	function handleSkuFocus() {
		if (!skuSubstring && title) {
			skuSubstring = generateSkuSubstring(title);
		}
	}

	// Validate SKU substring
	function validateSkuSubstring(value: string): boolean {
		return value.length >= 2 && value.length <= 4 && value === value.toUpperCase();
	}

	let skuError = '';
	$: skuError = skuSubstring && !validateSkuSubstring(skuSubstring) 
		? 'SKU substring must be 2-4 uppercase characters' 
		: '';

	// Force uppercase for SKU input
	function handleSkuInput(event: Event) {
		const input = event.target as HTMLInputElement;
		skuSubstring = input.value.toUpperCase();
	}

	async function handleAttributeSaved() {
		await queryClient.invalidateQueries({ queryKey: ['product-type-attributes', $page.params.type] });
	}

	function handleProductDialogChange(open: boolean) {
		showProductAttributeForm = open;
	}

	function handleVariationDialogChange(open: boolean) {
		showVariationAttributeForm = open;
	}

	function handleDeleteDialogChange(open: boolean) {
		showDeleteConfirmation = open;
		if (!open) {
			deleteConfirmText = '';
		}
	}

	async function handleAttributeDelete(attributeId: number) {
		try {
			await api().deleteProductTypeAttribute(attributeId);
			await queryClient.invalidateQueries({ queryKey: ['product-type-attributes', $page.params.type] });
		} catch (error) {
			console.error('Failed to delete attribute:', error);
		}
	}

	async function handleUpdateProductType() {
		if (skuSubstring && !validateSkuSubstring(skuSubstring)) {
			toast.error('Please fix the SKU substring before saving');
			return;
		}

		try {
			await api().updateProductType(parseInt($page.params.type), {
				title,
				sku_substring: skuSubstring,
				shippable,
				digital
			});
			toast.success('Product type updated successfully');
			await queryClient.invalidateQueries({ queryKey: ['product-type', $page.params.type] });
		} catch (error) {
			console.error('Failed to update product type:', error);
			toast.error('Failed to update product type');
		}
	}

	async function handleDeleteProductType() {
		if (!isDeleteConfirmValid) return;
		
		try {
			await api().deleteProductType(parseInt($page.params.type));
			toast.success('Product type deleted successfully');
			showDeleteConfirmation = false;
			goto(`/${$page.params.shop}/product-types`);
		} catch (error) {
			console.error('Failed to delete product type:', error);
			toast.error('Failed to delete product type');
		}
	}

	function handleEditAttribute(attribute: Attribute) {
		editingAttribute = attribute;
		if (attribute.applies_to === 'Product') {
			showProductAttributeForm = true;
		} else {
			showVariationAttributeForm = true;
		}
	}

	onMount(async () => {
		await Promise.all([
			queryClient.prefetchQuery({
				queryKey: [`shop-${$page.params.shop}-product-type`, $page.params.type],
				queryFn: () => api(authFetch).getProductTypeById(parseInt($page.params.type)),
			}),
			queryClient.prefetchQuery({
				queryKey: [`shop-${$page.params.shop}-product-type-attributes`, $page.params.type],
				queryFn: () => api(authFetch).getProductTypeAttributes(parseInt($page.params.type)),
			}),
		]);
	});
</script>

<main class="grid flex-1 items-start gap-4 p-4 sm:px-6 sm:py-0 md:gap-8">
	<div class="mx-auto grid max-w-7xl flex-1 auto-rows-max gap-4">
		<div class="flex items-center gap-4">
			<Button variant="outline" size="icon" class="h-7 w-7">
				<ChevronLeft class="h-4 w-4" />
				<span class="sr-only">Back</span>
			</Button>

			<div class="hidden items-center gap-2 md:ml-auto md:flex">
				<Button variant="outline" size="sm" on:click={() => showDeleteConfirmation = true}>Delete</Button>
				<Button size="sm" on:click={handleUpdateProductType}>Update Type</Button>
			</div>
		</div>
		<div class="grid gap-4 md:grid-cols-1 lg:gap-8">
			<div class="grid auto-rows-max items-start gap-4">
				<Card.Root>
					<Card.Header>
						<Card.Title>Product Type</Card.Title>
						<Card.Description>
							Manage the details and configuration for this product type. Set its name, attributes, and variations.
						</Card.Description>
					</Card.Header>
					<Card.Content>
						<div class="grid gap-6">
							<div class="grid gap-3">
								<Label for="name">Name</Label>
								<Input id="name" type="text" class="w-full" bind:value={title} />
							</div>
							<div class="grid gap-3">
								<Label for="skuSubstring">SKU Substring</Label>
								<div class="flex flex-col gap-1">
									<div class="flex gap-2">
										<Input 
											id="skuSubstring" 
											type="text" 
											maxlength={4}
											class="w-full {skuError ? 'border-red-500' : ''}" 
											bind:value={skuSubstring}
											on:input={handleSkuInput}
											on:focus={handleSkuFocus}
											placeholder="Will auto-generate when focused"
										/>
										<Button 
											variant="outline" 
											size="sm" 
											on:click={() => skuSubstring = generateSkuSubstring(title)}
											disabled={!title}
											title="Generate from title"
										>
											Generate
										</Button>
									</div>
									{#if skuError}
										<p class="text-sm text-red-500">{skuError}</p>
									{:else}
										<p class="text-sm text-gray-400">Used for product SKU generation (2-4 uppercase characters).</p>
									{/if}
								</div>
							</div>
							<div class="flex items-center space-x-2">
								<Label for="physical">Physical product</Label>
								<Switch id="physical" bind:checked={shippable} />
							</div>
							<div class="flex items-center space-x-2">
								<Label for="digital">Digital product</Label>
								<Switch id="digital" bind:checked={digital} />
							</div>
						</div>
					</Card.Content>
				</Card.Root>

				<Card.Root>
					<Card.Header>
						<Card.Title>Product Attributes</Card.Title>
						<Card.Description>Attributes that apply to the product as a whole</Card.Description>
					</Card.Header>
					<Card.Content class="max-h-[400px] overflow-y-auto">
						<Table.Root class="w-full">
							<Table.Header>
								<Table.Row>
									<Table.Head>Title</Table.Head>
									<Table.Head>Type</Table.Head>
									<Table.Head>Required</Table.Head>
									<Table.Head>Options</Table.Head>
									<Table.Head class="w-[50px]">Actions</Table.Head>
								</Table.Row>
							</Table.Header>
							<Table.Body>
								{#each productAttributes as attr}
									<Table.Row>
										<Table.Cell>{attr.title}</Table.Cell>
										<Table.Cell>{attr.data_type}</Table.Cell>
										<Table.Cell>
											{#if attr.required}
												<Badge variant="default">Required</Badge>
											{:else}
												<Badge variant="outline">Optional</Badge>
											{/if}
										</Table.Cell>
										<Table.Cell>
											{#if attr.options}
												<Select.Root>
													<Select.Trigger>
														<Select.Value placeholder="Select option" />
													</Select.Trigger>
													<Select.Content>
														{#each attr.options as option}
															<Select.Item value={option.attribute_option_id.toString()}>
																{#if attr.data_type === 'Color'}
																	<div class="flex items-center gap-2">
																		<div class="h-4 w-4 rounded-full" style="background-color: {option.value}"></div>
																		{option.value}
																	</div>
																{:else}
																	{option.value}
																{/if}
															</Select.Item>
														{/each}
													</Select.Content>
												</Select.Root>
											{/if}
										</Table.Cell>
										<Table.Cell>
											<div class="flex gap-2">
												<Button 
													variant="ghost" 
													size="icon" 
													class="h-8 w-8" 
													on:click={() => handleEditAttribute(attr)}
												>
													<Pencil class="h-4 w-4" />
													<span class="sr-only">Edit attribute</span>
												</Button>
												<Button 
													variant="ghost" 
													size="icon" 
													class="h-8 w-8" 
													on:click={() => handleAttributeDelete(attr.attribute_id)}
												>
													<Trash2 class="h-4 w-4" />
													<span class="sr-only">Delete attribute</span>
												</Button>
											</div>
										</Table.Cell>
									</Table.Row>
								{/each}
							</Table.Body>
						</Table.Root>
					</Card.Content>
					<Card.Footer class="justify-center border-t p-4">
							<Button 
								size="sm" 
								variant="ghost" 
								class="gap-1" 
								on:click={() => showProductAttributeForm = true}
							>
								<CirclePlus class="h-3.5 w-3.5" />
								Add Product Attribute
							</Button>
					</Card.Footer>
				</Card.Root>

				<Card.Root>
					<Card.Header>
						<Card.Title>Variation Attributes</Card.Title>
						<Card.Description>Attributes that can vary between different versions of the product</Card.Description>
					</Card.Header>
					<Card.Content class="max-h-[400px] overflow-y-auto">
						<Table.Root class="w-full">
							<Table.Header>
								<Table.Row>
									<Table.Head>Title</Table.Head>
									<Table.Head>Type</Table.Head>
									<Table.Head>Required</Table.Head>
									<Table.Head>Options</Table.Head>
									<Table.Head class="w-[50px]">Actions</Table.Head>
								</Table.Row>
							</Table.Header>
							<Table.Body>
								{#each variationAttributes as attr}
									<Table.Row>
										<Table.Cell>{attr.title}</Table.Cell>
										<Table.Cell>{attr.data_type}</Table.Cell>
										<Table.Cell>
											{#if attr.required}
												<Badge variant="default">Required</Badge>
											{:else}
												<Badge variant="outline">Optional</Badge>
											{/if}
										</Table.Cell>
										<Table.Cell>
											{#if attr.options}
												<Select.Root>
													<Select.Trigger>
														<Select.Value placeholder="Select option" />
													</Select.Trigger>
													<Select.Content>
														{#each attr.options as option}
															<Select.Item value={option.attribute_option_id.toString()}>
																{#if attr.data_type === 'Color'}
																	<div class="flex items-center gap-2">
																		<div class="h-4 w-4 rounded-full" style="background-color: {option.value}"></div>
																		{option.value}
																	</div>
																{:else}
																	{option.value}
																{/if}
															</Select.Item>
														{/each}
													</Select.Content>
												</Select.Root>
											{/if}
										</Table.Cell>
										<Table.Cell>
											<div class="flex gap-2">
												<Button 
													variant="ghost" 
													size="icon" 
													class="h-8 w-8" 
													on:click={() => handleEditAttribute(attr)}
												>
													<Pencil class="h-4 w-4" />
													<span class="sr-only">Edit attribute</span>
												</Button>
												<Button 
													variant="ghost" 
													size="icon" 
													class="h-8 w-8" 
													on:click={() => handleAttributeDelete(attr.attribute_id)}
												>
													<Trash2 class="h-4 w-4" />
													<span class="sr-only">Delete attribute</span>
												</Button>
											</div>
										</Table.Cell>
									</Table.Row>
								{/each}
							</Table.Body>
						</Table.Root>
					</Card.Content>
					<Card.Footer class="justify-center border-t p-4">
							<Button 
								size="sm" 
								variant="ghost" 
								class="gap-1"
								on:click={() => showVariationAttributeForm = true}
							>
								<CirclePlus class="h-3.5 w-3.5" />
								Add Variation Attribute
							</Button>
					</Card.Footer>
				</Card.Root>
			</div>
		</div>
		<div class="flex items-center justify-center gap-2 md:hidden">
			<Button variant="outline" size="sm" on:click={() => showDeleteConfirmation = true}>Delete</Button>
			<Button size="sm" on:click={handleUpdateProductType}>Save Type</Button>
		</div>
	</div>
</main>

{#if showProductAttributeForm}
	<Dialog.Root open={showProductAttributeForm} onOpenChange={handleProductDialogChange}>
		<AttributeForm 
			typeId={parseInt($page.params.type)}
			appliesTo="Product"
			attribute={editingAttribute ?? {}}
			onClose={() => {
				showProductAttributeForm = false;
				editingAttribute = null;
			}}
			onSave={handleAttributeSaved}
		/>
	</Dialog.Root>
{/if}

{#if showVariationAttributeForm}
	<Dialog.Root open={showVariationAttributeForm} onOpenChange={handleVariationDialogChange}>
		<AttributeForm 
			typeId={parseInt($page.params.type)}
			appliesTo="ProductVariation"
			attribute={editingAttribute ?? {}}
			onClose={() => {
				showVariationAttributeForm = false;
				editingAttribute = null;
			}}
			onSave={handleAttributeSaved}
		/>
	</Dialog.Root>
{/if}

{#if showDeleteConfirmation}
	<Dialog.Root open={showDeleteConfirmation} onOpenChange={handleDeleteDialogChange}>
		<Dialog.Content class="sm:max-w-md">
			<Dialog.Header>
				<Dialog.Title class="text-center">Delete Product Type</Dialog.Title>
				<Dialog.Description class="text-center">
					This action is irreversible and will permanently delete this product type and all associated products.
				</Dialog.Description>
			</Dialog.Header>
			<div class="flex flex-col gap-4 py-4">
				<div class="flex justify-center">
					<div class="flex h-12 w-12 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/20">
						<AlertTriangle class="h-6 w-6 text-red-600 dark:text-red-400" />
					</div>
				</div>
				
				<div class="text-center text-sm text-red-600 dark:text-red-400">
					<p>Deleting "{$productTypeQuery.data?.title}" will:</p>
					<ul class="mt-2 list-disc pl-5 text-left">
						<li>Permanently remove all products of this type</li>
						<li>Delete all associated attributes and variations</li>
						<li>Remove all product data from the database</li>
					</ul>
				</div>
				
				<div class="mt-2">
					<Label for="confirm-delete" class="mb-2 block text-sm">
						Type <span class="font-mono font-bold">{expectedDeleteConfirmText}</span> to confirm:
					</Label>
					<Input 
						id="confirm-delete" 
						bind:value={deleteConfirmText} 
						placeholder={expectedDeleteConfirmText}
						class="w-full"
					/>
				</div>
			</div>
			<Dialog.Footer class="flex flex-col-reverse sm:flex-row sm:justify-between">
				<Button variant="ghost" on:click={() => showDeleteConfirmation = false}>Cancel</Button>
				<Button 
					variant="destructive" 
					on:click={handleDeleteProductType} 
					disabled={!isDeleteConfirmValid}
				>
					Delete Product Type
				</Button>
			</Dialog.Footer>
		</Dialog.Content>
	</Dialog.Root>
{/if}
