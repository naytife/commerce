<script lang="ts">
  export let amount: number;
  export let currency: string = 'usd';
  export let config: any = {};
  export let onSuccess: (transactionId: string) => void = () => {};
  export let onError: (error: string) => void = () => {};

  // TODO: Implement Flutterwave JS SDK integration using config.public_key_flutterwave
  // For now, just a placeholder button
  function handlePay() {
    // Simulate success
    onSuccess('flutterwave_txn_' + Date.now());
  }
</script>

<div class="space-y-4">
  <button class="w-full bg-yellow-500 text-white py-3 rounded" on:click={handlePay}>
    Pay with Flutterwave (placeholder)
  </button>
  <p class="text-xs text-gray-500">Production: Integrate Flutterwave JS SDK here.</p>
</div>
