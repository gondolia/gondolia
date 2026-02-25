"use client";

import { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { MainLayout } from "@/components/MainLayout";
import { pimApiClient } from "@/lib/api/client";
import type { ProductType, ProductStatus } from "@/types/catalog";
import { ArrowLeft, ArrowRight, Check } from "lucide-react";
import Link from "next/link";

type Step = 1 | 2 | 3;

export default function NewProductPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const parentId = searchParams.get("parent");

  const [step, setStep] = useState<Step>(1);
  const [saving, setSaving] = useState(false);
  const [formData, setFormData] = useState({
    sku: "",
    productType: (parentId ? "variant" : "simple") as ProductType,
    status: "draft" as ProductStatus,
    nameDe: "",
    nameEn: "",
    descriptionDe: "",
    descriptionEn: "",
    manufacturer: "",
    manufacturerPartNumber: "",
  });

  const handleCreate = async () => {
    if (!formData.sku || !formData.nameDe) {
      alert("Bitte SKU und Name (DE) ausfüllen");
      return;
    }

    setSaving(true);
    try {
      const product = await pimApiClient.createProduct({
        sku: formData.sku,
        productType: formData.productType,
        status: formData.status,
        name: {
          de: formData.nameDe,
          en: formData.nameEn,
        },
        description: {
          de: formData.descriptionDe,
          en: formData.descriptionEn,
        },
        manufacturer: formData.manufacturer || undefined,
        manufacturerPartNumber: formData.manufacturerPartNumber || undefined,
        parentProductId: parentId || undefined,
      });

      router.push(`/products/${product.id}`);
    } catch (error: any) {
      alert(`Fehler beim Erstellen: ${error.message || "Unbekannter Fehler"}`);
      setSaving(false);
    }
  };

  const canProceedStep1 = formData.productType !== "";
  const canProceedStep2 = formData.sku !== "" && formData.nameDe !== "";

  return (
    <MainLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex items-center gap-4">
          <Link href="/products" className="rounded-lg p-2 hover:bg-gray-100">
            <ArrowLeft className="h-5 w-5" />
          </Link>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Neues Produkt erstellen</h1>
            <p className="mt-1 text-sm text-gray-500">
              {parentId ? `Variante für Produkt ${parentId}` : "Schritt für Schritt"}
            </p>
          </div>
        </div>

        {/* Progress Stepper */}
        <div className="flex items-center justify-center space-x-4">
          {[1, 2, 3].map((s) => (
            <div key={s} className="flex items-center">
              <div
                className={`flex h-10 w-10 items-center justify-center rounded-full ${
                  step === s
                    ? "bg-primary-600 text-white"
                    : step > s
                    ? "bg-green-600 text-white"
                    : "bg-gray-200 text-gray-500"
                }`}
              >
                {step > s ? <Check className="h-5 w-5" /> : s}
              </div>
              {s < 3 && (
                <div
                  className={`mx-2 h-1 w-16 ${
                    step > s ? "bg-green-600" : "bg-gray-200"
                  }`}
                />
              )}
            </div>
          ))}
        </div>

        {/* Step Content */}
        <div className="rounded-lg bg-white p-6 shadow-sm ring-1 ring-gray-200">
          {step === 1 && (
            <div className="space-y-6">
              <h2 className="text-lg font-semibold text-gray-900">Schritt 1: Produkttyp wählen</h2>
              <div className="grid gap-4 md:grid-cols-2">
                {["simple", "variant_parent", "variant", "bundle", "parametric"].map((type) => {
                  const isDisabled = parentId && type !== "variant";
                  return (
                    <button
                      key={type}
                      onClick={() => !isDisabled && setFormData({ ...formData, productType: type as ProductType })}
                      disabled={isDisabled}
                      className={`rounded-lg border-2 p-4 text-left transition-colors ${
                        formData.productType === type
                          ? "border-primary-600 bg-primary-50"
                          : isDisabled
                          ? "border-gray-200 bg-gray-100 cursor-not-allowed"
                          : "border-gray-200 hover:border-gray-300"
                      }`}
                    >
                      <p className="font-semibold text-gray-900">{type.replace("_", " ")}</p>
                      <p className="mt-1 text-sm text-gray-500">
                        {type === "simple" && "Einfaches Produkt ohne Varianten"}
                        {type === "variant_parent" && "Produkt mit Varianten (z.B. Farbe, Größe)"}
                        {type === "variant" && "Kind-Variante eines Varianten-Produkts"}
                        {type === "bundle" && "Bundle aus mehreren Produkten"}
                        {type === "parametric" && "Parametrisches Produkt"}
                      </p>
                    </button>
                  );
                })}
              </div>
            </div>
          )}

          {step === 2 && (
            <div className="space-y-6">
              <h2 className="text-lg font-semibold text-gray-900">Schritt 2: Stammdaten eingeben</h2>
              <div className="grid gap-6 md:grid-cols-2">
                <div>
                  <label className="block text-sm font-medium text-gray-700">
                    SKU <span className="text-red-600">*</span>
                  </label>
                  <input
                    type="text"
                    value={formData.sku}
                    onChange={(e) => setFormData({ ...formData, sku: e.target.value })}
                    className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
                    placeholder="z.B. PROD-001"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">Status</label>
                  <select
                    value={formData.status}
                    onChange={(e) => setFormData({ ...formData, status: e.target.value as ProductStatus })}
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
                    Name (DE) <span className="text-red-600">*</span>
                  </label>
                  <input
                    type="text"
                    value={formData.nameDe}
                    onChange={(e) => setFormData({ ...formData, nameDe: e.target.value })}
                    className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">Name (EN)</label>
                  <input
                    type="text"
                    value={formData.nameEn}
                    onChange={(e) => setFormData({ ...formData, nameEn: e.target.value })}
                    className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
                  />
                </div>
              </div>

              <div className="grid gap-6 md:grid-cols-2">
                <div>
                  <label className="block text-sm font-medium text-gray-700">Beschreibung (DE)</label>
                  <textarea
                    value={formData.descriptionDe}
                    onChange={(e) => setFormData({ ...formData, descriptionDe: e.target.value })}
                    rows={4}
                    className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">Beschreibung (EN)</label>
                  <textarea
                    value={formData.descriptionEn}
                    onChange={(e) => setFormData({ ...formData, descriptionEn: e.target.value })}
                    rows={4}
                    className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
                  />
                </div>
              </div>

              <div className="grid gap-6 md:grid-cols-2">
                <div>
                  <label className="block text-sm font-medium text-gray-700">Hersteller</label>
                  <input
                    type="text"
                    value={formData.manufacturer}
                    onChange={(e) => setFormData({ ...formData, manufacturer: e.target.value })}
                    className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700">Hersteller-Teilenummer</label>
                  <input
                    type="text"
                    value={formData.manufacturerPartNumber}
                    onChange={(e) => setFormData({ ...formData, manufacturerPartNumber: e.target.value })}
                    className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
                  />
                </div>
              </div>
            </div>
          )}

          {step === 3 && (
            <div className="space-y-6">
              <h2 className="text-lg font-semibold text-gray-900">Schritt 3: Bestätigen & Speichern</h2>
              <div className="rounded-lg bg-gray-50 p-4 space-y-3">
                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <span className="font-medium text-gray-700">SKU:</span>
                    <p className="text-gray-900">{formData.sku}</p>
                  </div>
                  <div>
                    <span className="font-medium text-gray-700">Typ:</span>
                    <p className="text-gray-900">{formData.productType}</p>
                  </div>
                  <div>
                    <span className="font-medium text-gray-700">Status:</span>
                    <p className="text-gray-900">{formData.status}</p>
                  </div>
                  <div>
                    <span className="font-medium text-gray-700">Name (DE):</span>
                    <p className="text-gray-900">{formData.nameDe}</p>
                  </div>
                  {formData.nameEn && (
                    <div>
                      <span className="font-medium text-gray-700">Name (EN):</span>
                      <p className="text-gray-900">{formData.nameEn}</p>
                    </div>
                  )}
                  {formData.manufacturer && (
                    <div>
                      <span className="font-medium text-gray-700">Hersteller:</span>
                      <p className="text-gray-900">{formData.manufacturer}</p>
                    </div>
                  )}
                </div>
              </div>
              <p className="text-sm text-gray-500">
                Nach dem Speichern können Sie Preise, Kategorien und Attribute hinzufügen.
              </p>
            </div>
          )}

          {/* Navigation Buttons */}
          <div className="mt-8 flex items-center justify-between border-t pt-6">
            {step > 1 ? (
              <button
                onClick={() => setStep((s) => (s - 1) as Step)}
                className="flex items-center gap-2 rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
              >
                <ArrowLeft className="h-4 w-4" />
                Zurück
              </button>
            ) : (
              <div />
            )}

            {step < 3 ? (
              <button
                onClick={() => setStep((s) => (s + 1) as Step)}
                disabled={(step === 1 && !canProceedStep1) || (step === 2 && !canProceedStep2)}
                className="flex items-center gap-2 rounded-lg bg-primary-600 px-4 py-2 text-sm font-semibold text-white hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Weiter
                <ArrowRight className="h-4 w-4" />
              </button>
            ) : (
              <button
                onClick={handleCreate}
                disabled={saving}
                className="flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-semibold text-white hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {saving ? "Speichert..." : "Produkt erstellen"}
                <Check className="h-4 w-4" />
              </button>
            )}
          </div>
        </div>
      </div>
    </MainLayout>
  );
}
