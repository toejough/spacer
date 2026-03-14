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
