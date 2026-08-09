---
title: Restore optional manual Pi common-live CI
status: implementation
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
