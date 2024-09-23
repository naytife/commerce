<script lang="ts">
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

	import { type SuperValidated, type Infer, superForm } from 'sveltekit-superforms';
	import { zodClient } from 'sveltekit-superforms/adapters';
	import { formSchema, type FormSchema } from './schema';

	export let data: SuperValidated<Infer<FormSchema>>;

	const form = superForm(data, {
		validators: zodClient(formSchema),
		dataType: 'json'
	});
	const { form: formData } = form;

	$: console.log('Current formData:', $formData);
</script>

<main class="grid flex-1 items-start gap-4 p-4 sm:px-6 sm:py-0 md:gap-8">
	<div class="mx-auto grid max-w-[59rem] flex-1 auto-rows-max gap-4">
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
						<Card.Description>Lipsum dolor sit amet, consectetur adipiscing elit</Card.Description>
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
							<Form.Fieldset {form} name="size">
								<Form.Legend>Screen</Form.Legend>
								{#each $formData.form.data.size as _, i}
									<Form.ElementField {form} name="size[{i}]">
										<Form.Description class={i !== 0 ? 'sr-only' : ''}>
											Add links to your website, blog, or social media profiles.
										</Form.Description>
										<Form.Control let:attrs>
											<Input {...attrs} bind:value={$formData.form.data.size[i]} />
										</Form.Control>
										<Form.FieldErrors />
									</Form.ElementField>
								{/each}
							</Form.Fieldset>
							<Button type="button" variant="outline" size="sm" class="mt-2">Add Option</Button>
						</div>
					</Card.Content>
				</Card.Root>
				<Card.Root>
					<Card.Header>
						<Card.Title>Variations</Card.Title>
						<Card.Description>Lipsum dolor sit amet, consectetur adipiscing elit</Card.Description>
					</Card.Header>
					<Card.Content>
						<Table.Root>
							<Table.Header>
								<Table.Row>
									<Table.Head class="w-[100px]">SKU</Table.Head>
									<Table.Head>Stock</Table.Head>
									<Table.Head>Price</Table.Head>
									<Table.Head class="w-[100px]">Size</Table.Head>
								</Table.Row>
							</Table.Header>
							<Table.Body>
								<Table.Row>
									<Table.Cell class="font-semibold">GGPC-001</Table.Cell>
									<Table.Cell>
										<Label for="stock-1" class="sr-only">Stock</Label>
										<Input id="stock-1" type="number" value="100" />
									</Table.Cell>
									<Table.Cell>
										<Label for="price-1" class="sr-only">Price</Label>
										<Input id="price-1" type="number" value="99.99" />
									</Table.Cell>
									<Table.Cell>
										<ToggleGroup.Root type="single" value="s" variant="outline">
											<ToggleGroup.Item value="s">S</ToggleGroup.Item>
											<ToggleGroup.Item value="m">M</ToggleGroup.Item>
											<ToggleGroup.Item value="l">L</ToggleGroup.Item>
										</ToggleGroup.Root>
									</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell class="font-semibold">GGPC-002</Table.Cell>
									<Table.Cell>
										<Label for="stock-2" class="sr-only">Stock</Label>
										<Input id="stock-2" type="number" value="143" />
									</Table.Cell>
									<Table.Cell>
										<Label for="price-2" class="sr-only">Price</Label>
										<Input id="price-2" type="number" value="99.99" />
									</Table.Cell>
									<Table.Cell>
										<ToggleGroup.Root type="single" value="m" variant="outline">
											<ToggleGroup.Item value="s">S</ToggleGroup.Item>
											<ToggleGroup.Item value="m">M</ToggleGroup.Item>
											<ToggleGroup.Item value="l">L</ToggleGroup.Item>
										</ToggleGroup.Root>
									</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell class="font-semibold">GGPC-003</Table.Cell>
									<Table.Cell>
										<Label for="stock-3" class="sr-only">Stock</Label>
										<Input id="stock-3" type="number" value="32" />
									</Table.Cell>
									<Table.Cell>
										<Label for="price-3" class="sr-only">Stock</Label>
										<Input id="price-3" type="number" value="99.99" />
									</Table.Cell>
									<Table.Cell>
										<ToggleGroup.Root type="single" value="s" variant="outline">
											<ToggleGroup.Item value="s">S</ToggleGroup.Item>
											<ToggleGroup.Item value="m">M</ToggleGroup.Item>
											<ToggleGroup.Item value="l">L</ToggleGroup.Item>
										</ToggleGroup.Root>
									</Table.Cell>
								</Table.Row>
							</Table.Body>
						</Table.Root>
					</Card.Content>
					<Card.Footer class="justify-center border-t p-4">
						<Button size="sm" variant="ghost" class="gap-1">
							<CirclePlus class="h-3.5 w-3.5" />
							Add Variant
						</Button>
					</Card.Footer>
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
