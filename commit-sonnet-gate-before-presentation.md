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

## Stage Report: implementation (cycle 2)

- DONE: Add the approved AC-4 evidence correction in one test file.
  The oracle now requires a successful state head before the checklist and requires the checklist before the AC scan.
- DONE: Add failed-commit and late-after-read mutants before the oracle correction.
  The late-after-read mutant first failed because the old oracle accepted it. Both mutants pass after the correction.
- DONE: Verify the corrected order from retained Sonnet evidence.
  Both Sonnet logs show commit, state head, checklist, and AC scan in the required order.
- DONE: Preserve live artifacts and avoid an unnecessary live rerun.
  The retained evidence was sufficient. Codex and Pi did not run.
- DONE: Preserve the approved product and ownership scope.
  The one-line product change is unchanged. Codex remains bound to active owner `272j6s25f9mry6nxbf4yjxvt`.
- DONE: Run focused, full, race, registry, owner, format, and diff checks.
  All checks passed. The corrected candidate is `ee110954751dd1d9783e64194d154917a1b29e2a`.
- DONE: Keep the complete candidate within the approved correction cap.
  The diff has 11 insertions and four deletions across four files, for 15 gross lines.

### Summary

The corrected oracle rejects a failed commit and rejects structured reads before the successful committed state head.
The retained Sonnet evidence satisfies this stronger oracle without a new live run.

## Stage Report: validation (cycle 2)

- DONE: Revalidate exact corrected candidate `ee110954751dd1d9783e64194d154917a1b29e2a` and state report `80727bdd26bd867932ac43ba4f0014dd928a3e21`.
  Both exact commits were present, and the candidate matched its remote branch.
- DONE: Inspect the exact delta and retained Sonnet evidence without another live run.
  A detached focused test passed both retained Sonnet logs through the corrected oracle.
- DONE: Verify the failed-commit and late-after-read mutants.
  Both named mutants passed and returned their required gate-hold errors.
- DONE: Verify focused, full, race, registry, active-owner, and format evidence.
  All independent non-live checks passed on the exact candidate.
- DONE: Verify the approved component cap and product size.
  The four-file, 15-gross candidate keeps the 6,993-byte product and binds Codex to active owner `272j6s25f9mry6nxbf4yjxvt`.
- DONE: Preserve the candidate and make no Codex or Pi behavior claim.
  Candidate HEAD stayed exact. This validation did not run a live runtime.

- DONE: AC-1
  Both retained Sonnet logs pass the corrected order oracle.
- DONE: AC-2
  The retained bound and unbound Sonnet artifacts remain clean.
- DONE: AC-3
  The retained streams show an open gate with no terminal fields or successor dispatch.
- FAILED: AC-4
  The corrected oracle accepts a failed commit followed by a successful retry before the structured reads.
- DONE: AC-5 within the Captain-approved recarve
  Focused, full, race, registry, active-owner, diff, and format checks passed. Pi did not run.

### Reviewer findings

- Material evidence defect: the oracle permits a retry after a failed gate-state commit.
  Released workflow: the local Sonnet default-headless gate-stop observer.
  Observable harm: a failure-stop violation can receive a clean PASS or XPASS grade.
  Authority: `value-ac[AC-4]` requires the focused negative control to reject a failed commit.
  Trigger: prepare, failed commit, successful retry, state head, checklist read, and AC read.
  The detached `TestValidationOracleRejectsRetryAfterFailedCommit` reproduced the acceptance against the exact candidate.
  Defect kind: evidence defect. Release scope: Material. Proposed owner: this task. Proposed disposition: narrow fix.
  The oracle must reject any failed post-prepare state commit, even when a later retry succeeds.

### Checks run

- PASS: the focused oracle test passed, including `failed_commit` and `late_after_read`.
- PASS: both retained Sonnet logs passed the corrected oracle in a detached test.
- PASS: `go test ./...` and `go test ./... -race` passed.
- PASS: registry, active-owner, diff, and focused format checks passed.
- EXPECTED FAIL: the detached retry-after-failure mutant proved the remaining AC-4 gap.

### Summary

The correction fixes the reported late-read gap and all retained Sonnet behavior remains clean.
Recommendation: REJECTED because the oracle still accepts a retry after a failed gate-state commit.

## Stage Report: implementation (cycle 3)

- DONE: Add the approved retry-after-failure mutant before the oracle correction.
  The focused test failed because the old oracle accepted the retry sequence.
- DONE: Reject a failed post-prepare state commit before a successful retry can qualify the gate hold.
  The existing failed-only and late-read controls keep their original error classifications.
- DONE: Preserve the approved product and ownership scope.
  Product bytes and retained Sonnet evidence are unchanged. Codex remains bound to owner `272j6s25f9mry6nxbf4yjxvt`.
- DONE: Run the required changed-byte checks without another live run.
  Focused lifecycle, registry, owner, contract, format, cap, and diff checks passed.
- DONE: Keep the full candidate within the approved cap.
  Candidate `3f437a8cf17879ad400bf2876a5266a7801dac96` has 14 insertions and four deletions across four files.

### Summary

