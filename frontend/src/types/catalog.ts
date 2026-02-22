// Catalog Types für Gondolia Webshop

/** Ensure relative URLs have a leading slash so they resolve from root, not relative to current path */
function ensureLeadingSlash(url: string | undefined): string | undefined {
  if (!url) return url;
  if (url.startsWith('/') || url.startsWith('http')) return url;
  return `/${url}`;
}

// Product Types
export type ProductType = 'simple' | 'variant_parent' | 'variant' | 'parametric' | 'bundle';

// Variant Axis Definition
export interface VariantAxis {
  id?: string;
  attributeCode: string;
  position: number;
  label: Record<string, string>; // i18n labels
  options: AxisOption[];
  inputType?: 'select' | 'range'; // 'select' (default) or 'range' for parametric
  minValue?: number;
  maxValue?: number;
  stepValue?: number;
  unit?: string;
}

// Parametric pricing configuration
export interface ParametricPricing {
  id: string;
  productId: string;
  formulaType: string; // 'fixed' | 'per_unit' | 'per_m2' | 'per_running_meter'
  basePrice: number;
  unitPrice?: number;
  currency: string;
  minOrderValue?: number;
}

// Parametric price calculation response
export interface ParametricPriceResponse {
  sku: string;
  unitPrice: number;
  totalPrice: number;
  currency: string;
  quantity: number;
  breakdown?: Record<string, number>;
}

// Bundle types
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
  defaultParameters?: Record<string, number> | null;
  createdAt: string;
  updatedAt: string;
}

export interface BundlePriceRequest {
  components: Array<{
    component_id: string; // ID of bundle_component record, NOT product ID
    quantity: number;
    parameters?: Record<string, number>;
    selections?: Record<string, string>;
  }>;
}

export interface BundlePriceResponse {
  priceMode: string; // 'computed' | 'fixed'
  components: Array<{
    componentId: string; // bundle_component ID
    productId: string; // product ID
    sku: string;
    quantity: number;
    unitPrice: number;
    lineTotal: number;
  }>;
  total: number;
  currency: string;
}

// Axis Option (einzelner Wert einer Achse)
export interface AxisOption {
  code: string;
  label: Record<string, string>; // i18n labels
  position: number;
  available?: boolean; // dynamisch: ist diese Option bei aktueller Auswahl verfügbar?
}

// Axis Value Entry (Achsenwert einer konkreten Variante)
export interface AxisValueEntry {
  axisAttributeCode: string;
  optionCode: string;
  axisLabel?: Record<string, string>;
  optionLabel?: Record<string, string>;
}

// Variant Price
export interface VariantPrice {
  net: number;
  currency: string;
}

// Variant Availability
export interface VariantAvailability {
  inStock: boolean;
  quantity?: number;
}

// Product Variant (kompakte Darstellung innerhalb des Parent)
export interface ProductVariant {
  id: string;
  sku: string;
  axisValues: Record<string, string>; // attribute_code -> option_code
  status: string;
  images?: Array<{ url: string; isPrimary?: boolean; sortOrder?: number }>;
  price?: VariantPrice;
  availability?: VariantAvailability;
}

// Price Range (für variant_parent)
export interface PriceRange {
  min: number;
  max: number;
  currency: string;
}

export interface Product {
  id: string;
  productType?: ProductType;
  parentId?: string;
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
  
  // Variant-specific (nur bei variant_parent)
  variantAxes?: VariantAxis[];
  variants?: ProductVariant[];
  priceRange?: PriceRange;
  variantCount?: number;
  variantSummary?: Record<string, string[]>; // axis_code -> list of labels
  
  // Variant-specific (nur bei variant)
  parentSummary?: { id: string; sku: string; name: Record<string, string> };
  axisValues?: AxisValueEntry[];
  
  // Parametric-specific (nur bei parametric)
  parametricPricing?: ParametricPricing;

  // Bundle-specific (nur bei bundle)
  bundleMode?: 'fixed' | 'configurable';
  bundlePriceMode?: 'computed' | 'fixed';
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
export interface ApiVariantAxis {
  id?: string;
  attribute_code: string;
  position: number;
  label?: Record<string, string>;
  options?: ApiAxisOption[];
  input_type?: 'select' | 'range';
  min_value?: number;
  max_value?: number;
  step_value?: number;
  unit?: string;
}

export interface ApiAxisOption {
  code: string;
  label: Record<string, string>;
  position: number;
  available?: boolean;
}

export interface ApiAxisValueEntry {
  axis_attribute_code: string;
  option_code: string;
  axis_label?: Record<string, string>;
  option_label?: Record<string, string>;
}

export interface ApiProductVariant {
  id: string;
  sku: string;
  axis_values: Record<string, string>;
  status: string;
  images?: Array<{ url: string; is_primary?: boolean; sort_order?: number }>;
  price?: {
    net: number;
    currency: string;
  };
  availability?: {
    in_stock: boolean;
    quantity?: number;
  };
}

export interface ApiProduct {
  id: string;
  tenant_id: string;
  product_type?: string; // 'simple' | 'variant_parent' | 'variant'
  parent_id?: string;
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
  images?: string[] | Array<{ url: string; is_primary?: boolean; sort_order?: number }>;
  attributes?: Array<{ key: string; type: string; value: string }>; // Backend structure
  created_at: string;
  updated_at: string;
  
