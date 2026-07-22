---
title: Classify control-plane work without forcing it into workflow entity state
status: backlog
source: "GitHub issue spacedock-dev/spacedock#555 (captain-filed) — the FO contract pressures agents to create meta-entities merely to obtain an authorized write/worktree/dispatch for control-plane work (governance, release coordination, portfolio/sprint shaping)."
issue: spacedock-dev/spacedock#555
id: j7st8gr31kq38w6t6vb884vc
---

Introduce an explicit work-plane classification (a `work.classify` decision BEFORE `write.classify`) so control-plane work is never turned into a workflow entity just to obtain write authority, a worktree, or a worker dispatch. Full problem, desired contract, and acceptance criteria in #555 — ideation reads the issue.

## Problem
{Ideation fills in from #555. Seed: path-based write.classify + entity-only dispatch can force control-plane work (owned by an active role) into a meta-entity whose deliverable is the shaping activity itself — recursive task state. The reusable fix belongs in the FO framework, not a project folder exception.}

## Acceptance criteria
{Ideation fills in from #555's Acceptance section: precedence stated between work-ontology / role-ownership / path-based write.classify / dispatch; a bounded control-plane route with no hard-coded folder allowlist; a fixture/contract test proving no meta-entity is created for shaping work; a second fixture proving real member outcomes still route through entity dispatch; missing authority still stops for captain direction.}
