## 1. Code removal

- [x] 1.1 In `archiveItem` (`srv/static/script.js`), remove the `if (!confirm('Abandon this item?')) return false;` line

## 2. Test updates

- [x] 2.1 Update `archiveItem allows abandoning completed todo` — remove or change the `confirms.length === 1` assertion to reflect no dialog fires
- [x] 2.2 Update `abandonFromModal allows archiving done todo` — same fix
- [x] 2.3 Run `make test-js` and `make test`

## 3. Verification

- [x] 3.1 Manual pass via the running app: abandon a todo from its card button and from the edit modal, abandon a note from its edit modal — confirm no dialog appears, the item becomes archived immediately, and Reopen still works (verified via headless browser against a local build on port 8000, isolated from the live app.service on 8080; all three paths archived immediately with no dialog interaction, and reopening a todo worked correctly)
- [x] 3.2 Run `openspec validate --all --strict` to confirm the `todo-list` delta applies cleanly
