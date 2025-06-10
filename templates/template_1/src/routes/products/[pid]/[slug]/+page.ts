import { load_ProductQuery } from '$houdini';

export const prerender = true;
export const ssr = true;
export const csr = true;

// Encode a Node global ID from type and raw ID, using btoa in the browser and Buffer on the server
function encodeGlobalId(type: string, id: string): string {
  const str = `${type}:${id}`;
  if (typeof btoa === 'function') {
    return btoa(str);
  }
  // server-side: Buffer is available on globalThis
  return globalThis.Buffer.from(str, 'utf8').toString('base64');
}

import type { PageLoad } from './$types';

export const load: PageLoad = async (event) => {
  const { pid } = event.params;
  const graphqlId = encodeGlobalId('Product', pid);
  return load_ProductQuery({ event, variables: { pid: graphqlId } });
}; 