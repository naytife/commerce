import { HoudiniClient } from '$houdini';

let url = 'https://default.naytife.com/api/query'; // Fallback URL

if (typeof window !== 'undefined') {
    const host = window.location.host; // e.g., store1.naytife.com
    const subdomain = host.split('.')[0]; // extracts 'store1'

    url = `https://${subdomain}.naytife.com/api/query`;
}

export default new HoudiniClient({
    url,
    // Uncomment this to configure the network call (for things like authentication)
    // for more information, please visit here: https://www.houdinigraphql.com/guides/authentication
    // fetchParams({ session }) {
    //     return {
    //         headers: {
    //             Authentication: `Bearer ${session.token}`,
    //         }
    //     }
    // }
});
