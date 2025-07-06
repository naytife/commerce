<script lang="ts">
  import * as Card from '$lib/components/ui/card';
  import { Button } from '$lib/components/ui/button';
  import { Label } from '$lib/components/ui/label';
  import { toast } from 'svelte-sonner';
  import { getContext, onMount } from 'svelte';
  import { page } from '$app/stores';
  import { api } from '$lib/api';
  import type { Shop } from '$lib/types';
  import { X, Upload } from 'lucide-svelte';

  let shop: Partial<Shop> = {};
  let loading = false;

  // State for tracking uploads
  let isLogoUploading = false;
  let isLogoDarkUploading = false;
  let isFaviconUploading = false;
  let isBannerUploading = false;
  let isBannerDarkUploading = false;

  // Type for Shop image fields
  type ShopImageKey = 'logo_url' | 'logo_url_dark' | 'favicon_url' | 'banner_url' | 'banner_url_dark';

  // Get authFetch and refetchShopData from context if available
  const authFetch: (input: RequestInfo | URL, init?: RequestInit | undefined) => Promise<Response> = getContext('authFetch');
  const refetchShopData: (() => Promise<void>) | undefined = getContext('refetchShopData');

  // Fetch shop data on mount
  onMount(async () => {
    loading = true;
    try {
      shop = await api(authFetch).getShop();
    } catch (e) {
      toast.error('Failed to load shop data');
    }
    loading = false;
  });

  // Handle file upload to R2
  async function uploadFileToR2(file: File, type: string) {
    try {
      const formData = new FormData();
      const timestamp = Date.now();
      const filename = `shops/${shop.shop_id}/images/${type}_${timestamp}_${file.name.replace(/\s+/g, '-')}`;
      formData.append('file', file);
      formData.append('filename', filename);
      const response = await fetch('/api/upload-to-r2', {
        method: 'POST',
        body: formData
      });
      const data = await response.json();
      if (!data.success) {
        throw new Error(data.error || 'Error uploading file');
      }
      return data.url;
    } catch (error) {
      console.error('Error uploading file:', error);
      toast.error('Failed to upload image');
      throw error;
    }
  }

  // Update shop with image URL
  async function updateShopImage(imageData: Partial<{ 
    logo_url: string;
    logo_url_dark: string;
    favicon_url: string;
    banner_url: string;
    banner_url_dark: string;
    cover_image_url?: string;
    cover_image_url_dark?: string;
  }>) {
    try {
      await api(authFetch).updateShopImages(imageData);
      if (!shop.images) {
        shop.images = {};
      }
      shop.images = { ...shop.images, ...imageData };
      if (refetchShopData) await refetchShopData();
    } catch (error) {
      console.error('Error updating shop image:', error);
      toast.error('Failed to update shop image');
      throw error;
    }
  }

  // Handle image deletion
  async function deleteImage(imageType: ShopImageKey) {
    try {
      const updateData: Partial<Record<ShopImageKey | 'cover_image_url' | 'cover_image_url_dark', string>> = {};
      updateData[imageType] = "";
      if (imageType === 'banner_url') {
        updateData['cover_image_url'] = "";
      } else if (imageType === 'banner_url_dark') {
        updateData['cover_image_url_dark'] = "";
      }
      await api(authFetch).updateShopImages(updateData);
      if (shop.images) {
        shop.images[imageType] = '';
        if (imageType === 'banner_url') {
          shop.images.cover_image_url = '';
        } else if (imageType === 'banner_url_dark') {
          shop.images.cover_image_url_dark = '';
        }
      }
      if (refetchShopData) await refetchShopData();
      toast.success('Image removed successfully');
    } catch (error) {
      console.error('Error removing image:', error);
      toast.error('Failed to remove image');
    }
  }

  // Upload handlers
  async function handleLogoUpload(event: Event) {
    const input = event.target as HTMLInputElement;
    if (!input.files || input.files.length === 0) return;
    const file = input.files[0];
    isLogoUploading = true;
    try {
      const imageUrl = await uploadFileToR2(file, 'logo');
      await updateShopImage({ logo_url: imageUrl });
      toast.success('Logo uploaded successfully');
    } catch (error) {
      console.error('Error uploading logo:', error);
    } finally {
      isLogoUploading = false;
      input.value = '';
    }
  }
  async function handleLogoDarkUpload(event: Event) {
    const input = event.target as HTMLInputElement;
    if (!input.files || input.files.length === 0) return;
    const file = input.files[0];
    isLogoDarkUploading = true;
    try {
      const imageUrl = await uploadFileToR2(file, 'logo_dark');
      await updateShopImage({ logo_url_dark: imageUrl });
      toast.success('Dark mode logo uploaded successfully');
    } catch (error) {
      console.error('Error uploading dark logo:', error);
    } finally {
      isLogoDarkUploading = false;
      input.value = '';
    }
  }
  async function handleFaviconUpload(event: Event) {
    const input = event.target as HTMLInputElement;
    if (!input.files || input.files.length === 0) return;
    const file = input.files[0];
    isFaviconUploading = true;
    try {
      const imageUrl = await uploadFileToR2(file, 'favicon');
      await updateShopImage({ favicon_url: imageUrl });
      toast.success('Favicon uploaded successfully');
    } catch (error) {
      console.error('Error uploading favicon:', error);
    } finally {
      isFaviconUploading = false;
      input.value = '';
    }
  }
  async function handleBannerUpload(event: Event) {
    const input = event.target as HTMLInputElement;
    if (!input.files || input.files.length === 0) return;
    const file = input.files[0];
    isBannerUploading = true;
    try {
      const imageUrl = await uploadFileToR2(file, 'banner');
      await updateShopImage({ banner_url: imageUrl, cover_image_url: imageUrl });
      toast.success('Banner uploaded successfully');
    } catch (error) {
      console.error('Error uploading banner:', error);
    } finally {
      isBannerUploading = false;
      input.value = '';
    }
  }
  async function handleBannerDarkUpload(event: Event) {
    const input = event.target as HTMLInputElement;
    if (!input.files || input.files.length === 0) return;
    const file = input.files[0];
    isBannerDarkUploading = true;
    try {
      const imageUrl = await uploadFileToR2(file, 'banner_dark');
      await updateShopImage({ banner_url_dark: imageUrl, cover_image_url_dark: imageUrl });
      toast.success('Dark mode banner uploaded successfully');
    } catch (error) {
      console.error('Error uploading dark banner:', error);
    } finally {
      isBannerDarkUploading = false;
      input.value = '';
    }
  }
