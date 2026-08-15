---
id: x2ezetxr82pztr4pqt1g4dhx
title: Rule the superseded verdict into or out of the schema vocabulary
status: backlog
source: "scope-validate-warnings ideation, 2026-08-15: 4 archived entities carry verdict superseded, a token the conventional enum [PASSED REJECTED] never admitted"
started:
completed:
verdict:
score:
worktree:
issue:
---

Writers intentionally emit verdict: superseded for superseded entities, but the schema enum admits only PASSED and REJECTED. Archived-scope warnings are silenced now, so the bite is forward-looking: the next active entity superseded on purpose warns as invalid, and any tool trusting the enum misreads the four archived records. Decide: admit superseded (and define its semantics) or route supersede through a different field and stop the writer.

## Problem

{Ideation fills this in.}

## Proposed approach

{Ideation fills this in.}

## Out of scope

Hand-editing archived entities (publish-only, per 6c45fd59c precedent).

## Expected surface and tolerance

Estimate net LOC change: small - schema line, writer, one test.

## Acceptance criteria

**AC-1 - An intentional supersede on an active entity produces zero validate warnings under the ruled vocabulary.**
Verified by: an active-entity fixture superseded through the supported path; baseline today warns.

**AC-2 - The ruling is enforced where verdicts are written, not only documented.**
Verified by: the writer-side test rejecting or accepting the token per the ruling.

**AC-3 - The suite stays green.**
Verified by: go test ./internal/status/ and the schema conformance tests.

## Test plan

One fixture per ruling arm, existing conformance suite as the floor.
