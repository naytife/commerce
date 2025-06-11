<script lang="ts">
	import File from 'lucide-svelte/icons/file';
	import ListFilter from 'lucide-svelte/icons/list-filter';
	import Ellipsis from 'lucide-svelte/icons/ellipsis';
	import CirclePlus from 'lucide-svelte/icons/circle-plus';
	import { 
		TrendingUp, 
		DollarSign, 
		ShoppingCart, 
		Users, 
		Eye,
		Edit,
		Trash2,
		Calendar,
		Clock,
		ArrowUpRight,
		BarChart3,
		Package
	} from 'lucide-svelte';
	import * as Table from '$lib/components/ui/table';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Tabs from '$lib/components/ui/tabs';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { getContext } from 'svelte';
	import { api } from '$lib/api';
	import type { Product } from '$lib/types';
	import { page } from '$app/stores';
	import { getCurrencySymbol, formatCurrencyWithLocale } from '$lib/utils/currency';

	const client = useQueryClient()
	const authFetch = getContext('authFetch')

	// Get shop data for currency
	const shopQuery = createQuery({
		queryKey: [`shop-${$page.params.shop}`],
		queryFn: () => api(authFetch as any).getShop($page.params.shop),
	});

	$: currencyCode = $shopQuery.data?.currency_code || 'USD';
	$: currencySymbol = getCurrencySymbol(currencyCode);

	let limit = 10

	const products = createQuery<
	Product[],
	Error
	>({
	queryKey: [`shop-${$page.params.shop}-products`, limit],
	queryFn: () => api(authFetch as any).getProducts(limit),
	})

	// Mock data for dashboard stats (replace with real API calls)
	const dashboardStats = {
		totalRevenue: 24580,
		revenueChange: '+12.5%',
		totalOrders: 1234,
		ordersChange: '+8.2%',
		totalCustomers: 567,
		customersChange: '+3.1%',
		conversionRate: '3.4%',
		conversionChange: '+0.8%'
	};

	const recentOrders = [
		{ id: '#3210', customer: 'John Doe', status: 'completed', amount: 156.00, date: '2 hours ago' },
		{ id: '#3209', customer: 'Jane Smith', status: 'processing', amount: 89.50, date: '4 hours ago' },
		{ id: '#3208', customer: 'Mike Johnson', status: 'pending', amount: 234.75, date: '6 hours ago' },
		{ id: '#3207', customer: 'Sarah Wilson', status: 'completed', amount: 67.25, date: '8 hours ago' },
	];
</script>

