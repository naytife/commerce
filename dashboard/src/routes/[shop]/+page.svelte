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
	import { createQuery, useQueryClient, type CreateQueryResult } from '@tanstack/svelte-query';
	import { getContext } from 'svelte';
	import { api } from '$lib/api';
	import type { Product, SalesSummary, CustomerSummary, OrdersOverTime } from '$lib/types';
	import { page } from '$app/stores';
	import { getCurrencySymbol, formatCurrencyWithLocale } from '$lib/utils/currency';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as CalendarUI from '$lib/components/ui/calendar';
	import { Popover } from '$lib/components/ui/popover';
	import { addDays, startOfToday, startOfWeek, startOfMonth, endOfToday, endOfWeek, endOfMonth, format as formatDate, isSameDay } from 'date-fns';
	import { CalendarDate } from '@internationalized/date';
	import { persistentWritable } from '$lib/stores/persistent';
	import { Chart, registerables } from 'chart.js';
	Chart.register(...registerables);
	import { Line } from 'svelte-chartjs';
	import { onMount } from 'svelte';

	const client = useQueryClient()
	const authFetch = getContext('authFetch')

	// Get shop data for currency
	const shopQuery = createQuery({
		queryKey: [`shop-${$page.params.shop}`],
		queryFn: () => api(authFetch as any).getShop(),
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

	// Time range state
	const TIME_RANGES = [
		{ key: 'today', label: 'Today' },
		{ key: 'week', label: 'This Week' },
		{ key: 'month', label: 'This Month' },
		{ key: 'custom', label: 'Custom Range...' },
	];

	// Conversion helpers
	function toCalendarDate(date: Date): CalendarDate {
		return new CalendarDate(date.getFullYear(), date.getMonth() + 1, date.getDate());
	}
	function toJSDate(cd: CalendarDate): Date {
		return new Date(cd.year, cd.month - 1, cd.day);
	}

	function getDefaultRange() {
		const today = new Date();
		return {
			key: 'month',
			label: 'Last 30 days',
			start: toCalendarDate(addDays(today, -29)),
			end: toCalendarDate(today),
		};
	}

	function timeRangeReviver(obj: any) {
		if (obj && obj.start && obj.end) {
			return {
				...obj,
				start: new CalendarDate(obj.start.year, obj.start.month, obj.start.day),
				end: new CalendarDate(obj.end.year, obj.end.month, obj.end.day),
			};
		}
		return obj;
	}

	export const timeRange = persistentWritable('dashboard-time-range', getDefaultRange(), timeRangeReviver);
	let showCustomDialog = false;
	let customStart = getDefaultRange().start;
	let customEnd = getDefaultRange().end;

	$: buttonLabel = getButtonLabel($timeRange);

	function setRange(key: string) {
		const today = new Date();
		if (key === 'today') {
			timeRange.set({
				key,
				label: 'Today',
				start: null,
				end: null,
			});
		} else if (key === 'week') {
			timeRange.set({
				key,
				label: 'This Week',
				start:null,
				end:null,
			});
		} else if (key === 'month') {
			timeRange.set({
				key,
				label: 'This Month',
				start:null,
				end:null,
			});
		} else if (key === 'custom') {
			showCustomDialog = true;
			const defaultStart = toCalendarDate(addDays(today, -29));
			const defaultEnd = toCalendarDate(today);
			customStart = $timeRange.start
				? new CalendarDate($timeRange.start.year, $timeRange.start.month, $timeRange.start.day)
				: defaultStart;
			customEnd = $timeRange.end
				? new CalendarDate($timeRange.end.year, $timeRange.end.month, $timeRange.end.day)
				: defaultEnd;
		}
	}

	function applyCustomRange() {
		timeRange.set({
			key: 'custom',
			label: `${customStart.toString()} - ${customEnd.toString()}`,
			start: customStart,
			end: customEnd,
		});
		showCustomDialog = false;
	}

	function getButtonLabel(timeRange: { key: string; label: string; start: CalendarDate; end: CalendarDate }) {
		if (timeRange.key === 'custom') {
			const start = toJSDate(timeRange.start);
			const end = toJSDate(timeRange.end);
			return `${formatDate(start, 'MMM d, yyyy')} - ${formatDate(end, 'MMM d, yyyy')}`;
		}
		if (timeRange.key === 'today') {
			return 'Today';
		}
		if (timeRange.key === 'week') {
			return 'This Week';
		}
		if (timeRange.key === 'month') {
			return 'This Month';
		}
		return timeRange.label;
	}

	// Update queries to use $timeRange
	let salesSummary: CreateQueryResult<SalesSummary, Error>;
	$: salesSummary = createQuery<SalesSummary, Error>({
		queryKey: [`shop-${$page.params.shop}-sales-summary`, $timeRange.key, $timeRange.start, $timeRange.end],
		queryFn: () => {
			const period = $timeRange.key;
			if (period === 'custom') {
				return api(authFetch as any).getSalesSummary({
					period,
					start_date: toJSDate($timeRange.start).toISOString().slice(0, 10),
					end_date: toJSDate($timeRange.end).toISOString().slice(0, 10),
				});
			} else {
				return api(authFetch as any).getSalesSummary({ period });
			}
		},
	});

	let customerSummary:CreateQueryResult<CustomerSummary, Error>;
	$: customerSummary = createQuery<CustomerSummary, Error>({
		queryKey: [`shop-${$page.params.shop}-customer-summary`, $timeRange.key, $timeRange.start, $timeRange.end],
		queryFn: () => {
			const period = $timeRange.key;
			if (period === 'custom') {
				return api(authFetch as any).getCustomerSummary({
					period,
					start_date: toJSDate($timeRange.start).toISOString().slice(0, 10),
					end_date: toJSDate($timeRange.end).toISOString().slice(0, 10),
				});
			} else {
				return api(authFetch as any).getCustomerSummary({ period });
			}
		},
	});

	const recentOrders = createQuery({
		queryKey: [`shop-${$page.params.shop}-recent-orders`],
		queryFn: () => api(authFetch as any).getOrders({ limit: 4}),
	});

	// --- Orders Over Time Chart State ---
	let primaryColor = '141, 70%, 37%'; // fallback to default

	// Interval state for chart
	const INTERVALS = [
		{ key: 'day', label: 'Day' },
		{ key: 'week', label: 'Week' },
		{ key: 'month', label: 'Month' }
	];
	let interval = 'day';

	onMount(() => {
		const style = getComputedStyle(document.documentElement);
		const cssPrimary = style.getPropertyValue('--primary').trim();
		if (cssPrimary) {
			primaryColor = cssPrimary;
		}
	});

	function getOrdersChartConfig(data: OrdersOverTime | null) {
		if (!data || !data.labels?.length || !data.orders?.length) return null;
		return {
			data: {
				labels: data.labels,
				datasets: [
					{
						label: 'Orders',
						data: data.orders,
						fill: false,
						borderColor: `hsl(${primaryColor})`,
						backgroundColor: `hsla(${primaryColor}, 0.2)`,
						tension: 0.1
					}
				]
			},
			options: {
				responsive: true,
				plugins: {
					legend: { display: true },
					title: { display: false }
				},
				scales: {
					x: { title: { display: true, text: 'Date' } },
					y: {
						title: { display: true, text: 'Orders' },
						beginAtZero: true,
						ticks: {
							callback: function(this: any, value: string | number, index: number, ticks: any[]) {
								if (typeof value === 'number' && Number.isInteger(value)) return value;
								return null;
							},
							stepSize: 1,
							precision: 0
						}
					}
				}
			}
		};
	}

	// Orders Over Time Query
	let ordersOverTime: CreateQueryResult<OrdersOverTime, Error>;
	$: ordersOverTime = createQuery<OrdersOverTime, Error>({
		queryKey: [
			`shop-${$page.params.shop}-orders-over-time`,
			$timeRange.key,
			$timeRange.start,
			$timeRange.end,
			interval
		],
		queryFn: () => {
			const period = $timeRange.key;
			if (period === 'custom') {
				return api(authFetch as any).getOrdersOverTime({
					period,
					interval,
					start_date: toJSDate($timeRange.start).toISOString().slice(0, 10),
					end_date: toJSDate($timeRange.end).toISOString().slice(0, 10),
				});
			} else {
				return api(authFetch as any).getOrdersOverTime({ period, interval });
			}
		},
	});

	$: ordersChartConfig = ($ordersOverTime.status === 'success' && $ordersOverTime.data && Array.isArray($ordersOverTime.data.labels) && $ordersOverTime.data.labels.length > 0 && Array.isArray($ordersOverTime.data.orders) && $ordersOverTime.data.orders.length > 0)
		? getOrdersChartConfig($ordersOverTime.data)
		: null;
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
			<DropdownMenu.Root>
				<DropdownMenu.Trigger asChild let:builder>
					<Button builders={[builder]} variant="outline" class="glass border-border/50 min-w-[160px] justify-between">
						<Calendar class="w-4 h-4 mr-2" />
						{buttonLabel}
					</Button>
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end" class="glass border-border/50 shadow-glass rounded-xl">
					{#each TIME_RANGES as range}
						<DropdownMenu.Item on:click={() => setRange(range.key)} class={$timeRange.key === range.key ? 'selected' : ''}>{range.label}</DropdownMenu.Item>
					{/each}
				</DropdownMenu.Content>
			</DropdownMenu.Root>
			<Button href="/{$page.params.shop}/product-types/create" class="btn-gradient shadow-brand">
				<CirclePlus class="w-4 h-4 mr-2" />
				Add Product
			</Button>
		</div>
	</div>

	<Dialog.Root bind:open={showCustomDialog}>
		<Dialog.Content class="max-w-full w-[95vw] sm:max-w-2xl">
			<Dialog.Header>
				<Dialog.Title>Custom Date Range</Dialog.Title>
				<Dialog.Description>Select a start and end date for your custom range.</Dialog.Description>
			</Dialog.Header>
			<div class="flex flex-col gap-4 mt-4">
				<div class="flex flex-col sm:flex-row gap-4">
					<div class="flex-1">
						<label class="block text-sm mb-1">Start Date</label>
						<CalendarUI.Calendar bind:value={customStart} />
					</div>
					<div class="flex-1">
						<label class="block text-sm mb-1">End Date</label>
						<CalendarUI.Calendar bind:value={customEnd} />
					</div>
				</div>
				<div class="flex flex-col sm:flex-row justify-end gap-2 mt-4">
					<Button variant="outline" on:click={() => (showCustomDialog = false)}>Cancel</Button>
					<Button on:click={applyCustomRange} disabled={!customStart || !customEnd || customEnd < customStart}>Apply</Button>
				</div>
			</div>
		</Dialog.Content>
	</Dialog.Root>

	<!-- Stats Cards -->
	<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
		<!-- Revenue Card -->
		<div class="card-interactive">
			<div class="flex items-center justify-between mb-4">
				<div class="w-12 h-12 bg-gradient-to-br from-success to-emerald-500 rounded-2xl flex items-center justify-center shadow-brand">
					<DollarSign class="w-6 h-6 text-white" />
				</div>
				<Badge variant="secondary" class={($salesSummary.data?.sales_change_pct ?? 0) < 0 ? 'text-destructive bg-destructive/10 border-destructive/20' : 'text-success bg-success/10 border-success/20'}>
					{#if $salesSummary.isLoading}
						...
					{:else if $salesSummary.isError}
						!
					{:else}
						{($salesSummary.data?.sales_change_pct ?? 0) > 0 ? '+' : ''}{($salesSummary.data?.sales_change_pct ?? 0).toFixed(0)}%
					{/if}
				</Badge>
			</div>
			<div class="space-y-1">
				<p class="text-sm text-muted-foreground">Total Revenue</p>
				<p class="text-2xl font-bold text-foreground">
					{#if $salesSummary.isLoading}
						...
					{:else if $salesSummary.isError}
						Error
					{:else}
						{formatCurrencyWithLocale($salesSummary.data?.total_sales || 0, currencyCode)}
					{/if}
				</p>
			</div>
		</div>

		<!-- Orders Card -->
		<div class="card-interactive">
			<div class="flex items-center justify-between mb-4">
				<div class="w-12 h-12 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center shadow-brand">
					<ShoppingCart class="w-6 h-6 text-white" />
				</div>
				<Badge variant="secondary" class={($salesSummary.data?.orders_change_pct ?? 0) < 0 ? 'text-destructive bg-destructive/10 border-destructive/20' : 'text-primary bg-primary/10 border-primary/20'}>
					{#if $salesSummary.isLoading}
						...
					{:else if $salesSummary.isError}
						!
					{:else}
						{($salesSummary.data?.orders_change_pct ?? 0) > 0 ? '+' : ''}{($salesSummary.data?.orders_change_pct ?? 0).toFixed(0)}%
					{/if}
				</Badge>
			</div>
			<div class="space-y-1">
				<p class="text-sm text-muted-foreground">Total Orders</p>
				<p class="text-2xl font-bold text-foreground">
					{#if $salesSummary.isLoading}
						...
					{:else if $salesSummary.isError}
						Error
					{:else}
						{($salesSummary.data?.total_orders || 0).toLocaleString()}
					{/if}
				</p>
			</div>
		</div>

		<!-- Customers Card -->
		<div class="card-interactive">
			<div class="flex items-center justify-between mb-4">
				<div class="w-12 h-12 bg-gradient-to-br from-accent to-secondary rounded-2xl flex items-center justify-center shadow-brand">
					<Users class="w-6 h-6 text-white" />
				</div>
				<Badge variant="secondary" class={($salesSummary.data?.customers_change_pct ?? 0) < 0 ? 'text-destructive bg-destructive/10 border-destructive/20' : 'text-accent bg-accent/10 border-accent/20'}>
					{#if $salesSummary.isLoading}
						...
					{:else if $salesSummary.isError}
						!
					{:else}
						{($salesSummary.data?.customers_change_pct ?? 0) > 0 ? '+' : ''}{($salesSummary.data?.customers_change_pct ?? 0).toFixed(0)}%
					{/if}
				</Badge>
			</div>
			<div class="space-y-1">
				<p class="text-sm text-muted-foreground">Total Customers</p>
				<p class="text-2xl font-bold text-foreground">
					{#if $customerSummary.isLoading}
						...
					{:else if $customerSummary.isError}
						Error
					{:else}
						{($customerSummary.data?.new_customers || 0) + ($customerSummary.data?.returning_customers || 0)}
					{/if}
				</p>
			</div>
		</div>

		<!-- Conversion Rate Card -->
		<div class="card-interactive">
			<div class="flex items-center justify-between mb-4">
				<div class="w-12 h-12 bg-gradient-to-br from-secondary to-primary rounded-2xl flex items-center justify-center shadow-brand">
					<TrendingUp class="w-6 h-6 text-white" />
				</div>
				<Badge variant="secondary" class="text-secondary bg-secondary/10 border-secondary/20">
					--
				</Badge>
			</div>
			<div class="space-y-1">
				<p class="text-sm text-muted-foreground">Conversion Rate</p>
				<p class="text-2xl font-bold text-foreground">--</p>
			</div>
		</div>
	</div>

	<!-- Main Content Grid -->
	<div class="grid gap-8 lg:grid-cols-3">
		<!-- Left Column: Chart + Recent Products -->
		<div class="space-y-8 lg:col-span-2">
			<!-- Orders Over Time Chart -->
			<div class="card-elevated mb-0">
				<div class="flex items-center justify-between mb-6">
					<div>
						<h2 class="text-xl font-semibold text-foreground">Orders Over Time</h2>
						<p class="text-sm text-muted-foreground">Visualize order trends for the selected period</p>
					</div>
					<!-- Interval Button Group -->
					<div class="flex gap-2">
						{#each INTERVALS as i}
							<Button
								on:click={() => interval = i.key}
								variant={interval === i.key ? 'default' : 'outline'}
								size="sm"
								class="min-w-[64px]"
							>
								{i.label}
							</Button>
						{/each}
					</div>
				</div>
				<div class="rounded-lg">
					{#if $ordersOverTime.status === 'pending'}
						<div class="flex justify-center items-center min-h-[200px] text-muted-foreground">Loading chart...</div>
					{:else if $ordersOverTime.status === 'error'}
						<div class="text-red-600 text-center py-8">{$ordersOverTime.error.message}</div>
					{:else if ordersChartConfig}
						<Line data={ordersChartConfig.data} options={ordersChartConfig.options} height={120} />
					{:else}
						<div class="text-center text-muted-foreground py-8">No order data available for chart.</div>
					{/if}
				</div>
			</div>

			<!-- Recent Products -->
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

		<!-- Right Column: Recent Orders + Quick Actions -->
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
					{#if $recentOrders.isLoading}
						{#each Array(4) as _}
							<div class="flex items-center justify-between p-4 bg-surface-elevated rounded-xl border border-border/50 animate-pulse">
								<div class="flex items-center gap-3">
									<div class="w-10 h-10 bg-gradient-to-br from-primary to-accent rounded-xl flex items-center justify-center">
										<ShoppingCart class="w-5 h-5 text-white" />
									</div>
									<div>
										<div class="h-4 w-20 bg-muted rounded mb-1"></div>
										<div class="h-3 w-16 bg-muted rounded"></div>
									</div>
								</div>
								<div class="text-right">
									<div class="h-6 w-16 bg-muted rounded mb-1"></div>
									<div class="h-4 w-12 bg-muted rounded"></div>
								</div>
							</div>
						{/each}
					{:else if $recentOrders.isError}
						<div class="text-center text-muted-foreground py-8">Error loading recent orders.</div>
					{:else if !$recentOrders.data?.data || $recentOrders.data.data.length === 0}
						<div class="text-center text-muted-foreground py-8">No recent orders.</div>
					{:else}
						{#each $recentOrders.data?.data ?? [] as order}
							<div class="flex items-center justify-between p-4 bg-surface-elevated rounded-xl border border-border/50 hover:bg-surface-muted/50 transition-colors">
								<div class="flex items-center gap-3">
									<div class="w-10 h-10 bg-gradient-to-br from-primary to-accent rounded-xl flex items-center justify-center">
										<ShoppingCart class="w-5 h-5 text-white" />
									</div>
									<div>
										<p class="font-medium text-foreground">#{order.order_id}</p>
										<p class="text-sm text-muted-foreground truncate max-w-[140px]">{order.customer_name || order.customer_email}</p>
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
					{/if}
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
