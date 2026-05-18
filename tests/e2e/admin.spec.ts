import { expect, test } from '@playwright/test';

test('admin login reveals contacts tab and data', async ({ page }) => {
  await page.route('**/api/v1/admin/login', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ token: 'test-token' }),
    });
  });

  await page.route('**/api/v1/admin/me', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ description: 'admin:test' }),
    });
  });

  await page.route('**/api/v1/contact', async (route) => {
    if (route.request().method() === 'GET') {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          submissions: [
            {
              name: 'Ada Lovelace',
              email: 'ada@example.com',
              message: 'Portfolio feedback',
              createdAt: new Date().toISOString(),
            },
          ],
          count: 1,
        }),
      });
      return;
    }

    await route.fulfill({
      status: 201,
      contentType: 'application/json',
      body: JSON.stringify({ message: 'submission received' }),
    });
  });

  await page.goto('/admin');

  await page.getByPlaceholder('Username').fill('admin');
  await page.getByPlaceholder('Password').fill('password');
  await page.getByRole('button', { name: /sign in/i }).click();

  await expect(page.getByText(/signed in as/i)).toBeVisible();
  await page.getByRole('button', { name: /^Contacts \(/ }).click();
  await expect(page.getByText('Ada Lovelace')).toBeVisible();
});