import { describe, it, expect } from "vitest";
import { effectScope } from "vue";
import { type ViewState, useLiveViewState } from "../../src/view-state";

describe("ViewState type", () => {
  // Given a ViewState in 'loaded' status
  // When narrowed via status check
  // Then data property is accessible
  it("narrows to loaded with data access", () => {
    const state: ViewState<string[]> = { status: "loaded", data: ["a"] };
    if (state.status === "loaded") {
      expect(state.data).toEqual(["a"]);
    }
  });

  // Given a ViewState in 'error' status
  // When narrowed via status check
  // Then message property is accessible
  it("narrows to error with message access", () => {
    const state: ViewState<string[]> = { status: "error", message: "fail" };
    if (state.status === "error") {
      expect(state.message).toBe("fail");
    }
  });
});

describe("useLiveViewState", () => {
  // Given a querier that resolves to undefined with notFoundWhen set
  // When useLiveViewState is called
  // Then it transitions to not-found
  it("transitions from loading to not-found when notFoundWhen matches", async () => {
    const scope = effectScope();
    let state!: ReturnType<typeof useLiveViewState<string | undefined>>;

    scope.run(() => {
      state = useLiveViewState(
        () => Promise.resolve(undefined),
        (d) => false,
        { notFoundWhen: (d) => d === undefined }
      );
    });

    await new Promise((r) => setTimeout(r, 50));

    expect(state.value).toEqual({ status: "not-found" });
    scope.stop();
  });

  // Given a querier that resolves to a non-empty array
  // When useLiveViewState is called with an isEmpty predicate
  // Then it starts as loading, then transitions to loaded
  it("transitions from loading to loaded", async () => {
    const scope = effectScope();
    let state!: ReturnType<typeof useLiveViewState<string[]>>;

    scope.run(() => {
      state = useLiveViewState(
        () => Promise.resolve(["a", "b"]),
        (d) => d.length === 0
      );
    });

    expect(state.value.status).toBe("loading");

    await new Promise((r) => setTimeout(r, 50));

    expect(state.value).toEqual({ status: "loaded", data: ["a", "b"] });
    scope.stop();
  });

  // Given a querier that resolves to an empty array
  // When useLiveViewState is called
  // Then it transitions to empty
  it("transitions from loading to empty", async () => {
    const scope = effectScope();
    let state!: ReturnType<typeof useLiveViewState<string[]>>;

    scope.run(() => {
      state = useLiveViewState(
        () => Promise.resolve([]),
        (d) => d.length === 0
      );
    });

    await new Promise((r) => setTimeout(r, 50));

    expect(state.value).toEqual({ status: "empty" });
    scope.stop();
  });

  // Given a querier that rejects
  // When useLiveViewState is called
  // Then it transitions to error with the message
  it("transitions from loading to error on rejection", async () => {
    const scope = effectScope();
    let state!: ReturnType<typeof useLiveViewState<string[]>>;

    scope.run(() => {
      state = useLiveViewState(
        () => Promise.reject(new Error("db failure")),
        (d) => d.length === 0
      );
    });

    await new Promise((r) => setTimeout(r, 50));

    expect(state.value).toEqual({ status: "error", message: "db failure" });
    scope.stop();
  });

  // Given useLiveViewState is running
  // When the scope is disposed
  // Then the subscription is cleaned up (no error thrown)
  it("cleans up subscription on scope dispose", async () => {
    const scope = effectScope();

    scope.run(() => {
      useLiveViewState(
        () => Promise.resolve(["a"]),
        (d) => d.length === 0
      );
    });

    // Should not throw
    scope.stop();
  });
});
