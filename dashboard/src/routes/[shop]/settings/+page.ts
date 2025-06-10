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

  // Create a custom fetch function with the Authorization header
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

  const { queryClient } = await parent()

  // Prefetch data with the authenticated fetch
  await queryClient.prefetchQuery({
    queryKey: [`shop-${params.shop}-shop`],
    queryFn: () => api(customFetch).getShop(),
  })

  return {  }
}