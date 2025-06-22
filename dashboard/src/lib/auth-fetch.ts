import { goto } from '$app/navigation';
import { browser } from '$app/environment';
import { toast } from 'svelte-sonner';

/**
 * Creates an authenticated fetch wrapper that automatically handles 401 responses
 * by redirecting to the login page
 */
export function createAuthenticatedFetch(baseFetch: any) {
  return async (input: RequestInfo | URL, init?: RequestInit): Promise<Response> => {
    const response = await baseFetch(input, init);
    
    // Handle 401 Unauthorized responses
    if (response.status === 401) {
      if (browser) {
        console.log('Received 401 response, redirecting to login');
        toast.error('Your session has expired. Please log in again.');
        goto('/login');
      }
      throw new Error('Unauthorized - session expired');
    }
    
    return response;
  };
}
