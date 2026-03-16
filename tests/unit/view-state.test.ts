import { describe, it, expect } from "vitest";
import { effectScope, nextTick } from "vue";
import { type ViewState, useLiveViewState } from "../../src/view-state";

describe("ViewState type", () => {
  it("narrows to loaded with data access", () => {
    const state: ViewState<string[]> = { status: "loaded", data: ["a"] };
    if (state.status === "loaded") {
      expect(state.data).toEqual(["a"]);
    }
  });

  it("narrows to error with message access", () => {
    const state: ViewState<string[]> = { status: "error", message: "fail" };
    if (state.status === "error") {
      expect(state.message).toBe("fail");
    }
  });
});

describe("useLiveViewState", () => {
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
