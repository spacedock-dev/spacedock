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
gates:
    version: 1
    records:
        - id: gate:71btbxdrken4kdmfsk0vptav:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:71btbxdrken4kdmfsk0vptav-backlog-1
              briefing:
                id: briefing:71btbxdrken4kdmfsk0vptav:backlog:attempt-1:revision-1
                digest: sha256:665a1ebd3a9258ea6dad9a58d7a1494000d82e6a2ed6b2d0a5f000c540943df5
                request-digest: sha256:b7a73f6a9d9f568a7d9c13110daf4286d155d40a011052c0b330e3434fe76b82
                room-ref: ./trim-dispatch-core-stale-prose/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:71btbxdrken4kdmfsk0vptav:backlog:1
                briefing: briefing:71btbxdrken4kdmfsk0vptav:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T02:53:55.975611Z"
                decision: approve
                reason: 'Captain ruling 2026-08-14 (dispatch them): approved into ideation'
              application:
                target-stage: ideation
                state: pending
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
