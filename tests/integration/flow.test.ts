import { describe, it, expect, beforeEach } from "vitest";
import { SpacerDB } from "../../src/db";
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
    // Create a deck
    const deckId = await db.decks.add({
      name: "Test Deck",
      createdAt: new Date(),
    } as any);

    // Add a card with default SM-2 state
    const now = new Date();
    const cardId = await db.cards.add({
      deckId,
      front: "What is 2+2?",
      back: "4",
      easeFactor: 2.5,
      interval: 0,
      repetitions: 0,
      nextReview: now,
    } as any);

    // Verify card is queryable by deckId
    const cards = await db.cards.where("deckId").equals(deckId).toArray();
    expect(cards).toHaveLength(1);
    expect(cards[0].front).toBe("What is 2+2?");

    // Filter due cards (nextReview <= now)
    const dueCards = cards.filter((c) => c.nextReview <= now);
    expect(dueCards).toHaveLength(1);

    // Review: apply SM-2 with quality=4 (good)
    const card = dueCards[0];
    const result = sm2(
      {
        easeFactor: card.easeFactor,
        interval: card.interval,
        repetitions: card.repetitions,
      },
      4
    );

    // Persist SM-2 result
    await db.cards.update(card.id, {
      easeFactor: result.easeFactor,
      interval: result.interval,
      repetitions: result.repetitions,
      nextReview: result.nextReview,
    });

    // Verify persisted state
    const updated = await db.cards.get(card.id);
    expect(updated!.repetitions).toBe(1);
    expect(updated!.interval).toBe(1);
    expect(updated!.easeFactor).toBeGreaterThanOrEqual(2.5);
    expect(updated!.nextReview.getTime()).toBeGreaterThan(now.getTime());

    // Card should no longer be due (nextReview is tomorrow)
    const stillDue = (await db.cards.where("deckId").equals(deckId).toArray())
      .filter((c) => c.nextReview <= now);
    expect(stillDue).toHaveLength(0);
  });

  it("handles failed review (quality < 3) by resetting repetitions", async () => {
    const deckId = await db.decks.add({
      name: "Fail Deck",
      createdAt: new Date(),
    } as any);

    await db.cards.add({
      deckId,
      front: "Hard question",
      back: "Hard answer",
      easeFactor: 2.5,
      interval: 6,
      repetitions: 2,
      nextReview: new Date(),
    } as any);

    const card = (await db.cards.where("deckId").equals(deckId).toArray())[0];
    const result = sm2(
      {
        easeFactor: card.easeFactor,
        interval: card.interval,
        repetitions: card.repetitions,
      },
      1 // fail
    );

    await db.cards.update(card.id, {
      easeFactor: result.easeFactor,
      interval: result.interval,
      repetitions: result.repetitions,
      nextReview: result.nextReview,
    });

    const updated = await db.cards.get(card.id);
    expect(updated!.repetitions).toBe(0);
    expect(updated!.interval).toBe(0);
  });
});
