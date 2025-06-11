<script lang="ts">
	import { onMount, getContext } from 'svelte';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { api } from '$lib/api';
	import type { Customer, CustomerSearchParams, PaginatedResponse } from '$lib/types';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Table from '$lib/components/ui/table';
	import * as Card from '$lib/components/ui/card';
	import * as Select from '$lib/components/ui/select';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Badge } from '$lib/components/ui/badge';
	import { Textarea } from '$lib/components/ui/textarea';
	import { Label } from '$lib/components/ui/label';
	import { Plus, Search, Filter, Edit, Trash2, Mail, Phone, Calendar, Users, UserPlus, Download, MoreHorizontal, Eye, MapPin } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';
	import { getCurrencySymbol, formatCurrencyWithLocale } from '$lib/utils/currency';
	import type { Shop } from '$lib/types';
	import { createQuery } from '@tanstack/svelte-query';

	const authFetch = getContext('authFetch') as typeof fetch;
	const apiWithAuth = api(authFetch);

	// Get shop currency
	const shopQuery = createQuery<Shop, Error>({
		queryKey: [`shop-${$page.params.shop}`],
		queryFn: () => apiWithAuth.getShop(),
		enabled: !!$page.params.shop
	});
	$: currencyCode = $shopQuery.data?.currency_code || 'USD';

	let customers: Customer[] = [];
	let loading = true;
	let searchQuery = '';
	let statusFilter: { value: string; label: string } | undefined = undefined;
	let currentPage = 1;
	let totalPages = 1;
	let totalCustomers = 0;
	let showCreateDialog = false;
	let showEditDialog = false;
	let selectedCustomer: Customer | null = null;

	// Create customer form
	let createForm = {
		first_name: '',
		last_name: '',
		email: '',
		phone: '',
		date_of_birth: '',
		gender: '',
		marketing_consent: false,
		preferred_language: 'en'
	};
	let createGenderSelect: { value: string; label: string } | undefined = undefined;

	// Edit customer form
	let editForm = {
		first_name: '',
		last_name: '',
		email: '',
		phone: '',
		date_of_birth: '',
		gender: '',
		marketing_consent: false,
		preferred_language: '',
		status: 'active' as 'active' | 'inactive' | 'suspended'
	};
	let editStatusSelect: { value: string; label: string } | undefined = undefined;
	let editGenderSelect: { value: string; label: string } | undefined = undefined;

	const loadCustomers = async (params: CustomerSearchParams = {}) => {
		try {
			loading = true;
			const response = await apiWithAuth.getCustomers({
				...params,
				page: currentPage,
				limit: 20
			});
			
			customers = response.data || [];
			totalPages = response.pagination?.total_pages || 1;
			totalCustomers = response.pagination?.total || 0;
		} catch (error) {
			console.error('Error loading customers:', error);
			toast.error('Failed to load customers');
			customers = [];
			totalPages = 1;
			totalCustomers = 0;
		} finally {
			loading = false;
		}
	};

	const handleSearch = () => {
		currentPage = 1;
		const params: CustomerSearchParams = {};
		
		if (searchQuery.trim()) {
			params.query = searchQuery.trim();
		}
		
		if (statusFilter?.value) {
			params.status = statusFilter.value;
		}
		
		loadCustomers(params);
	};

	const handleCreateCustomer = async () => {
		try {
			const formData = {
				...createForm,
				gender: (createGenderSelect?.value || createForm.gender) as 'male' | 'female' | 'other' | 'prefer_not_to_say' | undefined
			};
			await apiWithAuth.createCustomer(formData);
			showCreateDialog = false;
			resetCreateForm();
			loadCustomers();
			toast.success('Customer created successfully');
		} catch (error) {
			console.error('Error creating customer:', error);
			toast.error('Failed to create customer');
		}
	};

	const handleEditCustomer = async () => {
		if (!selectedCustomer) return;
		
		try {
			const formData = {
				...editForm,
				status: (editStatusSelect?.value || editForm.status) as 'active' | 'inactive' | 'suspended' | undefined,
				gender: (editGenderSelect?.value || editForm.gender) as 'male' | 'female' | 'other' | 'prefer_not_to_say' | undefined
			};
			await apiWithAuth.updateCustomer(selectedCustomer.customer_id, formData);
			showEditDialog = false;
			selectedCustomer = null;
			loadCustomers();
			toast.success('Customer updated successfully');
		} catch (error) {
			console.error('Error updating customer:', error);
			toast.error('Failed to update customer');
		}
	};

	const handleDeleteCustomer = async (customer: Customer) => {
		if (!confirm(`Are you sure you want to delete ${customer.first_name} ${customer.last_name}?`)) {
			return;
		}
		
		try {
			await apiWithAuth.deleteCustomer(customer.customer_id);
			loadCustomers();
		} catch (error) {
			console.error('Error deleting customer:', error);
		}
	};

	const openEditDialog = (customer: Customer) => {
		selectedCustomer = customer;
		editForm = {
			first_name: customer.first_name,
			last_name: customer.last_name,
			email: customer.email,
			phone: customer.phone || '',
			date_of_birth: customer.date_of_birth || '',
			gender: customer.gender || '',
			marketing_consent: customer.marketing_consent,
			preferred_language: customer.preferred_language || 'en',
			status: customer.status
		};
		
		// Set select values
		editStatusSelect = customer.status ? 
			{ value: customer.status, label: customer.status.charAt(0).toUpperCase() + customer.status.slice(1) } : 
			undefined;
		editGenderSelect = customer.gender ? 
			{ value: customer.gender, label: customer.gender === 'prefer_not_to_say' ? 'Prefer not to say' : customer.gender.charAt(0).toUpperCase() + customer.gender.slice(1) } : 
			undefined;
			
		showEditDialog = true;
	};

	const resetCreateForm = () => {
		createForm = {
			first_name: '',
			last_name: '',
			email: '',
			phone: '',
			date_of_birth: '',
			gender: '',
			marketing_consent: false,
			preferred_language: 'en'
		};
		createGenderSelect = undefined;
	};

	const formatDate = (dateString: string) => {
		if (!dateString) return 'N/A';
		return new Date(dateString).toLocaleDateString();
	};

	const formatCurrency = (amount: number) => {
		return formatCurrencyWithLocale(amount, currencyCode);
	};

	const getStatusBadgeVariant = (status: string) => {
		switch (status) {
			case 'active': return 'default';
			case 'inactive': return 'secondary';
			case 'suspended': return 'destructive';
			default: return 'outline';
		}
	};

	const getStatusClass = (status: string) => {
		switch (status) {
			case 'active': return 'status-active';
			case 'inactive': return 'status-inactive';
			case 'suspended': return 'status-destructive';
			default: return 'status-inactive';
		}
	};

	onMount(() => {
		loadCustomers();
	});

	// Reactive search - only run on client side
	$: if (browser && searchQuery === '' && !statusFilter?.value) {
		loadCustomers();
	}
