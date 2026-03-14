import 'fake-indexeddb/auto'
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
