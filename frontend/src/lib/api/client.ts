import type {
  LoginRequest,
  LoginResponse,
  MeResponse,
  ApiError,
  ApiMeResponse,
  User,
  Company,
} from "@/types";
import type {
  Product,
  Category,
  PriceScale,
  Manufacturer,
  ApiProduct,
  ApiCategory,
  ApiPriceScale,
  ApiManufacturer,
  ApiAttributeTranslation,
  AttributeLabels,
  PaginatedResponse,
  ProductSearchParams,
} from "@/types/catalog";
import { useAuthStore } from "@/lib/stores/authStore";

// Map API response (snake_case) to frontend types (camelCase)
function mapApiUser(apiUser: ApiMeResponse["user"], role: ApiMeResponse["role"], permissions: string[]): User {
  return {
    id: apiUser.id,
    email: apiUser.email,
    firstName: apiUser.firstname,
    lastName: apiUser.lastname,
    displayName: `${apiUser.firstname} ${apiUser.lastname}`,
    role: role.name,
    permissions,
  };
}

function mapApiCompany(apiCompany: ApiMeResponse["company"]): Company {
  return {
    id: apiCompany.id,
    name: apiCompany.name,
    sapNumber: apiCompany.sap_company_number,
    isActive: apiCompany.is_active,
  };
}

const API_URL = process.env.NEXT_PUBLIC_API_URL || "";
const TENANT_ID = process.env.NEXT_PUBLIC_TENANT_ID || "demo";

class ApiClient {
  private baseUrl: string;
  private tenantId: string;

  constructor(baseUrl: string, tenantId: string) {
    this.baseUrl = baseUrl;
    this.tenantId = tenantId;
  }

  private async request<T>(
    path: string,
    options: RequestInit = {}
  ): Promise<T> {
    const { accessToken, logout } = useAuthStore.getState();

    const headers: HeadersInit = {
      "Content-Type": "application/json",
      "X-Tenant-ID": this.tenantId,
      ...options.headers,
    };

    // Add auth header if we have access token
    if (accessToken) {
      (headers as Record<string, string>)["Authorization"] =
        `Bearer ${accessToken}`;
    }

    const response = await fetch(`${this.baseUrl}${path}`, {
      ...options,
      headers,
      credentials: "include", // Always include cookies
    });

    // Handle 401 - try to refresh token via cookie
    // Also attempt refresh when no accessToken (e.g. after page reload — token is in memory only)
    if (response.status === 401) {
      const refreshed = await this.refreshToken();
      if (refreshed) {
        // Retry the original request with new token
        (headers as Record<string, string>)["Authorization"] =
          `Bearer ${refreshed.accessToken}`;
        const retryResponse = await fetch(`${this.baseUrl}${path}`, {
          ...options,
          headers,
          credentials: "include",
        });
        if (!retryResponse.ok) {
          const error = await this.parseError(retryResponse);
          throw error;
        }
        return retryResponse.json();
      } else {
        logout();
        throw { code: "UNAUTHORIZED", message: "Session expired" } as ApiError;
      }
    }

    if (!response.ok) {
      const error = await this.parseError(response);
      throw error;
    }

    // Handle 204 No Content
    if (response.status === 204) {
      return {} as T;
    }

    return response.json();
  }

  private async parseError(response: Response): Promise<ApiError> {
    try {
      const data = await response.json();
      return data.error || { code: "UNKNOWN", message: "An error occurred" };
    } catch {
      return {
        code: response.status.toString(),
        message: response.statusText || "An error occurred",
      };
    }
  }

  // Guard against concurrent refresh attempts (React StrictMode double-invokes effects)
  private refreshPromise: Promise<{ accessToken: string; expiresIn: number } | null> | null = null;

  private async refreshToken(): Promise<{
    accessToken: string;
    expiresIn: number;
  } | null> {
    // If a refresh is already in flight, wait for it instead of starting another
    if (this.refreshPromise) {
      return this.refreshPromise;
    }

    this.refreshPromise = this.doRefresh();
    try {
      return await this.refreshPromise;
    } finally {
      this.refreshPromise = null;
    }
  }

  private async doRefresh(): Promise<{
    accessToken: string;
    expiresIn: number;
  } | null> {
    try {
      const response = await fetch(`${this.baseUrl}/api/v1/auth/refresh`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-Tenant-ID": this.tenantId,
        },
        credentials: "include", // Cookie is sent automatically
      });

      if (!response.ok) {
        return null;
      }

      const data = await response.json();
      const { setAccessToken } = useAuthStore.getState();
      setAccessToken(data.access_token, data.expires_in);

