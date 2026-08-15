---
id: 71btbxdrken4kdmfsk0vptav
title: Trim stale prose in the dispatch core contract
status: backlog
source: "0.27 cut audit (2026-08-14), adversarially verified; captain directed filing"
started:
completed:
verdict:
score:
worktree:
issue:
---

Two prose repairs in skills/first-officer/references/fo-dispatch-core.md.

1. Line 195: the replay sentence references "s4", a private dev-workflow entity ID no installed reader can resolve. Rewrite to name the behavior: `spacedock gate prepare` is replay-idempotent. Keep the prohibition (no attempt counter, retry token, cache, or alternate authority).
2. Lines 197 and 199: the stale abbreviated --next envelope shape plus the sentence that exists only to retract it. Fold the canonical ready_gates shape into line 197 and delete line 199. The canonical rule at line 191 already governs.

## Problem

{Ideation fills this in.}

## Proposed approach

{Ideation fills this in.}

## Out of scope

Any behavior change. The canonical envelope rule and its Go emitters.

## Expected surface and tolerance

Estimate net LOC change: -N, across 1 file. No observable semantics change: contract prose only, meaning preserved from the canonical sentences.

## Acceptance criteria

**AC-1 - No shipped skill text contains an unresolvable entity-id referent for this rule.**
Verified by: grep for "s4" in skills/ returns no gate-replay match.

**AC-2 - The --next envelope shape is stated once, correctly.**
Verified by: the retraction sentence is gone and contractlint reference-closure tests stay green.

**AC-3 - The suite stays green.**
Verified by: go test ./... passes.

## Test plan

Prose edit plus contractlint. No behavior tests needed.
