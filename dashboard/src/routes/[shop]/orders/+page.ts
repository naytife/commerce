import { api } from '$lib/api';
import type { PageLoad } from '../$types';
import type { Session as AuthSession } from '@auth/sveltekit';

interface CustomSession extends AuthSession {
  access_token?: string;
}

export const load: PageLoad = async ({ parent, fetch, params }: { parent: any; fetch: any; params: any }) => {
  // Get session from the parent layout
  const { session } = await parent();
  const customSession = session as CustomSession | null;
  const accessToken = customSession?.access_token;

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

  // Prefetch orders data
  await queryClient.prefetchQuery({
    queryKey: [`shop-${params.shop}-orders`],
    queryFn: () => api(customFetch).getOrders(),
  });

  return {};
}; 