  // Variant-specific
  variant_axes?: ApiVariantAxis[];
  variants?: ApiProductVariant[];
  price_range?: {
    min: number;
    max: number;
    currency: string;
  };
  variant_count?: number;
  variant_summary?: Record<string, string[]>;
  parent_summary?: { id: string; sku: string; name: Record<string, string> };
  axis_values?: ApiAxisValueEntry[];

  // Bundle-specific
  bundle_mode?: 'fixed' | 'configurable';
  bundle_price_mode?: 'computed' | 'fixed';
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

// Attribute Translation (i18n for product attribute keys)
export interface ApiAttributeTranslation {
  id: string;
  tenant_id: string;
  attribute_key: string;
  locale: string;
  display_name: string;
  unit?: string;
  created_at: string;
  updated_at: string;
}

// Mapped frontend type: attribute_key -> display label (incl. unit if present)
export type AttributeLabels = Record<string, string>;

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
  productType?: string;
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

  // Find the primary image URL from the images array (prefers is_primary=true, falls back to first)
  const getPrimaryImageUrl = (
    images: string[] | Array<{ url: string; is_primary?: boolean; sort_order?: number }> | undefined
  ): string | undefined => {
    if (!images || images.length === 0) return undefined;
    if (typeof images[0] === 'string') return images[0] as string;
    const objImages = images as Array<{ url: string; is_primary?: boolean; sort_order?: number }>;
    const primary = objImages.find(img => img.is_primary === true);
    const img = primary || objImages[0];
    return img?.url;
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

  // Map variant axes — extract options from variants if backend doesn't provide them
  const variantAxes = api.variant_axes?.map(axis => {
    const backendOptions = (axis.options || []).map(opt => ({
      code: opt.code,
      label: opt.label || {},
      position: opt.position,
      available: opt.available,
    }));

    // If no options from backend, derive them from variants' axis_values
    const options = backendOptions.length > 0 ? backendOptions : (() => {
      const uniqueCodes = new Set<string>();
      (api.variants || []).forEach(v => {
        const val = v.axis_values?.[axis.attribute_code];
        if (val) uniqueCodes.add(val);
      });
      return Array.from(uniqueCodes).sort().map((code, idx) => ({
        code,
        label: { de: code.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase()) } as Record<string, string>,
        position: idx,
        available: true,
      }));
    })();

    return {
      id: axis.id,
      attributeCode: axis.attribute_code,
      position: axis.position,
      label: axis.label || { de: axis.attribute_code.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase()) },
      options,
      inputType: axis.input_type || 'select',
      minValue: axis.min_value,
      maxValue: axis.max_value,
      stepValue: axis.step_value,
      unit: axis.unit,
    };
  });

  // Map variants
  const variants = api.variants?.map(v => ({
    id: v.id,
    sku: v.sku,
    axisValues: v.axis_values,
    status: v.status,
    images: v.images?.map(img => ({
      url: ensureLeadingSlash(img.url) || '',
      isPrimary: img.is_primary,
      sortOrder: img.sort_order,
    })),
    price: v.price,
    availability: v.availability ? {
      inStock: v.availability.in_stock,
      quantity: v.availability.quantity,
    } : undefined,
  }));

  // Map axis values (for variant products)
  const axisValues = api.axis_values?.map(av => ({
    axisAttributeCode: av.axis_attribute_code,
    optionCode: av.option_code,
    axisLabel: av.axis_label,
    optionLabel: av.option_label,
  }));

  return {
    id: api.id,
    productType: (api.product_type as ProductType) || 'simple',
    parentId: api.parent_id,
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
    imageUrl: ensureLeadingSlash(api.image_url || getPrimaryImageUrl(api.images)),
    images: api.images?.map((img: string | { url: string }) => ensureLeadingSlash(typeof img === 'string' ? img : img.url)).filter((url): url is string => url !== undefined),
    attributes,
    createdAt: api.created_at,
    updatedAt: api.updated_at,
    
    // Variant-specific fields
    variantAxes,
    variants,
    priceRange: api.price_range,
    variantCount: api.variant_count,
    variantSummary: api.variant_summary,
    parentSummary: api.parent_summary,
    axisValues,

    // Bundle-specific fields
    bundleMode: api.bundle_mode,
    bundlePriceMode: api.bundle_price_mode,
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
