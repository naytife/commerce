<script>
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { createShop } from '$lib/api';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import { getContext } from 'svelte';

	let title = '';
	let subdomain = '';
	// Slugify function for subdomain
	$: slug = subdomain
		.trim()
		.toLowerCase()
		.replace(/\s+/g, '-') // replace spaces with dashes
		.replace(/[^a-z0-9-]/g, ''); // remove invalid chars

	let loading = false;

	const authFetch = getContext('authFetch');

	async function handleCreate() {
		if (!title.trim() || !slug.trim()) {
			toast.error('Please provide both a title and subdomain');
			return;
		}
		loading = true;
		try {
			await createShop({ subdomain: slug, title }, authFetch);
			toast.success('Store created!');
			goto('/account');
		} catch (error) {
			console.error(error);
			toast.error('Failed to create store');
		} finally {
			loading = false;
		}
	}
</script>

<div class="container relative flex h-svh flex-col items-center justify-center lg:px-0">
	<Card.Root class="w-full max-w-sm">
		<Card.Header>
			<Card.Title class="text-2xl">Create Store</Card.Title>
			<Card.Description>Enter store details below to get started.</Card.Description>
		</Card.Header>
		<Card.Content class="grid gap-4">
			<div class="grid gap-2">
				<Label for="title">Title</Label>
				<Input id="title" type="text" placeholder="Naytife Clothing" required bind:value={title} />
			</div>
			<div class="grid gap-2">
				<Label for="subdomain">Subdomain</Label>
				<Input
					id="subdomain"
					type="text"
					placeholder="yourstore"
					required
					bind:value={subdomain}
				/>
				{#if slug}
					<span class="text-sm text-gray-500">
						Your store domain will be:
						<b class="text-green-600">{slug}.naytife.com</b>
					</span>
				{/if}
			</div>
		</Card.Content>
		<Card.Footer>
			<Button class="w-full" on:click={handleCreate} disabled={loading}>
				{loading ? 'Creating...' : 'Create'}
			</Button>
		</Card.Footer>
	</Card.Root>
</div>
