# 018 — PWA doesn't work offline

**Status:** done
**Type:** bug

## Observed Behavior

With the dev server stopped:
1. App was open and showing the deck list
2. Clicked on a deck — nothing happened (navigation failed silently)
3. Refreshed the page — app died completely

This is a website, not a PWA. A PWA should work offline after first load.

## Expected Behavior

With the server unavailable, the app should:
- Navigate between already-loaded views
- Continue to read/write to IndexedDB (Dexie works offline by nature)
- Survive a page refresh (serve from service worker cache)

## Root Cause

There is no service worker, no web app manifest, and no offline caching configured. The app has zero PWA infrastructure — it's a standard SPA that requires the server for every page load and lazy-loaded route chunk.

Specifically:
- No `manifest.json` / `manifest.webmanifest`
- No service worker (no `vite-plugin-pwa` or equivalent)
- Lazy-loaded route chunks (`() => import("./views/...")` in `main.ts`) fail when the server is unreachable — this is why clicking a deck did nothing
- On refresh, the browser requests `index.html` from the server, which is down — app dies

## Acceptance Criteria

- [ ] Web app manifest present with name, icons, theme color, `display: standalone`
- [ ] Service worker caches the app shell (HTML, JS, CSS) on first load
- [ ] All route chunks are precached so navigation works offline
- [ ] Page refresh works offline (serves index.html from SW cache)
- [ ] IndexedDB data remains accessible offline (should already work, just verify)
- [ ] Installable as a PWA on mobile and desktop
