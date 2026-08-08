## 1. Code removal

- [x] 1.1 Remove the `Hostname` field from the `Server` struct in `srv/server.go`
- [x] 1.2 Change `srv.New` to take no arguments (drop the `hostname string` parameter); update its body accordingly
- [x] 1.3 In `cmd/srv/main.go`'s `run()`, remove the `os.Hostname()` call, its error-fallback branch, and the now-unused `os` import if nothing else in the file needs it (kept the `os` import — `os.Stderr` is still used in `main()`)

## 2. Verification

- [x] 2.1 Run `go build ./...` to confirm it still compiles
- [x] 2.2 Run `make test`
- [x] 2.3 Run the binary and confirm it still starts and serves normally (`curl localhost:<port>/`) (verified on an isolated local build, port 8000; got HTTP 200; live production on 8080 untouched)
