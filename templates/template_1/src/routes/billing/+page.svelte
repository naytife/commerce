<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Badge } from '$lib/components/ui/badge';
	import { CreditCard, ArrowLeft, Plus, Edit, Trash2, Shield, Check, AlertTriangle } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';
	
	// Mock payment methods data
	interface PaymentMethod {
		id: number;
		type: string;
		brand: string;
		last4: string;
		expiryMonth: number;
		expiryYear: number;
		isDefault: boolean;
		name: string;
	}

	let paymentMethods: PaymentMethod[] = [
		{
			id: 1,
			type: 'card',
			brand: 'Visa',
			last4: '4242',
			expiryMonth: 12,
			expiryYear: 2025,
			isDefault: true,
			name: 'John Doe'
		},
		{
			id: 2,
			type: 'card',
			brand: 'Mastercard',
			last4: '8888',
			expiryMonth: 8,
			expiryYear: 2026,
			isDefault: false,
			name: 'John Doe'
		}
	];

	let showAddForm = false;
	let editingMethod: PaymentMethod | null = null;
	let isProcessing = false;
	let showDeleteDialog = false;
	let methodToDelete: PaymentMethod | null = null;

	// Form validation
	let formErrors = {
		cardNumber: '',
		cardName: '',
		expiryDate: '',
		cvv: '',
		billingAddress: '',
		billingCity: '',
		billingState: '',
		billingZip: ''
	};

	// Form data
	let formData = {
		cardNumber: '',
		cardName: '',
		expiryDate: '',
		cvv: '',
		billingAddress: '',
		billingCity: '',
		billingState: '',
		billingZip: ''
	};

	function addNewMethod() {
		showAddForm = true;
		editingMethod = null;
		resetForm();
	}

	function editMethod(method: PaymentMethod) {
		editingMethod = { ...method };
		showAddForm = true;
	}

	function deleteMethod(id: number) {
		paymentMethods = paymentMethods.filter(method => method.id !== id);
	}

	function setDefaultMethod(id: number) {
		paymentMethods = paymentMethods.map(method => ({
			...method,
			isDefault: method.id === id
		}));
	}

	function saveMethod(event: Event) {
		event.preventDefault();
		// In a real app, this would save to an API
		showAddForm = false;
		editingMethod = null;
	}

	function cancelForm() {
		showAddForm = false;
		editingMethod = null;
	}

	function getCardIcon(brand: string) {
		// In a real app, you'd return appropriate card icons
		return CreditCard;
	}

	function resetForm() {
		formData = {
			cardNumber: '',
			cardName: '',
			expiryDate: '',
			cvv: '',
			billingAddress: '',
			billingCity: '',
			billingState: '',
			billingZip: ''
		};
		formErrors = {
			cardNumber: '',
			cardName: '',
			expiryDate: '',
			cvv: '',
			billingAddress: '',
			billingCity: '',
			billingState: '',
			billingZip: ''
		};
	}

	function validateCardNumber(cardNumber: string): boolean {
		// Simple Luhn algorithm implementation
		const cleanNumber = cardNumber.replace(/\s/g, '');
		if (!/^\d{13,19}$/.test(cleanNumber)) return false;
		
		let sum = 0;
		let isEven = false;
		
		for (let i = cleanNumber.length - 1; i >= 0; i--) {
			let digit = parseInt(cleanNumber[i]);
			
			if (isEven) {
				digit *= 2;
				if (digit > 9) digit -= 9;
			}
			
			sum += digit;
			isEven = !isEven;
		}
		
		return sum % 10 === 0;
	}

	function validateExpiryDate(expiryDate: string): boolean {
		const match = expiryDate.match(/^(\d{2})\/(\d{2})$/);
		if (!match) return false;
		
		const month = parseInt(match[1]);
		const year = parseInt(match[2]) + 2000;
		const now = new Date();
		const currentYear = now.getFullYear();
		const currentMonth = now.getMonth() + 1;
		
		if (month < 1 || month > 12) return false;
		if (year < currentYear || (year === currentYear && month < currentMonth)) return false;
		
		return true;
	}

	function formatCardNumber(value: string): string {
		const cleaned = value.replace(/\s/g, '');
		const groups = cleaned.match(/.{1,4}/g) || [];
		return groups.join(' ').substr(0, 19);
	}

	function formatExpiryDate(value: string): string {
		const cleaned = value.replace(/\D/g, '');
		if (cleaned.length >= 2) {
			return cleaned.substr(0, 2) + '/' + cleaned.substr(2, 2);
		}
		return cleaned;
	}

	function detectCardBrand(cardNumber: string): string {
		const cleaned = cardNumber.replace(/\s/g, '');
		
		if (/^4/.test(cleaned)) return 'Visa';
		if (/^5[1-5]/.test(cleaned)) return 'Mastercard';
		if (/^3[47]/.test(cleaned)) return 'American Express';
		if (/^6011/.test(cleaned)) return 'Discover';
		
		return 'Unknown';
	}

	function validateForm(): boolean {
		let isValid = true;
		const errors = { ...formErrors };

		// Card number validation
		if (!formData.cardNumber) {
			errors.cardNumber = 'Card number is required';
			isValid = false;
		} else if (!validateCardNumber(formData.cardNumber)) {
			errors.cardNumber = 'Invalid card number';
			isValid = false;
		} else {
			errors.cardNumber = '';
		}

		// Cardholder name validation
		if (!formData.cardName.trim()) {
			errors.cardName = 'Cardholder name is required';
			isValid = false;
		} else {
			errors.cardName = '';
		}

		// Expiry date validation
		if (!formData.expiryDate) {
			errors.expiryDate = 'Expiry date is required';
			isValid = false;
		} else if (!validateExpiryDate(formData.expiryDate)) {
			errors.expiryDate = 'Invalid or expired date';
			isValid = false;
		} else {
			errors.expiryDate = '';
		}

		// CVV validation
		if (!formData.cvv) {
			errors.cvv = 'CVV is required';
			isValid = false;
		} else if (!/^\d{3,4}$/.test(formData.cvv)) {
			errors.cvv = 'Invalid CVV';
			isValid = false;
		} else {
			errors.cvv = '';
		}

		// Billing address validation
		if (!formData.billingAddress.trim()) {
			errors.billingAddress = 'Address is required';
			isValid = false;
		} else {
			errors.billingAddress = '';
		}

		if (!formData.billingCity.trim()) {
			errors.billingCity = 'City is required';
			isValid = false;
		} else {
			errors.billingCity = '';
		}

		if (!formData.billingState.trim()) {
			errors.billingState = 'State is required';
			isValid = false;
		} else {
			errors.billingState = '';
		}

		if (!formData.billingZip.trim()) {
			errors.billingZip = 'ZIP code is required';
			isValid = false;
		} else if (!/^\d{5}(-\d{4})?$/.test(formData.billingZip)) {
			errors.billingZip = 'Invalid ZIP code';
			isValid = false;
		} else {
			errors.billingZip = '';
		}

		formErrors = errors;
		return isValid;
	}

	function handleCardNumberInput(event: Event) {
		const target = event.target as HTMLInputElement;
		const formatted = formatCardNumber(target.value);
		formData.cardNumber = formatted;
		target.value = formatted;
	}

	function handleExpiryDateInput(event: Event) {
		const target = event.target as HTMLInputElement;
		const formatted = formatExpiryDate(target.value);
		formData.expiryDate = formatted;
		target.value = formatted;
	}

	async function handleSaveMethod(event: Event) {
		event.preventDefault();
		
		if (!validateForm()) {
			toast.error('Please fix the errors below');
			return;
		}

		isProcessing = true;
		
		try {
			// Simulate API call
			await new Promise(resolve => setTimeout(resolve, 2000));
			
			// Create new payment method
			const newMethod = {
				id: Date.now(),
				type: 'card',
				brand: detectCardBrand(formData.cardNumber),
				last4: formData.cardNumber.replace(/\s/g, '').slice(-4),
				expiryMonth: parseInt(formData.expiryDate.split('/')[0]),
				expiryYear: parseInt('20' + formData.expiryDate.split('/')[1]),
				isDefault: paymentMethods.length === 0, // First card is default
				name: formData.cardName
			};

			if (editingMethod) {
				// Update existing method
				paymentMethods = paymentMethods.map(method => 
					method.id === editingMethod!.id ? { ...newMethod, id: editingMethod!.id } : method
				);
				toast.success('Payment method updated successfully');
			} else {
				// Add new method
				paymentMethods = [...paymentMethods, newMethod];
				toast.success('Payment method added successfully');
			}

			showAddForm = false;
			editingMethod = null;
			resetForm();
		} catch (error) {
			toast.error('Failed to save payment method');
		} finally {
			isProcessing = false;
		}
	}

	function confirmDelete(method: PaymentMethod) {
		methodToDelete = method;
		showDeleteDialog = true;
	}

	async function handleDeleteConfirm() {
		if (!methodToDelete) return;
		
		try {
			paymentMethods = paymentMethods.filter(method => method.id !== methodToDelete!.id);
			toast.success('Payment method deleted successfully');
			
			// If deleted method was default, set first remaining as default
			if (methodToDelete!.isDefault && paymentMethods.length > 0) {
				paymentMethods[0].isDefault = true;
			}
		} catch (error) {
			toast.error('Failed to delete payment method');
		} finally {
			showDeleteDialog = false;
			methodToDelete = null;
		}
	}

	function handleDeleteCancel() {
		showDeleteDialog = false;
		methodToDelete = null;
	}

	// ...existing code...
