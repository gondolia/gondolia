"use client";

import { useEffect, useState } from "react";
import { MainLayout } from "@/components/MainLayout";
import { pimApiClient } from "@/lib/api/client";
import type { Category } from "@/types/catalog";
import {
  ChevronRight,
  ChevronDown,
  Plus,
  Edit,
  Trash2,
  FolderTree,
} from "lucide-react";

export default function CategoriesPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(null);
  const [expandedIds, setExpandedIds] = useState<Set<string>>(new Set());
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchCategories = async () => {
      setLoading(true);
      try {
        const data = await pimApiClient.getCategories();
        setCategories(data);
      } catch (error) {
        console.error("Failed to fetch categories:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchCategories();
  }, []);

  const toggleExpand = (id: string) => {
    const newExpanded = new Set(expandedIds);
    if (newExpanded.has(id)) {
      newExpanded.delete(id);
    } else {
      newExpanded.add(id);
    }
    setExpandedIds(newExpanded);
  };

  const renderCategory = (category: Category, level: number = 0) => {
    const hasChildren = category.children && category.children.length > 0;
    const isExpanded = expandedIds.has(category.id);
    const isSelected = selectedCategory?.id === category.id;

    return (
      <div key={category.id}>
        <div
          className={`flex items-center gap-2 rounded-lg px-3 py-2 transition-colors ${
            isSelected ? "bg-primary-100" : "hover:bg-gray-100"
          }`}
          style={{ paddingLeft: `${level * 1.5 + 0.75}rem` }}
        >
          <button
            onClick={() => hasChildren && toggleExpand(category.id)}
            className="flex-shrink-0"
          >
            {hasChildren ? (
              isExpanded ? (
                <ChevronDown className="h-4 w-4 text-gray-500" />
              ) : (
                <ChevronRight className="h-4 w-4 text-gray-500" />
              )
            ) : (
              <div className="h-4 w-4" />
            )}
          </button>

          <FolderTree className="h-4 w-4 text-gray-400" />

          <button
            onClick={() => setSelectedCategory(category)}
            className="flex-1 text-left"
          >
            <span className="text-sm font-medium text-gray-900">
              {category.name.de || category.name.en || category.code}
            </span>
            {category.productCount !== undefined && (
              <span className="ml-2 text-xs text-gray-500">
                ({category.productCount})
              </span>
            )}
          </button>
        </div>

        {hasChildren && isExpanded && (
          <div>
            {category.children!.map((child) =>
              renderCategory(child, level + 1)
            )}
          </div>
        )}
      </div>
    );
  };

  return (
    <MainLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Kategorien</h1>
            <p className="mt-1 text-sm text-gray-500">
              Kategorie-Hierarchie verwalten
            </p>
          </div>
          <button className="flex items-center gap-2 rounded-lg bg-primary-600 px-4 py-2 text-sm font-semibold text-white hover:bg-primary-700">
            <Plus className="h-4 w-4" />
            Neue Kategorie
          </button>
        </div>

        {/* Layout: Tree on left, Details on right */}
        <div className="grid gap-6 lg:grid-cols-3">
          {/* Category Tree */}
          <div className="lg:col-span-1">
            <div className="rounded-lg bg-white p-4 shadow-sm ring-1 ring-gray-200">
              <h2 className="mb-4 text-sm font-semibold uppercase text-gray-500">
                Kategorie-Baum
              </h2>
              {loading ? (
                <p className="text-sm text-gray-500">Lädt...</p>
              ) : categories.length === 0 ? (
                <p className="text-sm text-gray-500">Keine Kategorien vorhanden</p>
              ) : (
                <div className="space-y-1">
                  {categories.map((cat) => renderCategory(cat))}
                </div>
              )}
            </div>
          </div>

          {/* Category Details */}
          <div className="lg:col-span-2">
            <div className="rounded-lg bg-white p-6 shadow-sm ring-1 ring-gray-200">
              {selectedCategory ? (
                <div className="space-y-6">
                  <div className="flex items-center justify-between">
                    <h2 className="text-lg font-semibold text-gray-900">
                      Kategorie-Details
                    </h2>
                    <div className="flex gap-2">
                      <button className="flex items-center gap-2 rounded-lg border border-gray-300 px-3 py-1.5 text-sm font-medium text-gray-700 hover:bg-gray-50">
                        <Edit className="h-4 w-4" />
                        Bearbeiten
                      </button>
                      <button className="flex items-center gap-2 rounded-lg border border-red-300 px-3 py-1.5 text-sm font-medium text-red-700 hover:bg-red-50">
                        <Trash2 className="h-4 w-4" />
                        Löschen
                      </button>
                    </div>
                  </div>

                  <div className="space-y-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700">
                        Code
                      </label>
                      <p className="mt-1 text-sm text-gray-900">
                        {selectedCategory.code}
                      </p>
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700">
                        Name (DE)
                      </label>
                      <p className="mt-1 text-sm text-gray-900">
                        {selectedCategory.name.de || "-"}
                      </p>
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700">
                        Name (EN)
                      </label>
                      <p className="mt-1 text-sm text-gray-900">
                        {selectedCategory.name.en || "-"}
                      </p>
                    </div>

                    {selectedCategory.description && (
                      <div>
                        <label className="block text-sm font-medium text-gray-700">
                          Beschreibung (DE)
                        </label>
                        <p className="mt-1 text-sm text-gray-900">
                          {selectedCategory.description.de || "-"}
                        </p>
                      </div>
                    )}

                    <div className="grid gap-4 md:grid-cols-2">
                      <div>
                        <label className="block text-sm font-medium text-gray-700">
                          Status
                        </label>
                        <p className="mt-1 text-sm">
                          <span
                            className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${
                              selectedCategory.isActive
                                ? "bg-green-100 text-green-800"
                                : "bg-red-100 text-red-800"
                            }`}
                          >
                            {selectedCategory.isActive ? "Aktiv" : "Inaktiv"}
                          </span>
                        </p>
                      </div>

                      <div>
                        <label className="block text-sm font-medium text-gray-700">
                          Anzahl Produkte
                        </label>
                        <p className="mt-1 text-sm text-gray-900">
                          {selectedCategory.productCount || 0}
                        </p>
                      </div>
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700">
                        Sortierung
                      </label>
                      <p className="mt-1 text-sm text-gray-900">
                        {selectedCategory.sortOrder}
                      </p>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="flex h-64 items-center justify-center">
                  <p className="text-sm text-gray-500">
                    Wähle eine Kategorie aus dem Baum
                  </p>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </MainLayout>
  );
}
