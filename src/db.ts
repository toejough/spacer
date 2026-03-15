import Dexie, { type EntityTable } from "dexie";

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
