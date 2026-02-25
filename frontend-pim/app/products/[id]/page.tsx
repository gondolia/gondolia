"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
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
  Plus,
  Trash2,
  Edit,
  X,
} from "lucide-react";
import Link from "next/link";

type Tab = "master" | "prices" | "categories" | "attributes" | "variants" | "bundles";

export default function ProductDetailPage() {
  const params = useParams();
  const router = useRouter();
  const productId = params.id as string;

  const [product, setProduct] = useState<Product | null>(null);
  const [editedProduct, setEditedProduct] = useState<Partial<Product>>({});
  const [prices, setPrices] = useState<PriceScale[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [bundleComponents, setBundleComponents] = useState<BundleComponent[]>([]);
  const [variants, setVariants] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [activeTab, setActiveTab] = useState<Tab>("master");
  const [toast, setToast] = useState<{ type: "success" | "error"; message: string } | null>(null);

  const showToast = (type: "success" | "error", message: string) => {
    setToast({ type, message });
    setTimeout(() => setToast(null), 3000);
  };

  const fetchProduct = async () => {
    setLoading(true);
    try {
      const productData = await pimApiClient.getProduct(productId);
      setProduct(productData);
      setEditedProduct({});

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

      // Fetch variants if variant_parent
      if (productData.productType === "variant_parent") {
        const variantsData = await pimApiClient.getVariants(productId);
        setVariants(variantsData);
      }
    } catch (error) {
      console.error("Failed to fetch product:", error);
      showToast("error", "Fehler beim Laden des Produkts");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProduct();
  }, [productId]);

  const handleSave = async () => {
    if (!product) return;
    setSaving(true);
    try {
      await pimApiClient.updateProduct(productId, editedProduct);
      showToast("success", "Produkt erfolgreich gespeichert");
      await fetchProduct();
    } catch (error: any) {
      console.error("Failed to save product:", error);
      showToast("error", error.message || "Fehler beim Speichern");
    } finally {
      setSaving(false);
    }
  };

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

  const hasChanges = Object.keys(editedProduct).length > 0;

  return (
    <MainLayout>
      {toast && (
        <div className={`fixed top-4 right-4 z-50 rounded-lg px-4 py-3 shadow-lg ${
          toast.type === "success" ? "bg-green-600" : "bg-red-600"
        } text-white`}>
          {toast.message}
        </div>
      )}

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
          <button
            onClick={handleSave}
            disabled={!hasChanges || saving}
            className="flex items-center gap-2 rounded-lg bg-primary-600 px-4 py-2 text-sm font-semibold text-white hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <Save className="h-4 w-4" />
            {saving ? "Speichert..." : "Speichern"}
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
            <MasterDataTab
              product={product}
              editedProduct={editedProduct}
              onEdit={setEditedProduct}
            />
          )}
          {activeTab === "prices" && (
            <PricesTab
              productId={productId}
              prices={prices}
              onUpdate={fetchProduct}
              onToast={showToast}
            />
          )}
          {activeTab === "categories" && (
            <CategoriesTab
              productId={productId}
              categories={categories}
              onUpdate={fetchProduct}
              onToast={showToast}
            />
          )}
          {activeTab === "attributes" && (
            <AttributesTab
              productId={productId}
              product={product}
              onUpdate={fetchProduct}
              onToast={showToast}
            />
          )}
          {activeTab === "variants" && (
            <VariantsTab
              product={product}
              variants={variants}
            />
          )}
          {activeTab === "bundles" && (
            <BundlesTab
              productId={productId}
              components={bundleComponents}
              onUpdate={fetchProduct}
              onToast={showToast}
            />
          )}
        </div>
      </div>
    </MainLayout>
  );
}

