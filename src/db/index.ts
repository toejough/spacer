import Dexie, { type EntityTable } from 'dexie'
import type { Card, Deck } from './models'

const db = new Dexie('spacer') as Dexie & {
  cards: EntityTable<Card, 'id'>
  decks: EntityTable<Deck, 'id'>
}

db.version(1).stores({
  decks: 'id, name',
  cards: 'id, deckId, nextReview',
})

export { db }
