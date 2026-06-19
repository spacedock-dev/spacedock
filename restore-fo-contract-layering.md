---
id: zyhwwm5a9a08vv6htpef1c34
title: Restore core/adapter and core/mod layering in the FO contract
status: backlog
source: '0221 captain review (2026-06-19, CL): auditing the layered-FO contract surfaced two pre-existing layering leaks in the host-neutral fo-dispatch-core.md. Defer 72 (tier-delegation); land this cleanup + the prose-function notation on a clean foundation first.'
started:
completed:
verdict:
score: 0.5
worktree:
issue:
sprint: 0221-layered-fo
sprint-readiness: ready
group: cleanup
---

The host-neutral FO dispatch core has leaked host-specific and mod-specific detail; restore the core/adapter and core/mod boundaries before more contract is built on it.

## Problem

`skills/first-officer/references/fo-dispatch-core.md` is the host-neutral dispatch core, but two leaks violate that boundary:

1. **Claude model identifiers in the generic core.** Reuse-condition-4 names `opus`/`sonnet`/`haiku` and the `"opus[1m]"` fallback example directly. These are Claude-specific; the host-neutral core should refer generically to "the host's canonical model enum" and delegate the actual enum + the `[1m]` example to the Claude adapter (`claude-fo-dispatch.md`), where the model-for-member lookup and context-budget mapping already live. The behavioral contract (a stamped fallback value forces a one-time fresh re-stamp) must not change.

2. **A GitHub PR-pending scan in the generic event loop.** Event-loop step 1 ("Check PR-pending entities") runs `gh pr view` and advances merged PRs directly in the host-neutral core. PR lifecycle is the `pr-merge` mod's domain — its `idle` hook already does this scan (and its own prose calls itself "defense in depth" for the core's built-in scan). A workflow with no `pr-merge` mod (a `merge: local` or non-code workflow) should never reach for `gh` in its loop. The generic loop should only fire idle hooks; the `pr-merge` idle hook should own PR-pending entirely.

Both are pre-existing debt, but the layered-FO sprint is adding more contract on this foundation, so they are fixed first. A possible third item to assess in ideation: the `status --set` merge-hook guard falsely trips on a non-terminalizing `worktree=` clear (observed 2026-06-19) — confirm scope and fold in or split out.

## Out of scope

- The 72 (fo-tier-delegation) tier map's own model-name references — 72 is deferred; it rebuilds on this clean foundation in a later sprint.
- The prose-function `«»` notation restructure — that is the prose-function-restructure (czw) member, which lands after this cleanup.

## Acceptance criteria

Ideation fills these in — each an end-state property proven by behavior or a legitimate structural (contractlint) check against an independent source, never a prose-grep tautology.
