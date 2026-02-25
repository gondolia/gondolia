"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { apiClient } from "@/lib/api/client";
import { Panel, PanelHeader, PanelBody } from "@/components/ui/Panel";
import { Button } from "@/components/ui/Button";
import type { Order } from "@/types/cart";

export default function OrdersPage() {
  const router = useRouter();
  const [orders, setOrders] = useState<Order[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadOrders();
  }, []);

  const loadOrders = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const ordersData = await apiClient.getOrders();
      setOrders(ordersData);
    } catch (err) {
      const e = err as { message?: string };
      setError(e.message || "Fehler beim Laden der Bestellungen");
      console.error("Failed to load orders:", err);
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
    });
  };

  const getStatusLabel = (status: string) => {
    const labels: Record<string, string> = {
      pending: "Ausstehend",
      confirmed: "Bestätigt",
      processing: "In Bearbeitung",
      shipped: "Versandt",
      delivered: "Geliefert",
      cancelled: "Storniert",
    };
    return labels[status] || status;
  };

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      pending: "bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-200",
      confirmed: "bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-200",
      processing: "bg-purple-100 dark:bg-purple-900/30 text-purple-800 dark:text-purple-200",
      shipped: "bg-indigo-100 dark:bg-indigo-900/30 text-indigo-800 dark:text-indigo-200",
      delivered: "bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-200",
      cancelled: "bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-200",
    };
    return colors[status] || "bg-gray-100 dark:bg-gray-800 text-gray-800 dark:text-gray-200";
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

  return (
    <div className="max-w-7xl mx-auto py-8 px-4">
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Meine Bestellungen</h1>
        <Button variant="outline" onClick={() => router.push("/products")}>
          Weiter einkaufen
        </Button>
      </div>

      {error && (
        <div className="mb-6 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
          <p className="text-red-800 dark:text-red-200">{error}</p>
        </div>
      )}

      {orders.length === 0 ? (
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
                  d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                />
              </svg>
              <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
                Keine Bestellungen vorhanden
              </h2>
              <p className="text-gray-500 dark:text-gray-400 mb-6">
                Sie haben noch keine Bestellung aufgegeben
              </p>
              <Button variant="primary" onClick={() => router.push("/products")}>
                Produkte entdecken
              </Button>
            </div>
          </PanelBody>
        </Panel>
      ) : (
        <div className="space-y-4">
          {orders.map((order) => (
            <Link key={order.id} href={`/orders/${order.id}`}>
              <Panel className="hover:shadow-lg transition-shadow cursor-pointer">
                <PanelBody>
                  <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                    {/* Order Info */}
                    <div className="flex-1">
                      <div className="flex items-center gap-3 mb-2">
                        <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                          {order.orderNumber}
                        </h3>
                        <span
                          className={`px-2 py-1 text-xs font-medium rounded ${getStatusColor(
                            order.status
                          )}`}
                        >
                          {getStatusLabel(order.status)}
                        </span>
                      </div>
                      <p className="text-sm text-gray-500 dark:text-gray-400 mb-2">
                        Bestellt am {formatDate(order.createdAt)}
                      </p>
                      <div className="flex flex-wrap gap-2 text-sm text-gray-600 dark:text-gray-400">
                        <span>{order.itemCount ?? order.items?.length ?? 0} Artikel</span>
                        <span>•</span>
                        <span className="font-semibold text-gray-900 dark:text-white">
                          {formatPrice(order.total, order.currency)}
                        </span>
                      </div>
                    </div>

                    {/* Items Preview */}
                    <div className="flex -space-x-2 overflow-hidden">
                      {order.items && order.items.slice(0, 3).map((item) => (
                        <div
                          key={item.id}
                          className="inline-block w-12 h-12 rounded-full ring-2 ring-white dark:ring-gray-900 bg-gray-100 dark:bg-gray-700 flex items-center justify-center text-xs font-medium text-gray-600 dark:text-gray-400"
                        >
                          {item.quantity}×
                        </div>
                      ))}
                      {!order.items && (order.itemCount ?? 0) > 0 && (
                        <div className="inline-block w-12 h-12 rounded-full ring-2 ring-white dark:ring-gray-900 bg-gray-100 dark:bg-gray-700 flex items-center justify-center text-xs font-medium text-gray-600 dark:text-gray-400">
                          {order.itemCount}×
                        </div>
                      )}
                      {order.items && order.items.length > 3 && (
                        <div className="inline-block w-12 h-12 rounded-full ring-2 ring-white dark:ring-gray-900 bg-gray-200 dark:bg-gray-600 flex items-center justify-center text-xs font-medium text-gray-600 dark:text-gray-400">
                          +{order.items.length - 3}
                        </div>
                      )}
                    </div>

                    {/* Arrow */}
                    <div className="hidden sm:block">
                      <svg
                        className="w-6 h-6 text-gray-400"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M9 5l7 7-7 7"
                        />
                      </svg>
                    </div>
                  </div>
                </PanelBody>
              </Panel>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
