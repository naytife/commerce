import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const prerender = false;
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

export const load: PageLoad = async (event) => {
  const { pid, slug } = event.params;
  
  try {
    // Load product data from our static files
    const response = await event.fetch('/data/inventory.json');
    const inventory = await response.json();
    
    // Find the product by productId in the edges structure
    const productEdge = inventory.products.edges?.find((edge: any) => 
      edge.node.productId === parseInt(pid)
    );
    
    if (!productEdge) {
      throw error(404, 'Product not found');
    }
    
    const product = productEdge.node;

    // Create a data structure for backward compatibility
    return {
      ProductQuery: {
        data: {
          product: {
            id: product.id,
            productId: product.productId,
            title: product.title,
            description: product.description,
            slug: product.slug,
            images: product.images || [],
            variants: product.variants || [],
            defaultVariant: product.defaultVariant,
            attributes: product.attributes || []
          }
        }
      }
    };
  } catch (e) {
    console.error('Error loading product:', e);
    throw error(500, 'Failed to load product data');
  }
}; 