<!-- <script lang="ts">
	import { type Infer, type SuperValidated, superForm } from 'sveltekit-superforms';
	import { zodClient } from 'sveltekit-superforms/adapters';
	import * as Card from '$lib/components/ui/card';
	import * as Form from '$lib/components/ui/form';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea';
	import { get, writable } from 'svelte/store';
	import type { PageData } from './$types';
	import { generalSettingsFormSchema } from './schema';

	export let data: PageData;
	$: ({ Shop } = data);

	const shopData = $Shop;

	$: console.log($Shop?.data?.shop);
	// Check if the data is available and map it to form fields
	const initialFormData = $shopData
		? {
				title: $shopData?.data?.shop?.title || '',
				email: $shopData?.data?.shop?.contactEmail || '',
				description: $shopData?.shop?.shop?.about || '',
				phone: $shopData?.data?.shop?.contactPhone || ''
			}
		: {
				title: '',
				email: '',
				description: '',
				phone: ''
			};

	$: console.log(initialFormData);
	$: console.log($shopData);
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
</script> -->
<script lang="ts">
	import { type Infer, type SuperValidated, superForm } from 'sveltekit-superforms';
	import { zodClient } from 'sveltekit-superforms/adapters';
	import * as Card from '$lib/components/ui/card';
	import * as Form from '$lib/components/ui/form';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea';
	import { writable } from 'svelte/store';
	import type { PageData } from './$types';
	import { generalSettingsFormSchema } from './schema';

	export let data: PageData;
	$: ({ Shop } = data);

	// Define initial form data as empty strings
	const formData = writable({
		title: '',
		email: '',
		description: '',
		phone: ''
	});

	// Update `formData` reactively when `Shop` data becomes available
	$: if ($Shop?.data?.shop) {
		const shopData = $Shop?.data?.shop;
		formData.set({
			title: shopData.title || '',
			email: shopData.contactEmail || '',
			description: shopData.about || '',
			phone: shopData.contactPhone || ''
		});
	}

	// Initialize superForm with formData
	const form = superForm(formData, {
		validators: zodClient(generalSettingsFormSchema),
		dataType: 'json'
	});

	const { form: formFields, enhance } = form;
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
