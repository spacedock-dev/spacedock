---
id: kq8x081hr1x9cb9jptw8wz0c
title: 0202 pre-cut non-blockers — gofmt drift (2 files) + dedup stateHasOrigin
status: backlog
source: "0202 pre-cut antipattern audit (2026-06-13). Two NON-BLOCKER findings recorded; neither gated the v0.20.2 cut. Pre-existing, not 0202-introduced."
group: cleanup
---

Two non-blocking cleanups the 0202 pre-cut audit recorded. Mechanical.

## Items

1. **gofmt drift** on `internal/release/channel_agreement_guard_test.go` (from #357/#352) and a single trailing-blank-line drift on `internal/contract/contract_test.go` (#352). `gofmt -w` to clear. Neither was introduced by 0202 (the sprint cleaned `survey_sync_codex_test.go` instead).
2. **Duplicated `stateHasOrigin`** across `internal/status` (via `runGitCmd`) and `internal/dispatch/build.go` (local `exec.Command`). Both ask `git remote get-url origin`. Defensible today (package boundary; status's `runGitCmd` is unexported), but a candidate for a shared helper.

## Out of scope
Anything requiring a behavior change. Each item is a localized format/dedup fix.

## Acceptance criteria (sketch)

**AC-1 (sketch) — `gofmt -l` is clean on the two drifting files.**
Verified by: `gofmt -l internal/release/channel_agreement_guard_test.go internal/contract/contract_test.go` prints nothing.

**AC-2 (sketch) — `stateHasOrigin` has a single definition (if the dedup is taken).**
Verified by: a grep/test showing one shared implementation, or a recorded decision that the package boundary justifies keeping two.

## Notes
Recorded per the roadmap Close step (seed the next sprint with the pre-cut audit's deferred findings).
