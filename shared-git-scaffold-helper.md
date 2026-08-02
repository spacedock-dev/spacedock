---
id: zfmbm75wvmfj38h73wrtpqqy
title: Add a shared, persisting git-scaffold test helper to close the identity-config gap
status: ideation
source: "Found while dogfooding collapse-gate-approval-ceremony's CI run, 2026-08-02: TestStampCommitsInlineBeforeWorktreeCreation's fixture omitted git config user.name/user.email, passing locally (ambient git identity) but failing on the CI runner (none). A follow-up search found no shared, persisting git-scaffold-init helper anywhere in the repo -- roughly 15-20 fixture functions across internal/cli and internal/dispatch each hand-roll their own init+config sequence independently, including one file-local gitInit helper (internal/dispatch/parity_harness_test.go) that only sets identity via scoped -c flags on its own seed commit, not a persisted git config, so it doesn't even cover a later plain git commit in the same repo. Captain directed: file and fast-track."
started: 2026-08-02T14:31:58Z
completed:
verdict:
score: 0.5
worktree:
issue:
gates:
    version: 1
    current:
        gate: gate:zfmbm75wvmfj38h73wrtpqqy:ideation
    records:
        - id: gate:zfmbm75wvmfj38h73wrtpqqy:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:zfmbm75wvmfj38h73wrtpqqy-backlog-1
              briefing:
                id: briefing:zfmbm75wvmfj38h73wrtpqqy:backlog:attempt-1:revision-1
                digest: sha256:ed3ebaef25d66a395a77d778b29879d7ff5ad90f40fc8b65d157c3df486059fd
                digest-domain: canonical-bytes
                request-digest: sha256:2ee4ca7169ea0dbbe53f8d8a436d26bdd6bb84a7608e059a04ce7681d557efe2
                room-ref: ./shared-git-scaffold-helper/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:zfmbm75wvmfj38h73wrtpqqy:backlog:1
                briefing: briefing:zfmbm75wvmfj38h73wrtpqqy:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-02T14:31:53.236223Z"
                decision: approve
                reason: 'Captain directed in chat: file the git-scaffold-helper follow-up and dispatch it fast-track.'
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
        - id: gate:zfmbm75wvmfj38h73wrtpqqy:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:zfmbm75wvmfj38h73wrtpqqy-ideation-1
              briefing:
                id: briefing:zfmbm75wvmfj38h73wrtpqqy:ideation:attempt-1:revision-1
                digest: sha256:59cc466ce2cb6c8f36c103d8c766897d685bced367fe201841e6359f9e2fb05e
                digest-domain: canonical-bytes
                request-digest: sha256:098b5702b7624a02cfba8844f982aa3292d74c16bf74be81fdaa8ca0e615d9de
                room-ref: ./shared-git-scaffold-helper/review/ideation/briefing-1
---

Add one shared, exported test helper that initializes a scratch git repo with a persisted `user.name`/`user.email` config in a single call, so no future fixture can reproduce the "works locally on ambient identity, fails on a clean CI runner" bug class, and migrate the existing hand-rolled fixtures to it.

## Problem

No shared git-scaffold-init helper exists anywhere in this repo. The ideation survey corrected the stub's original estimate: the duplication spans **9 packages / ~28 test files / ~34 non-bare `git init` scaffold sites**, not "~15-20 across `internal/cli` and `internal/dispatch`".

More importantly, the survey found the sites do not merely duplicate one correct sequence -- they use **four different identity mechanisms, and only one of them actually works**:

1. **Persisted repo config** (`git config user.email …`) -- correct. Covers every later commit in that repo, including one run by the code under test. Used by `internal/gitsource`, `internal/release/notes_extract_test.go`, `internal/status/mutation_test.go`, `internal/status/state_origin_test.go`, `internal/dispatch/reconcile_test.go`, `internal/gates/prepare_test.go`, `internal/statesync/publish_test.go`.
2. **Process env vars** on the test's own wrapper (`GIT_AUTHOR_NAME`/`GIT_COMMITTER_*` set on `cmd.Env`) -- `internal/cli`'s `git()` helper (`state_init_test.go:34`). Covers git calls made *through the wrapper* only.
3. **Scoped `-c` flags** per invocation -- `internal/dispatch/parity_harness_test.go:81`, `internal/status/worktree_overlay_test.go:36`, `internal/ensigncycle/cycle_test.go:331`. Covers that one invocation only.
4. **No identity at all** -- `internal/status/boot_orphan_abs_test.go:90` (`gitC`).

