"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { MainLayout } from "@/components/MainLayout";
import { ArrowLeft, Save } from "lucide-react";
import Link from "next/link";
import type { ProductType } from "@/types/catalog";

export default function NewProductPage() {
  const router = useRouter();
  const [step, setStep] = useState(1);
  const [productType, setProductType] = useState<ProductType>("simple");
  const [formData, setFormData] = useState({
    sku: "",
    nameDe: "",
    nameEn: "",
    descriptionDe: "",
    descriptionEn: "",
    status: "draft" as const,
    basePrice: "",
    currency: "EUR",
  });

  const handleTypeSelect = (type: ProductType) => {
    setProductType(type);
    setStep(2);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    // TODO: Implement product creation when API is ready
    alert("Produkt-Erstellung wird implementiert, sobald die Write-API verfügbar ist.");
  };

  const productTypes = [
    {
      type: "simple" as ProductType,
      title: "Simple Product",
      description: "Ein einfaches Produkt ohne Varianten oder Konfiguration",
    },
    {
      type: "variant_parent" as ProductType,
      title: "Variant Parent",
      description: "Produkt mit Varianten (z.B. verschiedene Farben/Größen)",
    },
    {
      type: "bundle" as ProductType,
      title: "Bundle Product",
      description: "Produkt-Bundle aus mehreren Komponenten",
    },
    {
      type: "parametric" as ProductType,
      title: "Parametric Product",
      description: "Produkt mit konfigurierbaren Parametern",
    },
  ];

  return (
    <MainLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Link href="/products" className="rounded-lg p-2 hover:bg-gray-100">
              <ArrowLeft className="h-5 w-5" />
            </Link>
            <div>
              <h1 className="text-2xl font-bold text-gray-900">
                Neues Produkt erstellen
              </h1>
              <p className="mt-1 text-sm text-gray-500">
                Schritt {step} von 3
              </p>
            </div>
          </div>
        </div>

        {/* Step 1: Type Selection */}
        {step === 1 && (
          <div className="space-y-4">
            <h2 className="text-lg font-semibold text-gray-900">
              Produkttyp auswählen
            </h2>
            <div className="grid gap-4 md:grid-cols-2">
              {productTypes.map((pt) => (
                <button
                  key={pt.type}
                  onClick={() => handleTypeSelect(pt.type)}
                  className="rounded-lg border-2 border-gray-200 p-6 text-left transition-all hover:border-primary-500 hover:bg-primary-50"
                >
                  <h3 className="font-semibold text-gray-900">{pt.title}</h3>
                  <p className="mt-2 text-sm text-gray-600">{pt.description}</p>
                </button>
              ))}
            </div>
          </div>
        )}

        {/* Step 2: Master Data */}
        {step === 2 && (
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="rounded-lg bg-white p-6 shadow-sm ring-1 ring-gray-200">
              <h2 className="mb-6 text-lg font-semibold text-gray-900">
                Stammdaten
              </h2>

              <div className="space-y-6">
                <div className="grid gap-6 md:grid-cols-2">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      SKU *
                    </label>
                    <input
                      type="text"
                      required
                      value={formData.sku}
                      onChange={(e) =>
                        setFormData({ ...formData, sku: e.target.value })
                      }
                      className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
                      placeholder="z.B. PROD-001"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Status *
                    </label>
                    <select
                      required
                      value={formData.status}
                      onChange={(e) =>
                        setFormData({
                          ...formData,
                          status: e.target.value as "draft" | "active" | "inactive",
                        })
                      }
                      className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
                    >
                      <option value="draft">Draft</option>
                      <option value="active">Active</option>
                      <option value="inactive">Inactive</option>
                    </select>
                  </div>
                </div>

                <div className="grid gap-6 md:grid-cols-2">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Name (DE) *
                    </label>
                    <input
                      type="text"
                      required
                      value={formData.nameDe}
                      onChange={(e) =>
                        setFormData({ ...formData, nameDe: e.target.value })
                      }
                      className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Name (EN)
                    </label>
                    <input
                      type="text"
                      value={formData.nameEn}
                      onChange={(e) =>
                        setFormData({ ...formData, nameEn: e.target.value })
                      }
                      className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
                    />
                  </div>
                </div>

                <div className="grid gap-6 md:grid-cols-2">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Beschreibung (DE)
                    </label>
                    <textarea
                      rows={4}
                      value={formData.descriptionDe}
                      onChange={(e) =>
                        setFormData({
                          ...formData,
                          descriptionDe: e.target.value,
                        })
                      }
                      className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Beschreibung (EN)
                    </label>
                    <textarea
                      rows={4}
                      value={formData.descriptionEn}
                      onChange={(e) =>
                        setFormData({
                          ...formData,
                          descriptionEn: e.target.value,
                        })
                      }
                      className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
                    />
                  </div>
                </div>
              </div>
            </div>

            <div className="rounded-lg bg-white p-6 shadow-sm ring-1 ring-gray-200">
              <h2 className="mb-6 text-lg font-semibold text-gray-900">
                Basispreis
              </h2>

              <div className="grid gap-6 md:grid-cols-2">
                <div>
                  <label className="block text-sm font-medium text-gray-700">
                    Preis *
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    required
                    value={formData.basePrice}
                    onChange={(e) =>
                      setFormData({ ...formData, basePrice: e.target.value })
                    }
                    className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
                    placeholder="0.00"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">
                    Währung *
                  </label>
                  <select
                    required
                    value={formData.currency}
                    onChange={(e) =>
                      setFormData({ ...formData, currency: e.target.value })
                    }
                    className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
                  >
                    <option value="EUR">EUR</option>
                    <option value="USD">USD</option>
                    <option value="GBP">GBP</option>
                  </select>
                </div>
              </div>
            </div>

            <div className="flex justify-between">
              <button
                type="button"
                onClick={() => setStep(1)}
                className="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
              >
                Zurück
              </button>
              <button
                type="submit"
                className="flex items-center gap-2 rounded-lg bg-primary-600 px-4 py-2 text-sm font-semibold text-white hover:bg-primary-700"
              >
                <Save className="h-4 w-4" />
                Produkt erstellen
              </button>
            </div>
          </form>
        )}
      </div>
    </MainLayout>
  );
}
