---
id: am0xw6d0t4k06kjk7va0efcr
title: ensigncycle/streamwatch.go is test-supporting lib shipping in the production binary — rename to *_test.go to keep it out of release builds
status: validation
source: captain (2026-06-02) — observed 513 LOC in internal/ensigncycle/streamwatch.go that is test-only (zero non-test importers); cluster test-supporting lib so it does not ship in the binary
score: "0.20"
worktree: .worktrees/spacedock-ensign-ensigncycle-streamwatch-test-only
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

**Decision: single-file rename** (not split). The 513 LOC are one cohesive unit — the `streamWatcher` and its predicates/helpers are used together by the live tests; splitting buys nothing for the binary (both halves would be `*_test.go` either way) and only adds churn. Single `streamwatch_test.go` it is.

Rename `internal/ensigncycle/streamwatch.go` → `internal/ensigncycle/streamwatch_test.go`:
- The package declaration stays `package ensigncycle`.
- All symbols stay unexported and same-package-accessible from the existing `*_test.go` files.
- The Go test toolchain compiles `*_test.go` only for the test binary, so the production binary drops the 513 LOC at the next build.

**One additional edit is required — the rename is NOT pure.** The ideation spike (see Spike below) found a hidden dependency the original body missed: `internal/ensigncycle/live_budget_test.go:28` hardcodes the source filename in its AC-1 timeout-budget guard:

```go
var liveBudgetSources = []string{"streamwatch.go", "live_test.go"}
```

`TestNoTimeoutLiteralExceeds60s` opens that file by literal name to AST-parse it. After the rename, this fails with `parse streamwatch.go: open streamwatch.go: no such file or directory`. The implementation must update this literal to `"streamwatch_test.go"` in the same commit. This is the only non-rename change.

## Out of scope

- Moving any other production file into `*_test.go` form: the scan found no other candidates with zero non-test importers.
- Introducing a separate `internal/testlib/` package: that would still ship in the binary (Go has no test-only package concept), and the cross-package access would force exporting symbols that are currently unexported.
- Rewriting any of the watcher logic. This is a filename move only.

## Acceptance criteria

**AC-1 — The release binary stops shipping streamwatch.go's 513 LOC.**
Verified by: `find internal cmd -name '*.go' -not -name '*_test.go' | xargs wc -l | tail -1` shows the non-test total drops by ~513 LOC after the rename (modulo any other in-flight changes); `go test ./internal/ensigncycle/...` stays green; `go build -o spacedock ./cmd/spacedock` succeeds.

**AC-2 — No production code path silently depended on streamwatch symbols.**
Verified by: `go build ./...` succeeds after the rename (the compiler proves no non-test code referenced the symbols); the `internal/ensigncycle` suite stays green with no source edits beyond the rename and the single `live_budget_test.go:28` filename-literal update.

**AC-3 — The AC-1 timeout-budget guard tracks the renamed file.**
Verified by: `liveBudgetSources` in `internal/ensigncycle/live_budget_test.go` lists `"streamwatch_test.go"` (not `"streamwatch.go"`); `go test ./internal/ensigncycle/ -run TestNoTimeoutLiteralExceeds60s` passes (it would `t.Fatalf` on the missing old filename if the literal were stale).

## Test plan

- Single commit: the `git mv` rename plus the one-line `live_budget_test.go:28` literal update (`"streamwatch.go"` → `"streamwatch_test.go"`).
- Run `go build ./... && go build -o spacedock ./cmd/spacedock && go test ./internal/ensigncycle/` — all three pass. The package suite includes `TestNoTimeoutLiteralExceeds60s`, which is the regression net that the guard literal is correct.
- Cost: trivial. No new tests; the existing suite is the regression net.
- No live-workflow test needed — the claim is "the binary does not include this file" which is a build-time observation; the guard fix is covered by an existing offline test.

## Spike (ideation, 2026-06-03)

The riskiest unverified assumption was "the rename is pure — no other file depends on the literal filename `streamwatch.go`." Exercised it end-to-end on HEAD (`6526943e`): `git mv streamwatch.go streamwatch_test.go`, then `go build ./...` + `go build -o spacedock ./cmd/spacedock` (both clean — compiler confirms no non-test code referenced the symbols), then `go test ./internal/ensigncycle/`.

Result: the test run FAILED on `TestNoTimeoutLiteralExceeds60s` — `live_budget_test.go:28` hardcodes `"streamwatch.go"` and opens it by name. This is the hidden dependency the original body's "no other code changes needed" claim missed. Patching the literal to `"streamwatch_test.go"` turned the suite green (`ok internal/ensigncycle`). Spike reverted; working tree clean. The Go test-file exclusion mechanism itself is proven toolchain behavior (no novelty), but the filename-literal coupling was the real unknown and is now on the record. This finding seeds the implementation's first check: the rename commit must touch both files.

## Notes

