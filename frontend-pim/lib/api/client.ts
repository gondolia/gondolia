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
    // Product already has category_ids — fetch product and resolve categories
    const product = await this.getProduct(productId);
    if (!product || !product.categoryIds || product.categoryIds.length === 0) {
      return [];
    }
    // Fetch all categories and filter by product's category_ids
    const allCategories = await this.getCategories();
    return allCategories.filter(c => product.categoryIds.includes(c.id));
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

  // ==================== WRITE OPERATIONS ====================

  async createProduct(data: Partial<Product>): Promise<Product> {
    const payload = {
      sku: data.sku,
      name: data.name,
      description: data.description,
      product_type: data.productType,
      status: data.status,
      manufacturer: data.manufacturer,
      manufacturer_part_number: data.manufacturerPartNumber,
      attributes: data.attributes,
      parent_product_id: data.parentProductId,
      variant_axes: data.variantAxes,
    };

    const response = await this.request<{ data: ApiProduct }>(
      `${this.catalogUrl}/api/v1/products`,
      {
        method: "POST",
        body: JSON.stringify(payload),
      }
    );

    const { mapApiProduct } = await import("@/types/catalog");
    return mapApiProduct(response.data);
  }

  async updateProduct(id: string, data: Partial<Product>): Promise<Product> {
    const payload: Record<string, any> = {};
    if (data.name) payload.name = data.name;
    if (data.description) payload.description = data.description;
    if (data.status) payload.status = data.status;
    if (data.manufacturer !== undefined) payload.manufacturer = data.manufacturer;
    if (data.manufacturerPartNumber !== undefined) payload.manufacturer_part_number = data.manufacturerPartNumber;
    if (data.attributes) payload.attributes = data.attributes;

    const response = await this.request<{ data: ApiProduct }>(
      `${this.catalogUrl}/api/v1/products/${id}`,
      {
        method: "PUT",
        body: JSON.stringify(payload),
      }
    );

    const { mapApiProduct } = await import("@/types/catalog");
    return mapApiProduct(response.data);
  }

  async updateProductStatus(id: string, status: ProductStatus): Promise<void> {
    await this.request(
      `${this.catalogUrl}/api/v1/products/${id}/status`,
      {
        method: "PATCH",
        body: JSON.stringify({ status }),
      }
    );
  }

  async deleteProduct(id: string): Promise<void> {
    await this.request(
      `${this.catalogUrl}/api/v1/products/${id}`,
      {
        method: "DELETE",
      }
    );
  }

  async createPrice(productId: string, data: Omit<PriceScale, 'id' | 'productId'>): Promise<PriceScale> {
    const payload = {
      min_quantity: data.minQuantity,
      price: data.price,
      currency: data.currency,
    };

    const response = await this.request<{ data: ApiPriceScale }>(
      `${this.catalogUrl}/api/v1/products/${productId}/prices`,
      {
        method: "POST",
        body: JSON.stringify(payload),
      }
    );

    const { mapApiPriceScale } = await import("@/types/catalog");
    return mapApiPriceScale(response.data);
  }

  async updatePrice(productId: string, priceId: string, data: Partial<PriceScale>): Promise<PriceScale> {
    const payload: Record<string, any> = {};
    if (data.minQuantity !== undefined) payload.min_quantity = data.minQuantity;
    if (data.price !== undefined) payload.price = data.price;
    if (data.currency !== undefined) payload.currency = data.currency;

    const response = await this.request<{ data: ApiPriceScale }>(
      `${this.catalogUrl}/api/v1/products/${productId}/prices/${priceId}`,
      {
        method: "PUT",
        body: JSON.stringify(payload),
      }
    );

    const { mapApiPriceScale } = await import("@/types/catalog");
    return mapApiPriceScale(response.data);
  }

  async deletePrice(productId: string, priceId: string): Promise<void> {
    await this.request(
      `${this.catalogUrl}/api/v1/products/${productId}/prices/${priceId}`,
      {
        method: "DELETE",
      }
    );
  }

  async assignCategory(productId: string, categoryId: string): Promise<void> {
    await this.request(
      `${this.catalogUrl}/api/v1/categories/${categoryId}/products`,
      {
        method: "POST",
        body: JSON.stringify({ product_id: productId }),
      }
    );
  }

  async removeCategory(productId: string, categoryId: string): Promise<void> {
    await this.request(
      `${this.catalogUrl}/api/v1/categories/${categoryId}/products/${productId}`,
      {
        method: "DELETE",
      }
    );
  }

  async createAttribute(productId: string, key: string, attrType: string, value: any): Promise<void> {
    await this.request(
      `${this.catalogUrl}/api/v1/products/${productId}/attributes`,
      {
        method: "POST",
        body: JSON.stringify({ key, type: attrType, value }),
      }
    );
  }

  async updateAttribute(productId: string, key: string, value: any): Promise<void> {
    await this.request(
      `${this.catalogUrl}/api/v1/products/${productId}/attributes/${key}`,
      {
        method: "PUT",
        body: JSON.stringify({ value }),
      }
    );
  }

  async deleteAttribute(productId: string, key: string): Promise<void> {
    await this.request(
      `${this.catalogUrl}/api/v1/products/${productId}/attributes/${key}`,
      {
        method: "DELETE",
      }
    );
  }

  async addBundleComponent(productId: string, data: {
    componentProductId: string;
    quantity: number;
    minQuantity?: number;
    maxQuantity?: number;
  }): Promise<void> {
    await this.request(
      `${this.catalogUrl}/api/v1/products/${productId}/bundle-components`,
      {
        method: "POST",
        body: JSON.stringify({
          component_product_id: data.componentProductId,
          quantity: data.quantity,
          min_quantity: data.minQuantity,
          max_quantity: data.maxQuantity,
        }),
      }
    );
  }

  async removeBundleComponent(productId: string, componentId: string): Promise<void> {
    await this.request(
      `${this.catalogUrl}/api/v1/products/${productId}/bundle-components/${componentId}`,
      {
        method: "DELETE",
      }
    );
  }

  async getVariants(parentProductId: string): Promise<Product[]> {
    const response = await this.request<{ data: ApiProduct[] }>(
      `${this.catalogUrl}/api/v1/products/${parentProductId}/variants`
    );
    const { mapApiProduct } = await import("@/types/catalog");
    return (response.data || []).map(mapApiProduct);
  }

  async createCategory(data: Partial<Category>): Promise<Category> {
    const payload = {
      code: data.code,
      name: data.name,
      description: data.description,
      parent_id: data.parentId,
      sort_order: data.sortOrder,
      is_active: data.isActive,
    };

    const response = await this.request<{ data: ApiCategory }>(
      `${this.catalogUrl}/api/v1/categories`,
      {
        method: "POST",
        body: JSON.stringify(payload),
      }
    );

    const { mapApiCategory } = await import("@/types/catalog");
    return mapApiCategory(response.data);
  }

  async updateCategory(id: string, data: Partial<Category>): Promise<Category> {
    const payload: Record<string, any> = {};
    if (data.name) payload.name = data.name;
    if (data.description !== undefined) payload.description = data.description;
    if (data.sortOrder !== undefined) payload.sort_order = data.sortOrder;
    if (data.isActive !== undefined) payload.is_active = data.isActive;

    const response = await this.request<{ data: ApiCategory }>(
      `${this.catalogUrl}/api/v1/categories/${id}`,
      {
        method: "PUT",
        body: JSON.stringify(payload),
      }
    );

    const { mapApiCategory } = await import("@/types/catalog");
    return mapApiCategory(response.data);
  }

  async deleteCategory(id: string): Promise<void> {
    await this.request(
      `${this.catalogUrl}/api/v1/categories/${id}`,
      {
        method: "DELETE",
      }
    );
  }
}

export const pimApiClient = new PimApiClient(CATALOG_URL, TENANT_ID);
