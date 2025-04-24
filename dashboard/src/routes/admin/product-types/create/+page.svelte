<script lang="ts">
	import ChevronLeft from 'lucide-svelte/icons/chevron-left';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { api } from '$lib/api';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';

	// Define the ProductType interface to include the skuSubstring property
	interface ProductTypeCreate {
		title: string;
		sku_substring: string;
		shippable: boolean;
		digital: boolean;
	}

	let title = '';
	let skuSubstring = '';
	let shippable = false;
	let digital = false;

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

	async function handleSave() {
		if (skuSubstring && !validateSkuSubstring(skuSubstring)) {
			toast.error('Please fix the SKU substring before saving');
			return;
		}

		try {
			const productTypeData: ProductTypeCreate = {
				title,
				sku_substring: skuSubstring,
				shippable,
				digital
			};
			
			const response = await api().createProductType(productTypeData);
			toast.success('Product type created successfully');
			goto(`/admin/product-types/${response.id}/edit`);
		} catch (error) {
			toast.error('Failed to create product type');
			console.error('Error creating product type:', error);
		}
	}

	function handleDiscard() {
		goto('/admin/product-types');
	}

	// Force uppercase for SKU input
	function handleSkuInput(event: Event) {
		const input = event.target as HTMLInputElement;
		skuSubstring = input.value.toUpperCase();
	}
</script>

<main class="grid flex-1 items-start gap-4 p-4 sm:px-6 sm:py-0 md:gap-8">
	<div class="mx-auto grid max-w-7xl flex-1 auto-rows-max gap-4">
		<div class="flex items-center gap-4">
			<Button variant="outline" size="icon" class="h-7 w-7" on:click={handleDiscard}>
				<ChevronLeft class="h-4 w-4" />
				<span class="sr-only">Back</span>
			</Button>

			<div class="hidden items-center gap-2 md:ml-auto md:flex">
				<Button variant="outline" size="sm" on:click={handleDiscard}>Discard</Button>
				<Button size="sm" on:click={handleSave}>Save Type</Button>
			</div>
		</div>
		<div class="grid gap-4 md:grid-cols-1 lg:gap-8">
			<div class="grid auto-rows-max items-start gap-4">
				<Card.Root>
					<Card.Header>
						<Card.Title>Product Type</Card.Title>
						<Card.Description>
							Create a new product type by setting its basic properties. You can add attributes after creation.
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
								<Label for="physical">Shippable product?</Label>
								<Switch id="physical" bind:checked={shippable} />
							</div>
							<div class="flex items-center space-x-2">
								<Label for="digital">Digital product?</Label>
								<Switch id="digital" bind:checked={digital} />
							</div>
						</div>
					</Card.Content>
				</Card.Root>
			</div>
		</div>
		<div class="flex items-center justify-center gap-2 md:hidden">
			<Button variant="outline" size="sm" on:click={handleDiscard}>Discard</Button>
			<Button size="sm" on:click={handleSave}>Save Type</Button>
		</div>
	</div>
</main>
