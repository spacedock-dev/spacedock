---
id: rdjjq9hbv86skkw12z106z6q
title: Make merge-guard archive finalization durable across split-root hosts
status: ideation
source: "Roborev 2146 during durable-decisions 6y final implementation review, 2026-07-24"
started: 2026-07-24T15:38:16Z
completed:
verdict:
score: "1.0"
worktree:
issue:
sprint: durable-decisions
gates:
    version: 1
    current:
        gate: gate:docs-dev:rd:backlog
    records:
        - id: gate:docs-dev:rd:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:rd-backlog-1
              briefing:
                id: briefing:docs-dev:rd:backlog:attempt-1:revision-1
                digest: sha256:a6109cdc4013c201d8a268c32dc6795d354718cfdeae85ba0fd84546ed720659
                digest-domain: canonical-bytes
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:rd:backlog:1
                briefing: briefing:docs-dev:rd:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-07-24T15:38:08.065814Z"
                decision: approve
                reason: The split-root remote can retain an active merge sentinel after local archive finalization; a restart-safe supported publication path is a material durable-decisions prerelease requirement.
                adoption-note: you have the conn toward the sprint goal. authorized to approve gates, PR, relevant CI lanes, and merge.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

Make a successful split-root `merge guard` finalization durable on the configured state remote, including restart after interruption, so another First Officer cannot resurrect an already archived task from the last pushed sentinel.

## Problem

`merge guard` terminalizes and commits the path-scoped archive move in the state checkout, but it does not publish that commit. The preceding `pr=pr-merge:<N>` state commit is pushed, so the remote can remain at an active terminal sentinel while the local checkout alone contains the archive. A fresh host then sees stale active state; an interruption after the archive commit but before a manual push leaves no active task for the local host to resume through ordinary merge-guard selection.

This is a supported split-root durability defect exposed by 6y's terminal-gate resumption, not a reason to expand 6y's no-CLI/nine-file implementation. It is a durable-decisions prerelease blocker.

## Scope boundaries

- Own the supported state-sync mechanism needed after merge-guard archive commit and its restart path.
- Reuse the existing `state commit` non-fast-forward/rebase-conflict discipline; never force-push or invent a second sync policy.
- Preserve path-scoped archive commits and sibling dirt isolation.
- Do not change gate authority, recorder/application schemas, provider presentation, PR-host behavior, or 6y's lifecycle semantics.
- Do not paper over the gap with FO-authored raw `git push` prose; expose one supported binary-owned outcome.

## Acceptance criteria

**AC-1 (VALUE) - Successful split-root finalization is visible as the same archived terminal task from a fresh host.**
Verified by: a real two-clone fixture pushes the merge sentinel, runs finalization on host A, then boots host B and observes only the archived terminal task at the pushed archive commit, with no active sentinel row. Removing the publish step makes host B reproduce the stale active task.

**AC-2 - Interruption after the local archive commit but before publication has a supported idempotent resume.**
Verified by: a failure-injection fixture stops at that boundary, reruns the documented binary path without hand Git commands, publishes exactly one archive commit, and leaves the state checkout clean. A path that can resolve only active tasks fails this test.

**AC-3 - Peer synchronization retains the existing conflict safety.**
Verified by: a disjoint peer state commit is integrated and re-pushed with both changes present; a same-task conflict exits through the existing halt class with rebase aborted, no force push, and both sides recoverable.

**AC-4 - Archive publication remains path-scoped and does not sweep sibling dirt.**
Verified by: the final pushed commit contains only the task's live/archive paths while an unrelated dirty task remains uncommitted locally. Replacing path-scoped staging with a broad add makes the fixture fail.

**AC-5 - Inline/local-only workflows retain truthful existing behavior.**
Verified by: focused fixtures show no remote operation is required when no split-root remote is configured, while output distinguishes durable remote publication from local-only finalization. Silent success claiming remote durability without a push fails.

## Test expectations

Ideation must spike the exact crash seam first using the existing real-git merge-guard and state-sync fixtures: archive locally, leave origin at the sentinel, and determine the smallest existing binary entry point that can resume from the archived slug. Compare reuse/refactoring of the current state-sync implementation before proposing any new command. Implementation must add focused red/green tests before code, then run `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, strict docs, and the runtime/live lanes required by the actual diff. Validation includes a detached adversarial audit because this changes the merge/state mutation boundary.
