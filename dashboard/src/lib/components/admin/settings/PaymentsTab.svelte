<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { toast } from 'svelte-sonner';
	import type { Shop } from '$lib/types';
	import { api } from '$lib/api';
	import { getContext } from 'svelte';
	
	export let shop: Partial<Shop>;
	
	const authFetch: (input: RequestInfo | URL, init?: RequestInit | undefined) => Promise<Response> = getContext('authFetch');
	const refetchShopData: () => Promise<void> = getContext('refetchShopData');

	// Payment method management
	interface PaymentMethod {
		id: string;
		name: string;
		description: string;
		enabled: boolean;
		settings: PaymentMethodSetting[];
	}

	interface PaymentMethodSetting {
		id: string;
		label: string;
		type: 'text' | 'password' | 'checkbox';
		value: string | boolean;
		placeholder?: string;
		required: boolean;
	}
	
	// Example payment methods with their settings
	let paymentMethods: PaymentMethod[] = [
		{
			id: 'cash-on-delivery',
			name: 'Cash on Delivery',
			description: 'Accept cash payments upon delivery',
			enabled: false,
			settings: [
				{
					id: 'extra-fee',
					label: 'Extra fee',
					type: 'text',
					value: '0',
					placeholder: 'Additional fee for COD (e.g. 5)',
					required: false
				}
			]
		},
		{
			id: 'stripe',
			name: 'Stripe',
			description: 'Accept credit card payments via Stripe',
			enabled: false,
			settings: [
				{
					id: 'publishable-key',
					label: 'Publishable Key',
					type: 'text',
					value: '',
					placeholder: 'pk_test_...',
					required: true
				},
				{
					id: 'secret-key',
					label: 'Secret Key',
					type: 'password',
					value: '',
					placeholder: 'sk_test_...',
					required: true
				},
				{
					id: 'test-mode',
					label: 'Test Mode',
					type: 'checkbox',
					value: true,
					required: false
				}
			]
		},
		{
			id: 'paypal',
			name: 'PayPal',
			description: 'Accept payments via PayPal',
			enabled: false,
			settings: [
				{
					id: 'client-id',
					label: 'Client ID',
					type: 'text',
					value: '',
					placeholder: 'Your PayPal client ID',
					required: true
				},
				{
					id: 'client-secret',
					label: 'Client Secret',
					type: 'password',
					value: '',
					placeholder: 'Your PayPal client secret',
					required: true
				},
				{
					id: 'sandbox-mode',
					label: 'Sandbox Mode',
					type: 'checkbox',
					value: true,
					required: false
				}
			]
		},
		{
			id: 'flutterwave',
			name: 'Flutterwave',
			description: 'Accept payments via Flutterwave',
			enabled: false,
			settings: [
				{
					id: 'public-key',
					label: 'Public Key',
					type: 'text',
					value: '',
					placeholder: 'Your Flutterwave public key',
					required: true
				},
				{
					id: 'secret-key',
					label: 'Secret Key',
					type: 'password',
					value: '',
					placeholder: 'Your Flutterwave secret key',
					required: true
				},
				{
					id: 'encryption-key',
					label: 'Encryption Key',
					type: 'password',
					value: '',
					placeholder: 'Your Flutterwave encryption key',
					required: true
				},
				{
					id: 'test-mode',
					label: 'Test Mode',
					type: 'checkbox',
					value: true,
					required: false
				}
			]
		},
		{
			id: 'paystack',
			name: 'Paystack',
			description: 'Accept payments via Paystack',
			enabled: false,
			settings: [
				{
					id: 'public-key',
					label: 'Public Key',
					type: 'text',
					value: '',
					placeholder: 'Your Paystack public key',
					required: true
				},
				{
					id: 'secret-key',
					label: 'Secret Key',
					type: 'password',
					value: '',
					placeholder: 'Your Paystack secret key',
					required: true
				},
				{
					id: 'test-mode',
					label: 'Test Mode',
					type: 'checkbox',
					value: true,
					required: false
				}
			]
		}
	];

	function togglePaymentMethod(method: PaymentMethod) {
		method.enabled = !method.enabled;
		
		if (method.enabled) {
			// Only show a toast if there are missing required fields, but don't block enabling
			const missingSettings = method.settings
				.filter(setting => setting.required && setting.type !== 'checkbox' && !setting.value);
			
			if (missingSettings.length > 0) {
				toast.warning(`Please fill out the required settings for ${method.name}`);
			} else {
				toast.success(`${method.name} has been enabled`);
			}
		} else {
			toast.success(`${method.name} has been disabled`);
		}
	}

	// Serialize payment methods for API
	function serializePaymentMethods() {
		return paymentMethods.map(method => ({
			id: method.id,
			enabled: method.enabled,
			settings: method.settings.map(setting => ({
				id: setting.id,
				value: setting.value
			}))
		}));
	}

	async function savePaymentMethodSettings(method: PaymentMethod) {
		// Validate required fields
		const missingSettings = method.settings
			.filter(setting => setting.required && setting.type !== 'checkbox' && !setting.value);
		
		if (missingSettings.length > 0) {
			toast.error('Please fill out all required fields');
			return;
		}
		
		try {
			// In a real implementation, we should store payment methods in the shop object
			// Here we're saving the payment settings to the shop with a dedicated property
			// This would need to be properly implemented on the backend
			const updateData: Partial<Shop> & { payment_methods?: any } = {
				payment_methods: serializePaymentMethods()
			};
			
			// Call the API to update the shop data
			await api(authFetch).updateShop(updateData);
			
			// Refetch shop data to update the UI
			await refetchShopData();
			
			toast.success(`${method.name} settings saved successfully`);
		} catch (error) {
			console.error('Error updating payment settings:', error);
			toast.error(`Failed to update ${method.name} settings`);
		}
	}
	
	// Handle form submission
	async function handleFormSubmit(event: Event) {
		event.preventDefault();
		const form = event.target as HTMLFormElement;
		const methodIdInput = form.querySelector('input[name="method-id"]') as HTMLInputElement;
		const methodId = methodIdInput?.value;
		
		// Find the payment method
		const method = paymentMethods.find(m => m.id === methodId);
		if (!method) {
			toast.error('Payment method not found');
			return;
		}
		
		// Validate required fields
		const missingSettings = method.settings
			.filter(setting => setting.required && setting.type !== 'checkbox' && !setting.value);
		
		if (missingSettings.length > 0) {
			toast.error('Please fill out all required fields');
			return;
		}
		
		try {
			// In a real implementation, we should store payment methods in the shop object
			// Here we're saving the payment settings to the shop with a dedicated property
			const updateData: Partial<Shop> & { payment_methods?: any } = {
				payment_methods: serializePaymentMethods()
			};
			
			// Call the API to update the shop data
			await api(authFetch).updateShop(updateData);
			
			// Refetch shop data to update the UI
			await refetchShopData();
			
			toast.success(`${method.name} settings saved successfully`);
		} catch (error) {
			console.error('Error updating payment settings:', error);
			toast.error(`Failed to update ${method.name} settings`);
		}
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Payment Methods</Card.Title>
		<Card.Description>Configure the payment options your customers can use at checkout</Card.Description>
	</Card.Header>
	<Card.Content>
		<div class="space-y-6">
			{#each paymentMethods as method}
				<div class="border rounded-md overflow-hidden">
					<div class="flex items-center justify-between p-4 bg-muted/50">
						<div class="flex items-center gap-3">
							<div class="w-10 h-10 flex items-center justify-center bg-background rounded-md border">
								<span class="font-semibold">{method.name.charAt(0)}</span>
							</div>
							<div>
								<h3 class="font-medium">{method.name}</h3>
								<p class="text-sm text-muted-foreground">{method.description}</p>
							</div>
						</div>
						<div class="flex items-center space-x-2">
							<Checkbox 
								id={`enable-${method.id}`} 
								checked={method.enabled}
								onCheckedChange={() => togglePaymentMethod(method)}
							/>
							<Label
								for={`enable-${method.id}`}
								class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
							>
								{method.enabled ? 'Enabled' : 'Disabled'}
							</Label>
						</div>
					</div>
					
					<div class="p-4 space-y-4">
						<form method="POST" class="space-y-4" id={`payment-method-${method.id}-form`} on:submit={handleFormSubmit}>
							<input type="hidden" name="form-id" value="payment-method-form" />
							<input type="hidden" name="method-id" value={method.id} />
							
							{#each method.settings as setting}
								<div class="flex w-full flex-col gap-2">
									{#if setting.type === 'checkbox'}
										<div class="flex items-center space-x-2">
											<Checkbox 
												id={`${method.id}-${setting.id}`} 
												name={setting.id}
												checked={setting.value === true}
												onCheckedChange={(checked) => setting.value = checked}
											/>
											<Label
												for={`${method.id}-${setting.id}`}
												class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
											>
												{setting.label}
											</Label>
										</div>
									{:else}
										<Label for={`${method.id}-${setting.id}`}>{setting.label}</Label>
										<Input 
											id={`${method.id}-${setting.id}`} 
											name={setting.id}
											type={setting.type}
											bind:value={setting.value}
											placeholder={setting.placeholder || ''}
											required={setting.required}
										/>
									{/if}
								</div>
							{/each}
							
							<Button 
								type="submit"
								class="mt-2"
							>
								Save {method.name} Settings
							</Button>
						</form>
					</div>
				</div>
			{/each}
		</div>
	</Card.Content>
</Card.Root> 