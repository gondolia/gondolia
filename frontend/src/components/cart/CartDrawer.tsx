"use client";

import { useCart } from "@/context/CartContext";
import { CartItemRow } from "./CartItemRow";
import { Button } from "@/components/ui/Button";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

interface CartDrawerProps {
  isOpen: boolean;
  onClose: () => void;
}

export function CartDrawer({ isOpen, onClose }: CartDrawerProps) {
  const { cart, isLoading, clearCart } = useCart();
  const router = useRouter();

  // Prevent body scroll when drawer is open
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "";
    }
    return () => {
      document.body.style.overflow = "";
    };
  }, [isOpen]);

  const formatPrice = (price: number, currency?: string) => {
    // Fallback to CHF if no currency provided, or to first item's currency
    const currencyCode = currency || cart?.items[0]?.currency || "CHF";
    return new Intl.NumberFormat("de-CH", {
      style: "currency",
      currency: currencyCode,
    }).format(price);
  };

  const handleCheckout = () => {
    onClose();
    router.push("/checkout");
  };

  const handleClearCart = async () => {
    if (confirm("Möchten Sie den Warenkorb wirklich leeren?")) {
      try {
        await clearCart();
      } catch (err) {
        console.error("Failed to clear cart:", err);
      }
    }
  };

  if (!isOpen) return null;

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black bg-opacity-50 z-40 transition-opacity"
        onClick={onClose}
      />

      {/* Drawer */}
      <div className="fixed top-0 right-0 bottom-0 w-full sm:w-[450px] bg-white dark:bg-gray-900 shadow-2xl z-50 flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-700">
          <h2 className="text-lg font-bold text-gray-900 dark:text-white">
            Warenkorb
          </h2>
          <button
            onClick={onClose}
            className="p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
            aria-label="Schließen"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto p-4">
          {isLoading ? (
            <div className="flex items-center justify-center py-12">
              <svg
                className="w-8 h-8 animate-spin text-primary-600"
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
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
                />
              </svg>
            </div>
          ) : !cart || cart.items.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-12 text-center">
              <svg
                className="w-16 h-16 text-gray-300 dark:text-gray-600 mb-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"
                />
              </svg>
              <p className="text-gray-500 dark:text-gray-400 mb-4">
                Ihr Warenkorb ist leer
              </p>
              <Button variant="primary" onClick={onClose}>
                Weiter einkaufen
              </Button>
            </div>
          ) : (
            <div className="space-y-3">
              {cart.items.map((item) => (
                <CartItemRow key={item.id} item={item} />
              ))}
            </div>
          )}
        </div>

        {/* Footer */}
        {cart && cart.items.length > 0 && (
          <div className="border-t border-gray-200 dark:border-gray-700 p-4 space-y-4">
            {/* Subtotal */}
            <div className="flex items-center justify-between">
              <span className="text-base font-semibold text-gray-700 dark:text-gray-300">
                Zwischensumme
              </span>
              <span className="text-xl font-bold text-gray-900 dark:text-white">
                {formatPrice(cart.subtotal, cart.currency)}
              </span>
            </div>

            {/* Actions */}
            <div className="space-y-2">
              <Button
                variant="primary"
                size="lg"
                onClick={handleCheckout}
                className="w-full"
              >
                Zur Kasse
              </Button>
              <Button
                variant="outline"
                size="md"
                onClick={handleClearCart}
                className="w-full"
              >
                Warenkorb leeren
              </Button>
            </div>

            <p className="text-xs text-gray-500 dark:text-gray-400 text-center">
              Steuern und Versandkosten werden beim Checkout berechnet
            </p>
          </div>
        )}
      </div>
    </>
  );
}
