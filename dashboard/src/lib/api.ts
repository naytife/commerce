import type { Product, Shop, ProductType, ApiResponse, ProductTypeAttribute, ProductImage } from './types'
import { toast } from 'svelte-sonner'
import { writable } from 'svelte/store'

// Create a store to hold the shop ID
export const currentShopId = writable<number | null>(null)

// Function to fetch shop ID from subdomain
export const fetchShopIdFromSubdomain = async (customFetch = fetch) => {
  try {
    const hostname = window.location.hostname
    const subdomain = hostname.split('.')[0]
    
    // Skip for localhost development
    if (hostname === 'localhost') {
      currentShopId.set(2) // Default for development
      return 2
    }
    
    const response = await customFetch(
      `http://127.0.0.1:8000/v1/shops/subdomain/${subdomain}`,
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

export const api = (customFetch = fetch) => {
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
    
    return `http://127.0.0.1:8000/v1/shops/${shopId}`
  }

  return {
    getProducts: async (limit: number) => {
      const response = await customFetch(
        `${getShopUrl()}/products`,
      )
      const data = await response.json() as ApiResponse<Product[]>
      return data.data
    },
    getProductsByType: async (typeId: number, limit: number = 10, after: number = 0) => {
      const url = `${getShopUrl()}/product-types/${typeId}/products?after=${after}&limit=${limit}`
      const response = await customFetch(url)
      const data = await response.json() as ApiResponse<Product[]>
      return data.data
    },
    getProductById: async (id: number): Promise<Product> => {
      const response = await customFetch(
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
      const response = await customFetch(
        `${getShopUrl()}`,
      );
      const data = await response.json() as ApiResponse<Shop>;
      return data.data;
    },
    getProductTypes: async () => {
      const response = await customFetch(
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
        } else {
          toast.error('Failed to create product');
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
          return responseData.data;
        } else {
          toast.error(`Failed to update shop images: ${responseData.message || 'Unknown error'}`);
          throw new Error(responseData.message || 'Failed to update shop images');
        }
      } catch (error) {
        toast.error('An error occurred while updating shop images');
        throw error;
      }
    }
  }
}
