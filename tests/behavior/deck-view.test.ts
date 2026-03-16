import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";
import { createRouter, createMemoryHistory } from "vue-router";
import DeckView from "../../src/views/DeckView.vue";

vi.mock("../../src/db", () => ({
  db: {},
  getDeck: vi.fn(),
  getDeckCards: vi.fn(),
  getDueCards: vi.fn(),
  createCard: vi.fn(),
}));

vi.mock("dexie", async () => {
  const actual = await vi.importActual("dexie");
  return {
    ...actual,
    liveQuery: (querier: () => Promise<unknown>) => ({
      subscribe: (observer: { next?: (v: unknown) => void; error?: (e: unknown) => void }) => {
        querier().then(
          (v) => observer.next?.(v),
          (e) => observer.error?.(e)
        );
        return { unsubscribe: () => {} };
      },
    }),
  };
});

import { getDeck, getDeckCards, getDueCards } from "../../src/db";

async function createTestRouter(deckId: number) {
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: "/", component: { template: "<div />" } },
      { path: "/deck/:id", component: DeckView },
      { path: "/review/:deckId", component: { template: "<div />" } },
    ],
  });
  await router.push(`/deck/${deckId}`);
  return router;
}

describe("DeckView", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  // Given the deck query returns undefined
  // When DeckView renders
  // Then it shows the not-found state
  it("renders not-found when deck does not exist", async () => {
    vi.mocked(getDeck).mockResolvedValue(undefined);
    vi.mocked(getDeckCards).mockResolvedValue([]);
    vi.mocked(getDueCards).mockResolvedValue([]);

    const router = await createTestRouter(999);
    const wrapper = mount(DeckView, { global: { plugins: [router] } });
    await flushPromises();

    expect(wrapper.find("[data-testid='not-found']").exists()).toBe(true);
    expect(wrapper.text()).toContain("Deck not found");
  });

  // Given the deck exists but has no cards
  // When DeckView renders
  // Then it shows the deck header and empty-cards state
  it("renders empty cards state when deck has no cards", async () => {
    vi.mocked(getDeck).mockResolvedValue({ id: 1, name: "Spanish", createdAt: new Date() });
    vi.mocked(getDeckCards).mockResolvedValue([]);
    vi.mocked(getDueCards).mockResolvedValue([]);

    const router = await createTestRouter(1);
    const wrapper = mount(DeckView, { global: { plugins: [router] } });
    await flushPromises();

    expect(wrapper.find("[data-testid='not-found']").exists()).toBe(false);
    expect(wrapper.find("[data-testid='no-cards']").exists()).toBe(true);
    expect(wrapper.text()).toContain("Spanish");
  });

  // Given the deck exists with cards
  // When DeckView renders
  // Then it shows the card list
  it("renders loaded state with cards", async () => {
    vi.mocked(getDeck).mockResolvedValue({ id: 1, name: "Spanish", createdAt: new Date() });
    vi.mocked(getDeckCards).mockResolvedValue([
      { id: 10, deckId: 1, front: "hola", back: "hello", easeFactor: 2.5, interval: 1, repetitions: 0, nextReview: new Date() },
    ]);
    vi.mocked(getDueCards).mockResolvedValue([]);

    const router = await createTestRouter(1);
    const wrapper = mount(DeckView, { global: { plugins: [router] } });
    await flushPromises();

    expect(wrapper.find("[data-testid='no-cards']").exists()).toBe(false);
    expect(wrapper.find("[data-testid='card-10']").exists()).toBe(true);
  });

  // Given queries have not resolved yet
  // When DeckView first renders
  // Then it shows the loading skeleton
  it("renders loading skeleton initially", async () => {
    vi.mocked(getDeck).mockReturnValue(new Promise(() => {}));
    vi.mocked(getDeckCards).mockReturnValue(new Promise(() => {}));
    vi.mocked(getDueCards).mockReturnValue(new Promise(() => {}));

    const router = await createTestRouter(1);
    const wrapper = mount(DeckView, { global: { plugins: [router] } });

    expect(wrapper.find("[data-testid='loading-skeleton']").exists()).toBe(true);
  });
});
