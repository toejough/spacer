import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import MobileTodoList from 'src/components/MobileTodoList.vue'

describe('MobileTodoList', () => {
  it('initializes with demo todos and responds to actions', async () => {
    const wrapper = mount(MobileTodoList)
    expect(wrapper.text()).toContain('Write project spec')
    // trigger Done on first button
    const doneBtn = wrapper.find('button.btn-done')
    await doneBtn.trigger('click')
    // first todo should now be marked done
    expect(wrapper.html()).toContain('completed')
  })
})
