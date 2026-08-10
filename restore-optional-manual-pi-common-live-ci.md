---
title: Restore optional manual Pi common-live CI
status: validation
source: "Captain correction, 2026-08-09: PR #639 removed pi-live when it removed Pi from pull-request approvals. Restore manual CI Pi evidence without making Pi a merge requirement."
started: 2026-08-09T15:42:25Z
completed:
verdict:
score: 0.9
sprint: test-behavior-completeness
sprint-readiness: ready
group: live-ci-evidence
worktree: .worktrees/spacedock-ensign-restore-optional-manual-pi-common-live-ci
issue:
pr:
mod-block:
id: 0aqnm6v8ajns6cpsknxn9wf2
gates:
    version: 1
    records:
        - id: gate:0aqnm6v8ajns6cpsknxn9wf2:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:0aqnm6v8ajns6cpsknxn9wf2-backlog-1
              briefing:
                id: briefing:0aqnm6v8ajns6cpsknxn9wf2:backlog:attempt-1:revision-1
                digest: sha256:1ac393459063bbb598d4631c0ddda9acfdf0de48f2535f2b600fe53bd63263c8
                request-digest: sha256:cb99842a502664750d83cd144f580afc457ef1d73fd628ddb9540acc2c66ee06
                room-ref: ./restore-optional-manual-pi-common-live-ci/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:0aqnm6v8ajns6cpsknxn9wf2:backlog:1
                briefing: briefing:0aqnm6v8ajns6cpsknxn9wf2:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-09T15:42:06.443214Z"
                decision: approve
                reason: Captain directed dispatch; the task restores optional retained Pi CI evidence without changing pull-request requirements or blocking pre3.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:0aqnm6v8ajns6cpsknxn9wf2:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:0aqnm6v8ajns6cpsknxn9wf2-ideation-1
              briefing:
                id: briefing:0aqnm6v8ajns6cpsknxn9wf2:ideation:attempt-1:revision-1
                digest: sha256:38cfe1ae900ba2cbaede656e65dafa1c54b3a88355c2fb3d93a2cd52d6419e75
                request-digest: sha256:5ac2196d66b19f30d73abfbba974fb0d938a86a5559b3d90527ac32c9b45c205
                room-ref: ./restore-optional-manual-pi-common-live-ci/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:0aqnm6v8ajns6cpsknxn9wf2:ideation:1
                briefing: briefing:0aqnm6v8ajns6cpsknxn9wf2:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T15:52:25.300232Z"
                decision: approve
                reason: The direction provides retained Pi Luna/max CI evidence and makes manual Opus and Pi cadences runtime-exclusive, using existing workflow conditions and job-graph tests without changing pull-request requirements.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:0aqnm6v8ajns6cpsknxn9wf2:validation
          stage: validation
          attempts:
            - id: gate-attempt:0aqnm6v8ajns6cpsknxn9wf2-validation-1
              briefing:
                id: briefing:0aqnm6v8ajns6cpsknxn9wf2:validation:attempt-1:revision-1
                digest: sha256:726967ec584dcb98ea9063dd63cfd9c2292daa7a46086e686270d364a5d2e5f8
                request-digest: sha256:123ff303ed7debbd1797108104478d97390968d1ce57c336fd22dd2655cd9619
                room-ref: ./restore-optional-manual-pi-common-live-ci/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:0aqnm6v8ajns6cpsknxn9wf2:validation:1
                briefing: briefing:0aqnm6v8ajns6cpsknxn9wf2:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T21:53:57.842218Z"
                decision: revise
                reason: 'Exact-candidate run 31323946889 retained no journey-metrics files, failing AC-1. Captain authorized a bounded design reset: preserve AC-1 and the existing metrics format; add at most two files and +120 net lines; reuse journeymetrics APIs; prove the current no-op failure and exact repaired Pi artifact.'
            - id: gate-attempt:0aqnm6v8ajns6cpsknxn9wf2-validation-2
              briefing:
                id: briefing:0aqnm6v8ajns6cpsknxn9wf2:validation:attempt-2:revision-1
                digest: sha256:a4e0d87d65c4177acb055ce97e7e9455da2b9447b93913e84cb1ea34631af5ed
                request-digest: sha256:dd0c93da792868c7556b720df2015c9d77baceef2cc7cb786de5aa326a3b6525
                room-ref: ./restore-optional-manual-pi-common-live-ci/review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:0aqnm6v8ajns6cpsknxn9wf2:validation:2
                briefing: briefing:0aqnm6v8ajns6cpsknxn9wf2:validation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-09T23:49:20.318112Z"
                decision: approve
                reason: 'Approved under sprint conn: candidate 7f5f79aad passed local/offline validation, exact Pi retained-metrics and exclusivity evidence, and exact Opus exclusivity evidence. AC-2 remains an explicit merge-boundary blocker until the normal PR run passes Sonnet and Codex and skips Pi with zero steps.'
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:0aqnm6v8ajns6cpsknxn9wf2-validation-3
              briefing:
                id: briefing:0aqnm6v8ajns6cpsknxn9wf2:validation:attempt-3:revision-1
                digest: sha256:ece7aa0d2114c1f27bb75d2ef1977d2a4ff598a041e2037b3613d2f720558a0a
                request-digest: sha256:8aa95e57461fe487561f2bd493a0e18f88e5f04b1f3ea63f4965f2b161f37083
                room-ref: ./restore-optional-manual-pi-common-live-ci/review/validation/briefing-3
