---
id: 8dppn2p2e30vc8fewnda40m6
title: Raise the fo-gate-lifecycle byte cap to home the three boot-resident gate fns
status: backlog
source: "Follow-up adopted at the de-lecture-and-defer-fo-contract (wn) ideation gate, 2026-08-24: the gate-first engage path means gate.ac-cross-check, gate.assemble-verdict, and feedback.route cannot ride the dispatch-core read; their only always-loaded-at-gate home is fo-gate-lifecycle/SKILL.md, 55 B under its captain-set 7,700 B cap"
started:
completed:
verdict:
score:
worktree:
issue:
---

Move the three gate-path prose functions (gate.ac-cross-check, gate.assemble-verdict, feedback.route - 2,653 B) out of the boot-resident shared core into fo-gate-lifecycle/SKILL.md, which requires raising that file's contractlint cap (fo_function_reference_invariant_test.go, currently 7,700 B) to ~10,400 B. A captain cap decision on its own evidence - deliberately excluded from wn to keep that task a clean de-lecture plus one provably-safe move. Recovers ~1,020 tok per greet on top of wn's ~1,200. Depends on wn landing first (its moved-section test machinery is the AC-3 pattern to reuse; its pointer conventions apply).

## Problem

{Ideation fills this in. Seeded: the three fns fire only at gates, but the gate-first engage path loads fo-gate-lifecycle, not dispatch-core - so their only safe deferral target is the capped skill file. The cap exists to bound gate-time context; raising it trades gate-time bytes for boot-time bytes, net-positive for every session that greets without gating.}

## Proposed approach

{Ideation fills this in - lightly. The move + cap bump + wn's moved-section table extended to these three + pointers per wn's convention.}

## Out of scope

Any prose change beyond relocation. wn's deletions (landed). New triggers or files.

## Expected surface and tolerance

Estimate net LOC change: ~0 net across 4 files (shared-core, fo-gate-lifecycle SKILL, contractlint cap test, boot_resident_closure_test). Ideation refines.

## Acceptance criteria

Seeded; ideation refines lightly.

**AC-1 (value) - A greet-and-stop boot sheds the three gate fns' bytes on top of wn's reduction.**
Verified by: the same stack-base-vs-tip A/B byte accounting wn uses; fails if the relocation leaves restating pointers or the sections boot-resident.

**AC-2 - No gate path can reach a moved fn before fo-gate-lifecycle is loaded.**
Verified by: wn's moved-section table pattern plus the gate-first engage argument re-checked against the shipped contract at implementation time; the live gate journey is the behavioral form.

## Test plan

{Ideation fills this in. The cap bump is a one-line contractlint change the captain approves at the gate; lanes per the contract-surface rule.}
