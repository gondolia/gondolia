"use client";

import { useState, useMemo, useCallback, useEffect, useRef } from "react";
import type { VariantAxis, ParametricPricing, ParametricPriceResponse } from "@/types/catalog";
import { apiClient } from "@/lib/api/client";
import { AddToCartButton } from "@/components/cart/AddToCartButton";

interface ParametricConfiguratorProps {
  productId: string;
  axes: VariantAxis[];
  pricing?: ParametricPricing;
  currency?: string;
  className?: string;
  onSelectionsChange?: (selections: Record<string, string>) => void;
  onParametersChange?: (parameters: Record<string, number>) => void;
  /** When true, hides the standalone price display + add-to-cart (used inside bundles) */
  embedded?: boolean;
}

export function ParametricConfigurator({
  productId,
  axes,
  pricing,
  currency = "CHF",
  className = "",
  onSelectionsChange,
  onParametersChange,
  embedded = false,
}: ParametricConfiguratorProps) {
  const sortedAxes = useMemo(
    () => [...axes].sort((a, b) => a.position - b.position),
    [axes]
  );

  const selectAxes = useMemo(() => sortedAxes.filter((a) => a.inputType !== "range"), [sortedAxes]);
  const rangeAxes = useMemo(() => sortedAxes.filter((a) => a.inputType === "range"), [sortedAxes]);

  // State
  const [selections, setSelections] = useState<Record<string, string>>(() => {
    const initial: Record<string, string> = {};
    for (const axis of selectAxes) {
      if (axis.options.length > 0) {
        initial[axis.attributeCode] = axis.options[0].code;
      }
    }
    return initial;
  });

  const [parameters, setParameters] = useState<Record<string, number>>(() => {
    const initial: Record<string, number> = {};
    for (const axis of rangeAxes) {
      initial[axis.attributeCode] = axis.minValue ?? 0;
    }
    return initial;
  });

  // Notify parent of initial values
  useEffect(() => {
    onSelectionsChange?.(selections);
    onParametersChange?.(parameters);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []); // Only on mount

  const [quantity, setQuantity] = useState(1);
  const [priceResult, setPriceResult] = useState<ParametricPriceResponse | null>(null);
  const [isCalculating, setIsCalculating] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const debounceTimer = useRef<ReturnType<typeof setTimeout> | null>(null);

  const getLabel = (labels: Record<string, string>): string =>
    labels.de || labels.en || Object.values(labels)[0] || "";

  // Validation
  const validationErrors = useMemo(() => {
    const errors: Record<string, string> = {};
    for (const axis of rangeAxes) {
      const val = parameters[axis.attributeCode];
      if (val === undefined || val === null) {
        errors[axis.attributeCode] = "Pflichtfeld";
        continue;
      }
      if (axis.minValue !== undefined && val < axis.minValue) {
        errors[axis.attributeCode] = `Min. ${axis.minValue} ${axis.unit || ""}`;
      }
      if (axis.maxValue !== undefined && val > axis.maxValue) {
        errors[axis.attributeCode] = `Max. ${axis.maxValue} ${axis.unit || ""}`;
      }
    }
    return errors;
  }, [parameters, rangeAxes]);

  const isValid = Object.keys(validationErrors).length === 0;

  // Debounced price calculation
  const calculatePrice = useCallback(async () => {
    if (!isValid) {
      setPriceResult(null);
      return;
    }
    setIsCalculating(true);
    setError(null);
    try {
      const result = await apiClient.calculateParametricPrice(productId, parameters, selections, quantity);
      setPriceResult(result);
    } catch (err) {
      const e = err as { message?: string };
      setError(e.message || "Preisberechnung fehlgeschlagen");
      setPriceResult(null);
    } finally {
      setIsCalculating(false);
    }
  }, [productId, parameters, selections, quantity, isValid]);

  useEffect(() => {
    if (debounceTimer.current) clearTimeout(debounceTimer.current);
    debounceTimer.current = setTimeout(calculatePrice, 400);
    return () => {
      if (debounceTimer.current) clearTimeout(debounceTimer.current);
    };
  }, [calculatePrice]);

  const handleParameterChange = (code: string, value: number) => {
    setParameters((prev) => {
      const next = { ...prev, [code]: value };
      onParametersChange?.(next);
      return next;
    });
  };

  const handleSelectionChange = (code: string, value: string) => {
    setSelections((prev) => {
      const next = { ...prev, [code]: value };
      onSelectionsChange?.(next);
      return next;
    });
  };

  const formatPrice = (price: number) =>
    new Intl.NumberFormat("de-CH", { style: "currency", currency }).format(price);

  const formulaLabel: Record<string, string> = {
    per_m2: "pro m²",
    per_running_meter: "pro Laufmeter",
    per_unit: "pro Stück",
    fixed: "Festpreis",
  };

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Select axes (dropdowns) */}
      {selectAxes.length > 0 && (
        <div className="space-y-4">
          {selectAxes.map((axis) => (
            <div key={axis.attributeCode}>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                {getLabel(axis.label)}
              </label>
              <div className="flex flex-wrap gap-2">
                {axis.options.map((opt) => {
                  const isSelected = selections[axis.attributeCode] === opt.code;
                  return (
                    <button
                      key={opt.code}
                      type="button"
                      onClick={() => handleSelectionChange(axis.attributeCode, opt.code)}
                      className={`
                        px-4 py-2 text-sm rounded-lg border-2 transition-all
                        ${isSelected
                          ? "border-primary-500 bg-primary-50 dark:bg-primary-900/20 text-primary-700 dark:text-primary-300 font-medium shadow-sm"
                          : "border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600"
                        }
                      `}
                    >
                      {getLabel(opt.label)}
                    </button>
                  );
                })}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Divider */}
      {selectAxes.length > 0 && rangeAxes.length > 0 && (
        <div className="border-t border-gray-200 dark:border-gray-700" />
      )}

      {/* Range axes (numeric inputs) */}
      {rangeAxes.length > 0 && (
        <div className="space-y-4">
          <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-300 uppercase tracking-wide">
            Masse angeben
          </h3>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            {rangeAxes.map((axis) => {
              const val = parameters[axis.attributeCode] ?? axis.minValue ?? 0;
              const hasError = !!validationErrors[axis.attributeCode];

              return (
                <div key={axis.attributeCode}>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    {getLabel(axis.label)}
                    {axis.unit && (
                      <span className="text-gray-400 dark:text-gray-500 ml-1">({axis.unit})</span>
                    )}
                  </label>
                  <div className="relative">
                    <input
                      type="number"
                      min={axis.minValue}
                      max={axis.maxValue}
                      step={axis.stepValue || 1}
                      value={val}
                      onChange={(e) => handleParameterChange(axis.attributeCode, parseFloat(e.target.value) || 0)}
                      className={`
                        w-full px-3 py-2 border rounded-lg text-sm
                        ${hasError
                          ? "border-red-400 bg-red-50 dark:bg-red-900/10 text-red-700 dark:text-red-400 focus:ring-red-500"
                          : "border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-primary-500 focus:border-primary-500"
                        }
                        focus:ring-2
                      `}
                    />
                    {axis.unit && (
                      <span className="absolute right-3 top-2.5 text-xs text-gray-400 dark:text-gray-500 pointer-events-none">
                        {axis.unit}
                      </span>
                    )}
                  </div>
                  <div className="flex justify-between mt-1">
                    {hasError ? (
                      <span className="text-xs text-red-500">{validationErrors[axis.attributeCode]}</span>
                    ) : (
                      <span className="text-xs text-gray-400 dark:text-gray-500">
                        {axis.minValue !== undefined && axis.maxValue !== undefined
                          ? `${axis.minValue} – ${axis.maxValue} ${axis.unit || ""}`
                          : ""}
                      </span>
                    )}
                    {axis.stepValue && axis.stepValue !== 1 && (
                      <span className="text-xs text-gray-400 dark:text-gray-500">
                        Schritt: {axis.stepValue} {axis.unit || ""}
                      </span>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* Quantity (hidden in embedded/bundle mode) */}
      {!embedded && (
      <div className="border-t border-gray-200 dark:border-gray-700 pt-4">
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Menge
        </label>
        <input
          type="number"
          min={1}
          value={quantity}
          onChange={(e) => setQuantity(Math.max(1, parseInt(e.target.value) || 1))}
          className="w-24 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-sm focus:ring-2 focus:ring-primary-500"
        />
      </div>
      )}

      {/* Price result (hidden in embedded/bundle mode) */}
      {!embedded && (
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
            {/* Resolved SKU */}
            {priceResult.sku && (
              <div className="flex justify-between items-center text-sm">
                <span className="text-gray-500 dark:text-gray-400">Artikelnummer</span>
                <span className="font-mono text-xs bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded text-gray-800 dark:text-gray-200">
                  {priceResult.sku}
                </span>
              </div>
            )}

            {/* Breakdown */}
            {priceResult.breakdown && Object.keys(priceResult.breakdown).length > 0 && (
              <div className="space-y-1 text-xs text-gray-500 dark:text-gray-400">
                {priceResult.breakdown.area_m2 !== undefined && (
                  <div className="flex justify-between">
                    <span>Fläche</span>
                    <span>{(priceResult.breakdown.area_m2).toFixed(4)} m²</span>
                  </div>
                )}
                {priceResult.breakdown.length_m !== undefined && (
                  <div className="flex justify-between">
                    <span>Länge</span>
                    <span>{priceResult.breakdown.length_m} m</span>
                  </div>
                )}
                {priceResult.breakdown.thickness_factor !== undefined && (
                  <div className="flex justify-between">
                    <span>Stärkefaktor</span>
                    <span>×{priceResult.breakdown.thickness_factor}</span>
                  </div>
                )}
                {pricing && (
                  <div className="flex justify-between">
                    <span>Berechnungsart</span>
                    <span>{formulaLabel[pricing.formulaType] || pricing.formulaType}</span>
                  </div>
                )}
              </div>
            )}

            {/* Prices */}
            <div className="border-t border-gray-200 dark:border-gray-700 pt-3 space-y-1">
              <div className="flex justify-between text-sm">
                <span className="text-gray-600 dark:text-gray-400">Stückpreis</span>
                <span className="font-semibold text-gray-900 dark:text-white">
                  {formatPrice(priceResult.unitPrice)}
                </span>
              </div>
              {quantity > 1 && (
                <div className="flex justify-between text-lg font-bold">
                  <span className="text-gray-700 dark:text-gray-300">Gesamtpreis ({quantity}×)</span>
                  <span className="text-primary-600 dark:text-primary-400">
                    {formatPrice(priceResult.totalPrice)}
                  </span>
                </div>
              )}
              {quantity === 1 && (
                <div className="flex justify-between text-lg font-bold">
                  <span className="text-gray-700 dark:text-gray-300">Preis</span>
                  <span className="text-primary-600 dark:text-primary-400">
                    {formatPrice(priceResult.totalPrice)}
                  </span>
                </div>
              )}
            </div>
          </>
        )}

        {!priceResult && !isCalculating && !error && (
          <div className="text-sm text-gray-400 dark:text-gray-500">
            Bitte alle Masse angeben für die Preisberechnung
          </div>
        )}
      </div>
      )}

      {/* Add to Cart (hidden in embedded/bundle mode) */}
      {!embedded && (
      <AddToCartButton
        productId={productId}
        quantity={quantity}
        configuration={{
          parameters,
          selections,
        }}
        disabled={!priceResult || isCalculating}
        size="lg"
        className="w-full"
      >
        {priceResult
          ? `In den Warenkorb – ${formatPrice(priceResult.totalPrice)}${priceResult.sku ? ` (${priceResult.sku})` : ""}`
          : "Bitte konfigurieren"}
      </AddToCartButton>
      )}
    </div>
  );
}
