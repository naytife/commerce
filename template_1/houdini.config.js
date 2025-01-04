/// <references types="houdini-svelte">

/** @type {import('houdini').ConfigFile} */
const config = {
    "watchSchema": {
        "url": "http://newstore.localhost:8080/api/query"
    },
    "runtimeDir": ".houdini",
    "plugins": {
        "houdini-svelte": {static: true, }
    },
    "scalars": {
        DateTime: {                  // <- The GraphQL Scalar
            type: "String"  // <-  The TypeScript type
          }
    }, 
}

export default config
