import type {
  LoginRequest,
  LoginResponse,
  ApiError,
  PimUser,
} from "@/types";
import type {
  Product,
  Category,
  PriceScale,
  ApiProduct,
  ApiCategory,
  ApiPriceScale,
  PaginatedResponse,
  ProductSearchParams,
  BundleComponent,
  mapApiProduct,
  mapApiCategory,
  mapApiPriceScale,
} from "@/types/catalog";
import { useAuthStore } from "@/lib/stores/authStore";

const CATALOG_URL = process.env.NEXT_PUBLIC_CATALOG_URL || "http://catalog:8081";
const TENANT_ID = process.env.NEXT_PUBLIC_TENANT_ID || "demo";

class PimApiClient {
  private catalogUrl: string;
  private tenantId: string;

  constructor(catalogUrl: string, tenantId: string) {
    this.catalogUrl = catalogUrl;
    this.tenantId = tenantId;
  }

  private async request<T>(
    url: string,
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

    const response = await fetch(url, {
      ...options,
      headers,
      credentials: "include",
    });

    // Handle 401 - unauthorized
    if (response.status === 401) {
      logout();
      throw { code: "UNAUTHORIZED", message: "Session expired" } as ApiError;
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

  // ==================== AUTH API (Mock) ====================

  async login(credentials: LoginRequest): Promise<LoginResponse> {
    // Mock login - hardcoded admin/admin
    // In production, this would call a real auth service
    if (credentials.email === "admin" && credentials.password === "admin") {
      const user: PimUser = {
        id: "pim-admin-001",
        email: "admin@gondolia.com",
        firstName: "PIM",
        lastName: "Admin",
        displayName: "PIM Admin",
        role: "admin",
      };

      // Mock JWT token
      const mockToken = btoa(JSON.stringify({
        sub: user.id,
        email: user.email,
        exp: Date.now() + 8 * 60 * 60 * 1000
      }));

      return {
        accessToken: mockToken,
        expiresIn: 28800, // 8 hours
        user,
      };
    }

    throw {
      code: "INVALID_CREDENTIALS",
      message: "Invalid email or password",
    } as ApiError;
  }

  async logout(): Promise<void> {
    // Mock logout
    return Promise.resolve();
  }

  // ==================== CATALOG API (Read) ====================

  async getProducts(params?: ProductSearchParams): Promise<PaginatedResponse<Product>> {
    const searchParams = new URLSearchParams();
    if (params?.q) searchParams.set("q", params.q);
    if (params?.categoryId) searchParams.set("category_id", params.categoryId);
    if (params?.productType) searchParams.set("product_type", params.productType);
    if (params?.status) searchParams.set("status", params.status);

    const limit = params?.limit || 50;
    const offset = params?.offset !== undefined
      ? params.offset
      : params?.page ? (params.page - 1) * limit : 0;
    searchParams.set("limit", limit.toString());
    searchParams.set("offset", offset.toString());

    if (params?.sortBy) searchParams.set("sort_by", params.sortBy);
    if (params?.sortOrder) searchParams.set("sort_order", params.sortOrder);

    const query = searchParams.toString();
    const path = `${this.catalogUrl}/api/v1/products${query ? `?${query}` : ""}`;

    const response = await this.request<{
      data: ApiProduct[];
      total: number;
      limit: number;
      offset: number;
    }>(path);

    const { mapApiProduct } = await import("@/types/catalog");

    const page = Math.floor(response.offset / response.limit) + 1;
    const totalPages = Math.ceil(response.total / response.limit);

    const products = (response.data || []).map(mapApiProduct);

    // Enrich with prices
    const pricePromises = products.map(async (product) => {
      try {
        const prices = await this.getProductPrices(product.id);
        if (prices.length > 0) {
          const basePrice = prices.sort((a, b) => a.minQuantity - b.minQuantity)[0];
          product.basePrice = basePrice.price;
          product.currency = basePrice.currency;
        }
      } catch {
        // Ignore
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
    const response = await this.request<{ data: ApiProduct }>(
      `${this.catalogUrl}/api/v1/products/${id}`
    );
    const { mapApiProduct } = await import("@/types/catalog");
    return mapApiProduct(response.data);
  }

  async getProductPrices(productId: string): Promise<PriceScale[]> {
    const response = await this.request<{ data: ApiPriceScale[] }>(
      `${this.catalogUrl}/api/v1/products/${productId}/prices`
    );
    const { mapApiPriceScale } = await import("@/types/catalog");
    return (response.data || []).map(mapApiPriceScale);
  }

  async getCategories(): Promise<Category[]> {
    const response = await this.request<{ data: ApiCategory[] }>(
      `${this.catalogUrl}/api/v1/categories`
    );
    const { mapApiCategory } = await import("@/types/catalog");
    const flat = (response.data || []).map(mapApiCategory);

    // Build tree
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

  async getCategory(id: string): Promise<Category> {
    const response = await this.request<ApiCategory>(
      `${this.catalogUrl}/api/v1/categories/${id}`
    );
    const { mapApiCategory } = await import("@/types/catalog");
    return mapApiCategory(response);
  }

  async getProductCategories(productId: string): Promise<Category[]> {
    const response = await this.request<{ data: ApiCategory[] }>(
      `${this.catalogUrl}/api/v1/products/${productId}/categories`
    );
    const { mapApiCategory } = await import("@/types/catalog");
    return (response.data || []).map(mapApiCategory);
  }

  async getBundleComponents(productId: string): Promise<BundleComponent[]> {
    const response = await this.request<{
      components: Array<{
        id: string;
        tenant_id: string;
        bundle_product_id: string;
        component_product_id: string;
        product?: ApiProduct;
        quantity: number;
        min_quantity?: number;
        max_quantity?: number;
        sort_order: number;
        default_parameters?: Record<string, number> | null;
        created_at: string;
        updated_at: string;
      }>
    }>(`${this.catalogUrl}/api/v1/products/${productId}/bundle-components`);

    const { mapApiProduct } = await import("@/types/catalog");

    return response.components.map(c => ({
      id: c.id,
      tenantId: c.tenant_id,
      bundleProductId: c.bundle_product_id,
      componentProductId: c.component_product_id,
      product: c.product ? mapApiProduct(c.product) : undefined,
      quantity: c.quantity,
      minQuantity: c.min_quantity,
      maxQuantity: c.max_quantity,
      sortOrder: c.sort_order,
      defaultParameters: c.default_parameters,
      createdAt: c.created_at,
      updatedAt: c.updated_at,
    }));
  }

  // ==================== WRITE OPERATIONS (Stubs for future) ====================

  async createProduct(data: Partial<Product>): Promise<Product> {
    // TODO: Implement when catalog write API is ready
    throw { code: "NOT_IMPLEMENTED", message: "Create product not yet implemented" };
  }

  async updateProduct(id: string, data: Partial<Product>): Promise<Product> {
    // TODO: Implement when catalog write API is ready
    throw { code: "NOT_IMPLEMENTED", message: "Update product not yet implemented" };
  }

  async deleteProduct(id: string): Promise<void> {
    // TODO: Implement when catalog write API is ready
    throw { code: "NOT_IMPLEMENTED", message: "Delete product not yet implemented" };
  }

  async createPrice(productId: string, data: Omit<PriceScale, 'id' | 'productId'>): Promise<PriceScale> {
    // TODO: Implement when catalog write API is ready
    throw { code: "NOT_IMPLEMENTED", message: "Create price not yet implemented" };
  }

  async updatePrice(id: string, data: Partial<PriceScale>): Promise<PriceScale> {
    // TODO: Implement when catalog write API is ready
    throw { code: "NOT_IMPLEMENTED", message: "Update price not yet implemented" };
  }

  async deletePrice(id: string): Promise<void> {
    // TODO: Implement when catalog write API is ready
    throw { code: "NOT_IMPLEMENTED", message: "Delete price not yet implemented" };
  }

  async assignCategory(productId: string, categoryId: string): Promise<void> {
    // TODO: Implement when catalog write API is ready
    throw { code: "NOT_IMPLEMENTED", message: "Assign category not yet implemented" };
  }

  async removeCategory(productId: string, categoryId: string): Promise<void> {
    // TODO: Implement when catalog write API is ready
    throw { code: "NOT_IMPLEMENTED", message: "Remove category not yet implemented" };
  }

  async createCategory(data: Partial<Category>): Promise<Category> {
    // TODO: Implement when catalog write API is ready
    throw { code: "NOT_IMPLEMENTED", message: "Create category not yet implemented" };
  }

  async updateCategory(id: string, data: Partial<Category>): Promise<Category> {
    // TODO: Implement when catalog write API is ready
    throw { code: "NOT_IMPLEMENTED", message: "Update category not yet implemented" };
  }

  async deleteCategory(id: string): Promise<void> {
    // TODO: Implement when catalog write API is ready
    throw { code: "NOT_IMPLEMENTED", message: "Delete category not yet implemented" };
  }
}

export const pimApiClient = new PimApiClient(CATALOG_URL, TENANT_ID);
