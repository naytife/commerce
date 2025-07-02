import { fetchShopIdFromSubdomain } from '$lib/api';
import type { PageLoad } from './$types';
import type { Session as AuthSession } from '@auth/sveltekit';

interface CustomSession extends AuthSession {
  access_token?: string;
  user?: {
    name?: string;
    email?: string;
  };
}

export const load: PageLoad = async ({ parent, fetch, params }) => {
  const { session } = await parent();
  const accessToken = (session as CustomSession)?.access_token;

  // Create a custom fetch function with the Authorization header
  const customFetch = (input: RequestInfo | URL, init: RequestInit = {}) => {
    if (accessToken) {
      init.headers = {
        ...init.headers,
        Authorization: `Bearer ${accessToken}`,
      };
    }
    return fetch(input, init);
  };

  // Resolve the numeric shop ID using the authenticated fetch
  const shopId = await fetchShopIdFromSubdomain(params.shop, customFetch);

  return {
    shopId,
    accessToken,
    customFetch
  };
}; 