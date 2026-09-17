---
title: Honor workflow-declared same-stage revision without mandatory reviewer machinery
status: validation
source: Captain request after email-triage FO issue 792
issue: spacedock-dev/spacedock#792
score: 0.95
started: 2026-09-15T04:30:03Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-workflow-owned-same-stage-revision
pr: "#801"
mod-block: merge:pr-merge
id: zz1yqc2w2katp28wpa8nghx2
gates:
    version: 1
    records:
        - id: gate:zz1yqc2w2katp28wpa8nghx2:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:zz1yqc2w2katp28wpa8nghx2-backlog-1
              briefing:
                id: briefing:zz1yqc2w2katp28wpa8nghx2:backlog:attempt-1:revision-1
                digest: sha256:ea21af7990ca6e4a52dac540a434ad088b3bbecca35dcb03cbd1c70d39950ea6
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:zz1yqc2w2katp28wpa8nghx2:backlog:1
                briefing: briefing:zz1yqc2w2katp28wpa8nghx2:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-09-15T04:29:46.231251Z"
                decision: approve
                reason: Captain explicitly requested local filing and ideation dispatch to spike actual same-stage revision and inspect all related Roborev settings/history.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:zz1yqc2w2katp28wpa8nghx2:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:zz1yqc2w2katp28wpa8nghx2-ideation-1
              briefing:
                id: briefing:zz1yqc2w2katp28wpa8nghx2:ideation:attempt-1:revision-1
                digest: sha256:2641e880a198f5610f18f87c23b7c5ecb8fc87f393f2d2ed67f276bf3c5c3a7b
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:zz1yqc2w2katp28wpa8nghx2:ideation:1
                briefing: briefing:zz1yqc2w2katp28wpa8nghx2:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-09-15T15:51:58.903625Z"
                decision: approve
                reason: 'Captain approved the returned design and independent review in chat: consider it approve. Implement conditional workflow-owned correction/review obligations with existing guards and live proof.'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:zz1yqc2w2katp28wpa8nghx2:validation
          stage: validation
          attempts:
            - id: gate-attempt:zz1yqc2w2katp28wpa8nghx2-validation-1
              briefing:
                id: briefing:zz1yqc2w2katp28wpa8nghx2:validation:attempt-1:revision-1
                digest: sha256:24ad88ebe6562841eea4aab4abaa54b1be456a6c4dc0c6408bbe3f446215b4bd
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:zz1yqc2w2katp28wpa8nghx2:validation:1
                briefing: briefing:zz1yqc2w2katp28wpa8nghx2:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-09-16T02:13:15.314105Z"
                decision: approve
                reason: Captain approved validation in Subspace resolution:binding-1789524384326931000; accepts direction and presented evidence.
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:zz1yqc2w2katp28wpa8nghx2-validation-2
              briefing:
                id: briefing:zz1yqc2w2katp28wpa8nghx2:validation:attempt-2:revision-1
                digest: sha256:fe9a4b043e5dcff6c02c64cc4758d2acdbf88d16837f72e2621b1fb22996e232
                room-ref: '@review/validation/briefing-2'
              resolution:
                type: Resolution
                id: resolution:spacedock:zz1yqc2w2katp28wpa8nghx2:validation:2
                briefing: briefing:zz1yqc2w2katp28wpa8nghx2:validation:attempt-2:revision-1
                by: person:captain
                at: "2026-09-16T20:44:34.918799Z"
                decision: approve
                reason: Captain binding resolution binding-1789591360273000000 approves all three named corrected validation snapshots in /tmp/stack-corrections-review.md for stack publication and final-tip CI. Restack preserves all patches and tree. No merge authority or live failure waiver.
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:zz1yqc2w2katp28wpa8nghx2-validation-3
              briefing:
                id: briefing:zz1yqc2w2katp28wpa8nghx2:validation:attempt-3:revision-1
                digest: sha256:b24f90d366e513f50c50fab7ba6db4745284410ba8232f33818adde5797add84
                room-ref: '@review/validation/briefing-3'
              resolution:
                type: Resolution
                id: resolution:spacedock:zz1yqc2w2katp28wpa8nghx2:validation:3
                briefing: briefing:zz1yqc2w2katp28wpa8nghx2:validation:attempt-3:revision-1
                by: person:captain
                at: "2026-09-17T15:04:35.671918Z"
                decision: approve
                reason: Captain binding resolution:binding-1789657416015652000 approves validation attempt 3 and candidate 67c68536 for corrected-stack publication and full tip CI. Reviewed artifact sha256:090f392bf3872fe4169655453f82ff67e908d3c287051a65df70a7be0c2247ac. Merge and failure waiver remain unauthorized.
              application:
                target-stage: done
                state: pending
review-round:
    id: round:zz1yqc2w2katp28wpa8nghx2:validation:1
    stage: validation
    cycle: 1
    briefing:
        id: briefing:zz1yqc2w2katp28wpa8nghx2:implementation:round-1
        digest: sha256:cffdb429dbee76cb829c1603afe3fb976c3d0bf4a0c0910864f759801856a4f7
        room-ref: '@review/validation/round-1'
---

Make same-stage gate revision follow the declared workflow while preserving independent review where it is required.

## Problem

Issue https://github.com/spacedock-dev/spacedock/issues/792 reports an email-triage workflow with feedback-to equal to its own gated stage, no independent reviewer stage, and no canonical review-round log. The shared feedback-rejection-flow nevertheless mandates both round publication and reviewer rerun. This can manufacture absent machinery or stall legitimate correction.

The captain suspects same-stage revision already shipped alongside Roborev and the development workflow's implementation-stage review. Do not assume new binary machinery is needed. The skill is unchanged from v0.27.2 to v0.28.0-pre3; the five-step rewrite at 884b55af0 retained earlier assumptions.

## Proposed approach

Keep the shipped gate and dispatch mechanics. Change only the generic feedback skill's applicability rules.

1. Route the authorized correction to the declared `feedback-to` target.
2. Publish a round only when the workflow requires a canonical review round. A Feedback Cycles projection is a separate requirement.
3. Commit the corrected entity and all required round evidence before continuing.
4. Run independent review only when the workflow requires it. Stage equality alone cannot waive an independent review obligation.
5. Prepare and present one fresh gate after all applicable obligations complete. Stop with that gate open.

The self-feedback case skips absent machinery. Missing required evidence still stops the flow. Preserve authorization, review independence, cycle limits, and stale-attempt guards.

### Exact proposed skill change

In `skills/feedback-rejection-flow/SKILL.md`, replace unconditional step 2 publication with:

> When the active workflow requires a canonical correction round, record it exactly once with the existing recorder. Otherwise continue without creating a round.
> Apply the workflow's Feedback Cycles projection separately. Its absence does not waive a required canonical round.

Replace unconditional step 4 and step 5's reviewer-only precondition with:

> When the workflow requires independent review, obtain its fresh verdict before gate preparation. Never use the correction worker as its own reviewer.
> Otherwise, completed and committed correction work satisfies the review prerequisite. Same-stage feedback alone does not waive declared review requirements.
> Prepare and present exactly one fresh open gate after every applicable requirement completes. Stop without approving or consuming it.

Update each conditional step's completion condition to match its applicability. An absent obligation completes without manufacturing its evidence.
Preserve cycle-limit escalation even when no reviewer is declared. Repeated correction cannot bypass the limit by skipping reviewer reruns.

Keep step 1's distinct authorization requirement. A captain's concrete revise instruction supplies correction authority for workflows without a review-finding checkpoint.
Do not manufacture development-specific classification policy for those workflows.

## Risk evidence

