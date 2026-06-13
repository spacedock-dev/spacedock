---
id: xf7fft1hnj51eq7kagsc9833
title: On-demand polish — drop the standing-teammate lifecycle, move usage prose into the mod
status: backlog
source: "captain (2026-06-13, this session) — the standing-teammate lifecycle (discovery pass / lazy-spawn / declaration / team-scope teardown / first-boot-wins) is ~4 contract subsections of maintenance surface for an infrequently-used polisher; amortization doesn't pay for infrequent use. Captain chose approach A: on-demand one-shot polish dispatch, and 'the mod can add the prose for the standing team member' — feature-specific usage prose lives in the mod, the contract keeps only a generic hook. Taken into 0203."
started:
completed:
verdict:
score:
worktree:
issue:
sprint: 0203-fo-efficiency
---

Replace the standing-teammate lifecycle with on-demand one-shot polish, and move the feature's usage prose out of the FO contract into the mod that declares it.

## Problem

Supporting the standing prose-polisher (`comm-officer`) costs ~4 contract subsections — the shared-core `## Standing Teammates` concepts plus the runtime's discovery-pass / lazy-spawn / declaration-and-routing-mechanics (all relocated into `claude-fo-dispatch.md` by j9's P1 split) — plus the `using-claude-team` lifecycle notes. All of that machinery exists to buy ONE thing: amortization (keep the polisher resident so its skill loads once across many polishes). For an infrequently-used feature that trade does not pay — we maintain a long-lived, team-scoped, first-boot-wins, lazy-spawn, teardown lifecycle to save a skill-load that rarely happens.

## Proposed approach (seed — ideation designs it)

**Approach A — on-demand one-shot.** When the FO (or an ensign) has a deliberate draft to polish, it dispatches a one-shot polish agent only then; there is no resident teammate, no team-scope lifecycle. You pay one spawn per polish, which is rare.

**Mod-owned usage prose (the captain's key point).** The feature-specific prose moves OUT of the contract and INTO the `comm-officer` mod that declares it. The mod already carries the agent prompt (`## Agent Prompt`) and a routing block (`## Routing Usage` surfaced by `dispatch show-standing`/`dispatch build`); extend it to own the on-demand *usage* prose too (when/how to invoke polish, the four polish modes, the boundary rules). The FO/ensign contract keeps only a minimal generic hook: "registered polish mods declare on-demand polishers; the dispatch helper surfaces their usage; route deliberate drafts to them per the mod's declaration."

**What the contract sheds:** the discovery pass, the lazy-spawn pass, the standing-declaration layout, team-scope teardown of standing members, first-boot-wins — gone. **What stays / changes:** a binary path that emits a one-shot polish dispatch from the mod's prompt (repurpose `spawn-standing`'s prompt-extraction into a one-shot `dispatch polish`/spec-emit); the mod becomes self-describing.

## Out of scope

- The behavioral polish QUALITY (the elements-of-style skill usage) — unchanged.
- T3's FO-contract prose audit (sibling; this is a specific contract-surface reduction, T3 is the general audit).
- Approach B (spawn-on-first-route, keep residency) — recorded as the fallback if a live-polish-latency concern surfaces at ideation, but A is the chosen direction.

## Sequencing (gated, like T3)

This edits `claude-fo-dispatch.md` — the file j9's P1 split created and moved the standing-teammate prose into — plus the shared core, the `comm-officer` mod, the dispatch binary, and the `using-claude-team` skill. So it **sequences AFTER j9 P1 merges** (same ordering hazard as T3 / sr's prose retarget): dispatching prose work against today's `main` line-anchors would edit lines that no longer live there.

## Scaffolding guardrail

Touches shipped scaffolding (`skills/first-officer/references/`, the `_mods/comm-officer.md` mod, the dispatch binary, `skills/using-claude-team/`) — a dispatched worker under test, never an FO-direct edit.

## Acceptance criteria (seed — ideation defines external proofs; NO contract prose-grep)

The proof that the lifecycle is gone is behavioral, not a grep over the contract: a live drive where the FO composes a deliberate draft and polishes it via a one-shot on-demand dispatch, with durable on-disk state showing NO standing teammate in the team `config.json` across the session until a polish is actually needed (and the one-shot polisher present only transiently when polishing). The existing live scenarios stay green and the offline gate exits 0. The mod-owned usage is exercised by actually invoking polish per the mod's declaration (the FO/ensign reads the usage from the mod, not the contract). Ideation defines the concrete scenario + the before/after wording.
