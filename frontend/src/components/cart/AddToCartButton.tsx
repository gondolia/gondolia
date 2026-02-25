"use client";

import { useState } from "react";
import { Button } from "@/components/ui/Button";
import { Toast } from "@/components/ui/Toast";
import { useCart } from "@/context/CartContext";
import type { AddToCartRequest } from "@/types/cart";

interface AddToCartButtonProps {
  productId: string;
  variantId?: string;
  quantity?: number;
  configuration?: AddToCartRequest["configuration"];
  disabled?: boolean;
  children?: React.ReactNode;
  variant?: "primary" | "secondary" | "outline";
  size?: "sm" | "md" | "lg";
  className?: string;
  showIcon?: boolean;
}

export function AddToCartButton({
  productId,
  variantId,
  quantity = 1,
  configuration,
  disabled = false,
  children,
  variant = "primary",
  size = "md",
  className = "",
  showIcon = true,
}: AddToCartButtonProps) {
  const { addItem } = useCart();
  const [isAdding, setIsAdding] = useState(false);
  const [showSuccess, setShowSuccess] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const handleAddToCart = async () => {
    setIsAdding(true);
    setErrorMessage(null);
    try {
      await addItem({
        productId,
        variantId,
        quantity,
        configuration,
      });

      // Show success feedback
      setShowSuccess(true);
      setTimeout(() => setShowSuccess(false), 2000);
    } catch (err) {
      console.error("Failed to add to cart:", err);
      const message = err instanceof Error ? err.message : "Fehler beim Hinzufügen zum Warenkorb";
      setErrorMessage(message);
      setTimeout(() => setErrorMessage(null), 4000);
    } finally {
      setIsAdding(false);
    }
  };

  return (
    <>
      <Button
        variant={variant}
        size={size}
        onClick={handleAddToCart}
        disabled={disabled || isAdding}
        isLoading={isAdding}
        className={className}
      >
        {showSuccess ? (
          <>
            <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M5 13l4 4L19 7"
              />
            </svg>
            Hinzugefügt!
          </>
        ) : (
          <>
            {showIcon && !isAdding && (
              <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"
                />
              </svg>
            )}
            {children || "In den Warenkorb"}
          </>
        )}
      </Button>
      {errorMessage && (
        <Toast
          message={errorMessage}
          type="error"
          onClose={() => setErrorMessage(null)}
        />
      )}
    </>
  );
}
