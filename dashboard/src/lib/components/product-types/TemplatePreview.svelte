<script lang="ts">
	import type { PredefinedProductType } from '$lib/types';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Eye, Package, Layers } from 'lucide-svelte';

	export let template: PredefinedProductType;
	export let onCreateFromTemplate: (template: PredefinedProductType) => void;
	export let loading = false;

	let dialogOpen = false;

	// Group attributes by applies_to
	$: groupedAttributes = template.attributes.reduce((acc, attr) => {
		if (!acc[attr.applies_to]) {
			acc[attr.applies_to] = [];
		}
		acc[attr.applies_to].push(attr);
		return acc;
	}, {} as Record<string, typeof template.attributes>);

	function handleCreate() {
		dialogOpen = false;
		onCreateFromTemplate(template);
	}
</script>

<Dialog.Root bind:open={dialogOpen}>
	<Dialog.Trigger asChild let:builder>
		<Button builders={[builder]} variant="outline" size="sm" class="gap-2">
			<Eye class="h-4 w-4" />
			Preview
		</Button>
	</Dialog.Trigger>
	<Dialog.Content class="max-w-4xl max-h-[80vh] overflow-y-auto">
		<Dialog.Header>
			<Dialog.Title class="flex items-center gap-3">
				<span class="text-2xl">{template.icon}</span>
				{template.title}
			</Dialog.Title>
			<Dialog.Description>
				{template.description}
			</Dialog.Description>
		</Dialog.Header>

		<div class="space-y-6">
			<!-- Template Properties -->
			<Card.Root>
				<Card.Header>
					<Card.Title class="text-lg">Template Properties</Card.Title>
				</Card.Header>
				<Card.Content>
					<div class="grid gap-4 sm:grid-cols-2">
						<div>
							<span class="text-sm font-medium text-muted-foreground">Category</span>
							<p class="text-sm">{template.category}</p>
						</div>
						<div>
							<span class="text-sm font-medium text-muted-foreground">SKU Prefix</span>
							<p class="text-sm font-mono bg-muted px-2 py-1 rounded">{template.sku_substring}</p>
						</div>
						<div>
							<span class="text-sm font-medium text-muted-foreground">Product Type</span>
							<div class="flex gap-2 mt-1">
								{#if template.shippable}
									<Badge variant="secondary">Physical Product</Badge>
								{/if}
								{#if template.digital}
									<Badge variant="secondary">Digital Product</Badge>
								{/if}
							</div>
						</div>
					</div>
				</Card.Content>
			</Card.Root>

			<!-- Attributes -->
			{#each Object.entries(groupedAttributes) as [appliesTo, attributes]}
				<Card.Root>
					<Card.Header>
						<Card.Title class="text-lg flex items-center gap-2">
							{#if appliesTo === 'Product'}
								<Package class="h-5 w-5" />
								Product Attributes
							{:else}
								<Layers class="h-5 w-5" />
								Variant Attributes
							{/if}
							<Badge variant="outline" class="ml-auto">{attributes.length} attributes</Badge>
						</Card.Title>
						<Card.Description>
							{#if appliesTo === 'Product'}
								Attributes that apply to the product as a whole (e.g., brand, material, description)
							{:else}
								Attributes that can vary between product variants (e.g., size, color, storage)
							{/if}
						</Card.Description>
					</Card.Header>
					<Card.Content>
						<div class="grid gap-4 sm:grid-cols-2">
							{#each attributes as attribute}
								<div class="border rounded-lg p-3">
									<div class="flex items-start justify-between">
										<div class="flex-1">
											<h4 class="font-medium">{attribute.title}</h4>
											<div class="flex items-center gap-2 mt-1">
												<Badge variant="outline" class="text-xs">
													{attribute.data_type}
												</Badge>
												{#if attribute.required}
													<Badge variant="destructive" class="text-xs">Required</Badge>
												{/if}
												{#if attribute.unit}
													<Badge variant="secondary" class="text-xs">{attribute.unit}</Badge>
												{/if}
											</div>
										</div>
									</div>
									{#if attribute.options && attribute.options.length > 0}
										<div class="mt-2">
											<span class="text-xs text-muted-foreground">Options:</span>
											<div class="flex flex-wrap gap-1 mt-1">
												{#each attribute.options.slice(0, 5) as option}
													<Badge variant="outline" class="text-xs">{option.value}</Badge>
												{/each}
												{#if attribute.options.length > 5}
													<Badge variant="outline" class="text-xs">
														+{attribute.options.length - 5} more
													</Badge>
												{/if}
											</div>
										</div>
									{/if}
								</div>
							{/each}
						</div>
					</Card.Content>
				</Card.Root>
			{/each}
		</div>

		<Dialog.Footer class="flex justify-between">
			<div class="text-sm text-muted-foreground">
				This template will create a complete product type with all attributes and options ready to use.
			</div>
			<div class="flex gap-2">
				<Button variant="outline" on:click={() => dialogOpen = false}>
					Cancel
				</Button>
				<Button on:click={handleCreate} disabled={loading}>
					{loading ? 'Creating...' : `Create ${template.title} Product Type`}
				</Button>
			</div>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
