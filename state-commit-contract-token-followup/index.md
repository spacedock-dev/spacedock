---
title: Trim state.commit contract token load
status: implementation
source: "Captain follow-up after per-runtime contract token count (2026-06-20): shared-core grew +943 tokens since v0.22.0; biggest new body is «state.commit» at ~587 tokens. Determine and ship the leanest safe form."
started: 2026-06-20T18:36:31Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-state-commit-contract-token-followup
issue:
sprint: 0221-layered-fo
sprint-readiness:
id: 6cccykszyvxz5mxhrf270fa2
---

The v0.22.0-to-current contract token comparison showed `skills/first-officer/references/first-officer-shared-core.md` grew by about 943 `o200k_base` tokens. The largest new section is `«state.commit»(slug)` at roughly 587 tokens. This task should determine whether that body can collapse now that the binary/state-verb path exists, and then implement a safe reduction or produce a documented no-cut decision if the mechanics remain load-bearing.

## Problem

`«state.commit»` currently carries detailed split-root commit mechanics in the boot-resident shared core. That may be the right temporary prose-function form, but it is now the largest shared-core token-growth driver. The contract should keep the state safety guarantees while avoiding duplicated or obsolete mechanics that the binary already owns.

## Proposed approach

Ideation should compare the current `«state.commit»` body against the shipped `spacedock state commit` behavior and related state-management text. Identify which parts are still unique FO obligations, which are binary-owned, and which are duplicated elsewhere. Prefer a minimal shared-core body that names the intent, required invariant, shipped command, and failure handling, with detailed mechanics moved to binary docs/tests or removed when already enforced.

Coordinate with `trim-dispatch-adapter-prose` (`ad`) only where the same contract-hygiene principles apply. This task is specifically about state/gate shared-core token load, starting with `«state.commit»`; it is not the dispatch-adapter trim.

## Out of scope

Do not redesign state branch topology, merge hooks, PR completion, or dispatch capability bindings. Do not weaken path-scoped commit, non-fast-forward retry, rebase-conflict halt, or split-root safety guarantees.

## Acceptance criteria

**AC-1 - The `«state.commit»` shared-core body is either measurably reduced or explicitly justified as still load-bearing.**
Verified by: an implementation report with before/after `o200k_base` token counts for `first-officer-shared-core.md` and the `«state.commit»` section, plus a diff or no-cut rationale.

**AC-2 - State safety guarantees remain enforced by behavior or tests, not only prose.**
Verified by: relevant Go tests or command-level fixtures covering path-scoped state commit, push/rebase handling where practical, and rebase-conflict halt/refusal behavior; if an existing test is the proof, cite it and run it.

**AC-3 - The resulting contract text follows runtime-support principles.**
Verified by: contractlint or targeted review showing the shared core names lifecycle/state capabilities and avoids unnecessary duplicated mechanics, mutable step-number coupling, and runtime-specific tool bindings.

## Test plan

Run focused tests for the state commit command and contract lint touched by the change, then `go test ./internal/contractlint ./internal/status ./internal/cli` at minimum. If the implementation changes Go state-commit behavior, run `go test ./...`.

## Ideation recommendation

Recommend a measured trim, not a no-cut. The `spacedock state commit <slug>` command now owns the operational mechanics that the boot-resident shared core still repeats: entity path resolution, path-scoped commit, push, non-fast-forward `pull --rebase`, re-push, no-origin local-only behavior, clean no-op behavior, and same-entity rebase-conflict halt. The current `«state.commit»` plus `### Split-Root State Sync` block is 587 `o200k_base` tokens in `skills/first-officer/references/first-officer-shared-core.md`; the narrow `«state.commit»` function body alone is 143 tokens. A lean replacement can carry the unique FO obligations in about 166 tokens, for a target cut of at least 350 `o200k_base` tokens from the current 587-token block.

The implementation should delete the standalone `### Split-Root State Sync` mechanics from the boot-resident shared core and keep one compact `«state.commit»` body:

