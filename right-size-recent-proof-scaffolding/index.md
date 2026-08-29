---
title: Right-size recent proof scaffolding
status: implementation
source: "Captain review of PRs #767, #768, #776, #777, #780, and #781 on 2026-08-28."
started: 2026-08-28T22:07:02Z
completed:
verdict:
score: 0.98
worktree: .worktrees/spacedock-ensign-right-size-recent-proof-scaffolding
issue:
pr:
mod-block:
milestone: 0.28.0
id: 3nm832m6pcnm8008n3wt7h9s
gates:
    version: 1
    records:
        - id: gate:3nm832m6pcnm8008n3wt7h9s:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:3nm832m6pcnm8008n3wt7h9s-backlog-1
              briefing:
                id: briefing:3nm832m6pcnm8008n3wt7h9s:backlog:attempt-1:revision-1
                digest: sha256:110006f02ab3dbf144f9b908a67ff15f6eb208135613e3b6556d17ffb457c998
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:3nm832m6pcnm8008n3wt7h9s:backlog:1
                briefing: briefing:3nm832m6pcnm8008n3wt7h9s:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-28T22:06:23.935141Z"
                decision: approve
                reason: Captain approved the prepared backlog Briefing in Subspace.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:3nm832m6pcnm8008n3wt7h9s:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:3nm832m6pcnm8008n3wt7h9s-ideation-1
              briefing:
                id: briefing:3nm832m6pcnm8008n3wt7h9s:ideation:attempt-1:revision-1
                digest: sha256:36b334984f1673428edbf0730e97ec92cc1fdc21137fade8d50d62d14f800941
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:3nm832m6pcnm8008n3wt7h9s:ideation:1
                briefing: briefing:3nm832m6pcnm8008n3wt7h9s:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-28T22:55:13.742015Z"
                decision: approve
                reason: Captain approved the prepared ideation Briefing in Subspace.
              application:
                target-stage: implementation
                state: consumed
---

Remove redundant tests and fake proof from the current 0.28 stack. Keep one authoritative observation for each distinct failure mode.

This task is the cleanup layer above `spacedock-ensign/flat-entity-gate-room-durability` at `37f50588aa37cfa571a88e4aa87b8f5c8f1b39e8`. The implementation worktree must start from that commit. The adjacent `retire-non-authoritative-test-oracles` task removes older dead oracles and stays separate.

## Problem

Recent fixes added more proof code than product code. Some checks inspect prose or source text. Some repeat the same invariant across package, command, fixture, and live layers. PR #781 also installs the current checkout with a stable-looking version in every common live scenario. That run does not prove the published stable package.

The result is slow CI, large fixtures, and false confidence. At stack parent `37f50588aa37cfa571a88e4aa87b8f5c8f1b39e8`, #776 alone contributes three dedicated test files totaling 301 lines, plus 34 lines of source-reading checks and a second Pi live registration. #781 makes every common Claude or Codex journey install a copied checkout: 17 installs for the common journey set, 19 in the full Claude set (17 common plus two recovery runners), and 18 in the full Codex set (17 common plus the wait-matrix runner).

A validator can prove the same distinct failures with existing owners, one smaller Pi journey, and one-off install evidence whose source is named truthfully.

During implementation, the first Pi validation reported a material finding: the child committed the work marker but omitted the required Stage Report. The first officer classified the supported normal workflow as a fresh Pi initial dispatch with a work-only checklist, the observable harm as an incomplete completion artifact, and the authority as `value-ac[AC-3]`; the captain then ruled “ok do it” and authorized the smallest product correction. Root-cause investigation before mutation found that this run had selected a stale main-root `spacedock` binary whose dispatch artifact predated the stack parent's Pi first-action correction. A diagnostic run bound to the current worktree binary loaded `skills/ensign/SKILL.md`, `ensign-shared-core.md`, and the Pi adapter, then wrote and committed the exact Stage Report. The correction is therefore already present in the approved stack parent and must be preserved by this cleanup; duplicating the report protocol in generated dispatch text would recreate the scaffolding this task removes.

