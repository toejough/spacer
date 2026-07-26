import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import MobileTodoCard from 'src/components/MobileTodoCard.vue'

describe('MobileTodoCard', () => {
  it('renders open state with Done and Abandon buttons', () => {
    const todo = { id: '1', title: 'Task', status: 'open' }
    const wrapper = mount(MobileTodoCard, { props: { todo } })
    expect(wrapper.text()).toContain('Task')
    expect(wrapper.find('button.btn-done').exists()).toBe(true)
    expect(wrapper.find('button.btn-abandon').exists()).toBe(true)
  })

  it('renders Reopen on closed items', () => {
    const todo = { id: '2', title: 'Done', status: 'done' }
    const wrapper = mount(MobileTodoCard, { props: { todo } })
    expect(wrapper.find('button.btn-reopen').exists()).toBe(true)
  })
})
