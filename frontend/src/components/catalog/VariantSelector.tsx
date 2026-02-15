"use client";

import { useState, useEffect } from "react";
import type { VariantAxis, AxisOption } from "@/types/catalog";

interface VariantSelectorProps {
  axes: VariantAxis[];
  selectedValues: Record<string, string>;
  onSelect: (axisCode: string, optionCode: string) => void;
  availableValues?: Record<string, AxisOption[]>;
  className?: string;
}

export function VariantSelector({
  axes,
  selectedValues,
  onSelect,
  availableValues,
  className = "",
}: VariantSelectorProps) {
  // Sort axes by position
  const sortedAxes = [...axes].sort((a, b) => a.position - b.position);

  const getLabel = (labels: Record<string, string>): string => {
    return labels.de || labels.en || Object.values(labels)[0] || "";
  };

  const isOptionAvailable = (axisCode: string, optionCode: string): boolean => {
    if (!availableValues || !availableValues[axisCode]) {
      return true; // If no availability data, assume available
    }
    const option = availableValues[axisCode].find((opt) => opt.code === optionCode);
    return option?.available !== false;
  };

  return (
    <div className={`space-y-6 ${className}`}>
      {sortedAxes.map((axis) => {
        const axisLabel = getLabel(axis.label);
        const options = [...axis.options].sort((a, b) => a.position - b.position);
        const useDropdown = options.length > 6;

        return (
          <div key={axis.attributeCode}>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
              {axisLabel}
            </label>

            {useDropdown ? (
              // Dropdown für >6 Optionen
              <select
                value={selectedValues[axis.attributeCode] || ""}
                onChange={(e) => onSelect(axis.attributeCode, e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              >
                <option value="">-- Bitte wählen --</option>
                {options.map((option) => {
                  const available = isOptionAvailable(axis.attributeCode, option.code);
                  const label = getLabel(option.label);
                  return (
                    <option
                      key={option.code}
                      value={option.code}
                      disabled={!available}
                    >
                      {label} {!available ? "(Nicht verfügbar)" : ""}
                    </option>
                  );
                })}
              </select>
            ) : (
              // Buttons für ≤6 Optionen
              <div className="flex flex-wrap gap-2">
                {options.map((option) => {
                  const available = isOptionAvailable(axis.attributeCode, option.code);
                  const isSelected = selectedValues[axis.attributeCode] === option.code;
                  const label = getLabel(option.label);

                  return (
                    <button
                      key={option.code}
                      type="button"
                      onClick={() => available && onSelect(axis.attributeCode, option.code)}
                      disabled={!available}
                      className={`
                        px-4 py-2 rounded-md border-2 font-medium transition-all
                        ${
                          isSelected
                            ? "border-primary-600 bg-primary-50 dark:bg-primary-900/20 text-primary-700 dark:text-primary-300"
                            : available
                            ? "border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:border-primary-400 dark:hover:border-primary-500"
                            : "border-gray-200 dark:border-gray-700 bg-gray-100 dark:bg-gray-900 text-gray-400 dark:text-gray-600 cursor-not-allowed line-through"
                        }
                      `}
                      title={!available ? "Nicht verfügbar für diese Kombination" : ""}
                    >
                      {label}
                    </button>
                  );
                })}
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
}
