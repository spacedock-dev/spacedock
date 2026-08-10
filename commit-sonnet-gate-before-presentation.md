---
title: Commit Sonnet gate before presentation
status: validation
score: "0.90"
source: "n28 exact Claude default-headless finding, 2026-08-10"
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-evidence
id: kky8pg7wc8xgb985epwss092
gates:
    version: 1
    records:
        - id: gate:kky8pg7wc8xgb985epwss092:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:kky8pg7wc8xgb985epwss092-backlog-1
              briefing:
                id: briefing:kky8pg7wc8xgb985epwss092:backlog:attempt-1:revision-1
                digest: sha256:a45f5479086379c8c2fa589df282f75f9fd2e98afa472f24c553480dae7f0398
                request-digest: sha256:2315b198901a80e371d222f7a35ec38f11019774d5b8ef5237989e6c0fe5cf8d
                room-ref: ./commit-sonnet-gate-before-presentation/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:kky8pg7wc8xgb985epwss092:backlog:1
                briefing: briefing:kky8pg7wc8xgb985epwss092:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T18:44:24.797313Z"
                decision: approve
                reason: Captain created this active Sonnet owner to repair the target-external gate lifecycle defect.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:kky8pg7wc8xgb985epwss092:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:kky8pg7wc8xgb985epwss092-ideation-1
              briefing:
                id: briefing:kky8pg7wc8xgb985epwss092:ideation:attempt-1:revision-1
                digest: sha256:8e843810ce2f0f10b39cd9689c4dbc41b323f372695c969fbc7374a5a903d2da
                request-digest: sha256:21635b1a8292b76f2e135f53131b93f11c7969f8afbd02cf4cf9581023c6fd50
                room-ref: ./commit-sonnet-gate-before-presentation/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:kky8pg7wc8xgb985epwss092:ideation:1
                briefing: briefing:kky8pg7wc8xgb985epwss092:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T18:49:10.674133Z"
                decision: approve
                reason: Captain assigned this exact Sonnet gate-lifecycle repair. The one-file design preserves the approved boundary.
              application:
                target-stage: implementation
                state: consumed
started: 2026-08-10T18:44:40Z
worktree: .worktrees/spacedock-ensign-commit-sonnet-gate-before-presentation
---
## Problem

The exact local Sonnet default-headless journey prepared an open validation gate.
The host did not commit the gate state.
The First Officer then read and presented the gate.
The target reported `gate-hold-violation`.

## Value

After `gate prepare`, the Sonnet First Officer commits the durable gate state. It stops at the same clean open validation gate and dispatches no successor.

## Scope

- Repair the host-neutral gate lifecycle at the prepare-and-bind boundary.
- Use the exact n28 Claude artifact as the baseline.
- Remove the Sonnet `default-headless-gate-stop` binding after bound XPASS evidence.
- Transfer the Codex binding to `select-actionable-codex-default-headless-task`.
- Do not change n28 acknowledgment mechanics, Pi, or the zero-discovery repair.
- Use local Sonnet subscription authentication before required PR CI.

## Acceptance criteria

- AC-1: The command order is successful `gate prepare`, then `state commit`, then structured reads and presentation.
- AC-2: The exact local Sonnet default-headless target first reports bound XPASS-green and then passes normally after binding removal.
- AC-3: The final entity is a clean open validation gate with no terminal fields and no successor dispatch.
- AC-4: A focused negative control rejects a missing, failed, or late state commit.
- AC-5: Full, race, format, registry, active-owner, and required exact PR checks pass. Pi remains skipped.

## Baseline evidence

- Released user and workflow: local-subscription Sonnet default-headless gate stop without Captain authority.
- Observable harm: the prepared gate is not proven durable before presentation.
- Value authority: `skills/fo-gate-lifecycle/SKILL.md` requires prepare, commit, structured reads, and presentation in that order.
- Sonnet trigger: `/tmp/n284-happy-claude/claude-shared-scenarios/default-headless-gate-stop/command.log` has gate prepare, reads, and presentation, but no state commit.

## Captain-approved scope disposition

On 2026-08-10, the Captain expanded this task from Sonnet to the host-neutral Sonnet and Codex gate commit boundary.
Later Codex evidence did not reach the gate-commit behavior.
The Captain recarved accepted value to Sonnet and assigned Codex selection to `272j6s25f9mry6nxbf4yjxvt`.
The host-neutral product bytes, one-file design, and eight-line hard limit remain unchanged.

