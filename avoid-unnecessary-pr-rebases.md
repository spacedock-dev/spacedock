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

Prevent the `pr-merge` workflow from rewriting a validated candidate when Git can
merge that candidate cleanly with the integration trunk.

## Problem

The shipped mod rebases every branch before PR creation. The development mod
requires the First Officer to do the same work before the hook runs. The shared
reconcile result also tells an owner to `pull --rebase` for ancestry drift.

These rules use behind/ahead counts as proof of a conflict. A clean rebase changes
the candidate SHA and makes candidate-specific evidence stale. It can also race
another trunk change before a PR exists.

A clean merge tree is necessary, but it is not sufficient integration proof. A
candidate can remove an API while the new base adds a caller in another file. Git
merges those paths cleanly, but the merged program does not build.

## Observed delivery cases

The following observations use `origin/main` from 2026-08-05. They are examples,
not dependencies for this task.

| Case | Candidate | Behind/ahead | Merge result | Delivery fact |
|---|---|---:|---|---|
| SK | `cd76f52ab` | `5 1` | clean | PR #617 is mergeable at the same head. An earlier rebase withdrew validation attempt 1. |
| 26n | `2cebb23b8` | `66 11` | conflict | PR #583 now uses reconciled head `aeb3009b0`. |
| KD | `58888fddb` | `66 10` | conflict | The held candidate requires owning-branch reconciliation. |
| Z5 | `37abaa8b3` | `5 5` | clean | Base drift alone does not require a candidate rewrite. |

SK and Z5 show that ancestry counts do not decide mergeability. The 26n and KD
results show that the clean path cannot replace owner reconciliation.

## Spike results

Four real Git spikes established the revised mechanism:

1. A clean `1 1` graph produced a merge tree and left the candidate SHA unchanged.
2. A textual-clean graph produced tree `0e813e0329`. `go test ./...` on that tree
   failed because the base added a caller for an API that the candidate removed.
3. `-z --name-only --no-messages` preserved a filename with a tab and newline.
   A directory/file conflict returned exit 1 and a synthetic conflict path.
4. A branch moved after its SHA was recorded. An explicit SHA refspec still pushed
   the recorded SHA, not the new local branch tip.

Git 2.38 supplies `--write-tree`, NUL output, name-only conflict data, and the
0/1/error exit contract. This task does not use `--quiet` and does not require Git
2.50. An older or incompatible Git version fails closed with an upgrade diagnostic.

## Captain decisions

- **Git floor:** retain Git 2.38 compatibility. Remove `--quiet`. Do not add a Git
  2.50 floor or a rebase fallback.
- **Ownership:** EP2 owns mergeability detection, evidence keys, and pinned-head
  delivery. EP2 depends on G3 and D8 for generic recorded-owner handoff. It must
  not duplicate or replace that mechanism.

## Proposed approach

### 1. Add one production mergeability seam

Add this read-only command:

```text
spacedock dispatch mergeability --workflow-dir DIR ENTITY --json
```

The binary resolves the entity worktree, configured trunk, `origin/{trunk}` tip,
and candidate `HEAD`. It records both SHAs before it runs Git. The merge hook
fetches the trunk before each command call.

The command runs Git 2.38-compatible plumbing without `--quiet`:

```text
git merge-tree --write-tree -z --name-only --no-messages BASE_SHA CANDIDATE_SHA
```

The binary parses bytes, not lines. The first NUL field is the merge-tree OID.
The remaining fields are exact conflict paths. An exit status of 1 means conflict,
even when the path list is empty. Other statuses are command errors.

The stable JSON result contains:

```json
{
  "command": "mergeability",
  "base_sha": "...",
  "candidate_sha": "...",
  "merge_tree_oid": "...",
  "mergeable": true,
  "conflict_paths": [],
  "pathless_conflict": false
}
```

This mechanism serves AC-1 and AC-2. A prose command recipe is insufficient
because it has no production parser or stable output. A GitHub query is
insufficient because no PR exists at the first decision.

### 2. Make ancestry drift report-only

The shared reconcile result continues to report `stale-branch`, `behind`, `trunk`,
`worktree`, and `owned`. Its `reason` must not prescribe rebase. The Claude action
must report this row and defer the decision to the merge boundary.

The merge hook runs `dispatch mergeability` before it presents the PR draft. It
runs the command again after captain approval and a fresh trunk fetch. A clean
result permits PR preparation without a candidate rewrite. If the candidate is
already an ancestor of the base, the hook reports it as integrated and does not
create an empty PR.

This mechanism serves AC-1. Changing only the Claude instruction is insufficient
because the shared JSON would still prescribe `pull --rebase (owned)`.

### 3. Delegate conflicts without changing authority

If `mergeable` is false, EP2 passes the structured result to the generic G3/D8
owner-handoff path. The package includes the entity, PR, branch, worktree,
`base_sha`, `candidate_sha`, `merge_tree_oid`, exact conflict paths, and the
`pathless_conflict` flag.

