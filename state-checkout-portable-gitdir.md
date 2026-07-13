---
title: Make the split-root state checkout operable from foreign path surfaces (Cowork sandbox)
status: implementation
source: live Cowork dogfood 2026-07-13 (survey-claude-cowork-runtime-detection run-1 session)
started: 2026-07-13T02:54:30Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-state-checkout-portable-gitdir
issue:
id: qwp7vn4dnt0sx96wzz1wfy8b
---

Make every Spacedock-owned state operation work when the same repository tree is reached through a path surface different from where its state checkout was created. Observed live: a Cowork sandbox mounts the project at `/sessions/<name>/mnt/<folder>` while the linked checkout's `.git` file still points at the host-absolute administrative directory, so Git reports `not a git repository` even though the checkout and repository metadata are both present under the mounted project.

## Problem

By default, `git worktree add` records absolute paths in both the checkout's `.git` pointer file and the administrative directory's `gitdir` back-pointer. The split-root checkout lives inside the project (`docs/dev/.spacedock-state` backed by `.git/worktrees/<id>`), so a second path surface for the same tree — Cowork mount, bind mount, moved directory, or remounted disk — leaves all data reachable but makes those recorded paths stale. Observed live in Cowork: `spacedock state commit` failed, while an explicit `GIT_DIR=.git/worktrees/-spacedock-state GIT_WORK_TREE=docs/dev/.spacedock-state` commit succeeded on Git 2.34.

The persisted administrative name is not necessarily the checkout basename: Git may append a numeric suffix when names collide. Recovery therefore must preserve the administrative id parsed from the stale `.git` pointer and re-anchor it under the current repository common directory; it must not guess `<basename>` alone.

## Compatibility spikes

The risky mechanisms were exercised before selecting the design. All spike
repositories lived under `/tmp`; no project metadata was changed.

### Legacy absolute links and collision ids

Local Git 2.39.5 created two linked worktrees with the same basename. Git named
the second administrative directory `state1`, proving that recovery cannot
derive the id from the checkout basename. After moving the entire repository
from path A to path B:

```text
checkout pointer: gitdir: /private/tmp/.../A/repo/.git/worktrees/state1
admin back-pointer: /private/tmp/.../A/repo/two/state/.git
git -C B/repo/two/state status: exit 128
fatal: not a git repository: /private/tmp/.../A/repo/.git/worktrees/state1
GIT_DIR=B/repo/.git/worktrees/state1 GIT_WORK_TREE=B/repo/two/state git status:
## state-two
```

This reproduces the Git 2.34 Cowork failure shape on another pre-2.48 Git and
demonstrates the compatibility runner's core operation. Neither link file needed
rewriting.

### Relative-worktree compatibility consequence

A disposable Git 2.48.1 built from the official release tarball created a
worktree with `--relative-paths`, then the whole repository moved:

```text
checkout pointer: gitdir: ../../.git/worktrees/state
admin back-pointer: ../../../nested/state/.git
extensions.relativeWorktrees=true
Git 2.48.1 after move: exit 0, ## state
Git 2.39.5 on the same repository: exit 128
fatal: unknown repository extension found: relativeworktrees
```

