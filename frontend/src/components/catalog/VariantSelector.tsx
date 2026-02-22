"use client";

import { useState, useEffect, useMemo } from "react";
import type { VariantAxis, AxisOption, ProductVariant } from "@/types/catalog";

interface VariantSelectorProps {
  axes: VariantAxis[];
  selectedValues: Record<string, string>;
  onSelect: (axisCode: string, optionCode: string) => void;
  variants?: ProductVariant[]; // Array aller Varianten für intelligente Filterung
  availableValues?: Record<string, AxisOption[]>; // Optional: Backend-calculated availability (fallback)
  className?: string;
}

export function VariantSelector({
  axes,
  selectedValues,
  onSelect,
  variants,
  availableValues,
  className = "",
}: VariantSelectorProps) {
  // Sort axes by position
  const sortedAxes = [...axes].sort((a, b) => a.position - b.position);

  const getLabel = (labels: Record<string, string>): string => {
    return labels.de || labels.en || Object.values(labels)[0] || "";
  };

  // BUG 1 FIX: Berechne verfügbare Optionen basierend auf existierenden Varianten
  // und bereits gewählten Achsenwerten
  // Berechne verfügbare Optionen pro Achse.
  // Wichtig: Beim Filtern für eine Achse wird deren eigene Auswahl NICHT berücksichtigt,
  // damit man frei zwischen Optionen derselben Achse wechseln kann.
  const computedAvailableValues = useMemo(() => {
    if (!variants || variants.length === 0) {
      return undefined;
    }

    const available: Record<string, Set<string>> = {};
    const axisCodes = axes.map(a => a.attributeCode);

    for (const targetAxis of axisCodes) {
      available[targetAxis] = new Set();

      // Filtere nur nach den Auswahlen der ANDEREN Achsen
      const otherSelections = Object.entries(selectedValues).filter(
        ([code, v]) => code !== targetAxis && v !== "" && v != null
      );

      const filteredVariants = variants.filter((variant) =>
        otherSelections.every(([axisCode, optionCode]) =>
          variant.axisValues[axisCode] === optionCode
        )
      );

      filteredVariants.forEach((variant) => {
        const val = variant.axisValues[targetAxis];
        if (val) available[targetAxis].add(val);
      });
    }

    return available;
  }, [variants, selectedValues, axes]);

  const isOptionAvailable = (axisCode: string, optionCode: string): boolean => {
    // Verwende computed availability wenn variants vorhanden
    if (computedAvailableValues) {
      const availableCodes = computedAvailableValues[axisCode];
      return availableCodes ? availableCodes.has(optionCode) : false;
    }
    
    // Fallback: Backend-provided availability
    if (availableValues && availableValues[axisCode]) {
      const option = availableValues[axisCode].find((opt) => opt.code === optionCode);
      return option?.available !== false;
    }
    
    // Wenn keine Availability-Daten, alles verfügbar
    return true;
  };

  const allAxesSelected = sortedAxes.every((axis) => selectedValues[axis.attributeCode] && selectedValues[axis.attributeCode] !== "");

  // BUG 3 FIX: Prüfe ob die gewählte Kombination als Variante existiert
  const hasValidVariant = useMemo(() => {
    if (!allAxesSelected || !variants || variants.length === 0) {
      return true; // Noch nicht alle gewählt oder keine Varianten-Daten
    }

    return variants.some((variant) => {
      return Object.entries(selectedValues).every(([axisCode, optionCode]) => {
        return variant.axisValues[axisCode] === optionCode;
      });
    });
  }, [allAxesSelected, variants, selectedValues]);

  return (
    <div className={`space-y-5 ${className}`}>
      {/* Progress indicator + Reset */}
      <div className="flex items-center gap-2 text-sm">
        <div className="flex items-center gap-1.5">
          {sortedAxes.map((axis, idx) => {
            const isSelected = !!selectedValues[axis.attributeCode];
            return (
              <div key={axis.attributeCode} className="flex items-center gap-1.5">
                <div
                  className={`
                    w-6 h-6 rounded-full flex items-center justify-center text-xs font-semibold
                    ${
                      isSelected
                        ? "bg-primary-600 text-white"
                        : "bg-gray-200 dark:bg-gray-700 text-gray-500 dark:text-gray-400"
                    }
                  `}
                >
                  {isSelected ? "✓" : idx + 1}
                </div>
                {idx < sortedAxes.length - 1 && (
                  <div
                    className={`w-8 h-0.5 ${
                      isSelected
                        ? "bg-primary-600"
                        : "bg-gray-200 dark:bg-gray-700"
                    }`}
                  />
                )}
              </div>
            );
          })}
        </div>
        <span className="text-gray-600 dark:text-gray-400 ml-2">
          {Object.values(selectedValues).filter(v => v !== "" && v != null).length} / {sortedAxes.length} ausgewählt
        </span>
        {Object.values(selectedValues).some(v => v !== "" && v != null) && (
          <button
            type="button"
            onClick={() => {
              sortedAxes.forEach(axis => onSelect(axis.attributeCode, ""));
            }}
            className="ml-auto text-xs text-primary-600 dark:text-primary-400 hover:underline font-medium"
          >
            Alle zurücksetzen
          </button>
        )}
      </div>

      {/* Axis selectors */}
      <div className="space-y-4">
        {sortedAxes.map((axis, axisIdx) => {
          const axisLabel = getLabel(axis.label);
          const options = [...axis.options].sort((a, b) => a.position - b.position);
          const useDropdown = options.length > 6; // Design decision: >6 → Dropdown
          const isCurrentAxisSelected = !!selectedValues[axis.attributeCode];

          return (
            <div
              key={axis.attributeCode}
              className={`
                p-4 rounded-lg border-2 transition-all
                ${
                  isCurrentAxisSelected
                    ? "border-primary-500 bg-primary-50 dark:bg-primary-900/10"
                    : "border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800"
                }
              `}
            >
              {/* Axis header */}
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-3">
                  <div
                    className={`
                      w-7 h-7 rounded-full flex items-center justify-center text-sm font-bold
                      ${
                        isCurrentAxisSelected
                          ? "bg-primary-600 text-white"
                          : "bg-gray-300 dark:bg-gray-700 text-gray-600 dark:text-gray-400"
                      }
                    `}
                  >
                    {isCurrentAxisSelected ? "✓" : axisIdx + 1}
                  </div>
                  <label className="text-base font-semibold text-gray-900 dark:text-white">
                    {axisLabel}
                  </label>
                </div>
                {/* Zurücksetzen entfernt — Benutzer kann direkt zwischen Optionen wechseln */}
              </div>

              {useDropdown ? (
                // Dropdown für >8 Optionen
                <select
                  value={selectedValues[axis.attributeCode] || ""}
                  onChange={(e) => onSelect(axis.attributeCode, e.target.value)}
                  className="w-full px-4 py-3 border-2 border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-900 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-primary-500 font-medium"
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
                        {label} {!available ? "✕ (Nicht verfügbar)" : ""}
                      </option>
                    );
                  })}
                </select>
              ) : (
                // Button-Grid für ≤8 Optionen
                <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-2.5">
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
                          relative px-4 py-3 rounded-lg border-2 font-semibold text-sm transition-all
                          flex items-center justify-center gap-2 min-h-[3rem]
                          ${
                            isSelected
                              ? "border-primary-600 bg-primary-600 text-white shadow-md scale-105"
                              : available
                              ? "border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-900 text-gray-700 dark:text-gray-300 hover:border-primary-500 dark:hover:border-primary-400 hover:bg-gray-50 dark:hover:bg-gray-800 hover:scale-102"
                              : "border-gray-200 dark:border-gray-800 bg-gray-100 dark:bg-gray-900/50 text-gray-400 dark:text-gray-600 cursor-not-allowed opacity-50"
                          }
                        `}
                        title={!available ? "Diese Option ist nicht verfügbar für die gewählte Kombination" : label}
                      >
                        {isSelected && (
                          <span className="text-lg leading-none">✓</span>
                        )}
                        <span className={available ? "" : "line-through"}>
                          {label}
                        </span>
                      </button>
                    );
                  })}
                </div>
              )}
            </div>
          );
        })}
      </div>

      {/* Status message */}
      {allAxesSelected ? (
        hasValidVariant ? (
          <div className="flex items-center gap-2 p-3 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg">
            <svg
              className="w-5 h-5 text-green-600 dark:text-green-400 flex-shrink-0"
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
            <span className="text-sm font-medium text-green-800 dark:text-green-200">
              Alle Varianten-Merkmale ausgewählt
            </span>
          </div>
        ) : (
          // BUG 3 FIX: Zeige "nicht verfügbar" statt "bitte wählen" wenn ungültige Kombination
          <div className="flex items-center gap-2 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
            <svg
              className="w-5 h-5 text-red-600 dark:text-red-400 flex-shrink-0"
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
            <span className="text-sm font-medium text-red-800 dark:text-red-200">
              Diese Kombination ist nicht verfügbar
            </span>
          </div>
        )
      ) : (
        <div className="flex items-center gap-2 p-3 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
          <svg
            className="w-5 h-5 text-blue-600 dark:text-blue-400 flex-shrink-0"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <span className="text-sm text-blue-800 dark:text-blue-200">
            Bitte wählen Sie alle Merkmale aus, um das Produkt in den Warenkorb zu legen
          </span>
        </div>
      )}
    </div>
  );
}