The first local Codex run collided with a concurrent Sonnet dispatch artifact.
The authorized sequential Codex run selected the queued task as a gate and stopped before dispatch.
Its artifact remains at `/tmp/kky-bound-codex-sequential.CHE1eq`.

## Ideation requirements

- Name exact files and gross/net estimate before product edits.
- Preserve the existing gate grammar and authority.
- Define one focused falsifier for the missing commit.

## Proposed approach

Change only the open-gate resume rule in `skills/fo-gate-lifecycle/SKILL.md`.
The rule will require `state commit ENTITY` before the structured reads and presentation.
This rule also applies when a worker prepared or committed the open gate.
The First Officer will stop if the commit command fails.

The exact design replaces the `open → present` clause near line 74.
The replacement names the commit, the two structured reads, and the presentation order.
No other lifecycle branch changes.

The simplest alternative was a new binary command that combines commit and presentation.
That alternative is too large because `state commit` already supports a clean no-op.
The retained logger also observes its successful state head.

No spike is needed.
The clean no-op behavior exists in `TestStateCommitNoOpWhenClean`.
The command-order oracle exists in `assertRecordedGateHoldLog`.

## Expected surface

- `skills/fo-gate-lifecycle/SKILL.md`: 3 insertions and 1 deletion.
- Expected total: 1 file, 4 gross lines, and 2 net lines.
- Hard tolerance: 1 file and 8 gross lines.

The change alters host-neutral Sonnet and Codex runtime instruction behavior only.
It does not alter command grammar, stored formats, or authority.
It does not alter n28 dispatch acknowledgment, Pi, or zero-discovery behavior.

## Acceptance criteria and test plan

**AC-1: A successful open-gate presentation has a successful `state commit` after `gate prepare` and before structured reads.**
The exact local Sonnet journey verifies the command order in `command.log`.

**AC-2: The exact local Sonnet target first reports XPASS-green and then passes after its binding is removed.**
The live journey verifies the repaired behavior and the final open gate.

**AC-3: The final gate remains open without a Resolution, an Application, terminal fields, or successor dispatch.**
The existing `assertGateHeld` and `assertRecordedGateHoldLog` checks verify this state.

**AC-4: A missing, failed, or late commit does not qualify the gate hold.**
The focused `missing commit` case in `TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle` is the primary falsifier.
The same oracle rejects a commit before the successful prepare.

**AC-5: The product checks and exact required PR checks pass.**
Implementation runs the focused lifecycle tests, full tests, race tests, formatting, registry checks, and active-owner checks.
The required PR lanes run for Sonnet and Codex.
Pi remains skipped.

## Stage Report: ideation

- DONE: Read the complete task and exact n28 Claude baseline.
  The log shows successful `gate prepare`, structured reads, and presentation without a binary `state commit` after prepare.
- DONE: Inspected the prepare-and-bind contract and the Sonnet instruction path.
  The cold prepare sequence requires the commit, but the open resume clause routes directly to presentation.
- DONE: Defined the smallest repair and a missing-commit falsifier.
  One resume clause changes, and the retained command-log mutant rejects the exact omission.
- DONE: Named the exact file and line estimate before product edits.
  The design changes one file with 4 gross lines and a hard limit of 8 gross lines.
- DONE: Preserved the excluded behavior and authority.
  The design changes no grammar, format, authority, n28 mechanism, Codex behavior, Pi behavior, or zero-discovery behavior.
- DONE: Wrote the ideation report in Simplified English.
  This report gives the approach, expected surface, semantics, acceptance criteria, and test plan.

### Summary

The design adds the missing open-gate commit boundary before Sonnet presentation.
It uses one instruction file and the existing command-log falsifier.

## Stage Report: implementation

- DONE: Read the complete task and approved ideation report.
  The implementation followed the one-file product design and the eight-gross-line limit.
- DONE: Confirm the exact n28 Claude baseline and retained missing-commit falsifier.
  The baseline omitted `state commit`, and the focused oracle rejects the same omission.
- DONE: Change only `skills/fo-gate-lifecycle/SKILL.md` within eight gross lines.
  Product commit `9b561e4d1` replaces one line and keeps the skill at 6,993 bytes.
- DONE: Run the exact local bound Sonnet target and retain XPASS-green evidence.
  The target logged XPASS with no observed error after 411.66 seconds.
- DONE: Remove only the kky Sonnet binding after XPASS, then run exact local normal PASS.
  The unbound target passed after 512.24 seconds and logged a complete acknowledgment chain.
