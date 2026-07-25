---
title: Trim state.commit contract token load
status: done
source: "Captain follow-up after per-runtime contract token count (2026-06-20): shared-core grew +943 tokens since v0.22.0; biggest new body is «state.commit» at ~587 tokens. Determine and ship the leanest safe form."
started: 2026-06-20T18:36:31Z
completed: 2026-06-21T05:40:58Z
verdict: passed
score:
worktree: .worktrees/spacedock-ensign-state-commit-contract-token-followup
issue:
sprint: 0221-layered-fo
sprint-readiness:
id: 6cccykszyvxz5mxhrf270fa2
mod-block:
pr: pr-merge:419
archived: 2026-07-25T19:09:55Z
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

## Stage Report: implementation

- DONE: Implement the approved measured trim for `«state.commit»` in `skills/first-officer/references/first-officer-shared-core.md`, targeting at least a 350 `o200k_base` token cut from the current `«state.commit»` plus split-root sync prose while preserving the unique FO obligations.
  Before/after with `o200k_base`: shared core 7666 -> 7233 tokens; `«state.commit»` plus sync block 587 -> 166 tokens, a 421-token block cut; code commit `eb2be8d`.
- DONE: Update any directly affected contractlint tests so they backstop the compact obligation rather than preserving the old multi-step recipe.
  Updated `internal/contractlint/prose_function_backstop_test.go` to pin the compact "do not force-push or auto-resolve" halt anchor instead of the removed multi-step recipe.
- DONE: Verify or cite behavior-level tests for path-scoped state commit, no-origin/no-op handling, non-fast-forward sync, and exit-3 rebase-conflict halt; add focused tests only if the existing source/test reality is insufficient.
  Existing `internal/cli/state_commit_test.go` covers path-scoped commits, no-origin local-only JSON, clean no-op, non-fast-forward rebase/re-push, and exit-3 conflict halt; `internal/cli/state_ready_test.go` covers boot-side rebase halt.
- DONE: Keep the boundary with `ad`: do not touch dispatch adapters, runtime await/reuse/guardrail prose, worker capability bindings, or merge choreography unless strictly required by this state-commit trim.
  Changed only `skills/first-officer/references/first-officer-shared-core.md` and the directly affected contractlint backstop.
- DONE: Run focused verification for touched contractlint/state areas, and broader `go test ./...` if Go behavior changes.
  Ran `go test ./internal/contractlint ./internal/status ./internal/cli` (`986` passed), `go test ./internal/contractlint` (`49` passed), `go test ./...` (`1692` passed), `go test ./... -race` (`1692` passed), `gofmt -w ./cmd ./internal`, and `git diff --check`.
- DONE: Append a Stage Report: implementation to the entity file with before/after token counts, changed files, verification evidence, and commit all implementation outputs on the worktree branch plus the entity report path-scoped in state.
  This report is appended here; code commit is `eb2be8d`, and the state commit follows path-scoped in `docs/dev/.spacedock-state`.

### Summary

Trimmed the boot-resident state commit contract to the direct command call, exit-0 done condition, and exit-3 FO halt/surface obligation. The removed mechanics are already owned by `spacedock state commit` and covered by real-git state command tests, while contractlint now backstops the compact no-force/no-auto-resolve obligation.

## Stage Report: validation

- DONE: Review the implementation diff on branch `spacedock-ensign/state-commit-contract-token-followup` in worktree `.worktrees/spacedock-ensign-state-commit-contract-token-followup` against the ideation recommendation and acceptance criteria.
  Commit `eb2be8d` changes only `skills/first-officer/references/first-officer-shared-core.md` and `internal/contractlint/prose_function_backstop_test.go`, matching the recommended state-commit/shared-core scope.
- DONE: Verify AC-1 with reproducible `o200k_base` token counts for `first-officer-shared-core.md` and the `«state.commit»` section/block; confirm the cut meets or exceeds the target or report the exact miss.
  Reproduced with `uvx --from tiktoken python`: shared core 7666 -> 7233 tokens; `«state.commit»` plus removed sync block 587 -> 166 tokens, a 421-token block cut exceeding the 350-token target.
- DONE: Verify AC-2 by running or reproducing the relevant behavior-level tests for state commit safety and confirming the report's cited tests cover path-scoped commit, no-origin/no-op behavior, non-fast-forward sync, and exit-3 conflict halt.
  Ran `go test ./internal/cli -run 'TestStateCommit(HaltsOnSameEntityConflict|IsPathScoped|MultiWriterHappyPath|NoOriginLocalOnly|NoOpWhenClean)|TestStateReadyHaltsOnBootConflict'`: 6 passed, covering conflict halt, path scope, non-fast-forward rebase/re-push, no-origin, no-op, and boot conflict halt.
- DONE: Verify AC-3 by reviewing the resulting contract text and contractlint changes for the runtime-support principles: capability-oriented shared core, no unnecessary duplicated mechanics, no mutable step-number coupling, and no runtime-specific tool binding introduced by this task.
  Reviewed compact `«state.commit»` text: it names the `spacedock state commit <slug>` capability, removes duplicated git mechanics and numbered sync steps, and adds no runtime-specific tool binding; `go test ./internal/contractlint` passed 49 tests.
- DONE: Run the implementation's cited verification or a justified focused subset; at minimum run the contractlint/state tests needed for the touched files, and report exact commands/results.
  Ran `go test ./internal/contractlint ./internal/status ./internal/cli`: 986 passed; `go test ./internal/contractlint`: 49 passed; `git diff --check`: clean.
- DONE: Append a Stage Report: validation with a PASSED or REJECTED recommendation, cite evidence for every AC, and commit only the entity report path-scoped in state.
  Recommendation: PASSED. This report is appended here; state commit follows path-scoped for this entity only.

### Summary

Validation recommends PASSED. The implementation meets the token-reduction target, keeps state safety backed by behavior-level tests, and leaves the shared-core contract capability-oriented without preserving the old multi-step git recipe.
