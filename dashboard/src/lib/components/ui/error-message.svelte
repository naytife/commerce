<script lang="ts">
	import { AlertTriangle, RefreshCw } from 'lucide-svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';

	export let title = 'Something went wrong';
	export let message = 'An error occurred while loading the content.';
	export let onRetry: (() => void) | null = null;
</script>

<Card.Root class="border-destructive/20 bg-destructive/5">
	<Card.Content class="flex items-center gap-4 p-6">
		<div class="flex-shrink-0">
			<AlertTriangle class="h-8 w-8 text-destructive" />
		</div>
		<div class="flex-1 space-y-1">
			<h3 class="font-semibold text-destructive">{title}</h3>
			<p class="text-sm text-muted-foreground">{message}</p>
		</div>
		{#if onRetry}
			<Button variant="outline" size="sm" on:click={onRetry} class="gap-2">
				<RefreshCw class="h-4 w-4" />
				Retry
			</Button>
		{/if}
	</Card.Content>
</Card.Root>
