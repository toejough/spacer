# UI Foundations Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Establish design tokens via Tailwind v4 `@theme` and a `ViewState<T>` discriminated union for consistent loading/error/empty/not-found states across all views.

**Architecture:** Design tokens defined in `main.css` via `@theme`, consumed as Tailwind utility classes. `ViewState<T>` type + `useLiveViewState` composable in `src/view-state.ts` subscribes directly to Dexie's `liveQuery` observable with error handling. Each view switches on `state.status` for skeleton/loaded/empty/error/not-found rendering.

**Tech Stack:** Vue 3, TypeScript, Tailwind CSS v4, Dexie (liveQuery), Lucide Vue, Vitest, happy-dom

**Spec:** `docs/plans/2026-03-15-ui-foundations-design.md`

**Note:** All commits should use `/commit` skill which handles the `AI-Used: [claude]` trailer automatically. Inline `git commit` commands in this plan are simplified examples.

---

## File Structure

| File | Responsibility |
|------|---------------|
| `src/main.css` | `@theme` design tokens, `@keyframes skeleton-pulse` |
| `src/view-state.ts` | `ViewState<T>` type, `useLiveViewState` composable |
| `src/App.vue` | Shell with token-based styling (bg, header) |
| `src/views/HomeView.vue` | Deck list with view state pattern + tokens |
| `src/views/DeckView.vue` | Dual view state (deck + cards) + tokens |
| `src/views/ReviewView.vue` | Manual view state ref + tokens |
| `tests/unit/view-state.test.ts` | ViewState type narrowing, useLiveViewState transitions |
| `tests/behavior/home-view.test.ts` | HomeView skeleton/empty/loaded rendering |
| `tests/behavior/deck-view.test.ts` | DeckView skeleton/empty/not-found/loaded rendering |
| `tests/behavior/review-view.test.ts` | ReviewView skeleton/empty/loaded rendering |
| `tests/integration/view-states.test.ts` | Full data path → loaded state, nonexistent deck → not-found |

---

## Chunk 1: Foundation (tokens + ViewState type + composable)

### Task 1: Install dependencies

**Files:**
- Modify: `package.json`

- [ ] **Step 1: Install lucide-vue-next and @vue/test-utils**

Run: `npm install lucide-vue-next && npm install -D @vue/test-utils`

(`@vue/test-utils` was removed in cycle 7 but is needed again for behavioral view tests.)

- [ ] **Step 2: Verify installation**

Run: `node -e "require.resolve('lucide-vue-next'); require.resolve('@vue/test-utils')"`
Expected: exit 0

- [ ] **Step 3: Commit**

```bash
git add package.json package-lock.json
git commit -m "deps: add lucide-vue-next and @vue/test-utils"
```

### Task 2: Design tokens in main.css

**Files:**
- Modify: `src/main.css`

Reference the spec palette table in `docs/plans/2026-03-15-ui-foundations-design.md` § 1.

- [ ] **Step 1: Add @theme block and skeleton animation to main.css**

Replace the entire contents of `src/main.css` with:

```css
@import "tailwindcss";

@theme {
  --color-bg: #e8e4dc;
  --color-surface: #ddd8ce;
  --color-surface-raised: #ece8e0;
  --color-primary: #0891b2;
  --color-primary-hover: #0e7490;
  --color-secondary: #65a30d;
  --color-accent: #86195e;
  --color-accent-hover: #701a50;
  --color-text: #4a4a4a;
  --color-text-muted: #6b6560;
  --color-text-faint: #9b958e;
  --color-border: #d0cabe;
  --radius-sm: 6px;
  --radius-md: 8px;
  --shadow-card: 0 1px 3px rgba(0, 0, 0, 0.08);
}

@keyframes skeleton-pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.animate-skeleton {
  animation: skeleton-pulse 1.5s ease-in-out infinite;
}
```

- [ ] **Step 2: Verify the build still works**

