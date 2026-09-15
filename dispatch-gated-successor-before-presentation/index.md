---
title: Dispatch a gated successor before presenting its gate
status: ideation
source: Captain approved separate headless contract fix on 2026-09-15
started: 2026-09-15T14:19:00Z
completed:
verdict:
score: 0.95
worktree:
issue:
pr:
mod-block:
id: vcfwr5ptn6pedv32wq7cbsy6
gates:
    version: 1
    records:
        - id: gate:vcfwr5ptn6pedv32wq7cbsy6:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:vcfwr5ptn6pedv32wq7cbsy6-backlog-1
              briefing:
                id: briefing:vcfwr5ptn6pedv32wq7cbsy6:backlog:attempt-1:revision-1
                digest: sha256:b6057983672d6c156860ec510eeeb522745ea6e009b88a5035a4bab2e4142526
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:vcfwr5ptn6pedv32wq7cbsy6:backlog:1
                briefing: briefing:vcfwr5ptn6pedv32wq7cbsy6:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-09-15T14:18:52.714103Z"
                decision: approve
                reason: Captain approved adding the separate headless contract fix; CI only after complete stack individually verified locally.
              application:
                target-stage: ideation
                state: consumed
---

Remove the contradictory completion shortcut so a headless first officer runs a gated successor before preparing its review.

## Problem

The corrected pre-gate fixture dispatched and completed implementation, then attempted gate preparation while still at implementation. The loaded shared dispatch contract says to stop when the next stage is gated, conflicting with its Gate successor guard requiring that stage's dispatch and completed report first.

## Proposed approach

Remove the next-stage gate stop exception in skills/first-officer/references/fo-dispatch-core.md. Distinguish completed-stage gate routing from entering and dispatching a gated successor. Preserve the existing Gate successor guard, authority, freshness, and approval rules. Captain approved this separate stack layer; no new CLI behavior.

## Risk evidence

/tmp/spacedock-stack-headless-live/finding.md records exact failed candidate da50d61d6 and loaded conflicting text. TestLiveCommonDefaultHeadlessGateStop failed in 122.09 seconds after one implementation worker completed; no validation transition/worker occurred. GateGuardrail control passed. Wording causality is a hypothesis that requires the same targeted live journey after the correction.

## Out of scope

New state metadata, scheduler or CLI changes, changing assertions to accept an unfinished validation stage, and same-stage revision issue 792.

## Expected surface and tolerance

Estimate net LOC change: +2, across 1 file; 3 insertions and 1 deletion in `skills/first-officer/references/fo-dispatch-core.md`. Tolerance remains net +35 and 3 files. Existing proof already exercises the failing journey, so no additional prose-grep test or second file is planned. Layer branches from committed stack tip `e52ece67abdb282ab061c41ee897a8a341b22954`; preserve all lower stack commits.

Permitted observable change: the FO advances into and dispatches the gated successor after completed ungated work, then prepares review after successor completion. Command grammar, stored formats, worker/frontmatter authority, freshness, and captain approval/consumption rules remain unchanged.

## Acceptance criteria

**AC-1 — Headless completion enters and completes the gated successor before presenting its gate.**
Verified by: targeted local Codex TestLiveCommonDefaultHeadlessGateStop observes implementation completion, validation transition/worker report, then one open gate, with no approval or consumption. Restoring the shortcut is the falsifying intervention; do not assert causal certainty from prose presence.

**AC-2 — A gate-ready workflow still stops without spending captain authority.**
Verified by: targeted local Codex TestLiveCommonGateGuardrail remains passing; no successor dispatch or consumption occurs without approval.

**AC-3 — The correction preserves existing worker authority and delivery semantics.**
Verified by: focused existing dispatch/skill smoke and durable lifecycle controls; code diff contains no CLI/state schema or grader weakening. Current tip full normal/race results are recorded before completion, distinguishing any reproduced baseline defect.

## Test plan

