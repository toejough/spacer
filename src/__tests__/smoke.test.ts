import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createPinia } from 'pinia'
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
    const wrapper = mount(App, { global: { plugins: [router, createPinia()] } })
    expect(wrapper.text()).toContain('Spacer')
  })

  it('has a router-view for page content', async () => {
    const router = createTestRouter()
    router.push('/')
    await router.isReady()
    const wrapper = mount(App, { global: { plugins: [router, createPinia()] } })
    expect(wrapper.html()).toContain('Home')
  })
})
