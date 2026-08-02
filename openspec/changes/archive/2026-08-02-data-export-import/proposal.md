# Data export and import

## Why

All todo/note data lives in browser `localStorage`, keyed to a single origin. There is no login and no server-side persistence of items today. If the app moves to a new hosting URL, or a user switches browsers/devices, or clears site data, they lose everything with no way to recover it or move it. Users need a way to back up and restore their data without requiring an account.

Of the two options considered — (a) local export/import as a JSON file, and (b) letting each user configure their own cloud sync target (e.g. their own Google Drive) — option (a) is far simpler to build and ship now: no OAuth, no server-side storage, no per-user credentials, and it still fully satisfies "move hosting locations without losing data" and "no login required". Cloud-sync-to-user-owned-storage remains a good future enhancement but is out of scope for this change.

## What

Add a "Backup & Restore" section (in the Help tab) with two actions:

- **Export**: downloads all local data (todos + notes, including spaced-repetition and cloze state) as a single JSON file the user can save anywhere (local disk, their own cloud drive folder synced by the OS, email attachment, etc.).
- **Import**: lets the user pick a previously exported JSON file and load it back in. Import supports two modes:
  - **Replace**: wipes current local data and replaces it with the imported data (for moving to a new host/browser with a clean slate).
  - **Merge**: adds imported items that don't already exist (by id) to the current data, leaving existing items untouched (for combining data from two sources).

The exported file is versioned (`schema_version`) and includes the app's storage key/shape so future format changes can be migrated on import.

No server or account changes are required — this is entirely client-side, operating on the same `localStorage` data the app already uses.

## Impact

- Affected capability: `data-portability` (new)
- Affected code: `srv/static/script.js` (export/import logic), `srv/templates/index.html` (UI in Help tab), `srv/static/style.css` (minor styling), `srv/static/script.test.js` (tests)
