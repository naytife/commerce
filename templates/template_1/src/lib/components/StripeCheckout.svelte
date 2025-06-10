<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { StripeService, createPaymentIntent, fetchPaymentMethods, type PaymentMethodConfig } from '$lib/services/stripe';
  import { shopId } from '$lib/stores/shop';
  import { get } from 'svelte/store';

  export let amount: number;
  export let currency: string = 'usd';
  export let orderId: string = '';
  export let customerId: string = '';
  export let metadata: Record<string, string> = {};
  export let onSuccess: (paymentIntentId: string) => void = () => {};
  export let onError: (error: string) => void = () => {};

  let stripeService: StripeService | null = null;
  let paymentElement: HTMLDivElement;
  let loading = true;
  let processing = false;
  let errorMessage = '';
  let paymentMethods: PaymentMethodConfig[] = [];
  let stripeEnabled = false;

  onMount(async () => {
    await initializeStripe();
  });

  onDestroy(() => {
    if (stripeService) {
      stripeService.destroy();
    }
  });

  async function initializeStripe() {
    try {
      loading = true;
      errorMessage = '';
      
      const currentShopId = get(shopId);
      if (!currentShopId) {
        errorMessage = 'Shop not found';
        loading = false;
        return;
      }

      // Fetch payment methods to get Stripe configuration
      paymentMethods = await fetchPaymentMethods(currentShopId);
      const stripeMethod = paymentMethods.find(
        (method) => method.provider === 'stripe' && method.enabled
      );

      if (!stripeMethod?.config.publishable_key) {
        errorMessage = 'Stripe is not configured for this shop';
        loading = false;
        return;
      }

      stripeEnabled = true;

      // Create payment intent
      const paymentIntent = await createPaymentIntent({
        shopId: currentShopId,
        amount,
        currency,
        orderId,
        customerId,
        metadata: {
          frontend_version: '1.0.0',
          checkout_session: Date.now().toString(),
          ...metadata,
        },
      });
      if (!paymentIntent) {
        errorMessage = 'Failed to create payment intent';
        loading = false;
        return;
      }

      // Initialize Stripe
      stripeService = new StripeService();
      const initialized = await stripeService.initialize(
        stripeMethod.config.publishable_key,
        paymentIntent.client_secret
      );

      if (!initialized) {
        errorMessage = 'Failed to initialize Stripe';
        loading = false;
        return;
      }

      // Mount payment element
      if (paymentElement) {
        stripeService.createPaymentElement(paymentElement);
      }

      loading = false;
    } catch (error) {
      console.error('Error initializing Stripe:', error);
      errorMessage = 'Failed to initialize payment system';
      loading = false;
    }
  }

  async function handleSubmit() {
    if (!stripeService || processing) return;

    try {
      processing = true;
      errorMessage = '';

      const returnUrl = `${window.location.origin}/checkout/success`;
      const result = await stripeService.confirmPayment(returnUrl);

      if (result.success) {
        onSuccess('payment_confirmed');
      } else {
        errorMessage = result.error || 'Payment failed';
        onError(errorMessage);
      }
    } catch (error) {
      console.error('Payment error:', error);
      errorMessage = 'Payment processing failed';
      onError(errorMessage);
    } finally {
      processing = false;
    }
  }
</script>

{#if !stripeEnabled}
  <div class="bg-yellow-50 border-l-4 border-yellow-500 p-4 mb-6 dark:bg-yellow-900/20 dark:border-yellow-500">
    <p class="text-sm text-yellow-700 dark:text-yellow-400">
      Stripe payment processing is not available for this shop.
    </p>
  </div>
{:else if loading}
  <div class="flex items-center justify-center p-8">
    <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-700"></div>
    <span class="ml-3 text-gray-600 dark:text-gray-400">Loading payment form...</span>
  </div>
{:else if errorMessage}
  <div class="bg-red-50 border-l-4 border-red-500 p-4 mb-6 dark:bg-red-900/20 dark:border-red-500">
    <p class="text-sm text-red-700 dark:text-red-400">{errorMessage}</p>
    <button 
      on:click={initializeStripe} 
      class="mt-2 text-sm text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300 underline"
    >
      Try again
    </button>
  </div>
{:else}
  <div class="space-y-6">
    <!-- Stripe Payment Element -->
    <div bind:this={paymentElement} class="min-h-[200px]"></div>
    
    <!-- Submit Button -->
    <button
      on:click={handleSubmit}
      disabled={processing}
      class="w-full bg-primary-700 hover:bg-primary-800 disabled:bg-gray-400 text-white flex items-center justify-center py-3 px-6 font-medium transition-colors duration-200 focus:ring-4 focus:ring-primary-300 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
    >
      {#if processing}
        <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
        Processing...
      {:else}
        Pay ${(amount).toFixed(2)}
      {/if}
    </button>
    
    <!-- Powered by Stripe -->
    <div class="text-center">
      <p class="text-xs text-gray-500 dark:text-gray-400">
        Powered by <span class="font-medium">Stripe</span>
      </p>
    </div>
  </div>
{/if}
