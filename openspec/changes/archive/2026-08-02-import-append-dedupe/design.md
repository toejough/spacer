## Design notes

- Notes currently store their text in the `title` field (the `content` field exists in the data model but is unused by the UI today). "Duplicate note content" dedup compares `title` for items where `item_type === 'note'`, trimmed of surrounding whitespace, exact match (case-sensitive) for the initial implementation.
- Todos are matched only by `item_type === 'todo'`; they are excluded from the dedup check entirely and always appended.
- Because ids are reassigned on import, any references to old ids in the imported file's own data (there are none today) would need no special handling.
