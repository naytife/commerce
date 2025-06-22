<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { MapPin, ArrowLeft, Plus, Edit, Trash2 } from 'lucide-svelte';
	
	// Mock addresses data
	interface Address {
		id: number;
		name: string;
		fullName: string;
		address1: string;
		address2?: string;
		city: string;
		state: string;
		postalCode: string;
		country: string;
		phone: string;
		isDefault: boolean;
	}

	let addresses: Address[] = [
		{
			id: 1,
			name: 'Home',
			isDefault: true,
			fullName: 'John Doe',
			address1: '123 Main Street',
			address2: 'Apt 4B',
			city: 'New York',
			state: 'NY',
			postalCode: '10001',
			country: 'United States',
			phone: '+1 (555) 123-4567'
		},
		{
			id: 2,
			name: 'Office',
			isDefault: false,
			fullName: 'John Doe',
			address1: '456 Business Ave',
			address2: 'Suite 200',
			city: 'New York',
			state: 'NY',
			postalCode: '10002',
			country: 'United States',
			phone: '+1 (555) 987-6543'
		}
	];

	let showAddForm = false;
	let editingAddress: Address | null = null;

	function addNewAddress() {
		showAddForm = true;
		editingAddress = null;
	}

	function editAddress(address: Address) {
		editingAddress = { ...address };
		showAddForm = true;
	}

	function deleteAddress(id: number) {
		addresses = addresses.filter(addr => addr.id !== id);
	}

	function setDefaultAddress(id: number) {
		addresses = addresses.map(addr => ({
			...addr,
			isDefault: addr.id === id
		}));
	}

	function saveAddress(event: Event) {
		event.preventDefault();
		// In a real app, this would save to an API
		showAddForm = false;
		editingAddress = null;
	}

	function cancelForm() {
		showAddForm = false;
		editingAddress = null;
	}
</script>

