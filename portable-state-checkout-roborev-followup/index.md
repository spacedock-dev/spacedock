---
title: Close Roborev repository-boundary findings in portable state checkout
status: implementation
source: Follow-up to archived qw after Roborev thorough branch review job 488, 2026-07-14
started: 2026-07-13T16:28:38Z
completed:
verdict:
score: 0.95
worktree: .worktrees/spacedock-ensign-portable-state-checkout-roborev-followup
issue:
id: hc0fswq1kap2znjctzabx750
---

Repair the repository-boundary and commit-order defects found after `state-checkout-portable-gitdir` was archived prematurely following an invalid local merge attempt. That attempt was removed from `main` and never pushed. The reviewed range was `557f8df3e6a62d34987edda70533375fc48ba8f6..a70e9121f0707dfbee1e9d1341bac6acc951038e`; Roborev job `488` returned FAIL under the corrected repository guideline.

## Problem

The portable state-checkout change can still cross repository boundaries or leave state commits stranded. Split-root finalization may climb from a state directory missing `.git` into the code repository; an operable absolute gitfile can bypass identity and back-pointer validation and keep mutating the original repository after a copy; and `state commit` can create a local commit before discovering that `origin` is invalid, after which a retry returns early without pushing it.

## Compatibility spike

The highest-risk mechanism was reproduced in `/tmp` without changing project
metadata. A source repository had a linked `state` worktree and was copied while
the source remained reachable. The copied worktree's regular-file `.git` still
named the source administrative directory. Ordinary Git accepted it:

```text
git -C copy/state rev-parse --git-common-dir -> source/.git
git -C copy/state rev-parse --show-toplevel -> copy/state
```

Committing `copy/state/copied.txt` through that apparently operable gitfile
advanced the source `state` ref, left the copied repository's `state` ref
unchanged, and made the source worktree report `D copied.txt`. The source and
copied gitfiles, administrative directories, and indexes had different
filesystem identities. This proves that successful `git -C` and matching path
suffixes are insufficient. The regular-file lane must compare the recorded
administrative target with the candidate under the declared project's current
common directory, and compare an existing back-pointer target with the current
checkout gitfile, using filesystem identity. No further spike is needed: the
missing-metadata climb and origin-before-commit ordering are direct control-flow
defects visible in the reviewed code.

A second `/tmp` spike closed the trust-anchor question from staff review. The
declared project was itself a linked worktree, with a nested state worktree, and
the entire project was copied while its source and common repository remained
reachable. Both regular-file gitfiles in the copy pointed into the source common
directory, and ordinary Git succeeded from both copied paths. After removing
the source state gitfile, the state administrative `gitdir` back-pointer target
was nonexistent—the fallback case the moved-repository design permits—yet
ordinary Git still succeeded from the copied state path. The project
administrative back-pointer still resolved to the source project gitfile, whose
filesystem identity differed from the copied project gitfile. Therefore the
declared project must itself be validated before its common directory can serve
as the trust anchor; state validation alone cannot distinguish this copy.

## Proposed approach

### Strict state-checkout resolver

Make the declared project root and declared state checkout explicit inputs to
the state Git runner. Split-root mutation paths must never obtain their Git root
by walking upward from the entity directory. Resolve one runner before the
command's first mutation and reuse it for every subordinate Git operation.

Establish the project trust anchor without first trusting ordinary Git. A
project with an in-place `.git` directory anchors its own common directory. If
the declared project has a regular-file `.git`, parse and validate that linked
worktree metadata first: its recorded administrative target must exist, its
`commondir` must resolve to the same filesystem object as the candidate common
directory, and its existing administrative `gitdir` back-pointer target must be
the same filesystem object as the declared project's `.git` file. Do not permit
the stale-suffix fallback for the project trust anchor; without an existing
same-object back-pointer, ownership of an external common directory is not
proved. Only after this validation may the resulting common directory be used
to construct and validate the nested state candidate.

For a regular-file `.git`, remove the current "ordinary Git succeeds" shortcut.
Always parse the exact administrative id and construct the expected candidate
at `<declared-project-common-dir>/worktrees/<id>`. Apply these rules in order:

1. If the gitfile's recorded target exists, it must be the same filesystem
   object as the expected candidate. A copied checkout whose target still
   reaches the source therefore fails instead of being treated as a normal
   worktree.
2. If the recorded target is stale, the expected candidate may provide the
   existing moved-repository recovery lane. Its `HEAD`, `commondir`, `gitdir`,
   and administrative id remain mandatory and confined.
3. The resolved `commondir` must be the same filesystem object as the declared
   project's common directory. When the administrative `gitdir` back-pointer
   target exists, it must be the same object as the declared checkout's `.git`
   file. Only a nonexistent stale back-pointer may use the already-designed
   repository-relative suffix check.

Use filesystem identity rather than string equality for existing objects so a
bind mount, sandbox mount, or symlinked surface of the same repository remains
valid. Different copied objects fail even when their suffixes match. Preserve
the standalone `.git`-directory compatibility lane, but do not let a missing or
malformed declared `.git` fall through to the parent code repository.

