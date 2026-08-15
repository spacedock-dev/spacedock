---
id: ebgwr177kjjs6w5thhywz408
title: Trim dead gate model and projection surface
status: ideation
source: "0.27 cut audit (2026-08-14), adversarially verified; captain directed filing"
sprint: durable-decisions
started:
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:ebgwr177kjjs6w5thhywz408:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:ebgwr177kjjs6w5thhywz408-backlog-1
              briefing:
                id: briefing:ebgwr177kjjs6w5thhywz408:backlog:attempt-1:revision-1
                digest: sha256:fc72271948d2273a9e0ede89b05604c4741457bd98fadaf50b9ebd1e7457b14e
                request-digest: sha256:5f016522677589c77f88f557e4e5f1d306c9f507dbce176995941f0789e7da34
                room-ref: ./trim-dead-gate-model-surface/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ebgwr177kjjs6w5thhywz408:backlog:1
                briefing: briefing:ebgwr177kjjs6w5thhywz408:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T02:53:37.741016Z"
                decision: approve
                reason: 'Captain ruling 2026-08-14 (dispatch them): approved into ideation'
              application:
                target-stage: ideation
                state: consumed
---

Remove four verified-dead pieces of the gate model surface. No consumer exists for any of them.

1. The nine projected gate-* status columns (internal/status/discover.go:219-227). No skill, CI lane, or live-history query ever named one. Keep the gates.Read call at discover.go:216 - it feeds the load-bearing gate-readiness chain. Also drop the frontmatter-contract.md:13 sentence that names three of the nine.
2. Summary.Condition and Summary.Eligible (internal/gates/model.go:95-96). Debris from the eligibility cut (013c8729e). No writer, no reader.
3. Annotation.Target, .Kind, .Body (internal/gates/review.go:18-20). Decoded, never validated, projected, or read. The retained raw review log preserves the bytes.
4. The ReadWithWarnings and Diagnostic aliases (internal/gates/io.go). Zero callers, and an internal package permits no external caller.

Captain override recorded here: the archived decision remove-standalone-gate-eligibility.md:173 said to keep gate-state, gate-application, and gate-target-stage. That retention named no consumer, and none exists. Approval of this entity overrides that sentence.

## Problem

{Ideation fills this in.}

## Proposed approach

{Ideation fills this in.}

## Out of scope

The gate-readiness chain, gates.Read, filterApplicationMappings read tolerance, ReadDiagnostics and its status --validate consumer, resolution includes, mediaType.

## Expected surface and tolerance

Estimate net LOC change: -NNN, across ~5 files. No observable semantics change: no consumer exists for the removed surface.

## Acceptance criteria

**AC-1 - The change removes more lines than it adds.**
Verified by: cumulative line delta against origin/main is negative.

**AC-2 - No source references the removed fields or aliases.**
Verified by: grep for the nine column names, Summary.Condition, Summary.Eligible, ReadWithWarnings, and gates.Diagnostic returns no matches in cmd, internal, skills.

**AC-3 - The suite stays green.**
Verified by: go test ./... and go test ./... -race pass.

## Test plan

Deletion plus the existing suite. Own tests of the removed projections are deleted with them.
