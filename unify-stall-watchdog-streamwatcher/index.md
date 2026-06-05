---
id: n2xgf274sb9zk2jgg2r5w0xg
title: Unify the two live-test stall mechanisms + reconcile the 60s/120s budgets fleet-wide
status: backlog
source: "first-officer (2026-06-05) — surfaced during gq (feedback-nonhappy-live-coverage) timeout re-fix. The shared-runner stall-watchdog added there is a SECOND implementation of a mechanism the codebase already has (the mature streamWatcher), and the two now run at different budgets (120s vs 60s) with different guards. Flagged out-of-scope by the gq worker; filed to track, not churn gq."
score: "0.20"
started:
completed:
verdict: rejected
worktree:
issue:
---

`internal/ensigncycle` now carries TWO parallel no-progress stall mechanisms with the same "reset the deadline on every drained line" design:

1. **`streamWatcher`** — the mature Go port of the upstream Python FOStreamWatcher; `quietBudgetDefault = 60s`; wired into `TestLiveEnsignCycle` (`live_test.go`) via `newPipeLineSource`; guarded by `live_budget_test.go` (`TestNoTimeoutLiteralExceeds60s` + `TestBudgetConstantsAreUnder60s`, the standing AC-1 "no timeout literal >60s" guard).
2. **`streamWithStallWatchdog`** — added in gq for the shared-scenario runners (claude/codex `live_runner`), which previously had NO per-step discipline (just a monolithic basket, now removed); `stageStallTimeout = 120s`, pinned by `TestStageStallTimeoutIsCaptainApprovedException` as an audited captain-approved exception.

This is a DRY violation (two impls of one mechanism) and a fleet-wide inconsistency (60s vs 120s; two different guards).

## Why 120s is the right shared-runner budget (measured, do not regress)

The gq FO measurement established the 120s exception is evidence-justified, not arbitrary: opus rejection-flow runs ~9m with a max FO-stream-silence gap of **59.1s** (a sub-agent dispatch blocks the FO top-level stream; opus emits NO `task_progress` heartbeats during the dispatch, unlike sonnet's 32), so a 60s budget is razor-thin/flaky on opus. Sonnet max gap 28.3s. So the unified budget should be 120s (~2x margin), as an audited exception — NOT 60s.

## Open question this task must settle

Does `TestLiveEnsignCycle`'s `streamWatcher` (still at 60s) have the same opus exposure? It is a different scenario (single-cycle full-cycle smoke) and has run at 60s without a reported flake, so the 60s may be fine for ITS dispatch pattern — or it may be latent. Measure it (the gq FO's stream-gap parse is the reusable tool) before bumping; do not bump on speculation.

## Out of scope

The gq deliverable itself (the reuse fix + the shared-runner watchdog at 120s) — it ships as-is with the audited exception. This task is purely the consolidation.

## Acceptance criteria

**AC-1 — one stall mechanism.** The shared-scenario runners and `TestLiveEnsignCycle` share a single no-progress watcher implementation (keep the mature `streamWatcher`; remove the `streamWithStallWatchdog` duplicate, or vice-versa with justification).
Verified by: the duplicate type is gone (grep/AST) and both call sites route through the surviving one; offline `go test ./...` green.

**AC-2 — one budget + one guard, with the exception documented once.** The per-stage budget is consistent fleet-wide, and the standing AC-1 guard (`live_budget_test.go`) either keeps ≤60s for the paths that measure safe OR encodes the captain-approved 120s exception explicitly for the paths that need it — no path silently outside the guard.
Verified by: the AC-1 guard scans all live-runner stall budgets; a planted >budget literal in any covered file reds it; `TestLiveEnsignCycle`'s budget decision is backed by a measured stream-gap (per the open question), not speculation.

## Test plan

Offline Go refactor + the standing guard extension; mutation-confirm the guard reds on a drifted/over-budget literal in every covered file. If `TestLiveEnsignCycle`'s budget is bumped, a one-off opus stream-gap measurement justifies it. No new live model spend beyond at most one measurement run.
