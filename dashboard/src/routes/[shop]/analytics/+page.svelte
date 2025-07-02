<script lang="ts">
import { Chart, registerables } from 'chart.js';
Chart.register(...registerables);
import { api } from '$lib/api';
import { onMount, getContext } from 'svelte';
import { Line } from 'svelte-chartjs';
import { Badge } from '$lib/components/ui/badge';
import { Button } from '$lib/components/ui/button';
import * as Card from '$lib/components/ui/card';
import * as Table from '$lib/components/ui/table';
import { TrendingUp, DollarSign, ShoppingCart, Users, BarChart3, Package, AlertTriangle, ListFilter, File } from 'lucide-svelte';

// --- Types ---
interface SalesSummary {
  total_sales: number;
  total_orders: number;
  average_order_value: number;
}
interface OrdersOverTime {
  labels: string[];
  orders: number[];
}
interface TopProduct {
  product_id: string;
  name: string;
  units_sold: number;
  revenue: number;
}
interface CustomerSummary {
  new_customers: number;
  returning_customers: number;
  top_customers: TopCustomer[];
}
interface TopCustomer {
  customer_id: string;
  name: string;
  orders: number;
  total_spent: number;
  is_registered: boolean;
}
interface LowStockVariant {
  product_name: string;
  product_variation_id: number;
  sku: string;
  description: string;
  stock: number;
}

export let data: { shopId: number };

const authFetch = getContext('authFetch') as typeof fetch;
const analyticsApi = api(authFetch);

// --- State ---
let salesSummary: SalesSummary | null = null;
let ordersOverTime: OrdersOverTime | null = null;
let topProducts: TopProduct[] = [];
let customerSummary: CustomerSummary | null = null;
let lowStockVariants: LowStockVariant[] = [];
let loading = true;
let error = '';

let period: string = 'month';
let interval: string = 'day';
let threshold: number = 5;

// --- Chart config helper ---
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
          borderColor: 'rgb(75, 192, 192)',
          backgroundColor: 'rgba(75, 192, 192, 0.2)',
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

function formatCurrency(value: number | undefined) {
  return typeof value === 'number'
    ? value.toLocaleString(undefined, { style: 'currency', currency: 'USD', maximumFractionDigits: 2 })
    : '-';
}

// --- Data loading ---
async function loadAnalytics() {
  if (!data.shopId) return;
  loading = true;
  error = '';
  try {
    const [sales, orders, products, customers, lowStock] = await Promise.all([
      analyticsApi.getSalesSummary(data.shopId, { period }),
      analyticsApi.getOrdersOverTime(data.shopId, { interval, period }),
      analyticsApi.getTopProducts(data.shopId, { period }),
      analyticsApi.getCustomerSummary(data.shopId, { period }),
      analyticsApi.getLowStockVariants(data.shopId, { threshold })
    ]);
    salesSummary = sales.data;
    ordersOverTime = orders.data;
    topProducts = products.data;
    customerSummary = customers.data;
    lowStockVariants = lowStock.data;
  } catch (e: any) {
    error = e.message || 'Failed to load analytics data';
  } finally {
    loading = false;
  }
}

onMount(() => {
  if (data.shopId) loadAnalytics();
});

$: if (data.shopId && period && interval && threshold !== undefined) {
  loadAnalytics();
}

function handleThresholdInput(e: Event) {
  const target = e.target as HTMLInputElement;
  threshold = +target.value;
}

// Memoized chart config for orders over time
$: ordersChartConfig = (ordersOverTime && Array.isArray(ordersOverTime.labels) && ordersOverTime.labels.length > 0 && Array.isArray(ordersOverTime.orders) && ordersOverTime.orders.length > 0)
  ? getOrdersChartConfig(ordersOverTime)
  : null;
</script>

