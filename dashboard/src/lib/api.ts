import type { 
  Product, 
  Shop, 
  ProductType, 
  ApiResponse, 
  ProductTypeAttribute, 
  ProductImage, 
  Order,
  Customer,
  CustomerCreatePayload,
  CustomerUpdatePayload,
  CustomerSearchParams,
  InventoryItem,
  InventoryReport,
  LowStockVariant,
  StockMovement,
  StockUpdatePayload,
  StockMovementCreatePayload,
  InventorySearchParams,
  StockMovementSearchParams,
  PaginatedResponse,
  PaymentMethod,
  PaymentMethodConfig,
  PredefinedProductType,
  ProductTypeWithTemplateResponse
} from './types'
import { toast } from 'svelte-sonner'
import { writable } from 'svelte/store'
import { publishState } from './stores/publishState'
import { createAuthenticatedFetch } from './auth-fetch'

// Create a store to hold the shop ID
export const currentShopId = writable<number | null>(null)

// Function to fetch shop ID from path
export const fetchShopIdFromSubdomain = async (shopIdentifier: string, customFetch = fetch) => {
  try {
    if (!shopIdentifier) {
      throw new Error('No shop identifier provided')
    }
    
    const authenticatedFetch = createAuthenticatedFetch(customFetch);
    const response = await authenticatedFetch(
      `http://127.0.0.1:8080/v1/subdomains/${shopIdentifier}`,
    )
    const data = await response.json() as ApiResponse<Shop>
    if (data.data && data.data.shop_id) {
      currentShopId.set(data.data.shop_id)
      return data.data.shop_id
    } else {
      throw new Error('Shop not found')
    }
  } catch (error) {
    console.error('Error fetching shop ID:', error)
    toast.error('Failed to load shop information')
    return null
  }
}

// Fetch all shops (for account page, not shop-specific)
export const getAllShops = async (customFetch = fetch) => {
  try {
    const authenticatedFetch = createAuthenticatedFetch(customFetch);
    const response = await authenticatedFetch('http://127.0.0.1:8080/v1/shops');
    const data = await response.json() as ApiResponse<Shop[]>;
    return Array.isArray(data.data) ? data.data : [];
  } catch (error) {
    console.error('Error fetching shops:', error);
    toast.error('Failed to load shops');
    return [];
  }
}

// Create a new shop (store)
export const createShop = async (
  shop: { subdomain: string; title: string; template?: string },
  customFetch = fetch
) => {
  try {
    const authenticatedFetch = createAuthenticatedFetch(customFetch);
    const response = await authenticatedFetch('http://127.0.0.1:8080/v1/shops', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        subdomain: shop.subdomain,
        title: shop.title,
        template: shop.template || 'template_1', // Default to template_1 if not provided
        currency_code: 'USD', // Default currency
        status: 'DRAFT' // Default status
      }),
    });
    const data = await response.json() as ApiResponse<Shop>;
    if (response.ok) {
      toast.success('Store created successfully');
      return data.data;
    } else {
      toast.error(data.message || 'Failed to create store');
      throw new Error(data.message || 'Failed to create store');
    }
  } catch (error) {
    toast.error('An error occurred while creating the store');
    throw error;
  }
};

// Delete a shop (store)
export const deleteShop = async (
  shopId: number,
  customFetch = fetch
) => {
  try {
    const authenticatedFetch = createAuthenticatedFetch(customFetch);
    const response = await authenticatedFetch(`http://127.0.0.1:8080/v1/shops/${shopId}`, {
      method: 'DELETE',
    });
    if (response.ok) {
      toast.success('Store deleted successfully');
      return true;
    } else {
      const data = await response.json().catch(() => ({ message: 'Failed to delete store' }));
      toast.error(data.message || 'Failed to delete store');
      return false;
    }
  } catch (error) {
    toast.error('An error occurred while deleting the store');
    throw error;
  }
};

