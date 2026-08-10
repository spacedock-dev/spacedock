---
title: Continue Codex to validation after implementation
status: ideation
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
---
## Problem

The exact Codex auto-continue journey stopped after implementation. The First Officer did not dispatch or run fresh validation. No validation Stage Report appeared.

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

## Ideation requirements

- Name exact files and gross/net estimate before product edits.
- Identify the smallest behavior boundary and one falsifying control.
- Keep local failing baseline, repaired normal PASS, validation, PR, and merge flow.
