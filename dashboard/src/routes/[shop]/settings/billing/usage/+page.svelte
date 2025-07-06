<script lang="ts">
  import * as Card from '$lib/components/ui/card';
	import ComingSoon from '$lib/components/ui/coming-soon.svelte';
  import { onMount } from 'svelte';
  import { writable } from 'svelte/store';

  let loading = false;
  let error = '';
  let usage = writable({
    apiCalls: 1234,
    apiLimit: 5000,
    mediaStorage: 2.5, // GB
    mediaLimit: 10, // GB
  });

  $: apiPercent = Math.min(100, Math.round(($usage.apiCalls / $usage.apiLimit) * 100));
  $: mediaPercent = Math.min(100, Math.round(($usage.mediaStorage / $usage.mediaLimit) * 100));

  onMount(async () => {
    loading = true;
    // TODO: Fetch usage from `/api/shops/{shop_id}/usage` (if implemented)
    loading = false;
  });
</script>
<ComingSoon>
<div class="space-y-6 w-full px-4 md:px-8 mx-auto">
  <h2 class="text-2xl font-bold">Usage Overview</h2>
  <Card.Root class="w-full">
    <Card.Header>
      <Card.Title>Resource Usage</Card.Title>
      <Card.Description>Track API, media, and other usage limits</Card.Description>
    </Card.Header>
    <Card.Content>
      <div class="space-y-6">
        <div>
          <div class="flex justify-between mb-1">
            <span class="font-medium">API Calls</span>
            <span class="text-sm text-muted-foreground">{$usage.apiCalls} / {$usage.apiLimit} ({apiPercent}%)</span>
          </div>
          <div class="w-full bg-muted rounded h-3 overflow-hidden">
            <div class="bg-primary h-3 rounded" style="width: {apiPercent}%; transition: width 0.3s;"></div>
          </div>
        </div>
        <div>
          <div class="flex justify-between mb-1">
            <span class="font-medium">Media Storage</span>
            <span class="text-sm text-muted-foreground">{$usage.mediaStorage} GB / {$usage.mediaLimit} GB ({mediaPercent}%)</span>
          </div>
          <div class="w-full bg-muted rounded h-3 overflow-hidden">
            <div class="bg-primary h-3 rounded" style="width: {mediaPercent}%; transition: width 0.3s;"></div>
          </div>
        </div>
      </div>
    </Card.Content>
  </Card.Root>
</div>
</ComingSoon>