Make the identity decision testable behind a small injected `sameObject(a, b)`
oracle. Its deterministic test supplies two distinct, non-symlink lexical paths
that the oracle marks as the same object and proves the resolver accepts them;
the negative marks identical-looking suffixes as different and proves refusal.
Add a Linux capability-gated real bind-mount test that presents the project and
state checkout at a second mounted path, asserts distinct path strings plus
equal device/inode identity, and runs the public state command successfully.
The dedicated CI step enables this test and treats inability to create the bind
mount as failure; ordinary developer runs may skip it when the capability flag
is absent. This is the evidence for the sandbox/bind-mount claim—symlink
coverage alone is not.

### Finalization preflight

For split-root merge finalization, resolve the project root from the workflow
definition with ordinary Git, then resolve the declared entity directory with
the strict state runner before clearing `mod-block`, terminalizing, or moving
the entity. Pass that bound runner into archive commit and rollback. Compute
split-root pathspecs relative to the declared state checkout; do not call
`FindGitRoot(entityDir)` in either path. An absent state `.git` therefore returns
the stable recovery diagnostic with the live entity, code branch, state ref,
indexes, and archive tree untouched. Inline workflows retain their ordinary
Git path.

### Origin-before-commit ordering

In `state commit`, resolve the state runner and classify the exact named
`origin` before staging or committing the entity. A successfully enumerated
remote list without `origin` remains the explicit local-only case. A listed
`origin` whose URL cannot be read is an error before mutation. Retain that one
classification for the later local-only or push path instead of probing again
after the commit. Once the URL is repaired, the still-dirty entity follows the
normal path-scoped commit and push path exactly once.

### Branch and integration path

Keep `spacedock-ensign/state-checkout-portable-gitdir` fixed at reviewed head
`a70e9121f0707dfbee1e9d1341bac6acc951038e`. Create the follow-up implementation
branch from that exact commit and append repair commits; do not cherry-pick,
squash, rebase, or rewrite the original four commits. The eventual PR targets
`main`, so its range contains both the original implementation and the repairs.
Record local `main` at
`557f8df3e6a62d34987edda70533375fc48ba8f6` before work and require its OID to
remain exactly that value throughout this task. There is no local-merge
fallback. If pushing the workflow-bearing branch remains blocked by OAuth's
missing `workflow` scope, stop at the integration gate and report that blocker
rather than merging locally or weakening the diff.

After formatting and full/race gates pass in a clean worktree, run one
replacement thorough Roborev branch review against `origin/main` under the
corrected repository guideline. Resolve and record immutable `base_sha` and
`head_sha`, require the Roborev job's stored range to equal
`base_sha..head_sha`, and inspect the stored prompt for the corrected
`Compatibility posture`, `Trust boundaries`, `Behavioral proof`, and `Review
focus` guideline sections. PASS without that prompt evidence is invalid. Any
finding returns the same branch to implementation for another cycle; it is not
approval evidence. Push that same exact HEAD, verify the remote branch and PR
head equal it, and record the PR base OID. If the PR base OID differs from
`base_sha` at gate presentation or integration, rerun the full/race gates and a
replacement Roborev review against the new exact base before proceeding. Open
the PR and leave integration to the PR merge ceremony.

## Concrete documentation diff

In `docs/site/advanced/split-root-state.md`, replace:

> Recovery validates the state worktree's existing administrative id and does
> not rewrite repository metadata.

with:

> Recovery validates the state worktree's administrative metadata against the
> current project and checkout without rewriting it. A copied checkout whose
> Git links still target a reachable source repository is refused rather than
> allowed to mutate the source.

Replace the final sentence:

> Failure to resolve the checkout, enumerate its remotes, or read a declared
> origin URL is an error.

with:

> Failure to resolve the checkout, enumerate its remotes, or read a declared
> origin URL is an error before `state commit` creates a local commit.

## Out of scope

- Rewriting or deleting the archived `qw` record.
- Hiding the prior rejected validation or Roborev review.
- Removing the Git 2.48 workflow file to bypass the OAuth `workflow`-scope requirement.
- Broad refactoring outside state Git resolution, finalization, and commit ordering.

## Acceptance criteria

**AC-1 (VALUE) - Split-root operations never mutate a parent or source repository when the declared project or state checkout is missing, copied, moved, or backed by an operable stale absolute gitfile.**
Verified by: a public-command matrix covering missing state `.git`, copied ordinary project plus state, copied linked project plus nested state, moved checkout, and a same-object mounted alias. The copied-linked-project negative leaves the permitted state back-pointer target nonexistent while its project back-pointer still identifies the source, proving project trust-anchor validation rejects it. Each negative records exact HEADs and state refs, index-file bytes, worktree registrations and status, gitfile/admin-file bytes, entity bytes and mode, and archive contents before invocation; it exits nonzero with every snapshot unchanged. The moved and same-object-mounted positives still commit only to the declared state ref.

**AC-2 - Merge finalization with a declared split-root state directory lacking `.git` fails closed instead of staging or committing the archived entity on the code branch.**
Verified by: a public merge-guard/finalization regression with a merge sentinel that removes the state gitfile, then records code HEAD, code index bytes and staged paths, state ref, live entity bytes/mode, and archive tree. The command exits 1 before clearing or terminalizing fields; the code and state refs, both indexes, live bytes, and archive tree remain exact.

