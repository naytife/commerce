<script lang="ts">
  export let amount: number;
  export let currency: string = 'usd';
  export let config: any = {};
  export let onSuccess: (transactionId: string) => void = () => {};
  export let onError: (error: string) => void = () => {};

  // TODO: Implement Paystack JS SDK integration using config.public_key
  // For now, just a placeholder button
  function handlePay() {
    // Simulate success
    onSuccess('paystack_txn_' + Date.now());
  }
</script>

<div class="space-y-4">
  <button class="w-full bg-green-600 text-white py-3 rounded" on:click={handlePay}>
    Pay with Paystack (placeholder)
  </button>
  <p class="text-xs text-gray-500">Production: Integrate Paystack JS SDK here.</p>
</div>
