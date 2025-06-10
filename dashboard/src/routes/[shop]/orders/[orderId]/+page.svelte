<script lang="ts">
	import { ArrowLeft, MoreHorizontal } from 'lucide-svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { getContext } from 'svelte';
	import { api } from '$lib/api';
	import type { Order, Shop } from '$lib/types';
	import { format } from 'date-fns';
	import { toast } from 'svelte-sonner';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

	const authFetch = getContext<typeof fetch>('authFetch');
	const queryClient = useQueryClient();
	const orderId = $page.params.orderId;

	const order = createQuery<Order, Error>({
		queryKey: [`shop-${$page.params.shop}-order`, orderId],
		queryFn: () => api(authFetch).getOrderById(Number(orderId))
	});

	const shopQuery = createQuery<Shop, Error>({
		queryKey: [`shop-${$page.params.shop}`],
		queryFn: () => api(authFetch).getShop(),
		enabled: !!$page.params.shop
	});
	$: currencyCode = $shopQuery.data?.currency_code || 'USD';
	$: currencySymbol = currencyCode === 'NGN' ? '₦' : '$';
	$: formatCurrency = (amount: number) => `${currencySymbol}${amount.toFixed(2)}`;

	async function handleDeleteOrder() {
		try {
			await api(authFetch).deleteOrder(Number(orderId));
			toast.success('Order deleted successfully');
			window.location.href = `/${$page.params.shop}/orders`;
		} catch (error) {
			toast.error('Failed to delete order');
		}
	}

	async function handleUpdateStatus(status: string) {
		try {
			await api(authFetch).updateOrderStatus(Number(orderId), status);
			toast.success('Order status updated successfully');
			queryClient.invalidateQueries({ queryKey: ['order', orderId] });
		} catch (error) {
			toast.error('Failed to update order status');
		}
	}

	function getStatusBadgeVariant(status: string) {
		switch (status.toLowerCase()) {
			case 'pending':
				return 'outline';
			case 'processing':
				return 'secondary';
			case 'completed':
				return 'default';
			case 'cancelled':
				return 'destructive';
			case 'refunded':
				return 'secondary';
			default:
				return 'default';
		}
	}

	function getStatusDescription(status: string) {
		switch (status.toLowerCase()) {
			case 'pending':
				return 'Order placed but not processed';
			case 'processing':
				return 'Payment confirmed, preparing for shipping';
			case 'completed':
				return 'Order delivered and closed';
			case 'cancelled':
				return 'Order cancelled before completion';
			case 'refunded':
				return 'Full refund issued';
			default:
				return '';
		}
	}

	onMount(async () => {
		await queryClient.prefetchQuery({
			queryKey: [`shop-${$page.params.shop}-order`, $page.params.orderId],
			queryFn: () => api(authFetch).getOrderById(Number($page.params.orderId)),
		});
	});
</script>

