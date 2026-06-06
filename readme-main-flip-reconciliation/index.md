---
id: nbz0yjvmqm6gda6csw8yef7k
title: README reconciliation for main flip and 0.20.0 install paths
status: implementation
source: "captain (2026-06-06) - before flipping main, reconcile README/install docs; consider existing README PRs #213 and #220."
score: "0.39"
started: 2026-06-06T06:06:33Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-readme-main-flip-reconciliation
issue:
---

The main-flip milestone needs README and install-facing docs that describe the
post-flip world accurately: stable users install from `main`, while `next`
remains available as a dev-only channel. Current `origin/next` docs still carry
pre-flip assumptions such as "releases are cut from next", "origin/main is
vestigial", and `spacedock install` reinstalling from
`spacedock-dev/spacedock@next`. Those statements are true for the 0.19.x
pre-main line but become wrong once `next` is force-pushed to `main` and
`v0.20.0` becomes the first stable release.

The worker must reconcile the docs against the branch plan recorded in
`main-flip-0200-marketplace`:

- tag current `origin/main` as `v0-archived`;
- replace `origin/main` with the prepared `next` line;
- keep `next` as the dev-only release/publish branch;
- move stable release mechanics and marketplace install references to `main`;
- keep a clear dev-only `next` path for source builds and pre-stable publishing.

## Existing README PRs to consider

Review these open PRs for reusable reader-facing direction, examples, and copy
patterns. They are inputs, not sources of truth: both target the old `main`
history and must be reconciled against current `origin/next` and the 0.20.0 flip
mechanics before reusing any material.

- PR #213, `docs(readme): lead with the problem Spacedock solves`
  (`spacedock-readme-problem-rewrite`): problem-led README opening that names
  the pain, mechanism, and "you want Spacedock if" reader axes.
- PR #220, `docs: refactor README for newcomers (developer and non-developer)`
  (`docs/readme-refactor-newcomer-friendly`): larger newcomer-oriented docs
  suite with README, getting-started, usage, examples, and prompts. Consider its
  install clarity, glossary, examples, and non-developer framing, but do not
  blindly import stale command behavior or old-branch assumptions.

## Scope

Likely files:

- `README.md`
- `docs/install-journey.md`
- `docs/releasing.md`
- any nearby docs that still describe `next` as the stable marketplace/release
  source after the flip

This task is documentation-only unless a small test/documentation fixture is
needed to keep command examples honest. Do not implement the branch flip, release
mechanics change, or upgrade-path behavior here; this task prepares the docs for
that work.

## Acceptance criteria

**AC-1 - Stable install docs describe the post-flip `main` lane.**
Verified by: reviewing the changed docs and, where practical, command snippets
that name observable commands. The stable path must direct users to the released
binary and plugin from the stable `main` lane, not to `next`.

**AC-2 - Dev docs keep `next` as a dev-only channel.**
Verified by: docs distinguish stable `main` from dev-only `next`, including
source-build or dev publish guidance where relevant. `next` must not be deleted
from the story.

**AC-3 - Release docs match the intended branch mechanics.**
Verified by: `docs/releasing.md` or an equivalent release note no longer says
`origin/main` is vestigial for the stable lane; it records archiving pre-v1 main,
guarded replacement of `main` from `next`, and `0.20.0` tagging from `main` while
retaining a dev-only `next` path.

**AC-4 - Existing README PRs were considered.**
Verified by: the implementation report names what was reused, adapted, or
rejected from PR #213 and PR #220, with a short reason for each. The product
docs do not need to mention the PRs unless that helps readers.

**AC-5 - Upgrade-path docs do not promise behavior this task does not ship.**
Verified by: stale-plugin, missing-binary, and outdated-binary guidance is
accurate to current or explicitly planned behavior. If the current behavior is
not yet implemented, docs should point to the pending confirmation work rather
than claim it already works.

## Test plan

- Inspect current `origin/next` docs and the open README PRs #213 and #220.
- Update the docs so stable `main`, dev-only `next`, and the 0.20.0 flip plan are
  not contradictory.
- Run focused doc checks where available, then at least `go test ./...`.
- If tests are blocked by unrelated local state, report the exact blocker and
  run the narrowest relevant checks that still prove the docs build or references
  resolve.
