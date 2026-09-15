---
title: Honor workflow-declared same-stage revision without mandatory reviewer machinery
status: implementation
source: Captain request after email-triage FO issue 792
issue: spacedock-dev/spacedock#792
score: 0.95
started: 2026-09-15T04:30:03Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-workflow-owned-same-stage-revision
pr:
mod-block:
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

Current approved tolerance: net +236 LOC maximum and 5 files, bound by resolution:binding-1789513304916529000; the added owner is `docs/runtime-live-ci-registry.md` for explicit targeted-proof policy.

Expected owners: `skills/feedback-rejection-flow/SKILL.md`, `internal/ensigncycle/claude_live_runner_test.go`, and `internal/ensigncycle/claude_runtime_helpers_test.go`.
The existing `runClaudeRejectionFlowScenario` owns the rejection run; `shared_live_runner_test.go` registers it.
The topology grader lives in `claude_runtime_helpers_test.go`. Registration changes, if needed, fit the existing four-file tolerance.
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
