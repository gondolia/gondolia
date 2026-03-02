import { test, expect } from '@playwright/test';

test.describe('Search', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/login');
    await page.getByLabel('Email').fill('admin@demo.local');
    await page.getByLabel('Password').fill('admin123');
    await page.getByRole('button', { name: 'Sign in' }).click();
    await expect(page).toHaveURL('/dashboard');
  });

  test('search bar is visible in header', async ({ page }) => {
    const searchInput = page.getByPlaceholder('Produkte suchen...');
    await expect(searchInput).toBeVisible();
  });

  test('search with query navigates to results', async ({ page }) => {
    const searchInput = page.getByPlaceholder('Produkte suchen...');
    await searchInput.fill('Holz');
    await searchInput.press('Enter');

    // Should navigate to products page with search query
    await expect(page).toHaveURL(/\/products\?q=Holz/);
  });

  test('search shows suggestions while typing', async ({ page }) => {
    const searchInput = page.getByPlaceholder('Produkte suchen...');
    await searchInput.fill('Holz');

    // Wait for suggestions dropdown to appear (debounce is 300ms)
    const suggestions = page.locator('.absolute.z-50');
    await expect(suggestions).toBeVisible({ timeout: 5000 });

    // Suggestions should contain product entries with SKU
    await expect(suggestions.locator('button').first()).toBeVisible();
    await expect(suggestions.getByText('SKU:').first()).toBeVisible();
  });

  test('clicking a suggestion navigates to product detail', async ({ page }) => {
    const searchInput = page.getByPlaceholder('Produkte suchen...');
    await searchInput.fill('Holz');

    // Wait for suggestions
    const suggestions = page.locator('.absolute.z-50');
    await expect(suggestions).toBeVisible({ timeout: 5000 });

    // Click first suggestion
    await suggestions.locator('button').first().click();

    // Should navigate to product detail page
    await expect(page).toHaveURL(/\/products\/.+/);
  });

  test('search results page shows product count', async ({ page }) => {
    await page.goto('/search?q=Holz');

    // Should show result count
    await expect(page.getByText(/\d+ Produkt/)).toBeVisible({ timeout: 5000 });
  });

  test('search results page shows product cards', async ({ page }) => {
    await page.goto('/search?q=Holz');

    // Wait for results to load
    await expect(page.getByText(/\d+ Produkt/)).toBeVisible({ timeout: 5000 });

    // Should have product cards (or "Keine Produkte gefunden")
    const hasProducts = await page.locator('[class*="grid"] a, [class*="grid"] [class*="card"]').count();
    const noResults = await page.getByText('Keine Produkte gefunden').count();

    expect(hasProducts > 0 || noResults > 0).toBeTruthy();
  });

  test('search with no results shows empty state', async ({ page }) => {
    await page.goto('/search?q=xyznonexistent12345');

    // Should show "Keine Produkte gefunden"
    await expect(page.getByText('Keine Produkte gefunden')).toBeVisible({ timeout: 5000 });
  });

  test('search results page has filter sidebar', async ({ page }) => {
    await page.goto('/search?q=Holz');

    // Should show filter panel
    await expect(page.getByRole('heading', { name: 'Filter' })).toBeVisible({ timeout: 5000 });
  });

  test('empty search does not navigate', async ({ page }) => {
    const searchInput = page.getByPlaceholder('Produkte suchen...');
    await searchInput.fill('');
    await searchInput.press('Enter');

    // Should stay on dashboard (empty search doesn't submit)
    await expect(page).toHaveURL('/dashboard');
  });
});