## Required outcome

Delete or combine the redundant proof added by PRs #767, #768, #776, #777, #780, and #781. Preserve one primary proof owner for each distinct product failure mode.

Add a short workflow rule for test selection. The rule must not add a lint, gate, CI lane, required table, or recurring ceremony.

Implementation starts from stack parent `37f50588aa37cfa571a88e4aa87b8f5c8f1b39e8`. It changes proof selection and validation guidance while preserving the stack parent's smallest Pi behavior correction: a fresh worker is explicitly directed to load the ensign discipline and its shared Stage Report protocol before reading a work-only assignment. No additional product mechanism is justified unless the retained owner fails with the current worktree binary.

## Proposed approach and proof-owner audit

Use one committed behavioral owner for each failure mode. A second layer remains only when its falsifying edit is distinct from the primary owner's. The implementation makes no new test file, fixture, scenario, registry row, lane, lint, gate, or recurring validation step.

| Source | Distinct failure mode | Primary owner and action |
| --- | --- | --- |
| #767 | Completion guard collapses missing report, malformed checklist, untracked entity, or dirty entity into one diagnosis, or mutates bytes on refusal. | Retain `TestEnteredWorktreeStageFailureDiagnosticsAreDistinctAndByteClean` as the four-row classification owner. Retain the near-miss and cycle-suffix rows in `TestEnteredStageProjectionRequiresCommittedCompleteReport` for heading-token parsing. |
| #767 | Consumed nonterminal approval cannot complete through merge guard. | Retain the existing `TestTerminalDeliveryFailureReworkRoundTrip`, which already drives merge guard through the terminal journey. Delete `TestConsumedNonterminalApprovalAllowsMergeGuard`; it repeats that authority path. The separate `TestConsumedNonterminalApprovalAllowsOrdinaryTerminalFields` remains the direct `status --set` control. |
| #768 | Absolute, state-relative, and launch-relative artifact/reference forms do not bind to the same immutable source. | Retain `TestPrepareSelectedSourceFormsShareImmutableBinding` at the resolver layer. Retain `TestResolveSelectedSourceRefusesAmbiguityAbsenceAndProbeErrors` for ambiguous, absent, lexical-alias, and probe-error grammar. |
| #768 | Directory, symlink, unreadable, or foreign-root selected sources mutate state or cross trust roots. | Retain `TestPrepareSelectedSourceGuardsRemainByteClean`; these are distinct resolver refusals. |
| #768 | The CLI drops `LaunchDir` before calling the resolver. | Fold one launch-relative artifact case into `TestGatePrepareCLIPassesStateRelativeArtifactWithoutCwdJoin`. Delete the two-flag setup duplicate `TestGatePrepareCLIResolvesLaunchRelativeSelectedSources`; artifact/reference parity stays owned below the CLI. |
| #776 | The Pi dispatch artifact loses its host-specific first-action shape. | Retain `TestBuildPiHostPromptShape` as the deterministic artifact-output owner. Remove the duplicate `/skill:ensign` token assertion from generic JSON ergonomics coverage and delete `TestPiFirstActionInvokesEnsignSkill`, which reasserts the generated prompt and then claims worker behavior. |
| #776 | A fresh Pi worker fails to load ensign, write a structurally valid report, or commit the entity cleanly when its checklist does not self-describe the report format. | Modify existing `TestLivePiFrontDoorSmoke` so its checklist names only the work, then keep its child-transcript `assertPiEnsignBootContract`, report-shape, git-log, and clean-tree observations. Delete `TestLivePiNonSelfDescribingDispatch` and its two private fixture files; remove its CI/registry entry. |
| #776, adjacent | Source or prompt text is treated as proof that the live journey retained its graders or that the agent followed instructions. | Delete `TestPiLiveSmokePromptRequiresExactStageReportHeading` and `TestLivePiFrontDoorSmokeRetainsAllGraders`. Their claimed outcomes are observed by the retained live journey. |
| #777 | Any supported resume alias leaks the bootstrap, or an ordinary launch loses it. | Retain one table-driven `TestPiResumeSuppressesBootstrapPrompt` with `--resume`, `--resume=<id>`, `-r`, `--continue`, `-c`, and one non-resume control. Share one fixture/runner setup across rows and delete the second non-resume spelling. |
| #780 | Flat and folder entities encode/resolve canonical review refs differently, or malformed reserved refs spend authority. | Retain `TestReviewRoomRefsHaveOneMeaningAcrossEntityForms` as the resolver owner and `TestMalformedCanonicalRoomRefCannotSpendGateAuthority` as the one public authority-boundary propagation case. |
| #780 | A flat entity's round is not durably committed for a fresh host. | Retain `TestStateCommitMakesFlatRoundDurableInFreshHost`; no lower-level round test can catch the state-commit/fresh-clone failure. |
| #780 | Flat publication escapes the shared review home, or replay follows an artifact outside the allowed boundary. | Reduce `TestRoundPublishesFlatAndFolderThroughSharedReviewHome` to one flat publication case. Reduce `TestRoundReplayRefusesSiblingAndMutableEntityArtifacts` to one flat sibling-target replay refusal. Delete `TestRoundArtifactBoundaryIsByteCleanForFlatAndFolder`, the folder rows, the mutable-entity duplicates, and the declared-folder-policy subcase; existing `TestRoundRecordNeutralReplayAndRefusalsAreByteClean` owns folder publication/exact replay and the preflight matrix owns mutable-entity publication. |
| #780 | A historical `./review/...` pointer is rewritten or stops replaying. | Retain `TestRoundFrozenLegacyFolderPointerReplaysUnchanged`. Retain `TestValidateResolvesCanonicalFlatRoomWithoutConversionWarning` for the distinct status-warning surface. |
| #781 | Host inventory schemas parse incorrectly. | Retain `TestParsePluginInventoryLiveSchemas`. |
| #781 | Doctor misses an enabled sibling, treats a disabled sibling as conflict, or mutates installation. | Retain `TestDoctorSiblingInventory`. |
| #781 | Claude or Codex frontdoors launch with conflicting sibling scope, fail to heal when allowed, attempt installation under `--no-install`, or ignore inventory failure. | Retain `TestStable0271FrontDoorsHealEnabledSiblingBeforeLaunch`; its host table and refusal controls exercise product behavior without a package installation. |
| #781 | A copied checkout is presented as the published stable package. | Restore common Claude/Codex journeys to checkout-candidate `--plugin-dir` execution and remove per-runner installs. For release provenance, run one staged-candidate install and, only when the stable claim matters, one ordinary install from the real published release channel; record source, observed version, and launch behavior. |

