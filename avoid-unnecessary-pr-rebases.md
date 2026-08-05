---
id: ep2cz3zsb2qpyyh889nyeqpr
title: Avoid unnecessary rebases before opening mergeable PRs
status: ideation
source: GitHub issue #616; Captain intake 2026-08-05
started: 2026-08-04T17:00:42Z
completed:
verdict:
score: 1.0
worktree:
issue: "#616"
sprint: durable-decisions
gates:
    version: 1
    records:
        - id: gate:ep2cz3zsb2qpyyh889nyeqpr:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:ep2cz3zsb2qpyyh889nyeqpr-backlog-1
              briefing:
                id: briefing:ep2cz3zsb2qpyyh889nyeqpr:backlog:attempt-1:revision-1
                digest: sha256:33cb5c91df3346daadadbb7a6c3e3d7740e4af2ab9ca5a0b0cac93608cc932fd
                request-digest: sha256:fb514b4f3cb9a4ff9b05a114822c28c36b369b9f7dc9d449d055f0f63e020bd2
                room-ref: ./avoid-unnecessary-pr-rebases/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ep2cz3zsb2qpyyh889nyeqpr:backlog:1
                briefing: briefing:ep2cz3zsb2qpyyh889nyeqpr:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-04T17:00:18.701809Z"
                decision: approve
                reason: 'Captain directed intake and dispatch and made issue #616 desired behavior authoritative for this sprint session.'
              application:
                target-stage: ideation
                state: consumed
---

Remove the mandatory rebase from the `pr-merge` policy. Keep the approved
candidate commit when Git can merge it with the current integration trunk.

## Problem

The shipped mod rebases every candidate before it creates a PR. The development
workflow requires the same rebase before it invokes the mod.

Issue #616 records the smallest failure. The candidate is one commit behind and
one commit ahead of the trunk. Git can merge both commits, but the policy rewrites
the candidate SHA and invalidates its exact-head approval.

## Captain design reset

The Captain rejected the cycle-2 production-seam proposal. Issue #616 is a defect
in the `pr-merge` policy, not a missing dispatch API.

The rejected proposal added a command, JSON output, shared reconcile changes,
host machinery, and new evidence keys. Those changes did not need to correct the
mandatory rebase. They also expanded EP2 into G3/D8 ownership and integration
evidence policy.

Cycle 3 limits EP2 to the two mod files. Existing validation and PR review continue
to own integration proof. G3/D8 continues to own conflict reconciliation.
EP2 no longer promises that all ancestry drift is report-only. Thus, the prior
shared-core and reconcile-test finding is outside this reduced policy correction.

## Proposed policy

### 1. Record the approved candidate

Before the Captain reviews the PR draft, record this value:

```text
CANDIDATE_SHA=$(git -C {worktree} rev-parse HEAD)
```

Show the candidate SHA in the draft. The Captain approves this exact commit.

### 2. Exercise mergeability without a rebase

After approval, update or publish the integration trunk with the existing
variant-specific step. Resolve that current trunk tip as `BASE_SHA`.

Run this command from the code repository:

```text
git merge-tree --write-tree "$BASE_SHA" "$CANDIDATE_SHA" >/dev/null
```

Exit 0 means that Git found a clean merge. Continue PR delivery without a rebase.
Exit 1 means conflict. Stop PR and local-merge delivery, and refer the conflict to
the G3/D8 owner path. Preserve the pending delivery authority.

Any other exit means that mergeability is unknown. Stop delivery and report the
error. Do not rebase or use local merge as a fallback for this result.

The command does not change the candidate ref, index, or worktree. It is only the
delivery decision. It is not new integration evidence.

### 3. Push the approved commit

For a clean result, push the recorded commit with an explicit refspec:

```text
git push origin "$CANDIDATE_SHA:refs/heads/{branch}"
```

This one-line rule keeps the Captain-approved commit if the local branch name
moves. The task does not add remote-head or PR-head verification.

## Acceptance criteria

**AC-1 (VALUE) - A clean `1 1` candidate reaches PR creation with its approved SHA.**

Verified by a real Git graph with one base commit and one candidate commit after
their merge base. The mergeability command exits 0. The candidate SHA, refs,
index, and worktree do not change. The command log contains no rebase.

**AC-2 - Conflict and unknown results stop delivery without changing authority.**

Verified by conflict and command-error exercises. Neither result pushes a branch,
creates a PR, starts a local merge, or changes entity state. The conflict result
refers to G3/D8. The unknown result reports the command error.

## Expected surface

Expected changes are two files and approximately `+30/-10` lines, with a 1.5x
tolerance:

- `mods/pr-merge.md`: replace the mandatory rebase with the mergeability exercise
  and explicit candidate-SHA push.
- `docs/dev/_mods/pr-merge.md`: remove the pre-hook rebase requirement and apply
  the same decision to the split-root workflow.

No command grammar, JSON, stored state, shared core, reconcile behavior, host
adapter, evidence infrastructure, or Git-version policy changes.

