---
id: ep2cz3zsb2qpyyh889nyeqpr
title: Avoid unnecessary rebases before opening mergeable PRs
status: validation
source: GitHub issue #616; Captain intake 2026-08-05
started: 2026-08-04T17:00:42Z
completed:
verdict:
score: 1.0
worktree: .worktrees/spacedock-ensign-avoid-unnecessary-pr-rebases
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
        - id: gate:ep2cz3zsb2qpyyh889nyeqpr:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:ep2cz3zsb2qpyyh889nyeqpr-ideation-1
              briefing:
                id: briefing:ep2cz3zsb2qpyyh889nyeqpr:ideation:attempt-1:revision-1
                digest: sha256:5eefe3110c01a7a0675dfd03064a976d21683bdf80759e8002a6c20a9ebc3c8a
                request-digest: sha256:805c007f16d43aa151e85f77c044d41c23da74ae4df3925ee49c9b3928b2e65f
                room-ref: ./avoid-unnecessary-pr-rebases/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ep2cz3zsb2qpyyh889nyeqpr:ideation:1
                briefing: briefing:ep2cz3zsb2qpyyh889nyeqpr:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-05T02:11:03.127412Z"
                decision: approve
                reason: 'Captain approved the narrow two-mod issue #616 design after final Science Officer approval.'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:ep2cz3zsb2qpyyh889nyeqpr:validation
          stage: validation
          attempts:
            - id: gate-attempt:ep2cz3zsb2qpyyh889nyeqpr-validation-1
              briefing:
                id: briefing:ep2cz3zsb2qpyyh889nyeqpr:validation:attempt-1:revision-1
                digest: sha256:9d87e4f7f77c656e68e0c65c530fd86932902f3794af4eb5c81e7dbd363f4b4d
                request-digest: sha256:d7f9a328933b859a4e4cc44eac897ac3f276165cfffbfdbdb861061afe259ead
                room-ref: ./avoid-unnecessary-pr-rebases/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:ep2cz3zsb2qpyyh889nyeqpr:validation:1
                briefing: briefing:ep2cz3zsb2qpyyh889nyeqpr:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-05T06:41:48.81045Z"
                decision: revise
                reason: Fix the zsh exact-SHA refspec with braces in both approved mod files. Replace rigid exit-code semantics with First Officer inspection of merge-tree stdout, stderr, and context; actual conflict routes to G3/D8, command failure or uncertainty reports unknown and preserves authority. Add no machinery or preflight.
            - id: gate-attempt:ep2cz3zsb2qpyyh889nyeqpr-validation-2
              briefing:
                id: briefing:ep2cz3zsb2qpyyh889nyeqpr:validation:attempt-2:revision-1
                digest: sha256:87e9609b5ee1131cbce5daa1485b628a4b23b31f9de457f8146b0d11e39d4229
                request-digest: sha256:73e6b483241bdd8ea0d9d6a82e4ebc6a82121a90e4494ae76d2d047f777f8e1d
                room-ref: ./avoid-unnecessary-pr-rebases/review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:ep2cz3zsb2qpyyh889nyeqpr:validation:2
                briefing: briefing:ep2cz3zsb2qpyyh889nyeqpr:validation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-05T07:05:36.914037Z"
                decision: revise
                reason: 'Stamp the shipped pr-merge mod version 0.27.0. Apply issue #616 behavior only to mods/pr-merge.md. Restore docs/dev/_mods/pr-merge.md to its pre-EP2 customized bytes; the dev copy is reconciled later by refit, preserving its split-root customization.'
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

### Feedback Cycles

- Cycle 1: REJECTED — validation real-Git/zsh matrix; surface 2 files/34 LOC vs estimate 40 (85%); AC unchanged
- Cycle 2: REJECTED — Captain scope/version reset; surface 2 files/34 LOC vs estimate 40 (85%); AC unchanged

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

