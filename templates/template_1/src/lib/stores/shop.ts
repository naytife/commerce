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
    fetchShopBySubdomain: async (subdomain?: string): Promise<Shop | null> => {
      try {
        const actualSubdomain = subdomain || (browser ? window.location.hostname.split('.')[0] : null);
        if (!actualSubdomain) return null;

        const apiUrl = (import.meta.env.VITE_API_URL || 'http://ynt.localhost:8080').replace('/query', '');
        const response = await fetch(`${apiUrl}/shops/subdomain/${actualSubdomain}`);
        
        if (!response.ok) {
          console.error('Failed to fetch shop by subdomain:', response.statusText);
          return null;
        }
        
        const result = await response.json();
        const shop = result.data;
        
        if (shop) {
          set(shop);
          return shop;
        }
        
        return null;
      } catch (error) {
        console.error('Error fetching shop by subdomain:', error);
        return null;
      }
    },

    // Initialize shop from current context
    initialize: async (): Promise<void> => {
      if (!browser) return;
      
      const shopId = getCurrentShopId();
      if (shopId) {
        // Development mode with direct shop ID
        try {
          const apiUrl = (import.meta.env.VITE_API_URL || 'http://ynt.localhost:8080').replace('/query', '');
          const response = await fetch(`${apiUrl}/shops/${shopId}`);
          
          if (response.ok) {
            const result = await response.json();
            set(result.data);
          }
        } catch (error) {
          console.error('Error fetching shop by ID:', error);
        }
      } else {
        // Production mode with subdomain
        const store = createShopStore();
        await store.fetchShopBySubdomain();
      }
    }
  };
}

export const shop = createShopStore();

// Derived store for shop ID
export const shopId = derived(shop, ($shop) => $shop?.id || null);

// Helper function to get current shop ID
function getCurrentShopId(): string | null {
  if (!browser) return null;
  
  const hostname = window.location.hostname;
  
  // For development: localhost, 127.0.0.1, etc. - we'll use a default shop ID
  if (hostname === 'localhost' || hostname.startsWith('127.0.0.1') || hostname.startsWith('192.168')) {
    return '1'; // Default shop ID for development
  }
  
  return null; // Will be resolved by subdomain lookup
}
