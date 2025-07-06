<script lang="ts">
  import * as Card from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';
  import { Button } from '$lib/components/ui/button';
  import { Label } from '$lib/components/ui/label';
  import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '$lib/components/ui/select';
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
    shipping_enabled: boolean;
    free_shipping_threshold: number;
    default_shipping_cost: number;
    shipping_zones: Array<{
      name: string;
      countries: string[];
      cost: number;
      free_shipping_threshold: number;
    }>;
    weight_based_shipping: boolean;
    weight_unit: 'kg' | 'lb' | 'g' | 'oz';
    handling_fee: number;
    shipping_policy?: string;
  } = {
    shipping_enabled: false,
    free_shipping_threshold: 0,
    default_shipping_cost: 0,
    shipping_zones: [],
    weight_based_shipping: false,
    weight_unit: 'kg',
    handling_fee: 0
  };
  
  let initialForm = deepClone(form);
  let initialized = false;

  $: if (!initialized && shop && Object.keys(shop).length > 0) {
    // TODO: Load actual shipping settings from shop data
    form = {
      shipping_enabled: false,
      free_shipping_threshold: 0,
      default_shipping_cost: 0,
      shipping_zones: [],
      weight_based_shipping: false,
      weight_unit: 'kg',
      handling_fee: 0
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
      // await api(authFetch).updateShippingSettings(form);
      await refetchShopData();
      initialForm = deepClone(form);
      toast.success('Shipping settings updated!');
    } catch (error) {
      toast.error('Failed to update shipping settings');
    }
  }

  function addShippingZone() {
    form.shipping_zones = [...form.shipping_zones, {
      name: '',
      countries: [],
      cost: 0,
      free_shipping_threshold: 0
    }];
  }

  function removeShippingZone(index: number) {
    form.shipping_zones = form.shipping_zones.filter((_, i) => i !== index);
  }

  function updateShippingZone(index: number, field: string, value: any) {
    form.shipping_zones[index] = { ...form.shipping_zones[index], [field]: value };
    form.shipping_zones = [...form.shipping_zones];
  }
</script>

<Card.Root class="bg-gradient w-full max-w-screen-md mx-auto">
  <Card.Header>
    <Card.Title>Shipping Configuration</Card.Title>
    <Card.Description>Configure shipping rates, zones, and delivery options for your store.</Card.Description>
  </Card.Header>
  <Card.Content>
    {#if !shop || Object.keys(shop).length === 0}
      <p>Loading...</p>
    {:else}
      <form class="space-y-6 w-full" on:submit={handleFormSubmit}>
        <div class="flex items-center justify-between">
          <div class="space-y-0.5">
            <Label>Enable Shipping</Label>
            <p class="text-muted-foreground text-sm">Allow customers to select shipping options</p>
          </div>
          <Switch bind:checked={form.shipping_enabled} />
        </div>

        {#if form.shipping_enabled}
          <div class="space-y-4">
            <div>
              <Label for="default-shipping-cost">Default Shipping Cost</Label>
              <Input 
                id="default-shipping-cost" 
                type="number" 
                bind:value={form.default_shipping_cost} 
                placeholder="0.00" 
                min="0" 
                step="0.01"
              />
              <p class="text-muted-foreground text-sm">Default shipping cost for all orders</p>
            </div>

            <div>
              <Label for="free-shipping-threshold">Free Shipping Threshold</Label>
              <Input 
                id="free-shipping-threshold" 
                type="number" 
                bind:value={form.free_shipping_threshold} 
                placeholder="0.00" 
                min="0" 
                step="0.01"
              />
              <p class="text-muted-foreground text-sm">Order value at which shipping becomes free</p>
            </div>

            <div>
              <Label for="handling-fee">Handling Fee</Label>
              <Input 
                id="handling-fee" 
                type="number" 
                bind:value={form.handling_fee} 
                placeholder="0.00" 
                min="0" 
                step="0.01"
              />
              <p class="text-muted-foreground text-sm">Additional handling fee per order</p>
            </div>

            <div class="flex items-center justify-between">
              <div class="space-y-0.5">
                <Label>Weight-Based Shipping</Label>
                <p class="text-muted-foreground text-sm">Calculate shipping based on product weight</p>
              </div>
              <Switch bind:checked={form.weight_based_shipping} />
            </div>

            {#if form.weight_based_shipping}
              <div>
                <Label for="weight-unit">Weight Unit</Label>
                <select bind:value={form.weight_unit} class="w-full p-2 border rounded-md bg-background">
                  <option value="kg">Kilograms (kg)</option>
                  <option value="lb">Pounds (lb)</option>
                  <option value="g">Grams (g)</option>
                  <option value="oz">Ounces (oz)</option>
                </select>
                <p class="text-muted-foreground text-sm">Unit of measurement for weight-based calculations</p>
              </div>
            {/if}

            <div>
              <Label for="shipping-policy">Shipping Policy</Label>
              <Textarea 
                id="shipping-policy" 
                bind:value={form.shipping_policy} 
                placeholder="Enter your shipping policy details..."
                rows={3}
              />
              <p class="text-muted-foreground text-sm">Information about shipping times and policies</p>
            </div>

            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <Label>Shipping Zones</Label>
                <Button type="button" variant="outline" size="sm" on:click={addShippingZone}>
                  Add Zone
                </Button>
              </div>
              
              {#each form.shipping_zones as zone, index}
                <Card.Root class="p-4 border">
                  <div class="space-y-3">
                    <div class="flex items-center justify-between">
                      <Label>Zone {index + 1}</Label>
                      <Button 
                        type="button" 
                        variant="ghost" 
                        size="sm" 
                        on:click={() => removeShippingZone(index)}
                      >
                        Remove
                      </Button>
                    </div>
                    <div>
                      <Label>Zone Name</Label>
                      <Input 
                        value={zone.name}
                        placeholder="e.g., North America"
                        on:input={(e) => updateShippingZone(index, 'name', e.currentTarget.value)}
                      />
                    </div>
                    <div>
                      <Label>Countries (comma-separated)</Label>
                      <Input 
                        value={zone.countries.join(', ')}
                        placeholder="US, CA, MX"
                        on:input={(e) => updateShippingZone(index, 'countries', e.currentTarget.value.split(',').map(c => c.trim()).filter(Boolean))}
                      />
                    </div>
                    <div class="grid grid-cols-2 gap-3">
                      <div>
                        <Label>Shipping Cost</Label>
                        <Input 
                          type="number"
                          value={zone.cost}
                          placeholder="0.00"
                          min="0"
                          step="0.01"
                          on:input={(e) => updateShippingZone(index, 'cost', Number(e.currentTarget.value) || 0)}
                        />
                      </div>
                      <div>
                        <Label>Free Shipping Threshold</Label>
                        <Input 
                          type="number"
                          value={zone.free_shipping_threshold}
                          placeholder="0.00"
                          min="0"
                          step="0.01"
                          on:input={(e) => updateShippingZone(index, 'free_shipping_threshold', Number(e.currentTarget.value) || 0)}
                        />
                      </div>
                    </div>
                  </div>
                </Card.Root>
              {/each}
            </div>
          </div>
        {/if}

        <Button type="submit" disabled={!hasChanges} class="w-full">
          Update Shipping Settings
        </Button>
      </form>
    {/if}
  </Card.Content>
</Card.Root>
