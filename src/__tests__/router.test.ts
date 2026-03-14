import { describe, it, expect } from 'vitest'
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