---

Give maintainers one optional GitHub Actions command that runs the Pi common journeys and retains their evidence.

## Problem

PR #639 correctly removed Pi from pull-request approvals. It also deleted the `pi-live` job and replaced manual CI with a local-only command.

The local command uses subscription authentication, but it does not give the team a reproducible CI environment or retained artifacts.

## Required outcome

Add `pi` to the existing manual `live_cadence` input in `runtime-live-e2e.yml`.

For this choice, run the offline job and one approval-gated `pi-live` job. Do not create Claude or Codex jobs.

Keep pull-request behavior unchanged. Pull requests create only the Sonnet and Codex live jobs, so Pi is not a merge requirement.

Make the release cadences exclusive. Manual `live_cadence=opus-pre-release` runs offline plus Opus/max only. Manual `live_cadence=pi` runs offline plus Pi Luna/max only. Manual `live_cadence=sonnet` may preserve its current Sonnet-plus-Codex behavior.

Use the existing `CI-E2E-PI` environment and the pinned Pi packages. Run Pi with the stored `OPENAI_API_KEY`, model `openai/gpt-5.6-luna`, and maximum thinking. Run the canonical `^TestLiveCommon` selector and the Pi front-door substrate proof. Retain logs, artifacts, and journey metrics.

Do not add a second workflow, scheduler, registry, simulator, or product-behavior repair.

GitHub represents every statically declared job in the run graph. A false job-level condition can therefore appear as a skipped node. In this task, "do not create" means that the excluded host gets no runner, environment approval request, secret access, or executed step. A skipped graph node is acceptable; an executed or approval-waiting excluded-host job is not.

## Proposed approach

1. Extend the existing `live_cadence` choice from `sonnet|opus-pre-release` to `sonnet|opus-pre-release|pi`. Keep `sonnet` as the default.
2. Add job-level event conditions. Pull requests run Claude Sonnet/max plus Codex Luna/max. Manual `sonnet` may run Sonnet plus Codex as it does now. Manual `opus-pre-release` runs Opus/max and excludes Codex and Pi. Manual `pi` runs Pi Luna/max and excludes Claude and Codex. The offline job remains unconditional and is the only prerequisite of each selected live job.
3. Restore the previously proven Pi job in the same workflow instead of inventing a new runner. Keep its verified package integrity pins, current-checkout binary, Pi doctor check, isolated live-test environment, clean gotestsum logs, and `CI-E2E-PI` environment.
4. Set the Pi job's model to `openai/gpt-5.6-luna:max`. Pi treats the suffix as the `max` thinking level. Both common journeys and the front-door proof already consume the shared Pi model selector.
5. Add the Pi journey-metrics directory and upload it with the common-journey detail JSONL, front-door detail JSONL, current-checkout binary, Pi diagnostics, and Pi live artifacts. Do not add Pi to the pull-request journey-delta comment, because the Pi cadence is manual-only.
6. Change the workflow guards from "Pi must be absent" to the event matrix above. Parse the workflow structure and executable commands; do not prove routing with comments or prose search.

The Pi job is long because package integrity and runtime isolation are part of the retained evidence. The simplest alternative was a local command only. It is insufficient because it supplies neither a reproducible runner nor retained artifacts. A second workflow could avoid skipped nodes but would duplicate the security and offline-gate contract, so it is outside scope.

## Expected surface

