// src/routes/+page.ts
import { api, fetchShopIdFromSubdomain } from '$lib/api'
import type { PageLoad } from './$types'
import type { Session as AuthSession } from '@auth/sveltekit'

// Create a custom session type that includes access_token
interface CustomSession extends AuthSession {
  access_token?: string;
  user?: {
    name?: string;
    email?: string;
  };
}

export const load: PageLoad = async ({ parent, fetch, params }) => {
  // Get session from the parent layout (loaded in +layout.server.ts)
  const { session } = await parent()
  const accessToken = (session as CustomSession)?.access_token
  // console.log('accessToken page.ts', session)

  // Create a custom fetch function with the Authorization header
  const customFetch = async (input: RequestInfo | URL, init: RequestInit = {}) => {
    if (accessToken) {
      init.headers = {
        ...init.headers,
        Authorization: `Bearer ${accessToken}`,
      }
    }
    return fetch(input, init)
  }

  // Prefetch the shop ID using the shop param
  await fetchShopIdFromSubdomain(params.shop, customFetch)

  const { queryClient } = await parent()

  // Prefetch data with the authenticated fetch
  await queryClient.prefetchQuery({
    queryKey: [`shop-${params.shop}-products`, 10], // Updated
    queryFn: () => api(customFetch).getProducts(10),
  })

  return { /* Add any additional data to return */ }
}