- The 513 LOC concentration came up in the captain's "why is the codebase 10k LOC" inquiry on 2026-06-02. The honest answer was that streamwatch.go is recently-added live-e2e infrastructure (#267 `3g` + #271 `hd`); this entity converts that observation into a small permanent win — the file stays, it just stops shipping.
- Pairs with the sprint-notes "parsing modernization" follow-up only in spirit (both are post-bootstrap LOC hygiene); no dependency, can land independently.
- Sequence: low score, can wait behind `zj` / `n3` / `qs` / the install-ref-refresh entity; pick up in any quiet slot.

## Stage Report: ideation

- DONE: Confirm on current HEAD that internal/ensigncycle/streamwatch.go still has zero non-test importers (the rename precondition) and decide single-file vs split into multiple *_test.go files.
  On HEAD 6526943e: 513 LOC, all top-level symbols unexported (grep for `^(func|type|var|const) [A-Z]` empty), zero `.go` files import the package, every symbol consumer is a `*_test.go` (live_test.go, streamwatch_unit_test.go, live_budget_test.go). Decided single-file; rationale in Proposed approach.

### Summary
Confirmed the rename precondition empirically on HEAD and chose a single-file rename (split buys nothing). The throwaway spike caught a hidden dependency the original body missed — `live_budget_test.go:28` hardcodes the literal filename `"streamwatch.go"` and AST-parses it, so the rename is NOT pure: `TestNoTimeoutLiteralExceeds60s` fails until that literal becomes `"streamwatch_test.go"`. Updated the approach, added AC-3 for the guard, refined the test plan to a single two-file commit, and recorded the spike result. Spike reverted; working tree clean.

## Stage Report: implementation

- DONE: Rename internal/ensigncycle/streamwatch.go -> streamwatch_test.go and update the hardcoded filename literal that AST-parses it.
  `git mv` rename (pure, 0 content change) + `live_budget_test.go:28` literal `"streamwatch.go"` -> `"streamwatch_test.go"`; committed together as 91156966 on branch spacedock-ensign/ensigncycle-streamwatch-test-only. TDD: ran TestNoTimeoutLiteralExceeds60s RED first (`parse streamwatch.go: ...no such file`), then green. `go build ./...` Success, `go test ./internal/ensigncycle/...` 92 passed — no source edits beyond rename + that one literal.
- DONE: Produce AC evidence (AC-1 LOC drop + binary build; AC-3 guard test passes with renamed file).
  AC-1: non-test LOC 11924 -> 11266 (drop 658; the file grew to 658 LOC off origin/next vs the 513 measured at the older HEAD 6526943e the body cited — the "modulo in-flight changes" caveat); `go build -o spacedock ./cmd/spacedock` Success. AC-2: `go build ./...` Success proves no non-test caller. AC-3: `go test ./internal/ensigncycle/ -run TestNoTimeoutLiteralExceeds60s` 1 passed; `liveBudgetSources` now lists `"streamwatch_test.go"`.

### Summary
Single-commit rename of streamwatch.go to streamwatch_test.go plus the one required literal update in live_budget_test.go, so the 658 LOC of test-only watcher infrastructure no longer compiles into the release binary. Followed TDD on the guard test (RED on stale literal, green after update). Non-test LOC dropped 11924 -> 11266; all three builds (`./...`, binary, package suite of 92 tests) green. The only deviation from the body is the LOC figure: the file is now 658 LOC (was 513 at the older cited HEAD), which fully accounts for the larger-than-expected drop; no other files were moved. Note: the package already contained pre-existing streamwatch_unit_test.go and streamwatch_regression_test.go — these were untouched.

## Stage Report: validation

- DONE: Independently reproduce the evidence for each AC: AC-1 (non-test LOC drops materially vs origin/next AND `go build -o spacedock ./cmd/spacedock` succeeds), AC-2 (`go build ./...` succeeds), AC-3 (`go test ./internal/ensigncycle/ -run TestNoTimeoutLiteralExceeds60s` passes and liveBudgetSources lists "streamwatch_test.go").
  AC-1: non-test LOC 11924 (origin/next) -> 11266 (HEAD) = drop 658, exactly the renamed file's size; binary build OK. AC-2: `go build ./...` Success. AC-3: literal is `"streamwatch_test.go"` (live_budget_test.go:28); targeted test 1 passed.
- DONE: Confirm the change is MINIMAL and correct: `git diff origin/next...HEAD` shows ONLY the rename + the single live_budget_test.go literal update; full offline suite green. Recommend PASSED or REJECTED.
  Diff is exactly the rename (content byte-identical: `git diff origin/next:streamwatch.go HEAD:streamwatch_test.go` empty, exit 0) + one-line literal change. `go test ./...` 978 passed, exit 0. No detached audit needed (low blast radius: pure rename + one test literal).

### Summary
Independently reproduced all AC evidence on a fresh run in the worktree (HEAD 91156966 off origin/next). The rename is verifiably pure — the renamed file's content is byte-identical to origin/next's streamwatch.go (empty diff). The only content edit is the required live_budget_test.go:28 literal update, which the targeted guard test confirms tracks the renamed file. Full offline suite green (978 passed). Recommendation: PASSED.
