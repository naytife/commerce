/// <references types="houdini-svelte">

/** @type {import('houdini').ConfigFile} */
const config = {
	schemaPath: './schema.graphql',
	watchSchema: {
		url: 'https://apple.naytife.com/api/query'
	},

	scalars: {
		DateTime: {
			type: 'string'
		}
	},
	plugins: {
		'houdini-svelte': {}
	}
};

export default config;