## Risk evidence

The named owners above were audited at the exact stack parent. The decisive overlaps are executable, not thematic: the retained Pi frontdoor already reads the child transcript and checks the report, commit, and clean tree; the existing round test already publishes and exactly replays folder form; the existing terminal round-trip already crosses merge guard; resolver tests already exercise both selected-source flags and path grammars.

#781's temporary marketplace is valid candidate-package proof. Its `stableLiveRelease` helper, however, copies the checkout, reads its manifest version, stamps a newly built binary, and routes installation through that temporary marketplace. No observation in that path reaches the published stable channel.

No product-mechanism spike is needed: every cleanup target relies on an already exercised owner listed above. The workflow-rule effect must be checked during implementation, after the README has an uncommitted edit; the `behavior-diff` skill refuses to run without that change, and the captain explicitly kept README editing out of ideation.

## Out of scope

- Product behavior changes beyond preserving the authorized Pi first-action correction already present at the stack parent.
- New test frameworks, CI lanes, runtime controllers, or fixture protocols.
- The older dead-oracle cleanup owned by `retire-non-authoritative-test-oracles`.
- Release-process redesign.
- Editing `docs/dev/README.md` during ideation; implementation applies the approved diff below.

## Expected surface and tolerance

Estimate net LOC change: -560, across 18 files. Expected insertions: 65. Expected deletions: 625. Tolerance: 100 net lines and 2 files.

