"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import { apiClient } from "@/lib/api/client";
import type { Product, PriceScale, Category, AxisOption, AttributeLabels } from "@/types/catalog";
import { Panel, PanelHeader, PanelBody } from "@/components/ui/Panel";
import { Button } from "@/components/ui/Button";
import { VariantSelector } from "@/components/catalog/VariantSelector";
import { VariantMatrixView } from "@/components/catalog/VariantMatrixView";

export default function ProductDetailPage() {
  const params = useParams();
  const router = useRouter();
  const searchParams = useSearchParams();
  const productId = params.id as string;

  const [product, setProduct] = useState<Product | null>(null);
  const [priceScales, setPriceScales] = useState<PriceScale[]>([]);
  const [category, setCategory] = useState<Category | null>(null);
  const [categoryChain, setCategoryChain] = useState<Category[]>([]);
  const [attributeLabels, setAttributeLabels] = useState<AttributeLabels>({});
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [quantity, setQuantity] = useState(1);
  
  // Variant state
  const [selectedAxisValues, setSelectedAxisValues] = useState<Record<string, string>>({});
  const [availableAxisValues, setAvailableAxisValues] = useState<Record<string, AxisOption[]>>({});
  const [selectedVariant, setSelectedVariant] = useState<Product | null>(null);
  const [isLoadingVariant, setIsLoadingVariant] = useState(false);
  const [viewMode, setViewMode] = useState<'selector' | 'matrix'>('selector');

  useEffect(() => {
    loadProductDetails();
    // Load attribute translations once (locale DE)
    apiClient.getAttributeTranslations("de").then(setAttributeLabels).catch(() => {});
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [productId]);

  // Load available axis values when selection changes
  useEffect(() => {
    if (product?.productType === 'variant_parent' && Object.keys(selectedAxisValues).length > 0) {
      loadAvailableAxisValues();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedAxisValues]);

  // Try to select variant when all axes are selected
  useEffect(() => {
    if (product?.productType === 'variant_parent' && product.variantAxes) {
      const allAxesSelected = product.variantAxes.every(
        (axis) => selectedAxisValues[axis.attributeCode]
      );
      if (allAxesSelected) {
        selectVariant();
      } else {
        setSelectedVariant(null);
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedAxisValues]);

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

  const loadProductDetails = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const productData = await apiClient.getProduct(productId);
      setProduct(productData);

      // Load prices (only for simple or variant products, not variant_parent)
      if (productData.productType !== 'variant_parent') {
        const pricesData = await apiClient.getProductPrices(productId);
        setPriceScales(pricesData);
      }

      // Load category chain if available
      if (productData.categoryId) {
        try {
          const chain = await buildCategoryChain(productData.categoryId);
          setCategoryChain(chain);
          setCategory(chain[chain.length - 1] || null);
        } catch (err) {
          console.error("Failed to load category:", err);
        }
      }

      // Initialize variant selection from URL params if present
      if (productData.productType === 'variant_parent') {
        const initialSelection: Record<string, string> = {};
        const variantSku = searchParams.get('variant');
        
        if (variantSku && productData.variants) {
          // Find variant by SKU and pre-select its axis values
          const variant = productData.variants.find(v => v.sku === variantSku);
          if (variant) {
            Object.assign(initialSelection, variant.axisValues);
          }
        }
        
        setSelectedAxisValues(initialSelection);
      }
    } catch (err) {
      const error = err as { message?: string };
      setError(error.message || "Fehler beim Laden des Produkts");
    } finally {
      setIsLoading(false);
    }
  };

  const loadAvailableAxisValues = async () => {
    if (!product) return;
    
    try {
      const result = await apiClient.getAvailableAxisValues(productId, selectedAxisValues);
      setAvailableAxisValues(result.available);
    } catch (err) {
      console.error("Failed to load available axis values:", err);
    }
  };

  const selectVariant = async () => {
    if (!product) return;
    
    setIsLoadingVariant(true);
    try {
      const variant = await apiClient.selectVariant(productId, selectedAxisValues);
      setSelectedVariant(variant);
      
      // Load prices for the selected variant
      const pricesData = await apiClient.getProductPrices(variant.id);
      setPriceScales(pricesData);
      
      // Update URL with variant SKU (for sharing/bookmarking)
      const url = new URL(window.location.href);
      url.searchParams.set('variant', variant.sku);
      window.history.replaceState({}, '', url.toString());
    } catch (err) {
      console.error("Failed to select variant:", err);
      setSelectedVariant(null);
      setPriceScales([]);
    } finally {
      setIsLoadingVariant(false);
    }
  };

  const handleAxisSelection = (axisCode: string, optionCode: string) => {
    setSelectedAxisValues((prev) => ({
      ...prev,
      [axisCode]: optionCode,
    }));
  };

  // Resolve a human-readable label for an attribute key.
  // Priority: 1) API translation  2) prettified key (snake_case → Title Case)
  const getAttributeLabel = (key: string): string => {
    if (attributeLabels[key]) return attributeLabels[key];
    // Prettify: remove trailing _unit suffix variations and capitalise words
    return key
      .replace(/_/g, " ")
      .replace(/\b\w/g, (c) => c.toUpperCase());
  };

  const formatPrice = (price: number, currency: string) => {
    return new Intl.NumberFormat("de-CH", {
      style: "currency",
      currency: currency,
    }).format(price);
  };

  const getCurrentPrice = () => {
    // For variant_parent without selection, show price range
    if (product?.productType === 'variant_parent' && !selectedVariant) {
      return null;
    }

    const displayProduct = selectedVariant || product;
    if (!displayProduct) return null;

    // Find applicable price scale
    const applicableScale = priceScales
      .filter((scale) => quantity >= scale.minQuantity)
      .sort((a, b) => b.minQuantity - a.minQuantity)[0];

    return applicableScale ? applicableScale.price : displayProduct.basePrice;
  };

  const getDisplayProduct = (): Product | null => {
    // For variant_parent, display selected variant or parent
    if (product?.productType === 'variant_parent') {
      return selectedVariant || product;
    }
    return product;
  };

  const getDisplayImage = (): string | undefined => {
    // If variant has its own image, use it
    if (selectedVariant?.imageUrl) {
      return selectedVariant.imageUrl;
    }
    // Otherwise use parent/product image
    return product?.imageUrl;
  };

  const getDisplaySku = (): string | undefined => {
    if (selectedVariant) {
      return selectedVariant.sku;
    }
    return product?.sku;
  };

  const getStockInfo = (): { available: boolean; quantity: number } => {
    if (product?.productType === 'variant_parent') {
      if (selectedVariant) {
        return {
          available: selectedVariant.stockQuantity > 0,
          quantity: selectedVariant.stockQuantity,
        };
      }
      // No variant selected - check if any variant is in stock
      const hasStock = product.variants?.some(v => v.availability?.inStock);
      return { available: hasStock || false, quantity: 0 };
    }
    
    return {
      available: product?.stockQuantity ? product.stockQuantity > 0 : false,
      quantity: product?.stockQuantity || 0,
    };
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
  const displayProduct = getDisplayProduct();
  const displayImage = getDisplayImage();
  const displaySku = getDisplaySku();
  const stockInfo = getStockInfo();

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
        {categoryChain.map((cat) => (
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
        <span className="text-gray-900 dark:text-white">{product.name}</span>
      </nav>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Product Image */}
        <Panel>
          <div className="aspect-square bg-gray-100 dark:bg-gray-800 rounded-lg overflow-hidden">
            {displayImage ? (
              <img
                src={displayImage}
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
              Artikelnummer: {displaySku}
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

          {/* Variant Selector / Matrix Toggle (only for variant_parent) */}
          {product.productType === 'variant_parent' && product.variantAxes && (
            <Panel>
              <PanelHeader>
                <div className="flex items-center justify-between">
                  <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                    {viewMode === 'selector' ? 'Variante wählen' : 'Alle Varianten'}
                  </h2>
                  <div className="flex items-center bg-gray-100 dark:bg-gray-800 rounded-lg p-0.5">
                    <button
                      type="button"
                      onClick={() => setViewMode('selector')}
                      className={`px-3 py-1.5 text-xs font-medium rounded-md transition-colors ${
                        viewMode === 'selector'
                          ? 'bg-white dark:bg-gray-700 text-gray-900 dark:text-white shadow-sm'
                          : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300'
                      }`}
                    >
                      Auswahl
                    </button>
                    <button
                      type="button"
                      onClick={() => setViewMode('matrix')}
                      className={`px-3 py-1.5 text-xs font-medium rounded-md transition-colors ${
                        viewMode === 'matrix'
                          ? 'bg-white dark:bg-gray-700 text-gray-900 dark:text-white shadow-sm'
                          : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300'
                      }`}
                    >
                      Matrix
                    </button>
                  </div>
                </div>
              </PanelHeader>
              <PanelBody>
                {viewMode === 'selector' ? (
                  <>
                    <VariantSelector
                      axes={product.variantAxes}
                      selectedValues={selectedAxisValues}
                      onSelect={handleAxisSelection}
                      variants={product.variants}
                      availableValues={availableAxisValues}
                    />
                    {isLoadingVariant && (
                      <div className="mt-4 text-sm text-gray-500 dark:text-gray-400">
                        Variante wird geladen...
                      </div>
                    )}
                  </>
                ) : (
                  <VariantMatrixView
                    axes={product.variantAxes}
                    variants={product.variants || []}
                    currency={product.currency}
                    unit={product.unit}
                  />
                )}
              </PanelBody>
            </Panel>
          )}

          {/* Price & Stock */}
          <Panel>
            <PanelBody>
              <div className="flex items-end justify-between">
                <div>
                  <div className="text-sm text-gray-500 dark:text-gray-400 mb-1">
                    Preis
                  </div>
                  {product.productType === 'variant_parent' && !selectedVariant ? (
                    // Show price range for variant_parent
                    product.priceRange ? (
                      <div>
                        <div className="text-3xl font-bold text-primary-600 dark:text-primary-400">
                          ab {formatPrice(product.priceRange.min, product.priceRange.currency)}
                        </div>
                        {product.priceRange.max !== product.priceRange.min && (
                          <div className="text-sm text-gray-500 dark:text-gray-400">
                            bis {formatPrice(product.priceRange.max, product.priceRange.currency)}
                          </div>
                        )}
                      </div>
                    ) : (
                      <div className="text-lg text-gray-500 dark:text-gray-400">
                        Bitte wählen Sie eine Variante
                      </div>
                    )
                  ) : (
                    // Show specific price for simple or selected variant
                    <div>
                      <div className="text-3xl font-bold text-primary-600 dark:text-primary-400">
                        {currentPrice && displayProduct &&
                          formatPrice(currentPrice, displayProduct.currency)}
                      </div>
                      <div className="text-sm text-gray-500 dark:text-gray-400">
                        pro {displayProduct?.unit || product.unit}
                      </div>
                    </div>
                  )}
                </div>
                <div className="text-right">
                  {stockInfo.available ? (
                    <div>
                      <span className="inline-block px-3 py-1 text-sm font-medium text-green-800 dark:text-green-200 bg-green-100 dark:bg-green-900/30 rounded-full">
                        Auf Lager
                      </span>
                      {stockInfo.quantity > 0 && (
                        <div className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                          {stockInfo.quantity} {displayProduct?.unit || product.unit} verfügbar
                        </div>
                      )}
                    </div>
                  ) : product.productType === 'variant_parent' && !selectedVariant ? (
                    <span className="inline-block px-3 py-1 text-sm font-medium text-gray-600 dark:text-gray-400 bg-gray-100 dark:bg-gray-800 rounded-full">
                      Variante wählen
                    </span>
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
              <Button 
                className="w-full mt-4" 
                size="lg"
                disabled={
                  product.productType === 'variant_parent' && !selectedVariant
                }
              >
                {product.productType === 'variant_parent' && !selectedVariant
                  ? 'Bitte Variante wählen'
                  : 'In den Warenkorb'
                }
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
                    {getAttributeLabel(key)}
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
