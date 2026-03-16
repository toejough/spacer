import { describe, it, expect, beforeEach } from "vitest";
import { effectScope } from "vue";
import { SpacerDB, createDeck, createCard, getAllDecks, getDeck, getDeckCards } from "../../src/db";
import { useLiveViewState } from "../../src/view-state";

describe("useLiveViewState with real Dexie", () => {
  let db: SpacerDB;

  beforeEach(() => {
    db = new SpacerDB("test-view-states-" + Date.now());
  });

  // Given no decks exist
  // When useLiveViewState subscribes to getAllDecks
  // Then it transitions to empty
  it("empty database produces empty view state", async () => {
    const scope = effectScope();
    let state!: ReturnType<typeof useLiveViewState>;

    scope.run(() => {
      state = useLiveViewState(
        () => getAllDecks(db),
        (d) => d.length === 0
      );
    });

    await new Promise((r) => setTimeout(r, 50));
    expect(state.value.status).toBe("empty");
    scope.stop();
  });

  // Given a deck exists with cards
  // When useLiveViewState subscribes to getDeckCards
  // Then it transitions to loaded with the card data
  it("deck with cards produces loaded view state", async () => {
    const deckId = await createDeck(db, "Spanish");
    await createCard(db, deckId, "hola", "hello");

    const scope = effectScope();
    let state!: ReturnType<typeof useLiveViewState>;

    scope.run(() => {
      state = useLiveViewState(
        () => getDeckCards(db, deckId),
        (c) => c.length === 0
      );
    });

    await new Promise((r) => setTimeout(r, 50));
    expect(state.value.status).toBe("loaded");
    scope.stop();
  });

  // Given a nonexistent deck ID
  // When useLiveViewState subscribes to getDeck with notFoundWhen
  // Then it transitions to not-found
  it("nonexistent deck produces not-found view state", async () => {
    const scope = effectScope();
    let state!: ReturnType<typeof useLiveViewState>;

    scope.run(() => {
      state = useLiveViewState(
        () => getDeck(db, 99999),
        () => false,
        { notFoundWhen: (d) => d === undefined }
      );
    });

    await new Promise((r) => setTimeout(r, 50));
    expect(state.value.status).toBe("not-found");
    scope.stop();
  });
});