## Stage Report: implementation

- DONE: Change only the two approved pr-merge mod files to remove mandatory ancestry rebase and preserve the Captain-approved candidate SHA.
  Code commit `af64ffb01` changes only `mods/pr-merge.md` and `docs/dev/_mods/pr-merge.md`; both drafts expose the recorded SHA and neither delivery path rebases it.
- DONE: Implement the read-only clean/conflict/unknown merge-tree decision and exact-SHA refspec push while delegating conflicts to G3/D8.
  Fixture `/tmp/spacedock-pr-merge-cycle3.QOqi7m` would fail if clean were not exit 0, conflict not exit 1, unknown not exit 128, or the remote ref differed from the recorded SHA.
- DONE: Stay within about +30/-10 and prove the clean 1 1, conflict, unknown, and moved-local-branch cases without committed prose-grep tests.
  Commit `af64ffb01` is two files at +25/-9; the direct real-Git fixture preserved clean HEAD/refs/index/worktree and pushed `527936e` after the local branch moved to `5fe8104`.
- FAILED: Run the repository-wide normal and race suites without environmental fixture failures.
  Both commands passed every other package but `internal/gates` expected seven entity paths absent from the shared `.spacedock-state` checkout; focused `internal/contractlint` and `skills/integration` tests passed.

### Summary

The two pr-merge policies now record and display the approved candidate, decide
mergeability without mutation, fail closed on conflict or error, and push the
recorded SHA. No command, shared reconciliation, host, or evidence machinery changed.

## Stage Report: implementation

- DONE: Change only the two approved pr-merge mod files to remove mandatory ancestry rebase and preserve the Captain-approved candidate SHA.
  Candidate `af64ffb01` remains unchanged and differs from `main` only in `mods/pr-merge.md` and `docs/dev/_mods/pr-merge.md` at +25/-9.
- DONE: Implement the read-only clean/conflict/unknown merge-tree decision and exact-SHA refspec push while delegating conflicts to G3/D8.
  Fixture `/tmp/spacedock-pr-merge-cycle3.QOqi7m` falsifies wrong exit classification, candidate mutation, or a remote head different from the approved SHA.
- DONE: Stay within about +30/-10 and prove the clean 1 1, conflict, unknown, and moved-local-branch cases without committed prose-grep tests.
  The direct real-Git fixture observed clean=0, conflict=1, unknown=128, unchanged clean state, and remote `527936e` after the local branch moved to `5fe8104`.
- DONE: Re-run `go test ./...` and `go test ./... -race` from the current exact candidate after current state convergence and diagnose failures from this run.
  Auto-discovered live state reproduced one failing manifest test because eight fixed pilot paths had archived; every other package passed in both fresh runs.
- DONE: Use a supported fixture setup to make both suites genuinely pass and record exact evidence with no false green.
  With `SPACEDOCK_STATE_ROOT` set to a snapshot of state commit `73f41e2a2232ebb561710bce568641ec976d5f3d`, the focused 31-path test and both full suites exited 0.
- DONE: Keep the code surface exactly `mods/pr-merge.md` and `docs/dev/_mods/pr-merge.md`; do not expand it.
  Candidate inspection after all reruns still reports exactly those two files and no worktree changes.
- DONE: Commit/push the state report and candidate branch before completion.
  Exact candidate `af64ffb018a88679aac1972c5af157ce64a43502` is published at `origin/spacedock-ensign/avoid-unnecessary-pr-rebases`; this report is published path-scoped on `spacedock-state/dev`.

### Summary

The prior red was live-state drift in a fixed release-pilot manifest test, not a
candidate regression. The supported immutable state fixture makes both complete
suites green while the approved two-file candidate remains byte-for-byte unchanged.

## Stage Report: validation

- DONE: Validate commit `af64ffb01` only and keep the approved implementation surface unchanged.
  `af64ffb018a88679aac1972c5af157ce64a43502` remains a clean two-file `+25/-9` candidate; `gofmt`, `git diff --check`, and HEAD/status snapshots changed no candidate bytes.
