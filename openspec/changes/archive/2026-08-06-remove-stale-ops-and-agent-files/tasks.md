## 1. Delete stale files

- [x] 1.1 Delete `run`
- [x] 1.2 Delete `issues/0001-ui-abandon-visuals.md` and confirm `issues/` is now empty and remove it
- [x] 1.3 Delete `shelley-hooks/install.sh` and `shelley-hooks/slash/opsx`, confirm `shelley-hooks/` is now empty and remove it

## 2. Update AGENTS.md

- [x] 2.1 Remove the two Shelley hook lines ("A Shelley slash hook is installed..." and "If the hook is missing, install it with...")
- [x] 2.2 Remove or rewrite the "Queued work" section — `adopt-standard-app-layout` is archived, not current; confirm via `openspec list` that nothing is actually queued before deciding whether to delete the section or replace it with something accurate

## 3. Verify

- [x] 3.1 `grep -rn "shelley\|Shelley" --include="*.md" .` (excluding `node_modules/`, `.git/`, archived openspec changes, and `.pi/` which is unrelated) returns nothing
- [x] 3.2 `grep -rln "\./run\b\|/app/repo/run\b\|0001-ui-abandon-visuals" .` (excluding archived openspec changes) returns nothing
- [x] 3.3 `make build && make test` pass (sanity check — none of this should affect the build, but confirm)