<div>
	<main class="grid flex-1 items-start gap-4 p-4 sm:px-6 sm:py-0 md:gap-8">
		<div class="flex items-center gap-4">
			<Button variant="outline" size="icon" href={`/${$page.params.shop}/orders`}>
				<ArrowLeft class="h-4 w-4" />
				<span class="sr-only">Back</span>
			</Button>
			<h1 class="text-3xl font-bold">Order Details</h1>
		</div>

		{#if $order.status === 'pending'}
			<p>Loading...</p>
		{:else if $order.status === 'error'}
			<span>Error: {$order.error.message}</span>
		{:else}
			<div class="grid gap-4 md:grid-cols-2">
				<Card.Root>
					<Card.Header>
						<div class="flex items-center justify-between">
							<Card.Title>Order #{$order.data.order_id}</Card.Title>
							<DropdownMenu.Root>
								<DropdownMenu.Trigger asChild let:builder>
									<Button
										variant="ghost"
										class="h-8 w-8 p-0"
										builders={[builder]}
									>
										<span class="sr-only">Open menu</span>
										<MoreHorizontal class="h-4 w-4" />
									</Button>
								</DropdownMenu.Trigger>
								<DropdownMenu.Content align="end">
									<DropdownMenu.Label>Actions</DropdownMenu.Label>
									<DropdownMenu.Separator />
									<DropdownMenu.Sub>
										<DropdownMenu.SubTrigger>Update Status</DropdownMenu.SubTrigger>
										<DropdownMenu.SubContent>
											<DropdownMenu.Item on:click={() => handleUpdateStatus('pending')}>
												Pending
											</DropdownMenu.Item>
											<DropdownMenu.Item on:click={() => handleUpdateStatus('processing')}>
												Processing
											</DropdownMenu.Item>
											<DropdownMenu.Item on:click={() => handleUpdateStatus('completed')}>
												Completed
											</DropdownMenu.Item>
											<DropdownMenu.Item on:click={() => handleUpdateStatus('cancelled')}>
												Cancelled
											</DropdownMenu.Item>
											<DropdownMenu.Item on:click={() => handleUpdateStatus('refunded')}>
												Refunded
											</DropdownMenu.Item>
										</DropdownMenu.SubContent>
									</DropdownMenu.Sub>
									<DropdownMenu.Separator />
									<DropdownMenu.Item
										class="text-destructive"
										on:click={handleDeleteOrder}
									>
										Delete
									</DropdownMenu.Item>
								</DropdownMenu.Content>
							</DropdownMenu.Root>
						</div>
						<Card.Description>
							Order placed on {format(new Date($order.data.created_at), 'PPP')}
						</Card.Description>
					</Card.Header>
					<Card.Content class="grid gap-4">
						<div class="grid gap-2">
							<div class="font-medium">Status</div>
							<Badge variant={getStatusBadgeVariant($order.data.status)}>
								{$order.data.status}
							</Badge>
							<div class="text-xs text-muted-foreground">
								{getStatusDescription($order.data.status)}
							</div>
						</div>
						<div class="grid gap-2">
							<div class="font-medium">Payment Status</div>
							<Badge variant={$order.data.payment_status === 'paid' ? 'default' : 'outline'}>
								{$order.data.payment_status}
							</Badge>
						</div>
						<div class="grid gap-2">
							<div class="font-medium">Payment Method</div>
							<div>{$order.data.payment_method}</div>
						</div>
						<div class="grid gap-2">
							<div class="font-medium">Transaction ID</div>
							<div>{$order.data.transaction_id}</div>
						</div>
					</Card.Content>
				</Card.Root>

				<Card.Root>
					<Card.Header>
						<Card.Title>Customer Information</Card.Title>
					</Card.Header>
					<Card.Content class="grid gap-4">
						<div class="grid gap-2">
							<div class="font-medium">Name</div>
							<div>{$order.data.customer_name}</div>
						</div>
						<div class="grid gap-2">
							<div class="font-medium">Email</div>
							<div>{$order.data.customer_email}</div>
						</div>
						<div class="grid gap-2">
							<div class="font-medium">Phone</div>
							<div>{$order.data.customer_phone}</div>
						</div>
						<div class="grid gap-2">
							<div class="font-medium">Shipping Address</div>
							<div>{$order.data.shipping_address}</div>
						</div>
					</Card.Content>
				</Card.Root>

				<Card.Root class="md:col-span-2">
					<Card.Header>
						<Card.Title>Order Items</Card.Title>
					</Card.Header>
					<Card.Content>
						<div class="grid gap-4">
							{#each $order.data.items as item}
								<div class="flex items-center justify-between">
									<div>
										<div class="font-medium">Product Variation #{item.product_variation_id}</div>
										<div class="text-sm text-muted-foreground">
											Quantity: {item.quantity}
										</div>
									</div>
									<div class="font-medium">{formatCurrency(item.price)}</div>
								</div>
							{/each}
						</div>
						<div class="mt-4 grid gap-2 border-t pt-4">
							<div class="flex items-center justify-between">
								<div>Subtotal</div>
								<div>{formatCurrency($order.data.amount - $order.data.tax - $order.data.shipping_cost)}</div>
							</div>
							<div class="flex items-center justify-between">
								<div>Tax</div>
								<div>{formatCurrency($order.data.tax)}</div>
							</div>
							<div class="flex items-center justify-between">
								<div>Shipping</div>
								<div>{formatCurrency($order.data.shipping_cost)}</div>
							</div>
							<div class="flex items-center justify-between font-medium">
								<div>Total</div>
								<div>{formatCurrency($order.data.amount)}</div>
							</div>
						</div>
					</Card.Content>
				</Card.Root>
			</div>
		{/if}
	</main>
</div> 