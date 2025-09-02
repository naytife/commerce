<script lang="ts">
	import { onMount, getContext } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { api, currentShopId } from '$lib/api';
	import type { 
		InventoryItem, 
		InventoryReport, 
		LowStockProduct,
		StockMovement, 
		StockUpdatePayload,
		InventorySearchParams,
		StockMovementSearchParams,
		StockMovementsResponse
	} from '$lib/types';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Table from '$lib/components/ui/table';
	import * as Card from '$lib/components/ui/card';
	import * as Select from '$lib/components/ui/select';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Tabs from '$lib/components/ui/tabs';
	import { Badge } from '$lib/components/ui/badge';
	import { Textarea } from '$lib/components/ui/textarea';
	import { Label } from '$lib/components/ui/label';
	import { getCurrencySymbol, formatCurrencyWithLocale } from '$lib/utils/currency';
	import type { Shop } from '$lib/types';
	import { createQuery } from '@tanstack/svelte-query';
	import { 
		Package, 
		Search, 
		Filter, 
		Edit, 
		TrendingDown, 
		TrendingUp, 
		AlertTriangle,
		Archive,
		Plus,
		Minus,
		RefreshCw,
		BarChart3,
		Download,
		MoreHorizontal,
		Eye,
		Trash2,
		RotateCcw,
		ChevronDown,
		ChevronRight
	} from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	const authFetch = getContext('authFetch') as typeof fetch;
	const apiWithAuth = api(authFetch);

	// Get shop currency
	const shopQuery = createQuery<Shop, Error>({
		queryKey: [`shop-${$page.params.shop}`],
		queryFn: () => apiWithAuth.getShop(),
		enabled: !!$page.params.shop
	});
	$: currencyCode = $shopQuery.data?.currency_code || 'USD';

	let inventoryReport: InventoryReport | null = null;
	let lowStockVariants: LowStockProduct[] = [];
	let stockMovements: StockMovement[] = [];
	let loading = true;
	let loadingLowStock = false;
	let loadingMovements = false;
	let searchQuery = '';
	let lowStockOnly = false;
	let outOfStockOnly = false;
	let currentPage = 1;
	let totalPages = 1;
	let showStockUpdateDialog = false;
	let showMovementDialog = false;
	let selectedVariant: InventoryItem | null = null;
	let activeTab: string = 'overview';

	// Collapsible table state - track expanded products
	let expandedProducts = new Set<string>();
	let expandedLowStockProducts = new Set<string>();
	let expandedMovementProducts = new Set<string>();

	// Stock update form
	let stockUpdateForm = {
		quantity: 0,
		movement_type: 'adjustment' as 'adjustment' | 'purchase' | 'damage' | 'transfer',
		reason: '',
		cost_price: 0
	};
	let stockUpdateMovementTypeSelect: { value: string; label: string } | undefined = undefined;

	// Movement filters
	let movementFilters = {
		movement_type: '',
		date_from: '',
		date_to: '',
		page: 1
	};
	let movementTypeFilterSelect: { value: string; label: string } | undefined = undefined;

	const loadInventoryReport = async (params: InventorySearchParams = {}) => {
		try {
			loading = true;
			inventoryReport = await apiWithAuth.getInventoryReport({
				...params,
				page: currentPage,
				limit: 20
			});
		} catch (error) {
			console.error('Error loading inventory report:', error);
			toast.error('Failed to load inventory report');
		} finally {
			loading = false;
		}
	};

	const loadLowStockVariants = async () => {
		try {
			loadingLowStock = true;
			lowStockVariants = await apiWithAuth.getLowStockVariants();
		} catch (error) {
			console.error('Error loading low stock variants:', error);
			toast.error('Failed to load low stock variants');
		} finally {
			loadingLowStock = false;
		}
	};

	const loadStockMovements = async (params: StockMovementSearchParams = {}) => {
		try {
			loadingMovements = true;
			const filters = {
				...movementFilters,
				movement_type: movementTypeFilterSelect?.value || movementFilters.movement_type
			};
			const response = await apiWithAuth.getStockMovements({
				...params,
				...filters
			});
			stockMovements = response.movements || [];
			totalPages = Math.ceil(response.total / response.limit) || 1;
		} catch (error) {
			console.error('Error loading stock movements:', error);
			toast.error('Failed to load stock movements');
			stockMovements = [];
			totalPages = 1;
		} finally {
			loadingMovements = false;
		}
	};

	const handleSearch = () => {
		currentPage = 1;
		const params: InventorySearchParams = {};
		
		if (searchQuery.trim()) {
			params.query = searchQuery.trim();
		}
		
		if (lowStockOnly) {
			params.low_stock_only = true;
		}
		
		if (outOfStockOnly) {
			params.out_of_stock_only = true;
		}
		
		loadInventoryReport(params);
	};

	const openStockUpdateDialog = (variant: InventoryItem) => {
		selectedVariant = variant;
		stockUpdateForm = {
			quantity: variant.current_stock,
			movement_type: 'adjustment',
			reason: '',
			cost_price: 0
		};
		stockUpdateMovementTypeSelect = { value: 'adjustment', label: 'Adjustment' };
		showStockUpdateDialog = true;
	};

	const handleStockUpdate = async () => {
		if (!selectedVariant) return;
		
		try {
			// Extract movement type from select with proper typing
			const movementType = (stockUpdateMovementTypeSelect?.value || 'adjustment') as 'adjustment' | 'purchase' | 'damage' | 'transfer';
			const updatePayload: StockUpdatePayload = {
				quantity: Number(stockUpdateForm.quantity),
				movement_type: movementType,
				reason: stockUpdateForm.reason,
				cost_price: stockUpdateForm.cost_price ? Number(stockUpdateForm.cost_price) : undefined
			};
			
			await apiWithAuth.updateVariantStock(selectedVariant.variant_id, updatePayload);
			showStockUpdateDialog = false;
			selectedVariant = null;
			toast.success('Stock updated successfully');
			loadInventoryReport();
			loadStockMovements();
		} catch (error) {
			console.error('Error updating stock:', error);
			toast.error('Failed to update stock');
		}
	};

	const getStockStatusBadge = (item: InventoryItem): { variant: 'default' | 'destructive' | 'outline' | 'secondary'; text: string } => {
		if (item.current_stock === 0) {
			return { variant: 'destructive', text: 'Out of Stock' };
		} else if (item.current_stock <= item.low_stock_threshold) {
			return { variant: 'secondary', text: 'Low Stock' };
		} else {
			return { variant: 'default', text: 'In Stock' };
		}
	};

	const getMovementTypeIcon = (type: string) => {
		switch (type) {
			case 'adjustment': return RefreshCw;
			case 'sale': return Minus;
			case 'purchase': return Plus;
			case 'return': return TrendingUp;
			case 'damage': return AlertTriangle;
			case 'transfer': return Archive;
			default: return RefreshCw;
		}
	};

	const getMovementTypeBadge = (type: string): { variant: 'default' | 'destructive' | 'outline' | 'secondary'; text: string } => {
		switch (type) {
			case 'adjustment': return { variant: 'outline', text: 'Adjustment' };
			case 'sale': return { variant: 'destructive', text: 'Sale' };
			case 'purchase': return { variant: 'default', text: 'Purchase' };
			case 'return': return { variant: 'secondary', text: 'Return' };
			case 'damage': return { variant: 'destructive', text: 'Damage' };
			case 'transfer': return { variant: 'outline', text: 'Transfer' };
			default: return { variant: 'outline', text: type };
		}
	};

	const formatDate = (dateString: string) => {
		return new Date(dateString).toLocaleDateString();
	};

	const formatDateTime = (dateString: string) => {
		return new Date(dateString).toLocaleString();
	};

	const formatCurrency = (amount: number) => {
		return formatCurrencyWithLocale(amount, currencyCode);
	};

	// Group inventory items by product
	const groupInventoryByProduct = (items: InventoryItem[]) => {
		const groups = new Map<string, {
			product_id: number;
			product_title: string;
			variants: InventoryItem[];
			total_stock: number;
			total_variants: number;
			lowest_stock: number;
		}>();

		items.forEach(item => {
			const key = `${item.product_id}-${item.product_title}`;
			if (!groups.has(key)) {
				groups.set(key, {
					product_id: item.product_id,
					product_title: item.product_title,
					variants: [],
					total_stock: 0,
					total_variants: 0,
					lowest_stock: Infinity
				});
			}
			const group = groups.get(key)!;
			group.variants.push(item);
			group.total_stock += item.current_stock;
			group.total_variants += 1;
			group.lowest_stock = Math.min(group.lowest_stock, item.current_stock);
		});

		return Array.from(groups.values()).sort((a, b) => a.product_title.localeCompare(b.product_title));
	};

	// Group low stock products by product
	const groupLowStockByProduct = (items: LowStockProduct[]) => {
		const groups = new Map<string, {
			product_name: string;
			variants: LowStockProduct[];
			total_variants: number;
			lowest_stock: number;
		}>();

		items.forEach(item => {
			const key = item.product_name;
			if (!groups.has(key)) {
				groups.set(key, {
					product_name: item.product_name,
					variants: [],
					total_variants: 0,
					lowest_stock: Infinity
				});
			}
			const group = groups.get(key)!;
			group.variants.push(item);
			group.total_variants += 1;
			group.lowest_stock = Math.min(group.lowest_stock, item.stock);
		});

		return Array.from(groups.values()).sort((a, b) => a.product_name.localeCompare(b.product_name));
	};

	// Group stock movements by product
	const groupMovementsByProduct = (movements: StockMovement[]) => {
		const groups = new Map<string, {
			product_title: string;
			movements: StockMovement[];
			total_movements: number;
		}>();

		movements.forEach(movement => {
			const key = movement.product_title;
			if (!groups.has(key)) {
				groups.set(key, {
					product_title: movement.product_title,
					movements: [],
					total_movements: 0
				});
			}
			const group = groups.get(key)!;
			group.movements.push(movement);
			group.total_movements += 1;
		});

		return Array.from(groups.values()).sort((a, b) => a.product_title.localeCompare(b.product_title));
	};

	// Toggle product expansion
	const toggleProductExpansion = (productKey: string, table: 'inventory' | 'lowstock' | 'movements') => {
		switch (table) {
			case 'inventory':
				if (expandedProducts.has(productKey)) {
					expandedProducts.delete(productKey);
				} else {
					expandedProducts.add(productKey);
				}
				expandedProducts = expandedProducts;
				break;
			case 'lowstock':
				if (expandedLowStockProducts.has(productKey)) {
					expandedLowStockProducts.delete(productKey);
				} else {
					expandedLowStockProducts.add(productKey);
				}
				expandedLowStockProducts = expandedLowStockProducts;
				break;
			case 'movements':
				if (expandedMovementProducts.has(productKey)) {
					expandedMovementProducts.delete(productKey);
				} else {
					expandedMovementProducts.add(productKey);
				}
				expandedMovementProducts = expandedMovementProducts;
				break;
		}
	};

	// Expand/Collapse all products in a table
	const expandAllProducts = (table: 'inventory' | 'lowstock' | 'movements') => {
		switch (table) {
			case 'inventory':
				if (inventoryReport?.items) {
					const allProductKeys = groupInventoryByProduct(inventoryReport.items)
						.map(g => `${g.product_id}-${g.product_title}`);
					allProductKeys.forEach(key => expandedProducts.add(key));
					expandedProducts = expandedProducts;
				}
				break;
			case 'lowstock':
				const allLowStockKeys = groupLowStockByProduct(lowStockVariants)
					.map(g => g.product_name);
				allLowStockKeys.forEach(key => expandedLowStockProducts.add(key));
				expandedLowStockProducts = expandedLowStockProducts;
				break;
			case 'movements':
				const allMovementKeys = groupMovementsByProduct(stockMovements)
					.map(g => g.product_title);
				allMovementKeys.forEach(key => expandedMovementProducts.add(key));
				expandedMovementProducts = expandedMovementProducts;
				break;
		}
	};

	const collapseAllProducts = (table: 'inventory' | 'lowstock' | 'movements') => {
		switch (table) {
			case 'inventory':
				expandedProducts.clear();
				expandedProducts = expandedProducts;
				break;
			case 'lowstock':
				expandedLowStockProducts.clear();
				expandedLowStockProducts = expandedLowStockProducts;
				break;
			case 'movements':
				expandedMovementProducts.clear();
				expandedMovementProducts = expandedMovementProducts;
				break;
		}
	};

	onMount(() => {
		// Wait for shop ID to be available before loading data
		let unsubscribe: (() => void) | undefined;
		unsubscribe = currentShopId.subscribe((shopId) => {
			if (shopId !== null) {
				loadInventoryReport();
				loadLowStockVariants();
				loadStockMovements();
				// Unsubscribe after initial load
				if (unsubscribe) {
					unsubscribe();
				}
			}
		});
	});

	// Reactive search - only trigger if shop ID is available
	$: if ($currentShopId !== null && searchQuery === '' && !lowStockOnly && !outOfStockOnly) {
		loadInventoryReport();
	}

	// Load movements when filters change - only if shop ID is available
	$: {
		if ($currentShopId !== null && activeTab === 'movements') {
			loadStockMovements();
		}
	}
