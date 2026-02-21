"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import { apiClient } from "@/lib/api/client";
import type { Category, Product, PaginatedResponse } from "@/types/catalog";
import { Panel, PanelBody } from "@/components/ui/Panel";
import { ProductCard } from "@/components/catalog/ProductCard";
import { Pagination } from "@/components/catalog/Pagination";
import { Button } from "@/components/ui/Button";

export default function CategoryDetailPage() {
  const params = useParams();
  const router = useRouter();
  const categoryId = params.id as string;

  const [category, setCategory] = useState<Category | null>(null);
  const [categoryChain, setCategoryChain] = useState<Category[]>([]);
  const [products, setProducts] = useState<PaginatedResponse<Product> | null>(
    null
  );
  const searchParams = useSearchParams();
  const currentPage = parseInt(searchParams.get("page") || "1");
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadCategoryDetails();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [categoryId, currentPage]);

  // Build category chain from leaf to root
  const buildCategoryChain = async (categoryId: string): Promise<Category[]> => {
    const chain: Category[] = [];
    let currentId: string | undefined = categoryId;

    while (currentId) {
      try {
        const cat = await apiClient.getCategory(currentId);
        chain.unshift(cat); // Add to beginning
        currentId = cat.parentId;
      } catch (err) {
        console.error("Failed to load parent category:", err);
        break;
      }
    }

    return chain;
  };

  const loadCategoryDetails = async () => {
    setIsLoading(true);
    setError(null);

    try {
      // Build category chain first
      const chain = await buildCategoryChain(categoryId);
      setCategoryChain(chain);
      
      // Fetch category details
      const categoryData = await apiClient.getCategory(categoryId);
      
      // Fetch direct children with product counts (using List endpoint with parent_id filter)
      const children = await apiClient.getChildCategories(categoryId);
      categoryData.children = children;

      // Fetch products using server-side pagination with include_children=true
      const productsData = await apiClient.getCategoryProducts(categoryId, {
        page: currentPage,
        limit: 12,
        includeChildren: true, // Backend handles subcategory products recursively
      });
      
      // Use product total as category product count (includes all descendants)
      categoryData.productCount = productsData.total;
      
      setProducts(productsData);
      setCategory(categoryData);
    } catch (err) {
      const error = err as { message?: string };
      setError(error.message || "Fehler beim Laden der Kategorie");
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
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
            Kategorie wird geladen...
          </p>
        </div>
      </div>
    );
  }

  if (error || !category) {
    return (
      <div className="text-center py-12">
        <p className="text-red-600 dark:text-red-400">
          {error || "Kategorie nicht gefunden"}
        </p>
        <Button className="mt-4" onClick={() => router.push("/categories")}>
          Zurück zur Kategorieübersicht
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Breadcrumb */}
      <nav className="flex items-center text-sm text-gray-500 dark:text-gray-400">
        <Link
          href="/categories"
          className="hover:text-primary-600 dark:hover:text-primary-400"
        >
          Kategorien
        </Link>
        {categoryChain.slice(0, -1).map((cat) => (
          <span key={cat.id}>
            <span className="mx-2">/</span>
            <Link
              href={`/categories/${cat.id}`}
              className="hover:text-primary-600 dark:hover:text-primary-400"
            >
              {cat.name}
            </Link>
          </span>
        ))}
        <span className="mx-2">/</span>
        <span className="text-gray-900 dark:text-white">{category.name}</span>
      </nav>

      {/* Category Header */}
      <Panel>
        <PanelBody>
          {category.imageUrl && (
            <div className="aspect-video bg-gray-100 dark:bg-gray-800 rounded-lg overflow-hidden mb-6">
              <img
                src={category.imageUrl}
                alt={category.name}
                className="w-full h-full object-cover"
              />
            </div>
          )}
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">
            {category.name}
          </h1>
          {category.description && (
            <p className="text-lg text-gray-600 dark:text-gray-400">
              {category.description}
            </p>
          )}
          {category.productCount !== undefined && (
            <div className="mt-4 text-sm text-gray-500 dark:text-gray-400">
              {category.productCount} Produkte in dieser Kategorie
            </div>
          )}
        </PanelBody>
      </Panel>

      {/* Subcategories */}
      {category.children && category.children.length > 0 && (
        <div>
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
            Unterkategorien
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
            {category.children.map((subCategory) => (
              <Link
                key={subCategory.id}
                href={`/categories/${subCategory.id}`}
              >
                <Panel className="hover:shadow-lg transition-shadow cursor-pointer">
                  <PanelBody className="text-center">
                    <h3 className="font-semibold text-gray-900 dark:text-white mb-1">
                      {subCategory.name}
                    </h3>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      {subCategory.productCount || 0} Produkte
                    </p>
                  </PanelBody>
                </Panel>
              </Link>
            ))}
          </div>
        </div>
      )}

      {/* Products */}
      <div>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
          Produkte
        </h2>

        {products && products.items.length > 0 ? (
          <div className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
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
          <div className="text-center py-12 bg-gray-50 dark:bg-gray-800/50 rounded-lg">
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
              Keine Produkte in dieser Kategorie.
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
