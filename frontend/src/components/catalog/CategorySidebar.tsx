"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname, useSearchParams } from "next/navigation";
import type { Category } from "@/types/catalog";
import { Panel, PanelHeader, PanelBody } from "@/components/ui/Panel";
import { cn } from "@/lib/utils";

interface CategorySidebarProps {
  categories: Category[];
  selectedCategoryId?: string;
}

export function CategorySidebar({
  categories,
  selectedCategoryId,
}: CategorySidebarProps) {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  // Auto-expand categories that contain the selected category or ARE selected
  const getInitialExpanded = (): Set<string> => {
    const expanded = new Set<string>();
    if (!selectedCategoryId) return expanded;
    const findPath = (cats: Category[]): boolean => {
      for (const cat of cats) {
        if (cat.id === selectedCategoryId) {
          // If this category has children, expand it to show them
          if (cat.children && cat.children.length > 0) {
            expanded.add(cat.id);
          }
          return true;
        }
        if (cat.children && findPath(cat.children)) {
          expanded.add(cat.id);
          return true;
        }
      }
      return false;
    };
    findPath(categories);
    return expanded;
  };

  const [expandedCategories, setExpandedCategories] = useState<Set<string>>(
    getInitialExpanded
  );

  const toggleCategory = (categoryId: string) => {
    setExpandedCategories((prev) => {
      const next = new Set(prev);
      if (next.has(categoryId)) {
        next.delete(categoryId);
      } else {
        next.add(categoryId);
      }
      return next;
    });
  };

  const buildCategoryUrl = (categoryId: string) => {
    const params = new URLSearchParams(searchParams);
    params.set("category", categoryId);
    params.delete("page"); // Reset page when changing category
    return `${pathname}?${params.toString()}`;
  };

  const renderCategory = (category: Category, level = 0) => {
    const hasChildren = category.children && category.children.length > 0;
    const isExpanded = expandedCategories.has(category.id);
    const isSelected = category.id === selectedCategoryId;

    return (
      <div key={category.id}>
        <div
          className={cn(
            "flex items-center justify-between py-2 px-3 rounded-md transition-colors",
            level > 0 && "ml-4",
            isSelected
              ? "bg-primary-100 dark:bg-primary-900/30 text-primary-900 dark:text-primary-200 font-medium"
              : "hover:bg-gray-100 dark:hover:bg-gray-800 text-gray-700 dark:text-gray-300"
          )}
        >
          <Link
            href={buildCategoryUrl(category.id)}
            className="flex-1 text-sm"
          >
            {category.name}
            {category.productCount !== undefined && (
              <span className="ml-2 text-xs text-gray-500 dark:text-gray-400">
                ({category.productCount})
              </span>
            )}
          </Link>
          {hasChildren && (
            <button
              onClick={() => toggleCategory(category.id)}
              className="p-1 hover:bg-gray-200 dark:hover:bg-gray-700 rounded"
            >
              <svg
                className={cn(
                  "w-4 h-4 transition-transform",
                  isExpanded && "rotate-90"
                )}
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
            </button>
          )}
        </div>
        {hasChildren && isExpanded && (
          <div className="mt-1">
            {category.children!.map((child) => renderCategory(child, level + 1))}
          </div>
        )}
      </div>
    );
  };

  return (
    <Panel>
      <PanelHeader>
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
          Kategorien
        </h2>
      </PanelHeader>
      <PanelBody className="space-y-1">
        <Link
          href="/products"
          className={cn(
            "block py-2 px-3 rounded-md text-sm transition-colors",
            !selectedCategoryId
              ? "bg-primary-100 dark:bg-primary-900/30 text-primary-900 dark:text-primary-200 font-medium"
              : "hover:bg-gray-100 dark:hover:bg-gray-800 text-gray-700 dark:text-gray-300"
          )}
        >
          Alle Produkte
        </Link>
        {categories
          .filter((cat) => !cat.parentId)
          .map((category) => renderCategory(category))}
      </PanelBody>
    </Panel>
  );
}