- FAILED: Verify the two mod files preserve the Captain-approved candidate SHA, use read-only merge-tree classification, and push the exact SHA without rebasing.
  Both record the SHA and use non-mutating `merge-tree`, but their exact zsh push refspec becomes `<SHA>efs/heads/<branch>` and fails before PR creation; no rebase occurred.
- DONE: Reproduce clean, conflict, unknown-error, and moved-local-branch cases with real Git evidence and confirm no mutation of candidate/index/worktree.
  `/tmp/spacedock-validation-616.Lyum7i` observed clean=0 on a `1 1` graph, conflict=1, command-error=1, and unchanged HEAD/refs/index/worktree/authority/action markers for every merge-tree exercise.
- FAILED: Verify **AC-1 (VALUE) - A clean `1 1` candidate reaches PR creation with its approved SHA.**
  After the local branch moved, zsh interpreted `"$CANDIDATE_SHA:refs/heads/candidate"` via its `:r` modifier; real Git exited 1 with no remote branch, while brace-delimited control pushed the approved SHA exactly.
- FAILED: Verify **AC-2 - Conflict and unknown results stop delivery without changing authority.**
  Stops and mutation guards held, but an unavailable-object command error exited 1 with `not something we can merge`, so the prose falsely calls unknown failure a content conflict and refers it to G3/D8.
- DONE: Classify the AC-1 finding on both required axes.
  Material outcome defect: current zsh is an observed normal runtime; the promised clean PR is not created. Rework must brace the refspec as `${CANDIDATE_SHA}:refs/heads/{branch}` in both approved files.
- DONE: Classify the AC-2 finding on both required axes.
  Material outcome defect: exit 1 is not uniquely a content-conflict verdict. Per Captain clarification, rework must require FO inspection of stdout/stderr/context and classify uncertainty as unknown, without new machinery or rigid preflight.
- DONE: Run the required focused and repository-wide suites; diagnose the reported internal/gates fixture-path failures from this exact head and classify them from current evidence.
  Focused contract/integration suites passed; both full suites failed only the pilot manifest test on eight now-archived live-state paths, while its immutable state `73f41e2a2` control passed and its test blob equals `main`.
- FAILED: Recommend PASSED with no material finding.
  Recommend REJECTED: AC-1 and AC-2 each have a material outcome defect; the unrelated live-state fixture drift does not contribute to rejection.

### Summary

The read-only mergeability exercise preserves candidate authority, but the exact
zsh refspec prevents clean PR delivery and exit 1 cannot by itself establish a
content conflict. Rework the two mod files as directed, then rerun this real-Git matrix.

## Stage Report: implementation (cycle 2)

- DONE: Brace the recorded candidate variable in the exact-SHA push refspec in both approved mod files so supported zsh pushes the approved SHA exactly.
  Commit `a20331ef0` uses `${CANDIDATE_SHA}:refs/heads/{branch}` twice; fixture `/tmp/spacedock-ep2-zsh-matrix.YC9RZY` would fail if remote `dff1890` followed moved local head `e745af0`.
- DONE: Replace rigid exit-code semantics with First Officer inspection of merge-tree output and context: actual content conflict routes to G3/D8; command failure or uncertainty reports unknown and preserves authority.
  The same zsh fixture observed conflict and unavailable-object results both at exit 1 but distinguished `CONFLICT (content)` stdout from `not something we can merge` stderr, with refs/index/worktrees unchanged.
- DONE: Keep the two-file scope and rerun the real-Git zsh clean/conflict/unknown/moved-branch matrix plus focused, full, and race suites with exact evidence.
  Candidate `a20331ef0` remains two files at +25/-9; contract/integration, full, and race suites exited 0 using `SPACEDOCK_STATE_ROOT` snapshot `73f41e2a2` for the fixed 31-path pilot manifest.

