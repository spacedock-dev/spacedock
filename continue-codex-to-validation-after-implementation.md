---
title: Continue Codex to validation after implementation
status: validation
score: "0.90"
source: "PR #664 Codex auto-continue failure, 2026-08-10"
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-evidence
id: v8pcpdmrdfmq7emm65cjdc4p
gates:
    version: 1
    records:
        - id: gate:v8pcpdmrdfmq7emm65cjdc4p:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:v8pcpdmrdfmq7emm65cjdc4p-backlog-1
              briefing:
                id: briefing:v8pcpdmrdfmq7emm65cjdc4p:backlog:attempt-1:revision-1
                digest: sha256:ceb846a40def6bfd03d0986cab331263c5349ac078c6bbf35ca24de60708f1f0
                request-digest: sha256:0f8fe2d8e42b2052f020580ef332395f4710b57d3aa288cf7b380258df269c65
                room-ref: ./continue-codex-to-validation-after-implementation/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:v8pcpdmrdfmq7emm65cjdc4p:backlog:1
                briefing: briefing:v8pcpdmrdfmq7emm65cjdc4p:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T19:01:54.942236Z"
                decision: approve
                reason: Captain directed immediate ideation for the exact Codex auto-continue product repair.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:v8pcpdmrdfmq7emm65cjdc4p:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:v8pcpdmrdfmq7emm65cjdc4p-ideation-1
              briefing:
                id: briefing:v8pcpdmrdfmq7emm65cjdc4p:ideation:attempt-1:revision-1
                digest: sha256:454e049070257771e2d127c6c5cd945cedf45cb591496f4b64554b4bf08d3cd5
                request-digest: sha256:11cef063c866477dac13329e146e5669cacc88d6715660cf5e3203507aa26541
                room-ref: ./continue-codex-to-validation-after-implementation/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:v8pcpdmrdfmq7emm65cjdc4p:ideation:1
                briefing: briefing:v8pcpdmrdfmq7emm65cjdc4p:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T19:07:31.694893Z"
                decision: approve
                reason: Captain directed this exact Codex auto-continue repair. The one-file design preserves every excluded surface.
              application:
                target-stage: implementation
                state: consumed
started: 2026-08-10T19:02:21Z
worktree: .worktrees/spacedock-ensign-continue-codex-to-validation-after-implementation
---
## Problem

The exact Codex auto-continue journey stopped after the validation dispatch. The First Officer did not wait for validation to finish. No validation Stage Report appeared.

## Value

After implementation completes, Codex dispatches and runs fresh validation. Codex commits the validation Stage Report and then enters the validation gate.

## Scope

- Repair only the Codex auto-continue behavior after implementation.
- Use PR #664 artifact 9075649719 as the exact baseline.
- Keep the smallest product instruction or runtime surface plus focused proof.
- Do not change dvd, n28, Pi, XFAIL policy, or unrelated journeys.
- Do not add a permanent XFAIL.
- Use local Codex subscription authentication before required PR CI.

## Acceptance criteria

- AC-1: Exact local Codex `TestLiveCommonAutoContinueAfterImplementation` passes normally and retains artifacts.
- AC-2: Implementation completion dispatches and runs one fresh validation worker.
- AC-3: Validation commits a complete Stage Report before the workflow enters the validation gate.
- AC-4: Focused controls reject missing validator dispatch, missing validator run, missing report, or early gate entry.
- AC-5: Full, race, format, registry, active-owner, and required exact PR checks pass. Pi remains skipped.

## Baseline evidence

- Released user and workflow: Codex auto-continue after implementation.
- Observable harm: the workflow stops without fresh validation or its committed report.
- Value authority: `TestLiveCommonAutoContinueAfterImplementation` requires implementation-to-validation continuation.
- Trigger: run 31419396371, job 93556549760, artifact 9075649719. The target failed after 179.42 seconds with `no ## Stage Report: validation appeared`.

The split-root transcript gives the exact failure sequence:

- The First Officer advanced the task to validation.
- `dispatch build --stamp` completed successfully.
- The First Officer spawned the fresh validator, Bohr.
- The First Officer called the dispatch the requested bounded outcome.
- The First Officer ended the turn without a `wait_agent` call.
- The final message said that Bohr remained active.

The single-root transcript is the working comparison. It waited for the validator, read the validation report, and then entered gate handling.

## Proposed approach

Change only the Codex async-dispatch binding. Add these two sentences after the current async-dispatch rule:

> A fresh spawn does not satisfy a single-entity stop. When this worker is the only unresolved work, continue async monitoring until its completion signal.

The existing completion rule then requires report verification and the next workflow action. No new state, command, or scheduler mechanism is necessary.

The simplest alternative is a prompt change in the live fixture. That change is insufficient because normal single-entity prompts must get the same behavior.

## Expected surface

- `skills/first-officer/references/codex-first-officer-runtime.md`: two insertions, no deletions.
- Expected product total: one file, two gross lines, net plus two lines.
- Tolerance: one file and no more than four gross lines.

The design changes only Codex runtime behavior. A fresh worker spawn is no longer a valid single-entity stop.

The design does not change command grammar, stored formats, authority, Claude behavior, Pi behavior, or XFAIL policy.

## Test plan

- Run `TestAutoContinueNegativeStoppedAfterImplementation`. Its status-only case rejects early gate entry and a missing validation report.
- Run exact local Codex `TestLiveCommonAutoContinueAfterImplementation` with subscription authentication.
- Retain the local Codex artifacts. The split-root transcript must show one fresh validator, async monitoring, and its completion signal.
- Make sure that the validation Stage Report commit exists before gate handling begins.
- Run the full, race, format, registry, and active-owner checks.
- Run the exact required Codex and Sonnet PR checks. Keep Pi skipped.

The exact live journey is the falsifying control. It fails if dispatch, execution, report creation, or completion monitoring is missing.

## Risk evidence

Artifact 9075649719 exercised the highest-risk runtime handoff. It proves that dispatch succeeds and that the stop decision is the failed boundary.

## Ideation requirements

- Name exact files and gross/net estimate before product edits.
- Identify the smallest behavior boundary and one falsifying control.
- Keep local failing baseline, repaired normal PASS, validation, PR, and merge flow.

## Stage Report: ideation

- DONE: Read the complete task and exact PR #664 artifact.
  Artifact 9075649719 contains both auto-continue fixture transcripts and the exact failed job result.
- DONE: Inspect the Codex auto-continue contract and current instruction path.
  The Codex async binding requires monitoring, but it does not reject a fresh spawn as a single-entity stop.
- DONE: Define the smallest product repair and focused falsifying controls.
  The repair adds two sentences to the Codex runtime adapter. The existing exact live journey is the primary control.
- DONE: Name exact files and gross/net estimate before product edits.
  The estimate is one file, two gross lines, and net plus two lines, with a four-line tolerance.
- DONE: Preserve dvd, n28, Pi, XFAIL policy, and unrelated journeys.
  The design changes only the Codex single-entity stop boundary after a fresh spawn.
- DONE: Write and push a Simplified-English ideation Stage Report.
  This report records the evidence, design, surface, semantic boundary, and test plan.

### Summary

The exact failure is a premature Codex stop after successful validation dispatch. The smallest repair is a two-line Codex runtime rule.

## Stage Report: implementation

- DONE: Read the complete task and approved ideation report.
  The approved report limits the repair to one Codex runtime file and four gross lines.
- DONE: Confirm the exact PR #664 baseline and split-root failure sequence.
  Artifact 9075649719 stops after the validator spawn and contains no validation Stage Report.
- DONE: Change only skills/first-officer/references/codex-first-officer-runtime.md within four gross lines.
  Commit `c42b1dfdf` adds two lines to the Codex async-dispatch binding.
- DONE: Run the exact local Codex auto-continue target and retain normal PASS artifacts.
  The retained run passed in 357.10 seconds and fails if validation dispatch, monitoring, or the report is absent.
- DONE: Run focused, format, full, race, registry, and active-owner checks.
  The negative control, `gofmt`, both full suites, registry check, and active-owner check passed.
