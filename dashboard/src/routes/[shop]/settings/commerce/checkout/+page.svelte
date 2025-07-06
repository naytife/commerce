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
    require_account: boolean;
    require_shipping_address: boolean;
    require_billing_address: boolean;
    allow_notes: boolean;
    terms_required: boolean;
    terms_text?: string;
    privacy_policy_required: boolean;
    privacy_policy_url?: string;
    order_confirmation_email: boolean;
    order_confirmation_template?: string;
    abandoned_cart_emails: boolean;
    abandoned_cart_delay_hours: number;
    max_retry_attempts: number;
    retry_delay_minutes: number;
  } = {
    guest_checkout: true,
    require_account: false,
    require_shipping_address: true,
    require_billing_address: true,
    allow_notes: true,
    terms_required: true,
    privacy_policy_required: false,
    order_confirmation_email: true,
    abandoned_cart_emails: false,
    abandoned_cart_delay_hours: 24,
    max_retry_attempts: 3,
    retry_delay_minutes: 30
  };
  
  let initialForm = deepClone(form);
  let initialized = false;

  $: if (!initialized && shop && Object.keys(shop).length > 0) {
    // TODO: Load actual checkout settings from shop data
    form = {
      guest_checkout: true,
      require_account: false,
      require_shipping_address: true,
      require_billing_address: true,
      allow_notes: true,
      terms_required: true,
      privacy_policy_required: false,
      order_confirmation_email: true,
      abandoned_cart_emails: false,
      abandoned_cart_delay_hours: 24,
      max_retry_attempts: 3,
      retry_delay_minutes: 30
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

          <div class="flex items-center justify-between">
            <div class="space-y-0.5">
              <Label>Require Account Creation</Label>
              <p class="text-muted-foreground text-sm">Force customers to create an account before checkout</p>
            </div>
            <Switch bind:checked={form.require_account} />
          </div>
        </div>

        <!-- Address Requirements -->
        <div class="space-y-4">
          <h3 class="text-lg font-semibold">Address Requirements</h3>
          
          <div class="flex items-center justify-between">
            <div class="space-y-0.5">
              <Label>Require Shipping Address</Label>
              <p class="text-muted-foreground text-sm">Collect shipping address for all orders</p>
            </div>
            <Switch bind:checked={form.require_shipping_address} />
          </div>

          <div class="flex items-center justify-between">
            <div class="space-y-0.5">
              <Label>Require Billing Address</Label>
              <p class="text-muted-foreground text-sm">Collect billing address for payment processing</p>
            </div>
            <Switch bind:checked={form.require_billing_address} />
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

          <div class="flex items-center justify-between">
            <div class="space-y-0.5">
              <Label>Require Privacy Policy</Label>
              <p class="text-muted-foreground text-sm">Customers must accept privacy policy</p>
            </div>
            <Switch bind:checked={form.privacy_policy_required} />
          </div>

          {#if form.privacy_policy_required}
            <div>
              <Label for="privacy-policy-url">Privacy Policy URL</Label>
              <Input 
                id="privacy-policy-url" 
                bind:value={form.privacy_policy_url} 
                placeholder="https://yourstore.com/privacy"
                type="url"
              />
              <p class="text-muted-foreground text-sm">Link to your privacy policy</p>
            </div>
          {/if}
        </div>

        <!-- Email Notifications -->
        <div class="space-y-4">
          <h3 class="text-lg font-semibold">Email Notifications</h3>
          
          <div class="flex items-center justify-between">
            <div class="space-y-0.5">
              <Label>Order Confirmation Emails</Label>
              <p class="text-muted-foreground text-sm">Send confirmation emails when orders are placed</p>
            </div>
            <Switch bind:checked={form.order_confirmation_email} />
          </div>

          {#if form.order_confirmation_email}
            <div>
              <Label for="confirmation-template">Confirmation Email Template</Label>
              <Textarea 
                id="confirmation-template" 
                bind:value={form.order_confirmation_template} 
                placeholder="Customize your order confirmation email..."
                rows={4}
              />
              <p class="text-muted-foreground text-sm">Custom template for order confirmation emails</p>
            </div>
          {/if}

          <div class="flex items-center justify-between">
            <div class="space-y-0.5">
              <Label>Abandoned Cart Emails</Label>
              <p class="text-muted-foreground text-sm">Send reminder emails for abandoned carts</p>
            </div>
            <Switch bind:checked={form.abandoned_cart_emails} />
          </div>

          {#if form.abandoned_cart_emails}
            <div>
              <Label for="abandoned-cart-delay">Abandoned Cart Delay (hours)</Label>
              <Input 
                id="abandoned-cart-delay" 
                type="number" 
                value={form.abandoned_cart_delay_hours}
                placeholder="24" 
                min="1" 
                max="168"
                on:input={(e) => form.abandoned_cart_delay_hours = Number(e.currentTarget.value) || 24}
              />
              <p class="text-muted-foreground text-sm">Hours to wait before sending abandoned cart email</p>
            </div>
          {/if}
        </div>

        <!-- Payment Retry Settings -->
        <div class="space-y-4">
          <h3 class="text-lg font-semibold">Payment Retry Settings</h3>
          
          <div>
            <Label for="max-retry-attempts">Maximum Retry Attempts</Label>
            <Input 
              id="max-retry-attempts" 
              type="number" 
              value={form.max_retry_attempts}
              placeholder="3" 
              min="1" 
              max="10"
              on:input={(e) => form.max_retry_attempts = Number(e.currentTarget.value) || 3}
            />
            <p class="text-muted-foreground text-sm">Number of times to retry failed payments</p>
          </div>

          <div>
            <Label for="retry-delay">Retry Delay (minutes)</Label>
            <Input 
              id="retry-delay" 
              type="number" 
              value={form.retry_delay_minutes}
              placeholder="30" 
              min="5" 
              max="1440"
              on:input={(e) => form.retry_delay_minutes = Number(e.currentTarget.value) || 30}
            />
            <p class="text-muted-foreground text-sm">Minutes to wait between payment retry attempts</p>
          </div>
        </div>

        <Button type="submit" disabled={!hasChanges} class="w-full">
          Update Checkout Settings
        </Button>
      </form>
    {/if}
  </Card.Content>
</Card.Root>