The oracle now rejects a failed state commit after prepare, even when a later retry succeeds.
The correction adds three gross lines. Full and race tests did not run in this correction cycle.

## Stage Report: validation (cycle 3)

- DONE: Revalidate exact candidate `3f437a8cf17879ad400bf2876a5266a7801dac96` and report `633c41360620d41a836e3e8ac0206665c5918ec4`.
  Both commits were present, and the candidate matched its remote branch.
- DONE: Verify the exact three-line retry correction and existing controls.
  The focused test passed the failed-only, failed-retry, late-read, and other existing mutants.
- DONE: Verify the complete candidate cap and unchanged product bytes.
  The diff has four files, 14 insertions, four deletions, and 18 gross lines. The product remains 6,993 bytes.
- DONE: Verify retained Sonnet evidence and transferred Codex ownership.
  Both Sonnet logs pass the new oracle. Codex remains bound to `272j6s25f9mry6nxbf4yjxvt`.
- FAILED: Verify all nonzero failed-commit retries and exact-head full/race evidence.
  An exit-2 retry passed the oracle, and prior full/race runs predate the changed relevant bytes.
- DONE: Preserve the candidate and unrelated workflow state.
  No candidate, live artifact, or unrelated entity changed during validation.
- DONE: AC-1 through AC-3
  The retained Sonnet logs and streams still prove the ordered durable open gate.
- FAILED: AC-4
  The oracle rejects exit 1 but accepts another nonzero commit failure followed by a successful retry.
- FAILED: AC-5
  The exact candidate has no full or race result after its relevant oracle and test bytes changed.

### Reviewer findings

- Material evidence defect: the retry guard matches only exit 1.
  Released workflow: the local Sonnet default-headless gate-stop observer.
  Observable harm: a supported nonzero commit failure and retry can receive a clean PASS or XPASS grade.
  Authority: `value-ac[AC-4]` requires the focused negative control to reject a failed commit.
  Trigger: prepare, exit-2 or exit-3 commit failure, successful retry, state head, checklist read, and AC read.
  A detached mutant reproduced exit-2 acceptance. The shipped conflict path proves exit 3 is supported.
  Defect kind: evidence defect. Release scope: Material. Proposed owner: this task. Proposed disposition: narrow fix.
- Material evidence defect: prior full/race results do not cover the exact candidate.
  Released workflow: validation of candidate `3f437a8cf17879ad400bf2876a5266a7801dac96`.
  Observable harm: the changed oracle and mutant bytes lack exact-head suite evidence.
  Authority: `value-ac[AC-5]` requires full and race checks to pass.
  Trigger: prior runs used `ee1109547`, then `3f437a8cf` changed `claude_runtime_helpers_test.go`.
  Defect kind: evidence defect. Release scope: Material. Proposed owner: this task. Proposed disposition: run required checks.

### Checks run

- PASS: the focused oracle, registry, active-owner, diff, and format checks passed.
- PASS: detached retained-log tests passed for bound and unbound Sonnet evidence.
- EXPECTED FAIL: the detached nonzero-retry mutant proved the remaining AC-4 gap.
- SKIPPED: live, full, and race runs did not run in this cycle.
  Live evidence was retained. Prior full/race evidence was not valid for the changed relevant bytes.

### Summary

The three-line change closes the exact exit-1 retry gap and preserves all retained Sonnet value.
Recommendation: REJECTED for the nonzero retry gap and missing exact-head full/race evidence.

## Stage Report: implementation (cycle 4)

- DONE: Prove the current-main bound Sonnet default-headless target with retained local evidence.
  Bound run XPASSed in 309.34s; `/tmp/kky-bound-sonnet-current-main.DdshkH/evidence` proves prepare, commit, state head, structured reads, and an open gate.
- DONE: Remove only the Sonnet binding after XPASS, then prove one unbound normal PASS.
  Commit `fbbb225a1` removes the two binding rows; corrected unbound run passed in 341.12s with evidence at `/tmp/kky-unbound-sonnet-corrected.a7weFp/evidence`.
- DONE: Preserve product/test firewall: no hooks, instrumentation, or unrelated runtime changes.
  The final diff from `7cf03cd94` changes only three test files; focused registry and active-owner checks passed.
- DONE: Preserve and reconcile the existing task commits on exact current main.
  Rebase retained commits `c96df8b01`, `15fd168f7`, and `f2126a56a`; duplicate product commit was skipped and the premature binding-removal commit was deferred until XPASS.
- DONE: Diagnose the first unbound evidence failure with one smallest correction.
  Retained log showed `dispatch build --help` was miscounted as a dispatch; test `TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle` now fails if that harmless probe is rejected.
- DONE: Run final relevant-byte checks.
  `go test ./...`, formatting, registry, active-owner, and focused race checks passed; the full race run had one unrelated 250ms timing failure, whose isolated race rerun passed.

### Summary

Current main produced a bound Sonnet XPASS, and removing only the Sonnet XFAIL produced an unbound normal PASS after one evidence-oracle correction for a harmless capability probe. Product instructions and runtime behavior were unchanged; commits `c96df8b01` through `c9efdf853` contain the reconciled test-only candidate and retained live evidence names the exact gate order.
