---
title: Merge-guard verb — spacedock merge guard <slug> (atomic mod-block set→invoke→clear→terminalize)
status: backlog
source: 0205 carve (2026-06-17, captain "stamp them") — index DoD candidate "merge-finalize"; 2y MERGED unblocks it.
score: 0.5
sprint: 0205-layered-fo
group: verb-core
sprint-readiness: ready
id: mzmc0dgkq1nbazxx8j0mnfn6
---

`spacedock merge guard <slug>` enforces the mod-block set→invoke→clear→terminalize-only-after sequence atomically (set before the hook, detect completion by state delta, clear in a standalone `--set`, terminalize only after) — so the single highest Haiku merge risk (combining / skipping / reordering those steps) is owned by the binary. The status tool's existing refusal of terminal-with-mod-block-set is the backstop.

Spike must-build: «merge» (held mechanically by Haiku but carries judgment triggers). Build against the REAL post-2y host-neutral merge core (2y MERGED).

Ideation must resolve: the verb design vs the pr-merge / Ship-Local ceremonies; the AC boundary vs an optional `merge ceremony` / `worktree safe-remove` companion; whether it folds the `team-mode-verdict-omission` (re) atomicity fix (verdict-write atomic w.r.t. status + teardown); oracle-based ACs.