<div class="space-y-8">
	<!-- Dashboard Header -->
	<div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
		<div>
			<h1 class="text-3xl font-bold text-foreground mb-2">Dashboard Overview</h1>
			<p class="text-muted-foreground">
				Welcome back! Here's what's happening with your store today.
			</p>
		</div>
		<div class="flex items-center gap-3">
			<Button variant="outline" class="glass border-border/50">
				<Calendar class="w-4 h-4 mr-2" />
				Last 30 days
			</Button>
			<Button href="/{$page.params.shop}/product-types/create" class="btn-gradient shadow-brand">
				<CirclePlus class="w-4 h-4 mr-2" />
				Add Product
			</Button>
		</div>
	</div>

	<!-- Stats Cards -->
	<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
		<!-- Revenue Card -->
		<div class="card-interactive">
			<div class="flex items-center justify-between mb-4">
				<div class="w-12 h-12 bg-gradient-to-br from-success to-emerald-500 rounded-2xl flex items-center justify-center shadow-brand">
					<DollarSign class="w-6 h-6 text-white" />
				</div>
				<Badge variant="secondary" class="text-success bg-success/10 border-success/20">
					{dashboardStats.revenueChange}
				</Badge>
			</div>
			<div class="space-y-1">
				<p class="text-sm text-muted-foreground">Total Revenue</p>
				<p class="text-2xl font-bold text-foreground">{formatCurrencyWithLocale(dashboardStats.totalRevenue, currencyCode)}</p>
			</div>
		</div>

		<!-- Orders Card -->
		<div class="card-interactive">
			<div class="flex items-center justify-between mb-4">
				<div class="w-12 h-12 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center shadow-brand">
					<ShoppingCart class="w-6 h-6 text-white" />
				</div>
				<Badge variant="secondary" class="text-primary bg-primary/10 border-primary/20">
					{dashboardStats.ordersChange}
				</Badge>
			</div>
			<div class="space-y-1">
				<p class="text-sm text-muted-foreground">Total Orders</p>
				<p class="text-2xl font-bold text-foreground">{dashboardStats.totalOrders.toLocaleString()}</p>
			</div>
		</div>

		<!-- Customers Card -->
		<div class="card-interactive">
			<div class="flex items-center justify-between mb-4">
				<div class="w-12 h-12 bg-gradient-to-br from-accent to-secondary rounded-2xl flex items-center justify-center shadow-brand">
					<Users class="w-6 h-6 text-white" />
				</div>
				<Badge variant="secondary" class="text-accent bg-accent/10 border-accent/20">
					{dashboardStats.customersChange}
				</Badge>
			</div>
			<div class="space-y-1">
				<p class="text-sm text-muted-foreground">Total Customers</p>
				<p class="text-2xl font-bold text-foreground">{dashboardStats.totalCustomers.toLocaleString()}</p>
			</div>
		</div>

		<!-- Conversion Rate Card -->
		<div class="card-interactive">
			<div class="flex items-center justify-between mb-4">
				<div class="w-12 h-12 bg-gradient-to-br from-secondary to-primary rounded-2xl flex items-center justify-center shadow-brand">
					<TrendingUp class="w-6 h-6 text-white" />
				</div>
				<Badge variant="secondary" class="text-secondary bg-secondary/10 border-secondary/20">
					{dashboardStats.conversionChange}
				</Badge>
			</div>
			<div class="space-y-1">
				<p class="text-sm text-muted-foreground">Conversion Rate</p>
				<p class="text-2xl font-bold text-foreground">{dashboardStats.conversionRate}</p>
			</div>
		</div>
	</div>

	<!-- Main Content Grid -->
	<div class="grid gap-8 lg:grid-cols-3">
		<!-- Recent Products -->
		<div class="lg:col-span-2">
			<div class="card-elevated">
				<div class="flex items-center justify-between mb-6">
					<div>
						<h2 class="text-xl font-semibold text-foreground">Recent Products</h2>
						<p class="text-sm text-muted-foreground">Manage your product catalog</p>
					</div>
					<div class="flex items-center gap-2">
						<DropdownMenu.Root>
							<DropdownMenu.Trigger asChild let:builder>
								<Button builders={[builder]} variant="outline" size="sm" class="glass border-border/50">
									<ListFilter class="h-3.5 w-3.5 mr-2" />
									Filter
								</Button>
							</DropdownMenu.Trigger>
							<DropdownMenu.Content align="end" class="glass border-border/50 shadow-glass rounded-xl">
								<DropdownMenu.Label>Filter by Status</DropdownMenu.Label>
								<DropdownMenu.Separator />
								<DropdownMenu.CheckboxItem checked>Active</DropdownMenu.CheckboxItem>
								<DropdownMenu.CheckboxItem>Draft</DropdownMenu.CheckboxItem>
								<DropdownMenu.CheckboxItem>Archived</DropdownMenu.CheckboxItem>
							</DropdownMenu.Content>
						</DropdownMenu.Root>
						<Button size="sm" variant="outline" class="glass border-border/50">
							<File class="h-3.5 w-3.5 mr-2" />
							Export
						</Button>
					</div>
				</div>

				<div class="border border-border/50 rounded-xl overflow-hidden">
					<Table.Root>
						<Table.Header>
							<Table.Row class="bg-surface-elevated border-border/50">
								<Table.Head class="hidden w-[100px] sm:table-cell">Image</Table.Head>
								<Table.Head class="font-medium">Name</Table.Head>
								<Table.Head class="font-medium">Status</Table.Head>
								<Table.Head class="font-medium">Price</Table.Head>
								<Table.Head class="hidden md:table-cell font-medium">Sales</Table.Head>
								<Table.Head class="hidden md:table-cell font-medium">Created</Table.Head>
								<Table.Head class="w-[50px]"></Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#if $products.status === 'pending'}
								{#each Array(5) as _}
									<Table.Row class="border-border/50">
										<Table.Cell class="hidden sm:table-cell">
											<div class="w-12 h-12 bg-muted animate-pulse rounded-lg"></div>
										</Table.Cell>
										<Table.Cell><div class="h-4 bg-muted animate-pulse rounded w-3/4"></div></Table.Cell>
										<Table.Cell><div class="h-6 bg-muted animate-pulse rounded-full w-16"></div></Table.Cell>
										<Table.Cell><div class="h-4 bg-muted animate-pulse rounded w-20"></div></Table.Cell>
										<Table.Cell class="hidden md:table-cell"><div class="h-4 bg-muted animate-pulse rounded w-12"></div></Table.Cell>
										<Table.Cell class="hidden md:table-cell"><div class="h-4 bg-muted animate-pulse rounded w-24"></div></Table.Cell>
										<Table.Cell><div class="h-8 w-8 bg-muted animate-pulse rounded"></div></Table.Cell>
									</Table.Row>
								{/each}
							{:else if $products.status === 'error'}
								<Table.Row>
									<Table.Cell colspan={7} class="text-center py-8">
										<div class="text-muted-foreground">
											<Package class="w-8 h-8 mx-auto mb-2 opacity-50" />
											Error loading products: {$products.error.message}
										</div>
									</Table.Cell>
								</Table.Row>
							{:else if $products.data?.length === 0}
								<Table.Row>
									<Table.Cell colspan={7} class="text-center py-12">
										<div class="text-muted-foreground">
											<Package class="w-12 h-12 mx-auto mb-4 opacity-50" />
											<p class="text-lg font-medium mb-2">No products yet</p>
											<p class="text-sm mb-4">Start by adding your first product to your store</p>
											<Button href="/{$page.params.shop}/product-types/create" class="btn-gradient shadow-brand">
												<CirclePlus class="w-4 h-4 mr-2" />
												Add Product
											</Button>
										</div>
									</Table.Cell>
								</Table.Row>
							{:else}
								{#each $products.data as product}
									<Table.Row class="border-border/50 hover:bg-surface-elevated/50 transition-colors">
										<Table.Cell class="hidden sm:table-cell">
											<div class="w-12 h-12 bg-gradient-to-br from-muted to-muted/70 rounded-lg flex items-center justify-center">
												<Package class="w-6 h-6 text-muted-foreground" />
											</div>
										</Table.Cell>
										<Table.Cell class="font-medium">{product?.title || 'Untitled'}</Table.Cell>
										<Table.Cell>
											<Badge 
												variant="secondary" 
												class={product?.status === 'active' ? 'status-active' : 'status-inactive'}
											>
												{product?.status || 'Unknown'}
											</Badge>
										</Table.Cell>
										<Table.Cell class="font-medium">
											{product?.variants[0]?.price || 'N/A'}
										</Table.Cell>
										<Table.Cell class="hidden md:table-cell text-muted-foreground">25</Table.Cell>
										<Table.Cell class="hidden md:table-cell text-muted-foreground text-sm">
											{product?.created_at ? new Date(product.created_at).toLocaleDateString() : 'N/A'}
										</Table.Cell>
										<Table.Cell>
											<DropdownMenu.Root>
												<DropdownMenu.Trigger asChild let:builder>
													<Button
														aria-haspopup="true"
														size="icon"
														variant="ghost"
														builders={[builder]}
														class="h-8 w-8 rounded-lg hover:bg-surface-elevated"
													>
														<Ellipsis class="h-4 w-4" />
														<span class="sr-only">Toggle menu</span>
													</Button>
												</DropdownMenu.Trigger>
												<DropdownMenu.Content align="end" class="glass border-border/50 shadow-glass rounded-xl">
													<DropdownMenu.Label>Actions</DropdownMenu.Label>
													<DropdownMenu.Separator />
													<DropdownMenu.Item 
														href="/{$page.params.shop}/products/{product?.product_id || ''}"
														class="flex items-center gap-2"
													>
														<Edit class="w-4 h-4" />
														Edit
													</DropdownMenu.Item>
													<DropdownMenu.Item class="flex items-center gap-2">
														<Eye class="w-4 h-4" />
														View
													</DropdownMenu.Item>
													<DropdownMenu.Separator />
													<DropdownMenu.Item class="flex items-center gap-2 text-destructive">
														<Trash2 class="w-4 h-4" />
														Delete
													</DropdownMenu.Item>
												</DropdownMenu.Content>
											</DropdownMenu.Root>
										</Table.Cell>
									</Table.Row>
								{/each}
							{/if}
						</Table.Body>
					</Table.Root>
				</div>

				{#if $products.data && $products.data.length > 0}
					<div class="flex items-center justify-between mt-4 pt-4 border-t border-border/50">
						<div class="text-sm text-muted-foreground">
							Showing <strong>1-{Math.min(limit, $products.data.length)}</strong> of <strong>{$products.data.length}</strong> products
						</div>
						<Button 
							href="/{$page.params.shop}/product-types" 
							variant="outline" 
							size="sm"
							class="glass border-border/50"
						>
							View All Products
							<ArrowUpRight class="w-4 h-4 ml-2" />
						</Button>
					</div>
				{/if}
			</div>
		</div>

		<!-- Recent Orders -->
		<div class="space-y-6">
			<!-- Recent Orders Card -->
			<div class="card-elevated">
				<div class="flex items-center justify-between mb-6">
					<div>
						<h2 class="text-xl font-semibold text-foreground">Recent Orders</h2>
						<p class="text-sm text-muted-foreground">Latest customer orders</p>
					</div>
					<Button 
						href="/{$page.params.shop}/orders" 
						variant="outline" 
						size="sm"
						class="glass border-border/50"
					>
						View All
					</Button>
				</div>

				<div class="space-y-4">
					{#each recentOrders as order}
						<div class="flex items-center justify-between p-4 bg-surface-elevated rounded-xl border border-border/50 hover:bg-surface-muted/50 transition-colors">
							<div class="flex items-center gap-3">
								<div class="w-10 h-10 bg-gradient-to-br from-primary to-accent rounded-xl flex items-center justify-center">
									<ShoppingCart class="w-5 h-5 text-white" />
								</div>
								<div>
									<p class="font-medium text-foreground">{order.id}</p>
									<p class="text-sm text-muted-foreground">{order.customer}</p>
								</div>
							</div>
							<div class="text-right">
								<Badge 
									variant="secondary" 
									class={order.status === 'completed' ? 'status-active' : order.status === 'processing' ? 'status-warning' : 'status-inactive'}
								>
									{order.status}
								</Badge>
								<p class="text-sm font-medium text-foreground mt-1">{formatCurrencyWithLocale(order.amount, currencyCode)}</p>
							</div>
						</div>
					{/each}
				</div>
			</div>

			<!-- Quick Actions Card -->
			<div class="card-elevated">
				<h2 class="text-xl font-semibold text-foreground mb-4">Quick Actions</h2>
				<div class="space-y-3">
					<Button 
						href="/{$page.params.shop}/product-types/create" 
						class="w-full justify-start btn-gradient shadow-brand"
					>
						<CirclePlus class="w-4 h-4 mr-3" />
						Add New Product
					</Button>
					<Button 
						href="/{$page.params.shop}/orders" 
						variant="outline" 
						class="w-full justify-start glass border-border/50"
					>
						<ShoppingCart class="w-4 h-4 mr-3" />
						View Orders
					</Button>
					<Button 
						href="/{$page.params.shop}/customers" 
						variant="outline" 
						class="w-full justify-start glass border-border/50"
					>
						<Users class="w-4 h-4 mr-3" />
						Manage Customers
					</Button>
					<Button 
						href="/{$page.params.shop}/settings" 
						variant="outline" 
						class="w-full justify-start glass border-border/50"
					>
						<BarChart3 class="w-4 h-4 mr-3" />
						Analytics
					</Button>
				</div>
			</div>
		</div>
	</div>
</div>
