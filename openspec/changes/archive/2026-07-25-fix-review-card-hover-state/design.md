# Design: Fix review card hover/focus state carry-over

## Problem

When a user clicks a rating on one review card, `submitReview` removes that item from the due list and calls `loadReview()` to re-render the remaining cards. Because the browser transfers CSS `:hover` state to whatever element is under the pointer after a DOM swap, the next card's rating button at the same coordinates immediately appears visually selected. If the active element was a review button, some browsers also shift focus to the next focusable button at the same coordinates. The result is that the next card looks like it has a preselected rating (e.g., the "Good" option already filled in), which can mislead the user and cause accidental clicks.

## Solution

Remove the CSS `:hover` rule entirely for review buttons and replace it with a JavaScript-driven `hovered` class managed by `mouseenter`/`mouseleave` listeners. Because newly rendered buttons never have the `hovered` class, they can never look preselected when they appear under a stationary pointer.

### Changes in `srv/static/script.js`

1. Add `attachReviewButtonHover(container)` helper that selects `.review-btn` elements inside the given container and adds `mouseenter`/`mouseleave` listeners to toggle the `hovered` class.

2. Call `attachReviewButtonHover(container)` inside `loadReview()` after rendering the cards, so every freshly rendered card gets its own independent listeners.

3. Call `attachReviewButtonHover(card)` inside `revealCloze()` after replacing the buttons, so cloze reveal buttons also get independent listeners.

4. In `submitReview()`, after saving the updated item, blur `document.activeElement` if it is a review button. This prevents focus from jumping to the next card's button at the same coordinates. Then call `loadReview()` as before.

### Changes in `srv/static/style.css`

Replace all `:hover` rules for `.review-btn` with equivalent `.hovered` rules. For example:

```css
.review-btn.hovered{color:#fff}
.review-btn.r0.hovered{background:var(--danger)}
.review-btn.r1.hovered{background:#e57373}
...
```

Remove the `.review-pointer-locked` rule and any `:hover` rules on `.review-btn`.

### Tests

Add a regression test in `srv/static/script.test.js` that calls `submitReview` with a stubbed DOM and asserts that the re-rendered HTML does not contain the `hovered` class. Existing SM-2 and review-entry tests continue to pass.

## Out of scope

- The inactive `web-app/` Quasar prototype is not changed.
- Review scheduling, SM-2 logic, and data model are unchanged.
