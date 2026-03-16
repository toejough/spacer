import Dexie, { type EntityTable } from "dexie";
import { newSM2State, type SM2Result } from "./sm2";

export interface Deck {
  id: number;
  name: string;
  createdAt: Date;
}

export interface Card {
  id: number;
  deckId: number;
  front: string;
  back: string;
  easeFactor: number;
  interval: number;
  repetitions: number;
  nextReview: Date;
}

export class SpacerDB extends Dexie {
  decks!: EntityTable<Deck, "id">;
  cards!: EntityTable<Card, "id">;

  constructor(name = "spacer") {
    super(name);
    this.version(1).stores({
      decks: "++id, name",
      cards: "++id, deckId, nextReview",
    });
  }
}

export const db = new SpacerDB();

// Queries

export function getAllDecks(db: SpacerDB): Promise<Deck[]> {
  return db.decks.toArray();
}

export function getDeck(db: SpacerDB, id: number): Promise<Deck | undefined> {
  return db.decks.get(id);
}

export function getDeckCards(db: SpacerDB, deckId: number): Promise<Card[]> {
  return db.cards.where("deckId").equals(deckId).toArray();
}

export async function getDueCards(db: SpacerDB, deckId: number, now = new Date()): Promise<Card[]> {
  const cards = await db.cards.where("deckId").equals(deckId).toArray();
  return cards.filter((c) => c.nextReview <= now);
}

// Mutations

export async function createDeck(db: SpacerDB, name: string): Promise<number> {
  return db.decks.add({ name, createdAt: new Date() });
}

export async function createCard(db: SpacerDB, deckId: number, front: string, back: string): Promise<number> {
  const sm2 = newSM2State();
  return db.cards.add({
    deckId,
    front,
    back,
    ...sm2,
  });
}

export async function updateCardReview(db: SpacerDB, cardId: number, result: SM2Result): Promise<void> {
  await db.cards.update(cardId, {
    easeFactor: result.easeFactor,
    interval: result.interval,
    repetitions: result.repetitions,
    nextReview: result.nextReview,
  });
}
