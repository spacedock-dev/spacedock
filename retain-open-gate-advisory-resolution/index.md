---
title: Retain advisory results without closing an open gate
status: backlog
source: "Captain-approved local task from debrief issue 3, 2026-08-02"
score: 0.7
worktree:
issue:
id: vp4f3wpf9bpht578yd64b12d
---

The gate system must retain an advisory review result without applying a binding gate decision.

## Problem

The normal open-gate path has no durable record for an advisory approval or rejection from a standing reviewer. The correction-round recorder is not a general record for an open gate. This forces advisory evidence into chat or into a binding decision path.

## Proposed approach

Add a workflow-neutral advisory result record for an open gate. Store the reviewer, result, evidence reference, and timestamp. Keep the gate open. Keep Captain approval or rejection as the only binding gate decision. Coordinate with `workflow-neutral-advisory-round-recorder` so the two tasks do not duplicate structural storage.

## Out of scope

Do not auto-advance or close a gate from advisory evidence. Do not replace the existing correction-round recorder. Do not add compatibility wrappers for the unreleased v1 surface.

## Acceptance criteria

**AC-1 - An advisory approval or rejection persists on an open gate without changing stage status.**
Verified by: a CLI or fixture test that records each advisory result, reads the durable gate state, and asserts the result is present while the entity remains in its current stage.

**AC-2 - A later Captain decision remains independent of the advisory result.**
Verified by: a test that records advisory approval, then records a Captain rejection, and asserts the rejection routes the gate without treating the advisory result as binding authority.

## Test plan

Run focused gate-record tests and the existing gate lifecycle tests. Run `go test ./...` for the full suite.