- `.github/workflows/runtime-live-e2e.yml`: restore one conditional Pi job and add cadence routing, model, metrics, and artifact settings; approximately 230-260 insertions and fewer than 15 deletions.
- `internal/release/runtime_live_evidence_workflow_test.go`: replace the Pi-absence oracle with parsed event/cadence/job-routing and named-evidence assertions; approximately 55-85 changed lines.
- `internal/release/workflow_exec_guard_test.go`: require Pi metrics production and upload only on the manual Pi job; approximately 15-30 changed lines.
- `internal/release/journey_workflow_test.go`: replace the paid-Pi rejection with manual-only Pi producer and artifact assertions; approximately 20-40 changed lines.
- `docs/runtime-live-ci.md`: document the manual Pi cadence, Luna/max setting, retained artifacts, and unchanged pull-request lanes; approximately 8-16 changed lines.
- `docs/site/contributing/architecture-notes.md`: replace the local-only Pi CI statement with the manual cadence boundary; approximately 3-6 changed lines.

Tolerance: at most 7 files and 380 inserted lines. A small helper-test edit in `internal/release` is allowed if the existing YAML parser cannot expose job conditions. No production Go package, live journey implementation, fixture, registry, skill, or release workflow may change.

Observable semantics allowed to change: the manual workflow input accepts `pi`; manual Opus and Pi dispatches become runtime-exclusive; the Pi CI model/thinking level and retained artifacts change. Command grammar, product runtime behavior, workflow-state formats, gate authority, pull-request required checks, manual Sonnet behavior, and the 17-journey registry do not change.

## Mechanism proof

The restored runtime mechanism and package pins were exercised by the former `pi-live` job before its cadence removal. A local Pi 0.82.1 probe on 2026-08-09 confirmed that `--model` accepts provider-qualified models, `--thinking max` is supported, and an available `OPENAI_API_KEY` exposes `openai/gpt-5.6-luna`. Inspection of the pinned 0.80.10 package catalog confirmed the same direct-OpenAI model, `max` thinking mapping, and exact existing package integrity.

The remaining unverified mechanism is the new GitHub event routing as a whole. Implementation must first make the structural workflow test fail on the current two-choice/no-Pi workflow, then pass it after the change. The final proof is one real exact-cadence dispatch; no simulator can replace it.

## Acceptance criteria

**AC-1 (VALUE) - A maintainer can obtain retained Pi CI evidence on demand.**
Verified by: one manual `live_cadence=pi` run passes the offline and Pi jobs, records `openai/gpt-5.6-luna` with maximum thinking, and uploads common-journey JSONL, substrate JSONL, Pi diagnostics, and journey metrics.

**AC-2 - Pi does not become a pull-request merge requirement.**
Verified by: one pull-request workflow run executes Sonnet and Codex live jobs, allocates no Pi runner, requests no `CI-E2E-PI` approval, and leaves any statically declared Pi node skipped.

**AC-3 - A Pi-only dispatch spends no Claude or Codex lane.**
Verified by: the manual Pi run allocates no Claude or Codex runner, requests only the `CI-E2E-PI` approval, and leaves any statically declared Claude or Codex nodes skipped.

**AC-4 (VALUE) - Each manual release cadence spends only its selected runtime.**
Verified by: Actions job and deployment records for exact-candidate `opus-pre-release` and `pi` dispatches show offline plus only Opus or Pi respectively; excluded hosts allocate no runner, request no environment approval, receive no secret, and execute no step. A skipped static graph node is permitted. Manual `sonnet` is not a release cadence and may retain Sonnet plus Codex.

**AC-5 - The desired journey registry stays unchanged.**
Verified by: registry reconciliation passes with the same current 17 common journeys, fixtures, TODO owners, and canonical selector.

## Test plan

Start with a failing structural test for the event matrix: pull request selects Sonnet plus Codex; manual `sonnet` may select Sonnet plus Codex; manual `opus-pre-release` selects Opus and excludes Codex plus Pi; manual `pi` selects Pi and excludes Claude plus Codex. Assert the Pi job uses `CI-E2E-PI`, the stored OpenAI key, `openai/gpt-5.6-luna:max`, both exact test selectors, and retained metrics/artifacts. Assert the Opus job uses only `CI-E2E-OPUS`. Use the existing parsed-YAML and executable-step inspection helpers, not prose grep, a new parser, or a simulator for GitHub Actions.

Run `go test ./...`, `go test ./... -race`, formatting, and the registry reconciliation. Dispatch one real `live_cadence=pi` workflow on the exact candidate commit, approve only `CI-E2E-PI`, and inspect the Actions job list, requested environment, model record, and downloaded artifact contents. For Opus exclusivity, use the next exact-candidate `opus-pre-release` run and verify that it requests only `CI-E2E-OPUS`; if no release run exists for that commit, dispatch one. Inspect one pull-request run on the same workflow shape to confirm that the existing Sonnet and Codex lanes still run and Pi consumes no runner or approval.

