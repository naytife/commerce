// Order service for template_1 frontend
import { apiClient } from './api';
import { decodeShopId } from '../utils';

export interface CreateOrderRequest {
  customer_name: string;
  customer_email?: string;
  customer_phone?: string;
  shipping_address: string;
  shipping_method: string;
  payment_method: string;
  discount?: number;
  shipping_cost?: number;
  tax?: number;
  items: Array<{
    product_variation_id: string;
    quantity: number;
    price: number;
  }>;
}

export interface OrderResponse {
  order_id: string;
  status: string;
  // ...other fields as needed
}

export async function createOrder(shopId: string, order: CreateOrderRequest): Promise<OrderResponse> {
  // Use decodeShopId to get the numeric shop id for the API endpoint
  const decodedShopId = decodeShopId(shopId);
  return apiClient.restPost(`/v1/shops/${decodedShopId}/orders`, order);
}

export {};
