<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import * as Form from '$lib/components/ui/form';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import type { PageData } from './$types';
	import { api } from '$lib/api';
	import type { Shop } from '$lib/types';
	import { getContext } from 'svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { toast } from 'svelte-sonner';
	import { page } from '$app/stores';
	export let data: PageData;
	const authFetch: (input: RequestInfo | URL, init?: RequestInit | undefined) => Promise<Response> =
		getContext('authFetch');

	$: shopQuery = createQuery<Shop, Error>({
		queryKey: [`shop-${$page.params.shop}`],
		queryFn: () => api(authFetch).getShop(),
		enabled: !!$page.params.shop
	});

	// Initialize shop data from query
	let shop: Partial<Shop> = {};
	$: if ($shopQuery.data) {
		shop = { ...$shopQuery.data };
	}

	// Payment method form
	let addPaymentMethodOpen = false;
	let cardNumber = '';
	let cardName = '';
	let expiryDate = '';
	let cvc = '';

	function handleAddPaymentMethod() {
		// Validate form
		if (!cardNumber || !cardName || !expiryDate || !cvc) {
			toast.error('Please fill in all payment details');
			return;
		}

		// In a real app, this would call a payment processing API
		toast.success('Payment method added successfully');
		addPaymentMethodOpen = false;
		
		// Reset form
		cardNumber = '';
		cardName = '';
		expiryDate = '';
		cvc = '';
	}
</script>

<div class="space-y-6">
	<h2 class="text-2xl font-bold">Billing Settings</h2>
	
	<Card.Root>
		<Card.Header>
			<Card.Title>Current Plan</Card.Title>
			<Card.Description>Your current subscription plan and features</Card.Description>
		</Card.Header>
		<Card.Content>
			<div class="p-4 border rounded-md">
				<div class="flex justify-between items-center">
					<div>
						<h3 class="font-medium">Free Plan</h3>
						<p class="text-sm text-muted-foreground">Basic features for small stores</p>
					</div>
					<span class="text-sm font-semibold bg-primary/10 text-primary px-2 py-1 rounded">Active</span>
				</div>
			</div>
		</Card.Content>
	</Card.Root>
	
	<Card.Root>
		<Card.Header>
			<Card.Title>Available Plans</Card.Title>
			<Card.Description>Upgrade your plan to unlock more features</Card.Description>
		</Card.Header>
		<Card.Content>
			<div class="grid gap-4 md:grid-cols-2">
				<div class="p-4 border rounded-md hover:border-primary cursor-pointer">
					<div class="flex justify-between items-start">
						<div>
							<h4 class="font-medium">Standard Plan</h4>
							<p class="text-sm text-muted-foreground mt-1">For growing businesses</p>
							<ul class="text-sm mt-2 space-y-1">
								<li>• Unlimited products</li>
								<li>• Advanced analytics</li>
								<li>• Custom domain</li>
								<li>• 24/7 support</li>
							</ul>
						</div>
						<div class="text-right">
							<span class="text-xl font-bold">$19</span>
							<span class="text-sm text-muted-foreground">/month</span>
						</div>
					</div>
					<Button class="mt-4 w-full" variant="outline">Select Plan</Button>
				</div>
				
				<div class="p-4 border rounded-md hover:border-primary cursor-pointer">
					<div class="flex justify-between items-start">
						<div>
							<h4 class="font-medium">Premium Plan</h4>
							<p class="text-sm text-muted-foreground mt-1">For established businesses</p>
							<ul class="text-sm mt-2 space-y-1">
								<li>• Everything in Standard</li>
								<li>• Abandoned cart recovery</li>
								<li>• Advanced SEO tools</li>
								<li>• Priority support</li>
							</ul>
						</div>
						<div class="text-right">
							<span class="text-xl font-bold">$49</span>
							<span class="text-sm text-muted-foreground">/month</span>
						</div>
					</div>
					<Button class="mt-4 w-full" variant="outline">Select Plan</Button>
				</div>
			</div>
		</Card.Content>
	</Card.Root>
	
	<Card.Root>
		<Card.Header>
			<Card.Title>Payment Methods</Card.Title>
			<Card.Description>Manage your payment information</Card.Description>
		</Card.Header>
		<Card.Content>
			<div class="p-4 border rounded-md text-center">
				<p class="text-muted-foreground">No payment methods added yet</p>
				<Button 
					variant="outline" 
					class="mt-2" 
					on:click={() => addPaymentMethodOpen = true}
				>
					+ Add Payment Method
				</Button>
			</div>
		</Card.Content>
	</Card.Root>
	
	<Card.Root>
		<Card.Header>
			<Card.Title>Billing History</Card.Title>
			<Card.Description>Your recent invoices and payments</Card.Description>
		</Card.Header>
		<Card.Content>
			<div class="border rounded-md overflow-hidden">
				<table class="w-full text-sm">
					<thead class="bg-muted">
						<tr>
							<th class="text-left p-2">Date</th>
							<th class="text-left p-2">Description</th>
							<th class="text-right p-2">Amount</th>
							<th class="text-right p-2">Status</th>
						</tr>
					</thead>
					<tbody>
						<tr class="border-t">
							<td colspan="4" class="p-4 text-center text-muted-foreground">
								No billing history available
							</td>
						</tr>
					</tbody>
				</table>
			</div>
		</Card.Content>
	</Card.Root>
</div>

<!-- Add Payment Method Dialog -->
<Dialog.Root bind:open={addPaymentMethodOpen}>
	<Dialog.Content class="max-w-md">
		<Dialog.Header>
			<Dialog.Title>Add Payment Method</Dialog.Title>
			<Dialog.Description>
				Enter your card details to add a new payment method
			</Dialog.Description>
		</Dialog.Header>
		<div class="grid gap-4 py-4">
			<div class="grid gap-2">
				<Label for="card-number">Card Number</Label>
				<Input 
					id="card-number" 
					placeholder="1234 5678 9012 3456" 
					bind:value={cardNumber}
					on:input={(e) => {
						// Format card number with spaces after every 4 digits
						e.currentTarget.value = e.currentTarget.value
							.replace(/\s/g, '')
							.replace(/(\d{4})/g, '$1 ')
							.trim();
						cardNumber = e.currentTarget.value;
					}}
				/>
			</div>
			<div class="grid gap-2">
				<Label for="card-name">Name on Card</Label>
				<Input 
					id="card-name" 
					placeholder="John Doe" 
					bind:value={cardName}
				/>
			</div>
			<div class="grid grid-cols-2 gap-4">
				<div class="grid gap-2">
					<Label for="expiry">Expiry Date</Label>
					<Input 
						id="expiry" 
						placeholder="MM/YY" 
						bind:value={expiryDate}
						on:input={(e) => {
							// Format expiry date as MM/YY
							let val = e.currentTarget.value.replace(/\D/g, '');
							if (val.length > 2) {
								val = val.substring(0, 2) + '/' + val.substring(2, 4);
							}
							expiryDate = val;
							e.currentTarget.value = val;
						}}
					/>
				</div>
				<div class="grid gap-2">
					<Label for="cvc">CVC</Label>
					<Input 
						id="cvc" 
						placeholder="123" 
						bind:value={cvc}
						type="password"
					/>
				</div>
			</div>
		</div>
		<Dialog.Footer>
			<Button variant="outline" on:click={() => addPaymentMethodOpen = false}>
				Cancel
			</Button>
			<Button on:click={handleAddPaymentMethod}>
				Add Payment Method
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root> 