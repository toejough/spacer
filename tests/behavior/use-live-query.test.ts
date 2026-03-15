import { describe, it, expect, vi } from "vitest";
import { nextTick } from "vue";
import { useLiveQuery } from "../../src/use-live-query";
import { effectScope } from "vue";

describe("useLiveQuery", () => {
  it("starts with the initial value", () => {
    const scope = effectScope();
    scope.run(() => {
      const data = useLiveQuery(() => Promise.resolve(["a", "b"]), [] as string[]);
      expect(data.value).toEqual([]);
    });
    scope.stop();
  });

  it("updates when the querier resolves", async () => {
    const scope = effectScope();
    let result!: ReturnType<typeof useLiveQuery<string[]>>;

    scope.run(() => {
      result = useLiveQuery(() => Promise.resolve(["a", "b"]), [] as string[]);
    });

    // liveQuery is async — give it time to resolve
    await new Promise((r) => setTimeout(r, 50));

    expect(result.value).toEqual(["a", "b"]);
    scope.stop();
  });
});
