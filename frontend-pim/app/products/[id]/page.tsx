"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { MainLayout } from "@/components/MainLayout";
import { pimApiClient } from "@/lib/api/client";
import type { Product, Category, PriceScale, BundleComponent } from "@/types/catalog";
import {
  ArrowLeft,
  Save,
  Package,
  Euro,
  FolderTree,
  FileText,
  Grid3x3,
  Boxes,
} from "lucide-react";
import Link from "next/link";

type Tab = "master" | "prices" | "categories" | "attributes" | "variants" | "bundles";

export default function ProductDetailPage() {
  const params = useParams();
  const productId = params.id as string;

  const [product, setProduct] = useState<Product | null>(null);
  const [prices, setPrices] = useState<PriceScale[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [bundleComponents, setBundleComponents] = useState<BundleComponent[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<Tab>("master");

  useEffect(() => {
    const fetchProduct = async () => {
      setLoading(true);
      try {
        const productData = await pimApiClient.getProduct(productId);
        setProduct(productData);

        // Fetch related data based on product type
        const pricesData = await pimApiClient.getProductPrices(productId);
        setPrices(pricesData);

        const categoriesData = await pimApiClient.getProductCategories(productId);
        setCategories(categoriesData);

        // Fetch bundle components if bundle type
        if (productData.productType === "bundle") {
          const bundleData = await pimApiClient.getBundleComponents(productId);
          setBundleComponents(bundleData);
        }
      } catch (error) {
        console.error("Failed to fetch product:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchProduct();
  }, [productId]);

  if (loading || !product) {
    return (
      <MainLayout>
        <div className="flex items-center justify-center py-12">
          <p className="text-gray-500">Lädt...</p>
        </div>
      </MainLayout>
    );
  }

  const tabs = [
    { id: "master" as Tab, label: "Stammdaten", icon: Package },
    { id: "prices" as Tab, label: "Preise", icon: Euro },
    { id: "categories" as Tab, label: "Kategorien", icon: FolderTree },
    { id: "attributes" as Tab, label: "Attribute", icon: FileText },
  ];

  if (product.productType === "variant_parent") {
    tabs.push({ id: "variants" as Tab, label: "Varianten", icon: Grid3x3 });
  }

  if (product.productType === "bundle") {
    tabs.push({ id: "bundles" as Tab, label: "Bundles", icon: Boxes });
  }

  return (
    <MainLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Link
              href="/products"
              className="rounded-lg p-2 hover:bg-gray-100"
            >
              <ArrowLeft className="h-5 w-5" />
            </Link>
            <div>
              <h1 className="text-2xl font-bold text-gray-900">
                {product.name.de || product.name.en || "Unbenannt"}
              </h1>
              <p className="mt-1 text-sm text-gray-500">SKU: {product.sku}</p>
            </div>
          </div>
          <button className="flex items-center gap-2 rounded-lg bg-primary-600 px-4 py-2 text-sm font-semibold text-white hover:bg-primary-700">
            <Save className="h-4 w-4" />
            Speichern
          </button>
        </div>

        {/* Tabs */}
        <div className="border-b border-gray-200">
          <nav className="-mb-px flex gap-8">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center gap-2 border-b-2 px-1 py-4 text-sm font-medium ${
                  activeTab === tab.id
                    ? "border-primary-600 text-primary-600"
                    : "border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700"
                }`}
              >
                <tab.icon className="h-4 w-4" />
                {tab.label}
              </button>
            ))}
          </nav>
        </div>

        {/* Tab Content */}
        <div className="rounded-lg bg-white p-6 shadow-sm ring-1 ring-gray-200">
          {activeTab === "master" && (
            <MasterDataTab product={product} />
          )}
          {activeTab === "prices" && (
            <PricesTab prices={prices} />
          )}
          {activeTab === "categories" && (
            <CategoriesTab categories={categories} />
          )}
          {activeTab === "attributes" && (
            <AttributesTab product={product} />
          )}
          {activeTab === "variants" && (
            <VariantsTab product={product} />
          )}
          {activeTab === "bundles" && (
            <BundlesTab components={bundleComponents} />
          )}
        </div>
      </div>
    </MainLayout>
  );
}

// Master Data Tab
function MasterDataTab({ product }: { product: Product }) {
  return (
    <div className="space-y-6">
      <div className="grid gap-6 md:grid-cols-2">
        <div>
          <label className="block text-sm font-medium text-gray-700">
            Name (DE)
          </label>
          <input
            type="text"
            value={product.name.de || ""}
            readOnly
            className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700">
            Name (EN)
          </label>
          <input
            type="text"
            value={product.name.en || ""}
            readOnly
            className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm"
          />
        </div>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700">
          Beschreibung (DE)
        </label>
        <textarea
          value={product.description?.de || ""}
          readOnly
          rows={4}
          className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm"
        />
      </div>

      <div className="grid gap-6 md:grid-cols-3">
        <div>
          <label className="block text-sm font-medium text-gray-700">SKU</label>
          <input
            type="text"
            value={product.sku}
            readOnly
            className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700">Typ</label>
          <input
            type="text"
            value={product.productType}
            readOnly
            className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700">Status</label>
          <input
            type="text"
            value={product.status}
            readOnly
            className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm"
          />
        </div>
      </div>
    </div>
  );
}

// Prices Tab
function PricesTab({ prices }: { prices: PriceScale[] }) {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">Staffelpreise</h3>
        <button className="rounded-lg bg-primary-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-primary-700">
          Preis hinzufügen
        </button>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="border-b border-gray-200 bg-gray-50">
            <tr>
              <th className="px-4 py-2 text-left text-xs font-medium uppercase text-gray-500">
                Min. Menge
              </th>
              <th className="px-4 py-2 text-left text-xs font-medium uppercase text-gray-500">
                Preis
              </th>
              <th className="px-4 py-2 text-left text-xs font-medium uppercase text-gray-500">
                Währung
              </th>
              <th className="px-4 py-2 text-right text-xs font-medium uppercase text-gray-500">
                Aktionen
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {prices.length === 0 ? (
              <tr>
                <td colSpan={4} className="px-4 py-8 text-center text-sm text-gray-500">
                  Keine Preise definiert
                </td>
              </tr>
            ) : (
              prices.map((price) => (
                <tr key={price.id}>
                  <td className="px-4 py-3 text-sm text-gray-900">{price.minQuantity}</td>
                  <td className="px-4 py-3 text-sm text-gray-900">{price.price.toFixed(2)}</td>
                  <td className="px-4 py-3 text-sm text-gray-900">{price.currency}</td>
                  <td className="px-4 py-3 text-right text-sm">
                    <button className="text-red-600 hover:text-red-800">Löschen</button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}

// Categories Tab
function CategoriesTab({ categories }: { categories: Category[] }) {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">Zugewiesene Kategorien</h3>
        <button className="rounded-lg bg-primary-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-primary-700">
          Kategorie hinzufügen
        </button>
      </div>
      {categories.length === 0 ? (
        <p className="text-sm text-gray-500">Keine Kategorien zugewiesen</p>
      ) : (
        <div className="space-y-2">
          {categories.map((cat) => (
            <div
              key={cat.id}
              className="flex items-center justify-between rounded-lg border border-gray-200 p-3"
            >
              <div>
                <p className="font-medium text-gray-900">
                  {cat.name.de || cat.name.en}
                </p>
                <p className="text-sm text-gray-500">Code: {cat.code}</p>
              </div>
              <button className="text-red-600 hover:text-red-800">Entfernen</button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

// Attributes Tab
function AttributesTab({ product }: { product: Product }) {
  const attributes = product.attributes || {};
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">Produktattribute</h3>
        <button className="rounded-lg bg-primary-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-primary-700">
          Attribut hinzufügen
        </button>
      </div>
      {Object.keys(attributes).length === 0 ? (
        <p className="text-sm text-gray-500">Keine Attribute definiert</p>
      ) : (
        <div className="space-y-2">
          {Object.entries(attributes).map(([key, value]) => (
            <div
              key={key}
              className="flex items-center justify-between rounded-lg border border-gray-200 p-3"
            >
              <div>
                <p className="font-medium text-gray-900">{key}</p>
                <p className="text-sm text-gray-500">{JSON.stringify(value)}</p>
              </div>
              <button className="text-red-600 hover:text-red-800">Löschen</button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

// Variants Tab
function VariantsTab({ product }: { product: Product }) {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">Kind-Varianten</h3>
        <button className="rounded-lg bg-primary-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-primary-700">
          Variante hinzufügen
        </button>
      </div>
      <p className="text-sm text-gray-500">
        Varianten-Verwaltung wird geladen... (Feature in Entwicklung)
      </p>
    </div>
  );
}

// Bundles Tab
function BundlesTab({ components }: { components: BundleComponent[] }) {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">Bundle-Komponenten</h3>
        <button className="rounded-lg bg-primary-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-primary-700">
          Komponente hinzufügen
        </button>
      </div>
      {components.length === 0 ? (
        <p className="text-sm text-gray-500">Keine Komponenten definiert</p>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="border-b border-gray-200 bg-gray-50">
              <tr>
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-gray-500">
                  Produkt
                </th>
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-gray-500">
                  Menge
                </th>
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-gray-500">
                  Min/Max
                </th>
                <th className="px-4 py-2 text-right text-xs font-medium uppercase text-gray-500">
                  Aktionen
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {components.map((comp) => (
                <tr key={comp.id}>
                  <td className="px-4 py-3 text-sm text-gray-900">
                    {comp.product?.name.de || comp.componentProductId}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-900">{comp.quantity}</td>
                  <td className="px-4 py-3 text-sm text-gray-500">
                    {comp.minQuantity || "-"} / {comp.maxQuantity || "-"}
                  </td>
                  <td className="px-4 py-3 text-right text-sm">
                    <button className="text-red-600 hover:text-red-800">Löschen</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
