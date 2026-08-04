---
id: ep2cz3zsb2qpyyh889nyeqpr
title: Avoid unnecessary rebases before opening mergeable PRs
status: backlog
source: GitHub issue #616; Captain intake 2026-08-05
started:
completed:
verdict:
score: 1.0
worktree:
issue: "#616"
sprint: durable-decisions
---

Prevent the `pr-merge` workflow from rewriting an already validated candidate merely because `origin/main` advanced when the branch remains cleanly mergeable.

## Problem

The current pre-PR contract requires rebasing onto the latest integration trunk. A cleanly mergeable branch can therefore receive a new head, lose exact-head validation and gate evidence, repeat expensive checks, and race the next upstream landing without resolving any conflict.

## Proposed approach

Ideation must define a falsifiable mergeability decision before rebase, the narrow conditions that require owning-branch reconciliation, and the exact evidence invalidated when a candidate head truly changes.

## Out of scope

This task does not bypass real conflicts, weaken exact-head evidence, authorize force pushes, or move semantic reconciliation from the owning branch to the First Officer or Captain.

## Acceptance criteria

**AC-1 - A cleanly mergeable branch can reach PR creation without changing its validated head solely because the integration trunk advanced.**
Verified by: a regression fixture starts from a clean ahead/behind branch, exercises the merge-policy front door, and asserts that the candidate SHA is unchanged while PR preparation remains eligible.

**AC-2 - A genuinely conflicting branch is routed to its recorded owner without burning pending gate or merge authority.**
Verified by: a conflict fixture observes the owner-reconciliation route and unchanged PR, mod-block, and pending approval state.

**AC-3 - When reconciliation changes the candidate head, the contract names and requires only the evidence made stale by that change.**
Verified by: command or fixture output distinguishes unchanged-head delivery from changed-head reconciliation and exposes the required refresh set.

## Test plan

Ideation will identify the smallest existing mergeability primitive, add focused Go/fixture coverage for clean and conflicting branches, retain stable command output fixtures, and require `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` before completion. Live-host testing is not required unless the proposed mechanism changes runtime-host behavior.