Primary live owners at exact base `e52ece67abdb282ab061c41ee897a8a341b22954`: `internal/ensigncycle/shared_live_runner_test.go:128` (`TestLiveCommonDefaultHeadlessGateStop`) and `:166` (`TestLiveCommonGateGuardrail`), both through `runGateStopScenario` in `internal/ensigncycle/claude_live_runner_test.go:249`. The pre-gate setup and already-gated control remain unchanged.

AC-1: `assertImplementationWorkerLifecycle` in `internal/ensigncycle/claude_runtime_helpers_test.go` checks correlated implementation completion before validation transition and one parsed DONE implementation report. The retained baseline fails this owner with `spawns=1, completed=130, validation=-1`. Restoring the shortcut is the proposed behavioral falsifier; a fresh pass supports the correction but does not by itself prove wording causality.

Existing automation does not separately require a validation worker/report before prepare. Independent validation must inspect the same run's native lifecycle events, command log, and state Git history for the full order: implementation completion/report → validation status and worker dispatch → validation completion/report → gate prepare. Skipping the validation worker, borrowing the implementation report, or preparing before validation completion fails AC-1 even if the existing live test passes. This is a bounded one-off review of retained artifacts, not a new observer or prose test.

AC-2: `assertGateHeld` in `internal/ensigncycle/gate_assert_impl_test.go` and `assertRecordedGateHoldLog` in `internal/ensigncycle/claude_runtime_helpers_test.go` require the bound open validation gate and reject approval/application, successor status/dispatch, repeated prepare, or post-prepare mutation. Preserve these assertions. Existing deterministic owners `TestAssertGateHeld` and `TestAssertGateHeldAcceptsPreparedFixtureBinding` live in `internal/ensigncycle/gate_assert_test.go`; `TestGateGuardrailNegativeBrokenStateTransition` lives in `internal/ensigncycle/shared_scenarios_negative_test.go`.

AC-3: before editing, run existing dispatch integration and lifecycle negative controls. `TestSplitRootFolderWorktreeDispatch` and `TestFlatEntitySlugUnchanged` in `skills/integration/dispatch_test.go` exercise stamped dispatch paths and entity identity; a wrong state/worktree path or changed slug falsifies them. `TestImplementationLifecycleAndObserverNegativeControls` and `TestCodexNativeLifecycleUsesCorrelatedSessionHandle` in `internal/ensigncycle/claude_runtime_helpers_test.go` reject skipped/early or uncorrelated completion evidence. These are seconds-scale deterministic controls, not proof that an LLM follows the revised wording.

Focused commands from the eventual candidate checkout:

```bash
go test ./skills/integration -run 'Test(SplitRootFolderWorktreeDispatch|FlatEntitySlugUnchanged)$' -count=1
go test -tags live ./internal/ensigncycle -run 'Test(PreGateWorkflowIsStageCoherent|AssertGateHeld|AssertGateHeldAcceptsPreparedFixtureBinding|GateGuardrailNegativeBrokenStateTransition|ImplementationLifecycleAndObserverNegativeControls|CodexNativeLifecycleUsesCorrelatedSessionHandle|AssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle)$' -count=1
```

After implementation, FO schedules one serialized local Codex run per affected journey at the exact final tip (minutes; no CI until every stack layer is locally verified):

```bash
go build -o /tmp/vcfwr5ptn6-candidate-spacedock ./cmd/spacedock
SPACEDOCK_BIN=/tmp/vcfwr5ptn6-candidate-spacedock SPACEDOCK_LIVE_RUNTIME=codex SPACEDOCK_LIVE_ARTIFACT_DIR=/tmp/vcfwr5ptn6-codex-live go test -tags live ./internal/ensigncycle -run '^TestLiveCommon(DefaultHeadlessGateStop|GateGuardrail)$' -count=1 -parallel=1 -timeout=20m -v
```

Keep `SPACEDOCK_CODEX_LIVE_REQUIRED` unset. Use the proven environment/shim retained under `/tmp/spacedock-stack-headless-live` and the existing isolated-home/auth harness documented by the source task. Require observed execution with no SKIP; record Codex version, exact candidate SHA, exit status, artifact root, and durable state commits. Read `docs/runtime-support.md` and its assume-it-already-works operating prompt before treating first-contact setup failures as host limitations. Independent validation must also reassess source task `separate-headless-pregate-fixture-instructions.md` AC-2 from this run. Run `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` at final stack tip before claiming implementation complete, distinguishing any reproduced baseline failure.