**AC-3 - Every regular-file `.git` target is validated against the declared project even when raw Git commands would succeed.**
Verified by: copy-while-source-remains fixtures first prove raw `git -C copy status` and `git -C copy/state status` succeed. The linked-project variant removes the source state gitfile so the state back-pointer target is missing, then invokes public `state commit`; project trust validation still exits 1 before state fallback or staging. Source and copy code/state refs, both state index files, worktree registrations/status, pointer and back-pointer bytes, and entity bytes remain exact; the copied edit remains dirty. An injected distinct-path/same-object oracle and a capability-gated real Linux bind mount prove filesystem identity accepts a second path surface of the same checkout.

**AC-4 - `state commit` validates origin before local mutation and cannot strand a commit behind the no-op retry path.**
Verified by: a repository whose remote enumeration contains `origin` but whose URL lookup fails. The first public call exits 1 with local HEAD, state ref, index bytes, staged paths, entity bytes, and remote ref unchanged while preserving the dirty edit. The fixture fixes author/committer names, emails, and timestamps and invokes the retry with an explicit message. After adding the URL, the second call creates exactly one commit whose sole parent is the pre-failure HEAD, whose author/committer identities and dates equal the fixture values, whose message equals the explicit message byte-for-byte, whose changed-path set is exactly the requested entity path, and whose blob equals the dirty entity bytes. Local HEAD and the remote state ref equal that commit, with no staged residue. A separate no-origin fixture still creates one local-only commit with unchanged output and exit contract.

**AC-5 - The follow-up preserves the original reviewed history and can integrate only through a PR.**
Verified by: the original branch ref still equals `a70e9121f0707dfbee1e9d1341bac6acc951038e`; that commit is an ancestor of the follow-up HEAD; the four original commits and their trees are unchanged; local `main` remains exactly `557f8df3e6a62d34987edda70533375fc48ba8f6`; and the pushed remote branch plus PR head equal the recorded follow-up HEAD. A missing `workflow` OAuth scope blocks at push/PR rather than selecting a local fallback.

**AC-6 - The corrected exact range passes independent review and required repository gates.**
Verified by: record exact `base_sha` and `head_sha`; `gofmt -w ./cmd ./internal` leaves that head clean; `go test ./...` and `go test ./... -race` pass there; and a replacement thorough Roborev job records PASS for exactly `base_sha..head_sha`. The stored job prompt contains all four corrected-guideline sections. Remote branch and PR head equal `head_sha`, and the PR base OID equals `base_sha`; any base drift invalidates the evidence and requires fresh full/race plus Roborev proof. A FAIL is routed back to implementation and cannot advance or finalize the task.

## Test plan

Implementation cost is medium and fixture-heavy. Start with the three Roborev
findings as failing public-command regressions on the follow-up branch before
changing production code. Use real Git repositories for the boundary tests and
small injected filesystem/runner seams only for identity-error exhaustiveness.

The missing-metadata fixture drives merge finalization, not a helper. The copy
fixtures keep the source live and demonstrate that raw Git accepts both a copied
ordinary project and a copied linked project before Spacedock refuses them. In
the linked-project case, remove the source state gitfile so the state
back-pointer target is missing; the project trust-anchor check must still reject
the copy before that fallback. Snapshot helpers record SHA-256 bytes for every
index and metadata file in addition to exact refs and porcelain status.

Test the injected identity oracle with distinct, non-symlink path strings so the
resolver's decision is proven independently of host path canonicalization. Add
a capability-gated Linux integration test using a real bind mount and a
dedicated CI step that enables it and fails rather than skips when mount setup is
unavailable. Assert the two surfaces have distinct strings but the compared
objects have equal filesystem identity, then run the public command and verify
the expected state ref only.

The origin fixture creates a listed remote without a URL, proves failure before
the index or HEAD changes, repairs the URL, then proves the exact parent,
author, committer, timestamps, message, changed path, blob, local ref, and remote
ref of one durable pushed commit. Keep moved A-to-B recovery, no-origin
local-only, malformed metadata, and normal unmoved command tests green.

Run focused tests for `internal/stategit`, `internal/status` merge guard, and
`internal/cli` state commit first, then `gofmt -w ./cmd ./internal`,
`go test ./...`, and `go test ./... -race`. Confirm the worktree is clean and
record the exact SHA before the replacement Roborev review. Review
the immutable `base_sha..head_sha` as the complete branch behavior, not only the
repair commits. Inspect the stored job range and prompt, then fetch and compare
the PR's base OID before gate presentation and integration. Base drift requires
fresh full/race and Roborev evidence. Only atomic PASS evidence permits push and
PR creation; local `main` must retain its exact starting OID and no local merge
is part of this plan.

### Feedback Cycles

**Cycle 1 (ideation staff review).** Independent review rejected the first
follow-up design on four false-green or incomplete boundaries:

- The resolver validates the state checkout but still trusts a declared
  project's regular-file `.git` through ordinary Git. Add the copied linked
  project case, including the allowed missing state back-pointer fallback, and
  prove the project trust anchor cannot cross into the source repository while
  same-object path aliases remain valid.
