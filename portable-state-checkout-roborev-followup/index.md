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

## Proposed approach

- Resolve split-root Git strictly from the declared entity/state directory before any mutation; never accept a parent code repository when state `.git` metadata is absent.
- Validate every regular-file `.git` target against the declared project, including operable absolute targets. Check administrative identity, common directory, and back-pointer with filesystem identity where mount aliases require it.
- Resolve and validate `origin` before creating a state commit, while preserving explicit local-only behavior and reliable retry/push handling.
- Add the three adversarial regressions requested by Roborev, then rerun the exact branch review after the fixes.

## Out of scope

- Rewriting or deleting the archived `qw` record.
- Hiding the prior rejected validation or Roborev review.
- Removing the Git 2.48 workflow file to bypass the OAuth `workflow`-scope requirement.
- Broad refactoring outside state Git resolution, finalization, and commit ordering.

## Acceptance criteria

**AC-1 (VALUE) - Split-root operations never mutate a parent or source repository when the declared state checkout is missing, copied, moved, or backed by an operable stale absolute gitfile.**
Verified by: a public-command matrix covering missing state `.git`, copy-while-source-remains, moved checkout, and mount-alias identity; each negative fails before mutation and asserts both repositories' refs, indexes, worktrees, and entity bytes remain unchanged.

**AC-2 - Merge finalization with a declared split-root state directory lacking `.git` fails closed instead of staging or committing the archived entity on the code branch.**
Verified by: a merge-guard/finalization regression that records code HEAD, index, status, archive paths, and state before and after the refused operation.

**AC-3 - Every regular-file `.git` target is validated against the declared project even when raw Git commands would succeed.**
Verified by: a copy-while-source-remains test whose copied checkout points at the still-operable source metadata; the command refuses and the source worktree metadata and refs remain byte-for-byte unchanged.

**AC-4 - `state commit` validates origin before local mutation and cannot strand a commit behind the no-op retry path.**
Verified by: a repository with a named `origin` lacking a usable URL; the first call fails with unchanged HEAD and entity dirt preserved, then a corrected origin allows one commit and push. Preserve the documented no-origin local-only case.

**AC-5 - The corrected exact range passes independent Roborev review and required repository gates.**
Verified by: a replacement thorough Roborev job on the follow-up branch records PASS at its exact head, followed by `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.

## Test plan

Start with the three Roborev findings as failing public-command regressions. Exercise missing state metadata, copy-with-live-source, named-origin-without-URL, no-origin local-only, moved relative checkout, and mount-alias cases. Assert exact identity, refs, index, worktree registrations, entity bytes, exit codes, and cleanup state—not only error substrings. Run focused packages, full and race suites, then a replacement thorough Roborev branch review under the repository guideline.