      return {
        accessToken: data.access_token,
        expiresIn: data.expires_in,
      };
    } catch {
      return null;
    }
  }

  async login(credentials: LoginRequest): Promise<LoginResponse> {
    // 1. Call login endpoint to get access token (refresh token is set as HttpOnly cookie)
    const response = await fetch(`${this.baseUrl}/api/v1/auth/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Tenant-ID": this.tenantId,
      },
      body: JSON.stringify(credentials),
      credentials: "include", // Receive refresh token cookie
    });

    if (!response.ok) {
      const error = await this.parseError(response);
      throw error;
    }

    // API returns snake_case (only access token in body now)
    const tokenData = await response.json();
    const accessToken = tokenData.access_token;
    const expiresIn = tokenData.expires_in;

    // 2. Fetch user info with the new token
    const meResponse = await fetch(`${this.baseUrl}/api/v1/auth/me`, {
      headers: {
        "Content-Type": "application/json",
        "X-Tenant-ID": this.tenantId,
        "Authorization": `Bearer ${accessToken}`,
      },
      credentials: "include",
    });

    if (!meResponse.ok) {
      const error = await this.parseError(meResponse);
      throw error;
    }

    const meData: ApiMeResponse = await meResponse.json();

    return {
      accessToken,
      expiresIn,
      user: mapApiUser(meData.user, meData.role, meData.permissions),
      company: mapApiCompany(meData.company),
    };
  }

  async logout(): Promise<void> {
    try {
      // Server will clear the refresh token cookie
      await this.request<void>("/api/v1/auth/logout", {
        method: "POST",
      });
    } catch {
      // Ignore errors on logout
    }
  }

  async getMe(): Promise<MeResponse> {
    const apiResponse = await this.request<ApiMeResponse>("/api/v1/auth/me");

    // Map available companies from user's company memberships
    const availableCompanies: Company[] = apiResponse.user.companies
      ? apiResponse.user.companies.map((uc) => mapApiCompany(uc.company))
      : [mapApiCompany(apiResponse.company)];

    return {
      user: mapApiUser(apiResponse.user, apiResponse.role, apiResponse.permissions),
      company: mapApiCompany(apiResponse.company),
      availableCompanies,
    };
  }

  async healthCheck(): Promise<{ status: string }> {
    const response = await fetch(`${this.baseUrl}/health/ready`, {
      headers: {
        "X-Tenant-ID": this.tenantId,
      },
    });
    if (!response.ok) {
      throw new Error("Health check failed");
    }
    return response.json();
  }

  // ==================== CATALOG API ====================

  // Products
  async getProducts(params?: ProductSearchParams): Promise<PaginatedResponse<Product>> {
    const searchParams = new URLSearchParams();
    if (params?.q) searchParams.set("q", params.q);
    if (params?.categoryId) searchParams.set("category_id", params.categoryId);
    if (params?.manufacturerId) searchParams.set("manufacturer_id", params.manufacturerId);
    if (params?.minPrice) searchParams.set("min_price", params.minPrice.toString());
    if (params?.maxPrice) searchParams.set("max_price", params.maxPrice.toString());
    
    // Backend uses offset/limit, not page
    const limit = params?.limit || 50;
    const offset = params?.page ? (params.page - 1) * limit : 0;
    searchParams.set("limit", limit.toString());
    searchParams.set("offset", offset.toString());
    
    if (params?.sortBy) searchParams.set("sort_by", params.sortBy);
    if (params?.sortOrder) searchParams.set("sort_order", params.sortOrder);

    const query = searchParams.toString();
    const path = `/api/v1/products${query ? `?${query}` : ""}`;
    
    // Backend returns {data: [], limit, offset, total}
    const response = await this.request<{
      data: ApiProduct[];
      total: number;
      limit: number;
      offset: number;
    }>(path);

    // Import mappers dynamically
    const { mapApiProduct } = await import("@/types/catalog");

    // Convert offset-based to page-based pagination
    const page = Math.floor(response.offset / response.limit) + 1;
    const totalPages = Math.ceil(response.total / response.limit);

    const products = (response.data || []).map(mapApiProduct);

    // Enrich products with base price from prices endpoint
    const pricePromises = products.map(async (product) => {
      try {
        const prices = await this.getProductPrices(product.id);
        if (prices.length > 0) {
          // Use lowest min_quantity price as base price
          const basePrice = prices
            .sort((a, b) => a.minQuantity - b.minQuantity)[0];
          product.basePrice = basePrice.price;
          product.currency = basePrice.currency;
        }
      } catch {
        // Ignore price fetch errors
      }
      return product;
    });

    const enrichedProducts = await Promise.all(pricePromises);

    return {
      items: enrichedProducts,
      total: response.total,
      page,
      limit: response.limit,
      totalPages,
    };
  }

  async getProduct(id: string): Promise<Product> {
    const response = await this.request<{ data: ApiProduct; parametric_pricing?: any }>(`/api/v1/products/${id}`);
    const { mapApiProduct } = await import("@/types/catalog");
    const product = mapApiProduct(response.data);
    // Attach parametric pricing if present
    if (response.parametric_pricing) {
      const pp = response.parametric_pricing;
      product.parametricPricing = {
        id: pp.id,
        productId: pp.product_id,
        formulaType: pp.formula_type,
        basePrice: pp.base_price,
        unitPrice: pp.unit_price,
        currency: pp.currency,
        minOrderValue: pp.min_order_value,
      };
    }
    return product;
  }

  async calculateParametricPrice(productId: string, parameters: Record<string, number>, selections?: Record<string, string>, quantity?: number): Promise<import("@/types/catalog").ParametricPriceResponse> {
    const raw = await this.request<Record<string, unknown>>(`/api/v1/products/${productId}/calculate-price`, {
      method: 'POST',
      body: JSON.stringify({ parameters, selections: selections || {}, quantity: quantity || 1 }),
    });
    // Map snake_case API response to camelCase
    return {
      sku: raw.sku as string,
      unitPrice: raw.unit_price as number,
      totalPrice: raw.total_price as number,
      currency: raw.currency as string,
      quantity: raw.quantity as number,
      breakdown: raw.breakdown as Record<string, number> | undefined,
    };
  }

  async searchProducts(query: string): Promise<Product[]> {
    // Backend doesn't have /search endpoint, use getProducts with q param
    const response = await this.getProducts({ q: query, limit: 50 });
    return response.items;
  }

  // Categories
  async getCategories(): Promise<Category[]> {
    // Backend returns {data: []} as a flat list
    const response = await this.request<{ data: ApiCategory[] }>("/api/v1/categories");
    const { mapApiCategory } = await import("@/types/catalog");
    const flat = (response.data || []).map(mapApiCategory);

    // Build tree: assign children to their parents
    const byId = new Map<string, Category>();
    for (const cat of flat) {
      cat.children = [];
      byId.set(cat.id, cat);
    }
    const roots: Category[] = [];
    for (const cat of flat) {
      if (cat.parentId && byId.has(cat.parentId)) {
        byId.get(cat.parentId)!.children!.push(cat);
      } else if (!cat.parentId) {
        roots.push(cat);
      }
    }
    return roots;
  }

  async getChildCategories(parentId: string): Promise<Category[]> {
    // Backend /categories/list endpoint with parent_id filter and product_count
    const response = await this.request<{ data: ApiCategory[]; total: number }>(
      `/api/v1/categories/list?parent_id=${parentId}`
    );
    const { mapApiCategory } = await import("@/types/catalog");
    return (response.data || []).map(mapApiCategory);
  }

  async getCategory(id: string): Promise<Category> {
    const response = await this.request<ApiCategory>(`/api/v1/categories/${id}`);
    const { mapApiCategory } = await import("@/types/catalog");
    return mapApiCategory(response);
  }

  async getCategoryProducts(
    categoryId: string,
    params?: { page?: number; limit?: number; includeChildren?: boolean }
  ): Promise<PaginatedResponse<Product>> {
    const searchParams = new URLSearchParams();
    
    // Backend uses offset/limit, not page
    const limit = params?.limit || 50;
    const offset = params?.page ? (params.page - 1) * limit : 0;
    searchParams.set("limit", limit.toString());
    searchParams.set("offset", offset.toString());
    
    // Add include_children parameter if specified
    if (params?.includeChildren) {
      searchParams.set("include_children", "true");
    }

    const query = searchParams.toString();
    const path = `/api/v1/categories/${categoryId}/products${query ? `?${query}` : ""}`;
    
    // Backend returns {data: [], limit, offset, total}
    const response = await this.request<{
      data: ApiProduct[];
      total: number;
      limit: number;
      offset: number;
    }>(path);

    // Import mappers dynamically
    const { mapApiProduct } = await import("@/types/catalog");

    // Convert offset-based to page-based pagination
    const page = Math.floor(response.offset / response.limit) + 1;
    const totalPages = Math.ceil(response.total / response.limit);

    const products = (response.data || []).map(mapApiProduct);

    // Enrich products with base price from prices endpoint
    const pricePromises = products.map(async (product) => {
      try {
        const prices = await this.getProductPrices(product.id);
        if (prices.length > 0) {
          // Use lowest min_quantity price as base price
          const basePrice = prices
            .sort((a, b) => a.minQuantity - b.minQuantity)[0];
          product.basePrice = basePrice.price;
          product.currency = basePrice.currency;
        }
      } catch {
        // Ignore price fetch errors
      }
      return product;
    });

    const enrichedProducts = await Promise.all(pricePromises);

    return {
      items: enrichedProducts,
      total: response.total,
      page,
      limit: response.limit,
      totalPages,
    };
  }

  // Prices
  async getProductPrices(productId: string): Promise<PriceScale[]> {
    // Backend returns {data: []}
    const response = await this.request<{ data: ApiPriceScale[] }>(
      `/api/v1/products/${productId}/prices`
    );
    const { mapApiPriceScale } = await import("@/types/catalog");
    return (response.data || []).map(mapApiPriceScale);
  }

  // Manufacturers (optional, if backend supports)
  async getManufacturers(): Promise<Manufacturer[]> {
    try {
      // Backend might return {data: []} format
      const response = await this.request<{ data?: ApiManufacturer[] } | ApiManufacturer[]>("/api/v1/manufacturers");
      const { mapApiManufacturer } = await import("@/types/catalog");
      
      // Handle both {data: []} and [] responses
      const items = Array.isArray(response) ? response : (response.data || []);
      return items.map(mapApiManufacturer);
    } catch {
      // Fallback if endpoint doesn't exist yet
      return [];
    }
  }

  // Attribute Translations
  async getAttributeTranslations(locale = "de"): Promise<AttributeLabels> {
    try {
      const response = await this.request<{
        data: Record<string, ApiAttributeTranslation>;
      }>(`/api/v1/attribute-translations/by-locale/${locale}`);

      // Convert to simple key -> label map, appending unit if present
      const labels: AttributeLabels = {};
      for (const [key, translation] of Object.entries(response.data || {})) {
        labels[key] = translation.unit
          ? `${translation.display_name} (${translation.unit})`
          : translation.display_name;
      }
      return labels;
    } catch {
      // Fallback: return empty map, caller will use raw keys
      return {};
    }
  }

  // ==================== VARIANT API ====================

  // Get all variants of a parent product
  async getProductVariants(productId: string): Promise<import("@/types/catalog").ProductVariant[]> {
    const response = await this.request<{ data: import("@/types/catalog").ApiProductVariant[] }>(
      `/api/v1/products/${productId}/variants`
    );
    
    return (response.data || []).map(v => ({
      id: v.id,
      sku: v.sku,
      axisValues: v.axis_values,
      status: v.status,
      images: v.images?.map(img => ({
        url: img.url,
        isPrimary: img.is_primary,
        sortOrder: img.sort_order,
      })),
      price: v.price,
      availability: v.availability ? {
        inStock: v.availability.in_stock,
        quantity: v.availability.quantity,
      } : undefined,
    }));
  }

  // Select a specific variant based on axis values
  async selectVariant(
    productId: string, 
    axisValues: Record<string, string>
  ): Promise<Product> {
    const params = new URLSearchParams(axisValues);
    const response = await this.request<{ data: import("@/types/catalog").ApiProduct }>(
      `/api/v1/products/${productId}/variants/select?${params.toString()}`
    );
    
    const { mapApiProduct } = await import("@/types/catalog");
    return mapApiProduct(response.data);
  }

  // Get available axis values based on current selection
  async getAvailableAxisValues(
    productId: string,
    selectedAxes: Record<string, string>
  ): Promise<{
    selected: Record<string, string>;
    available: Record<string, import("@/types/catalog").AxisOption[]>;
  }> {
    const params = new URLSearchParams(selectedAxes);
    const response = await this.request<{
      selected: Record<string, string>;
      available: Record<string, import("@/types/catalog").ApiAxisOption[]>;
    }>(`/api/v1/products/${productId}/variants/available?${params.toString()}`);
    
    // Map available options
    const available: Record<string, import("@/types/catalog").AxisOption[]> = {};
    for (const [axisCode, options] of Object.entries(response.available)) {
      available[axisCode] = options.map(opt => ({
        code: opt.code,
        label: opt.label || {},
        position: opt.position,
        available: opt.available,
      }));
    }
    
    return {
      selected: response.selected,
      available,
    };
  }
}

export const apiClient = new ApiClient(API_URL, TENANT_ID);
