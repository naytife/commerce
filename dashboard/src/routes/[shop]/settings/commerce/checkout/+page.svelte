<script lang="ts">
  import * as Card from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';
  import { Button } from '$lib/components/ui/button';
  import { Label } from '$lib/components/ui/label';
  import { Switch } from '$lib/components/ui/switch';
  import { Textarea } from '$lib/components/ui/textarea';
  import { toast } from 'svelte-sonner';
  import type { Shop } from '$lib/types';
  import { api } from '$lib/api';
  import { getContext } from 'svelte';
  import { deepEqual, deepClone } from '$lib/utils/deepEqual';
  import type { Writable } from 'svelte/store';

  const shopStore = getContext('shop') as Writable<Partial<Shop>>;
  const authFetch = getContext('authFetch') as typeof fetch;
  const refetchShopData = getContext('refetchShopData') as () => Promise<void>;

  let shop: Partial<Shop> = {};
  $: shop = $shopStore;

  let form: {
    guest_checkout: boolean;
    allow_notes: boolean;
    terms_required: boolean;
    terms_text?: string;
  } = {
    guest_checkout: true,
    allow_notes: true,
    terms_required: true,
  };
  
  let initialForm = deepClone(form);
  let initialized = false;

  $: if (!initialized && shop && Object.keys(shop).length > 0) {
    // TODO: Load actual checkout settings from shop data
    form = {
      guest_checkout: true,
      allow_notes: true,
      terms_required: true,
    };
    initialForm = deepClone(form);
    initialized = true;
  }

  let hasChanges = false;
  $: hasChanges = Object.keys(form).length > 0 && !deepEqual(form, initialForm);

  async function handleFormSubmit(event: Event) {
    event.preventDefault();
    try {
      // TODO: Implement actual API call
      // await api(authFetch).updateCheckoutSettings(form);
      await refetchShopData();
      initialForm = deepClone(form);
      toast.success('Checkout settings updated!');
    } catch (error) {
      toast.error('Failed to update checkout settings');
    }
  }
</script>

<Card.Root class="bg-gradient w-full max-w-screen-md mx-auto">
  <Card.Header>
    <Card.Title>Checkout Configuration</Card.Title>
    <Card.Description>Configure checkout flow, requirements, and customer experience settings.</Card.Description>
  </Card.Header>
  <Card.Content>
    {#if !shop || Object.keys(shop).length === 0}
      <p>Loading...</p>
    {:else}
      <form class="space-y-6 w-full" on:submit={handleFormSubmit}>
        <!-- Customer Account Settings -->
        <div class="space-y-4">
          <h3 class="text-lg font-semibold">Customer Account</h3>
          
          <div class="flex items-center justify-between">
            <div class="space-y-0.5">
              <Label>Guest Checkout</Label>
              <p class="text-muted-foreground text-sm">Allow customers to checkout without creating an account</p>
            </div>
            <Switch bind:checked={form.guest_checkout} />
          </div>
        </div>

        <!-- Order Options -->
        <div class="space-y-4">
          <h3 class="text-lg font-semibold">Order Options</h3>
          
          <div class="flex items-center justify-between">
            <div class="space-y-0.5">
              <Label>Allow Order Notes</Label>
              <p class="text-muted-foreground text-sm">Let customers add notes to their orders</p>
            </div>
            <Switch bind:checked={form.allow_notes} />
          </div>

          <div class="flex items-center justify-between">
            <div class="space-y-0.5">
              <Label>Require Terms Acceptance</Label>
              <p class="text-muted-foreground text-sm">Customers must accept terms and conditions</p>
            </div>
            <Switch bind:checked={form.terms_required} />
          </div>

          {#if form.terms_required}
            <div>
              <Label for="terms-text">Terms & Conditions Text</Label>
              <Textarea 
                id="terms-text" 
                bind:value={form.terms_text} 
                placeholder="Enter your terms and conditions..."
                rows={4}
              />
              <p class="text-muted-foreground text-sm">Text that customers must accept</p>
            </div>
          {/if}
        </div>
        <Button type="submit" disabled={!hasChanges} class="w-full">
          Update Checkout Settings
        </Button>
      </form>
    {/if}
  </Card.Content>
</Card.Root>
