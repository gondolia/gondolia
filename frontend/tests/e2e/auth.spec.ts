import { test, expect } from '@playwright/test';

test.describe('Authentication', () => {
  test('login page is reachable', async ({ page }) => {
    await page.goto('/login');

    // Check login page loads
    await expect(page).toHaveURL('/login');
    await expect(page.getByRole('heading', { name: 'Login' })).toBeVisible();
    await expect(page.getByText('Sign in to your account')).toBeVisible();
  });

  test('successful login with valid credentials', async ({ page }) => {
    await page.goto('/login');

    // Fill in credentials
    await page.getByLabel('Email').fill('admin@demo.local');
    await page.getByLabel('Password').fill('admin123');

    // Submit login form
    await page.getByRole('button', { name: 'Sign in' }).click();

    // Should redirect to dashboard after successful login
    await expect(page).toHaveURL('/dashboard');

    // Verify dashboard content is visible
    await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible();
    await expect(page.getByRole('heading', { name: 'User Information' })).toBeVisible();
    await expect(page.getByRole('heading', { name: 'Company Information' })).toBeVisible();
  });

  test('login shows error with invalid credentials', async ({ page }) => {
    await page.goto('/login');

    // Fill in wrong credentials
    await page.getByLabel('Email').fill('wrong@example.com');
    await page.getByLabel('Password').fill('wrongpassword');

    // Submit login form
    await page.getByRole('button', { name: 'Sign in' }).click();

    // Should stay on login page and show error
    await expect(page).toHaveURL('/login');
    // Wait for error message to appear
    await expect(page.locator('.bg-red-50, .bg-red-900\\/30')).toBeVisible({ timeout: 10000 });
  });

  test('logout functionality works', async ({ page }) => {
    // First login
    await page.goto('/login');
    await page.getByLabel('Email').fill('admin@demo.local');
    await page.getByLabel('Password').fill('admin123');
    await page.getByRole('button', { name: 'Sign in' }).click();

    // Wait for dashboard
    await expect(page).toHaveURL('/dashboard');

    // Click logout button
    await page.getByRole('button', { name: 'Logout' }).click();

    // Should redirect to login page
    await expect(page).toHaveURL('/login');
  });

  test('protected routes redirect to login when not authenticated', async ({ page }) => {
    // Try to access protected route
    await page.goto('/dashboard');

    // Should redirect to login
    await expect(page).toHaveURL('/login');
  });
});
