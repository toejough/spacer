## Why

This app serves named routes and an explicit static directory, so it does not publish its repository
the way the platform's other apps once did. That made it look already-correct during the fleet
migration. Checking rather than assuming found two real gaps:

```
200  /static/                  a browsable directory listing
200  /static/script.test.js    a test file, on the web
```

The static directory is enumerable and holds whatever happens to be beside the site's assets. That
is the same defect the rest of the fleet just fixed, one directory further in.

There is a second, sharper problem. `srv/server.go` resolves its template and static directories
from `runtime.Caller(0)` — the **path the source occupied when it was compiled**. When the platform
moved this app's repository from `/app` to `/app/repo`, the existing binary kept looking at the old
path and 404'd every page. It only recovered on a rebuild. A binary that depends on where it was
built will break again on the next move, and it fails in a way that looks like a routing bug.

The app also builds at startup, and that build had been failing silently under systemd for some
time: systemd provides no `HOME`, so `go build` could not resolve `GOPATH` and the unit fell through
to whatever binary already existed. The fix is in `/app/run` and belongs in this repository.

## What Changes

- **The published set becomes `/app/public`**, the same directory every app on this platform serves
  from. Assets are symlinked into it from this repository, so publishing is deliberate and a file
  nobody linked — a test file, a stray note — is unreachable.
- **No directory listing.** `/static/` stops being browsable.
- **Templates and static content stop depending on the compile-time source path.** Embed them, or
  resolve from a path known at run time.
- **The build environment is explicit**, so a rebuild under systemd succeeds instead of silently
  reusing a stale binary.

## Impact

- `srv/server.go` — `StaticDir`, `TemplatesDir`, listing behaviour.
- `/app/run` — already carries the `HOME` fix; record it here so it survives a rebuild of the env.
- `/app/public/` — currently empty; becomes the published set.
- No platform-side change. The platform starts `/app/run`; this app keeps its own server.
