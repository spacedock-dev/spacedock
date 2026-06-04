---
id: gq9g4vrz03kgd8w46cvf09k7
title: Live-scenario coverage for the non-happy feedback-rejection paths
status: ideation
source: "captain (2026-06-04) — a9 detached audit surfaced that the feedback-rejection non-happy-path guarantees are guarded only by review + the single-cycle live scenario. Investigation: the old tests/test_rejection_flow.py drove 2 full cycles + reviewer-reuse; the current Go rejection-flow scenario simplified to a single route-back; and NEITHER era ever drove the 3rd-cycle escalation or the budget-probe fail-safe. Use the existing prose-based shared-scenario runner to exercise these."
score: "0.30"
started: 2026-06-04T20:04:39Z
completed:
verdict:
worktree:
issue:
---

The feedback-rejection procedure carries non-happy-path behavioral guarantees that no standing test exercises. The a9 detached audit (commit 98629283) demonstrated each is gut-able while the suite stays green, because the static oracles assert only substring-presence of the prose, not the behavior. Add live-scenario coverage via the existing prose-based shared-scenario runner.

## The gaps (concrete)

- **3-cycle escalation.** The contract says: route findings back to implementation on cycles 1 and 2, and **on cycle 3 escalate to the human instead of auto-bouncing a 4th time**. Drift → the FO loops forever (reject → re-implement → reject …), burning tokens. NEVER tested: the live `rejection-flow` scenario drives one cycle; the old `tests/test_rejection_flow.py` drove two (impl → val REJECTED → feedback → cycle-2 impl-fix → cycle-2 reval via reviewer-reuse) but never a third. Highest blast-radius (runaway cost).
- **Budget-probe fail-safe (reuse condition 0).** Before reusing the kept-alive ensign for rework, the FO must consult the context-budget probe and **fresh-dispatch if the ensign is over budget or the probe is unavailable**. Drift → reuse of a blown-context ensign → degraded/failed work. Never tested either era.
- **Coverage regression in the Go port.** The current Go `rejection-flow` scenario *starts* from an already-rejected report and asserts a single route-back — it dropped the live 2-cycle trajectory + the SendMessage reviewer-reuse the Python test actually exercised. Restoring that multi-cycle trajectory is part of this task.

## Why this, why now

This generalizes across the 0.19.6 decomposition line (t3/a9/wm/p2): moving contract prose into lazy skills leaves these behavioral guarantees guarded only by review + sparse live scenarios. a9 did not regress it (the prose was never standing-behaviorally-tested) — but the guarantees are real, with real failure modes, so the lever is expanding live-scenario coverage, prioritized by blast-radius.

## Direction (for ideation to flesh out)

- Add a `feedback-3-cycle-escalation` shared/live scenario using the existing prose runner (`Use $spacedock:first-officer`): a fixture seeded with two prior rejections so the next is the 3rd; assert the durable end-state is an **escalation to the human** (an escalation marker / no 4th auto-bounce), graded on durable state, never transcript phrasing (per scenario-testing-principles).
- Evaluate restoring the full multi-cycle trajectory + reviewer-reuse the Go port simplified away.
- Decide whether the budget-probe fail-safe warrants its own scenario or is better served by the binary-gate task (see `feedback-guarantee-binary-gate`) — some of these guarantees may be cheaper to enforce in the binary than to drive live.
- Triage by blast-radius; not every non-happy path earns a (costly) live scenario.

## Notes

Sibling: `feedback-guarantee-binary-gate` (the "promote prose → binary code gate" lever — the stronger fix where a guarantee is mechanizable). Provenance: a9 (`feedback-rejection-flow-skill-extraction`) detached audit, 2026-06-04.
