<script lang="ts">
  import * as Card from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';
  import { Button } from '$lib/components/ui/button';
  import { Label } from '$lib/components/ui/label';
  import { Switch } from '$lib/components/ui/switch';
  import { Badge } from '$lib/components/ui/badge';
  import { Textarea } from '$lib/components/ui/textarea';
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

  const socialPlatforms = [
    {
      id: 'facebook',
      name: 'Facebook',
      description: 'Connect your Facebook page and share products',
      icon: '📘',
      color: 'bg-blue-500',
      configFields: [
        { key: 'page_id', label: 'Page ID', type: 'text', placeholder: '123456789012345', required: true },
        { key: 'access_token', label: 'Access Token', type: 'password', placeholder: 'EAAG...', required: true },
        { key: 'auto_post', label: 'Auto-post new products', type: 'boolean', required: false }
      ]
    },
    {
      id: 'instagram',
      name: 'Instagram',
      description: 'Share products and stories on Instagram',
      icon: '📷',
      color: 'bg-pink-500',
      configFields: [
        { key: 'business_account_id', label: 'Business Account ID', type: 'text', placeholder: '17841412345678901', required: true },
        { key: 'access_token', label: 'Access Token', type: 'password', placeholder: 'IGQVJ...', required: true },
        { key: 'auto_post', label: 'Auto-post new products', type: 'boolean', required: false }
      ]
    },
    {
      id: 'twitter',
      name: 'Twitter',
      description: 'Tweet about new products and promotions',
      icon: '🐦',
      color: 'bg-blue-400',
      configFields: [
        { key: 'api_key', label: 'API Key', type: 'text', placeholder: 'your_api_key', required: true },
        { key: 'api_secret', label: 'API Secret', type: 'password', placeholder: 'your_api_secret', required: true },
        { key: 'access_token', label: 'Access Token', type: 'password', placeholder: 'your_access_token', required: true },
        { key: 'access_token_secret', label: 'Access Token Secret', type: 'password', placeholder: 'your_access_token_secret', required: true },
        { key: 'auto_tweet', label: 'Auto-tweet new products', type: 'boolean', required: false }
      ]
    },
    {
      id: 'pinterest',
      name: 'Pinterest',
      description: 'Pin products to Pinterest boards',
      icon: '📌',
      color: 'bg-red-500',
      configFields: [
        { key: 'app_id', label: 'App ID', type: 'text', placeholder: 'your_app_id', required: true },
        { key: 'access_token', label: 'Access Token', type: 'password', placeholder: 'your_access_token', required: true },
        { key: 'board_id', label: 'Board ID', type: 'text', placeholder: 'board_name', required: true },
        { key: 'auto_pin', label: 'Auto-pin new products', type: 'boolean', required: false }
      ]
    },
    {
      id: 'tiktok',
      name: 'TikTok',
      description: 'Share product videos on TikTok',
      icon: '🎵',
      color: 'bg-black',
      configFields: [
        { key: 'client_key', label: 'Client Key', type: 'text', placeholder: 'your_client_key', required: true },
        { key: 'client_secret', label: 'Client Secret', type: 'password', placeholder: 'your_client_secret', required: true },
        { key: 'access_token', label: 'Access Token', type: 'password', placeholder: 'your_access_token', required: true }
      ]
    }
  ];

  let form: {
    social_integration_enabled: boolean;
    platforms: Array<{
      id: string;
      enabled: boolean;
      config: Record<string, string | boolean>;
    }>;
    auto_sharing: boolean;
    share_new_products: boolean;
    share_promotions: boolean;
    custom_message_template?: string;
    hashtags?: string[];
  } = {
    social_integration_enabled: false,
    platforms: socialPlatforms.map(platform => ({
      id: platform.id,
      enabled: false,
      config: {}
    })),
    auto_sharing: false,
    share_new_products: false,
    share_promotions: false
  };
  
  let initialForm = deepClone(form);
  let initialized = false;

  $: if (!initialized && shop && Object.keys(shop).length > 0) {
    // TODO: Load actual social media settings from shop data
    form = {
      social_integration_enabled: false,
      platforms: socialPlatforms.map(platform => ({
        id: platform.id,
        enabled: false,
        config: {}
      })),
      auto_sharing: false,
      share_new_products: false,
      share_promotions: false
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
      // await api(authFetch).updateSocialMediaSettings(form);
      await refetchShopData();
      initialForm = deepClone(form);
      toast.success('Social media settings updated!');
    } catch (error) {
      toast.error('Failed to update social media settings');
    }
  }

  function togglePlatform(platformId: string) {
    const platform = form.platforms.find(p => p.id === platformId);
    if (platform) {
      platform.enabled = !platform.enabled;
      form.platforms = [...form.platforms];
    }
  }

  function updatePlatformConfig(platformId: string, key: string, value: string | boolean) {
    const platform = form.platforms.find(p => p.id === platformId);
    if (platform) {
      platform.config[key] = value;
      form.platforms = [...form.platforms];
    }
  }

  function getPlatformConfig(platformId: string) {
    return form.platforms.find(p => p.id === platformId)?.config || {};
  }

  function isPlatformEnabled(platformId: string) {
    return form.platforms.find(p => p.id === platformId)?.enabled || false;
  }

  function handleHashtagsChange(value: string) {
    form.hashtags = value.split(',').map(tag => tag.trim().replace('#', '')).filter(Boolean);
  }
