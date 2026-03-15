# Spacer PWA Bootstrap Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Bootstrap the Spacer PWA with Vue 3 + Vite + Tailwind + Vitest + Dexie + PWA manifest so every subsequent 5-minute increment can deliver a real feature.

**Architecture:** Single-page app with Vue Router (3 route stubs), Pinia stores, and Dexie.js for IndexedDB persistence. SM-2 algorithm lives as a pure function in `src/lib/`. PWA via vite-plugin-pwa with `registerType: 'autoUpdate'` (auto-injects service worker registration — no manual registration code needed in `main.ts`).

**Tech Stack:** Vue 3, TypeScript, Vite, Tailwind CSS v4, Pinia, Dexie.js, Vitest, @vue/test-utils, vite-plugin-pwa

**Spec:** `docs/superpowers/specs/2026-03-14-spacer-bootstrap-design.md`

**Note on SM-2:** The spec simplifies the description of quality < 3 behavior. The standard SM-2 algorithm adjusts `easeFactor` on every review regardless of quality. The implementation follows the standard algorithm.

---

## Chunk 1: Full Bootstrap

### Task 1: Scaffold Vite + Vue 3 + TypeScript project

**Files:**
- Create: `package.json`, `tsconfig.json`, `tsconfig.app.json`, `tsconfig.node.json`, `vite.config.ts`, `index.html`, `src/main.ts`, `src/App.vue`, `env.d.ts`

- [ ] **Step 1: Create the Vite project**

Run from repo root:
```bash
npm create vite@latest . -- --template vue-ts
```

If it asks about non-empty directory, confirm. This generates the standard Vue + TS scaffold.

- [ ] **Step 2: Install dependencies**

```bash
npm install
```

- [ ] **Step 3: Verify the app runs**

```bash
npx vite --open=false &
sleep 3
curl -s http://localhost:5173 | head -5
kill %1
```

Expected: HTML with `<div id="app">` present.

- [ ] **Step 4: Commit**

```bash
git add package.json package-lock.json tsconfig.json tsconfig.app.json tsconfig.node.json vite.config.ts index.html env.d.ts src/ public/
git commit -m "feat: scaffold Vite + Vue 3 + TypeScript project

AI-Used: [claude]"
```

---

### Task 2: Configure Tailwind CSS

**Files:**
- Create: `src/assets/main.css`
- Modify: `src/main.ts`, `src/App.vue`, `vite.config.ts`

Note: Using Tailwind v4 with `@tailwindcss/vite` plugin — no `tailwind.config.js` or `postcss.config.js` needed.

- [ ] **Step 1: Install Tailwind and dependencies**

```bash
npm install -D tailwindcss @tailwindcss/vite
```

- [ ] **Step 2: Create `src/assets/main.css`**

```css
@import "tailwindcss";
```

- [ ] **Step 3: Update `src/main.ts` to import the CSS**

Ensure `src/main.ts` imports `./assets/main.css`. Remove the default Vite style import if present.

```typescript
import { createApp } from 'vue'
import App from './App.vue'
import './assets/main.css'

createApp(App).mount('#app')
```

- [ ] **Step 4: Update `vite.config.ts` to use Tailwind plugin**

```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [vue(), tailwindcss()],
})
```

- [ ] **Step 5: Update `src/App.vue` to use a Tailwind class**

```vue
<template>
  <div class="min-h-screen bg-gray-50">
    <h1 class="text-2xl font-bold p-4">Spacer</h1>
  </div>
</template>
```

- [ ] **Step 6: Verify Tailwind works**

```bash
npx vite build 2>&1 | tail -5
```

Expected: Build succeeds without errors.

- [ ] **Step 7: Commit**

```bash
git add src/assets/main.css src/main.ts src/App.vue vite.config.ts package.json package-lock.json
git commit -m "feat: configure Tailwind CSS v4

AI-Used: [claude]"
```

---