Run: `targ build`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add src/main.css
git commit -m "feat(tokens): add Warm Stone design tokens via Tailwind v4 @theme (#033)"
```

### Task 3: Apply tokens to App.vue shell

**Files:**
- Modify: `src/App.vue`

- [ ] **Step 1: Update App.vue to use token classes**

Replace the current template. Change:
- `bg-gray-50` → `bg-bg` (page background)
- `bg-indigo-600` → `bg-primary` (header)
- Keep `text-white`, `p-4`, `text-xl`, `font-bold`, `max-w-lg`, `mx-auto` (Tailwind defaults, not tokenized)
- Add `text-text` to the outer div for default text color

```vue
<template>
  <div class="min-h-screen bg-bg text-text">
    <header class="bg-primary text-white p-4">
      <router-link to="/" class="text-xl font-bold">Spacer</router-link>
    </header>
    <main class="max-w-lg mx-auto p-4">
      <router-view />
    </main>
  </div>
</template>
```

- [ ] **Step 2: Verify dev server renders with new colors**

Run: `targ dev` (visually check — cream background, cyan header)

- [ ] **Step 3: Commit**

```bash
git add src/App.vue
git commit -m "feat(tokens): apply design tokens to App shell (#033)"
```

### Task 4: ViewState type + useLiveViewState composable — tests first

**Files:**
- Create: `tests/unit/view-state.test.ts`
- Create: `src/view-state.ts`

- [ ] **Step 1: Write failing tests for ViewState type narrowing and useLiveViewState transitions**

Create `tests/unit/view-state.test.ts`:

```ts
import { describe, it, expect } from "vitest";
import { effectScope, nextTick } from "vue";
import { type ViewState, useLiveViewState } from "../../src/view-state";

describe("ViewState type", () => {
  // Given a ViewState in 'loaded' status
  // When narrowed via status check
  // Then data property is accessible
  it("narrows to loaded with data access", () => {
    const state: ViewState<string[]> = { status: "loaded", data: ["a"] };
    if (state.status === "loaded") {
      expect(state.data).toEqual(["a"]);
    }
  });

  // Given a ViewState in 'error' status
  // When narrowed via status check
  // Then message property is accessible
  it("narrows to error with message access", () => {
    const state: ViewState<string[]> = { status: "error", message: "fail" };
    if (state.status === "error") {
      expect(state.message).toBe("fail");
    }
  });
});

describe("useLiveViewState", () => {
  // Given a querier that resolves to undefined with notFoundWhen set
  // When useLiveViewState is called
  // Then it transitions to not-found
  it("transitions from loading to not-found when notFoundWhen matches", async () => {
    const scope = effectScope();
    let state!: ReturnType<typeof useLiveViewState<string | undefined>>;

    scope.run(() => {
      state = useLiveViewState(
        () => Promise.resolve(undefined),
        (d) => false,
        { notFoundWhen: (d) => d === undefined }
      );
    });

    await new Promise((r) => setTimeout(r, 50));

    expect(state.value).toEqual({ status: "not-found" });
    scope.stop();
  });

  // Given a querier that resolves to a non-empty array
  // When useLiveViewState is called with an isEmpty predicate
  // Then it starts as loading, then transitions to loaded
  it("transitions from loading to loaded", async () => {
    const scope = effectScope();
    let state!: ReturnType<typeof useLiveViewState<string[]>>;

    scope.run(() => {
      state = useLiveViewState(
        () => Promise.resolve(["a", "b"]),
        (d) => d.length === 0
      );
    });

    expect(state.value.status).toBe("loading");

    await new Promise((r) => setTimeout(r, 50));

    expect(state.value).toEqual({ status: "loaded", data: ["a", "b"] });
    scope.stop();
  });

  // Given a querier that resolves to an empty array
  // When useLiveViewState is called
  // Then it transitions to empty
  it("transitions from loading to empty", async () => {
    const scope = effectScope();
    let state!: ReturnType<typeof useLiveViewState<string[]>>;

    scope.run(() => {
      state = useLiveViewState(
        () => Promise.resolve([]),
        (d) => d.length === 0
      );
    });

    await new Promise((r) => setTimeout(r, 50));

    expect(state.value).toEqual({ status: "empty" });
    scope.stop();
  });

  // Given a querier that rejects
  // When useLiveViewState is called
  // Then it transitions to error with the message
  it("transitions from loading to error on rejection", async () => {
    const scope = effectScope();
    let state!: ReturnType<typeof useLiveViewState<string[]>>;

    scope.run(() => {
      state = useLiveViewState(
        () => Promise.reject(new Error("db failure")),
        (d) => d.length === 0
      );
    });

    await new Promise((r) => setTimeout(r, 50));

    expect(state.value).toEqual({ status: "error", message: "db failure" });
    scope.stop();
  });

  // Given useLiveViewState is running
  // When the scope is disposed
  // Then the subscription is cleaned up (no error thrown)
  it("cleans up subscription on scope dispose", async () => {
    const scope = effectScope();

    scope.run(() => {
      useLiveViewState(
        () => Promise.resolve(["a"]),
        (d) => d.length === 0
      );
    });

    // Should not throw
    scope.stop();
  });
});
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `npx vitest run tests/unit/view-state.test.ts`
Expected: FAIL — module `../../src/view-state` not found

