// Product types
export type ProductType = 'simple' | 'variant_parent' | 'variant' | 'bundle' | 'parametric';
export type ProductStatus = 'active' | 'inactive' | 'draft';

// Axis value for variants (frontend format)
export interface AxisValue {
  variantId: string;
  axisId: string;
  axisAttributeCode: string;
  optionCode: string;
  axisLabel: Record<string, string>;
  optionLabel: Record<string, string>;
}

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
  axisValues?: AxisValue[]; // Detailed axis values for variants
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

// Attribute from API (array format)
export interface ApiAttribute {
  key: string;
  type: string;
  value: any;
}

// Axis value from API for variants
export interface ApiAxisValue {
  variant_id: string;
  axis_id: string;
  axis_attribute_code: string;
  option_code: string;
  axis_label: Record<string, string>;
  option_label: Record<string, string>;
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
  attributes?: ApiAttribute[]; // Array format from API
  created_at: string;
  updated_at: string;
  parent_id?: string; // Note: API uses parent_id not parent_product_id
  variant_axes?: Record<string, string>;
  axis_values?: ApiAxisValue[]; // For variant products
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
  active: boolean; // Note: API uses 'active' not 'is_active'
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
  // Convert attributes from array to Record format
  const attributesRecord: Record<string, any> = {};
  if (Array.isArray(api.attributes)) {
    for (const attr of api.attributes) {
      attributesRecord[attr.key] = attr.value;
    }
  } else if (api.attributes) {
    // Fallback if already Record format
    Object.assign(attributesRecord, api.attributes);
  }

  // Map axis_values to axisValues
  const axisValues: AxisValue[] | undefined = api.axis_values?.map(av => ({
    variantId: av.variant_id,
    axisId: av.axis_id,
    axisAttributeCode: av.axis_attribute_code,
    optionCode: av.option_code,
    axisLabel: av.axis_label,
    optionLabel: av.option_label,
  }));

  // Build variantAxes from axis_values if not provided directly
  let variantAxes = api.variant_axes;
  if (!variantAxes && api.axis_values && api.axis_values.length > 0) {
    variantAxes = {};
    for (const av of api.axis_values) {
      variantAxes[av.axis_label.de || av.axis_label.en || av.axis_attribute_code] =
        av.option_label.de || av.option_label.en || av.option_code;
    }
  }

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
    attributes: Object.keys(attributesRecord).length > 0 ? attributesRecord : undefined,
    createdAt: api.created_at,
    updatedAt: api.updated_at,
    parentProductId: api.parent_id,
    variantAxes,
    axisValues,
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
    isActive: api.active, // API uses 'active' not 'is_active'
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
