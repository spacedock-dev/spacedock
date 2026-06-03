---
id: am0xw6d0t4k06kjk7va0efcr
title: ensigncycle/streamwatch.go is test-supporting lib shipping in the production binary — rename to *_test.go to keep it out of release builds
status: ideation
source: captain (2026-06-02) — observed 513 LOC in internal/ensigncycle/streamwatch.go that is test-only (zero non-test importers); cluster test-supporting lib so it does not ship in the binary
score: "0.20"
worktree:
started: 2026-06-03T07:09:59Z
completed:
verdict:
issue:
---

`internal/ensigncycle/streamwatch.go` is the Go port of the upstream Python `FOStreamWatcher` (scripts/test_lib.py) and is used exclusively by the live-cycle tests in the same package. Every symbol is unexported. A whole-codebase scan confirms it is **the only non-test file with zero non-test importers** — every other production file is reached by some non-test caller.

It currently ships in the release binary at 513 LOC (~5% of the 9.7k-LOC non-test footprint). Renaming it to `*_test.go` keeps the test surface intact and drops it from the production build.

## Problem

Go has no test-only package concept — only test-only files. A non-`_test.go` filename means the file is compiled into every consumer of the package, including the production binary. `streamwatch.go` has been a non-test file by historical accident (it was a Go port of a Python test helper; the conversion did not tighten the filename), so 513 LOC of test-supporting infrastructure ships in every release.

The scan (the captain's question, answered):
- Symbol audit: every `func` / `type` / `var` in `streamwatch.go` is lowercase (unexported), so no cross-package production code can reference them.
- Importer audit: `grep -rln '"github.com/spacedock-dev/spacedock/internal/ensigncycle"'` matches only `*_test.go` files in the same package and no other production files. Same-package access does not require an import line, so the test files reach the symbols directly without an import statement.
- Whole-codebase pass: `streamwatch.go` is the only `internal/**/*.go` non-test file with zero non-test importers.

## Proposed approach

Rename `internal/ensigncycle/streamwatch.go` → `internal/ensigncycle/streamwatch_test.go`. No other code changes needed:
- The package declaration stays `package ensigncycle`.
- All symbols stay unexported and same-package-accessible from the existing `*_test.go` files.
- The Go test toolchain compiles `*_test.go` only for the test binary, so the production binary drops the 513 LOC at the next build.

If the captain prefers to split the watcher itself from its predicates/helpers, that can land as two `*_test.go` files (e.g., `streamwatch_test.go` + `streamwatch_helpers_test.go`) — equivalent outcome for the binary, finer-grained file organization. Ideation picks one.

## Out of scope

- Moving any other production file into `*_test.go` form: the scan found no other candidates with zero non-test importers.
- Introducing a separate `internal/testlib/` package: that would still ship in the binary (Go has no test-only package concept), and the cross-package access would force exporting symbols that are currently unexported.
- Rewriting any of the watcher logic. This is a filename move only.

## Acceptance criteria

**AC-1 — The release binary stops shipping streamwatch.go's 513 LOC.**
Verified by: `find internal cmd -name '*.go' -not -name '*_test.go' | xargs wc -l | tail -1` shows the non-test total drops by ~513 LOC after the rename (modulo any other in-flight changes); `go test ./internal/ensigncycle/...` stays green; `go build -o spacedock ./cmd/spacedock` succeeds.

**AC-2 — No production code path silently depended on streamwatch symbols.**
Verified by: `go build ./...` succeeds after the rename (the compiler proves no non-test code referenced the symbols); the existing `cycle_test.go` / `live_test.go` / `streamwatch_unit_test.go` suite stays green without source edits beyond the rename.

## Test plan

- Single rename commit. Run `go build ./... && go test ./...` and confirm both pass without other source changes.
- Cost: trivial. No new tests; the existing suite is the regression net.
- No live-workflow test needed — the claim is "the binary does not include this file" which is a build-time observation.

## Notes

- The 513 LOC concentration came up in the captain's "why is the codebase 10k LOC" inquiry on 2026-06-02. The honest answer was that streamwatch.go is recently-added live-e2e infrastructure (#267 `3g` + #271 `hd`); this entity converts that observation into a small permanent win — the file stays, it just stops shipping.
- Pairs with the sprint-notes "parsing modernization" follow-up only in spirit (both are post-bootstrap LOC hygiene); no dependency, can land independently.
- Sequence: low score, can wait behind `zj` / `n3` / `qs` / the install-ref-refresh entity; pick up in any quiet slot.
