import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";
import { createRouter, createMemoryHistory } from "vue-router";
import ReviewView from "../../src/views/ReviewView.vue";

vi.mock("../../src/db", () => ({
  db: {},
  getDeck: vi.fn(),
  getDueCards: vi.fn(),
  updateCardReview: vi.fn(),
}));

vi.mock("../../src/sm2", () => ({
  sm2: vi.fn(() => ({
    easeFactor: 2.5,
    interval: 1,
    repetitions: 1,
    nextReview: new Date(),
  })),
}));

import { getDeck, getDueCards } from "../../src/db";

function createTestRouter(deckId: number) {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: "/", component: { template: "<div />" } },
      { path: "/deck/:id", component: { template: "<div />" } },
      { path: "/review/:deckId", component: ReviewView },
    ],
  });
}

describe("ReviewView", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  // Given due cards have not loaded yet
  // When ReviewView first renders
  // Then it shows the loading skeleton
  it("renders loading skeleton before data loads", async () => {
    vi.mocked(getDeck).mockReturnValue(new Promise(() => {}));
    vi.mocked(getDueCards).mockReturnValue(new Promise(() => {}));

    const router = createTestRouter(1);
    await router.push("/review/1");
    await router.isReady();
    const wrapper = mount(ReviewView, { global: { plugins: [router] } });

    expect(wrapper.find("[data-testid='loading-skeleton']").exists()).toBe(true);
  });

  // Given getDeck returns undefined (nonexistent deck)
  // When ReviewView loads
  // Then it shows the not-found state
  it("renders not-found when deck does not exist", async () => {
    vi.mocked(getDeck).mockResolvedValue(undefined);

    const router = createTestRouter(999);
    await router.push("/review/999");
    await router.isReady();
    const wrapper = mount(ReviewView, { global: { plugins: [router] } });
    await flushPromises();

    expect(wrapper.find("[data-testid='not-found']").exists()).toBe(true);
    expect(wrapper.text()).toContain("Deck not found");
  });

  // Given getDueCards returns an empty array
  // When ReviewView loads
  // Then it shows the empty/done state
  it("renders empty state when no due cards", async () => {
    vi.mocked(getDeck).mockResolvedValue({ id: 1, name: "Spanish", createdAt: new Date() });
    vi.mocked(getDueCards).mockResolvedValue([]);

    const router = createTestRouter(1);
    await router.push("/review/1");
    await router.isReady();
    const wrapper = mount(ReviewView, { global: { plugins: [router] } });
    await flushPromises();

    expect(wrapper.find("[data-testid='review-done']").exists()).toBe(true);
  });

  // Given getDueCards returns cards
  // When ReviewView loads
  // Then it shows the first card
  it("renders first card when due cards exist", async () => {
    vi.mocked(getDeck).mockResolvedValue({ id: 1, name: "Spanish", createdAt: new Date() });
    vi.mocked(getDueCards).mockResolvedValue([
      { id: 1, deckId: 1, front: "hola", back: "hello", easeFactor: 2.5, interval: 1, repetitions: 0, nextReview: new Date() },
    ]);

    const router = createTestRouter(1);
    await router.push("/review/1");
    await router.isReady();
    const wrapper = mount(ReviewView, { global: { plugins: [router] } });
    await flushPromises();

    expect(wrapper.find("[data-testid='card-front']").text()).toBe("hola");
  });

  // Given getDueCards rejects
  // When ReviewView loads
  // Then it shows the error state
  it("renders error state on load failure", async () => {
    vi.mocked(getDeck).mockResolvedValue({ id: 1, name: "Spanish", createdAt: new Date() });
    vi.mocked(getDueCards).mockRejectedValue(new Error("db failure"));

    const router = createTestRouter(1);
    await router.push("/review/1");
    await router.isReady();
    const wrapper = mount(ReviewView, { global: { plugins: [router] } });
    await flushPromises();

    expect(wrapper.find("[data-testid='error-state']").exists()).toBe(true);
    expect(wrapper.text()).toContain("Something went wrong");
  });
});
