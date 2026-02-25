"use client";

import { useState } from "react";
import type { CartItem } from "@/types/cart";
import { useCart } from "@/context/CartContext";

interface CartItemRowProps {
  item: CartItem;
}

export function CartItemRow({ item }: CartItemRowProps) {
  const { updateQuantity, removeItem } = useCart();
  const [isUpdating, setIsUpdating] = useState(false);

  const formatPrice = (price: number, currency: string) => {
    return new Intl.NumberFormat("de-CH", {
      style: "currency",
      currency: currency || "CHF",
    }).format(price);
  };

  const handleQuantityChange = async (newQuantity: number) => {
    if (newQuantity < 1) return;
    setIsUpdating(true);
    try {
      await updateQuantity(item.id, newQuantity);
    } catch (err) {
      console.error("Failed to update quantity:", err);
    } finally {
      setIsUpdating(false);
    }
  };

  const handleRemove = async () => {
    setIsUpdating(true);
    try {
      await removeItem(item.id);
    } catch (err) {
      console.error("Failed to remove item:", err);
    } finally {
      setIsUpdating(false);
    }
  };

  // Format configuration summary
  const getConfigurationSummary = () => {
    if (!item.configuration) return null;

    const parts: string[] = [];

    // Bundle components
    if (item.configuration.bundleComponents && item.configuration.bundleComponents.length > 0) {
      parts.push(`${item.configuration.bundleComponents.length} Komponenten`);
    }

    // Parametric parameters
    if (item.configuration.parameters && Object.keys(item.configuration.parameters).length > 0) {
      const params = Object.entries(item.configuration.parameters)
        .map(([key, value]) => `${key}: ${value}`)
        .join(", ");
      parts.push(params);
    }

    // Parametric selections
    if (item.configuration.selections && Object.keys(item.configuration.selections).length > 0) {
      const selections = Object.entries(item.configuration.selections)
        .map(([key, value]) => `${key}: ${value}`)
        .join(", ");
      parts.push(selections);
    }

    return parts.length > 0 ? parts.join(" • ") : null;
  };

  const configSummary = getConfigurationSummary();

  return (
    <div className="flex gap-3 p-3 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
      {/* Product Image */}
      <div className="w-20 h-20 flex-shrink-0 bg-gray-100 dark:bg-gray-700 rounded overflow-hidden">
        {item.imageUrl ? (
          <img
            src={item.imageUrl}
            alt={item.productName}
            className="w-full h-full object-cover"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-gray-400">
            <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
              />
            </svg>
          </div>
        )}
      </div>

      {/* Product Info */}
      <div className="flex-1 min-w-0">
        <div className="flex items-start justify-between gap-2">
          <div className="flex-1 min-w-0">
            <h4 className="font-semibold text-sm text-gray-900 dark:text-white line-clamp-2">
              {item.productName}
            </h4>
            <div className="text-xs text-gray-500 dark:text-gray-400 font-mono mt-0.5">
              {item.sku}
            </div>
            {configSummary && (
              <div className="text-xs text-gray-600 dark:text-gray-400 mt-1 line-clamp-1">
                {configSummary}
              </div>
            )}
            {item.productType === "bundle" && (
              <div className="text-xs bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-300 px-2 py-0.5 rounded inline-block mt-1">
                Bundle
              </div>
            )}
            {item.productType === "parametric" && (
              <div className="text-xs bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 px-2 py-0.5 rounded inline-block mt-1">
                Parametrisch
              </div>
            )}
          </div>

          {/* Remove Button */}
          <button
            onClick={handleRemove}
            disabled={isUpdating}
            className="text-gray-400 hover:text-red-600 dark:hover:text-red-400 transition-colors p-1 disabled:opacity-50"
            aria-label="Entfernen"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
              />
            </svg>
          </button>
        </div>

        {/* Quantity and Price */}
        <div className="flex items-center justify-between mt-2">
          <div className="flex items-center gap-2">
            <button
              onClick={() => handleQuantityChange(item.quantity - 1)}
              disabled={isUpdating || item.quantity <= 1}
              className="w-7 h-7 rounded border border-gray-300 dark:border-gray-600 flex items-center justify-center text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
            >
              −
            </button>
            <span className="text-sm font-medium text-gray-900 dark:text-white min-w-[2rem] text-center">
              {item.quantity}
            </span>
            <button
              onClick={() => handleQuantityChange(item.quantity + 1)}
              disabled={isUpdating}
              className="w-7 h-7 rounded border border-gray-300 dark:border-gray-600 flex items-center justify-center text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
            >
              +
            </button>
          </div>

          <div className="text-right">
            <div className="text-sm font-bold text-gray-900 dark:text-white">
              {formatPrice(item.totalPrice, item.currency)}
            </div>
            {item.quantity > 1 && (
              <div className="text-xs text-gray-500 dark:text-gray-400">
                {formatPrice(item.unitPrice, item.currency)} / Stk.
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
