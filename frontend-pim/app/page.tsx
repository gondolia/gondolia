"use client";

import { useEffect, useState } from "react";
import { MainLayout } from "@/components/MainLayout";
import { pimApiClient } from "@/lib/api/client";
import { Package, FolderTree, TrendingUp, Clock } from "lucide-react";

interface DashboardStats {
  totalProducts: number;
  productsByType: Record<string, number>;
  totalCategories: number;
  recentChanges: number;
}

export default function DashboardPage() {
  const [stats, setStats] = useState<DashboardStats>({
    totalProducts: 0,
    productsByType: {},
    totalCategories: 0,
    recentChanges: 0,
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        // Fetch products to get counts
        const productsResponse = await pimApiClient.getProducts({ limit: 1 });
        const totalProducts = productsResponse.total;

        // Fetch all products to count by type (in real app, backend would provide this)
        const allProductsResponse = await pimApiClient.getProducts({ limit: 100 });
        const productsByType: Record<string, number> = {};
        allProductsResponse.items.forEach((p) => {
          productsByType[p.productType] = (productsByType[p.productType] || 0) + 1;
        });

        // Fetch categories
        const categories = await pimApiClient.getCategories();
        const flattenCategories = (cats: typeof categories): number => {
          return cats.reduce((acc, cat) => {
            return acc + 1 + (cat.children ? flattenCategories(cat.children) : 0);
          }, 0);
        };
        const totalCategories = flattenCategories(categories);

        setStats({
          totalProducts,
          productsByType,
          totalCategories,
          recentChanges: allProductsResponse.items.length, // Placeholder
        });
      } catch (error) {
        console.error("Failed to fetch dashboard stats:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, []);

  const statCards = [
    {
      label: "Gesamt Produkte",
      value: stats.totalProducts,
      icon: Package,
      color: "bg-blue-500",
    },
    {
      label: "Kategorien",
      value: stats.totalCategories,
      icon: FolderTree,
      color: "bg-green-500",
    },
    {
      label: "Letzte Änderungen",
      value: stats.recentChanges,
      icon: Clock,
      color: "bg-orange-500",
    },
  ];

  return (
    <MainLayout>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
          <p className="mt-1 text-sm text-gray-500">
            Übersicht über Ihr Produktinformationsmanagement
          </p>
        </div>

        {/* Stats Cards */}
        <div className="grid gap-6 md:grid-cols-3">
          {statCards.map((card) => (
            <div
              key={card.label}
              className="rounded-lg bg-white p-6 shadow-sm ring-1 ring-gray-200"
            >
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">
                    {card.label}
                  </p>
                  <p className="mt-2 text-3xl font-semibold text-gray-900">
                    {loading ? "..." : card.value}
                  </p>
                </div>
                <div className={`rounded-lg ${card.color} p-3`}>
                  <card.icon className="h-6 w-6 text-white" />
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* Products by Type */}
        <div className="rounded-lg bg-white p-6 shadow-sm ring-1 ring-gray-200">
          <h2 className="text-lg font-semibold text-gray-900">
            Produkte nach Typ
          </h2>
          <div className="mt-4 space-y-3">
            {loading ? (
              <p className="text-sm text-gray-500">Lädt...</p>
            ) : Object.keys(stats.productsByType).length === 0 ? (
              <p className="text-sm text-gray-500">Keine Produkte vorhanden</p>
            ) : (
              Object.entries(stats.productsByType).map(([type, count]) => (
                <div key={type} className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="h-2 w-2 rounded-full bg-primary-500" />
                    <span className="text-sm font-medium capitalize text-gray-700">
                      {type.replace("_", " ")}
                    </span>
                  </div>
                  <span className="text-sm font-semibold text-gray-900">
                    {count}
                  </span>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </MainLayout>
  );
}
