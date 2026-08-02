# Tasks

- [x] Update `applyImportedItems`/`importDataFromText` in script.js to always append, dropping the `mode` (merge/replace) parameter
- [x] Reassign a fresh id to every imported item (todos and notes) before appending
- [x] Add note-content dedup: skip imported note if its title/content exactly matches an existing note's; existing note (and its review metadata) wins
- [x] Ensure todos are never deduped — always appended regardless of title match
- [x] Remove Replace-mode UI (single "Import" button instead of Merge/Replace buttons) and its confirm() dialog
- [x] Update index.html Help tab copy to describe the new always-append + note-dedup behavior
- [x] Update script.test.js: replace merge/replace-mode tests with append + note-dedup tests; remove obsolete replace-mode tests
- [x] Bump asset version and run `make test` / `node srv/static/script.test.js`
- [x] Commit