// Check subdomain availability
export const checkSubdomainAvailability = async (
  subdomain: string,
  customFetch = fetch
) => {
  try {
    const authenticatedFetch = createAuthenticatedFetch(customFetch);
    const response = await authenticatedFetch(`http://127.0.0.1:8080/v1/subdomains/${subdomain}/check`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    const data = await response.json() as ApiResponse<{ subdomain: string; available: boolean; message: string }>;
    if (response.ok) {
      return data.data;
    } else {
      throw new Error(data.message || 'Failed to check subdomain availability');
    }
  } catch (error) {
    console.error('Error checking subdomain availability:', error);
    throw error;
  }
};

export const api = (customFetch = fetch) => {
  // Create authenticated fetch wrapper
  const authenticatedFetch = createAuthenticatedFetch(customFetch);
  
  // Helper to get base URL with current shop ID
  const getShopUrl = () => {
    let shopId: number | null = null
    
    // Use a synchronous approach to get the current value
    currentShopId.subscribe(id => {
      shopId = id
    })()
    
    if (shopId === null) {
      throw new Error('Shop ID not set. Call fetchShopIdFromSubdomain first.')
    }
    
    return `http://127.0.0.1:8080/v1/shops/${shopId}`
  }

  return {
    getProducts: async (limit: number) => {
      const response = await authenticatedFetch(
        `${getShopUrl()}/products`,
      )
      const data = await response.json() as ApiResponse<Product[]>
      return data.data
    },
    getProductsByType: async (typeId: number, limit: number = 10, after: number = 0) => {
      const url = `${getShopUrl()}/product-types/${typeId}/products?after=${after}&limit=${limit}`
      const response = await authenticatedFetch(url)
      const data = await response.json() as ApiResponse<Product[]>
      return data.data
    },
    getProductById: async (id: number): Promise<Product> => {
      const response = await authenticatedFetch(
        `${getShopUrl()}/products/${id}`,
      );
      const data = await response.json() as ApiResponse<Product>;
      return data.data;
    },
    getProductImages: async (productId: string | number): Promise<ProductImage[]> => {
      const response = await customFetch(
        `${getShopUrl()}/products/${productId}/images`,
      );
      const data = await response.json() as ApiResponse<{ images: ProductImage[] }>;
      return data.data.images;
    },
    getProductTypeAttributes: async (typeId: number): Promise<ProductTypeAttribute[]> => {
      const response = await customFetch(
        `${getShopUrl()}/product-types/${typeId}/attributes`,
      );
      const data = await response.json() as ApiResponse<ProductTypeAttribute[]>;
      return data.data;
    },
    getShop: async () => {
      const response = await authenticatedFetch(
        `${getShopUrl()}`,
      );
      const data = await response.json() as ApiResponse<Shop>;
      return data.data;
    },
    getProductTypes: async () => {
      const response = await authenticatedFetch(
        `${getShopUrl()}/product-types`,
      );
      const data = await response.json() as ApiResponse<ProductType[]>;
      return data.data;
    },
    getProductTypeById: async (typeId: number): Promise<ProductType> => {
      const response = await customFetch(
        `${getShopUrl()}/product-types/${typeId}`,
      );
      const data = await response.json() as ApiResponse<ProductType>;
      return data.data;
    },
    createProductTypeAttribute: async (typeId: number, attribute: ProductTypeAttribute) => {
      try {
        const response = await customFetch(
          `${getShopUrl()}/product-types/${typeId}/attributes`,
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(attribute),
          }
        );
        const data = await response.json() as ApiResponse<ProductTypeAttribute>;
        
        if (response.ok) {
          toast.success('Attribute created successfully');
        } else {
          toast.error('Failed to create attribute');
        }
        
        return data.data;
      } catch (error) {
        toast.error('An error occurred while creating the attribute');
        throw error;
      }
    },
    deleteProductTypeAttribute: async (attributeId: number) => {
      try {
        const response = await customFetch(
          `${getShopUrl()}/attributes/${attributeId}`,
          {
            method: 'DELETE',
          }
        );
        const data = await response.json() as ApiResponse<ProductTypeAttribute>;
        
        if (response.ok) {
          toast.success('Attribute deleted successfully');
        } else {
          toast.error('Failed to delete attribute');
        }
        
        return data.data;
      } catch (error) {
        toast.error('An error occurred while deleting the attribute');
        throw error;
      }
    },
    createProductType: async (data: { title: string; sku_substring: string; shippable: boolean; digital: boolean }) => {
      try {
        const response = await customFetch(
          `${getShopUrl()}/product-types`,
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
          }
        );
        const responseData = await response.json() as ApiResponse<ProductType>;
        return responseData.data;
      } catch (error) {
        console.error('Error creating product type:', error);
        throw error;
      }
    },
    // Predefined product type functions
    getPredefinedProductTypes: async (): Promise<PredefinedProductType[]> => {
      try {
        const response = await customFetch('http://127.0.0.1:8080/v1/predefined-product-types');
        const data = await response.json() as ApiResponse<PredefinedProductType[]>;
        return data.data;
      } catch (error) {
        console.error('Error fetching predefined product types:', error);
        throw error;
      }
    },
    getPredefinedProductType: async (templateId: string): Promise<PredefinedProductType> => {
      try {
        const response = await customFetch(`http://127.0.0.1:8080/v1/predefined-product-types/${templateId}`);
        const data = await response.json() as ApiResponse<PredefinedProductType>;
        return data.data;
      } catch (error) {
        console.error('Error fetching predefined product type:', error);
        throw error;
      }
    },
    createProductTypeFromTemplate: async (templateId: string): Promise<ProductTypeWithTemplateResponse> => {
      try {
        const response = await customFetch(
          `${getShopUrl()}/product-types/from-template`,
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ template_id: templateId }),
          }
        );
        const responseData = await response.json() as ApiResponse<ProductTypeWithTemplateResponse>;
        
        if (!response.ok) {
          throw new Error(responseData.message || 'Failed to create product type from template');
        }
        
        return responseData.data;
      } catch (error) {
        console.error('Error creating product type:', error);
        throw error;
      }
    },
    updateProductTypeAttribute: async (typeId: number, attributeId: number, data: Partial<ProductTypeAttribute>) => {
      try {
        const response = await customFetch(
          `${getShopUrl()}/attributes/${attributeId}`,
          {
            method: 'PUT',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
          }
        );
        const responseData = await response.json() as ApiResponse<ProductTypeAttribute>;
        return responseData.data;
      } catch (error) {
        console.error('Error updating attribute:', error);
        throw error;
      }
    },
    deleteProductType: async (typeId: number): Promise<void> => {
      await customFetch(
        `${getShopUrl()}/product-types/${typeId}`,
        {
          method: 'DELETE',
        }
      );
    },
    updateProductType: async (typeId: number, data: { title: string; sku_substring: string; shippable: boolean; digital: boolean }) => {
      try {
        const response = await customFetch(
          `${getShopUrl()}/product-types/${typeId}`,
          {
            method: 'PUT',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
          }
        );
        const responseData = await response.json() as ApiResponse<ProductType>;
        return responseData.data;
      } catch (error) {
        console.error('Error updating product type:', error);
        throw error;
      }
    },
    createProduct: async (typeId: number, productData: any) => {
      try {
        const response = await customFetch(
          `${getShopUrl()}/product-types/${typeId}/products`,
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(productData),
          }
        );
        const responseData = await response.json();
        if (response.ok) {
          toast.success('Product created successfully');
          
          // Track the change for publish state
          publishState.recordChange({
            type: 'product_create',
            entity: `product:${productData.title}`,
            description: `Product "${productData.title}" created`
          });
          
        } else {
          toast.error(responseData.message || 'Failed to create product');
        }
        return responseData.data;
      } catch (error) {
        toast.error('An error occurred while creating the product');
        throw error;
      }
    },
    updateProduct: async (productId: number, productData: any) => {
      try {
        const response = await customFetch(
          `${getShopUrl()}/products/${productId}`,
          {
            method: 'PUT',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(productData),
          }
        );
        const responseData = await response.json();
        if (response.ok) {
          toast.success('Product updated successfully');
          
          // Track the change for publish state
          publishState.recordChange({
            type: 'product_update',
            entity: `product:${productId}`,
            description: `Product updated`
          });
          
        } else {
          toast.error('Failed to update product');
        }
        return responseData.data;
      } catch (error) {
        toast.error('An error occurred while updating the product');
        throw error;
      }
    },
    addProductImage: async (productId: string | number, imageData: { url: string, alt: string, filename?: string, is_primary?: boolean }) => {
      try {
        // Convert to number if it's a string
        const id = typeof productId === 'string' ? parseInt(productId) : productId;
        
        if (isNaN(id)) {
          toast.error('Invalid product ID');
          throw new Error('Invalid product ID');
        }

        const response = await customFetch(
          `${getShopUrl()}/products/${id}/images`,
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              url: imageData.url,
              alt: imageData.alt
            }),
          }
        );
        const responseData = await response.json();
        if (response.ok) {
          toast.success('Image added successfully');
        } else {
          toast.error(`Failed to add image: ${responseData.message || 'Unknown error'}`);
          throw new Error(responseData.message || 'Failed to add image');
        }
        return responseData.data;
      } catch (error) {
        toast.error('An error occurred while adding the image');
        throw error;
      }
    },
    deleteProductImage: async (productId: string | number, imageId: number) => {
      try {
        // Convert to number if it's a string
        const id = typeof productId === 'string' ? parseInt(productId) : productId;
        
        if (isNaN(id)) {
          toast.error('Invalid product ID');
          throw new Error('Invalid product ID');
        }

        if (!imageId || typeof imageId !== 'number' || isNaN(imageId)) {
          toast.error('Invalid image ID');
          throw new Error('Invalid image ID');
        }

        const response = await customFetch(
          `${getShopUrl()}/products/${id}/images/${imageId}`,
          {
            method: 'DELETE',
          }
        );
        
        if (response.ok) {
          toast.success('Image deleted successfully');
          return true;
        } else {
          const errorData = await response.json().catch(() => ({ message: 'Unknown error' }));
          toast.error(`Failed to delete image: ${errorData.message || 'Unknown error'}`);
          return false;
        }
      } catch (error) {
        toast.error('An error occurred while deleting the image');
        throw error;
      }
    },
    setProductPrimaryImage: async (productId: string | number, imageId: number) => {
      try {
        // Convert to number if it's a string
        const id = typeof productId === 'string' ? parseInt(productId) : productId;
        
        if (isNaN(id)) {
          toast.error('Invalid product ID');
          throw new Error('Invalid product ID');
        }

        if (!imageId || typeof imageId !== 'number' || isNaN(imageId)) {
          toast.error('Invalid image ID');
          throw new Error('Invalid image ID');
        }

        const response = await customFetch(
          `${getShopUrl()}/products/${id}/images/${imageId}/primary`,
          {
            method: 'PUT',
          }
        );
        
        if (response.ok) {
          toast.success('Primary image updated');
          return true;
        } else {
          const errorData = await response.json().catch(() => ({ message: 'Unknown error' }));
          toast.error(`Failed to update primary image: ${errorData.message || 'Unknown error'}`);
          return false;
        }
      } catch (error) {
        toast.error('An error occurred while updating the primary image');
        throw error;
      }
    },
    updateShop: async ( data: Partial<Shop>) => {
      try {
        const response = await customFetch(
          `${getShopUrl()}`,
          {
            method: 'PUT',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
          }
        );
        
        const responseData = await response.json();
        
        if (response.ok) {
          toast.success('Shop settings updated successfully');
          
          // Track the change for publish state
          publishState.recordChange({
            type: 'shop_update',
            entity: 'shop',
            description: 'Store settings updated'
          });
          
          return responseData.data;
        } else {
          toast.error(`Failed to update shop settings: ${responseData.message || 'Unknown error'}`);
          throw new Error(responseData.message || 'Failed to update shop settings');
        }
      } catch (error) {
        toast.error('An error occurred while updating shop settings');
        throw error;
      }
    },
    updateShopImages: async (imageData: {
      banner_url?: string;
      banner_url_dark?: string;
      cover_image_url?: string;
      cover_image_url_dark?: string;
      favicon_url?: string;
      logo_url?: string;
      logo_url_dark?: string;
    }) => {
      try {
        const response = await customFetch(
          `${getShopUrl()}/images`,
          {
            method: 'PUT',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(imageData),
          }
        );
        
        const responseData = await response.json();
        
        if (response.ok) {
          toast.success('Shop images updated successfully');
          
          // Track the change for publish state
          publishState.recordChange({
            type: 'image_update',
            entity: 'shop',
            description: 'Store images updated'
          });
          
          return responseData.data;
        } else {
          toast.error(`Failed to update shop images: ${responseData.message || 'Unknown error'}`);
          throw new Error(responseData.message || 'Failed to update shop images');
        }
      } catch (error) {
        toast.error('An error occurred while updating shop images');
        throw error;
      }
    },
    // Order related functions
    getOrders: async () => {
      const response = await authenticatedFetch(
        `${getShopUrl()}/orders`,
      );
      const data = await response.json() as ApiResponse<Order[]>;
      return data.data;
    },
    getOrderById: async (orderId: number): Promise<Order> => {
      const response = await authenticatedFetch(
        `${getShopUrl()}/orders/${orderId}`,
      );
      const data = await response.json() as ApiResponse<Order>;
      return data.data;
    },
    deleteOrder: async (orderId: number): Promise<void> => {
      try {
        const response = await customFetch(
          `${getShopUrl()}/orders/${orderId}`,
          {
            method: 'DELETE',
          }
        );
        
        if (!response.ok) {
          const errorData = await response.json().catch(() => ({ message: 'Unknown error' }));
          throw new Error(errorData.message || 'Failed to delete order');
        }
      } catch (error) {
        toast.error('An error occurred while deleting the order');
        throw error;
      }
    },
    updateOrderStatus: async (orderId: number, status: string): Promise<Order> => {
      try {
        const response = await customFetch(
          `${getShopUrl()}/orders/${orderId}/status`,
          {
            method: 'PATCH',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ status }),
          }
        );
        
        if (!response.ok) {
          const errorData = await response.json().catch(() => ({ message: 'Unknown error' }));
          throw new Error(errorData.message || 'Failed to update order status');
        }
        
        const data = await response.json() as ApiResponse<Order>;
        return data.data;
      } catch (error) {
        toast.error('An error occurred while updating the order status');
        throw error;
      }
    },

    // Customer Management API Functions
    getCustomers: async (params?: CustomerSearchParams): Promise<PaginatedResponse<Customer>> => {
      const queryParams = new URLSearchParams();
      if (params) {
        Object.entries(params).forEach(([key, value]) => {
          if (value !== undefined && value !== null) {
            queryParams.append(key, value.toString());
          }
        });
      }
      
      const url = `${getShopUrl()}/customers${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
      const response = await authenticatedFetch(url);
      const data = await response.json() as ApiResponse<PaginatedResponse<Customer>>;
      return data.data;
    },

    getCustomerById: async (customerId: number): Promise<Customer> => {
      const response = await authenticatedFetch(`${getShopUrl()}/customers/${customerId}`);
      const data = await response.json() as ApiResponse<Customer>;
      return data.data;
    },

    createCustomer: async (customerData: CustomerCreatePayload): Promise<Customer> => {
      try {
        const response = await authenticatedFetch(`${getShopUrl()}/customers`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(customerData),
        });
        
        const data = await response.json() as ApiResponse<Customer>;
        
        if (response.ok) {
          toast.success('Customer created successfully');
          return data.data;
        } else {
          toast.error(data.message || 'Failed to create customer');
          throw new Error(data.message || 'Failed to create customer');
        }
      } catch (error) {
        toast.error('An error occurred while creating the customer');
        throw error;
      }
    },

    updateCustomer: async (customerId: number, customerData: CustomerUpdatePayload): Promise<Customer> => {
      try {
        const response = await authenticatedFetch(`${getShopUrl()}/customers/${customerId}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(customerData),
        });
        
        const data = await response.json() as ApiResponse<Customer>;
        
        if (response.ok) {
          toast.success('Customer updated successfully');
          return data.data;
        } else {
          toast.error(data.message || 'Failed to update customer');
          throw new Error(data.message || 'Failed to update customer');
        }
      } catch (error) {
        toast.error('An error occurred while updating the customer');
        throw error;
      }
    },

    deleteCustomer: async (customerId: number): Promise<void> => {
      try {
        const response = await customFetch(`${getShopUrl()}/customers/${customerId}`, {
          method: 'DELETE',
        });
        
        if (response.ok) {
          toast.success('Customer deleted successfully');
        } else {
          const errorData = await response.json().catch(() => ({ message: 'Unknown error' }));
          toast.error(errorData.message || 'Failed to delete customer');
          throw new Error(errorData.message || 'Failed to delete customer');
        }
      } catch (error) {
        toast.error('An error occurred while deleting the customer');
        throw error;
      }
    },

    searchCustomers: async (query: string, limit: number = 20): Promise<Customer[]> => {
      const params = new URLSearchParams({
        query,
        limit: limit.toString()
      });
      
      const response = await customFetch(`${getShopUrl()}/customers/search?${params.toString()}`);
      const data = await response.json() as ApiResponse<Customer[]>;
      return data.data;
    },

    // Inventory Management API Functions
    getInventoryReport: async (params?: InventorySearchParams): Promise<InventoryReport> => {
      const queryParams = new URLSearchParams();
      if (params) {
        Object.entries(params).forEach(([key, value]) => {
          if (value !== undefined && value !== null) {
            queryParams.append(key, value.toString());
          }
        });
      }
      
      const url = `${getShopUrl()}/inventory${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
      const response = await customFetch(url);
      const data = await response.json() as ApiResponse<InventoryReport>;
      return data.data;
    },

    getLowStockVariants: async (threshold?: number): Promise<LowStockVariant[]> => {
      const params = threshold ? `?threshold=${threshold}` : '';
      const response = await customFetch(`${getShopUrl()}/inventory/low-stock${params}`);
      const data = await response.json() as ApiResponse<LowStockVariant[]>;
      return data.data;
    },

    updateVariantStock: async (variantId: number, stockData: StockUpdatePayload): Promise<InventoryItem> => {
      try {
        const response = await customFetch(`${getShopUrl()}/inventory/variants/${variantId}/stock`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(stockData),
        });
        
        const data = await response.json() as ApiResponse<InventoryItem>;
        
        if (response.ok) {
          toast.success('Stock updated successfully');
          return data.data;
        } else {
          toast.error(data.message || 'Failed to update stock');
          throw new Error(data.message || 'Failed to update stock');
        }
      } catch (error) {
        toast.error('An error occurred while updating stock');
        throw error;
      }
    },

    getStockMovements: async (params?: StockMovementSearchParams): Promise<PaginatedResponse<StockMovement>> => {
      const queryParams = new URLSearchParams();
      if (params) {
        Object.entries(params).forEach(([key, value]) => {
          if (value !== undefined && value !== null) {
            queryParams.append(key, value.toString());
          }
        });
      }
      
      const url = `${getShopUrl()}/inventory/movements${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
      const response = await customFetch(url);
      const data = await response.json() as ApiResponse<PaginatedResponse<StockMovement>>;
      return data.data;
    },

    createStockMovement: async (movementData: StockMovementCreatePayload): Promise<StockMovement> => {
      try {
        const response = await customFetch(`${getShopUrl()}/inventory/movements`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(movementData),
        });
        
        const data = await response.json() as ApiResponse<StockMovement>;
        
        if (response.ok) {
          toast.success('Stock movement recorded successfully');
          return data.data;
        } else {
          toast.error(data.message || 'Failed to record stock movement');
          throw new Error(data.message || 'Failed to record stock movement');
        }
      } catch (error) {
        toast.error('An error occurred while recording stock movement');
        throw error;
      }
    },

    // Payment Methods API
    getPaymentMethods: async (): Promise<PaymentMethod[]> => {
      try {
        const response = await authenticatedFetch(`${getShopUrl()}/payment-methods`);
        const data = await response.json() as ApiResponse<PaymentMethod[]>;
        if (response.ok) {
          return data.data || [];
        } else {
          toast.error(data.message || 'Failed to load payment methods');
          throw new Error(data.message || 'Failed to load payment methods');
        }
      } catch (error) {
        console.error('Error fetching payment methods:', error);
        toast.error('An error occurred while loading payment methods');
        return [];
      }
    },

    updatePaymentMethod: async (methodType: string, config: PaymentMethodConfig): Promise<PaymentMethod> => {
      try {
        const response = await customFetch(`${getShopUrl()}/payment-methods/${methodType}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(config),
        });
        
        const data = await response.json() as ApiResponse<PaymentMethod>;
        
        if (response.ok) {
          toast.success('Payment method updated successfully');
          return data.data;
        } else {
          toast.error(data.message || 'Failed to update payment method');
          throw new Error(data.message || 'Failed to update payment method');
        }
      } catch (error) {
        toast.error('An error occurred while updating payment method');
        throw error;
      }
    },

    updatePaymentMethodStatus: async (methodType: string, isEnabled: boolean): Promise<PaymentMethod> => {
      try {
        const response = await customFetch(`${getShopUrl()}/payment-methods/${methodType}/status`, {
          method: 'PATCH',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ is_enabled: isEnabled }),
        });
        
        const data = await response.json() as ApiResponse<PaymentMethod>;
        
        if (response.ok) {
          toast.success(`Payment method ${isEnabled ? 'enabled' : 'disabled'} successfully`);
          return data.data;
        } else {
          toast.error(data.message || 'Failed to update payment method status');
          throw new Error(data.message || 'Failed to update payment method status');
        }
      } catch (error) {
        toast.error('An error occurred while updating payment method status');
        throw error;
      }
    },

    deletePaymentMethod: async (methodType: string): Promise<void> => {
      try {
        const response = await customFetch(`${getShopUrl()}/payment-methods/${methodType}`, {
          method: 'DELETE',
        });
        
        if (response.ok) {
          toast.success('Payment method deleted successfully');
        } else {
          const data = await response.json() as ApiResponse<null>;
          toast.error(data.message || 'Failed to delete payment method');
          throw new Error(data.message || 'Failed to delete payment method');
        }
      } catch (error) {
        toast.error('An error occurred while deleting payment method');
        throw error;
      }
    },

    testPaymentMethod: async (methodType: string): Promise<any> => {
      try {
        const response = await customFetch(`${getShopUrl()}/payment-methods/${methodType}/test`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
        });
        
        const data = await response.json() as ApiResponse<any>;
        
        if (response.ok) {
          toast.success('Payment method test successful');
          return data.data;
        } else {
          toast.error(data.message || 'Payment method test failed');
          throw new Error(data.message || 'Payment method test failed');
        }
      } catch (error) {
        toast.error('An error occurred while testing payment method');
        throw error;
      }
    },

    // Publish API
    publishStore: async (templateName: string = 'template_1') => {
      try {
        const response = await authenticatedFetch(`${getShopUrl()}/publish`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            template_name: templateName,
            changes: [], // Will be set by PublishButton component
          }),
        });
        
        const data = await response.json() as ApiResponse<any>;
        
        if (response.ok) {
          toast.success('Store publish initiated successfully');
          return data.data;
        } else {
          toast.error(data.message || 'Failed to publish store');
          throw new Error(data.message || 'Failed to publish store');
        }
      } catch (error) {
        toast.error('An error occurred while publishing the store');
        throw error;
      }
    },

    getPublishStatus: async (jobId: string) => {
      try {
        const response = await authenticatedFetch(`${getShopUrl()}/publish/status?job_id=${jobId}`);
        const data = await response.json() as ApiResponse<any>;
        return data.data;
      } catch (error) {
        console.error('Error fetching publish status:', error);
        throw error;
      }
    },

    getPublishHistory: async (limit: number = 10, offset: number = 0) => {
      try {
        const response = await authenticatedFetch(`${getShopUrl()}/publish/history?limit=${limit}&offset=${offset}`);
        const data = await response.json() as ApiResponse<any[]>;
        return data.data || [];
      } catch (error) {
        console.error('Error fetching publish history:', error);
        return [];
      }
    },

    // Template Management API
    getTemplates: async () => {
      try {
        const response = await authenticatedFetch(`http://127.0.0.1:8080/v1/templates`);
        const data = await response.json() as ApiResponse<any[]>;
        
        if (response.ok) {
          return data.data || [];
        } else {
          toast.error(data.message || 'Failed to fetch templates');
          throw new Error(data.message || 'Failed to fetch templates');
        }
      } catch (error) {
        toast.error('An error occurred while fetching templates');
        throw error;
      }
    },

    getTemplateVersions: async (templateName: string) => {
      try {
        const response = await authenticatedFetch(`http://127.0.0.1:8080/v1/templates/${templateName}/versions`);
        const data = await response.json() as ApiResponse<any[]>;
        
        if (response.ok) {
          return data.data || [];
        } else {
          toast.error(data.message || 'Failed to fetch template versions');
          throw new Error(data.message || 'Failed to fetch template versions');
        }
      } catch (error) {
        toast.error('An error occurred while fetching template versions');
        throw error;
      }
    },

    getLatestTemplateVersion: async (templateName: string) => {
      try {
        const response = await authenticatedFetch(`http://127.0.0.1:8080/v1/templates/${templateName}/latest`);
        const data = await response.json() as ApiResponse<any>;
        
        if (response.ok) {
          return data.data;
        } else {
          toast.error(data.message || 'Failed to fetch latest template version');
          throw new Error(data.message || 'Failed to fetch latest template version');
        }
      } catch (error) {
        toast.error('An error occurred while fetching latest template version');
        throw error;
      }
    },

    buildTemplate: async (request: {
      template_name: string;
      git_commit?: string;
      force?: boolean;
    }) => {
      try {
        const response = await authenticatedFetch(`http://127.0.0.1:8080/v1/templates/build`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(request),
        });
        
        const data = await response.json() as ApiResponse<any>;
        
        if (response.ok) {
          toast.success('Template build initiated successfully');
          return data.data;
        } else {
          toast.error(data.message || 'Failed to initiate template build');
          throw new Error(data.message || 'Failed to initiate template build');
        }
      } catch (error) {
        toast.error('An error occurred while building template');
        throw error;
      }
    },

    deployStore: async (request: {
      template_name: string;
      version?: string;
      data_override?: Record<string, string>;
    }) => {
      try {
        const response = await authenticatedFetch(`${getShopUrl()}/deploy`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(request),
        });
        
        const data = await response.json() as ApiResponse<any>;
        
        if (response.ok) {
          toast.success('Store deployment initiated successfully');
          return data.data;
        } else {
          toast.error(data.message || 'Failed to initiate store deployment');
          throw new Error(data.message || 'Failed to initiate store deployment');
        }
      } catch (error) {
        toast.error('An error occurred while deploying store');
        throw error;
      }
    },

    updateStoreData: async (request: {
      data_type?: string;
      incremental?: boolean;
      changes?: Record<string, string>;
    }) => {
      try {
        const response = await authenticatedFetch(`${getShopUrl()}/update-data`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(request),
        });
        
        const data = await response.json() as ApiResponse<any>;
        
        if (response.ok) {
          toast.success('Store data updated successfully');
          return data.data;
        } else {
          toast.error(data.message || 'Failed to update store data');
          throw new Error(data.message || 'Failed to update store data');
        }
      } catch (error) {
        toast.error('An error occurred while updating store data');
        throw error;
      }
    },

    getDeploymentStatus: async () => {
      try {
        const response = await authenticatedFetch(`${getShopUrl()}/deployment-status`);
        const data = await response.json() as ApiResponse<any>;
        
        if (response.ok) {
          return data.data;
        } else {
          // Don't show error toast for deployment status as it might not exist yet
          return null;
        }
      } catch (error) {
        console.error('Error fetching deployment status:', error);
        return null;
      }
    }
  }
}