- [ ] **Step 3: Write minimal implementation**

Create `src/view-state.ts`:

```ts
import { ref, onScopeDispose, type Ref } from "vue";
import { liveQuery } from "dexie";

export type ViewState<T> =
  | { status: "loading" }
  | { status: "loaded"; data: T }
  | { status: "empty" }
  | { status: "error"; message: string }
  | { status: "not-found" };

export interface LiveViewStateOptions<T> {
  notFoundWhen?: (data: T) => boolean;
}

export function useLiveViewState<T>(
  querier: () => Promise<T>,
  isEmpty: (data: T) => boolean,
  options?: LiveViewStateOptions<T>
): Ref<ViewState<T>> {
  const state = ref<ViewState<T>>({ status: "loading" }) as Ref<ViewState<T>>;

  const subscription = liveQuery(querier).subscribe({
    next: (value) => {
      if (options?.notFoundWhen?.(value)) {
        state.value = { status: "not-found" };
      } else {
        state.value = isEmpty(value)
          ? { status: "empty" }
          : { status: "loaded", data: value };
      }
    },
    error: (err) => {
      state.value = {
        status: "error",
        message: err instanceof Error ? err.message : String(err),
      };
    },
  });

  onScopeDispose(() => subscription.unsubscribe());

  return state;
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `npx vitest run tests/unit/view-state.test.ts`
Expected: all 7 tests PASS

- [ ] **Step 5: Run full test suite to check nothing broke**

Run: `targ test`
Expected: all tests PASS

- [ ] **Step 6: Commit**

```bash
git add src/view-state.ts tests/unit/view-state.test.ts
git commit -m "feat(view-state): add ViewState<T> type and useLiveViewState composable (#036)"
```

---

## Chunk 2: HomeView migration

### Task 5: HomeView — behavioral tests first

**Files:**
- Create: `tests/behavior/home-view.test.ts`
- Modify: `src/views/HomeView.vue`

HomeView is the simplest view — single query, two states (empty, loaded). Good proof of the pattern.

- [ ] **Step 1: Write failing behavioral tests for HomeView states**

Create `tests/behavior/home-view.test.ts`:

```ts
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `npx vitest run tests/behavior/home-view.test.ts`
Expected: FAIL — HomeView doesn't have skeleton or new empty state markup yet

- [ ] **Step 3: Rewrite HomeView with tokens + view state pattern**

Replace `src/views/HomeView.vue` with:

