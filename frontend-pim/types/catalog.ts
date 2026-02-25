// Product types
export type ProductType = 'simple' | 'variant_parent' | 'variant' | 'bundle' | 'parametric';
export type ProductStatus = 'active' | 'inactive' | 'draft';

export interface Product {
  id: string;
  tenantId: string;
  sku: string;
  name: Record<string, string>; // Localized names
  description?: Record<string, string>; // Localized descriptions
  productType: ProductType;
  status: ProductStatus;
  manufacturer?: string;
  manufacturerPartNumber?: string;
  images?: ProductImage[];
  attributes?: Record<string, any>;
  createdAt: string;
  updatedAt: string;
  // Variant-specific fields
  parentProductId?: string;
  variantAxes?: Record<string, string>; // e.g., {"color": "red", "size": "L"}
  // Pricing
  basePrice?: number;
  currency?: string;
  prices?: PriceScale[];
  // Categories
  categoryIds?: string[];
  categories?: Category[];
}

export interface ProductImage {
  url: string;
  isPrimary: boolean;
  sortOrder: number;
}

export interface PriceScale {
  id: string;
  productId: string;
  minQuantity: number;
  price: number;
  currency: string;
}

export interface Category {
  id: string;
  tenantId: string;
  code: string;
  name: Record<string, string>;
  description?: Record<string, string>;
  parentId?: string;
  sortOrder: number;
  isActive: boolean;
  productCount?: number;
  children?: Category[];
}

export interface BundleComponent {
  id: string;
  tenantId: string;
  bundleProductId: string;
  componentProductId: string;
  product?: Product;
  quantity: number;
  minQuantity?: number;
  maxQuantity?: number;
  sortOrder: number;
  defaultParameters?: Record<string, number>;
  createdAt: string;
  updatedAt: string;
}

// API response types (snake_case)
export interface ApiProduct {
  id: string;
  tenant_id: string;
  sku: string;
  name: Record<string, string>;
  description?: Record<string, string>;
  product_type: ProductType;
  status: ProductStatus;
  manufacturer?: string;
  manufacturer_part_number?: string;
  images?: Array<{
    url: string;
    is_primary: boolean;
    sort_order: number;
  }>;
  attributes?: Record<string, any>;
  created_at: string;
  updated_at: string;
  parent_product_id?: string;
  variant_axes?: Record<string, string>;
  category_ids?: string[];
}

export interface ApiCategory {
  id: string;
  tenant_id: string;
  code: string;
  name: Record<string, string>;
  description?: Record<string, string>;
  parent_id?: string;
  sort_order: number;
  is_active: boolean;
  product_count?: number;
}

export interface ApiPriceScale {
  id: string;
  product_id: string;
  min_quantity: number;
  price: number;
  currency: string;
}

// Pagination
export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

export interface ProductSearchParams {
  q?: string;
  categoryId?: string;
  productType?: ProductType;
  status?: ProductStatus;
  page?: number;
  offset?: number;
  limit?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

// Mappers
export function mapApiProduct(api: ApiProduct): Product {
  return {
    id: api.id,
    tenantId: api.tenant_id,
    sku: api.sku,
    name: api.name,
    description: api.description,
    productType: api.product_type,
    status: api.status,
    manufacturer: api.manufacturer,
    manufacturerPartNumber: api.manufacturer_part_number,
    images: api.images?.map(img => ({
      url: img.url,
      isPrimary: img.is_primary,
      sortOrder: img.sort_order,
    })),
    attributes: api.attributes,
    createdAt: api.created_at,
    updatedAt: api.updated_at,
    parentProductId: api.parent_product_id,
    variantAxes: api.variant_axes,
    categoryIds: api.category_ids || [],
  };
}

export function mapApiCategory(api: ApiCategory): Category {
  return {
    id: api.id,
    tenantId: api.tenant_id,
    code: api.code,
    name: api.name,
    description: api.description,
    parentId: api.parent_id,
    sortOrder: api.sort_order,
    isActive: api.is_active,
    productCount: api.product_count,
  };
}

export function mapApiPriceScale(api: ApiPriceScale): PriceScale {
  return {
    id: api.id,
    productId: api.product_id,
    minQuantity: api.min_quantity,
    price: api.price,
    currency: api.currency,
  };
}
