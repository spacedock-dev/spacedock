---
id: j903f6f1vgckk3kj6j6zbmt3
title: Lazy-TeamCreate + shallow-boot-then-greet — reach interactive readiness fast
status: backlog
source: "captain (2026-06-04) — boot analysis (this session) measured boot at ~7:36, ~80% model-compose; TeamCreate is the single largest write (89k cache_creation — the whole prompt prefix re-cached to the 1h cache when team-mode activates) and sits on the critical path BEFORE the gate that decides whether dispatch even happens."
score: "0.34"
started:
completed:
verdict:
worktree:
issue:
---

The FO boot front-loads all deep work before greeting the captain: contract gate → sync/build → state pull → status → README → skills → mod read → TeamCreate → reconcile → entity reads → a long compose → greeting. Most of that is not needed to be interactive. The boot analysis measured the single biggest write as TeamCreate (89k cache-rewrite), and it fires before the captain has decided whether to dispatch at all — wasted entirely if they steer elsewhere.

## Problem

- Time-to-interactive is gated on the full deep-boot, not on a fast orient pass.
- TeamCreate's 89k cache-rewrite is on the critical path before the dispatch gate — paid even when no dispatch follows.
- The largest wallclock sink is model-compose over an ever-growing context (entity reads + a full multi-entity plan before the first word), not tool latency.

## Proposed approach

Split boot into a fast **orient + greet** critical path and a lazy **act** path:

- **Critical path (target: interactive in <30s):** contract version gate (safety) + `status --boot` (one call: mods / orphans / dispatchables / state-backend + the split-root halt-gate) + a `status` overview + read the handoff. Then a TIGHT greeting ("here's the state, what's the move?") — NOT a full multi-entity plan.
- **Lazy / deferred (on first need):** **TeamCreate** → defer to the first team-mode dispatch (consistent with standing-teammate lazy-spawn; "TeamCreate before any Agent" still holds at first dispatch). Reconcile → background after the greeting. Entity-body reads → only when acting on that entity. Deep planning → after the captain picks a thread.
- **Keep on the critical path (cheap + safety-critical):** the contract gate and the split-root halt-gate. These are not deferred.

## Acceptance criteria (ideation formalizes; per the proof-policy the proof is a MEASUREMENT, not a presence oracle)

**AC-1 (measured) — a cold boot reaches an interactive greeting within the target, with TeamCreate NOT on the critical path.**
Verified by: a measured boot trace (the per-turn jsonl token/wallclock parse from the boot analysis) on a real cold boot showing (a) time-to-greeting under the target budget, and (b) no TeamCreate call before the greeting — a before/after vs the current boot. (NOT a check that the contract text says "defer TeamCreate.")

**AC-2 (behavioral) — deferring TeamCreate does not break team-mode dispatch.**
Verified by: a live drive (or the existing live-cycle smoke) showing the first team-mode dispatch lazily creates the team and dispatches correctly — the sequencing rule holds at first dispatch.

## Test plan

- The measured boot trace (jsonl parse) — the AC-1 before/after. Needs a real cold boot.
- AC-2 rides the live-cycle smoke (first dispatch lazily creates the team).
- High-stakes (the FO's own Startup procedure + the runtime "TeamCreate at startup" rule) → detached adversarial audit before merge.

## Notes

Provenance: boot analysis (this session) — 7:36 boot, TeamCreate 89k cache-rewrite (the biggest single write), on the critical path before the dispatch gate. The deferred-TeamCreate note is already in the binary-simplification roadmap. Caveat: some boot steps are SAFETY (contract gate, split-root halt) — keep them on the critical path; only the expensive informational + team work defers.