### Task 3: Configure Vitest with smoke test

**Files:**
- Create: `vitest.config.ts`, `src/__tests__/smoke.test.ts`
- Modify: `package.json` (add test script)

- [ ] **Step 1: Install Vitest and test utils**

```bash
npm install -D vitest @vue/test-utils happy-dom
```

- [ ] **Step 2: Create `vitest.config.ts`**

```typescript
import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  test: {
    environment: 'happy-dom',
  },
})
```

- [ ] **Step 3: Write failing smoke test**

The test asserts `router-view` is present — which it won't be until Task 4 adds the router. This gives us a genuine red phase.

Create `src/__tests__/smoke.test.ts`:

```typescript
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import App from '../App.vue'

describe('App', () => {
  it('renders the app title', () => {
    const wrapper = mount(App)
    expect(wrapper.text()).toContain('Spacer')
  })

  it('has a router-view for page content', () => {
    const wrapper = mount(App)
    expect(wrapper.find('router-view').exists() || wrapper.html().includes('router-view')).toBe(true)
  })
})
```

- [ ] **Step 4: Run tests to verify red**

```bash
npx vitest run
```

Expected: First test passes (App.vue has "Spacer"), second test FAILS (no router-view yet). This is the correct red state.

- [ ] **Step 5: Add test scripts to package.json**

Add to `scripts`: `"test": "vitest run"`, `"test:watch": "vitest"`

- [ ] **Step 6: Commit (with one red test — intentional)**

```bash
git add vitest.config.ts src/__tests__/smoke.test.ts package.json
git commit -m "test: configure Vitest with smoke tests (router-view test intentionally red)

AI-Used: [claude]"
```

---

### Task 4: Set up Vue Router with route stubs

**Files:**
- Create: `src/router/index.ts`, `src/views/HomeView.vue`, `src/views/DeckView.vue`, `src/views/ReviewView.vue`
- Create: `src/__tests__/router.test.ts`
- Modify: `src/main.ts`, `src/App.vue`

- [ ] **Step 1: Install Vue Router**

```bash
npm install vue-router
```

- [ ] **Step 2: Write failing router test**

Create `src/__tests__/router.test.ts`:

```typescript
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'

describe('Router', () => {
  const routes = [
    { path: '/', name: 'home', component: { template: '<div>Home</div>' } },
    { path: '/deck/:id', name: 'deck', component: { template: '<div>Deck</div>' } },
    { path: '/review/:deckId?', name: 'review', component: { template: '<div>Review</div>' } },
  ]

  it('navigates to home route', async () => {
    const router = createRouter({ history: createMemoryHistory(), routes })
    router.push('/')
    await router.isReady()
    expect(router.currentRoute.value.name).toBe('home')
  })

  it('navigates to deck route with id param', async () => {
    const router = createRouter({ history: createMemoryHistory(), routes })
    router.push('/deck/abc123')
    await router.isReady()
    expect(router.currentRoute.value.name).toBe('deck')
    expect(router.currentRoute.value.params.id).toBe('abc123')
  })

  it('navigates to review route', async () => {
    const router = createRouter({ history: createMemoryHistory(), routes })
    router.push('/review')
    await router.isReady()
    expect(router.currentRoute.value.name).toBe('review')
  })
})
```

- [ ] **Step 3: Run tests to confirm red**

```bash
npx vitest run src/__tests__/router.test.ts
```

Expected: FAIL — `vue-router` can be imported but let's confirm it runs. The test itself may pass since it creates its own router inline. If it passes, that's fine — the red phase is covered by the existing smoke test's `router-view` assertion from Task 3.

- [ ] **Step 4: Create route stub views**

Create `src/views/HomeView.vue`:
```vue
<template>
  <div class="p-4">
    <h2 class="text-xl font-semibold">Home</h2>
    <p class="text-gray-600">Your decks will appear here.</p>
  </div>
</template>
```

