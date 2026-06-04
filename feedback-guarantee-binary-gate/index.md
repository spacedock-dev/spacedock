---
id: xae5tx4hhyce916x034y3q9x
title: Evaluate promoting feedback-rejection guarantees from prose to binary-enforced gates
status: backlog
source: "captain (2026-06-04) — the a9 detached audit's dominant finding is that FO-behavioral guarantees living as contract prose (3-cycle escalation, budget-probe fail-safe) are immune to static tests and only sparsely covered by live scenarios. Where a guarantee is mechanizable, enforcing it in the binary (a guard / a tracked field) is the stronger fix than a live scenario — the 'code gate over prose-only rule' working principle. A third decomposition lever next to ceremony->binary (p2) and judgment->skill (t3/a9)."
score: "0.24"
started:
completed:
verdict:
worktree:
issue:
---

The audit exposed that some FO-behavioral guarantees are inherently un-static-testable (a text oracle can't prove the FO obeys "escalate on cycle 3"), and live-scenario coverage for them is expensive. But several of these guarantees are **mechanizable** — the binary could track and enforce them deterministically, the way `status --set`/`--archive` already refuse a terminal transition without `pr`/`mod-block`. That eliminates the body-vs-label gap entirely (the guarantee is no longer prose).

This is the **guarantee → code-gate** lever — a third token-efficiency/robustness lever alongside the 0.19.6 line's ceremony→binary (p2) and judgment→lazy-skill (t3/a9).

## Candidates to evaluate

- **3-cycle escalation.** Could be a binary-tracked feedback-cycle counter that the system enforces — refuse a 4th auto-bounce / force an escalation marker on the 3rd rejection. Eliminates the infinite-loop regression class. Analogous to the existing `mod-block` terminal guard.
- **Budget-probe fail-safe.** The probe is already a binary command (`spacedock dispatch context-budget`). Question: should the system *force* the consult (e.g. the dispatch/reuse path refuses to reuse an over-budget member) rather than instructing the FO in prose to consult it?

## Scope of this task

Ideation/decision: for each candidate, decide whether it is genuinely mechanizable (a deterministic guard over real state) vs inherently a model judgment (stays prose + live scenario). Where mechanizable, scope the binary change; where not, defer to `feedback-nonhappy-live-coverage`. Output is a decision + (for the mechanizable ones) a concrete guard design — not docs-only; per dev policy, a decision with nothing shipped belongs in the roadmap, so this either produces a guard worth implementing or records the determination that these stay prose-guarded.

## Notes

Sibling: `feedback-nonhappy-live-coverage` (the live-scenario lever for the inherently-behavioral guarantees). Provenance: a9 (`feedback-rejection-flow-skill-extraction`) detached audit, 2026-06-04.
