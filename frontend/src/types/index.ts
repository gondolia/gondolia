export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  displayName: string;
  role: string;
  permissions: string[];
}

export interface Company {
  id: string;
  name: string;
  sapNumber: string;
  isActive: boolean;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  accessToken: string;
  expiresIn: number;
  user: User;
  company: Company;
}

// Raw API response types (snake_case from backend)
export interface ApiUser {
  id: string;
  email: string;
  firstname: string;
  lastname: string;
  is_active: boolean;
  companies?: ApiUserCompany[];
}

export interface ApiUserCompany {
  company: ApiCompany;
  role: ApiRole;
}

export interface ApiCompany {
  id: string;
  name: string;
  sap_company_number: string;
  is_active: boolean;
}

export interface ApiRole {
  id: string;
  name: string;
  permissions: Record<string, boolean>;
}

export interface ApiMeResponse {
  user: ApiUser;
  company: ApiCompany;
  role: ApiRole;
  permissions: string[];
}

export interface MeResponse {
  user: User;
  company: Company;
  availableCompanies: Company[];
}

export interface ApiError {
  code: string;
  message: string;
  traceId?: string;
  requestId?: string;
}

export interface ApiResponse<T> {
  data?: T;
  error?: ApiError;
}
