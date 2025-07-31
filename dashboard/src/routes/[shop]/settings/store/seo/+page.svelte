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
	import ComingSoon from '$lib/components/ui/coming-soon.svelte';

  const shopStore = getContext('shop') as Writable<Partial<Shop>>;
  const authFetch = getContext('authFetch') as typeof fetch;
  const refetchShopData = getContext('refetchShopData') as () => Promise<void>;

  let shop: Partial<Shop> = {};
  $: shop = $shopStore;

  let form: { seo_title?: string; seo_description?: string; seo_keywords?: string[] } = {};
  let initialForm: { seo_title?: string; seo_description?: string; seo_keywords?: string[] } = {};
  let initialized = false;

  $: if (!initialized && shop && Object.keys(shop).length > 0) {
    form = {
      seo_title: shop.seo_title || '',
      seo_description: shop.seo_description || '',
      seo_keywords: Array.isArray(shop.seo_keywords) ? shop.seo_keywords : (shop.seo_keywords ? [shop.seo_keywords] : [])
    };
    initialForm = deepClone(form);
    initialized = true;
  }

  let hasChanges = false;
  $: hasChanges = Object.keys(form).length > 0 && !deepEqual(form, initialForm);

  async function handleFormSubmit(event: Event) {
    event.preventDefault();
    try {
      await api(authFetch).updateShop({
        seo_title: form.seo_title,
        seo_description: form.seo_description,
        seo_keywords: form.seo_keywords
      });
      await refetchShopData();
      initialForm = deepClone(form);
      toast.success('SEO settings updated!');
    } catch (error) {
      toast.error('Failed to update SEO settings');
    }
  }
</script>
<ComingSoon>
<Card.Root class="bg-gradient w-full max-w-screen-md mx-auto">
  <Card.Header>
    <Card.Title>SEO & Visibility</Card.Title>
    <Card.Description>Add meta tags, sitemap config, and control search engine indexing.</Card.Description>
  </Card.Header>
  <Card.Content>
    {#if !shop || Object.keys(shop).length === 0}
      <p>Loading...</p>
    {:else}
      <form class="space-y-4 w-full" on:submit={handleFormSubmit}>
        <div>
          <Label for="seo-title">SEO Title</Label>
          <Input id="seo-title" bind:value={form.seo_title} placeholder="Your store SEO title" />
          <p class="text-muted-foreground text-sm">Title that appears in search engine results.</p>
        </div>
        <div>
          <Label for="seo-description">SEO Description</Label>
          <Textarea id="seo-description" bind:value={form.seo_description} placeholder="Brief description of your store for search engines" />
          <p class="text-muted-foreground text-sm">Description that appears in search engine results.</p>
        </div>
        <div>
          <Label for="seo-keywords">SEO Keywords</Label>
          <Input id="seo-keywords"
            value={form.seo_keywords ? form.seo_keywords.join(', ') : ''}
            placeholder="Comma-separated keywords"
            on:input={(e) => {
              form.seo_keywords = e.currentTarget.value.split(',').map(k => k.trim()).filter(Boolean);
            }}
          />
          <p class="text-muted-foreground text-sm">Keywords to help with search engine indexing.</p>
        </div>
        <Button type="submit" disabled={!hasChanges}>Update SEO Settings</Button>
      </form>
    {/if}
  </Card.Content>
</Card.Root> 
</ComingSoon>