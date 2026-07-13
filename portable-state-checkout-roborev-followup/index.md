---
title: Close Roborev repository-boundary findings in portable state checkout
status: ideation
source: Follow-up to archived qw after Roborev thorough branch review job 488, 2026-07-14
started: 2026-07-13T16:28:38Z
completed:
verdict:
score: 0.95
worktree:
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

## Proposed approach

### Strict state-checkout resolver

Make the declared project root and declared state checkout explicit inputs to
the state Git runner. Split-root mutation paths must never obtain their Git root
by walking upward from the entity directory. Resolve one runner before the
command's first mutation and reuse it for every subordinate Git operation.

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
There is no local-merge fallback. If pushing the workflow-bearing branch remains
blocked by OAuth's missing `workflow` scope, stop at the integration gate and
report that blocker rather than merging locally or weakening the diff.

After formatting and full/race gates pass in a clean worktree, run one
replacement thorough Roborev branch review against `origin/main` under the
corrected repository guideline. Record the exact reviewed HEAD and require
PASS. Any finding returns the same branch to implementation for another cycle;
it is not approval evidence. Push that same exact HEAD, verify the remote branch
and PR head equal it, open the PR, and leave integration to the PR merge
ceremony.

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

**AC-1 (VALUE) - Split-root operations never mutate a parent or source repository when the declared state checkout is missing, copied, moved, or backed by an operable stale absolute gitfile.**
Verified by: a public-command matrix covering missing state `.git`, copy-while-source-remains, moved checkout, and a same-object alias. Each negative records exact HEADs and state refs, index-file bytes, worktree registrations and status, gitfile/admin-file bytes, entity bytes and mode, and archive contents before invocation; it exits nonzero with every snapshot unchanged. The moved and same-object-alias positives still commit only to the declared state ref.

**AC-2 - Merge finalization with a declared split-root state directory lacking `.git` fails closed instead of staging or committing the archived entity on the code branch.**
Verified by: a public merge-guard/finalization regression with a merge sentinel that removes the state gitfile, then records code HEAD, code index bytes and staged paths, state ref, live entity bytes/mode, and archive tree. The command exits 1 before clearing or terminalizing fields; the code and state refs, both indexes, live bytes, and archive tree remain exact.

**AC-3 - Every regular-file `.git` target is validated against the declared project even when raw Git commands would succeed.**
Verified by: a copy-while-source-remains fixture first proves raw `git -C copy/state status` succeeds, then edits a copied entity and invokes public `state commit`. It exits 1 before staging. Source and copy code/state refs, both state index files, worktree registrations/status, pointer and back-pointer bytes, and entity bytes remain exact; the copied edit remains dirty. A same-object alias positive proves filesystem identity does not reject a second path surface of the same checkout.

**AC-4 - `state commit` validates origin before local mutation and cannot strand a commit behind the no-op retry path.**
Verified by: a repository whose remote enumeration contains `origin` but whose URL lookup fails. The first public call exits 1 with local HEAD, state ref, index bytes, staged paths, entity bytes, and remote ref unchanged while preserving the dirty edit. After adding the URL, the second call creates exactly one commit whose tree contains the exact entity bytes, pushes local HEAD to the remote state ref, and leaves no staged residue. A separate no-origin fixture still creates one local-only commit with unchanged output and exit contract.

**AC-5 - The follow-up preserves the original reviewed history and can integrate only through a PR.**
Verified by: the original branch ref still equals `a70e9121f0707dfbee1e9d1341bac6acc951038e`; that commit is an ancestor of the follow-up HEAD; the four original commits and their trees are unchanged; no merge commit is created on local `main`; and the pushed remote branch plus PR head equal the recorded follow-up HEAD. A missing `workflow` OAuth scope blocks at push/PR rather than selecting a local fallback.

**AC-6 - The corrected exact range passes independent review and required repository gates.**
Verified by: `gofmt -w ./cmd ./internal` leaves the exact head clean, `go test ./...` and `go test ./... -race` pass there, and a replacement thorough Roborev branch job against `origin/main` under the corrected guideline records PASS for that same SHA. A FAIL is routed back to implementation and cannot advance or finalize the task.

## Test plan

Implementation cost is medium and fixture-heavy. Start with the three Roborev
findings as failing public-command regressions on the follow-up branch before
changing production code. Use real Git repositories for the boundary tests and
small injected filesystem/runner seams only for identity-error exhaustiveness.

The missing-metadata fixture drives merge finalization, not a helper. The copy
fixture keeps the source live and demonstrates that raw Git accepts the bad
pointer before Spacedock refuses it. Its snapshot helper records SHA-256 bytes
for each index and metadata file in addition to exact refs and porcelain status.
The origin fixture creates a listed remote without a URL, proves failure before
the index or HEAD changes, repairs the URL, then proves one durable pushed
commit. Keep moved A-to-B recovery, no-origin local-only, malformed metadata,
and normal unmoved command tests green. Exercise a symlinked same-object surface
in default CI; a bind-mount live check is useful where available but is not the
only acceptance proof.

Run focused tests for `internal/stategit`, `internal/status` merge guard, and
`internal/cli` state commit first, then `gofmt -w ./cmd ./internal`,
`go test ./...`, and `go test ./... -race`. Confirm the worktree is clean and
record the exact SHA before the replacement Roborev review. Review
`origin/main..<exact-head>` as the complete branch behavior, not only the repair
commits. Only PASS at that SHA permits push and PR creation; no local merge is
part of this plan.

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
