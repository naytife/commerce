<script lang="ts">
  import * as Card from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';
  import { Button } from '$lib/components/ui/button';
  import { Label } from '$lib/components/ui/label';
  import { Switch } from '$lib/components/ui/switch';
  import { Badge } from '$lib/components/ui/badge';
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

  const analyticsProviders = [
    {
      id: 'google_analytics',
      name: 'Google Analytics',
      description: 'Track website traffic and user behavior',
      icon: '📊',
      configFields: [
        { key: 'tracking_id', label: 'Tracking ID', type: 'text', placeholder: 'G-XXXXXXXXXX', required: true },
        { key: 'measurement_id', label: 'Measurement ID', type: 'text', placeholder: 'G-XXXXXXXXXX', required: false }
      ]
    },
    {
      id: 'facebook_pixel',
      name: 'Facebook Pixel',
      description: 'Track conversions and build custom audiences',
      icon: '📱',
      configFields: [
        { key: 'pixel_id', label: 'Pixel ID', type: 'text', placeholder: '123456789012345', required: true }
      ]
    },
    {
      id: 'google_tag_manager',
      name: 'Google Tag Manager',
      description: 'Manage all your tracking tags in one place',
      icon: '🏷️',
      configFields: [
        { key: 'container_id', label: 'Container ID', type: 'text', placeholder: 'GTM-XXXXXXX', required: true }
      ]
    },
    {
      id: 'hotjar',
      name: 'Hotjar',
      description: 'Understand user behavior with heatmaps and recordings',
      icon: '🔥',
      configFields: [
        { key: 'site_id', label: 'Site ID', type: 'text', placeholder: '1234567', required: true }
      ]
    },
    {
      id: 'mixpanel',
      name: 'Mixpanel',
      description: 'Advanced analytics and user behavior tracking',
      icon: '📈',
      configFields: [
        { key: 'project_token', label: 'Project Token', type: 'text', placeholder: 'your_project_token', required: true }
      ]
    }
  ];

  let form: {
    analytics_enabled: boolean;
    providers: Array<{
      id: string;
      enabled: boolean;
      config: Record<string, string>;
    }>;
    enhanced_ecommerce: boolean;
    conversion_tracking: boolean;
    custom_events: boolean;
  } = {
    analytics_enabled: false,
    providers: analyticsProviders.map(provider => ({
      id: provider.id,
      enabled: false,
      config: {}
    })),
    enhanced_ecommerce: false,
    conversion_tracking: false,
    custom_events: false
  };
  
  let initialForm = deepClone(form);
  let initialized = false;

  $: if (!initialized && shop && Object.keys(shop).length > 0) {
    // TODO: Load actual analytics settings from shop data
    form = {
      analytics_enabled: false,
      providers: analyticsProviders.map(provider => ({
        id: provider.id,
        enabled: false,
        config: {}
      })),
      enhanced_ecommerce: false,
      conversion_tracking: false,
      custom_events: false
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
      // await api(authFetch).updateAnalyticsSettings(form);
      await refetchShopData();
      initialForm = deepClone(form);
      toast.success('Analytics settings updated!');
    } catch (error) {
      toast.error('Failed to update analytics settings');
    }
  }

  function toggleProvider(providerId: string) {
    const provider = form.providers.find(p => p.id === providerId);
    if (provider) {
      provider.enabled = !provider.enabled;
      form.providers = [...form.providers];
    }
  }

  function updateProviderConfig(providerId: string, key: string, value: string) {
    const provider = form.providers.find(p => p.id === providerId);
    if (provider) {
      provider.config[key] = value;
      form.providers = [...form.providers];
    }
  }

  function getProviderConfig(providerId: string) {
    return form.providers.find(p => p.id === providerId)?.config || {};
  }

  function isProviderEnabled(providerId: string) {
    return form.providers.find(p => p.id === providerId)?.enabled || false;
  }
