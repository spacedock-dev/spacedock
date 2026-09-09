---
title: Honor slug identity in correction-round recording
status: backlog
source: Captain report of email triage correction-round failure, 2026-09-09
started:
completed:
verdict:
score: 0.7
worktree:
issue:
pr:
id: zjk8yfarexrde9bj7d2z92p0
---

Correction-round recording must accept valid slug-style tasks without a stored id. The captain requested filing only; no dispatch or implementation is authorized by this filing.

## Problem

An email triage correction was routed to intake. The correction-round recorder then refused with "entity has no identity for correction round" because the task had a blank id. Its workflow declares id-style: slug, so the slug already supplies stable identity. Adding an id is a workaround, not a repair to invalid workflow state.

## Risk evidence

- internal/status/new.go intentionally leaves slug-style seeds without a stored id.
- internal/gates/prepare.go entityIdentity resolves slug-style identity from the entity path and workflow configuration.
- internal/gates/round.go resolveRound directly reads the id frontmatter key and rejects an absent or blank value. It also validates the round pointer against that directly read value.

These contracts disagree for supported slug-style entities. This finding is based on source inspection and the reported runtime refusal; no new reproduction has been run in this filing.

## Proposed approach

Use the canonical workflow-aware identity resolution for correction-round recording and its replay, read, and validation paths. Preserve existing stable-ID behavior and immutable round semantics. Do not require users to add redundant IDs or change workflow identity style.

## Acceptance criteria

**AC-1 — Valid slug-style entities can record and replay correction rounds without a stored id.**
Verified by: a fixture creates a folder-form slug entity through spacedock new, records a valid closed correction round, and replays it successfully. Restoring direct nonblank-id enforcement must fail the test. Assert identity is the slug, never index.

**AC-2 — Round identity checks remain consistent across supported entity forms and ID styles.**
Verified by: existing round fixtures cover stored-ID entities; extend the primary owner for absent and blank slug id, flat and folder forms, pointer identity mismatch, and read/validation of the retained round. Wrong pointer identity must still refuse.

**AC-3 — Recording preserves unrelated workflow state and immutable replay behavior.**
Verified by: existing round tests assert unchanged status, gate authority, and body outside the round pointer; exact replay remains idempotent and divergent replay remains refused.

## Test plan

Start with the existing internal/gates/round_test.go proof owner and the smallest real new-to-round CLI fixture for workflow resolution. Write the failing regression before implementation. Exercise behavior and resulting bytes, not source-text presence. Run required Go and race checks for the implementation; do not mutate live email batches.

## Out of scope

Email policy changes, adding redundant IDs to live tasks, identity migration, the separate stage-completion design, and a new correction-round protocol.