</script>

<h2 class="text-2xl font-semibold mb-2">Branding</h2>
<p class="mb-6 text-muted-foreground">Upload logo, cover image, and customize your store's branding images.</p>

<div class="space-y-6">
  <!-- Logo Upload -->
  <Card.Root class="bg-gradient">
    <Card.Header>
      <Card.Title>Logo</Card.Title>
      <Card.Description>Upload your store's logo. Recommended size: 200x200px</Card.Description>
    </Card.Header>
    <Card.Content>
      <div class="grid gap-6 sm:grid-cols-2">
        <!-- Standard Logo -->
        <div class="space-y-3">
          <Label>Standard Logo</Label>
          <div class="relative flex justify-center aspect-square w-full max-w-[200px] h-48 border rounded-md overflow-hidden">
            {#if shop.images?.logo_url}
              <img
                src={shop.images.logo_url}
                alt="Store Logo"
                class="object-contain w-full h-full p-2"
              />
              <button 
                class="absolute top-2 right-2 flex h-6 w-6 items-center justify-center rounded-full bg-destructive text-white hover:bg-destructive/90"
                on:click={() => deleteImage('logo_url')}
              >
                <X class="h-4 w-4" />
              </button>
            {:else}
              <label
                class="flex h-full w-full cursor-pointer flex-col items-center justify-center gap-1 rounded-md border border-dashed p-4 hover:bg-muted/50"
              >
                <Upload class="h-6 w-6 text-muted-foreground" />
                <span class="text-xs text-muted-foreground">Upload Logo</span>
                <input 
                  type="file" 
                  accept="image/*" 
                  class="hidden" 
                  on:change={handleLogoUpload}
                  disabled={isLogoUploading}
                />
              </label>
            {/if}
          </div>
          {#if shop.images?.logo_url}
            <div class="flex justify-end">
              <label class="inline-flex items-center gap-2 cursor-pointer text-sm text-primary hover:underline">
                <span>Change Logo</span>
                <input 
                  type="file" 
                  accept="image/*" 
                  class="hidden" 
                  on:change={handleLogoUpload}
                  disabled={isLogoUploading}
                />
              </label>
            </div>
          {/if}
        </div>
        <!-- Dark Mode Logo -->
        <div class="space-y-3">
          <Label>Dark Mode Logo</Label>
          <div class="relative flex justify-center aspect-square w-full max-w-[200px] h-48 border rounded-md overflow-hidden bg-zinc-900">
            {#if shop.images?.logo_url_dark}
              <img
                src={shop.images.logo_url_dark}
                alt="Dark Mode Logo"
                class="object-contain w-full h-full p-2"
              />
              <button 
                class="absolute top-2 right-2 flex h-6 w-6 items-center justify-center rounded-full bg-destructive text-white hover:bg-destructive/90"
                on:click={() => deleteImage('logo_url_dark')}
              >
                <X class="h-4 w-4" />
              </button>
            {:else}
              <label
                class="flex h-full w-full cursor-pointer flex-col items-center justify-center gap-1 rounded-md border border-dashed border-gray-700 p-4 hover:bg-zinc-800"
              >
                <Upload class="h-6 w-6 text-gray-400" />
                <span class="text-xs text-gray-400">Upload Dark Logo</span>
                <input 
                  type="file" 
                  accept="image/*" 
                  class="hidden" 
                  on:change={handleLogoDarkUpload}
                  disabled={isLogoDarkUploading}
                />
              </label>
            {/if}
          </div>
          {#if shop.images?.logo_url_dark}
            <div class="flex justify-end">
              <label class="inline-flex items-center gap-2 cursor-pointer text-sm text-primary hover:underline">
                <span>Change Dark Logo</span>
                <input 
                  type="file" 
                  accept="image/*" 
                  class="hidden" 
                  on:change={handleLogoDarkUpload}
                  disabled={isLogoDarkUploading}
                />
              </label>
            </div>
          {/if}
        </div>
      </div>
    </Card.Content>
  </Card.Root>
  <!-- Favicon Upload -->
  <Card.Root class="bg-gradient">
    <Card.Header>
      <Card.Title>Favicon</Card.Title>
      <Card.Description>Upload your store's favicon. Recommended size: 32x32px or 64x64px</Card.Description>
    </Card.Header>
    <Card.Content>
      <div class="space-y-3">
        <div class="relative flex justify-center aspect-square w-16 h-16 border rounded-md overflow-hidden">
          {#if shop.images?.favicon_url}
            <img
              src={shop.images.favicon_url}
              alt="Favicon"
              class="object-contain w-full h-full p-1"
            />
            <button 
              class="absolute top-1 right-1 flex h-4 w-4 items-center justify-center rounded-full bg-destructive text-white hover:bg-destructive/90"
              on:click={() => deleteImage('favicon_url')}
            >
              <X class="h-2 w-2" />
            </button>
          {:else}
            <label
              class="flex h-full w-full cursor-pointer flex-col items-center justify-center gap-1 rounded-md border border-dashed p-1 hover:bg-muted/50"
            >
              <Upload class="h-4 w-4 text-muted-foreground" />
              <input 
                type="file" 
                accept="image/*" 
                class="hidden" 
                on:change={handleFaviconUpload}
                disabled={isFaviconUploading}
              />
            </label>
          {/if}
        </div>
        {#if shop.images?.favicon_url}
          <div class="flex">
            <label class="inline-flex items-center gap-2 cursor-pointer text-sm text-primary hover:underline">
              <span>Change Favicon</span>
              <input 
                type="file" 
                accept="image/*" 
                class="hidden" 
                on:change={handleFaviconUpload}
                disabled={isFaviconUploading}
              />
            </label>
          </div>
        {/if}
      </div>
    </Card.Content>
  </Card.Root>
  <!-- Banner/Cover Images -->
  <Card.Root class="bg-gradient">
    <Card.Header>
      <Card.Title>Banner/Cover Images</Card.Title>
      <Card.Description>Upload your store's banner and cover images. Recommended size: 1600x400px</Card.Description>
    </Card.Header>
    <Card.Content>
      <div class="grid gap-6 sm:grid-cols-2">
        <!-- Standard Banner -->
        <div class="space-y-3">
          <Label>Standard Banner</Label>
          <div class="relative flex justify-center w-full h-[200px] border rounded-md overflow-hidden">
            {#if shop.images?.banner_url}
              <img
                src={shop.images.banner_url}
                alt="Store Banner"
                class="object-cover w-full h-full"
              />
              <button 
                class="absolute top-2 right-2 flex h-6 w-6 items-center justify-center rounded-full bg-destructive text-white hover:bg-destructive/90"
                on:click={() => deleteImage('banner_url')}
              >
                <X class="h-4 w-4" />
              </button>
            {:else}
              <label
                class="flex h-full w-full cursor-pointer flex-col items-center justify-center gap-1 rounded-md border border-dashed p-4 hover:bg-muted/50"
              >
                <Upload class="h-6 w-6 text-muted-foreground" />
                <span class="text-xs text-muted-foreground">Upload Banner</span>
                <input 
                  type="file" 
                  accept="image/*" 
                  class="hidden" 
                  on:change={handleBannerUpload}
                  disabled={isBannerUploading}
                />
              </label>
            {/if}
          </div>
          {#if shop.images?.banner_url}
            <div class="flex justify-end">
              <label class="inline-flex items-center gap-2 cursor-pointer text-sm text-primary hover:underline">
                <span>Change Banner</span>
                <input 
                  type="file" 
                  accept="image/*" 
                  class="hidden" 
                  on:change={handleBannerUpload}
                  disabled={isBannerUploading}
                />
              </label>
            </div>
          {/if}
        </div>
        <!-- Dark Mode Banner -->
        <div class="space-y-3">
          <Label>Dark Mode Banner</Label>
          <div class="relative flex justify-center w-full h-[200px] border rounded-md overflow-hidden bg-zinc-900">
            {#if shop.images?.banner_url_dark}
              <img
                src={shop.images.banner_url_dark}
                alt="Dark Mode Banner"
                class="object-cover w-full h-full"
              />
              <button 
                class="absolute top-2 right-2 flex h-6 w-6 items-center justify-center rounded-full bg-destructive text-white hover:bg-destructive/90"
                on:click={() => deleteImage('banner_url_dark')}
              >
                <X class="h-4 w-4" />
              </button>
            {:else}
              <label
                class="flex h-full w-full cursor-pointer flex-col items-center justify-center gap-1 rounded-md border border-dashed border-gray-700 p-4 hover:bg-zinc-800"
              >
                <Upload class="h-6 w-6 text-gray-400" />
                <span class="text-xs text-gray-400">Upload Dark Banner</span>
                <input 
                  type="file" 
                  accept="image/*" 
                  class="hidden" 
                  on:change={handleBannerDarkUpload}
                  disabled={isBannerDarkUploading}
                />
              </label>
            {/if}
          </div>
          {#if shop.images?.banner_url_dark}
            <div class="flex justify-end">
              <label class="inline-flex items-center gap-2 cursor-pointer text-sm text-primary hover:underline">
                <span>Change Dark Banner</span>
                <input 
                  type="file" 
                  accept="image/*" 
                  class="hidden" 
                  on:change={handleBannerDarkUpload}
                  disabled={isBannerDarkUploading}
                />
              </label>
            </div>
          {/if}
        </div>
      </div>
      <p class="mt-2 text-sm text-muted-foreground">These banner images will also be used as cover images for your store.</p>
    </Card.Content>
  </Card.Root>
</div>