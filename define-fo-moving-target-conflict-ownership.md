---
title: Define FO ownership and recovery for moving-target PR conflicts
status: validation
source: "Captain follow-up after the 2026-08-03 durable-decisions conflict diagnosis."
started: 2026-08-03T16:01:22Z
completed:
verdict:
score: 0.95
worktree: .worktrees/spacedock-ensign-define-fo-moving-target-conflict-ownership
issue:
sprint: durable-decisions
group: fo-contract
milestone: 0.27.0
id: g3912c6f6jkgd0yjmyg6h7yn
gates:
    version: 1
    records:
        - id: gate:g3912c6f6jkgd0yjmyg6h7yn:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:g3912c6f6jkgd0yjmyg6h7yn-backlog-1
              briefing:
                id: briefing:g3912c6f6jkgd0yjmyg6h7yn:backlog:attempt-1:revision-1
                digest: sha256:fbe944ff1e7905c23765219faa581bc2c1cba4980f68664253038adacb2afafa
                request-digest: sha256:7c706e3f3ab394a0653b087e1007f2aab9246a08c95e8c8947bb976a60d2736a
                room-ref: ./define-fo-moving-target-conflict-ownership/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:g3912c6f6jkgd0yjmyg6h7yn:backlog:1
                briefing: briefing:g3912c6f6jkgd0yjmyg6h7yn:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-03T15:59:40.23286Z"
                decision: approve
                reason: Captain conn authorizes approval when SO concurs. SO concurs with parallel EJ/G3 ideation; keep EJ-before-G3 implementation landing because they share the FO dispatch and merge contract.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:g3912c6f6jkgd0yjmyg6h7yn:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:g3912c6f6jkgd0yjmyg6h7yn-ideation-1
              briefing:
                id: briefing:g3912c6f6jkgd0yjmyg6h7yn:ideation:attempt-1:revision-1
                digest: sha256:4d10f045deb6106fb8b30200e4e726677c24b73aa534c4cdffe7a86d3fba85ce
                request-digest: sha256:4459edfb7e1c00dca121c198716a0b3fa42dd76a97fa07d5d8479ad75ff9e69e
                room-ref: ./define-fo-moving-target-conflict-ownership/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:g3912c6f6jkgd0yjmyg6h7yn:ideation:1
                briefing: briefing:g3912c6f6jkgd0yjmyg6h7yn:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-03T16:19:06.88638Z"
                decision: approve
                reason: Captain conn authorizes approval when SO concurs. The ideation fixture proves pending authority preservation, conflict abort and ownership handoff, exact-head freshness, and keep-moving without a resolver. Hold G3 implementation until EJ lands; SO concurs with EJ-before-G3 implementation ordering.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:g3912c6f6jkgd0yjmyg6h7yn:validation
          stage: validation
          attempts:
            - id: gate-attempt:g3912c6f6jkgd0yjmyg6h7yn-validation-1
              briefing:
                id: briefing:g3912c6f6jkgd0yjmyg6h7yn:validation:attempt-1:revision-1
                digest: sha256:dcbf98f0298e3d33367275130829e79d68bd5abde0570b3f4c04ac76a268f85d
                request-digest: sha256:1824eb5c77179ac9e3e053c6146ff87fcfbdb564bd19ef4a193105178dc2ae0b
                room-ref: ./define-fo-moving-target-conflict-ownership/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:g3912c6f6jkgd0yjmyg6h7yn:validation:1
                briefing: briefing:g3912c6f6jkgd0yjmyg6h7yn:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-06T00:44:23.012474Z"
                decision: revise
                reason: Detached validation proved AC-5 can pass while equivalent force, auto-resolution, and side-discard spellings are accepted. Return to implementation for the authorized narrow oracle and mutation-matrix fix only.
            - id: gate-attempt:g3912c6f6jkgd0yjmyg6h7yn-validation-2
              briefing:
                id: briefing:g3912c6f6jkgd0yjmyg6h7yn:validation:attempt-2:revision-1
                digest: sha256:a0f74c17ae324e559cecb54f62d476c8882279d06a8aadff5be1bc3d5374b392
                request-digest: sha256:deb11ba04282cf83df6132218841c75ea64380ccf513dc12bebcd670a91834c9
                room-ref: ./define-fo-moving-target-conflict-ownership/review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:g3912c6f6jkgd0yjmyg6h7yn:validation:2
                briefing: briefing:g3912c6f6jkgd0yjmyg6h7yn:validation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-06T01:07:21.458184Z"
                decision: revise
                reason: Cycle-2 detached validation proved standard merge/rebase ours-strategy forms still evade AC-5. Return to implementation for the authorized narrow existing-oracle and matrix fix only.
            - id: gate-attempt:g3912c6f6jkgd0yjmyg6h7yn-validation-3
              briefing:
                id: briefing:g3912c6f6jkgd0yjmyg6h7yn:validation:attempt-3:revision-1
                digest: sha256:922aef932c1cae15d40f0b5dea0be7b3f83432f0529f297351278c0ff36d9e87
                request-digest: sha256:aa3b58c0f40321424501c76abc4528751ca1e184e30cab6a092c7ecf53e8943a
                room-ref: ./define-fo-moving-target-conflict-ownership/review/validation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:g3912c6f6jkgd0yjmyg6h7yn:validation:3
                briefing: briefing:g3912c6f6jkgd0yjmyg6h7yn:validation:attempt-3:revision-1
                by: person:captain
                at: "2026-08-06T03:08:28.159189Z"
                decision: revise
                reason: Remove the test command parser and command-spelling matrix. A parser in a test is a red flag. Capture actual Git arguments at the fixture boundary and compare those values directly.
