// Type definitions for the optimized e-commerce template

export interface ProductAttribute {
	title: string;
	value: string;
}

export interface ProductVariant {
	id: string;
	variationId: number;
	description: string;
	price: number;
	availableQuantity: number;
	stockStatus: string;
	isDefault: boolean;
	attributes: ProductAttribute[];
}

export interface ProductImage {
	url: string;
	altText: string;
}

export interface Product {
	id: string;
	productId: number;
	title: string;
	description: string;
	slug: string;
	createdAt: string;
	updatedAt: string;
	defaultVariant: ProductVariant;
	variants: ProductVariant[];
	images: ProductImage[];
	attributes: ProductAttribute[];
}

// GraphQL Edge structure for maintaining compatibility with original design
export interface ProductEdge {
	cursor: string;
	node: Product;
}

export interface Category {
	id: string;
	name: string;
	description?: string;
	slug: string;
	image?: string;
}

export interface SiteMetadata {
	shopName: string;
	description: string;
	settings?: {
		currency?: string;
		locale?: string;
		theme?: string;
	};
	contact?: {
		email?: string;
		phone?: string;
	};
	social?: {
		facebook?: string;
		instagram?: string;
		twitter?: string;
	};
}

export interface CartItem {
	id: string;
	title: string;
	price: number;
	quantity: number;
	image: string;
	slug: string;
}

export interface SearchFilters {
	category?: string;
	priceRange?: [number, number];
	inStock?: boolean;
	featured?: boolean;
}
