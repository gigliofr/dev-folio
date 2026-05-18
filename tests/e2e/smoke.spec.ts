import { expect, test } from '@playwright/test';

test('public pages render primary navigation', async ({ page }) => {
  await page.goto('/');

  await expect(page.getByRole('heading', { name: /a portfolio system designed to be rebranded/i })).toBeVisible();
  await expect(page.getByRole('link', { name: 'View projects' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Contact me' })).toBeVisible();
  await expect(page.getByRole('navigation', { name: /primary navigation/i }).getByRole('link', { name: 'Projects' })).toBeVisible();
  await expect(page.getByRole('navigation', { name: /primary navigation/i }).getByRole('link', { name: 'Blog' })).toBeVisible();
  await expect(page.getByRole('navigation', { name: /primary navigation/i }).getByRole('link', { name: 'Contact' })).toBeVisible();
});

test('projects page supports searching and filtering', async ({ page }) => {
  await page.goto('/projects');

  await expect(page.getByRole('heading', { level: 2, name: /a structured project archive ready for filters/i })).toBeVisible();
  await page.getByPlaceholder('Search projects, technologies, or years').fill('content');
  await expect(page.getByText('Content Hub')).toBeVisible();
});