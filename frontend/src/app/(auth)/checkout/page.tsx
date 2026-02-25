"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useCart } from "@/context/CartContext";
import { apiClient } from "@/lib/api/client";
import { Panel, PanelHeader, PanelBody } from "@/components/ui/Panel";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import type { Address } from "@/types/cart";

export default function CheckoutPage() {
  const router = useRouter();
  const { cart, isLoading: cartLoading } = useCart();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [shippingAddress, setShippingAddress] = useState<Address>({
    company: "",
    firstName: "",
    lastName: "",
    street: "",
    city: "",
    postalCode: "",
    country: "CH",
    phone: "",
  });

  const [billingAddress, setBillingAddress] = useState<Address>({
    company: "",
    firstName: "",
    lastName: "",
    street: "",
    city: "",
    postalCode: "",
    country: "CH",
    phone: "",
  });

  const [useSameAddress, setUseSameAddress] = useState(true);
  const [notes, setNotes] = useState("");

  const formatPrice = (price: number, currency: string) => {
    return new Intl.NumberFormat("de-CH", {
      style: "currency",
      currency: currency || "CHF",
    }).format(price);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError(null);

    try {
      const result = await apiClient.checkout({
        shippingAddress,
        billingAddress: useSameAddress ? shippingAddress : billingAddress,
        notes: notes || undefined,
      });

      // Validate result before navigating
      if (!result || !result.order || !result.order.id) {
        throw new Error("Ungültige Antwort vom Server");
      }

      router.push(`/checkout/confirmation/${result.order.id}`);
    } catch (err) {
      const e = err as { message?: string; code?: string };
      // Provide more helpful error messages based on error type
      let errorMessage = "Fehler beim Abschließen der Bestellung";
      if (e.message) {
        if (e.message.includes("cart is empty") || e.message.includes("Warenkorb ist leer")) {
          errorMessage = "Ihr Warenkorb ist leer. Bitte fügen Sie Produkte hinzu.";
        } else if (e.message.includes("authentication") || e.message.includes("unauthorized")) {
          errorMessage = "Bitte melden Sie sich erneut an.";
        } else {
          errorMessage = e.message;
        }
      }
      setError(errorMessage);
      console.error("Checkout failed:", err);
    } finally {
      setIsSubmitting(false);
    }
  };

  if (cartLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <svg className="w-8 h-8 animate-spin text-primary-600" fill="none" viewBox="0 0 24 24">
          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
        </svg>
      </div>
    );
  }

  if (!cart || cart.items.length === 0) {
    return (
      <div className="max-w-4xl mx-auto py-12">
        <Panel>
          <PanelBody>
            <div className="text-center py-12">
              <svg
                className="w-16 h-16 text-gray-300 dark:text-gray-600 mx-auto mb-4"
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
              <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
                Ihr Warenkorb ist leer
              </h2>
              <p className="text-gray-500 dark:text-gray-400 mb-6">
                Fügen Sie Produkte hinzu, bevor Sie zur Kasse gehen
              </p>
              <Button variant="primary" onClick={() => router.push("/products")}>
                Weiter einkaufen
              </Button>
            </div>
          </PanelBody>
        </Panel>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto py-8 px-4">
      <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-8">
        Zur Kasse
      </h1>

      {error && (
        <div className="mb-6 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
          <p className="text-red-800 dark:text-red-200">{error}</p>
        </div>
      )}

      <form onSubmit={handleSubmit}>
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Left Column - Address Forms */}
          <div className="lg:col-span-2 space-y-6">
            {/* Shipping Address */}
            <Panel>
              <PanelHeader>
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                  Lieferadresse
                </h2>
              </PanelHeader>
              <PanelBody>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <Input
                    label="Firma"
                    value={shippingAddress.company}
                    onChange={(e) => setShippingAddress({ ...shippingAddress, company: e.target.value })}
                  />
                  <div className="sm:col-span-2" />
                  <Input
                    label="Vorname"
                    required
                    value={shippingAddress.firstName}
                    onChange={(e) => setShippingAddress({ ...shippingAddress, firstName: e.target.value })}
                  />
                  <Input
                    label="Nachname"
                    required
                    value={shippingAddress.lastName}
                    onChange={(e) => setShippingAddress({ ...shippingAddress, lastName: e.target.value })}
                  />
                  <Input
                    label="Straße und Hausnummer"
                    required
                    className="sm:col-span-2"
                    value={shippingAddress.street}
                    onChange={(e) => setShippingAddress({ ...shippingAddress, street: e.target.value })}
                  />
                  <Input
                    label="Postleitzahl"
                    required
                    value={shippingAddress.postalCode}
                    onChange={(e) => setShippingAddress({ ...shippingAddress, postalCode: e.target.value })}
                  />
                  <Input
                    label="Stadt"
                    required
                    value={shippingAddress.city}
                    onChange={(e) => setShippingAddress({ ...shippingAddress, city: e.target.value })}
                  />
                  <Input
                    label="Land"
                    required
                    value={shippingAddress.country}
                    onChange={(e) => setShippingAddress({ ...shippingAddress, country: e.target.value })}
                  />
                  <Input
                    label="Telefon"
                    value={shippingAddress.phone}
                    onChange={(e) => setShippingAddress({ ...shippingAddress, phone: e.target.value })}
                  />
                </div>
              </PanelBody>
            </Panel>

            {/* Billing Address */}
            <Panel>
              <PanelHeader>
                <div className="flex items-center justify-between">
                  <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                    Rechnungsadresse
                  </h2>
                  <label className="flex items-center gap-2 text-sm">
                    <input
                      type="checkbox"
                      checked={useSameAddress}
                      onChange={(e) => setUseSameAddress(e.target.checked)}
                      className="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                    />
                    <span className="text-gray-700 dark:text-gray-300">
                      Gleich wie Lieferadresse
                    </span>
                  </label>
                </div>
              </PanelHeader>
              {!useSameAddress && (
                <PanelBody>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <Input
                      label="Firma"
                      value={billingAddress.company}
                      onChange={(e) => setBillingAddress({ ...billingAddress, company: e.target.value })}
                    />
                    <div className="sm:col-span-2" />
                    <Input
                      label="Vorname"
                      required
                      value={billingAddress.firstName}
                      onChange={(e) => setBillingAddress({ ...billingAddress, firstName: e.target.value })}
                    />
                    <Input
                      label="Nachname"
                      required
                      value={billingAddress.lastName}
                      onChange={(e) => setBillingAddress({ ...billingAddress, lastName: e.target.value })}
                    />
                    <Input
                      label="Straße und Hausnummer"
                      required
                      className="sm:col-span-2"
                      value={billingAddress.street}
                      onChange={(e) => setBillingAddress({ ...billingAddress, street: e.target.value })}
                    />
                    <Input
                      label="Postleitzahl"
                      required
                      value={billingAddress.postalCode}
                      onChange={(e) => setBillingAddress({ ...billingAddress, postalCode: e.target.value })}
                    />
                    <Input
                      label="Stadt"
                      required
                      value={billingAddress.city}
                      onChange={(e) => setBillingAddress({ ...billingAddress, city: e.target.value })}
                    />
                    <Input
                      label="Land"
                      required
                      value={billingAddress.country}
                      onChange={(e) => setBillingAddress({ ...billingAddress, country: e.target.value })}
                    />
                    <Input
                      label="Telefon"
                      value={billingAddress.phone}
                      onChange={(e) => setBillingAddress({ ...billingAddress, phone: e.target.value })}
                    />
                  </div>
                </PanelBody>
              )}
            </Panel>

            {/* Notes */}
            <Panel>
              <PanelHeader>
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                  Anmerkungen (optional)
                </h2>
              </PanelHeader>
              <PanelBody>
                <textarea
                  value={notes}
                  onChange={(e) => setNotes(e.target.value)}
                  rows={3}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                  placeholder="Zusätzliche Hinweise zur Bestellung..."
                />
              </PanelBody>
            </Panel>
          </div>

          {/* Right Column - Order Summary */}
          <div className="lg:col-span-1">
            <Panel>
              <PanelHeader>
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                  Bestellübersicht
                </h2>
              </PanelHeader>
              <PanelBody>
                <div className="space-y-3">
                  {cart.items.map((item) => (
                    <div key={item.id} className="flex gap-3 pb-3 border-b border-gray-200 dark:border-gray-700 last:border-0">
                      <div className="w-16 h-16 flex-shrink-0 bg-gray-100 dark:bg-gray-700 rounded overflow-hidden">
                        {item.imageUrl ? (
                          <img
                            src={item.imageUrl}
                            alt={item.productName}
                            className="w-full h-full object-cover"
                          />
                        ) : (
                          <div className="w-full h-full flex items-center justify-center text-gray-400">
                            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
                      <div className="flex-1 min-w-0">
                        <h4 className="text-sm font-medium text-gray-900 dark:text-white line-clamp-2">
                          {item.productName}
                        </h4>
                        <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                          {item.quantity}× {formatPrice(item.unitPrice, item.currency)}
                        </p>
                        <p className="text-sm font-semibold text-gray-900 dark:text-white mt-1">
                          {formatPrice(item.totalPrice, item.currency)}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>

                <div className="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700 space-y-2">
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-600 dark:text-gray-400">Zwischensumme</span>
                    <span className="font-medium text-gray-900 dark:text-white">
                      {formatPrice(cart.subtotal, cart.currency)}
                    </span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-600 dark:text-gray-400">MwSt.</span>
                    <span className="font-medium text-gray-900 dark:text-white">
                      wird berechnet
                    </span>
                  </div>
                  <div className="flex justify-between text-lg font-bold pt-2 border-t border-gray-200 dark:border-gray-700">
                    <span className="text-gray-900 dark:text-white">Gesamt</span>
                    <span className="text-primary-600 dark:text-primary-400">
                      {formatPrice(cart.subtotal, cart.currency)}
                    </span>
                  </div>
                </div>

                <Button
                  type="submit"
                  variant="primary"
                  size="lg"
                  className="w-full mt-6"
                  isLoading={isSubmitting}
                  disabled={isSubmitting}
                >
                  Jetzt kostenpflichtig bestellen
                </Button>
              </PanelBody>
            </Panel>
          </div>
        </div>
      </form>
    </div>
  );
}
