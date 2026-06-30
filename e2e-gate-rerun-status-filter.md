---
title: "e2e-gate --status success filter hides re-run-to-green Runtime Live E2E runs from the release gate"
status: backlog
score: 0.7
source: "v0.23.0 stable cut, 2026-06-30. The release e2e-gate blocked the cut: gh run list --workflow 'Runtime Live E2E' --status success -c <commit> returned [] for a run that was re-run to green (run_attempt 2), because GitHub's --status success filter only matches first-attempt successes. The genuine green matrix (run 28429490220, all 5 lanes success on the tagged SHA e3f85ec3) was invisible to the gate. Worked around the cut by dispatching a fresh run_attempt-1 green run; this fix removes the trap."
id: q3v95nmwkwh656p5m6rhher1
---

The release e2e-gate's `gh run list` query uses `--status success`, which on GitHub only matches a run whose FIRST attempt concluded success. A run that flaked on one lane and was re-run to green via `gh run rerun --failed` has run_attempt >= 2 and is EXCLUDED by `--status success`, even though its overall `conclusion` is `success` on the exact tagged SHA. This directly contradicts the proof-policy / handoff guidance to "re-run the failed lane to green, never block the tag" — re-running to green produces a run the gate cannot see, so the cut blocks at e2e-gate (goreleaser skipped).

## Problem
`cmd/spacedock-release/e2e_gate.go` (`ghRunListForCommit`) builds `gh run list --workflow "Runtime Live E2E" --status success -c <commit> --json databaseId,headSha,conclusion,status`. The `--status success` pre-filter is both redundant AND harmful: `internal/release/e2egate.go` (`EvaluateE2EGate`) ALREADY checks `run.Conclusion == "success" && run.HeadSha == releaseCommit`. Live evidence from the v0.23.0 cut: `--status success -c e3f85ec3` -> []; `-c e3f85ec3` (no --status) -> the run with conclusion=success, run_attempt=2.

## Proposed approach
Drop `--status success` from `ghRunListForCommit`. The `-c <commit>` filter scopes to the tagged commit and `EvaluateE2EGate` filters for conclusion==success, so the predicate stays correct and now sees re-run-to-green runs.

## Acceptance criteria
- **AC-1 (VALUE)** — a Runtime Live E2E run that was re-run to green (overall conclusion success on the tagged SHA, run_attempt >= 2) PASSES the e2e-gate. Verified by: a fixture-driven test where the run-list JSON's only success entry is a re-run-to-green run, asserting the gate accepts it; RED on the current `--status success` query shape, GREEN after dropping it.
- **AC-2** — a parked/waiting run (empty conclusion) still does NOT satisfy the gate (no regression of the original spike intent: parked offline-only runs must not qualify).