### Summary

The corrected mods preserve zsh refspec syntax and treat merge-tree status as one
signal rather than a verdict. Actual conflict evidence routes to G3/D8, while tool
failure or ambiguity stops as unknown without changing pending authority.

## Stage Report: validation (cycle 2)

- DONE: Validate corrected candidate `a20331ef0` only and preserve the approved two-file surface.
  `a20331ef0bceeca6fba34b5c607a63b772b7dc8e` remains clean at `+25/-9`; `gofmt`, `git diff --check`, HEAD, and status snapshots changed no candidate bytes.
- DONE: Verify the two mod files preserve the Captain-approved candidate SHA, use read-only merge-tree classification, and push the exact SHA without rebasing.
  Both record the SHA, expose merge-tree stdout/stderr/status/context for semantic judgment, forbid rebase, and brace `${CANDIDATE_SHA}:refs/heads/{branch}` for zsh.
- DONE: Verify **AC-1 (VALUE) - A clean `1 1` candidate reaches PR creation with its approved SHA.**
  `/tmp/spacedock-validation-616-cycle2.XaIkKI` observed clean exit 0 with unchanged HEAD/refs/index/worktree, then zsh pushed recorded `772d1a81` while the local branch stood at `3e128394`.
- DONE: Verify **AC-2 - Conflict and unknown results stop delivery without changing authority.**
  Equal exit-1 cases were distinguished atomically: `CONFLICT (content)` routed to G3/D8, while unavailable-object stderr reported unknown; neither changed refs/index/worktree/authority/actions.
- DONE: Reproduce clean, conflict, unknown-error, and moved-local-branch cases with real Git evidence and confirm no mutation of candidate/index/worktree.
  The matrix asserted exact identity and terminal action state across all four variants; only the authorized clean exact-SHA push created a remote branch.
- DONE: Run the required focused and repository-wide suites; diagnose the reported internal/gates fixture-path failures from this exact head and classify them from current evidence.
  Contract/integration and immutable-state pilot checks passed; `go test ./...` and `go test ./... -race` passed with state `73f41e2a2`; live state alone lacks eight archived manifest paths.
- DONE: Recommend PASSED with deferred risks separated from material findings.
  Recommend PASSED: both prior material outcome defects are corrected, all AC evidence is valid, and no material, deferred-risk, or polish finding remains.

### Summary

The corrected two-file policy preserves the approved SHA through mergeability
judgment and zsh delivery. AC-1 and AC-2 pass the full adversarial matrix, and the
remaining live-state test red is independently proven fixture drift.

## Stage Report: implementation (cycle 3)

- DONE: Change only mods/pr-merge.md: preserve the validated mergeability-first exact-SHA behavior and stamp version 0.27.0.
  Commit `2acb3630b` leaves the final branch diff at one shipped mod, +13/-4, with semantic behavior unchanged from validated `a20331ef0` and version `0.27.0`.
- DONE: Restore docs/dev/_mods/pr-merge.md byte-for-byte to its pre-EP2 customized baseline; leave reconciliation to a later refit.
  Working and `main` baseline blobs both hash to `44e7b345eee634b35780efe3106ad827b424661b45dd75aaeaec9fd90f38c4f2`; the dev mod is absent from the final diff.
- DONE: Prove the final branch changes exactly one shipped mod and rerun the real-Git zsh matrix plus focused, full, and race suites.
  Fixture `/tmp/spacedock-ep2-zsh-final.ZA50IQ` would fail on mutation, wrong conflict/unknown judgment, or moved-head delivery; focused, full, and race suites all exited 0 with state snapshot `73f41e2a2`.

### Summary

The final candidate applies issue #616 only to the shipped `pr-merge` mod and
stamps it `0.27.0`. The customized development mod is restored exactly and remains
owned by the later refit path.
