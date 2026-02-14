// Catalog Types für Gondolia Webshop

export interface Product {
  id: string;
  sku: string;
  name: string;
  description: string;
  shortDescription?: string;
  categoryId: string;
  categoryName?: string;
  manufacturerId?: string;
  manufacturerName?: string;
  basePrice: number;
  currency: string;
  unit: string;
  minOrderQuantity: number;
  stockQuantity: number;
  isActive: boolean;
  imageUrl?: string;
  images?: string[];
  attributes?: Record<string, string>;
  createdAt: string;
  updatedAt: string;
}

export interface Category {
  id: string;
  name: string;
  description?: string;
  parentId?: string;
  level: number;
  path: string;
  isActive: boolean;
  productCount?: number;
  imageUrl?: string;
  children?: Category[];
  createdAt: string;
  updatedAt: string;
}

export interface PriceScale {
  id: string;
  productId: string;
  minQuantity: number;
  maxQuantity?: number;
  price: number;
  currency: string;
  discountPercent?: number;
  createdAt: string;
  updatedAt: string;
}

export interface Manufacturer {
  id: string;
  name: string;
  description?: string;
  logoUrl?: string;
  website?: string;
  isActive: boolean;
}

// API Response Types (snake_case from backend)
export interface ApiProduct {
  id: string;
  tenant_id: string;
  sku: string;
  name: string | { de?: string; en?: string }; // Multilingual support
  description: string | { de?: string; en?: string }; // Multilingual support
  short_description?: string | { de?: string; en?: string };
  category_ids?: string[]; // Backend uses array
  category_id?: string; // Computed for compat
  category_name?: string;
  manufacturer_id?: string;
  manufacturer_name?: string;
  base_price?: number; // Optional in backend
  currency?: string;
  unit?: string;
  min_order_quantity?: number;
  stock_quantity?: number;
  status?: string; // Backend uses "active" | "inactive" | "discontinued"
  is_active?: boolean;
  image_url?: string;
  images?: string[];
  attributes?: Array<{ key: string; type: string; value: string }>; // Backend structure
  created_at: string;
  updated_at: string;
}

export interface ApiCategory {
  id: string;
  tenant_id: string;
  code: string;
  name: string | { de?: string; en?: string }; // Multilingual support
  description?: string | { de?: string; en?: string };
  parent_id?: string;
  sort_order: number;
  active: boolean;
  level?: number; // Computed
  path?: string; // Computed
  is_active?: boolean; // Compat
  product_count?: number;
  image_url?: string;
  children?: ApiCategory[];
  created_at: string;
  updated_at: string;
}

export interface ApiPriceScale {
  id: string;
  product_id: string;
  min_quantity: number;
  max_quantity?: number;
  price: number;
  currency: string;
  discount_percent?: number;
  created_at: string;
  updated_at: string;
}

export interface ApiManufacturer {
  id: string;
  name: string;
  description?: string;
  logo_url?: string;
  website?: string;
  is_active: boolean;
}

// Paginated Response
export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

// Search & Filter Params
export interface ProductSearchParams {
  q?: string;
  categoryId?: string;
  manufacturerId?: string;
  minPrice?: number;
  maxPrice?: number;
  page?: number;
  limit?: number;
  sortBy?: "name" | "price" | "created_at";
  sortOrder?: "asc" | "desc";
}

// Mappers
export function mapApiProduct(api: ApiProduct): Product {
  // Extract language-specific values (prefer 'de', fallback to 'en' or string)
  const getName = (val: string | { de?: string; en?: string } | undefined): string => {
    if (!val) return '';
    if (typeof val === 'string') return val;
    return val.de || val.en || '';
  };

  // Convert attributes array to record
  const attributes: Record<string, string> = {};
  if (Array.isArray(api.attributes)) {
    api.attributes.forEach(attr => {
      attributes[attr.key] = attr.value;
    });
  } else if (api.attributes && typeof api.attributes === 'object') {
    Object.assign(attributes, api.attributes);
  }

  return {
    id: api.id,
    sku: api.sku,
    name: getName(api.name),
    description: getName(api.description),
    shortDescription: api.short_description ? getName(api.short_description) : undefined,
    categoryId: api.category_id || (api.category_ids && api.category_ids[0]) || '',
    categoryName: api.category_name,
    manufacturerId: api.manufacturer_id,
    manufacturerName: api.manufacturer_name,
    basePrice: api.base_price || 0,
    currency: api.currency || 'CHF',
    unit: api.unit || 'Stk',
    minOrderQuantity: api.min_order_quantity || 1,
    stockQuantity: api.stock_quantity || 0,
    isActive: api.is_active ?? (api.status === 'active'),
    imageUrl: api.image_url || (api.images && api.images.length > 0 ? (typeof api.images[0] === 'string' ? api.images[0] : api.images[0]?.url) : undefined),
    images: api.images?.map((img: string | { url: string }) => typeof img === 'string' ? img : img.url),
    attributes,
    createdAt: api.created_at,
    updatedAt: api.updated_at,
  };
}

export function mapApiCategory(api: ApiCategory): Category {
  // Extract language-specific values (prefer 'de', fallback to 'en' or string)
  const getName = (val: string | { de?: string; en?: string } | undefined): string => {
    if (!val) return '';
    if (typeof val === 'string') return val;
    return val.de || val.en || '';
  };

  // Compute level from parent_id (0 if root, otherwise needs to be calculated)
  const level = api.level ?? (api.parent_id ? 1 : 0);
  
  // Compute path from code or id
  const path = api.path ?? `/${api.code || api.id}`;

  return {
    id: api.id,
    name: getName(api.name),
    description: api.description ? getName(api.description) : undefined,
    parentId: api.parent_id,
    level,
    path,
    isActive: api.is_active ?? api.active,
    productCount: api.product_count,
    imageUrl: api.image_url,
    children: api.children?.map(mapApiCategory),
    createdAt: api.created_at,
    updatedAt: api.updated_at,
  };
}

export function mapApiPriceScale(api: ApiPriceScale): PriceScale {
  return {
    id: api.id,
    productId: api.product_id,
    minQuantity: api.min_quantity,
    maxQuantity: api.max_quantity,
    price: api.price,
    currency: api.currency,
    discountPercent: api.discount_percent,
    createdAt: api.created_at,
    updatedAt: api.updated_at,
  };
}

export function mapApiManufacturer(api: ApiManufacturer): Manufacturer {
  return {
    id: api.id,
    name: api.name,
    description: api.description,
    logoUrl: api.logo_url,
    website: api.website,
    isActive: api.is_active,
  };
}