Mechanisms 2, 3, and 4 all leave the *same* hole: when the `spacedock` binary under test runs a plain `git commit` inside the fixture repo, it resolves identity from the repo/global config -- which the fixture never persisted. That is exactly the observed defect. `TestStampCommitsInlineBeforeWorktreeCreation` (currently on the `collapse-gate-approval-ceremony` worktree branch, `internal/dispatch/build_stamp_test.go:439`) scaffolds with mechanism 3, passed on every local run, and failed on the GitHub Actions runner. It was fixed by a two-line patch *at the call site*, carrying a comment that explains the trap for the next person -- which is the signal that the trap belongs in a helper, not in a comment.

**Why local runs did not catch it (spike finding, see below):** scrubbing `HOME` is not sufficient to reproduce. Git falls back to auto-detecting `user.email` from `username@hostname`, which succeeds on a developer Mac and failed on the CI runner. Deterministic local reproduction requires `user.useConfigOnly=true`.

## Proposed approach

### The helper

New package `internal/testgit`, single non-test file `internal/testgit/testgit.go`:

```go
// Package testgit scaffolds throwaway git repositories for tests.
package testgit

// InitRepo initializes a git repo at dir and persists user.name/user.email in
// the repo's own config, so a later plain `git commit` -- including one run by
// the code under test -- resolves an identity without the host's ambient git
// config. Extra args are appended to `git init` (e.g. "-b", "main").
func InitRepo(t testing.TB, dir string, initArgs ...string)
```

Decisions, each with the alternative considered:

- **A non-test `.go` file, not `_test.go`.** Go cannot import a `_test.go` symbol across packages, and the callers span 9 packages. *Alternative:* one copy per package -- that is the status quo being removed.
- **Package `internal/testgit`.** Matches the repo's small purpose-named package convention (`internal/gitsource`, `internal/statesync`). *Alternative:* a vaguer `internal/testlib` -- it would attract unrelated helpers and re-create the grab-bag this entity exists to prevent.
- **Importing `testing` from a non-test file does not bloat the binary.** Verified: `go list -deps ./cmd/spacedock` returns `internal/safehouse` (real production code) but *not* `internal/livescenario` or `internal/contractlint`, the repo's existing test-only support packages. `internal/testgit` will be imported only by `_test.go` files, so it stays out of the binary's dep graph -- re-checkable with that same one-line command. Stdlib precedent: `net/http/httptest`.
- **`testing.TB`, not `*testing.T`.** Free generality; callable from benchmark and helper contexts.
- **Variadic `initArgs`, so one function covers every site.** Real call sites need `-q`, `-b main`, and `-b next`. *Alternative:* separate `InitRepo`/`InitRepoMain` variants -- more surface for no gain.
- **No `InitRepoWithCommit` variant.** A seed commit is not part of the identity bug class, and callers that want one already have a package-local `git()` wrapper and two lines of `add`/`commit`. Adding a second exported function would be a second mechanism serving no AC.
- **Identity values `Spacedock Test` / `spacedock@example.invalid`,** matching `internal/gitsource` and the call-site fix already merged on the collapse-gate branch. `.invalid` is the RFC 2606 reserved TLD, so it can never route. Safe to standardize: **no test in the repo asserts on commit author or committer** -- every `git log` use is `--format=%H`, `--oneline`, `-S` pickaxe, `--name-only`, or `%s`; none uses `%an`/`%ae`.
- **The helper does not return anything.** No existing caller consumes a value from its init sequence.

### The migration

Replace the identity mechanism at every **non-bare** `git init` scaffold site with a call to `testgit.InitRepo`. Two shapes:

**Group A -- file-local init helpers whose bodies collapse to one call** (6): `internal/dispatch/parity_harness_test.go:81` `gitInit` (mechanism 3), `internal/status/worktree_overlay_test.go:36` `gitInitWorktreeFixture` (3), `internal/ensigncycle/cycle_test.go:323` `gitInit` (3), `internal/status/mutation_test.go:104` `gitInit` (1), `internal/status/state_origin_test.go:13` `gitInitNoRemote` (1), `internal/dispatch/reconcile_test.go:596` `repoGitInit` (1).

**Group B -- inline scaffold sites** (~28) across `internal/cli` (`gate_test.go:162,175`, `state_init_test.go:243`, `state_new_test.go:227,270`, `state_from_root_test.go:170,222`, `state_ready_test.go:136`, `decoupling_behavior_test.go:107`, `merge_test.go:50`, `terminal_consume_test.go:60`), `internal/status` (`boot_sandbox_test.go:38`, `boot_identify_test.go:394`, `entered_stage_test.go:287`, `boot_state_remote_test.go:18`, `merge_guard_foreign_cwd_test.go:44,56`, `boot_orphan_abs_test.go:60,128,207`), `internal/gates/prepare_test.go:1077,1093`, `internal/gitsource/source_test.go:244`, `internal/statesync/publish_test.go:31`, `internal/release` (`notes_extract_test.go:34,87`, `edge_reconcile_test.go:82`, `edge_advance_noregress_test.go:67`, `claude_candidate_binary_workflow_test.go:177`), and `internal/ensigncycle` (`liveassert_unit_test.go:95`, `live_test.go:416`).

Sites already on mechanism 1 are migrated too. They are not broken, but leaving them hand-rolled leaves the next author a working example of the pattern the guard forbids.

`internal/cli`'s `git()` wrapper (mechanism 2) keeps its `GIT_*` env vars -- they are harmless once identity is persisted, and stripping them is unrelated churn.

## Spike result

The checklist asked whether migrating `internal/dispatch`'s file-local helpers changes existing test behavior. **Spiked rather than assumed, because the answer was not obvious and two of the three findings contradicted the stub.**

1. **The migration is behavior-neutral.** Patched `parity_harness_test.go`'s `gitInit` from scoped `-c` to persisted config; `go test ./internal/dispatch -count=1` passed (14.8s). Reverted.
2. **The `internal/dispatch` gap is latent, not active.** Under a fully scrubbed identity, the package passes *both* patched and unpatched -- no later plain `commit` runs in those fixtures today. The migration there is preventive.
3. **`main` is currently green under identity scrub.** Full `go test ./...` with the scrub env exits 0. So "the suite passes on a clean runner" cannot serve as the value AC -- it already does. The value is preventing the *next* fixture from re-opening the hole, which is what drives AC-1's guard below.
4. **The stub's proposed AC-3 reproduction does not work.** `HOME`-scrubbing alone exits **0** -- git auto-detects `username@hostname`. The deterministic reproduction is `user.useConfigOnly=true`, which exits **128** (`Author identity unknown`). Fully env-injectable, so a test can apply it to a subprocess:

       GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null \
       GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=user.useConfigOnly GIT_CONFIG_VALUE_0=true

   This also explains the original escape: `TestStampCommitsInlineBeforeWorktreeCreation` *already* did `t.Setenv("HOME", …)` and still passed locally.

No further spike needed. The remaining mechanisms are proven: cross-package import of a non-test helper package (used today by `internal/livescenario`), and `filepath.WalkDir` repo sweeps in `internal/contractlint` (used today by `boundary_guard_test.go:75`).

## Expected surface

