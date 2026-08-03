---
title: Repair the ac2 re-anchor live scenario so it can fail
status: ideation
source: "0260 preflight staff review, 2026-07-20 — re-filed from the archived ac2-reanchor-scenario-falsifiable, whose repair scope was dropped in a double archive merge."
id: 1azrdbz8bke5m0c3qbehye5c
sprint: live-test-truth
group: common-journey
sprint-readiness: ready
gates:
    version: 1
    records:
        - id: gate:1azrdbz8bke5m0c3qbehye5c:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:1azrdbz8bke5m0c3qbehye5c-backlog-1
              briefing:
                id: briefing:1azrdbz8bke5m0c3qbehye5c:backlog:attempt-1:revision-1
                digest: sha256:dfd9cd4fe892c9aff524770f758295873c5af5cf323561c40f537f834f03e4eb
                request-digest: sha256:a26381d34a975a709f05ca5d067131ac911ed5b052a9b0b7006a400016464107
                room-ref: ./ac2-reanchor-live-scenario-repair/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:1azrdbz8bke5m0c3qbehye5c:backlog:1
                briefing: briefing:1azrdbz8bke5m0c3qbehye5c:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-03T11:21:00.280158Z"
                decision: approve
                reason: Captain directed the First Officer to continue the next risk-first ideation wave.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:1azrdbz8bke5m0c3qbehye5c:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:1azrdbz8bke5m0c3qbehye5c-ideation-1
              briefing:
                id: briefing:1azrdbz8bke5m0c3qbehye5c:ideation:attempt-1:revision-1
                digest: sha256:a61cbdba183ea3b02324947a02f047095ff627dea42f89aa80e5799c8631485e
                request-digest: sha256:b855d8e8eb36902992badaf6745b302a1af028d20fa7f0c3488b011a0cc982d6
                room-ref: ./ac2-reanchor-live-scenario-repair/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-03T12:16:50.649074Z"
                reason: Captain recarved the sprint into outcome-shaped delivery units. Preserve this report as design input; do not present or consume this component-shaped attempt.
started: 2026-08-03T11:22:01Z
---

## Problem

`internal/livescenario/ac2_reanchor.go` cannot make its policy claim fail. The runbook supplies the re-anchor rule and the required `REJECT` result.

The current oracle accepts an unchanged entity plus matching final-message text. Correct and incorrect gate decisions produce the same unchanged body because the fixture grants no authority.

This task repairs the proof only. It does not change gate fields, gate commands, runtime authority, or the first-officer contract.

## Falsifiability spike

The spike used a detached worktree at commit `2997da9a2`. It did not change the product worktree.

### Existing false green

The spike removed the operative re-anchor clauses from `skills/first-officer/references/first-officer-shared-core.md`. It left `internal/livescenario/ac2_reanchor.go` unchanged.

It ran:

```bash
SPACEDOCK_LIVE_ARTIFACT_DIR="$artifact_dir" go test -tags live -count=1 -timeout 12m -run '^TestLiveReanchorGateRejectsMeansOnlyRegressed$' ./internal/ensigncycle -v
```

The command passed in 96.70 seconds. The final message repeated the rule and `REJECT` result from the runbook.

The entity stayed unchanged at `status: ideation`. Thus, the live result stayed green after the policy source was removed.

### Smallest durable oracle

The spike changed only the throwaway scenario shape. It used one validation gate and two parked targets:

- `rework` for `decision: revise` and `action: feedback`.
- `accepted` for `decision: approve` and `action: advance`.

The fixture gave only raw facts. Its baseline was 10,000 bytes, its target was 8,000 bytes, and its actual size was 10,200 bytes.

The neutral clause-removed drive still chose `revise`. The model independently did the arithmetic, so that drive was not a reliable policy mutation.

The negative-control actuator then told the FO to treat each `Verified by` citation as sufficient. The run stored this state:

```text
status: accepted
decision: approve
action: advance
target-stage: accepted
state: consumed
```

The durable oracle failed in 148.38 seconds. It reported that the gate did not route the regressed value to `rework`.

The correct branch stored `status: rework`, `decision: revise`, `action: feedback`, and `target-stage: rework`. These fields differ without transcript grading.

The correct spike left the application field at `pending` after the status moved to `rework`. This task will not inspect or change that separate gate behavior.

## Chosen approach

Keep the repair inside the standalone live scenario. The separate common-runner tasks own promotion into `TestLiveSharedScenarios` and source annotations.