---

Define the ownership, durable hold, evidence, and keep-moving behavior when a pending PR meets a moving merge target. The task closes the contract ambiguity without adding a resolver worker or changing the workflow definition.

## Problem

The First Officer contract correctly forbids automatic conflict resolution, force-push, and discarding either side, but it does not name the owner or the durable transition when an open PR becomes dirty because the merge target moved. A First Officer can therefore mistake a filtered empty `status --next` for a session stop, ask the wrong worker to resolve a two-writer conflict, clear a pending terminal approval, or let stale validation evidence authorize a new head. The current PR and state must remain truthful while unrelated sprint work continues.

## Proposed approach

Make ownership explicit. The code/worktree owner may rebase and commit only within its owned branch when the Captain authorizes that ordinary branch maintenance; the Captain owns semantic reconciliation of two-writer conflicts and decides whether to rework, rebase, or close/reopen. The First Officer owns transport and state safety: abort any conflicted local rebase, surface the PR number, old and new base/head commits, and conflict paths, and preserve the affected task's current status, bare `pr:#N`, `mod-block: merge:pr-merge`, and pending terminal approval. No conflict path clears or consumes authority, terminalizes, auto-resolves, force-pushes, or reuses old validation evidence.

After the Captain or owning branch author produces a new exact head, the First Officer records the new head as the delivery candidate, reruns the required CI and fresh validation/gate evidence, and only then resumes the merge ceremony. The event loop holds only the affected task and keeps independently ready tasks dispatching or presenting. The contract must name this as a per-entity conflict hold, not a generic resolver route.

## Out of scope

- Automatic conflict classification, conflict resolution, `git checkout --ours/--theirs`, force-push, or force-finalization.
- A standing resolver worker, new conflict database, new frontmatter field, or workflow-specific finding taxonomy.
- Changes to the commissioned workflow definition, PR provider, branch protection, or GitHub merge policy.
- Reusing a superseded approval, old-head CI, or prior gate Briefing as proof for a rewritten candidate.

## Acceptance criteria

**AC-1 (VALUE) — A moving-target conflict cannot finalize or lose a pending decision.**
Verified by: a local Git/PR fixture that changes the merge target after an approval is pending; the observed state keeps the affected task at its current stage with `pr:#N`, `mod-block: merge:pr-merge`, and the pending terminal application, while the command log records no terminalization, consume, worker dispatch, force operation, or auto-resolution and surfaces exact conflict paths plus peer/base/head evidence.