- **New:** `internal/testgit/testgit.go` (~35 LOC), `internal/testgit/testgit_test.go` (~70 LOC), one guard test in `internal/contractlint` (~60 LOC).
- **Modified:** ~28 test files across 9 packages, typically -3/+1 lines each.
- **Estimate:** ~31 files touched, ~165 insertions, ~90 deletions.
- **Observable semantics changed: none.** No command grammar, no stored format, no authority model, no runtime behavior, no user-visible output -- so no doc diff is owed. Test-only code; `cmd/spacedock`'s dependency graph is unchanged, checkable with `go list -deps ./cmd/spacedock`.
- **Tolerance:** ±10 files, ±80 LOC. If a fixture genuinely needs a non-persisted or differently-identified setup, leave it, add it to the guard's explicit allowlist with a reason, and report it rather than forcing the migration.

## Common scenarios and expected behavior

- A new fixture that uses the shared helper never has to remember the three-step init sequence, and cannot omit persisted identity by construction.
- Existing fixtures, once migrated, produce byte-identical test behavior -- this is a refactor, not a functional change to any test's assertions.
- The helper works identically regardless of the host machine's ambient git identity (or lack of it), which is exactly the gap that let this bug through undetected until a live CI run.

## Out of scope

Changing any test's assertions, coverage, or behavior beyond the init/identity mechanism. Any production (non-test) code change. Consolidating other kinds of test-fixture duplication not related to git identity/init.

Also explicitly out of scope:

- **`git init --bare` upstream fixtures** (`statesync/publish_test.go:41`, `dispatch/parity_harness_test.go:113`, `status/merge_guard_test.go:20`, `status/state_origin_test.go:44`, `status/boot_state_remote_test.go:22`, `cli/state_sync_test.go:76,124`, `cli/state_init_test.go:130,173`, `cli/state_from_root_test.go:106`, `cli/state_commit_test.go:24`). A bare repo never authors a commit, so it is not in the bug class. The AC-1 guard excludes them.
- **A CI scrubbed-identity lane.** Running the whole suite under `user.useConfigOnly=true` in CI would catch this class behaviorally as well as structurally, and the spike shows it is a one-line env prefix. But it is a CI/CD change and a separate decision from this refactor -- recommended as a follow-up entity, not assumed here.
- **Migrating `internal/cli`'s `git()` wrapper off its `GIT_*` env vars.** Unrelated churn once identity is persisted.

## Acceptance criteria

**AC-1 (VALUE) - The count of test-created git repos that can omit a persisted identity goes from ~34 to 0, and cannot climb back silently.**
Measured by a guard test in `internal/contractlint` that walks `internal/**/*_test.go` and fails on any non-bare `git … "init"` invocation that is not a `testgit.InitRepo` call, printing each offender as `file:line`. The number is the independent, wrong-way-movable baseline: **the guard fails today listing ~34 sites, and passes at 0 after migration**; a future hand-rolled fixture moves it back up and reds the build. Bare upstreams (`init --bare`) are excluded by the guard -- they never author commits, so they are not in the bug class.
*Alternative considered:* recording a one-shot grep count in the stage report. Insufficient: the count would be true on the day of the commit and unenforced thereafter, and "the bug class recurs silently" is the exact failure this entity exists to stop -- it already recurred once with no local or review signal. The guard is the only form of this AC that survives the commit that satisfies it.

**AC-2 - Existing tests are unaffected.**
Verified by `go test ./... -count=1` and `go test ./... -race` passing after migration, with no changes to any test's assertions or fixture content beyond the init/identity call substitution (reviewable in the diff: every hunk is a scaffold-site replacement).

**AC-3 - A repo scaffolded by the helper survives a plain commit with no ambient identity available.**
Verified by a focused test in `internal/testgit` that calls `InitRepo`, then runs a plain `git commit` as a subprocess carrying `GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=user.useConfigOnly GIT_CONFIG_VALUE_0=true`, asserting exit 0. **Falsifiability is pinned by a paired negative case in the same test:** the identical subprocess against a bare `git init` repo (no `InitRepo`) must exit non-zero. Without the negative case the positive one passes on any developer machine via hostname auto-detection -- measured above -- which is precisely how the original bug escaped. The scrub is applied to the subprocess only; no CI workflow change is in scope.

## Test plan

Offline only; no live workflow, no fixture goldens, no CLI behavior fixtures. Three pieces, in dependency order:

