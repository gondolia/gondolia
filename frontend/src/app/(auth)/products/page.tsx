"use client";

import { useEffect, useState, useRef } from "react";
import { useSearchParams } from "next/navigation";
import { apiClient } from "@/lib/api/client";
import type { Product, Category } from "@/types/catalog";
import { ProductCard } from "@/components/catalog/ProductCard";
import { CategorySidebar } from "@/components/catalog/CategorySidebar";
import { Input } from "@/components/ui/Input";
import { Panel } from "@/components/ui/Panel";

const PAGE_SIZE = 12;

export default function ProductsPage() {
  const searchParams = useSearchParams();

  // Accumulated product list for load-more
  const [allProducts, setAllProducts] = useState<Product[]>([]);
  // Total count from backend (used to show progress & decide if more exist)
  const [totalProducts, setTotalProducts] = useState(0);
  // Backend offset for the *next* fetch (tracks how many raw items we've already requested)
  const [nextOffset, setNextOffset] = useState(0);

  const [categories, setCategories] = useState<Category[]>([]);
  const [searchQuery, setSearchQuery] = useState(searchParams.get("q") || "");
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Derive current filter values from URL
  const categoryId = searchParams.get("category") || undefined;
  const productTypeFilter = searchParams.get("type") || undefined;
  const q = searchParams.get("q") || undefined;

  // Use a ref to cancel stale async results when filters change mid-flight
  const loadContextRef = useRef(0);

  useEffect(() => {
    loadCategories();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Reset accumulated list and reload from offset 0 whenever filters change
  useEffect(() => {
    const ctx = ++loadContextRef.current;

    setAllProducts([]);
    setNextOffset(0);
    setTotalProducts(0);
    setError(null);
    setIsLoading(true);

    (async () => {
      try {
        const { filtered, total, rawFetched } = await doFetch(0, categoryId, productTypeFilter, q);
        if (ctx !== loadContextRef.current) return; // stale result — discard
        setAllProducts(filtered);
        setTotalProducts(total);
        setNextOffset(rawFetched);
      } catch (err) {
        if (ctx !== loadContextRef.current) return;
        const e = err as { message?: string };
        setError(e.message || "Fehler beim Laden der Produkte");
      } finally {
        if (ctx === loadContextRef.current) {
          setIsLoading(false);
        }
      }
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [categoryId, productTypeFilter, q]);

  const loadCategories = async () => {
    try {
      const data = await apiClient.getCategories();
      setCategories(data);
    } catch (err) {
      console.error("Failed to load categories:", err);
    }
  };

  /**
   * Fetch one page of products from the backend at the given offset.
   * Returns filtered items (variant children excluded), the backend total, and
   * the raw fetched count (= offset + limit, used as next offset).
   */
  const doFetch = async (
    offset: number,
    catId: string | undefined,
    typeFilter: string | undefined,
    searchQ: string | undefined
  ) => {
    let data;
    if (catId) {
      data = await apiClient.getCategoryProducts(catId, {
        limit: PAGE_SIZE,
        offset,
        includeChildren: true, // Backend handles subcategories server-side
      });
    } else {
      data = await apiClient.getProducts({
        q: searchQ,
        productType: typeFilter,
        limit: PAGE_SIZE,
        offset,
      });
    }

    // Exclude raw variant items (only simple, variant_parent, parametric, bundle shown)
    const filtered = data.items.filter((p) => p.productType !== "variant");

    return {
      filtered,
      total: data.total,
      // Next offset = current offset + how many raw items the backend returned
      // (PAGE_SIZE or less on the last page)
      rawFetched: offset + data.items.length,
    };
  };

  const handleLoadMore = async () => {
    setIsLoadingMore(true);
    try {
      const { filtered, rawFetched } = await doFetch(
        nextOffset,
        categoryId,
        productTypeFilter,
        q
      );
      setAllProducts((prev) => [...prev, ...filtered]);
      setNextOffset(rawFetched);
    } catch (err) {
      console.error("Failed to load more products:", err);
    } finally {
      setIsLoadingMore(false);
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

  // More items exist as long as we haven't fetched everything from the backend
  const hasMore = nextOffset < totalProducts;

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
              ${
                productTypeFilter === f.value ||
                (!productTypeFilter && !f.value)
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
          ) : allProducts.length > 0 ? (
            <div className="space-y-6">
              {/* Total count info */}
              <div className="text-sm text-gray-500 dark:text-gray-400">
                {allProducts.length === totalProducts ? (
                  <span>
                    {totalProducts}{" "}
                    {totalProducts === 1 ? "Produkt" : "Produkte"}
                  </span>
                ) : (
                  <span>
                    {allProducts.length} von {totalProducts} Produkten geladen
                  </span>
                )}
              </div>

              {/* Product grid */}
              <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                {allProducts.map((product) => (
                  <ProductCard key={product.id} product={product} />
                ))}
              </div>

              {/* "Mehr laden" button — only shown while there are more products */}
              {hasMore && (
                <div className="flex justify-center pt-2 border-t border-gray-200 dark:border-gray-700">
                  <button
                    onClick={handleLoadMore}
                    disabled={isLoadingMore}
                    className="mt-4 px-8 py-3 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium flex items-center gap-2"
                  >
                    {isLoadingMore ? (
                      <>
                        <svg
                          className="w-4 h-4 animate-spin"
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
                        Wird geladen…
                      </>
                    ) : (
                      `Mehr Produkte laden (${allProducts.length} von ${totalProducts})`
                    )}
                  </button>
                </div>
              )}
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