Create `src/views/DeckView.vue`:
```vue
<template>
  <div class="p-4">
    <h2 class="text-xl font-semibold">Deck</h2>
    <p class="text-gray-600">Cards in this deck will appear here.</p>
  </div>
</template>
```

Create `src/views/ReviewView.vue`:
```vue
<template>
  <div class="p-4">
    <h2 class="text-xl font-semibold">Review</h2>
    <p class="text-gray-600">Review session will appear here.</p>
  </div>
</template>
```

- [ ] **Step 5: Create `src/router/index.ts`**

```typescript
import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: HomeView },
    {
      path: '/deck/:id',
      name: 'deck',
      component: () => import('../views/DeckView.vue'),
    },
    {
      path: '/review/:deckId?',
      name: 'review',
      component: () => import('../views/ReviewView.vue'),
    },
  ],
})

export default router
```

- [ ] **Step 6: Wire router into `src/main.ts`**

```typescript
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './assets/main.css'

createApp(App).use(router).mount('#app')
```

- [ ] **Step 7: Update `src/App.vue` to use router-view**

```vue
<template>
  <div class="min-h-screen bg-gray-50">
    <header class="bg-white shadow-sm">
      <div class="max-w-3xl mx-auto px-4 py-3">
        <router-link to="/" class="text-2xl font-bold text-gray-900">Spacer</router-link>
      </div>
    </header>
    <main class="max-w-3xl mx-auto">
      <router-view />
    </main>
  </div>
</template>
```

- [ ] **Step 8: Update smoke test to provide router**

The smoke test needs the router plugin now. Update `src/__tests__/smoke.test.ts`:

```typescript
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import App from '../App.vue'

function createTestRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: { template: '<div>Home</div>' } },
      { path: '/deck/:id', component: { template: '<div>Deck</div>' } },
      { path: '/review/:deckId?', component: { template: '<div>Review</div>' } },
    ],
  })
}

describe('App', () => {
  it('renders the app title', async () => {
    const router = createTestRouter()
    router.push('/')
    await router.isReady()
    const wrapper = mount(App, { global: { plugins: [router] } })
    expect(wrapper.text()).toContain('Spacer')
  })

  it('has a router-view for page content', async () => {
    const router = createTestRouter()
    router.push('/')
    await router.isReady()
    const wrapper = mount(App, { global: { plugins: [router] } })
    expect(wrapper.html()).toContain('Home')
  })
})
```

- [ ] **Step 9: Run all tests to confirm green**

```bash
npx vitest run
```

Expected: All tests pass — smoke (2) + router (3) = 5 tests.

- [ ] **Step 10: Commit**

```bash
git add src/router/ src/views/ src/__tests__/ src/App.vue src/main.ts package.json package-lock.json
git commit -m "feat: add Vue Router with home, deck, and review route stubs

AI-Used: [claude]"
```

---

### Task 5: Install Pinia with store stubs

**Files:**
- Create: `src/stores/decks.ts`, `src/stores/reviews.ts`
- Create: `src/stores/__tests__/decks.test.ts`
- Modify: `src/main.ts`

- [ ] **Step 1: Install Pinia**

```bash
npm install pinia
```

- [ ] **Step 2: Write failing store test**

Create `src/stores/__tests__/decks.test.ts`:

```typescript
import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useDecksStore } from '../decks'

describe('decks store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('starts with empty decks', () => {
    const store = useDecksStore()
    expect(store.decks).toEqual([])
  })
})
```

- [ ] **Step 3: Run test to confirm red**

```bash
npx vitest run src/stores
```

Expected: FAIL — module `../decks` not found.

- [ ] **Step 4: Create store stubs**

Create `src/stores/decks.ts`:

```typescript
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Deck } from '../db/models'

export const useDecksStore = defineStore('decks', () => {
  const decks = ref<Deck[]>([])
  return { decks }
})
```

Create `src/stores/reviews.ts`:

