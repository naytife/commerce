import { loadStripe, type Stripe, type StripeElements, type StripePaymentElement } from '@stripe/stripe-js';
import { browser } from '$app/environment';
import { apiClient } from './api';

let stripePromise: Promise<Stripe | null> | null = null;

export const getStripe = (publishableKey: string): Promise<Stripe | null> => {
  if (!browser) return Promise.resolve(null);
  
  if (!stripePromise) {
    stripePromise = loadStripe(publishableKey);
  }
  return stripePromise;
};

export interface PaymentMethodConfig {
  id: string;
  name: string;
  provider: string;
  enabled: boolean;
  config: {
    publishable_key?: string;
    test_mode?: boolean;
  };
}

export interface CreatePaymentIntentRequest {
  amount: number;
  currency: string;
  shopId: string;
  orderId?: string;
  customerId?: string;
  metadata?: Record<string, string>;
}

export interface PaymentIntentResponse {
  client_secret: string;
  payment_intent_id: string;
}

export class StripeService {
  private stripe: Stripe | null = null;
  private elements: StripeElements | null = null;
  private paymentElement: StripePaymentElement | null = null;
  private clientSecret: string | null = null;

  async initialize(publishableKey: string, clientSecret: string): Promise<boolean> {
    try {
      this.stripe = await getStripe(publishableKey);
      if (!this.stripe) return false;

      this.clientSecret = clientSecret;
      this.elements = this.stripe.elements({
        clientSecret,
        appearance: {
          theme: 'stripe',
          variables: {
            colorPrimary: '#3b82f6',
            fontFamily: 'system-ui, sans-serif',
          },
        },
      });

      return true;
    } catch (error) {
      console.error('Failed to initialize Stripe:', error);
      return false;
    }
  }

  createPaymentElement(container: HTMLElement): void {
    if (!this.elements) {
      throw new Error('Stripe Elements not initialized');
    }

    this.paymentElement = this.elements.create('payment');
    this.paymentElement.mount(container);
  }

  async confirmPayment(returnUrl: string): Promise<{ success: boolean; error?: string }> {
    if (!this.stripe || !this.paymentElement) {
      return { success: false, error: 'Stripe not properly initialized' };
    }

    try {
      const { error } = await this.stripe.confirmPayment({
        elements: this.elements!,
        confirmParams: {
          return_url: returnUrl,
        },
      });

      if (error) {
        return { success: false, error: error.message };
      }

      return { success: true };
    } catch (error) {
      console.error('Payment confirmation failed:', error);
      return { success: false, error: 'Payment processing failed' };
    }
  }

  destroy(): void {
    if (this.paymentElement) {
      this.paymentElement.destroy();
      this.paymentElement = null;
    }
    this.elements = null;
    this.stripe = null;
    this.clientSecret = null;
  }
}

export async function fetchPaymentMethods(shopId: string): Promise<PaymentMethodConfig[]> {
  try {
    // Try REST endpoint first (if available)
    try {
      return await apiClient.restGet<PaymentMethodConfig[]>(`/shops/${shopId}/payment-methods`);
    } catch (restError) {
      console.log('REST endpoint not available, trying GraphQL...');
      
      // Fallback to GraphQL query
      const query = `
        query GetShopPaymentMethods($shopId: ID!) {
          shop(id: $shopId) {
            paymentMethods {
              id
              name
              provider
              enabled
              config
            }
          }
        }
      `;

      const result = await apiClient.query<{ shop: { paymentMethods: PaymentMethodConfig[] } }>(
        query,
        { shopId }
      );

      return result.shop?.paymentMethods || [];
    }
  } catch (error) {
    console.error('Error fetching payment methods:', error);
    
    // Development fallback - return mock Stripe configuration
    if (import.meta.env.VITE_ENVIRONMENT === 'development') {
      console.log('Using development fallback payment methods');
      return [
        {
          id: 'stripe-dev',
          name: 'Stripe',
          provider: 'stripe',
          enabled: true,
          config: {
            publishable_key: import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY || 'pk_test_fallback',
            test_mode: true,
          },
        },
      ];
    }
    
    return [];
  }
}

export async function createPaymentIntent(
  request: CreatePaymentIntentRequest
): Promise<PaymentIntentResponse | null> {
  try {
    // Try REST endpoint first (if available)
    try {
      return await apiClient.restPost<PaymentIntentResponse>(`/shops/${request.shopId}/payment/intent`, {
        amount: Math.round(request.amount * 100), // Convert to cents
        currency: request.currency,
        order_id: request.orderId,
        customer_id: request.customerId,
        metadata: {
          shop_id: request.shopId,
          order_id: request.orderId || '',
          customer_id: request.customerId || '',
          ...request.metadata,
        },
      });
    } catch (restError) {
      console.log('REST endpoint not available, trying GraphQL...');
      
      // Fallback to GraphQL mutation
      const mutation = `
        mutation CreatePaymentIntent($input: CreatePaymentIntentInput!) {
          createPaymentIntent(input: $input) {
            clientSecret
            paymentIntentId
          }
        }
      `;

      const result = await apiClient.query<{ 
        createPaymentIntent: { 
          clientSecret: string; 
          paymentIntentId: string; 
        } 
      }>(mutation, {
        input: {
          shopId: request.shopId,
          amount: Math.round(request.amount * 100),
          currency: request.currency,
          orderId: request.orderId,
          customerId: request.customerId,
          metadata: {
            shop_id: request.shopId,
            order_id: request.orderId || '',
            customer_id: request.customerId || '',
            ...request.metadata,
          },
        },
      });

      return {
        client_secret: result.createPaymentIntent.clientSecret,
        payment_intent_id: result.createPaymentIntent.paymentIntentId,
      };
    }
  } catch (error) {
    console.error('Error creating payment intent:', error);
    
    // Development fallback - create mock payment intent
    if (import.meta.env.VITE_ENVIRONMENT === 'development') {
      console.log('Using development fallback payment intent');
      return {
        client_secret: `pi_test_${Date.now()}_secret_test`,
        payment_intent_id: `pi_test_${Date.now()}`,
      };
    }
    
    return null;
  }
}