```markdown
## «state.commit»(slug): record one entity's change durably and concurrency-safe

- **effect:** invoke `spacedock state commit <slug>` after each state mutation. The command resolves the split-root entity, commits only that entity path, syncs with `origin` when present, and reports local-only/no-op as needed.
- **done-when:** the command exits 0 with the entity committed, and pushed when an origin exists.
- **block:** exit 3 means a same-entity rebase conflict was aborted by `spacedock state commit`; HALT, surface the named conflicting path(s) to the captain, and do not force-push or auto-resolve.
- → **shipped**: `` `spacedock state commit <slug>` `` — invoke it directly.
```

The safety guarantees remain covered by shipped command behavior and tests. `internal/cli/state_commit_test.go` covers same-entity conflict halt with exit 3, conflicting path naming, clean checkout after abort, no force-push, path-scoped commits that do not sweep sibling files, multi-writer non-fast-forward rebase/re-push for disjoint entities, no-origin local-only JSON, and clean no-op commits. `internal/cli/state_ready_test.go` covers the boot-side rebase halt and integration path. The existing `internal/contractlint/prose_function_backstop_test.go` currently pins two prohibition strings in the shared core; implementation should either update that test to the new compact exit-3/FO-obligation wording or replace it with a backstop that proves the compact block still names "do not force-push or auto-resolve" without preserving the old multi-step recipe.

No spike needed: the riskiest mechanism is already exercised by the real-git command tests named above, and the token measurement was run against the live file with `uvx --from tiktoken python` using `o200k_base` (`first-officer-shared-core.md` current total: 7666; current `«state.commit»` plus sync block: 587; proposed compact block: 166).

Boundary with `trim-dispatch-adapter-prose` (`ad`): this ticket owns only `skills/first-officer/references/first-officer-shared-core.md` state-commit/shared-core load and any directly necessary contractlint updates for that trim. It should not touch `fo-dispatch-core.md`, `claude-fo-dispatch.md`, runtime adapter await/reuse/guardrail prose, worker capability bindings, or merge-guard choreography; those remain `ad`'s scope.

## Stage Report: ideation

- DONE: Assess whether the current «state.commit» shared-core body is still load-bearing or can be reduced now that state commit behavior is binary-owned.
  Determined it can be reduced: current 587-token `«state.commit»` plus sync block mostly restates shipped `spacedock state commit` mechanics; compact target is about 166 `o200k_base` tokens.
- DONE: Compare the current shared-core state prose against shipped state command behavior and tests, naming any duplicated or obsolete mechanics.
  Duplicated mechanics: path resolution, path-scoped commit, push, non-ff `pull --rebase`, re-push, no-origin local-only, no-op, and conflict abort/halt are in `internal/cli/state_sync.go` and covered by `internal/cli/state_commit_test.go` / `state_ready_test.go`.
- DONE: Produce a concrete recommended implementation scope with acceptance criteria/test updates, including before/after token-count targets or a no-cut rationale.
  Added recommended replacement text and targets: current shared core 7666 tokens; current block 587; proposed block 166; implementation should cut at least 350 tokens and update the contractlint backstop to the compact obligation.
- DONE: Coordinate boundary with existing task ad/trim-dispatch-adapter-prose so this ticket owns state.commit/shared-core token load and does not duplicate dispatch-adapter trim.
  Recorded boundary: this ticket touches state-commit shared-core and direct contractlint only; `ad` owns dispatch adapters, await/reuse/guardrail prose, worker capabilities, and merge choreography.
- DONE: Append a Stage Report: ideation to the entity file and commit only that entity path in the state checkout.
  This report is appended here; commit follows path-scoped in `docs/dev/.spacedock-state`.

### Summary

Ideation recommends trimming the boot-resident `«state.commit»` prose now. The command and tests own the operational safety mechanics; the shared core should keep only the intent, direct command call, exit-0 done condition, and exit-3 FO halt/surface obligation.
