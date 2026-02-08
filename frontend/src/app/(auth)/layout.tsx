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
  const { isAuthenticated, tokens, setAuth, setLoading, logout } = useAuthStore();
  const [isChecking, setIsChecking] = useState(true);

  useEffect(() => {
    const checkAuth = async () => {
      // Wait for hydration
      await new Promise((resolve) => setTimeout(resolve, 0));

      const { tokens: storedTokens, isAuthenticated: isAuth } = useAuthStore.getState();

      if (!storedTokens?.accessToken) {
        setLoading(false);
        router.replace("/login");
        return;
      }

      // If we have tokens but no user info, fetch it
      if (storedTokens && !isAuth) {
        try {
          const me = await apiClient.getMe();
          setAuth(
            me.user,
            me.company,
            me.availableCompanies,
            storedTokens
          );
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

  if (!isAuthenticated && !tokens) {
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