## Documentation diff proposed at ideation

In `docs/runtime-live-ci.md`, replace:

> Pi live evidence runs locally with `pi login`. The offline job keeps the registry reconciliation. An explicit local API key is also supported.

with:

> An explicit `live_cadence=pi` dispatch runs the 17 common Pi journeys and the Pi front-door proof with `openai/gpt-5.6-luna` at maximum thinking. It waits only for `CI-E2E-PI` approval and retains Pi logs, diagnostics, journey metrics, and session artifacts. Pull requests still run only Sonnet and Codex; Pi is optional and is not a merge requirement. Local Pi execution remains supported with `pi login` or an API key.

Also change the Opus cadence sentence to:

> An explicit `live_cadence=opus-pre-release` dispatch runs offline plus `claude-opus-4-8` at maximum effort. It allocates no Codex or Pi runner and requests only `CI-E2E-OPUS` approval.

In `docs/site/contributing/architecture-notes.md`, replace:

> **Pi evidence** uses the local subscription path with `pi login`. Pull requests keep only the free registry reconciliation.

with:

> **Manual release evidence** is runtime-exclusive: `opus-pre-release` runs only Opus/max behind `CI-E2E-OPUS`, and `pi` runs only Pi Luna/max behind `CI-E2E-PI` while retaining the common-journey and front-door artifacts. Pull requests continue to run Sonnet/max and Codex Luna/max; they do not run Pi.

## Stage Report: ideation

- DONE: Design one manual Pi-only cadence in the existing workflow that uses the stored OpenAI key and retained evidence.
  The design adds one `pi` choice, restores the proven conditional job, selects direct-OpenAI Luna/max, and retains detail, diagnostic, session, and journey-metric artifacts.
- DONE: Keep pull-request Sonnet and Codex behavior unchanged, with no Pi merge requirement or unrelated lane spend.
  Parsed event conditions preserve both PR lanes; Pi manual dispatch requests only `CI-E2E-PI`, with excluded static nodes permitted only as skipped records.
- DONE: Specify the smallest workflow and test changes plus one real exact-cadence validation run.
  The expected surface reuses the deleted Pi job, changes three focused workflow guards and two docs, and requires a real exact-commit `live_cadence=pi` run.

### Summary

The design restores retained Pi evidence without restoring Pi as a pull-request gate. It also reconciles the task with D8: the unchanged desired registry now contains 17 common journeys, not 16. The only remaining mechanism proof is the real GitHub cadence run after implementation.

## Stage Report: ideation

- DONE: Design one manual Pi-only cadence in the existing workflow that uses the stored OpenAI key and retained evidence.
  The original Pi design remains; Pi uses direct-OpenAI Luna/max and retains common, substrate, diagnostic, session, and journey-metric evidence.
- DONE: Keep pull-request Sonnet and Codex behavior unchanged, with no Pi merge requirement or unrelated lane spend.
  The corrected matrix keeps PR Sonnet/max plus Codex Luna/max, makes Opus manual dispatch Opus-only, and makes Pi manual dispatch Pi-only.
- DONE: Specify the smallest workflow and test changes plus one real exact-cadence validation run.
  Existing job conditions and YAML test helpers own exclusivity; exact-candidate Pi and Opus Actions records prove that excluded runtimes consume no runner, approval, secret, or step.

### Summary

The Captain-approved fold makes both manual release cadences exclusive. `opus-pre-release` spends only Opus, `pi` spends only Pi, and pull requests retain Sonnet plus Codex; no new workflow, parser, simulator, registry, or release mechanism is introduced.

## Implementation design reset

Captain-conn approval on 2026-08-09 raises the file tolerance from 7 to 8 without changing semantics or LOC. The two required surfaces missing from the estimate are the existing live-registry reconciliation guard and clean-log/gotestsum guard; both must change because they currently encode the deliberate absence of Pi CI. The authorized eight files are the workflow, four existing workflow/registry test files, the clean-log test, and the two approved runtime CI documents. No workflow, parser, simulator, registry entry, live journey, fixture, release mechanism, or other product surface is added.

## Stage Report: implementation

- DONE: Deliver the approved eight-file manual Pi cadence candidate within the +380 net-line limit, preserving pull-request Sonnet and Codex behavior and exclusive manual Opus/Pi routing.
  Candidate `e838fba693d5562fa4cee99cfbfb682014331a4f` changes eight approved files and adds 244 net lines.
