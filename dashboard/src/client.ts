import { HoudiniClient } from '$houdini';

const currentHost = window.location.host; // gets the current host (e.g., store1.naytife.com)
const url = `https://${currentHost}/query`; // constructs the query URL dynamically

export default new HoudiniClient({
	url

	// uncomment this to configure the network call (for things like authentication)
	// for more information, please visit here: https://www.houdinigraphql.com/guides/authentication
	// fetchParams({ session }) {
	//     return {
	//         headers: {
	//             Authentication: `Bearer ${session.token}`,
	//         }
	//     }
	// }
});
