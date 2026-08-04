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

Prevent the `pr-merge` workflow from rewriting a validated candidate when the
candidate still merges cleanly with the integration trunk.

## Problem

The shipped mod rebases every branch before PR creation. The development mod
requires the First Officer to do the same work before the merge hook runs. The
Claude reconcile adapter also treats every owned `stale-branch` row as a rebase
instruction. These rules use ancestry drift as proof of a conflict.

A rebase changes the candidate commit even when Git can merge both tips. This
change makes exact-head checks stale without correcting a conflict. Issue #616
records the smallest case: the candidate is one commit behind and one commit
ahead of `origin/main`, and Git can merge it without a conflict.

## Observed delivery cases

The following observations use `origin/main` from 2026-08-05. They are examples,
not dependencies for this task.

| Case | Candidate | Behind/ahead | `git merge-tree` | Delivery fact |
|---|---|---:|---|---|
| SK | `cd76f52ab` | `5 1` | exit 0 | PR #617 is mergeable at the same head. An earlier clean rebase withdrew validation attempt 1. |
| 26n | `2cebb23b8` | `66 11` | exit 1 | The original branch conflicts. PR #583 now uses rewritten head `aeb3009b0`. |
| KD | `58888fddb` | `66 10` | exit 1 | The held candidate conflicts and has no PR. The owning branch must reconcile it. |
| Z5 | `37abaa8b3` | `5 5` | exit 0 | The held candidate is cleanly mergeable. Base drift alone does not require a rewrite. |

SK and Z5 prove that ahead/behind counts do not decide mergeability. The 26n and
KD results prove that the clean path cannot replace owner reconciliation.

## Spike result

A real Git fixture at `/tmp/spacedock-616-mergeability.tntUm9` created two
independent graphs. The clean graph was exactly `1 1`. This command exited 0 and
left the candidate SHA unchanged:

```sh
git -C "$WORKTREE" merge-tree --write-tree --quiet "origin/$BASE" "$CANDIDATE"
```

The conflict graph changed `shared.txt` on both sides. The probe exited 1 and left
the candidate SHA unchanged. A diagnostic run with `--name-only` named
`shared.txt`. These operations write temporary Git objects. They do not change
refs, the index, or the worktree. No new helper needs a spike.

## Proposed approach

### 1. Make mergeability the delivery decision

At the merge boundary, fetch the configured trunk. Record its tip as `BASE_SHA`.
Record the worktree tip as `CANDIDATE_SHA`. Then run this read-only probe:

```sh
git -C "$WORKTREE" merge-tree --write-tree --quiet "$BASE_SHA" "$CANDIDATE_SHA"
```

The exit status has one meaning at this boundary:

- Exit 0 means that Git produced a clean merge tree. Keep `CANDIDATE_SHA` and
  continue PR preparation if the candidate still has commits outside the base.
- Exit 1 means that Git found a conflict. Run the same command with `--name-only`
  to collect conflict paths. Then enter the owner route in section 2.
- Any other exit means that the decision is unknown. Stop PR creation and report
  the command error. Do not rebase as a fallback.

If the candidate is already an ancestor of the base, report that it is already
integrated. Do not create an empty PR.

After captain approval, fetch the trunk again. If `BASE_SHA` changed, run the
probe again. If `CANDIDATE_SHA` changed, use the evidence rule in section 3.
Push and open the PR only after the final probe exits 0.

This mechanism serves AC-1. The simpler `behind == 0` rule rejects SK and Z5.
GitHub mergeability is insufficient because no PR exists before this decision.
A new Spacedock command is unnecessary because Git already gives one exit-status
decision without worktree mutation.

### 2. Route conflicts to the recorded owner

An ancestry-only `stale-branch` row becomes report-only. The reconcile sweep must
not rebase a branch. The merge hook owns the delivery-time probe.

If the probe exits 1, preserve the entity stage, `pr`, merge `mod-block`, and
pending terminal application. The detection path must not consume, clear, or
supersede authority. It must not run `merge guard --rework`.

Identify the owner from the workflow dispatch identity, current entity stage,
branch, and `worktree` field. Do not use the GitHub author or Git credential as
the owner. Send the conflict package to the live owner. If that handle is absent,
dispatch a fresh worker for the same recorded stage and worktree.

The package contains the entity, PR reference, branch, worktree, old and new base
SHAs, candidate SHA, and conflict paths. The owner changes only its branch after
captain authorization. The First Officer guards state and transport. It does not
resolve conflicts, select `ours` or `theirs`, or force-push.

This mechanism serves AC-2. A generic resolver is insufficient because it has no
semantic ownership. `merge guard --rework` is insufficient because it burns the
pending terminal application and clears truthful delivery state.

### 3. Refresh evidence by candidate SHA

