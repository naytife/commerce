<script lang="ts" context="module">
	import { z } from 'zod';
	export const generalSettingsFormSchema = z.object({
		title: z
			.string()
			.min(2, 'title must be at least 2 characters.')
			.max(30, 'title must not be longer than 30 characters'),
		email: z.string({ required_error: 'Please enter an email address' }).email(),
		description: z.string().min(4).max(160).default('Our amazing shop.'),
		phone: z
			.string()
			.min(1, 'Please enter a valid phone number')
			.max(15, 'Please enter a valid phone number')
	});
	export type GeneralSettingsFormSchema = typeof generalSettingsFormSchema;
</script>

<script lang="ts">
	import { type Infer, type SuperValidated, superForm } from 'sveltekit-superforms';
	import { zodClient } from 'sveltekit-superforms/adapters';
	import * as Card from '$lib/components/ui/card';
	import * as Form from '$lib/components/ui/form';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea';
	import { get, writable } from 'svelte/store';

	// Assuming this is your Houdini's store data
	export let data: {
		ShopQuery: {
			store: any; // Replace 'any' with the correct type if known
		};
	};

	// Extract the data from the Houdini store
	const shopData = get(data);
	console.log(shopData);

	// Check if the data is available and map it to form fields
	const initialFormData = shopData
		? {
				title: shopData.data.shop.title || '',
				email: shopData.data.shop.contactEmail || '',
				description: shopData.shop.shop.about || '',
				phone: shopData.data.shop.contactPhone || ''
			}
		: {
				title: '',
				email: '',
				description: '',
				phone: ''
			};

	// Initialize the writable store with the form data
	const formData = writable(initialFormData);

	const form = superForm(formData, {
		validators: zodClient(generalSettingsFormSchema),
		dataType: 'json'
	});

	const { form: formFields, enhance } = form;

	// function addPhone() {
	// 	$formData.phones = [...$formData.phones, ''];

	// 	tick().then(() => {
	// 		const phoneInputs = Array.from(
	// 			document.querySelectorAll<HTMLElement>("#general-settings-form input[name='phones']")
	// 		);
	// 		const lastInput = phoneInputs[phoneInputs.length - 1];
	// 		lastInput && lastInput.focus();
	// 	});
	// }
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>General Settings</Card.Title>
	</Card.Header>
	<Card.Content>
		<form method="POST" class="space-y-8" use:enhance id="general-settings-form">
			<Form.Field {form} name="title">
				<Form.Control let:attrs>
					<Form.Label>title</Form.Label>
					<Input placeholder="My Shop" {...attrs} bind:value={$formData.title} />
				</Form.Control>
				<Form.Description>This is the title of your shop.</Form.Description>
				<Form.FieldErrors />
			</Form.Field>

			<Form.Field {form} name="email">
				<Form.Control let:attrs>
					<Form.Label>contact email</Form.Label>
					<Input placeholder="info@myshop.com" {...attrs} bind:value={$formData.email} />
				</Form.Control>
				<Form.Description
					>This is the email your customers will use to contact you.</Form.Description
				>
				<Form.FieldErrors />
			</Form.Field>
			<Form.Field {form} name="phone">
				<Form.Control let:attrs>
					<Form.Label>contact phone</Form.Label>
					<Input placeholder="(+234) 00 000 0000" {...attrs} bind:value={$formData.phone} />
				</Form.Control>
				<Form.Description
					>This is the phone number your customers will use to contact you.</Form.Description
				>
				<Form.FieldErrors />
			</Form.Field>

			<Form.Field {form} name="description">
				<Form.Control let:attrs>
					<Form.Label>Short description</Form.Label>
					<Textarea {...attrs} bind:value={$formData.description} />
				</Form.Control>
				<Form.Description>This will be displayed on the about page of the shop.</Form.Description>
				<Form.FieldErrors />
			</Form.Field>
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