Expected files: `docs/dev/README.md`; `.github/workflows/runtime-live-e2e.yml`; `docs/runtime-live-ci-registry.md`; three CLI tests (`gate_test.go`, `pi_frontdoor_test.go`, `terminal_consume_test.go`); two dispatch tests (`build_json_ergonomics_test.go`, deleted `build_stage_report_protocol_test.go`); five Pi/live proof files (`pi_live_runner_test.go`, `pi_live_controls_test.go`, `pi_evidence_grade_impl_test.go`, and deleted `pi_nonself_describing_{build,live}_test.go`); `round_test.go`; and the four shared Claude/Codex live-runner files changed by #781.

Count baselines: relevant Pi journeys fall from two to one. The common journey inventory remains 17 per host; package installs fall from 17 per common host set (19 full Claude, 18 full Codex) to zero committed live-suite installs. Release-source validation uses at most two one-off installs: one local candidate and one real published stable package.

Allowed semantic change: test selection, test setup, live-run count, and validation guidance. Command grammar, stored formats, authority, and supported runtime behavior must not change.

Corrected implementation estimate after the captain's authorization remains +65/-625 = -560 net across 18 files: root-cause evidence shows the required product behavior is already in the stack parent, so the cleanup adds no duplicate production lines or files. The current implementation surface before final validation is +87/-603 = -516 net across 17 files, within both approved tolerances. If a current-binary owner failure requires a new product edit, recompute this surface before mutation and stop if it exceeds either tolerance.

## Exact proposed workflow diff

Implementation applies this diff; ideation does not edit the workflow README.

```diff
diff --git a/docs/dev/README.md b/docs/dev/README.md
--- a/docs/dev/README.md
+++ b/docs/dev/README.md
@@
 - **Evidence must be able to fail.** Each AC's cited evidence names the concrete change that would flip it — the falsifying edit. A criterion whose author cannot name that change does not count. The gate reads the falsifying change, not a pass count.
+- **One proof owner per failure mode.** Reuse or modify an existing behavioral test before adding one. Add another committed check only when a distinct falsifying edit would escape the primary owner; otherwise combine or delete it. Use one-off manual validation for release provenance or external wiring that a committed test cannot reproduce truthfully.
@@
-  - Test plans should state what verifies the implementation, estimated cost/complexity, and whether fixture, CLI, or live workflow tests are needed.
+  - Test plans should name the existing primary proof owner (or explain why none exists), the distinct falsifying edit for each additional check, estimated cost/complexity, and whether deterministic, live, or one-off manual validation is needed.
```

This is guidance inside the existing Proof policy and ideation output. It creates no gate, lint, lane, table, fixture, or recurring ceremony.

## Incident-derived behavior diff

This is a single-role rule: an ideation worker chooses its own proof plan, with no fabricated FO/worker handoff. Implementation must mint a `base` capsule with the skill's `make-capsule.sh` and proceed only after it prints `CAPSULE OK`; the source repository is not a hand-built fixture. The before/after pair differs only in `docs/dev/README.md`.

Exact behavior-diff parameters for captain confirmation before the six-run comparison:

```text
--file docs/dev/README.md
--task "Using the docs/dev workflow's ideation policy, review the proof surface added for the Pi initial-dispatch bootstrap change: a fresh worker now loads the ensign skill before reading its assignment. Propose the validation plan you would approve, including what to run and what committed checks, if any, to keep, change, remove, or add. Do not modify files."
--agent codex --vocab spacedock
```

The task names the incident and decision point but not the desired one-owner outcome. Run three fresh trials per variant (six headless runs), keep the report local, and report the flow divergence without an automatic pass/fail claim. Identical flows mean the task did not expose the decision and require a sharper task, not that the rule passed or failed.

## Acceptance criteria