## Exact instruction and documentation change

The shared dispatch reference is the behavior documentation being corrected. No separate CLI or docs-site contract changes are proposed. In `skills/first-officer/references/fo-dispatch-core.md:55`, replace the following exact base paragraph:

**A completed non-gated, non-terminal stage is not a stopping point.** After verifying the report, the FO MUST advance the entity to the next stage and dispatch it (reuse-or-fresh per the dispatch module's reuse conditions) BEFORE ending its turn. Only these conditions legitimately halt the turn here: the next stage is `gate: true` (present the gate and wait), the entity is terminal (run the merge/cleanup ceremony), an explicit blocker (a `«halt.rebase-conflict»`, an unmet clarification), or a captain decision the contract requires. Absent one, stopping after a completion-only report is a contract violation.

with:

**A completed non-gated, non-terminal stage is not a stopping point.** After verifying the report, the FO MUST advance the entity to the next stage and dispatch it (reuse-or-fresh per the dispatch module's reuse conditions) BEFORE ending its turn. Only these conditions legitimately halt the turn here: the entity is terminal (run the merge/cleanup ceremony), an explicit blocker (a `«halt.rebase-conflict»`, an unmet clarification), or a captain decision the contract requires. Absent one, stopping after a completion-only report is a contract violation.

If the successor is gated, enter and dispatch it under the Gate successor guard; prepare and present its gate only after that stage's completion signal and verified report. Gate routing at completion applies to the completed stage, not an unentered successor.

Preserve Completion step 4 and the existing Gate successor guard byte-for-byte. Removing only the exception is the simplest alternative; the added sentence states which completion owns gate routing and explicitly points to the existing guard, addressing the observed premature prepare without introducing a new mechanism. This correction serves AC-1; no new CLI, scheduler, metadata, or host adapter is needed.

## Ideation proof and limits

No new spike needed: retained failed-run commands and state already exercise implementation dispatch/completion and the real binary's premature-prepare refusal; the gate-ready control passed. The source task's deterministic setup/control proof and exact-tip source reads establish existing proof owners. Whether the wording correction changes the model's choice remains unverified until the scheduled live intervention; ideation adds no live run or architecture exploration.

Retained evidence: `/tmp/spacedock-stack-headless-live/finding.md`, `commands-summary.json`, `retained-rollouts/`, `retained-workflows/TestLiveCommonDefaultHeadlessGateStop4234871294/`, and `codex-shared-scenarios/default-headless-gate-stop/` beneath the same artifact root. The source entity is `docs/dev/.spacedock-state/separate-headless-pregate-fixture-instructions.md` (flat file). No product source was edited during ideation.

### Feedback Cycles

## Stage Report: ideation

- DONE: Turn the approved bounded wording correction into exact before/after and verify existing proof-owner paths using retained failed-run evidence; no new architecture.
  Exact base `e52ece67abdb282ab061c41ee897a8a341b22954`, shared dispatch paragraph, source task, and retained finding identify the conflicting shortcut and existing live/deterministic owners.
- DONE: Confirm smallest surface and existing targeted local headless/gate-ready proof plan, preserving lower stack layers and all authority assertions.
  Plan is one shared-reference file (+3/-1, net +2), keeps all assertions, and requires independent native-event/state review for validation completion beyond the current implementation-lifecycle assertion.
- SKIPPED: New live execution and implementation checks during ideation.
  This stage owns the task body only; existing retained failure is the baseline, and candidate checks/live intervention are scheduled after implementation at exact final tip.

### Summary

Specified the exact correction and reduced the expected surface to one file because existing tests own the failing journey. Recorded the automation's validation-worker gap explicitly and assigned full-sequence evidence review to independent validation; no live success or wording causality is claimed.
