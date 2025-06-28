<script lang="ts">
  export let amount: number;
  export let currency: string = 'usd';
  export let config: any = {};
  export let onSuccess: (transactionId: string) => void = () => {};
  export let onError: (error: string) => void = () => {};

  // TODO: Implement PayPal JS SDK integration using config.client_id
  // For now, just a placeholder button
  function handlePay() {
    // Simulate success
    onSuccess('paypal_txn_' + Date.now());
  }
</script>

<div class="space-y-4">
  <button class="w-full bg-blue-600 text-white py-3 rounded" on:click={handlePay}>
    Pay with PayPal (placeholder)
  </button>
  <p class="text-xs text-gray-500">Production: Integrate PayPal JS SDK here.</p>
</div>
