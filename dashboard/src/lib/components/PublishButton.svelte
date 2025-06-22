<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { publishState } from '$lib/stores/publishState';
	import { toast } from 'svelte-sonner';
	import { getContext } from 'svelte';
	import { Loader2, Upload, CheckCircle, AlertCircle } from 'lucide-svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import { currentShopId } from '$lib/api';

	export let variant: 'default' | 'outline' | 'secondary' = 'default';
	export let size: 'sm' | 'default' | 'lg' = 'default';

	const authFetch = getContext('authFetch') as typeof fetch;
	let showConfirmDialog = false;

	$: canPublish = $publishState.hasUnpublishedChanges && !$publishState.isPublishing;
	$: buttonText = $publishState.isPublishing 
		? 'Publishing...' 
		: $publishState.hasUnpublishedChanges 
			? 'Publish Changes' 
			: 'Published';

	async function handlePublish() {
		if (!canPublish) return;
		
		showConfirmDialog = true;
	}

	async function confirmPublish() {
		try {
			publishState.setPublishing(true);
			showConfirmDialog = false;

			// Get shop ID from the store
			let shopId: number | null = null;
			const unsubscribe = currentShopId.subscribe(id => {
				shopId = id;
			});
			unsubscribe();

			if (!shopId) {
				throw new Error('Shop ID not found. Please refresh the page and try again.');
			}

			// Call publish API endpoint using the shop ID
			const response = await authFetch(`http://127.0.0.1:8080/v1/shops/${shopId}/publish`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({
					template_name: 'template_1', // Default template for now
					changes: $publishState.changesSince,
				}),
			});

			if (!response.ok) {
				const errorData = await response.json().catch(() => ({ message: 'Unknown error' }));
				throw new Error(errorData.message || 'Failed to trigger publish');
			}

			const result = await response.json();
			
			// Mark as published
			publishState.markAsPublished();
			
			toast.success('Store published successfully! Changes are now live.');
			
		} catch (error) {
			console.error('Publish error:', error);
			const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';
			publishState.setPublishing(false, errorMessage);
			toast.error(`Failed to publish: ${errorMessage}`);
		}
	}

	function getChangesSummary(): string {
		if ($publishState.changesSince.length === 0) return 'No changes to publish';
		
		const counts = $publishState.changesSince.reduce((acc, change) => {
			acc[change.type] = (acc[change.type] || 0) + 1;
			return acc;
		}, {} as Record<string, number>);
		
		const parts = [];
		if (counts.shop_update) parts.push(`${counts.shop_update} store update${counts.shop_update > 1 ? 's' : ''}`);
		if (counts.product_create) parts.push(`${counts.product_create} product${counts.product_create > 1 ? 's' : ''} added`);
		if (counts.product_update) parts.push(`${counts.product_update} product${counts.product_update > 1 ? 's' : ''} updated`);
		if (counts.product_delete) parts.push(`${counts.product_delete} product${counts.product_delete > 1 ? 's' : ''} deleted`);
		if (counts.image_update) parts.push(`${counts.image_update} image${counts.image_update > 1 ? 's' : ''} updated`);
		
		return parts.join(', ') || 'Various updates';
	}
</script>

<Button 
	{variant}
	{size}
	disabled={!canPublish}
	on:click={handlePublish}
	class="relative {canPublish ? 'animate-pulse shadow-lg' : ''}"
>
	{#if $publishState.isPublishing}
		<Loader2 class="w-4 h-4 mr-2 animate-spin" />
	{:else if $publishState.hasUnpublishedChanges}
		<Upload class="w-4 h-4 mr-2" />
	{:else}
		<CheckCircle class="w-4 h-4 mr-2 text-green-500" />
	{/if}
	{buttonText}
</Button>

{#if $publishState.publishError}
	<div class="flex items-center gap-2 mt-2 text-sm text-red-600">
		<AlertCircle class="w-4 h-4" />
		{$publishState.publishError}
	</div>
{/if}

<!-- Confirmation Dialog -->
<Dialog.Root bind:open={showConfirmDialog}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>Publish Changes</Dialog.Title>
			<Dialog.Description>
				Are you sure you want to publish your changes? This will make them live on your store.
			</Dialog.Description>
		</Dialog.Header>
		
		<div class="py-4">
			<h4 class="font-medium mb-2">Changes to be published:</h4>
			<div class="bg-muted p-3 rounded-md text-sm">
				{getChangesSummary()}
			</div>
			
			{#if $publishState.changesSince.length > 0}
				<div class="mt-3 max-h-32 overflow-y-auto">
					<h5 class="text-xs font-medium text-muted-foreground mb-2">Recent changes:</h5>
					{#each $publishState.changesSince.slice(-5) as change}
						<div class="text-xs py-1 border-b border-border/30 last:border-0">
							<span class="font-medium">{change.description}</span>
							<span class="text-muted-foreground ml-1">
								({change.timestamp.toLocaleTimeString()})
							</span>
						</div>
					{/each}
				</div>
			{/if}
		</div>
		
		<Dialog.Footer>
			<Button variant="outline" on:click={() => showConfirmDialog = false}>
				Cancel
			</Button>
			<Button on:click={confirmPublish}>
				<Upload class="w-4 h-4 mr-2" />
				Publish Now
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
