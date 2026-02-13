"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/lib/stores/authStore";
import { apiClient } from "@/lib/api/client";
import { Header } from "@/components/layout/Header";

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const { isAuthenticated, accessToken, setAuth, setLoading, logout } = useAuthStore();
  const [isChecking, setIsChecking] = useState(true);

  useEffect(() => {
    const checkAuth = async () => {
      // Wait for hydration
      await new Promise((resolve) => setTimeout(resolve, 0));

      const { accessToken: storedToken, isAuthenticated: isAuth } = useAuthStore.getState();

      // No access token means user needs to login
      // Refresh token is in HttpOnly cookie, so we'll try to get user info
      // and if that fails (401), the client will auto-refresh via cookie
      if (!storedToken && !isAuth) {
        try {
          // Try to get user info - this will auto-refresh if cookie exists
          const me = await apiClient.getMe();
          const { accessToken: newToken, tokenExpiresAt } = useAuthStore.getState();
          
          if (newToken && tokenExpiresAt) {
            const expiresIn = Math.floor((tokenExpiresAt - Date.now()) / 1000);
            setAuth(
              me.user,
              me.company,
              me.availableCompanies,
              newToken,
              expiresIn
            );
          }
        } catch {
          logout();
          router.replace("/login");
          return;
        }
      }

      setIsChecking(false);
      setLoading(false);
    };

    checkAuth();
  }, [router, setAuth, setLoading, logout]);

  if (isChecking) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-950">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600" />
      </div>
    );
  }

  if (!isAuthenticated && !accessToken) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-950">
      <Header />
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {children}
      </main>
    </div>
  );
}