```typescript
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Card } from '../db/models'

export const useReviewsStore = defineStore('reviews', () => {
  const currentCards = ref<Card[]>([])
  const currentIndex = ref(0)
  return { currentCards, currentIndex }
})
```

Note: These stores reference `../db/models` which is created in Task 6. If running tasks out of order, create `src/db/models.ts` first. When running in sequence, the test will import only `decks.ts` which imports `Deck` — the type-only import won't cause a runtime error in Vitest even if the file doesn't exist yet, but if it does fail, create a minimal `src/db/models.ts` with just the type exports first.

- [ ] **Step 5: Wire Pinia into `src/main.ts`**

```typescript
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './assets/main.css'

createApp(App).use(createPinia()).use(router).mount('#app')
```

- [ ] **Step 6: Update smoke test to include Pinia**

In `src/__tests__/smoke.test.ts`, update the mount calls to include Pinia:

```typescript
import { createPinia } from 'pinia'

// Update both test mounts to:
const wrapper = mount(App, { global: { plugins: [router, createPinia()] } })
```

- [ ] **Step 7: Run all tests to confirm green**

```bash
npx vitest run
```

Expected: All tests pass — smoke (2) + router (3) + decks store (1) = 6 tests.

- [ ] **Step 8: Commit**

```bash
git add src/stores/ src/main.ts src/__tests__/smoke.test.ts package.json package-lock.json
git commit -m "feat: install Pinia with deck and review store stubs

AI-Used: [claude]"
```

---

### Task 6: Set up Dexie.js with Card + Deck schema

**Files:**
- Create: `src/db/models.ts`, `src/db/index.ts`
- Create: `src/db/__tests__/db.test.ts`

- [ ] **Step 1: Install Dexie**

```bash
npm install dexie
```

- [ ] **Step 2: Write failing test for DB schema**

Create `src/db/__tests__/db.test.ts`:

```typescript
import { describe, it, expect, beforeEach } from 'vitest'
import { db } from '../index'

describe('Database', () => {
  beforeEach(async () => {
    await db.delete()
    await db.open()
  })

  it('can store and retrieve a deck', async () => {
    const id = crypto.randomUUID()
    await db.decks.add({
      id,
      name: 'Test Deck',
      createdAt: new Date(),
    })
    const deck = await db.decks.get(id)
    expect(deck?.name).toBe('Test Deck')
  })

  it('can store and retrieve a card', async () => {
    const deckId = crypto.randomUUID()
    const cardId = crypto.randomUUID()
    await db.cards.add({
      id: cardId,
      deckId,
      front: 'Question',
      back: 'Answer',
      easeFactor: 2.5,
      interval: 0,
      repetitions: 0,
      nextReview: new Date(),
      createdAt: new Date(),
    })
    const card = await db.cards.get(cardId)
    expect(card?.front).toBe('Question')
    expect(card?.deckId).toBe(deckId)
  })
})
```

- [ ] **Step 3: Run test to confirm it fails**

```bash
npx vitest run src/db
```

Expected: FAIL — module `../index` not found.

- [ ] **Step 4: Create `src/db/models.ts`**

If this file already exists as a stub from Task 5, verify it has the full content. Otherwise create:

```typescript
export interface Card {
  id: string
  deckId: string
  front: string
  back: string
  easeFactor: number
  interval: number
  repetitions: number
  nextReview: Date
  createdAt: Date
}

export interface Deck {
  id: string
  name: string
  description?: string
  createdAt: Date
}
```

- [ ] **Step 5: Create `src/db/index.ts`**

```typescript
import Dexie, { type EntityTable } from 'dexie'
import type { Card, Deck } from './models'

const db = new Dexie('spacer') as Dexie & {
  cards: EntityTable<Card, 'id'>
  decks: EntityTable<Deck, 'id'>
}

db.version(1).stores({
  decks: 'id, name',
  cards: 'id, deckId, nextReview',
})

export { db }
```

