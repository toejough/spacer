## Context

The platform forwards a port and starts `/app/run`. What an app publishes is the app's business,
which is why this is a change here rather than something the platform does to this environment.

The guarantee to meet is the platform's: only what has been deliberately published is reachable, and
the public surface is not enumerable. How an app meets it is up to the app — this one keeps its own
Go server, because it serves dynamic routes and a database-backed page that a static file server
cannot.

What has to change is narrower than the other apps needed: not the whole layout, but where the
static content comes from and how the binary finds it.

## Goals / Non-Goals

**Goals:**

- Only deliberately published assets are reachable; `/static/` is not enumerable.
- The binary works wherever the repository sits, so a move cannot 404 the app again.
- A rebuild under systemd actually rebuilds.

**Non-Goals:**

- Switching to the platform's default server. It is a static file server and this app has routes
  and a database.
- Changing the app's routes, templates, or behaviour. Nothing a user sees should change.
- Moving the repository. The platform already did that; `/app/repo` is where it lives.

## Decisions

### 1. Serve `/app/public`, and publish into it by symlink

`StaticDir` becomes `/app/public` rather than a directory inside the source tree. Assets are
published by linking them:

```sh
ln -s ../repo/srv/static/script.js  /app/public/script.js
ln -s ../repo/srv/static/style.css  /app/public/style.css
```

Two things follow. The published set is one directory that can be read at a glance, and it is the
same directory every other app on this platform publishes from — so the answer to "what does this
app expose" does not depend on which app you are looking at. And `script.test.js` stops being served
without anyone having to remember it, because nobody will link it.

The link resolves at request time, so editing an asset in the repository is live on the next request.

### 2. Stop deriving paths from the compile-time source location

`runtime.Caller(0)` returns the path the source file had when the binary was built. That is not
where the source is at run time, and the difference is invisible until something moves — at which
point every page 404s and the cause looks like routing.

Embedding templates with `embed.FS` is the durable answer: they travel inside the binary and cannot
be missing. If they must stay on disk, resolve from a path known at run time — an absolute path, or
one relative to the executable — never from where the compiler saw the source.

### 3. Make the build environment explicit

`/app/run` exports `HOME` before building. Without it, systemd's empty environment leaves `go build`
unable to resolve `GOPATH`; the build fails, `make` reports the error, and the unit starts whatever
binary is already there. Nothing looks broken until the stale binary is wrong about something —
which is exactly how the path problem above stayed hidden.

## Risks / Trade-offs

- **The app has a database and real state** → nothing here touches it, but verify the page renders
  and data is intact after the change, not just that the process starts.
- **A missed asset 404s the page's styling** → enumerate what `/static/` currently serves and link
  each deliberately, rather than discovering omissions in the browser.
- **Embedding templates changes the build** → it also removes a whole class of "works here, not
  there". Worth it, but it is a real change to how the binary is produced.
- **Publishing an asset becomes a second step** → the friction is the guarantee.

## Open Questions

- Whether `/static/` should remain the URL prefix once content comes from `/app/public`. Keeping it
  avoids touching templates; dropping it makes the URL match the directory.
