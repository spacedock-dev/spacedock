# 0203 — FO Efficiency: shallow boot + lazy contract

**Milestone:** 0.20.3
**Status:** ideation complete, at the ideation gate (2026-06-13). The full spec lives in the j9 entity (`docs/dev/.spacedock-state/lazy-teamcreate-shallow-boot/`); this doc is the sprint index.
**Theme:** make the first officer cheap to boot and run.

## Why

Boot forensics on a live FO session (`/tmp/boot-analysis-spacedock-v1.md`) measured **~160k peak context and ~13.6 min** to reach an interactive greet — **with no team created and no worker dispatched.** The cost is structural, not a bug: the FO reads its entire contract (both reference files ~16k), the workflow README, and both mod files up front, then renders the full status table — most of it unused on a boot that never dispatches. Generation latency scales with loaded context, so the wall-clock is dominated by thinking at 100k+.

## Goal (success criterion)

An FO reaches interactive readiness — greet + state summary + *able to present a gate* — in seconds at **< ~60k** context, versus today's minutes at 126k+. Proven by a live FO-boot drive that observes correct behavior with the deferred modules unloaded — never by a grep over the restructured contract.

## Cost levers (ranked — why j9 is the backbone)

| Lever | ~boot cost removed | Needs the split? |
|-------|-------------------:|------------------|
| Lazy-TeamCreate (defer the team-mode prefix re-cache) | **~89k cache-creation** | no |
| Defer contract reads at greet | ~16k | yes (minimal) |
| Defer the human status-table render | ~8.7k | no |
| Defer mod-file reads | ~6.5k | no |

The 89k lazy-TeamCreate dwarfs the rest and sits on the critical path *before the dispatch gate* — so **`j9` (lazy-TeamCreate + shallow-boot) is the backbone.** The contract split is the enabling refactor for the one lever that needs it, and doubles as the contract-cleanup ask.

## The cut — reorganize the contract by *when* it is needed

**Boot-resident** (read on every FO start): contract-gate/startup, discovery, `status --boot --json`, status viewer, ID styles, single-entity mode, write scope, captain interaction, event loop, gate-presentation entry, clarification, working principles. Exactly enough to greet, report state, and present a gate.

**Deferred — loaded only when its phase begins:**
- **Dispatch/team module** (first dispatch): team creation, dispatch adapter, worker resolution, reuse conditions, standing teammates, degraded mode, context budget. ≈70% of `claude-first-officer-runtime.md` plus shared-core's reuse/standing-teammate sections. The bulk — and the biggest pure boot waste today.
- **Merge module** (terminal boundary): merge-and-cleanup, ship-local, teardown, mod-block enforcement.
- **Already lazy** (the precedent we extend): `present-gate`, `feedback-rejection-flow`, `using-claude-team` skills.

## Mechanism

The `spacedock:first-officer` skill reads only the boot-resident core at startup. The team/dispatch content folds behind the existing lazily-loaded `using-claude-team` skill (already invoked at first dispatch); the merge content becomes a lazily-loaded reference invoked at terminalization. No new pattern — we extend the one already in the codebase.

## Boot flow (the j9 shallow-boot)

contract-gate → discovery → `status --boot --json` → **greet and stop for input.** Team creation, mod-file reads, the dispatch/merge modules, and the human status table all defer to the moment they are needed.

## Test plan (honors the proof policy)

- **New live shared-runtime scenario `shallow-boot`:** the FO boots, greets, and presents a gate with the dispatch/merge modules *not* loaded — verified through the live `internal/ensigncycle` harness on durable behavior. The win is correct behavior at lower loaded context; behavioral and live, not a contract grep.
- **Regression:** existing `gate-guardrail` / `rejection-flow` / `merge-hook-guardrail` scenarios still pass — the deferred modules load correctly when a real dispatch/merge happens.
- **Structural guard** (in the allowed `internal/contractlint` quarantine): a reference-closure check that the boot-resident core has no dependency on deferred-only content. Structural, not prose-grep.

## Tasks

**`j9` is the backbone — one task, three phases.** Reshape the existing backlog entity; do not fragment it. The contract split runs first because it is the behavior-preserving enabler (and the contract-audit ask); the 89k lazy-TeamCreate is the headline lever.

**j9 — Lazy-TeamCreate + shallow-boot-then-greet:**
- **Phase 1 — contract structural split** (enabler + the "audit and cleanup the fo contract" ask): extract boot-resident vs deferred into a lazy dispatch ref + a lazy merge ref; slim the boot-resident core + the skill loader. Behavior-preserving. *Proof: existing `gate-guardrail`/`rejection-flow`/`merge-hook-guardrail` live scenarios still pass + a `internal/contractlint` reference-closure guard.* Spike verdict VIABLE, ~70% boot-read cut (full spec in the j9 entity).
- **Phase 2 — lazy-TeamCreate**: defer the `TeamCreate` call (the ~89k cache-creation) off the boot/greet path to first-dispatch-need. Needs no split.
- **Phase 3 — shallow-boot-then-greet**: greet off `status --boot --json`; defer mod-reads, the human status-table render, and the (now-split) deferred contract modules. Folds C3 (mod-defer) + C4 (status-render discipline). *Proof: new live `shallow-boot` scenario.*

**T3 — residual-prose audit + comm-officer polish** (file along, post-Phase-1; the cut-list does not exist until the split lands).

Boot-report habits (scope greps to headings; delegate bulk reconciliation reads to a subagent) tighten existing "Probe and Ideation Discipline" prose and ride along in Phase 1 / T3.

## Contract content fixes (captain audit, 2026-06-13)

Two fixes fold into Phase 1's contract cleanup, beyond the structural split:

1. **Drop the unnecessary `agents/first-officer.md` cross-reference** from `first-officer-shared-core.md` (line 3, "Keep aligned with…") — not load-bearing.
2. **Add a top Operating-principles (ethos) section** the shipped skill lacks today — its absence lets Codex drift from the `agents/first-officer.md` ethos. Combine the existing `## Working Principles` under it. Verbatim:

   > You are dispatcher and responsible for making sure the work is done by the crew. What awesome looks like for the crew:
   > - Begin with the end, be clear about the value.
   > - Do the hardest things first, de-risk when it is cheap.
   > - Communicate and act concisely, choose the simplest approach, JFDI.

This is boot-resident **guidance content** (behavior-shaping, not a testable AC): proof is the existing live scenarios still passing + review — not a "drift reduced" metric. The principles also govern *how* the contract is simplified (lead with value, hardest-first, simplest/concise).

## Out of scope (parked, not 0.20.3)

- **p2 / vc** — `spacedock pr complete` + `reconcile --act`: the binary-simplification line (higher ROI, heavier lift) → 0.20.4.
- **xp** — cross-session FO↔Commander comms: the coordination infra that makes multi-FO safe. Its own track.
- **ey** — proof-policy port to shipped scaffolding: adjacent, separate.

## Operational landing

This doc lives in the main repo (collision-free). Filing T1 and reshaping `j9` into `docs/dev/.spacedock-state` are coordinated with the active Commander session via path-scoped commits (disjoint entities → safe under the multi-writer protocol).
