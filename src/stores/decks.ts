import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Deck } from '../db/models'

export const useDecksStore = defineStore('decks', () => {
  const decks = ref<Deck[]>([])
  return { decks }
})