**AC-2 — Ownership and the required transition are unambiguous.**
Verified by: a behavior fixture that exercises an FO-held conflict, Captain-authorized branch rebase, and new-head handoff; the trace attributes semantic resolution to the Captain/owning branch and transport/state guarding to the FO, aborts failed rebases, and requires a fresh exact-head validation before merge guard can finalize.

**AC-3 — Stale evidence cannot authorize a rewritten head.**
Verified by: an exact-head fixture that changes the PR head after a green validation result and asserts that the old CI/gate evidence is rejected or held until the new head's required lanes and fresh gate evidence pass.

**AC-4 — Unrelated sprint work keeps moving.**
Verified by: a multi-entity event-loop fixture with one conflict-held PR and two independent ready entities; the conflicted task remains held with its evidence while both unrelated entities reach their next dispatch or gate action in the same loop sequence.

**AC-5 — The generic contract does not create a resolver escape hatch.**
Verified by: a detached adversarial audit that attempts auto-resolution, `-X ours`/`-X theirs`, force-push, direct frontmatter clearing, and a generic resolver spawn; each attempt fails or remains outside the accepted trace, while the supported Captain-owned reconciliation path remains executable.

## Expected surface and semantic boundaries

Expected change is limited to the generic First Officer conflict/merge/dispatch contract and its behavior fixtures: `skills/first-officer/references/first-officer-shared-core.md`, `skills/first-officer/references/fo-merge-core.md`, `skills/first-officer/references/fo-dispatch-core.md`, relevant host adapter notes, and focused conflict/event-loop tests. Estimate 4–9 files and +150/-30 lines, with a 2x tolerance. Allowed semantic change: owner attribution, evidence, and per-entity hold/re-entry ordering. Stored state schema, workflow README, PR provider behavior, and merge authority remain unchanged.

## Test plan

Start with a throwaway two-branch fixture whose base advances after an approved PR is recorded. Capture the failed rebase paths and peer commit, abort cleanly, and assert byte-preserved state. Add command-log tests for the Captain/branch-owner handoff, exact-head freshness, and keep-moving matrix. Run focused status/dispatch/merge tests, a detached adversarial audit on the contract/guard surface, then `go test ./...`, `go test ./... -race`, formatting, and all required host lanes if the runtime adapters change.

## Stage-specific test gates

- Ideation must replay one moving-target conflict and record the exact pre/post state and ownership handoff before selecting wording or a helper.
- Implementation must preserve pending authority bytes, prove no auto-resolution or force path, and add no resolver worker or state field.
- Validation must reproduce the conflict/abort/new-head sequence, stale-evidence rejection, and unrelated-task continuation, plus focused/full/race/format and the required detached audit/live evidence.

## Stage Report: ideation

- DONE: Replay one moving-target PR conflict and record the exact pre/post state, conflict path, and Captain/branch-owner/FO handoff.
  Evidence: a throwaway Git/PR fixture at `/tmp/spacedock-g3-conflict.oWmCsG` used `main` as the merge target and `pr-7` as the open PR. Before the target moved, the approved-but-not-consumed state was `status: validation`, `pr: '#7'`, `mod-block: merge:pr-merge`, and `terminal-application: pending`; its exact state-file SHA-256 was `67ed11855aacfbd6e6254ed9c6629c7bdb33398516606a822402b169ca05da48`. The old base/merge-base was `1ee5119d9828811d7ce26adac08d020bc3a65640`; the old PR head was `7918cba8d1f765a1805efd7e016366b73cd2b9f3`. After `main` advanced to base `0e97600e663d56beb807067e2b175ff6f739e6f1`, `git rebase main` exited 1 with `UU contract.txt`, `REBASE_HEAD=7918cba8d1f765a1805efd7e016366b73cd2b9f3`, and conflict path `contract.txt` (the peer/base commit was `0e97600…`). The FO ran `git rebase --abort`; afterward `pr-7` was back at `7918cba…`, `main` remained `0e97600…`, the worktree was clean, no unmerged paths remained, and the state-file SHA-256 was unchanged. No terminalization, gate consume, worker dispatch, force operation, or automatic resolver command occurred. The Captain then authorized ordinary branch maintenance; the owning branch author manually reconciled `contract.txt`, producing exact new head `e54865db2dff623ab9d7bfe7ba069ec2aae6226c` on base `0e97600…` (ahead/behind `0/1`).

