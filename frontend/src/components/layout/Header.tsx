"use client";

import { useRouter, usePathname } from "next/navigation";
import Link from "next/link";
import { useAuthStore } from "@/lib/stores/authStore";
import { apiClient } from "@/lib/api/client";
import { Button } from "@/components/ui/Button";
import { cn } from "@/lib/utils";

export function Header() {
  const router = useRouter();
  const pathname = usePathname();
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

  const isActive = (path: string) => {
    if (path === "/dashboard") {
      return pathname === path;
    }
    return pathname.startsWith(path);
  };

  const navItems = [
    { label: "Dashboard", href: "/dashboard" },
    { label: "Produkte", href: "/products" },
    { label: "Kategorien", href: "/categories" },
  ];

  return (
    <header className="bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-700">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center gap-8">
            <Link href="/dashboard">
              <span className="text-xl font-bold text-primary-600 cursor-pointer hover:text-primary-700">
                Gondolia
              </span>
            </Link>

            {/* Navigation */}
            <nav className="hidden md:flex items-center gap-1">
              {navItems.map((item) => (
                <Link key={item.href} href={item.href}>
                  <span
                    className={cn(
                      "px-4 py-2 rounded-md text-sm font-medium transition-colors cursor-pointer",
                      isActive(item.href)
                        ? "bg-primary-100 dark:bg-primary-900/30 text-primary-900 dark:text-primary-200"
                        : "text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800"
                    )}
                  >
                    {item.label}
                  </span>
                </Link>
              ))}
            </nav>
          </div>

          <div className="flex items-center gap-4">
            {company && (
              <span className="hidden sm:block text-sm text-gray-500 dark:text-gray-400">
                {company.name}
              </span>
            )}
            {user && (
              <div className="hidden sm:block text-sm text-gray-700 dark:text-gray-300">
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
