"use client";

import { useEffect, ReactNode } from "react";
import { useRouter, usePathname } from "next/navigation";
import { useAuthStore } from "@/lib/stores/authStore";

interface AuthGuardProps {
  children: ReactNode;
}

export function AuthGuard({ children }: AuthGuardProps) {
  const router = useRouter();
  const pathname = usePathname();
  const { isAuthenticated } = useAuthStore();

  useEffect(() => {
    // Redirect to login if not authenticated and not on login page
    if (!isAuthenticated && pathname !== "/login") {
      router.push("/login");
    }
    // Redirect to home if authenticated and on login page
    if (isAuthenticated && pathname === "/login") {
      router.push("/");
    }
  }, [isAuthenticated, pathname, router]);

  return <>{children}</>;
}
