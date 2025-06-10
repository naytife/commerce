import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { env } from '$env/dynamic/private';
import Stripe from 'stripe';

// This would normally be imported from the backend or a shared package
const stripe = new Stripe(env.STRIPE_SECRET_KEY || '', {
  apiVersion: '2025-05-28.basil',
});

export const POST: RequestHandler = async ({ request }) => {
  const body = await request.text();
  const signature = request.headers.get('stripe-signature');

  if (!signature) {
    return json({ error: 'Missing Stripe signature' }, { status: 400 });
  }

  const webhookSecret = env.STRIPE_WEBHOOK_SECRET;
  if (!webhookSecret) {
    console.error('STRIPE_WEBHOOK_SECRET is not configured');
    return json({ error: 'Webhook secret not configured' }, { status: 500 });
  }

  let event: Stripe.Event;

  try {
    event = stripe.webhooks.constructEvent(body, signature, webhookSecret);
  } catch (err) {
    console.error('Webhook signature verification failed:', err);
    return json({ error: 'Invalid signature' }, { status: 400 });
  }

  console.log(`Received webhook: ${event.type}`);

  try {
    switch (event.type) {
      case 'payment_intent.succeeded':
        await handlePaymentSucceeded(event.data.object as Stripe.PaymentIntent);
        break;
        
      case 'payment_intent.payment_failed':
        await handlePaymentFailed(event.data.object as Stripe.PaymentIntent);
        break;
        
      case 'payment_intent.canceled':
        await handlePaymentCanceled(event.data.object as Stripe.PaymentIntent);
        break;
        
      case 'charge.dispute.created':
        await handleChargeDispute(event.data.object as Stripe.Dispute);
        break;
        
      default:
        console.log(`Unhandled event type: ${event.type}`);
    }

    return json({ received: true });
  } catch (error) {
    console.error('Error processing webhook:', error);
    return json({ error: 'Webhook processing failed' }, { status: 500 });
  }
};

async function handlePaymentSucceeded(paymentIntent: Stripe.PaymentIntent) {
  console.log('Payment succeeded:', paymentIntent.id);
  
  // Extract order information from metadata
  const orderId = paymentIntent.metadata.order_id;
  const shopId = paymentIntent.metadata.shop_id;
  
  if (!orderId || !shopId) {
    console.error('Missing order_id or shop_id in payment intent metadata');
    return;
  }

  try {
    // Update order status in the backend
    const response = await fetch(`${env.BACKEND_API_URL}/orders/${orderId}/status`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${env.BACKEND_API_TOKEN}`, // If using API tokens
      },
      body: JSON.stringify({
        payment_status: 'paid',
        status: 'processing',
        transaction_id: paymentIntent.id,
      }),
    });

    if (!response.ok) {
      throw new Error(`Failed to update order status: ${response.statusText}`);
    }

    console.log(`Order ${orderId} payment confirmed and status updated`);
    
    // Here you could also:
    // - Send confirmation email to customer
    // - Update inventory
    // - Trigger fulfillment process
    // - Update analytics

  } catch (error) {
    console.error('Error updating order after payment success:', error);
  }
}

async function handlePaymentFailed(paymentIntent: Stripe.PaymentIntent) {
  console.log('Payment failed:', paymentIntent.id);
  
  const orderId = paymentIntent.metadata.order_id;
  const shopId = paymentIntent.metadata.shop_id;
  
  if (!orderId || !shopId) {
    console.error('Missing order_id or shop_id in payment intent metadata');
    return;
  }

  try {
    // Update order status to reflect payment failure
    const response = await fetch(`${env.BACKEND_API_URL}/orders/${orderId}/status`, {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${env.BACKEND_API_TOKEN}`,
      },
      body: JSON.stringify({
        payment_status: 'failed',
        status: 'cancelled',
        transaction_id: paymentIntent.id,
      }),
    });

    if (!response.ok) {
      throw new Error(`Failed to update order status: ${response.statusText}`);
    }

    console.log(`Order ${orderId} marked as payment failed`);

  } catch (error) {
    console.error('Error updating order after payment failure:', error);
  }
}

async function handlePaymentCanceled(paymentIntent: Stripe.PaymentIntent) {
  console.log('Payment canceled:', paymentIntent.id);
  
  const orderId = paymentIntent.metadata.order_id;
  
  if (orderId) {
    try {
      const response = await fetch(`${env.BACKEND_API_URL}/orders/${orderId}/status`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${env.BACKEND_API_TOKEN}`,
        },
        body: JSON.stringify({
          payment_status: 'pending',
          status: 'cancelled',
        }),
      });

      if (!response.ok) {
        throw new Error(`Failed to update order status: ${response.statusText}`);
      }

      console.log(`Order ${orderId} marked as cancelled`);

    } catch (error) {
      console.error('Error updating order after payment cancellation:', error);
    }
  }
}

async function handleChargeDispute(dispute: Stripe.Dispute) {
  console.log('Charge dispute created:', dispute.id);
  
  // Handle dispute notifications
  // This could involve:
  // - Notifying administrators
  // - Updating order status
  // - Logging for dispute management
  
  console.log(`Dispute for charge ${dispute.charge} - Reason: ${dispute.reason}`);
}