</script>
<ComingSoon>
<Card.Root class="bg-gradient w-full max-w-screen-md mx-auto">
  <Card.Header>
    <Card.Title>Analytics & Tracking</Card.Title>
    <Card.Description>Configure analytics tools and tracking codes for your store.</Card.Description>
  </Card.Header>
  <Card.Content>
    {#if !shop || Object.keys(shop).length === 0}
      <p>Loading...</p>
    {:else}
      <form class="space-y-6 w-full" on:submit={handleFormSubmit}>
        <!-- General Analytics Settings -->
        <div class="space-y-4">
          <h3 class="text-lg font-semibold">General Settings</h3>
          
          <div class="flex items-center justify-between">
            <div class="space-y-0.5">
              <Label>Enable Analytics</Label>
              <p class="text-muted-foreground text-sm">Track visitor behavior and store performance</p>
            </div>
            <Switch bind:checked={form.analytics_enabled} />
          </div>

          {#if form.analytics_enabled}
            <div class="flex items-center justify-between">
              <div class="space-y-0.5">
                <Label>Enhanced E-commerce</Label>
                <p class="text-muted-foreground text-sm">Track detailed e-commerce events and conversions</p>
              </div>
              <Switch bind:checked={form.enhanced_ecommerce} />
            </div>

            <div class="flex items-center justify-between">
              <div class="space-y-0.5">
                <Label>Conversion Tracking</Label>
                <p class="text-muted-foreground text-sm">Track sales, signups, and other conversions</p>
              </div>
              <Switch bind:checked={form.conversion_tracking} />
            </div>

            <div class="flex items-center justify-between">
              <div class="space-y-0.5">
                <Label>Custom Events</Label>
                <p class="text-muted-foreground text-sm">Track custom user interactions and events</p>
              </div>
              <Switch bind:checked={form.custom_events} />
            </div>
          {/if}
        </div>

        <!-- Analytics Providers -->
        {#if form.analytics_enabled}
          <div class="space-y-4">
            <h3 class="text-lg font-semibold">Analytics Providers</h3>
            
            {#each analyticsProviders as provider}
              <Card.Root class="border">
                <Card.Header class="pb-3">
                  <div class="flex items-center justify-between">
                    <div class="flex items-center space-x-3">
                      <span class="text-2xl">{provider.icon}</span>
                      <div>
                        <Card.Title class="text-lg">{provider.name}</Card.Title>
                        <Card.Description>{provider.description}</Card.Description>
                      </div>
                    </div>
                    <div class="flex items-center space-x-2">
                      {#if isProviderEnabled(provider.id)}
                        <Badge variant="default">Active</Badge>
                      {:else}
                        <Badge variant="secondary">Inactive</Badge>
                      {/if}
                      <Switch 
                        checked={isProviderEnabled(provider.id)} 
                        onCheckedChange={() => toggleProvider(provider.id)}
                      />
                    </div>
                  </div>
                </Card.Header>
                
                {#if isProviderEnabled(provider.id)}
                  <Card.Content class="pt-0">
                    <div class="space-y-4">
                      {#each provider.configFields as field}
                        <div>
                          <Label for="{provider.id}-{field.key}">{field.label}</Label>
                          <Input 
                            id="{provider.id}-{field.key}"
                            type="text"
                            value={getProviderConfig(provider.id)[field.key] || ''}
                            placeholder={field.placeholder}
                            required={field.required}
                            on:input={(e) => updateProviderConfig(provider.id, field.key, e.currentTarget.value)}
                          />
                          {#if field.required}
                            <p class="text-muted-foreground text-sm">Required</p>
                          {/if}
                        </div>
                      {/each}
                    </div>
                  </Card.Content>
                {/if}
              </Card.Root>
            {/each}
          </div>
        {/if}

        <!-- Tracking Events -->
        {#if form.analytics_enabled && form.custom_events}
          <div class="space-y-4">
            <h3 class="text-lg font-semibold">Custom Events</h3>
            <Card.Root class="border">
              <Card.Content class="pt-4">
                <div class="space-y-3">
                  <div>
                    <h4 class="font-medium">Available Events</h4>
                    <ul class="text-sm text-muted-foreground space-y-1">
                      <li>• <code>product_view</code> - When a product is viewed</li>
                      <li>• <code>add_to_cart</code> - When an item is added to cart</li>
                      <li>• <code>remove_from_cart</code> - When an item is removed from cart</li>
                      <li>• <code>begin_checkout</code> - When checkout process starts</li>
                      <li>• <code>purchase</code> - When an order is completed</li>
                      <li>• <code>sign_up</code> - When a customer creates an account</li>
                      <li>• <code>login</code> - When a customer logs in</li>
                    </ul>
                  </div>
                  <div>
                    <h4 class="font-medium">Event Properties</h4>
                    <p class="text-sm text-muted-foreground">
                      Each event includes relevant data like product details, order values, and user information.
                    </p>
                  </div>
                </div>
              </Card.Content>
            </Card.Root>
          </div>
        {/if}

        <!-- Data Privacy -->
        <div class="space-y-4">
          <h3 class="text-lg font-semibold">Data Privacy</h3>
          <Card.Root class="border">
            <Card.Content class="pt-4">
              <div class="space-y-3">
                <div>
                  <h4 class="font-medium">GDPR Compliance</h4>
                  <p class="text-sm text-muted-foreground">
                    Analytics tracking respects user privacy preferences and GDPR requirements.
                  </p>
                </div>
                <div>
                  <h4 class="font-medium">Cookie Consent</h4>
                  <p class="text-sm text-muted-foreground">
                    Analytics cookies are only set after user consent is obtained.
                  </p>
                </div>
                <div>
                  <h4 class="font-medium">Data Retention</h4>
                  <p class="text-sm text-muted-foreground">
                    Analytics data is retained according to your privacy policy and applicable regulations.
                  </p>
                </div>
              </div>
            </Card.Content>
          </Card.Root>
        </div>

        <Button type="submit" disabled={!hasChanges} class="w-full">
          Update Analytics Settings
        </Button>
      </form>
    {/if}
  </Card.Content>
</Card.Root>
</ComingSoon>