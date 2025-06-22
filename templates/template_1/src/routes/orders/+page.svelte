<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Package, Truck, CheckCircle, Clock, ArrowLeft } from 'lucide-svelte';
	
	// Mock order data
	const orders = [
		{
			id: 'ORD-001',
			date: '2024-01-15',
			status: 'delivered',
			total: 89.99,
			items: [
				{ name: 'Wireless Headphones', quantity: 1, price: 79.99 },
				{ name: 'Phone Case', quantity: 1, price: 10.00 }
			]
		},
		{
			id: 'ORD-002',
			date: '2024-01-20',
			status: 'shipped',
			total: 129.50,
			items: [
				{ name: 'Smart Watch', quantity: 1, price: 129.50 }
			]
		},
		{
			id: 'ORD-003',
			date: '2024-01-22',
			status: 'processing',
			total: 45.00,
			items: [
				{ name: 'Coffee Mug', quantity: 2, price: 22.50 }
			]
		}
	];

	function getStatusIcon(status: string) {
		switch (status) {
			case 'delivered':
				return CheckCircle;
			case 'shipped':
				return Truck;
			case 'processing':
				return Clock;
			default:
				return Package;
		}
	}

	function getStatusColor(status: string) {
		switch (status) {
			case 'delivered':
				return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300';
			case 'shipped':
				return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300';
			case 'processing':
				return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300';
			default:
				return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300';
		}
	}
</script>

<svelte:head>
	<title>My Orders</title>
	<meta name="description" content="View and track your order history" />
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-6xl">
	<div class="mb-8">
		<Button variant="ghost" class="mb-4">
			<a href="/account">
				<ArrowLeft class="h-4 w-4 mr-2" />
				Back to Account
			</a>
		</Button>
		<h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">My Orders</h1>
		<p class="text-gray-600 dark:text-gray-400">Track and manage your order history</p>
	</div>

	<div class="space-y-6">
		{#each orders as order}
			<Card.Root>
				<Card.Header>
					<div class="flex items-center justify-between">
						<div>
							<Card.Title class="text-lg">Order {order.id}</Card.Title>
							<p class="text-sm text-gray-600 dark:text-gray-400">
								Placed on {new Date(order.date).toLocaleDateString()}
							</p>
						</div>
						<div class="flex items-center gap-4">
							<Badge class={getStatusColor(order.status)}>
								<svelte:component this={getStatusIcon(order.status)} class="h-3 w-3 mr-1" />
								{order.status.charAt(0).toUpperCase() + order.status.slice(1)}
							</Badge>
							<div class="text-right">
								<p class="font-semibold">${order.total.toFixed(2)}</p>
							</div>
						</div>
					</div>
				</Card.Header>
				<Card.Content>
					<div class="space-y-3">
						{#each order.items as item}
							<div class="flex items-center justify-between py-2 border-b border-gray-200 dark:border-gray-700 last:border-b-0">
								<div>
									<h4 class="font-medium text-gray-900 dark:text-white">{item.name}</h4>
									<p class="text-sm text-gray-600 dark:text-gray-400">Quantity: {item.quantity}</p>
								</div>
								<p class="font-medium">${item.price.toFixed(2)}</p>
							</div>
						{/each}
					</div>
					<div class="mt-4 flex gap-2">
						<Button variant="outline" size="sm">View Details</Button>
						{#if order.status === 'delivered'}
							<Button variant="outline" size="sm">Reorder</Button>
						{:else if order.status === 'shipped'}
							<Button variant="outline" size="sm">Track Package</Button>
						{/if}
					</div>
				</Card.Content>
			</Card.Root>
		{/each}

		{#if orders.length === 0}
			<Card.Root>
				<Card.Content class="text-center py-12">
					<Package class="h-12 w-12 mx-auto text-gray-400 mb-4" />
					<h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">No orders yet</h3>
					<p class="text-gray-600 dark:text-gray-400 mb-4">When you place orders, they'll appear here</p>
					<Button>
						<a href="/">Start Shopping</a>
					</Button>
				</Card.Content>
			</Card.Root>
		{/if}
	</div>
</div>
