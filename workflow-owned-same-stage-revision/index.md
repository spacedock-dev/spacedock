---
title: Honor workflow-declared same-stage revision without mandatory reviewer machinery
status: ideation
source: Captain request after email-triage FO issue 792
issue: spacedock-dev/spacedock#792
score: 0.95
started: 2026-09-15T04:30:03Z
completed:
verdict:
worktree:
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
Tolerance: net +90 LOC maximum and 4 files. Larger scope requires a design reset.

Expected owners: `skills/feedback-rejection-flow/SKILL.md`, the existing shared live rejection fixture/test owner, and its existing result grader.
Use `internal/ensigncycle/shared_promoted_live_test.go` and the discovered rejection fixture owner after reading their current boundaries.
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

The existing shared rejection live runner owns agent routing and completion proof. Extend its workflow inputs and durable-state grader, rather than creating another runner.
Before implementation, add the smallest failing self-feedback fixture. Run current and candidate skills through the same serialized local runtime lane.
A valid paraphrase must pass. Restoring unconditional recorder/reviewer requirements must fail the self-feedback case.
Removing the required-review condition must fail the separate-reviewer and same-stage-required-review controls.
Conflating projection absence with round absence must fail a required-round/no-projection control.
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
