---
title: "Merge-boundary constraints surface at dispatch time instead of loading at the terminal boundary"
status: backlog
source: "Principle-derived contract review, 0250 Commander session 2026-07-07 (captain-requested filing). Finding: the deferred-load architecture optimizes tokens cheapest-first but defers RISK discovery — fo-merge-core (PR guardrails, rebase-before-hook, PR-body shapes sourced from stage reports) loads only at the terminal boundary, after implementation and validation choices are locked. A merge-time constraint that invalidates earlier work is discovered last, when rework costs whole feedback cycles. No incident this sprint; ranked a real-but-low-frequency residual."
started:
completed:
verdict:
score: 0.3
worktree:
issue:
id: 73s70j4h8p217kqjwe8r7ewx
---

Hardest-first says the constraints most able to invalidate work should be visible while the work is cheap; the merge module's late load inverts that for merge-boundary rules. Direction for ideation (options open): a one-line merge-policy digest carried in the dispatch file the build helper already assembles (stage defs already ride it — zero boot-resident cost); or merge-relevant obligations named in the worktree stages' README definitions; or a short forward-pointer in the dispatch module. Constraint: must not grow the boot-resident contract (the 0250 leanness ceiling discipline stands). Acceptance sketch: value — an implementation dispatched on a fixture whose merge policy imposes a report-shape obligation can satisfy it without any terminal-boundary rework (baseline: the obligation is only discoverable at merge-core load); mechanism — the chosen surface ships with a fixture test.
