<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { tick } from 'svelte';
  import { StripeService, createPaymentIntent, fetchPaymentMethods, type PaymentMethodConfig } from '$lib/services/stripe';
  import { shopId } from '$lib/stores/shop';
  import { get } from 'svelte/store';

  export let amount: number;
  export let currency: string = 'usd';
  export let orderId: string = '';
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
  let paymentElementMounted = false;

  onDestroy(() => {
    if (stripeService) {
      stripeService.destroy();
    }
  });

  // Only call initializeStripe() when the component is mounted and orderId is present
  onMount(() => {
    if (orderId) {
      initializeStripe();
    }
  });

  $: if (orderId && !stripeService && !loading) {
    initializeStripe();
  }

  // Only fetch payment methods and set stripeEnabled when initializing Stripe, not on every prop change
  let paymentMethodsLoaded = false;

  // Only show the Stripe error if we've actually tried to initialize Stripe (not on first render)
  let triedToInitStripe = false;

  async function initializeStripe() {
    triedToInitStripe = true;
    try {
      loading = true;
      errorMessage = '';
      paymentElementMounted = false;
      paymentMethodsLoaded = false;
      const currentShopId = get(shopId);
      console.log('[StripeCheckout] initializeStripe called', { currentShopId, orderId, amount, currency, metadata });
      if (!currentShopId) {
        errorMessage = 'Shop not found';
        loading = false;
        return;
      }
      paymentMethods = await fetchPaymentMethods(currentShopId);
      paymentMethodsLoaded = true;
      const stripeMethod = paymentMethods.find(
        (method) => method.provider === 'stripe' && method.enabled
      );
      console.log('[StripeCheckout] stripeMethod', stripeMethod);
      if (!stripeMethod?.config.publishable_key) {
        errorMessage = 'Stripe is not configured for this shop';
        stripeEnabled = false;
        loading = false;
        return;
      }
      stripeEnabled = true;
      // Call createPaymentIntent with minimal payload
      if (!orderId) {
        errorMessage = 'Order ID is required for payment.';
        loading = false;
        return;
      }
      // Call createPaymentIntent with minimal payload
      const paymentIntent = await createPaymentIntent({
        shopId: currentShopId,
        orderId,
        paymentMethodType: 'stripe'
      });
      console.log('[StripeCheckout] createPaymentIntent result', paymentIntent);
      if (!paymentIntent) {
        errorMessage = 'Failed to create payment intent';
        loading = false;
        return;
      }
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
      // Do not mount payment element here; let the reactive block handle it
      loading = false;
    } catch (error) {
      errorMessage = 'Failed to initialize payment system';
      loading = false;
    }
  }

  // Helper to extract payment_intent_id from client_secret
  function getPaymentIntentIdFromClientSecret(secret: string | null): string {
    if (!secret) return '';
    const idx = secret.indexOf('_secret');
    return idx > 0 ? secret.substring(0, idx) : secret;
  }

  async function handleSubmit() {
    if (!stripeService || processing) return;

    try {
      processing = true;
      errorMessage = '';

      const currentShopId = get(shopId);
      if (!orderId) {
        errorMessage = 'Order ID missing for payment confirmation.';
        processing = false;
        return;
      }
      // Confirm payment with Stripe Elements
      const returnUrl = `${window.location.origin}/checkout/success`;
      const result = await stripeService.confirmPayment(returnUrl);

      if (result.success) {
        // Call backend to confirm payment and update order status
        try {
          const paymentIntentId = getPaymentIntentIdFromClientSecret(stripeService['clientSecret']);
          const confirmResp = await fetch(`/api/v1/payments/${currentShopId}/confirm`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              payment_intent_id: paymentIntentId,
              order_id: orderId
            })
          });
          const confirmData = await confirmResp.json();
          if (confirmResp.ok && (confirmData.status === 'succeeded' || confirmData.status === 'completed')) {
            onSuccess(paymentIntentId);
          } else {
            errorMessage = confirmData.message || 'Payment confirmation failed.';
            onError(errorMessage);
          }
        } catch (err) {
          errorMessage = 'Failed to confirm payment with backend.';
          onError(errorMessage);
        }
      } else {
        errorMessage = result.error || 'Payment failed';
        onError(errorMessage);
      }
    } catch (error) {
      errorMessage = 'Payment processing failed';
      onError(errorMessage);
    } finally {
      processing = false;
    }
  }

  $: if (stripeService && !loading && !paymentElementMounted && paymentElement) {
    try {
      stripeService.createPaymentElement(paymentElement);
      paymentElementMounted = true;
    } catch (e) {
      errorMessage = 'Failed to mount Stripe payment element';
    }
  }
</script>

{#if !stripeEnabled && triedToInitStripe}
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
