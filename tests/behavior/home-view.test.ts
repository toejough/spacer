import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount, flushPromises } from "@vue/test-utils";
import { createRouter, createMemoryHistory } from "vue-router";
import HomeView from "../../src/views/HomeView.vue";

// Mock db module — control what queries return
vi.mock("../../src/db", () => ({
  db: {},
  getAllDecks: vi.fn(),
  createDeck: vi.fn(),
}));

// Mock dexie's liveQuery to emit mock data synchronously
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

import { getAllDecks } from "../../src/db";

function createTestRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: "/", component: HomeView },
      { path: "/deck/:id", component: { template: "<div />" } },
    ],
  });
}

describe("HomeView", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  // Given the query returns an empty array
  // When HomeView renders
  // Then it shows the empty state with icon and message
  it("renders empty state when no decks exist", async () => {
    vi.mocked(getAllDecks).mockResolvedValue([]);
    const router = createTestRouter();
    const wrapper = mount(HomeView, { global: { plugins: [router] } });
    await flushPromises();

    const empty = wrapper.find("[data-testid='empty-state']");
    expect(empty.exists()).toBe(true);
    expect(empty.text()).toContain("No decks yet");
  });

  // Given the query returns decks
  // When HomeView renders
  // Then it shows the deck list, not the empty state
  it("renders loaded state with deck list", async () => {
    vi.mocked(getAllDecks).mockResolvedValue([
      { id: 1, name: "Spanish", createdAt: new Date() },
      { id: 2, name: "Biology", createdAt: new Date() },
    ]);
    const router = createTestRouter();
    const wrapper = mount(HomeView, { global: { plugins: [router] } });
    await flushPromises();

    expect(wrapper.find("[data-testid='empty-state']").exists()).toBe(false);
    expect(wrapper.find("[data-testid='deck-1']").exists()).toBe(true);
    expect(wrapper.find("[data-testid='deck-2']").exists()).toBe(true);
  });

  // Given the view is in its initial state (before query resolves)
  // When HomeView first renders
  // Then it shows the loading skeleton
  it("renders loading skeleton initially", () => {
    vi.mocked(getAllDecks).mockReturnValue(new Promise(() => {})); // never resolves
    const router = createTestRouter();
    const wrapper = mount(HomeView, { global: { plugins: [router] } });

    expect(wrapper.find("[data-testid='loading-skeleton']").exists()).toBe(true);
  });
});