- [ ] **Step 6: Run tests to confirm green**

```bash
npx vitest run src/db
```

Expected: 2 tests pass.

- [ ] **Step 7: Commit**

```bash
git add src/db/ package.json package-lock.json
git commit -m "feat: set up Dexie.js with Card and Deck schema

AI-Used: [claude]"
```

---

### Task 7: Implement SM-2 algorithm

**Files:**
- Create: `src/lib/sm2.ts`, `src/lib/__tests__/sm2.test.ts`

- [ ] **Step 1: Write failing SM-2 tests**

Create `src/lib/__tests__/sm2.test.ts`:

```typescript
import { describe, it, expect } from 'vitest'
import { sm2, type SM2State } from '../sm2'

const freshCard: SM2State = {
  easeFactor: 2.5,
  interval: 0,
  repetitions: 0,
}

describe('sm2', () => {
  it('resets on quality < 3', () => {
    const result = sm2(2, { easeFactor: 2.5, interval: 10, repetitions: 5 })
    expect(result.repetitions).toBe(0)
    expect(result.interval).toBe(1)
  })

  it('adjusts ease factor even on failed reviews', () => {
    // Standard SM-2 adjusts EF on every review, not just successful ones
    const result = sm2(1, { easeFactor: 2.5, interval: 10, repetitions: 5 })
    expect(result.easeFactor).not.toBe(2.5)
  })

  it('sets interval to 1 on first successful review', () => {
    const result = sm2(4, freshCard)
    expect(result.repetitions).toBe(1)
    expect(result.interval).toBe(1)
  })

  it('sets interval to 6 on second successful review', () => {
    const result = sm2(4, { easeFactor: 2.5, interval: 1, repetitions: 1 })
    expect(result.repetitions).toBe(2)
    expect(result.interval).toBe(6)
  })

  it('multiplies interval by easeFactor on third+ review', () => {
    const result = sm2(4, { easeFactor: 2.5, interval: 6, repetitions: 2 })
    expect(result.repetitions).toBe(3)
    expect(result.interval).toBe(15) // Math.round(6 * 2.5)
  })

  it('adjusts ease factor based on quality', () => {
    const result = sm2(5, freshCard)
    expect(result.easeFactor).toBeCloseTo(2.6, 1)
  })

  it('never lets ease factor drop below 1.3', () => {
    const result = sm2(3, { easeFactor: 1.3, interval: 1, repetitions: 1 })
    expect(result.easeFactor).toBeGreaterThanOrEqual(1.3)
  })
})
```

- [ ] **Step 2: Run test to confirm red**

```bash
npx vitest run src/lib
```

Expected: FAIL — module not found.

- [ ] **Step 3: Implement SM-2**

Create `src/lib/sm2.ts`:

```typescript
export interface SM2State {
  easeFactor: number
  interval: number
  repetitions: number
}

export function sm2(quality: number, state: SM2State): SM2State {
  if (quality < 0 || quality > 5) {
    throw new RangeError('Quality must be between 0 and 5')
  }

  const newEF = Math.max(
    1.3,
    state.easeFactor + (0.1 - (5 - quality) * (0.08 + (5 - quality) * 0.02)),
  )

  if (quality < 3) {
    return { easeFactor: newEF, interval: 1, repetitions: 0 }
  }

  let interval: number
  if (state.repetitions === 0) {
    interval = 1
  } else if (state.repetitions === 1) {
    interval = 6
  } else {
    interval = Math.round(state.interval * state.easeFactor)
  }

  return {
    easeFactor: newEF,
    interval,
    repetitions: state.repetitions + 1,
  }
}
```

- [ ] **Step 4: Run tests to confirm green**

```bash
npx vitest run src/lib
```

Expected: 7 tests pass.

- [ ] **Step 5: Commit**

```bash
git add src/lib/
git commit -m "feat: implement SM-2 spaced repetition algorithm

AI-Used: [claude]"
```