- DONE: Keep dvd, n28, Pi, XFAIL policy, and unrelated journeys unchanged.
  The candidate changes only the Codex single-entity stop rule.
- DONE: Commit and push the exact candidate and a Simplified-English implementation Stage Report.
  Branch `spacedock-ensign/continue-codex-to-validation-after-implementation` contains pushed commit `c42b1dfdf`.

### Summary

Codex now continues monitoring after a fresh validation spawn when no other work is unresolved. The retained evidence is in `.spacedock-evidence/v8pcpdmrdf-local-codex`.

## Stage Report: validation

- DONE: Confirm exact candidate, remote match, base, and clean tracked worktree.
  Local and remote HEAD are `c42b1dfdf`. The base is `0bbe9d46c`. Only retained evidence is untracked.
- DONE: Inspect the exact one-file, four-gross-line product scope.
  The diff adds two lines to `codex-first-officer-runtime.md` and changes no other tracked file.
- DONE: Inspect retained local Codex PASS artifacts for validator monitoring, completion, report commit, and gate entry.
  Both fixtures waited and reached committed validation reports before gate handling.
- DONE: Run focused negative, registry, owner, format, and diff checks without repeating owned full, race, or live runs.
  The focused negative, registry, owner, format, and diff checks passed.
- DONE: Report each finding with four workflow evidence fields.
  Finding V8-1 records all four fields below.
- DONE: Write and push a Simplified-English PASSED or REJECTED validation Stage Report.
  This report recommends REJECTED because the focused proof does not establish AC-2 or AC-4.

### Acceptance evidence

- AC-1: The retained local Codex run passed both fixture variants in 357.10 seconds.
- AC-2: The artifacts contain no structured `spawn_agent` event for the validator.
- AC-3: Each fixture contains a committed validation report before gate handling.
- AC-4: The negative oracle rejects missing state, but it accepts a report without worker evidence.
- AC-5: The local format, registry, owner, and diff checks passed. Required PR checks remain pending.

### Material findings

- V8-1: Evidence defect, Material, current-task ownership, Needs decision.
  Released user and normal workflow: supported Codex auto-continue after implementation.
  Observable harm: the oracle accepts an advanced state without a structured fresh validator dispatch or run.
  Authority: `value-ac[AC-2]` requires one fresh validation worker to dispatch and run.
  Trigger evidence: `single-root/codex-exec.jsonl` has zero `spawn_agent` events and one empty-target `wait` event.
  The retained `split-root/codex-exec.jsonl` file has the same event counts and empty target.

### Summary

The two-line instruction change is within the approved scope. The retained journey reaches validation, but its proof cannot establish the fresh worker boundary.

REJECTED. The First Officer must authorize a proof-surface change or request a design decision.

## Stage Report: implementation (cycle 2)

- DONE: Rebase the product repair onto the exact n28 landing.
  Rebased head `46368d1c1` preserves the one-file, two-line product change on `5fca1605d`.
- DONE: Add the approved three-file evidence reset within its cap.
  The uncommitted test reset changes three existing files and has exactly 80 gross lines.
- FAILED: Prove validation-stage pending, armed, and consumed acknowledgment states.
  The exact single-root run retained only the `pending` audit ref for epoch `30c04c28fc1482a5b240406bb4a2e121`.
- SKIPPED: Prove the later complete-and-committed predicate and open gate.
  The receipt-chain check failed first, so later evidence cannot repair the missing worker start.
- DONE: Preserve the product head, evidence reset, and all artifacts without claiming a solution.
  New evidence remains at `.spacedock-evidence/v8pcpdmrdf-n28-reset`; earlier evidence also remains unchanged.
- DONE: Avoid transcript parsing, Pi work, and new production files.
  The check used immutable n28 acknowledgment refs and stopped at the first missing state.

### Summary

HOLD. The landed n28 mechanism did not attest the validation worker start in the exact Codex journey.

A separate mechanical-continuation prerequisite owns the missing host primitive. The v8 product repair and all evidence remain preserved for that prerequisite.

## Stage Report: implementation (cycle 3)

