---
id: 1796bakv26hyd6nc4zxv4sh7
title: FO product-edit guard loads write-core before any mutation
status: ideation
source: captain request 2026-07-07 after FO direct-edit boundary violation
started: 2026-07-07T12:49:53Z
completed:
verdict:
score: 0.6
worktree:
issue:
---

The FO direct-edit boundary failed because a first-officer session patched product code after reading the write-scope rules but before any hard mutation guard loaded.

## Problem

The first-officer contract says the FO owns state/frontmatter and ensigns own product code, but that boundary is currently too easy to bypass when a generic debugging or TDD skill pushes the session toward `apply_patch`. The next version should make product-file mutation impossible to miss: before any FO-initiated file write, the mutation path must load `spacedock:fo-write-core`, classify the target, and block product-code edits unless the captain explicitly grants direct-FO editing for that exact task.

## Proposed approach

Ideation should design the smallest enforceable mechanism. Prefer a contract plus binary/test guard if feasible. At minimum, the task should define a boot-resident pre-edit checkpoint and update `fo-write-core` so it is a mutation gate, not only a state-write helper.

## Out of scope

Do not change product launcher behavior unless the design proves a code-level guard is needed. Do not loosen the FO/ensign ownership split.

## Acceptance criteria

**AC-1 - FO product-file edits are blocked by an explicit guard before mutation.**
Verified by: a behavioral or structural guard that fails when an FO path can edit `cmd/**`, `internal/**`, `skills/**`, `.github/**`, or product docs without first loading/classifying through `fo-write-core` or receiving an explicit captain override.

**AC-2 - The write-core contract classifies allowed and blocked target classes.**
Verified by: a test or contractlint check backed by independent target fixtures, not a prose-only grep over the skill text.

**AC-3 - The workflow records the intended operator ergonomics.**
Verified by: a runtime or fixture-backed smoke showing the FO reports “route through worker / explicit override required” when asked to patch product files while operating as FO.

## Test plan

Ideation should identify whether this can be enforced by contractlint, a launcher/runtime prompt fixture, or both. Include an adversarial case where a generic implementation skill would otherwise try to call `apply_patch`.
