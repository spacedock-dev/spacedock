---
id: aw0sk3hagmjqs5cxrp732axe
title: merge guard archive commit discards git stderr and has no index.lock retry
status: backlog
source: "Fable ensign incident diagnostic, 2026-08-03, investigating a merge guard CRITICAL failure during collapse-gate-approval-ceremony's finalize. Root cause: .git/index.lock contention from concurrent ensign writers on the shared state checkout hit commitArchiveMove's git add (internal/status/merge.go:783) and the rollback's git reset (merge.go:739) -- both index-writing operations, both exit 128, while the rollback's plain os.Rename reversals succeeded (matching the lock signature, not disk pressure -- confirmed by a successful retry under identical disk conditions). runGitCmd (internal/status/handlers.go:677) uses cmd.Output(), so Go captures git's real stderr into ExitError.Stderr, but callers %v-format the error and silently discard the captured stderr -- the tool had the answer in memory and never printed it. Separately, rollbackArchive (merge.go:694) emits CRITICAL whenever any rollback step errors, without verifying the end state -- in this incident rollback had actually fully succeeded (nothing was left to unstage), so the CRITICAL warning was a false alarm that could have masked a genuinely critical case in a differently-shaped failure."
started:
completed:
verdict:
score: 0.4
worktree:
issue:
---

Two related hardening fixes to `merge guard`'s archive-commit git write path, both using existing in-repo patterns rather than new mechanisms:

1. **Surface real git stderr instead of discarding it.** `runGitCmd`'s underlying `cmd.Output()` already captures git's stderr into `exec.ExitError.Stderr`; `commitArchiveMove`/`rollbackArchive`'s callers currently `%v`-format the wrapped error and lose it. `internal/cli/state.go:360`'s `runGitOutput` already returns the captured stderr correctly -- reuse that pattern here so a future incident is self-diagnosing instead of an opaque "exit status 128."
2. **Retry on `.git/index.lock` contention.** `internal/cli/state_sync.go:317`'s `runGitRetryLock` already retries this exact class of failure (4x, 500ms) for a documented reason ("a sibling writer's add/commit can hold the lock"). Apply the same retry to `commitArchiveMove` and `rollbackArchive` -- the retry alone would have prevented this specific incident.
3. **Verify rollback's actual end-state before emitting CRITICAL.** After a rollback attempt, check the entity is genuinely back at its pre-attempt snapshot (path, bytes, `git diff --cached --quiet` on the touched pathspecs) before deciding whether to warn. Emit CRITICAL only on a verified inconsistency; otherwise report "rollback verified complete; original failure: <real stderr>" so a false alarm doesn't erode trust in a genuinely critical warning.

## Out of scope

A disk-space preflight check -- ruled out by the diagnosis as not the proximate cause of this incident, and disk-free-space thresholds are noisy on APFS (purgeable space). Worktree/scratch-dir cleanup itself is a separate, already-filed entity.

## Acceptance criteria

**AC-1** -- A simulated `.git/index.lock` collision during `commitArchiveMove` or `rollbackArchive` is retried and succeeds without a false CRITICAL, verified by a test that holds the lock file briefly during the operation.

**AC-2** -- A genuine (non-lock) git failure during the archive commit surfaces the real git stderr text in the reported error, not just an exit code.

**AC-3** -- Rollback verification: a test asserting that after a successful rollback, no CRITICAL is emitted, and the reported message names the original failure's real stderr.

## Test plan

Offline only. Fixture-based: simulate lock contention via a held `index.lock` file, simulate a genuine git failure via a broken repo state, assert on captured output/error text and CRITICAL-vs-clean messaging in each case.