- DONE: Reconcile pending authority, exact-head freshness, keep-moving, and no-resolver constraints against the existing merge and dispatch contract, and name the remaining risk.
  Evidence: `fo-merge-core.md` makes the terminal-target application remain `pending` after `gate consume`; only `merge guard` may consume it after proven delivery, while retryable delivery trouble leaves the same approval pending. That preserves `pr:#7`, `mod-block: merge:pr-merge`, and the current stage during the conflict. `first-officer-shared-core.md` already requires aborting a conflicted rebase and surfacing paths/peer evidence, with no force-push, `-X ours/theirs`, or discard; `fo-merge-core.md` applies the same no-auto-resolve rule to local merge. `fo-dispatch-core.md` checks blocked entities before `status --next`, so the affected entity can be held while unrelated ready entities continue; the owned-worktree rule in the Claude adapter permits a rebase only when `drift.owned == true` and makes an unowned branch report-only. The old green evidence was tied to head `7918cba…`; after the owner handoff the exact-head check rejected it against `e54865d…`, so required CI and fresh gate evidence must run on the new head before merge guard resumes. The unresolved risk is scope ambiguity: the current generic word “HALT” can be read as stopping the whole event loop, and it does not identify the Captain/owning-branch semantic owner or the exact old/new base/head tuple. Implementation must distinguish a code-worktree PR conflict (per-entity hold, keep unrelated work moving) from a split-root state-sync conflict (global state safety halt), and re-check the PR head immediately before fresh validation and delivery.

- DONE: Select the smallest implementation direction and recommend entry to the implementation gate.
  Evidence: no helper, resolver worker, frontmatter field, workflow change, or provider change is needed. Add contract wording and focused behavior fixtures on the existing shared-core, merge-core, dispatch-core, and host-adapter surfaces; make the fixture assert byte-preserved pending authority, abort/path/peer evidence, explicit Captain/branch-owner handoff, exact-head invalidation, and independent dispatch. Recommend APPROVE ideation and enter implementation with this per-entity conflict-hold/re-entry contract, preserving the existing sole terminal consumer and merge authority.

### Summary

The real moving-target replay shows that a dirty PR must remain truthful and pending while the FO aborts and reports the conflict. Semantic reconciliation belongs to the Captain and owning branch; the FO owns transport/state safety and fresh exact-head evidence. Unrelated ready work may continue, but a rewritten head cannot reuse old validation. Proceed to implementation with the existing contract surfaces and no resolver escape hatch.

## Stage Report: implementation

- DONE: Commit the minimal First Officer contract and fixtures that preserve pending authority and exact conflict evidence during a moving-target PR conflict.
  Code commit `686ebd8e8` changes six files: the shared/merge/dispatch/Claude contracts plus a real-Git fixture; removing the abort, byte equality, exact tuple, or pending-state behavior makes `TestMovingTargetConflictOwnershipAndRecovery` fail.
- DONE: Prove Captain/branch-owner/First-Officer responsibility, conflicted-rebase abort, fresh exact-head evidence, and unrelated-work continuation without any resolver path.
  The fixture conflicts on `contract.txt`, aborts to the old head, distinguishes semantic owner from transport guard, rejects old-head CI/gate evidence, and requires both unrelated dispatch and gate actions; isolated adversarial mutations red on terminalization, consume, ours/theirs, force, resolver dispatch, state mutation, or whole-loop hold.