- Integration evidence must require local `main`'s exact OID to remain unchanged;
  “no merge commit” does not exclude fast-forward, squash, cherry-pick, or direct
  ref movement. The Roborev record must bind exact base and head SHAs and prove
  the job prompt contains the corrected guideline; rerun when the PR base differs.
- A symlink does not prove bind/sandbox-mount compatibility. Require a
  deterministic distinct-path same-object identity oracle plus a
  capability-gated real mount check, or narrow the compatibility claim.
- The state-commit proof must pin exact parent, author/committer attribution,
  message, and changed-path set so sibling-sweeping or misattributed commits
  cannot pass.

Routed back to ideation: revise the design, AC-4 through AC-6, and the test plan
together, then return for another independent staff review before the captain
gate.

## Stage Report: ideation

- DONE: Reproduce the operable-stale-gitfile boundary crossing and identify the discriminating validation rule.
  A `/tmp` copy-while-source-remains spike showed ordinary Git accepting the copied worktree, advancing the source state ref, and corrupting the source worktree view. Filesystem identity distinguishes the source administrative target from the copied project's candidate while preserving aliases of the same objects.
- DONE: Design repairs for all three exact Roborev job 488 findings.
  The design makes project and state roots explicit, validates every regular-file gitfile before use, preflights split-root finalization before any entity mutation, and classifies origin before staging or committing.
- DONE: Define semantic adversarial tests that assert exact refs, indexes, bytes, attribution, and cleanup state.
  The missing-metadata, live-source-copy, and invalid-origin journeys drive public commands and distinguish safe refusal from false-green error-substring tests; moved, alias, and local-only positives protect compatibility.
- DONE: Specify a PR-only correction path that preserves the original reviewed branch history.
  The follow-up descends from immutable head `a70e9121`, requires full/race and corrected-guideline Roborev PASS at one exact SHA, pushes that SHA, and opens a PR. OAuth failure blocks; local merge, squash, rebase, and history rewriting are excluded.

### Summary

The repair validates a state checkout against its declared project even when
ordinary Git can follow a stale absolute pointer, and it moves both merge
finalization and origin inspection ahead of mutation. Public adversarial tests
pin exact repository state, while an append-only follow-up branch and exact-SHA
Roborev gate provide a PR-only route to integration.

## Stage Report: ideation (cycle 2)

- DONE: Close the declared-project trust-anchor gap, including the state fallback that could otherwise mask it.
  A copied linked-project spike kept the source common directory live, removed the source state gitfile so the nested state back-pointer target was missing, and proved ordinary Git still accepted both copied paths. The revised design validates a regular-file project gitfile and its existing same-object back-pointer before trusting its common directory; the state stale-back-pointer fallback cannot run first.
- DONE: Replace path-shape alias evidence with deterministic identity and real-mount proof.
  The test plan now requires an injected distinct-path/same-object oracle plus a capability-gated Linux bind-mount test in a dedicated CI step that fails rather than skips when enabled. Symlinks are not the acceptance oracle for sandbox-mount compatibility.
- DONE: Make origin-retry and integration evidence semantically exact.
  AC-4 pins the successful commit's sole parent, author/committer identities and timestamps, exact message, sole changed path, blob, and local/remote refs. AC-5 and AC-6 keep local main at its exact starting OID and bind gates, corrected-guideline prompt, Roborev range, remote head, PR head, and PR base to recorded SHAs, with mandatory rerun on base drift.

### Summary

The cycle-2 design validates the project before using it as the state resolver's
trust anchor, proves mount aliases through filesystem identity and a real bind
mount, and makes both the origin retry commit and PR/Roborev evidence atomic.
The immutable qwp head and PR-only, no-local-merge path remain unchanged.

## Stage Report: implementation

- DONE: Implement the project/state identity, finalization-preflight, and origin-before-commit repairs with public adversarial regressions for every Roborev job 488 finding and the staff-review trust-anchor gap.
  Commit `2ca163b8ff0ff350799ac412b5fdeda46925e504` validates linked-project and state metadata by filesystem identity, preflights split-root finalization, and classifies the exact named origin before staging. Public regressions cover live-source copies, missing state metadata, origin failure/retry, moved repositories, injected same-object aliases, and a capability-gated real bind mount.
- DONE: Preserve immutable qwp history and exact local main OID while committing only the follow-up branch; make commit attribution, changed paths, refs, output, and cleanup assertions semantically exact.
  Original head `a70e9121f0707dfbee1e9d1341bac6acc951038e` remains the follow-up's ancestor and the original branch still names it; local `main` and `origin/main` both remain `557f8df3e6a62d34987edda70533375fc48ba8f6`. Tests pin source/copy refs, indexes, metadata and entity bytes, retry parent and attribution, exact message and sole path, remote ref, and staged cleanup.
- DONE: Apply the documented compatibility behavior and run gofmt, focused tests, go test ./..., and go test ./... -race, leaving a committed exact head ready for fresh validation and Roborev review.
  The split-root guide and dedicated Linux bind-mount CI lane were updated. `gofmt -w ./cmd ./internal` left exact head `2ca163b8ff0ff350799ac412b5fdeda46925e504` clean; focused suites, `go test ./...`, and `go test ./... -race` passed.

