---
id: gwttp4rhr07bqtejzebn7hha
title: Reconcile Class D/E must never destroy committed or unpushed work — rethink repo-hygiene, consider an ownership lease
status: ideation
source: captain + FO (2026-06-14, this session) — the reconcile sweep flagged local main (ahead by 1, the unpushed README cut) as Class E "stale local main" and prescribed `git reset --hard origin/main`, which would have deleted the commit. Sibling to #369 (separated reconcile team-management A/B/C from repo-hygiene D/E, fixed trunk detection) and #370 (pr-merge trunk refit off pre-flip `next`); this addresses the D/E remedy-safety gap those left.
started: 2026-06-14T21:10:17Z
completed:
verdict:
score: 0.42
worktree:
issue:
sprint: 0203-fo-efficiency
---

The reconcile sweep's repo-hygiene classes (D: stale worktree branch; E: stale local main) can destroy work. Rethink them so reconcile never destroys committed or unpushed work and never mutates a worktree the running session does not own. Keep it simple; if needed, model ownership as a lease.

## Problem

The event-loop reconcile sweep (contract: `claude-fo-dispatch.md` Event Loop step 0; helper: `internal/dispatch/reconcile.go`) classifies repo-hygiene drift into Class D (worktree branch behind `origin/{trunk}` → `pull --rebase`) and Class E (local main "stale" → `git fetch && git reset --hard origin/{trunk} && rebuild`).

Three faults, found 2026-06-14 (after the #369 trunk fix was already live):

- **Class E detects the wrong direction.** `classE` counts `git rev-list --count origin/{trunk}..main` — commits local main has that origin lacks, i.e. AHEAD / unpushed — stores it in an `Ahead` field, then emits `reset main->origin/{trunk}`. So Class E fires precisely when the FO holds unpushed commits on main (an FO-direct README/state edit awaiting a push), and the remedy discards them. The class is labelled "stale local main" but detects unpushed-ahead, not behind.
- **The remedy is destructive.** `reset --hard` against a branch that legitimately carries local commits (FO write scope commits the workflow README and state transitions to main, then pushes on the captain's word) can lose committed work. The contract should never prescribe a hard reset that can destroy committed work.
- **Class D mutates un-owned worktrees.** It auto-rebases any active entity's worktree behind trunk, including a worktree another session (the Commander) is actively driving — a cross-session collision on someone else's branch.

## Proposed approach

{Ideation fills. Captain direction: rethink the D/E classes; keep it simple; if necessary introduce an ownership-lease semantic so a session only auto-syncs state it owns. A candidate shape to evaluate, not prescribe: E distinguishes behind (`merge --ff-only` / `pull --rebase` + rebuild) from ahead/unpushed (push or report, never reset) and reports rather than auto-resolves a diverged main; D only auto-rebases a worktree the current team/session owns — the ownership signal that classes A/B/C already use via team identity / `leadSessionId` — and reports a peer-owned worktree instead of touching it.}

## Out of scope

{Ideation fills. Likely: the A/B/C team-management classes and the trunk-detection fix, both already handled by #369.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 — Reconcile never prescribes or performs a reset that can lose committed work.** For a local main carrying unpushed commits the remedy is non-destructive (push or report); a diverged main is reported, not auto-reset.
Verified by: {a reconcile unit test over a seeded repo whose local main is ahead of origin, asserting no `reset --hard` remedy is emitted and the reported class/remedy matches the non-destructive contract. Ideation pins the fixture.}

**AC-2 — Reconcile does not auto-mutate a worktree the running session does not own.** A worktree owned by another session is reported, not rebased.
Verified by: {a test over a roster/ownership fixture asserting a peer-owned worktree yields a report-only drift item, not a rebase action. Ideation pins it.}

## Test plan

{Ideation fills. Likely Go unit tests over `internal/dispatch/reconcile.go` with seeded git fixtures for the ahead/behind/diverged main cases and the owned/un-owned worktree cases; the contract-side wording change is verified by the behavior it points at, not by a prose match.}
