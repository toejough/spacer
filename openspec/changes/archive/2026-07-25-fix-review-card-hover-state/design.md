# Design: Fix review card hover/focus state carry-over

## Problem

When a user clicks a rating on one review card, `submitReview` removes that item from the due list and calls `loadReview()` to re-render the remaining cards. If the mouse pointer is still over the rating button area, the next card's button at the same screen coordinates immediately receives the browser's `:hover` state. If the active element was a review button, some browsers also shift focus to the next focusable button at the same coordinates. The result is that the next card appears with a rating button already visually selected (e.g., the "Good" option), which looks like a preselected value and can cause accidental clicks.

## Solution

Break the browser's hover/focus association with the old button before the DOM swap, and suppress hit-testing on the new buttons during the brief re-render window.

### Changes in `srv/static/script.js`

1. Add a helper `lockReviewPointerEvents()` that adds the class `review-pointer-locked` to `#reviewCards` and removes it on the next `pointermove` or after a 50 ms timeout (whichever comes first). This ensures the lock is short-lived and does not interfere with normal interaction once the pointer moves.

2. In `submitReview()`, after saving the updated item and before calling `loadReview()`:
   - Blur `document.activeElement` if it is a review button, preventing focus from jumping to the next card's button.
   - Call `lockReviewPointerEvents()` so the re-rendered buttons do not receive pointer events while the DOM is being swapped.

### Changes in `srv/static/style.css`

Add a single rule:

```css
.review-pointer-locked .review-buttons { pointer-events: none; }
```

This disables pointer events for all rating buttons inside the locked review container while still leaving the rest of the page responsive.

### Tests

Add a regression test in `srv/static/script.test.js` that calls `submitReview` with a stubbed DOM and asserts that the `#reviewCards` container receives the `review-pointer-locked` class. Existing SM-2 and review-entry tests continue to pass.

## Out of scope

- The inactive `web-app/` Quasar prototype is not changed.
- Review scheduling, SM-2 logic, and data model are unchanged.