1. `internal/testgit/testgit_test.go` -- AC-3's positive/negative pair, plus one test asserting `git config user.email` resolves from the repo's own config after a single `InitRepo` call (i.e. persisted, not inherited). Written first; it is the spike exercise from finding 4 turned into a real test.
2. `internal/contractlint/git_scaffold_guard_test.go` -- AC-1's sweep. Follows the existing `boundary_guard_test.go:75` `filepath.WalkDir` + offender-list pattern. Expected to fail on first write and pass once the migration lands; that red-to-green transition is the AC-1 measurement and should be recorded in the implementation stage report with both counts.
3. Migration, then `go test ./... -count=1` and `go test ./... -race` for AC-2.

Cost: low. The helper and both new tests are small; the migration is mechanical but touches ~28 files, so the bulk of the effort is careful, repetitive substitution rather than design. Recommended `-race` budget ~10 min based on the observed ~15s for `internal/dispatch` alone.

**Coordination note for implementation:** `TestStampCommitsInlineBeforeWorktreeCreation` -- the fixture that motivated this entity -- currently lives on the `collapse-gate-approval-ceremony` worktree branch, not `main`. It is not in this entity's migration list. Whichever branch lands second should migrate it to `testgit.InitRepo` and drop the two-line call-site `git config` patch; otherwise the guard will flag it on merge.

## Stage Report: ideation

- DONE: Confirm and tighten the design already drafted in the backlog stub (shared helper name/location, migration scope) -- this is fast-tracked, so don't re-derive from scratch, just verify it's right and make it concrete (exact function signature, exact package/file location, exact list of fixtures to migrate).
  Pinned `internal/testgit/testgit.go` with `func InitRepo(t testing.TB, dir string, initArgs ...string)`; migration list enumerated to `file:line` in Proposed approach.
- DONE: Confirm no spike is needed (this is a mechanical Go test-helper refactor with no unverified mechanism), or spike the one thing that might be non-obvious: whether migrating the file-local gitInit/repoGitInit helpers in internal/dispatch changes any existing test's behavior.
  Spiked. Answer: no behavior change (`go test ./internal/dispatch` green patched and reverted) -- but the spike overturned two stub assumptions; see Spike result.
- DONE: Refine Acceptance criteria and Test plan only if the backlog stub's versions need correction; otherwise confirm them as-is.
  All three ACs rewritten -- the stub's AC-1 grep and AC-3 reproduction were both measurably wrong.

### Corrections the survey and spike forced

1. **Migration scope was understated ~2x.** Stub said "~15-20 fixtures across `internal/cli` and `internal/dispatch`". Actual: ~34 non-bare scaffold sites across **9 packages / ~28 files**.
2. **The duplication is not one repeated sequence but four different identity mechanisms, three of which are broken** (persisted config; `GIT_*` env on the test's wrapper; scoped `-c` flags; nothing). Only persisted repo config survives a plain `git commit` run by the code under test. This reframes the entity from DRY cleanup to correctness.
3. **Stub AC-3's reproduction does not reproduce.** `HOME`-scrub alone exits **0** (git auto-detects `username@hostname`); `user.useConfigOnly=true` exits **128**. Measured both. The stub's AC-3 test would have passed vacuously on any developer machine -- the same blind spot that let the original bug ship.
4. **`main` is already green under identity scrub** (full `go test ./...` exit 0), so "suite passes on a clean runner" is not a usable value metric. AC-1 was re-anchored on a `contractlint` guard whose offender count is the moving baseline (~34 -> 0, climbs back on regression).

### Summary

Confirmed the stub's core idea and made it concrete: one exported `testgit.InitRepo` in a new test-only `internal/testgit` package, migrating ~34 scaffold sites across 9 packages. The spike answered the assigned question (migration is behavior-neutral) and, more usefully, invalidated the stub's AC-3 reproduction method and its value metric -- both now measured rather than assumed. Two judgment calls left visible for the gate: the `contractlint` guard is a new mechanism justified in AC-1 against the weaker grep-in-report alternative, and a CI scrubbed-identity lane is recommended as a follow-up rather than pulled into this refactor.