- DONE: Prove the candidate with the required offline, race, formatting, registry, active-owner, and exact workflow-structure checks; record any unavailable live evidence honestly.
  `gofmt`, both full suites, registry reconciliation, active-owner proof, and focused workflow guards passed.
  The relative active-owner command did not find the split-root checkout. The same test passed with its assigned absolute state path.
  The local Pi proof stopped before execution. The `openai-codex` OAuth refresh returned `401 refresh_token_reused`.
  No protected environment received an approval request. Exact Pi, Opus, and pull-request Actions evidence remains unavailable.
- DONE: Commit a complete implementation Stage Report that identifies the candidate tip, changed surfaces, test results, and readiness for independent validation.
  The candidate changes one workflow, five guard files, and two documents. It is ready for fresh validation with live evidence outstanding.

### Summary

The adopted candidate restores an optional Pi Luna/max cadence and retains its evidence. Manual Opus and Pi routing is exclusive, while pull requests retain Sonnet and Codex. The candidate is clean and ready for independent validation at `e838fba693d5562fa4cee99cfbfb682014331a4f`.

## Stage Report: validation

- FAILED: Independently reproduce AC-1 through AC-5 against exact candidate `e838fba693d5562fa4cee99cfbfb682014331a4f`, including event-matrix, model, artifact, exclusivity, registry, owner, formatting, full-suite, and race evidence.
  AC-1 failed because run `31323946889` retained no journey-metric file.
- FAILED: AC-1 (VALUE) - A maintainer can obtain retained Pi CI evidence on demand.
  The exact Pi run passed, but its downloaded artifact contained zero files under `live-artifacts/journey-metrics`.
- SKIPPED: AC-2 - Pi does not become a pull-request merge requirement.
  No pull-request run exists for the exact candidate. The parsed condition excludes Pi from pull requests.
- DONE: AC-3 - A Pi-only dispatch spends no Claude or Codex lane.
  Run `31323946889` ran `offline` and `pi-live`. GitHub skipped `claude-live` and `codex-live` without steps.
- FAILED: AC-4 (VALUE) - Each manual release cadence spends only its selected runtime.
  The Pi cadence used only `CI-E2E-PI`. No exact-candidate Opus run exists, so Opus evidence is unavailable.
- DONE: AC-5 - The desired journey registry stays unchanged.
  Registry reconciliation passed for the same 17 journeys, fixtures, owners, and selector.
- DONE: Verify the approved eight-file/+380 boundary and absence of scope creep.
  The candidate changes eight approved files and adds 244 net lines. No live journey, fixture, registry, product, or release workflow changed.
- DONE: Verify the existing metrics format and artifact records.
  Both Go detail JSONL files and all session JSONL files passed `jq`. The metrics packages and format files are unchanged.
- DONE: Run formatting, focused guards, registry reconciliation, the mutable owner check, the full suite, and the race suite.
  `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` passed.
  Focused `./internal/release ./internal/contractlint` tests passed. The owner test passed with the exact `SPACEDOCK_LIVE_STATE_DIR` path.
- DONE: Run the required detached adversarial audit for CI/release machinery on a throwaway checkout and show that the candidate guards fail under a claim-breaking edit.
  A detached edit let Opus start Codex. `TestRuntimeLiveWorkflowHasOneExplicitClaudeCadence` failed with the changed condition.
- DONE: Report a PASSED or REJECTED recommendation with every live lane classified from current evidence; do not call unavailable, skipped, unapproved, stale, or different-candidate evidence green.
  Recommendation: REJECTED. Pi evidence is current; exact-candidate Opus and pull-request evidence is unavailable.
- FAILED: Material AC-1 outcome defect — retained Pi journey metrics are absent.
  The normal trigger is `live_cadence=pi`. The missing files remove promised evidence from the maintainer artifact.
  Authority: `value-ac[AC-1]` requires retained journey metrics. Run `31323946889` and the empty metrics path prove the trigger.
  The First Officer selected `ROUTE FOR DECISION`. The likely emitter repair is outside the approved eight-file surface.

### Summary

The exact Pi cadence passed with Luna/max and exclusive `CI-E2E-PI` use. The artifact did not contain promised journey metrics, so AC-1 failed. The recommendation is REJECTED, and the candidate remains unchanged.

## Implementation correction design reset

Captain approval adds two allowed files: `internal/ensigncycle/pi_shared_live_runner_test.go` and `internal/ensigncycle/journey_metrics_live_test.go`.
The expected correction has 35-80 gross insertions, at most 10 deletions, and at most +120 net lines.
The full candidate retains the original eight-file limit plus these two exact files.

The end-user value remains unchanged. A successful manual Pi cadence retains at least one non-empty record for each completed journey.
These records use the existing journey-metrics artifact and format. AC-1, runtime targets, lane policy, and result format do not change.

