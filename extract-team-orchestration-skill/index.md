---
id: zdwt24era9rdg622e683x91v
title: Extract generic Claude-team orchestration out of the FO contract into a reusable `using-claude-team` skill; the FO contract keeps only spacedock-specific decisions
status: ideation
source: captain (2026-06-03) — "can we move the team related thing into a separate skill? 10 steps, error-prone, there must be an easier way or some missing abstraction"
score: "0.30"
worktree:
started: 2026-06-04T05:24:55Z
completed:
verdict:
issue:
---

The Claude-team lifecycle (~10 steps) is inlined in the FO contract (`first-officer-shared-core.md` + `claude-first-officer-runtime.md`), mixing **generic Claude-team-harness discipline** with **spacedock-specific decisions**. `claude-first-officer-runtime.md` is the single largest contract file (~8.8k tokens) and most of it is team mechanics. The ceremony is error-prone — this session the FO tripped the TeamCreate sequencing rule, the deferred-tool `SendMessage` hop, and the supersede-shutdown dance.

## Problem

The generic team-orchestration discipline has zero spacedock content but is welded into the spacedock FO contract, so: (a) it can't be reused by any other Claude multi-agent workflow; (b) it bloats the FO boot read; (c) the error-prone ceremony is spread across a long contract rather than concentrated in one tested place.

## Proposed approach

Two-layer split:

**Layer 1 → a `using-claude-team` skill (generic, reusable, zero spacedock content):**
- TeamCreate (naming convention, returned-name handling, fresh-suffix recovery ladder)
- the **Awaiting-Completion discipline** (the don't-prematurely-teardown rules — the most error-prone, most generic block)
- supersede-shutdown + terminal teardown
- Degraded Mode (triggers / effects / cooperative-shutdown sweep)
- SendMessage routing, including the deferred-tool `ToolSearch select:SendMessage` hop
- standing-teammate spawn pattern

**Layer 2 → the FO contract stays, but thin (spacedock-specific decisions only):**
- dispatch via `spacedock dispatch build`
- reuse conditions (stage-model match, worktree-state routing — spacedock stage semantics)
- reconcile Class A/B → teardown mapping
- the `spacedock dispatch context-budget` probe

The FO contract `@`-references the skill for the lifecycle and keeps only the spacedock decision points.

**Load-bearing nuance — why this is a SKILL, not a binary command:** team ops are Claude *tools* (`TeamCreate`/`Agent`/`SendMessage`), not shell commands, so the binary cannot own them (it can't call Claude tools). The binary already owns everything shell-able (dispatch build, reconcile, context-budget). So the missing abstraction is a *skill* sitting beside the existing binary helpers — NOT a new binary command. Do not attempt to push team-tool-orchestration into the binary.

**Cross-runtime:** either a Claude `using-claude-team` skill + a Codex analog, or one `using-agent-teams` skill with per-runtime sections (the codex FO runtime adapter has the parallel structure).

## Out of scope

- Moving team-tool orchestration into the binary (impossible — see the nuance above).
- The spacedock-specific decision logic (reuse conditions, reconcile mapping) — stays in the FO contract.
- The ensign-side completion-signal mechanics (related but a different surface).

## Acceptance criteria

**AC-1 — A `using-claude-team` skill carries the generic team lifecycle, and the FO contract references it instead of inlining.**
Verified by: an oracle asserting the generic blocks (TeamCreate/recovery, Awaiting-Completion, supersede/teardown, Degraded Mode, SendMessage routing) live in the new skill file AND no longer appear inlined in `claude-first-officer-runtime.md`; the FO runtime file's token/line count drops materially.

**AC-2 — The generic/specific boundary is clean.**
Verified by: an oracle asserting spacedock-specific decisions (reuse conditions referencing stage model/worktree, reconcile Class mapping, `dispatch build`/`context-budget`) remain in the FO contract and did NOT leak into the generic skill.

**AC-3 — No behavior regression on the decomposed contract.**
Verified by: a live FO drives a full workflow cycle (dispatch → gate → merge) on the decomposed contract and completes cleanly — the Phase 0.A bar (load-bearing meaning preserved across a real drive).

## Test plan

- Instruction-text oracles for AC-1/AC-2 (the proven contract-oracle pattern; extend `internal/hostneutrality` / `skills/integration`). Cost: low.
- One live FO cycle for AC-3. Cost: medium.
- High-stakes surface (the FO operating contract itself) → detached adversarial audit before merge.

## Notes

- Lived friction this session: TeamCreate sequencing trip, deferred `SendMessage` hop, supersede-shutdown dance — concrete evidence the abstraction is missing.
- Sibling to `ensign-contract-dev-leakage` (ep) — both decompose the contract along generic-vs-specific lines. Sibling to `dispatch-build-json-ergonomics` (mt) and `launcher-binary-path-passthrough` (fc) as a contract/binary ergonomics cluster.
- Not test-themed → a contract-decomposition-track candidate (next sprint, or alongside 0.19.5 if the captain wants it now).
