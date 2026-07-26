const { test, expect } = require('@playwright/test');
const AxeBuilder = require('@axe-core/playwright').default;

test('accessibility: todos tab should have no detectable a11y violations', async ({ page }) => {
  await page.goto('http://localhost:9201/#/');
  await page.getByRole('tab', { name: 'todos' }).click();
  const results = await new AxeBuilder({ page }).include('body').analyze();
  expect(results.violations).toEqual([]);
});
