"use client";

import { useState, useEffect, useCallback, useRef, useMemo } from "react";
import type { BundleComponent, BundlePriceResponse } from "@/types/catalog";
import { apiClient } from "@/lib/api/client";
import { ParametricConfigurator } from "./ParametricConfigurator";
import { AddToCartButton } from "@/components/cart/AddToCartButton";

interface BundleConfiguratorProps {
  productId: string;
  components: BundleComponent[];
  bundleMode: 'fixed' | 'configurable';
  bundlePriceMode: 'computed' | 'fixed';
  currency?: string;
  className?: string;
}

interface ComponentQuantity {
  componentId: string; // bundle_component ID
  componentProductId: string; // product ID
  quantity: number;
  parameters?: Record<string, number>;
  selections?: Record<string, string>;
}

export function BundleConfigurator({
  productId,
  components,
  bundleMode,
  bundlePriceMode,
  currency = "CHF",
  className = "",
}: BundleConfiguratorProps) {
  // Sort components by sortOrder
  const sortedComponents = [...components].sort((a, b) => a.sortOrder - b.sortOrder);

  // Initialize quantities from component defaults
  const [quantities, setQuantities] = useState<Record<string, ComponentQuantity>>(() => {
    const initial: Record<string, ComponentQuantity> = {};
    sortedComponents.forEach(comp => {
      initial[comp.componentProductId] = {
        componentId: comp.id,
        componentProductId: comp.componentProductId,
        quantity: comp.quantity,
        parameters: comp.defaultParameters || {},
      };
    });
    return initial;
  });

  const [priceResult, setPriceResult] = useState<BundlePriceResponse | null>(null);
  const [isCalculating, setIsCalculating] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const debounceTimer = useRef<ReturnType<typeof setTimeout> | null>(null);

  // Validation - memoized to prevent unnecessary recalculation
  const { validationErrors, isValid } = useMemo(() => {
    const errors: Record<string, string> = {};
    sortedComponents.forEach(comp => {
      const qty = quantities[comp.componentProductId]?.quantity || 0;
      if (comp.minQuantity !== undefined && qty < comp.minQuantity) {
        errors[comp.componentProductId] = `Min. ${comp.minQuantity}`;
      }
      if (comp.maxQuantity !== undefined && qty > comp.maxQuantity) {
        errors[comp.componentProductId] = `Max. ${comp.maxQuantity}`;
      }
    });
    return {
      validationErrors: errors,
      isValid: Object.keys(errors).length === 0,
    };
  }, [sortedComponents, quantities]);

  // Debounced price calculation
  const calculatePrice = useCallback(async () => {
    if (!isValid) {
      setPriceResult(null);
      return;
    }

    setIsCalculating(true);
    setError(null);

    try {
      const request = {
        components: Object.values(quantities).map(q => ({
          component_id: q.componentId, // Use bundle_component ID, not product ID
          quantity: q.quantity,
          parameters: q.parameters && Object.keys(q.parameters).length > 0 ? q.parameters : undefined,
          selections: q.selections && Object.keys(q.selections).length > 0 ? q.selections : undefined,
        })),
      };

      const result = await apiClient.calculateBundlePrice(productId, request);
      setPriceResult(result);
    } catch (err) {
      const e = err as { message?: string };
      setError(e.message || "Preisberechnung fehlgeschlagen");
      setPriceResult(null);
    } finally {
      setIsCalculating(false);
    }
  }, [productId, quantities, isValid]);

  useEffect(() => {
    if (debounceTimer.current) clearTimeout(debounceTimer.current);
    debounceTimer.current = setTimeout(calculatePrice, 400);
    return () => {
      if (debounceTimer.current) clearTimeout(debounceTimer.current);
    };
  }, [calculatePrice]);

  const handleQuantityChange = (componentProductId: string, quantity: number) => {
    setQuantities(prev => ({
      ...prev,
      [componentProductId]: {
        ...prev[componentProductId],
        quantity,
      },
    }));
  };

  const handleParametersChange = (componentProductId: string, parameters: Record<string, number>) => {
    setQuantities(prev => ({
      ...prev,
      [componentProductId]: {
        ...prev[componentProductId],
        parameters,
      },
    }));
  };

  const handleSelectionsChange = (componentProductId: string, selections: Record<string, string>) => {
    setQuantities(prev => ({
      ...prev,
      [componentProductId]: {
        ...prev[componentProductId],
        selections,
      },
    }));
  };

  const formatPrice = (price: number) =>
    new Intl.NumberFormat("de-CH", { style: "currency", currency }).format(price);

  // Prepare cart configuration
  const getCartConfiguration = () => {
    return {
      bundleComponents: Object.values(quantities).map(q => ({
        componentId: q.componentId,
        quantity: q.quantity,
        parameters: q.parameters && Object.keys(q.parameters).length > 0 ? q.parameters : undefined,
        selections: q.selections && Object.keys(q.selections).length > 0 ? q.selections : undefined,
      })),
    };
  };

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Bundle Mode Info */}
      <div className="flex items-center gap-2 text-sm">
        <span className="px-3 py-1 bg-blue-100 dark:bg-blue-900/20 text-blue-800 dark:text-blue-200 rounded-full font-medium">
          Bundle: {bundleMode === 'fixed' ? 'Feste Zusammenstellung' : 'Konfigurierbar'}
        </span>
        {bundlePriceMode === 'fixed' && (
          <span className="px-3 py-1 bg-green-100 dark:bg-green-900/20 text-green-800 dark:text-green-200 rounded-full font-medium">
            Festpreis
          </span>
        )}
      </div>

      {/* Components List */}
      <div className="space-y-4">
        <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-300 uppercase tracking-wide">
          Bundle-Komponenten
        </h3>

        <div className="space-y-4">
          {sortedComponents.map((comp, idx) => {
            const product = comp.product;
            const qty = quantities[comp.componentProductId]?.quantity || comp.quantity;
            const hasError = !!validationErrors[comp.componentProductId];

            if (!product) {
              return (
                <div key={comp.id} className="p-4 border border-gray-200 dark:border-gray-700 rounded-lg">
                  <div className="text-sm text-gray-500 dark:text-gray-400">
                    Komponente wird geladen... (ID: {comp.componentProductId})
                  </div>
                </div>
              );
            }

            return (
              <div
                key={comp.id}
                className="p-4 border-2 border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800"
              >
                {/* Component Header */}
                <div className="flex items-start gap-4 mb-4">
                  {/* Thumbnail */}
                  <div className="w-20 h-20 flex-shrink-0 bg-gray-100 dark:bg-gray-700 rounded overflow-hidden">
                    {product.imageUrl ? (
                      <img
                        src={product.imageUrl}
                        alt={product.name}
                        className="w-full h-full object-cover"
                      />
                    ) : (
                      <div className="w-full h-full flex items-center justify-center text-gray-400">
                        <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                        </svg>
                      </div>
                    )}
                  </div>

                  {/* Product Info */}
                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between gap-2">
                      <div>
                        <div className="flex items-center gap-2 mb-1">
                          <span className="inline-flex items-center justify-center w-6 h-6 bg-primary-100 dark:bg-primary-900/30 text-primary-700 dark:text-primary-300 rounded-full text-xs font-bold">
                            {idx + 1}
                          </span>
                          <h4 className="font-semibold text-gray-900 dark:text-white">
                            {product.name}
                          </h4>
                        </div>
                        <div className="text-xs text-gray-500 dark:text-gray-400 font-mono">
                          SKU: {product.sku}
                        </div>
                        {product.shortDescription && (
                          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                            {product.shortDescription}
                          </p>
                        )}
                      </div>
                    </div>
                  </div>
                </div>

                {/* Quantity Input (for configurable bundles or parametric products) */}
                {(bundleMode === 'configurable' || product.productType === 'parametric') && (
                  <div className="mb-4">
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Menge {product.unit && `(${product.unit})`}
                    </label>
                    <div className="flex items-center gap-3">
                      <input
                        type="number"
                        min={comp.minQuantity || 1}
                        max={comp.maxQuantity}
                        value={qty}
                        onChange={(e) => handleQuantityChange(comp.componentProductId, parseInt(e.target.value) || 0)}
                        disabled={bundleMode === 'fixed' && product.productType !== 'parametric'}
                        className={`
                          w-32 px-3 py-2 border rounded-lg text-sm
                          ${hasError
                            ? "border-red-400 bg-red-50 dark:bg-red-900/10 text-red-700 dark:text-red-400"
                            : "border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                          }
                          disabled:bg-gray-100 dark:disabled:bg-gray-700 disabled:cursor-not-allowed
                        `}
                      />
                      {(comp.minQuantity !== undefined || comp.maxQuantity !== undefined) && (
                        <span className="text-xs text-gray-500 dark:text-gray-400">
                          {comp.minQuantity !== undefined && comp.maxQuantity !== undefined
                            ? `${comp.minQuantity} – ${comp.maxQuantity}`
                            : comp.minQuantity !== undefined
                            ? `Min. ${comp.minQuantity}`
                            : `Max. ${comp.maxQuantity}`}
                        </span>
                      )}
                      {hasError && (
                        <span className="text-xs text-red-500">{validationErrors[comp.componentProductId]}</span>
                      )}
                    </div>
                  </div>
                )}

                {/* Fixed quantity display */}
                {bundleMode === 'fixed' && product.productType !== 'parametric' && (
                  <div className="mb-4 p-3 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
                    <div className="flex items-center justify-between text-sm">
                      <span className="text-gray-600 dark:text-gray-400">Enthaltene Menge:</span>
                      <span className="font-semibold text-gray-900 dark:text-white">
                        {comp.quantity} {product.unit}
                      </span>
                    </div>
                  </div>
                )}

                {/* Parametric Configurator (if component is parametric) */}
                {product.productType === 'parametric' && product.variantAxes && (
                  <div className="mt-4 p-4 bg-gray-50 dark:bg-gray-700/30 rounded-lg">
                    <h5 className="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">
                      Parametrische Konfiguration
                    </h5>
                    <ParametricConfigurator
                      productId={product.id}
                      axes={product.variantAxes}
                      pricing={product.parametricPricing}
                      currency={currency}
                      className="bg-transparent"
                      embedded
                      onSelectionsChange={(sel) => handleSelectionsChange(comp.componentProductId, sel)}
                      onParametersChange={(params) => handleParametersChange(comp.componentProductId, params)}
                    />
                  </div>
                )}

                {/* Component Price (from price result if available) */}
                {priceResult && priceResult.components.find(c => c.componentId === comp.id) && (
                  <div className="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
                    <div className="flex items-center justify-between text-sm">
                      <span className="text-gray-600 dark:text-gray-400">Komponentenpreis:</span>
                      <span className="font-semibold text-gray-900 dark:text-white">
                        {formatPrice(priceResult.components.find(c => c.componentId === comp.id)!.lineTotal)}
                      </span>
                    </div>
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>

      {/* Price Result */}
      <div className="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 space-y-3">
        {isCalculating && (
          <div className="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
            <svg className="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
            </svg>
            Preis wird berechnet…
          </div>
        )}

        {error && (
          <div className="text-sm text-red-600 dark:text-red-400">{error}</div>
        )}

        {priceResult && !isCalculating && (
          <>
            {/* Component Breakdown */}
            {priceResult.components.length > 0 && (
              <div className="space-y-2 text-sm">
                {priceResult.components.map((comp) => {
                  const component = sortedComponents.find(c => c.id === comp.componentId);
                  const product = component?.product;
                  return (
                    <div key={comp.componentId} className="flex justify-between items-center">
                      <span className="text-gray-600 dark:text-gray-400">
                        {product?.name || comp.sku} ({comp.quantity}×)
                      </span>
                      <span className="font-medium text-gray-900 dark:text-white">
                        {formatPrice(comp.lineTotal)}
                      </span>
                    </div>
                  );
                })}
              </div>
            )}

            {/* Total Price */}
            <div className="border-t border-gray-200 dark:border-gray-700 pt-3">
              <div className="flex justify-between items-center text-lg font-bold">
                <span className="text-gray-700 dark:text-gray-300">Gesamtpreis Bundle:</span>
                <span className="text-primary-600 dark:text-primary-400 text-2xl">
                  {formatPrice(priceResult.total)}
                </span>
              </div>
              {bundlePriceMode === 'fixed' && (
                <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                  Festpreis (unabhängig von Einzelpreisen)
                </div>
              )}
            </div>
          </>
        )}

        {!priceResult && !isCalculating && !error && (
          <div className="text-sm text-gray-400 dark:text-gray-500">
            Konfiguration wird berechnet…
          </div>
        )}
      </div>

      {/* Add to Cart Button */}
      <AddToCartButton
        productId={productId}
        quantity={1}
        configuration={getCartConfiguration()}
        disabled={!priceResult || isCalculating || !isValid}
        size="lg"
        className="w-full"
      >
        {priceResult && isValid
          ? `In den Warenkorb – ${formatPrice(priceResult.total)}`
          : "Bitte konfigurieren"}
      </AddToCartButton>
    </div>
  );
}
