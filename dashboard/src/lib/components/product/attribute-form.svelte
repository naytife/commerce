<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import { Switch } from '$lib/components/ui/switch';
	import { api } from '$lib/api';
	import { toast } from 'svelte-sonner';
	import type { ProductTypeAttribute } from '$lib/types';
	import { X } from 'lucide-svelte';
	import { createEventDispatcher } from 'svelte';
	
	export let typeId: number;
	export let appliesTo: 'Product' | 'ProductVariation';
	export let onClose: () => void;
	export let onSave: (attribute: ProductTypeAttribute) => void;
	export let attribute: Partial<ProductTypeAttribute> = {}; // For edit mode

	let title = attribute.title || '';
	let dataType = attribute.data_type || 'Text';
	let required = attribute.required || false;
	let unit = attribute.unit || '';
	let options: string[] = attribute.options?.map(opt => opt.value) || [];
	let newOption = '';

	const dispatch = createEventDispatcher();

	function addOption() {
		if (newOption.trim()) {
			options = [...options, newOption.trim()];
			newOption = '';
		}
	}

	function removeOption(index: number) {
		options = options.filter((_, i) => i !== index);
	}

	async function handleSubmit() {
		try {
			const attributeData = {
				title,
				data_type: dataType,
				required,
				unit: unit || undefined,
				applies_to: appliesTo,
				options: options.map(value => ({ value })),
				product_type_id: typeId,
				attribute_id: attribute.attribute_id
			} as ProductTypeAttribute;

			let result: ProductTypeAttribute;
			
			if (attribute.attribute_id) {
				// Edit mode
				result = await api().updateProductTypeAttribute(typeId, attribute.attribute_id, attributeData);
				toast.success('Attribute updated successfully');
			} else {
				// Create mode
				result = await api().createProductTypeAttribute(typeId, attributeData);
				toast.success('Attribute created successfully');
			}

			onSave(result);
			onClose();
		} catch (error) {
			toast.error(attribute.attribute_id ? 'Failed to update attribute' : 'Failed to create attribute');
			console.error('Error saving attribute:', error);
		}
	}
</script>

<Dialog.Content class="sm:max-w-[425px]">
	<Dialog.Header>
		<Dialog.Title>{attribute.attribute_id ? 'Edit' : 'Add'} {appliesTo} Attribute</Dialog.Title>
		<Dialog.Description>
			{attribute.attribute_id ? 'Update' : 'Create'} an attribute that will be used to define properties of your {appliesTo.toLowerCase()}.
		</Dialog.Description>
	</Dialog.Header>
	<form on:submit|preventDefault={handleSubmit}>
		<div class="grid gap-4 py-4">
			<div class="grid gap-2">
				<Label for="title">Title</Label>
				<Input id="title" bind:value={title} required />
			</div>
			<div class="grid gap-2">
				<Label for="type">Type</Label>
				<Select.Root>
					<Select.Trigger>
						<Select.Value placeholder="Select type">
							{dataType}
						</Select.Value>
					</Select.Trigger>
					<Select.Content>
						<Select.Item value="Text" on:click={() => dataType = "Text"}>Text</Select.Item>
						<Select.Item value="Number" on:click={() => dataType = "Number"}>Number</Select.Item>
						<Select.Item value="Date" on:click={() => dataType = "Date"}>Date</Select.Item>
						<Select.Item value="Option" on:click={() => dataType = "Option"}>Option</Select.Item>
						<Select.Item value="Color" on:click={() => dataType = "Color"}>Color</Select.Item>
					</Select.Content>
				</Select.Root>
			</div>
			{#if dataType === 'Number'}
				<div class="grid gap-2">
					<Label for="unit">Unit</Label>
					<Select.Root>
						<Select.Trigger>
							<Select.Value>
								{unit || 'Select unit'}
							</Select.Value>
						</Select.Trigger>
						<Select.Content>
							<Select.Item value="KG" on:click={() => unit = "KG"}>KG</Select.Item>
							<Select.Item value="GB" on:click={() => unit = "GB"}>GB</Select.Item>
							<Select.Item value="INCH" on:click={() => unit = "INCH"}>INCH</Select.Item>
						</Select.Content>
					</Select.Root>
				</div>
			{/if}
			{#if dataType === 'Option'}
				<div class="grid gap-2">
					<Label>Options</Label>
					<div class="flex gap-2">
						<Input 
							placeholder="Enter option" 
							value={newOption}
							on:input={(e) => newOption = e.currentTarget.value}
							on:keydown={(e) => e.key === 'Enter' && (e.preventDefault(), addOption())}
						/>
						<Button type="button" variant="outline" on:click={addOption}>Add</Button>
					</div>
					<div class="mt-2 flex flex-wrap gap-2">
						{#each options as option, i}
							<div class="flex items-center gap-1 rounded-full bg-secondary px-3 py-1 text-sm">
								{option}
								<button
									type="button"
									class="ml-1 rounded-full p-1 hover:bg-secondary-foreground/10"
									on:click={() => removeOption(i)}
								>
									<X class="h-3 w-3" />
								</button>
							</div>
						{/each}
					</div>
				</div>
			{/if}
			{#if dataType === 'Color'}
				<div class="grid gap-2">
					<Label>Color Options</Label>
					<div class="flex gap-2">
						<Input 
							type="color"
							bind:value={newOption}
							class="h-10 w-20"
						/>
						<Button type="button" variant="outline" on:click={addOption}>Add Color</Button>
					</div>
					<div class="mt-2 flex flex-wrap gap-2">
						{#each options as option, i}
							<div class="flex items-center gap-1 rounded-full px-3 py-1 text-sm" style="background-color: {option}">
								<div class="w-4 h-4 rounded-full" style="background-color: {option}"></div>
								{option}
								<button
									type="button"
									class="ml-1 rounded-full p-1 hover:bg-secondary-foreground/10"
									on:click={() => removeOption(i)}
								>
									<X class="h-3 w-3" />
								</button>
							</div>
						{/each}
					</div>
				</div>
			{/if}
			<div class="flex items-center gap-2">
				<Switch id="required" bind:checked={required} />
				<Label for="required">Required</Label>
			</div>
		</div>
		<Dialog.Footer>
			<Button type="button" variant="outline" on:click={onClose}>Cancel</Button>
			<Button type="submit">{attribute.attribute_id ? 'Update' : 'Create'}</Button>
		</Dialog.Footer>
	</form>
</Dialog.Content>