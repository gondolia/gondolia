import type {
  LoginRequest,
  LoginResponse,
  MeResponse,
  ApiError,
  ApiMeResponse,
  User,
  Company,
} from "@/types";
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
    const { tokens, logout } = useAuthStore.getState();

    const headers: HeadersInit = {
      "Content-Type": "application/json",
      "X-Tenant-ID": this.tenantId,
      ...options.headers,
    };

    // Add auth header if we have tokens
    if (tokens?.accessToken) {
      (headers as Record<string, string>)["Authorization"] =
        `Bearer ${tokens.accessToken}`;
    }

    const response = await fetch(`${this.baseUrl}${path}`, {
      ...options,
      headers,
      credentials: "include",
    });

    // Handle 401 - try to refresh token
    if (response.status === 401 && tokens?.refreshToken) {
      const refreshed = await this.refreshToken(tokens.refreshToken);
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

  private async refreshToken(refreshToken: string): Promise<{
    accessToken: string;
    refreshToken: string;
    expiresIn: number;
  } | null> {
    try {
      const response = await fetch(`${this.baseUrl}/api/v1/auth/refresh`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-Tenant-ID": this.tenantId,
        },
        body: JSON.stringify({ refreshToken }),
        credentials: "include",
      });

      if (!response.ok) {
        return null;
      }

      const data = await response.json();
      const { setTokens } = useAuthStore.getState();
      setTokens({
        accessToken: data.accessToken,
        refreshToken: data.refreshToken,
        expiresAt: Date.now() + data.expiresIn * 1000,
      });

      return data;
    } catch {
      return null;
    }
  }

  async login(credentials: LoginRequest): Promise<LoginResponse> {
    // 1. Call login endpoint to get tokens
    const response = await fetch(`${this.baseUrl}/api/v1/auth/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Tenant-ID": this.tenantId,
      },
      body: JSON.stringify(credentials),
      credentials: "include",
    });

    if (!response.ok) {
      const error = await this.parseError(response);
      throw error;
    }

    // API returns snake_case
    const tokenData = await response.json();
    const accessToken = tokenData.access_token;
    const refreshToken = tokenData.refresh_token;
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
      refreshToken,
      expiresIn,
      user: mapApiUser(meData.user, meData.role, meData.permissions),
      company: mapApiCompany(meData.company),
    };
  }

  async logout(): Promise<void> {
    try {
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
}

export const apiClient = new ApiClient(API_URL, TENANT_ID);