- DONE: Stay within the approved 4-9 file surface and tolerance; run focused, formatting, full, race, and required detached or live checks.
  Surface is 6 files and +212/-4 lines; focused and focused-race pass, `gofmt -w ./cmd ./internal` and `git diff --check` pass, and detached commit `686ebd8e8` passes both behavior/adversarial tests. Full and race runs exercised all packages; their sole failure is `internal/gates` reading the same nine absent manifest paths from the shared `.spacedock-state`, while every code/worktree package passes.

### Summary

The First Officer now holds only the conflicted PR entity, preserves its pending authority and exact conflict evidence, and keeps unrelated work moving. Captain-authorized branch ownership and exact-new-head revalidation are explicit, with no resolver, force path, provider change, state field, or workflow edit.

### Verification Addendum

Verification pinned `SPACEDOCK_STATE_ROOT` to clean detached state snapshot `fa0804a6a8a5f98858992c6b2f1b49bf609aa5c5`. At unchanged clean candidate head `686ebd8e88b5e8a49a6f05f15bfc5a6dd4d0816d`, both `go test ./...` and `go test ./... -race` exercised all packages; every package except `internal/gates` passed and no race was reported.

The exact remaining failure in both commands is `TestV1PilotManifestReadsAndValidates`: its manifest names nine paths absent from the immutable snapshot — `bind-post-rework-briefing-at-rejection-regate.md`, `codex-launch-multi-agent-v2.md`, `collapse-gate-approval-ceremony/index.md`, `cut-workflow-specific-round-recorder-from-v1/index.md`, `gate-agent-ergonomics.md`, `minimize-v1-gate-application-schema/index.md`, `shared-git-scaffold-helper.md`, `status-pagination-and-default-sorting.md`, and `status-where-robust-and-discoverable.md`. The identical stable result rules out concurrent archival during either run; candidate bytes and HEAD remained unchanged and clean.

## Review-finding disposition

- Reviewer observation: detached audit `/tmp/spacedock-g3-audit.fiwXP6` at exact candidate `686ebd8e88b5e8a49a6f05f15bfc5a6dd4d0816d` added only a throwaway test; `gradeMovingTargetTrace` accepted `git rebase -X theirs main`, `git checkout --ours contract.txt`, and `git push -f origin HEAD`.
- Released user and normal workflow: a First Officer handling the supported moving-target PR conflict can express auto-resolution, side discard, or force with these ordinary Git spellings.
- Observable harm: the adversarial oracle remains green despite prohibited conflict actions, so AC-5 can appear proven while the accepted trace violates the shipped safety contract.
- Authority: `contract[skills/first-officer/references/first-officer-shared-core.md#haltrebase-conflictpaths-abort-surface-stop]` requires never force, auto-resolve, discard either side, or dispatch a resolver.
- Trigger evidence: `go test ./internal/ensigncycle -run '^TestDetachedAuditRejectsEquivalentForbiddenSpellings$' -count=1` failed with all three commands reported as accepted forbidden actions.
- Worker proposal: Material / evidence defect / G3-owned; narrowly extend the existing behavior oracle and mutation matrix to reject semantically equivalent auto-resolution, side-discard, and force spellings.
- First Officer authorization: `FIX`, narrowly limited to that oracle/matrix correction; validator must not mutate the candidate.
- Validator recommendation: `REJECTED`; candidate stayed clean and unchanged at `686ebd8e88b5e8a49a6f05f15bfc5a6dd4d0816d`.

## Stage Report: validation

- DONE: Reproduce AC-1 through AC-5 at exact candidate 686ebd8e88b5e8a49a6f05f15bfc5a6dd4d0816d, including byte-preserved pending authority, exact conflict tuple and abort, ownership split, stale-evidence rejection, and unrelated-work continuation.
  AC-1 through AC-4 reproduce in the real-Git fixture: abort restores the old head and byte-identical pending state, the trace carries old/new base, old head, conflict path, distinct Captain/branch-owner and FO roles, rejects stale exact-head evidence, and continues task B dispatch plus task C gate presentation. AC-5 is not validly established because the detached audit found equivalent forbidden spellings the oracle accepts.
