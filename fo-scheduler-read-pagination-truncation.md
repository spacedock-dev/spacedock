---
id: pgdyphtaqfx1zn0h7ax31e5h
title: FO scheduler reads must not silently truncate at the page limit
status: backlog
source: "trim-dispatch-core-stale-prose ideation and validation, 2026-08-15; captain 2026-08-15: necessary before 0.27"
started:
completed:
verdict:
score: "0.95"
worktree:
issue:
---

The FO event loop's scheduler reads (status --where, status --next) return paginated JSON with a default limit of 25, and the FO contract reads no pagination field. Past 25 mod-blocked or matching entities the loop silently never sees the rest: work is dropped with no signal. Captain: needed before 0.27 stable.

Direction space for ideation: the binary could serve FO-internal reads unpaginated (or with an explicit --all), or the contract could instruct consuming has_next - prefer the binary owning it so no prose depends on paging mechanics.

## Problem

{Ideation fills this in.}

## Proposed approach

{Ideation fills this in.}

## Out of scope

Captain-facing table rendering.

## Expected surface and tolerance

Estimate net LOC change: small, across binary and possibly one contract sentence.

## Acceptance criteria

**AC-1 - A scheduler read over more than 25 matching entities observes all of them.**
Verified by: a fixture workflow with 30 matching entities; baseline today observes 25.

**AC-2 - No silent cap remains: either no cap, or the output names the remainder.**
Verified by: the same fixture asserts the complete set or an explicit continuation signal.

**AC-3 - The suite stays green.**
Verified by: go test ./internal/status/ and contractlint.

## Test plan

One 30-entity fixture, both reads, plus the existing suite.
