---
title: Dispatch a gated successor before presenting its gate
status: validation
source: Captain approved separate headless contract fix on 2026-09-15
started: 2026-09-15T14:19:00Z
completed:
verdict:
score: 0.95
worktree: .worktrees/spacedock-ensign-dispatch-gated-successor-before-presentation
issue:
pr: "#798"
mod-block: merge:pr-merge
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
        - id: gate:vcfwr5ptn6pedv32wq7cbsy6:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:vcfwr5ptn6pedv32wq7cbsy6-ideation-1
              briefing:
                id: briefing:vcfwr5ptn6pedv32wq7cbsy6:ideation:attempt-1:revision-1
                digest: sha256:239ad8e6ead865fdd7e23be6ce8053a4b7988179e495b91d49c9641ef4ea586f
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:vcfwr5ptn6pedv32wq7cbsy6:ideation:1
                briefing: briefing:vcfwr5ptn6pedv32wq7cbsy6:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-09-15T14:23:15.974746Z"
                decision: approve
                reason: Captain approved this separate contract fix. Exact bounded one-file change removes conflicting stop exception; existing local live and independent sequence proof covers all criteria, no new mechanism.
                conn:
                    quote: approve.
                    source: Captain reply approving separate headless contract-fix scope on 2026-09-15; original goal authorizes dispatch and stacked PRs
              application:
                target-stage: implementation
                state: consumed
        - id: gate:vcfwr5ptn6pedv32wq7cbsy6:validation
          stage: validation
          attempts:
            - id: gate-attempt:vcfwr5ptn6pedv32wq7cbsy6-validation-1
              briefing:
                id: briefing:vcfwr5ptn6pedv32wq7cbsy6:validation:attempt-1:revision-1
                digest: sha256:b33c2b1691ff037ee6c725e8f3d0607c86bdd1c63ce4cf35418dfc8fb7f12a2a
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:vcfwr5ptn6pedv32wq7cbsy6:validation:1
                briefing: briefing:vcfwr5ptn6pedv32wq7cbsy6:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-09-15T14:44:59.739735Z"
                decision: approve
                reason: 'All three ACs independently verified at474cdb8c5: both local live journeys pass and full native validation-before-prepare order holds. Fullnormal/race completed with only proven baseline resolverfailure, no races. Stack locally verified before CI.'
                conn:
                    quote: approve.
                    source: Captain approval of separate headless contract fix on 2026-09-15; active goal grants stacked PR and Codex CI authority
              application:
                target-stage: done
                state: pending
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

## Stage Report: implementation

- DONE: Apply exact approved one-file completion clarification on committed stack parent e52ece67a, preserving lower layers and existing guards.
  Candidate `474cdb8c515ffe0ad36c9aeb3b65b3b7cecf3f80` has parent `e52ece67abdb282ab061c41ee897a8a341b22954`; exact approved shared-reference change is +3/-1, with Completion step 4 and Gate successor guard unchanged.
- DONE: Run existing focused checks before and after editing, preserve exact diff and evidence; no prose-grep or new mechanism.
  Both named commands in Test plan passed before/after (dispatch: 1.014s/0.467s; lifecycle: 4.013s/4.103s); logs and `candidate.patch` are retained under `/tmp/vcfwr5ptn6-implementation-evidence/`.
  Dispatch controls detect wrong split-root/worktree paths or entity slugs; lifecycle/gate controls reject skipped, early, or uncorrelated completion and unauthorized gate application; they do not prove model compliance with wording.
- DONE: Commit candidate and complete report for independent targeted live validation; no CI or live runs by implementation.
  Candidate is committed and owned worktree clean; independent validation owns both targeted Codex journeys and full native-event/state ordering review required by AC-1.
- DONE: Run gofmt and preserve the bounded candidate surface.
  `gofmt -w ./cmd ./internal` succeeded; only inherited unrelated alignment in `internal/release/runtime_live_evidence_workflow_test.go` was excluded, with exact diff saved as `inherited-gofmt.patch` in the evidence root.
- SKIPPED: Full normal/race checks during candidate preparation.
  FO assigned one sequential final-tip normal/race run at `474cdb8c5` after this implementation handoff; results will follow in an addendum. No full-suite or live success is claimed here.

### Summary

Removed the contradictory next-stage gate halt exception and explicitly required gated successor dispatch and verified completion before gate preparation. The approved wording-only candidate is ready for independent validation; live behavioral acceptance and final-tip full-suite results remain pending with their assigned owners.

## Stage Report: implementation — final-tip full-suite addendum