## Stage Report: implementation (cycle 2)

- DONE: Before candidate edits, update and commit the 0a task body with the authorized two-file design reset, gross/net estimate, +120 net tolerance, preserved AC-1 value, and unchanged metrics format.
  State commit `dc3f25d1e` records both files, the 35-80 gross estimate, the +120 net limit, and the unchanged value and format.
- DONE: Add the minimum Pi metrics producer path plus adjacent proof, reuse existing `journeymetrics` APIs, and show the focused proof fails on the current no-op before it passes on the repair.
  The initial focused proof failed with `Pi journey metric files = 0, want 1` against `e838fba69`.
  The durable example now passes and proves one non-empty Pi record with its scenario, runtime, model, and duration.
  The repair uses `BuildRecord` and `EmitRecord` in the existing `shared-scenarios` directory.
- DONE: Commit the correction and a complete implementation Stage Report with offline/race/registry evidence and exact repaired-candidate readiness for the Pi run.
  Candidate `f5c4b32d48f036fddb52233e1c2ca9a1fbf2baa1` is pushed and ready for the exact Pi run.
  `gofmt`, the focused proof, both full suites, registry reconciliation, and the mutable owner check passed.
  The correction has 51 insertions, 2 deletions, and +49 net lines. The full candidate has ten approved files and +293 net lines.

### Summary

The Pi driver now emits one existing-format metrics record after each completed common journey. The focused proof detects the prior empty artifact and passes on the repair. Paid CI remains for fresh validation.

## Stage Report: validation (cycle 2)

- FAILED: Re-review exact candidate `f5c4b32d48f036fddb52233e1c2ca9a1fbf2baa1`, reproduce the focused no-op/fixed metrics proof, and verify the bounded reset stayed within its two-file/+120 limits and existing format.
  The fixed example passed, but the same example also passed after the detached audit restored the driver no-op.
- DONE: Verify the authorized correction boundary and existing record format.
  The correction changes two approved files, adds 51 lines, deletes 2 lines, and adds 49 net lines.
  The helper uses `journeymetrics.BuildRecord` and `journeymetrics.EmitRecord`. Its record has the expected scenario, runtime, model, duration, and bytes.
- FAILED: Reproduce the required focused no-op failure and corrected pass.
  `Example_emitPiScenarioMetrics` passed on the candidate and passed after `piSharedLiveDriver.emitMetrics` became an exact no-op.
- SKIPPED: Run the exact repaired manual Pi cadence and inspect the downloaded artifact for non-empty existing-format journey-metrics records for every completed common journey; classify any auth or protected-environment wait honestly.
  The First Officer stopped live execution after the Material evidence defect. No protected environment received an approval request.
- SKIPPED: AC-1 (VALUE) - A maintainer can obtain retained Pi CI evidence on demand.
  No exact-candidate Pi artifact exists. The focused proof cannot detect the original no-op regression.
- SKIPPED: AC-2 - Pi does not become a pull-request merge requirement.
  No exact-candidate pull-request run exists. The unchanged parsed condition excludes Pi from pull requests.
- SKIPPED: AC-3 - A Pi-only dispatch spends no Claude or Codex lane.
  No exact-candidate Pi run exists. Structural conditions are unchanged, but they are not live evidence.
- SKIPPED: AC-4 (VALUE) - Each manual release cadence spends only its selected runtime.
  Exact-candidate Pi, Opus, and pull-request evidence is unavailable.
- DONE: AC-5 - The desired journey registry stays unchanged.
  Registry reconciliation and the mutable owner test passed for the same 17 journeys.
- SKIPPED: Re-run the applicable offline/race/registry/adversarial evidence and report AC-1 through AC-5 with a PASSED or REJECTED recommendation, never treating unavailable exact-candidate evidence as green.
  Registry, owner, formatting, focused, and adversarial checks ran. The First Officer stopped the full and race suites before the authorized fix.
- FAILED: Material AC-1 evidence defect — the focused proof does not exercise the Pi driver seam.
  The normal user runs the manual Pi cadence and expects retained metrics.
  The guard permits the exact no-op that removes all Pi journey metrics.
  Authority: `value-ac[AC-1]` requires retained journey metrics.
  The detached no-op edit passed `Example_emitPiScenarioMetrics`, which proves the trigger.
  The First Officer selected `FIX`. The next proof must call `piSharedLiveDriver.emitMetrics` and fail on the no-op.
- DONE: Report a PASSED or REJECTED recommendation with all live lanes classified from current exact-candidate evidence.
  Recommendation: REJECTED. All exact-candidate live lanes remain unavailable, and none is classified as green.

