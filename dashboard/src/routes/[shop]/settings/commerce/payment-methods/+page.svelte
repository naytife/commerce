<script lang="ts">
  import * as Card from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';
  import { Button } from '$lib/components/ui/button';
  import { Label } from '$lib/components/ui/label';
  import { Checkbox } from '$lib/components/ui/checkbox';
  import { toast } from 'svelte-sonner';
  import type { PaymentMethod, PaymentMethodType, PaymentMethodConfig, Shop } from '$lib/types';
  import { api } from '$lib/api';
  import { getContext } from 'svelte';
  import { deepEqual, deepClone } from '$lib/utils/deepEqual';
  import type { Writable } from 'svelte/store';

  // --- Breadcrumbs ---
  const breadcrumbs = [
    { label: 'Settings' },
    { label: 'Commerce' },
    { label: 'Payment Methods' }
  ];

  // --- Shop context ---
  const shopStore = getContext('shop') as Writable<Partial<Shop>>;
  const authFetch: (input: RequestInfo | URL, init?: RequestInit | undefined) => Promise<Response> = getContext('authFetch');
  let shop: Partial<Shop> = {};
  $: shop = $shopStore;

  // Payment method management
  interface PaymentMethodSetting {
    id: string;
    label: string;
    type: 'text' | 'password' | 'checkbox';
    value: string | boolean;
    placeholder?: string;
    required: boolean;
  }

  interface UIPaymentMethod {
    id: PaymentMethodType;
    name: string;
    description: string;
    enabled: boolean;
    settings: PaymentMethodSetting[];
  }

  // Payment methods configurations
  const paymentMethodConfigs: Record<PaymentMethodType, Omit<UIPaymentMethod, 'enabled'>> = {
    stripe: {
      id: 'stripe',
      name: 'Stripe',
      description: 'Accept credit card payments via Stripe',
      settings: [
        {
          id: 'publishable_key',
          label: 'Publishable Key',
          type: 'text',
          value: '',
          placeholder: 'pk_test_...',
          required: true
        },
        {
          id: 'secret_key',
          label: 'Secret Key',
          type: 'password',
          value: '',
          placeholder: 'sk_test_...',
          required: true
        },
        {
          id: 'test_mode',
          label: 'Test Mode',
          type: 'checkbox',
          value: true,
          required: false
        }
      ]
    },
    paypal: {
      id: 'paypal',
      name: 'PayPal',
      description: 'Accept payments via PayPal',
      settings: [
        {
          id: 'client_id',
          label: 'Client ID',
          type: 'text',
          value: '',
          placeholder: 'Your PayPal client ID',
          required: true
        },
        {
          id: 'client_secret',
          label: 'Client Secret',
          type: 'password',
          value: '',
          placeholder: 'Your PayPal client secret',
          required: true
        },
        {
          id: 'sandbox_mode',
          label: 'Sandbox Mode',
          type: 'checkbox',
          value: true,
          required: false
        }
      ]
    },
    flutterwave: {
      id: 'flutterwave',
      name: 'Flutterwave',
      description: 'Accept payments via Flutterwave',
      settings: [
        {
          id: 'public_key',
          label: 'Public Key',
          type: 'text',
          value: '',
          placeholder: 'Your Flutterwave public key',
          required: true
        },
        {
          id: 'secret_key',
          label: 'Secret Key',
          type: 'password',
          value: '',
          placeholder: 'Your Flutterwave secret key',
          required: true
        },
        {
          id: 'encryption_key',
          label: 'Encryption Key',
          type: 'password',
          value: '',
          placeholder: 'Your Flutterwave encryption key',
          required: true
        },
        {
          id: 'test_mode',
          label: 'Test Mode',
          type: 'checkbox',
          value: true,
          required: false
        }
      ]
    },
    paystack: {
      id: 'paystack',
      name: 'Paystack',
      description: 'Accept payments via Paystack',
      settings: [
        {
          id: 'public_key',
          label: 'Public Key',
          type: 'text',
          value: '',
          placeholder: 'Your Paystack public key',
          required: true
        },
        {
          id: 'secret_key',
          label: 'Secret Key',
          type: 'password',
          value: '',
          placeholder: 'Your Paystack secret key',
          required: true
        },
        {
          id: 'test_mode',
          label: 'Test Mode',
          type: 'checkbox',
          value: true,
          required: false
        }
      ]
    }
  };

  let paymentMethods: UIPaymentMethod[] = [];
  let loading = true;
  let initialPaymentMethods: UIPaymentMethod[] = [];
  let hasChanges = false;

  function resetInitialPaymentMethods() {
    initialPaymentMethods = deepClone(paymentMethods);
  }

  $: hasChanges = !deepEqual(paymentMethods, initialPaymentMethods);

  // --- Load payment methods only after shop is loaded ---
  $: if (shop && Object.keys(shop).length > 0) {
    loadPaymentMethods();
  }

  async function loadPaymentMethods() {
    try {
      loading = true;
      // No shop.id or shop.slug needed for API call
      const apiResponse = await api(authFetch).getPaymentMethods();
      const apiPaymentMethods = (apiResponse || []).map(pm => ({
        method_type: pm.id, // map 'id' to 'method_type'
        is_enabled: pm.enabled, // map 'enabled' to 'is_enabled'
        config: pm.config || {},
      }));
      paymentMethods = Object.keys(paymentMethodConfigs).map(methodType => {
        const config = paymentMethodConfigs[methodType as PaymentMethodType];
        const existingMethod = apiPaymentMethods.find(pm => pm.method_type === methodType);
        const settings = config.settings.map(setting => ({
          ...setting,
          value: existingMethod?.config?.[setting.id] ?? setting.value
        }));
        return {
          ...config,
          enabled: existingMethod?.is_enabled ?? false,
          settings
        };
      });
      resetInitialPaymentMethods();
    } catch (error) {
      console.error('Error loading payment methods:', error);
      toast.error('Failed to load payment methods');
    } finally {
      loading = false;
    }
  }

  async function togglePaymentMethod(method: UIPaymentMethod) {
    const newStatus = !method.enabled;
    try {
      await api(authFetch).updatePaymentMethodStatus(method.id, newStatus);
      method.enabled = newStatus;
      if (newStatus) {
        const missingSettings = method.settings.filter(setting => setting.required && setting.type !== 'checkbox' && !setting.value);
        if (missingSettings.length > 0) {
          toast.warning(`${method.name} is enabled but please configure the required settings`);
        } else {
          toast.success(`${method.name} has been enabled`);
        }
      } else {
        toast.success(`${method.name} has been disabled`);
      }
      resetInitialPaymentMethods();
    } catch (error) {
      console.error('Error toggling payment method:', error);
      toast.error(`Failed to ${newStatus ? 'enable' : 'disable'} ${method.name}`);
    }
  }

  async function savePaymentMethodSettings(method: UIPaymentMethod) {
    const missingSettings = method.settings.filter(setting => setting.required && setting.type !== 'checkbox' && !setting.value);
    if (missingSettings.length > 0) {
      toast.error('Please fill out all required fields');
      return;
    }
    try {
      const config: Record<string, any> = {};
      method.settings.forEach(setting => {
        config[setting.id] = setting.value;
      });
      const paymentMethodConfig: PaymentMethodConfig = {
        method_type: method.id,
        is_enabled: method.enabled,
        config
      };
      await api(authFetch).updatePaymentMethod(method.id, paymentMethodConfig);
      toast.success(`${method.name} settings saved successfully`);
      resetInitialPaymentMethods();
    } catch (error) {
      console.error('Error updating payment settings:', error);
      toast.error(`Failed to update ${method.name} settings`);
    }
  }

  async function handleFormSubmit(event: Event) {
    event.preventDefault();
    const form = event.target as HTMLFormElement;
    const methodIdInput = form.querySelector('input[name="method-id"]') as HTMLInputElement;
    const methodId = methodIdInput?.value;
    const method = paymentMethods.find(m => m.id === methodId);
    if (!method) {
      toast.error('Payment method not found');
      return;
    }
    await savePaymentMethodSettings(method);
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
    <Card.Title>Payment Methods</Card.Title>
    <Card.Description>Configure the payment options your customers can use at checkout</Card.Description>
  </Card.Header>
  <Card.Content>
    {#if loading}
      <div class="flex items-center justify-center py-8">
        <div class="text-muted-foreground">Loading payment methods...</div>
      </div>
    {:else}
      <div class="space-y-6">
        {#each paymentMethods as method}
          <div class="border rounded-md overflow-hidden">
            <div class="flex items-center justify-between p-4 bg-muted/50">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 flex items-center justify-center bg-background rounded-md border">
                  <span class="font-semibold">{method.name.charAt(0)}</span>
                </div>
                <div>
                  <h3 class="font-medium">{method.name}</h3>
                  <p class="text-sm text-muted-foreground">{method.description}</p>
                </div>
              </div>
              <div class="flex items-center space-x-2">
                <Checkbox 
                  id={`enable-${method.id}`} 
                  checked={method.enabled}
                  onCheckedChange={() => togglePaymentMethod(method)}
                />
                <Label
                  for={`enable-${method.id}`}
                  class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                  {method.enabled ? 'Enabled' : 'Disabled'}
                </Label>
              </div>
            </div>
            <div class="p-4 space-y-4">
              <form method="POST" class="space-y-4" id={`payment-method-${method.id}-form`} on:submit={handleFormSubmit}>
                <input type="hidden" name="form-id" value="payment-method-form" />
                <input type="hidden" name="method-id" value={method.id} />
                {#each method.settings as setting}
                  <div class="flex w-full flex-col gap-2">
                    {#if setting.type === 'checkbox'}
                      <div class="flex items-center space-x-2">
                        <Checkbox 
                          id={`${method.id}-${setting.id}`} 
                          name={setting.id}
                          checked={setting.value === true}
                          onCheckedChange={(checked) => setting.value = checked}
                        />
                        <Label
                          for={`${method.id}-${setting.id}`}
                          class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                        >
                          {setting.label}
                        </Label>
                      </div>
                    {:else}
                      <Label for={`${method.id}-${setting.id}`}>{setting.label}</Label>
                      <Input 
                        id={`${method.id}-${setting.id}`} 
                        name={setting.id}
                        type={setting.type}
                        bind:value={setting.value}
                        placeholder={setting.placeholder || ''}
                        required={setting.required}
                      />
                    {/if}
                  </div>
                {/each}
                <Button 
                  type="submit"
                  class="mt-2"
                  disabled={!hasChanges}
                >
                  Save {method.name} Settings
                </Button>
              </form>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </Card.Content>
</Card.Root>
