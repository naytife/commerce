import { api } from '$lib/api'
import type { PageLoad } from './$types'
import type { Session as AuthSession } from '@auth/sveltekit'

interface CustomSession extends AuthSession {
  access_token?: string;
}

export const load: PageLoad = async ({ parent, fetch, params }) => {
  // Get session from the parent layout
  const { session } = await parent()
  const accessToken = (session as CustomSession)?.access_token

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

  // Prefetch templates and deployment status
  await Promise.all([
    queryClient.prefetchQuery({
      queryKey: ['templates'],
      queryFn: () => api(customFetch).getTemplates(),
    }),
    queryClient.prefetchQuery({
      queryKey: ['deployment-status'],
      queryFn: () => api(customFetch).getDeploymentStatus(),
    })
  ])

  return {}
}
