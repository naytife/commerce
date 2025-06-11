import { HoudiniClient } from '$houdini';

export default new HoudiniClient({
    // VITE_API_URL="http://example.com/api/query" npm run dev
    url: import.meta.env.VITE_API_URL || 'http://ynt.localhost:8080/query'

    // uncomment this to configure the network call (for things like authentication)
    // for more information, please visit here: https://www.houdinigraphql.com/guides/authentication
    // fetchParams({ session }) {
    //     return {
    //         headers: {
    //             Authentication: `Bearer ${session.token}`,
    //         }
    //     }
    // }
})