- DONE: Run focused, formatting, full, and race evidence; compare the stable internal/gates pilot-manifest failure with current origin/main and classify its candidate ownership without editing either G3 or the unrelated manifest.
  Focused and focused-race pass; `gofmt -d ./cmd ./internal`, `git diff --check`, detached contract lint, and skill integration pass. Full and race runs against state `fa0804a6a8a5f98858992c6b2f1b49bf609aa5c5` fail only the nine-path pilot-manifest test with no race report; current `origin/main` `8b5af99baa5c37fe7c969a904819041688420e22` fails the same test for only two absent paths, so the output is not identical but the failure is unrelated and not G3-owned.
- FAILED: Perform the required detached adversarial audit on the shipped First Officer contract and fixtures, verify no resolver/force/provider/state-schema escape hatch, and recommend PASSED or REJECTED with exact evidence.
  The shipped prose and six-file surface add no resolver, provider, or state-schema route, and named mutations pass, but the throwaway exact-head audit proves the oracle accepts `-X theirs`, `checkout --ours`, and `push -f`; recommend REJECTED for the authorized narrow evidence fix.

### Summary

The exact candidate preserves pending authority, conflict identity, ownership separation, fresh-head gating, and unrelated-work continuation, while the unrelated pilot-manifest failure remains outside G3. Validation rejects because the adversarial observation boundary accepts semantically equivalent forbidden Git actions; revise only that existing oracle and mutation matrix, then rerun focused, detached, full, and race evidence.

### Feedback Cycles

- Cycle 1: REJECTED — detached validation / equivalent forbidden Git spellings; surface 6 files/+212/-4 vs estimate 4-9 files/+150/-30; AC unchanged
- Cycle 2: REJECTED — detached re-review / ours strategy forms; surface 6 files/+263/-4 vs estimate 4-9 files/+150/-30; AC unchanged

## Stage Report: implementation (cycle 2)

- DONE: Make the existing behavior oracle and mutation matrix reject the demonstrated equivalent force, auto-resolution, and side-discard spellings, including `git rebase -X theirs main`, `git checkout --ours contract.txt`, and `git push -f origin HEAD`.
  Commit `d4a3590dd` replaces the incomplete substring list with action/token checks; the matrix independently reds the three exact findings plus attached/separate strategy options, checkout/restore sides, force/lease, bundled `-fu`, and forced `+refspec` variants.
- DONE: Keep the supported Captain-owned reconciliation and clean abort/new-head trace green without expanding product contract, provider, schema, resolver, or workflow surfaces.
  `TestMovingTargetConflictOwnershipAndRecovery` remains green on the real conflict/abort, byte-preserved authority, Captain handoff, and exact-new-head evidence; correction changes only the two pre-existing oracle/matrix test files.
- DONE: Run focused and focused-race tests, reproduce the detached adversarial mutations, run formatting/full/race evidence, and report exact head and surface; do not repair the unrelated pilot manifest.
  Focused and focused-race pass; detached exact head `d4a3590ddd7ad7e8b3aec022e462e5ecfda6b9fb` passes supported/adversarial tests; `gofmt` and `git diff --check` pass. Candidate surface is 6 files/+263/-4 total (cycle correction 2 files/+60/-9). Full and race on immutable state `8d31ebd4e` pass every package except the unchanged nine-path `internal/gates` pilot-manifest failure; no race is reported.

### Summary

The corrected oracle now rejects the validator's exact escape spellings and adjacent Git equivalents without changing the First Officer contract or supported reconciliation trace. The candidate is clean at `d4a3590dd`; only the explicitly out-of-scope pilot manifest remains red in aggregate suites.

## Review-finding disposition (cycle 2)