```vue
<script setup lang="ts">
import { ref } from "vue";
import { db, getAllDecks, createDeck as createDeckDb } from "../db";
import { useLiveViewState } from "../view-state";
import { Library } from "lucide-vue-next";

const decksState = useLiveViewState(() => getAllDecks(db), (d) => d.length === 0);
const newDeckName = ref("");

async function createDeck() {
  const name = newDeckName.value.trim();
  if (!name) return;
  await createDeckDb(db, name);
  newDeckName.value = "";
}
</script>

<template>
  <div>
    <h2 class="text-lg font-semibold mb-4">Your Decks</h2>

    <form @submit.prevent="createDeck" class="flex gap-2 mb-4">
      <input
        v-model="newDeckName"
        placeholder="New deck name"
        class="flex-1 border border-border rounded-sm bg-surface-raised px-3 py-2 text-text placeholder:text-text-faint"
        data-testid="deck-name-input"
      />
      <button
        type="submit"
        class="bg-primary hover:bg-primary-hover text-white px-4 py-2 rounded-sm"
        data-testid="create-deck-btn"
      >
        Create
      </button>
    </form>

    <!-- Loading skeleton -->
    <div v-if="decksState.status === 'loading'" data-testid="loading-skeleton" class="space-y-2">
      <div v-for="i in 3" :key="i" class="h-12 bg-surface rounded-md animate-skeleton" />
    </div>

    <!-- Empty state -->
    <div
      v-else-if="decksState.status === 'empty'"
      data-testid="empty-state"
      class="text-center py-12"
    >
      <Library class="mx-auto mb-2 text-text-faint" :size="36" />
      <p class="text-text-muted font-medium">No decks yet</p>
      <p class="text-text-faint text-sm">Create your first deck to start studying</p>
    </div>

    <!-- Error state -->
    <div
      v-else-if="decksState.status === 'error'"
      data-testid="error-state"
      class="text-center py-12"
    >
      <p class="text-text-muted font-medium">Something went wrong</p>
      <p class="text-text-faint text-sm">{{ decksState.message }}</p>
    </div>

    <!-- Loaded state -->
    <ul v-else-if="decksState.status === 'loaded'" class="space-y-2">
      <li v-for="deck in decksState.data" :key="deck.id">
        <router-link
          :to="`/deck/${deck.id}`"
          class="block p-3 bg-surface rounded-md shadow-card hover:shadow-md border-l-[3px] border-secondary"
          :data-testid="`deck-${deck.id}`"
        >
          {{ deck.name }}
        </router-link>
      </li>
    </ul>
  </div>
</template>
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `npx vitest run tests/behavior/home-view.test.ts`
Expected: all 3 tests PASS

- [ ] **Step 5: Run full test suite**

Run: `targ test`
Expected: all tests PASS

- [ ] **Step 6: Commit**

```bash
git add src/views/HomeView.vue tests/behavior/home-view.test.ts package.json package-lock.json
git commit -m "feat(home): apply design tokens and view state pattern to HomeView (#033, #036)"
```

---

## Chunk 3: DeckView migration

### Task 6: DeckView — behavioral tests first

**Files:**
- Create: `tests/behavior/deck-view.test.ts`
- Modify: `src/views/DeckView.vue`

DeckView is the most complex — dual view state (deck query + cards query), not-found handling.

- [ ] **Step 1: Write failing behavioral tests for DeckView states**

Create `tests/behavior/deck-view.test.ts`:

```ts
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

function createTestRouter(deckId: number) {
  const router = createRouter({
    history: createMemoryHistory(`/deck/${deckId}`),
    routes: [
      { path: "/", component: { template: "<div />" } },
      { path: "/deck/:id", component: DeckView },
      { path: "/review/:deckId", component: { template: "<div />" } },
    ],
  });
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

    const router = createTestRouter(999);
    await router.isReady();
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

    const router = createTestRouter(1);
    await router.isReady();
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

    const router = createTestRouter(1);
    await router.isReady();
    const wrapper = mount(DeckView, { global: { plugins: [router] } });
    await flushPromises();

    expect(wrapper.find("[data-testid='no-cards']").exists()).toBe(false);
    expect(wrapper.find("[data-testid='card-10']").exists()).toBe(true);
  });

  // Given queries have not resolved yet
  // When DeckView first renders
  // Then it shows the loading skeleton
  it("renders loading skeleton initially", () => {
    vi.mocked(getDeck).mockReturnValue(new Promise(() => {}));
    vi.mocked(getDeckCards).mockReturnValue(new Promise(() => {}));
    vi.mocked(getDueCards).mockReturnValue(new Promise(() => {}));

    const router = createTestRouter(1);
    const wrapper = mount(DeckView, { global: { plugins: [router] } });

    expect(wrapper.find("[data-testid='loading-skeleton']").exists()).toBe(true);
  });
});
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `npx vitest run tests/behavior/deck-view.test.ts`
Expected: FAIL — DeckView doesn't have the new markup

