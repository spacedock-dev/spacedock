---
title: "e2e-gate --status success filter hides re-run-to-green Runtime Live E2E runs from the release gate"
status: ideation
score: 0.7
source: "v0.23.0 stable cut, 2026-06-30. The release e2e-gate blocked the cut: gh run list --workflow 'Runtime Live E2E' --status success -c <commit> returned [] for a run that was re-run to green (run_attempt 2), because GitHub's --status success filter only matches first-attempt successes. The genuine green matrix (run 28429490220, all 5 lanes success on the tagged SHA e3f85ec3) was invisible to the gate. Worked around the cut by dispatching a fresh run_attempt-1 green run; this fix removes the trap."
id: q3v95nmwkwh656p5m6rhher1
started: 2026-07-02T01:23:59Z
---

The release e2e-gate's `gh run list` query uses `--status success`. At the v0.23.0 cut this filter returned `[]` for a run that had been re-run to green minutes earlier (run_attempt 2, overall `conclusion: success` on the tagged SHA), while the same query WITHOUT `--status` returned that run with conclusion=success — so the gate blocked a genuinely green matrix and goreleaser was skipped. This directly contradicts the proof-policy / handoff guidance to "re-run the failed lane to green, never block the tag". A 2026-07-02 re-exercise (see Spike determination) shows the filter later catches up, so the trap is a consistency window right after a re-run flips a run to green — exactly the window in which a blocked cut gets retried.

## Problem
`cmd/spacedock-release/e2e_gate.go` (`ghRunListForCommit`) builds `gh run list --workflow "Runtime Live E2E" --status success -c <commit> --json databaseId,headSha,conclusion,status`. The `--status success` pre-filter is both redundant AND harmful: `internal/release/e2egate.go` (`EvaluateE2EGate`) ALREADY checks `run.Conclusion == "success" && run.HeadSha == releaseCommit`, and the filter can lag `conclusion` after a re-run (v0.23.0 gate-time evidence: `--status success -c e3f85ec3` -> []; `-c e3f85ec3` (no --status) -> the run with conclusion=success, run_attempt=2). The function's own doc comment justifies the flag as the parked-run defense, but that defense actually lives in the predicate (empty conclusion never matches), so the flag buys nothing and costs green re-runs.

## Proposed approach
Drop `--status success` from `ghRunListForCommit`. The `-c <commit>` filter scopes to the tagged commit and `EvaluateE2EGate` filters for conclusion==success, so the predicate stays correct and now sees re-run-to-green runs. Alongside the flag removal:
- Rewrite `ghRunListForCommit`'s doc comment: the parked-run exclusion belongs to `EvaluateE2EGate`, and the query is deliberately unfiltered by status because the status filter lags `conclusion` after re-runs.
- Update the query description in the comment at `internal/release/e2egate_workflow_test.go:16` (it spells out the old query shape).
- Apply the doc diff below to `docs/releasing.md`.

Edge case considered: without `--status success` the run list for the commit also carries failed/cancelled/parked entries, all counting against `gh run list`'s default `--limit 20`. Runtime Live E2E is workflow_dispatch-only and re-runs share one run entry, so runs-per-commit stays in single digits; no limit change needed. `EvaluateE2EGate` already scans for ANY matching success among mixed entries (`TestE2EGatePicksMatchingRunAmongMany`).

## Spike determination
The riskiest mechanism — `gh run list` WITHOUT `--status` returning the re-run-to-green run with conclusion=success — is proven; no further spike needed:
- **2026-06-30 gate time** (recorded in the task source): `--status success -c e3f85ec3` -> `[]`; `-c e3f85ec3` -> the run with conclusion=success, run_attempt=2. The re-run went green at 08:44:25Z and the fresh workaround run was only created at 08:54:33Z, so the filtered query was still empty at least ~10 minutes after green while the unfiltered query already saw the success.
- **2026-07-02 re-exercise** (this ideation): BOTH queries now return run 28429490220 (`gh run view` confirms attempt=2, conclusion=success, headSha=e3f85ec3b1...). So "`--status success` only matches first-attempt successes" is NOT a permanent GitHub rule — the observed failure is a filter-consistency lag after re-run-to-green. The fix does not depend on which explanation is right: at gate time `conclusion` was correct and the unfiltered query saw it (proven on both dates), so dropping the redundant pre-filter removes the trap under either.

## Acceptance criteria
- **AC-1 (VALUE)** — a Runtime Live E2E run that was re-run to green (overall conclusion success on the tagged SHA, run_attempt >= 2) PASSES the e2e-gate. Verified by test T1: a fixture-driven test where the run-list JSON's only success entry is a re-run-to-green run, asserting the gate accepts it; RED on the current `--status success` query shape, GREEN after dropping it.
- **AC-2** — a parked/waiting run (empty conclusion) still does NOT satisfy the gate (no regression of the original spike intent: parked offline-only runs must not qualify). Verified by test T2 end-to-end under the new query shape, plus the existing predicate tests staying green.

