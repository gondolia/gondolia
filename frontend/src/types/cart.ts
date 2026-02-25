// Cart & Order Types für Gondolia Webshop

export type ProductType = 'simple' | 'variant' | 'bundle' | 'parametric';

export interface CartConfiguration {
  // For bundle products: component selections and quantities
  bundleComponents?: Array<{
    componentId: string;
    quantity: number;
    parameters?: Record<string, number>;
    selections?: Record<string, string>;
  }>;
  // For parametric products: parameters and selections
  parameters?: Record<string, number>;
  selections?: Record<string, string>;
}

export interface CartItem {
  id: string;
  cartId: string;
  productId: string;
  variantId?: string;
  productType: ProductType;
  productName: string;
  sku: string;
  imageUrl?: string;
  quantity: number;
  unitPrice: number;
  totalPrice: number;
  currency: string;
  configuration?: CartConfiguration;
  createdAt: string;
  updatedAt: string;
}

export interface Cart {
  id: string;
  tenantId: string;
  userId?: string;
  sessionId?: string;
  status: 'active' | 'merged' | 'completed';
  items: CartItem[];
  subtotal: number;
  currency: string;
  createdAt: string;
  updatedAt: string;
}

export interface AddToCartRequest {
  productId: string;
  variantId?: string;
  quantity: number;
  configuration?: CartConfiguration;
}

export interface UpdateCartItemRequest {
  quantity: number;
}

// Order Types

export type OrderStatus = 'pending' | 'confirmed' | 'processing' | 'shipped' | 'delivered' | 'cancelled';

export interface Address {
  street: string;
  city: string;
  postalCode: string;
  country: string;
  company?: string;
  firstName?: string;
  lastName?: string;
  phone?: string;
}

export interface OrderItem {
  id: string;
  orderId: string;
  productId: string;
  variantId?: string;
  productType: ProductType;
  productName: string;
  sku: string;
  quantity: number;
  unitPrice: number;
  totalPrice: number;
  currency: string;
  configuration?: CartConfiguration;
  createdAt: string;
}

export interface Order {
  id: string;
  tenantId: string;
  userId: string;
  orderNumber: string;
  status: OrderStatus;
  subtotal: number;
  taxAmount: number;
  total: number;
  currency: string;
  shippingAddress: Address;
  billingAddress: Address;
  notes?: string;
  itemCount?: number;
  items: OrderItem[];
  createdAt: string;
  updatedAt: string;
}

export interface CheckoutRequest {
  shippingAddress: Address;
  billingAddress: Address;
  notes?: string;
}

export interface CheckoutResponse {
  order: Order;
}

// API Response Types (snake_case from backend)

export interface ApiCartItem {
  id: string;
  cart_id: string;
  product_id: string;
  variant_id?: string;
  product_type: string;
  product_name: string;
  sku: string;
  image_url?: string;
  quantity: number;
  unit_price: number;
  total_price: number;
  currency: string;
  configuration?: Record<string, unknown>;
  created_at: string;
  updated_at: string;
}

export interface ApiCart {
  id: string;
  tenant_id: string;
  user_id?: string;
  session_id?: string;
  status: string;
  items: ApiCartItem[];
  subtotal: number;
  currency: string;
  created_at: string;
  updated_at: string;
}

export interface ApiOrderItem {
  id: string;
  order_id: string;
  product_id: string;
  variant_id?: string;
  product_type: string;
  product_name: string;
  sku: string;
  quantity: number;
  unit_price: number;
  total_price: number;
  currency: string;
  configuration?: Record<string, unknown>;
  created_at: string;
}

export interface ApiOrder {
  id: string;
  tenant_id: string;
  user_id: string;
  order_number: string;
  status: string;
  subtotal: number;
  tax_amount: number;
  total: number;
  currency: string;
  shipping_address: Record<string, unknown>;
  billing_address: Record<string, unknown>;
  notes?: string;
  item_count?: number;
  items?: ApiOrderItem[];
  created_at: string;
  updated_at: string;
}

// Mapper Functions

export function mapApiCartItem(item: ApiCartItem): CartItem {
  return {
    id: item.id,
    cartId: item.cart_id,
    productId: item.product_id,
    variantId: item.variant_id,
    productType: item.product_type as ProductType,
    productName: item.product_name,
    sku: item.sku,
    imageUrl: item.image_url,
    quantity: item.quantity,
    unitPrice: item.unit_price,
    totalPrice: item.total_price,
    currency: item.currency,
    configuration: item.configuration as CartConfiguration | undefined,
    createdAt: item.created_at,
    updatedAt: item.updated_at,
  };
}

export function mapApiCart(cart: ApiCart): Cart {
  return {
    id: cart.id,
    tenantId: cart.tenant_id,
    userId: cart.user_id,
    sessionId: cart.session_id,
    status: cart.status as Cart['status'],
    items: (cart.items || []).map(mapApiCartItem),
    subtotal: cart.subtotal,
    currency: cart.currency,
    createdAt: cart.created_at,
    updatedAt: cart.updated_at,
  };
}

export function mapApiOrderItem(item: ApiOrderItem): OrderItem {
  return {
    id: item.id,
    orderId: item.order_id,
    productId: item.product_id,
    variantId: item.variant_id,
    productType: item.product_type as ProductType,
    productName: item.product_name,
    sku: item.sku,
    quantity: item.quantity,
    unitPrice: item.unit_price,
    totalPrice: item.total_price,
    currency: item.currency,
    configuration: item.configuration as CartConfiguration | undefined,
    createdAt: item.created_at,
  };
}

export function mapApiOrder(order: ApiOrder): Order {
  return {
    id: order.id,
    tenantId: order.tenant_id,
    userId: order.user_id,
    orderNumber: order.order_number,
    status: order.status as OrderStatus,
    subtotal: order.subtotal,
    taxAmount: order.tax_amount,
    total: order.total,
    currency: order.currency,
    shippingAddress: order.shipping_address as unknown as Address,
    billingAddress: order.billing_address as unknown as Address,
    notes: order.notes,
    itemCount: order.item_count,
    items: (order.items || []).map(mapApiOrderItem),
    createdAt: order.created_at,
    updatedAt: order.updated_at,
  };
}