**AC-1 (VALUE): The stack keeps one primary proof for each distinct failure mode and removes at least 300 net test or harness lines.**
Verified by: compare `git diff --numstat` from stack parent `37f50588aa37cfa571a88e4aa87b8f5c8f1b39e8` and audit the owner map above. A duplicate owner, an owner without a distinct falsifying edit, or a reduction below 300 net lines fails this criterion.

**AC-2: No committed test claims agent or runtime behavior from prose, generated prompt text, or Go source text.**
Verified by: inspect the changed tests and execute their retained behavioral owners. A retained assertion whose expected behavior comes only from instruction/source text fails; deterministic artifact-output shape remains legitimate only where that shape itself is the claim.

**AC-3: A fresh Pi worker given a work-only, non-self-describing checklist produces the required Stage Report, and one live journey observes the real child transcript, ensign skill load, report, and clean commit.**
Verified by: build the current worktree binary and run only modified `TestLivePiFrontDoorSmoke` with `SPACEDOCK_BIN` bound explicitly to that binary. The transcript must show the worker load the ensign skill and shared core before producing the structurally valid report and clean commit. A stale/ambient launcher, missing report, second journey, fixture, registry row, or CI lane for the same behavior fails this criterion.

**AC-4: Package-source claims name the real source. A copied checkout is never labeled as published stable.**
Verified by: confirm common live runners use checkout-candidate overrides with zero installs. Run one candidate frontdoor install from the staged marketplace and inspect its source. When published stable is relevant, manually install with no marketplace override from the real release channel and record source, version, and launch behavior. A temporary marketplace presented as published stable fails this criterion.

**AC-5: Resolver and round coverage remains at its owning layer, with at most one public-boundary propagation check for each invariant.**
Verified by: run focused resolver, `gate prepare`, round publication, and replay tests. Repeating one invariant across cross-product forms without a distinct failure mode fails this criterion.

**AC-6: The workflow guides authors toward one primary proof owner, and the rule has before/after behavioral evidence without new enforcement machinery.**
Verified by: apply the exact README diff uncommitted, run the six-trial incident-derived `behavior-diff` comparison after captain confirmation, and present its flow report for judgment. A new lint, gate, CI lane, required table, fixture, or recurring validation step fails this criterion.

