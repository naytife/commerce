import { api } from '$lib/api'
import type { PageLoad } from './$types'


export const load: PageLoad = async ({ parent, fetch, params }) => {
  // Get session from the parent layout (loaded in +layout.server.ts)
  const { session } = await parent()
  const accessToken = session?.access_token

  // Create a custom fetch function with the Authorization header
  const customFetch = async (url: string, options: RequestInit = {}) => {
    if (accessToken) {
      options.headers = {
        ...options.headers,
        Authorization: `Bearer ${accessToken}`,
      }
    }    
    return fetch(url, options)
  }

  const { queryClient } = await parent()

  // Prefetch data with the authenticated fetch
  await queryClient.prefetchQuery({
    queryKey: ['shop', 'gossip'], // Adjusted to match api.ts
    queryFn: () => api(customFetch).getShop("gossip"),
  })

  return {  }
}