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

## Problem statement

A relative `--workflow-dir` is joined against the process cwd (`resolveRoots`, `internal/status/roots.go:41-47`; `dir` is `os.Getwd()` from `internal/cli/cli.go:555`). From a foreign cwd this mis-resolution ends in the same misleading error two ways:

1. **Resolved dir does not exist** (cwd outside the project): a missing README defaults `workflowIDStyle` to `sequential` (`internal/status/identity.go:74-78`), the entity scan of the nonexistent dir finds nothing, and `resolveMutationEntity` reports `Error: entity not found: <slug>` (`internal/status/native_runner.go:558-560`).
2. **Resolved dir exists but is the wrong copy** (the reported repro — cwd is a linked agent worktree of the project): `{worktree}/docs/<wf>/README.md` is tracked and present, its `state:` field diverges the entity dir to `{worktree}/docs/<wf>/.spacedock-state` — but the split-root state checkout is a separate untracked git checkout that exists only under the main working copy. The entity scan of the missing state dir finds nothing → the same `entity not found`.

In both cases the operator gets zero information that workflow resolution — not the slug — is the problem, and no recovery path short of knowing to `cd` back to the project root. Spike-verified end-to-end (see Spike record).

## Proposed approach

Keep the resolution rule as-is (a relative `--workflow-dir` stays cwd-relative — all state verbs and `merge guard` deliberately share this semantics, `internal/cli/state_sync.go:448-451`) and make `merge guard` **fail closed with an actionable diagnostic** when the resolved roots cannot possibly hold entities. Insert a roots guard in `MergeGuard` between `resolveRoots` and `resolveMutationEntity` (`internal/status/merge.go:82-87`), in the same spirit as the existing `validateRootsOrExit` used by `--validate` (`internal/status/roots.go:73-95`):

- **Resolved definition dir does not exist** → refuse (exit 1):
  `merge guard: --workflow-dir docs/dev resolves to /abs/cwd/docs/dev (a relative --workflow-dir resolves against the current directory), which does not exist`
- **Split-root declared, state checkout missing** → refuse (exit 1):
  `merge guard: workflow at /abs/cwd/docs/dev declares state: .spacedock-state but the state checkout is missing at /abs/cwd/docs/dev/.spacedock-state`
- **Corrective hint**, appended to either refusal when it validates: derive the main-checkout root as the parent of `git rev-parse --git-common-dir` run from the invocation dir (resolving a relative output such as `.git` against that dir — see Spike record), join the as-passed relative spelling to it, and when that candidate has a README whose roots resolve with an existing entity dir, append:
  `did you mean --workflow-dir /abs/mainroot/docs/dev? (the current directory is a linked git worktree; a relative --workflow-dir resolves against it)`
  When git fails (non-repo cwd, bare) or the candidate does not validate, the refusal stands without the hint. The hint is only computed for a relative as-passed spelling.

A valid resolvable workflow with a wrong slug keeps the existing `entity not found` error — that message is correct there.

