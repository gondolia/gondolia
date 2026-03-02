import { test, expect } from '@playwright/test';

test.describe('Navigation', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/login');
    await page.getByLabel('Email').fill('admin@demo.local');
    await page.getByLabel('Password').fill('admin123');
    await page.getByRole('button', { name: 'Sign in' }).click();
    await expect(page).toHaveURL('/dashboard');
  });

  test('products page is reachable', async ({ page }) => {
    // Click on Produkte link in navigation
    await page.getByRole('link', { name: 'Produkte' }).click();

    // Check we're on the products page
    await expect(page).toHaveURL('/products');
    await expect(page.getByRole('heading', { name: 'Produkte' })).toBeVisible();

    // Check product type filters are visible
    await expect(page.getByRole('button', { name: 'Alle' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Einfach' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Varianten' })).toBeVisible();
  });

  test('categories page is reachable', async ({ page }) => {
    // Click on Kategorien link in navigation
    await page.getByRole('link', { name: 'Kategorien' }).click();

    // Check we're on the categories page
    await expect(page).toHaveURL('/categories');
    await expect(page.getByRole('heading', { name: 'Kategorien' })).toBeVisible();
  });

  test('cart drawer opens', async ({ page }) => {
    // Find and click the cart button
    await page.getByRole('button', { name: 'Warenkorb' }).click();

    // Check that cart drawer is visible (CartDrawer component)
    await expect(page.getByRole('heading', { name: 'Warenkorb' })).toBeVisible();
  });

  test('dashboard is reachable from header logo', async ({ page }) => {
    // Navigate away from dashboard first
    await page.getByRole('link', { name: 'Produkte' }).click();
    await expect(page).toHaveURL('/products');

    // Click on logo/brand name to go back to dashboard
    await page.getByText('Gondolia').first().click();

    // Should be on dashboard
    await expect(page).toHaveURL('/dashboard');
  });

  test('main navigation items are visible', async ({ page }) => {
    // Check all main navigation items are visible
    await expect(page.getByRole('link', { name: 'Dashboard' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Produkte' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Kategorien' })).toBeVisible();

    // Check cart icon is visible
    await expect(page.getByRole('button', { name: 'Warenkorb' })).toBeVisible();

    // Check logout button is visible
    await expect(page.getByRole('button', { name: 'Logout' })).toBeVisible();
  });

  test('navigation highlights active page', async ({ page }) => {
    // On dashboard, Dashboard link should be active
    const dashboardLink = page.getByRole('link', { name: 'Dashboard' }).locator('span');
    await expect(dashboardLink).toHaveClass(/bg-primary-100|bg-primary-900/);

    // Navigate to products
    await page.getByRole('link', { name: 'Produkte' }).click();
    await expect(page).toHaveURL('/products');

    // Products link should now be active
    const productsLink = page.getByRole('link', { name: 'Produkte' }).locator('span');
    await expect(productsLink).toHaveClass(/bg-primary-100|bg-primary-900/);
  });
});
