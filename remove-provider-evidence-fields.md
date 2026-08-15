---
id: 7c4w88fnmnbtc0tgkrvx0vxj
title: Remove the provider-evidence gate fields
status: backlog
source: "Captain directive, 2026-08-14: value review found zero value; no writer, no retained bytes, no verifier"
started:
completed:
verdict:
score:
worktree:
issue:
sprint: durable-decisions
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
