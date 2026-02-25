"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { MainLayout } from "@/components/MainLayout";
import { pimApiClient } from "@/lib/api/client";
import type { Product, ProductType, ProductStatus } from "@/types/catalog";
import {
  Plus,
  Search,
  Filter,
  ChevronLeft,
  ChevronRight,
  MoreVertical,
} from "lucide-react";

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [filterType, setFilterType] = useState<ProductType | "">("");
  const [filterStatus, setFilterStatus] = useState<ProductStatus | "">("");
  const [sortBy, setSortBy] = useState("name");
  const [sortOrder, setSortOrder] = useState<"asc" | "desc">("asc");
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);
  const limit = 20;

  useEffect(() => {
    const fetchProducts = async () => {
      setLoading(true);
      try {
        const response = await pimApiClient.getProducts({
          q: searchQuery || undefined,
          productType: filterType || undefined,
          status: filterStatus || undefined,
          sortBy,
          sortOrder,
          page,
          limit,
        });
        setProducts(response.items);
        setTotalPages(response.totalPages);
        setTotal(response.total);
      } catch (error) {
        console.error("Failed to fetch products:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchProducts();
  }, [searchQuery, filterType, filterStatus, sortBy, sortOrder, page]);

  const handleSearch = (query: string) => {
    setSearchQuery(query);
    setPage(1); // Reset to first page on search
  };

  return (
    <MainLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Produkte</h1>
            <p className="mt-1 text-sm text-gray-500">
              {total} Produkte gesamt
            </p>
          </div>
          <Link
            href="/products/new"
            className="flex items-center gap-2 rounded-lg bg-primary-600 px-4 py-2 text-sm font-semibold text-white hover:bg-primary-700"
          >
            <Plus className="h-4 w-4" />
            Neues Produkt
          </Link>
        </div>

        {/* Filters */}
        <div className="rounded-lg bg-white p-4 shadow-sm ring-1 ring-gray-200">
          <div className="grid gap-4 md:grid-cols-4">
            {/* Search */}
            <div className="md:col-span-2">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
                <input
                  type="text"
                  placeholder="Produkt suchen..."
                  value={searchQuery}
                  onChange={(e) => handleSearch(e.target.value)}
                  className="w-full rounded-lg border border-gray-300 pl-10 pr-4 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
                />
              </div>
            </div>

            {/* Type Filter */}
            <div>
              <select
                value={filterType}
                onChange={(e) => {
                  setFilterType(e.target.value as ProductType | "");
                  setPage(1);
                }}
                className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
              >
                <option value="">Alle Typen</option>
                <option value="simple">Simple</option>
                <option value="variant_parent">Variant Parent</option>
                <option value="variant">Variant</option>
                <option value="bundle">Bundle</option>
                <option value="parametric">Parametric</option>
              </select>
            </div>

            {/* Status Filter */}
            <div>
              <select
                value={filterStatus}
                onChange={(e) => {
                  setFilterStatus(e.target.value as ProductStatus | "");
                  setPage(1);
                }}
                className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
              >
                <option value="">Alle Status</option>
                <option value="active">Active</option>
                <option value="inactive">Inactive</option>
                <option value="draft">Draft</option>
              </select>
            </div>
          </div>
        </div>

        {/* Products Table */}
        <div className="rounded-lg bg-white shadow-sm ring-1 ring-gray-200">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="border-b border-gray-200 bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
                    SKU
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
                    Name
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
                    Typ
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
                    Preis
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
                    Geändert
                  </th>
                  <th className="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500">
                    Aktionen
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200 bg-white">
                {loading ? (
                  <tr>
                    <td colSpan={7} className="px-6 py-8 text-center text-sm text-gray-500">
                      Lädt...
                    </td>
                  </tr>
                ) : products.length === 0 ? (
                  <tr>
                    <td colSpan={7} className="px-6 py-8 text-center text-sm text-gray-500">
                      Keine Produkte gefunden
                    </td>
                  </tr>
                ) : (
                  products.map((product) => (
                    <tr key={product.id} className="hover:bg-gray-50">
                      <td className="whitespace-nowrap px-6 py-4 text-sm font-medium text-gray-900">
                        {product.sku}
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-900">
                        <Link
                          href={`/products/${product.id}`}
                          className="font-medium hover:text-primary-600"
                        >
                          {product.name.de || product.name.en || "Unbenannt"}
                        </Link>
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
                        <span className="inline-flex rounded-full bg-gray-100 px-2 py-1 text-xs font-medium text-gray-800">
                          {product.productType.replace("_", " ")}
                        </span>
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-sm">
                        <span
                          className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${
                            product.status === "active"
                              ? "bg-green-100 text-green-800"
                              : product.status === "inactive"
                              ? "bg-red-100 text-red-800"
                              : "bg-yellow-100 text-yellow-800"
                          }`}
                        >
                          {product.status}
                        </span>
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-900">
                        {product.basePrice
                          ? `${product.basePrice.toFixed(2)} ${product.currency}`
                          : "-"}
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
                        {new Date(product.updatedAt).toLocaleDateString("de-DE")}
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-right text-sm">
                        <button className="text-gray-400 hover:text-gray-600">
                          <MoreVertical className="h-5 w-5" />
                        </button>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between border-t border-gray-200 px-6 py-3">
              <div className="text-sm text-gray-500">
                Seite {page} von {totalPages}
              </div>
              <div className="flex gap-2">
                <button
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                  disabled={page === 1}
                  className="rounded-lg border border-gray-300 px-3 py-1 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50"
                >
                  <ChevronLeft className="h-4 w-4" />
                </button>
                <button
                  onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                  disabled={page === totalPages}
                  className="rounded-lg border border-gray-300 px-3 py-1 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50"
                >
                  <ChevronRight className="h-4 w-4" />
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </MainLayout>
  );
}