The retained [command transcript](artifacts/commands.json) records actual subprocess arguments, exits, output, and working directories.
[Git bundles](artifacts/self.bundle) preserve the disposable workflow history. Artifact folders retain the final files and frozen gate rooms.
The three Python drivers retain the exact local exercise. Their absolute paths describe this run, rather than a portable test harness.

- Binary: installed `spacedock 0.28.0-pre3`; inspected source: `origin/main` at `2a7b8719843e40b79545f0bb4def6609cdd9ebbf`.
- Self-feedback: `plan` declares `feedback-to: plan`, with no separate reviewer, round, or Feedback Cycles projection.
- Real `gate record --decision revise` closed attempt 1. Real `dispatch build --feedback-reflow` emitted a correction assignment for `plan`.
- The manual spike operator corrected and committed the plan. Commit `a8d4a6a` contains the correction and completion report.
- `gate consume` refused the rejected attempt. A later approval request failed with `attempt gate-attempt:task-plan-1 is frozen closed`.
- Fresh preparation produced attempt 2. Status showed exactly one `awaiting-captain` row. The final tree was clean at `7d9d772`.
- Frozen input and old briefing bytes remained identical. The corrected plan matched the separately seeded frozen input.
- No approval or consume command targeted the new attempt. No email or external execution occurred.
- Advisory implementation control copied the shipped Roborev fixture. Publication retained five entries and left the implementation stage unchanged.
- Missing canonical log and an unanswered reviewer rejection both failed. The unanswered rejection reported that nothing was recorded.
- Separate-reviewer control used the dev stage topology. Correction and validation dispatch envelopes targeted different stages.
- Negative control: the binary accepted premature preparation without fresh validation. The probe gate was withdrawn before manual validation and fresh preparation.

The negative control confirms a proof boundary: workflow prose owns independent-review enforcement. CLI success cannot prove that an agent follows that prose.
This ideation contains manual CLI execution and operator reports, not live spawned-worker completion or an independent human review.
Implementation must prove the changed skill with a live workflow drive. No new binary validator follows from this boundary.

### Roborev configuration and history

Current `.roborev.toml` defines correctness (`codex`, `gpt-5.6-sol`) and product (`claude-code`, `opus`) reviewers.
Its `branch_final` panel combines both reviewers with Codex synthesis. Review guidelines cover compatibility, trust boundaries, behavioral proof, adversarial review, and release-scope classification.
Only Material findings and decisions enter `## Review Findings`. These settings govern review evidence, not gate routing.

`roborev quickstart` reports a running daemon, missing post-commit hook, no configured default agent, and no Codex agent hook or installed fix/refine skills.
Repository registration returned daemon HTTP 500. Claude hook configuration and global Roborev configuration were unreadable due to OS permissions.
These are audit limits, not evidence that the hidden configuration is absent. No installation or configuration mutation occurred.
The actual Git hook directory contains a Spacedock state-validation pre-commit hook and no post-commit hook.
The workflow mods contain communication and PR-merge behavior, without Roborev routing.

Current dev topology is ungated implementation, then fresh gated validation with `feedback-to: implementation`.
Its review-finding disposition section requires separate FO authorization, workflow-owned Cycle projection, and canonical round recording.
Roborev corrections can occur within implementation. This is an advisory correction loop, not an implementation self-feedback gate.
The sprint's `docs/roadmap/0260-proportionality/commander-debrief.md` explicitly requires `branch_final` before validation.

History separates these mechanisms:

- `f06cce04a` introduced apply-once recorded gate approvals.
- `c355fbe44` introduced complete advisory rounds, including the Roborev implementation fixture and dev workflow integration.
- `ea6723acb` removed workflow policy from the recorder boundary. The related backlog task still owns remaining recorder semantics.
- `884b55af0` compressed rejection handling to five steps. It retained unconditional publication and reviewer rerun.
- Current installed feedback skill matches this unconditional shape. The file remains unchanged between the cited release endpoints.

Exact status queries found `roborev-validation-hook` and `roborev-workflow-setup-skill` in ideation.
`portable-state-checkout-roborev-followup` is in implementation, and `workflow-neutral-advisory-round-recorder` is in backlog.
`simplify-feedback-rejection-flow` is done, PR 718.
The setup task's retained draft proposes quick/code_completion panels, bounded convergence, and implementation-exit evidence. Those proposals are not the current dev workflow.

## Out of scope

