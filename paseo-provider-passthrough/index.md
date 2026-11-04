---
title: Preserve orchestrator control in launcher passthrough mode
status: backlog
source: "Captain authorized issue 813 on 2026-11-03: Proceed with the 3. You have the conn to approve gate and open pr"
issue: spacedock-dev/spacedock#813
sprint: paseo-integration
id: xerfqfkayjawte0gjem7v8v1
gates:
    version: 1
    records:
        - id: gate:xerfqfkayjawte0gjem7v8v1:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:xerfqfkayjawte0gjem7v8v1-backlog-1
              briefing:
                id: briefing:xerfqfkayjawte0gjem7v8v1:backlog:attempt-1:revision-1
                digest: sha256:2feac6d7c7831a290eb2866b9ca6bbd0d50eb165039916586156b5fb8a6b4e0d
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:xerfqfkayjawte0gjem7v8v1:backlog:1
                briefing: briefing:xerfqfkayjawte0gjem7v8v1:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-04T07:48:27.359557Z"
                decision: approve
                reason: The seed defines the reported orchestrator failures, retains sandbox behavior, protects ordinary launches, and requires each host boundary to be exercised.
                conn:
                    quote: Proceed with the 3. You have the conn to approve gate and open pr
                    source: Captain message in this Codex conversation on 2026-11-03
              application:
                target-stage: ideation
                state: pending
---
Paseo can use Spacedock as a provider command while retaining control of its protocol, agent identity, prompts, and permission UI.

## Problem

Issue https://github.com/spacedock-dev/spacedock/issues/813 records bootstrap prompts breaking Codex app-server and injected permission flags bypassing Paseo approvals.
The launcher must retain sandbox wrapping while the orchestrator owns the host session.

## Proposed approach

Design an explicit passthrough mode for Claude, Codex, and Pi.
Retain existing sandbox construction and binary propagation. Skip interactive identity, bootstrap, permission, and compatibility injection in this mode.
Support per-session sandbox environment settings and the common compatibility-bypass flag.
Ideation must reproduce the reported argv and protocol boundaries before selecting the smallest change.

## Risk evidence

Issue 813 includes tested provider commands and failure output.
Probe actual host CLI protocol startup where available. Use exact argv fixtures for each host and wrapping branch.
Follow docs/runtime-support.md and the debugging skill before fixes.

## Out of scope

A new Paseo provider, a replacement process supervisor, a release, and unrelated launcher refactoring.
The optional exec-in-place suggestion is not required.

## Acceptance criteria

**AC-1**: Passthrough Codex app-server and Claude/Pi protocol launches receive valid host argv and a clean stdout channel.
Verified by: command fixtures and a supported live or fixture-backed protocol startup. Adding the bootstrap prompt or banner to stdout must fail.

**AC-2**: The orchestrator retains agent identity, prompting, and permission settings in wrapped and unwrapped passthrough launches.
Verified by: argv fixtures and the Claude stream init permission mode where available. Injecting a first-officer or permission override must fail.

**AC-3**: Passthrough keeps sandbox wrapping, capability selection, and SPACEDOCK_BIN propagation, including per-session SPACEDOCK_SAFEHOUSE and SPACEDOCK_SAFEHOUSE_ENABLE settings.
Verified by: environment and argv fixtures with explicit flag precedence. Losing wrapping or capability propagation must fail.

**AC-4**: Every host accepts --skip-compat-check, passthrough skips compatibility gating, and ordinary interactive behavior remains stable.
Verified by: existing launcher regressions and focused cases for every host. An interactive behavior change must fail.

**AC-5**: Help and product documentation explain the supported mode and environment settings.
Verified by: runnable command examples and review of the changed documentation.

## Test plan

Add focused failing tests before implementation. Reuse existing launcher and safehouse proof owners.
Run go test ./..., go test ./... -race, and gofmt -w ./cmd ./internal.
Validation must perform the required detached adversarial audit for launcher changes.
Declare the expected files, net LOC, tolerance, semantic changes, and concrete documentation diff during ideation.

### Feedback Cycles
