# Import is always append, with note-content dedup

## Why

The original data-portability change offered two import modes: Merge (add by id) and Replace (destructive wipe+restore). In practice, id-based merge is fragile — ids are per-browser autoincrement values, so two independently-created exports can easily collide on id purely by coincidence, causing legitimate distinct items to be silently dropped or, worse, comparing unrelated items as "the same". Destructive Replace is also risky to keep as a default option given how easy it is to click.

What actually matters for this app: notes are knowledge you don't want duplicated (the same fact reviewed twice is confusing and pollutes spaced-repetition stats), while todos are fine to have duplicates of (e.g. a recurring todo re-added from an old backup is harmless, and there's no reliable way to tell "the same todo" from "a new todo with the same title").

## What

Simplify to a single import behavior: **Import is always an append.**

- All imported items are added to local storage; nothing existing is ever removed or replaced.
- Item ids are reassigned on import so they never collide with existing local ids.
- The only de-duplication is for **notes**: if an imported note's content exactly matches an existing note's content, the note is *not* added — the existing local note (with its own review metadata: next_review, ease_factor, interval, repetitions, cloze_data, etc.) wins and is left untouched.
- Todos are never de-duplicated — every imported todo is appended, even if its title matches an existing todo.
- Remove the "Replace" import mode and its confirmation dialog; there is only one Import action now.

## Impact

- Affected capability: `data-portability`
- Affected code: `srv/static/script.js` (`applyImportedItems`, `triggerImport`, remove replace-mode logic), `srv/templates/index.html` (single Import button instead of Merge/Replace), `srv/static/script.test.js`
