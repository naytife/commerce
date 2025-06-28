// Page load function for landing page with optimized data integration
import type { PageLoad } from './$types';

export const prerender = false;
export const ssr = true;
export const csr = true;

export const load: PageLoad = async ({ fetch }) => {
	try {
		// Load inventory data directly (no transformation needed)
		const inventoryResponse = await fetch('/data/products.json');
		const inventory = await inventoryResponse.json();
		
		// Load metadata
		const metadataResponse = await fetch('/data/shop.json');
		const metadata = await metadataResponse.json();

		// Return data directly from JSON structure
		return {
			ProductsQuery: {
				data: {
					products: inventory.products || { edges: [] }
				}
			},
			ShopQuery: {
				data: {
					shop: {
						currencyCode: metadata.currencyCode || 'USD',
						name: metadata.title || 'Shop'
					}
				}
			}
		};
	} catch (error) {
		console.error('Failed to load page data:', error);
		// Return empty structure to prevent crashes
		return {
			ProductsQuery: {
				data: {
					products: {
						edges: []
					}
				}
			},
			ShopQuery: {
				data: {
					shop: {
						currencyCode: 'USD',
						name: 'Shop'
					}
				}
			}
		};
	}
};
