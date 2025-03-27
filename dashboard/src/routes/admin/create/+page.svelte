<script>
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
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
				<Input id="title" type="title" placeholder="Naytife Clothing" required />
			</div>
			<div class="grid gap-2">
				<Label for="domain">Domain</Label>
				<Input id="domain" type="domain" placeholder="Wintheiser.naytife.com" required />
				<span>Domain availabe/Unavailabe</span>
			</div>
		</Card.Content>
		<Card.Footer>
			<Button class="w-full">Create</Button>
		</Card.Footer>
	</Card.Root>
</div>