- DONE: Run final-tip full normal and race suites once, sequentially.
  Exact clean candidate `474cdb8c515ffe0ad36c9aeb3b65b3b7cecf3f80`; `go test ./...` then `go test ./... -race` used `GOFLAGS=-p=2`, `/Users/clkao/go/bin` prepended to PATH, and unset `SPACEDOCK_BIN`/`SPACEDOCK_REPO_ROOT`.
- FAILED: Full-suite green result.
  Both commands exited 1 solely on `TestCodexResolveManifestAgainstInstalledHost`: installed-host detection disagreed with resolver returning the stale `0.28.0-pre0` manifest; FO supplied independent baseline reproduction at `2a7b87198`, so no redundant baseline run was made.
- DONE: Preserve complete current-tip evidence and candidate.
  `/tmp/spacedock-codex-stack-474-full/summary.md` links sibling `normal.log`, `normal.exit`, `race.log`, and `race.exit`; all other packages passed, no race detector findings, and candidate HEAD/worktree remained unchanged and clean.

### Summary

Completed the required normal/race checks at the exact candidate; the only failure matches the independently reproduced baseline installed-host resolver defect. No candidate edits, retries, or live/CI runs were performed by this worker; full-suite success is not claimed.

## Stage Report: validation

- DONE: Independently verify exact474cdb8c5 has only approved one-file contract change and preserves lower layers, completed-stage gate ownership and freshness.
  Exact candidate `474cdb8c515ffe0ad36c9aeb3b65b3b7cecf3f80`, parent `e52ece67abdb282ab061c41ee897a8a341b22954`; +3/-1 in shared dispatch reference only, with existing Completion step 4, successor guard, authority and freshness rules unchanged; no CLI/schema/grader change.
- DONE: Run targeted local Codex headless and gate-ready journeys once serially; inspect native worker completions and durable report order before gate preparation, beyond automated assertions.
  Exact candidate-built binary, Codex 0.154.0 Luna/max, isolated local OAuth, REQUIRED unset: headless PASS 307.57s, gate-ready PASS 127.00s, total 434.840s, exit 0/no SKIP; no retry/CI/global auth changes.
- DONE: Report AC1–3 with exact evidence and source fixture AC2 disposition; no CI, candidate changes or retries without diagnosis. Full suites owned by implementation worker.
  AC-1/AC-2 and source fixture AC-2 pass from this new run; focused AC-3 controls pass; final full-suite baseline limitation is recorded below and in implementation's addendum.
- DONE: AC-1 — Headless completion enters and completes the gated successor before presenting its gate.
  Implementation report `adc8b3df` + native completion 14:30:32.120Z → validation status 14:30:48.225Z/stamp `63aa6761` → separate fresh validation report `f0001501` + native completion 14:32:32.564Z → prepare 14:32:58.341Z/final gate commit `4bd499c9`.
  Parent `01a0a578-2d41…` correlates distinct implementation `01a0a579-5d1d…` and validation `01a0a57a-7db2…` session identities; each worker appends only its own body report. No borrowed report, skipped validator, early prepare, approval/application, or handoff.
- DONE: AC-2 — A gate-ready workflow still stops without spending captain authority.
  GateGuardrail parent `01a0a57c-e1e1…` has zero worker spawn/dispatch/status mutation, one prepare/state commit `f367d0ec`; final validation entity has one open briefing and no resolution/application or handoff.
- DONE: AC-3 — The correction preserves existing worker authority and delivery semantics.
  Independent dispatch controls passed 0.537s and fixture/lifecycle/gate controls 4.424s; wrong state/worktree/slug, skipped/early/uncorrelated completion, or post-prepare authority mutation falsify them. Full normal/race results: both exited 1 solely at known `TestCodexResolveManifestAgainstInstalledHost`; all other packages passed, no race warning (`/tmp/spacedock-codex-stack-474-full/{normal,race}.log`); implementation owns the full-suite addendum.
- DONE: Reassess source fixture `separate-headless-pregate-fixture-instructions` AC-2 without editing its prior report.
  Recommend PASSED for the exact missing validation-transition boundary using the full same-run order above; original reviewer owns reconciliation. A single fresh pass supports this intervention but does not establish unique wording causality.
- DONE: Semantic adversarial pass and scoped runtime limitation.
  Ungated completion→gated successor and completed-gated→hold both exercised; terminal/blocker/freshness branches unchanged. Validation's unavailable `gate validate` invocation (exit 2, unknown subcommand) is honestly retained in its report and human package; FO scoped it as observed friction outside this task's proven order/hold ACs.

### Summary

Recommend **PASSED** for the bounded one-file candidate and source fixture AC-2; no new material finding or candidate mutation. Full normal/race suites retain the documented installed-host compatibility limitation, so this is not an all-suites-green claim. Detailed review, exact native line/timestamp references, commands, state Git history, and final tracked-state archives are in `/tmp/spacedock-stack-headless-474-live/review.md` and that artifact root; periodic backup directories may retain setup-deleted files, so Git and committed-state archives are authoritative, and no auth is retained.

