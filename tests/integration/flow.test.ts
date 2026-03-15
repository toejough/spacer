import { describe, it, expect, beforeEach } from "vitest";
import {
  SpacerDB,
  createDeck,
  createCard,
  getDeckCards,
  getDueCards,
  updateCardReview,
} from "../../src/db";
import { sm2 } from "../../src/sm2";

let testDbCounter = 0;

function createTestDB() {
  return new SpacerDB(`test-spacer-${++testDbCounter}-${Date.now()}`);
}

describe("Full flow: create deck → add card → review → SM-2 update", () => {
  let db: SpacerDB;

  beforeEach(() => {
    db = createTestDB();
  });

  it("completes the entire user journey", async () => {
    const deckId = await createDeck(db, "Test Deck");
    const now = new Date();

    await createCard(db, deckId, "What is 2+2?", "4");

    const cards = await getDeckCards(db, deckId);
    expect(cards).toHaveLength(1);
    expect(cards[0].front).toBe("What is 2+2?");

    const dueCards = await getDueCards(db, deckId, now);
    expect(dueCards).toHaveLength(1);

    const card = dueCards[0];
    const result = sm2(
      {
        easeFactor: card.easeFactor,
        interval: card.interval,
        repetitions: card.repetitions,
      },
      4
    );

    await updateCardReview(db, card.id, result);

    const updated = await db.cards.get(card.id);
    expect(updated!.repetitions).toBe(1);
    expect(updated!.interval).toBe(1);
    expect(updated!.easeFactor).toBeGreaterThanOrEqual(2.5);
    expect(updated!.nextReview.getTime()).toBeGreaterThan(now.getTime());

    const stillDue = await getDueCards(db, deckId, now);
    expect(stillDue).toHaveLength(0);
  });

  it("handles failed review (quality < 3) by resetting repetitions", async () => {
    const deckId = await createDeck(db, "Fail Deck");
    const cardId = await createCard(db, deckId, "Hard question", "Hard answer");

    // Push card to look like it has prior review history
    await db.cards.update(cardId, { interval: 6, repetitions: 2 });

    const card = (await db.cards.get(cardId))!;
    const result = sm2(
      {
        easeFactor: card.easeFactor,
        interval: card.interval,
        repetitions: card.repetitions,
      },
      1
    );

    await updateCardReview(db, cardId, result);

    const updated = await db.cards.get(cardId);
    expect(updated!.repetitions).toBe(0);
    expect(updated!.interval).toBe(0);
  });
});
