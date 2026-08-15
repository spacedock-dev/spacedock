---
id: 7c4w88fnmnbtc0tgkrvx0vxj
title: Remove the provider-evidence gate fields
status: ideation
source: "Captain directive, 2026-08-14: value review found zero value; no writer, no retained bytes, no verifier"
started: 2026-08-15T02:55:32Z
completed:
verdict:
score:
worktree:
issue:
sprint: durable-decisions
gates:
    version: 1
    records:
        - id: gate:7c4w88fnmnbtc0tgkrvx0vxj:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:7c4w88fnmnbtc0tgkrvx0vxj-backlog-1
              briefing:
                id: briefing:7c4w88fnmnbtc0tgkrvx0vxj:backlog:attempt-1:revision-1
                digest: sha256:163027efa9408cfc6b01e548d89f063de4a6cbc6257f1022c08a4b536bc2b27f
                request-digest: sha256:6107a996b958133a0620aa02c9a5ee9071804f657c4ae61347c9c4b190ec1ccd
                room-ref: ./remove-provider-evidence-fields/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:7c4w88fnmnbtc0tgkrvx0vxj:backlog:1
                briefing: briefing:7c4w88fnmnbtc0tgkrvx0vxj:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T02:53:31.794778Z"
                decision: approve
                reason: 'Captain ruling 2026-08-14 (dispatch them): approved into ideation'
              application:
                target-stage: ideation
                state: consumed
---

Remove `ProviderEvidence` from the gate model: the struct, the `provider-evidence` field, the `Validate` branches that police it on open and withdrawn attempts, and their tests. Zero provider-closed attempts exist across 424 recorded attempts, so no stored record carries the field. Re-introduce the fields only together with a real provider integration: writer, retention, and verifier.

## Problem

{Ideation fills this in.}

## Proposed approach

{Ideation fills this in.}

## Out of scope

The retained-authority checks for request.json. They have a live verifier. Keep them.

## Expected surface and tolerance

Estimate net LOC change: -NNN, across ~3 files. No observable semantics change: no stored record carries the field.

## Acceptance criteria

**AC-1 - No source or test references the fields.**
Verified by: `grep -rl "provider-evidence\|ProviderEvidence" internal skills cmd` returns no matches.

**AC-2 - Existing state still decodes.**
Verified by: the suite stays green, plus a one-off gate read over docs/dev/.spacedock-state recorded in the report.

**AC-3 - The change removes more lines than it adds.**
Verified by: cumulative line delta of the diff against origin/main is negative.

## Test plan

Deletion plus the existing suite. One state-decode pass as one-off evidence.