### Summary

State Git operations now bind one validated runner to the declared project and checkout instead of trusting an operable stale pointer or walking into a parent repository. The repair is committed only on the append-only follow-up branch, with immutable qwp ancestry and local main unchanged; it is ready for fresh independent validation and exact-range Roborev review, not integration.

### Feedback Cycles

**Cycle 2 (validation Roborev and detached adversarial audit).** Exact-range
Roborev job `716` rejected
`557f8df3e6a62d34987edda70533375fc48ba8f6..2ca163b8ff0ff350799ac412b5fdeda46925e504`
under the corrected guideline:

- High: `stategit.Resolve` accepts a checkout with a standalone `.git`
  directory before validating its relationship to the declared project. A
  workflow-local symlink can therefore select an unrelated repository for
  `state commit` or merge finalization. The detached validation audit added an
  uncommitted public-command regression in a throwaway checkout; changing the
  workflow's state path to a symlink into an unrelated repository made
  `state commit` exit 0 and advance that unrelated HEAD. Add a fail-closed
  symlink-escape regression that snapshots both repositories and validate
  containment before accepting the standalone-directory lane.
- Medium: the existing-checkout path in `state init` now turns fetch or pull
  failures into exit 1, where the reviewed base treated refresh as best-effort
  and retained idempotent success. Preserve the existing public exit contract,
  or explicitly scope and prove a compatibility change.

The same detached audit's copied-linked-project public-command journey was
clean: with the source live and the nested state back-pointer target removed,
the command rejected before mutation and preserved exact source/copy refs,
gitfiles, and dirty entity bytes. No validation edit entered the implementation
worktree. The audit worktree was removed after the run.

**Cycle 3 (Roborev-first validation, job 758; captain escalation).** The
corrected-guideline first gate rejected exact
`557f8df3e6a62d34987edda70533375fc48ba8f6..f91863463a7058cf1b98f20ffdae868173c37695`
before tests, push, PR creation, or CI:

- HIGH: linked-worktree containment is lexical; a workflow-local symlink can
  select an external linked worktree and mutate it.
- HIGH: a checkout-local `.git/commondir` can redirect the standalone lane to
  another repository's refs and objects.
- MEDIUM: `status --boot` discards state resolver/origin errors and exits without
  a diagnostic.

Per the three-cycle feedback rule, this rejection is escalated to the captain
instead of automatically routing another implementation round. The follow-up
branch remains local at `f9186346`; local `main` is unchanged; no PR or CI exists
for this branch.

Captain decision: send cycle 3 back without reframing. The linked-worktree
symlink and standalone `commondir` escapes are direct violations of the existing
declared-project containment invariant; the silent boot diagnostic is an
adjacent regression of the same resolver routing. Fix and test locally. The next
fresh validator must again run corrected exact-range Roborev first; no push, PR,
or CI is permitted before PASS.

**Cycle 4 (Roborev-first validation, job 855).** Corrected-guideline Roborev
rejected exact
`557f8df3e6a62d34987edda70533375fc48ba8f6..f63f71eb4ee41b227e24ca288d0aaffc0366fb54`
before tests, push, PR creation, or CI:

- HIGH: inherited `GIT_DIR`, `GIT_WORK_TREE`, `GIT_COMMON_DIR`, and
  `GIT_INDEX_FILE` can redirect subprocesses after path validation.
- HIGH: effective inherited, included, global, or system `core.worktree` can
  redirect the standalone-repository runner.
- MEDIUM: project-gitfile validation rejects valid separate-git-dir and
  submodule layouts even though only the linked-worktree lane requires the
  stricter relationship check.
- LOW: boot can print a returned `computeNextID` error twice.

Captain standing decision: send these findings back without reframing. They
remain within the existing trust-boundary and compatibility scope. Sanitize Git
routing authority, validate configuration provenance without rejecting supported
layouts, de-duplicate the boot error, and add exact zero-mutation regressions.
The next fresh validator must again run corrected exact-range Roborev first; no
push, PR, or CI is permitted before PASS.

**Cycle 5 (Roborev-first validation, job 894).** Corrected-guideline thorough
Roborev rejected exact
`557f8df3e6a62d34987edda70533375fc48ba8f6..3c842d1713411b41ae7f074a3ced8e48d4ac820c`
before tests, push, PR creation, or CI:

- MEDIUM: `Runner.Combined` does not set `cmd.Dir` to the validated checkout,
  so relative remote URLs resolve from the caller's process directory.
- MEDIUM: writable Git routing variables including `GIT_OBJECT_DIRECTORY` and
  `GIT_SHALLOW_FILE` remain inherited.
- MEDIUM: `Origin` can map checkout disappearance after `Resolve` to
  no-origin/local-only success instead of failing closed.

Captain standing decision: send these findings back without reframing. They are
direct continuations of the existing routing-authority, compatibility, and
fail-closed scope. Fix them with public-command zero-mutation regressions. The
next fresh validator must again run corrected exact-range Roborev first; no
push, PR, or CI is permitted before PASS.

## Stage Report: validation

- FAILED: AC-1 — refuse every cross-repository state target before mutation.
  The missing-metadata and live-source copy matrices passed, but the detached
  public-command audit reproduced the standalone-`.git` symlink escape from
  Roborev job `716`: an unrelated repository's HEAD advanced.
