"use client";

import { useState, useMemo, useCallback, useEffect } from "react";
import type { VariantAxis, ProductVariant, PriceScale } from "@/types/catalog";

interface VariantMatrixViewProps {
  axes: VariantAxis[];
  variants: ProductVariant[];
  currency?: string;
  unit?: string;
  onSelectVariant?: (variantId: string, axisValues: Record<string, string>) => void;
  onAddToCart?: (items: Array<{ variantId: string; sku: string; quantity: number }>) => void;
  getVariantPrices?: (variantId: string) => Promise<PriceScale[]>;
  className?: string;
}

export function VariantMatrixView({
  axes,
  variants,
  currency = "CHF",
  unit = "Stk",
  onSelectVariant,
  onAddToCart,
  getVariantPrices,
  className = "",
}: VariantMatrixViewProps) {
  const [quantities, setQuantities] = useState<Record<string, number>>({});
  const [sortBy, setSortBy] = useState<string>("sku");
  const [sortDir, setSortDir] = useState<"asc" | "desc">("asc");
  const [filterText, setFilterText] = useState("");
  const [priceScalesCache, setPriceScalesCache] = useState<Record<string, PriceScale[]>>({});

  // Fetch tier prices for variants that have quantities > 0
  useEffect(() => {
    if (!getVariantPrices) return;
    const variantIds = Object.keys(quantities).filter(
      (id) => quantities[id] > 0 && !priceScalesCache[id]
    );
    if (variantIds.length === 0) return;
    
    let cancelled = false;
    Promise.all(
      variantIds.map((id) =>
        getVariantPrices(id).then((scales) => ({ id, scales })).catch(() => ({ id, scales: [] as PriceScale[] }))
      )
    ).then((results) => {
      if (cancelled) return;
      setPriceScalesCache((prev) => {
        const next = { ...prev };
        for (const { id, scales } of results) {
          next[id] = scales;
        }
        return next;
      });
    });
    return () => { cancelled = true; };
  }, [quantities, getVariantPrices, priceScalesCache]);

  // Resolve the best price for a variant given a quantity
  const getEffectivePrice = useCallback(
    (variant: ProductVariant, qty: number): number | null => {
      const scales = priceScalesCache[variant.id];
      if (scales && scales.length > 0 && qty > 0) {
        const applicable = scales
          .filter((s) => qty >= s.minQuantity)
          .sort((a, b) => b.minQuantity - a.minQuantity)[0];
        if (applicable) return applicable.price;
      }
      return variant.price?.net ?? null;
    },
    [priceScalesCache]
  );

  const sortedAxes = useMemo(
    () => [...axes].sort((a, b) => a.position - b.position),
    [axes]
  );

  const getLabel = (labels: Record<string, string>): string =>
    labels.de || labels.en || Object.values(labels)[0] || "";

  // Build axis option label lookup
  const axisOptionLabels = useMemo(() => {
    const lookup: Record<string, Record<string, string>> = {};
    for (const axis of axes) {
      lookup[axis.attributeCode] = {};
      for (const opt of axis.options) {
        lookup[axis.attributeCode][opt.code] = getLabel(opt.label);
      }
    }
    return lookup;
  }, [axes]);

  const getAxisValueLabel = (axisCode: string, optionCode: string): string =>
    axisOptionLabels[axisCode]?.[optionCode] || optionCode;

  // Filter & sort variants
  const displayVariants = useMemo(() => {
    let filtered = variants.filter((v) => v.status === "active");

    if (filterText) {
      const q = filterText.toLowerCase();
      filtered = filtered.filter((v) => {
        if (v.sku.toLowerCase().includes(q)) return true;
        return Object.entries(v.axisValues).some(([axis, opt]) =>
          getAxisValueLabel(axis, opt).toLowerCase().includes(q)
        );
      });
    }

    filtered.sort((a, b) => {
      let cmp = 0;
      if (sortBy === "sku") {
        cmp = a.sku.localeCompare(b.sku);
      } else if (sortBy === "price") {
        cmp = (a.price?.net || 0) - (b.price?.net || 0);
      } else if (sortBy === "stock") {
        cmp = (a.availability?.quantity || 0) - (b.availability?.quantity || 0);
      } else {
        // Sort by axis value
        const aVal = getAxisValueLabel(sortBy, a.axisValues[sortBy] || "");
        const bVal = getAxisValueLabel(sortBy, b.axisValues[sortBy] || "");
        cmp = aVal.localeCompare(bVal);
      }
      return sortDir === "asc" ? cmp : -cmp;
    });

    return filtered;
  }, [variants, filterText, sortBy, sortDir, axisOptionLabels]);

  const handleSort = (column: string) => {
    if (sortBy === column) {
      setSortDir((d) => (d === "asc" ? "desc" : "asc"));
    } else {
      setSortBy(column);
      setSortDir("asc");
    }
  };

  const setQuantity = (variantId: string, qty: number) => {
    setQuantities((prev) => {
      if (qty <= 0) {
        const next = { ...prev };
        delete next[variantId];
        return next;
      }
      return { ...prev, [variantId]: qty };
    });
  };

  const totalItems = Object.values(quantities).reduce((s, q) => s + q, 0);
  const totalPrice = Object.entries(quantities).reduce((sum, [id, qty]) => {
    const v = variants.find((v) => v.id === id);
    if (!v) return sum;
    const price = getEffectivePrice(v, qty) ?? 0;
    return sum + price * qty;
  }, 0);

  const handleAddAll = () => {
    const items = Object.entries(quantities)
      .filter(([_, qty]) => qty > 0)
      .map(([id, qty]) => {
        const v = variants.find((v) => v.id === id)!;
        return { variantId: id, sku: v.sku, quantity: qty };
      });
    if (items.length > 0) {
      onAddToCart?.(items);
    }
  };

  const formatPrice = (price: number) =>
    new Intl.NumberFormat("de-CH", { style: "currency", currency }).format(price);

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Header with filter and summary */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3">
        <div className="relative w-full sm:w-72">
          <input
            type="text"
            value={filterText}
            onChange={(e) => setFilterText(e.target.value)}
            placeholder="Varianten filtern…"
            className="w-full pl-9 pr-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-sm text-gray-900 dark:text-white placeholder-gray-400 focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          />
          <svg className="absolute left-3 top-2.5 w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        <div className="text-sm text-gray-500 dark:text-gray-400">
          {displayVariants.length} Variante{displayVariants.length !== 1 ? "n" : ""}
        </div>
      </div>

      {/* Sort controls */}
      <div className="flex items-center gap-2 text-xs text-gray-500 dark:text-gray-400">
        <span>Sortieren:</span>
        {[
          { key: "sku", label: "SKU" },
          ...sortedAxes.map((a) => ({ key: a.attributeCode, label: getLabel(a.label) })),
          { key: "price", label: "Preis" },
          { key: "stock", label: "Lager" },
        ].map((col) => (
          <button
            key={col.key}
            onClick={() => handleSort(col.key)}
            className={`px-2 py-1 rounded transition-colors ${
              sortBy === col.key
                ? "bg-primary-100 dark:bg-primary-900/30 text-primary-700 dark:text-primary-300 font-medium"
                : "hover:bg-gray-100 dark:hover:bg-gray-800"
            }`}
          >
            {col.label}
            {sortBy === col.key && (
              <span className="ml-0.5">{sortDir === "asc" ? "↑" : "↓"}</span>
            )}
          </button>
        ))}
      </div>

      {/* Compact variant rows */}
      <div className="border border-gray-200 dark:border-gray-700 rounded-lg divide-y divide-gray-100 dark:divide-gray-800">
        {displayVariants.map((variant) => {
          const inStock = variant.availability?.inStock ?? true;
          const stockQty = variant.availability?.quantity;
          const qty = quantities[variant.id] || 0;
          const effectivePrice = getEffectivePrice(variant, qty);
          const basePrice = variant.price?.net;
          const hasTierDiscount = effectivePrice != null && basePrice != null && effectivePrice < basePrice;

          return (
            <div
              key={variant.id}
              className={`
                px-3 py-2 transition-colors
                ${qty > 0
                  ? "bg-primary-50/50 dark:bg-primary-900/10"
                  : "hover:bg-gray-50 dark:hover:bg-gray-800/40"
                }
                ${!inStock ? "opacity-60" : ""}
              `}
            >
              {/* Line 1: SKU + axis values as tags */}
              <div className="flex items-center gap-2 flex-wrap">
                <span className="font-mono text-xs shrink-0 text-gray-500 dark:text-gray-400">
                  {variant.sku}
                </span>
                <span className="text-gray-300 dark:text-gray-600 shrink-0">·</span>
                {sortedAxes.map((axis, i) => (
                  <span key={axis.attributeCode} className="text-sm text-gray-900 dark:text-white">
                    <span className="text-gray-400 dark:text-gray-500 text-xs">{getLabel(axis.label)}:</span>{" "}
                    <span className="font-medium">{getAxisValueLabel(axis.attributeCode, variant.axisValues[axis.attributeCode] || "")}</span>
                    {i < sortedAxes.length - 1 && <span className="text-gray-300 dark:text-gray-600 ml-2">·</span>}
                  </span>
                ))}
              </div>

              {/* Line 2: Price + Stock + Quantity */}
              <div className="flex items-center justify-between mt-1">
                <div className="flex items-center gap-3">
                  {/* Price */}
                  <span className="text-sm">
                    {effectivePrice == null ? (
                      <span className="text-gray-400">–</span>
                    ) : (
                      <>
                        {hasTierDiscount && (
                          <span className="text-xs text-gray-400 line-through mr-1">
                            {formatPrice(basePrice)}
                          </span>
                        )}
                        <span className={`font-semibold ${hasTierDiscount ? "text-green-600 dark:text-green-400" : "text-gray-900 dark:text-white"}`}>
                          {formatPrice(effectivePrice)}
                        </span>
                      </>
                    )}
                  </span>

                  {/* Stock */}
                  {inStock ? (
                    <span className="inline-flex items-center gap-1">
                      <span className="w-1.5 h-1.5 rounded-full bg-green-500" />
                      {stockQty != null && (
                        <span className="text-xs text-gray-500 dark:text-gray-400">{stockQty}</span>
                      )}
                    </span>
                  ) : (
                    <span className="inline-flex items-center gap-1">
                      <span className="w-1.5 h-1.5 rounded-full bg-red-400" />
                      <span className="text-xs text-gray-400">0</span>
                    </span>
                  )}
                </div>

                {/* Quantity */}
                <div onClick={(e) => e.stopPropagation()}>
                  <input
                    type="number"
                    min={0}
                    value={qty || ""}
                    onChange={(e) => setQuantity(variant.id, parseInt(e.target.value) || 0)}
                    disabled={!inStock}
                    placeholder="0"
                    className={`
                      w-16 px-2 py-1 text-center text-sm border rounded-md
                      ${qty > 0
                        ? "border-primary-500 bg-primary-50 dark:bg-primary-900/20 text-primary-700 dark:text-primary-300 font-semibold"
                        : "border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                      }
                      disabled:opacity-40 disabled:cursor-not-allowed
                      focus:ring-2 focus:ring-primary-500 focus:border-primary-500
                    `}
                  />
                </div>
              </div>
            </div>
          );
        })}
        {displayVariants.length === 0 && (
          <div className="px-4 py-8 text-center text-gray-400 dark:text-gray-500">
            Keine Varianten gefunden
          </div>
        )}
      </div>

      {/* Footer with totals and action */}
      {totalItems > 0 && (
        <div className="flex items-center justify-between p-4 bg-primary-50 dark:bg-primary-900/20 border border-primary-200 dark:border-primary-800 rounded-lg">
          <div className="space-y-0.5">
            <div className="text-sm text-gray-600 dark:text-gray-400">
              <span className="font-semibold text-gray-900 dark:text-white">{totalItems}</span> {unit} in{" "}
              <span className="font-semibold text-gray-900 dark:text-white">
                {Object.keys(quantities).filter((k) => quantities[k] > 0).length}
              </span>{" "}
              Variante{Object.keys(quantities).filter((k) => quantities[k] > 0).length !== 1 ? "n" : ""}
            </div>
            <div className="text-lg font-bold text-primary-600 dark:text-primary-400">
              Total: {formatPrice(totalPrice)}
            </div>
          </div>
          <button
            onClick={handleAddAll}
            className="px-6 py-3 bg-primary-600 hover:bg-primary-700 text-white font-semibold rounded-lg shadow-sm transition-colors focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
          >
            Alle in den Warenkorb
          </button>
        </div>
      )}
    </div>
  );
}
