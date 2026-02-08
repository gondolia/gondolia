"use client";

import type { HTMLAttributes } from "react";
import { cn } from "@/lib/utils";

export interface PanelProps extends HTMLAttributes<HTMLDivElement> {
  variant?: "default" | "elevated";
}

export function Panel({
  className,
  variant = "default",
  children,
  ...props
}: PanelProps) {
  return (
    <div
      className={cn(
        "rounded-lg bg-white dark:bg-gray-900",
        variant === "default" && "border border-gray-200 dark:border-gray-700",
        variant === "elevated" && "shadow-lg",
        className
      )}
      {...props}
    >
      {children}
    </div>
  );
}

export function PanelHeader({
  className,
  children,
  ...props
}: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        "border-b border-gray-200 dark:border-gray-700 px-6 py-4",
        className
      )}
      {...props}
    >
      {children}
    </div>
  );
}

export function PanelBody({
  className,
  children,
  ...props
}: HTMLAttributes<HTMLDivElement>) {
  return (
    <div className={cn("px-6 py-4", className)} {...props}>
      {children}
    </div>
  );
}
