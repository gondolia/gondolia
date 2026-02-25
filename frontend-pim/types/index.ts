// Auth types
export interface PimUser {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  displayName: string;
  role: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  accessToken: string;
  expiresIn: number;
  user: PimUser;
}

export interface ApiError {
  code: string;
  message: string;
}
