"use client";

import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import { apiClient } from "@/lib/api/client";
import type { Product, Category, PaginatedResponse } from "@/types/catalog";
import { ProductCard } from "@/components/catalog/ProductCard";
import { CategorySidebar } from "@/components/catalog/CategorySidebar";
import { Pagination } from "@/components/catalog/Pagination";
import { Input } from "@/components/ui/Input";
import { Panel } from "@/components/ui/Panel";

export default function ProductsPage() {
  const searchParams = useSearchParams();
  const [products, setProducts] = useState<PaginatedResponse<Product> | null>(
    null
  );
  const [categories, setCategories] = useState<Category[]>([]);
  const [searchQuery, setSearchQuery] = useState(
    searchParams.get("q") || ""
  );
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const currentPage = parseInt(searchParams.get("page") || "1");
  const categoryId = searchParams.get("category") || undefined;
  const productTypeFilter = searchParams.get("type") || undefined;

  useEffect(() => {
    loadCategories();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    loadProducts();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchParams]);

  const loadCategories = async () => {
    try {
      const data = await apiClient.getCategories();
      setCategories(data);
    } catch (err) {
      console.error("Failed to load categories:", err);
    }
  };

  const loadProducts = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const q = searchParams.get("q") || undefined;

      // If a category is selected, use getCategoryProducts with include_children
      if (categoryId) {
        const data = await apiClient.getCategoryProducts(categoryId, {
          page: currentPage,
          limit: 12,
          includeChildren: true, // Backend handles subcategories
        });
        
        // Filter out variant products (only show simple and variant_parent)
        const filteredItems = data.items.filter(
          p => p.productType !== 'variant'
        );
        setProducts({
          ...data,
          items: filteredItems,
        });
      } else {
        // No category filter — use general product search
        const data = await apiClient.getProducts({
          q,
          productType: productTypeFilter,
          page: currentPage,
          limit: 12,
        });
        
        // Filter out variant products (only show simple and variant_parent)
        const filteredItems = data.items.filter(
          p => p.productType !== 'variant'
        );
        setProducts({
          ...data,
          items: filteredItems,
        });
      }
    } catch (err) {
      const error = err as { message?: string };
      setError(error.message || "Fehler beim Laden der Produkte");
    } finally {
      setIsLoading(false);
    }
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    const params = new URLSearchParams(searchParams);
    if (searchQuery) {
      params.set("q", searchQuery);
    } else {
      params.delete("q");
    }
    params.delete("page");
    window.location.href = `/products?${params.toString()}`;
  };

  const handleTypeFilter = (type: string | undefined) => {
    const params = new URLSearchParams(searchParams);
    if (type) {
      params.set("type", type);
    } else {
      params.delete("type");
    }
    params.delete("page");
    window.location.href = `/products?${params.toString()}`;
  };

  const typeFilters = [
    { label: "Alle", value: undefined },
    { label: "Einfach", value: "simple" },
    { label: "Varianten", value: "variant_parent" },
    { label: "Parametrisch", value: "parametric" },
    { label: "Bundles", value: "bundle" },
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
          Produkte
        </h1>
      </div>

      {/* Search */}
      <Panel className="p-4">
        <form onSubmit={handleSearch} className="flex gap-3">
          <Input
            type="search"
            placeholder="Produkte suchen..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="flex-1"
          />
          <button
            type="submit"
            className="px-6 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 transition-colors font-medium"
          >
            Suchen
          </button>
        </form>
      </Panel>

      {/* Product Type Filter */}
      <div className="flex flex-wrap gap-2">
        {typeFilters.map((f) => (
          <button
            key={f.value ?? "all"}
            onClick={() => handleTypeFilter(f.value)}
            className={`
              px-4 py-2 text-sm rounded-full font-medium transition-colors
              ${productTypeFilter === f.value || (!productTypeFilter && !f.value)
                ? "bg-primary-600 text-white shadow-sm"
                : "bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700"
              }
            `}
          >
            {f.label}
          </button>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        {/* Sidebar */}
        <div className="lg:col-span-1">
          <CategorySidebar
            categories={categories}
            selectedCategoryId={categoryId}
          />
        </div>

        {/* Product Grid */}
        <div className="lg:col-span-3">
          {isLoading ? (
            <div className="flex items-center justify-center h-64">
              <div className="text-center">
                <svg
                  className="mx-auto h-12 w-12 animate-spin text-primary-600"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    className="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth="4"
                  />
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  />
                </svg>
                <p className="mt-4 text-gray-600 dark:text-gray-400">
                  Produkte werden geladen...
                </p>
              </div>
            </div>
          ) : error ? (
            <div className="text-center py-12">
              <p className="text-red-600 dark:text-red-400">{error}</p>
            </div>
          ) : products && products.items.length > 0 ? (
            <div className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                {products.items.map((product) => (
                  <ProductCard key={product.id} product={product} />
                ))}
              </div>
              <Pagination
                currentPage={products.page}
                totalPages={products.totalPages}
                totalItems={products.total}
              />
            </div>
          ) : (
            <div className="text-center py-12">
              <svg
                className="mx-auto h-12 w-12 text-gray-400 dark:text-gray-600"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"
                />
              </svg>
              <p className="mt-4 text-gray-600 dark:text-gray-400">
                Keine Produkte gefunden.
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