- DONE: AC-2 — fail split-root merge finalization before entity or repository mutation when state `.git` is missing.
  `TestMergeGuardMissingStateGitFailsBeforeFinalization` passed and pinned code
  and state refs, both indexes, entity bytes and mode, worktree registrations,
  and the absent archive tree.
- DONE: AC-3 — reject operable stale regular-file gitlinks against the declared project.
  The focused ordinary-copy and copied-linked-project tests passed. A detached
  public `state commit` journey independently covered the linked-project case
  with the permitted missing nested back-pointer fallback and exact cleanup
  snapshots.
- DONE: AC-4 — classify the exact named origin before committing and produce one exact retry commit.
  `TestStateCommitReadsDeclaredOriginBeforeMutation` and the no-origin fixture
  passed, including exact parent, attribution, timestamps, message, sole path,
  blob, local/remote refs, and empty staged state.
- FAILED: AC-5 — integrate only through an exact-SHA PR while preserving history.
  The local invariants passed: original qwp remains `a70e9121`, follow-up HEAD
  remains `2ca163b8`, and both local `main` and `origin/main` remain `557f8df3`.
  Roborev failed, so the branch was not pushed and no PR was opened; the remote
  branch remains absent and no local-merge fallback was attempted.
- FAILED: AC-6 — pass every repository gate and corrected exact-range review at one head.
  `gofmt -w ./cmd ./internal` left the implementation worktree clean;
  focused adversarial tests, `go test ./...`, and `go test ./... -race` passed.
  The real bind-mount tests are correctly capability-gated to the dedicated
  Ubuntu workflow and cannot run on this macOS validator; that remote lane was
  not entered because the review failed before push. Roborev job `716` stored
  the exact `base_sha..head_sha` range and a prompt containing `Compatibility
  posture`, `Trust boundaries`, `Behavioral proof`, and `Review focus`, but
  returned FAIL with the high and medium findings above.

### Summary

Independent validation rejects exact head `2ca163b8`. Core copied-checkout,
missing-metadata, origin-ordering, formatting, full, and race evidence is green,
but the detached audit confirms a cross-repository symlink escape and Roborev
also identifies an existing-checkout `state init` compatibility regression.
No push, PR, ref movement, implementation edit, or local merge occurred.

## Stage Report: implementation (cycle 2)

- DONE: Reject standalone .git directories that escape the declared project through symlinks or unrelated repositories, with public-command snapshots proving zero mutation in both repositories.
  Commit `f91863463a7058cf1b98f20ffdae868173c37695` resolves the declared workflow, checkout, and standalone Git directory through filesystem identity before accepting the compatibility lane. `TestStateCommitRejectsStandaloneRepositorySymlinkEscape` uses an unrelated nested repository outside the workflow root and pins both HEADs, indexes, statuses, entity bytes and mode, stdout, and the symlink; `TestResolveRejectsStandaloneGitDirectorySymlinkEscape` covers the adjacent `.git`-only symlink.
- DONE: Restore the existing-checkout state init best-effort refresh exit contract unless an explicit compatibility change is separately scoped and proven.
  Existing-checkout metadata still fails closed, while fetch failure skips pull and pull failure is ignored exactly as the reviewed base did. `TestStateInitExistingCheckoutRefreshIsBestEffort` proves both paths retain exit 0, empty stderr, exact success stdout, and the expected fetch/pull cardinality.
- DONE: Run focused, full, and race tests and commit locally without pushing; leave an exact clean head for fresh corrected-guideline Roborev-first validation before any CI engagement.
  Focused `internal/cli`, `internal/stategit`, and adjacent `internal/status` tests passed; `gofmt -w ./cmd ./internal` left exact head `f91863463a7058cf1b98f20ffdae868173c37695` clean; `go test ./...` and `go test ./... -race` passed. The code branch was not pushed and no PR, CI, merge, approval, or integration action was attempted.

### Summary

The standalone compatibility lane now accepts only a state checkout and Git directory physically contained by the declared workflow, closing both whole-checkout and `.git`-only symlink escapes. The prior idempotent `state init` refresh contract is restored, and exact local head `f9186346` is ready for fresh corrected-guideline Roborev-first validation.

## Stage Report: validation (cycle 2)

- FAILED: Run corrected-guideline Roborev first on exact range 557f8df3..f9186346 and inspect its stored range/prompt; on any finding stop REJECTED without push, PR, or CI.
  Job `758` stored exact range `557f8df3e6a62d34987edda70533375fc48ba8f6..f91863463a7058cf1b98f20ffdae868173c37695`, included all four corrected-guideline sections, and returned FAIL with two high and one medium findings. Job `751` was discarded because positional inclusive-range syntax stored `aeeadf45..f9186346` instead of the required Git difference.
- SKIPPED: Only after Roborev PASS independently reproduce the symlink-escape and best-effort-init fixes, AC-1 through AC-6, gofmt cleanliness, focused tests, full/race suites, real-mount evidence, and detached audit.
  Roborev failed the first gate, so no local test, formatting, real-mount, or detached-audit work ran in this cycle; prior-cycle evidence was not promoted across the failed exact head.