The candidate SHA is the invalidation key. A base-only change makes only the
mergeability result stale. Run the probe again, but keep validation, CI, gate,
and PR-draft evidence when `CANDIDATE_SHA` is unchanged.

If owner reconciliation changes `CANDIDATE_SHA`, refresh these items:

1. Deterministic validation and required live checks that ran on the old SHA.
2. CI results whose recorded head is the old SHA.
3. Validation report claims, Briefing content, and terminal approval that cite or
   authorize the old SHA.
4. PR draft evidence that was derived from those stale results.
5. The mergeability probe against the current trunk tip.

Keep the ideation baseline, accepted scope, entity identity, issue, branch owner,
and evidence that is not head-specific. Conflict detection makes no authority
write. If the owner returns a new head, replace old-head gate authority through
the normal recorded-gate path. Never consume it as delivery proof.

This mechanism serves AC-3. Refreshing all evidence on every trunk change repeats
the defect. Reusing old exact-head evidence after a candidate change weakens the
existing validation contract.

## Out of scope

- Automatic conflict resolution, force-push, or a generic resolver worker.
- A new command, stored field, state schema, PR provider, or merge strategy.
- A requirement that every PR branch is current with the trunk.
- Changes to state-checkout rebase safety. A same-entity state conflict remains a
  global state-safety halt.
- Changes to GitHub branch protection or CI merge-queue policy.

## Acceptance criteria

**AC-1 (VALUE) - A clean `1 1` candidate reaches PR preparation with zero candidate-head rewrites.**

Verified by: a real Git fixture runs the exact mod probe on a clean branch that is
one commit behind and one ahead. The fixture records zero rebase commands, exit 0,
the same candidate SHA, and an eligible PR-preparation result.

**AC-2 - A conflicting candidate reaches its recorded owner with delivery authority unchanged.**

Verified by: a conflict fixture records exit 1 and exact conflict paths. It routes
the package to the recorded stage, branch, and worktree owner. Entity status, `pr`,
`mod-block`, and the pending terminal application remain byte-identical. The
command log contains no consume, clear, supersede, auto-resolution, or force action.

**AC-3 - Evidence refresh follows candidate-head changes, not base ancestry changes.**

Verified by: a table-driven fixture first changes only the base and requires only
a new mergeability probe. It then changes the candidate and requires the five-item
refresh set from section 3. A mutation that refreshes exact-head checks for the
base-only case, or reuses old-head evidence, must fail.

## Expected surface and semantic boundaries

Expected changes are 6-8 files and approximately +300/-20 lines:

- `mods/pr-merge.md`: replace the mandatory rebase with the probe and routes.
- `docs/dev/_mods/pr-merge.md`: apply the same rule to the split-root variant.
- `skills/first-officer/references/claude-fo-dispatch.md`: make ancestry-only
  `stale-branch` rows report-only.
- `skills/integration/pr_merge_policy_test.go`: add real Git and state fixtures.
- `internal/ensigncycle/` shared fixture and Claude runner files: prove that the
  live adapter reports ancestry drift without rebasing the clean candidate.
- `docs/site/advanced/mods-and-standing-teammates.md`: document the user-visible
  no-rewrite policy.

Tolerance is 2x for files and insertions. Command grammar and stored formats do
not change. Authority remains with the recorded gate and owning branch. Runtime
behavior changes only at stale-branch handling and PR delivery.

The proposed documentation change is:

```diff
 The canonical example is the `pr-merge` mod: it opens the code-branch PR at merge,
 records the PR on the entity, and holds the terminal transition until the PR merges.
+It keeps a validated branch head when Git can merge it cleanly. A real conflict
+returns to the recorded branch owner without consuming pending merge authority.
```

## Test plan

Add the focused smoke test before the instruction changes. The test extracts the
declared probe from both mod files and runs it in temporary real Git repositories.

The clean fixture creates exactly one trunk commit and one candidate commit after
their merge base. It checks exit 0, a stable candidate SHA, clean index/worktree,
and no rebase command. The conflict fixture changes one path on both sides. It
checks exit 1, the named path, unchanged refs and state bytes, and owner routing.

Add a table for base-only and candidate-head changes. The table checks the exact
refresh set and rejects both over-refresh and stale-head reuse. Add a contract test
that both mod variants use the same decision and that the Claude adapter does not
rebase an ancestry-only `stale-branch` row.

Run the focused integration package, `go test ./...`, `go test ./... -race`,
`gofmt -w ./cmd ./internal`, `git diff --check`, and the existing contract lints.
Run one isolated Claude fixture with a clean `1 1` branch. The fixture must show
that the reconcile pass does not change its head. It must then reach the clean
merge-boundary decision. Other host lanes are unchanged and do not need a run.

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