The EP2 detection path preserves status, `pr`, merge `mod-block`, and the pending
terminal application. It does not consume, clear, supersede, or rework authority.
It does not identify an owner from Git credentials. G3/D8 owns live-owner reuse,
fresh-owner dispatch, and the no-resolver boundary.

This mechanism serves AC-2. EP2 depends on G3 and D8 before implementation can
claim the owner route. Duplicating their dispatch mechanism would create two
generic handoff contracts.

### 4. Key evidence to what can make it stale

Use two evidence keys:

- **Candidate key:** `candidate_sha`. Unit, behavior, and live checks that inspect
  only candidate bytes use this key.
- **Integration key:** `(base_sha, candidate_sha, merge_tree_oid)`. Build,
  integration, and merge-result checks use this key.

If only `base_sha` changes, rerun mergeability and integration-keyed proof. Keep
candidate-keyed proof. The pending terminal application cannot finalize delivery
until the new integration proof is green.

If `candidate_sha` changes, rerun candidate-keyed proof and integration-keyed
proof. Refresh report, Briefing, approval, and PR-draft claims that cite the old
candidate. Keep the ideation baseline, accepted scope, entity identity, issue,
owner, and proof that does not use either changed key.

This mechanism serves AC-3. Candidate-only invalidation is insufficient because a
clean merge can fail integration. Full refresh on every base change repeats the
original churn.

### 5. Push the approved SHA, not a moving branch name

After the final clean decision, push the recorded candidate to an explicit ref:

```text
git push origin CANDIDATE_SHA:refs/heads/BRANCH
```

Do not use force. Before PR creation, `git ls-remote` must report that exact SHA
for the remote branch. After PR creation, `gh pr view --json headRefOid` must
report the same SHA before EP2 records `pr:`. A mismatch stops delivery and leaves
state authority unchanged.

This mechanism serves AC-4. Rechecking local `HEAD` is insufficient because the
current mod pushes a branch name after the check.

## Out of scope

- Generic owner discovery, live-handle reuse, or fresh owner dispatch. G3/D8 owns
  these mechanisms.
- Automatic conflict resolution, force-push, or a generic resolver worker.
- A new stored field, state schema, PR provider, or merge strategy.
- Changes to state-checkout rebase safety. A same-entity state conflict remains a
  global state-safety halt.
- Changes to GitHub branch protection or merge-queue policy.

## Acceptance criteria

**AC-1 (VALUE) - A clean `1 1` candidate reaches PR preparation with zero candidate-head rewrites.**

Verified by: a real Git fixture invokes production `dispatch mergeability`. It
records a clean result, zero rebase commands, the same candidate SHA, and eligible
PR preparation. Shared reconcile and a focused Claude run leave the head unchanged.

**AC-2 - Every conflict result is NUL-safe and reaches the generic owner route without changing delivery authority.**

Verified by: production command fixtures cover ordinary paths, tabs, newlines,
directory/file conflicts, and a conflict with no path rows. G3/D8 receives the
structured package. Status, `pr`, `mod-block`, and pending authority stay
byte-identical. No resolver, consume, clear, supersede, rework, or force action runs.

**AC-3 - Integration evidence changes when any integration-key component changes.**

Verified by: a clean textual merge removes an API while its new base adds a caller.
The fixture requires new integration proof and observes the merged tree fail to
build. A matrix keeps candidate-only proof for a base-only change, but rejects old
integration proof. A candidate change invalidates both evidence classes.

**AC-4 - PR delivery publishes the approved candidate SHA even if the local branch moves.**

Verified by: an adversarial real-remote fixture changes the local branch after the
decision. The explicit refspec publishes only `candidate_sha`. Remote verification
and `headRefOid` match it before `pr:` can change. Branch-name push or any mismatch
must fail.

## Expected surface and semantic boundaries

The next Captain gate can approve this proposal: 11-14 files, approximately
`+450/-40` lines, with a 1.5x tolerance. This estimate is not pre-approved
implementation scope.

- `internal/dispatch/mergeability.go` and focused real-Git tests: command,
  Git-version error, byte parser, stable JSON, and evidence keys.
- `internal/dispatch/dispatch.go`: route and help for `dispatch mergeability`.
- `internal/dispatch/reconcile.go`, `reconcile_test.go`, and trunk tests: retain
  drift data but remove the shared rebase prescription.
- `mods/pr-merge.md` and `docs/dev/_mods/pr-merge.md`: call the production seam,
  delegate conflicts, refresh evidence by key, and push the pinned SHA.
- `skills/first-officer/references/claude-fo-dispatch.md`: make ancestry drift
  report-only.
- `internal/contractlint/`: check structural references and forbidden rebase
  prescriptions. It must not extract or execute instruction prose.
- `internal/ensigncycle/` shared fixture and Claude runner files: prove the live
  no-rebase result and G3/D8 handoff boundary.
- `docs/site/advanced/mods-and-standing-teammates.md`: document the user-visible
  mergeability, evidence-key, and pinned-head policy.

