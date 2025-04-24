export type Product = {
    product_id: number
    title: string
    description: string
    product_type_id: number
    status: string
    variants: ProductVariant[]
    attributes: ProductAttribute[]
    images?: ProductImage[]
}
export type Shop = {
    shop_id: number
    title: string
    about: string
    address: string
    created_at: string
    currency_code: string
    custom_domain: string
    email: string
    facebook_link: string
    instagram_link: string
    phone_number: string
    seo_description: string
    seo_keywords: string[]
    seo_title: string
    status: string
    subdomain: string
    whatsapp_link: string
    whatsapp_phone_number: string
    updated_at?: string
    images?: {
        banner_url?: string
        banner_url_dark?: string
        cover_image_url?: string
        cover_image_url_dark?: string
        favicon_url?: string
        logo_url?: string
        logo_url_dark?: string
        shop_image_id?: number
    }
}
export type ProductType = {
    id: number
    title: string
    shippable: boolean
    digital: boolean
    sku_substring: string
  }
  export interface AttributeOption {
    attribute_option_id: number;
    value: string;
  }
  
 export interface ProductTypeAttribute {
    attribute_id: number;
    title: string;
    data_type: string;
    options?: AttributeOption[];
    required: boolean;
    applies_to: "Product" | "ProductVariation";
    product_type_id: number;
    unit?: string;
  }
  
  export interface ApiResponse<T> {
    data: T;
    message: string;
  }

  export interface AttributeOption {
		attribute_option_id: number;
		value: string;
	}

	export interface Attribute {
		attribute_id: number;
		title: string;
		data_type: string;
		required: boolean;
		applies_to: "Product" | "ProductVariation";
		product_type_id: number;
		options?: AttributeOption[];
	}

  export interface ProductAttribute {
    attribute_id: number;
    attribute_option_id?: number;
    value: string;
  }

  export interface ProductVariant {
    id?: number;
    attributes: ProductAttribute[];
    available_quantity: number;
    description: string;
    is_default: boolean;
    price: number;
    seo_description: string | null;
    seo_keywords: string[] | null;
    seo_title: string | null;
    sku?: string;
    slug?: string;
    created_at?: string | null;
    updated_at?: string | null;
  }

  export interface ProductCreatePayload {
    attributes: ProductAttribute[];
    description: string;
    title: string;
    variants: ProductVariant[];
  }

  export interface ProductImage {
    id?: number
    product_id?: number
    url: string
    alt?: string
    filename?: string
    is_primary?: boolean
    created_at?: string
    updated_at?: string
  }
