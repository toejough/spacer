import { test, expect } from '@playwright/test';

test('test', async ({ page }) => {
  // go to the dev page
  await page.goto('http://localhost:9201/#/');
  // expect tabs to be visible
  await expect(page.getByRole('tab', { name: 'notes' })).toBeVisible();
  await expect(page.getByRole('tab', { name: 'flashcards' })).toBeVisible();
  // expect input to be visible with button
  await expect(page.getByPlaceholder('Enter a new note here')).toBeVisible();
  await expect(page.getByRole('button')).toBeVisible();
  // expect no notes
  // expect no flashcards
  await page.locator('button').filter({ hasText: 'remove' }).click();
  // add a note
  await page.getByPlaceholder('Enter a new note here').click();
  await page.getByPlaceholder('Enter a new note here').fill('new note with text');
  await page.getByRole('button').click();
  // expect a note
  await expect(page.locator('div').filter({ hasText: /^new note with text$/ }).first()).toBeVisible();
  await expect(page.getByText('drag_indicator')).toBeVisible();
  // click the note
  await page.locator('div').filter({ hasText: /^new note with text$/ }).first().click();
  await page.locator('div').filter({ hasText: /^new note with text$/ }).first().click();
  // expect a toggle button, an edit field with the text
  await page.getByRole('button', { name: 'Toggle flashcard with BOLD' }).click();
  await page.getByText('new note with text').click();
  // select text and click the flashcard toggle
  // expect a flashcard preview
  await page.getByText('new note with ____').click();
  await page.getByText('(text)').click();
  // go to the flashcards tab & expect a card there
  await page.getByRole('tab', { name: 'flashcards' }).locator('div').nth(1).click();
  await expect(page.locator('div').filter({ hasText: /^Next Due: 2025-01-04$/ }).first()).toBeVisible();
  await page.locator('div').filter({ hasText: /^drag_indicator$/ }).click();
  await page.getByText('new note with ____').click();
  await expect(page.getByRole('button', { name: 'Show Answer' })).toBeVisible();
  // click the show answer button
  await page.getByRole('button', { name: 'Show Answer' }).click();
  // expect the answer to show up
  await expect(page.getByText('text')).toBeVisible();
  // click remembered, and the show answer button to come back, answer to disappear, and a new date
  await page.getByRole('button', { name: 'Remembered' }).click();
  await expect(page.getByText('Next Due: 2025-01-')).toBeVisible();
  await expect(page.getByRole('button', { name: 'Show Answer' })).toBeVisible();
  // click show answer and forgot, expect show answer to come back, answer to disappear, and today's date
  await page.getByRole('button', { name: 'Show Answer' }).click();
  await expect(page.getByRole('button', { name: 'Forgot' })).toBeVisible();
  await page.getByRole('button', { name: 'Forgot' }).click();
  await expect(page.getByText('Next Due: 2025-01-')).toBeVisible();
  // go back to notes & add a new note
  await page.getByRole('tab', { name: 'notes' }).click();
  await page.getByPlaceholder('Enter a new note here').click();
  await page.getByPlaceholder('Enter a new note here').fill('another new note');
  await page.getByPlaceholder('Enter a new note here').press('Enter');
  // expect new note
  await expect(page.locator('div').filter({ hasText: /^another new note$/ }).first()).toBeVisible();
  // drag note down & expect new order
  await page.getByText('drag_indicator').nth(1).click();
  await page.getByText('drag_indicator').nth(1).click();
  await page.locator('div').filter({ hasText: /^drag_indicator$/ }).nth(1).click();
  await expect(page.getByText('new note with text')).toBeVisible();
  await page.locator('div').filter({ hasText: /^another new note$/ }).first().click();
  await expect(page.getByRole('listitem').nth(1)).toBeVisible();
});
