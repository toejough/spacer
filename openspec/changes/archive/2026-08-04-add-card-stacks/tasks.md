## 1. Data model

- [x] 1.1 Add `stack_id` field to item objects (default `null`), threaded through wherever items are created/loaded/saved in `srv/static/script.js`
- [x] 1.2 Add `stacks` localStorage key and `loadStacks()`/`saveStacks()` helpers, mirroring the existing `loadItems()`/`saveItems()` pattern
- [x] 1.3 Add a helper to compute a stack's members from `stack_id` (derived, not stored as a list) and its visible-type count for a given tab/type filter
- [x] 1.4 Add garbage-collection of stacks with zero members whenever membership changes (remove-from-stack, delete item)

## 2. Sort/render integration

- [x] 2.1 Extend the per-tab render path (Todos, Notes, Search) to rank all visible-type items via the existing relevance comparator, then collapse consecutive-by-rank stack members into a single collapsed-tile entry at the first (most urgent) occurrence
- [x] 2.2 Build the collapsed stack tile component: name, visible-member count, expand/collapse control
- [x] 2.3 Wire expand/collapse state (per stack, per view) so expanding a stack renders its visible-type member cards individually, using the existing `renderTodoCard()`/`renderNoteCard()`
- [x] 2.4 Ensure a stack with zero visible-type members is omitted from a given tab (Todos/Notes), matching the type-filtered rendering rule
- [x] 2.5 Extend search rendering to show mixed-type stacks as collapsed tiles when any member matches, expanding to show all matching members regardless of type

## 3. Drag-and-drop infrastructure

- [x] 3.1 Add a drag-handle element to each rendered card, visually and functionally distinct from the swipe-to-delete hit area
- [x] 3.2 Implement pointer-event-based drag tracking from the handle (pointerdown/pointermove/pointerup), producing a ghost/drag-preview element, following the existing swipe-gesture's pointer-event conventions
- [x] 3.3 Implement drop-target detection: resolve the element under the pointer on release to either "another unstacked card," "a stack's collapsed tile," "a stack member card," or "outside any stack"
- [x] 3.4 Wire drag-and-drop for desktop pointers as well as touch, verifying no regression to the existing swipe-to-delete gesture on Done/Abandoned cards — fixed a real conflict found during review: swipe's `pointerdown` handler now ignores presses that originate on `.drag-handle`, so grabbing the handle on a Done/Abandoned card no longer also arms swipe-to-delete.

## 4. Stack creation, membership, naming

- [x] 4.1 Build the stack-name prompt modal (reusing the existing modal-overlay pattern), with non-empty-name validation and Cancel
- [x] 4.2 Wire "drop card onto unstacked card" → open stack-name prompt → on confirm, create stack and set both cards' `stack_id`; on cancel, no-op
- [x] 4.3 Wire "drop card onto a stack (or a card already in a stack)" → add dropped card to that stack, no prompt
- [x] 4.4 Wire "drag a stacked card out and drop outside its stack" → clear that card's `stack_id`, triggering stack GC if it was the last member
- [x] 4.5 Confirm dragging one stack's tile onto another stack's tile is a no-op (explicitly not implemented as a merge) — satisfied by construction: only cards render a drag handle, collapsed stack tiles don't, so a stack can never be the dragged element.
- [x] 4.6 Add inline rename control on the collapsed stack tile (edit-in-place with confirm/cancel)

## 5. Edit modal integration

- [x] 5.1 Add a stack field to the edit modal showing the card's current stack (or "no stack")
- [x] 5.2 Populate the stack field's options from existing stacks; allow selecting "no stack" to remove membership
- [x] 5.3 On save, apply the selected stack change to the card's `stack_id` (triggering stack GC if applicable) and refresh the relevant list
- [x] 5.4 Add stack rename access from the edit modal for cards that are stack members

## 6. Styling

- [x] 6.1 Style the drag handle, collapsed stack tile, expand/collapse control, and inline rename input in `srv/static/style.css`, consistent with existing card styling
- [x] 6.2 Style drag-in-progress state (ghost element, valid/invalid drop-target highlighting)

## 7. Verification

- [x] 7.1 Manually verify all scenarios in `specs/card-stacks/spec.md`, `specs/todo-list/spec.md`, `specs/note-taking/spec.md`, `specs/search/spec.md` — verified in a real headless browser via the new `.claude/skills/run-remember-everything` Playwright driver: added two todos, dragged one onto the other (name prompt appeared, cancel/confirm both work), confirmed the resulting stack collapses to a tile with a visible-member count, expands to show both member cards, dragged a member back out and confirmed it un-stacked (count dropped to 1, card became a normal loose card again), and confirmed the edit modal's Stack field lists "No stack" plus existing stacks. No console errors in any of these flows.
- [x] 7.2 Verify existing swipe-to-delete behavior is unaffected by the new drag handle — verified in-browser: a leftward mouse drag starting on `.drag-handle` no longer deletes the card (fix from 3.4 confirmed live), while the same leftward drag starting on the card body still swipes-to-delete and shows the Undo toast, exactly as before.
- [x] 7.3 Verify existing relevance-sort behavior for unstacked items is unaffected — confirmed via the full existing JS test suite passing unmodified (`make test-js`), plus new tests asserting unstacked items still appear as individual entries in `buildDisplayList`.
- [x] 7.4 Run `make test` and add/update Go tests only if server-side code changed (expected: none, per design's no-backend-changes decision) — `make test` and `make test-js` both pass; no Go/backend files were touched.
