// Enhanced data store for template with backend integration
import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';

// Types aligned with updated JSON format
export interface ProductAttribute {
  title: string;
  value: string;
}

export interface Product {
  id: string;
  productId: number;
  title: string;
  description: string;
  slug: string;
  createdAt: string;
  updatedAt: string;
  defaultVariant: ProductVariant;
  variants: ProductVariant[];
  images: ProductImage[];
  attributes: ProductAttribute[];
}

export interface ProductVariant {
  id: string;
  variationId: number;
  description: string;
  price: number;
  availableQuantity: number;
  stockStatus: string;
  isDefault: boolean;
  attributes: ProductAttribute[];
}

export interface ProductImage {
  url: string;
  altText: string;
}

export interface Shop {
  id: string;
  shop_id: number;
  title: string;
  about: string;
  contactEmail: string;
  contactPhone: string;
  currencyCode: string;
  defaultDomain: string;
  address: {
    address: string;
  };
  images: {
    siteLogo: {
      url: string;
      altText: string | null;
    } | null;
    siteLogoDark: any | null;
    banner: any | null;
    bannerDark: any | null;
    coverImage: any | null;
    coverImageDark: any | null;
    favicon: any | null;
  };
  facebookLink: string | null;
  instagramLink: string | null;
  whatsAppLink: string | null;
  whatsAppNumber: string | null;
  paymentMethods: any[];
  seoTitle: string | null;
  seoDescription: string | null;
  seoKeywords: any[];
  shopProductsCategory: any | null;
}

export interface InventoryData {
  products: {
    edges: Array<{
      node: Product;
    }>;
  };
}

// Configuration
// All API calls are now obsolete; only static JSON files are used for data.
// Remove API_BASE_URL, FALLBACK_SHOP_ID, and related backend config.

const CACHE_DURATION = 5 * 60 * 1000; // 5 minutes

// Cache management
interface CacheEntry<T> {
  data: T;
  timestamp: number;
  ttl: number;
}

class DataCache {
  private cache = new Map<string, CacheEntry<any>>();

  set<T>(key: string, data: T, ttl: number = CACHE_DURATION): void {
	this.cache.set(key, {
	  data,
	  timestamp: Date.now(),
	  ttl
	});
  }

  get<T>(key: string): T | null {
	const entry = this.cache.get(key);
	if (!entry) return null;

	if (Date.now() - entry.timestamp > entry.ttl) {
	  this.cache.delete(key);
	  return null;
	}

	return entry.data;
  }

  clear(): void {
	this.cache.clear();
  }
}

const cache = new DataCache();

// Data stores
export const inventoryData = writable<InventoryData | null>(null);
export const shopData = writable<Shop | null>(null);
export const loadingState = writable<boolean>(false);
export const errorState = writable<string | null>(null);

// Derived stores
export const allProducts = derived(inventoryData, ($data) => 
  $data?.products?.edges?.map(edge => edge.node) || []
);

export const availableProducts = derived(allProducts, ($products) => 
  $products.filter(product => product.defaultVariant.availableQuantity > 0)
);

export const featuredProducts = derived(allProducts, ($products) => 
  $products.filter(product => 
    product.attributes.some(attr => attr.title === 'Featured' && attr.value === 'true')
  )
);

export const categories = derived(allProducts, ($products) => {
  const categorySet = new Set<string>();
  $products.forEach(product => {
    // Extract categories from attributes
    const categoryAttr = product.attributes.find(attr => attr.title === 'Category');
    if (categoryAttr) {
      categorySet.add(categoryAttr.value);
    }
  });
  return Array.from(categorySet).sort();
});

// Products grouped by category
export const productsByCategory = derived(allProducts, ($products) => {
  const categoryMap = new Map<string, Product[]>();
  $products.forEach(product => {
    const categoryAttr = product.attributes.find(attr => attr.title === 'Category');
    const category = categoryAttr?.value || 'Uncategorized';
    if (!categoryMap.has(category)) {
      categoryMap.set(category, []);
    }
    categoryMap.get(category)!.push(product);
  });
  return categoryMap;
});

// Data manager class
class DataManager {
  private static instance: DataManager;

  private constructor() {}

  static getInstance(): DataManager {
	if (!DataManager.instance) {
	  DataManager.instance = new DataManager();
	}
	return DataManager.instance;
  }

  async fetchFromEndpoint<T>(endpoint: string, ttl: number = CACHE_DURATION): Promise<T | null> {
    // This method is obsolete and should not be used. Always return null.
    return null;
  }

  async fetchFromStatic<T>(path: string): Promise<T | null> {
	const cacheKey = `static:${path}`;
	const cached = cache.get<T>(cacheKey);
	
	if (cached) {
	  return cached;
	}

	try {
	  const response = await fetch(path);
	  if (!response.ok) {
		throw new Error(`HTTP ${response.status}: ${response.statusText}`);
	  }

	  const data = await response.json();
	  cache.set(cacheKey, data);
	  return data;
	} catch (error) {
	  console.error(`Failed to fetch static ${path}:`, error);
	  return null;
	}
  }

  // Remove transformation methods since data is already in correct format
  
  getShopIdFromHostname(): number | null {
    // No fallback needed; always expect static JSON files to exist.
    return null;
  }

  async loadAllData(shopId?: number): Promise<void> {
	loadingState.set(true);
	errorState.set(null);

	try {
	  // Only use static files for shop and product data
	  const shopMetadata = await this.fetchFromStatic('/data/shop.json');
	  if (shopMetadata) {
		shopData.set(shopMetadata as Shop);
	  } else {
		throw new Error('No shop data available (static /data/shop.json missing)');
	  }

	  const staticInventory = await this.fetchFromStatic('/data/products.json');
	  if (staticInventory) {
		inventoryData.set(staticInventory as InventoryData);
	  } else {
		// Create empty structure if static file missing
		inventoryData.set({ products: { edges: [] } });
	  }

      // Defensive: check if staticInventory is InventoryData
      let productCount = 0;
      const inv = staticInventory as Partial<InventoryData>;
      if (inv && inv.products && Array.isArray(inv.products.edges)) {
        productCount = inv.products.edges.length;
      }
      console.log(`Loaded ${productCount} products from static data`);
	} catch (error) {
	  console.error('Failed to load data:', error);
	  errorState.set(`Failed to load data: ${error instanceof Error ? error.message : 'Unknown error'}`);
	} finally {
	  loadingState.set(false);
	}
  }

  async refreshData(): Promise<void> {
	cache.clear();
	await this.loadAllData();
  }
}

// Create singleton instance
const dataManager = DataManager.getInstance();

// Export data manager
export { dataManager };

// Additional exports for compatibility with layout
export const initializeData = () => dataManager.loadAllData();
export const isDataLoading = loadingState;
export const dataLoadError = errorState;
export const siteMetadata = shopData;
export const dataLastUpdated = derived(shopData, ($shop) => new Date().toISOString());

// Search functionality
export const searchQuery = writable('');
export const searchResults = derived(
  [allProducts, searchQuery],
  ([$products, $query]) => {
	if (!$query.trim()) return $products;
	
	const searchTerm = $query.toLowerCase();
	return $products.filter(product => 
	  product.title.toLowerCase().includes(searchTerm) ||
	  product.description.toLowerCase().includes(searchTerm) ||
	  product.slug.toLowerCase().includes(searchTerm)
	);
  }
);

// Auto-refresh and initial load
if (browser) {
  // Initial data load
  dataManager.loadAllData();
  
  // Refresh data every 5 minutes
  setInterval(() => {
	dataManager.refreshData();
  }, 5 * 60 * 1000);
}