## Stage Report: implementation — routine restack addendum

- DONE: Restack the owned tip onto the approved lower fixture correction without changing its patch.
  New candidate `0f5ee4862f42416bda9275ddaeff6dd67c547115` has parent `0fb8dd634d397e59ba169ec81d5d602594b02800`; clean one-commit rebase from `474cdb8c5`/`e52ece67a`, range-diff `=`, own-layer diff byte-identical (one file +3/-1), lower correction included as ancestor.
- DONE: Run required full normal then race once at the restacked tip.
  Commands `go test ./...` and `go test ./... -race` ran sequentially with `GOFLAGS=-p=2`, `/Users/clkao/go/bin` prepended to PATH, and stale `SPACEDOCK_BIN`/`SPACEDOCK_REPO_ROOT` unset; all packages except the known resolver owner passed (some Go-cached), with no race findings.
- FAILED: Full-suite green result.
  Both exited 1 only on `TestCodexResolveManifestAgainstInstalledHost`, with the identical stale `0.28.0-pre0` manifest mismatch already observed at `474cdb8c5` and independently reproduced at baseline `2a7b87198` per FO evidence; no redundant baseline run or new failure.
- DONE: Required formatting command and clean exact candidate.
  `gofmt -w ./cmd ./internal` succeeded; only the previously known unrelated two-line alignment in `internal/release/runtime_live_evidence_workflow_test.go` appeared, retained as `inherited-gofmt.patch` and excluded under existing scope authorization. Candidate HEAD remains `0f5ee4862`, clean, with no new code edits.
- DONE: Preserve evidence and existing validation/gate state.
  `/tmp/spacedock-stack-799-restacked-full/summary.md`, `range-diff.txt`, `own-layer.patch`, normal/race logs and exit files retain results; no new live/CI runs, code push, frontmatter changes, or replacement of prior validation evidence.

### Summary

Routine restack includes the lower fixture correction while preserving this task's exact approved wording patch. Both required full suites completed with only the known baseline resolver limitation; existing live validation remains evidence at its recorded prior SHA, and no new live result is claimed for the restacked tip.

## Stage Report: implementation — two-phase parent restack addendum

- DONE: Restack the owned wording patch onto the corrected approved two-phase parent.
  Final candidate `380bd7bf34b3dd1de07a125d75035c593f5ff930` atop `faef711575c3407a01c4aaa1c87d64d8acabf1bd`; clean one-commit rebase from saved `37d9a6c38`/`71d1f6b1` boundary, range-diff `=`, byte-identical own-layer patch (one file +3/-1), new parent ancestry verified.
- DONE: Retain superseded normal evidence and honor race hold.
  Intermediate `37d9a6c38` normal exited 1 with known resolver failure plus new `TestRuntimeLiveRegistryReconciliation` rejection of unregistered `TestSmallestMechanismPhaseMetrics`; `/tmp/spacedock-stack-799-twophase-full/summary.md` retains it. No race ran there; lower owner corrected the test wrapper in `faef711` before final restack.
- DONE: Run final-tip full normal then sole full race once, sequentially.
  `go test ./...` and `go test ./... -race` used `GOFLAGS=-p=2`, PATH prepended `/Users/clkao/go/bin`, and unset stale `SPACEDOCK_BIN`/`SPACEDOCK_REPO_ROOT`. Contractlint now passes (0.585s normal/3.769s race); every package except known resolver passes (some Go-cached), and no race findings occurred.
- FAILED: Full-suite green result.
  Both final commands exited 1 only on `TestCodexResolveManifestAgainstInstalledHost` with identical stale `0.28.0-pre0` manifest mismatch, matching prior logs and FO-provided independent baseline `2a7b87198` reproduction; no baseline rerun or additional failure.
- DONE: Required formatting check and exact clean final candidate.
  `gofmt -w ./cmd ./internal` succeeded; prior byte-identical unrelated release workflow alignment delta was retained as `inherited-gofmt.patch` and excluded under existing authorization. HEAD remains `380bd7bf3`, clean, with no additional candidate edits.
- DONE: Preserve final evidence and existing validation/gate data.
  `/tmp/spacedock-stack-799-twophase-final-full/summary.md` links boundary, range-diff, patch, normal/race logs and exits. No code push, CI, new live run, or frontmatter/gate changes; prior live reports retain their original exact-tip attribution.

### Summary

The final restack preserves this task's approved wording while including the corrected two-phase lower layer. Final normal/race suites completed with only the established installed-host resolver limitation; the intermediate registry failure is retained separately and absent from final results.
