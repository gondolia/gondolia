"use client";

import Link from "next/link";
import type { Product } from "@/types/catalog";
import { Panel } from "@/components/ui/Panel";

interface ProductCardProps {
  product: Product;
}

export function ProductCard({ product }: ProductCardProps) {
  const formatPrice = (price: number, currency: string) => {
    return new Intl.NumberFormat("de-CH", {
      style: "currency",
      currency: currency,
    }).format(price);
  };

  // Determine stock status for variant_parent
  const getStockStatus = () => {
    if (product.productType === 'variant_parent') {
      const hasStock = product.variants?.some(v => v.availability?.inStock);
      return hasStock ? 'available' : 'unavailable';
    }
    return product.stockQuantity > 0 ? 'available' : 'unavailable';
  };

  const stockStatus = getStockStatus();

  return (
    <Link href={`/products/${product.id}`}>
      <Panel className="h-full hover:shadow-lg transition-shadow cursor-pointer">
        <div className="aspect-square bg-gray-100 dark:bg-gray-800 rounded-t-lg overflow-hidden relative">
          {product.imageUrl ? (
            <img
              src={product.imageUrl}
              alt={product.name}
              className="w-full h-full object-cover"
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center text-gray-400 dark:text-gray-600">
              <svg
                className="w-16 h-16"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                />
              </svg>
            </div>
          )}
          {/* Bundle Icon Overlay */}
          {product.productType === 'bundle' && (
            <div className="absolute top-2 right-2 bg-purple-600 text-white p-2 rounded-full shadow-lg">
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </div>
          )}
        </div>
        <div className="p-4 space-y-2">
          <div className="flex items-center justify-between">
            <div className="text-xs text-gray-500 dark:text-gray-400 font-mono">
              {product.sku}
            </div>
            {product.productType === 'variant_parent' && product.variantCount && (
              <div className="text-xs text-primary-600 dark:text-primary-400 font-medium">
                {product.variantCount} Varianten
              </div>
            )}
            {product.productType === 'bundle' && (
              <div className="text-xs bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-300 px-2 py-1 rounded font-medium">
                Bundle
              </div>
            )}
            {product.productType === 'parametric' && (
              <div className="text-xs bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 px-2 py-1 rounded font-medium">
                Parametrisch
              </div>
            )}
          </div>
          <h3 className="font-semibold text-gray-900 dark:text-white line-clamp-2 min-h-[3rem]">
            {product.name}
          </h3>
          {product.productType === 'variant_parent' && product.variantSummary && (
            <div className="flex flex-wrap gap-1">
              {Object.entries(product.variantSummary).slice(0, 2).map(([axis, values]) => (
                <span
                  key={axis}
                  className="text-xs bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400 px-2 py-0.5 rounded"
                >
                  {values.slice(0, 4).join(", ")}
                  {values.length > 4 && ` +${values.length - 4}`}
                </span>
              ))}
            </div>
          )}
          {product.shortDescription && (
            <p className="text-sm text-gray-600 dark:text-gray-400 line-clamp-2">
              {product.shortDescription}
            </p>
          )}
          <div className="flex items-end justify-between pt-2">
            <div>
              {product.productType === 'bundle' ? (
                // Bundle: show price or "ab" price depending on mode
                <div>
                  {product.basePrice > 0 ? (
                    <>
                      <div className="text-lg font-bold text-primary-600 dark:text-primary-400">
                        {product.bundleMode === 'configurable' ? 'ab ' : ''}{formatPrice(product.basePrice, product.currency)}
                      </div>
                      <div className="text-xs text-gray-500 dark:text-gray-400">
                        {product.bundleMode === 'configurable' ? 'Mengen anpassbar' : 'Komplettpaket'}
                      </div>
                    </>
                  ) : (
                    <>
                      <div className="text-lg font-bold text-primary-600 dark:text-primary-400">
                        Bundle
                      </div>
                      <div className="text-xs text-gray-500 dark:text-gray-400">
                        Preis je nach Konfiguration
                      </div>
                    </>
                  )}
                </div>
              ) : product.productType === 'variant_parent' && product.priceRange ? (
                // Show price range for variant_parent
                <div>
                  <div className="text-lg font-bold text-primary-600 dark:text-primary-400">
                    ab {formatPrice(product.priceRange.min, product.priceRange.currency)}
                  </div>
                  {product.priceRange.max !== product.priceRange.min && (
                    <div className="text-xs text-gray-500 dark:text-gray-400">
                      bis {formatPrice(product.priceRange.max, product.priceRange.currency)}
                    </div>
                  )}
                </div>
              ) : (
                // Show specific price for simple products
                <div>
                  <div className="text-lg font-bold text-primary-600 dark:text-primary-400">
                    {formatPrice(product.basePrice, product.currency)}
                  </div>
                  <div className="text-xs text-gray-500 dark:text-gray-400">
                    pro {product.unit}
                  </div>
                </div>
              )}
            </div>
            {stockStatus === 'available' ? (
              <span className="text-xs text-green-600 dark:text-green-400 font-medium">
                Auf Lager
              </span>
            ) : (
              <span className="text-xs text-red-600 dark:text-red-400 font-medium">
                Nicht verfügbar
              </span>
            )}
          </div>
        </div>
      </Panel>
    </Link>
  );
}
