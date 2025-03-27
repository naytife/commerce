import { redirect, type Handle } from '@sveltejs/kit';
import { handle as authenticationHandle } from './auth';
import { sequence } from '@sveltejs/kit/hooks';


async function authorizationHandle({ event, resolve }) {
  // Protect any routes under /authenticated
  if (event.url.pathname.startsWith('/admin')) {
    const session = await event.locals.auth();
    if (!session || session.error==='RefreshTokenError') {
      // Redirect to the signin page
      throw redirect(303, '/');
    }
  }
 
  // If the request is still here, just proceed as normally
  return resolve(event);
}
 
// First handle authentication, then authorization
// Each function acts as a middleware, receiving the request handle
// And returning a handle which gets passed to the next function
export const handle: Handle = sequence(authenticationHandle, authorizationHandle)