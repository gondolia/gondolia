"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { apiClient } from "@/lib/api/client";
import type { Product, PriceScale, Category } from "@/types/catalog";
import { Panel, PanelHeader, PanelBody } from "@/components/ui/Panel";
import { Button } from "@/components/ui/Button";

export default function ProductDetailPage() {
  const params = useParams();
  const router = useRouter();
  const productId = params.id as string;

  const [product, setProduct] = useState<Product | null>(null);
  const [priceScales, setPriceScales] = useState<PriceScale[]>([]);
  const [category, setCategory] = useState<Category | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [quantity, setQuantity] = useState(1);

  useEffect(() => {
    loadProductDetails();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [productId]);

  const loadProductDetails = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const [productData, pricesData] = await Promise.all([
        apiClient.getProduct(productId),
        apiClient.getProductPrices(productId),
      ]);

      setProduct(productData);
      setPriceScales(pricesData);

      // Load category if available
      if (productData.categoryId) {
        try {
          const categoryData = await apiClient.getCategory(
            productData.categoryId
          );
          setCategory(categoryData);
        } catch (err) {
          console.error("Failed to load category:", err);
        }
      }
    } catch (err) {
      const error = err as { message?: string };
      setError(error.message || "Fehler beim Laden des Produkts");
    } finally {
      setIsLoading(false);
    }
  };

  const formatPrice = (price: number, currency: string) => {
    return new Intl.NumberFormat("de-CH", {
      style: "currency",
      currency: currency,
    }).format(price);
  };

  const getCurrentPrice = () => {
    if (!product) return null;

    // Find applicable price scale
    const applicableScale = priceScales
      .filter((scale) => quantity >= scale.minQuantity)
      .sort((a, b) => b.minQuantity - a.minQuantity)[0];

    return applicableScale ? applicableScale.price : product.basePrice;
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
            Produkt wird geladen...
          </p>
        </div>
      </div>
    );
  }

  if (error || !product) {
    return (
      <div className="text-center py-12">
        <p className="text-red-600 dark:text-red-400">
          {error || "Produkt nicht gefunden"}
        </p>
        <Button className="mt-4" onClick={() => router.push("/products")}>
          Zurück zur Produktliste
        </Button>
      </div>
    );
  }

  const currentPrice = getCurrentPrice();

  return (
    <div className="space-y-6">
      {/* Breadcrumb */}
      <nav className="flex items-center text-sm text-gray-500 dark:text-gray-400">
        <Link
          href="/products"
          className="hover:text-primary-600 dark:hover:text-primary-400"
        >
          Produkte
        </Link>
        {category && (
          <>
            <span className="mx-2">/</span>
            <Link
              href={`/products?category=${category.id}`}
              className="hover:text-primary-600 dark:hover:text-primary-400"
            >
              {category.name}
            </Link>
          </>
        )}
        <span className="mx-2">/</span>
        <span className="text-gray-900 dark:text-white">{product.name}</span>
      </nav>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Product Image */}
        <Panel>
          <div className="aspect-square bg-gray-100 dark:bg-gray-800 rounded-lg overflow-hidden">
            {product.imageUrl ? (
              <img
                src={product.imageUrl}
                alt={product.name}
                className="w-full h-full object-contain p-8"
              />
            ) : (
              <div className="w-full h-full flex items-center justify-center text-gray-400 dark:text-gray-600">
                <svg
                  className="w-32 h-32"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={1}
                    d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                  />
                </svg>
              </div>
            )}
          </div>
        </Panel>

        {/* Product Info */}
        <div className="space-y-6">
          <div>
            <div className="text-sm text-gray-500 dark:text-gray-400 font-mono mb-2">
              Artikelnummer: {product.sku}
            </div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">
              {product.name}
            </h1>
            {product.shortDescription && (
              <p className="text-lg text-gray-600 dark:text-gray-400">
                {product.shortDescription}
              </p>
            )}
          </div>

          {/* Price & Stock */}
          <Panel>
            <PanelBody>
              <div className="flex items-end justify-between">
                <div>
                  <div className="text-sm text-gray-500 dark:text-gray-400 mb-1">
                    Preis
                  </div>
                  <div className="text-3xl font-bold text-primary-600 dark:text-primary-400">
                    {currentPrice &&
                      formatPrice(currentPrice, product.currency)}
                  </div>
                  <div className="text-sm text-gray-500 dark:text-gray-400">
                    pro {product.unit}
                  </div>
                </div>
                <div className="text-right">
                  {product.stockQuantity > 0 ? (
                    <div>
                      <span className="inline-block px-3 py-1 text-sm font-medium text-green-800 dark:text-green-200 bg-green-100 dark:bg-green-900/30 rounded-full">
                        Auf Lager
                      </span>
                      <div className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                        {product.stockQuantity} {product.unit} verfügbar
                      </div>
                    </div>
                  ) : (
                    <span className="inline-block px-3 py-1 text-sm font-medium text-red-800 dark:text-red-200 bg-red-100 dark:bg-red-900/30 rounded-full">
                      Nicht verfügbar
                    </span>
                  )}
                </div>
              </div>

              {/* Quantity Selector */}
              <div className="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Menge
                </label>
                <div className="flex items-center gap-3">
                  <input
                    type="number"
                    min={product.minOrderQuantity}
                    value={quantity}
                    onChange={(e) =>
                      setQuantity(Math.max(1, parseInt(e.target.value) || 1))
                    }
                    className="w-24 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                  />
                  <span className="text-gray-600 dark:text-gray-400">
                    {product.unit}
                  </span>
                  {product.minOrderQuantity > 1 && (
                    <span className="text-sm text-gray-500 dark:text-gray-400">
                      (Mindestmenge: {product.minOrderQuantity})
                    </span>
                  )}
                </div>
              </div>

              {/* Total Price */}
              {currentPrice && (
                <div className="mt-4 p-4 bg-gray-50 dark:bg-gray-800/50 rounded-md">
                  <div className="flex items-center justify-between">
                    <span className="text-gray-700 dark:text-gray-300 font-medium">
                      Gesamtpreis:
                    </span>
                    <span className="text-2xl font-bold text-primary-600 dark:text-primary-400">
                      {formatPrice(currentPrice * quantity, product.currency)}
                    </span>
                  </div>
                </div>
              )}

              {/* Add to Cart Button (Placeholder) */}
              <Button className="w-full mt-4" size="lg">
                In den Warenkorb
              </Button>
            </PanelBody>
          </Panel>

          {/* Price Scales */}
          {priceScales.length > 0 && (
            <Panel>
              <PanelHeader>
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                  Staffelpreise
                </h2>
              </PanelHeader>
              <PanelBody>
                <div className="space-y-2">
                  {priceScales.map((scale) => (
                    <div
                      key={scale.id}
                      className="flex items-center justify-between py-2 px-3 rounded-md bg-gray-50 dark:bg-gray-800/50"
                    >
                      <div className="text-sm text-gray-700 dark:text-gray-300">
                        ab {scale.minQuantity} {product.unit}
                        {scale.maxQuantity && ` bis ${scale.maxQuantity}`}
                      </div>
                      <div className="flex items-center gap-3">
                        {scale.discountPercent && (
                          <span className="text-xs text-green-600 dark:text-green-400 font-medium">
                            -{scale.discountPercent}%
                          </span>
                        )}
                        <span className="font-semibold text-gray-900 dark:text-white">
                          {formatPrice(scale.price, scale.currency)}
                        </span>
                      </div>
                    </div>
                  ))}
                </div>
              </PanelBody>
            </Panel>
          )}
        </div>
      </div>

      {/* Product Description */}
      {product.description && (
        <Panel>
          <PanelHeader>
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
              Produktbeschreibung
            </h2>
          </PanelHeader>
          <PanelBody>
            <div className="prose dark:prose-invert max-w-none">
              <p className="text-gray-700 dark:text-gray-300 whitespace-pre-wrap">
                {product.description}
              </p>
            </div>
          </PanelBody>
        </Panel>
      )}

      {/* Product Attributes */}
      {product.attributes && Object.keys(product.attributes).length > 0 && (
        <Panel>
          <PanelHeader>
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
              Technische Daten
            </h2>
          </PanelHeader>
          <PanelBody>
            <dl className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {Object.entries(product.attributes).map(([key, value]) => (
                <div
                  key={key}
                  className="border-b border-gray-200 dark:border-gray-700 pb-3"
                >
                  <dt className="text-sm text-gray-500 dark:text-gray-400 mb-1">
                    {key}
                  </dt>
                  <dd className="text-gray-900 dark:text-white font-medium">
                    {value}
                  </dd>
                </div>
              ))}
            </dl>
          </PanelBody>
        </Panel>
      )}
    </div>
  );
}
