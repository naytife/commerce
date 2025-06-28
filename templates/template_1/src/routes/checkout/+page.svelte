<script lang="ts">
  import { onMount } from 'svelte';
  import { cart } from '$lib/stores/cart';
  import { shop, shopId } from '$lib/stores/shop';
  import { derived, get } from 'svelte/store';
  import { goto } from '$app/navigation';
  import { CreditCard, Zap, DollarSign, Globe, ShoppingBag, Truck, ArrowRight } from 'lucide-svelte';
  import * as Card from '$lib/components/ui/card';
  import { currencySymbol } from '$lib/stores/currency';
  import { fetchPaymentMethods, type PaymentMethodConfig } from '$lib/services/stripe';
  import { apiClient } from '$lib/services/api';
  import StripeCheckout from '$lib/components/StripeCheckout.svelte';
  import PayPalCheckout from '$lib/components/PayPalCheckout.svelte';
  import PaystackCheckout from '$lib/components/PaystackCheckout.svelte';
  import FlutterwaveCheckout from '$lib/components/FlutterwaveCheckout.svelte';
  import { createOrder, type CreateOrderRequest } from '$lib/services/order';

  let name = '';
  let email = '';
  let address = '';
  let city = '';
  let postalCode = '';
  let country = '';
  let paymentProvider = '';

  const total = derived(cart, items =>
    items.reduce((sum, item) => sum + item.price * item.quantity, 0)
  );

  let availablePaymentMethods: PaymentMethodConfig[] = [];
  let paymentOptions: Array<{ value: string; name: string; icon: any; enabled: boolean }> = [];

  const shippingOptions = [
    { value: 'STANDARD', name: 'Standard Shipping', cost: 0, description: '5-7 business days' },
    { value: 'EXPRESS', name: 'Express Shipping', cost: 5, description: '2-3 business days' },
    { value: 'OVERNIGHT', name: 'Overnight Shipping', cost: 20, description: 'Next business day' }
  ];

  // initialize shipping and form state
  let shippingMethod = shippingOptions[0].value;
  let phoneNumber = '';
  let errorMessage = '';
  let shippingCost: number;
  // --- STRIPE ORDER-FIRST FLOW REFACTOR ---
  // Remove auto-initialization of Stripe on mount. Only initialize when user clicks "Complete Order".
  let stripeOrderId: string = '';
  let showStripeCheckout = false;
  let creatingStripeOrder = false;
  $: shippingCost = shippingOptions.find(o => o.value === shippingMethod)?.cost ?? 0;
  $: orderTotal = $total + shippingCost;

  async function loadPaymentMethods() {
    try {
      const currentShopId = get(shopId);
      if (!currentShopId) {
        errorMessage = 'Shop ID is not available. Cannot load payment methods.';
        paymentOptions = [];
        availablePaymentMethods = [];
        paymentProvider = '';
        return;
      }

      availablePaymentMethods = await fetchPaymentMethods(currentShopId);
      console.log('DEBUG: availablePaymentMethods', availablePaymentMethods);
      if (!Array.isArray(availablePaymentMethods)) {
        errorMessage = 'Payment methods data is not an array!';
        return;
      }
      if (availablePaymentMethods.length === 0) {
        errorMessage = 'No payment methods found in shop.json!';
      }
      // Build payment options based on available methods
      paymentOptions = [
        { 
          value: 'paypal', 
          name: 'PayPal', 
          icon: CreditCard, 
          enabled: availablePaymentMethods.some(m => m.provider === 'paypal' && m.enabled) 
        },
        { 
          value: 'stripe', 
          name: 'Stripe', 
          icon: Zap, 
          enabled: availablePaymentMethods.some(m => m.provider === 'stripe' && m.enabled) 
        },
        { 
          value: 'paystack', 
          name: 'Paystack', 
          icon: DollarSign, 
          enabled: availablePaymentMethods.some(m => m.provider === 'paystack' && m.enabled) 
        },
        { 
          value: 'flutterwave', 
          name: 'Flutterwave', 
          icon: Globe, 
          enabled: availablePaymentMethods.some(m => m.provider === 'flutterwave' && m.enabled) 
        }
      ];
      console.log('DEBUG: paymentOptions before filter', paymentOptions);
      paymentOptions = paymentOptions.filter(option => option.enabled);
      console.log('DEBUG: paymentOptions after filter', paymentOptions);
      // Select first available payment method
      if (paymentOptions.length > 0 && !paymentProvider) {
        paymentProvider = paymentOptions[0].value;
      }
    } catch (error) {
      console.error('Failed to load payment methods:', error);
      errorMessage = 'Failed to load payment methods.';
      paymentOptions = [];
      availablePaymentMethods = [];
      paymentProvider = '';
    }
  }

  function handlePaymentSelection(provider: string) {
    paymentProvider = provider;
    // Reset Stripe state if user changes payment method
    showStripeCheckout = false;
    stripeOrderId = '';
  }

  function handleStripeSuccess(paymentIntentId: string) {
    cart.clear();
    goto('/checkout/success');
  }

  function handlePayPalSuccess(transactionId: string) {
    processOrder(transactionId);
  }

  function handlePaystackSuccess(transactionId: string) {
    processOrder(transactionId);
  }

  function handleFlutterwaveSuccess(transactionId: string) {
    processOrder(transactionId);
  }

  async function processOrder(transactionId?: string) {
    try {
      const currentShopId = get(shopId);
      if (!currentShopId) {
        errorMessage = 'Shop ID not found';
        return;
      }

      // Calculate shipping cost
      const shippingOption = shippingOptions.find(o => o.value === shippingMethod);
      const calculatedShippingCost = shippingOption?.cost ?? 0;
      
      // Build checkout request for backend (optimized)
      const checkoutRequest = {
        order_id: stripeOrderId, // use the created order id
        shop_id: currentShopId,
        payment_method_type: paymentProvider
      };

      // 1. Create payment session
      const checkoutResponse = await apiClient.restPost(`/v1/payments/checkout`, checkoutRequest);
      if (!checkoutResponse || !checkoutResponse.payment_intent_id) {
        throw new Error('Failed to create payment session');
      }
      const paymentIntentId = checkoutResponse.payment_intent_id;

      // 2. Confirm payment
      const confirmPayload: any = { payment_intent_id: paymentIntentId };
      if (transactionId) confirmPayload.transaction_id = transactionId;
      const confirmResponse = await apiClient.restPost(`/v1/payments/${currentShopId}/confirm`, confirmPayload);
      if (confirmResponse.status === 'completed' || confirmResponse.status === 'succeeded') {
        cart.clear();
        goto('/checkout/success');
      } else {
        // 3. Optionally poll for status
        let status = confirmResponse.status;
        let attempts = 0;
        while (status !== 'completed' && status !== 'succeeded' && attempts < 5) {
          await new Promise(r => setTimeout(r, 1500));
          const statusResp = await apiClient.restGet(`/v1/payments/${currentShopId}/status/${paymentIntentId}`);
          status = statusResp.status;
          attempts++;
        }
        if (status === 'completed' || status === 'succeeded') {
          cart.clear();
          goto('/checkout/success');
        } else {
          throw new Error('Payment processing failed');
        }
      }
    } catch (error: any) {
      console.error(error);
      errorMessage = (error instanceof Error ? error.message : '') || 'An error occurred while placing your order.';
    }
  }

  async function handlePlaceOrder() {
    errorMessage = '';
    if (!phoneNumber) {
      errorMessage = 'Please enter a phone number.';
      return;
    }
    if (!paymentProvider) {
      errorMessage = 'Please select a payment provider.';
      return;
    }
    if (!shippingMethod) {
      errorMessage = 'Please select a shipping method.';
      return;
    }
    if (paymentProvider === 'stripe') {
      creatingStripeOrder = true;
      try {
        const currentShopId = get(shopId) || '';
        const orderReq = {
          customer_name: name,
          customer_email: email,
          customer_phone: phoneNumber,
          shipping_address: `${address},${city},${postalCode},${country}`,
          shipping_method: shippingMethod,
          payment_method: 'stripe',
          discount: 0,
          shipping_cost: shippingCost,
          tax: 0,
          items: get(cart).map(item => ({
            product_variation_id: item.id,
            quantity: item.quantity,
            price: item.price
          }))
        };
        const orderResp = await createOrder(currentShopId, orderReq);
        if (orderResp && orderResp.order_id) {
          stripeOrderId = orderResp.order_id;
          showStripeCheckout = true;
          console.log('[Checkout] Order created, setting showStripeCheckout:', showStripeCheckout, 'stripeOrderId:', stripeOrderId);
        } else {
          errorMessage = 'Failed to create order for payment.';
          showStripeCheckout = false;
          stripeOrderId = '';
        }
      } catch (err) {
        errorMessage = 'Failed to create order for payment.';
        showStripeCheckout = false;
        stripeOrderId = '';
      } finally {
        creatingStripeOrder = false;
      }
      return;
    }
    // For other providers, process order as before
    await processOrder();
  }

  function handlePayError(error: string) {
    errorMessage = `Payment failed: ${error}`;
  }

  function handleStripeError(error: string) {
    errorMessage = `Payment failed: ${error}`;
  }

  onMount(async () => {
    
    // Initialize shop context
    await shop.initialize();
    
    // Load available payment methods
    await loadPaymentMethods();
  });
