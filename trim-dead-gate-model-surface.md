---
id: ebgwr177kjjs6w5thhywz408
title: Trim dead gate model and projection surface
status: backlog
source: "0.27 cut audit (2026-08-14), adversarially verified; captain directed filing"
sprint: durable-decisions
started:
completed:
verdict:
score:
worktree:
issue:
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
