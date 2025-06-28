import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';

export interface Shop {
  id: string;
  title: string;
  subdomain: string;
  currency_code: string;
  status: string;
}

function createShopStore() {
  const { subscribe, set, update } = writable<Shop | null>(null);

  return {
    subscribe,
    set,
    update,
    
    // Always load local shop.json in dev/static mode
    initialize: async (): Promise<void> => {
      if (!browser) return;
      try {
        const response = await fetch('/data/shop.json');
        if (response.ok) {
          const shopData = await response.json();
          set(shopData);
        } else {
          console.error('Failed to load local shop.json:', response.statusText);
        }
      } catch (error) {
        console.error('Error loading local shop.json:', error);
      }
    },

    // Get shop ID from subdomain
    getCurrentShopId: (): string | null => {
      if (!browser) return null;
      
      const hostname = window.location.hostname;
      const parts = hostname.split('.');
      
      // For development: localhost, 127.0.0.1, etc. - we'll use a default shop ID
      if (hostname === 'localhost' || hostname.startsWith('127.0.0.1') || hostname.startsWith('192.168')) {
        return '1'; // Default shop ID for development
      }
      
      // For production: extract subdomain
      if (parts.length >= 3) {
        // subdomain.domain.com -> use subdomain to get shop ID
        // We'll need to fetch the shop by subdomain
        return null; // Will be resolved by fetchShopBySubdomain
      }
      
      return null;
    },

    // Fetch shop data by subdomain
    fetchShopBySubdomain: async (subdomain?: string): Promise<any> => {
      try {
        if (!browser) return null;
        const response = await fetch('/data/shop.json');
        if (!response.ok) {
          console.error('Failed to load local shop.json:', response.statusText);
          return null;
        }
        const shopData = await response.json();
        set(shopData);
        return shopData;
      } catch (error) {
        console.error('Error loading local shop.json:', error);
        return null;
      }
    }
  };
}

export const shop = createShopStore();

// Derived store for shop ID
export const shopId = derived(shop, ($shop) => $shop?.id || null);

// Helper function to get current shop ID
function getCurrentShopId(): string | null {
  return null;
}
