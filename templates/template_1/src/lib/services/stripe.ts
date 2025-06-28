import { loadStripe, type Stripe, type StripeElements, type StripePaymentElement } from '@stripe/stripe-js';
import { browser } from '$app/environment';
import { apiClient } from './api';
import { decodeShopId } from '../utils';

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
  shopId: string;
  paymentMethodType: string;
  items: Array<{
    product_id: string;
    name: string;
    description?: string;
    amount: number;
    currency: string;
    quantity: number;
  }>;
  customer: {
    email: string;
    name: string;
    phone: string;
  };
  shipping: {
    name: string;
    address: {
      line1: string;
      line2?: string;
      city: string;
      state?: string;
      postal_code: string;
      country: string;
    };
  };
  currency: string;
  success_url: string;
  cancel_url: string;
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
    // Read from static JSON file instead of API
    const res = await fetch('/data/shop.json');
    if (!res.ok) throw new Error('Failed to load shop.json');
    const shopData = await res.json();
    if (!Array.isArray(shopData.paymentMethods)) return [];
    // Normalize config keys to match PaymentMethodConfig
    return shopData.paymentMethods.map((pm: any) => ({
      id: pm.id,
      name: pm.name,
      provider: pm.provider,
      enabled: pm.enabled,
      config: {
        publishable_key: pm.config?.publishableKey ?? undefined,
        test_mode: pm.config?.testMode ?? undefined,
        client_id: pm.config?.clientId ?? undefined,
        sandbox_mode: pm.config?.sandboxMode ?? undefined,
        public_key: pm.config?.publicKey ?? undefined,
        test_mode_paystack: pm.config?.testModePaystack ?? undefined,
        public_key_flutterwave: pm.config?.publicKeyFlutterwave ?? undefined,
        test_mode_flutterwave: pm.config?.testModeFlutterwave ?? undefined,
      }
    }));
  } catch (error) {
    console.error('Error fetching payment methods from shop.json:', error);
    return [];
  }
}

export async function createPaymentIntent(
  request: { shopId: string, orderId: string, paymentMethodType: string }
): Promise<PaymentIntentResponse | null> {
  try {
    const decodedShopId = Number(decodeShopId(request.shopId));
    const payload = {
      order_id: request.orderId,
      shop_id: decodedShopId,
      payment_method_type: request.paymentMethodType
    };
    const response = await apiClient.restPost<PaymentIntentResponse>(`/v1/payments/checkout`, payload);
    return response;
  } catch (error) {
    console.error('Error creating payment intent:', error);
    if (import.meta.env.VITE_ENVIRONMENT === 'development') {
      return {
        client_secret: `pi_test_${Date.now()}_secret_test`,
        payment_intent_id: `pi_test_${Date.now()}`,
      };
    }
    return null;
  }
}
