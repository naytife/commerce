import { api } from '$lib/api'
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
  const productId = parseInt(params.productId)
  const productTypeId = parseInt(params.type)

  // Create a custom fetch function with the Authorization header
  const customFetch = (input: RequestInfo | URL, init: RequestInit = {}) => {
    if (accessToken) {
      init.headers = {
        ...init.headers,
        Authorization: `Bearer ${accessToken}`,
      }
    }    
    return fetch(input, init)
  }

  const { queryClient } = await parent()

  // Prefetch data with the authenticated fetch
  await queryClient.prefetchQuery({
    queryKey: ['product', productId],
    queryFn: () => api(customFetch).getProductById(productId),
  })
  await queryClient.prefetchQuery({
    queryKey: ['product-type-attributes', productTypeId],
    queryFn: () => api(customFetch).getProductTypeAttributes(productTypeId),
  })

  return { productId, productTypeId }
}