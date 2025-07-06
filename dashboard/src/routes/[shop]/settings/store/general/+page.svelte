<script lang="ts">
  import * as Card from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';
  import { Textarea } from '$lib/components/ui/textarea';
  import { Button } from '$lib/components/ui/button';
  import { Label } from '$lib/components/ui/label';
  import { toast } from 'svelte-sonner';
  import type { Shop } from '$lib/types';
  import { api } from '$lib/api';
  import { getContext } from 'svelte';
  import { deepEqual, deepClone } from '$lib/utils/deepEqual';
  import type { Writable } from 'svelte/store';

  // --- Breadcrumbs ---
  const breadcrumbs = [
    { label: 'Settings' },
    { label: 'Store' },
    { label: 'General' }
  ];

  const shopStore = getContext('shop') as Writable<Partial<Shop>>;
  const authFetch = getContext('authFetch') as typeof fetch;
  const refetchShopData = getContext('refetchShopData') as () => Promise<void>;

  let shop: Partial<Shop> = {};
  $: shop = $shopStore;

  let form: Partial<Shop> = {};
  let initialShop: Partial<Shop> = {};
  let initialized = false;

  $: if (!initialized && shop && Object.keys(shop).length > 0) {
    form = deepClone(shop);
    initialShop = deepClone(shop);
    initialized = true;
  }

  let hasChanges = false;
  $: hasChanges = Object.keys(form).length > 0 && !deepEqual(form, initialShop);

  async function handleFormSubmit(event: Event) {
    event.preventDefault();
    try {
      await api(authFetch).updateShop(form);
      await refetchShopData();
      initialShop = deepClone(form);
      toast.success('Store info updated!');
    } catch (error) {
      toast.error('Failed to update store info');
    }
  }
</script>

<!-- Breadcrumbs UI (not navigable) -->
<nav class="mb-6 text-sm text-muted-foreground flex items-center gap-2">
  {#each breadcrumbs as crumb, i}
    <span>{crumb.label}</span>{#if i < breadcrumbs.length - 1}<span>/</span>{/if}
  {/each}
</nav>

<Card.Root class="bg-gradient w-full max-w-screen-md mx-auto">
  <Card.Header>
    <Card.Title>Store Info</Card.Title>
    <Card.Description>Set your store name, contact info, and basic preferences.</Card.Description>
  </Card.Header>
  <Card.Content>
    {#if !shop || Object.keys(shop).length === 0}
      <p>Loading...</p>
    {:else}
      <form class="space-y-4 w-full" on:submit={handleFormSubmit}>
        <div>
          <Label for="store-name">Store Name</Label>
          <Input id="store-name" bind:value={form.title} required />
        </div>
        <div>
          <Label for="email">Email</Label>
          <Input id="email" type="email" bind:value={form.email} />
        </div>
        <div>
          <Label for="phone">Phone</Label>
          <Input id="phone" bind:value={form.phone_number} />
        </div>
        <div>
          <Label for="address">Address</Label>
          <Input id="address" bind:value={form.address} />
        </div>
        <div>
          <Label for="about">About</Label>
          <Textarea id="about" bind:value={form.about} />
        </div>
        <div>
          <Label for="currency">Currency</Label>
          <Input id="currency" bind:value={form.currency_code} />
        </div>
        <Button type="submit" disabled={!hasChanges}>Update Store Info</Button>
      </form>
    {/if}
  </Card.Content>
</Card.Root> 