</script>

<section class="bg-white py-12 antialiased dark:bg-gray-900">
  <div class="mx-auto max-w-7xl px-6">
    <div class="border-b border-gray-200 dark:border-gray-700 pb-6 mb-8">
      <h1 class="text-2xl font-medium text-gray-900 dark:text-white tracking-tight">Checkout</h1>
    </div>
    
    <div class="flex flex-col lg:flex-row gap-12">
      <div class="w-full lg:w-7/12">
        {#if errorMessage}
          <div class="bg-red-50 border-l-4 border-red-500 p-4 mb-6 dark:bg-red-900/20 dark:border-red-500">
            <p class="text-sm text-red-700 dark:text-red-400">{errorMessage}</p>
          </div>
        {/if}
        
        <div class="border border-gray-200 dark:border-gray-700 p-6 mb-8">
          <h2 class="text-lg font-medium text-gray-900 dark:text-white mb-6">Contact Information</h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Full Name*</label>
              <input 
                id="name"
                type="text" 
                bind:value={name} 
                placeholder="John Doe" 
                class="w-full p-3 border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:border-primary-500 dark:focus:border-primary-500 focus:outline-none" 
                required
              />
            </div>
            <div>
              <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Email Address*</label>
              <input 
                id="email"
                type="email" 
                bind:value={email} 
                placeholder="you@example.com" 
                class="w-full p-3 border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:border-primary-500 dark:focus:border-primary-500 focus:outline-none" 
                required
              />
            </div>
            <div>
              <label for="phone" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Phone Number*</label>
              <input 
                id="phone"
                type="tel" 
                bind:value={phoneNumber} 
                placeholder="123-456-7890" 
                class="w-full p-3 border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:border-primary-500 dark:focus:border-primary-500 focus:outline-none" 
                required
              />
            </div>
          </div>
        </div>
        
        <div class="border border-gray-200 dark:border-gray-700 p-6 mb-8">
          <h2 class="text-lg font-medium text-gray-900 dark:text-white mb-6">Shipping Address</h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="md:col-span-2">
              <label for="address" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Address*</label>
              <input 
                id="address"
                type="text" 
                bind:value={address} 
                placeholder="1234 Main St" 
                class="w-full p-3 border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:border-primary-500 dark:focus:border-primary-500 focus:outline-none" 
                required
              />
            </div>
            <div>
              <label for="city" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">City*</label>
              <input 
                id="city"
                type="text" 
                bind:value={city} 
                placeholder="City" 
                class="w-full p-3 border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:border-primary-500 dark:focus:border-primary-500 focus:outline-none" 
                required
              />
            </div>
            <div>
              <label for="postal" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Postal Code*</label>
              <input 
                id="postal"
                type="text" 
                bind:value={postalCode} 
                placeholder="Postal Code" 
                class="w-full p-3 border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:border-primary-500 dark:focus:border-primary-500 focus:outline-none" 
                required
              />
            </div>
            <div>
              <label for="country" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Country*</label>
              <input 
                id="country"
                type="text" 
                bind:value={country} 
                placeholder="Country" 
                class="w-full p-3 border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:border-primary-500 dark:focus:border-primary-500 focus:outline-none" 
                required
              />
            </div>
          </div>
        </div>
        
        <div class="border border-gray-200 dark:border-gray-700 p-6 mb-8">
          <h2 class="text-lg font-medium text-gray-900 dark:text-white mb-6">Shipping Method</h2>
          <div class="space-y-4">
            {#each shippingOptions as option}
              <label class="flex cursor-pointer border border-gray-200 dark:border-gray-700 p-4 hover:border-gray-300 dark:hover:border-gray-600 transition-colors duration-200 {shippingMethod === option.value ? 'border-primary-700 dark:border-primary-500' : ''}">
                <input 
                  type="radio" 
                  bind:group={shippingMethod} 
                  value={option.value} 
                  class="sr-only" 
                />
                <div class="flex-1">
                  <div class="flex justify-between items-center mb-1">
                    <span class="text-sm font-medium text-gray-900 dark:text-white">{option.name}</span>
                    <span class="text-sm font-medium text-gray-900 dark:text-white">
                      {option.cost === 0 ? 'Free' : `${$currencySymbol}${option.cost.toFixed(2)}`}
                    </span>
                  </div>
                  <p class="text-sm text-gray-500 dark:text-gray-400">{option.description}</p>
                </div>
                <Truck class="w-5 h-5 ml-4 text-gray-400 {shippingMethod === option.value ? 'text-primary-700 dark:text-primary-500' : ''}" />
              </label>
            {/each}
          </div>
        </div>
        
        <div class="border border-gray-200 dark:border-gray-700 p-6">
          <h2 class="text-lg font-medium text-gray-900 dark:text-white mb-6">Payment Method</h2>
          
          {#if paymentOptions.length === 0}
            <div class="bg-yellow-50 border-l-4 border-yellow-500 p-4 dark:bg-yellow-900/20 dark:border-yellow-500">
              <p class="text-sm text-yellow-700 dark:text-yellow-400">
                No payment methods are currently configured for this shop.
              </p>
            </div>
          {:else}
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-6">
              {#each paymentOptions as option}
                <label class="flex items-center border border-gray-200 dark:border-gray-700 p-4 cursor-pointer hover:border-gray-300 dark:hover:border-gray-600 transition-colors duration-200 {paymentProvider === option.value ? 'border-primary-700 dark:border-primary-500' : ''}">
                  <input 
                    type="radio" 
                    bind:group={paymentProvider} 
                    value={option.value}
                    on:change={() => handlePaymentSelection(option.value)}
                    class="sr-only" 
                  />
                  <svelte:component this={option.icon} class="w-5 h-5 text-gray-400 {paymentProvider === option.value ? 'text-primary-700 dark:text-primary-500' : ''}" />
                  <span class="ml-3 text-sm font-medium text-gray-900 dark:text-white">{option.name}</span>
                </label>
              {/each}
            </div>
          {/if}
        </div>
      </div>
      
      <div class="w-full lg:w-5/12 mt-8 lg:mt-0">
        <div class="border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 sticky top-6">
          <div class="p-6 border-b border-gray-200 dark:border-gray-700">
            <h2 class="text-lg font-medium text-gray-900 dark:text-white mb-6">Order Summary</h2>
            <div class="max-h-80 overflow-y-auto pr-2">
              {#each $cart as item (item.id)}
                <div class="flex py-4 border-b border-gray-200 dark:border-gray-700 last:border-b-0">
                  <div class="w-16 h-16 bg-gray-100 dark:bg-gray-700 mr-4 flex-shrink-0">
                    {#if item.image}
                      <img src={item.image} alt={item.title} class="w-full h-full object-cover" />
                    {/if}
                  </div>
                  <div class="flex-1">
                    <h3 class="text-sm font-medium text-gray-900 dark:text-white">{item.title}</h3>
                    <p class="text-sm text-gray-500 dark:text-gray-400">Qty: {item.quantity}</p>
                  </div>
                  <div class="ml-4 text-right">
                    <p class="text-sm font-medium text-gray-900 dark:text-white">{$currencySymbol}{(item.price * item.quantity).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</p>
                  </div>
                </div>
              {/each}
            </div>
          </div>
          
          <div class="p-6">
            <div class="space-y-3 mb-4">
              <div class="flex justify-between text-sm">
                <span class="text-gray-600 dark:text-gray-400">Subtotal</span>
                <span class="text-gray-900 dark:text-white font-medium">{$currencySymbol}{$total.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
              </div>
              <div class="flex justify-between text-sm">
                <span class="text-gray-600 dark:text-gray-400">Shipping</span>
                <span class="text-gray-900 dark:text-white font-medium">{shippingCost === 0 ? 'Free' : `${$currencySymbol}${shippingCost.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`}</span>
              </div>
              <div class="pt-3 mt-3 border-t border-gray-200 dark:border-gray-700 flex justify-between">
                <span class="text-base font-medium text-gray-900 dark:text-white">Total</span>
                <span class="text-base font-semibold text-gray-900 dark:text-white">{$currencySymbol}{orderTotal.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
              </div>
            </div>
            
            <!-- --- UI: Only show StripeCheckout after order is created and showStripeCheckout is true --- -->
            {#if showStripeCheckout && paymentProvider === 'stripe'}
              <div class="mt-6 p-4 border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800">
                <h3 class="text-md font-medium text-gray-900 dark:text-white mb-4">Complete Payment</h3>
                <StripeCheckout 
                  amount={orderTotal}
                  currency={$shop?.currency_code?.toLowerCase() || 'usd'}
                  orderId={stripeOrderId}
                  metadata={{
                    customer_name: name,
                    customer_email: email,
                    customer_phone: phoneNumber,
                    shipping_method: shippingMethod,
                    shipping_address: `${address},${city},${postalCode},${country}`,
                    cart_items: JSON.stringify($cart.map(item => ({ id: item.id, title: item.title, quantity: item.quantity })))
                  }}
                  onSuccess={handleStripeSuccess}
                  onError={handleStripeError}
                />
              </div>
            {:else}
              <!-- Show the Complete Order button for all providers, including Stripe (before order is created) -->
              <button 
                on:click={handlePlaceOrder} 
                disabled={!paymentProvider || paymentOptions.length === 0 || creatingStripeOrder}
                class="w-full bg-primary-700 hover:bg-primary-800 disabled:bg-gray-400 disabled:cursor-not-allowed text-white flex items-center justify-center py-4 px-6 font-medium uppercase tracking-wider transition-colors duration-200 focus:ring-4 focus:ring-primary-300 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
              >
                {#if creatingStripeOrder}
                  <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                  Creating Order...
                {:else}
                  {paymentProvider === 'stripe' ? 'Continue to Payment' : 'Complete Order'}
                {/if}
                <ArrowRight class="w-4 h-4 ml-2" />
              </button>
            {/if}
          </div>
        </div>
      </div>
    </div>
  </div>
</section>