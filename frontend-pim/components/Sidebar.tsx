"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  LayoutDashboard,
  Package,
  FolderTree,
  LogOut,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { useAuthStore } from "@/lib/stores/authStore";
import { pimApiClient } from "@/lib/api/client";

const navigation = [
  { name: "Dashboard", href: "/", icon: LayoutDashboard },
  { name: "Produkte", href: "/products", icon: Package },
  { name: "Kategorien", href: "/categories", icon: FolderTree },
];

export function Sidebar() {
  const pathname = usePathname();
  const { user, logout } = useAuthStore();

  const handleLogout = async () => {
    await pimApiClient.logout();
    logout();
  };

  return (
    <div className="flex h-screen w-64 flex-col bg-sidebar">
      {/* Logo */}
      <div className="flex h-16 items-center justify-center border-b border-sidebar-hover">
        <h1 className="text-xl font-bold text-white">Gondolia PIM</h1>
      </div>

      {/* Navigation */}
      <nav className="flex-1 space-y-1 p-4">
        {navigation.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link
              key={item.name}
              href={item.href}
              className={cn(
                "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                isActive
                  ? "bg-primary-600 text-white"
                  : "text-gray-300 hover:bg-sidebar-hover hover:text-white"
              )}
            >
              <item.icon className="h-5 w-5" />
              {item.name}
            </Link>
          );
        })}
      </nav>

      {/* User Info & Logout */}
      <div className="border-t border-sidebar-hover p-4">
        {user && (
          <div className="mb-3 text-sm text-gray-300">
            <div className="font-medium text-white">{user.displayName}</div>
            <div className="text-xs">{user.email}</div>
          </div>
        )}
        <button
          onClick={handleLogout}
          className="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-gray-300 transition-colors hover:bg-sidebar-hover hover:text-white"
        >
          <LogOut className="h-5 w-5" />
          Abmelden
        </button>
      </div>
    </div>
  );
}