</script>

<ComingSoon>
<Card.Root class="bg-gradient w-full max-w-screen-md mx-auto">
  <Card.Header>
    <Card.Title>Social Media Integration</Card.Title>
    <Card.Description>Connect your social media accounts and automate product sharing.</Card.Description>
  </Card.Header>
  <Card.Content>
    {#if !shop || Object.keys(shop).length === 0}
      <p>Loading...</p>
    {:else}
      <form class="space-y-6 w-full" on:submit={handleFormSubmit}>
        <!-- General Social Media Settings -->
        <div class="space-y-4">
          <h3 class="text-lg font-semibold">General Settings</h3>
          
          <div class="flex items-center justify-between">
            <div class="space-y-0.5">
              <Label>Enable Social Media Integration</Label>
              <p class="text-muted-foreground text-sm">Connect and manage your social media accounts</p>
            </div>
            <Switch bind:checked={form.social_integration_enabled} />
          </div>

          {#if form.social_integration_enabled}
            <div class="flex items-center justify-between">
              <div class="space-y-0.5">
                <Label>Auto-sharing</Label>
                <p class="text-muted-foreground text-sm">Automatically share content to connected platforms</p>
              </div>
              <Switch bind:checked={form.auto_sharing} />
            </div>

            {#if form.auto_sharing}
              <div class="flex items-center justify-between">
                <div class="space-y-0.5">
                  <Label>Share New Products</Label>
                  <p class="text-muted-foreground text-sm">Automatically share when new products are added</p>
                </div>
                <Switch bind:checked={form.share_new_products} />
              </div>

              <div class="flex items-center justify-between">
                <div class="space-y-0.5">
                  <Label>Share Promotions</Label>
                  <p class="text-muted-foreground text-sm">Automatically share discounts and promotions</p>
                </div>
                <Switch bind:checked={form.share_promotions} />
              </div>
            {/if}
          {/if}
        </div>

        <!-- Social Media Platforms -->
        {#if form.social_integration_enabled}
          <div class="space-y-4">
            <h3 class="text-lg font-semibold">Connected Platforms</h3>
            
            {#each socialPlatforms as platform}
              <Card.Root class="border">
                <Card.Header class="pb-3">
                  <div class="flex items-center justify-between">
                    <div class="flex items-center space-x-3">
                      <div class="w-10 h-10 rounded-full {platform.color} flex items-center justify-center text-white text-xl">
                        {platform.icon}
                      </div>
                      <div>
                        <Card.Title class="text-lg">{platform.name}</Card.Title>
                        <Card.Description>{platform.description}</Card.Description>
                      </div>
                    </div>
                    <div class="flex items-center space-x-2">
                      {#if isPlatformEnabled(platform.id)}
                        <Badge variant="default">Connected</Badge>
                      {:else}
                        <Badge variant="secondary">Disconnected</Badge>
                      {/if}
                      <Switch 
                        checked={isPlatformEnabled(platform.id)} 
                        onCheckedChange={() => togglePlatform(platform.id)}
                      />
                    </div>
                  </div>
                </Card.Header>
                
                {#if isPlatformEnabled(platform.id)}
                  <Card.Content class="pt-0">
                    <div class="space-y-4">
                      {#each platform.configFields as field}
                        <div>
                          <Label for={platform.id + '-' + field.key}>{field.label}</Label>
                          {#if field.type === 'password'}
                            <Input 
                              id={platform.id + '-' + field.key}
                              type="password"
                              value={getPlatformConfig(platform.id)[field.key] || ''}
                              placeholder={field.placeholder}
                              required={field.required}
                              on:input={(e) => updatePlatformConfig(platform.id, field.key, e.currentTarget.value)}
                            />
                          {:else if field.type === 'boolean'}
                            <div class="flex items-center justify-between">
                              <span class="text-sm text-muted-foreground">{field.label}</span>
                              <Switch 
                                checked={!!getPlatformConfig(platform.id)[field.key]}
                                onCheckedChange={(checked) => updatePlatformConfig(platform.id, field.key, checked)}
                              />
                            </div>
                          {:else}
                            <Input 
                              id={platform.id + '-' + field.key}
                              type="text"
                              value={getPlatformConfig(platform.id)[field.key] || ''}
                              placeholder={field.placeholder}
                              required={field.required}
                              on:input={(e) => updatePlatformConfig(platform.id, field.key, e.currentTarget.value)}
                            />
                          {/if}
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

        <!-- Content Settings -->
        {#if form.social_integration_enabled && form.auto_sharing}
          <div class="space-y-4">
            <h3 class="text-lg font-semibold">Content Settings</h3>
            
            <div>
              <Label for="message-template">Message Template</Label>
              <Textarea 
                id="message-template" 
                bind:value={form.custom_message_template} 
                placeholder="Check out our new product: &#123;product_name&#125; - &#123;product_description&#125; #shopping #ecommerce"
                rows={3}
              />
              <p class="text-muted-foreground text-sm">
                Use placeholders like <code>&#123;product_name&#125;</code>, <code>&#123;product_description&#125;</code>, <code>&#123;price&#125;</code>, <code>&#123;url&#125;</code>
              </p>
            </div>

            <div>
              <Label for="hashtags">Default Hashtags</Label>
              <Input 
                id="hashtags"
                value={form.hashtags ? form.hashtags.join(', ') : ''}
                placeholder="shopping, ecommerce, fashion"
                on:input={(e) => handleHashtagsChange(e.currentTarget.value)}
              />
              <p class="text-muted-foreground text-sm">Comma-separated hashtags to include in posts</p>
            </div>
          </div>
        {/if}

        <!-- Analytics & Insights -->
        <div class="space-y-4">
          <h3 class="text-lg font-semibold">Analytics & Insights</h3>
          <Card.Root class="border">
            <Card.Content class="pt-4">
              <div class="space-y-3">
                <div>
                  <h4 class="font-medium">Social Media Performance</h4>
                  <p class="text-sm text-muted-foreground">
                    Track engagement, reach, and conversions from social media posts.
                  </p>
                </div>
                <div>
                  <h4 class="font-medium">Best Posting Times</h4>
                  <p class="text-sm text-muted-foreground">
                    Get recommendations for optimal posting times based on your audience.
                  </p>
                </div>
                <div>
                  <h4 class="font-medium">Content Performance</h4>
                  <p class="text-sm text-muted-foreground">
                    Analyze which products and content perform best on each platform.
                  </p>
                </div>
              </div>
            </Card.Content>
          </Card.Root>
        </div>

        <Button type="submit" disabled={!hasChanges} class="w-full">
          Update Social Media Settings
        </Button>
      </form>
    {/if}
  </Card.Content>
</Card.Root>
</ComingSoon>