## Test plan

Use real temporary Git repositories. Do not extract or execute commands from mod
prose.

1. Create a clean `1 1` graph. Run the selected Git command directly. Check exit
   0, an unchanged candidate SHA, unchanged refs, and a clean index and worktree.
2. Create a content conflict. Check exit 1 and no branch, PR, local-merge, or state
   action.
3. Run an invalid command environment. Check the unknown route and the same
   no-action result.
4. Move a local branch name after draft approval. Push the recorded SHA to a bare
   remote and check that the remote branch contains that SHA.
5. Run `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, and
   `git diff --check`.

## Real Git evidence

The cycle-3 fixture is `/tmp/spacedock-616-cycle3.HcC2gz`. Its ahead/behind result
is exactly `1 1`. The proposed mergeability command exited 0 and returned tree
`4b825dc642cb6eb9a060e54bf8d69288fbee4904`.

The candidate stayed at `af53ce22b2fb43475596977ac6ca2831053b3aea`.
The index and worktree stayed clean. No ref changed and no rebase ran.

## Out of scope

- A `dispatch mergeability` command or stable JSON result.
- Shared reconcile, shared-core, or host-adapter behavior.
- New integration-evidence keys or base-sensitive validation policy.
- Generic owner discovery, dispatch, or conflict resolution. G3/D8 owns that path.
- Remote-head and PR-head verification after the exact-SHA push.
- A new minimum Git version or compatibility fallback.

## Stage Report: ideation

- DONE: Define the smallest falsifiable mergeability decision that prevents clean-head rewrites and serves AC-1.
  A real `1 1` graph preserved the candidate SHA at exit 0. SK and Z5 show the same live clean-merge result.
- DONE: Specify recorded-owner reconciliation and authority preservation for genuinely conflicting branches under AC-2.
  The conflict route uses the recorded stage, branch, and worktree owner, preserves delivery state, and forbids rework, automatic resolution, and force operations.
- DONE: Bound exact-head evidence invalidation and regression coverage, including the clean ahead/behind case from issue #616, under AC-3.
  The design names the exact refresh set and defines real-Git, state-byte, and live Claude regression fixtures.

### Summary

The design replaces mandatory rebase with one read-only Git mergeability decision.
Clean candidates keep their validated SHA. Conflicting candidates return to their
recorded owner while the First Officer preserves state and pending authority.

## Stage Report: ideation (cycle 2)

- DONE: Correct AC-3 with separate candidate and integration evidence keys.
  The integration key contains the base SHA, candidate SHA, and merge-tree OID.
  A clean-merge/build-red spike proves that base-only drift can require new
  integration proof even when Git finds no text conflict.
- DONE: Replace the invalid quiet diagnostic with one production command.
  The command keeps Git 2.38 support and parses NUL-separated tree and conflict
  fields. Its tests include unusual filenames, directory/file conflicts, and a
  pathless conflict result.
- DONE: Keep behavioral proof outside instruction prose.
  Production command tests and shared live fixtures prove behavior. Contract lint
  checks only structural references and forbidden rebase prescriptions.
- DONE: Include shared reconcile and Claude stale-branch behavior in the surface.
  Both layers report ancestry drift without prescribing or running a rebase.
- DONE: Close the probe-to-push race with exact-head delivery.
  The merge hook pushes the recorded candidate SHA. It checks the remote ref and
  PR head before it records `pr:`. An adversarial fixture moves the local branch.
- DONE: Preserve generic owner-handoff authority.
  EP2 delegates conflict packages to G3/D8. It does not duplicate owner discovery,
  dispatch, rework, or authority mutation.
- DONE: Record both Captain decisions in the contract.
  The design retains Git 2.38 compatibility and makes G3/D8 a dependency for the
  generic owner route.

### Summary

The revised design keeps clean candidate heads without treating textual mergeability
as complete integration proof. It also preserves exact-head authority from the
decision through PR creation and delegates real conflicts to the existing owner path.

## Stage Report: ideation (cycle 3)

- DONE: Reduce the correction to the mod-owned issue #616 policy.
  The expected surface contains only the shipped and development `pr-merge` mods.
- DONE: Remove the unconditional ancestry rebase and keep the approved candidate.
  A direct `merge-tree --write-tree` exercise decides clean, conflict, or unknown.
- DONE: Keep genuine conflict reconciliation in the existing owner path.
  Conflict stops delivery and refers to G3/D8 without new dispatch or state logic.
- DONE: Preserve one falsifiable clean `1 1` exercise.
  Fixture `/tmp/spacedock-616-cycle3.HcC2gz` exited 0 with an unchanged candidate,
  unchanged refs, and a clean index and worktree.
- DONE: Record the Captain design reset and remove the larger mechanism.
  The task adds no command, JSON, shared reconcile, host, or evidence infrastructure.

### Summary

Cycle 3 corrects issue #616 in the two policy files that cause the rewrite. The
mod keeps the approved SHA on a clean result and fails closed on conflict or error.