- SKIPPED: Only after every local gate passes push the exact reviewed head, open and verify the PR, then engage required CI; OAuth workflow-scope or CI failure blocks with no local integration.
  No push, PR creation, CI trigger or approval, merge, rebase, ref movement, or implementation edit was attempted.

### Summary

Validation recommends REJECTED for exact head `f91863463a7058cf1b98f20ffdae868173c37695`. Roborev job `758` found that a symlink to an external linked worktree can escape the declared project, a checkout-local `.git/commondir` can redirect standalone operations into another repository, and `status --boot` can discard resolver errors and exit silently. AC-1 and AC-6 therefore fail; AC-2 through AC-5 were not revalidated because the mandated first gate stopped the cycle. Local `main` remained exactly `557f8df3e6a62d34987edda70533375fc48ba8f6`, and the implementation worktree remained clean at the reviewed head.

## Stage Report: implementation (cycle 3)

- DONE: Canonicalize and identity-check declared project/worktree containment so symlinks to external linked worktrees fail before mutation while real same-object mount aliases remain supported.
  Commits `bee4ddba2bdf68f5927cc706e68947da9578cb2d` and `f63f71eb4ee41b227e24ca288d0aaffc0366fb54` canonicalize the original declared workflow and checkout before containment while retaining filesystem identity for administrative metadata and bind aliases. `TestStateCommitRejectsLinkedWorktreeSymlinkEscape` drives the public command and pins code/external HEADs, indexes, statuses, gitfile, entity bytes, stdout, and error output.
- DONE: Validate standalone checkout commondir/worktree ownership so checkout-local metadata cannot redirect refs or objects outside the declared checkout, with public zero-mutation regressions.
  Standalone resolution requires the common directory to be the same filesystem object as `.git` and any configured `core.worktree` to be the checkout itself; normal configs avoid a multiplicative Git subprocess. `TestStateCommitRejectsStandaloneExternalCommondir` proves raw Git follows the external common directory before the public command refuses with both repositories and local metadata unchanged; `TestResolveRejectsStandaloneExternalCoreWorktree` covers the adjacent worktree redirect.
- DONE: Render boot resolver/origin failures deterministically, run focused/full/race tests, and commit locally without push; leave an exact head for fresh Roborev-first validation before PR or CI.
  Text and JSON boot now emit exact `Error: ...` stderr with exit 1 and empty stdout. Focused suites and `go test ./...` passed at exact clean head `f63f71eb4ee41b227e24ca288d0aaffc0366fb54`; the aggregate race run passed every package except one unrelated pre-existing `TestSonnetTeamDeleteHangReplay` replay flake, which then passed both isolated under `-race` and on a full `internal/ensigncycle -race` package rerun. No code push, PR, CI, merge, approval, or integration action occurred.

### Summary

Linked and standalone state checkouts can no longer redirect mutation through path or Git-metadata aliases outside the declared workflow, and boot surfaces resolver failures instead of failing silently. The append-only local head is ready for fresh corrected-guideline Roborev-first validation; the race replay flake and its successful exact-head reruns are recorded explicitly rather than hidden.

## Stage Report: validation (cycle 3)

- FAILED: Run corrected-guideline Roborev first on exact range 557f8df3..f63f71eb and inspect its stored range/prompt; on any finding stop REJECTED without tests, push, PR, or CI.
  Roborev job `855` stored exact range `557f8df3e6a62d34987edda70533375fc48ba8f6..f63f71eb4ee41b227e24ca288d0aaffc0366fb54`, reviewed all eight commits in that Git difference, and included `Compatibility posture`, `Trust boundaries`, `Behavioral proof`, and `Review focus`. It returned FAIL with two high, one medium, and one low finding: inherited Git-routing environment variables can redirect validated operations; effective inherited or included `core.worktree` can redirect standalone operations; linked-worktree-only project-gitfile validation rejects separate-git-dir and submodule layouts; and boot can print a returned `computeNextID` error twice.
- SKIPPED: Only after Roborev PASS independently reproduce both external-redirect refusals and boot diagnostics, all ACs, gofmt, focused/full/race suites, real-mount evidence, and detached mutation audit.
  Roborev failed the mandatory first gate. No local test, formatting, real-mount, or detached-audit command ran, and the implementation report's earlier evidence was not promoted to this exact-head validation.
- SKIPPED: Only after every local gate passes push exact head f63f71eb, open and verify the PR, then engage required CI; report actual outcomes and block on OAuth, waiting, skipped, or failed evidence.
  No code push, PR creation, CI trigger or approval, merge, rebase, ref movement, or implementation edit was attempted. The remote follow-up branch remains absent; local `main` and `origin/main` remain exactly `557f8df3e6a62d34987edda70533375fc48ba8f6`, immutable qwp remains `a70e9121f0707dfbee1e9d1341bac6acc951038e`, and the clean implementation worktree remains at `f63f71eb4ee41b227e24ca288d0aaffc0366fb54`.

### Summary

