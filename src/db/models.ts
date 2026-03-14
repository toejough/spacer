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
