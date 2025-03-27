<script lang="ts">
	import * as Form from '$lib/components/ui/form';
	import { Input } from '$lib/components/ui/input';
	import { formSchema, type FormSchema } from './schema';
	import { superForm, type SuperValidated, type Infer } from 'sveltekit-superforms';
	import { zodClient } from 'sveltekit-superforms/adapters';
	import { graphql } from '$houdini';

	const CreateCategoryMutation = graphql(`
		mutation AddCategory($parentID: ID!, $title: String!, $description: String) {
			createCategory(category: { parentID: $parentID, title: $title, description: $description }) {
				...category
			}
		}
	`);

	export let data: SuperValidated<Infer<FormSchema>>;

	const form = superForm(data, {
		validators: zodClient(formSchema),
		dataType: 'json'
	});

	const { form: formData, enhance } = form;

	let errorMessage = '';
	let successMessage = '';

	$: if (form?.form?.error) {
		errorMessage = form.form.error;
	} else if (form?.form?.success) {
		successMessage = 'Category created successfully!';
	}

	async function handleSubmitWithValidation(event: SubmitEvent) {
		event.preventDefault();

		if (!form.form.valid) {
			return; // Prevent submission if form is invalid
		}

		try {
			const response = await CreateCategoryMutation.mutate({
				title: $formData.name,
				description: $formData.description,
				parentID: 'Q2F0ZWdvcnk6Nw=='
			});

			if (response.errors) {
				errorMessage = 'An unexpected error occurred.';
			} else {
				successMessage = 'Category created successfully!';
			}
		} catch (error) {
			errorMessage = 'An unexpected error occurred.';
		}
	}
	async function handleSubmitOnSubmit(event: SubmitEvent) {
		event.preventDefault();

		try {
			const response = await CreateCategoryMutation.mutate({
				title: $formData.name,
				description: $formData.description,
				parentID: 'Q2F0ZWdvcnk6Nw=='
			});
			if (response.errors) {
				errorMessage = 'An unexpected error occurred.';
			} else {
				successMessage = 'Category created successfully!';
			}
		} catch (error) {
			errorMessage = 'An unexpected error occurred.';
		}
	}
</script>

<form method="post" use:enhance on:submit={handleSubmitOnSubmit}>
	<div class="grid gap-6">
		<!-- Name Field -->
		<div class="grid gap-3">
			<Form.Field {form} name="name">
				<Form.Control let:attrs>
					<Form.Label for="name" aria-required="true">Name</Form.Label>
					<Input type="text" class="w-full" bind:value={$formData.name} {...attrs} />
				</Form.Control>
				<Form.FieldErrors />
			</Form.Field>
		</div>

		<!-- Description Field -->
		<div class="grid gap-3">
			<Form.Field {form} name="description">
				<Form.Control let:attrs>
					<Form.Label for="description" aria-required="true">Description</Form.Label>
					<Input type="text" class="w-full" bind:value={$formData.description} {...attrs} />
				</Form.Control>
				<Form.FieldErrors />
			</Form.Field>
		</div>

		<!-- Submit Button -->
		<div class="hidden w-full items-center gap-2 md:ml-auto md:flex">
			<Form.Button type="submit" size="sm">Save Category</Form.Button>
		</div>

		<!-- Error and Success Messages -->
		{#if errorMessage}
			<div class="text-red-300">{errorMessage}</div>
		{/if}
		{#if successMessage}
			<div class="text-green-300">{successMessage}</div>
		{/if}
	</div>
</form>
