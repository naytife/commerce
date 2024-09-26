<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Store } from 'lucide-svelte';
	import { get } from 'svelte/store';

	/* @type { import('./$houdini').PageData } */
	export let data;
	$: ({ AccountQuery } = data);
</script>

<div class="container relative flex h-svh flex-col items-center justify-center lg:px-0">
	<Card.Root class="w-full max-w-sm">
		<Card.Header>
			<Card.Title class="text-2xl">My Store</Card.Title>
			<Card.Description>Select your store</Card.Description>
		</Card.Header>
		<Card.Content class="grid gap-4">
			<div>
				{#each $AccountQuery.data.myShops as shop, idx (idx)}
					<div class="mb-4 grid grid-cols-[25px_1fr] items-start gap-1 pb-4 last:mb-0 last:pb-0">
						<Store class="h-5 w-5" />
						<a href="/admin/{shop.id}" class="space-y-1">
							<p class="text-sm font-medium leading-none">
								{shop.title}
							</p>
							<p class="text-sm text-muted-foreground">
								{shop.domain}
							</p>
						</a>
					</div>
				{/each}
			</div>
		</Card.Content>
		<Card.Footer>
			<Button class="w-full" href="/account/s/create">Create a new store</Button>
		</Card.Footer>
	</Card.Root>
</div>
