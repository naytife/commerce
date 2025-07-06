<script lang="ts">
  import * as Card from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';
  import { Button } from '$lib/components/ui/button';
  import { Label } from '$lib/components/ui/label';
  import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '$lib/components/ui/select';
  import { Switch } from '$lib/components/ui/switch';
  import { toast } from 'svelte-sonner';
  import type { Shop } from '$lib/types';
  import { api } from '$lib/api';
  import { getContext } from 'svelte';
  import { deepEqual, deepClone } from '$lib/utils/deepEqual';
  import type { Writable } from 'svelte/store';
	import ComingSoon from '$lib/components/ui/coming-soon.svelte';

  const shopStore = getContext('shop') as Writable<Partial<Shop>>;
  const authFetch = getContext('authFetch') as typeof fetch;
  const refetchShopData = getContext('refetchShopData') as () => Promise<void>;

  let shop: Partial<Shop> = {};
  $: shop = $shopStore;

  let form: {
    tax_enabled: boolean;
    tax_type: string;
    tax_rate: number;
    tax_inclusive: boolean;
    tax_exempt_products: string[];
    tax_registration_number?: string;
  } = {
    tax_enabled: false,
    tax_type: 'none',
    tax_rate: 0,
    tax_inclusive: false,
    tax_exempt_products: []
  };
  
  let initialForm = deepClone(form);
  let initialized = false;

  $: if (!initialized && shop && Object.keys(shop).length > 0) {
    // TODO: Load actual tax settings from shop data
    form = {
      tax_enabled: false,
      tax_type: 'none',
      tax_rate: 0,
      tax_inclusive: false,
      tax_exempt_products: []
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
      // await api(authFetch).updateTaxSettings(form);
      await refetchShopData();
      initialForm = deepClone(form);
      toast.success('Tax settings updated!');
    } catch (error) {
      toast.error('Failed to update tax settings');
    }
  }

  function handleTaxExemptProductsChange(value: string) {
    form.tax_exempt_products = value.split(',').map(item => item.trim()).filter(Boolean);
  }
</script>
<ComingSoon>
<Card.Root class="bg-gradient w-full max-w-screen-md mx-auto">
  <Card.Header>
    <Card.Title>Tax Configuration</Card.Title>
    <Card.Description>Configure tax rates, rules, and exemptions for your store.</Card.Description>
  </Card.Header>
  <Card.Content>
    {#if !shop || Object.keys(shop).length === 0}
      <p>Loading...</p>
    {:else}
      <form class="space-y-6 w-full" on:submit={handleFormSubmit}>
        <div class="flex items-center justify-between">
          <div class="space-y-0.5">
            <Label>Enable Tax Collection</Label>
            <p class="text-muted-foreground text-sm">Collect and remit taxes on sales</p>
          </div>
          <Switch bind:checked={form.tax_enabled} />
        </div>

        {#if form.tax_enabled}
          <div class="space-y-4">
            <div>
              <Label for="tax-type">Tax Type</Label>
              <Select >
                <SelectTrigger>
                  <SelectValue placeholder="Select tax type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">No Tax</SelectItem>
                  <SelectItem value="vat">VAT (Value Added Tax)</SelectItem>
                  <SelectItem value="gst">GST (Goods & Services Tax)</SelectItem>
                  <SelectItem value="sales">Sales Tax</SelectItem>
                  <SelectItem value="custom">Custom Tax</SelectItem>
                </SelectContent>
              </Select>
              <p class="text-muted-foreground text-sm">Select the type of tax you need to collect</p>
            </div>

            <div>
              <Label for="tax-rate">Tax Rate (%)</Label>
              <Input 
                id="tax-rate" 
                type="number" 
                bind:value={form.tax_rate} 
                placeholder="0.00" 
                min="0" 
                max="100" 
                step="0.01"
              />
              <p class="text-muted-foreground text-sm">Percentage rate to apply to taxable items</p>
            </div>

            <div>
              <Label for="tax-registration">Tax Registration Number</Label>
              <Input 
                id="tax-registration" 
                bind:value={form.tax_registration_number} 
                placeholder="Enter your tax registration number"
              />
              <p class="text-muted-foreground text-sm">Your official tax registration number</p>
            </div>

            <div class="flex items-center justify-between">
              <div class="space-y-0.5">
                <Label>Tax Inclusive Pricing</Label>
                <p class="text-muted-foreground text-sm">Display prices with tax included</p>
              </div>
              <Switch bind:checked={form.tax_inclusive} />
            </div>

            <div>
              <Label for="tax-exempt">Tax Exempt Products</Label>
              <Input 
                id="tax-exempt"
                value={form.tax_exempt_products.join(', ')}
                placeholder="Product IDs separated by commas"
                on:input={(e) => handleTaxExemptProductsChange(e.currentTarget.value)}
              />
              <p class="text-muted-foreground text-sm">Product IDs that should be exempt from tax</p>
            </div>
          </div>
        {/if}

        <Button type="submit" disabled={!hasChanges} class="w-full">
          Update Tax Settings
        </Button>
      </form>
    {/if}
  </Card.Content>
</Card.Root>
</ComingSoon>