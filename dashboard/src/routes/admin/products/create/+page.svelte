<script lang="ts">
	import Trash2 from 'lucide-svelte/icons/trash-2';
	import ChevronLeft from 'lucide-svelte/icons/chevron-left';
	import Upload from 'lucide-svelte/icons/upload';

	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Form from '$lib/components/ui/form';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import * as Table from '$lib/components/ui/table';
	import { Textarea } from '$lib/components/ui/textarea';
	import * as ToggleGroup from '$lib/components/ui/toggle-group';
	import { CirclePlus } from 'lucide-svelte';
	import { getContext } from 'svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { api } from '$lib/api';
	import type { Product } from '$lib/types';
	import { derived } from 'svelte/store';
	import { ScrollArea } from '$lib/components/ui/scroll-area';

	export let data: PageData
	const authFetch = getContext('authFetch')

	$:typeAttributesQuery = createQuery<Product>({
		queryKey: ['product-type-attributes', 4],
		queryFn: () => api(authFetch).getProductById(4),
	})

	// Generate all combinations of variant attributes
	function generateCombinations(attributes) {
		if (attributes.length === 0) return [[]];
		const [first, ...rest] = attributes;
		const combinations = generateCombinations(rest);
		return first.options.flatMap(option =>
			combinations.map(combination => [{ title: first.title, value: option.value }, ...combination])
		);
	}

	$: variantAttributes = $typeAttributesQuery.data?.filter(attr => attr.applies_to === 'ProductVariation') || [];
	$: combinations = generateCombinations(variantAttributes);

	// Track disabled state for each variation
	let disabledVariations = [];

	// Separate lists for active and disabled variations
	let activeCombinations = [];
	let disabledCombinations = [];

	$: activeCombinations = combinations.filter((_, index) => !disabledVariations[index]);
	$: disabledCombinations = combinations.filter((_, index) => disabledVariations[index]);

	// Toggle the disabled state of a variation and move it between lists
	function toggleVariation(index, isDisabledList = false) {
		const sourceList = isDisabledList ? disabledCombinations : activeCombinations;
		const targetList = isDisabledList ? activeCombinations : disabledCombinations;

		// Move the variation between lists
		const [movedCombination] = sourceList.splice(index, 1);
		targetList.push(movedCombination);

		// Update the disabledVariations array
		const globalIndex = combinations.indexOf(movedCombination);
		disabledVariations[globalIndex] = !disabledVariations[globalIndex];
	}
</script>

