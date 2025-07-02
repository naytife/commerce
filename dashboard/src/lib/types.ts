export type Product = {
    product_id: number
    title: string
    description: string
    product_type_id: number
    status: string
    variants: ProductVariant[]
    attributes: ProductAttribute[]
    images?: ProductImage[]
    created_at?: string
    updated_at?: string
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

  // Predefined Product Type Templates
  export interface PredefinedProductType {
    id: string;
    title: string;
    description: string;
    sku_substring: string;
    shippable: boolean;
    digital: boolean;
    category: string;
    icon: string;
    attributes: PredefinedAttributeTemplate[];
  }

  export interface PredefinedAttributeTemplate {
    title: string;
    data_type: string; // Text, Number, Date, Option, Color
    unit?: string;
    required: boolean;
    applies_to: string; // Product, ProductVariation
    options?: PredefinedAttributeOption[];
  }

  export interface PredefinedAttributeOption {
    value: string;
  }

  export interface CreateProductTypeFromTemplateParams {
    template_id: string;
  }

  export interface ProductTypeWithTemplateResponse {
    product_type: ProductType;
    attributes: Attribute[];
  }

  // Store Template Types
  export interface StoreTemplate {
    name: string;
    title: string;
    description: string;
    preview_url?: string;
    thumbnail_url?: string;
    category: string;
    features: string[];
    version: string;
    created_at: string;
    updated_at: string;
  }

  export interface StoreTemplateVersion {
    version: string;
    build_id: string;
    git_commit: string;
    status: string;
    assets_url?: string;
    created_at: string;
  }

  export interface OrderItem {
    created_at: string;
    order_id: number;
    order_item_id: number;
    price: number;
    product_variation_id: number;
    quantity: number;
    shop_id: number;
    updated_at: string;
  }

  export interface Order {
    amount: number;
    created_at: string;
    customer_email: string;
    customer_name: string;
    customer_phone: string;
    discount: number;
    items: OrderItem[];
    order_id: number;
    payment_method: string;
    payment_status: string;
    shipping_address: string;
    shipping_cost: number;
    shipping_method: string;
    shipping_status: string;
    shop_id: number;
    status: string;
    tax: number;
    transaction_id: string;
    updated_at: string;
    user_id: string;
    username: string;
  }

  // Customer Management Types
  export interface Customer {
    customer_id: number;
    user_id: string;
    first_name: string;
    last_name: string;
    email: string;
    phone?: string;
    date_of_birth?: string;
    gender?: string;
    marketing_consent: boolean;
    preferred_language?: string;
    customer_group_id?: number;
    total_spent: number;
    order_count: number;
    last_order_date?: string;
    status: 'active' | 'inactive' | 'suspended';
    created_at: string;
    updated_at: string;
  }

  export interface CustomerGroup {
    customer_group_id: number;
    name: string;
    description?: string;
    discount_percentage?: number;
    created_at: string;
    updated_at: string;
  }

  export interface CustomerAddress {
    address_id: number;
    customer_id: number;
    type: 'billing' | 'shipping';
    first_name: string;
    last_name: string;
    company?: string;
    address_line_1: string;
    address_line_2?: string;
    city: string;
    state: string;
    postal_code: string;
    country: string;
    phone?: string;
    is_default: boolean;
    created_at: string;
    updated_at: string;
  }

  export interface CustomerCreatePayload {
    first_name: string;
    last_name: string;
    email: string;
    phone?: string;
    date_of_birth?: string;
    gender?: string;
    marketing_consent: boolean;
    preferred_language?: string;
    customer_group_id?: number;
  }

  export interface CustomerUpdatePayload {
    first_name?: string;
    last_name?: string;
    email?: string;
    phone?: string;
    date_of_birth?: string;
    gender?: string;
    marketing_consent?: boolean;
    preferred_language?: string;
    customer_group_id?: number;
    status?: 'active' | 'inactive' | 'suspended';
  }

  export interface CustomerSearchParams {
    query?: string;
    status?: string;
    customer_group_id?: number;
    created_after?: string;
    created_before?: string;
    order_count_min?: number;
    order_count_max?: number;
    total_spent_min?: number;
    total_spent_max?: number;
    page?: number;
    limit?: number;
  }

  // Inventory Management Types
  export interface InventoryItem {
    variant_id: number;
    product_id: number;
    product_title: string;
    variant_title: string;
    sku: string;
    current_stock: number;
    reserved_stock: number;
    available_stock: number;
    low_stock_threshold: number;
    cost_price?: number;
    location?: string;
    last_updated: string;
  }

  export interface InventoryReport {
    total_products: number;
    total_variants: number;
    low_stock_count: number;
    out_of_stock_count: number;
    total_inventory_value: number;
    items: InventoryItem[];
  }

  export interface LowStockVariant {
    variant_id: number;
    product_id: number;
    product_title: string;
    variant_title: string;
    sku: string;
    current_stock: number;
    low_stock_threshold: number;
    days_of_stock_remaining?: number;
  }

  export interface StockMovement {
    movement_id: number;
    variant_id: number;
    product_title: string;
    variant_title: string;
    sku: string;
    movement_type: 'adjustment' | 'sale' | 'purchase' | 'return' | 'damage' | 'transfer';
    quantity_change: number;
    previous_quantity: number;
    new_quantity: number;
    reason?: string;
    reference_type?: string;
    reference_id?: number;
    user_id?: string;
    created_at: string;
  }

  export interface StockUpdatePayload {
    quantity: number;
    movement_type: 'adjustment' | 'purchase' | 'damage' | 'transfer';
    reason?: string;
    cost_price?: number;
  }

  export interface StockMovementCreatePayload {
    variant_id: number;
    movement_type: 'adjustment' | 'purchase' | 'damage' | 'transfer';
    quantity_change: number;
    reason?: string;
    reference_type?: string;
    reference_id?: number;
  }

  export interface InventorySearchParams {
    query?: string;
    low_stock_only?: boolean;
    out_of_stock_only?: boolean;
    product_type_id?: number;
    location?: string;
    page?: number;
    limit?: number;
  }

  export interface StockMovementSearchParams {
    variant_id?: number;
    movement_type?: string;
    date_from?: string;
    date_to?: string;
    user_id?: string;
    page?: number;
    limit?: number;
  }

  // Pagination Types
  export interface PaginatedResponse<T> {
    data: T[];
    pagination: {
      page: number;
      limit: number;
      total: number;
      total_pages: number;
      has_next: boolean;
      has_previous: boolean;
    };
  }

  export interface PaginationParams {
    page?: number;
    limit?: number;
  }

  // Payment Method Types
  export type PaymentMethodType = 'stripe' | 'paypal' | 'paystack' | 'flutterwave';

  export interface PaymentMethod {
    id: String;
    shop_id?: number;
    method_type: PaymentMethodType;
    enabled: boolean;
    config?: Record<string, any>;
    created_at?: string;
    updated_at?: string;
  }

  export interface PaymentMethodConfig {
    method_type: PaymentMethodType;
    is_enabled: boolean;
    config: Record<string, any>;
  }

  export interface PaymentMethodsRequest {
    payment_methods: PaymentMethodConfig[];
  }

  // Analytics Types
  export interface SalesSummary {
    total_sales: number;
    total_orders: number;
    average_order_value: number;
    sales_change_pct: number;
    orders_change_pct: number;
    customers_change_pct: number;
  }

  export interface OrdersOverTime {
    labels: string[];
    orders: number[];
  }

  export interface TopProduct {
    product_variation_id: number;
    product_name: string;
    units_sold: number;
    revenue: number;
  }

  export interface TopCustomer {
    customer_id: string | null;
    name: string;
    orders: number;
    total_spent: number;
    is_registered: boolean;
  }

  export interface CustomerSummary {
    new_customers: number;
    returning_customers: number;
    top_customers: TopCustomer[];
  }

  export interface LowStockProduct {
    product_variation_id: number;
    product_name: string;
    sku: string;
    description: string;
    stock: number;
  }
