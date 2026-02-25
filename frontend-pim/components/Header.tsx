"use client";

import { useAuthStore } from "@/lib/stores/authStore";

export function Header() {
  const { user } = useAuthStore();

  return (
    <header className="flex h-16 items-center justify-between border-b bg-white px-8">
      <div className="flex items-center gap-4">
        <h2 className="text-lg font-semibold text-gray-900">
          {/* Page title will be set by individual pages */}
        </h2>
      </div>

      <div className="flex items-center gap-4">
        {/* Tenant Switcher - Placeholder for now */}
        <div className="rounded-lg border border-gray-300 px-3 py-1.5 text-sm">
          <span className="font-medium">Tenant:</span>{" "}
          <span className="text-gray-600">Demo</span>
        </div>

        {user && (
          <div className="text-sm">
            <span className="font-medium">{user.displayName}</span>
            <span className="ml-2 rounded-full bg-primary-100 px-2 py-0.5 text-xs text-primary-700">
              {user.role}
            </span>
          </div>
        )}
      </div>
    </header>
  );
}
