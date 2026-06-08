---
id: 69rk6t1vbehsd4fwxnsnjwma
title: survey — work around agentsview not persisting Codex cwd (scoping fallback + hint), do NOT fix agentsview
status: ideation
source: "captain (2026-06-08) — agentsview ingests Codex sessions and derives project but does NOT persist Codex cwd into sessions.cwd, so the survey's cwd-scoped repo-identity query misses all Codex sessions. Decision: work around it in OUR survey skill with a fallback/hint, not fix agentsview."
score: "0.28"
started: 2026-06-08T15:29:12Z
completed:
verdict:
worktree:
issue:
sprint: 0198-pre-flip-hardening
group: survey
sprint-readiness: ready
---

The survey's repo-identity scoping (the `git rev-parse --git-common-dir` coalesce, xn) scopes by `sessions.cwd`. agentsview does not persist Codex `cwd`, so cwd-scoped queries miss all Codex sessions. Work around it survey-side; do not fix agentsview.

## Problem

- agentsview ingests Codex sessions, derives `project='Foo'`, but leaves `sessions.cwd` unset for Codex.
- The survey's scoping `cwd = :repo_root OR cwd LIKE :repo_root || '/%'` therefore misses all Codex sessions.
- Filtering by `project` catches them but loses repo-root identity (can't distinguish same-basename repos, subdirs, or worktrees).

## Proposed approach (ideation firms)

A survey-side workaround — NOT an agentsview fix:
- Include Codex sessions via a `project` fallback when `cwd` is absent (e.g. union the cwd-scoped set with Codex sessions whose `project` matches the repo basename), AND/OR
- Surface a hint in the survey output when Codex sessions exist but lack `cwd` — so the reader knows Codex work may be under-scoped / matched only by basename.
- Ideation decides the blend (fallback query vs hint vs both) and how to bound the basename-collision risk the fallback reintroduces.

## Out of scope

Fixing agentsview to persist Codex `cwd` (upstream `kenn-io/agentsview`) — explicitly NOT this task.

## Acceptance criteria (sketch)

- The survey accounts for Codex sessions under a repo even when their `cwd` is unset — verified by a query-smoke fixture with Codex rows (`project` set, `cwd` null) asserting they are included (or flagged), plus a live drive showing the hint when applicable.

## Notes

Survey-skill workaround; same area as xn. Root cause is upstream agentsview (`sessions.cwd` unset for Codex); we deliberately work around it rather than block on the upstream fix. Candidate for a survey-followup / 0.19.8 sprint.