### Summary

The correction is inside the authorized boundary and uses the existing record format. Its focused proof does not detect the exact driver no-op regression. The recommendation is REJECTED, and the candidate remains unchanged.

## Stage Report: implementation (cycle 3)

- DONE: Replace the focused proof so it executes the supported `piSharedLiveDriver.emitMetrics` seam and fails when that method is restored to the exact no-op regression.
  The fixed seed calls the driver seam. The detached no-op failed with `Pi journey metric files = 0, want 1`.
  The exact candidate passed the same focused command for `FuzzPiSharedLiveDriverEmitsJourneyMetric/seed#0`.
- DONE: Keep the correction inside the same two files, existing metrics format, and +120 net limit; do not add another helper, schema, format, or producer path.
  The correction has 57 insertions, 2 deletions, and +55 net lines across the same two files.
  The proof replaces the bypassing example, and the producer still uses the existing `BuildRecord` and `EmitRecord` path.
- DONE: Commit the proof correction and a complete Simplified-English implementation Stage Report with focused no-op/fixed output and full offline/race/registry evidence.
  Candidate `7f5f79aadb90b30e72eb243fb91732a4cf6063a7` is pushed and ready for fresh validation.
  `gofmt`, both full suites, registry reconciliation, and the mutable owner test passed.

### Summary

The focused proof now calls the same Pi driver seam as the common journeys. It fails on the original no-op and passes on the repaired candidate. The full candidate changes ten approved files and adds 299 net lines.

## Stage Report: validation (cycle 3)

- DONE: Re-review exact candidate `7f5f79aadb90b30e72eb243fb91732a4cf6063a7` and verify the correction boundary.
  The correction changes the same two approved files and adds 55 net lines.
  The producer still uses `journeymetrics.BuildRecord` and `journeymetrics.EmitRecord`.
- DONE: Reproduce the focused Pi driver proof and the detached no-op failure.
  The exact seed passed through `piSharedLiveDriver.emitMetrics`.
  The detached no-op failed with `Pi journey metric files = 0, want 1`.
- DONE: Run the required local and offline checks without changing candidate bytes.
  `gofmt`, `go test ./...`, and `go test ./... -race` passed.
  Registry reconciliation and the mutable owner check passed for the same 17 journeys.
  Local Pi stopped before execution because the OpenAI API key was unavailable.
- DONE: AC-1 (VALUE) - A maintainer can obtain retained Pi CI evidence on demand.
  Exact run `31338875783` passed on serial attempt 2 with Pi Luna/max.
  Artifact `9045701917` has 12 non-empty schema-v2 records for all 12 completed executions.
  Each record names Pi, Luna/max, the scenario, a positive duration, and run `31338875783`.
- SKIPPED: AC-2 - Pi does not become a pull-request merge requirement.
  The parsed pull-request conditions exclude Pi. No exact pull-request run exists before the normal merge boundary.
  The Captain permits validation approval with AC-2 pending. This approval does not satisfy AC-2.
  Before merge, Sonnet and Codex must pass. Pi must skip with zero steps and no approval.
- DONE: AC-3 - A Pi-only dispatch spends no Claude or Codex lane.
  Exact run `31338875783` used only `CI-E2E-PI` after offline tests.
  Claude and Codex skipped with zero steps. No unrelated environment requested approval.
- DONE: AC-4 (VALUE) - Each manual release cadence spends only its selected runtime.
  Exact Opus run `31340713337` requested only `CI-E2E-OPUS` after offline tests.
  Its Pi and Codex jobs skipped with zero steps. Exact Pi evidence proves the opposite exclusive route.
- DONE: AC-5 - The desired journey registry stays unchanged.
  Registry reconciliation and active-owner checks passed for all 17 registered journeys.
- DONE: Inspect exact artifacts, job records, deployment records, and excluded job shapes.
  The Pi artifact JSONL, session files, model files, and metrics files parsed successfully.
  The Opus artifact `9046295047` retained its failed scenario streams and final messages.
- DONE: Record the first Pi attempt without assigning an unsupported candidate repair.
  Attempt 1 failed before metrics emission. The same recorded-gate journey passed at baseline `e838fba`.
  The authorized serial rerun passed. The First Officer declined a candidate change.
- DONE: Record two Material sprint findings from the exact Opus run.
  `default-headless-gate-stop` crossed its no-authority boundary before validation.
  `rejection-flow` did not invoke `gate record --round validation/1`.
  The second failure also occurred in run `31321450831`.
  These are existing product gaps outside this task. The First Officer declined a task change and rerun.
