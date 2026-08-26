---
id: mex0e7x7tw00ggpdd0773nd8
title: Patch-line stable delivery via a merge-commit advance
status: backlog
source: "Deferred from patch-release-line-support's lean descope (captain ruling 2026-08-25, 'lean'); direction and evidence from the 2026-08-25 independent staff review of that task's ideation"
started:
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
archived: 2026-08-26T04:54:08Z
---

Ship an old-line patch release (a v0.X.Z while main is a later line) to stable-channel users. Deferred until a patch is genuinely demanded: the lean cut blocks old-line cuts pre-publication, so `stable` never diverges and delivery machinery is not yet needed.

## Problem

{Ideation fills this in. Seeded: once a patch commit moves `stable` off main's history, every later latest-line release push is non-fast-forward and the stable channel freezes. The lean cut prevents the divergence instead of handling it; this task builds the handling when a real patch need arrives.}

## Proposed approach

{Ideation fills this in. Recorded direction from the staff review: a merge-commit advance — CI runs `git commit-tree "$RELEASE_COMMIT^{tree}" -p "$STABLE_SHA" -p "$RELEASE_COMMIT"` and plain-pushes the result. Always a fast-forward by construction across every line shape; no force authority; a concurrent-run race fails closed as a plain non-FF rejection; the consumer-cache risk of a force-moved ref is moot. Cost: stable's tip SHA no longer equals the tagged SHA (the tree is byte-identical and the second parent IS the tagged commit, so provenance survives). The version gate (a stable-advance-decision comparator over shipped ComparePreVersion) is still required. Review-verified fact: refs/heads/stable has no branch protection today, so a force alternative would normalize existing capability, not add new capability — the merge-commit form is preferred anyway for its fail-closed race behavior.}

## Risk evidence

{Backlog: the 2026-08-25 staff-review spike matrix (7-test bare-origin fixture: FF/non-FF/lease semantics) and its merge-commit evaluation decide the direction. Ideation must spike: both hosts' installed-plugin update paths across a stable ref whose tip is a merge commit; and the release ritual's patch-branch steps (branch off the prior stable tag, not origin/main).}

## Out of scope

The pre-publication old-line block, the main-stamp move, and the pre-stamp automation (patch-release-line-support ships them). Marketplace display fields. The e2e-gate waiver mechanics.

## Expected surface and tolerance

Estimate net LOC change: +200, across 7 files. {Backlog seed; ideation refines — the replay harness is the bulk.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - After a patch cut off the stable line and a later latest-line cut, stable-channel users receive both releases in order, with no force-push in the release workflow.**
Verified by: {ideation refines — seed: bare-origin replay of the real steps across the two-cut sequence asserting stable's tree equals each tagged tree in turn and the push is never forced; falsifying edit: replace the merge-commit advance with a plain ref push — the second cut reds non-FF.}

## Test plan

{Ideation fills this in. Seeded: extend the stable_advance replay harness; host update-path spike; .github/** is release machinery — the detached adversarial audit applies. All comments and user-facing doc text follow ASD-STE100.}

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}

## Superseded

Captain ruling 2026-08-26: no patch-release machinery. The stable-channel line-crossing is handled as a documented one-time manual force-with-lease push at the next latest-line cut (docs/releasing.md, commit 8ccebceb5). The merge-commit design and the staff review's comparison stay recorded here.
