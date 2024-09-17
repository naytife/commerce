// // src/routes/admin/+layout.ts
// import { keycloak } from '$lib/keycloakConfig';
// import type { LayoutLoad } from './$types';

// export const load: LayoutLoad = async ({ url }) => {
// 	if (typeof window === 'undefined') {
// 		// Skip SSR and return nothing to avoid errors
// 		return {};
// 	}

// 	// Initialize Keycloak only on the client side
// 	// if (!keycloak.authenticated) {
// 	// 	await keycloak.init({ onLoad: 'login-required' });
// 	// 	return;
// 	// }

// 	// return {
// 	// 	user: keycloak.tokenParsed
// 	// };
// };
