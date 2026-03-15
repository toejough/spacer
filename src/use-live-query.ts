import { ref, onScopeDispose, type Ref } from "vue";
import { liveQuery } from "dexie";

export function useLiveQuery<T>(
  querier: () => Promise<T>,
  initialValue: T
): Ref<T> {
  const data = ref(initialValue) as Ref<T>;
  const subscription = liveQuery(querier).subscribe({
    next: (value) => {
      data.value = value;
    },
  });
  onScopeDispose(() => subscription.unsubscribe());
  return data;
}
