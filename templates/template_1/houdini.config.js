/// <references types="houdini-svelte">

/** @type {import('houdini').ConfigFile} */
const config = {
    "watchSchema": {
        "url":  'http://gossip.localhost:8080/query'
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
