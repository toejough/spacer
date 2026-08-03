## 1. Data layer

- [x] 1.1 Add `data-id` attributes to rendered todo and note cards (needed for gesture delegation to resolve the target item)
- [x] 1.2 Add a pending-delete buffer: on commit, hold the full item object + a timer instead of immediately removing it from the saved `localStorage` array
- [x] 1.3 Add `undoDelete(id)` — restores the item from the pending buffer and re-renders
- [x] 1.4 On timer expiry, actually remove the item from the saved array

## 2. Swipe gesture

- [x] 2.1 Add delegated pointer/touch event handling on `#todoList`/`#noteList` (not per-card listeners, since cards are rebuilt on every refresh)
- [x] 2.2 Implement horizontal-vs-vertical disambiguation (deadzone + ratio threshold) so the gesture doesn't fight list scrolling
- [x] 2.3 Only arm the swipe gesture on cards whose item is Done or Abandoned (todos) / Abandoned (notes) — Open/active cards never engage it
- [x] 2.4 On threshold-cross, commit the pending delete and show the Undo toast
- [x] 2.5 Visual drag feedback (reveal/color treatment) during an in-progress swipe, and a snap-back if released before threshold — background tints red as soon as the swipe is recognized (`.swiping`), intensifying to a stronger red + danger border once past the threshold (`.swipe-armed`)
- [x] 2.6 Reveal-behind panel: deletable cards (Done/Abandoned todos, Abandoned notes) are wrapped in `.item-card-wrapper` with an absolutely-positioned `.swipe-reveal` (red background + trash icon) behind them, exposed in the space the card vacates as it slides — matches the standard "sliding drawer" swipe-to-delete pattern (researched: LogRocket's swipe-to-delete/swipe-to-reveal guidance, DesignMonks' delete-button UI practices). Ineligible cards render unwrapped, with no reveal markup at all.
- [x] 2.7 Fix: `.swipe-reveal` bled through at rest because `.item-card.done`/`.abandoned` (the only cards eligible for the reveal) already carry reduced `opacity`, making the whole card translucent. Fixed by hiding `.swipe-reveal` (`opacity:0`) except when the wrapper carries a `.swipe-active` class, toggled by JS only during an active drag, and forcing the dragged card to full opacity while `.swiping`/`.swipe-armed`.
- [x] 2.8 Committed to swipe-left-only (the standard convention — reveals actions on the right); a rightward drag is now ignored like a vertical scroll. This was needed for the trash icon to be positioned consistently: it's now centered within the final revealed strip (width = `SWIPE_THRESHOLD_PX`, wired to CSS via `--swipe-threshold-px` so JS/CSS can't drift), not the full card width, so it lands centered in the exposed area exactly at the swipe threshold.
- [x] 2.9 Fix: any vertical scroll on an already-armed card was committing the delete. Cause: `pointercancel` (fired when the browser hands the gesture off to native scroll) was wired to the same handler as `pointerup`, so an interrupted drag past threshold committed exactly like a deliberate release. Split into `releaseGesture` (pointerup only — commits if past threshold, else snaps back) and `abortGesture` (pointercancel/pointerleave — always snaps back, never commits, regardless of drag distance). Verified by dispatching a real `pointercancel` event mid-drag past threshold in a live browser and confirming the item survives and the card snaps back.
- [x] 2.10 Fix: swipe still canceled immediately on any vertical motion even after 2.9. Root cause was deeper — `touch-action` was only toggled to `none` by JS after our own deadzone/ratio logic decided a gesture was horizontal, but browsers decide scroll-vs-JS-gesture ownership from the very first pixels of a touch sequence based on whatever `touch-action` the element already has; by the time our JS reacted, the browser had often already committed to native scroll and fired `pointercancel`. Fixed by setting `touch-action: pan-y` **statically** in CSS (not toggled) on swipeable cards — this reserves vertical panning for native scroll and excludes horizontal panning from native handling, so the browser's own arbitration does the disambiguation instead of a race with our JS. Verified with real touch input via CDP's `Input.dispatchTouchEvent` (synthetic DOM events don't trigger this arbitration, so a real touch pipeline was required): a pure vertical touch drag now never arms/engages the gesture, and a real horizontal touch swipe still deletes correctly end-to-end (immediate visual removal + undo toast, permanent removal after the undo window).

## 3. Undo toast

- [x] 3.1 Add toast markup/CSS with an Undo action and a visible or implicit expiry
- [x] 3.2 Wire the toast's Undo button to `undoDelete`
- [x] 3.3 Toast auto-dismisses and the delete becomes permanent when the window elapses

## 4. Keyboard/screen-reader parity

- [x] 4.1 Add a "Delete Forever" button to the edit modal markup in `srv/templates/index.html`
- [x] 4.2 Add a modal update function (sibling to `updateModalAbandonButton`) that shows the button only when the item is Done/Abandoned, hides it for Open/active
- [x] 4.3 Wire the button to commit the delete — reuses the same pending-delete/undo-toast path as the swipe, for one consistent delete mechanism and to keep the undo safety net available from the modal path too — and closes the modal

## 5. Verify

- [x] 5.1 Swipe-delete works on a Done todo, an Abandoned todo, and an Abandoned note — verified the underlying `deleteItemPending`/list-filtering logic via a VM-harness script, AND separately verified the real pointer-drag DOM interaction visually: installed Playwright/Chromium in this environment, ran the actual server, and drove real mouse-drag gestures against the live page, screenshotting at rest / mid-swipe / armed / post-delete-with-toast. Caught and fixed two visual bugs this way (2.7, 2.8) that the VM-harness alone couldn't have found. Real touch-device (as opposed to mouse-emulated) behavior still wasn't checked.
- [x] 5.2 Open todos and active notes never arm the swipe gesture or show "Delete Forever" — verified `itemAllowsDelete` gating and modal button visibility for both allowed and disallowed states
- [x] 5.3 Undo restores the exact prior item state, including `cloze_data` for multi-cloze items — verified: storage is byte-identical after delete+undo, since nothing is written until the timer actually fires
- [x] 5.4 After the toast expires, the item is actually gone from `localStorage` (not just visually removed) — verified via `commitDelete` directly and a real 5-second `setTimeout` end-to-end
- [x] 5.5 "Delete Forever" in the edit modal works via keyboard alone (tab + enter), no pointer/touch required — verified the button's visibility and click-invocation logic (equivalent to what Enter/Space would trigger on a focused button); real keyboard-navigation/focus-order still needs a manual check in an actual browser
- [x] 5.6 `make test-js` passes — all 24 existing tests pass
- [x] 5.7 Bump asset version in `srv/templates/index.html` (v33 -> v34)
