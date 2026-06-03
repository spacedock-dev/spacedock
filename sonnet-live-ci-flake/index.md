---
id: yyqez6npx8qb11b5v7fgwjtf
title: Sonnet live-CI fails with a repeatable shape — FO subprocess ends its turn after shutdown_request before the streamwatcher's expected dispatch-close fires
status: ideation
source: session-10 — identical failure shape on PR #275 (n3 frontmatter-hash-quoting) and PR #277 (2a require-external-proof-guard); both offline-green + opus-green + sonnet-red. Two matching failures = mechanism issue, not random flake. Blocks merging n3 + 2a into 0.19.5.
score: "0.19"
worktree:
started: 2026-06-03T07:09:59Z
completed:
verdict:
issue:
---

The live-e2e CI cycle passes on offline and on opus but **fails on sonnet with a repeatable shape** across two unrelated PRs (n3 #275, 2a #277). The live FO subprocess ends its turn after sending `shutdown_request` to its ensign — waiting for the ensign's shutdown approval — but the streamwatcher's expected dispatch-close assertion does not fire before the FO turn ends. The FO subprocess self-reports success (`terminal_reason: completed`, `is_error: false`) while the test's bounded stop condition was never reached. Two identical failures on different PRs means this is a streamwatcher-vs-sonnet-FO-interpretation issue, not a flake.

## Problem

The live cycle test asserts a bounded sequence of observable stream events (dispatch open → … → dispatch close) within a budget. On sonnet:
- The FO subprocess reaches the point of tearing down the ensign and emits `shutdown_request`.
- Per the FO contract's `## Awaiting Completion` rules, it then ends its turn empty, waiting for the ensign's cooperative shutdown to complete.
- The streamwatcher's `expectDispatchClose` (or equivalent bounded-stop assertion) does not observe the close event before the FO subprocess's turn ends, so the test fails its bound — even though the FO subprocess itself terminated with `is_error: false`.

Two non-exclusive hypotheses (ideation picks/refutes):
1. **Streamwatcher budget mismatch on sonnet's verbosity.** Sonnet's interpretation is more verbose, so the dispatch-close event arrives outside the event/time budget the streamwatcher allows. The budget was tuned against opus's terser stream.
2. **Contract-side ambiguity about turn-ending after `shutdown_request`.** The FO ends its turn before the close is observable, and the contract does not guarantee the close emits within the same observable window. The fix may be a contract clause (emit/await the close before ending the turn) rather than a budget bump.

## Proposed approach

1. **Reproduce from the archived evidence first.** Session 10's live-e2e CI archives the agent session JSONL on failure (shipped this sprint). Pull the JSONL from the failed sonnet runs on #275 and #277 and trace the exact event ordering around `shutdown_request` / dispatch-close — confirm which hypothesis holds before touching code. This is the riskiest unknown; exercise it first.
2. Based on the trace, either:
   - widen / correct the streamwatcher's bounded-stop budget or its close-event matcher (`internal/ensigncycle/streamwatch.go` and its tests), or
   - tighten the FO contract's turn-ending-after-shutdown clause so the close is observable within the asserted window.
3. Add a regression test reproducing the sonnet ordering so the fix is locked in.
4. Re-trigger live CI on #275 and #277; both must pass sonnet + opus + offline.

## Out of scope

- Merging n3 / 2a themselves — that is the follow-on once this fix lands (tracked on those PRs).
- The `am` streamwatch rename (test-only file → `*_test.go`) — orthogonal; do not couple.
- Broader live-CI redesign; this is a targeted fix for the observed ordering bug.

## Acceptance criteria

Proof must run the behavior (a CI run is the only true oracle for this runtime-observable claim — `looks-right-here ≠ runs-right-everywhere`). Ideation refines.

**AC-1 — The sonnet live cycle reaches its bounded stop condition on the failing surface.**
Verified by: a live-e2e CI run on sonnet (on n3 #275 and 2a #277, or an equivalent reproduction) passing the dispatch-close bound that currently fails — the run is the oracle.

**AC-2 — A regression test reproduces the sonnet event ordering and fails before the fix, passes after.**
Verified by: a Go test in `internal/ensigncycle` (or the streamwatcher's test suite) that encodes the archived sonnet ordering around `shutdown_request` / dispatch-close; red on the pre-fix code, green after.

**AC-3 — Opus and offline cycles stay green (no regression).**
Verified by: the opus live run + `go test ./...` still passing after the fix.

## Test plan

- Diagnostic: read the archived JSONL from the failed sonnet runs (#275, #277) — cheap, must precede any code change.
- Regression: a Go unit/fixture test over the streamwatcher's bounded-stop logic against the captured ordering.
- Confirmation: live-e2e CI on sonnet + opus. Cost: medium-high (live cycles), but the JSONL replay should make most of the loop offline-reproducible.
- High-stakes surface (CI/release machinery) → detached adversarial audit required before merge per the dev README's validation stage.

## Notes

- Failure evidence: PR #275 and PR #277 sonnet live-CI runs (session 10). Both: offline ✅, opus ✅, sonnet ❌. The captain cancelled the n3 retrigger at session-10 close; the failures are reproducible, not transient.
- This is the gating fix that carries n3 + 2a back into mergeable shape for 0.19.5.
- Related context: `internal/ensigncycle/streamwatch.go` is the Go port of the upstream Python `FOStreamWatcher`; the `am` entity proposes renaming it to `*_test.go` but does not change its logic — keep the two entities decoupled.
