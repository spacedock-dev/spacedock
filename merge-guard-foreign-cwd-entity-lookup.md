---
id: bmexmw5bdjffr8n67gdp6vra
title: "spacedock merge guard: entity not found from a foreign cwd even with --workflow-dir"
status: ideation
source: "GitHub issue #485 (spacedock-dev/spacedock#485), filed by clkao 2026-07-07, from the same live split-root dogfooding session as #484 (5 stages, ~6 entities, sd-b32 ids) — may share a root cause."
started: 2026-07-07T22:59:51Z
completed:
verdict:
score:
worktree:
issue: spacedock-dev/spacedock#485
---

`spacedock merge guard <slug> --verdict passed --workflow-dir docs/<workflow>` reports `Error: entity not found: <slug>` (exit 1) when run from a foreign cwd — inside an agent worktree of the project — even though `--workflow-dir` is passed explicitly. The identical command succeeds once `cd`'d back to the project root. The explicit `--workflow-dir` flag not compensating for a foreign cwd suggests the relative value resolves against cwd rather than the enclosing repo root; either resolving it against `git rev-parse --show-toplevel`, or documenting that the flag must be absolute and erroring more clearly than `entity not found` when the resolved dir does not exist, would fix the confusion. Filed alongside #484 (cwd-sensitivity in `state ready`) as a possibly shared root cause. Full repro is in the linked issue.