</script>

<svelte:head>
	<title>Billing & Payment Methods</title>
	<meta name="description" content="Manage your payment methods and billing information" />
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
				<h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">Billing & Payment</h1>
				<p class="text-gray-600 dark:text-gray-400">Manage your payment methods and billing information</p>
			</div>
			<Button on:click={addNewMethod}>
				<Plus class="h-4 w-4 mr-2" />
				Add Payment Method
			</Button>
		</div>
	</div>

	<div class="space-y-6">
		<!-- Security Notice -->
		<Card.Root class="border-blue-200 bg-blue-50 dark:border-blue-800 dark:bg-blue-950">
			<Card.Content class="pt-6">
				<div class="flex items-center gap-3">
					<Shield class="h-5 w-5 text-blue-600 dark:text-blue-400" />
					<div>
						<p class="text-sm font-medium text-blue-900 dark:text-blue-100">
							Your payment information is secure
						</p>
						<p class="text-sm text-blue-700 dark:text-blue-300">
							We use industry-standard encryption to protect your data
						</p>
					</div>
				</div>
			</Card.Content>
		</Card.Root>

		<!-- Payment Methods List -->
		{#if !showAddForm}
			{#each paymentMethods as method}
				<Card.Root>
					<Card.Header>
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-3">
								<svelte:component this={getCardIcon(method.brand)} class="h-6 w-6" />
								<div>
									<Card.Title class="text-lg">
										{method.brand} ••••••••••••{method.last4}
									</Card.Title>
									<p class="text-sm text-gray-600 dark:text-gray-400">
										Expires {method.expiryMonth.toString().padStart(2, '0')}/{method.expiryYear}
									</p>
								</div>
								{#if method.isDefault}
									<Badge>Default</Badge>
								{/if}
							</div>
							<div class="flex items-center gap-2">
								<Button variant="ghost" size="icon" on:click={() => editMethod(method)}>
									<Edit class="h-4 w-4" />
								</Button>
								<Button variant="ghost" size="icon" on:click={() => confirmDelete(method)}>
									<Trash2 class="h-4 w-4" />
								</Button>
							</div>
						</div>
					</Card.Header>
					<Card.Content>
						<div class="space-y-2">
							<p class="text-gray-600 dark:text-gray-400">Cardholder: {method.name}</p>
						</div>
						{#if !method.isDefault}
							<div class="mt-4">
								<Button variant="outline" size="sm" on:click={() => setDefaultMethod(method.id)}>
									Set as Default
								</Button>
							</div>
						{/if}
					</Card.Content>
				</Card.Root>
			{/each}

			{#if paymentMethods.length === 0}
				<Card.Root>
					<Card.Content class="text-center py-12">
						<CreditCard class="h-12 w-12 mx-auto text-gray-400 mb-4" />
						<h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">No payment methods</h3>
						<p class="text-gray-600 dark:text-gray-400 mb-4">
							Add a payment method to make checkout faster and easier
						</p>
						<Button on:click={addNewMethod}>
							<Plus class="h-4 w-4 mr-2" />
							Add Your First Payment Method
						</Button>
					</Card.Content>
				</Card.Root>
			{/if}
		{/if}

		<!-- Add/Edit Payment Method Form -->
		{#if showAddForm}
			<Card.Root>
				<Card.Header>
					<Card.Title>
						{editingMethod ? 'Edit Payment Method' : 'Add New Payment Method'}
					</Card.Title>
				</Card.Header>
				<Card.Content>
					<form on:submit={handleSaveMethod} class="space-y-4">
						<div>
							<Label htmlFor="cardNumber">Card Number</Label>
							<div class="relative">
								<Input
									id="cardNumber"
									bind:value={formData.cardNumber}
									on:input={handleCardNumberInput}
									placeholder="1234 5678 9012 3456"
									maxlength={19}
									class={formErrors.cardNumber ? 'border-red-500' : ''}
									required
								/>
								{#if detectCardBrand(formData.cardNumber)}
									<div class="absolute right-3 top-1/2 transform -translate-y-1/2">
										<Badge variant="outline" class="text-xs">
											{detectCardBrand(formData.cardNumber)}
										</Badge>
									</div>
								{/if}
							</div>
							{#if formErrors.cardNumber}
								<p class="text-sm text-red-600 mt-1 flex items-center gap-1">
									<AlertTriangle class="h-3 w-3" />
									{formErrors.cardNumber}
								</p>
							{/if}
						</div>

						<div>
							<Label htmlFor="cardName">Cardholder Name</Label>
							<Input
								id="cardName"
								bind:value={formData.cardName}
								on:blur={() => validateForm()}
								placeholder="Name on card"
								class={formErrors.cardName ? 'border-red-500' : ''}
								required
							/>
							{#if formErrors.cardName}
								<p class="text-sm text-red-600 mt-1 flex items-center gap-1">
									<AlertTriangle class="h-3 w-3" />
									{formErrors.cardName}
								</p>
							{/if}
						</div>

						<div class="grid gap-4 md:grid-cols-2">
							<div>
								<Label htmlFor="expiryDate">Expiry Date</Label>
								<Input
									id="expiryDate"
									bind:value={formData.expiryDate}
									on:input={handleExpiryDateInput}
									placeholder="MM/YY"
									maxlength={5}
									class={formErrors.expiryDate ? 'border-red-500' : ''}
									required
								/>
								{#if formErrors.expiryDate}
									<p class="text-sm text-red-600 mt-1 flex items-center gap-1">
										<AlertTriangle class="h-3 w-3" />
										{formErrors.expiryDate}
									</p>
								{/if}
							</div>
							<div>
								<Label htmlFor="cvv" class="flex items-center gap-2">
									CVV
									<button
										type="button"
										class="text-gray-400 hover:text-gray-600"
										title="3-digit code on the back of your card (4 digits for American Express)"
									>
										<AlertTriangle class="h-3 w-3" />
									</button>
								</Label>
								<Input
									id="cvv"
									bind:value={formData.cvv}
									on:blur={() => validateForm()}
									placeholder={detectCardBrand(formData.cardNumber) === 'American Express' ? '1234' : '123'}
									maxlength={detectCardBrand(formData.cardNumber) === 'American Express' ? 4 : 3}
									class={formErrors.cvv ? 'border-red-500' : ''}
									type="password"
									required
								/>
								{#if formErrors.cvv}
									<p class="text-sm text-red-600 mt-1 flex items-center gap-1">
										<AlertTriangle class="h-3 w-3" />
										{formErrors.cvv}
									</p>
								{/if}
							</div>
						</div>

						<div class="border-t pt-4">
							<h4 class="font-medium mb-3">Billing Address</h4>
							<div class="space-y-3">
								<div>
									<Label htmlFor="billingAddress">Address</Label>
									<Input 
										id="billingAddress" 
										bind:value={formData.billingAddress}
										on:blur={() => validateForm()}
										placeholder="Street address" 
										class={formErrors.billingAddress ? 'border-red-500' : ''}
										required 
									/>
									{#if formErrors.billingAddress}
										<p class="text-sm text-red-600 mt-1 flex items-center gap-1">
											<AlertTriangle class="h-3 w-3" />
											{formErrors.billingAddress}
										</p>
									{/if}
								</div>
								<div class="grid gap-3 md:grid-cols-3">
									<div>
										<Label htmlFor="billingCity">City</Label>
										<Input 
											id="billingCity" 
											bind:value={formData.billingCity}
											on:blur={() => validateForm()}
											placeholder="City" 
											class={formErrors.billingCity ? 'border-red-500' : ''}
											required 
										/>
										{#if formErrors.billingCity}
											<p class="text-sm text-red-600 mt-1 flex items-center gap-1">
												<AlertTriangle class="h-3 w-3" />
												{formErrors.billingCity}
											</p>
										{/if}
									</div>
									<div>
										<Label htmlFor="billingState">State</Label>
										<Input 
											id="billingState" 
											bind:value={formData.billingState}
											on:blur={() => validateForm()}
											placeholder="State" 
											class={formErrors.billingState ? 'border-red-500' : ''}
											required 
										/>
										{#if formErrors.billingState}
											<p class="text-sm text-red-600 mt-1 flex items-center gap-1">
												<AlertTriangle class="h-3 w-3" />
												{formErrors.billingState}
											</p>
										{/if}
									</div>
									<div>
										<Label htmlFor="billingZip">ZIP Code</Label>
										<Input 
											id="billingZip" 
											bind:value={formData.billingZip}
											on:blur={() => validateForm()}
											placeholder="12345" 
											class={formErrors.billingZip ? 'border-red-500' : ''}
											required 
										/>
										{#if formErrors.billingZip}
											<p class="text-sm text-red-600 mt-1 flex items-center gap-1">
												<AlertTriangle class="h-3 w-3" />
												{formErrors.billingZip}
											</p>
										{/if}
									</div>
								</div>
							</div>
						</div>

						<div class="flex gap-2 pt-4">
							<Button type="submit" disabled={isProcessing}>
								{#if isProcessing}
									<div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
								{/if}
								{editingMethod ? 'Update Payment Method' : 'Save Payment Method'}
							</Button>
							<Button type="button" variant="outline" on:click={cancelForm} disabled={isProcessing}>
								Cancel
							</Button>
						</div>
					</form>
				</Card.Content>
			</Card.Root>
		{/if}
	</div>
</div>