<main class="grid flex-1 items-start gap-4 p-4 sm:px-6 sm:py-0 md:gap-8">
	<div class="mx-auto grid max-w-7xl flex-1 auto-rows-max gap-4">
		<div class="flex items-center gap-4">
			<Button variant="outline" size="icon" class="h-7 w-7">
				<ChevronLeft class="h-4 w-4" />
				<span class="sr-only">Back</span>
			</Button>

			<div class="hidden items-center gap-2 md:ml-auto md:flex">
				<Button variant="outline" size="sm">Discard</Button>
				<Button size="sm">Save Product</Button>
			</div>
		</div>
		<div class="grid gap-4 md:grid-cols-[1fr_250px] lg:grid-cols-3 lg:gap-8">
			<div class="grid auto-rows-max items-start gap-4 lg:col-span-2 lg:gap-8">
				<Card.Root>
					<Card.Header>
						<Card.Title>Product Details</Card.Title>
						<Card.Description>Enter general product information</Card.Description>
					</Card.Header>
					<Card.Content>
						<div class="grid gap-6">
							<div class="grid gap-3">
								<Label for="name">Name</Label>
								<Input id="name" type="text" class="w-full" value="Gamer Gear Pro Controller" />
							</div>
							<div class="grid gap-3">
								<Label for="description">Description</Label>
								<Textarea
									id="description"
									value="Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam auctor, nisl nec ultricies ultricies, nunc nisl ultricies nunc, nec ultricies nunc nisl nec nunc."
									class="min-h-32"
								/>
							</div>
						</div>
					</Card.Content>
				</Card.Root>
				<Card.Root>
					<Card.Header>
						<Card.Title>Attributes</Card.Title>
						<Card.Description>Possible attributes for this product</Card.Description>
					</Card.Header>
					<Card.Content>
						<div>
							<form class="grid gap-4">
								{#if $typeAttributesQuery.data}
									{#each $typeAttributesQuery.data as attribute}
									{#if attribute.applies_to === 'Product'}
										<div class="flex w-full flex-col gap-3">
											<Label for={attribute.title}>{attribute.title}</Label>
											{#if attribute.options}
												<Select.Root>
													<Select.Trigger id={attribute.title} aria-label={attribute.title}>
														<Select.Value placeholder={`Select ${attribute.title}`} />
													</Select.Trigger>
													<Select.Content>
														{#each attribute.options as option}
															<Select.Item value={option.value} label={option.value}>
																{option.value}
															</Select.Item>
														{/each}
													</Select.Content>
												</Select.Root>
											{:else}
												<Input 
													type={attribute.data_type === 'Number' ? 'number' : 'text'}
													id={attribute.title}
													placeholder={attribute.title}
													class="w-full"
													required={attribute.required}
												/>
											{/if}
										</div>
									{/if}
									{/each}
								{/if}
							</form>

						</div>
					</Card.Content>
				</Card.Root>
				<Card.Root>
					<Card.Header>
						<Card.Title>Variations</Card.Title>
						<Card.Description>Define variations for all attribute combinations</Card.Description>
					</Card.Header>
					<Card.Content class="max-h-60 max-w-[53rem] overflow-y-auto">
						<h3 class="text-lg font-semibold">Active Variations</h3>
						<Table.Root class="relative w-full table-fixed">
							<Table.Header class="sticky top-0 z-50 bg-gray-900">
							  <Table.Row>
								<Table.Head class="w-[100px]">SKU</Table.Head>
								<Table.Head class="w-[100px]">Stock</Table.Head>
								<Table.Head class="w-[100px]">Price</Table.Head>
								{#each variantAttributes as attribute}
								  <Table.Head class="w-[100px]">{attribute.title}</Table.Head>
								{/each}
								<Table.Head class="w-[50px]">Action</Table.Head>
							  </Table.Row>
							</Table.Header>
							<Table.Body>
							  {#each activeCombinations as combination, index}
								<Table.Row>
								  <Table.Cell class="font-semibold whitespace-nowrap">
									<Label for={`sku-active-${index}`} class="sr-only">SKU</Label>
									<Input id={`sku-active-${index}`} type="text" placeholder="SKU" />
								  </Table.Cell>
								  <Table.Cell class="whitespace-nowrap">
									<Label for={`stock-active-${index}`} class="sr-only">Stock</Label>
									<Input id={`stock-active-${index}`} type="number" placeholder="Stock" />
								  </Table.Cell>
								  <Table.Cell class="whitespace-nowrap">
									<Label for={`price-active-${index}`} class="sr-only">Price</Label>
									<Input id={`price-active-${index}`} type="number" placeholder="Price" />
								  </Table.Cell>
								  {#each combination as attribute}
									<Table.Cell class="whitespace-nowrap">{attribute.value}</Table.Cell>
								  {/each}
								  <Table.Cell class="whitespace-nowrap">
									<Button
									  size="sm"
									  variant="ghost"
									  on:click={() => toggleVariation(index)}
									  class="gap-1"
									>
									  <Trash2 class="h-3.5 w-3.5" />
									</Button>
								  </Table.Cell>
								</Table.Row>
							  {/each}
							  {#each disabledCombinations as combination, index}
							  <Table.Row class="opacity-50 line-through">
								<Table.Cell class="font-semibold whitespace-nowrap">
								  <Label for={`sku-disabled-${index}`} class="sr-only">SKU</Label>
								  <Input id={`sku-disabled-${index}`} type="text" placeholder="SKU" disabled />
								</Table.Cell>
								<Table.Cell class="whitespace-nowrap">
								  <Label for={`stock-disabled-${index}`} class="sr-only">Stock</Label>
								  <Input id={`stock-disabled-${index}`} type="number" placeholder="Stock" disabled />
								</Table.Cell>
								<Table.Cell class="whitespace-nowrap">
								  <Label for={`price-disabled-${index}`} class="sr-only">Price</Label>
								  <Input id={`price-disabled-${index}`} type="number" placeholder="Price" disabled />
								</Table.Cell>
								{#each combination as attribute}
								  <Table.Cell class="whitespace-nowrap">{attribute.value}</Table.Cell>
								{/each}
								<Table.Cell class="whitespace-nowrap">
								  <Button
									size="sm"
									variant="ghost"
									on:click={() => toggleVariation(index, true)}
									class="gap-1"
								  >
									<CirclePlus class="h-3.5 w-3.5" />
								  </Button>
								</Table.Cell>
							  </Table.Row>
							{/each}
							</Table.Body>
						  </Table.Root>
					  
					  </Card.Content>
				</Card.Root>
				<Card.Root>
					<Card.Header>
						<Card.Title>Product Category</Card.Title>
					</Card.Header>
					<Card.Content>
						<div class="grid gap-6 sm:grid-cols-3">
							<div class="grid gap-3">
								<Label for="category">Category</Label>
								<Select.Root>
									<Select.Trigger id="category" aria-label="Select category">
										<Select.Value placeholder="Select category" />
									</Select.Trigger>
									<Select.Content>
										<Select.Item value="clothing" label="Clothing">Clothing</Select.Item>
										<Select.Item value="electronics" label="Electronics">Electronics</Select.Item>
										<Select.Item value="accessories" label="Accessories">Accessories</Select.Item>
									</Select.Content>
								</Select.Root>
							</div>
							<div class="grid gap-3">
								<Label for="subcategory">Subcategory (optional)</Label>
								<Select.Root>
									<Select.Trigger id="subcategory" aria-label="Select subcategory">
										<Select.Value placeholder="Select subcategory" />
									</Select.Trigger>
									<Select.Content>
										<Select.Item value="t-shirts" label="T-Shirts">T-Shirts</Select.Item>
										<Select.Item value="hoodies" label="Hoodies">Hoodies</Select.Item>
										<Select.Item value="sweatshirts" label="Sweatshirts">Sweatshirts</Select.Item>
									</Select.Content>
								</Select.Root>
							</div>
						</div>
					</Card.Content>
				</Card.Root>
			</div>
			<div class="grid auto-rows-max items-start gap-4 lg:gap-8">
				<Card.Root>
					<Card.Header>
						<Card.Title>Product Status</Card.Title>
					</Card.Header>
					<Card.Content>
						<div class="grid gap-6">
							<div class="grid gap-3">
								<Label for="status">Status</Label>
								<Select.Root>
									<Select.Trigger id="status" aria-label="Select status">
										<Select.Value placeholder="Select status" />
									</Select.Trigger>
									<Select.Content>
										<Select.Item value="draft" label="Draft">Draft</Select.Item>
										<Select.Item value="published" label="Active">Active</Select.Item>
										<Select.Item value="archived" label="Archived">Archived</Select.Item>
									</Select.Content>
								</Select.Root>
							</div>
						</div>
					</Card.Content>
				</Card.Root>
				<Card.Root class="overflow-hidden">
					<Card.Header>
						<Card.Title>Product Images</Card.Title>
						<Card.Description>Lipsum dolor sit amet, consectetur adipiscing elit</Card.Description>
					</Card.Header>
					<Card.Content>
						<div class="grid gap-2">
							<img
								alt="Product"
								class="aspect-square w-full rounded-md object-cover"
								height="300"
								src="/images/placeholder.png"
								width="300"
							/>
							<div class="grid grid-cols-3 gap-2">
								<button>
									<img
										alt="Product"
										class="aspect-square w-full rounded-md object-cover"
										height="84"
										src="/images/placeholder.png"
										width="84"
									/>
								</button>
								<button>
									<img
										alt="Product"
										class="aspect-square w-full rounded-md object-cover"
										height="84"
										src="/images/placeholder.png"
										width="84"
									/>
								</button>
								<button
									class="flex aspect-square w-full items-center justify-center rounded-md border border-dashed"
								>
									<Upload class="h-4 w-4 text-muted-foreground" />
									<span class="sr-only">Upload</span>
								</button>
							</div>
						</div>
					</Card.Content>
				</Card.Root>
			</div>
		</div>
		<div class="flex items-center justify-center gap-2 md:hidden">
			<Button variant="outline" size="sm">Discard</Button>
			<Button size="sm">Save Product</Button>
		</div>
	</div>
</main>
