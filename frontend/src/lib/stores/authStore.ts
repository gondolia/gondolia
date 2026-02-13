"use client";

import { create } from "zustand";
import type { User, Company } from "@/types";

interface AuthState {
  user: User | null;
  company: Company | null;
  availableCompanies: Company[];
  accessToken: string | null;
  tokenExpiresAt: number | null;
  isAuthenticated: boolean;
  isLoading: boolean;

  // Actions
  setAuth: (
    user: User,
    company: Company,
    availableCompanies: Company[],
    accessToken: string,
    expiresIn: number
  ) => void;
  setAccessToken: (accessToken: string, expiresIn: number) => void;
  setUser: (user: User) => void;
  setCompany: (company: Company) => void;
  setLoading: (loading: boolean) => void;
  logout: () => void;
  isTokenExpired: () => boolean;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  company: null,
  availableCompanies: [],
  accessToken: null,
  tokenExpiresAt: null,
  isAuthenticated: false,
  isLoading: true,

  setAuth: (user, company, availableCompanies, accessToken, expiresIn) => {
    set({
      user,
      company,
      availableCompanies,
      accessToken,
      tokenExpiresAt: Date.now() + expiresIn * 1000,
      isAuthenticated: true,
      isLoading: false,
    });
  },

  setAccessToken: (accessToken, expiresIn) => {
    set({
      accessToken,
      tokenExpiresAt: Date.now() + expiresIn * 1000,
    });
  },

  setUser: (user) => {
    set({ user });
  },

  setCompany: (company) => {
    set({ company });
  },

  setLoading: (loading) => {
    set({ isLoading: loading });
  },

  logout: () => {
    set({
      user: null,
      company: null,
      availableCompanies: [],
      accessToken: null,
      tokenExpiresAt: null,
      isAuthenticated: false,
      isLoading: false,
    });
  },

  isTokenExpired: () => {
    const { tokenExpiresAt } = get();
    if (!tokenExpiresAt) return true;
    // Add 30 second buffer
    return Date.now() >= tokenExpiresAt - 30000;
  },
}));
