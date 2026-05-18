import { expect, test } from '@playwright/test';

test('contact form submits successfully with accessible feedback', async ({ page }) => {
  await page.route('**/api/v1/contact', async (route) => {
    if (route.request().method() === 'POST') {
      await route.fulfill({
        status: 201,
        contentType: 'application/json',
        body: JSON.stringify({ message: 'submission received' }),
      });
      return;
    }

    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify([]),
    });
  });

  await page.goto('/contact');

  await page.getByLabel('Name').fill('Ada Lovelace');
  await page.getByLabel('Email').fill('ada@example.com');
  await page.getByLabel('Message').fill('I would like to discuss a new portfolio project.');
  await page.getByRole('button', { name: /send message/i }).click();

  await expect(page.getByRole('status')).toContainText(/thank you/i);
});