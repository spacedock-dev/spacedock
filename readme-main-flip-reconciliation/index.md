---
id: nbz0yjvmqm6gda6csw8yef7k
title: README reconciliation for main flip and 0.20.0 install paths
status: validation
source: "captain (2026-06-06) - before flipping main, reconcile README/install docs; consider existing README PRs #213 and #220."
score: "0.39"
started: 2026-06-06T06:06:33Z
completed: 2026-06-06T06:25:18Z
verdict: PASSED
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

## Stage Report: implementation

- DONE: Stable install/release docs describe the post-flip `main` lane and no longer present `next` as the stable marketplace source.
  Evidence: code commit `9fae02f7` updates `README.md`, `docs/install-journey.md`, and `docs/releasing.md`; focused stale-language `rg` returned no matches in those docs.
- DONE: Dev docs preserve `next` as a dev-only channel and release docs record the archive-current-main, guarded replacement, and 0.20.0 main-tag mechanics.
  Evidence: `docs/releasing.md` now records `v0-archived`, `--force-with-lease=main:<sha>`, `v0.20.0` tagging from `main`, and a dev-only `next` publishing section.
- DONE: Implementation report explicitly says what was reused, adapted, or rejected from README PRs #213 and #220, and docs do not overclaim unshipped upgrade-path behavior.
  Evidence: reused/adapted PR #213's problem-led opener; adapted PR #220's newcomer lane/example framing; rejected both PRs' stale old-main install commands and marketplace behavior; docs state upgrade-path confirmation is still owed.
- DONE: Verification.
  Evidence: `gofmt -w ./cmd ./internal`; `go test ./...` passed 1138 tests in 16 packages; `go test ./... -race` passed 1138 tests in 16 packages; `git diff --check` passed.

### Summary

Docs now describe stable installs and releases as post-flip `main` lane behavior while keeping `next` for development-only source builds and pre-stable publishing. No release-mechanics code, branch flip, or upgrade-path behavior was implemented; the docs call out that those release gates remain owed before `0.20.0`.

## Stage Report: validation

- DONE: AC-1 stable install docs describe the post-flip `main` lane.
  Evidence: commit `9fae02f7` changes only `README.md`, `docs/install-journey.md`, and `docs/releasing.md`. `README.md:24-35` says tagged releases, Homebrew artifacts, and marketplace plugin installs come from `main`, with `spacedock install --host claude` resolving `spacedock-dev/spacedock` on `main`, not `next`. `docs/install-journey.md:7-8`, `28-67` name the stable `main` lane and the brew/plugin install path from `main`. Focused stale-language scan over the three changed docs found no residual `spacedock-dev/spacedock@next`, `clkao/spacedock`, `vestigial`, or `NEVER push` stable-install wording.
- DONE: AC-2 dev docs keep `next` as a dev-only channel.
  Evidence: `README.md:37-45` keeps the source-build lane on `next` with `--plugin-dir` and explicitly says `@next` is not the stable install path. `docs/install-journey.md:9-10`, `98-115`, and `117-165` keep the dev-only source routes, including local `next` checkout, `go install ...@next`, dev snapshot, and `--plugin-dir`. `docs/releasing.md:133-141` keeps `next` for source builds and deliberate `next-publish`, then says commands or manifests using `@next` are dev-only.
- DONE: AC-3 release docs match the intended branch mechanics.
  Evidence: `docs/releasing.md:3-6` now says stable releases start from `main` at `v0.20.0` and `next` is dev-only. `docs/releasing.md:14-20` records archiving current pre-v1 `origin/main` as `v0-archived`; `32-40` records guarded replacement with `--force-with-lease=main:"$preflip_main"` while preserving `next`; `42-50` records cutting `v0.20.0` from the `main` release line. `docs/releasing.md:52-55` explicitly leaves `.github/workflows/release.yml` and `.goreleaser.yaml` retargeting as pending main-flip release work, so the unchanged `.goreleaser.yaml` `@next` comment is not presented as shipped stable behavior. `docs/releasing.md:74-123` changes later stable release mechanics to `origin/main` and `release/X.Y.Z:main`.
- DONE: AC-4 existing README PRs were considered without importing stale old-main behavior.
  Evidence: `gh pr view 213 --json number,title,headRefName,baseRefName,body,url,state` returned open PR #213, branch `spacedock-readme-problem-rewrite`, title `docs(readme): lead with the problem Spacedock solves`, targeting old `main`; `gh pr diff 213 --patch --color never` showed reusable pain/mechanism/state-outside-agent README opening material. The resulting `README.md:3-5` adapts that mechanism into the conservative current README. `gh pr view 220` returned open PR #220, branch `docs/readme-refactor-newcomer-friendly`, title `docs: refactor README for newcomers (developer and non-developer)`, also targeting old `main`; `gh pr diff 220` showed newcomer scenario/example framing but also old install commands such as `claude plugin marketplace add clkao/spacedock && claude plugin install spacedock`. The product docs reject that stale install behavior and instead use brew plus `spacedock install --host claude` on `main` (`README.md:27-35`, `docs/install-journey.md:34-67`). The implementation report records the same reuse/adapt/reject split at lines 108-109 of this entity.
- DONE: AC-5 upgrade-path docs do not promise unshipped behavior.
  Evidence: `README.md:51-61` says stale-plugin recovery is `spacedock doctor` plus `spacedock install --host claude`, then explicitly says old-plugin/no-binary and binary/plugin-skew journeys still need release-gate confirmation before the `0.20.0` flip. `docs/install-journey.md:80-96` makes `doctor` the compatibility source of truth, says missing binary means install Homebrew first, and does not promise automatic recovery for old `0.12.x` no-binary or `0.19.x` skew cases. No release-mechanics code, branch flip, or upgrade-path behavior was added in commit `9fae02f7`.
- DONE: Verification evidence reproduced.
  Evidence: `gofmt -w ./cmd ./internal` exited 0; `git diff --check` exited 0; `go test ./...` passed 1138 tests in 16 packages; `go test ./... -race` passed 1138 tests in 16 packages. Product worktree status before this state-report edit was clean except for branch `spacedock-ensign/readme-main-flip-reconciliation` being one commit ahead with implementation commit `9fae02f7`.

### Summary

Recommendation: PASSED. All five ACs are satisfied by the product diff and reproduced checks. The docs now send stable users to `main`, preserve `next` as dev-only, document the `v0-archived` / guarded-main-replacement / `v0.20.0` main-tag mechanics, account for PRs #213 and #220 without importing their old-main install behavior, and avoid promising unshipped upgrade recovery. Detached adversarial audit was not triggered: this is a docs-only validation and the release machinery itself remains explicitly pending.