- DONE: Reproduce current-main Codex auto-continue with isolated-home native lifecycle evidence.
  Exact unbound local target passed both fixtures in 392.15 seconds; artifacts are under `.spacedock-evidence/v8pcpdmrdf-minimized-unbound-pass`.
- DONE: Prove one validator spawn, correlated handle/completion, committed validation report, and open gate.
  The isolated-home parent rollout oracle requires exactly one validation spawn, its returned task path, matching completion before gate preparation, a report-bearing commit, and canonical open gate per fixture.
- DONE: Remove only the auto-continue binding after XPASS, then prove one unbound normal PASS.
  Commit `db7d5de72` removes the one Codex XFAIL row and its one registry expectation after the bound evidence, then the exact unbound target passed both variants.
- DONE: Keep the correction narrow and avoid the rejected acknowledgment mechanism.
  Product scope is two added Codex runtime lines; no dispatch acknowledgments, instrumentation, state protocol, fixture coaching, Pi behavior, or production Go code changed.
- DONE: Run final format, diff, and registry checks on the exact committed bytes.
  `gofmt -w ./cmd ./internal`, `git diff --check`, and `TestRuntimeLiveRegistryReconciliation` passed before commit.
- DONE: Commit and push the minimized candidate.
  Branch `spacedock-ensign/continue-codex-to-validation-after-implementation` is pushed at `db7d5de72` on exact base `7cf03cd94`.

### Summary

Codex now monitors the sole fresh validator through completion and uses the reported active entity as the gate artifact with committed `README.md` as the distinct reference. Product scope is +2/-0 lines; isolated-home observation and shared live plumbing are +65/-11 lines; XFAIL binding and registry expectation removal are +1/-2 lines, for +67/-13 versus exact current main.

## Stage Report: validation (cycle 2)

- DONE: Verify the minimized exact candidate proves native validator spawn, handle, completion, committed report, and open gate.
  Exact HEAD `db7d5de72` passed both live fixtures in 386.91 seconds; the oracle fails if one native spawn, its returned handle, matching completion order, durable report commit, or canonical open gate changes.
- DONE: Reject stale n28 machinery, instrumentation, hooks, broad registry churn, or test-driven product behavior.
  Diff from base `7cf03cd94` has no n28 acknowledgments, production Go, hooks, or fixture coaching; product behavior is only +2 Codex runtime lines, with native observation confined to tests.
- DONE: Confirm only the Codex auto-continue XFAIL was removed after bound XPASS and exact-candidate PASS.
  Retained bound run at `598149370` predates commit `db7d5de72`; that commit removes only owner `v8pcpdmrdfmq7emm65cjdc4p` and its matching registry expectation, and unbound exact HEAD passed.

### Acceptance evidence

- AC-1: `TestLiveCommonAutoContinueAfterImplementation` passed single-root and split-root normally at exact HEAD; retained artifacts are under `.spacedock-evidence/v8pcpdmrdf-validation-exact-head`.
- AC-2: The live oracle consumed the isolated Codex parent rollout and required exactly one validation `spawn_agent`, its returned task path, and a matching completion author before gate preparation.
- AC-3: The oracle resolved the active report entity, required one DONE validation section, found its Git commit with `git log -S`, and read the canonical validation gate as open.
- AC-4: Native lifecycle controls reject public-only output, missing or late completion, missing returned handle, and ambiguous rollout; `TestAutoContinueNegativeStoppedAfterImplementation` rejects missing state or report.
- AC-5: `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, registry reconciliation, active-owner, and diff checks passed. Codex and Sonnet PR checks remain required before merge; no PR exists. Pi remains skipped.

### Reviewer findings

- No material finding. The former V8-1 evidence defect is closed by native rollout correlation and exact live PASS.
- Deferred delivery check: run required Codex and Sonnet PR lanes after PR creation; promote to Material if either fails or does not exercise commit `db7d5de72`.

### Summary

PASSED for gate and PR progression. The minimized candidate proves the complete native validation lifecycle without the rejected n28 mechanism or product instrumentation; required PR CI remains a pre-merge obligation.
