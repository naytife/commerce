<script lang="ts">
  import { onMount } from 'svelte';
  import * as Card from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';
  import { Button } from '$lib/components/ui/button';
  import { Label } from '$lib/components/ui/label';
  import { CheckCircle, Globe, Lock } from 'lucide-svelte';
  import { writable } from 'svelte/store';
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

  let form: { subdomain?: string; custom_domain?: string } = {};
  let initialForm: { subdomain?: string; custom_domain?: string } = {};
  let initialized = false;

  $: if (!initialized && shop && Object.keys(shop).length > 0) {
    form = {
      subdomain: shop.subdomain || '',
      custom_domain: shop.custom_domain || ''
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
        subdomain: form.subdomain,
        custom_domain: form.custom_domain
      });
      await refetchShopData();
      initialForm = deepClone(form);
      toast.success('Domain settings updated!');
    } catch (error) {
      toast.error('Failed to update domain settings');
    }
  }

  // Mock data for domains
  $: liveDomain = `${shop.subdomain || 'yourstore'}.naytife.com`;
  let liveStatus = 'Secure & Live';
  let customDomains = [
    { domain: 'merchantbrand.com', status: 'Verifying', ssl: 'CNAME', actions: 'Pending' },
    { domain: 'www.awesomebrand.shop', status: 'Active', ssl: 'Active', actions: 'Active' }
  ];

  let showAddForm = false;
  let newDomain = '';
  let dnsCname = 'aweombrand.myplatform.com';
  let dnsTxtName = '_verify.myplatfrm.com';
  let dnsTxtValue = 'abc1234-verification-token';

  function openAddForm() {
    showAddForm = true;
  }
  function closeAddForm() {
    showAddForm = false;
    newDomain = '';
  }
  function submitDomain() {
    // Placeholder for domain submission logic
    showAddForm = false;
    newDomain = '';
    // Show toast or feedback here
  }
</script>

<ComingSoon>
  <Card.Root class="bg-gradient w-full max-w-screen-md mx-auto">
    <Card.Header>
      <div class="flex items-center gap-2 mb-2">
        <Lock class="w-5 h-5 text-muted-foreground" />
        <Card.Title>Custom Domain Configuration</Card.Title>
      </div>
      <Card.Description>
        <div class="flex items-center gap-2 mb-2">
          <Globe class="w-4 h-4 text-muted-foreground" />
          <span>Your store is live at:</span>
        </div>
        <div class="ml-6 text-sm font-mono">https://{liveDomain}</div>
        <div class="flex items-center gap-2 ml-6 mt-1">
          <CheckCircle class="w-4 h-4 text-green-600" />
          <span class="text-green-700">{liveStatus}</span>
        </div>
        <div class="flex items-center gap-2 mt-4">
          <Globe class="w-4 h-4 text-muted-foreground" />
          <span>Want to use your own domain?</span>
        </div>
        <Button class="mt-2" on:click={openAddForm} disabled={showAddForm}>+ Add Custom Domain</Button>
      </Card.Description>
    </Card.Header>
    <Card.Content>
      <div class="mb-6">
        <h3 class="font-semibold mb-2">Add Custom Domain</h3>
        <table class="w-full text-sm mb-4">
          <thead>
            <tr class="border-b">
              <th class="text-left py-1">Custom Domain</th>
              <th class="text-left py-1">Status</th>
              <th class="text-left py-1">SSL</th>
              <th class="text-left py-1">Actions</th>
            </tr>
          </thead>
          <tbody>
            {#each customDomains as d}
              <tr class="border-b last:border-0">
                <td class="py-1">{d.domain}</td>
                <td class="py-1">
                  {#if d.status === 'Active'}
                    <span class="inline-flex items-center gap-1 text-green-700"><CheckCircle class="w-4 h-4" /> Active</span>
                  {:else if d.status === 'Verifying'}
                    <span class="inline-flex items-center gap-1 text-yellow-600"><Globe class="w-4 h-4" /> Verifying</span>
                  {:else}
                    <span>{d.status}</span>
                  {/if}
                </td>
                <td class="py-1">{d.ssl}</td>
                <td class="py-1">{d.actions}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>

      {#if showAddForm}
        <form class="space-y-4 max-w-xl mb-6" on:submit|preventDefault={submitDomain}>
          <div>
            <Label for="custom-domain">Custom Domain</Label>
            <Input id="custom-domain" bind:value={newDomain} placeholder="yourdomain.com" />
          </div>
          <div class="bg-muted rounded p-3 text-xs">
            <div class="mb-2 font-semibold">Add the following DNS records to your domain:</div>
            <ul class="list-disc ml-5">
              <li class="mb-1">
                <span class="font-semibold">CNAME Record</span><br />
                Name <span class="font-mono">CNAME</span> <span class="font-mono">{dnsCname}</span><br />
                Type <span class="font-mono">CNAME</span>
              </li>
              <li>
                <span class="font-semibold">TXT Record (for verification)</span><br />
                Name <span class="font-mono">{dnsTxtName}</span><br />
                Value <span class="font-mono">{dnsTxtValue}</span>
              </li>
            </ul>
            <div class="mt-2 text-muted-foreground">Changes may take up to 24 hours to propagate.</div>
          </div>
          <div class="flex gap-2">
            <Button type="submit">Submit Domain</Button>
            <Button type="button" variant="secondary" on:click={closeAddForm}>Cancel</Button>
          </div>
        </form>
      {/if}
    </Card.Content>
  </Card.Root>
</ComingSoon>