Setting `worktree.useRelativePaths=true` produced absolute links and no
extension on Git 2.39.5, but relative links plus the extension on Git 2.48.1.
Thus the config alone silently changes meaning by Git version. Spacedock must
probe the option when the config opts in, and relative creation must remain
off by default while Git 2.34 consumers are supported. This matches Git's
[`git-worktree` documentation](https://git-scm.com/docs/git-worktree), which
warns that the extension is incompatible with older Git.

## Proposed approach

Ship two complementary lanes, with legacy recovery first.

### 1. Shared state-checkout Git runner

Add a small `internal/stategit` package with one resolved runner for a project
checkout plus its nested state checkout. Three present consumers justify the
package: state verbs in `internal/cli`, split-root Git probes/commits in
`internal/status`, and assignment guidance in `internal/dispatch`.

The runner preserves current behavior on the normal path. It resolves the
checkout's `.git` pointer; when its target exists, commands remain ordinary
`git -C <state-checkout> ...` calls with unchanged stdout, stderr, and exit
semantics.

Fallback is narrow and read-only:

1. It activates only for a linked checkout whose `.git` is a well-formed gitfile
   and whose recorded target is missing. An arbitrary Git command failure never
   triggers recovery.
2. It parses the exact final element after `worktrees/` as the administrative id.
   The id must be one clean path element; empty, `.`, `..`, separators, and
   escaping forms fail.
3. It asks the still-operable project checkout for `--git-common-dir`, resolves
   that path under the current surface, and constructs
   `<common-dir>/worktrees/<id>`. A `filepath.Rel` confinement check keeps the
   candidate below the current `worktrees` directory.
4. It requires the candidate directory plus `HEAD`, `commondir`, and `gitdir`.
   `commondir` must resolve back to the current common directory, and the stale
   `gitdir` back-pointer must end in the current repository-relative state path
   plus `/.git`. A missing or mismatched check returns one recovery diagnostic;
   it never scans for or guesses another administrative directory.
5. It runs Git with the validated candidate as `GIT_DIR` and the current state
   checkout as `GIT_WORK_TREE`. It never runs `git worktree repair` and never
   rewrites either link file.

Resolve the runner once per command so every subordinate operation shares the
same mode. Route all current state-checkout Git calls through it: `state commit`
and its lock/rebase/conflict helpers, `state ready`, the present-checkout branch
of `state init`, boot's origin probe, split-root archive commit/rollback in merge
guard, and dispatch's state-commit guidance probe. Dispatch must propagate a
runner or origin-classification error and produce no misleading local-only
guidance.

Give the shared runner one tri-state origin query rather than independent
boolean helpers. It first runs `git remote`; failure to enumerate remotes is an
error. If that successful result does not contain the exact remote name
`origin`, the checkout is genuinely local-only. If `origin` is listed, it then
runs `git remote get-url origin`; failure is an error and success means origin
is present. Thus only successfully proved absence selects `remote: none` or
local-only guidance; a declared origin with an unreadable URL is fatal in CLI,
boot, and dispatch.

For an already-present checkout, `state init` resolves the runner and uses that
same origin query. A proved-absent origin preserves today's initialized success
path. A present origin must complete both fetch and rebase pull before the
success line is printed; either failure is returned instead of being silently
discarded. Resolver validation failure returns the stable recovery diagnostic
without running fetch or pull.

`state sweep` remains unchanged because it reads state files and invokes Git only
for the operable main checkout. Keep it in the end-to-end fixture to prove the
public convergence sequence rather than adding a recovery call where none is
needed.

### 2. One creation policy for `state new` and `state init`

Keep persisted state worktrees absolute by default. This is the compatibility
policy while Git 2.34 remains a supported consumer; lane 1 makes those checkouts
portable through Spacedock commands.

Use Git's repository-local config as the explicit opt-in, avoiding a new
Spacedock flag or README field. Read only the local scope with a typed boolean;
an invalid local value is an error. Inherited system or global config never
opts a repository into the compatibility-breaking extension:

- local unset/false `worktree.useRelativePaths` -> every creation-time
  `worktree add` runs with `-c worktree.useRelativePaths=false`, overriding an
  inherited true value and forcing absolute links on both old and new Git;
- local true -> before any birth/resume mutation, probe `git worktree add -h`
  for the `--relative-paths` option token; if absent, fail with a diagnostic
  that names the Git 2.48 floor and tells the operator to upgrade or unset the
  local config;
- local true with support -> every creation-time `worktree add` passes
  `--relative-paths`.

Apply the selected policy to the temporary orphan-birth worktree as well as the
persisted checkout. On Git 2.48 an inherited true value on that temporary add
can itself enable `extensions.relativeWorktrees`; forcing false in the default
lane prevents a repository-wide compatibility change before the final add.

The option probe avoids vendor-version parsing and prevents old Git from silently
ignoring the opt-in. Documentation must state that opting in declares every
consumer to be Git 2.48+ because Git sets `extensions.relativeWorktrees` and
older Git refuses the entire repository. Existing absolute checkouts are not
converted; lane 1 remains required.

## Concrete documentation diff

Append this section to `docs/site/advanced/split-root-state.md`:

> ## Moved and sandbox-mounted repositories
>
> Spacedock state commands continue to operate when a legacy state checkout's
> absolute Git link names another path surface, as long as the main project
> checkout and its `.git` directory are reachable together. Recovery validates
> the state worktree's existing administrative id and does not rewrite repository
> metadata. Bare `git -C <state-checkout>` may still fail on that legacy checkout;
> use `spacedock state ready`, `state sweep`, and `state commit` for workflow
> convergence.
>
> New state checkouts keep absolute Git links by default for compatibility with
> older Git. To opt into movable relative links, first ensure every host and
> sandbox that opens the repository uses Git 2.48 or newer, then run
> `git config --local worktree.useRelativePaths true` before `spacedock state
> new` or `spacedock state init`. Only a repository-local true value opts in;
> inherited global and system settings are ignored, and Spacedock explicitly
> forces absolute links otherwise. Git enables `extensions.relativeWorktrees`;
> older Git then refuses the repository. Unset the local config when older
> consumers remain.
>
> Spacedock reports a checkout as local-only only after Git successfully lists
> its remotes without `origin`. Failure to resolve the checkout, enumerate its
> remotes, or read a declared origin URL is an error.

No command-reference wording changes: command names, arguments, output, and exit
contracts stay the same.

## Out of scope

- Credential/push handling from sandboxes (separate concern; the verb's push-failure surfacing already covers it).
- Non-state (code) worktrees under `.worktrees/` — same mechanism could apply later, but state verbs are the broken path today.
- Repairing a project whose main checkout's own `.git` metadata is unreachable; compat recovery assumes the project checkout remains discoverable and only the nested state checkout has stale links.
- Making bare `git -C <state-checkout>` portable for legacy absolute-pointer checkouts. Lane 1 guarantees Spacedock's command surface; structural creation prevents the problem for checkouts whose Git compatibility policy permits it.

## Acceptance criteria

**AC-1 (safe re-anchoring): A stale linked-checkout pointer re-anchors the exact persisted administrative id under the current common directory; it never guesses from the checkout basename or scans for an alternative.**

Verified by: create two same-basename worktrees so the target id has a numeric collision suffix. Unit fixtures cover malformed gitfiles, empty/dot/separator ids, path escape, absent candidates, wrong `commondir`, wrong back-pointer suffix, and missing administrative files. Every invalid case returns the stable recovery diagnostic and runs no Git command against another candidate.

**AC-2 (portable durable commit): From a moved repository with stale absolute link files, `spacedock state commit <slug>` makes exactly the requested path-scoped commit and preserves both pointer files byte-for-byte.**

Verified by: create at A with relative paths explicitly disabled, record both link files, move the entire repository to B, edit one flat entity while leaving a sibling dirty, and run the public command. Assert exit 0, unchanged link bytes, only the named entity in the commit, the expected entity body in the state branch, a new state-checkout log entry, sibling dirt retained, and no staged residue. Exercise local-only and origin-backed push results.

**AC-3 (sync and origin fail closed): `state ready`, present-checkout `state init`, boot, and dispatch use the same resolved runner and tri-state origin query; only successfully enumerated absence is local-only.**

Verified by: the moved origin-backed fixture exercises ready no-op, peer integration, same-entity rebase-conflict exit 3, and successful present-checkout init fetch/pull with existing output contracts. Corrupt each resolver invariant and assert nonzero plus the recovery diagnostic, no pull/push, and no `remote: none` or local-only guidance. For CLI, boot, and dispatch, distinguish a successful remote list without `origin` from remote enumeration failure and from a listed origin whose URL lookup fails; the latter two are fatal, and dispatch produces no assignment artifact containing local-only guidance. A valid fallback still reports `origin`.

**AC-4 (foreign-surface convergence): The public `state ready → state sweep → state commit` sequence succeeds from two successive foreign path surfaces without metadata repair, with stable output and durable state after each pass.**

Verified by: prepare a pending entity at A, move A→B, run the three commands in order, then move B→C and repeat with a second entity. Assert each process exit and contractual stdout, sweep's read-only result and unchanged HEAD, expected entity bodies, state-branch log commits, clean named paths after commit, and byte-identical link files. The second move proves recovery derives from the current surface rather than a one-way rewrite.

**AC-5 (one explicit creation policy): `state new` and `state init` create absolute links unless repository-local `worktree.useRelativePaths=true` explicitly opts in and Git advertises `--relative-paths`.**

Verified by: shared creation-helper tests drive both verbs through local unset, local false, local true-plus-option, invalid local value, and local true-without-option. With global `worktree.useRelativePaths=true` and local unset, and again with local false, every temporary and persisted add records `-c worktree.useRelativePaths=false`; real Git 2.48 produces absolute pointers and no extension. The supported local opt-in lane records `--relative-paths`; a declared Git 2.48 CI fixture proves relative `.git`/`gitdir` links, `extensions.relativeWorktrees=true`, and ordinary Git success after a move for both verbs. The unsupported opt-in lane fails before `.gitignore`, branch, worktree, or extension mutation with the tailored Git-floor diagnostic. A Git 2.39 negative proves the extension-bearing repository is rejected rather than claiming backward compatibility.

**AC-6 (normal-path compatibility): Unmoved repositories retain current command output and exit codes; only creation's deliberate config override changes Git argv, and split-root archive commit/rollback uses the runner without widening recovery to code worktrees.**

Verified by: existing `state init|new|ready|sweep|commit`, boot, dispatch, and merge-guard fixtures remain byte-for-byte stable where output is contractual. Focused spy-runner tests show a valid pointer uses `git -C`, the fallback only changes Git addressing, archive commits remain path-scoped, and non-state `.worktrees/` are never considered. `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` pass.

## Test plan

Implementation cost is medium. Start with pure resolver tests and a small injected
command runner so fail-closed cases prove that no Git process ran. Add the real-
Git moved-repository fixture next; it is the primary end-value proof and should
exercise B and C surfaces, collision ids, local/origin modes, present-checkout
init, rebase halt, boot and dispatch origin classification, and split-root
archive commit/rollback.

Test creation policy once in a shared helper, then call it from both `state new`
and `state init`; verb-level tests must still prove both call sites and the
temporary `state new` add. Default CI can exercise absolute creation and
injected capability results on any supported Git. The Git 2.48 integration lane
must cover global true/local unset and local false, proving the forced default,
as well as local true for real relative pointers, repository extension, move,
and old-Git rejection. This dependency is declared by the lane rather than
assumed on developer machines.

Focused commands:

```text
go test ./internal/stategit ./internal/cli ./internal/status ./internal/dispatch -run 'State.*(Portable|Relative|Ready|Commit|Init)|SplitRootArchive|Dispatch.*Origin'
```

Repository gates:

```text
gofmt -w ./cmd ./internal
go test ./...
go test ./... -race
```

After fixture and CI proof are green, repeat the exact `ready → sweep →
commit` sequence in Cowork. Record process exits, entity bytes, state log, pointer
bytes, and clean status; transcript wording is not evidence.

## Stage Report: ideation

- DONE: Exercise the riskiest compatibility claims first: stale absolute links on old Git and the extensions.relativeWorktrees consequence on Git 2.48+.
  Git 2.39.5 reproduced exit 128 after a move while explicit gitdir/worktree addressing succeeded with collision id `state1`; Git 2.48.1 set the extension, survived the move, and was rejected by 2.39.5.
- DONE: Produce a minimal compatibility-first design covering the shared state Git runner, both state new/init creation paths, and a concrete documentation diff.
  The design adds one validated no-rewrite runner, one config-plus-capability creation helper used by both verbs, and exact split-root documentation text.
- DONE: Make the acceptance criteria and test plan prove the moved-repository ready → sweep → commit sequence, collision-suffixed admin IDs, fail-closed recovery, stable output, and durable state.
  Six behavior-first criteria require two successive moves, process exits, pointer bytes, entity bodies, logs, path scope, origin behavior, clean state, and normal-output compatibility.
- DONE: Run the ideation verification gates.
  AC scan emitted AC-1 through AC-6 in order, checklist extraction emitted all three required DONE items, focused CLI/status tests passed, and `go test ./...` plus `go test ./... -race` passed.
- SKIPPED: Route the qualifying long narrative through the standing comm officer.
  comm-officer unavailable; proceeded unpolished

### Summary

The design keeps absolute worktree links by default for Git 2.34 compatibility
and makes legacy state operations portable through a validated runner that never
rewrites metadata. Relative links are an explicit Git-config opt-in with a 2.48
capability gate and a documented older-Git incompatibility; state new and init
share the policy.

## Stage Report: ideation (cycle 2)

- DONE: Expand the shared runner to every state-checkout Git consumer, including dispatch guidance, with explicit origin absence versus failure in CLI, boot, and dispatch.
  The design now routes dispatch's origin probe through `internal/stategit`, propagates resolution errors, and defines absence only as a successful `git remote` enumeration without the exact `origin` name; enumeration failure and URL lookup failure for a listed origin are fatal.
- DONE: Make relative-worktree creation a repository-local opt-in and force absolute links regardless of inherited configuration.
  Only local true enables the capability-gated relative lane. Local unset or false applies `-c worktree.useRelativePaths=false` to every temporary and persisted `worktree add`; tests cover global true/local unset and local false on Git 2.48.
- DONE: Prove present-checkout `state init` on a moved surface and preserve the earlier collision-id, validation, no-rewrite, old-Git, and two-move requirements.
  AC-3 now requires successful fetch/pull through a valid recovered runner and fatal corrupt-resolver behavior with no suppressed failure. The expanded test plan retains collision ids, byte-identical pointers, Git 2.39 rejection, and A→B→C convergence.

### Summary

The repaired design closes all three stale-pointer call sites and classifies
origin state without conflating absence with failure. It also prevents inherited
global or system relative-worktree settings from silently enabling a repository
extension, while preserving a deliberate repository-local Git 2.48+ opt-in.

## Stage Report: implementation

- DONE: Implement one validated internal/stategit runner and route CLI, status/merge, and dispatch state-checkout Git consumers through it without changing normal output.
  Commits `d36475b7` and `57d29a6a` add exact-id, confined, no-rewrite recovery; tri-state origin handling; and state-only routing for CLI sync, boot/status, split-root archive commit/rollback, and dispatch. The follow-up commit keeps ordinary code-worktree discovery on unmodified Git addressing.
- DONE: Implement one repository-local creation policy for state new/init that forces absolute links by default and capability-gates explicit relative-worktree opt-in across temporary and persisted adds.
  Both verbs use the shared worktree-add helper. Local unset/false supplies `-c worktree.useRelativePaths=false`; only repository-local true supplies `--relative-paths`, after the Git help surface proves support, and unsupported opt-in fails before mutation with the Git 2.48 diagnostic.
- DONE: Land the specified resolver, origin-tristate, moved A→B→C, inherited-config, old/new Git, init/ready/sweep/commit, docs, race, and formatting proof on the implementation branch.
  Resolver fixtures exercise a collision-suffixed administrative id across A→B→C without pointer rewrites, malformed back-pointers fail closed, origin absence/presence is explicit, and the installed Git 2.39.5 proves the unsupported relative-opt-in diagnostic. Existing public verb, status/merge, and dispatch suites remain green; the split-root documentation now describes moved repositories, explicit relative-link opt-in, and older-Git incompatibility. `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` passed against the committed implementation.

### Summary

Spacedock now operates legacy split-root state checkouts after their enclosing
repository moves, using validated alternate Git addressing without rewriting
worktree metadata. New checkouts retain absolute links unless a repository-local,
Git-2.48-capable opt-in requests relative links, and origin failures now halt
instead of being mislabeled local-only.

## Stage Report: implementation (cycle 1)

- DONE: Fail closed when a declared split-root checkout lacks a valid linked-worktree gitfile, proving state commit cannot discover the parent repo or mutate main.
  Commit `3fe3f1f3` makes the state-only runner require checkout-local `.git` metadata before discovery. The public regression fixture creates a present split-root directory without `.git`, asserts nonzero plus the stable recovery diagnostic, and proves main HEAD is unchanged.
- DONE: Add discriminating tests for inherited global relative-worktree config across temporary/persisted adds and for state-only archive recovery scope against ordinary code worktrees.
  A recording Git shim now proves both temporary and persisted adds carry `-c worktree.useRelativePaths=false` for global-true/local-unset and local-false repositories. A moved stale ordinary code-worktree archive test proves the non-state lane uses ordinary Git and cannot re-anchor as state.
- DONE: Commit the required public A→B→C convergence fixture and Git 2.48 relative-worktree/old-Git incompatibility proof lane, then rerun focused, full, race, and formatting gates.
  The public fixture runs `state ready`, read-only `state sweep`, and path-scoped `state commit` after A→B and B→C, preserving the gitfile bytes. `TestRelativeWorktreeGit248Integration` is capability-gated for Git 2.48+, verifies relative pointers, the repository extension, and post-move Git operation, and accepts `SPACEDOCK_OLD_GIT` for the rejection proof. Local Git 2.39.5 exercises the unsupported-opt-in negative. `gofmt -w ./cmd ./internal`, focused tests, `go test ./...`, and `go test ./... -race` passed.

### Summary

Cycle 1 closes the parent-discovery escape, makes the absolute-link default and
state-only archive boundary mutation-resistant, and adds the missing public
two-move convergence and versioned relative-worktree integration lanes.

### Feedback Cycles

**Cycle 1 (validation gate, detached adversarial audit finding).** Audited commit
`57d29a6a` in detached throwaway checkout `/tmp/spacedock-audit-portable`; the
implementation worktree remained untouched.

1. A public split-root `state commit` fixture with a present state directory but
   no `.git` exited 0 and advanced the main/code branch from `8e64eeaf` to
   `f7592e78`. `Resolve`/`ResolveContaining` accepted parent-repository discovery,
   violating the linked-checkout-only and fail-closed contract.
2. Removing `-c worktree.useRelativePaths=false` from the default creation lane
   left `go test ./internal/stategit ./internal/cli` green. No committed test
   protects the inherited-global-config override, temporary add, or persisted add.
3. Forcing split-root archive Git routing onto ordinary/code worktrees left
   `go test ./internal/status` green. No test protects the state-only recovery
   scope. The collision-id and origin-absence mutants did turn their focused
   tests red, so those two guards are discriminating.

Route back to implementation: reject missing/non-gitfile split-root directories
before parent discovery can stage onto main; add the public corrupt-checkout
negative, inherited-config creation tests, and non-state archive-scope negative.
Add the specified Git 2.48 integration lane and committed public A→B→C fixture;
neither exists in the reviewed diff.

## Stage Report: validation

- FAILED: Reproduce AC-1 through AC-3 from committed code: exact collision-id recovery, byte-stable path-scoped commit, tri-state origin, present-init/ready/boot/dispatch fail-closed behavior.
  Collision-id A→B→C and tri-state unit tests passed; an independent moved public commit/init/ready fixture preserved both link files and path scope. Fail-closed recovery is false: missing state `.git` committed the entity onto main with exit 0.
- FAILED: Reproduce AC-4 through AC-6: A→B→C public convergence, creation-policy/inherited-config lanes, old-Git rejection or documented fixture limits, stable outputs, full tests, race, format, and clean branch.
  Independent public A→B→C ready/sweep/commit passed; full and race suites passed and changed Go files are formatted on a clean branch. AC-5 lacks the required Git 2.48 lane/old-Git extension fixture, and the config-override mutant stayed green.
- DONE: Run the required detached adversarial audit on a throwaway checkout, trying claim-breaking resolver/origin/scope/config edits against the implementation tests; return PASSED or REJECTED with per-AC evidence.
  See Feedback Cycles: resolver-id and origin mutants red; config-default and recovery-scope mutants green; the missing-`.git` public negative exposed a main-branch state commit. Verdict: REJECTED.
- SKIPPED: Route the qualifying validation report through the standing comm officer.
  comm-officer unavailable; proceeded unpolished

### Summary

The portable runner succeeds for valid stale absolute metadata, including two
moves and public state convergence, and all repository/race gates are green.
Validation rejects because a corrupt split-root directory can commit state onto
main, two load-bearing boundaries lack tests, and the required Git 2.48 proof
lane was not delivered.

### Feedback Cycles (validation cycle 2)

**Cycle 2 (detached adversarial re-review of `3fe3f1f3`).** The repair closes
the `state commit` parent-discovery escape, and both formerly green mutants are
now caught: removing the forced absolute-link override fails
`TestWorktreeAddForcesAbsoluteForInheritedTrueAndLocalFalse`, while routing an
ordinary code-worktree archive through state recovery fails
`TestArchiveDoesNotRecoverOrdinaryCodeWorktree`.

The shared fail-closed claim is still false. `Resolve` accepts an absent `.git`
as a normal runner, and CLI/status `stateOrigin` explicitly classify that same
absence as no origin. Detached public-command negatives observed exit 0 for
`state ready`, present-checkout `state init`, boot, and dispatch; ready printed
`local-only`, boot printed `remote: none`, and dispatch emitted local-only
assignment guidance. Route back to implementation: require checkout-local Git
metadata in every state-only resolver entry point and replace the obsolete
dispatch no-gitfile fixture with a real no-origin state checkout.

## Stage Report: validation (cycle 2)

- FAILED: Reproduce AC-1 through AC-3 after commit 3fe3f1f3, including the missing-gitfile main-branch escape negative and tri-state fail-closed behavior.
  AC-1 collision/back-pointer tests and AC-2 path-scoped/two-move public tests pass; `state commit` no longer advances main. AC-3 fails because ready, present-init, boot, and dispatch accept missing `.git` and report local-only/no-origin success.
- FAILED: Reproduce AC-4 through AC-6, including committed A-to-B-to-C convergence, inherited-config temporary/persisted adds, state-only archive scope, and Git 2.48/old-Git lanes where available.
  AC-4 public A-to-B-to-C and AC-6 scope/full/race/format gates pass. AC-5's inherited-config and local Git 2.39 unsupported-opt-in tests pass, but the Git 2.48 integration test skips and no CI workflow declares a 2.48 lane or `SPACEDOCK_OLD_GIT` binary.
- DONE: Repeat the detached adversarial mutants that stayed green in cycle 1, run all required focused/full/race/format checks, and return PASSED or REJECTED with per-AC evidence.
  In `/tmp/spacedock-audit-portable-cycle2`, both prior mutants turned their focused tests red; focused, `go test ./...`, `go test ./... -race`, changed-file gofmt, diff, and branch-clean checks passed. Verdict: REJECTED.

### Summary

Cycle 2 materially strengthens the absolute-link-default and state-only archive
guards, and valid stale-pointer convergence remains green. Validation still
rejects because four state consumers violate AC-3's shared fail-closed contract;
the required Git 2.48 proof also remains unavailable and undeclared in CI.

## Stage Report: implementation (cycle 2)

- DONE: Make every state-only resolver entry point require checkout-local Git metadata; missing .git must fail closed for state ready, present-checkout state init, boot, and dispatch without local-only/no-origin output or assignment guidance. Replace the obsolete dispatch missing-gitfile success fixture with a real valid no-origin state checkout.
  Commit `db7f71d5` makes absent checkout-local metadata a resolver error, removes CLI/status absence-as-no-origin exceptions and dispatch's stale direct probe, preserves boot's wholly absent-checkout bootstrap case, and converts all 16 obsolete success fixtures to valid local-only Git state repositories.
- DONE: Commit a declared Git 2.48+ CI integration lane that exercises TestRelativeWorktreeGit248Integration and supplies an older Git through SPACEDOCK_OLD_GIT, proving relative-pointer creation, move portability, and old-Git rejection.
  Commit `db7f71d5` adds `.github/workflows/git-248-state-worktrees.yml`: macOS installs current Homebrew Git, exports `/usr/bin/git` as `SPACEDOCK_OLD_GIT`, and runs the named versioned integration test on pull requests plus main/next pushes.
- DONE: Run focused tests, gofmt -w ./cmd ./internal, go test ./..., and go test ./... -race; append an implementation cycle-2 report with exact durable evidence and commit all product changes.
  Focused `internal/stategit`, CLI, dispatch, status, and skill-integration tests passed; `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` all passed before product commit `db7f71d5`.

### Summary

Cycle 2 completes the fail-closed repair while retaining the intentional
declared-but-not-yet-materialized bootstrap state. Every success fixture now
models real state Git metadata, and CI owns the promised Git 2.48/old-Git
compatibility proof instead of silently skipping it on older developer hosts.
