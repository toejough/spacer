import { test, expect } from "@playwright/test";

test.describe("PWA offline support", () => {
  test("navigating to a deck works offline", async ({ page, context }) => {
    // Load the app and create a deck so there's something to click
    await page.goto("/");
    await page.getByTestId("deck-name-input").fill("Offline Test Deck");
    await page.getByTestId("create-deck-btn").click();
    await page.getByTestId("deck-1").waitFor();

    // Wait for the service worker to be installed and activated
    await page.evaluate(async () => {
      const reg = await navigator.serviceWorker.ready;
      if (reg.active?.state !== "activated") {
        await new Promise<void>((resolve) => {
          reg.active?.addEventListener("statechange", () => {
            if (reg.active?.state === "activated") resolve();
          });
          // Already activated
          if (reg.active?.state === "activated") resolve();
        });
      }
    });

    // Go offline
    await context.setOffline(true);

    // Click the deck — this should navigate to /deck/1
    await page.getByTestId("deck-1").click();

    // DeckView should render with the deck name
    await expect(page.locator("text=Offline Test Deck")).toBeVisible();
    await expect(page.getByTestId("card-front-input")).toBeVisible();
  });

  test("refreshing the page works offline", async ({ page, context }) => {
    // Load the app
    await page.goto("/");
    await expect(page.locator("text=Spacer")).toBeVisible();

    // Wait for the service worker to be installed and activated
    await page.evaluate(async () => {
      const reg = await navigator.serviceWorker.ready;
      if (reg.active?.state !== "activated") {
        await new Promise<void>((resolve) => {
          reg.active?.addEventListener("statechange", () => {
            if (reg.active?.state === "activated") resolve();
          });
          if (reg.active?.state === "activated") resolve();
        });
      }
    });

    // Go offline
    await context.setOffline(true);

    // Refresh the page
    await page.reload();

    // App should still render
    await expect(page.locator("text=Spacer")).toBeVisible();
    await expect(page.locator("text=Your Decks")).toBeVisible();
  });
});
