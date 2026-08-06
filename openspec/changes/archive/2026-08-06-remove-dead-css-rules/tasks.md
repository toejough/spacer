## 1. Delete

- [x] 1.1 Delete `.item-content` from `srv/static/style.css`
- [x] 1.2 Delete `.item-tag` from `srv/static/style.css`
- [x] 1.3 Delete `.item-date` from `srv/static/style.css`
- [x] 1.4 Delete `.item-priority` and its `.p1`/`.p2`/`.p3` modifiers from `srv/static/style.css`
- [x] 1.5 (found during implementation, not in the original audit) Delete `.review-card .item-content` — a scoped variant of the same dead `item-content` class, missed because the original audit's grep only checked `srv/templates/index.html`/`srv/static/script.js`, not `style.css` itself for compound selectors reusing the class

## 2. Verify

- [x] 2.1 `grep -n "item-content\|item-tag\|item-date\|item-priority" srv/templates/index.html srv/static/script.js srv/static/style.css` still returns nothing (re-confirm after deletion, not just before; now also checking style.css itself for any remaining compound-selector reuse)
- [x] 2.2 Load the app in a browser: todo cards, note cards, and stacks all still render correctly with no visual regression
- [x] 2.3 `make build && make test` and `make test-js` pass