<div class="space-y-8">
  <!-- Header -->
  <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
    <div>
      <h1 class="text-3xl font-bold text-foreground mb-2">Analytics Dashboard</h1>
      <p class="text-muted-foreground">Track your store's performance, sales, and customer insights in real time.</p>
    </div>
    <div class="flex items-center gap-3">
      <Button variant="outline" class="glass border-border/50">
        <BarChart3 class="w-4 h-4 mr-2" />
        Export
      </Button>
      <Button variant="outline" class="glass border-border/50">
        <ListFilter class="w-4 h-4 mr-2" />
        Filters
      </Button>
    </div>
  </div>

  <!-- KPI Cards -->
  <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
    <div class="card-interactive">
      <div class="flex items-center justify-between mb-4">
        <div class="w-12 h-12 bg-gradient-to-br from-success to-emerald-500 rounded-2xl flex items-center justify-center shadow-brand">
          <DollarSign class="w-6 h-6 text-white" />
        </div>
      </div>
      <div class="space-y-1">
        <p class="text-sm text-muted-foreground">Total Sales</p>
        <p class="text-2xl font-bold text-foreground">{formatCurrency(salesSummary?.total_sales)}</p>
      </div>
    </div>
    <div class="card-interactive">
      <div class="flex items-center justify-between mb-4">
        <div class="w-12 h-12 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center shadow-brand">
          <ShoppingCart class="w-6 h-6 text-white" />
        </div>
      </div>
      <div class="space-y-1">
        <p class="text-sm text-muted-foreground">Total Orders</p>
        <p class="text-2xl font-bold text-foreground">{Number(salesSummary?.total_orders ?? 0)}</p>
      </div>
    </div>
    <div class="card-interactive">
      <div class="flex items-center justify-between mb-4">
        <div class="w-12 h-12 bg-gradient-to-br from-accent to-secondary rounded-2xl flex items-center justify-center shadow-brand">
          <Users class="w-6 h-6 text-white" />
        </div>
      </div>
      <div class="space-y-1">
        <p class="text-sm text-muted-foreground">New Customers</p>
        <p class="text-2xl font-bold text-foreground">{Number(customerSummary?.new_customers ?? 0)}</p>
      </div>
    </div>
    <div class="card-interactive">
      <div class="flex items-center justify-between mb-4">
        <div class="w-12 h-12 bg-gradient-to-br from-warning to-orange-500 rounded-2xl flex items-center justify-center shadow-brand">
          <TrendingUp class="w-6 h-6 text-white" />
        </div>
      </div>
      <div class="space-y-1">
        <p class="text-sm text-muted-foreground">Avg. Order Value</p>
        <p class="text-2xl font-bold text-foreground">{formatCurrency(salesSummary?.average_order_value)}</p>
      </div>
    </div>
  </div>

  <!-- Chart & Tables Grid -->
  <div class="grid gap-8 lg:grid-cols-3">
    <!-- Orders Over Time Chart -->
    <div class="lg:col-span-2">
      <div class="card-elevated">
        <div class="flex items-center justify-between mb-6">
          <div>
            <h2 class="text-xl font-semibold text-foreground">Orders Over Time</h2>
            <p class="text-sm text-muted-foreground">Visualize order trends for the selected period</p>
          </div>
        </div>
        <div class="bg-card rounded-lg">
          {#if loading}
            <div class="flex justify-center items-center min-h-[200px] text-muted-foreground">Loading chart...</div>
          {:else if error}
            <div class="text-red-600 text-center py-8">{error}</div>
          {:else if ordersChartConfig}
            <Line data={ordersChartConfig.data} options={ordersChartConfig.options} />
          {:else}
            <div class="text-center text-muted-foreground py-8">No order data available for chart.</div>
          {/if}
        </div>
      </div>
    </div>
    <!-- Top Products Table -->
    <div class="space-y-6">
      <div class="card-elevated">
        <div class="flex items-center justify-between mb-6">
          <div>
            <h2 class="text-xl font-semibold text-foreground">Top Products</h2>
            <p class="text-sm text-muted-foreground">Best selling products for the selected period</p>
          </div>
        </div>
        <div class="bg-card rounded-lg">
          <Table.Root>
            <Table.Header>
              <Table.Row>
                <Table.Head>Product</Table.Head>
                <Table.Head>Units Sold</Table.Head>
                <Table.Head>Revenue</Table.Head>
              </Table.Row>
            </Table.Header>
            <Table.Body>
              {#if loading}
                <Table.Row><Table.Cell colspan={3}>Loading...</Table.Cell></Table.Row>
              {:else if topProducts.length === 0}
                <Table.Row><Table.Cell colspan={3}>No data</Table.Cell></Table.Row>
              {:else}
                {#each topProducts as p}
                  <Table.Row>
                    <Table.Cell>{p.name}</Table.Cell>
                    <Table.Cell>{Number(p.units_sold ?? 0)}</Table.Cell>
                    <Table.Cell>{formatCurrency(p.revenue)}</Table.Cell>
                  </Table.Row>
                {/each}
              {/if}
            </Table.Body>
          </Table.Root>
        </div>
      </div>
      <!-- Top Customers Table -->
      <div class="card-elevated">
        <div class="flex items-center justify-between mb-6">
          <div>
            <h2 class="text-xl font-semibold text-foreground">Top Customers</h2>
            <p class="text-sm text-muted-foreground">Most valuable customers for the selected period</p>
          </div>
        </div>
        <div class="bg-card rounded-lg">
          <Table.Root>
            <Table.Header>
              <Table.Row>
                <Table.Head>Name</Table.Head>
                <Table.Head>Orders</Table.Head>
                <Table.Head>Total Spent</Table.Head>
                <Table.Head>Registered</Table.Head>
              </Table.Row>
            </Table.Header>
            <Table.Body>
              {#if loading}
                <Table.Row><Table.Cell colspan={4}>Loading...</Table.Cell></Table.Row>
              {:else if !customerSummary || !customerSummary.top_customers || customerSummary.top_customers.length === 0}
                <Table.Row><Table.Cell colspan={4}>No data</Table.Cell></Table.Row>
              {:else}
                {#each customerSummary.top_customers as c}
                  <Table.Row>
                    <Table.Cell>{c.name}</Table.Cell>
                    <Table.Cell>{Number(c.orders ?? 0)}</Table.Cell>
                    <Table.Cell>{formatCurrency(c.total_spent)}</Table.Cell>
                    <Table.Cell>{c.is_registered ? 'Yes' : 'No'}</Table.Cell>
                  </Table.Row>
                {/each}
              {/if}
            </Table.Body>
          </Table.Root>
        </div>
      </div>
      <!-- Low Stock Table -->
      <div class="card-elevated">
        <div class="flex items-center justify-between mb-6">
          <div>
            <h2 class="text-xl font-semibold text-foreground">Low Stock Variants</h2>
            <p class="text-sm text-muted-foreground">Variants at or below the low stock threshold</p>
          </div>
        </div>
        <div class="bg-card rounded-lg">
          <Table.Root>
            <Table.Header>
              <Table.Row>
                <Table.Head>Product</Table.Head>
                <Table.Head>Stock</Table.Head>
              </Table.Row>
            </Table.Header>
            <Table.Body>
              {#if loading}
                <Table.Row><Table.Cell colspan={2}>Loading...</Table.Cell></Table.Row>
              {:else if lowStockVariants.length === 0}
                <Table.Row><Table.Cell colspan={2}>No data</Table.Cell></Table.Row>
              {:else}
                {#each lowStockVariants as v}
                  <Table.Row>
                    <Table.Cell>{v.product_name}</Table.Cell>
                    <Table.Cell>{Number(v.stock ?? 0)}</Table.Cell>
                  </Table.Row>
                {/each}
              {/if}
            </Table.Body>
          </Table.Root>
        </div>
      </div>
    </div>
  </div>
</div> 