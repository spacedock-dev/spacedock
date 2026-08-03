---
id: m8m6s2v604cbysqf83dzfcfg
title: Stale worktree and scratch-dir reclamation (disk at 100%)
status: backlog
source: "Fable ensign incident diagnostic, 2026-08-03, investigating a merge guard incident. Disk is at 100% capacity with 2.3Gi free (worse than the 3.8Gi observed at incident time, days prior). git worktree prune finds 0 stale entries -- every registered worktree directory still physically exists. Full scope: 83 registered git worktrees. Sizes: .worktrees/ 3.4G across 32 dirs, .claude/worktrees/ 321M across 31 dirs (Task/Workflow-tool isolation worktrees from past sessions), /private/tmp/spacedock-* 5.4G across 272 dirs -- of which 3.4G is a single sessions.db in /private/tmp/spacedock-survey.VSXdX8 (a daemon home from survey tooling, not a worktree; may hold recorded sessions worth preserving -- captain call, not an automatic delete). Reclaimable universe is roughly 9G against 2.3G currently free."
started:
completed:
verdict:
score: 0.5
worktree:
issue:
---

One-time triage and reclamation of stale git worktrees and scratch directories accumulated across many past sessions, plus one policy decision about ongoing accumulation.

## Safety classes (from the incident diagnostic's survey)

- **Unambiguously safe now:** 8 worktrees on branches already merged to `main` -- `release/0.25.0-pre2`, `agy-runtime-support`, `codex-context-budget-probe-observer-evidence`, `gate-feature-prerelease-pilot`, `live-lanes-red-on-every-branch`, `restore-fo-merge-write-lazy-loading-reset`, `survey-claude-cowork-runtime-detection`, `codex-ci-luna-20260801`. Two of these carry 1-4 uncommitted files -- inspect before removing, don't blind `--force`. Caveat: a merged-ancestor check undercounts anything landed via squash-merge (verify each one's actual merge status, not just ancestry, before removing).
- **Probably safe:** ~15 detached-HEAD audit/spike worktrees under `/private/tmp/spacedock-*` (read-only audit checkouts, no uncommitted work expected) plus 6 `audit-*` worktrees under `.worktrees/`, plus the 31 dead-session `.claude/worktrees/wf_*`/`agent-*` directories from completed past Task/Workflow-tool runs.
- **Must exclude:** any worktree owned by a currently-active agent/session (e.g. one already observed with live teammates attached). Verify no live process/session references a worktree before removing it.
- **Captain decision, not automatic:** the 3.4G `sessions.db` under `/private/tmp/spacedock-survey.VSXdX8` -- unrelated survey-tooling daemon state, not a worktree, may hold data worth preserving.

## Policy question for ideation to resolve

Should `merge guard`'s finalize path (or the FO's own teardown step) automatically reap the finalized entity's own worktree as part of finalizing, rather than relying on the FO to remember the manual `git worktree remove` + `git branch -d` steps each time? This is exactly why the count keeps growing -- worktrees only get cleaned up when an FO happens to run the manual cleanup steps documented in the pr-merge mod, and evidently that hasn't consistently happened across many past sessions.

## Acceptance criteria

**AC-1 (VALUE)** -- Disk free space increases measurably (target: at least the ~9G reclaimable universe identified, safety classes permitting) and `git worktree list` no longer carries worktrees for merged/dead branches.
Verified by: `df` before/after, and a `git worktree list` diff showing only active-work and currently-in-use entries remain.

**AC-2** -- No active work is lost. Every removed worktree is independently verified (not just ancestry-checked) to be either fully merged with no uncommitted changes, or explicitly captain-approved as expendable scratch/audit output.
Verified by: a per-worktree log of what was removed and why (merged-branch confirmation or captain approval), committed as part of this entity's stage report.

**AC-3** -- The policy question is resolved and, if automatic reaping is chosen, implemented and tested; if not, the manual step is made more reliable (e.g. a boot-time or idle-hook reminder/lint) rather than left purely to FO memory.

## Test plan

Mostly operational/manual verification (this is a cleanup task, not new mechanism) -- if AC-3 lands automatic reaping, add a focused test for that specific code path. `go test ./...` must stay green throughout since worktree removal is filesystem-level, not code-level, but verify nothing in-repo assumes any of these paths exist.

## Out of scope

The `sessions.db` file's actual content/deletion -- flag it to the captain, do not delete without explicit direction.
