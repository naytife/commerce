// src/routes/+page.ts
import { api } from '$lib/api'
import type { PageLoad } from './$types'

export const load: PageLoad = async ({ parent, fetch, params }) => {
  // Get session from the parent layout (loaded in +layout.server.ts)
  const { session } = await parent()
  // Using type assertion as we don't have the exact type info
  const accessToken = (session as any)?.access_token
  const typeId = parseInt(params.type)
  // console.log('accessToken page.ts', session)

  // Create a custom fetch function with the Authorization header
  const customFetch = (input: RequestInfo | URL, init?: RequestInit) => {
    const options = init || {}
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
    queryKey: ['products', 'type', typeId, 10, 0],
    queryFn: () => api(customFetch).getProductsByType(typeId, 10, 0),
  })

  return { typeId }
}