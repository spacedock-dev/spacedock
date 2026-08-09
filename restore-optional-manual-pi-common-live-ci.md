---
title: Restore optional manual Pi common-live CI
status: ideation
source: "Captain correction, 2026-08-09: PR #639 removed pi-live when it removed Pi from pull-request approvals. Restore manual CI Pi evidence without making Pi a merge requirement."
started: 2026-08-09T15:42:25Z
completed:
verdict:
score: 0.9
sprint: test-behavior-completeness
sprint-readiness: ready
group: live-ci-evidence
worktree:
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
---

Give maintainers one optional GitHub Actions command that runs the Pi common journeys and retains their evidence.

## Problem

PR #639 correctly removed Pi from pull-request approvals. It also deleted the `pi-live` job and replaced manual CI with a local-only command.

The local command uses subscription authentication, but it does not give the team a reproducible CI environment or retained artifacts.

## Required outcome

Add `pi` to the existing manual `live_cadence` input in `runtime-live-e2e.yml`.

For this choice, run the offline job and one approval-gated `pi-live` job. Do not create Claude or Codex jobs.

Keep pull-request behavior unchanged. Pull requests create only the Sonnet and Codex live jobs, so Pi is not a merge requirement.

Use the existing `CI-E2E-PI` environment and the pinned Pi packages. Run the canonical `^TestLiveCommon` selector and the Pi front-door substrate proof. Retain logs, artifacts, and journey metrics.

Do not add a second workflow, scheduler, registry, simulator, or product-behavior repair.

## Acceptance criteria

**AC-1 (VALUE) - A maintainer can obtain retained Pi CI evidence on demand.**
Verified by: one manual `live_cadence=pi` run passes the offline and Pi jobs and uploads the common-journey and substrate evidence.

**AC-2 - Pi does not become a pull-request merge requirement.**
Verified by: one pull-request workflow run creates Sonnet and Codex live jobs and creates no Pi job.

**AC-3 - A Pi-only dispatch spends no Claude or Codex lane.**
Verified by: the manual Pi run creates no Claude or Codex live job and requests only the `CI-E2E-PI` approval.

**AC-4 - The desired journey registry stays unchanged.**
Verified by: registry reconciliation passes with the same 16 common journeys, fixtures, TODO owners, and canonical selector.

## Test plan

Update existing workflow expectations only where the command changed. Do not add a simulator for GitHub Actions.

Run offline, race, formatting, and registry checks. Then run one real manual Pi workflow and inspect its jobs and retained artifacts.