- Reviewer observation: detached audit `/tmp/spacedock-g3-audit-cycle2.GR01pk` at exact corrected candidate `d4a3590ddd7ad7e8b3aec022e462e5ecfda6b9fb` added only a throwaway test; the corrected oracle rejected the original three commands, attached/separate `-X` and `--strategy-option` forms, force-with-lease, and forced `+refspec`, but accepted `git merge -s ours main` and `git rebase --strategy=ours main`.
- Released user and normal workflow: a First Officer handling the supported conflict can invoke these standard Git merge/rebase strategy options.
- Observable harm: the oracle remains green when the `ours` strategy auto-resolves or discards peer/rebased changes, so AC-5 can appear proven while observable conflict safety is wrong.
- Authority: `contract[skills/first-officer/references/first-officer-shared-core.md#haltrebase-conflictpaths-abort-surface-stop]` requires never auto-resolve or discard either side.
- Trigger evidence: `go test ./internal/ensigncycle -run '^TestCycle2Detached(ForbiddenVariantMatrix|SupportedTraceHasNoFalsePositive)$' -count=1` failed only because those two forbidden strategy commands were accepted; the supported trace was not false-rejected.
- Worker proposal: Material / evidence defect / G3-owned; extend only the existing oracle/matrix to reject merge/rebase strategy forms selecting `ours` while preserving the supported trace.
- First Officer authorization: `FIX`, narrowly limited to the existing oracle/matrix; validator must not mutate the candidate.
- Validator recommendation: `REJECTED`; candidate stayed clean and unchanged at `d4a3590ddd7ad7e8b3aec022e462e5ecfda6b9fb`.

## Stage Report: validation (cycle 2)

- FAILED: Re-review AC-1 through AC-5 at exact corrected head d4a3590ddd7ad7e8b3aec022e462e5ecfda6b9fb and confirm the supported conflict/abort/authority/new-head/unrelated-work trace remains green.
  The detached package proves the supported trace remains accepted and the correction does not touch AC-1 through AC-4 contract surfaces, but AC-5 remains unproven because two valid `ours` strategy forms pass the oracle.
- FAILED: Re-run the detached adversarial audit for the three original escape spellings and adjacent option/refspec variants; prove the corrected oracle rejects them without false-rejecting the supported trace.
  The three originals and tested `-X`, `--strategy-option`, force-with-lease, and `+refspec` variants reject correctly, and `git push origin HEAD` plus the supported trace remain green; `git merge -s ours` and `git rebase --strategy=ours` are accepted, so the adjacent-option matrix fails.
- SKIPPED: Run focused, formatting, full, and race evidence, classify only current-head findings, and recommend PASSED or REJECTED without repairing the unrelated pilot manifest.
  The current-head material finding entered the checkpoint before expensive reruns; under the authorized `FIX` disposition the validator kept candidate bytes unchanged and stopped. Recommend REJECTED; the unrelated pilot manifest was not edited or used to classify G3.

### Summary

Cycle 2 closes the original three escape spellings without false-rejecting the supported trace, but its token oracle still misses valid `ours` strategy selection through `-s` and `--strategy`. Revise only those existing oracle/matrix paths, then rerun the exact supported, adversarial, focused, formatting, full, and race evidence.

## Stage Report: implementation (cycle 3)

- DONE: Make the existing oracle and matrix reject merge/rebase strategy selection of `ours`, including `git merge -s ours main` and `git rebase --strategy=ours main`, without false-rejecting the supported trace.
  Commit `c1409bff4` adds independent short/long and separate/attached matrix cases for `-s ours`, `-sours`, `--strategy ours`, and `--strategy=ours`; all four now red while `TestMovingTargetConflictOwnershipAndRecovery` remains green.
- DONE: Preserve the cycle-1 coverage for force, side checkout/restore, strategy-option, and forced refspec variants; do not change contract or authority semantics.
  The complete prior mutation matrix stays green, and the cycle correction changes only the existing oracle/matrix test files by +14/-0; no contract, provider, schema, resolver, or workflow bytes changed.
- DONE: Run focused, focused-race, detached adversarial, formatting, full, and race evidence and report the clean exact head and surface; leave the unrelated pilot manifest untouched.
  Focused and focused-race pass; detached exact head `c1409bff47a1313deb07ce93b62b2f61e965a651` passes supported/adversarial tests; `gofmt` and `git diff --check` pass. Candidate surface is 6 files/+277/-4 total. Full and race on immutable state `461605d09` pass every package except the unchanged nine-path `internal/gates` pilot-manifest failure; no race is reported.