**AC-7: The full deterministic suite remains green, and required live evidence is smaller than the pre-cleanup stack.**
Verified by: run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`. Run only modified `TestLivePiFrontDoorSmoke`; record the 2-to-1 relevant Pi journey change, unchanged 17 common journeys, zero committed live installs, and any one-off candidate/stable installs.

## Test plan

The selection rule is concrete: use a committed deterministic test for repeatable product behavior, at the lowest layer that owns the failure; reuse or modify the nearest owner first. Add another committed check only when a named falsifying edit escapes the primary owner. Use one existing live journey only when agent/runtime behavior cannot be observed deterministically. Use manual validation when truth depends on a mutable external source, especially release-channel provenance or one-time installation wiring.

Implementation sequence: perform the exact deletions/combinations in the owner audit; restore checkout-candidate common live runners; apply the approved README diff; run the behavior-diff comparison before committing that rule; format and run focused plus full deterministic suites; then build the current worktree launcher and run the one Pi frontdoor journey with `SPACEDOCK_BIN` explicitly bound to it, followed by the one-off package installations justified by AC-4. Inspect the retained journey's child transcript for the ensign skill and shared-core reads as well as its report, commit, and clean-tree graders. Treat a failure from an ambient or stale launcher as invalid evidence and rerun with the bound candidate before proposing mutation.

No audit row exposes an uncovered failure mode. Therefore implementation adds no test file, live scenario, fixture, registry row, CI lane, lint, gate, or harness code.

### Feedback Cycles

## Stage Report: ideation

- DONE: Audit the named recent tests and assign one primary proof owner to each distinct failure mode; identify exact checks to delete, combine, or retain.
  The owner table audits #767, #768, #776, #777, #780, and #781 at stack parent `37f50588`, with named retained, combined, and deleted checks.
- DONE: Refine the net LOC, file-count, live-journey, and package-install estimates against stack parent 37f50588aa37cfa571a88e4aa87b8f5c8f1b39e8.
  Estimate is +65/-625 = -560 net across 18 files; relevant Pi journeys fall 2-to-1 and committed install counts fall from 19 Claude/18 Codex full-suite calls to zero.
- DONE: Produce a concrete test-versus-manual validation rule and implementation plan with no new lint, gate, lane, fixture, scenario, or recurring ceremony.
  The plan selects the lowest behavioral owner, reserves live proof for runtime-only claims, uses manual release provenance, and includes the exact README and neutral six-trial behavior-diff design.

### Summary

Ideation reduces the recent proof stack by deleting duplicate live, source-reading, cross-product, and install scaffolding while preserving a named owner for every distinct failure. It adds only concise workflow guidance, with its effect measured by an incident-derived `behavior-diff` run during implementation rather than by new enforcement machinery.

## Stage Report: implementation

- DONE: Apply the approved proof-owner audit: delete or combine the named duplicate checks, retain every distinct behavioral owner, and change no product behavior.
  Commit `de61c1baf` removes 603 lines and adds 87 across 17 files; focused status, gate, dispatch, round, terminal, runner, and Pi owners passed, and deleting any retained classification, trust-boundary, fresh-clone, or launch-argv behavior would fail its named owner.
- DONE: Apply the exact docs/dev/README.md rule uncommitted, run Engram behavior-diff with three trials per variant, and commit the rule only after reporting a real flow divergence or an honest inconclusive result.
  Base capsule printed `CAPSULE OK`; six valid REVIEW trials produced an honest inconclusive result because repo-shell reads were denied, at `/private/tmp/spacedock-engram-behavior-diff-3nm832/runs/diff-20260828-164836/report.md`, before the README entered `de61c1baf`.
- DONE: Verify focused owners plus gofmt, go test ./..., go test ./... -race, one targeted Pi frontdoor journey, package-source observations, and actual LOC/files/install counts against the -560/18 estimate.
  Gofmt/diff checks and focused owners passed; isolated fresh-CODEX_HOME full and race suites passed, while the first ambient runs failed only on a stale local Codex plugin cache that the isolated reruns removed.

### Validation evidence

- The single retained `TestLivePiFrontDoorSmoke` passed in 198.29s with the current launcher and declared pi-subagents 0.53.0: one child, five reads, ensign skill at read #2, shared-core protocol, valid report, required commit, and clean state. Removing the first-action skill load or report protocol makes it fail.
- The first Pi attempt found a marker-only commit and no Stage Report. The FO recorded normal fresh-Pi/work-only use, incomplete-artifact harm, and `value-ac[AC-3]`; the captain authorized FIX with “ok do it.”
- Root-cause investigation showed that attempt used the ambient stale main-root binary and its pre-correction dispatch. A current-binary diagnostic produced the report; the remaining metadata-location failure came from local pi-subagents 0.35.1, and the declared 0.53.0 rerun passed. No duplicate product or harness edit was justified.
- Candidate install evidence used an isolated CODEX_HOME and staged local marketplace: one enabled `spacedock@spacedock-local` 0.28.0-pre1, doctor-compatible with the candidate binary.
- Published evidence used ordinary stable installation with no override: marketplace `https://github.com/spacedock-dev/marketplace.git`, plugin `https://github.com/spacedock-dev/spacedock.git` ref `stable`, version 0.27.2, compatible doctor, and successful `spacedock codex --no-install -- --version` launch.
- Actual surface is +87/-603 = -516 net across 17 files, within -560/18 tolerance; relevant Pi journeys are 2-to-1, common journeys stay 17 per host, committed suite installs are zero, and validation used one candidate plus one published-stable install.

### Summary

Implementation keeps one behavioral owner per distinct failure, removes duplicate/source-reading/live-install proof, and adds the approved concise workflow rule. It preserves the stack parent's authorized Pi first-action correction and proves a fresh worker can derive the Stage Report from a work-only checklist without adding another test, scenario, fixture, registry row, lane, lint, gate, table, harness, or product mechanism.
