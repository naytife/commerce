<script lang="ts">
	import { ArrowUpDown, MoreHorizontal, ShoppingCart, Filter, Search, Download, Plus, Eye, Trash2, Edit3 } from 'lucide-svelte';
	import * as Table from '$lib/components/ui/table';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Input } from '$lib/components/ui/input';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { getContext } from 'svelte';
	import { api } from '$lib/api';
	import type { Order, Shop, PaginatedResponse } from '$lib/types';
	import type { QueryObserverResult } from '@tanstack/svelte-query';
	import { format } from 'date-fns';
	import { toast } from 'svelte-sonner';
	import ListFilter from 'lucide-svelte/icons/list-filter';
	import { page as pageStore } from '$app/stores';
	import { onMount } from 'svelte';
	import { getCurrencySymbol, formatCurrencyWithLocale } from '$lib/utils/currency';
	import { get, writable } from 'svelte/store';

	const authFetch = getContext<typeof fetch>('authFetch');
	const queryClient = useQueryClient();

	// Access route parameters through the page store
	const shopParam = get(pageStore).params.shop;

	// Pagination state as Svelte stores
	const currentPage = writable(1);
	const limit = writable(10);
	let total = 0;
	let totalPages = 1;

	// Orders query with pagination (reactive)
	let ordersQuery;
	$: ordersQuery = createQuery<PaginatedResponse<Order>, Error>({
	    queryKey: [`shop-${shopParam}-orders`, $currentPage, $limit],
	    queryFn: () => api(authFetch).getOrders({ page: $currentPage, limit: $limit })
	});

	$: ordersData = $ordersQuery?.data?.data ?? [];
	$: pagination = $ordersQuery?.data?.pagination ?? { page: 1, limit: 10, total: 0, total_pages: 1 };
	$: total = pagination.total;
	$: totalPages = pagination.total_pages;

	// Query shop to get currency code and define a currency formatter
	const shopQuery = createQuery<Shop, Error>({
		queryKey: [`shop-${shopParam}`],
		queryFn: () => api(authFetch).getShop(),
	});
	$: currencyCode = $shopQuery.data?.currency_code || 'USD';
	$: currencySymbol = getCurrencySymbol(currencyCode);
	$: formatCurrency = (amount: number) => formatCurrencyWithLocale(amount, currencyCode);

	let statusFilter: Record<string, boolean> = { pending: true, processing: true, completed: true, cancelled: true, refunded: true };
	let searchQuery = '';
	
	// Filtering and search (client-side on current page)
	$: filteredOrders = ordersData.filter((order: Order) => {
	    const matchesStatus = statusFilter[order.status.toLowerCase()];
	    const matchesSearch = !searchQuery || 
	        order.customer_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
	        order.customer_email.toLowerCase().includes(searchQuery.toLowerCase()) ||
	        order.order_id.toString().includes(searchQuery);
	    return matchesStatus && matchesSearch;
	});

	function goToPage(newPage: number) {
	    if (newPage >= 1 && newPage <= totalPages) {
	        currentPage.set(newPage);
	    }
	}

	async function handleDeleteOrder(orderId: number) {
		try {
			await api(authFetch).deleteOrder(orderId);
			toast.success('Order deleted successfully');
			queryClient.invalidateQueries({ queryKey: ['orders'] });
		} catch (error) {
			toast.error('Failed to delete order');
		}
	}

	async function handleUpdateStatus(orderId: number, status: string) {
		try {
			await api(authFetch).updateOrderStatus(orderId, status);
			toast.success('Order status updated successfully');
			queryClient.invalidateQueries({ queryKey: ['orders'] });
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

	function getStatusClass(status: string) {
		switch (status.toLowerCase()) {
			case 'pending':
				return 'status-warning';
			case 'processing':
				return 'status-info';
			case 'completed':
				return 'status-active';
			case 'cancelled':
				return 'status-destructive';
			case 'refunded':
				return 'status-secondary';
			default:
				return 'status-inactive';
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

	function getStatusBgClass(status: string) {
		switch (status.toLowerCase()) {
			case 'pending':
				return 'bg-warning text-white';
			case 'processing':
				return 'bg-info text-white';
			case 'completed':
				return 'bg-success text-white';
			case 'cancelled':
				return 'bg-destructive text-white';
			case 'refunded':
				return 'bg-secondary text-white';
			default:
				return 'bg-muted text-foreground';
		}
	}

	// Calculate order statistics
	$: orderStats = {
		total: filteredOrders.length,
		pending: filteredOrders.filter((o: Order) => o.status.toLowerCase() === 'pending').length,
		processing: filteredOrders.filter((o: Order) => o.status.toLowerCase() === 'processing').length,
		completed: filteredOrders.filter((o: Order) => o.status.toLowerCase() === 'completed').length,
		totalRevenue: filteredOrders
			.filter((o: Order) => o.status.toLowerCase() === 'completed')
			.reduce((sum: number, order: Order) => sum + order.amount, 0),
		averageOrder: filteredOrders.length > 0 
			? filteredOrders.reduce((sum: number, order: Order) => sum + order.amount, 0) / filteredOrders.length 
			: 0
	};

	onMount(async () => {
		await queryClient.prefetchQuery({
			queryKey: [`shop-${shopParam}-orders`],
			queryFn: () => api(authFetch).getOrders(),
		});
	});
</script>

<div class="space-y-8">
	<!-- Header Section -->
	<div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
		<div>
			<h1 class="text-3xl font-bold text-foreground mb-2">Orders</h1>
			<p class="text-muted-foreground">
				Manage customer orders and track their fulfillment status
			</p>
		</div>
		<div class="flex items-center gap-3">
			<Button class="btn-gradient shadow-brand">
				<Plus class="w-4 h-4 mr-2" />
				Create Order
			</Button>
		</div>
	</div>

	<!-- Enhanced Summary Cards -->
	<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
		<!-- Total Orders -->
		<div class="card-interactive text-center">
			<div class="w-12 h-12 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-brand">
				<ShoppingCart class="w-6 h-6 text-white" />
			</div>
			<div class="text-2xl font-bold text-foreground mb-1">{orderStats.total}</div>
			<div class="text-sm text-muted-foreground">Total Orders</div>
		</div>
		
		<!-- Pending Orders -->
		<div class="card-interactive text-center">
			<div class="w-12 h-12 bg-gradient-to-br from-warning to-orange-500 rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-brand">
				<Filter class="w-6 h-6 text-white" />
			</div>
			<div class="text-2xl font-bold text-foreground mb-1">{orderStats.pending}</div>
			<div class="text-sm text-muted-foreground">Pending Orders</div>
		</div>
		
		<!-- Completed Orders -->
		<div class="card-interactive text-center">
			<div class="w-12 h-12 bg-gradient-to-br from-success to-emerald-500 rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-brand">
				<Eye class="w-6 h-6 text-white" />
			</div>
			<div class="text-2xl font-bold text-foreground mb-1">{orderStats.completed}</div>
			<div class="text-sm text-muted-foreground">Completed Orders</div>
		</div>
		
		<!-- Total Revenue -->
		<div class="card-interactive text-center">
			<div class="w-12 h-12 bg-gradient-to-br from-secondary to-accent rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-brand">
				<ArrowUpDown class="w-6 h-6 text-white" />
			</div>
			<div class="text-2xl font-bold text-foreground mb-1">{formatCurrency(orderStats.totalRevenue)}</div>
			<div class="text-sm text-muted-foreground">Total Revenue</div>
		</div>
	</div>

	<!-- Orders Management -->
	<div class="card-elevated">
		{#if $ordersQuery?.isLoading}
			<!-- Enhanced Loading State -->
			<div class="space-y-6">
				<div class="flex flex-col lg:flex-row gap-4">
					<div class="flex-1">
						<div class="h-10 bg-muted animate-pulse rounded-xl"></div>
					</div>
					<div class="flex gap-2">
						<div class="h-10 w-24 bg-muted animate-pulse rounded-xl"></div>
						<div class="h-10 w-24 bg-muted animate-pulse rounded-xl"></div>
					</div>
				</div>

				<div class="border border-border/50 rounded-xl overflow-hidden">
					<div class="bg-surface-elevated border-b border-border/50 p-4">
						<div class="grid grid-cols-7 gap-4">
							{#each Array(7) as _}
								<div class="h-4 bg-muted animate-pulse rounded"></div>
							{/each}
						</div>
					</div>
					{#each Array(5) as _}
						<div class="border-b border-border/50 p-4">
							<div class="grid grid-cols-7 gap-4">
								{#each Array(7) as _}
									<div class="h-6 bg-muted animate-pulse rounded"></div>
								{/each}
							</div>
						</div>
					{/each}
				</div>
			</div>
		{:else if $ordersQuery?.isError}
			<!-- Enhanced Error State -->
			<div class="flex flex-col items-center justify-center py-20">
				<div class="w-16 h-16 bg-gradient-to-br from-destructive to-red-500 rounded-full flex items-center justify-center mb-6 shadow-lg">
					<ShoppingCart class="w-8 h-8 text-white" />
				</div>
				<h3 class="text-xl font-semibold mb-2 text-foreground">Failed to load orders</h3>
				<p class="text-muted-foreground mb-6 text-center max-w-md">
					{$ordersQuery.error.message || 'Unable to fetch order data. Please try again.'}
				</p>
				<Button 
					on:click={() => queryClient.invalidateQueries({ queryKey: [`shop-${shopParam}-orders`] })} 
					variant="outline" 
					class="glass border-border/50"
				>
					Try Again
				</Button>
			</div>
		{:else if ordersData.length === 0}
			<!-- Enhanced Empty State -->
			<div class="flex flex-col items-center justify-center py-20">
				<div class="w-20 h-20 bg-gradient-to-br from-primary to-accent rounded-full flex items-center justify-center mb-8 shadow-brand">
					<ShoppingCart class="w-10 h-10 text-white" />
				</div>
				<h3 class="text-2xl font-bold mb-4 text-foreground">No orders yet</h3>
				<p class="text-lg text-muted-foreground mb-8 text-center max-w-2xl leading-relaxed">
					Once customers start placing orders, they'll appear here. You can track order status, manage fulfillment, and analyze sales data.
				</p>
				<div class="flex flex-col sm:flex-row gap-4">
					<Button class="btn-gradient shadow-brand">
						<Plus class="w-5 h-5 mr-2" />
						Create Test Order
					</Button>
					<Button variant="outline" class="glass border-border/50">
						View Products
					</Button>
				</div>
			</div>
		{:else}
			<!-- Orders Management Interface -->
			<div class="space-y-6">
				<!-- Search and Filters -->
				<div class="flex flex-col lg:flex-row gap-4">
					<div class="flex-1 relative">
						<Search class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-muted-foreground" />
						<Input
							bind:value={searchQuery}
							placeholder="Search orders by customer, email, or order ID..."
							class="pl-10 h-12 glass border-border/50 focus:border-primary/50"
						/>
					</div>
					<div class="flex gap-2">
						<DropdownMenu.Root>
							<DropdownMenu.Trigger asChild let:builder>
								<Button builders={[builder]} variant="outline" class="glass border-border/50 h-12">
									<Filter class="w-4 h-4 mr-2" />
									Filter
								</Button>
							</DropdownMenu.Trigger>
							<DropdownMenu.Content align="end" class="w-48">
								<DropdownMenu.Label>Filter by Status</DropdownMenu.Label>
								<DropdownMenu.Separator />
								<DropdownMenu.CheckboxItem bind:checked={statusFilter.pending}>
									<span class="flex items-center gap-2">
										<div class="w-2 h-2 rounded-full bg-warning"></div>
										Pending
									</span>
								</DropdownMenu.CheckboxItem>
								<DropdownMenu.CheckboxItem bind:checked={statusFilter.processing}>
									<span class="flex items-center gap-2">
										<div class="w-2 h-2 rounded-full bg-info"></div>
										Processing
									</span>
								</DropdownMenu.CheckboxItem>
								<DropdownMenu.CheckboxItem bind:checked={statusFilter.completed}>
									<span class="flex items-center gap-2">
										<div class="w-2 h-2 rounded-full bg-success"></div>
										Completed
									</span>
								</DropdownMenu.CheckboxItem>
								<DropdownMenu.CheckboxItem bind:checked={statusFilter.cancelled}>
									<span class="flex items-center gap-2">
										<div class="w-2 h-2 rounded-full bg-destructive"></div>
										Cancelled
									</span>
								</DropdownMenu.CheckboxItem>
								<DropdownMenu.CheckboxItem bind:checked={statusFilter.refunded}>
									<span class="flex items-center gap-2">
										<div class="w-2 h-2 rounded-full bg-secondary"></div>
										Refunded
									</span>
								</DropdownMenu.CheckboxItem>
							</DropdownMenu.Content>
						</DropdownMenu.Root>
						
					</div>
				</div>

				<!-- Orders Table -->
				<div class="border border-border/50 rounded-xl overflow-hidden">
					<Table.Root>
						<Table.Header>
							<Table.Row class="bg-surface-elevated border-border/50 hover:bg-surface-elevated">
								<Table.Head class="font-medium">Order</Table.Head>
								<Table.Head class="font-medium">Customer</Table.Head>
								<Table.Head class="font-medium">Date</Table.Head>
								<Table.Head class="font-medium">Amount</Table.Head>
								<Table.Head class="font-medium">Status</Table.Head>
								<Table.Head class="font-medium">Payment</Table.Head>
								<Table.Head class="w-[50px]"></Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each filteredOrders as order}
								<Table.Row class="border-border/50 hover:bg-surface-elevated/50 transition-colors">
									<Table.Cell class="font-medium">
										<div class="flex items-center gap-3">
											<div class="w-10 h-10 bg-gradient-to-br from-primary to-accent rounded-lg flex items-center justify-center shadow-sm">
												<ShoppingCart class="w-5 h-5 text-white" />
											</div>
											<div>
												<div class="font-semibold text-foreground">#{order.order_id}</div>
												<div class="text-xs text-muted-foreground">Order ID</div>
											</div>
										</div>
									</Table.Cell>
									<Table.Cell>
										<div>
											<div class="font-medium text-foreground">{order.customer_name}</div>
											<div class="text-sm text-muted-foreground">{order.customer_email}</div>
										</div>
									</Table.Cell>
									<Table.Cell>
										<div class="text-sm">
											<div class="font-medium text-foreground">
												{format(new Date(order.created_at), 'MMM d, yyyy')}
											</div>
											<div class="text-muted-foreground">
												{format(new Date(order.created_at), 'h:mm a')}
											</div>
										</div>
									</Table.Cell>
									<Table.Cell>
										<div class="font-semibold text-foreground">
											{formatCurrency(order.amount)}
										</div>
									</Table.Cell>
									<Table.Cell>
										<Badge class={getStatusBgClass(order.status)}>
											{order.status}
										</Badge>
									</Table.Cell>
									<Table.Cell>
										<Badge class={order.payment_status === 'paid' ? 'status-active' : 'status-warning'}>
											{order.payment_status}
										</Badge>
									</Table.Cell>
									<Table.Cell>
										<DropdownMenu.Root>
											<DropdownMenu.Trigger asChild let:builder>
												<Button
													variant="ghost"
													size="sm"
													class="h-8 w-8 p-0 hover:bg-surface-elevated"
													builders={[builder]}
												>
													<span class="sr-only">Open menu</span>
													<MoreHorizontal class="h-4 w-4" />
												</Button>
											</DropdownMenu.Trigger>
											<DropdownMenu.Content align="end" class="w-48">
												<DropdownMenu.Label>Actions</DropdownMenu.Label>
												<DropdownMenu.Item 
													href={`/${shopParam}/orders/${order.order_id}`}
													class="flex items-center gap-2"
												>
													<Eye class="w-4 h-4" />
													View Details
												</DropdownMenu.Item>
												<DropdownMenu.Item class="flex items-center gap-2">
													<Edit3 class="w-4 h-4" />
													Edit Order
												</DropdownMenu.Item>
												<DropdownMenu.Separator />
												<DropdownMenu.Sub>
													<DropdownMenu.SubTrigger>Update Status</DropdownMenu.SubTrigger>
													<DropdownMenu.SubContent>
														<DropdownMenu.Item on:click={() => handleUpdateStatus(order.order_id, 'pending')}>
															<span class="flex items-center gap-2">
																<div class="w-2 h-2 rounded-full bg-warning"></div>
																Pending
															</span>
														</DropdownMenu.Item>
														<DropdownMenu.Item on:click={() => handleUpdateStatus(order.order_id, 'processing')}>
															<span class="flex items-center gap-2">
																<div class="w-2 h-2 rounded-full bg-info"></div>
																Processing
															</span>
														</DropdownMenu.Item>
														<DropdownMenu.Item on:click={() => handleUpdateStatus(order.order_id, 'completed')}>
															<span class="flex items-center gap-2">
																<div class="w-2 h-2 rounded-full bg-success"></div>
																Completed
															</span>
														</DropdownMenu.Item>
														<DropdownMenu.Item on:click={() => handleUpdateStatus(order.order_id, 'cancelled')}>
															<span class="flex items-center gap-2">
																<div class="w-2 h-2 rounded-full bg-destructive"></div>
																Cancelled
															</span>
														</DropdownMenu.Item>
														<DropdownMenu.Item on:click={() => handleUpdateStatus(order.order_id, 'refunded')}>
															<span class="flex items-center gap-2">
																<div class="w-2 h-2 rounded-full bg-secondary"></div>
																Refunded
															</span>
														</DropdownMenu.Item>
													</DropdownMenu.SubContent>
												</DropdownMenu.Sub>
												<DropdownMenu.Separator />
												<DropdownMenu.Item
													class="flex items-center gap-2 text-destructive"
													on:click={() => handleDeleteOrder(order.order_id)}
												>
													<Trash2 class="w-4 h-4" />
													Delete
												</DropdownMenu.Item>
											</DropdownMenu.Content>
										</DropdownMenu.Root>
									</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				</div>

				<!-- Pagination & Summary -->
				{#if total > 0}
					<div class="flex items-center justify-between pt-4 border-t border-border/50">
						<div class="text-sm text-muted-foreground">
							Showing <strong>{($currentPage - 1) * $limit + 1}</strong> to <strong>{Math.min($currentPage * $limit, total)}</strong> of <strong>{total}</strong> orders
						</div>
						<div class="flex items-center gap-2">
							<Button variant="outline" size="sm" class="glass border-border/50" on:click={() => goToPage($currentPage - 1)} disabled={$currentPage === 1}>
								Previous
							</Button>
							<span>Page {$currentPage} of {totalPages}</span>
							<Button variant="outline" size="sm" class="glass border-border/50" on:click={() => goToPage($currentPage + 1)} disabled={$currentPage === totalPages}>
								Next
							</Button>
						</div>
					</div>
				{/if}
			</div>
		{/if}
	</div>
</div>