import { ref, onScopeDispose, type Ref } from "vue";
import { liveQuery } from "dexie";

export type ViewState<T> =
  | { status: "loading" }
  | { status: "loaded"; data: T }
  | { status: "empty" }
  | { status: "error"; message: string }
  | { status: "not-found" };

export interface LiveViewStateOptions<T> {
  notFoundWhen?: (data: T) => boolean;
}

export function useLiveViewState<T>(
  querier: () => Promise<T>,
  isEmpty: (data: T) => boolean,
  options?: LiveViewStateOptions<T>
): Ref<ViewState<T>> {
  const state = ref<ViewState<T>>({ status: "loading" }) as Ref<ViewState<T>>;

  const subscription = liveQuery(querier).subscribe({
    next: (value) => {
      if (options?.notFoundWhen?.(value)) {
        state.value = { status: "not-found" };
      } else {
        state.value = isEmpty(value)
          ? { status: "empty" }
          : { status: "loaded", data: value };
      }
    },
    error: (err) => {
      state.value = {
        status: "error",
        message: err instanceof Error ? err.message : String(err),
      };
    },
  });

  onScopeDispose(() => subscription.unsubscribe());

  return state;
}
