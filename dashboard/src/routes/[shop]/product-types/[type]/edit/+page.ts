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
  const { session } = await parent()
  const accessToken = (session as CustomSession)?.access_token

  const customFetch = async (input: RequestInfo | URL, init: RequestInit = {}) => {
    if (accessToken) {
      init.headers = {
        ...init.headers,
        Authorization: `Bearer ${accessToken}`,
      }
    }    
    return fetch(input, init)
  }

  const { queryClient } = await parent()

  // Prefetch both product type and its attributes
  await Promise.all([
    queryClient.prefetchQuery({
      queryKey: [`shop-${params.shop}-product-type`, params.type],
      queryFn: () => api(customFetch as any).getProductTypeById(parseInt(params.type)),
    }),
    queryClient.prefetchQuery({
      queryKey: [`shop-${params.shop}-product-type-attributes`, params.type],
      queryFn: () => api(customFetch as any).getProductTypeAttributes(parseInt(params.type)),
    })
  ])

  return {}
}