Rejected alternatives:
- **Resolve relative `--workflow-dir` against `git rev-parse --show-toplevel`** (the issue's first suggestion): spike-proven insufficient — inside a linked worktree, `--show-toplevel` returns the *worktree* root, whose `docs/<wf>` still lacks the state checkout. Fixing it for real would need the common-dir parent, silently changing relative-flag semantics for every verb that shares the resolver.
- **Auto-retry against the main-checkout root**: makes the command "just work" but violates the resolver's never-guess discipline (cf. `discoverWorkflowDownward` refusing on ambiguity), and `merge guard` is a terminal mutation (terminalize + archive + commit) — silently redirecting its I/O to a directory the operator did not name is the wrong place to guess. The hint names the exact corrected invocation instead.

### Documentation diff

`internal/cli/help.go` merge help (the docs site's command reference defers to `--help` as the source of truth). Before:

    --workflow-dir DIR          Target this workflow explicitly (skips auto-discovery)

After:

    --workflow-dir DIR          Target this workflow explicitly (skips auto-discovery).
                                A relative DIR resolves against the current directory;
                                from anywhere else (e.g. an agent worktree) pass an
                                absolute path.

## Out of scope

- Issue #484's `state ready` defects (cwd-relative fallback-hint paths, `origin`-remote hard-fail) — sibling entity `state-ready-cwd-path-resolution`.
- Changing the cwd-relative resolution rule itself, for any verb.
- Applying the same fail-closed roots guard to the other mutation surfaces (`status --set`, `--archive`, `new`) — same confusion shape exists there; candidate follow-up once this pattern is proven, not this entity.
- `new`/`status` help-text wording for `--workflow-dir` (only the merge help changes here).

## Acceptance criteria

- **AC-1 (measures the end-value).** From a linked-worktree cwd, the issue-#485 repro invocation (`merge guard <slug> --verdict passed --workflow-dir docs/<wf>` with the split-root state checkout absent in the worktree) produces a refusal whose suggested corrected invocation actually works: stderr names the resolved absolute path, the missing state-checkout path, and the exact absolute `--workflow-dir` to use; re-running with that suggested value exits 0 and lands the mutation in the main checkout's state dir (frontmatter delta on disk). Baseline that can move the wrong way: today the identical sequence dead-ends at `Error: entity not found: <slug>` naming zero recovery paths; the fixed behavior names the corrected invocation and that invocation is verified to succeed. Verified by: a Go test that stages a git repo + commissioned split-root workflow + real `git worktree add`, calls `MergeGuard` with `dir` set to the worktree path, asserts the diagnostic contents, then re-runs with the suggested path and asserts the phase signal and the entity's on-disk frontmatter mutation.
- **AC-2.** A `--workflow-dir` resolving to a nonexistent directory is refused with a diagnostic naming the as-passed spelling and the resolved absolute path — `entity not found` no longer appears for this shape. Verified by: unit test from a temp cwd with a bogus relative flag; assert exit 1 and stderr content.
- **AC-3.** A resolved split-root workflow whose declared state checkout is missing is refused before entity resolution, naming the missing state-checkout path. Verified by: unit test on a split-root fixture with the state dir removed; assert exit 1 and stderr names `{definitionDir}/{state}`.
- **AC-4.** A wrong slug against a valid, resolvable workflow still reports `Error: entity not found: <slug>`. Verified by: unit test on a valid fixture with a bogus slug.
- **AC-5.** Existing `merge guard` behavior on valid invocations is unchanged. Verified by: the existing suites pass unmodified (`go test ./internal/status/... ./internal/cli/...`).
- **AC-6.** The merge help text states the cwd-relative resolution rule per the documentation diff above. Verified by: the help-output test covering `setMergeHelp`.

## Test plan

Go unit tests in `internal/status` following the existing `merge_guard_test.go` conventions (`stageFixture` + git-initialized fixtures; `MergeGuard(args, dir, ...)` takes `dir` explicitly, so a foreign cwd is simulated by passing a different `dir` — no chdir needed). The AC-1 test additionally runs real `git init` + `git worktree add` + a nested state-checkout `git init` in a `t.TempDir()` (git is already a test dependency via `runGitCmd` fixtures). One help-text assertion in `internal/cli` for AC-6. Estimated cost: small — one new test file (~4 tests), a roots-guard helper + call in `merge.go`, help-text edit; no fixtures beyond temp-dir construction, no live workflow tests (the claim is command-level behavior, fully covered by driving `MergeGuard` directly).

## Spike record

Riskiest paths exercised first against a built binary (`go build ./cmd/spacedock`) on a scratch repo with a commissioned split-root workflow (`state: .spacedock-state` as a separate untracked git checkout) and a linked `git worktree add .worktrees/agent-x`:

1. **Repro confirmed**: from the worktree cwd, `merge guard my-task --verdict passed --workflow-dir docs/dev` → `Error: entity not found: my-task`, exit 1. From the main root, the identical command resolves the entity (proceeds past entity lookup; the minimal fixture then stops at an unrelated no-terminal-stage refusal, confirming lookup succeeded). The worktree contains `docs/dev/README.md` but no `docs/dev/.spacedock-state`, exactly as analyzed.
2. **`--show-toplevel` is the wrong base**: from a worktree subdir it returns the worktree root (`.../spike-repo/.worktrees/agent-x`) — the issue's suggested fix would not repair the reported repro.
3. **Common-dir parent is the right base, with a path-shape caveat**: `git rev-parse --git-common-dir` from a worktree subdir returns the main checkout's absolute `.git`; from the main root it returns the relative `.git` (git 2.39.5) — so the implementation must resolve the output against the invocation dir rather than relying on `--path-format=absolute` (git ≥ 2.31 only). Parent of the resolved common dir = main-checkout root, where `docs/dev/.spacedock-state` exists.

## Stage Report: ideation

- DONE: Problem statement, proposed approach, and out-of-scope boundary written into the entity body
  Root cause traced to resolveRoots' cwd-join (internal/status/roots.go:41-47) with both failure shapes (nonexistent dir, worktree-shadowed split-root); approach is a fail-closed roots guard in MergeGuard with a corrective did-you-mean hint; two alternatives rejected with reasons; four items fenced out of scope.
- DONE: Entity-level acceptance criteria (each with how it's verified) plus a test plan added
  AC-1..AC-6, each with its verification; AC-1 measures the end-value (repro refusal names a corrected invocation that is then verified to succeed with an on-disk mutation) against today's zero-recovery-path baseline; test plan sizes the work as Go unit tests on existing merge_guard_test.go conventions plus one real-worktree test.
- DONE: If the design rests on an unverified mechanism, the riskiest path is spiked first and recorded (or "no spike needed" with the proven mechanisms)
  Spike record in the body: repro confirmed end-to-end against a built binary + real linked worktree; the issue's suggested --show-toplevel base disproven; git rev-parse --git-common-dir proven as the hint's main-root source with its relative-output caveat (git 2.39.5).

### Summary

Ideated the fix for issue #485: keep cwd-relative --workflow-dir semantics, add a fail-closed roots guard to merge guard that names the resolved path and the missing state checkout, plus a spike-verified did-you-mean hint derived from the git common-dir parent. Spiked the riskiest mechanisms first — the repro reproduces exactly as filed, and the spike disproved the issue's own --show-toplevel suggestion, which materially changed the proposed approach.