- [ ] **Step 3: Rewrite DeckView with tokens + dual view state**

Replace `src/views/DeckView.vue` with:

```vue
<script setup lang="ts">
import { ref, computed } from "vue";
import { useRoute } from "vue-router";
import { db, getDeck, getDeckCards, getDueCards, createCard as createCardDb } from "../db";
import { useLiveViewState } from "../view-state";
import { useLiveQuery } from "../use-live-query";
import { Layers, Search, AlertTriangle } from "lucide-vue-next";

const route = useRoute();
const deckId = Number(route.params.id);

// Primary state: deck existence (not-found vs loaded)
const deckState = useLiveViewState(
  () => getDeck(db, deckId),
  () => false,
  { notFoundWhen: (d) => d === undefined }
);

// Secondary state: card list (empty vs loaded)
const cardsState = useLiveViewState(
  () => getDeckCards(db, deckId),
  (c) => c.length === 0
);

// Due cards — just a count for the review button, no view state needed
const dueCards = useLiveQuery(() => getDueCards(db, deckId), []);
const dueCount = computed(() => dueCards.value.length);

const front = ref("");
const back = ref("");

async function addCard() {
  const f = front.value.trim();
  const b = back.value.trim();
  if (!f || !b) return;
  await createCardDb(db, deckId, f, b);
  front.value = "";
  back.value = "";
}

</script>

<template>
  <div>
    <!-- Loading skeleton -->
    <div v-if="deckState.status === 'loading'" data-testid="loading-skeleton">
      <div class="flex items-center justify-between mb-4">
        <div class="h-6 w-32 bg-surface rounded-sm animate-skeleton" />
        <div class="h-9 w-24 bg-surface rounded-sm animate-skeleton" />
      </div>
      <div class="space-y-2 mb-4">
        <div class="h-10 bg-surface rounded-sm animate-skeleton" />
        <div class="h-10 bg-surface rounded-sm animate-skeleton" />
        <div class="h-9 w-24 bg-surface rounded-sm animate-skeleton" />
      </div>
      <div class="space-y-2">
        <div v-for="i in 2" :key="i" class="h-16 bg-surface rounded-md animate-skeleton" />
      </div>
    </div>

    <!-- Not found -->
    <div
      v-else-if="deckState.status === 'not-found'"
      data-testid="not-found"
      class="text-center py-12"
    >
      <Search class="mx-auto mb-2 text-text-faint" :size="36" />
      <p class="text-text-muted font-medium">Deck not found</p>
      <p class="text-text-faint text-sm mb-4">This deck may have been deleted</p>
      <router-link to="/" class="text-primary text-sm">&larr; Back to decks</router-link>
    </div>

    <!-- Error -->
    <div
      v-else-if="deckState.status === 'error'"
      data-testid="error-state"
      class="text-center py-12"
    >
      <AlertTriangle class="mx-auto mb-2 text-text-faint" :size="36" />
      <p class="text-text-muted font-medium">Something went wrong</p>
      <p class="text-text-faint text-sm">{{ deckState.message }}</p>
    </div>

    <!-- Deck loaded -->
    <div v-else-if="deckState.status === 'loaded'">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-lg font-semibold">{{ deckState.data.name }}</h2>
        <router-link
          :to="`/review/${deckId}`"
          v-if="dueCount > 0"
          class="bg-accent hover:bg-accent-hover text-white px-4 py-2 rounded-sm text-sm"
          data-testid="start-review-btn"
        >
          Review ({{ dueCount }} due)
        </router-link>
      </div>

      <form @submit.prevent="addCard" class="space-y-2 mb-4">
        <input
          v-model="front"
          placeholder="Front (question)"
          class="w-full border border-border rounded-sm bg-surface-raised px-3 py-2 text-text placeholder:text-text-faint"
          data-testid="card-front-input"
        />
        <input
          v-model="back"
          placeholder="Back (answer)"
          class="w-full border border-border rounded-sm bg-surface-raised px-3 py-2 text-text placeholder:text-text-faint"
          data-testid="card-back-input"
        />
        <button
          type="submit"
          class="bg-primary hover:bg-primary-hover text-white px-4 py-2 rounded-sm"
          data-testid="add-card-btn"
        >
          Add Card
        </button>
      </form>

      <!-- Empty cards -->
      <div
        v-if="cardsState.status === 'empty'"
        data-testid="no-cards"
        class="text-center py-8"
      >
        <Layers class="mx-auto mb-2 text-text-faint" :size="36" />
        <p class="text-text-muted font-medium">No cards yet</p>
        <p class="text-text-faint text-sm">Add cards using the form above</p>
      </div>

      <!-- Card list -->
      <ul v-else-if="cardsState.status === 'loaded'" class="space-y-2">
        <li
          v-for="card in cardsState.data"
          :key="card.id"
          class="p-3 bg-surface rounded-md shadow-card"
          :data-testid="`card-${card.id}`"
        >
          <div class="font-medium">{{ card.front }}</div>
          <div class="text-text-muted text-sm">{{ card.back }}</div>
        </li>
      </ul>

      <router-link to="/" class="inline-block mt-4 text-primary text-sm">
        &larr; Back to decks
      </router-link>
    </div>
  </div>
</template>
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `npx vitest run tests/behavior/deck-view.test.ts`
Expected: all 4 tests PASS

- [ ] **Step 5: Run full test suite**

Run: `targ test`
Expected: all tests PASS

- [ ] **Step 6: Commit**

```bash
git add src/views/DeckView.vue tests/behavior/deck-view.test.ts
git commit -m "feat(deck): apply design tokens and dual view state to DeckView (#033, #036)"
```

---

## Chunk 4: ReviewView migration

### Task 7: ReviewView — behavioral tests first

**Files:**
- Create: `tests/behavior/review-view.test.ts`
- Modify: `src/views/ReviewView.vue`

ReviewView uses manual `ViewState` ref (no `useLiveViewState`) since it loads imperatively.

- [ ] **Step 1: Write failing behavioral tests for ReviewView states**

Create `tests/behavior/review-view.test.ts`:

```ts
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
  const router = createRouter({
    history: createMemoryHistory(`/review/${deckId}`),
    routes: [
      { path: "/", component: { template: "<div />" } },
      { path: "/deck/:id", component: { template: "<div />" } },
      { path: "/review/:deckId", component: ReviewView },
    ],
  });
  return router;
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
    await router.isReady();
    const wrapper = mount(ReviewView, { global: { plugins: [router] } });
    await flushPromises();

    expect(wrapper.find("[data-testid='error-state']").exists()).toBe(true);
    expect(wrapper.text()).toContain("Something went wrong");
  });
});
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `npx vitest run tests/behavior/review-view.test.ts`
Expected: FAIL — ReviewView doesn't have skeleton or error markup

