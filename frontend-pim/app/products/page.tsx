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
  const [openMenuId, setOpenMenuId] = useState<string | null>(null);
  const [toast, setToast] = useState<{ type: "success" | "error"; message: string } | null>(null);

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

  const showToast = (type: "success" | "error", message: string) => {
    setToast({ type, message });
    setTimeout(() => setToast(null), 3000);
  };

  const handleSearch = (query: string) => {
    setSearchQuery(query);
    setPage(1); // Reset to first page on search
  };

  const handleChangeStatus = async (productId: string, newStatus: ProductStatus) => {
    try {
      await pimApiClient.updateProductStatus(productId, newStatus);
      showToast("success", "Status erfolgreich geändert");
      setOpenMenuId(null);
      // Refresh products
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
    } catch (error: any) {
      showToast("error", error.message || "Fehler beim Ändern des Status");
    }
  };

  const handleDelete = async (productId: string) => {
    if (!confirm("Produkt wirklich löschen?")) return;
    try {
      await pimApiClient.deleteProduct(productId);
      showToast("success", "Produkt erfolgreich gelöscht");
      setOpenMenuId(null);
      // Refresh products
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
    } catch (error: any) {
      showToast("error", error.message || "Fehler beim Löschen");
    }
  };

  return (
    <MainLayout>
      {toast && (
        <div className={`fixed top-4 right-4 z-50 rounded-lg px-4 py-3 shadow-lg ${
          toast.type === "success" ? "bg-green-600" : "bg-red-600"
        } text-white`}>
          {toast.message}
        </div>
      )}

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
                      <td className="whitespace-nowrap px-6 py-4 text-right text-sm relative">
                        <button
                          onClick={() => setOpenMenuId(openMenuId === product.id ? null : product.id)}
                          className="text-gray-400 hover:text-gray-600"
                        >
                          <MoreVertical className="h-5 w-5" />
                        </button>
                        {openMenuId === product.id && (
                          <div className="absolute right-0 top-8 z-10 w-48 rounded-lg bg-white shadow-lg ring-1 ring-gray-200">
                            <div className="py-1">
                              <Link
                                href={`/products/${product.id}`}
                                className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                                onClick={() => setOpenMenuId(null)}
                              >
                                Bearbeiten
                              </Link>
                              <div className="border-t border-gray-100">
                                <button
                                  onClick={() => handleChangeStatus(product.id, "active")}
                                  className="block w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100"
                                  disabled={product.status === "active"}
                                >
                                  Status: Active
                                </button>
                                <button
                                  onClick={() => handleChangeStatus(product.id, "inactive")}
                                  className="block w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100"
                                  disabled={product.status === "inactive"}
                                >
                                  Status: Inactive
                                </button>
                                <button
                                  onClick={() => handleChangeStatus(product.id, "draft")}
                                  className="block w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100"
                                  disabled={product.status === "draft"}
                                >
                                  Status: Draft
                                </button>
                              </div>
                              <div className="border-t border-gray-100">
                                <button
                                  onClick={() => handleDelete(product.id)}
                                  className="block w-full px-4 py-2 text-left text-sm text-red-700 hover:bg-red-50"
                                >
                                  Löschen
                                </button>
                              </div>
                            </div>
                          </div>
                        )}
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