- DONE: Apply the Captain-approved Codex disposition.
  Commit `43ce24b6d` transfers the Codex binding and mirrored row to `272j6s25f9mry6nxbf4yjxvt`.
- DONE: Preserve local live evidence and use subscription authentication.
  Sonnet evidence remains in `/tmp/kky-bound-sonnet.Ny8xiL` and `/tmp/kky-unbound-sonnet.wXJRIp`.
- DONE: Preserve the two Codex finding artifacts without another live run.
  The collision and queued-selection artifacts remain in `/tmp/kky-bound-codex.PTOBkb` and `/tmp/kky-bound-codex-sequential.CHE1eq`.
- DONE: Run focused, format, full, race, registry, and active-owner checks.
  All final checks passed after the authorized stale-cache cleanup removed the disk blocker.
- DONE: Keep the final surface within the approved boundaries.
  The final diff has three insertions and three deletions across one product file and two authorized binding files.
- DONE: Commit and push the exact candidate and a Simplified-English implementation Stage Report.
  Candidate `43ce24b6d16e0287ab115bcbc4b774d08a710e05` matches its remote branch.

### Summary

The shared resume rule now commits every open gate before the structured reads and presentation.
The exact Sonnet journey passed with and without its expected-failure binding.
Codex selection remains bound to its new active owner, and Pi did not run.

## Stage Report: validation

- DONE: Independently verify the one-line Sonnet gate-commit rule delivers a durable clean open gate before presentation.
  Both retained Sonnet logs show prepare, commit, state head, structured reads, and presentation in the required order.
- DONE: Inspect retained bound XPASS and unbound Sonnet PASS evidence, component cap, exact candidate scope, and transferred Codex ownership without rerunning owned live/full/race checks.
  The evidence maps to commits `9b561e4d1` and `43ce24b6d`. The final diff is three one-line replacements.
- DONE: Confirm registry/active-owner/format evidence and no Codex or Pi behavior claim is made by kky.
  The focused registry and active-owner tests passed, and both changed Go files have standard formatting.

### Acceptance evidence

- DONE: AC-1
  Each Sonnet log has one successful prepare, then a successful commit, state head, checklist read, AC read, and presentation.
- DONE: AC-2
  The bound artifact has clean Sonnet semantics with the old binding. The unbound artifact has clean semantics after binding removal.
- DONE: AC-3
  Both streams report `state=open`. The entity has no completion or verdict. Each log has no successor dispatch after prepare.
- FAILED: AC-4
  A detached mutant proved that the oracle accepts a commit after structured reads and accepts continuation after a failed commit.
- DONE: AC-5 within the Captain-approved recarve
  Focused format, registry, and active-owner checks passed. Implementation owns the retained full and race results. Pi did not run.

### Reviewer findings

- Material evidence defect: the AC-4 oracle does not enforce the complete commit boundary.
  Released workflow: the local Sonnet default-headless gate stop.
  Observable harm: a wrong command order can pass the focused oracle.
  Authority: `value-ac[AC-4]` requires rejection of a missing, failed, or late commit.
  Trigger: prepare, structured read, successful commit, and state head return no oracle error.
  A second trigger continues after a failed commit, reads evidence, then commits successfully.
  Defect kind: evidence defect. Release scope: Material. Proposed owner: this task. Proposed disposition: narrow fix.
  The exact boundary is `assertRecordedGateHoldLog`. It must reject reads after prepare and before the successful commit.

### Checks run

- PASS: `TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle` passed its current mutation set.
  Removing the successful post-prepare commit makes this test fail.
- PASS: `TestRuntimeLiveRegistryReconciliation` and `TestRuntimeLiveTODOOwnersAreActive` passed.
  A stale Sonnet binding or an inactive Codex owner makes these tests fail.
- PASS: `git diff --check` and focused `gofmt -d` checks returned no output.
  Whitespace damage or nonstandard Go formatting makes these checks fail.
- EXPECTED FAIL: detached `TestValidationOracleRejectsLateCommit` failed both adversarial cases.
  The candidate oracle accepted late commit and failed-commit continuation logs.

### Summary

The retained Sonnet evidence proves a durable clean open gate and supports AC-1 through AC-3.
The candidate stays within scope and transfers Codex ownership without a Codex or Pi behavior claim.
Recommendation: REJECTED until the narrow AC-4 evidence defect is fixed.
