# Tasks

- [x] Add `exportData()` in script.js: build `{schema_version, exported_at, items}` from `loadItems()` and trigger a file download
- [x] Add `importData(file, mode)` in script.js: parse JSON, validate shape, then either replace or merge into `localStorage`
- [x] Add id-collision merge logic (existing item wins) and post-import `refreshCurrent()` / count updates
- [x] Add Export/Import UI (buttons + hidden file input + mode choice) in the Help tab of index.html
- [x] Add confirmation dialog for destructive Replace import
- [x] Add user-visible error handling for invalid/malformed import files
- [x] Add CSS for the new Backup & Restore section if needed
- [x] Add/extend tests in script.test.js for export shape, merge behavior, replace behavior, and invalid-file handling
- [x] Bump asset version and run `make test`
- [x] Update README/help text to mention backup & restore
- [x] Commit
