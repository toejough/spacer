## 1. Code removal

- [x] 1.1 In `quickAdd` (`srv/static/script.js`), remove `content`, `priority`, `due_date`, and `tags` from the item object literal
- [x] 1.2 In `doSearch` (`srv/static/script.js`), remove the `|| i.content.toLowerCase().includes(q)` clause; match on title only
- [x] 1.3 In `baseItem` (`srv/static/script.test.js`), remove `content`, `priority`, `due_date`, and `tags` from the fixture

## 2. Test updates

- [x] 2.1 Update or remove any test assertions that reference `priority` (e.g. the import-dedup test that currently asserts `priority` is preserved) — use a different still-live metadata field to prove the same "existing local item wins" point if the test's purpose still needs one (swapped to `ease_factor`)
- [x] 2.2 Confirm no test asserts on `content`, `due_date`, or `tags`; remove any that do (also fixed a stale comment above `itemContentKey` naming the removed fields)
- [x] 2.3 Run `make test-js` and `make test`

## 3. Verification

- [x] 3.1 Manual pass via the running app: add a todo and a note, search for their titles (still matches), confirm nothing in the UI referenced these fields to begin with (edit modal, cards) so no visible change is expected (verified via headless browser against a local build on port 8000, isolated from the live app.service on 8080; new items carry no content/priority/due_date/tags keys, title search still matches, unrelated queries correctly return no results)
- [x] 3.2 Run `openspec validate --all --strict` to confirm the `search` and `data-portability` deltas apply cleanly
