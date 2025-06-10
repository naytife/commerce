import { redirect, type Handle } from '@sveltejs/kit';
import { handle as authenticationHandle } from './auth';
import { sequence } from '@sveltejs/kit/hooks';
import type { RequestEvent } from '@sveltejs/kit';

async function authorizationHandle({ event, resolve }: { event: RequestEvent, resolve: any }) {
  const publicRoutes = ['/', '/login', '/signin'];
  const isPublic = publicRoutes.includes(event.url.pathname);

  if (!isPublic) {
    const session = await event.locals.auth() as any;
    // Check for missing session or refresh error
    if (!session || session.error === 'RefreshTokenError') {
      throw redirect(303, '/login');
    }
    // Check for expired access token
    if (session.access_token_expires && Date.now() > Number(session.access_token_expires) * 1000) {
      throw redirect(303, '/login');
    }
  }

  // If the request is still here, just proceed as normally
  return resolve(event);
}
 
// First handle authentication, then authorization
// Each function acts as a middleware, receiving the request handle
// And returning a handle which gets passed to the next function
export const handle: Handle = sequence(authenticationHandle, authorizationHandle)