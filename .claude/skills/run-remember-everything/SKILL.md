---
name: run-remember-everything
description: Build, run, and drive the Remember Everything web app (Go server + vanilla-JS SPA) in a headless browser — launch the server, navigate/click/fill/drag with a Playwright-based driver, take screenshots, check console errors. Use when asked to run, test, screenshot, or verify UI behavior (including drag-and-drop) for this app.
---

Remember Everything is a Go HTTP server (`cmd/srv`) that serves a
vanilla-JS single-page app (`srv/static/`, `srv/templates/index.html`).
All app data lives in browser `localStorage` — there is no backend API
for items. It's driven here with a small Playwright-based headless
Chromium REPL (`driver.mjs`) rather than `chromium-cli`, which isn't
installed in this environment; the driver's command set mirrors it
closely (`nav`, `wait-for`, `click`, `fill`, `press`, `screenshot`,
`console-errors`), plus one app-specific addition: `drag`, needed
because this app's drag-and-drop is hand-rolled pointer events, not
HTML5 DnD — Playwright's built-in drag helpers don't fire the events
this app listens for.

All paths below are relative to the repo root (`/app`).

## Prerequisites

None beyond what's already in this container — Node.js, and a cached
Playwright Chromium at `~/.cache/ms-playwright/chromium-1234` (verified
present; if missing on a fresh machine, run
`npx playwright install chromium`).

## Build

```bash
make build      # go build -o todo-srv ./cmd/srv
```

## Run (agent path)

1. Install the driver's dependencies once (already done in this repo,
   committed via `.claude/skills/run-remember-everything/package.json`
   / `package-lock.json` — re-run only if `node_modules` here is
   missing):

   ```bash
   cd .claude/skills/run-remember-everything && npm install --no-audit --no-fund
   ```

2. Launch the server in the background and wait for it to serve:

   ```bash
   ./todo-srv >/tmp/srv.log 2>&1 &
   echo $! > /tmp/srv.pid
   timeout 20 bash -c 'until curl -sf http://localhost:8080/ >/dev/null; do sleep 0.5; done'
   ```

3. Drive it by piping a command script to the driver:

   ```bash
   cd .claude/skills/run-remember-everything
   node driver.mjs <<'EOF'
   nav http://localhost:8080/
   click .tab[data-tab="todos"]
   wait-for #todoAddInput
   fill #todoAddInput Buy milk
   click .quick-add button
   screenshot 01-todo-added
   console-errors
   EOF
   ```

   Screenshots land in `.claude/skills/run-remember-everything/screenshots/<name>.png`.

4. Stop the server when done:

   ```bash
   kill $(cat /tmp/srv.pid)
   ```

### Driver commands

| Command | Effect |
|---|---|
| `nav <url>` | Navigate |
| `wait-for <selector>` | Wait for visible (prefix `text=` for text match) |
| `click <selector>` | Click |
| `fill <selector> <text...>` | Fill an input/textarea |
| `press <key>` | Keyboard press (e.g. `Enter`, `Escape`) |
| `screenshot <name>` | Save `screenshots/<name>.png` |
| `console-errors` | Print collected console/page errors as JSON |
| `eval <js>` | `page.evaluate(...)`, prints JSON result |
| `sleep <ms>` | Pause |
| `drag <fromSelector> <toSelector>` | App-specific: simulates a real mouse `down` → `move` (several steps) → `up` between two elements' centers. Required for this app's stack drag-and-drop (`.drag-handle[data-drag-id="N"]` → `.item-card[data-id="M"]` or `[data-stack-id="N"]`), since it's pointer-event-based, not HTML5 DnD. |

Blocking modals (stack-name prompt, edit modal) are plain DOM elements
— `fill`/`click` them directly, e.g. `fill #stackNameInput Errands`
then `click #stackNameModal .btn-primary`.

## Run (human path)

```bash
./todo-srv
# open http://localhost:8080 in a browser
```

## Test

```bash
make test      # Go tests
make test-js   # node srv/static/script.test.js — no browser needed
```

## Gotchas

- **State doesn't persist across separate `driver.mjs` invocations.**
  Each run launches a fresh, non-persistent browser context, so
  `localStorage` (where all app data lives) resets every time. Do a
  whole scenario — add items, drag, verify — in **one** script/session,
  not across multiple `node driver.mjs` calls.
- **Playwright's built-in `dragTo`/HTML5 drag events don't work here.**
  The app implements stacking drag-and-drop with raw `pointerdown` /
  `pointermove` / `pointerup` listeners (see
  `srv/static/script.js` `attachDragHandlers`), not the HTML5 Drag and
  Drop API. Use the driver's `drag` command (manual `mouse.down` →
  `mouse.move` in steps → `mouse.up`), not Playwright's `dragTo`.
- **The drag handle is a distinct hit-target from the rest of the
  card.** Target `.drag-handle[data-drag-id="N"]` specifically, not
  `.item-card[data-id="N"]`, when simulating a drag start — dragging
  from elsewhere on a Done/Abandoned card triggers swipe-to-delete
  instead (that's intentional; see the todo-list spec's swipe scenarios).
- **`npm install` in this skill dir needs
  `PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD` unset only if the Chromium cache
  is missing.** It's already cached in this container; installing
  plain `playwright` (not `@playwright/test`) reused it with no
  download.

## Troubleshooting

- `Error: Cannot find module '.../driver.mjs'` when using a heredoc
  file argument — write the script with the `Write` tool (or `cat >
  file <<'EOF'` as its own, single Bash call) before invoking `node
  driver.mjs <file>`; chaining `cat > file <<EOF ... && node ...` in
  one composite command intermittently produced a missing file in this
  environment. Piping directly via stdin (`node driver.mjs <<'EOF' ...
  EOF`) is the more reliable form and is what's recommended above.
- `drag: selector not found` — the target selector doesn't exist yet
  (e.g. you tried to drag onto a card before adding it, or state was
  lost because the script ran in a fresh `driver.mjs` invocation — see
  Gotchas above).
