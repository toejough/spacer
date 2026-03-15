import { describe, it, expect } from "vitest";
import { newSM2State } from "../../src/sm2";

describe("newSM2State", () => {
  it("returns canonical SM-2 initial values", () => {
    const now = new Date("2025-01-15T10:00:00Z");
    const state = newSM2State(now);

    expect(state).toEqual({
      easeFactor: 2.5,
      interval: 0,
      repetitions: 0,
      nextReview: now,
    });
  });

  it("defaults nextReview to current time", () => {
    const before = new Date();
    const state = newSM2State();
    const after = new Date();

    expect(state.nextReview.getTime()).toBeGreaterThanOrEqual(before.getTime());
    expect(state.nextReview.getTime()).toBeLessThanOrEqual(after.getTime());
  });
});