- DONE: Report the validation recommendation from current exact-candidate evidence.
  Recommendation: PASSED. AC-1, AC-3, AC-4, and AC-5 passed. Structural pull-request exclusion passed.
  AC-2 remains pending at the normal merge boundary under the Captain's ruling.

### Summary

The exact Pi cadence retained complete metrics and used only its selected runtime. The exact Opus cadence also proved runtime exclusion. Validation recommends PASSED, with AC-2 pending at merge.

## Stage Report: implementation (cycle 4)

- DONE: Rebase the registered 0a branch onto exact `origin/main` tip `944f323ac`, preserve only the approved candidate changes, and stop without guessing if any content conflict occurs.
  The rebase completed without conflicts and produced candidate `3e1147f50` on `944f323ac`.
  The remote branch required an ancestry-only merge for a normal push. That merge was conflict-free and kept the exact tree.
- DONE: Compare the rebased patch to candidate `7f5f79aadb90b30e72eb243fb91732a4cf6063a7`, name every changed product byte or source binding, and keep the ten-file/+299 approved boundary.
  No product byte or source binding changed. Both patches have stable patch ID `7c0cea6b`.
  The rebased patch changes the same ten files, with 342 insertions, 43 deletions, and +299 net lines.
- DONE: Run formatting, focused, full, race, registry, and owner checks; push the rebased candidate and commit a complete Simplified-English implementation Stage Report.
  `gofmt`, the focused Pi proof, both full suites, registry reconciliation, and the mutable owner test passed.
  Candidate `20910e6351a42f0f2ca30a92c269ddffbc069582` is pushed with the rebased tree and preserved remote ancestry.

### Summary

The candidate now includes exact integration tip `944f323ac` without a content change. All checks pass on the rebased tree. The normal push preserved the old remote ancestry and the approved patch.

## Stage Report: validation (cycle 4)

- DONE: Independently verify candidate `20910e6351a42f0f2ca30a92c269ddffbc069582` contains `origin/main` tip `944f323ac` and preserves the approved pre-rebase patch byte-for-byte with no source-binding change.
  `944f323ac` is an ancestor, and the first-parent tree equals the candidate tree.
  Both complete plain patches are byte-identical and have stable patch ID `7c0cea6b`.
- DONE: Re-run focused, formatting, full, race, registry, owner, and adversarial checks on the new candidate; confirm the ten-file/+299 boundary.
  `gofmt -d ./cmd ./internal`, both full suites, and all focused checks passed.
  The registry and owner checks passed for 17 journeys. The patch changes ten approved files and adds 299 net lines.
  The detached no-op failed with `Pi journey metric files = 0, want 1`.
- DONE: Report PASSED or REJECTED for fresh validation, carry exact Pi/Opus evidence only if patch identity justifies it, and retain AC-2 as a mandatory exact-PR merge blocker.
  Recommendation: PASSED. Exact patch identity permits the prior live evidence to remain current.
- DONE: AC-1 (VALUE) - A maintainer can obtain retained Pi CI evidence on demand.
  Pi run `31338875783` attempt 2 passed. Artifact `9045701917` retained 12 valid metrics records for 12 completed executions.
- SKIPPED: AC-2 - Pi does not become a pull-request merge requirement.
  No exact pull-request run exists for `20910e635`. AC-2 remains a mandatory merge blocker under the Captain's ruling.
  The exact PR run must pass Sonnet and Codex. Pi must skip with zero steps and no approval.
- DONE: AC-3 - A Pi-only dispatch spends no Claude or Codex lane.
  Pi run `31338875783` used `CI-E2E-PI`. Its Claude and Codex jobs skipped with zero steps.
- DONE: AC-4 (VALUE) - Each manual release cadence spends only its selected runtime.
  Opus run `31340713337` used only `CI-E2E-OPUS`. Its Pi and Codex jobs skipped with zero steps.
  The run retained the prior task-external failures. The First Officer declined a task change and rerun.
- DONE: AC-5 - The desired journey registry stays unchanged.
  Registry reconciliation and the mutable owner join passed for the same 17 journeys.
- DONE: Preserve candidate bytes and source bindings during fresh validation.
  The worktree stayed clean at `20910e6351a42f0f2ca30a92c269ddffbc069582`.
- DONE: Explain why no paid live rerun is necessary after the rebase.
  The rebase changed no task byte, source binding, or test behavior. A paid rerun cannot add candidate-specific evidence.

### Summary

Fresh validation passed on the rebased candidate. The patch and all local checks match the approved candidate.
The Pi and Opus evidence remains valid because the task patch is byte-identical. AC-2 remains pending at the exact pull-request merge boundary.
