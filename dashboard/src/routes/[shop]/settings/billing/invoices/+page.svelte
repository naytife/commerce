<script lang="ts">
  import * as Card from '$lib/components/ui/card';
  import { Button } from '$lib/components/ui/button';
  import { onMount } from 'svelte';
  import { writable } from 'svelte/store';
	import ComingSoon from '$lib/components/ui/coming-soon.svelte';

  let loading = false;
  let error = '';
  let invoices = writable([
    { id: 'INV-001', date: '2024-06-01', amount: 49.99, status: 'Paid', url: '#' },
    { id: 'INV-002', date: '2024-05-01', amount: 49.99, status: 'Paid', url: '#' },
    { id: 'INV-003', date: '2024-04-01', amount: 49.99, status: 'Unpaid', url: '#' },
  ]);

  onMount(async () => {
    loading = true;
    // TODO: Fetch invoices from `/api/shops/{shop_id}/invoices` (if implemented)
    loading = false;
  });
</script>
<ComingSoon>
<div class="space-y-6 w-full px-4 md:px-8 mx-auto">
  <h2 class="text-2xl font-bold">Invoice History</h2>
  <Card.Root class="w-full bg-gradient">
    <Card.Header>
      <Card.Title>Invoices</Card.Title>
      <Card.Description>View and download past payment records</Card.Description>
    </Card.Header>
    <Card.Content>
      <div class="border rounded-md overflow-x-auto">
        <table class="w-full min-w-[600px] md:min-w-[900px] text-sm">
          <thead class="bg-muted">
            <tr>
              <th class="text-left p-2">Invoice ID</th>
              <th class="text-left p-2">Date</th>
              <th class="text-right p-2">Amount</th>
              <th class="text-right p-2">Status</th>
              <th class="text-right p-2">Download</th>
            </tr>
          </thead>
          <tbody>
            {#if $invoices.length === 0}
              <tr class="border-t">
                <td colspan="5" class="p-4 text-center text-muted-foreground">
                  No invoices available
                </td>
              </tr>
            {:else}
              {#each $invoices as invoice}
                <tr class="border-t">
                  <td class="p-2">{invoice.id}</td>
                  <td class="p-2">{invoice.date}</td>
                  <td class="p-2 text-right">${invoice.amount.toFixed(2)}</td>
                  <td class="p-2 text-right">{invoice.status}</td>
                  <td class="p-2 text-right">
                    <Button href={invoice.url} size="sm" variant="outline" download>Download</Button>
                  </td>
                </tr>
              {/each}
            {/if}
          </tbody>
        </table>
      </div>
    </Card.Content>
  </Card.Root>
</div>
</ComingSoon>