---
id: 2wm8dw4vgysthafbw1kmj3jd
title: Captain decisions pool into one non-blocking gate instead of halting the drive
status: backlog
source: "Commander session driving 0260-proportionality, 2026-07-20. Three captain decisions arose mid-drive (az's consent-gated Edit D; 841's uncovered-token ruling; 841's tolerance reconfirm-vs-park). Each halted the FO turn: ask in chat, wait, resume. All three were low-bandwidth rulings among 2-4 named options with an FO recommendation attached — not design conversations. Captain direction: make clarification non-blocking, and pool decisions into one lightweight gate once 3k and the subspace wiring are done."
started:
completed:
verdict:
score:
worktree:
issue:
---

Captain decisions raised while driving a sprint should accumulate into one pooled, non-blocking review rather than interrupting the drive once per decision.

## Problem

Driving a sprint surfaces captain-reserved decisions at unpredictable moments. Today each one is a hard stop: the FO composes a question, halts, and waits for an answer before continuing — even when other members are independently drivable and even when the decision is a one-word ruling among options the FO has already researched and recommended.

Observed in the 0260 Commander drive (2026-07-20), three in a single session:

- az's Edit D — a consent-gated edit whose gate record contradicted itself (approval reason said "awaits its own explicit yes/no"; a later block recorded "approved" via a blanket "agree with all recommendations" ruling plus a flag-if-not). Needed one explicit yes.
- 841's uncovered runtime tokens — whether to delete a check and record an honest gap, or retain a phrase check nothing exercises. A ruling, with the FO's recommendation already argued.
- 841's tolerance breach — reconfirm / re-scope / park, where the prior reconfirm's justification had been invalidated.

Each cost a full round trip. None needed more than a selection. Meanwhile the sprint's own rule requires a RECORDED decision before further repair dispatch, so the decisions cannot simply be skipped — they must be cheap, not absent.

The cost is not the captain's time answering; it is the drive stopping while the answer is pending, and the captain being interrupted N times instead of once.

## Proposed approach

Two properties, in this order — the first is the load-bearing one and does not depend on the second:

1. **Non-blocking.** A pending captain decision blocks only the member that needs it. The FO records the decision as pending, holds that member from advancing, and keeps driving every unaffected member. This is already the contract's stated posture ("keep dispatching other ready entities when one blocks"); what is missing is a durable place to park a pending decision so the drive does not have to hold it in conversation.

2. **Pooled.** Pending decisions accumulate and are presented together at a natural boundary as ONE captain-facing review with N annotated items, each carrying the FO's recommendation and grounds, rather than N separate interruptions.

## Dependencies

- **3k (gate-resolution recorder).** 0260 used its frontmatter convention as a dry-run; the binary is not shipped. The convention already proved the consume path works end to end: a closed attempt with `application: {action, target-stage, state: pending}` is applied once through the normal transition path and stamped `consumed`, with a drift check before consumption. A pooled decision gate is the same record shape with more than one open attempt at a time.
- **Subspace wiring.** Multi-artifact briefing packages are designed but absent from subspace-tui 0.9.0's surface, which forced single-file gate-summary+entity concatenation during 0260 shaping. The working float ritual is the `subspace-r-working-copy` skill's `scripts/review-local-zellij` (exact-tip local build, probe with a throwaway file first). A pooled gate is inherently multi-artifact — one briefing, N decision items — so it wants the packaging that is currently missing.

## Out of scope

Making decisions the captain has NOT reserved into gates. This is about the latency and batching of captain-reserved rulings, not about widening what requires captain sign-off. Also out of scope: any mechanism that infers consent from silence or from a blanket ruling — 0260 recorded a live instance of that going wrong, and the pooled gate must make each item's answer explicit and individually attributable.

## Acceptance criteria

Ideation fills these in. The value to measure is the drive not stopping: a run in which a captain decision is raised and other members continue advancing while it is pending, against today's baseline where the turn halts.
