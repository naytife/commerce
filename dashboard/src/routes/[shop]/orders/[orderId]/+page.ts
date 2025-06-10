import { api } from '$lib/api';
import type { PageLoad } from './$types';
import type { Session as AuthSession } from '@auth/sveltekit';

interface CustomSession extends AuthSession {
  access_token?: string;
}

interface OrderParams {
  orderId: string;
}

export const load: PageLoad = async ({ parent, params, fetch }) => {
  // Get session from the parent layout
  const { session } = await parent();
  const customSession = session as CustomSession | null;
  const accessToken = customSession?.access_token;
  
  // The route parameter from [orderId]
  const { orderId } = params as OrderParams;

  // Create a custom fetch function with the Authorization header
  const customFetch = async (input: RequestInfo | URL, init: RequestInit = {}) => {
    if (accessToken) {
      init.headers = {
        ...init.headers,
        Authorization: `Bearer ${accessToken}`,
      };
    }    
    return fetch(input, init);
  };

  const { queryClient } = await parent();

  // Prefetch order data
  await queryClient.prefetchQuery({
    queryKey: [`shop-${params.shop}-order`, orderId],
    queryFn: () => api(customFetch).getOrderById(Number(orderId)),
  });

  return {};
}; 