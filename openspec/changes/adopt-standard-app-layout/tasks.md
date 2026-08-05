# Tasks

This app has a database and is the one the platform's fleet migration deliberately did not rewrite.
Change the serving, not the app.

## 1. Record what is served now

- [ ] 1.1 Enumerate every asset `/static/` currently serves, before changing anything — that list is
  what has to be published, and discovering an omission in the browser is worse than reading it off
  now. `/static/` returns a listing today, which makes this easy and is itself the bug
- [ ] 1.2 Confirm the current gaps by request so the fix can be measured: `/static/` returns 200
  with an index, and `/static/script.test.js` returns 200

## 2. Publish deliberately

- [ ] 2.1 `StaticDir` becomes `/app/public`
- [ ] 2.2 Every asset from 1.1 that belongs on the web is symlinked into `/app/public`. Test files
  are not — `script.test.js` should stop being reachable because nobody linked it, not because
  something refused it
- [ ] 2.3 No directory listing for any path. Go's `http.FileServer` lists a directory with no
  `index.html`; that is what makes `/static/` browsable
- [ ] 2.4 `Cache-Control: no-store` still goes out on assets. The platform sits behind a CDN that
  caches static extensions for hours, which hides edits

## 3. Make the binary independent of where it was built

- [ ] 3.1 Templates and static paths stop coming from `runtime.Caller(0)`. Embed templates with
  `embed.FS`, or resolve from an absolute path or one relative to the executable
- [ ] 3.2 Prove it: build the binary, move the repository to a different directory, run it from
  there, and confirm pages still render. This exact failure already happened once — the repository
  moved from `/app` to `/app/repo` and the running binary 404'd everything until it was rebuilt
- [ ] 3.3 `/app/run` exports `HOME` before `make build`, and the repo carries that so it survives a
  rebuilt environment. Without it, systemd's empty environment leaves `go build` unable to resolve
  `GOPATH`, the build fails, and the unit silently starts whatever stale binary exists

## 4. Verify

- [ ] 4.1 `/static/` returns no index; `/static/script.test.js` is 404; every published asset is 200
- [ ] 4.2 The page renders with styling intact, and the database still has its data
- [ ] 4.3 A file added to the repository and not published is unreachable — the property is
  "unpublished", not "refused"
- [ ] 4.4 Editing a published asset is live on the next request, with no rebuild or restart