### Summary

The oracle now rejects every authorized `ours` merge-strategy spelling while retaining the supported trace and all earlier adversarial coverage. The candidate is clean at `c1409bff4`; this cycle required no contract or authority change.

## Review-finding disposition (cycle 3)

- Reviewer observation: detached audit `/tmp/spacedock-g3-audit-cycle3.eamjJ9` at exact candidate `c1409bff47a1313deb07ce93b62b2f61e965a651` added only a throwaway test; every cycle-1/cycle-2 unquoted form rejects and the supported matrix remains accepted, but `git rebase -X 'theirs' main`, `git checkout '--ours' contract.txt`, and `git push origin '+HEAD:refs/heads/main'` pass the whitespace-token oracle.
- Released user and normal workflow: ordinary shell quoting of Git option and refspec arguments is valid on the supported conflict path.
- Observable harm: the oracle stays green for auto-resolution, side discard, and forced-update commands solely because quote bytes remain in the trace string, so AC-5 can appear proven while observable behavior is unsafe.
- Authority: `contract[skills/first-officer/references/first-officer-shared-core.md#haltrebase-conflictpaths-abort-surface-stop]` requires never force, auto-resolve, or discard either side.
- Trigger evidence: `go test ./internal/ensigncycle -run '^TestCycle3Detached(ForbiddenMatrix|SupportedMatrix)$' -count=1` failed exactly the three quoted forbidden commands; supported `rebase --abort`, non-force push, and non-resolving merge were not false-rejected.
- Worker proposal: Material / evidence defect / G3-owned; the repeated lexical escapes show the free-form string mechanism needs a Captain-owned design decision rather than another automatic spelling patch.
- First Officer disposition: `ROUTE FOR DECISION`; no automatic fix, candidate mutation, or aggregate rerun at the cycle-3 escalation boundary.
- Decision boundary: either preserve semantic AC-5 by replacing/normalizing the free-form command string at a structured argv or canonical shell-token observation boundary, or have the Captain explicitly narrow accepted AC-5 coverage to literal spellings; validation recommends the structured observation reset.
- Validator recommendation: `REJECTED`; candidate stayed clean and unchanged at `c1409bff47a1313deb07ce93b62b2f61e965a651`.

## Stage Report: validation (cycle 3)

- FAILED: Re-review AC-1 through AC-5 at exact corrected head c1409bff47a1313deb07ce93b62b2f61e965a651 and confirm the supported conflict/abort/authority/new-head/unrelated-work trace remains green.
  The supported detached matrix remains accepted and cycle 3 changes no AC-1 through AC-4 contract surface, but AC-5 remains invalid because quoted shell equivalents bypass the oracle.
- FAILED: Re-run the complete detached forbidden-action matrix, including all cycle-1 and cycle-2 spellings and supported-trace false-positive checks; report any remaining semantic escape before candidate mutation.
  All prior unquoted force, refspec, side-selection, strategy-option, and `ours` strategy variants reject, and three supported commands remain accepted; quoted `theirs`, `--ours`, and forced `+refspec` arguments are the exact remaining semantic escapes.
- SKIPPED: Run focused, formatting, full, and race evidence, classify current-head results, and recommend PASSED or REJECTED; keep the unrelated pilot manifest out of G3 scope.
  The material finding entered the checkpoint before aggregate reruns; `ROUTE FOR DECISION` forbids further rerun or mutation. Recommend REJECTED, with the unrelated pilot manifest unchanged and excluded from G3 classification.

### Summary

Cycle 3 preserves the supported trace and closes every previously named unquoted spelling, but free-form whitespace tokenization cannot establish semantic command safety across normal shell quoting. Escalate for a Captain decision between a structured/canonical command observation reset and an explicit AC scope change; do not continue automatic lexical patches.
