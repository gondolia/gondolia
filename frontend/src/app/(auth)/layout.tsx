"use client";

import { useEffect, useState, useRef } from "react";
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
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const accessToken = useAuthStore((s) => s.accessToken);
  const [isChecking, setIsChecking] = useState(true);
  const checkingRef = useRef(false);

  useEffect(() => {
    // Prevent double-invocation from React StrictMode
    if (checkingRef.current) return;
    checkingRef.current = true;

    const checkAuth = async () => {
      const { accessToken: storedToken, isAuthenticated: isAuth } = useAuthStore.getState();

      // Already authenticated — no check needed
      if (storedToken && isAuth) {
        setIsChecking(false);
        useAuthStore.getState().setLoading(false);
        return;
      }

      // No access token — try to restore session via refresh token cookie
      try {
        const me = await apiClient.getMe();
        const { accessToken: newToken, tokenExpiresAt } = useAuthStore.getState();

        if (newToken && tokenExpiresAt) {
          const expiresIn = Math.floor((tokenExpiresAt - Date.now()) / 1000);
          useAuthStore.getState().setAuth(
            me.user,
            me.company,
            me.availableCompanies,
            newToken,
            expiresIn
          );
        }
        setIsChecking(false);
        useAuthStore.getState().setLoading(false);
      } catch {
        useAuthStore.getState().logout();
        router.replace("/login");
      }
    };

    checkAuth();

    return () => {
      // Allow re-check if component remounts after real unmount (not StrictMode)
      // Small delay to distinguish StrictMode double-invoke from real remount
      setTimeout(() => { checkingRef.current = false; }, 100);
    };
  }, [router]);

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
