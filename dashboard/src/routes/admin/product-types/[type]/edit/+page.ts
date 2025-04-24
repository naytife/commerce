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
  const { session, queryClient } = await parent()
  const accessToken = (session as CustomSession)?.access_token
  const typeId = parseInt(params.type)

  const customFetch = (input: RequestInfo | URL, init?: RequestInit) => {
    const options: RequestInit = init || {}
    if (accessToken) {
      options.headers = {
        ...options.headers,
        Authorization: `Bearer ${accessToken}`,
      }
    }    
    return fetch(input, options)
  }

  // Prefetch both product type and its attributes
  await Promise.all([
    queryClient.prefetchQuery({
      queryKey: ['product-type', typeId],
      queryFn: () => api(customFetch).getProductTypeById(typeId),
    }),
    queryClient.prefetchQuery({
      queryKey: ['product-type-attributes', typeId],
      queryFn: () => api(customFetch).getProductTypeAttributes(typeId),
    })
  ])

  return { typeId }
}