- [ ] **Step 3: Rewrite ReviewView with tokens + manual view state**

Replace `src/views/ReviewView.vue` with:

```vue
<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { db, getDeck, getDueCards, updateCardReview, type Card } from "../db";
import { sm2 } from "../sm2";
import { type ViewState } from "../view-state";
import { AlertTriangle, Search } from "lucide-vue-next";

const route = useRoute();
const router = useRouter();
const deckId = Number(route.params.deckId);

const loadState = ref<ViewState<Card[]>>({ status: "loading" });
const currentIndex = ref(0);
const flipped = ref(false);
const done = ref(false);

const cards = computed(() =>
  loadState.value.status === "loaded" ? loadState.value.data : []
);
const current = computed(() => cards.value[currentIndex.value]);

async function load() {
  try {
    const deck = await getDeck(db, deckId);
    if (!deck) {
      loadState.value = { status: "not-found" };
      return;
    }
    const dueCards = await getDueCards(db, deckId);
    if (dueCards.length === 0) {
      loadState.value = { status: "empty" };
      done.value = true;
    } else {
      loadState.value = { status: "loaded", data: dueCards };
    }
  } catch (err) {
    loadState.value = {
      status: "error",
      message: err instanceof Error ? err.message : String(err),
    };
  }
}

function flip() {
  flipped.value = true;
}

async function rate(quality: number) {
  const card = current.value;
  if (!card) return;

  const result = sm2(
    {
      easeFactor: card.easeFactor,
      interval: card.interval,
      repetitions: card.repetitions,
    },
    quality
  );

  await updateCardReview(db, card.id, result);

  flipped.value = false;
  if (currentIndex.value < cards.value.length - 1) {
    currentIndex.value++;
  } else {
    done.value = true;
  }
}

onMounted(load);
</script>

<template>
  <div>
    <!-- Loading skeleton -->
    <div v-if="loadState.status === 'loading'" data-testid="loading-skeleton">
      <div class="h-4 w-24 bg-surface rounded-sm animate-skeleton mb-2" />
      <div class="bg-surface rounded-md min-h-[200px] animate-skeleton mb-4" />
      <div class="flex justify-center">
        <div class="h-10 w-32 bg-surface rounded-sm animate-skeleton" />
      </div>
    </div>

    <!-- Not found -->
    <div
      v-else-if="loadState.status === 'not-found'"
      data-testid="not-found"
      class="text-center py-12"
    >
      <Search class="mx-auto mb-2 text-text-faint" :size="36" />
      <p class="text-text-muted font-medium">Deck not found</p>
      <p class="text-text-faint text-sm mb-4">This deck may have been deleted</p>
      <router-link to="/" class="text-primary text-sm">&larr; Back to decks</router-link>
    </div>

    <!-- Error state -->
    <div
      v-else-if="loadState.status === 'error'"
      data-testid="error-state"
      class="text-center py-12"
    >
      <AlertTriangle class="mx-auto mb-2 text-text-faint" :size="36" />
      <p class="text-text-muted font-medium">Something went wrong</p>
      <p class="text-text-faint text-sm mb-4">{{ loadState.message }}</p>
      <router-link :to="`/deck/${deckId}`" class="text-primary text-sm">
        &larr; Back to deck
      </router-link>
    </div>

    <!-- Review complete / empty -->
    <div v-else-if="done" data-testid="review-done">
      <h2 class="text-lg font-semibold mb-2">Review Complete</h2>
      <p class="text-text-muted mb-4">
        Reviewed {{ cards.length }} card{{ cards.length === 1 ? "" : "s" }}.
      </p>
      <button
        @click="router.push(`/deck/${deckId}`)"
        class="bg-primary hover:bg-primary-hover text-white px-4 py-2 rounded-sm"
      >
        Back to Deck
      </button>
    </div>

    <!-- Active review -->
    <div v-else-if="current">
      <p class="text-sm text-text-muted mb-2" data-testid="review-progress">
        Card {{ currentIndex + 1 }} of {{ cards.length }}
      </p>

      <div class="bg-surface rounded-md shadow-card p-6 min-h-[200px] flex items-center justify-center">
        <div class="text-center">
          <p class="text-xl" data-testid="card-front">{{ current.front }}</p>
          <p
            v-if="flipped"
            class="mt-4 text-text-muted border-t border-border pt-4"
            data-testid="card-back"
          >
            {{ current.back }}
          </p>
        </div>
      </div>

      <div v-if="!flipped" class="mt-4 text-center">
        <button
          @click="flip"
          class="bg-accent hover:bg-accent-hover text-white px-6 py-2 rounded-sm"
          data-testid="flip-btn"
        >
          Show Answer
        </button>
      </div>

      <div v-else class="mt-4 flex justify-center gap-2">
        <button
          v-for="q in [1, 2, 3, 4, 5]"
          :key="q"
          @click="rate(q)"
          class="px-4 py-2 rounded-sm border border-border hover:bg-surface-raised"
          :data-testid="`rate-${q}`"
        >
          {{ q }}
        </button>
      </div>
    </div>
  </div>
</template>
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `npx vitest run tests/behavior/review-view.test.ts`
Expected: all 5 tests PASS

- [ ] **Step 5: Run full test suite**

Run: `targ test`
Expected: all tests PASS

- [ ] **Step 6: Commit**

```bash
git add src/views/ReviewView.vue tests/behavior/review-view.test.ts
git commit -m "feat(review): apply design tokens and manual view state to ReviewView (#033, #036)"
```

---

## Chunk 5: Integration tests + cleanup

### Task 8: Integration tests for view state transitions

**Files:**
- Create: `tests/integration/view-states.test.ts`

- [ ] **Step 1: Write integration tests using real Dexie (fake-indexeddb)**

Create `tests/integration/view-states.test.ts`:

```ts
import { describe, it, expect, beforeEach } from "vitest";
import { effectScope } from "vue";
import { SpacerDB, createDeck, createCard, getAllDecks, getDeck, getDeckCards } from "../../src/db";
import { useLiveViewState } from "../../src/view-state";