</script>

<div class="space-y-8">
	<!-- Header Section -->
	<div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
		<div>
			<h1 class="text-3xl font-bold text-foreground mb-2">Customers</h1>
			<p class="text-muted-foreground">
				Manage your customer relationships and track engagement
			</p>
		</div>
		<div class="flex items-center gap-3">
			<Button variant="outline" class="glass border-border/50">
				<Download class="w-4 h-4 mr-2" />
				Export
			</Button>
			<Button on:click={() => (showCreateDialog = true)} class="btn-gradient shadow-brand">
				<UserPlus class="w-4 h-4 mr-2" />
				Add Customer
			</Button>
		</div>
	</div>

	<!-- Enhanced Stats Cards -->
	<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
		<div class="card-interactive text-center">
			<div class="w-12 h-12 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-brand">
				<Users class="w-6 h-6 text-white" />
			</div>
			<div class="text-2xl font-bold text-foreground mb-1">{totalCustomers.toLocaleString()}</div>
			<div class="text-sm text-muted-foreground">Total Customers</div>
		</div>
		
		<div class="card-interactive text-center">
			<div class="w-12 h-12 bg-gradient-to-br from-success to-emerald-500 rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-brand">
				<Mail class="w-6 h-6 text-white" />
			</div>
			<div class="text-2xl font-bold text-foreground mb-1">{customers.filter(c => c.status === 'active').length.toLocaleString()}</div>
			<div class="text-sm text-muted-foreground">Active Customers</div>
		</div>
		
		<div class="card-interactive text-center">
			<div class="w-12 h-12 bg-gradient-to-br from-secondary to-accent rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-brand">
				<Calendar class="w-6 h-6 text-white" />
			</div>
			<div class="text-2xl font-bold text-foreground mb-1">
				{customers.filter(c => new Date(c.created_at).getMonth() === new Date().getMonth()).length.toLocaleString()}
			</div>
			<div class="text-sm text-muted-foreground">New This Month</div>
		</div>
		
		<div class="card-interactive text-center">
			<div class="w-12 h-12 bg-gradient-to-br from-accent to-orange-500 rounded-2xl flex items-center justify-center mx-auto mb-4 shadow-brand">
				<Phone class="w-6 h-6 text-white" />
			</div>
			<div class="text-2xl font-bold text-foreground mb-1">
				{customers.length > 0 ? Math.round(customers.reduce((sum, c) => sum + c.order_count, 0) / customers.length) : 0}
			</div>
			<div class="text-sm text-muted-foreground">Avg. Orders</div>
		</div>
	</div>

	<!-- Customers Management -->
	<div class="card-elevated">
		{#if loading}
			<!-- Enhanced Loading State -->
			<div class="space-y-6">
				<div class="flex flex-col lg:flex-row gap-4">
					<div class="flex-1">
						<div class="h-12 bg-muted animate-pulse rounded-xl"></div>
					</div>
					<div class="flex gap-2">
						<div class="h-12 w-32 bg-muted animate-pulse rounded-xl"></div>
						<div class="h-12 w-24 bg-muted animate-pulse rounded-xl"></div>
					</div>
				</div>

				<div class="border border-border/50 rounded-xl overflow-hidden">
					<div class="bg-surface-elevated border-b border-border/50 p-4">
						<div class="grid grid-cols-8 gap-4">
							{#each Array(8) as _}
								<div class="h-4 bg-muted animate-pulse rounded"></div>
							{/each}
						</div>
					</div>
					{#each Array(5) as _}
						<div class="border-b border-border/50 p-4">
							<div class="grid grid-cols-8 gap-4">
								{#each Array(8) as _}
									<div class="h-6 bg-muted animate-pulse rounded"></div>
								{/each}
							</div>
						</div>
					{/each}
				</div>
			</div>
		{:else if customers.length === 0}
			<!-- Enhanced Empty State -->
			<div class="flex flex-col items-center justify-center py-20">
				<div class="w-20 h-20 bg-gradient-to-br from-primary to-accent rounded-full flex items-center justify-center mb-8 shadow-brand">
					<Users class="w-10 h-10 text-white" />
				</div>
				<h3 class="text-2xl font-bold mb-4 text-foreground">No customers yet</h3>
				<p class="text-lg text-muted-foreground mb-8 text-center max-w-2xl leading-relaxed">
					Start building your customer base by adding your first customer. Track their orders, preferences, and engagement over time.
				</p>
				<div class="flex flex-col sm:flex-row gap-4">
					<Button on:click={() => (showCreateDialog = true)} class="btn-gradient shadow-brand">
						<UserPlus class="w-5 h-5 mr-2" />
						Add First Customer
					</Button>
					<Button variant="outline" class="glass border-border/50">
						Import Customers
					</Button>
				</div>
			</div>
		{:else}
			<!-- Customer Management Interface -->
			<div class="space-y-6">
				<!-- Search and Filters -->
				<div class="flex flex-col lg:flex-row gap-4">
					<div class="flex-1 relative">
						<Search class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-muted-foreground" />
						<Input
							placeholder="Search customers by name, email, or phone..."
							bind:value={searchQuery}
							class="pl-10 h-12 glass border-border/50 focus:border-primary/50"
							on:keydown={(e) => e.key === 'Enter' && handleSearch()}
						/>
					</div>
					<div class="flex gap-2">
						<Select.Root bind:selected={statusFilter}>
							<Select.Trigger class="w-[180px] h-12 glass border-border/50">
								<Select.Value placeholder="Filter by status" />
							</Select.Trigger>
							<Select.Content>
								<Select.Item value="" label="All Status">All Status</Select.Item>
								<Select.Item value="active" label="Active">
									<span class="flex items-center gap-2">
										<div class="w-2 h-2 rounded-full bg-success"></div>
										Active
									</span>
								</Select.Item>
								<Select.Item value="inactive" label="Inactive">
									<span class="flex items-center gap-2">
										<div class="w-2 h-2 rounded-full bg-muted-foreground"></div>
										Inactive
									</span>
								</Select.Item>
								<Select.Item value="suspended" label="Suspended">
									<span class="flex items-center gap-2">
										<div class="w-2 h-2 rounded-full bg-destructive"></div>
										Suspended
									</span>
								</Select.Item>
							</Select.Content>
						</Select.Root>
						<Button on:click={handleSearch} variant="outline" class="glass border-border/50 h-12">
							<Filter class="w-4 h-4 mr-2" />
							Apply
						</Button>
					</div>
				</div>

				<!-- Customers Table -->
				<div class="border border-border/50 rounded-xl overflow-hidden">
					<Table.Root>
						<Table.Header>
							<Table.Row class="bg-surface-elevated border-border/50 hover:bg-surface-elevated">
								<Table.Head class="font-medium">Customer</Table.Head>
								<Table.Head class="font-medium">Contact</Table.Head>
								<Table.Head class="font-medium">Status</Table.Head>
								<Table.Head class="font-medium">Orders</Table.Head>
								<Table.Head class="font-medium">Total Spent</Table.Head>
								<Table.Head class="font-medium">Joined</Table.Head>
								<Table.Head class="w-[50px]"></Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each customers as customer}
								<Table.Row class="border-border/50 hover:bg-surface-elevated/50 transition-colors">
									<Table.Cell>
										<div class="flex items-center gap-3">
											<div class="w-10 h-10 bg-gradient-to-br from-primary to-accent rounded-full flex items-center justify-center text-white font-semibold text-sm">
												{customer.first_name.charAt(0)}{customer.last_name.charAt(0)}
											</div>
											<div>
												<div class="font-semibold text-foreground">{customer.first_name} {customer.last_name}</div>
												<div class="text-sm text-muted-foreground">ID: {customer.customer_id}</div>
											</div>
										</div>
									</Table.Cell>
									<Table.Cell>
										<div class="space-y-1">
											<div class="flex items-center gap-2 text-sm">
												<Mail class="w-3 h-3 text-muted-foreground" />
												<span class="text-foreground">{customer.email}</span>
											</div>
											{#if customer.phone}
												<div class="flex items-center gap-2 text-sm">
													<Phone class="w-3 h-3 text-muted-foreground" />
													<span class="text-muted-foreground">{customer.phone}</span>
												</div>
											{/if}
										</div>
									</Table.Cell>
									<Table.Cell>
										<Badge class={getStatusClass(customer.status)}>
											{customer.status}
										</Badge>
									</Table.Cell>
									<Table.Cell>
										<div class="font-medium text-foreground">{customer.order_count}</div>
									</Table.Cell>
									<Table.Cell>
										<div class="font-semibold text-foreground">{formatCurrency(customer.total_spent)}</div>
									</Table.Cell>
									<Table.Cell>
										<div class="text-sm">
											<div class="font-medium text-foreground">{formatDate(customer.created_at)}</div>
										</div>
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
												<DropdownMenu.Item class="flex items-center gap-2">
													<Eye class="w-4 h-4" />
													View Profile
												</DropdownMenu.Item>
												<DropdownMenu.Item 
													on:click={() => openEditDialog(customer)}
													class="flex items-center gap-2"
												>
													<Edit class="w-4 h-4" />
													Edit Customer
												</DropdownMenu.Item>
												<DropdownMenu.Separator />
												<DropdownMenu.Item
													class="flex items-center gap-2 text-destructive"
													on:click={() => handleDeleteCustomer(customer)}
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

				<!-- Enhanced Pagination -->
				{#if totalPages > 1}
					<div class="flex items-center justify-between pt-4 border-t border-border/50">
						<div class="text-sm text-muted-foreground">
							Showing <strong>{((currentPage - 1) * 20) + 1}-{Math.min(currentPage * 20, totalCustomers)}</strong> of <strong>{totalCustomers}</strong> customers
						</div>
						<div class="flex items-center gap-2">
							<Button
								variant="outline"
								size="sm"
								class="glass border-border/50"
								disabled={currentPage === 1}
								on:click={() => {currentPage--; loadCustomers();}}
							>
								Previous
							</Button>
							<span class="text-sm text-muted-foreground px-3">
								Page {currentPage} of {totalPages}
							</span>
							<Button
								variant="outline"
								size="sm"
								class="glass border-border/50"
								disabled={currentPage === totalPages}
								on:click={() => {currentPage++; loadCustomers();}}
							>
								Next
							</Button>
						</div>
					</div>
				{/if}
			</div>
		{/if}
	</div>
</div>

<!-- Create Customer Dialog -->
<Dialog.Root bind:open={showCreateDialog}>
	<Dialog.Content class="sm:max-w-[600px] glass border-border/50 backdrop-blur-xl">
		<Dialog.Header class="space-y-3">
			<div class="w-12 h-12 bg-gradient-to-br from-primary to-accent rounded-2xl flex items-center justify-center mx-auto shadow-brand">
				<UserPlus class="w-6 h-6 text-white" />
			</div>
			<Dialog.Title class="text-2xl font-bold text-center text-foreground">Add New Customer</Dialog.Title>
			<Dialog.Description class="text-center text-muted-foreground">
				Create a new customer account with their basic information and preferences.
			</Dialog.Description>
		</Dialog.Header>
		<div class="grid gap-6 py-6">
			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-2">
					<Label for="firstName" class="text-sm font-medium text-foreground">First Name *</Label>
					<Input 
						id="firstName" 
						bind:value={createForm.first_name} 
						placeholder="John" 
						class="form-input glass border-border/50 focus:border-primary/50"
					/>
				</div>
				<div class="space-y-2">
					<Label for="lastName" class="text-sm font-medium text-foreground">Last Name *</Label>
					<Input 
						id="lastName" 
						bind:value={createForm.last_name} 
						placeholder="Doe" 
						class="form-input glass border-border/50 focus:border-primary/50"
					/>
				</div>
			</div>
			<div class="space-y-2">
				<Label for="email" class="text-sm font-medium text-foreground">Email Address *</Label>
				<Input 
					id="email" 
					type="email" 
					bind:value={createForm.email} 
					placeholder="john@example.com"
					class="form-input glass border-border/50 focus:border-primary/50"
				/>
			</div>
			<div class="space-y-2">
				<Label for="phone" class="text-sm font-medium text-foreground">Phone Number</Label>
				<Input 
					id="phone" 
					bind:value={createForm.phone} 
					placeholder="+1 (555) 123-4567"
					class="form-input glass border-border/50 focus:border-primary/50"
				/>
			</div>
			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-2">
					<Label for="dateOfBirth" class="text-sm font-medium text-foreground">Date of Birth</Label>
					<Input 
						id="dateOfBirth" 
						type="date" 
						bind:value={createForm.date_of_birth}
						class="form-input glass border-border/50 focus:border-primary/50"
					/>
				</div>
				<div class="space-y-2">
					<Label for="gender" class="text-sm font-medium text-foreground">Gender</Label>
					<Select.Root bind:selected={createGenderSelect}>
						<Select.Trigger class="w-full glass border-border/50 focus:border-primary/50">
							<Select.Value placeholder="Select gender" />
						</Select.Trigger>
						<Select.Content class="glass border-border/50 backdrop-blur-xl">
							<Select.Item value="prefer_not_to_say" label="Prefer not to say">Prefer not to say</Select.Item>
							<Select.Item value="male" label="Male">Male</Select.Item>
							<Select.Item value="female" label="Female">Female</Select.Item>
							<Select.Item value="other" label="Other">Other</Select.Item>
						</Select.Content>
					</Select.Root>
				</div>
			</div>
			<div class="flex items-start space-x-3 p-4 rounded-xl glass border border-border/50">
				<input
					type="checkbox"
					id="marketingConsent"
					bind:checked={createForm.marketing_consent}
					class="mt-0.5 w-4 h-4 text-primary bg-transparent border-border/50 rounded focus:ring-primary/20 focus:ring-2"
				/>
				<div>
					<Label for="marketingConsent" class="text-sm font-medium text-foreground cursor-pointer">
						Marketing Communications
					</Label>
					<p class="text-xs text-muted-foreground mt-1">
						Customer consents to receive promotional emails, SMS, and other marketing materials.
					</p>
				</div>
			</div>
		</div>
		<Dialog.Footer class="flex gap-3 pt-6 border-t border-border/20">
			<Button 
				variant="outline" 
				on:click={() => (showCreateDialog = false)}
				class="flex-1 glass border-border/50 hover:bg-surface-elevated"
			>
				Cancel
			</Button>
			<Button 
				on:click={handleCreateCustomer}
				class="flex-1 btn-gradient shadow-brand"
			>
				<UserPlus class="w-4 h-4 mr-2" />
				Create Customer
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<!-- Edit Customer Dialog -->
<Dialog.Root bind:open={showEditDialog}>
	<Dialog.Content class="sm:max-w-[600px] glass border-border/50 backdrop-blur-xl">
		<Dialog.Header class="space-y-3">
			<div class="w-12 h-12 bg-gradient-to-br from-accent to-primary rounded-2xl flex items-center justify-center mx-auto shadow-brand">
				<Edit class="w-6 h-6 text-white" />
			</div>
			<Dialog.Title class="text-2xl font-bold text-center text-foreground">Edit Customer</Dialog.Title>
			<Dialog.Description class="text-center text-muted-foreground">
				Update customer information and account status.
			</Dialog.Description>
		</Dialog.Header>
		<div class="grid gap-6 py-6">
			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-2">
					<Label for="editFirstName" class="text-sm font-medium text-foreground">First Name *</Label>
					<Input 
						id="editFirstName" 
						bind:value={editForm.first_name} 
						class="form-input glass border-border/50 focus:border-primary/50"
					/>
				</div>
				<div class="space-y-2">
					<Label for="editLastName" class="text-sm font-medium text-foreground">Last Name *</Label>
					<Input 
						id="editLastName" 
						bind:value={editForm.last_name} 
						class="form-input glass border-border/50 focus:border-primary/50"
					/>
				</div>
			</div>
			<div class="space-y-2">
				<Label for="editEmail" class="text-sm font-medium text-foreground">Email Address *</Label>
				<Input 
					id="editEmail" 
					type="email" 
					bind:value={editForm.email} 
					class="form-input glass border-border/50 focus:border-primary/50"
				/>
			</div>
			<div class="space-y-2">
				<Label for="editPhone" class="text-sm font-medium text-foreground">Phone Number</Label>
				<Input 
					id="editPhone" 
					bind:value={editForm.phone} 
					class="form-input glass border-border/50 focus:border-primary/50"
				/>
			</div>
			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-2">
					<Label for="editStatus" class="text-sm font-medium text-foreground">Account Status</Label>
					<Select.Root bind:selected={editStatusSelect}>
						<Select.Trigger class="w-full glass border-border/50 focus:border-primary/50">
							<Select.Value />
						</Select.Trigger>
						<Select.Content class="glass border-border/50 backdrop-blur-xl">
							<Select.Item value="active" label="Active">
								<span class="flex items-center gap-2">
									<div class="w-2 h-2 rounded-full bg-success"></div>
									Active
								</span>
							</Select.Item>
							<Select.Item value="inactive" label="Inactive">
								<span class="flex items-center gap-2">
									<div class="w-2 h-2 rounded-full bg-muted-foreground"></div>
									Inactive
								</span>
							</Select.Item>
							<Select.Item value="suspended" label="Suspended">
								<span class="flex items-center gap-2">
									<div class="w-2 h-2 rounded-full bg-destructive"></div>
									Suspended
								</span>
							</Select.Item>
						</Select.Content>
					</Select.Root>
				</div>
				<div class="space-y-2">
					<Label for="editGender" class="text-sm font-medium text-foreground">Gender</Label>
					<Select.Root bind:selected={editGenderSelect}>
						<Select.Trigger class="w-full glass border-border/50 focus:border-primary/50">
							<Select.Value placeholder="Select gender" />
						</Select.Trigger>
						<Select.Content class="glass border-border/50 backdrop-blur-xl">
							<Select.Item value="prefer_not_to_say" label="Prefer not to say">Prefer not to say</Select.Item>
							<Select.Item value="male" label="Male">Male</Select.Item>
							<Select.Item value="female" label="Female">Female</Select.Item>
							<Select.Item value="other" label="Other">Other</Select.Item>
						</Select.Content>
					</Select.Root>
				</div>
			</div>
			<div class="flex items-start space-x-3 p-4 rounded-xl glass border border-border/50">
				<input
					type="checkbox"
					id="editMarketingConsent"
					bind:checked={editForm.marketing_consent}
					class="mt-0.5 w-4 h-4 text-primary bg-transparent border-border/50 rounded focus:ring-primary/20 focus:ring-2"
				/>
				<div>
					<Label for="editMarketingConsent" class="text-sm font-medium text-foreground cursor-pointer">
						Marketing Communications
					</Label>
					<p class="text-xs text-muted-foreground mt-1">
						Customer consents to receive promotional emails, SMS, and other marketing materials.
					</p>
				</div>
			</div>
		</div>
		<Dialog.Footer class="flex gap-3 pt-6 border-t border-border/20">
			<Button 
				variant="outline" 
				on:click={() => (showEditDialog = false)}
				class="flex-1 glass border-border/50 hover:bg-surface-elevated"
			>
				Cancel
			</Button>
			<Button 
				on:click={handleEditCustomer}
				class="flex-1 btn-gradient shadow-brand"
			>
				<Edit class="w-4 h-4 mr-2" />
				Update Customer
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