// Master Data Tab
function MasterDataTab({
  product,
  editedProduct,
  onEdit,
}: {
  product: Product;
  editedProduct: Partial<Product>;
  onEdit: (data: Partial<Product>) => void;
}) {
  const getValue = (field: keyof Product) => {
    if (editedProduct[field] !== undefined) {
      return editedProduct[field];
    }
    return product[field];
  };

  const getCurrentName = () => {
    return editedProduct.name || product.name;
  };

  const getCurrentDescription = () => {
    return editedProduct.description || product.description;
  };

  return (
    <div className="space-y-6">
      <div className="grid gap-6 md:grid-cols-2">
        <div>
          <label className="block text-sm font-medium text-gray-700">
            Name (DE)
          </label>
          <input
            type="text"
            value={getCurrentName()?.de || ""}
            onChange={(e) =>
              onEdit({
                ...editedProduct,
                name: { ...getCurrentName(), de: e.target.value },
              })
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
            value={getCurrentName()?.en || ""}
            onChange={(e) =>
              onEdit({
                ...editedProduct,
                name: { ...getCurrentName(), en: e.target.value },
              })
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
            value={getCurrentDescription()?.de || ""}
            onChange={(e) =>
              onEdit({
                ...editedProduct,
                description: {
                  ...getCurrentDescription(),
                  de: e.target.value,
                },
              })
            }
            rows={4}
            className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700">
            Beschreibung (EN)
          </label>
          <textarea
            value={getCurrentDescription()?.en || ""}
            onChange={(e) =>
              onEdit({
                ...editedProduct,
                description: {
                  ...getCurrentDescription(),
                  en: e.target.value,
                },
              })
            }
            rows={4}
            className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
          />
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-3">
        <div>
          <label className="block text-sm font-medium text-gray-700">SKU</label>
          <input
            type="text"
            value={product.sku}
            readOnly
            className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm bg-gray-50"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700">Typ</label>
          <input
            type="text"
            value={product.productType}
            readOnly
            className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm bg-gray-50"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700">Status</label>
          <select
            value={getValue("status") as string}
            onChange={(e) => onEdit({ ...editedProduct, status: e.target.value as any })}
            className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
          >
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
            <option value="draft">Draft</option>
          </select>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <div>
          <label className="block text-sm font-medium text-gray-700">Hersteller</label>
          <input
            type="text"
            value={(getValue("manufacturer") as string) || ""}
            onChange={(e) => onEdit({ ...editedProduct, manufacturer: e.target.value })}
            className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700">Hersteller-Teilenummer</label>
          <input
            type="text"
            value={(getValue("manufacturerPartNumber") as string) || ""}
            onChange={(e) => onEdit({ ...editedProduct, manufacturerPartNumber: e.target.value })}
            className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
          />
        </div>
      </div>
    </div>
  );
}

// Prices Tab
function PricesTab({
  productId,
  prices,
  onUpdate,
  onToast,
}: {
  productId: string;
  prices: PriceScale[];
  onUpdate: () => void;
  onToast: (type: "success" | "error", message: string) => void;
}) {
  const [showAddForm, setShowAddForm] = useState(false);
  const [newPrice, setNewPrice] = useState({ minQuantity: 1, price: 0, currency: "EUR" });
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editData, setEditData] = useState<Partial<PriceScale>>({});

  const handleAdd = async () => {
    try {
      await pimApiClient.createPrice(productId, newPrice);
      onToast("success", "Preis erfolgreich hinzugefügt");
      setShowAddForm(false);
      setNewPrice({ minQuantity: 1, price: 0, currency: "EUR" });
      onUpdate();
    } catch (error: any) {
      onToast("error", error.message || "Fehler beim Hinzufügen");
    }
  };

  const handleUpdate = async (priceId: string) => {
    try {
      await pimApiClient.updatePrice(productId, priceId, editData);
      onToast("success", "Preis erfolgreich aktualisiert");
      setEditingId(null);
      setEditData({});
      onUpdate();
    } catch (error: any) {
      onToast("error", error.message || "Fehler beim Aktualisieren");
    }
  };

  const handleDelete = async (priceId: string) => {
    if (!confirm("Preis wirklich löschen?")) return;
    try {
      await pimApiClient.deletePrice(productId, priceId);
      onToast("success", "Preis erfolgreich gelöscht");
      onUpdate();
    } catch (error: any) {
      onToast("error", error.message || "Fehler beim Löschen");
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">Staffelpreise</h3>
        <button
          onClick={() => setShowAddForm(!showAddForm)}
          className="flex items-center gap-2 rounded-lg bg-primary-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-primary-700"
        >
          {showAddForm ? <X className="h-4 w-4" /> : <Plus className="h-4 w-4" />}
          {showAddForm ? "Abbrechen" : "Preis hinzufügen"}
        </button>
      </div>

      {showAddForm && (
        <div className="rounded-lg border border-gray-200 p-4">
          <div className="grid gap-4 md:grid-cols-4">
            <div>
              <label className="block text-sm font-medium text-gray-700">Min. Menge</label>
              <input
                type="number"
                value={newPrice.minQuantity}
                onChange={(e) => setNewPrice({ ...newPrice, minQuantity: Number(e.target.value) })}
                className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Preis</label>
              <input
                type="number"
                step="0.01"
                value={newPrice.price}
                onChange={(e) => setNewPrice({ ...newPrice, price: Number(e.target.value) })}
                className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Währung</label>
              <select
                value={newPrice.currency}
                onChange={(e) => setNewPrice({ ...newPrice, currency: e.target.value })}
                className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm"
              >
                <option value="EUR">EUR</option>
                <option value="USD">USD</option>
                <option value="GBP">GBP</option>
              </select>
            </div>
            <div className="flex items-end">
              <button
                onClick={handleAdd}
                className="w-full rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700"
              >
                Hinzufügen
              </button>
            </div>
          </div>
        </div>
      )}

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
                  {editingId === price.id ? (
                    <>
                      <td className="px-4 py-3">
                        <input
                          type="number"
                          value={editData.minQuantity ?? price.minQuantity}
                          onChange={(e) => setEditData({ ...editData, minQuantity: Number(e.target.value) })}
                          className="w-full rounded border border-gray-300 px-2 py-1 text-sm"
                        />
                      </td>
                      <td className="px-4 py-3">
                        <input
                          type="number"
                          step="0.01"
                          value={editData.price ?? price.price}
                          onChange={(e) => setEditData({ ...editData, price: Number(e.target.value) })}
                          className="w-full rounded border border-gray-300 px-2 py-1 text-sm"
                        />
                      </td>
                      <td className="px-4 py-3">
                        <select
                          value={editData.currency ?? price.currency}
                          onChange={(e) => setEditData({ ...editData, currency: e.target.value })}
                          className="w-full rounded border border-gray-300 px-2 py-1 text-sm"
                        >
                          <option value="EUR">EUR</option>
                          <option value="USD">USD</option>
                          <option value="GBP">GBP</option>
                        </select>
                      </td>
                      <td className="px-4 py-3 text-right">
                        <button
                          onClick={() => handleUpdate(price.id)}
                          className="mr-2 text-green-600 hover:text-green-800"
                        >
                          Speichern
                        </button>
                        <button
                          onClick={() => {
                            setEditingId(null);
                            setEditData({});
                          }}
                          className="text-gray-600 hover:text-gray-800"
                        >
                          Abbrechen
                        </button>
                      </td>
                    </>
                  ) : (
                    <>
                      <td className="px-4 py-3 text-sm text-gray-900">{price.minQuantity}</td>
                      <td className="px-4 py-3 text-sm text-gray-900">{price.price.toFixed(2)}</td>
                      <td className="px-4 py-3 text-sm text-gray-900">{price.currency}</td>
                      <td className="px-4 py-3 text-right text-sm">
                        <button
                          onClick={() => {
                            setEditingId(price.id);
                            setEditData({ minQuantity: price.minQuantity, price: price.price, currency: price.currency });
                          }}
                          className="mr-3 text-primary-600 hover:text-primary-800"
                        >
                          <Edit className="inline h-4 w-4" />
                        </button>
                        <button
                          onClick={() => handleDelete(price.id)}
                          className="text-red-600 hover:text-red-800"
                        >
                          <Trash2 className="inline h-4 w-4" />
                        </button>
                      </td>
                    </>
                  )}
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
function CategoriesTab({
  productId,
  categories,
  onUpdate,
  onToast,
}: {
  productId: string;
  categories: Category[];
  onUpdate: () => void;
  onToast: (type: "success" | "error", message: string) => void;
}) {
  const [showAddForm, setShowAddForm] = useState(false);
  const [allCategories, setAllCategories] = useState<Category[]>([]);
  const [selectedCategoryId, setSelectedCategoryId] = useState("");
  const [loading, setLoading] = useState(false);

  const loadAllCategories = async () => {
    setLoading(true);
    try {
      const cats = await pimApiClient.getCategories();
      const flatten = (cats: Category[]): Category[] => {
        return cats.flatMap(c => [c, ...(c.children ? flatten(c.children) : [])]);
      };
      setAllCategories(flatten(cats));
    } catch (error: any) {
      onToast("error", "Fehler beim Laden der Kategorien");
    } finally {
      setLoading(false);
    }
  };

  const handleAdd = async () => {
    if (!selectedCategoryId) return;
    try {
      await pimApiClient.assignCategory(productId, selectedCategoryId);
      onToast("success", "Kategorie erfolgreich zugewiesen");
      setShowAddForm(false);
      setSelectedCategoryId("");
      onUpdate();
    } catch (error: any) {
      onToast("error", error.message || "Fehler beim Zuweisen");
    }
  };

  const handleRemove = async (categoryId: string) => {
    if (!confirm("Kategorie-Zuweisung wirklich entfernen?")) return;
    try {
      await pimApiClient.removeCategory(productId, categoryId);
      onToast("success", "Kategorie erfolgreich entfernt");
      onUpdate();
    } catch (error: any) {
      onToast("error", error.message || "Fehler beim Entfernen");
    }
  };

  useEffect(() => {
    if (showAddForm && allCategories.length === 0) {
      loadAllCategories();
    }
  }, [showAddForm]);

  const assignedIds = new Set(categories.map(c => c.id));
  const availableCategories = allCategories.filter(c => !assignedIds.has(c.id));

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">Zugewiesene Kategorien</h3>
        <button
          onClick={() => setShowAddForm(!showAddForm)}
          className="flex items-center gap-2 rounded-lg bg-primary-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-primary-700"
        >
          {showAddForm ? <X className="h-4 w-4" /> : <Plus className="h-4 w-4" />}
          {showAddForm ? "Abbrechen" : "Kategorie hinzufügen"}
        </button>
      </div>

      {showAddForm && (
        <div className="rounded-lg border border-gray-200 p-4">
          <div className="flex gap-4">
            <select
              value={selectedCategoryId}
              onChange={(e) => setSelectedCategoryId(e.target.value)}
              className="flex-1 rounded-lg border border-gray-300 px-3 py-2 text-sm"
              disabled={loading}
            >
              <option value="">Kategorie auswählen...</option>
              {availableCategories.map((cat) => (
                <option key={cat.id} value={cat.id}>
                  {cat.name.de || cat.name.en} ({cat.code})
                </option>
              ))}
            </select>
            <button
              onClick={handleAdd}
              disabled={!selectedCategoryId}
              className="rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 disabled:opacity-50"
            >
              Hinzufügen
            </button>
          </div>
        </div>
      )}

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
              <button
                onClick={() => handleRemove(cat.id)}
                className="flex items-center gap-1 text-red-600 hover:text-red-800"
              >
                <Trash2 className="h-4 w-4" />
                Entfernen
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

// Attributes Tab
function AttributesTab({
  productId,
  product,
  onUpdate,
  onToast,
}: {
  productId: string;
  product: Product;
  onUpdate: () => void;
  onToast: (type: "success" | "error", message: string) => void;
}) {
  const attributes = product.attributes || {};
  const [showAddForm, setShowAddForm] = useState(false);
  const [newAttr, setNewAttr] = useState({ key: "", type: "text", value: "" });
  const [editingKey, setEditingKey] = useState<string | null>(null);
  const [editValue, setEditValue] = useState<any>("");

  const handleAdd = async () => {
    if (!newAttr.key) {
      onToast("error", "Bitte einen Key eingeben");
      return;
    }
    try {
      await pimApiClient.createAttribute(productId, newAttr.key, newAttr.type, parseValue(newAttr.value, newAttr.type));
      onToast("success", "Attribut erfolgreich hinzugefügt");
      setShowAddForm(false);
      setNewAttr({ key: "", type: "text", value: "" });
      onUpdate();
    } catch (error: any) {
      onToast("error", error.message || "Fehler beim Hinzufügen");
    }
  };

  const handleUpdate = async (key: string) => {
    try {
      const parsedValue = typeof editValue === "string" ? tryParseJSON(editValue) : editValue;
      await pimApiClient.updateAttribute(productId, key, parsedValue);
      onToast("success", "Attribut erfolgreich aktualisiert");
      setEditingKey(null);
      setEditValue("");
      onUpdate();
    } catch (error: any) {
      onToast("error", error.message || "Fehler beim Aktualisieren");
    }
  };

  const handleDelete = async (key: string) => {
    if (!confirm(`Attribut "${key}" wirklich löschen?`)) return;
    try {
      await pimApiClient.deleteAttribute(productId, key);
      onToast("success", "Attribut erfolgreich gelöscht");
      onUpdate();
    } catch (error: any) {
      onToast("error", error.message || "Fehler beim Löschen");
    }
  };

  const parseValue = (val: string, type: string) => {
    if (type === "number") return Number(val);
    if (type === "boolean") return val === "true";
    if (type === "json") return JSON.parse(val);
    return val;
  };

  const tryParseJSON = (val: string) => {
    try {
      return JSON.parse(val);
    } catch {
      return val;
    }
  };

  const formatAttributeValue = (value: any): string => {
    if (value === null || value === undefined) return "-";
    if (typeof value === "object") {
      // Check if it's a metric value
      if (value.value !== undefined && value.unit !== undefined) {
        return `${value.value} ${value.unit}`;
      }
      // Check if it's an i18n value
      if (value.de !== undefined || value.en !== undefined) {
        const parts = [];
        if (value.de) parts.push(`DE: ${value.de}`);
        if (value.en) parts.push(`EN: ${value.en}`);
        return parts.join(" | ");
      }
      return JSON.stringify(value);
    }
    return String(value);
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">Produktattribute</h3>
        <button
          onClick={() => setShowAddForm(!showAddForm)}
          className="flex items-center gap-2 rounded-lg bg-primary-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-primary-700"
        >
          {showAddForm ? <X className="h-4 w-4" /> : <Plus className="h-4 w-4" />}
          {showAddForm ? "Abbrechen" : "Attribut hinzufügen"}
        </button>
      </div>

      {showAddForm && (
        <div className="rounded-lg border border-gray-200 p-4">
          <div className="grid gap-4 md:grid-cols-4">
            <div>
              <label className="block text-sm font-medium text-gray-700">Key</label>
              <input
                type="text"
                value={newAttr.key}
                onChange={(e) => setNewAttr({ ...newAttr, key: e.target.value })}
                className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm"
                placeholder="z.B. color"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Typ</label>
              <select
                value={newAttr.type}
                onChange={(e) => setNewAttr({ ...newAttr, type: e.target.value })}
                className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm"
              >
                <option value="text">Text</option>
                <option value="number">Number</option>
                <option value="boolean">Boolean</option>
                <option value="metric">Metric (Wert + Einheit)</option>
                <option value="i18n">I18n (mehrsprachig)</option>
                <option value="select">Select</option>
                <option value="json">JSON</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Wert</label>
              <input
                type="text"
                value={newAttr.value}
                onChange={(e) => setNewAttr({ ...newAttr, value: e.target.value })}
                className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm"
                placeholder={
                  newAttr.type === "metric"
                    ? '{"value": 10, "unit": "kg"}'
                    : newAttr.type === "i18n"
                    ? '{"de": "Text", "en": "Text"}'
                    : newAttr.type === "json"
                    ? '{"key": "value"}'
                    : "Wert"
                }
              />
            </div>
            <div className="flex items-end">
              <button
                onClick={handleAdd}
                className="w-full rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700"
              >
                Hinzufügen
              </button>
            </div>
          </div>
        </div>
      )}

      {Object.keys(attributes).length === 0 ? (
        <p className="text-sm text-gray-500">Keine Attribute definiert</p>
      ) : (
        <div className="space-y-2">
          {Object.entries(attributes).map(([key, value]) => (
            <div
              key={key}
              className="flex items-center justify-between rounded-lg border border-gray-200 p-3"
            >
              {editingKey === key ? (
                <>
                  <div className="flex-1">
                    <p className="font-medium text-gray-900 mb-2">{key}</p>
                    <input
                      type="text"
                      value={typeof editValue === "object" ? JSON.stringify(editValue) : String(editValue)}
                      onChange={(e) => setEditValue(e.target.value)}
                      className="w-full rounded border border-gray-300 px-2 py-1 text-sm"
                    />
                  </div>
                  <div className="ml-4">
                    <button
                      onClick={() => handleUpdate(key)}
                      className="mr-2 text-green-600 hover:text-green-800"
                    >
                      Speichern
                    </button>
                    <button
                      onClick={() => {
                        setEditingKey(null);
                        setEditValue("");
                      }}
                      className="text-gray-600 hover:text-gray-800"
                    >
                      Abbrechen
                    </button>
                  </div>
                </>
              ) : (
                <>
                  <div>
                    <p className="font-medium text-gray-900">{key}</p>
                    <p className="text-sm text-gray-500">
                      {formatAttributeValue(value)}
                    </p>
                  </div>
                  <div>
                    <button
                      onClick={() => {
                        setEditingKey(key);
                        setEditValue(value);
                      }}
                      className="mr-3 text-primary-600 hover:text-primary-800"
                    >
                      <Edit className="inline h-4 w-4" />
                    </button>
                    <button
                      onClick={() => handleDelete(key)}
                      className="text-red-600 hover:text-red-800"
                    >
                      <Trash2 className="inline h-4 w-4" />
                    </button>
                  </div>
                </>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

// Variants Tab
function VariantsTab({ product, variants }: { product: Product; variants: Product[] }) {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">Kind-Varianten</h3>
        <Link
          href={`/products/new?parent=${product.id}`}
          className="flex items-center gap-2 rounded-lg bg-primary-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-primary-700"
        >
          <Plus className="h-4 w-4" />
          Variante hinzufügen
        </Link>
      </div>

      {variants.length === 0 ? (
        <p className="text-sm text-gray-500">Keine Varianten definiert</p>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="border-b border-gray-200 bg-gray-50">
              <tr>
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-gray-500">
                  SKU
                </th>
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-gray-500">
                  Name
                </th>
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-gray-500">
                  Varianten-Achsen
                </th>
                <th className="px-4 py-2 text-left text-xs font-medium uppercase text-gray-500">
                  Status
                </th>
                <th className="px-4 py-2 text-right text-xs font-medium uppercase text-gray-500">
                  Aktionen
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {variants.map((variant) => (
                <tr key={variant.id}>
                  <td className="px-4 py-3 text-sm font-medium text-gray-900">{variant.sku}</td>
                  <td className="px-4 py-3 text-sm text-gray-900">
                    {variant.name.de || variant.name.en}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-500">
                    {variant.variantAxes ? Object.entries(variant.variantAxes).map(([k, v]) => `${k}: ${v}`).join(", ") : "-"}
                  </td>
                  <td className="px-4 py-3 text-sm">
                    <span
                      className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${
                        variant.status === "active"
                          ? "bg-green-100 text-green-800"
                          : variant.status === "inactive"
                          ? "bg-red-100 text-red-800"
                          : "bg-yellow-100 text-yellow-800"
                      }`}
                    >
                      {variant.status}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-right">
                    <Link
                      href={`/products/${variant.id}`}
                      className="text-primary-600 hover:text-primary-800"
                    >
                      Details
                    </Link>
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

// Bundles Tab
function BundlesTab({
  productId,
  components,
  onUpdate,
  onToast,
}: {
  productId: string;
  components: BundleComponent[];
  onUpdate: () => void;
  onToast: (type: "success" | "error", message: string) => void;
}) {
  const [showAddForm, setShowAddForm] = useState(false);
  const [newComponent, setNewComponent] = useState({
    componentProductId: "",
    quantity: 1,
    minQuantity: undefined as number | undefined,
    maxQuantity: undefined as number | undefined,
  });
  const [searchResults, setSearchResults] = useState<Product[]>([]);
  const [searching, setSearching] = useState(false);

  const searchProducts = async (query: string) => {
    if (query.length < 2) {
      setSearchResults([]);
      return;
    }
    setSearching(true);
    try {
      const result = await pimApiClient.getProducts({ q: query, limit: 10 });
      setSearchResults(result.items);
    } catch (error) {
      onToast("error", "Fehler bei der Produktsuche");
    } finally {
      setSearching(false);
    }
  };

  const handleAdd = async () => {
    if (!newComponent.componentProductId) {
      onToast("error", "Bitte ein Produkt auswählen");
      return;
    }
    try {
      await pimApiClient.addBundleComponent(productId, newComponent);
      onToast("success", "Komponente erfolgreich hinzugefügt");
      setShowAddForm(false);
      setNewComponent({
        componentProductId: "",
        quantity: 1,
        minQuantity: undefined,
        maxQuantity: undefined,
      });
      setSearchResults([]);
      onUpdate();
    } catch (error: any) {
      onToast("error", error.message || "Fehler beim Hinzufügen");
    }
  };

  const handleRemove = async (componentId: string) => {
    if (!confirm("Komponente wirklich entfernen?")) return;
    try {
      await pimApiClient.removeBundleComponent(productId, componentId);
      onToast("success", "Komponente erfolgreich entfernt");
      onUpdate();
    } catch (error: any) {
      onToast("error", error.message || "Fehler beim Entfernen");
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">Bundle-Komponenten</h3>
        <button
          onClick={() => setShowAddForm(!showAddForm)}
          className="flex items-center gap-2 rounded-lg bg-primary-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-primary-700"
        >
          {showAddForm ? <X className="h-4 w-4" /> : <Plus className="h-4 w-4" />}
          {showAddForm ? "Abbrechen" : "Komponente hinzufügen"}
        </button>
      </div>

      {showAddForm && (
        <div className="rounded-lg border border-gray-200 p-4 space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Produkt suchen</label>
            <input
              type="text"
              onChange={(e) => searchProducts(e.target.value)}
              className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm"
              placeholder="SKU oder Name eingeben..."
            />
            {searchResults.length > 0 && (
              <div className="mt-2 max-h-40 overflow-y-auto border border-gray-200 rounded">
                {searchResults.map((p) => (
                  <button
                    key={p.id}
                    onClick={() => {
                      setNewComponent({ ...newComponent, componentProductId: p.id });
                      setSearchResults([]);
                    }}
                    className="w-full px-3 py-2 text-left text-sm hover:bg-gray-100"
                  >
                    {p.sku} - {p.name.de || p.name.en}
                  </button>
                ))}
              </div>
            )}
            {newComponent.componentProductId && (
              <p className="mt-2 text-sm text-green-600">Produkt ausgewählt: {newComponent.componentProductId}</p>
            )}
          </div>

          <div className="grid gap-4 md:grid-cols-4">
            <div>
              <label className="block text-sm font-medium text-gray-700">Menge</label>
              <input
                type="number"
                value={newComponent.quantity}
                onChange={(e) => setNewComponent({ ...newComponent, quantity: Number(e.target.value) })}
                className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Min. Menge</label>
              <input
                type="number"
                value={newComponent.minQuantity || ""}
                onChange={(e) => setNewComponent({ ...newComponent, minQuantity: e.target.value ? Number(e.target.value) : undefined })}
                className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Max. Menge</label>
              <input
                type="number"
                value={newComponent.maxQuantity || ""}
                onChange={(e) => setNewComponent({ ...newComponent, maxQuantity: e.target.value ? Number(e.target.value) : undefined })}
                className="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm"
              />
            </div>
            <div className="flex items-end">
              <button
                onClick={handleAdd}
                className="w-full rounded-lg bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700"
              >
                Hinzufügen
              </button>
            </div>
          </div>
        </div>
      )}

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
                    {comp.product?.name.de || comp.product?.name.en || comp.componentProductId}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-900">{comp.quantity}</td>
                  <td className="px-4 py-3 text-sm text-gray-500">
                    {comp.minQuantity || "-"} / {comp.maxQuantity || "-"}
                  </td>
                  <td className="px-4 py-3 text-right text-sm">
                    <button
                      onClick={() => handleRemove(comp.id)}
                      className="text-red-600 hover:text-red-800"
                    >
                      <Trash2 className="inline h-4 w-4" />
                    </button>
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