describe("useLiveViewState with real Dexie", () => {
  let db: SpacerDB;

  beforeEach(() => {
    db = new SpacerDB("test-view-states-" + Date.now());
  });

  // Given no decks exist
  // When useLiveViewState subscribes to getAllDecks
  // Then it transitions to empty
  it("empty database produces empty view state", async () => {
    const scope = effectScope();
    let state!: ReturnType<typeof useLiveViewState>;

    scope.run(() => {
      state = useLiveViewState(
        () => getAllDecks(db),
        (d) => d.length === 0
      );
    });

    await new Promise((r) => setTimeout(r, 50));
    expect(state.value.status).toBe("empty");
    scope.stop();
  });

  // Given a deck exists with cards
  // When useLiveViewState subscribes to getDeckCards
  // Then it transitions to loaded with the card data
  it("deck with cards produces loaded view state", async () => {
    const deckId = await createDeck(db, "Spanish");
    await createCard(db, deckId, "hola", "hello");

    const scope = effectScope();
    let state!: ReturnType<typeof useLiveViewState>;

    scope.run(() => {
      state = useLiveViewState(
        () => getDeckCards(db, deckId),
        (c) => c.length === 0
      );
    });

    await new Promise((r) => setTimeout(r, 50));
    expect(state.value.status).toBe("loaded");
    scope.stop();
  });

  // Given a nonexistent deck ID
  // When useLiveViewState subscribes to getDeck with notFoundWhen
  // Then it transitions to not-found
  it("nonexistent deck produces not-found view state", async () => {
    const scope = effectScope();
    let state!: ReturnType<typeof useLiveViewState>;

    scope.run(() => {
      state = useLiveViewState(
        () => getDeck(db, 99999),
        () => false,
        { notFoundWhen: (d) => d === undefined }
      );
    });

    await new Promise((r) => setTimeout(r, 50));
    expect(state.value.status).toBe("not-found");
    scope.stop();
  });
});
```

- [ ] **Step 2: Run integration tests**

Run: `npx vitest run tests/integration/view-states.test.ts`
Expected: all 3 tests PASS

- [ ] **Step 3: Run full test suite to confirm nothing is broken**

Run: `targ test`
Expected: all tests PASS

- [ ] **Step 4: Commit**

```bash
git add tests/integration/view-states.test.ts
git commit -m "test: add integration tests for view state data paths (#036)"
```

### Task 9: Visual verification + final commit

- [ ] **Step 1: Start dev server and visually verify all views**

Run: `targ dev`

Check each view:
- Home (no decks): cream background, cyan header, centered Library icon + "No decks yet"
- Home (with decks): deck list items with green left border, surface background, card shadow
- Deck (nonexistent ID `/deck/999`): Search icon + "Deck not found" + cyan back link
- Deck (with cards): form inputs with raised surface bg, border, card list with surface bg
- Deck (no cards): Layers icon + "No cards yet"
- Review: magenta "Show Answer" button, surface card, border divider on flip, rating buttons with border

- [ ] **Step 2: Run `targ build` to verify production build**

Run: `targ build`
Expected: no errors, no type errors

- [ ] **Step 3: Run `targ check` for full quality gate**

Run: `targ check`
Expected: type-check passes, all tests pass
