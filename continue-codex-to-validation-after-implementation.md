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