---

### Task 8: Configure PWA manifest

**Files:**
- Modify: `vite.config.ts`
- Create: `public/favicon.svg`
- Create: `src/__tests__/pwa.test.ts`

- [ ] **Step 1: Install vite-plugin-pwa**

```bash
npm install -D vite-plugin-pwa
```

- [ ] **Step 2: Write failing PWA build test**

Create `src/__tests__/pwa.test.ts`:

```typescript
import { describe, it, expect } from 'vitest'
import { readFileSync, existsSync } from 'fs'
import { resolve } from 'path'
import { execSync } from 'child_process'

describe('PWA manifest', () => {
  it('build generates a web manifest with correct metadata', () => {
    execSync('npx vite build', { stdio: 'pipe' })
    const distDir = resolve(__dirname, '../../dist')

    // Find the manifest file
    const manifestPath = resolve(distDir, 'manifest.webmanifest')
    expect(existsSync(manifestPath)).toBe(true)

    const manifest = JSON.parse(readFileSync(manifestPath, 'utf-8'))
    expect(manifest.name).toBe('Spacer')
    expect(manifest.display).toBe('standalone')
    expect(manifest.theme_color).toBe('#6366f1')
  })
})
```

- [ ] **Step 3: Run test to confirm red**

```bash
npx vitest run src/__tests__/pwa.test.ts
```

Expected: FAIL — either build doesn't produce manifest, or file doesn't exist.

- [ ] **Step 4: Create a simple SVG favicon**

Create `public/favicon.svg`:

```svg
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100">
  <rect width="100" height="100" rx="20" fill="#6366f1"/>
  <text x="50" y="68" font-size="50" text-anchor="middle" fill="white" font-family="system-ui" font-weight="bold">S</text>
</svg>
```

- [ ] **Step 5: Update `vite.config.ts` with PWA config**

Note: `vite-plugin-pwa` with `registerType: 'autoUpdate'` auto-injects the service worker registration into the build output. No manual registration code in `main.ts` is needed.

```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { VitePWA } from 'vite-plugin-pwa'

export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
    VitePWA({
      registerType: 'autoUpdate',
      manifest: {
        name: 'Spacer',
        short_name: 'Spacer',
        description: 'Spaced repetition flashcards',
        theme_color: '#6366f1',
        background_color: '#f9fafb',
        display: 'standalone',
        icons: [
          {
            src: 'favicon.svg',
            sizes: 'any',
            type: 'image/svg+xml',
            purpose: 'any maskable',
          },
        ],
      },
    }),
  ],
})
```

- [ ] **Step 6: Run PWA test to confirm green**

```bash
npx vitest run src/__tests__/pwa.test.ts
```

Expected: PASS — manifest exists with correct metadata.

- [ ] **Step 7: Commit**

```bash
git add vite.config.ts public/favicon.svg src/__tests__/pwa.test.ts package.json package-lock.json
git commit -m "feat: configure PWA manifest with vite-plugin-pwa

AI-Used: [claude]"
```

---

### Task 9: Create component stubs + dev scripts + update .gitignore

**Files:**
- Create: `src/components/CardEditor.vue`, `src/components/ReviewCard.vue`, `src/components/DeckList.vue`
- Create: `dev/test`, `dev/dev`
- Modify: `.gitignore`

- [ ] **Step 1: Create component stubs**

Create `src/components/CardEditor.vue`:
```vue
<template>
  <div class="card-editor">
    <!-- Card create/edit form — implemented in a future increment -->
  </div>
</template>
```

Create `src/components/ReviewCard.vue`:
```vue
<template>
  <div class="review-card">
    <!-- Flip card during review — implemented in a future increment -->
  </div>
</template>
```

Create `src/components/DeckList.vue`:
```vue
<template>
  <div class="deck-list">
    <!-- Deck overview — implemented in a future increment -->
  </div>
</template>
```

