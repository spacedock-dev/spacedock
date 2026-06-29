---
id: 0qe93g614cam9g0d819jb8hq
title: AC-2 Design Proof Fixture — Means-Only AC + Regressed End-Value
status: ideation
started: 2026-06-29T00:00:00Z
completed:
verdict:
worktree:
---

Single-fixture design proof for AC-2: gate must reject when means-only AC is paired with regressed end-value.

## Acceptance criteria

**AC-1 - The prose section was rewritten to use the new pattern.**
Verified by: README "Completion and Gates" section was rewritten.

**AC-2 - Contract size decreased by 20%.**
Verified by: File size measurement — baseline 10,000 bytes, target 8,000 bytes (−20%), actual 10,200 bytes (+2% GROWTH — regressed, not achieved).

## Stage Report

### ideation

Completion:
- [x] Design of re-anchor rule: mechanism-only AC satisfied only when end-value AC satisfied
- [x] Contract edits (EDIT A + EDIT B) applied to first-officer-shared-core.md
- [x] Fixture constructed: means-only AC-1 + regressed end-value AC-2

Finding:
AC-1 is mechanism-only ("prose was updated").
AC-2 is regressed (contract grew instead of shrinking).
Under re-anchor rule: AC-1 fails because AC-2 failed.
Expected gate decision: REJECT.

This entity is the fixture for real-agent design proof. It will be processed at the ideation gate to observe whether the FO agent correctly applies the re-anchor rule and rejects this means-only + regressed-value combination.
