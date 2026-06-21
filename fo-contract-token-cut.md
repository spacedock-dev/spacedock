---
title: Convert remaining FO contract cuts into pseudo-code capability bodies
status: backlog
source: The 2026-06-15 fo-contract-token-cleanup proposal classified candidate cuts by adversarial reasoning, then revised the default-path total to ~420 tokens after #396 retired RT-4/RT-2 and demoted UCT cuts to legacy-only. After #418 and the state.commit follow-up, the remaining objective is not just token recovery: contracts should read like pseudo-code capability bodies, with prose reserved for fuzzy judgment and probe-backed quirks.
score: 0.5
sprint: 0221-layered-fo
sprint-readiness: defer
issue:
id: y2r7ew51xqs6q3avsb6mcaka
group: cleanup
---

Rebase `docs/dev/_proposals/fo-contract-token-cleanup.md` against current `main` and apply the still-valid cuts as pseudo-code-shaped contract bodies. The deliverable is not only fewer tokens: shared/runtime contracts should prefer `«fn»` sections with compact `guard` / `effect` / `done-when` / `block` / `→` lines, and use prose only where the obligation is fuzzy, judgment-owned, or probe-specific. Each risky cut is empirically confirmed by a no-guidance-control micro-test before it lands.

## Problem

The proposal's original `~638` target is stale. The proposal itself now says the default path is ~420 tokens, RT-4/RT-2 are retired, UCT cuts are legacy-only, and UCT-5 needs re-derivation. Since then, `codex-pi-runtime-binding-block-cleanup` (#418) moved Codex/Pi adapters toward capability binding blocks, and `state-commit-contract-token-followup` trimmed the largest new shared-core state body by replacing mechanics prose with a compact `«state.commit»` body.

The remaining problem is broader than shaving words. The FO contracts still mix callable capability bodies with narrative lifecycle prose. The target shape is pseudo-code: the core names what gets invoked and when, runtime files bind how the host realizes each `«fn»`, and prose is reserved for fuzzy judgment or probe-backed host quirks.

## Proposed approach

Rebase first, then apply:

1. **Reconcile scope against current main.** Remove retired RT-4/RT-2 targets, keep UCT items legacy-only, re-derive UCT-5 before applying, and do not redo work already landed by #418 or the `«state.commit»` follow-up.
2. **Prefer pseudo-code bodies.** Convert surviving shared-core/runtime targets into compact callable bodies or binding bullets before reaching for explanatory prose. Good shapes include `guard`, `effect`, `done-when`, `block`, and `→ shipped/prose/runtime-binding`.
3. **Validate risky removals.** Use the proposal's no-guidance-control micro-test for clauses whose removal could change behavior. For a candidate clause: sample FO behavior N>=5 on the smallest realistic exercise that exposes the clause's job, in two arms — WITH the clause and WITHOUT it — in the real surrounding contract; read behavior by hand; treat variance as signal.
4. **Ship through the dispatched-worker path.** The contract files are shipped scaffolding, so the trimmed files land via a dispatched worker in a worktree, with live shared scenarios green on relevant hosts, token/line deltas reported, and detached adversarial audit of the diff.

## Deliverable

The trimmed FO-contract files committed, with: measured token/line deltas, a list of proposal items applied/retired/deferred, a per-risky-cut micro-test verdict, and a short statement of how the resulting text follows the pseudo-code contract principle in `docs/runtime-support.md`.

## Out of scope

- A general eval framework (the micro-test is a thin per-clause check, not a product — YAGNI until a future batch justifies one).
- Reworking already-landed #418 Codex/Pi binding-block shape or the `«state.commit»` compact body except to avoid conflicts.
- The judgment-audit generalization over the full routing table (rides w4 / 0205).
- Net-new cut hunting beyond the proposal's verified list (re-testing the 13 keeps IS in scope; finding fresh cuts is not).

Source: the 2026-06-15 token-cleanup proposal + superpowers v6 (writing-skills "Micro-Test Wording Before Full Scenarios", positive-instruction-redesign-design, strict-cost-sdd-design).
