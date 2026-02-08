"use client";

import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import type { User, Company, AuthTokens } from "@/types";

interface AuthState {
  user: User | null;
  company: Company | null;
  availableCompanies: Company[];
  tokens: AuthTokens | null;
  isAuthenticated: boolean;
  isLoading: boolean;

  // Actions
  setAuth: (
    user: User,
    company: Company,
    availableCompanies: Company[],
    tokens: AuthTokens
  ) => void;
  setTokens: (tokens: AuthTokens) => void;
  setUser: (user: User) => void;
  setCompany: (company: Company) => void;
  setLoading: (loading: boolean) => void;
  logout: () => void;
  isTokenExpired: () => boolean;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      company: null,
      availableCompanies: [],
      tokens: null,
      isAuthenticated: false,
      isLoading: true,

      setAuth: (user, company, availableCompanies, tokens) => {
        set({
          user,
          company,
          availableCompanies,
          tokens,
          isAuthenticated: true,
          isLoading: false,
        });
      },

      setTokens: (tokens) => {
        set({ tokens });
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
          tokens: null,
          isAuthenticated: false,
          isLoading: false,
        });
      },

      isTokenExpired: () => {
        const { tokens } = get();
        if (!tokens) return true;
        // Add 30 second buffer
        return Date.now() >= tokens.expiresAt - 30000;
      },
    }),
    {
      name: "auth-storage",
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        tokens: state.tokens,
        user: state.user,
        company: state.company,
        availableCompanies: state.availableCompanies,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);
