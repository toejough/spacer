import { test, expect } from '@playwright/test';

test('mobile equal-access: Done, Abandon, Reopen, visual states', async ({ page }) => {
  await page.goto('http://localhost:9201/#/');
  // open todos tab
  await page.getByRole('tab', { name: 'todos' }).click();
  // expect two open items
  await expect(page.getByRole('listitem').filter({ hasText: 'Write project spec' })).toBeVisible();
  // Click Done on first item
  const first = page.getByRole('listitem').filter({ hasText: 'Write project spec' }).first();
  await first.getByRole('button', { name: 'Done' }).click();
  // Expect it to have done styling (background change)
  const bgDone = await first.evaluate((el) => getComputedStyle(el).backgroundColor);
  // Click Reopen on closed item (Publish release notes)
  const closed = page.getByRole('listitem').filter({ hasText: 'Publish release notes' }).first();
  await closed.getByRole('button', { name: 'Reopen' }).click();
  const bgReopen = await closed.evaluate((el) => getComputedStyle(el).backgroundColor);
  // Abandon the second open item
  const second = page.getByRole('listitem').filter({ hasText: 'Refactor API' }).first();
  await second.getByRole('button', { name: 'Abandon' }).click();
  // confirm dialog — accept
  await page.on('dialog', async dialog => { await dialog.accept(); });
  // check visual difference: abandoned vs done
  const abandoned = page.getByRole('listitem').filter({ hasText: 'Experimental feature X' }).first();
  const bgAbandoned = await abandoned.evaluate((el) => getComputedStyle(el).backgroundColor);
  expect(bgDone).not.toEqual(bgAbandoned);
  expect(bgReopen).not.toEqual(bgAbandoned);
});
