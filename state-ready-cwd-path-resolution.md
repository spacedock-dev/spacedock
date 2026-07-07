---
id: e6j9adxnn5hgv4hd7g5edr3t
title: "spacedock state ready resolves init paths relative to cwd and requires an origin remote"
status: backlog
source: "GitHub issue #484 (spacedock-dev/spacedock#484), filed by clkao 2026-07-07, from a live split-root dogfooding session on this exact workflow shape (5 stages, ~6 entities, sd-b32 ids)."
started:
completed:
verdict:
score:
worktree:
issue: spacedock-dev/spacedock#484
---

`spacedock state ready` has two cwd/remote-sensitivity defects observed while recovering a deleted split-root state checkout. (1) When invoked from inside an agent worktree (`.worktrees/<worker>-<entity>/`), its manual-fallback hint proposes re-adding the state checkout at a path nested UNDER that worktree instead of the project root — the command should resolve paths against the project root (or the workflow dir's declared location) independent of cwd, since hooks/scripts may invoke it from anywhere. (2) When the repo has no `origin` remote (deliberate local-only project), the command hard-fails on `git fetch origin <state-branch>` even though the state branch exists locally and a plain `git worktree add` from it succeeds — it should fall back to the local branch when the fetch fails or no remote is configured, matching `state commit`'s documented local-only behavior. Full repro, expected behavior, and the manual workaround used are in the linked issue.
