"use client";

import { ReactNode } from "react";
import { Sidebar } from "./Sidebar";
import { Header } from "./Header";
import { AuthGuard } from "./AuthGuard";

interface MainLayoutProps {
  children: ReactNode;
}

export function MainLayout({ children }: MainLayoutProps) {
  return (
    <AuthGuard>
      <div className="flex h-screen overflow-hidden bg-gray-50">
        <Sidebar />
        <div className="flex flex-1 flex-col overflow-hidden">
          <Header />
          <main className="flex-1 overflow-y-auto p-8">{children}</main>
        </div>
      </div>
    </AuthGuard>
  );
}
