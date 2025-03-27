<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import * as Form from '$lib/components/ui/form';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea';
	import type { PageData } from './$types';
	import { api } from '$lib/api'
	import type { Product } from '$lib/types'
	import { getContext } from 'svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';

	export let data: PageData
	const authFetch:(input: RequestInfo | URL, init?: RequestInit | undefined) => Promise<Response> = getContext('authFetch')

	$:productQuery = createQuery<Product>({
		queryKey: ['shop', "gossip"],
		queryFn: () => api(authFetch).getShop('gossip'),
	})
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>General Settings</Card.Title>
	</Card.Header>
	<Card.Content>
		<form method="POST" class="space-y-8"  id="general-settings-form">
			<div class="flex w-full flex-col gap-3">
					<Label for="store-name">Store Name</Label>
					<Input />
					<p class="text-muted-foreground text-sm">Enter your shop name.</p>
			</div>
			<div class="flex w-full flex-col gap-3">
				<Label for="email">Email</Label>
				<Input type="email" id="email" />
				<p class="text-muted-foreground text-sm">Store contact email address.</p>
			</div>

			<div class="flex w-full flex-col gap-3">
				<Label for="address">Address</Label>
				<Textarea id="address" />
				<p class="text-muted-foreground text-sm">Store physical address.</p>
			</div>

			<div class="flex w-full flex-col gap-3">
				<Label for="about">About</Label>
				<Textarea id="about" />
				<p class="text-muted-foreground text-sm">Tell customers about your store.</p>
			</div>

			<div class="flex w-full flex-col gap-3">
				<Label for="phone">Contact Phone Number</Label>
				<Input id="phone" type="tel" />
				<p class="text-muted-foreground text-sm">Store contact number.</p>
			</div>
			<div class="flex w-full flex-col gap-3">
				<Label for="currency">Currency</Label>
				<Select.Root>
					<Select.Trigger id="currency" aria-label="currency">
						<Select.Value  />
					</Select.Trigger>
					<Select.Content>
							<Select.Item value="NGN" label="NGN">
								NGN
							</Select.Item>
							<Select.Item value="USD" label="USD">
								USD
							</Select.Item>
					</Select.Content>
				</Select.Root>
			
				<p class="text-muted-foreground text-sm">Store's primary currency.</p>
			</div>
			<!-- <div>
		<Form.Fieldset {form} name="phones">
			<Form.Legend>Contact Phones</Form.Legend>
			{#each $formData.phones as _, i}
				<Form.ElementField {form} name="phones[{i}]">
					<Form.Description>Add store contact phone numbers</Form.Description>
					<Form.Control let:attrs>
						<Input {...attrs} bind:value={$formData.phones[i]} />
					</Form.Control>
					<Form.FieldErrors />
				</Form.ElementField>
			{/each}
		</Form.Fieldset>
		<Button type="button" variant="outline" size="sm" class="mt-2" on:click={addPhone}>
			Add Phone
		</Button>
	</div> -->

			<Form.Button>Update</Form.Button>
		</form>
	</Card.Content>
	<Card.Footer></Card.Footer>
</Card.Root>