## Test plan
- **T1 (AC-1, the RED/GREEN test)** — in `cmd/spacedock-release`, a test that exercises the REAL `ghRunListForCommit` (not the `fakeRunLister` seam, which bypasses the args and cannot see the bug) against a fake `gh` shim on PATH (`t.Setenv("PATH", dir)` — the established pattern of `internal/status/live_prstate_pin_test.go`). The shim pins the 2026-06-30 gate-time observation (the adversarial case, not GitHub's caught-up steady state): its fixture's only success entry is the re-run-to-green run (databaseId 28429490220, conclusion success, headSha = release commit; attempt 2 recorded in the shim fixture); when its argv contains `--status`, it prints `[]`; otherwise it prints the fixture projected to the requested `--json` fields. The test runs `runE2EGate([]string{commit}, ghRunListForCommit)` and asserts exit 0 and the step summary citing the matched run. RED today (args carry `--status success` -> shim returns `[]` -> exit 1), GREEN after dropping the flag.
- **T2 (AC-2)** — same shim harness, fixture containing ONLY a parked run (conclusion "", status waiting) for the release commit; assert `runE2EGate` exits 1. This proves the parked-run protection survives removal of the pre-filter end-to-end; the predicate-level guarantees stay covered by the existing `TestE2EGateBlocksParkedRun` / `TestE2EGateCommandBlocksOnParkedRun`.
- **T3 (regression)** — `go test ./internal/release/... ./cmd/spacedock-release/...` stays green, including `e2egate_workflow_test.go` (release.yml still gates goreleaser on the e2e-gate job).

Cost: small — two Go tests plus a ~10-line shim helper in `cmd/spacedock-release`; fixture/CLI level only. No live workflow test needed: the live mechanism is proven twice (Spike determination above), and the shim encodes exactly those recorded observations.

## Documentation diff
`docs/releasing.md`, step "4. Green that exact commit" — append one sentence after "nothing greens the commit unless you dispatch it:":

Before:
> Capture the release SHA, then dispatch Runtime Live E2E on it and wait for a `conclusion: success` run — the `e2e-gate` matches a green run to the tagged SHA, and the workflow is `workflow_dispatch`-only, so nothing greens the commit unless you dispatch it:

After:
> Capture the release SHA, then dispatch Runtime Live E2E on it and wait for a `conclusion: success` run — the `e2e-gate` matches a green run to the tagged SHA, and the workflow is `workflow_dispatch`-only, so nothing greens the commit unless you dispatch it. A lane that flakes can be re-run to green (`gh run rerun <run-id> --failed`); the re-run-to-green run satisfies the gate — no fresh dispatch needed:

## Out of scope
`.github/workflows/release.yml:56` (journey-ledger artifact download) uses the same `gh run list --status success` shape but different semantics: `--limit 1 --jq '.[0]'` picks the LATEST success, and naively dropping the flag there would pick the latest run regardless of conclusion. It is best-effort (a miss warns and skips, never blocks a cut). Same trap, separate fix shape — worth its own seed, not this task.

## Stage Report: ideation

- DONE: AC-1's fixture-driven test plan is concrete and falsifiable: the run-list JSON fixture's only success entry is a re-run-to-green run (run_attempt >= 2), asserted RED under the current --status success query shape and GREEN after dropping it.
  Test plan T1: fake `gh` shim on PATH drives the REAL `ghRunListForCommit` (the existing `fakeRunLister` seam bypasses the args and cannot see the bug); shim returns `[]` when argv carries `--status`, the attempt-2 success run otherwise.
- DONE: AC-2 (a parked/waiting run with empty conclusion still does NOT satisfy the gate) has its own explicit test in the plan — the original spike intent cannot regress silently.
  Test plan T2: same shim harness, parked-only fixture, assert exit 1 end-to-end under the new query shape, alongside the existing predicate tests.
- DONE: The riskiest mechanism — gh run list without --status returning the re-run-to-green run — is recorded as proven (the v0.23.0 live evidence in the task source) or re-exercised, not asserted; the spike/no-spike determination is on the record.
  Re-exercised live 2026-07-02: both queries now return run 28429490220 (gh run view: attempt=2, conclusion=success); Spike determination section records both observations and the no-further-spike call.

### Summary

Fleshed out problem, approach, ACs with tests, test plan, doc diff, and spike determination. The live re-exercise materially corrected the mechanism story: `--status success` does NOT permanently exclude attempt>=2 runs — the v0.23.0 miss was a filter-consistency lag right after re-run-to-green (the filtered query was still empty ~10 minutes after green while the unfiltered one saw the success), which is exactly the retry window; the fix is unchanged and valid under either explanation. Also surfaced two adjacent touch points: the stale query-shape comment in e2egate_workflow_test.go, and the same `--status success` shape in release.yml's journey-ledger step (different semantics, recorded as out of scope for a separate seed).