Allowed command grammar change: add `dispatch mergeability`. Allowed JSON change:
add its new stable result only. Stored formats do not change. Authority remains
with recorded gates and G3/D8 owner handoff. Runtime changes are limited to
ancestry-drift reporting, mergeability decisions, evidence refresh, and pinned
delivery.

The proposed documentation change is:

```diff
 The canonical example is the `pr-merge` mod: it opens the code-branch PR at merge,
 records the PR on the entity, and holds the terminal transition until the PR merges.
+It keeps a validated branch head when Git can merge it cleanly. Integration proof
+tracks the base, candidate, and merge tree. Conflicts return through the generic
+owner handoff, and PR delivery pushes the approved candidate SHA.
```

## Test plan

Add production tests before instruction changes. Do not extract or execute commands
from mod or skill prose.

1. Build real Git graphs for clean `1 1`, ordinary conflict, unusual filename,
   directory/file conflict, and pathless logical conflict cases. Invoke production
   `dispatch mergeability`. Check all JSON bytes, exit behavior, OID handling, and
   unchanged refs, index, and worktree.
2. Build the clean-merge/build-red API fixture. Materialize the returned merge tree
   and run `go test ./...`. Check that a base-only change invalidates integration
   proof while candidate-only proof remains valid.
3. Update shared reconcile fixtures. Check that `stale-branch` keeps ownership and
   drift facts but contains no rebase action or prescription.
4. Run the G3/D8 handoff fixture with ordinary and pathless conflicts. Check
   byte-identical delivery authority and no EP2-owned dispatch mechanism.
5. Move a local branch after the final decision. Push the recorded SHA to a real
   bare remote. Check the remote ref and a stubbed `headRefOid` before any `pr:`
   write. Reject a branch-name push and both mismatch mutations.
6. Run one isolated Claude fixture with a clean `1 1` branch. Check that reconcile
   does not change the head and that the merge boundary uses production JSON.
7. Run focused packages, `go test ./...`, `go test ./... -race`,
   `gofmt -w ./cmd ./internal`, `git diff --check`, and contract lints.

## Science Officer findings and authorized disposition

The five findings below are preserved unchanged from the Science Officer review.

1. **AC-3 is unsound:** base-only changes can merge textually but break integration (for example candidate removes an API while new base adds another-file caller). Candidate-only checks may key by candidate SHA, but integration/build proof must key by (base SHA,candidate SHA) or merge-tree OID; add a clean-merge/base-only integration-failure fixture.
2. **Conflict path command is invalid:** reusing `--quiet` with `--name-only` suppresses output; non-quiet output includes tree OID/messages. Specify distinct structured diagnostic parsing, ideally `-z`, `--no-messages`, explicit tree-OID handling; test unusual filenames/logical conflicts.
3. **Proof plan violates instruction-file quarantine** by extracting/comparing commands from mod prose. Use structural contractlint only for references and prove the decision through a production seam or runnable shared/live behavior.
4. **Expected surface omits shared reconcile behavior** at `internal/dispatch/reconcile.go` and its tests, which still prescribe `pull --rebase (owned)`. Include shared production/tests or justify why that reason remains.
5. **Exact-head authority has a probe-to-push race** because current delivery pushes a branch name. Push pinned SHA to explicit ref or verify remote head equals candidate before PR open/update; add adversarial branch-movement proof.

Authorized evidence ledger:

| Finding | Released user and normal workflow | Observable harm | Affected authority | Trigger evidence | Materiality | Owner | Disposition |
|---|---|---|---|---|---|---|---|
| 1 | Candidate delivery after trunk movement | Textual merge passes while merged program fails | AC-3 and terminal integration proof | Clean tree `0e813e0329`; build fails with `undefined: Old` | Material | EP2 | Fix with candidate and integration keys |
| 2 | Conflicting delivery with any valid Git path | Quiet output loses paths; line parsing corrupts names | AC-2 owner package | Tab/newline and directory/file spikes | Material | EP2 production seam | Remove quiet; parse NUL fields and tree OID |
| 3 | Merge-policy implementation and validation | Prose becomes its own behavioral oracle | AC-1 through AC-4 proof policy | Prior plan extracted mod commands | Material evidence defect | EP2 implementation/validation | Use production tests plus structural lint only |
| 4 | Shared and Claude stale-branch reconcile | Automatic rebase still rewrites clean candidates | AC-1 zero-rewrite value | `reconcile.go` reason and Claude action prescribe rebase | Material | EP2 shared reconcile surface | Make ancestry drift report-only in both layers |
| 5 | Captain-approved PR creation | Moving branch name can publish an unapproved head | AC-4 exact-head delivery | Explicit-SHA push retained the recorded head after local movement | Material | EP2 delivery surface | Pin SHA and verify remote plus `headRefOid` |

Captain decision A selects Git 2.38 compatibility without `--quiet`. Captain
decision B makes G3/D8 a dependency and forbids an EP2 owner-handoff duplicate.

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