Validation recommends REJECTED for exact head `f63f71eb`. Roborev job `855`
shows that validated paths can still be overridden by inherited Git routing and
configuration, while the project-gitfile lane regresses supported Git layouts
and one boot error can be duplicated. AC-1, AC-3, and AC-6 therefore fail; the
remaining criteria were not revalidated because the mandated first gate stopped
the cycle. No tests, push, PR, CI, implementation mutation, or local merge
occurred.

## Stage Report: implementation (cycle 4)

- DONE: Sanitize inherited Git repository-routing authority at every state-Git subprocess boundary and prove hostile values cannot mutate another repository.
  Commit `3c842d1713411b41ae7f074a3ced8e48d4ac820c` removes inherited `GIT_DIR`, `GIT_WORK_TREE`, `GIT_COMMON_DIR`, `GIT_INDEX_FILE`, and command-scoped Git config before invoking Git, then supplies only the checkout and Git directory validated by `Resolve`. `TestStateCommitSanitizesHostileGitRoutingAndConfig` drives the public command separately under each routing variable and under command-scoped, global, and system `core.worktree`; each case commits only the intended state checkout while pinning the unrelated repository's HEAD, index, status, and bytes plus the code repository's HEAD, index, and status.
- DONE: Validate standalone `core.worktree` provenance without rejecting valid separate-git-dir or submodule project layouts.
  Standalone resolution uses Git's canonical include-aware parser only when local config can name `worktree` or an include, rejects a `core.worktree` sourced outside the validated Git directory, and requires every value to resolve to the checkout. Adjacent tests accept direct and checkout-local included values, reject an external include even when its value names the checkout, and exercise operational runners from real separate-git-dir and submodule project gitfiles. The linked-worktree commondir/back-pointer relationship check now remains specific to the linked lane; self-contained project Git directories retain structural validation.
- DONE: De-duplicate boot's returned next-ID diagnostic, run formatting and all local gates, and keep the implementation local for fresh Roborev-first validation.
  `gatherBoot` now returns minting errors without printing them, leaving the text and JSON callers as the single diagnostic owners; a focused test pins empty callee stderr and the exact returned error. `gofmt -w ./cmd ./internal`, all focused suites, and exact-head `go test ./...` passed. The aggregate `go test ./... -race` passed every package except the same unrelated `TestSonnetTeamDeleteHangReplay` replay-fixture flake recorded in cycle 3; that test then passed in isolation under race, the full `internal/ensigncycle -race` package rerun passed, and the changed `internal/stategit` package passed a final exact-head race rerun. No code push, PR, CI, merge, approval, or integration action occurred.

### Summary

Every state Git subprocess now operates through repository and work-tree paths owned by the validated runner, so inherited routing and effective `core.worktree` cannot supersede the trust decision. Include provenance is checked without regressing separate-git-dir or submodule projects, and boot errors have one diagnostic owner. Clean local head `3c842d17` is ready for a fresh corrected-guideline exact-range Roborev-first validation; the branch remains unpushed and no PR or CI exists for it.

## Stage Report: validation (cycle 4)

- FAILED: Run corrected-guideline Roborev first on exact range 557f8df3..3c842d17 and inspect the stored range and all four guideline sections; on any finding stop REJECTED without tests, push, PR, or CI.
  Roborev job `894` stored exact range `557f8df3e6a62d34987edda70533375fc48ba8f6..3c842d1713411b41ae7f074a3ced8e48d4ac820c`, used thorough reasoning, and included `Compatibility posture`, `Trust boundaries`, `Behavioral proof`, and `Review focus`. It returned FAIL with three medium findings: `Runner.Combined` no longer gives relative remote URLs the checkout as their resolution directory; writable routing variables including `GIT_OBJECT_DIRECTORY` and `GIT_SHALLOW_FILE` remain inherited; and `Origin` can classify a checkout that disappears after resolution as local-only success.
- SKIPPED: Only after Roborev PASS independently reproduce hostile Git environment/config refusal, supported separate-git-dir/submodule layouts, single boot diagnostics, every AC, gofmt, focused/full/race suites, and detached zero-mutation evidence.
  Roborev failed the mandatory first gate. No local test, formatting, real-mount, race, or detached-audit command ran, and no prior-cycle test or flake attribution was promoted to this exact head.
- SKIPPED: Only after every local gate passes push exact head 3c842d17, open and verify the PR, then engage required CI; report actual outcomes and block on OAuth, waiting, skipped, unapproved, or failed evidence.
  No code push, PR creation, CI trigger or approval, merge, rebase, ref movement, or implementation edit was attempted. The remote follow-up branch remains absent; local `main` and `origin/main` remain exactly `557f8df3e6a62d34987edda70533375fc48ba8f6`, immutable qwp remains `a70e9121f0707dfbee1e9d1341bac6acc951038e`, and the clean implementation worktree remains at `3c842d1713411b41ae7f074a3ced8e48d4ac820c`.

### Summary

Validation recommends REJECTED for exact head `3c842d17`. Roborev job `894`
found that the sanitized runner still mishandles relative remotes, inherits
writable object-routing authority, and can silently downgrade a vanished
checkout to local-only. AC-1, AC-4, and AC-6 therefore fail; the remaining
criteria were not revalidated because the mandatory first gate stopped the
cycle. No tests, push, PR, CI, implementation mutation, or local merge occurred.