Codex release-test stack fixes, archive durability (#790), terminal output clarity (#791), AC scanner naming (#793), neutral-recorder schema/taxonomy changes, and live email execution.

## Expected surface and tolerance

Estimate net LOC change: +45, across 3 files; approximately 65 insertions and 20 deletions.
Original tolerance: net +90 LOC maximum and 4 files.
Earlier amended tolerance: net +200 LOC maximum and 4 files, explicitly approved by the captain on 2026-09-15 to cover the complete behavioral matrix. The original +45 LOC estimate remains planning history; scope and acceptance criteria are unchanged.

Earlier final tolerance: net +220 LOC maximum and 4 files, explicitly approved by the captain on 2026-09-15 after the complete proof matrix and authorized fixture fixes.

Historical targeted-proof tolerance: net +236 LOC maximum and 5 files, bound by resolution:binding-1789513304916529000. The captain revoked that coverage exemption on 2026-09-16; it is not current acceptance policy.

Expected owners: `skills/feedback-rejection-flow/SKILL.md`, `internal/ensigncycle/claude_live_runner_test.go`, and `internal/ensigncycle/claude_runtime_helpers_test.go`.
The conventional rejection runner remains unchanged. `runSameStageRevisionJourney` owns the six same-stage variants through the existing common-journey architecture.
Current authorized owners also include `shared_live_runner_test.go`, `scheduled_live_test.go`, `docs/runtime-live-ci-registry.md`, and `internal/contractlint/live_registry_reconciliation_test.go`. The 2026-09-16 captain correction authorizes routine registration and an existing-checker coverage guard: seven total layer files, estimated +100–150 correction net lines above the retained +252 candidate. Current correction is +141 net; total layer is +393 net.
No command grammar, persisted format, authorization rule, recorder taxonomy, or runtime adapter changes are planned.
The observable change is that workflows without round/reviewer obligations can complete same-stage correction and re-gate.
The exact skill wording above is the documentation diff. No additional command-reference change is needed.

## Acceptance criteria

**AC-1 — A self-feedback gate can correct and re-present without invented reviewer or round artifacts.**
Verified by: a live drive observes same-stage worker completion, a committed corrected plan, and exactly one fresh open gate.
Frozen input and old attempts remain unchanged. The new attempt has no resolution or application, and no external action occurs.
The independent baseline is the frozen input, not the generated plan or instruction text.

**AC-2 — Declared review and round obligations remain enforced independently of stage equality or projection presence.**
Verified by: conventional separate-reviewer and same-stage-with-review controls stop when required evidence is missing.
A required round without Feedback Cycles still records. An absent round contract never creates recorder artifacts.
The current dev/Roborev advisory control preserves its required review entries and unchanged stage.

**AC-3 — The fix changes only generic routing applicability.**
Verified by: before/after live runs differ on the no-reviewer self-feedback outcome while required-review controls remain strict.
The existing CLI guards continue to refuse stale authority and incomplete canonical rounds.
All six same-stage variants must be registered for routine supported-runtime CI. Prior local evidence remains historical evidence; final live acceptance comes from updated stack-tip CI. A targeted implementation-only category cannot exempt a test from selection. Intentional non-gating experiments require an explicit reason.
No new flags, state fields, recorder interpretation, or synthetic reviewer stage are introduced.

## Test plan

`runClaudeRejectionFlowScenario` in `internal/ensigncycle/claude_live_runner_test.go` owns agent routing and completion proof.
Extend its workflow inputs and the topology grader in `internal/ensigncycle/claude_runtime_helpers_test.go`.
`internal/ensigncycle/shared_live_runner_test.go` registers the runner. No new runner is planned.
Before implementation, add the smallest failing self-feedback fixture. Run current and candidate skills through the same serialized local runtime lane.
A valid paraphrase must pass. Restoring unconditional recorder/reviewer requirements must fail the self-feedback case.
Removing the required-review condition must fail the separate-reviewer and same-stage-required-review controls.
Conflating projection absence with round absence must fail a required-round/no-projection control.
Leaving unconditional completion conditions must fail the self-feedback case. Removing escalation must fail a no-reviewer cycle-limit control.
Assertions inspect committed plan bytes, dispatch/completion identity, attempt count, resolution absence, and retained required evidence.
Transcript phrases and instruction substring checks are not behavioral oracles.

Reuse `internal/ensigncycle/feedback_test.go` for dispatch transport and `internal/gates/round_test.go` for canonical round rejection.
Reuse gate lifecycle/eligibility tests for immutable old attempts. These tests already own their binary failure modes.
The retained CLI spike seeds fixtures; it is not an additional permanent framework.
Expected cost: one small fixture extension and a bounded live matrix, approximately 30–60 minutes plus runtime queue time.
Implementation runs focused checks, `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`.
The shipped-contract change also requires the existing detached adversarial audit.

### Feedback Cycles

- Cycle 1: REJECTED — V1 evidence defect / retained validator; surface 5 files / +252 net LOC vs estimate 3 files / +45 net LOC (460%); AC unchanged. Captain approved amended +252/five ceiling and eight replacement drives in resolution:binding-1789517330520748000; producer correction d46f74245 is complete, independent re-review pending.

## Stage Report: ideation

- DONE: Trace relevant Roborev settings, shipped same-stage implementation/revise behavior and history, distinguishing existing mechanics from pending designs.
  Current source/config audit and exact status queries distinguish shipped advisory rounds from the setup draft; global configuration access limits are explicit.
- DONE: Spike actual on-disk same-stage correction/re-gating with real CLI and retained state, including no-reviewer/log and required-review controls, without production email writes.
  `artifacts/commands.json` and three Git bundles preserve real CLI exits, committed corrections, old/new attempts, and negative controls.
- DONE: Propose the smallest evidence-backed generic-flow correction with precise scope, AC proof, and implementation test plan; stop at ideation review.
  The design changes only applicability in the shared skill and extends existing behavioral proof owners; live agent compliance remains implementation proof.
- SKIPPED: Live spawned-worker and independent-reviewer run during ideation.
  The CLI established mechanics; the negative reviewer control identifies the exact skill proof that implementation must supply.
- SKIPPED: Full Go suite and formatting during ideation.
  No shipped Go, skill, or production workflow changed. Implementation owns the required repository checks.

### Summary

Existing CLI mechanics already support same-stage correction and fresh open gates. The generic skill imposes reviewer and recorder obligations without checking the workflow.
The smallest fix makes those obligations conditional while retaining declared review requirements. Ideation stops for captain review.

## Staff review disposition

The independent staff review approved the minimal applicability change with no blocking findings.
The First Officer authorized these planning corrections; implementation remains unapproved.

- Accepted: correct proof ownership. The rejection runner is `runClaudeRejectionFlowScenario`, registered by `shared_live_runner_test.go`.
  The design now names its implementation file and topology grader. The three-file estimate and four-file tolerance remain unchanged.
- Accepted: make completion conditions conditional with their steps. Otherwise absent recorder or reviewer obligations can still block completion.
  The proposed wording and live control now cover that failure.
- Accepted: preserve cycle-limit escalation without a reviewer. Skipping a reviewer rerun cannot permit unlimited correction cycles.
  The design now requires a no-reviewer cycle-limit control.

No shipped files changed. The retained CLI spike was not rerun.

## Implementation proof findings and FO authorizations

All four findings below were proposed as Material, owned by this task, with disposition FIX. The First Officer sent a distinct FIX authorization before each remedy. No skill semantic expansion resulted.

- **Selected revision proof:** Normal workflow: a same-stage correction followed by gate preparation. Harm: checking the latest plan and a fresh attempt independently can accept a gate prepared before the correction. Authority: value-ac[AC-1] the fresh gate must represent the committed correction. Trigger: the focused old-current-file rule returned success for a stale selected revision and an unrelated selected artifact. FO authorized checking the exact intended plan source with existing canonical Git helpers, then explicitly allowed canonical context References as well as artifacts. The fixed test rejects stale/unrelated selections and accepts both valid selection forms; retained after-run Reference pins corrected bytes at `31902baa85ecca70a76ae02d4768a6d205fb0587`.
- **Fixture identity:** Normal workflow: a required canonical correction round. Harm: a fixture without explicit identity stops at setup rather than exercising its round obligation. Authority: value-ac[AC-2] required round evidence must be meaningfully exercised. Trigger: the first round-required drive returned `entity has no identity for correction round`. FO authorized adding `id: recorded-gate-task` before seed commits and rerunning affected round controls. The revised fixture passes focused real preparation/recording setup tests.
- **Projection false rejection:** Normal workflow: a worker reports a legitimate hold without an optional projection. Harm: an ordinary sentence was misreported as an invented projection. Authority: value-ac[AC-2] absence of optional machinery must be evaluated from actual workflow state. Trigger: `No Feedback Cycles projection or finding classifications were invented.` matched a broad prose substring despite no corresponding heading. FO authorized using existing Markdown heading matching and retaining that sentence as a focused negative; both the focused negative and retained cycle-limit heading regrade pass.

### Final tolerance and fixture alignment authorization — 2026-09-15

The captain explicitly approved +220 net lines across four files, superseding the earlier +200 tolerance while preserving the original estimate as history. FO authorized FIX for the owned AC-1 fixture-alignment finding.

- Materiality: Material; normal required-round re-entry selected reports and omitted the intended corrected plan.
- Ownership: Owned fixture defect.
- Harm: An inherited successor-dispatch acceptance criterion directs review away from the actual plan correction.
- Authority and trigger: value-ac[AC-1]; retained identity-fixed round-required evidence omitted the plan and recommended rejection against the unrelated successor criterion.
- Disposition: FIX, distinctly authorized by FO after captain approval; replace that acceptance criterion with plan-versus-frozen-input acceptance and explicitly declare selected/plan.md as the intended review source. Preserve all skill, gate-policy and unique falsifier semantics. Rerun only the affected required-round live control.

### Pending registration finding — 2026-09-15

- Normal workflow: running the required full normal and race suites on the frozen candidate.
- Harm: the new standalone live test is absent from the live registry and both suites reject it.
- Authority: repository live-test registration contract and the implementation checklist requiring local validation.
- Trigger: both suites fail TestRuntimeLiveRegistryReconciliation because TestLiveCommonSameStageRevision lacks an adjacent registered live-journey binding.
- Classification: Material, owned; proposed FIX awaits a distinct FO disposition. Move the six variants into the existing registered TestLiveCommonRejectionFlow while retaining its canonical liveJourney call, runner and every assertion. This avoids a fifth registry file and reduces the delta. A name-only exemption would still violate the standalone registration contract and is not proposed.
- Current evidence: the final required-round live run passes in 172.33 seconds; both full suites otherwise fail only the baseline installed-host resolver. All processes have completed. Public logs and the passing live bundle are retained under artifacts/implementation. No code edits occurred during checks.

### Registry FIX disposition and scheduling consequence — 2026-09-15

FO distinctly authorized integrating all six variants into the registered rejection-flow journey, preserving its canonical liveJourney call and every assertion. The implementation moves only the wrapper: fixtures, runner calls and assertions remain byte-identical. The final source delta is +217 net lines across four files, within the captain-approved +220 cap.

Runtime applicability was checked before editing: liveDriverForRuntime selects Claude, Codex or Pi from SPACEDOCK_LIVE_RUNTIME, and the runner chooses the corresponding native lifecycle extractor. No Codex-specific invocation is added to Claude/Pi. The unchanged canonical liveJourney calls parent t.Parallel once and checks its builder/assertion counters before the six synchronous child t.Run calls. The children add no t.Parallel and each owns a fresh fixture and driver.

The six controls now run in normal rejection-flow CI selection on each applicable runtime; local Codex evidence implies approximately 20–25 added serialized minutes, besides roughly nine minutes for the conventional case. Claude/Pi timing is unmeasured. Existing suite limits remain 90 minutes for Claude, 40 for Codex and 50 for Pi; this change consumes more of that budget and does not change those limits. Targeting a child subtest still executes the canonical parent-level drive first. The existing Pi XFAIL remains attached to the canonical call; no new exemptions or local-only classifications were added.

Focused registry reconciliation, timeout, fail-fast and gap-binding contract checks pass; the live-tagged entry point compiles and is discovered as TestLiveCommonRejectionFlow. Prior per-control live evidence is retained because only the wrapper moved. Full normal/race checks are rerunning against the frozen final wrapper.

### CI cost hold and proposed registration surface — 2026-09-15

FO paused the integrated wrapper after the +20–25-minute per-runtime cost was reported. The wrapper move had already occurred under the earlier FIX authorization; no CI, new live drive or code/state commit occurred. The current normal/race lane is only local verification.

Read-only registry inspection found no existing targeted acceptance-proof category: common behavior requires all runtime lanes, host-specific proofs cannot cover common semantics, and non-gating experiments cannot be release evidence. A name-only exemption or experiment relabel would misstate the proof.

Proposed minimal alternative, pending explicit policy/surface authorization: restore the standalone wrapper as TestLiveSameStageRevision; add a Targeted implementation proofs section in docs/runtime-live-ci-registry.md and clarify its distinction from release evidence. This adds a fifth file and 16 net documentation lines, for +236 net lines across five files. Passing controls remain mandatory for this bounded task; no all-runtime or CI coverage is claimed. Promote to the common registry before any release-coverage claim. No unsupported source annotation, runtime machinery, or CI selector change is proposed.

### Held-candidate local verification outcome — 2026-09-15

The local normal/race lane is finished; both exit 1. Registry reconciliation and suite-policy checks pass. Both suites retain the known installed-host resolver failure and reach cumulative ten-minute package timeouts in internal/cli and internal/ensigncycle. The race run additionally fails TestCodexProcessActivityResetsQuietBudget, TestCodexProcessQuietTimeoutPreservesFaultEvidence, TestCodexProcessRecognizesTerminalTurnBeforeOSExit and TestCodexProcessRequiresFinalMessageForTerminalTurn against 250ms no-progress budgets. These new failures are undiagnosed and are not waived; the earlier complete suites had only the resolver and subsequently fixed registry omission. Logs: artifacts/implementation/held-integrated-normal.log and held-integrated-race.log. All processes are stopped. No new live drives or CI were run, and no code/state commits have been made. The held wrapper remains +217 net/four files; the +236/five-file policy alternative is unapplied.

### Binding targeted-proof approval and final implementation — 2026-09-15

Captain approval resolution:binding-1789513304916529000 binds /private/tmp/inflight-implementation-decisions.md at sha256:3e4c7f21b5d53b1abfcc5796d19b7bbe7d90ff2f5c5398415f1e11b245fe9c81. It explicitly approves the targeted implementation-proof registry policy and +236 net lines across five files. FO separately authorized restoring the standalone wrapper and adding the minimal policy documentation, with no CI/all-runtime parity claim. This supersedes the earlier +220/four-file tolerance and the held +217/four-file integrated wrapper; earlier estimates remain history.

The final candidate is exactly +236 net lines across five files: skill net0; live runner+87; fixture/grader helpers+124; standalone registration+9; registry documentation+16. The held expensive common-journey expansion is removed. TestLiveSameStageRevision retains all six controls with unchanged fixtures, runtime selection and assertions; the existing common rejection-flow entry point is restored exactly. Targeted passing evidence is mandatory task acceptance, explicitly selected locally, and does not establish normal-CI or all-runtime parity.

Focused registry/suite-policy checks, stale/unrelated selected-revision falsifiers, worker-completion obligations and live entry-point compilation/discovery pass. The FO reserved the broad verification lane for sequential final normal/race checks. Prior live controls are reused under explicit FO authorization because wrapper relocation changes no fixture, selected runtime or grading behavior. No restacking, code push or CI is authorized or performed.

## Stage Report: implementation

- DONE: Make correction-round recording, independent-review steps, and completion conditions conditional on workflow requirements, preserving authorization and cycle limits.
  Candidate `5e04442a9c89ef644135721522dbf9bfcc8b7de0` changes generic skill applicability only; concrete captain authorization, declared review/round obligations and cycle-3 escalation remain required.
- DONE: Prove same-stage correction with an actual workflow drive and meaningful required-review, missing-evidence, and cycle-limit controls; retain existing Roborev behavior.
  [Proof index](artifacts/implementation/proof-index.md) retains same-launcher before/after: baseline completes correction but stalls on an absent round; candidate commits the frozen-input correction and selects it in one fresh unresolved/unconsumed gate without a reviewer or round.
  Required-review and separate-review controls hold on missing evidence; required-round publishes all four entries without projection, while missing-round refuses publication; weakening a declared obligation makes these controls fail.
  Cycle-limit retained-state regrade proves escalation after three frozen attempts without a fourth; conventional re-review preserves four entries and a fresh gate; existing canonical-round tests preserve Roborev's five entries and unchanged stage.
- DONE: Commit the smallest approved implementation, run required local validation, and report exact candidate and acceptance evidence for independent review.
  Five-file candidate is +236 net (247 insertions/11 deletions), within binding approval `resolution:binding-1789513304916529000`; assigned branch is clean, with no restack, code push or CI.
- DONE: Exercise focused falsifiers and final registry checks.
  Worker tests reject absent/duplicate completion and self-review; selected-plan tests reject premature and unrelated revisions while accepting the exact corrected plan as artifact or canonical context Reference; removing those guards makes the negatives fail.
  Final registry checks reject missing entry registration or suite-policy drift; the targeted entry compiles and is discoverable. It is mandatory local task acceptance, without a normal-CI or all-runtime parity claim.
- DONE: Perform required normal and race validation and report their exact outcomes.
  `go test ./...` and `go test ./... -race` each exit 1, solely for pre-existing `TestCodexResolveManifestAgainstInstalledHost`: spacedock@spacedock is absent while the resolver finds the spacedock-local manifest. All other packages pass; no final timeouts or data-race reports. These suites are not wholly green; see unchanged [normal](artifacts/implementation/final-normal.log) and [race](artifacts/implementation/final-race.log).
  Resolver finding disposition: DECLINED by FO as pre-existing and outside this task’s scope; disposition is resolved by that prior decline. An entirely green installed environment was not an assigned acceptance promise. This reporting correction preserves the failure evidence and does not waive any new regression.
- DONE: Preserve and resolve the earlier verification uncertainty without changing unrelated code.
  Earlier cumulative package timeouts and four 250ms quiet-budget failures are retained; every affected check passed serial isolation, then final reserved-lane suites had no recurrence. Formatting and whitespace checks are clean.
- SKIPPED: Perform independent validation or publish code.
  FO explicitly owns the next independent validation dispatch and later stack delivery; implementation performs neither. Only the path-scoped state report/evidence is published.

### Summary

Workflow-declared same-stage revision now completes without invented reviewer or round machinery, while declared obligations and frozen gate authority remain enforced. The approved targeted-proof registry policy keeps all six controls mandatory for task acceptance without expanding normal CI; the committed candidate is ready for independent validation with the sole known local resolver baseline documented.

## Review-finding disposition

### V1 — Missing retained native worker lifecycle evidence (validation, 2026-09-15)

- Exact observation: public `codex-exec.jsonl` omits native spawn/completion; `rejection-topology.tsv` retains only derived rows. Producer confirms isolated-home cleanup deleted the correlated parent rollouts. Native call-ID pairing, returned task-name identity, and author-attributed `FINAL_ANSWER`/`Done:` cannot be independently replayed.
- Released user and normal workflow: the approved targeted implementation proof requires native worker completion retained for all six same-stage controls; this validation is the intended evidence consumer.
- Observable harm: independently reproduced state outcomes cannot establish native worker ownership/completion or reviewer independence from surviving source bytes; the promised proof cannot be completely audited.
- Affected authority: value-ac[AC-1] live same-stage worker completion must be established; the current retained evidence cannot independently replay that native claim (AC-2 reviewer identity is also limited).
- Trigger evidence: [audit](artifacts/validation/audit.md), seven retained bundle regrades, public after stream's four empty wait results, and producer confirmation that all five original isolated-home roots are empty.
- Reviewer recommendation: **evidence defect / Material**, task-owned observation retention; **FIX proposed, not authorized**. No workflow outcome defect is alleged. Existing `nativeLifecycleStream` can retain already-read source bytes before cleanup; estimated one retention line plus 15–25 focused test lines exceeds the current exact +236 cap, so captain tolerance decision is required before repair.
- Proposed bounded replacement: six mandated candidate variants plus causal baseline; include the conventional re-review run for independently replayable AC-2 attribution (eight total). Retain public/native stream, exact task/call identity, attributed completion, bundle, and topology; no new lifecycle framework or broad/all-runtime rerun is proposed.
- FO disposition: pending; candidate remains unchanged. No automatic correction cycle or reviewer rerun is authorized by this report.

## Stage Report: validation

- DONE: Independently assess all three acceptance criteria against committed candidate and retained causal/native control evidence; do not duplicate owned green runs.
  AC-1: bundle regrade confirms exact committed correction, immutable frozen input/old gate, exactly one fresh open gate; native completion provenance remains unmet under V1. AC-1 is not fully validated.
  AC-2: required-review/missing-round hold state, canonical four-entry no-projection round, cycle-3 limit, and existing five-entry Roborev/lifecycle tests pass; native independent-review attribution remains limited by V1.
  AC-3: five-file/+236 committed diff changes only generic applicability plus test/registry proof; before/after state supports the causal outcome, with no new flags, persisted fields, recorder policy, or synthetic stage. Latest approved targeted-proof policy is reflected.
- DONE: Perform the required detached adversarial audit of skill authority and gate-selected artifact/reviewer/round controls, preserving exact findings and scope.
  [Audit](artifacts/validation/audit.md) and [temporary test](artifacts/validation/detached_validation_audit_test.go) preserve seven bundle regrades, exact selection identity/revision/type/EOF and worker order/ownership negatives; [results](artifacts/validation/detached-audit.log) pass. Weakening exact source or completion checks makes negatives fail.
  [Canonical round controls](artifacts/validation/canonical-round-audit.log) refuse incomplete/dangling logs without publication and preserve complete Roborev entries/status; accepting incomplete rounds makes them fail. V1 records missing native source evidence separately from state outcomes.
- DONE: Report PASSED or REJECTED with per-AC evidence and resolver baseline limitation; do not approve gates or push code/CI.
  **REJECTED — evidence defect V1**; no observed product outcome defect. FO/captain disposition is required before bounded repair because current approved tolerance is exact.
  Retained final `go test ./...` and `go test ./... -race` each exit 1 solely on the pre-existing installed-host resolver failure; prior FO decline remains valid, and neither suite is claimed wholly green. No expensive/full/live runs duplicated; no code mutation, gate approval or CI/code push.
- SKIPPED: Candidate repairs and replacement live proof.
  Validator has observation/recommendation authority only; minimal repair and bounded rerun proposal are recorded in V1. No deferred-risk or polish finding is substituted for the material retention gap.

### Summary

The candidate's durable behavior and detached adversarial state checks pass, but required native worker evidence was deleted after grading and cannot be independently replayed. Validation recommends REJECTED for that evidence defect, with a small existing-observer retention repair and bounded evidence replacement proposed; candidate bytes remain unchanged.

### V1 correction authorization and bounded implementation — 2026-09-16

Captain resolution:binding-1789517330520748000 (artifact sha256:32821e2071d9c315c7a6dfc09ba4c80c14e034bb07bb75958c136abd3a35bf30) approves final +252 net lines/five files and eight serialized native evidence replacements. FO separately authorized FIX. The proposal adds no lifecycle controller: existing correlated stream loading persists exact source bytes before cleanup and fails on write failure. Focused red was observed before persistence; byte replay after cleanup and existing parser negatives now pass. The earlier candidate and evidence remain history. Canonical advisory round sources are artifacts/implementation-cycle-1/briefing.json and briefing.review.jsonl; producer completion will be appended after all approved proof is complete. FO owns projection/publication and the retained validator owns independent review.

## Stage Report: implementation (cycle 1)

- DONE: Persist exact correlated lifecycle bytes and prove retention/write-failure/parser controls within252net/5files.
  Candidate `d46f74245c1c18c8a376e9a6d2a508c774e32759` retains the existing correlated public/native stream before cleanup, refuses write failure and retains conventional matching Git state; final task delta is +252net/five files under `resolution:binding-1789517330520748000`.
  Focused red first exposed missing persistence failure handling; green replay deletes the source home, requires exact byte equality and exercises existing missing-completion/wrong-owner negatives. Omitting persistence, accepting write failure or misattributing completion makes these controls fail.
- DONE: Replace the eight approved native evidence runs with matching retained raw/state and exact regrade, preserving old evidence.
  [Correction proof index](artifacts/implementation-cycle-1/proof-index.md) retains expected baseline red, six candidate PASS variants and conventional PASS. Same launcher/fixture/harness causal pair differs only in skill bytes; all eight source streams reproduce exact topology and correlate public parent identity.
  [Native/state manifest](artifacts/implementation-cycle-1/native-state-regrade.json) pins every raw/public/bundle digest and state HEAD; bundle regrades preserve corrected frozen-input plans, intended selected revisions, exact gate/round cardinalities and conventional reviewer independence/reuse. Old evidence is unchanged.
  Exactly8unique drives, no retries/duplicates:1842.17s summed test runtime versus4556.9645s wall; [timing](artifacts/implementation-cycle-1/timing.json) localizes2541.502s unexplained wall/package discrepancy to round-missing. Existing timeout limits were preserved.
- DONE: Commit correction, complete canonical correction round producer entries, and report for retained independent validator.
  Code committed as `d46f74245c1c18c8a376e9a6d2a508c774e32759`; canonical sources are [Briefing](artifacts/implementation-cycle-1/briefing.json) and [review log](artifacts/implementation-cycle-1/briefing.review.jsonl), including reviewer V1, producer proposal, distinct FO authorizations and producer completion/closing advisory Resolution.
- SKIPPED: Repeat broad suites, publish the neutral round, or independently validate the correction.
  FO explicitly limits replacement to focused retention checks and eight native drives, owns projection/publication, and will dispatch the retained validator. Prior final normal/race logs each exit1 only for the FO-declined out-of-scope resolver baseline; no wholly green-suite claim is made.

### Summary

V1 is repaired at the existing retention seam without a controller or workflow-semantic change. All approved replacement drives now preserve source-level native attribution and matching state; the committed correction and completed producer round are ready for independent validation, with original failed evidence retained unchanged. No code push, CI, restacking, gate approval or frontmatter mutation occurred.

### V1 independent correction recheck — 2026-09-16

V1 is **satisfied** by candidate `d46f74245c1c18c8a376e9a6d2a508c774e32759` after captain amendment `resolution:binding-1789517330520748000` and distinct FO FIX. The existing observer now retains correlated source bytes before cleanup, fails on write failure, and all eight approved replacement drives retain replayable native attribution and matching state. [Independent detached audit](artifacts/validation-cycle-1/audit.md) reproduces every topology row and rejects deleted completion events; no new finding, candidate mutation or second round publication occurred.

## Stage Report: validation (cycle 1)

- DONE: Independently assess all three acceptance criteria against committed candidate and retained causal/native control evidence; do not duplicate owned green runs.
  **AC-1 PASSED:** baseline/after raw native completion replays against corresponding bundles; exact committed frozen-input correction moves from one rejected attempt to one fresh open attempt without reviewer/round, resolution, application or successor. Gate-selected Git revision/digest resolves the corrected plan.
  **AC-2 PASSED:** both required-review controls hold with missing evidence; distinct correction/reviewer owners replay where review runs. Required-round publishes four canonical entries without projection; missing-round and cycle-3 hold. Conventional eight-event re-review and four-entry round regrade pass; prior Roborev five-entry/status controls remain valid.
  **AC-3 PASSED:** same launcher/fixture causal pair differs in skill bytes and observed re-gating outcome; retained CLI stale-authority/incomplete-round checks remain valid. Final approved +252/five-file surface adds evidence retention to generic routing applicability, without new flags/state schema/recorder policy/synthetic reviewer.
- DONE: Perform the required detached adversarial audit of skill authority and gate-selected artifact/reviewer/round controls, preserving exact findings and scope.
  [Detached replay](artifacts/validation-cycle-1/detached-replay.log) passes eight raw/public parent-ID correlations and exact TSV route replay, seven selected-state regrades, conventional round/gate checks, source-selection and worker-order negatives. Deleting native completion records invalidates all eight traces; stale/wrong selected bytes fail.
  Prior authority/canonical-round audit is reused; retained focused proof exercises write failure and exact byte replay after deleting source home. V1 is closed by evidence, with no new material/deferred-risk/polish finding.
- DONE: Report PASSED or REJECTED with per-AC evidence and resolver baseline limitation; do not approve gates or push code/CI.
  **PASSED.** Independent native provenance now meets the previously unmet requirement. Final prior normal/race each exit1 only for the FO-declined installed-host resolver baseline; neither is described as wholly green. No new full/native runs, candidate edits, gate approval or code/CI push.
  Timing retained truthfully: eight unique drives/no retries, 30m42.17s summed cases versus 75m56.96s wall; 2541.502s wall/package discrepancy in round-missing remains unexplained. Local targeted proof establishes no all-runtime/CI parity.
- SKIPPED: Publish another correction round or repeat owned green suites.
  FO already recorded validation/1 with exactly ten entries; reviewer recheck does not republish it. Only local detached artifact replay ran.

### Summary

All three acceptance criteria now have replayable behavioral evidence, and V1 is satisfied within the captain-approved surface. Validation recommends PASSED with the existing resolver baseline and unexplained timing discrepancy explicitly preserved; candidate and gate authority remain unchanged.


## Captain correction: routine live coverage — 2026-09-16

The captain directed: "no live test should be targeted implementation proof only," then asked how omission could require an explicit choice. FO authorized common-journey registration for all six variants and coverage enforcement in the existing reconciliation owner. This supersedes the targeted-proof exemption and release-acceptance statements in earlier reports. Those reports and logs remain historical records, not the current policy.

`TestLiveCommonSameStageRevision` now binds six fixture IDs to the real builder and worker assertion. Its existing driver runs plain, review-required, separate-review-required, round-required, round-missing, and cycle-limit serially. Each has a separate workflow and artifact path. The canonical conventional rejection journey remains unchanged. Claude and Codex select the journey once through the existing scheduler; Pi selects it once through the common-suite selector. No skip/XFAIL is added. Existing scheduler parallel suppression applies, so serial child tests introduce no nested parallel completion boundary.

The scheduler hint is 1,400 seconds for the six Codex variants, rounded from retained timings, and a conservative 1,800-second Claude estimate. These are ordering hints, not new timeouts or measured Claude evidence. Routine CI now incurs these six native drives per supported runtime target. There are no local native/model reruns in this correction; final runtime outcomes belong to updated stack-tip CI.

The registry no longer admits Targeted implementation proofs. Reconciliation now requires each live declaration to have an actual routine lane selection, including actual scheduler membership, or an explicit Non-gating live experiments entry with a nonblank reason. Missing selector, missing membership, duplicate selection, unclassified/targeted category and blank reason fixtures exercise the same checker. The actual repository first failed on its old standalone same-stage entry, then passed after promotion. The three existing experiments are unchanged. Upper-layer #802 promotion remains its owner's work: Codex semantic handles must become a runtime proof, and deterministic stamped identity belongs offline. This branch does not claim those upper edits are already incorporated.

All moved scenario grading statements compare identically after whitespace normalization and the counted assertion-parameter rename. Frozen attempts, corrected selected revision, unrelated/premature selection negatives, required round cardinality, native byte retention and write-failure/parser checks remain intact. Retained local native evidence is preserved without a new all-runtime claim.

Verification completed: focused registry and live-tagged deterministic checks pass. Required normal/race suites each exit 1 solely for the known installed-host resolver baseline; all other packages pass. Logs and exact command timestamps are under artifacts/implementation-ci-registration. No code push, CI trigger or workflow-frontmatter/gate mutation occurred.


## Stage Report: implementation

- DONE: Register every same-stage variant in routine supported-runtime CI and remove the targeted-only exception.
  Code commit `bf64ebeca949667c2d60e40a9c595ebdf62375b2` registers all six fixtures through `TestLiveCommonSameStageRevision`, one scheduler row, and the existing runtime adapters. Original assertions and unique falsifiers remain intact. The three explicit experiments are unchanged; #802 owns its separate promotions after restack.
- DONE: Require actual lane coverage or explicit exemption through the existing reconciliation checker.
  Real pre-promotion reconciliation failed on the old targeted entry. Focused checks now pass and reject missing selectors, missing scheduler membership, duplicate selection, unknown/targeted classification, and blank exemption reasons. Runtime-specific obligations come from the registry Lane and actual conditional scheduler rows, without a hardcoded test-name exception.
- DONE: Perform required local verification, retain honest limitations, and commit the correction for independent review.
  `go test ./...` and `go test ./... -race` each exit 1 only at `TestCodexResolveManifestAgainstInstalledHost`: the named marketplace plugin is absent but resolves to installed spacedock-local/pre0. This is the existing FO-declined out-of-scope baseline, disposition resolved; neither suite is claimed wholly green. `internal/ensigncycle` passes in 356.523s and 356.418s respectively. Focused registry and live-tagged deterministic checks pass. `gofmt -w ./cmd ./internal` and diff-check completed; unrelated existing formatting was preserved.

### Summary

Correction: 256 insertions / 115 deletions = +141 net across six modified files. Full task layer above `4ce49f1ea`: 416 insertions / 23 deletions = +393 net across seven files, within the FO-authorized +100–150 correction estimate and seven-file scope. This is a test-registration/enforcement correction; production code, skill behavior and workflow authority are unchanged.

Evidence: `artifacts/implementation-ci-registration/red-proof.txt`, `focused-registry.log`, `focused-runtime.log`, `normal.log`, `race.log`, `checks.json`, and `scope-and-preservation.txt`. Earlier native/state evidence remains unchanged. No local native/model run establishes this expanded runtime coverage; updated stack-tip CI is required for final live acceptance. The retained independent validator owns review. FO owns publication/restacking; no code/state push, CI trigger, gate decision or independent self-validation occurred.

## Review-finding disposition — routine CI correction

### V2a — Scheduler membership checker accepts an omitted runtime journey

- Exact evidence: [wrong-runtime-filter.log](artifacts/validation-ci-registration/wrong-runtime-filter.log) exits0 after filtering TestLiveCommonSameStageRevision from Codex's row-to-jobs loop. Actual scheduler execution with logging-only callable stubs changes from one invocation to zero; [baseline](artifacts/validation-ci-registration/offline-schedule-baseline.log), [omission](artifacts/validation-ci-registration/offline-schedule-codex-omission.log). No native host was launched.
- Released user and normal workflow: maintainers rely on reconciliation to reject routine-lane omissions unless explicitly exempted; runtime filtering is an ordinary scheduler change.
- Observable harm: declared row literals satisfy the guard while executable membership omits all six variants for Codex, without an exemption.
- Affected authority: captain-ruling[2026-09-16] omission from routine live coverage must require an explicit choice; actual selection must be enforced.
- Trigger: four-field AST rows are counted as both-host membership independently of the executed row-to-jobs filter.
- Proposal: **Material / evidence-enforcement defect / task-owned / FIX recommended**. Observe actual scheduler callback selection through existing scheduler/checker owners with non-native stubs; do not add a general Go interpreter. FO disposition required before candidate repair.

### V2b — Commented live command still satisfies coverage

- Exact evidence: [commented-selector.log](artifacts/validation-ci-registration/commented-selector.log) exits0 after prefixing the actual Claude gotestsum scheduler command with `#`; [reproduction](artifacts/validation-ci-registration/checker-audit.py).
- Released user and normal workflow: maintainers can comment out a routine workflow command while editing CI; coverage reconciliation is supposed to detect omission.
- Observable harm: no shell invocation remains for that command, but the lane is reported selected without exemption.
- Affected authority: captain-ruling[2026-09-16] missing routine live execution requires an explicit exemption rather than dead selector text.
- Trigger: liveTestCoverageError scans raw YAML lines and counts selector text in shell comments as execution.
- Proposal: **Material / evidence-enforcement defect / task-owned / FIX recommended**. Reuse existing yaml.v3/mappingValue workflow parsing and ignore shell-comment lines in jobs.steps.run. FO disposition required before candidate repair.

## Stage Report: validation (cycle 2)

- DONE: Verify all six same-stage variants are selected by routine lanes exactly once with preserved assertions and no targeted-only exception.
  Candidate bf64ebeca has one Claude/Codex scheduler row, Pi common selector, six serial variants and no new skip/XFAIL. Moved grading is semantically identical and conventional rejection body byte-identical; three named exemptions remain unchanged. Existing native/state proof remains valid, but updated required-host tip CI is pending.
- DONE: Attack the existing coverage checker with missing selectors, wrong runtime membership, unclassified tests and invalid exemptions; require actual failures.
  [Detached mutation matrix](artifacts/validation-ci-registration/audit.md) confirms missing selector/row, unclassified targeted category and missing exemption reason fail. Wrong runtime filtering and a commented command incorrectly pass; V2a/V2b preserve exact actual failures of enforcement.
- DONE: Assess local checks and scope against the captain amendment, preserving explicit pending live tip CI and any findings.
  **REJECTED** for V2a/V2b. Proposed bounded existing-owner corrections require distinct FO disposition; candidate remains unchanged. Earlier PASSED reports do not supersede this new captain-required enforcement scope.
  Required actual normal/race results each exit1 solely on known FO-declined installed-host resolver; neither is wholly green. +141 correction/six files and +393 full task/seven files match producer scope. No duplicate broad/native/model runs, upper-layer edits, code push, PR or CI.
- SKIPPED: Claim final live acceptance or repair candidate.
  New stack-tip CI must execute the promoted journey; old tip and retained local Codex evidence cannot establish expanded runtime coverage. Validator records findings only.

### Summary

Routine registration preserves all six behavioral controls and the conventional journey, but its new omission checker has two demonstrated false greens. Validation recommends REJECTED until actual runtime omission and commented-command controls fail as required, with new tip live CI still pending.

### V2a/V2b FO disposition — 2026-09-16

FO separately authorized FIX: the producer will observe actual callable identity/count by executing exact scheduler and ordering source offline with logging stubs, and parse existing YAML steps.run while discarding comments. Scope is one checker file plus report, with no models or new framework. Validator preserves the above findings and waits for producer completion before correction recheck; no self-repair is authorized.


## Stage Report: implementation

- DONE: Repair V2a with executed scheduler selection, not inferred row membership.
  Commit `342edfa8903e0e822ec545d2ee45401352eb6c60` runs the exact current scheduler and ordering source in a temporary, self-contained Go test, replacing only discovered live callable bodies with logging stubs. The actual rows-to-jobs filtering, slot loop and dispatch produce callable identity/count evidence. Reconciliation runs this once per Claude/Codex runtime, then reuses that inventory for all entries. Each subprocess has a one-minute timeout. No repository runtime adapters, auth setup or model drivers enter the temporary program.
- DONE: Repair V2b and retain both demonstrated mutation controls.
  Workflow selectors now come from yaml.v3-parsed jobs.steps.run scalars; shell-comment lines do not count. The retained runtime-continue mutation executes and observes Claude=1/Codex=0 for same-stage revision, which coverage rejects. The actual Claude command is then commented in the parsed workflow and the same unaffected filing inventory proves that coverage rejects it. Original missing-selector, membership, duplicate, classification and exemption controls remain. Validator red evidence is retained under artifacts/validation-ci-registration; the comment case was also reproduced locally before the remedy.
- DONE: Verify the changed package and commit the bounded correction for retained independent review.
  `go test ./internal/contractlint -count=1` passes in 2.632s; `go test ./internal/contractlint -race -count=1` passes in 8.696s. Focused checks pass in 2.447s. Current baseline reconciliation takes 1.16s, including offline scheduler calls of 0.556s and 0.401s; the two-mutation control takes 0.97s with 0.552s and 0.399s probes. Metadata fixtures reuse inventory and do not launch probes individually. Changed file is gofmt-clean and diff-check passes. No local model/native workflow, full-suite rerun, CI trigger, push, rebase or authority mutation occurred.

### Summary

Both task-owned coverage evidence defects were fixed under distinct FO authorization using the existing checker. Final correction is 133 insertions / 89 deletions = +44 net in one existing file; full task layer is +437 net across seven files. There is no new maintained harness file, scheduler framework, AST interpreter or coverage exception. The evidence boundary is executable scheduler dispatch plus declared YAML run commands with comment exclusion; it does not claim to interpret arbitrary shell programs.

Logs and exact command/duration metadata are under artifacts/implementation-ci-selection. Prior full normal/race results belong to candidate `bf64ebeca`, each with only the known FO-declined resolver baseline; they are not new-head full-suite results. FO explicitly authorized changed-package normal/race verification now and will run the required full commands once on the final combined local tip before final delivery. Final live evidence still belongs to updated stack-tip CI. Retained independent validation is pending; this producer does not validate its own work or advance state.

### V2a/V2b independent correction recheck — 2026-09-16

Both findings are **satisfied** at `342edfa8903e0e822ec545d2ee45401352eb6c60` under the recorded distinct FO FIX. The unchanged detached mutation driver now gets actual reconciliation failure for Codex runtime omission and the commented Claude command; baseline passes and all prior negative controls still fail. [Recheck evidence](artifacts/validation-ci-selection/audit.md) preserves exact source, outputs and limits. No new finding or candidate change occurred.

## Stage Report: validation (cycle 3)

- DONE: Verify all six same-stage variants are selected by routine lanes exactly once with preserved assertions and no targeted-only exception.
  Candidate `342edfa8903e0e822ec545d2ee45401352eb6c60` changes only the checker; prior six-variant assertion/fixture preservation and conventional byte-equivalence evidence remain valid. Actual scheduler callback counts establish one selection per Claude/Codex lane; parsed Pi common selector covers the common entry. Three explicit experiments remain unchanged; expanded host execution still needs new tip CI.
- DONE: Attack the existing coverage checker with missing selectors, wrong runtime membership, unclassified tests and invalid exemptions; require actual failures.
  [Exact unchanged mutation matrix](artifacts/validation-ci-selection/results.json) now passes expectations: baseline exit0; missing selector/row, runtime filter, commented command, unclassified entry and missing exemption reason all exit1. V2a reports same-stage selected0 in Codex; V2b reports selected0 in Claude. Existing duplicate/explicit-exemption controls remain intact.
  Counts come from executed exact scheduler/ordering code invoking declaration-derived callback stubs, not synthetic registry membership; no runtime/model driver enters the probe. YAML jobs.steps.run parsing excludes the demonstrated shell-comment command. V2a/V2b are satisfied with no new finding.
- DONE: Assess local checks and scope against the captain amendment, preserving explicit pending live tip CI and any findings.
  **PASSED for this bounded correction.** Changed-package actual normal/race checks pass; reused prior full normal/race belong to bf64ebeca and retain only the FO-declined installed-manifest resolver baseline. Required final combined-tip full suites and new tip host CI remain pending with FO; neither old CI nor local native evidence proves promoted runtime coverage.
  Correction is +44net in one existing checker file, without a new framework or coverage exception. No candidate/frontmatter edits, new model/native drives, duplicate suites, code push or CI; prior authority and finding history preserved.

### Summary

The checker now rejects both demonstrated omission paths using actual scheduler callback execution and parsed workflow commands. Validation recommends PASSED for the correction while explicitly leaving final combined-tip checks and promoted host CI outstanding; no new material finding remains.


## Stage Report: implementation

- DONE: Correct V3a/V3b at the native identity and completed-turn grading boundary.
  Candidate `67c68536fd6a49461624f1c6d86081d22bef4811` correlates Claude Agent tool IDs with returned native task IDs using the existing native metadata parser. Names remain routing/display labels, not independent-worker identity. Completion must match the opened tool and, when available, native task owner. The grader preserves ordered dispatch/completion pairing and allows completed reuse only of the current worker; fresh duplicate identities or extra workers fail. The caller infers reviewer presence from fresh spawns rather than completed turns. Focused controls first failed on both captured shapes and label-based self-review, then passed after correction. Full retained streams regrade at the actual grader: Claude workers at94/166 and274/331 are distinct; Codex spawn112/completion136 and followup151/completion163 share one native worker. Missing/wrong-owner completions, same native identity under changed labels or tool IDs, incomplete followup, self-review, and extra fresh workers remain rejected. Existing invalid selected-gate revision controls remain unchanged and pass.
- DONE: Repair adjacent AC-2 publication evidence using actual recorder argv/exit observations.
  Under separate FO FIX authorization, conventional Claude now uses the same existing observer shim and neutral-named recorder publication parser as Codex. Pi remains on its existing path. A focused test first reproduced the false count of two publications, then proves the actual launcher pins its binary and the existing shell startup shim restores interception inside a non-model Claude stand-in. An incomplete round exits1 under a successful shell wrapper; the corrected round exits0 under another successful shell wrapper; only one successful publication counts. Existing echo-only/zero, missing-round-room, duplicate real success and second-round success negatives remain rejected. Absent recorder observations fail closed. No model-authored exit echo or shell success flag supplies recorder truth. The premature first recorder call remains observed conduct; this fix does not claim perfect model compliance.
- DONE: Complete applicable local verification, retain prior evidence and report for independent review.
  Focused live-tagged normal checks pass7.003s, race checks pass8.092s, and registry reconciliation passes1.005s. Checks cover native identity/reuse, reviewer topology, same-stage gate selection, actual recorder execution, and shell startup interception. Full captured native streams pass the actual extractor/grader in a temporary replay test; source hashes and output are retained. Gofmt and diff-check pass; unrelated existing formatting was preserved. No local model/native workflow, broad full-suite rerun, code/state push, CI trigger, rebase, gate decision or frontmatter mutation occurred. FO directed full normal/race verification once on the final combined tip and owns the final live CI run.

### Summary

V3a/V3b and the separately authorized publication-count evidence defect are corrected in existing owners. This correction is161 additions/27 deletions = +134 net across five existing files: +93 for native identity/reuse and controls, +41 for actual recorder interception/publication proof. FO explicitly accepted +41 against the approximate20–40 estimate. Cumulative layer above4ce49f1ea is +571 net across10 files; the three additional existing evidence owners are shared_reviewer_reuse_test.go, rejection_round_execution_test.go and shared_round_recording_test.go. No skill, CLI, fixture-outcome, workflow policy, controller or logging framework changed.

Evidence is in artifacts/implementation-native-identity: captured-regrade.log, captured-sources.json, red-proof.txt, normal.log, race.log, registry.log and checks.json. Prior failures and raw/state bundles are unchanged. Earlier full normal/race results apply tobf64, not this new head; the known resolver baseline remains explicitly limited to those earlier runs. Current local results establish the deterministic correction only. Final combined-tip broad verification and live CI remain pending FO coordination, and the retained independent validator owns the review recommendation.

## Stage Report: validation (cycle 4)

- DONE: Independently verify native identity and reuse corrections against captured failures and strict wrong-owner/self-review/extra-worker negatives.
  Candidate `67c68536fd6a49461624f1c6d86081d22bef4811` satisfies V3a/V3b: retained full Claude/Codex replay maps two distinct native workers and one completed reused worker correctly. [Independent captured-stream falsifiers](artifacts/validation-native-identity/falsifiers.log) reject native-ID self-review, missing reviewer completion and an unspawned followup owner; owned wrong-owner/extra-worker/gate-selection negatives remain valid.
- DONE: Verify recorder success uses actual intercepted argv/exit and rejects shell-success masking, absent evidence, and duplicate true successes.
  Actual launcher/Claude-stand-in interception proof is retained in implementation-native-identity; the runner uses the same shell startup shim. Failed recorder exit1 and successful recorder exit0 under successful enclosing shells count one publication. Independent absent/masked/wrong-command/wrong-entity/duplicate observation controls reject; failed-then-one-success passes. Missing-room and duplicate-real-success execution negatives remain intact.
- DONE: Assess unchanged outcome/authority checks, exact scope and current local evidence; report recommendation with combined full suites and final tip CI pending.
  **PASSED for the deterministic correction; no new material finding.** [Review](artifacts/validation-native-identity/review.md) confirms unchanged durable/gate/review/cycle outcome guards and no skill/CLI/policy change. Exact correction+134net/five files and cumulative+571net/ten match approved scope.
  Independent focused falsifiers pass0.317s; reuse owned normal7.003s/race8.092s and registry1.005s without duplicate suites. Final combined broad normal/race and final host CI remain FO-owned pending Pi correction. Earlier broad logs, including known resolver failure, and old tip results do not establish new-head full/runtime acceptance. No model/native workflows, candidate/frontmatter edits, pushes or CI.

### Summary

The corrected evidence layer distinguishes native workers from display names and completed turns, and counts recorder success from the actual intercepted process. Validation recommends PASSED for this bounded correction while final combined suites and host CI remain pending; prior findings and failed-run evidence are preserved.
