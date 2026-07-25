# Fix review card hover/focus state carry-over

## Why

When a user clicks a rating on a review card, the card is removed and the next card renders immediately under the same pointer location. The browser transfers the hover and focus state to the new card's rating button at that position, making the "Good" option on the next card appear already selected. This is misleading and can lead to accidental ratings.

## What Changes

- `srv/static/script.js`: After a review submission, blur the active rating button and temporarily lock pointer events on the `#reviewCards` container while it re-renders. This prevents hover/focus state from carrying over to the freshly rendered card. The lock is released on the next `pointermove` or after 50 ms.
- `srv/static/style.css`: Add a `.review-pointer-locked .review-buttons` rule that disables pointer events while the lock is active.
- `srv/static/script.test.js`: Add a regression test verifying that `submitReview` applies the pointer-event lock to the review container.
- `srv/templates/index.html` and `srv/static/sw.js`: Bump version from `v18` to `v19` so the service worker and footer reflect the updated assets.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `spaced-review`: Adds a requirement that each review card renders independently and does not inherit hover/focus state from a previously reviewed card.

## Impact

- Affected code: `srv/static/script.js`, `srv/static/style.css`, `srv/static/script.test.js`, `srv/templates/index.html`, `srv/static/sw.js`.
- No API or database changes.
- No impact on the inactive `web-app` prototype.
