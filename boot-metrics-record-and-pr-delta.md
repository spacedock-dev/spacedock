---
id: e7fx68wkhhrs8tsvpv0j8s6x
title: Boot-window metrics — record turns/cache-creation instead of gating on magic thresholds, surface per-PR delta
status: backlog
source: captain conversation 2026-07-02, journey-costs ledger review (see docs/roadmap/0203-fo-efficiency/index.md task j9 / AC-6 for the constants' origin)
started:
completed:
verdict:
score:
worktree:
issue:
---

AC-6's shallow-boot measured-saving oracle (`internal/ensigncycle/shallow_boot_measure_test.go`, `assertShallowBootMeasured`) hard-fails CI on two absolute constants — `greetContextCeiling = 60000` and `teamRecacheSpikeFloor = 60000` — calibrated once against the pre-optimization ~160k/~89k baseline in task j9 (`0203-fo-efficiency`, shipped v0.20.3) and never revisited since. They no longer mean anything as a pass/fail gate. Replace them with recorded telemetry (turns-to-greet, greet-turn `cache_creation` tokens) that rides the existing `journeymetrics` ledger pipe, and close the per-PR visibility gap identified in the same conversation: the published `journey-costs-vX.json` ledger is release-cadence only and currently sourced from whichever manual `workflow_dispatch` run on `main` happened to run last (see `release.yml`'s `journey-ledger` job, `gh run list --branch main --limit 1`) — it never reflects the PR runs that already execute continuously (`runtime-live-e2e.yml` triggers on every PR into main once the live jobs' environment is approved) and already upload the exact raw stream needed (`live-artifacts/claude/<model>/claude-shared-scenarios/shallow-boot/claude-stream.jsonl`, ~90-day artifact retention).

## Problem

- The two AC-6 thresholds are magic numbers frozen against a baseline that predates most of the current contract; a red/green verdict against them carries no information about whether today's boot got better or worse.
- No historical trend exists for boot cost (turns, cache-creation tokens) at all — the release ledger is a single snapshot per cut, not a series.
- Per-PR live runs already produce the raw data needed but nothing reads or surfaces it; a regression introduced in a PR is invisible until (if ever) someone manually dispatches a main run before the next release.

## Proposed approach

Captain-agreed direction from conversation (ideation should firm up the specifics, not relitigate the shape):

1. **Record, don't gate.** In `runClaudeShallowBootScenario` (`claude_live_runner_test.go`), where `assertShallowBootMeasured` already computes the turn list and the greet-turn index, emit a `journeymetrics.Record` (via `journeymetrics.EmitRecord`, into the same `SPACEDOCK_JOURNEY_METRICS_DIR/shared-scenarios` dir the whole-run `shallow-boot` scenario record already uses) carrying `Turns: greet+1` and the greet turn's `TokenTotals` (specifically `CacheCreation`) as a distinct boot-window observation. Drop (or drastically loosen to a sanity-only tripwire, not the current tight absolute number) the two `t.Fatalf`s. Structural checks that are NOT threshold-based (`assertNoTeamCreateBeforeGreet`, `assertGreetInvokesNoDeferredFOSkill`) are unaffected and stay as hard gates.
2. **Widen release-ledger aggregation.** Change the `journey-ledger` job's discovery query in `release.yml` from "single latest `branch:main` run" to pulling journey-metrics from every successful Runtime Live E2E run (PR and main both) since the previous release tag, so the published ledger becomes a real aggregate/trend instead of one snapshot. Reuses artifacts that already exist; no new upload plumbing needed.
3. **Per-PR delta comment.** Add a step (triggered from a PR's own live run, once it completes) that reads the previously *published* release's `journey-costs-vX.json` (the last `gh release` asset, e.g. via the same `gh api` release-asset lookup used in this conversation), computes a delta per scenario/model against the current PR run's freshly emitted records, and posts it as a single updating PR comment (not a new comment per push). Scope this to whatever the ledger already tracks generically (turns, tokens, cost) rather than one-off special-casing the boot-window metric — the diff mechanism doesn't care which scenario it's comparing.

## Out of scope

- Retroactively backfilling the ~90-day window of already-existing PR artifacts into historical trend data — worth calling out as a possible fast-follow, not required here.
- A live/real-time dashboard beyond the PR comment.
- Expanding tracked scenarios beyond the current six.

## Acceptance criteria

Ideation must firm these up with concrete "Verified by" evidence per the README's AC rules (end-value measurement, not just "the code changed"); draft shape:

**AC-1 - The shallow-boot live scenario no longer fails CI on the absolute 60k/60k constants; turns-to-greet and greet-turn cache-creation are recorded as a journeymetrics observation instead.**
Verified by: a live (or offline fixture, per the existing `shallow_boot_measure_unit_test.go` pattern) run producing an emitted record with real Turns/Tokens.CacheCreation values, and the old `t.Fatalf` ceiling/spike paths no longer reachable as hard failures.

**AC-2 - The published release ledger reflects PR-run history, not a single manual main dispatch.**
Verified by: `journey-costs-vX.json` for a release built from a fixture/mock run set containing multiple PR-origin runs shows multiple observations per scenario spanning those runs, not just one.

**AC-3 - A PR gets a visible delta of its own live-run numbers against the previously published release ledger.**
Verified by: a fixture PR run + a stub previous-release ledger produces a single updating PR comment with a correct per-scenario delta (not a new comment appended each run).

## Test plan

{Ideation to fill in: unit tests for the new emission logic against offline fixture streams (mirroring `shallow_boot_measure_unit_test.go`), a fixture/mock test for the widened ledger-aggregation query, and a fixture test for the delta-computation + sticky-comment formatting. Live-workflow cost/complexity to be estimated in ideation.}
