"use client";

import React, { createContext, useContext, useState, useEffect, useCallback } from "react";
import type { Cart, CartItem, AddToCartRequest } from "@/types/cart";
import { apiClient } from "@/lib/api/client";

interface CartContextValue {
  cart: Cart | null;
  isLoading: boolean;
  error: string | null;
  itemCount: number;
  addItem: (request: AddToCartRequest) => Promise<void>;
  updateQuantity: (itemId: string, quantity: number) => Promise<void>;
  removeItem: (itemId: string) => Promise<void>;
  clearCart: () => Promise<void>;
  refreshCart: () => Promise<void>;
}

const CartContext = createContext<CartContextValue | undefined>(undefined);

export function CartProvider({ children }: { children: React.ReactNode }) {
  const [cart, setCart] = useState<Cart | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isMounted, setIsMounted] = useState(false);

  // Calculate item count
  const itemCount = cart?.items.reduce((sum, item) => sum + item.quantity, 0) || 0;

  // Load cart on mount (only client-side)
  const refreshCart = useCallback(async () => {
    // Only refresh cart on client-side
    if (typeof window === 'undefined') return;

    setIsLoading(true);
    setError(null);
    try {
      const loadedCart = await apiClient.getCart();
      setCart(loadedCart);
    } catch (err) {
      const e = err as { message?: string };
      setError(e.message || "Fehler beim Laden des Warenkorbs");
      console.error("Failed to load cart:", err);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    setIsMounted(true);
    refreshCart();
  }, [refreshCart]);

  // Add item to cart
  const addItem = useCallback(async (request: AddToCartRequest) => {
    setIsLoading(true);
    setError(null);
    try {
      const updatedCart = await apiClient.addToCart(request);
      setCart(updatedCart);
    } catch (err) {
      const e = err as { message?: string };
      setError(e.message || "Fehler beim Hinzufügen zum Warenkorb");
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Update item quantity
  const updateQuantity = useCallback(async (itemId: string, quantity: number) => {
    setIsLoading(true);
    setError(null);
    try {
      const updatedCart = await apiClient.updateCartItem(itemId, { quantity });
      setCart(updatedCart);
    } catch (err) {
      const e = err as { message?: string };
      setError(e.message || "Fehler beim Aktualisieren der Menge");
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Remove item from cart
  const removeItem = useCallback(async (itemId: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const updatedCart = await apiClient.removeCartItem(itemId);
      setCart(updatedCart);
    } catch (err) {
      const e = err as { message?: string };
      setError(e.message || "Fehler beim Entfernen aus dem Warenkorb");
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Clear entire cart
  const clearCart = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      await apiClient.clearCart();
      // Clear local state immediately - cart items are deleted on server
      // Set cart to empty state rather than fetching from server
      setCart(prevCart => prevCart ? {
        ...prevCart,
        items: [],
        subtotal: 0,
      } : null);
    } catch (err) {
      const e = err as { message?: string };
      setError(e.message || "Fehler beim Leeren des Warenkorbs");
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const value: CartContextValue = {
    cart,
    isLoading,
    error,
    itemCount,
    addItem,
    updateQuantity,
    removeItem,
    clearCart,
    refreshCart,
  };

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
}

export function useCart() {
  const context = useContext(CartContext);
  if (context === undefined) {
    throw new Error("useCart must be used within a CartProvider");
  }
  return context;
}
