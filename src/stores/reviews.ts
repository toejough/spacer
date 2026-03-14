import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Card } from '../db/models'

export const useReviewsStore = defineStore('reviews', () => {
  const currentCards = ref<Card[]>([])
  const currentIndex = ref(0)
  return { currentCards, currentIndex }
})
