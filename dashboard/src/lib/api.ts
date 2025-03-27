import type { S } from 'vitest/dist/chunks/config.DCnyCTbs.js';
import type { Product, Shop } from './types'

export const api = (customFetch = fetch) => ({
    getProducts: async (limit: number) => {
      const response = await customFetch(
        'http://127.0.0.1:8080/api/v1/shops/2/products',
      );
      const data = (await response.json()) as Array<Product>;
      return data.data;
    },
    getProductById: async (id: number): Promise<Product> => {
      const response = await customFetch(
        `http://127.0.0.1:8080/api/v1/shops/2/products/${id}`,
      );
      const data = (await response.json()) as Product;
      return data.data;
    },
    getProductTypeAttributes: async (typeId: number) => {
      const response = await customFetch(
        `http://127.0.0.1:8080/api/v1/shops/2/product-types/${typeId}/attributes`,
      );
      const data = (await response.json()) as Array<Product>;
      console.log('data', data)
      return data.data;
    },
    getShop: async (shopId: String) => {
      const response = await customFetch(
        `http://127.0.0.1:8080/api/v1/shops/subdomain/${shopId}`,
      );
      const data = (await response.json()) as Shop;
      return data.data;
    }

  })