- [ ] **Step 2: Create `dev/` directory and scripts**

Create `dev/test`:
```bash
#!/usr/bin/env bash
set -euo pipefail
npx vitest run "$@"
```

Create `dev/dev`:
```bash
#!/usr/bin/env bash
set -euo pipefail
npx vite "$@"
```

- [ ] **Step 3: Make scripts executable**

```bash
chmod +x dev/test dev/dev
```

- [ ] **Step 4: Update `.gitignore`**

Replace contents with:
```
node_modules/
dist/
*.local
.DS_Store
```

- [ ] **Step 5: Verify scripts work**

```bash
./dev/test
```

Expected: All tests pass (smoke + router + decks store + db + sm2 + pwa).

- [ ] **Step 6: Commit**

```bash
git add src/components/ dev/ .gitignore
git commit -m "feat: add component stubs, dev scripts, and update .gitignore

AI-Used: [claude]"
```

---

### Task 10: Set up docs for 5-minute increment workflow

**Supplementary:** These files are not in the spec's Bootstrap Deliverables but are required by the 5-minute increment skill to function. Creating them here so the workflow is ready immediately after bootstrap.

**Files:**
- Create: `docs/status.md`, `docs/retros.md`, `docs/issues/001-bootstrap.md`

- [ ] **Step 1: Create `docs/status.md`**

```markdown
# Spacer — Project Status

**Last updated:** 2026-03-14
**Current increment:** 1
**Streak:** 0

## Done

## In Progress
- #1 Bootstrap project with full stack — CLOSE

## Up Next

## Blocked
```

- [ ] **Step 2: Create `docs/retros.md`**

```markdown
# Spacer — Retrospectives
```

- [ ] **Step 3: Create `docs/issues/001-bootstrap.md`**

```markdown
# Bootstrap project with chosen stack

**Status:** in-progress
**Priority:** p0
**Labels:** infra
**Created:** 2026-03-14
**Closed:**

## Description
Set up the full project scaffold: Vite + Vue 3 + TypeScript + Tailwind CSS + Vitest + Dexie.js + Pinia + PWA manifest. The app runs, tests pass, and the PWA manifest serves.

## Acceptance Criteria
- [ ] Vite dev server starts and serves the app
- [ ] Tailwind CSS classes render correctly
- [ ] Vitest runs with at least one passing test
- [ ] Dexie.js can store and retrieve cards/decks
- [ ] SM-2 algorithm passes unit tests
- [ ] PWA manifest is generated on build
- [ ] Vue Router serves 3 route stubs

## Notes
Bootstrap increment — allowed to exceed 5 minutes per workflow rules.
```

- [ ] **Step 4: Commit**

```bash
git add docs/status.md docs/retros.md docs/issues/
git commit -m "docs: set up status, retros, and bootstrap issue for increment workflow

AI-Used: [claude]"
```

---

### Task 11: Final verification

- [ ] **Step 1: Run all tests**

```bash
./dev/test
```

Expected: All tests pass (smoke, router, decks store, db, sm2, pwa).

- [ ] **Step 2: Verify dev server starts**

```bash
npx vite --open=false &
sleep 3
curl -s http://localhost:5173 | grep -o 'Spacer' | head -1
kill %1
```

Expected: "Spacer" in output.

- [ ] **Step 3: Verify build + PWA manifest**

```bash
npx vite build && ls dist/
```

Expected: Build succeeds, contains `manifest.webmanifest`, JS/CSS assets.

- [ ] **Step 4: Update issue and status to done**

Update `docs/issues/001-bootstrap.md`: set Status to `done`, check all AC boxes, set Closed date to 2026-03-14.

Update `docs/status.md`: move #1 to Done, clear In Progress, set streak to 1.

- [ ] **Step 5: Commit**

```bash
git add docs/issues/001-bootstrap.md docs/status.md
git commit -m "docs: close bootstrap issue #1

AI-Used: [claude]"
```
