import { describe, it, expect, beforeEach } from "vitest";
import { SpacerDB } from "../../src/db";
import {
  getAllDecks,
  getDeck,
  getDeckCards,
  getDueCards,
  createDeck,
  createCard,
  updateCardReview,
} from "../../src/db";
import { sm2, newSM2State } from "../../src/sm2";

let testDbCounter = 0;

function createTestDB() {
  return new SpacerDB(`test-queries-${++testDbCounter}-${Date.now()}`);
}

describe("query functions", () => {
  let db: SpacerDB;

  beforeEach(() => {
    db = createTestDB();
  });

  describe("getAllDecks", () => {
    it("returns empty array when no decks exist", async () => {
      expect(await getAllDecks(db)).toEqual([]);
    });

    it("returns all decks", async () => {
      await createDeck(db, "Deck A");
      await createDeck(db, "Deck B");
      const decks = await getAllDecks(db);
      expect(decks).toHaveLength(2);
      expect(decks.map((d) => d.name)).toEqual(["Deck A", "Deck B"]);
    });
  });

  describe("getDeck", () => {
    it("returns undefined for nonexistent deck", async () => {
      expect(await getDeck(db, 999)).toBeUndefined();
    });

    it("returns the deck by id", async () => {
      const id = await createDeck(db, "My Deck");
      const deck = await getDeck(db, id);
      expect(deck?.name).toBe("My Deck");
    });
  });

  describe("getDeckCards", () => {
    it("returns only cards for the given deck", async () => {
      const d1 = await createDeck(db, "Deck 1");
      const d2 = await createDeck(db, "Deck 2");
      await createCard(db, d1, "Q1", "A1");
      await createCard(db, d2, "Q2", "A2");
      await createCard(db, d1, "Q3", "A3");

      const cards = await getDeckCards(db, d1);
      expect(cards).toHaveLength(2);
      expect(cards.map((c) => c.front)).toEqual(["Q1", "Q3"]);
    });
  });

  describe("getDueCards", () => {
    it("returns cards with nextReview <= now", async () => {
      const deckId = await createDeck(db, "Due Test");
      await createCard(db, deckId, "Due", "Yes");

      const now = new Date();
      const due = await getDueCards(db, deckId, now);
      expect(due).toHaveLength(1);
      expect(due[0].front).toBe("Due");
    });

    it("excludes cards with nextReview in the future", async () => {
      const deckId = await createDeck(db, "Future Test");
      const cardId = await createCard(db, deckId, "Future", "No");

      // Push nextReview into the future
      const tomorrow = new Date();
      tomorrow.setDate(tomorrow.getDate() + 1);
      await db.cards.update(cardId, { nextReview: tomorrow });

      const due = await getDueCards(db, deckId, new Date());
      expect(due).toHaveLength(0);
    });
  });

  describe("createDeck", () => {
    it("creates a deck with name and createdAt", async () => {
      const id = await createDeck(db, "New Deck");
      const deck = await db.decks.get(id);
      expect(deck?.name).toBe("New Deck");
      expect(deck?.createdAt).toBeInstanceOf(Date);
    });
  });

  describe("createCard", () => {
    it("creates a card with canonical SM-2 defaults", async () => {
      const deckId = await createDeck(db, "Card Test");
      const cardId = await createCard(db, deckId, "Front", "Back");
      const card = await db.cards.get(cardId);

      expect(card?.deckId).toBe(deckId);
      expect(card?.front).toBe("Front");
      expect(card?.back).toBe("Back");
      expect(card?.easeFactor).toBe(2.5);
      expect(card?.interval).toBe(0);
      expect(card?.repetitions).toBe(0);
      expect(card?.nextReview).toBeInstanceOf(Date);
    });
  });

  describe("updateCardReview", () => {
    it("persists SM-2 result to the card", async () => {
      const deckId = await createDeck(db, "Review Test");
      const cardId = await createCard(db, deckId, "Q", "A");

      const card = (await db.cards.get(cardId))!;
      const result = sm2(
        { easeFactor: card.easeFactor, interval: card.interval, repetitions: card.repetitions },
        4
      );
      await updateCardReview(db, cardId, result);

      const updated = await db.cards.get(cardId);
      expect(updated?.repetitions).toBe(1);
      expect(updated?.interval).toBe(1);
      expect(updated?.easeFactor).toBeGreaterThanOrEqual(2.5);
      expect(updated!.nextReview.getTime()).toBeGreaterThan(card.nextReview.getTime());
    });
  });
});
