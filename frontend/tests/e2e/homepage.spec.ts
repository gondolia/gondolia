import { test, expect } from '@playwright/test';

test.describe('Homepage', () => {
  test('homepage loads and redirects to login when not authenticated', async ({ page }) => {
    await page.goto('/');

    // Should redirect to login page when not authenticated
    await expect(page).toHaveURL('/login');
  });

  test('login page displays correctly', async ({ page }) => {
    await page.goto('/login');

    // Check for login form elements
    await expect(page.getByRole('heading', { name: 'Login' })).toBeVisible();
    await expect(page.getByLabel('Email')).toBeVisible();
    await expect(page.getByLabel('Password')).toBeVisible();
    await expect(page.getByRole('button', { name: 'Sign in' })).toBeVisible();
  });

  test('dashboard shows products after login', async ({ page }) => {
    // First login
    await page.goto('/login');
    await page.getByLabel('Email').fill('admin@demo.local');
    await page.getByLabel('Password').fill('admin123');
    await page.getByRole('button', { name: 'Sign in' }).click();

    // Wait for redirect to dashboard
    await expect(page).toHaveURL('/dashboard');

    // Navigate to products page
    await page.getByRole('link', { name: 'Produkte' }).click();
    await expect(page).toHaveURL('/products');

    // Wait for products to load (check for heading)
    await expect(page.getByRole('heading', { name: 'Produkte' })).toBeVisible();
  });
});
