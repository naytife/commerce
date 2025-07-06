<script lang="ts">
  import * as Card from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';
  import { Button } from '$lib/components/ui/button';
  import { Label } from '$lib/components/ui/label';
  import { Switch } from '$lib/components/ui/switch';
  import { Textarea } from '$lib/components/ui/textarea';
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

  let form: {
    api_enabled: boolean;
    api_key?: string;
    webhook_enabled: boolean;
    webhook_url?: string;
    webhook_secret?: string;
    custom_apps: Array<{
      id: string;
      name: string;
      description: string;
      api_key: string;
      permissions: string[];
      webhook_url?: string;
      enabled: boolean;
    }>;
  } = {
    api_enabled: false,
    webhook_enabled: false,
    custom_apps: []
  };
  
  let initialForm = deepClone(form);
  let initialized = false;

  $: if (!initialized && shop && Object.keys(shop).length > 0) {
    // TODO: Load actual custom apps settings from shop data
    form = {
      api_enabled: false,
      webhook_enabled: false,
      custom_apps: []
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
      // await api(authFetch).updateCustomAppsSettings(form);
      await refetchShopData();
      initialForm = deepClone(form);
      toast.success('Custom apps settings updated!');
    } catch (error) {
      toast.error('Failed to update custom apps settings');
    }
  }

  function addCustomApp() {
    form.custom_apps = [...form.custom_apps, {
      id: `app_${Date.now()}`,
      name: '',
      description: '',
      api_key: '',
      permissions: [],
      enabled: true
    }];
  }

  function removeCustomApp(index: number) {
    form.custom_apps = form.custom_apps.filter((_, i) => i !== index);
  }

  function updateCustomApp(index: number, field: string, value: any) {
    form.custom_apps[index] = { ...form.custom_apps[index], [field]: value };
    form.custom_apps = [...form.custom_apps];
  }

  function generateApiKey() {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let result = 'sk_';
    for (let i = 0; i < 32; i++) {
      result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return result;
  }

  const availablePermissions = [
    'products:read',
    'products:write',
    'orders:read',
    'orders:write',
    'customers:read',
    'customers:write',
    'analytics:read',
    'webhooks:manage'
  ];
</script>
<ComingSoon>
<Card.Root class="bg-gradient w-full max-w-screen-md mx-auto">
  <Card.Header>
    <Card.Title>Custom Apps & API</Card.Title>
    <Card.Description>Manage API access, webhooks, and custom integrations for your store.</Card.Description>
  </Card.Header>
  <Card.Content>
    {#if !shop || Object.keys(shop).length === 0}
      <p>Loading...</p>
    {:else}
      <form class="space-y-6 w-full" on:submit={handleFormSubmit}>
        <!-- API Settings -->
        <div class="space-y-4">
          <h3 class="text-lg font-semibold">API Configuration</h3>
          
          <div class="flex items-center justify-between">
            <div class="space-y-0.5">
              <Label>Enable API Access</Label>
              <p class="text-muted-foreground text-sm">Allow external applications to access your store data</p>
            </div>
            <Switch bind:checked={form.api_enabled} />
          </div>

          {#if form.api_enabled}
            <div>
              <Label for="api-key">API Key</Label>
              <div class="flex space-x-2">
                <Input 
                  id="api-key" 
                  type="password" 
                  bind:value={form.api_key} 
                  placeholder="sk_live_..."
                  readonly
                />
                <Button 
                  type="button" 
                  variant="outline" 
                  on:click={() => form.api_key = generateApiKey()}
                >
                  Generate
                </Button>
              </div>
              <p class="text-muted-foreground text-sm">Your store's API key for external integrations</p>
            </div>
          {/if}
        </div>

        <!-- Webhook Settings -->
        <div class="space-y-4">
          <h3 class="text-lg font-semibold">Webhook Configuration</h3>
          
          <div class="flex items-center justify-between">
            <div class="space-y-0.5">
              <Label>Enable Webhooks</Label>
              <p class="text-muted-foreground text-sm">Send real-time notifications to external services</p>
            </div>
            <Switch bind:checked={form.webhook_enabled} />
          </div>

          {#if form.webhook_enabled}
            <div>
              <Label for="webhook-url">Webhook URL</Label>
              <Input 
                id="webhook-url" 
                bind:value={form.webhook_url} 
                placeholder="https://your-app.com/webhook"
                type="url"
              />
              <p class="text-muted-foreground text-sm">URL where webhook events will be sent</p>
            </div>

            <div>
              <Label for="webhook-secret">Webhook Secret</Label>
              <Input 
                id="webhook-secret" 
                type="password" 
                bind:value={form.webhook_secret} 
                placeholder="whsec_..."
              />
              <p class="text-muted-foreground text-sm">Secret key to verify webhook authenticity</p>
            </div>
          {/if}
        </div>

        <!-- Custom Apps -->
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold">Custom Applications</h3>
            <Button type="button" variant="outline" size="sm" on:click={addCustomApp}>
              Add App
            </Button>
          </div>
          
          {#each form.custom_apps as app, index}
            <Card.Root class="border">
              <Card.Header class="pb-3">
                <div class="flex items-center justify-between">
                  <div class="flex items-center space-x-2">
                    <Label>App {index + 1}</Label>
                    {#if app.enabled}
                      <Badge variant="default">Active</Badge>
                    {:else}
                      <Badge variant="secondary">Inactive</Badge>
                    {/if}
                  </div>
                  <Button 
                    type="button" 
                    variant="ghost" 
                    size="sm" 
                    on:click={() => removeCustomApp(index)}
                  >
                    Remove
                  </Button>
                </div>
              </Card.Header>
              
              <Card.Content class="pt-0">
                <div class="space-y-4">
                  <div>
                    <Label>App Name</Label>
                    <Input 
                      value={app.name}
                      placeholder="My Custom App"
                      on:input={(e) => updateCustomApp(index, 'name', e.currentTarget.value)}
                    />
                  </div>
                  
                  <div>
                    <Label>Description</Label>
                    <Textarea 
                      value={app.description}
                      placeholder="What does this app do?"
                      rows={2}
                      on:input={(e) => updateCustomApp(index, 'description', e.currentTarget.value)}
                    />
                  </div>
                  
                  <div>
                    <Label>API Key</Label>
                    <div class="flex space-x-2">
                      <Input 
                        type="password"
                        value={app.api_key}
                        placeholder="sk_app_..."
                        on:input={(e) => updateCustomApp(index, 'api_key', e.currentTarget.value)}
                      />
                      <Button 
                        type="button" 
                        variant="outline" 
                        on:click={() => updateCustomApp(index, 'api_key', generateApiKey())}
                      >
                        Generate
                      </Button>
                    </div>
                  </div>
                  
                  <div>
                    <Label>Webhook URL (Optional)</Label>
                    <Input 
                      value={app.webhook_url || ''}
                      placeholder="https://app.com/webhook"
                      type="url"
                      on:input={(e) => updateCustomApp(index, 'webhook_url', e.currentTarget.value)}
                    />
                  </div>
                  
                  <div>
                    <Label>Permissions</Label>
                    <div class="grid grid-cols-2 gap-2">
                      {#each availablePermissions as permission}
                        <label class="flex items-center space-x-2">
                          <input 
                            type="checkbox" 
                            checked={app.permissions.includes(permission)}
                            on:change={(e) => {
                              if (e.currentTarget.checked) {
                                updateCustomApp(index, 'permissions', [...app.permissions, permission]);
                              } else {
                                updateCustomApp(index, 'permissions', app.permissions.filter(p => p !== permission));
                              }
                            }}
                          />
                          <span class="text-sm">{permission}</span>
                        </label>
                      {/each}
                    </div>
                  </div>
                  
                  <div class="flex items-center justify-between">
                    <Label>Enabled</Label>
                    <Switch 
                      checked={app.enabled} 
                      onCheckedChange={(checked) => updateCustomApp(index, 'enabled', checked)}
                    />
                  </div>
                </div>
              </Card.Content>
            </Card.Root>
          {/each}
        </div>

        <!-- Documentation -->
        <div class="space-y-4">
          <h3 class="text-lg font-semibold">API Documentation</h3>
          <Card.Root class="border">
            <Card.Content class="pt-4">
              <div class="space-y-3">
                <div>
                  <h4 class="font-medium">Base URL</h4>
                  <code class="text-sm bg-muted p-2 rounded block">https://api.yourstore.com/v1</code>
                </div>
                <div>
                  <h4 class="font-medium">Authentication</h4>
                  <p class="text-sm text-muted-foreground">Include your API key in the Authorization header:</p>
                  <code class="text-sm bg-muted p-2 rounded block">Authorization: Bearer YOUR_API_KEY</code>
                </div>
                <div>
                  <h4 class="font-medium">Available Endpoints</h4>
                  <ul class="text-sm text-muted-foreground space-y-1">
                    <li>• GET /products - List all products</li>
                    <li>• GET /orders - List all orders</li>
                    <li>• GET /customers - List all customers</li>
                    <li>• POST /webhooks - Create webhook</li>
                  </ul>
                </div>
              </div>
            </Card.Content>
          </Card.Root>
        </div>

        <Button type="submit" disabled={!hasChanges} class="w-full">
          Update Custom Apps Settings
        </Button>
      </form>
    {/if}
  </Card.Content>
</Card.Root>
</ComingSoon>