<svelte:head>
	<title>My Addresses</title>
	<meta name="description" content="Manage your delivery addresses" />
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-4xl">
	<div class="mb-8">
		<Button variant="ghost" class="mb-4">
			<a href="/account">
				<ArrowLeft class="h-4 w-4 mr-2" />
				Back to Account
			</a>
		</Button>
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">My Addresses</h1>
				<p class="text-gray-600 dark:text-gray-400">Manage your delivery addresses</p>
			</div>
			<Button on:click={addNewAddress}>
				<Plus class="h-4 w-4 mr-2" />
				Add Address
			</Button>
		</div>
	</div>

	<div class="space-y-6">
		<!-- Address List -->
		{#if !showAddForm}
			{#each addresses as address}
				<Card.Root>
					<Card.Header>
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-2">
								<Card.Title class="text-lg">{address.name}</Card.Title>
								{#if address.isDefault}
									<span class="inline-flex items-center rounded-full bg-primary px-2.5 py-0.5 text-xs font-medium text-primary-foreground">Default</span>
								{/if}
							</div>
							<div class="flex items-center gap-2">
								<Button variant="ghost" size="icon" on:click={() => editAddress(address)}>
									<Edit class="h-4 w-4" />
								</Button>
								<Button variant="ghost" size="icon" on:click={() => deleteAddress(address.id)}>
									<Trash2 class="h-4 w-4" />
								</Button>
							</div>
						</div>
					</Card.Header>
					<Card.Content>
						<div class="space-y-2">
							<p class="font-medium">{address.fullName}</p>
							<p class="text-gray-600 dark:text-gray-400">{address.address1}</p>
							{#if address.address2}
								<p class="text-gray-600 dark:text-gray-400">{address.address2}</p>
							{/if}
							<p class="text-gray-600 dark:text-gray-400">
								{address.city}, {address.state} {address.postalCode}
							</p>
							<p class="text-gray-600 dark:text-gray-400">{address.country}</p>
							<p class="text-gray-600 dark:text-gray-400">{address.phone}</p>
						</div>
						{#if !address.isDefault}
							<div class="mt-4">
								<Button variant="outline" size="sm" on:click={() => setDefaultAddress(address.id)}>
									Set as Default
								</Button>
							</div>
						{/if}
					</Card.Content>
				</Card.Root>
			{/each}

			{#if addresses.length === 0}
				<Card.Root>
					<Card.Content class="text-center py-12">
						<MapPin class="h-12 w-12 mx-auto text-gray-400 mb-4" />
						<h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">No addresses saved</h3>
						<p class="text-gray-600 dark:text-gray-400 mb-4">
							Add a delivery address to make checkout faster
						</p>
						<Button on:click={addNewAddress}>
							<Plus class="h-4 w-4 mr-2" />
							Add Your First Address
						</Button>
					</Card.Content>
				</Card.Root>
			{/if}
		{/if}

		<!-- Add/Edit Address Form -->
		{#if showAddForm}
			<Card.Root>
				<Card.Header>
					<Card.Title>
						{editingAddress ? 'Edit Address' : 'Add New Address'}
					</Card.Title>
				</Card.Header>
				<Card.Content>
					<form on:submit={saveAddress} class="space-y-4">
						<div>
							<label for="addressName" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Address Name</label>
							<input
								id="addressName"
								class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-3 py-2 text-gray-900 dark:text-gray-100 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
								placeholder="e.g., Home, Work, etc."
								value={editingAddress?.name || ''}
								required
							/>
						</div>

						<div>
							<label for="fullName" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Full Name</label>
							<input
								id="fullName"
								class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-3 py-2 text-gray-900 dark:text-gray-100 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
								placeholder="Enter recipient's full name"
								value={editingAddress?.fullName || ''}
								required
							/>
						</div>

						<div>
							<label for="address1" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Address Line 1</label>
							<input
								id="address1"
								class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-3 py-2 text-gray-900 dark:text-gray-100 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
								placeholder="Street address"
								value={editingAddress?.address1 || ''}
								required
							/>
						</div>

						<div>
							<label for="address2" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Address Line 2 (Optional)</label>
							<input
								id="address2"
								class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-3 py-2 text-gray-900 dark:text-gray-100 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
								placeholder="Apartment, suite, unit, etc."
								value={editingAddress?.address2 || ''}
							/>
						</div>

						<div class="grid gap-4 md:grid-cols-3">
							<div>
								<label for="city" class="block text-sm font-medium text-gray-700 dark:text-gray-300">City</label>
								<input
									id="city"
									class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-3 py-2 text-gray-900 dark:text-gray-100 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
									placeholder="City"
									value={editingAddress?.city || ''}
									required
								/>
							</div>
							<div>
								<label for="state" class="block text-sm font-medium text-gray-700 dark:text-gray-300">State/Province</label>
								<input
									id="state"
									class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-3 py-2 text-gray-900 dark:text-gray-100 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
									placeholder="State"
									value={editingAddress?.state || ''}
									required
								/>
							</div>
							<div>
								<label for="postalCode" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Postal Code</label>
								<input
									id="postalCode"
									class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-3 py-2 text-gray-900 dark:text-gray-100 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
									placeholder="ZIP/Postal code"
									value={editingAddress?.postalCode || ''}
									required
								/>
							</div>
						</div>

						<div>
							<label for="country" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Country</label>
							<input
								id="country"
								class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-3 py-2 text-gray-900 dark:text-gray-100 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
								placeholder="Country"
								value={editingAddress?.country || ''}
								required
							/>
						</div>

						<div>
							<label for="phone" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Phone Number</label>
							<input
								id="phone"
								type="tel"
								class="mt-1 block w-full rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-3 py-2 text-gray-900 dark:text-gray-100 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary"
								placeholder="Phone number"
								value={editingAddress?.phone || ''}
								required
							/>
						</div>

						<div class="flex gap-2 pt-4">
							<Button type="submit">
								{editingAddress ? 'Update Address' : 'Save Address'}
							</Button>
							<Button type="button" variant="outline" on:click={cancelForm}>
								Cancel
							</Button>
						</div>
					</form>
				</Card.Content>
			</Card.Root>
		{/if}
	</div>
</div>
