---
id: jsp6penfxp3h5n9z08vy0294
title: Single-entity (-p) feedback-cycle reviewer continuity — reusable kept-alive reviewer under lazy-TeamCreate
status: backlog
source: "j9 validation cycle 1 (2026-06-13). AC-3 rejection-flow surfaced a PRE-EXISTING contradiction: single-entity mode mandates bare dispatch, so the #141 keepalive reviewer cannot be reused across feedback cycles; P2 lazy-TeamCreate unmasked it (the eager 'create a team at startup' instruction was accidentally masking it). j9 landed via option (b) — the rejection-flow assertion was corrected to the bare contract. THIS is the spun-off option (a): the real continuity fix. Deterministic repro committed on j9's branch (da8e5912). Sibling: e3z (bare-mode-coverage-baseline)."
started:
completed:
verdict:
score:
worktree:
issue:
sprint:
---

Give single-entity (`-p`) feedback cycles real reviewer continuity: let the kept-alive cycle-1 reviewer be reused on cycle-2, so the #141 keepalive guarantee holds in `-p`, not just team mode.

## Problem

The contract has a latent contradiction, unmasked by j9's P2 lazy-TeamCreate:
- `single-entity → skip team creation → bare-mode dispatch` (since the original vendoring, `83c73494`).
- The #141 keepalive contract keeps the feedback-to reviewer alive across cycles for continuity — which requires a REUSABLE (team-mode) reviewer (reuse-condition-1: "not in bare mode").

So in `-p` mode the cycle-1 reviewer is bare → unreusable → cycle-2 gets a fresh reviewer that re-reads cold, losing cycle-1 context. Pre-P2 this was accidentally masked (the salient "create a team at startup" instruction led the model to create a team in `-p` despite the skip-clause); P2 removed that instruction and the latent bug surfaced. **Codex is unaffected** (it reuses via `send_input` to a persistent thread — no team_name lifecycle), so this is Claude-specific: team-gated reuse-condition-1 meeting single-entity bare mode.

There is also a non-determinism wrinkle observed at validation: the single-entity rejection-flow produced 2 spawns on one run and 0 spawns on another (the FO reusing the implementation ensign *through* validation, which violates `fresh: true`). The single-entity feedback behavior is under-specified, not merely "bare = fresh."

## Proposed approach (seed — ideation designs the fix)

Option (a) from the j9 recon: **let the `-p`/single-entity feedback flow create a team at FIRST dispatch** so the reviewer is reusable across cycles. This preserves the interactive boot win (the 45.9k/36.5k greet never dispatches, so an interactive greet-and-stop still creates no team — this is a team at the first `-p` dispatch, NOT an eager boot team). Make single-entity cycle-2 deterministic (always reuse the kept-alive reviewer when present; never reuse the impl ensign as its own validator). Blast radius: FO single-entity dispatch policy + the rejection-flow assertion (which j9 corrected to the bare contract — this task would move it to expect reuse in `-p`).

## Out of scope

- The j9 restructuring/test-correction (already landed via option b).
- The broader bare-mode coverage modeling — that is `e3z`; coordinate (this is the concrete reviewer-continuity instance of e3z's gap).

## Acceptance criteria (seed — ideation defines external proofs)

A live single-entity rejection-flow drive that asserts the cycle-1 reviewer is REUSED on cycle-2 (1 reviewer, by agentId) — the inverse of j9's corrected assertion — plus the existing offline repro (`da8e5912`) flipped to a team-mode-continuity control. Deterministic across runs (no 2-vs-0-spawn flake).
