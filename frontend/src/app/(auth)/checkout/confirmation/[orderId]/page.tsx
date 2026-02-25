"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { apiClient } from "@/lib/api/client";
import { Panel, PanelHeader, PanelBody } from "@/components/ui/Panel";
import { Button } from "@/components/ui/Button";
import type { Order } from "@/types/cart";

export default function CheckoutConfirmationPage() {
  const params = useParams();
  const router = useRouter();
  const orderId = params.orderId as string;

  const [order, setOrder] = useState<Order | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadOrder();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [orderId]);

  const loadOrder = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const orderData = await apiClient.getOrder(orderId);
      setOrder(orderData);
    } catch (err) {
      const e = err as { message?: string };
      setError(e.message || "Fehler beim Laden der Bestellung");
      console.error("Failed to load order:", err);
    } finally {
      setIsLoading(false);
    }
  };

  const formatPrice = (price: number, currency: string) => {
    return new Intl.NumberFormat("de-CH", {
      style: "currency",
      currency: currency || "CHF",
    }).format(price);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("de-CH", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <svg className="w-8 h-8 animate-spin text-primary-600" fill="none" viewBox="0 0 24 24">
          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
        </svg>
      </div>
    );
  }

  if (error || !order) {
    return (
      <div className="max-w-4xl mx-auto py-12">
        <Panel>
          <PanelBody>
            <div className="text-center py-12">
              <svg
                className="w-16 h-16 text-red-300 dark:text-red-600 mx-auto mb-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
                Bestellung nicht gefunden
              </h2>
              <p className="text-gray-500 dark:text-gray-400 mb-6">{error}</p>
              <Button variant="primary" onClick={() => router.push("/orders")}>
                Zu meinen Bestellungen
              </Button>
            </div>
          </PanelBody>
        </Panel>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto py-8 px-4">
      {/* Success Message */}
      <div className="mb-8 p-6 bg-green-50 dark:bg-green-900/20 border-2 border-green-200 dark:border-green-800 rounded-lg">
        <div className="flex items-start gap-4">
          <div className="flex-shrink-0">
            <svg
              className="w-12 h-12 text-green-600 dark:text-green-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          </div>
          <div className="flex-1">
            <h1 className="text-2xl font-bold text-green-900 dark:text-green-100 mb-2">
              Bestellung erfolgreich!
            </h1>
            <p className="text-green-800 dark:text-green-200 mb-4">
              Vielen Dank für Ihre Bestellung. Sie erhalten in Kürze eine Bestätigungsmail.
            </p>
            <div className="flex flex-wrap gap-2 text-sm">
              <span className="text-green-700 dark:text-green-300">
                Bestellnummer: <strong>{order.orderNumber}</strong>
              </span>
              <span className="text-green-600 dark:text-green-400">•</span>
              <span className="text-green-700 dark:text-green-300">
                {formatDate(order.createdAt)}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Order Details */}
      <div className="space-y-6">
        {/* Items */}
        <Panel>
          <PanelHeader>
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
              Bestellte Artikel
            </h2>
          </PanelHeader>
          <PanelBody>
            <div className="space-y-4">
              {order.items.map((item) => (
                <div
                  key={item.id}
                  className="flex gap-4 pb-4 border-b border-gray-200 dark:border-gray-700 last:border-0"
                >
                  <div className="flex-1">
                    <h4 className="font-semibold text-gray-900 dark:text-white">
                      {item.productName}
                    </h4>
                    <p className="text-sm text-gray-500 dark:text-gray-400 font-mono mt-1">
                      SKU: {item.sku}
                    </p>
                    <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                      {item.quantity}× {formatPrice(item.unitPrice, item.currency)}
                    </p>
                  </div>
                  <div className="text-right">
                    <p className="font-bold text-gray-900 dark:text-white">
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
                  {formatPrice(order.subtotal, order.currency)}
                </span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-600 dark:text-gray-400">MwSt.</span>
                <span className="font-medium text-gray-900 dark:text-white">
                  {formatPrice(order.taxAmount, order.currency)}
                </span>
              </div>
              <div className="flex justify-between text-xl font-bold pt-2 border-t border-gray-200 dark:border-gray-700">
                <span className="text-gray-900 dark:text-white">Gesamt</span>
                <span className="text-primary-600 dark:text-primary-400">
                  {formatPrice(order.total, order.currency)}
                </span>
              </div>
            </div>
          </PanelBody>
        </Panel>

        {/* Addresses */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Panel>
            <PanelHeader>
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                Lieferadresse
              </h2>
            </PanelHeader>
            <PanelBody>
              <div className="text-sm text-gray-700 dark:text-gray-300 space-y-1">
                {order.shippingAddress.company && <p className="font-semibold">{order.shippingAddress.company}</p>}
                <p>
                  {order.shippingAddress.firstName} {order.shippingAddress.lastName}
                </p>
                <p>{order.shippingAddress.street}</p>
                <p>
                  {order.shippingAddress.postalCode} {order.shippingAddress.city}
                </p>
                <p>{order.shippingAddress.country}</p>
                {order.shippingAddress.phone && <p>{order.shippingAddress.phone}</p>}
              </div>
            </PanelBody>
          </Panel>

          <Panel>
            <PanelHeader>
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                Rechnungsadresse
              </h2>
            </PanelHeader>
            <PanelBody>
              <div className="text-sm text-gray-700 dark:text-gray-300 space-y-1">
                {order.billingAddress.company && <p className="font-semibold">{order.billingAddress.company}</p>}
                <p>
                  {order.billingAddress.firstName} {order.billingAddress.lastName}
                </p>
                <p>{order.billingAddress.street}</p>
                <p>
                  {order.billingAddress.postalCode} {order.billingAddress.city}
                </p>
                <p>{order.billingAddress.country}</p>
                {order.billingAddress.phone && <p>{order.billingAddress.phone}</p>}
              </div>
            </PanelBody>
          </Panel>
        </div>

        {/* Notes */}
        {order.notes && (
          <Panel>
            <PanelHeader>
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                Anmerkungen
              </h2>
            </PanelHeader>
            <PanelBody>
              <p className="text-sm text-gray-700 dark:text-gray-300">{order.notes}</p>
            </PanelBody>
          </Panel>
        )}

        {/* Actions */}
        <div className="flex flex-col sm:flex-row gap-4">
          <Button variant="primary" onClick={() => router.push("/orders")}>
            Zu meinen Bestellungen
          </Button>
          <Button variant="outline" onClick={() => router.push("/products")}>
            Weiter einkaufen
          </Button>
        </div>
      </div>
    </div>
  );
}
