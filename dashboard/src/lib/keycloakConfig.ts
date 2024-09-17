let keycloak: any;

if (typeof window !== 'undefined') {
	const Keycloak = (await import('keycloak-js')).default;

	const keycloakConfig = {
		url: 'http://localhost:8080/',
		realm: 'Ashgrit',
		clientId: 'Ashgrit'
	};

	keycloak = new Keycloak(keycloakConfig);
}

export { keycloak };
