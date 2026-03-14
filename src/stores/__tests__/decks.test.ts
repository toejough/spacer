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
