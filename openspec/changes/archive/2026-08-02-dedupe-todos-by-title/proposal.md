# Dedupe todos by content on import, same as notes

## Why

The previous change made import always-append for todos, on the assumption that todos have no reliable notion of "the same item." That was wrong: the same rule that applies to notes applies to todos too. In this app, an item's **content** is its title (including any cloze markup) plus its type (todo vs note) — everything else (done/archived state, priority, due date, tags, spaced-repetition metadata, timestamps) is just metadata about that content. Two items with the same type and the same title are the same item, whether they're todos or notes.

## What

Apply the existing note dedup rule uniformly to both item types on import:

- An imported item is skipped if there is already a local item with the same `item_type` and the same title (trimmed, exact match) — this now applies to todos as well as notes.
- When skipped, the existing local item wins and keeps all of its own metadata (done/archived state, priority, tags, due date, spaced-repetition fields, timestamps, etc.), unchanged.
- Items that don't match anything existing (by type+title) are appended with a freshly assigned id, as before.
- Export/import remains login-free and purely client-side; no server changes.

## Impact

- Affected capability: `data-portability`
- Affected code: `srv/static/script.js` (`applyImportedItems`, rename/generalize `noteContentKey` to a type+title content key), `srv/static/script.test.js`, `srv/templates/index.html` (Help tab copy)
