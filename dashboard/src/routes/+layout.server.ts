import type { LayoutServerLoad } from './$types';
 
export const load: LayoutServerLoad = async (event) => {
    // Use cached session if available from authorization middleware
    let session;
    if (event.locals.cachedSession) {
        session = event.locals.cachedSession;
    } else {
        session = await event.locals.auth();
        event.locals.cachedSession = session;
    }
    
    return {
        session: session
    };
};