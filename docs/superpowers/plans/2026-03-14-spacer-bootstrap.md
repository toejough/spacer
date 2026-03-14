# Spacer PWA Bootstrap Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Bootstrap the Spacer PWA with Vue 3 + Vite + Tailwind + Vitest + Dexie + PWA manifest so every subsequent 5-minute increment can deliver a real feature.

**Architecture:** Single-page app with Vue Router (3 route stubs), Pinia stores, and Dexie.js for IndexedDB persistence. SM-2 algorithm lives as a pure function in `src/lib/`. PWA via vite-plugin-pwa.

**Tech Stack:** Vue 3, TypeScript, Vite, Tailwind CSS, Pinia, Dexie.js, Vitest, @vue/test-utils, vite-plugin-pwa

**Spec:** `docs/superpowers/specs/2026-03-14-spacer-bootstrap-design.md`

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
git add -A
git commit -m "feat: scaffold Vite + Vue 3 + TypeScript project"
```

---

### Task 2: Configure Tailwind CSS

**Files:**
- Create: `tailwind.config.js`, `postcss.config.js`
- Modify: `src/assets/main.css` (replace contents)

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
git add -A
git commit -m "feat: configure Tailwind CSS"
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
})
```

- [ ] **Step 4: Run test to verify it passes**

```bash
npx vitest run
```

Expected: 1 test passes (it should pass because App.vue already contains "Spacer" from Task 2).

- [ ] **Step 5: Add test script to package.json**

Add to `scripts`: `"test": "vitest run"`, `"test:watch": "vitest"`

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "feat: configure Vitest with smoke test"
```

---

### Task 4: Set up Vue Router with route stubs

**Files:**
- Create: `src/router/index.ts`, `src/views/HomeView.vue`, `src/views/DeckView.vue`, `src/views/ReviewView.vue`
- Modify: `src/main.ts`, `src/App.vue`

- [ ] **Step 1: Install Vue Router**

```bash
npm install vue-router
```

- [ ] **Step 2: Create route stub views**

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

- [ ] **Step 3: Create `src/router/index.ts`**

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

- [ ] **Step 4: Wire router into `src/main.ts`**

```typescript
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './assets/main.css'

createApp(App).use(router).mount('#app')
```

- [ ] **Step 5: Update `src/App.vue` to use router-view**

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

- [ ] **Step 6: Run tests to verify nothing broke**

```bash
npx vitest run
```

Expected: Smoke test may need updating since App.vue now uses router-link. Update the smoke test to provide the router:

```typescript
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import App from '../App.vue'

describe('App', () => {
  it('renders the app title', async () => {
    const router = createRouter({
      history: createWebHistory(),
      routes: [{ path: '/', component: { template: '<div>Home</div>' } }],
    })
    router.push('/')
    await router.isReady()

    const wrapper = mount(App, { global: { plugins: [router] } })
    expect(wrapper.text()).toContain('Spacer')
  })
})
```

- [ ] **Step 7: Commit**

```bash
git add -A
git commit -m "feat: add Vue Router with home, deck, and review route stubs"
```

---

### Task 5: Install Pinia

**Files:**
- Modify: `src/main.ts`

- [ ] **Step 1: Install Pinia**

```bash
npm install pinia
```

- [ ] **Step 2: Wire Pinia into `src/main.ts`**

```typescript
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './assets/main.css'

createApp(App).use(createPinia()).use(router).mount('#app')
```

- [ ] **Step 3: Run tests**

```bash
npx vitest run
```

Expected: All tests pass. Update smoke test to include Pinia if needed:
```typescript
import { createPinia } from 'pinia'
// in mount: global: { plugins: [router, createPinia()] }
```

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "feat: install and wire Pinia"
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
git add -A
git commit -m "feat: set up Dexie.js with Card and Deck schema"
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

Expected: 6 tests pass.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat: implement SM-2 spaced repetition algorithm"
```

---

### Task 8: Configure PWA manifest

**Files:**
- Modify: `vite.config.ts`
- Create: `public/favicon.svg`

- [ ] **Step 1: Install vite-plugin-pwa**

```bash
npm install -D vite-plugin-pwa
```

- [ ] **Step 2: Create a simple SVG favicon**

Create `public/favicon.svg`:

```svg
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100">
  <rect width="100" height="100" rx="20" fill="#6366f1"/>
  <text x="50" y="68" font-size="50" text-anchor="middle" fill="white" font-family="system-ui" font-weight="bold">S</text>
</svg>
```

- [ ] **Step 3: Update `vite.config.ts` with PWA config**

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

- [ ] **Step 4: Build and verify PWA manifest is generated**

```bash
npx vite build 2>&1 | tail -10
ls dist/*.webmanifest 2>/dev/null || ls dist/manifest.webmanifest 2>/dev/null || echo "Check dist/ for manifest"
```

Expected: Build succeeds, manifest file exists in `dist/`.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat: configure PWA manifest with vite-plugin-pwa"
```

---

### Task 9: Create dev scripts + update .gitignore

**Files:**
- Create: `dev/test`, `dev/dev`
- Modify: `.gitignore`

- [ ] **Step 1: Create `dev/` directory and scripts**

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

- [ ] **Step 2: Make scripts executable**

```bash
chmod +x dev/test dev/dev
```

- [ ] **Step 3: Update `.gitignore`**

Replace contents with:
```
node_modules/
dist/
*.local
.DS_Store
```

- [ ] **Step 4: Verify scripts work**

```bash
./dev/test
```

Expected: All tests pass (smoke + db + sm2).

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat: add dev scripts and update .gitignore"
```

---

### Task 10: Set up docs for 5-minute increment workflow

**Files:**
- Create: `docs/status.md`, `docs/retros.md`, `docs/issues/`

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

- [ ] **Step 3: Create `docs/issues/` directory with bootstrap issue**

Create `docs/issues/001-bootstrap.md`:

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
git add -A
git commit -m "docs: set up status, retros, and bootstrap issue for increment workflow"
```

---

### Task 11: Final verification

- [ ] **Step 1: Run all tests**

```bash
./dev/test
```

Expected: All tests pass (smoke, db, sm2).

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

Expected: Build succeeds, contains `manifest.webmanifest` (or similar), JS/CSS assets.

- [ ] **Step 4: Update issue and status to done**

Update `docs/issues/001-bootstrap.md`: set Status to `done`, check all AC boxes, set Closed date.

Update `docs/status.md`: move #1 to Done, clear In Progress, set streak to 1.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "docs: close bootstrap issue #1"
```