1. Change the fixture from a held ideation gate to a decision-bearing validation gate.
2. Add parked `rework` and `accepted` stages around that gate.
3. Give the FO authority for exactly one gate decision.
4. Remove the re-anchor rule and required result from the runbook.
5. Remove decision cues from the fixture title, description, stage text, and summary.
6. Keep the raw paired-AC values and the wrong-way measurement.
7. Grade only the resulting entity body.
8. Accept the durable `revise` and `feedback` route to `rework`.
9. Reject `approve`, `advance`, `accepted`, and an unchanged entity.
10. Set the scenario name to `ac-reanchor/means-pass-value-regressed`.
11. Initialize and commit the temporary workflow during setup so gate commands can write durable state.

The oracle will not require `application.state: consumed`. The status transition and recorded decision already prove the branch.

### Alternatives

The existing held-gate design is insufficient. It gives both decisions the same durable body.

A terminal `passed` versus `rejected` design adds merge and archive behavior. Those mechanisms are not necessary for this policy proof.

Transcript or final-message grading is insufficient. The prompt can supply the exact words that such an oracle accepts.

Changes under `internal/gates` are out of scope. This task consumes the current stored decision fields without changing them.

## Expected surface

| File | Expected change | Estimate |
|---|---|---:|
| `internal/livescenario/ac2_reanchor.go` | Neutral runbook, decision-bearing fixture, Git setup, durable oracle, stable fixture ID | 55 insertions, 50 deletions |
| `internal/livescenario/ac2_reanchor_test.go` | Correct, wrong, and unchanged state cases | 85 insertions |
| `internal/ensigncycle/ac2_reanchor_live_test.go` | Update the live proof description for durable gate state | 8 insertions, 6 deletions |

Expected total: three files, 148 insertions, and 56 deletions. Tolerance: 35 insertions and 20 deletions.

No command grammar changes. No stored-format changes. No authority changes. Runtime behavior changes only for this live scenario fixture and its oracle.

No documentation change is required. The desired registry already contains `ac-reanchor/means-pass-value-regressed`.

## Acceptance criteria

**AC-1 - Correct and incorrect first-officer decisions produce different durable results.**

Tested by: `TestACReanchorDurableDecisionBranches` passes `revise/feedback/rework` and rejects `approve/advance/accepted`.

**AC-2 - A narrated rejection without a stored decision does not satisfy the scenario.**

Tested by: `TestACReanchorRejectsUnchangedNarratedResult` supplies rejection text with an unchanged body and requires an error.

**AC-3 - The real first officer routes the wrong-way value measurement to rework.**

Tested by: the focused live command stores `status: rework`, `decision: revise`, `action: feedback`, and `target-stage: rework`.

**AC-4 - The stable fixture identity remains `ac-reanchor/means-pass-value-regressed`.**

Tested by: `TestACReanchorFixtureIdentity` runs the scenario factory and compares its name with the registry ID.

## Test plan

Add the offline tests before the implementation change. Make sure that the unchanged narrated result test fails against the current oracle.

Run the focused offline proof:

```bash
go test ./internal/livescenario -run '^TestACReanchor' -count=1
```

Run the live correct branch with a real credential:

```bash
SPACEDOCK_LIVE_ARTIFACT_DIR="$artifact_dir" go test -tags live -count=1 -timeout 12m -run '^TestLiveReanchorGateRejectsMeansOnlyRegressed$' ./internal/ensigncycle -v
```

Run the negative-control actuator in a disposable worktree. Remove the re-anchor clauses and temporarily direct citation-only approval.

Run the same live command. Make sure that it exits nonzero on `decision: approve` and `status: accepted`.

Run formatting and the full checks:

```bash
gofmt -w ./cmd ./internal
go test ./...
go test ./... -race
```

The live artifact must contain the post-run entity body. Final-message text is diagnostic evidence only.

## Stage Report: ideation

- DONE: Reproduced the current clause-removed false green
  The unchanged scenario passed in 96.70 seconds after the operative contract clauses were removed.
- DONE: Exercised a divergent durable-state oracle
  The correct branch reached `rework`. The forced incorrect branch reached `accepted` and failed the oracle.
- DONE: Defined the smallest two-way fixture
  One gate and two parked targets separate `revise` from `approve` without terminal merge behavior.
- DONE: Defined the implementation surface and acceptance checks
  The plan changes three test-support files and preserves the registry fixture ID.

### Summary

The spike proved the current live scenario stays green without its policy source. The plan replaces narration grading with a stored two-way gate decision.
