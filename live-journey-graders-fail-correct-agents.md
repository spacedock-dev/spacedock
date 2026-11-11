---
title: "Four live journeys fail agents that followed the contract"
source: "Root-cause read of Runtime Live E2E runs 42883048364, 42883041637, 42916100601 and 42916110241 during the 0.28.0 release preparation, 2026-11-11"
id: 9wpk5vftfth48rcaxft6kw27
---

Four live journeys intermittently fail on runs where the first officer did what the contract requires. Each failure is a defect in the grader, fixture, or scenario runbook, so a red lane says nothing about the product, and the release gate needs reruns to go green. Target: 0.28.1.

## Defects

1. **Codex `default-headless-gate-stop` reports `gate-not-held` on a held gate.** `internal/ensigncycle/gate_assert_impl_test.go` matches `completedSet` and `verdictSetFM` (case-insensitive, multiline) against the whole entity file. A validation stage report whose Summary starts a line with "Verdict: ..." matches, and the journey fails although status stayed `validation` with no record, consume, or dispatch. Seen in runs 42883048364 and 42916100601 attempt 2; this journey failed 4 of 14 Codex jobs between 2026-11-04 and 2026-11-10 and sits inside the release gate.

2. **Pi `recorded-gate-lifecycle` asks the first officer to approve a gate whose report leaves a criterion unevidenced.** The fixture built by `recordedGateEntity()` (`internal/ensigncycle/recorded_gate_lifecycle_test.go`) has a validation report that never cites AC-1, so `status --read --ac-scan` returns it `unevidenced` with no citations. The shared core requires every criterion to be cited or named unmet, while the runbook says to approve. Runs where the first officer stops after prepare or records `revise` fail; runs where it approves anyway pass. Failed 2 of the 3 Pi jobs that reached it.

3. **Haiku `ac-value-reanchor` requires a status move the contract does not define.** The runbook in `internal/livescenario/ac2_reanchor.go` says to apply the configured route until the entity reaches the parked stage without dispatching workers, and the grader requires `status == feedback-to`. The rejection flow's first step is a worker dispatch, no step sets status to the `feedback-to` target, and the gate lifecycle forbids `status --set` to advance a gate. A first officer that records `revise` and stops fails; one that improvises `status --set` passes.

4. **Pi `filing` reports `filing-command-not-observed` when the command ran.** `commandFilesViaNew` in `internal/ensigncycle/shared_filing_test.go` requires the slug directly after `new`. `spacedock new --workflow-dir <dir> <slug>` filed the entity and still failed the journey (run 42916110241 attempt 1).

## Outcome

Each of the four journeys grades the behaviour it exists to prove, and passes when the first officer follows the contract. A lane failure on any of them means the agent or the binary did something wrong.

## Unknowns for ideation

- Defect 3 has two remedies: grade the closed `revise` with no application and status not advanced, or define a revise-to-`feedback-to` status move in the contract or binary. The second changes product behaviour and is a captain decision.
- Whether the same whole-file regex pattern in defect 1 affects other gate assertions.
- The exact citation syntax the AC scan accepts, for the defect 2 fixture.
- Two of the four Codex failures and the haiku failure on commit `eef2d13` were not read; confirm they share these causes.

## Out of scope

The Pi common-journey suite running at its 50-minute `go test` deadline.
