"use client";

import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import Link from "next/link";
import { apiClient } from "@/lib/api/client";
import { Panel } from "@/components/ui/Panel";
import { ProductCard } from "@/components/catalog/ProductCard";
import { Pagination } from "@/components/catalog/Pagination";

interface Product {
  id: string;
  sku: string;
  name: Record<string, string>;
  description?: Record<string, string>;
  status: string;
  product_type?: string;
}

interface SearchResult {
  hits: Product[];
  total_hits: number;
  offset: number;
  limit: number;
  facets?: Record<string, Record<string, number>>;
}

export default function SearchPage() {
  const searchParams = useSearchParams();
  const query = searchParams.get("q") || "";
  const type = searchParams.get("type") || "";
  const status = searchParams.get("status") || "";
  const page = parseInt(searchParams.get("page") || "1", 10);
  const limit = 12;

  const [result, setResult] = useState<SearchResult | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchResults = async () => {
      setLoading(true);
      setError(null);

      try {
        const params = new URLSearchParams();
        if (query) params.set("q", query);
        if (type) params.set("type", type);
        if (status) params.set("status", status);
        params.set("offset", ((page - 1) * limit).toString());
        params.set("limit", limit.toString());

        const data = await apiClient.get<SearchResult>(`/api/v1/search?${params}`);
        setResult(data);
      } catch (err) {
        console.error("Search failed:", err);
        setError("Suche fehlgeschlagen. Bitte versuchen Sie es erneut.");
      } finally {
        setLoading(false);
      }
    };

    fetchResults();
  }, [query, type, status, page]);

  const totalPages = result ? Math.ceil(result.total_hits / limit) : 0;

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
          Suchergebnisse
        </h1>
        {query && (
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Suche nach: <span className="font-medium">{query}</span>
          </p>
        )}
        {result && (
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            {result.total_hits} Produkt{result.total_hits !== 1 ? "e" : ""} gefunden
          </p>
        )}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        {/* Filters Sidebar */}
        <aside className="lg:col-span-1">
          <Panel>
            <h2 className="font-semibold text-gray-900 dark:text-gray-100 mb-4">Filter</h2>

            {/* Product Type Filter */}
            {result?.facets?.product_type && (
              <div className="mb-6">
                <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Produkttyp
                </h3>
                <div className="space-y-2">
                  {Object.entries(result.facets.product_type).map(([value, count]) => (
                    <Link
                      key={value}
                      href={`/search?q=${query}&type=${value}`}
                      className="flex justify-between text-sm text-gray-600 dark:text-gray-400 hover:text-primary-600"
                    >
                      <span>{value}</span>
                      <span className="text-gray-400">({count})</span>
                    </Link>
                  ))}
                </div>
              </div>
            )}

            {/* Status Filter */}
            {result?.facets?.status && (
              <div className="mb-6">
                <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Status
                </h3>
                <div className="space-y-2">
                  {Object.entries(result.facets.status).map(([value, count]) => (
                    <Link
                      key={value}
                      href={`/search?q=${query}&status=${value}`}
                      className="flex justify-between text-sm text-gray-600 dark:text-gray-400 hover:text-primary-600"
                    >
                      <span>{value}</span>
                      <span className="text-gray-400">({count})</span>
                    </Link>
                  ))}
                </div>
              </div>
            )}

            {(type || status) && (
              <Link
                href={`/search?q=${query}`}
                className="text-sm text-primary-600 hover:text-primary-700"
              >
                Alle Filter zurücksetzen
              </Link>
            )}
          </Panel>
        </aside>

        {/* Results */}
        <div className="lg:col-span-3">
          {loading ? (
            <div className="text-center py-12">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"></div>
              <p className="mt-4 text-gray-600 dark:text-gray-400">Suche läuft...</p>
            </div>
          ) : error ? (
            <Panel>
              <div className="text-center py-12 text-red-600">{error}</div>
            </Panel>
          ) : result && result.hits.length > 0 ? (
            <>
              <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                {result.hits.map((product) => (
                  <ProductCard
                    key={product.id}
                    product={{
                      id: product.id,
                      sku: product.sku,
                      name: typeof product.name === 'object' ? (product.name.de || product.name.en || '') : product.name,
                      description: typeof product.description === 'object' ? (product.description.de || product.description.en || '') : (product.description || ''),
                      status: product.status,
                    }}
                  />
                ))}
              </div>

              {totalPages > 1 && (
                <div className="mt-8">
                  <Pagination
                    currentPage={page}
                    totalPages={totalPages}
                    onPageChange={(newPage) => {
                      const params = new URLSearchParams();
                      if (query) params.set("q", query);
                      if (type) params.set("type", type);
                      if (status) params.set("status", status);
                      params.set("page", newPage.toString());
                      window.location.href = `/search?${params}`;
                    }}
                  />
                </div>
              )}
            </>
          ) : (
            <Panel>
              <div className="text-center py-12">
                <svg
                  className="mx-auto h-12 w-12 text-gray-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                  />
                </svg>
                <h3 className="mt-4 text-lg font-medium text-gray-900 dark:text-gray-100">
                  Keine Produkte gefunden
                </h3>
                <p className="mt-2 text-gray-500 dark:text-gray-400">
                  Versuchen Sie es mit anderen Suchbegriffen.
                </p>
              </div>
            </Panel>
          )}
        </div>
      </div>
    </div>
  );
}
