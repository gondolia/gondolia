import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { PimUser } from "@/types";

interface AuthState {
  user: PimUser | null;
  accessToken: string | null;
  isAuthenticated: boolean;
  setAuth: (user: PimUser, accessToken: string) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      accessToken: null,
      isAuthenticated: false,
      setAuth: (user, accessToken) =>
        set({ user, accessToken, isAuthenticated: true }),
      logout: () =>
        set({ user: null, accessToken: null, isAuthenticated: false }),
    }),
    {
      name: "pim-auth",
    }
  )
);
