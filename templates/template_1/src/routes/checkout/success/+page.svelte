<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { CheckCircle, ShoppingBag, ArrowRight } from 'lucide-svelte';

  let paymentStatus = '';
  let orderDetails: any = null;

  onMount(() => {
    // Check URL parameters for payment status
    const urlParams = new URLSearchParams(window.location.search);
    paymentStatus = urlParams.get('payment_intent_status') || 'succeeded';
    
    // TODO: Fetch order details from backend based on payment intent or order ID
    // For now, we'll show a generic success message
  });
</script>

<svelte:head>
  <title>Order Confirmation - Thank You!</title>
</svelte:head>

<section class="bg-white py-12 antialiased dark:bg-gray-900 min-h-screen">
  <div class="mx-auto max-w-3xl px-6">
    <div class="text-center">
      {#if paymentStatus === 'succeeded'}
        <!-- Success State -->
        <div class="flex justify-center mb-8">
          <CheckCircle class="w-24 h-24 text-green-500" />
        </div>
        
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-4">
          Thank you for your order!
        </h1>
        
        <p class="text-lg text-gray-600 dark:text-gray-400 mb-8">
          Your payment has been successfully processed and your order is being prepared.
        </p>
        
        <div class="bg-green-50 border border-green-200 p-6 mb-8 dark:bg-green-900/20 dark:border-green-800">
          <h2 class="text-lg font-semibold text-green-800 dark:text-green-400 mb-2">
            Order Confirmed
          </h2>
          <p class="text-green-700 dark:text-green-300">
            You will receive an email confirmation shortly with your order details and tracking information.
          </p>
        </div>
        
      {:else if paymentStatus === 'processing'}
        <!-- Processing State -->
        <div class="flex justify-center mb-8">
          <div class="animate-spin rounded-full h-24 w-24 border-b-4 border-primary-700"></div>
        </div>
        
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-4">
          Processing your payment...
        </h1>
        
        <p class="text-lg text-gray-600 dark:text-gray-400 mb-8">
          Please wait while we confirm your payment. This may take a few moments.
        </p>
        
      {:else}
        <!-- Error State -->
        <div class="flex justify-center mb-8">
          <div class="w-24 h-24 rounded-full bg-red-100 dark:bg-red-900/20 flex items-center justify-center">
            <span class="text-red-500 text-4xl">⚠️</span>
          </div>
        </div>
        
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-4">
          Payment Issue
        </h1>
        
        <p class="text-lg text-gray-600 dark:text-gray-400 mb-8">
          There was an issue processing your payment. Please try again or contact support.
        </p>
        
        <div class="bg-red-50 border border-red-200 p-6 mb-8 dark:bg-red-900/20 dark:border-red-800">
          <h2 class="text-lg font-semibold text-red-800 dark:text-red-400 mb-2">
            Payment Failed
          </h2>
          <p class="text-red-700 dark:text-red-300">
            Your order was not completed. Please return to checkout and try again.
          </p>
        </div>
      {/if}
      
      <!-- Action Buttons -->
      <div class="space-y-4 sm:space-y-0 sm:space-x-4 sm:flex sm:justify-center">
        <button 
          on:click={() => goto('/')}
          class="w-full sm:w-auto bg-primary-700 hover:bg-primary-800 text-white flex items-center justify-center py-3 px-6 font-medium transition-colors duration-200 focus:ring-4 focus:ring-primary-300 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
        >
          Continue Shopping
          <ArrowRight class="w-4 h-4 ml-2" />
        </button>
        
        {#if paymentStatus !== 'succeeded'}
          <button 
            on:click={() => goto('/checkout')}
            class="w-full sm:w-auto border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 flex items-center justify-center py-3 px-6 font-medium transition-colors duration-200"
          >
            <ShoppingBag class="w-4 h-4 mr-2" />
            Return to Checkout
          </button>
        {/if}
      </div>
      
      <!-- Additional Information -->
      <div class="mt-12 text-left">
        <div class="border-t border-gray-200 dark:border-gray-700 pt-8">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
            What happens next?
          </h3>
          
          <div class="space-y-4">
            <div class="flex items-start">
              <div class="flex-shrink-0 w-6 h-6 bg-primary-100 dark:bg-primary-900/20 rounded-full flex items-center justify-center mr-3 mt-0.5">
                <span class="text-primary-700 dark:text-primary-400 text-sm font-semibold">1</span>
              </div>
              <div>
                <h4 class="font-medium text-gray-900 dark:text-white">Order Confirmation</h4>
                <p class="text-gray-600 dark:text-gray-400 text-sm">
                  You'll receive an email confirmation with your order details.
                </p>
              </div>
            </div>
            
            <div class="flex items-start">
              <div class="flex-shrink-0 w-6 h-6 bg-primary-100 dark:bg-primary-900/20 rounded-full flex items-center justify-center mr-3 mt-0.5">
                <span class="text-primary-700 dark:text-primary-400 text-sm font-semibold">2</span>
              </div>
              <div>
                <h4 class="font-medium text-gray-900 dark:text-white">Order Processing</h4>
                <p class="text-gray-600 dark:text-gray-400 text-sm">
                  We'll prepare your items for shipment within 1-2 business days.
                </p>
              </div>
            </div>
            
            <div class="flex items-start">
              <div class="flex-shrink-0 w-6 h-6 bg-primary-100 dark:bg-primary-900/20 rounded-full flex items-center justify-center mr-3 mt-0.5">
                <span class="text-primary-700 dark:text-primary-400 text-sm font-semibold">3</span>
              </div>
              <div>
                <h4 class="font-medium text-gray-900 dark:text-white">Shipping & Tracking</h4>
                <p class="text-gray-600 dark:text-gray-400 text-sm">
                  You'll receive tracking information once your order ships.
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</section>
