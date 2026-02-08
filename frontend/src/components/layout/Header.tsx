"use client";

import { useRouter } from "next/navigation";
import { useAuthStore } from "@/lib/stores/authStore";
import { apiClient } from "@/lib/api/client";
import { Button } from "@/components/ui/Button";

export function Header() {
  const router = useRouter();
  const { user, company, logout: storeLogout } = useAuthStore();

  const handleLogout = async () => {
    try {
      await apiClient.logout();
    } catch {
      // Ignore errors
    }
    storeLogout();
    router.push("/login");
  };

  return (
    <header className="bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-700">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center">
            <span className="text-xl font-bold text-primary-600">
              Webshop
            </span>
            {company && (
              <span className="ml-4 text-sm text-gray-500 dark:text-gray-400">
                {company.name}
              </span>
            )}
          </div>

          <div className="flex items-center gap-4">
            {user && (
              <div className="text-sm text-gray-700 dark:text-gray-300">
                <span className="font-medium">{user.displayName}</span>
                <span className="text-gray-500 dark:text-gray-400 ml-2">
                  ({user.role})
                </span>
              </div>
            )}
            <Button variant="outline" size="sm" onClick={handleLogout}>
              Logout
            </Button>
          </div>
        </div>
      </div>
    </header>
  );
}
