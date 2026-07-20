---
id: 0qe93g614cam9g0d819jb8hq
title: AC-2 Design Proof Fixture — Means-Only AC + Regressed End-Value
status: ideation
started: 2026-06-29T00:00:00Z
completed:
verdict:
worktree:
sprint: 0260-proportionality
group: test-cleanups
gates:
  version: 1
  current:
    gate: gate:docs-dev:0qe:ideation
    attempt: gate-attempt:0qe-ideation-1
  records:
    - id: gate:docs-dev:0qe:ideation
      stage: ideation
      current-attempt: gate-attempt:0qe-ideation-1
      attempts:
        - id: gate-attempt:0qe-ideation-1
          sequence: 1
          state: open
          briefing:
            id: briefing:0qe-ideation-1
            digest: sha256:10581a261293596e17f359b840d8f610f431141daadb4f09763d5185fe674670
          note: "Captain hold via float, 2026-07-20 (resolution:actor-1784524018982597000): the FO briefing lacked orientation — no plain statement of what the entity is or what AC-2 refers to. Attempt stays open; re-present with real context or fold into the lure-scenario catalog per the pending captain choice."
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

## Merged scope (adopted cross-review re-lock, 2026-07-20)

Absorbs `ac2-reanchor-scenario-falsifiable` — strengthening the PR-441 AC-2 re-anchor live scenario so it can fail on the regression it polices serves the same end as this fixture (the means-only-AC / regressed-end-value detection actually detecting); one member, one falsifiability outcome. Banked ideation here is not re-ideated; the absorbed scope folds in at this entity's gate.
