## Context

See proposal.md - Why. Relevant current state (from codebase investigation):

- Items (todos/notes) are plain JS objects in a flat array, persisted via `loadItems()`/`saveItems()` to `localStorage` under `remember_everything_items` (`srv/static/script.js`).
- Each tab (Todos, Notes, Search, Review) renders its own filtered/sorted view via `renderTodoCard()`/`renderNoteCard()` and a shared relevance-sort (`getReviewUrgencyRank`, state > review urgency > recency).
- The only existing gesture is swipe-to-reveal-delete (`wrapWithSwipeReveal()`), a horizontal pointer-drag limited to Done/Abandoned cards, ending in delete or snap-back. No general drag-and-drop, no drop-target detection, exists today.
- There is no framework and no build step: everything is hand-written vanilla JS/CSS/HTML served by the Go backend, which does not touch item data at all.
- A prior "column-board-with-swimlanes" proposal (automatic grouping by status+time, with drag-to-change-status) was drafted and deleted before implementation. It is not being resurrected; stacks are user-defined manual grouping, unrelated to status columns.

## Goals / Non-Goals

**Goals:**
- Add stacks as a client-side-only concept, consistent with the existing localStorage-only persistence model (no backend/API changes).
- Introduce drag-and-drop (handle-based) without breaking or conflicting with the existing swipe-to-delete gesture.
- Keep stack rendering/sorting integrated into the existing per-tab relevance-sort rather than as a separate mechanism.

**Non-Goals:**
- Cross-device sync of stacks (items aren't synced today either).
- Stack-onto-stack merging (explicitly out of scope per proposal).
- Reordering members within an expanded stack (not requested; membership and stack position are in scope, intra-stack member order is not specified beyond insertion order).
- Touch-and-hold / long-press drag initiation — the drag handle is a dedicated, always-visible affordance, so gesture disambiguation with swipe is by hit-target (handle vs. card body), not by gesture heuristics.

## Decisions

### Data model: `stack_id` on items + a separate `stacks` list
Add `stack_id: string | null` to the item object shape, and a new localStorage key (e.g. `remember_everything_stacks`) holding `{ id, name, created_at, updated_at }[]`. A stack's membership is derived by filtering items on `stack_id`, not stored as an ordered array on the stack — this avoids two sources of truth for membership (item.stack_id vs. a list-of-ids) that could drift. Insertion order for display purposes falls out of each member's own `updated_at`/creation order, consistent with how the rest of the app already treats ordering.

*Alternative considered:* stack owns an ordered array of item ids (Option B in the explore conversation). Rejected for v1 because it requires keeping two records in sync on every add/remove and this app has no transaction/atomicity guarantees beyond a single `saveItems()` call — a dangling id or an orphaned membership becomes possible on partial failures (e.g. export/import, manual localStorage edits). Item-owns-stack_id is self-healing: deleting an item automatically drops it from the stack with no cleanup step.

### Sort integration: compute stack position, don't create a parallel list
Reuse the existing relevance-sort function unmodified for ranking individual items. Build the rendered list per-tab in two steps: (1) rank all *visible-type* items (stacked or not) using the existing comparator, (2) walk that ranked list top to bottom, replacing each item with its stack's collapsed tile the first time any of that stack's members is encountered, and skipping that stack's remaining members thereafter (they are not yet-unranked; they simply already have a tile position). This makes "stack sorts by its most urgent visible member" fall out mechanically from the existing per-item ranking, rather than requiring a second bespoke stack-ranking function to be kept consistent with the item ranking.

### Drag-and-drop: native HTML5 DnD, not a library
Given the app has zero build step and no dependencies today, adding a DnD library (e.g. SortableJS) means introducing a new build/dependency-management story (the codebase's git history shows a `vuedraggable` era that was fully removed). Native HTML5 drag-and-drop (`draggable`, `dragstart`/`dragover`/`drop`) works without a build step and without new dependencies, at the cost of needing custom touch-device polyfilling (HTML5 DnD has no native touch support). Given this app is mobile-first (per the swipe gesture), touch dragging will be implemented as pointer-event-based drag (pointerdown on the handle → pointermove tracks a ghost element → pointerup resolves against the element under the pointer), mirroring the existing swipe-gesture's pointer-event approach for consistency, rather than mixing HTML5 DnD (desktop) with a separate touch implementation.

### Gesture disambiguation: dedicated drag handle
The drag handle is a separate, small hit-target rendered on each card (distinct from the swipe-to-delete zone, which is the card body). Pointer-down on the handle starts a drag; pointer-down elsewhere on the card behaves as today (swipe zone for Done/Abandoned cards, tap-to-edit otherwise). This avoids heuristic gesture conflict resolution (e.g. distinguishing "swipe" from "drag" by direction/velocity) entirely.

### Stack-name prompt UI
Reuse the existing modal-overlay pattern already used for the edit modal, rather than a native `prompt()`, for visual/accessibility consistency and to allow validation (non-empty name) with inline feedback instead of the OS-native dialog. Canceling this modal aborts the merge with no side effects (no stack created, no item's `stack_id` changed).

## Risks / Trade-offs

- [Native pointer-based DnD is more code to write and test than a library] → Mitigated by the app's existing precedent (swipe-to-delete) already being a hand-rolled pointer-event gesture; the same patterns (deadzone, direction ratio, threshold) extend naturally to a second gesture.
- [Derived membership (filter by `stack_id`) means listing a stack's members is O(n) over all items rather than O(1) via a stored list] → Acceptable: this is a personal single-user app with a small localStorage item count; no pagination or indexing exists anywhere else in the app either.
- [Two independent gestures (swipe zone vs. drag handle) on the same card increase UI surface/clutter, especially on small screens] → Mitigated by making the handle small and only relevant when the card is part of a stack-eligible interaction; exact visual treatment is left to implementation/tasks, not blocking design.
- [Deleting the last member of a stack leaves an empty, orphaned stack record in the `stacks` list] → Mitigate by garbage-collecting stacks with zero members at the point membership changes (on remove/delete), same principle as the self-healing membership model above.

## Migration Plan

No server/database migration — this is a purely client-side, additive change to the localStorage schema. Existing items without a `stack_id` are treated as unstacked (`stack_id: null` or absent, treated equivalently). No existing data needs to be transformed; the app already tolerates missing fields on load given its duck-typed item handling. No rollback concerns beyond reverting the frontend code, since no destructive migration of existing localStorage data occurs.
