# Tasks

- [x] Generalize `noteContentKey` in script.js to a content key based on `item_type + title` for any item, not just notes
- [x] Update `applyImportedItems` to dedupe both todos and notes against existing items by that content key
- [x] Ensure existing local item (with all its metadata) always wins on a content match, regardless of item type
- [x] Fresh id assignment for non-duplicate imported items stays as-is
- [x] Update Help tab copy in index.html to describe dedup applying to todos and notes alike
- [x] Update script.test.js: adjust/replace the "todos always append" tests with todo-dedup tests; keep note-dedup tests
- [x] Bump asset version and run `make test` / `node srv/static/script.test.js`
- [x] Commit