</script>

<div class="space-y-8">
	<!-- Header Section -->
	<div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
		<div>
			<h1 class="text-3xl font-bold text-foreground mb-2">Inventory Management</h1>
			<p class="text-muted-foreground">
				Track stock levels, movements, and inventory health across all products
			</p>
		</div>
		<div class="flex items-center gap-3">
			<!-- <Button variant="outline" class="glass border-border/50">
				<Download class="w-4 h-4 mr-2" />
				Export
			</Button>
			<Button variant="outline" class="glass border-border/50">
				<BarChart3 class="w-4 h-4 mr-2" />
				Reports
			</Button>
			<Button on:click={() => (showMovementDialog = true)} class="btn-gradient shadow-brand">
				<Plus class="w-4 h-4 mr-2" />
				Record Movement
			</Button> -->
		</div>
	</div>

	<!-- Enhanced Stats Cards -->
	{#if inventoryReport}
		<!-- <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
			<div class="card-interactive text-center">
				<div class="w-12 h-12 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-brand">
					<Package class="w-6 h-6 text-white" />
				</div>
				<div class="text-2xl font-bold text-foreground mb-1">{inventoryReport.total_products.toLocaleString()}</div>
				<div class="text-sm text-muted-foreground">Total Products</div>
			</div>
			
			<div class="card-interactive text-center">
				<div class="w-12 h-12 bg-gradient-to-br from-success to-emerald-500 rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-brand">
					<Archive class="w-6 h-6 text-white" />
				</div>
				<div class="text-2xl font-bold text-foreground mb-1">{inventoryReport.total_variants.toLocaleString()}</div>
				<div class="text-sm text-muted-foreground">Total Variants</div>
			</div>
			
			<div class="card-interactive text-center">
				<div class="w-12 h-12 bg-gradient-to-br from-warning to-orange-500 rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-brand">
					<AlertTriangle class="w-6 h-6 text-white" />
				</div>
				<div class="text-2xl font-bold text-foreground mb-1">{inventoryReport.low_stock_count.toLocaleString()}</div>
				<div class="text-sm text-muted-foreground">Low Stock Items</div>
			</div>
			
			<div class="card-interactive text-center">
				<div class="w-12 h-12 bg-gradient-to-br from-destructive to-red-500 rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-brand">
					<TrendingDown class="w-6 h-6 text-white" />
				</div>
				<div class="text-2xl font-bold text-foreground mb-1">{inventoryReport.out_of_stock_count.toLocaleString()}</div>
				<div class="text-sm text-muted-foreground">Out of Stock</div>
			</div>
		</div> -->
	{/if}

	<!-- Enhanced Tabs Navigation -->
	<div class="card-elevated">
		<Tabs.Root value={activeTab} onValueChange={(value) => activeTab = value || 'overview'}>
			<Tabs.List class="grid w-full grid-cols-3 glass border-border/50">
				<Tabs.Trigger value="overview" class="data-[state=active]:bg-primary/10 data-[state=active]:text-primary">
					<Package class="w-4 h-4 mr-2" />
					Inventory Overview
				</Tabs.Trigger>
				<Tabs.Trigger value="low-stock" class="data-[state=active]:bg-primary/10 data-[state=active]:text-primary">
					<AlertTriangle class="w-4 h-4 mr-2" />
					Low Stock Alerts
				</Tabs.Trigger>
				<Tabs.Trigger value="movements" class="data-[state=active]:bg-primary/10 data-[state=active]:text-primary">
					<RefreshCw class="w-4 h-4 mr-2" />
					Stock Movements
				</Tabs.Trigger>
			</Tabs.List>

		<!-- Inventory Overview Tab -->
		<Tabs.Content value="overview" class="space-y-6">
			<!-- Enhanced Search and Filters -->
			<div class="space-y-6">
				<div class="flex flex-col lg:flex-row gap-4">
					<div class="flex-1 relative">
						<Search class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-muted-foreground" />
						<Input
							placeholder="Search products by name, SKU, or location..."
							bind:value={searchQuery}
							class="pl-10 h-12 glass border-border/50 focus:border-primary/50"
							on:keydown={(e) => e.key === 'Enter' && handleSearch()}
						/>
					</div>
					<div class="flex flex-wrap gap-3">
						<div class="flex items-center space-x-3 p-3 rounded-xl glass border border-border/50">
							<input
								type="checkbox"
								id="lowStockFilter"
								bind:checked={lowStockOnly}
								class="w-4 h-4 text-primary bg-transparent border-border/50 rounded focus:ring-primary/20 focus:ring-2"
							/>
							<Label for="lowStockFilter" class="text-sm font-medium text-foreground cursor-pointer">
								Low Stock Only
							</Label>
						</div>
						<div class="flex items-center space-x-3 p-3 rounded-xl glass border border-border/50">
							<input
								type="checkbox"
								id="outOfStockFilter"
								bind:checked={outOfStockOnly}
								class="w-4 h-4 text-primary bg-transparent border-border/50 rounded focus:ring-primary/20 focus:ring-2"
							/>
							<Label for="outOfStockFilter" class="text-sm font-medium text-foreground cursor-pointer">
								Out of Stock Only
							</Label>
						</div>
						<Button on:click={handleSearch} variant="outline" class="glass border-border/50 h-12">
							<Filter class="w-4 h-4 mr-2" />
							Apply Filters
						</Button>
					</div>
				</div>
			</div>

			<!-- Enhanced Inventory Table -->
			<div class="space-y-4">
				<!-- Table Controls -->
				{#if inventoryReport?.items && inventoryReport.items.length > 0}
					{@const groupedInventory = groupInventoryByProduct(inventoryReport.items)}
					{#if groupedInventory.length > 1}
						<div class="flex justify-between items-center">
							<p class="text-sm text-muted-foreground">
								{groupedInventory.length} products with {inventoryReport.items.length} variants total
							</p>
							<div class="flex gap-2">
								<Button variant="outline" size="sm" on:click={() => expandAllProducts('inventory')}>
									<ChevronDown class="w-4 h-4 mr-1" />
									Expand All
								</Button>
								<Button variant="outline" size="sm" on:click={() => collapseAllProducts('inventory')}>
									<ChevronRight class="w-4 h-4 mr-1" />
									Collapse All
								</Button>
							</div>
						</div>
					{/if}
				{/if}
				
				<div class="border border-border/50 rounded-xl overflow-hidden">
				{#if loading}
					<!-- Enhanced Loading State -->
					<div class="space-y-6 p-8">
						<div class="bg-surface-elevated border-b border-border/50 p-4">
							<div class="grid grid-cols-9 gap-4">
								{#each Array(9) as _}
									<div class="h-4 bg-muted animate-pulse rounded"></div>
								{/each}
							</div>
						</div>
						{#each Array(5) as _}
							<div class="border-b border-border/50 p-4">
								<div class="grid grid-cols-9 gap-4">
									{#each Array(9) as _}
										<div class="h-6 bg-muted animate-pulse rounded"></div>
									{/each}
								</div>
							</div>
						{/each}
					</div>
				{:else if !inventoryReport || !inventoryReport.items || inventoryReport.items.length === 0}
					<!-- Enhanced Empty State -->
					<div class="flex flex-col items-center justify-center py-20">
						<div class="w-20 h-20 bg-gradient-to-br from-primary to-accent rounded-full flex items-center justify-center mb-8 shadow-brand">
							<Package class="w-10 h-10 text-white" />
						</div>
						<h3 class="text-2xl font-bold mb-4 text-foreground">No inventory items found</h3>
						<p class="text-lg text-muted-foreground mb-8 text-center max-w-2xl leading-relaxed">
							Start managing your inventory by adding products to your store. Track stock levels, set alerts, and monitor movements.
						</p>
						<div class="flex flex-col sm:flex-row gap-4">
							<Button class="btn-gradient shadow-brand">
								<Plus class="w-5 h-5 mr-2" />
								Add Products
							</Button>
							<Button variant="outline" class="glass border-border/50">
								Import Inventory
							</Button>
						</div>
					</div>
				{:else}
					{@const groupedInventory = groupInventoryByProduct(inventoryReport.items)}
					<Table.Root>
						<Table.Header>
							<Table.Row class="bg-surface-elevated border-border/50 hover:bg-surface-elevated">
								<Table.Head class="font-medium w-8"></Table.Head>
								<Table.Head class="font-medium">Product</Table.Head>
								<Table.Head class="font-medium">SKU</Table.Head>
								<Table.Head class="font-medium">Current Stock</Table.Head>
								<Table.Head class="font-medium">Available</Table.Head>
								<Table.Head class="font-medium">Reserved</Table.Head>
								<Table.Head class="font-medium">Status</Table.Head>
								<Table.Head class="font-medium">Location</Table.Head>
								<Table.Head class="font-medium">Last Updated</Table.Head>
								<Table.Head class="w-[50px]"></Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each groupedInventory as productGroup}
								{@const productKey = `${productGroup.product_id}-${productGroup.product_title}`}
								{@const isExpanded = expandedProducts.has(productKey)}
								<!-- Product Header Row -->
								<Table.Row class="bg-gradient-to-r from-slate-50 to-slate-100 dark:from-slate-900/50 dark:to-slate-800/50 border-l-4 border-primary/60 hover:border-primary hover:shadow-sm hover:-translate-y-px transition-all duration-200 cursor-pointer font-medium">
									<Table.Cell on:click={() => toggleProductExpansion(productKey, 'inventory')}>
										<div class="transition-transform duration-200 {isExpanded ? 'rotate-0' : '-rotate-90'}">
											{#if isExpanded}
												<ChevronDown class="w-4 h-4 text-muted-foreground" />
											{:else}
												<ChevronRight class="w-4 h-4 text-muted-foreground" />
											{/if}
										</div>
									</Table.Cell>
									<Table.Cell class="font-semibold text-foreground" on:click={() => toggleProductExpansion(productKey, 'inventory')}>
										<div class="flex items-center gap-2">
											<Package class="w-4 h-4 text-primary" />
											<span>{productGroup.product_title}</span>
											<Badge variant="outline" class="text-xs">
												{productGroup.total_variants} variant{productGroup.total_variants !== 1 ? 's' : ''}
											</Badge>
										</div>
									</Table.Cell>
									<Table.Cell class="text-muted-foreground">
										<span class="text-xs">Multiple SKUs</span>
									</Table.Cell>
									<Table.Cell>
										<div class="flex items-center gap-2">
											<span class="font-semibold text-foreground">{productGroup.total_stock}</span>
											{#if productGroup.lowest_stock <= 10}
												<AlertTriangle class="w-4 h-4 text-warning" />
											{/if}
										</div>
									</Table.Cell>
									<Table.Cell colspan={6} class="text-muted-foreground text-xs">
										<div class="flex items-center justify-between">
											<div>{isExpanded ? 'Click to collapse variants' : 'Click to expand variants'}</div>
											<Button
												variant="outline"
												size="sm"
												on:click={(e) => { e.stopPropagation(); goto(`/${$page.params.shop}/products/${productGroup.product_id}`); }}
											>
												Edit Product
											</Button>
										</div>
									</Table.Cell>
								</Table.Row>
								
								<!-- Variant Rows (shown when expanded) -->
								{#if isExpanded}
									{#each productGroup.variants as item}
										<Table.Row class="border-l-2 border-border/30 bg-background/50 hover:bg-surface-elevated/30 hover:border-l-primary/50 transition-all duration-200">
											<Table.Cell class="pl-8">
												<!-- Empty space for indentation -->
											</Table.Cell>
											<Table.Cell class="font-medium pl-6">
												<div class="space-y-1">
													<div class="text-sm text-muted-foreground font-normal">↳ {item.variant_title}</div>
												</div>
											</Table.Cell>
											<Table.Cell>
												<Badge variant="outline" class="font-mono text-xs">
													{item.sku}
												</Badge>
											</Table.Cell>
											<Table.Cell>
												<div class="flex items-center gap-2">
													<span class="font-semibold text-foreground">{item.current_stock}</span>
													{#if item.current_stock <= item.low_stock_threshold}
														<AlertTriangle class="w-4 h-4 text-warning" />
													{/if}
												</div>
											</Table.Cell>
											<Table.Cell>
												<span class="text-foreground">{item.available_stock}</span>
											</Table.Cell>
											<Table.Cell>
												<span class="text-muted-foreground">{item.reserved_stock}</span>
											</Table.Cell>
											<Table.Cell>
												{@const status = getStockStatusBadge(item)}
												<Badge class={status.variant === 'destructive' ? 'status-destructive' : 
															   status.variant === 'secondary' ? 'status-warning' : 
															   'status-active'}>
													{status.text}
												</Badge>
											</Table.Cell>
											<Table.Cell>
												<span class="text-muted-foreground">{item.location || 'Default'}</span>
											</Table.Cell>
											<Table.Cell>
												<div class="text-sm text-muted-foreground">
													{formatDate(item.last_updated)}
												</div>
											</Table.Cell>
											<Table.Cell>
												<Button
													variant="ghost"
													size="icon"
													class="h-8 w-8 hover:bg-surface-elevated"
													on:click={() => openStockUpdateDialog(item)}
												>
													<Edit class="w-4 h-4" />
												</Button>
											</Table.Cell>
										</Table.Row>
									{/each}
								{/if}
							{/each}
						</Table.Body>
					</Table.Root>
				{/if}
				</div>
			</div>
		</Tabs.Content>

		<!-- Low Stock Alerts Tab -->
		<Tabs.Content value="low-stock" class="space-y-6">
			<Card.Root>
				<Card.Content class="p-6">
					{#if loadingLowStock}
						<div class="flex justify-center items-center h-64">
							<div class="text-center">
								<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
								<p class="mt-4 text-muted-foreground">Loading low stock alerts...</p>
							</div>
						</div>
					{:else if lowStockVariants.length === 0}
						<div class="text-center py-12">
							<TrendingUp class="mx-auto h-12 w-12 text-green-500" />
							<h3 class="mt-4 text-lg font-semibold">All stock levels are healthy!</h3>
							<p class="text-muted-foreground">No products are currently below their low stock threshold.</p>
						</div>
					{:else}
						{@const groupedLowStock = groupLowStockByProduct(lowStockVariants)}
						
						<!-- Table Controls -->
						{#if groupedLowStock.length > 1}
							<div class="flex justify-between items-center mb-4">
								<p class="text-sm text-muted-foreground">
									{groupedLowStock.length} products with {lowStockVariants.length} low stock variants
								</p>
								<div class="flex gap-2">
									<Button variant="outline" size="sm" on:click={() => expandAllProducts('lowstock')}>
										<ChevronDown class="w-4 h-4 mr-1" />
										Expand All
									</Button>
									<Button variant="outline" size="sm" on:click={() => collapseAllProducts('lowstock')}>
										<ChevronRight class="w-4 h-4 mr-1" />
										Collapse All
									</Button>
								</div>
							</div>
						{/if}
						
						<Table.Root>
							<Table.Header>
								<Table.Row>
									<Table.Head class="w-8"></Table.Head>
									<Table.Head>Product</Table.Head>
									<Table.Head>SKU</Table.Head>
									<Table.Head>Current Stock</Table.Head>
									<Table.Head>Threshold</Table.Head>
									<Table.Head>Days Remaining</Table.Head>
									<Table.Head class="text-right">Actions</Table.Head>
								</Table.Row>
							</Table.Header>
							<Table.Body>
								{#each groupedLowStock as productGroup}
									{@const productKey = productGroup.product_name}
									{@const isExpanded = expandedLowStockProducts.has(productKey)}
									<!-- Product Header Row -->
									<Table.Row class="bg-gradient-to-r from-orange-50 to-orange-100 dark:from-orange-950/30 dark:to-orange-900/30 border-l-4 border-orange-500/60 hover:border-orange-500 hover:shadow-sm hover:-translate-y-px transition-all duration-200 cursor-pointer font-medium">
										<Table.Cell on:click={() => toggleProductExpansion(productKey, 'lowstock')}>
											<div class="transition-transform duration-200 {isExpanded ? 'rotate-0' : '-rotate-90'}">
												{#if isExpanded}
													<ChevronDown class="w-4 h-4 text-muted-foreground" />
												{:else}
													<ChevronRight class="w-4 h-4 text-muted-foreground" />
												{/if}
											</div>
										</Table.Cell>
										<Table.Cell class="font-semibold text-foreground" on:click={() => toggleProductExpansion(productKey, 'lowstock')}>
											<div class="flex items-center gap-2">
												<AlertTriangle class="w-4 h-4 text-orange-500" />
												<span>{productGroup.product_name}</span>
												<Badge variant="outline" class="text-xs bg-orange-100 dark:bg-orange-950">
													{productGroup.total_variants} variant{productGroup.total_variants !== 1 ? 's' : ''}
												</Badge>
											</div>
										</Table.Cell>
										<Table.Cell class="text-muted-foreground">
											<span class="text-xs">Multiple SKUs</span>
										</Table.Cell>
										<Table.Cell>
											<div class="flex items-center gap-2">
												<span class="font-semibold text-orange-600">Lowest: {productGroup.lowest_stock}</span>
												<AlertTriangle class="w-4 h-4 text-orange-500" />
											</div>
										</Table.Cell>
										<Table.Cell colspan={3} class="text-muted-foreground text-xs">
											{isExpanded ? 'Click to collapse variants' : 'Click to expand variants'}
										</Table.Cell>
									</Table.Row>
									
									<!-- Variant Rows (shown when expanded) -->
									{#if isExpanded}
										{#each productGroup.variants as variant}
											<Table.Row class="border-l-2 border-orange-200/50 bg-orange-25/50 dark:bg-orange-975/10 hover:bg-orange-50/80 dark:hover:bg-orange-950/20 hover:border-l-orange-400 transition-all duration-200">
												<Table.Cell class="pl-8">
													<!-- Empty space for indentation -->
												</Table.Cell>
												<Table.Cell class="font-medium pl-6">
													<div class="space-y-1">
														<div class="text-sm text-muted-foreground font-normal">↳ {variant.description}</div>
													</div>
												</Table.Cell>
												<Table.Cell>
													<Badge variant="outline" class="font-mono text-xs">
														{variant.sku}
													</Badge>
												</Table.Cell>
												<Table.Cell>
													<span class="font-semibold text-orange-600">{variant.stock}</span>
													<AlertTriangle class="w-4 h-4 text-orange-500 inline ml-1" />
												</Table.Cell>
												<Table.Cell>-</Table.Cell>
												<Table.Cell>
													<span class="text-muted-foreground">N/A</span>
												</Table.Cell>
												<Table.Cell class="text-right">
													<Button
														variant="outline"
														size="sm"
														on:click={() => openStockUpdateDialog({
															variant_id: variant.product_variation_id,
															product_id: 0, // Not available in LowStockProduct
															product_title: variant.product_name,
															variant_title: variant.description,
															sku: variant.sku,
															current_stock: variant.stock,
															reserved_stock: 0,
															available_stock: variant.stock,
															low_stock_threshold: 0, // Not available in LowStockProduct
															last_updated: new Date().toISOString()
														})}
													>
														Restock
													</Button>
												</Table.Cell>
											</Table.Row>
										{/each}
									{/if}
								{/each}
							</Table.Body>
						</Table.Root>
					{/if}
				</Card.Content>
			</Card.Root>
		</Tabs.Content>

		<!-- Stock Movements Tab -->
		<Tabs.Content value="movements" class="space-y-6">
			<!-- Movement Filters -->
			<Card.Root>
				<Card.Content class="p-6">
					<div class="grid grid-cols-1 md:grid-cols-4 gap-4">
						<div class="space-y-2">
							<Label>Movement Type</Label>
							<Select.Root bind:selected={movementTypeFilterSelect}>
								<Select.Trigger>
									<Select.Value placeholder="All types" />
								</Select.Trigger>
								<Select.Content>
									<Select.Item value="" label="All Types">All Types</Select.Item>
									<Select.Item value="adjustment" label="Adjustment">Adjustment</Select.Item>
									<Select.Item value="sale" label="Sale">Sale</Select.Item>
									<Select.Item value="purchase" label="Purchase">Purchase</Select.Item>
									<Select.Item value="return" label="Return">Return</Select.Item>
									<Select.Item value="damage" label="Damage">Damage</Select.Item>
									<Select.Item value="transfer" label="Transfer">Transfer</Select.Item>
								</Select.Content>
							</Select.Root>
						</div>
						<div class="space-y-2">
							<Label for="dateFrom">From Date</Label>
							<Input
								id="dateFrom"
								type="date"
								bind:value={movementFilters.date_from}
							/>
						</div>
						<div class="space-y-2">
							<Label for="dateTo">To Date</Label>
							<Input
								id="dateTo"
								type="date"
								bind:value={movementFilters.date_to}
							/>
						</div>
						<div class="space-y-2">
							<Label>&nbsp;</Label>
							<Button on:click={() => loadStockMovements()} class="w-full">
								<Filter class="w-4 h-4 mr-2" />
								Apply Filters
							</Button>
						</div>
					</div>
				</Card.Content>
			</Card.Root>

			<!-- Stock Movements Table -->
			<Card.Root>
				<Card.Content class="p-6">
					{#if loadingMovements}
						<div class="flex justify-center items-center h-64">
							<div class="text-center">
								<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
								<p class="mt-4 text-muted-foreground">Loading stock movements...</p>
							</div>
						</div>
					{:else if stockMovements.length === 0}
						<div class="text-center py-12">
							<Archive class="mx-auto h-12 w-12 text-muted-foreground" />
							<h3 class="mt-4 text-lg font-semibold">No stock movements found</h3>
							<p class="text-muted-foreground">Stock movements will appear here as inventory changes.</p>
						</div>
					{:else}
						{@const groupedMovements = groupMovementsByProduct(stockMovements)}
						
						<!-- Table Controls -->
						{#if groupedMovements.length > 1}
							<div class="flex justify-between items-center mb-4">
								<p class="text-sm text-muted-foreground">
									{groupedMovements.length} products with {stockMovements.length} stock movements
								</p>
								<div class="flex gap-2">
									<Button variant="outline" size="sm" on:click={() => expandAllProducts('movements')}>
										<ChevronDown class="w-4 h-4 mr-1" />
										Expand All
									</Button>
									<Button variant="outline" size="sm" on:click={() => collapseAllProducts('movements')}>
										<ChevronRight class="w-4 h-4 mr-1" />
										Collapse All
									</Button>
								</div>
							</div>
						{/if}
						
						<Table.Root>
							<Table.Header>
								<Table.Row>
									<Table.Head class="w-8"></Table.Head>
									<Table.Head>Product</Table.Head>
									<Table.Head>Type</Table.Head>
									<Table.Head>Change</Table.Head>
									<Table.Head>Previous</Table.Head>
									<Table.Head>New</Table.Head>
									<Table.Head>Reason</Table.Head>
									<Table.Head>Date</Table.Head>
								</Table.Row>
							</Table.Header>
							<Table.Body>
								{#each groupedMovements as productGroup}
									{@const productKey = productGroup.product_title}
									{@const isExpanded = expandedMovementProducts.has(productKey)}
									<!-- Product Header Row -->
									<Table.Row class="bg-gradient-to-r from-slate-50 to-slate-100 dark:from-slate-900/50 dark:to-slate-800/50 border-l-4 border-primary/60 hover:border-primary hover:shadow-sm hover:-translate-y-px transition-all duration-200 cursor-pointer font-medium">
										<Table.Cell on:click={() => toggleProductExpansion(productKey, 'movements')}>
											<div class="transition-transform duration-200 {isExpanded ? 'rotate-0' : '-rotate-90'}">
												{#if isExpanded}
													<ChevronDown class="w-4 h-4 text-muted-foreground" />
												{:else}
													<ChevronRight class="w-4 h-4 text-muted-foreground" />
												{/if}
											</div>
										</Table.Cell>
										<Table.Cell class="font-semibold text-foreground" on:click={() => toggleProductExpansion(productKey, 'movements')}>
											<div class="flex items-center gap-2">
												<RefreshCw class="w-4 h-4 text-primary" />
												<span>{productGroup.product_title}</span>
												<Badge variant="outline" class="text-xs">
													{productGroup.total_movements} movement{productGroup.total_movements !== 1 ? 's' : ''}
												</Badge>
											</div>
										</Table.Cell>
										<Table.Cell colspan={6} class="text-muted-foreground text-xs">
											{isExpanded ? 'Click to collapse movements' : 'Click to expand movements'}
										</Table.Cell>
									</Table.Row>
									
									<!-- Movement Rows (shown when expanded) -->
									{#if isExpanded}
										{#each productGroup.movements as movement}
											<Table.Row class="border-l-2 border-border/30 bg-background/50 hover:bg-surface-elevated/30 hover:border-l-primary/50 transition-all duration-200">
												<Table.Cell class="pl-8">
													<!-- Empty space for indentation -->
												</Table.Cell>
												<Table.Cell class="font-medium pl-6">
													<div class="space-y-1">
														<div class="text-sm text-muted-foreground font-normal">↳ {movement.variant_title}</div>
														<div class="text-xs text-muted-foreground font-mono">{movement.sku}</div>
													</div>
												</Table.Cell>
												<Table.Cell>
													{@const type = getMovementTypeBadge(movement.movement_type)}
													<Badge variant={type.variant}>
														{type.text}
													</Badge>
												</Table.Cell>
												<Table.Cell>
													<span class="font-semibold {movement.quantity > 0 ? 'text-green-600' : 'text-red-600'}">
														{movement.quantity > 0 ? '+' : ''}{movement.quantity}
													</span>
												</Table.Cell>
												<Table.Cell>{movement.previous_stock}</Table.Cell>
												<Table.Cell>{movement.new_stock}</Table.Cell>
												<Table.Cell>{movement.reason || 'N/A'}</Table.Cell>
												<Table.Cell>{formatDateTime(movement.created_at)}</Table.Cell>
											</Table.Row>
										{/each}
									{/if}
								{/each}
							</Table.Body>
						</Table.Root>

						<!-- Pagination -->
						{#if totalPages > 1}
							<div class="flex justify-center items-center gap-2 mt-6">
								<Button
									variant="outline"
									size="sm"
									disabled={movementFilters.page === 1}
									on:click={() => {
										movementFilters.page--;
										loadStockMovements();
									}}
								>
									Previous
								</Button>
								<span class="text-sm text-muted-foreground">
									Page {movementFilters.page} of {totalPages}
								</span>
								<Button
									variant="outline"
									size="sm"
									disabled={movementFilters.page === totalPages}
									on:click={() => {
										movementFilters.page++;
										loadStockMovements();
									}}
								>
									Next
								</Button>
							</div>
						{/if}
					{/if}
				</Card.Content>
			</Card.Root>
		</Tabs.Content>
	</Tabs.Root>
</div>

<!-- Enhanced Stock Update Dialog -->
<Dialog.Root bind:open={showStockUpdateDialog}>
	<Dialog.Content class="sm:max-w-[600px] glass border-border/50 backdrop-blur-xl">
		<Dialog.Header class="space-y-3">
			<div class="w-12 h-12 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center mx-auto shadow-brand">
				<Edit class="w-6 h-6 text-white" />
			</div>
			<Dialog.Title class="text-2xl font-bold text-center text-foreground">Update Stock</Dialog.Title>
			<Dialog.Description class="text-center text-muted-foreground">
				Adjust stock levels for {selectedVariant?.product_title} - {selectedVariant?.variant_title}
			</Dialog.Description>
		</Dialog.Header>
		<div class="grid gap-6 py-6">
			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-2">
					<Label for="currentStock" class="text-sm font-medium text-foreground">Current Stock</Label>
					<Input
						id="currentStock"
						value={selectedVariant?.current_stock || 0}
						disabled
						class="form-input glass border-border/50 bg-surface-elevated text-muted-foreground"
					/>
				</div>
				<div class="space-y-2">
					<Label for="newQuantity" class="text-sm font-medium text-foreground">New Quantity *</Label>
					<Input
						id="newQuantity"
						type="number"
						bind:value={stockUpdateForm.quantity}
						min="0"
						class="form-input glass border-border/50 focus:border-primary/50"
					/>
				</div>
			</div>
			<div class="space-y-2">
				<Label for="movementType" class="text-sm font-medium text-foreground">Movement Type</Label>
				<Select.Root bind:selected={stockUpdateMovementTypeSelect}>
					<Select.Trigger class="w-full glass border-border/50 focus:border-primary/50">
						<Select.Value />
					</Select.Trigger>
					<Select.Content class="glass border-border/50 backdrop-blur-xl">
						<Select.Item value="adjustment" label="Adjustment">
							<span class="flex items-center gap-2">
								<RefreshCw class="w-4 h-4" />
								Adjustment
							</span>
						</Select.Item>
						<Select.Item value="purchase" label="Purchase">
							<span class="flex items-center gap-2">
								<Plus class="w-4 h-4" />
								Purchase
							</span>
						</Select.Item>
						<Select.Item value="damage" label="Damage">
							<span class="flex items-center gap-2">
								<AlertTriangle class="w-4 h-4" />
								Damage
							</span>
						</Select.Item>
						<Select.Item value="transfer" label="Transfer">
							<span class="flex items-center gap-2">
								<Archive class="w-4 h-4" />
								Transfer
							</span>
						</Select.Item>
					</Select.Content>
				</Select.Root>
			</div>
			<div class="space-y-2">
				<Label for="reason" class="text-sm font-medium text-foreground">Reason</Label>
				<Textarea
					id="reason"
					bind:value={stockUpdateForm.reason}
					placeholder="Reason for stock adjustment..."
					rows={3}
					class="form-input glass border-border/50 focus:border-primary/50 resize-none"
				/>
			</div>
			{#if stockUpdateMovementTypeSelect?.value === 'purchase'}
				<div class="space-y-2">
					<Label for="costPrice" class="text-sm font-medium text-foreground">Cost Price</Label>
					<Input
						id="costPrice"
						type="number"
						bind:value={stockUpdateForm.cost_price}
						min="0"
						step="0.01"
						placeholder="0.00"
						class="form-input glass border-border/50 focus:border-primary/50"
					/>
				</div>
			{/if}
		</div>
		<Dialog.Footer class="flex gap-3 pt-6 border-t border-border/20">
			<Button 
				variant="outline" 
				on:click={() => (showStockUpdateDialog = false)}
				class="flex-1 glass border-border/50 hover:bg-surface-elevated"
			>
				Cancel
			</Button>
			<Button 
				on:click={handleStockUpdate}
				class="flex-1 btn-gradient shadow-brand"
			>
				<Edit class="w-4 h-4 mr-2" />
				Update Stock
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
</div>
