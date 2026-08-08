## 1. Export

- [x] 1.1 Bump `EXPORT_SCHEMA_VERSION` to `2`
- [x] 1.2 Update `buildExportPayload` to include `stacks: loadStacks()` in the payload

## 2. Import — shape upgrade

- [x] 2.1 Update `upgradeExportData` to return both `items` and `stacks` (defaulting `stacks` to `[]` when the field is absent, e.g. a pre-v2 export file)
- [x] 2.2 Update `importDataFromText` to pass both `items` and `stacks` through to the import routine

## 3. Import — stack id remapping and merge-by-name

- [x] 3.1 Build the `oldStackId -> localStackId` map: for each imported stack, find an existing local stack with an exact (trimmed) name match; reuse its id if found, otherwise create a new local stack via `getNextStackId` and use the new id
- [x] 3.2 When appending a new item during import, rewrite its `stack_id` through the map from 3.1 (items with no `stack_id` are unaffected; items skipped as content-duplicates are unaffected — the existing local item and its `stack_id` are left untouched)
- [x] 3.3 Persist any newly created stacks alongside the existing ones (merged-into stacks are not modified — name/timestamps stay as they are locally)

## 4. Import — result summary

- [x] 4.1 Have the import routine return counts: items added, items skipped as duplicates, stacks added, stacks merged
- [x] 4.2 Show an `alert()` with those counts after a successful import, matching the existing failure-path `alert()` idiom in `handleImportFileChange`

## 5. Tests (`srv/static/script.test.js`)

- [x] 5.1 `exportData`/`buildExportPayload` includes a `stacks` array reflecting current local stacks
- [x] 5.2 `buildExportPayload` includes an empty `stacks` array when no stacks exist
- [x] 5.3 Import restores a stack that has no local counterpart (fresh id assigned, member items linked to it)
- [x] 5.4 Import merges into an existing local stack with the same name (no duplicate stack created, member items linked to the existing stack's id)
- [x] 5.5 Import of an item skipped as a content-duplicate does not alter the existing local item's `stack_id`
- [x] 5.6 Import of a file with no `stacks` field (pre-v2 export) succeeds and imports items normally, with zero stacks created
- [x] 5.7 Import result includes correct counts for items added, items skipped, stacks added, and stacks merged, across a mixed scenario (some of each)

## 6. Verification

- [x] 6.1 Run `make test` (and `make test-js`, which is where this change's coverage actually lives)
- [x] 6.2 Manual pass via the running app: export data with at least one multi-member stack, clear local storage, import the file, confirm the stack and its membership reappear and the summary alert reports correct counts (verified via headless browser against a local build on port 8000, isolated from the live app.service on 8080; summary returned `{itemsAdded:2, itemsSkipped:0, stacksAdded:1, stacksMerged:0}` and the "Errands" stack rendered correctly with both members)
