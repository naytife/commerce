import { api } from '$lib/api'
import type { PageLoad } from './$types'
import type { Session as AuthSession } from '@auth/sveltekit';

interface CustomSession extends AuthSession {
  access_token?: string;
}


export const load: PageLoad = async ({ parent, fetch, params }) => {
  // Get session from the parent layout (loaded in +layout.server.ts)
  const { session } = await parent()
  const accessToken = (session as CustomSession)?.access_token
  const typeId = parseInt(params.type)

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

  const { queryClient } = await parent()

  // Prefetch data with the authenticated fetch
  await queryClient.prefetchQuery({
    queryKey: [`shop-${params.shop}-product-type-attributes`, params.type],
    queryFn: () => api(customFetch as any).getProductTypeAttributes(parseInt(params.type)),
  })

  return { typeId }
}