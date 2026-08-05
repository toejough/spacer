# Tasks

**Two things move under your feet while you run this**, so read them first:

- Your shell already lands in `/app/repo` — that is where this repository lives now, and where
  `openspec` finds its root.
- The platform starts `/app/server/bin/serve` with no arguments. `/app/run` still works today but is
  being retired; task 3.3 moves this app's build-and-start script to the new path.

This app has a database and is the one the platform's fleet migration deliberately did not rewrite.
Change the serving, not the app.

## 1. Record what is served now

- [x] 1.1 Enumerate every asset `/static/` currently serves, before changing anything — that list is
  what has to be published, and discovering an omission in the browser is worse than reading it off
  now. `/static/` returns a listing today, which makes this easy and is itself the bug
- [x] 1.2 Confirm the current gaps by request so the fix can be measured: `/static/` returns 200
  with an index, and `/static/script.test.js` returns 200

## 2. Publish deliberately

- [x] 2.1 `StaticDir` becomes `/app/public`
- [x] 2.2 Every asset from 1.1 that belongs on the web is symlinked into `/app/public`. Test files
  are not — `script.test.js` should stop being reachable because nobody linked it, not because
  something refused it
- [x] 2.3 No directory listing for any path. Go's `http.FileServer` lists a directory with no
  `index.html`; that is what makes `/static/` browsable
- [x] 2.4 `Cache-Control: no-store` still goes out on assets. The platform sits behind a CDN that
  caches static extensions for hours, which hides edits
- [x] 2.5 `srv/server.go` — `StaticDir` and `TemplatesDir` fields, plus any `Dir` on the file server,
  are updated to use `/app/public` with no directory listing

## 3. Make the binary independent of where it was built

- [x] 3.1 Templates are embedded with `embed.FS` so they travel inside the binary. Static paths
  similarly stop using `runtime.Caller(0)` and instead embed or resolve from a runtime-known path
- [x] 3.2 Prove it: build the binary, move the repository to a different directory, run it from
  there, and confirm pages still render. This exact failure already happened once — the repository
  moved from `/app` to `/app/repo` and the running binary 404'd everything until it was rebuilt
- [x] 3.3 **The build+start script moves to `/app/server/bin/serve`.** The platform starts that path
  with no arguments and is retiring `/app/run`; a script there is a first-class answer, not a
  workaround. The four lines move unchanged:

  ```sh
  #!/bin/sh
  export HOME=/root      # systemd provides none, and `go build` needs it for GOPATH
  cd /app/repo
  make build
  exec ./todo-srv -listen 0.0.0.0:8080
  ```

  `HOME` is the load-bearing line: without it the build fails, `make` reports the error, and the unit
  silently starts whatever stale binary already exists — which is how a binary compiled at the old
  repository path kept 404ing after the move. **Leave `/app/run` in place** pointing at
  `exec /app/server/bin/serve`, so both work until the platform's `bin-serve-is-the-entry-point`
  change lands. That change is blocked on this task
- [x] 3.4 Replacing `bin/serve` means the host's `fleet/check-baseline` will report this app's copy as
  **locally modified** from now on, and `fleet/offer` will refuse to overwrite it. Correct and
  expected: the fleet is recording that this app runs its own server

## 4. Verify

- [x] 4.1 `/static/` returns no index; `/static/script.test.js` is 404; every published asset is 200
- [x] 4.2 The page renders with styling intact, and the database still has its data
- [x] 4.3 A file added to the repository and not published is unreachable — the property is
  "unpublished", not "refused"
- [x] 4.4 Editing a published asset is live on the next request, with no rebuild or restart
