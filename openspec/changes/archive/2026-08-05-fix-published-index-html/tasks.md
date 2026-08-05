## 1. Publish the entry page correctly

- [x] 1.1 Replace `/app/public/index.html` (hand-written stub) with a symlink: `ln -sf ../repo/srv/templates/index.html /app/public/index.html`
- [x] 1.2 Confirm the other four published assets are untouched and still symlinked: `manifest.json`, `script.js`, `style.css`, `sw.js`

## 2. Restore a compiling, correct `srv/server.go`

- [x] 2.1 Remove unused imports (`io/fs`, `os`, `path/filepath`) left by the abandoned `fix-root-handler-404` attempt
- [x] 2.2 Restore the root route (`GET /{$}`) to an explicit `handleIndex` handler that serves `TemplatesDir + "/index.html"` with `Cache-Control: no-store`, removing the `http.FileServer(http.Dir(TemplatesDir))` no-op
- [x] 2.3 Keep the `/app/public` repointing (`StaticDir`/`TemplatesDir` consts) from `adopt-standard-app-layout` unchanged
- [x] 2.4 `go build ./...` succeeds with no errors

## 3. Verify by loading the page, not by checking status codes

- [x] 3.1 `GET /` returns the full page (line count and byte size matching `srv/templates/index.html`, not the 13-line stub)
- [x] 3.2 `Cache-Control: no-store` is present on the response for `/`
- [x] 3.3 `/static/style.css`, `/static/script.js`, `/static/manifest.json`, `/sw.js` all return 200
- [x] 3.4 `/openspec/`, `/repo/`, `/index.html`, `/static/script.test.js` all still return 404
- [x] 3.5 Load the page in a real (or headless) browser: the app UI renders, `<link rel="manifest">` is present, the service worker registers without console errors, and a todo can be added and persisted
- [x] 3.6 Confirm `app.service` is running the rebuilt binary (not a stale one from before the fix)
