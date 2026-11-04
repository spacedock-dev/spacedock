---
title: Expose complete gate history as machine-readable JSON
status: backlog
source: "Captain authorized issue 814 on 2026-11-03: Proceed with the 3. You have the conn to approve gate and open pr"
issue: spacedock-dev/spacedock#814
sprint: paseo-integration
id: 8v4gxdyfenqxz3t9wpc2c474
gates:
    version: 1
    records:
        - id: gate:8v4gxdyfenqxz3t9wpc2c474:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:8v4gxdyfenqxz3t9wpc2c474-backlog-1
              briefing:
                id: briefing:8v4gxdyfenqxz3t9wpc2c474:backlog:attempt-1:revision-1
                digest: sha256:1eff4abee71d034d2bb6e6bb157e4b2ecc422e69e87742e9773adb551d07a908
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:8v4gxdyfenqxz3t9wpc2c474:backlog:1
                briefing: briefing:8v4gxdyfenqxz3t9wpc2c474:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-04T07:47:02.951714Z"
                decision: approve
                reason: The committed seed addresses issue 814 with complete stored-tree fidelity, read-only behavior, output stability, and focused proof requirements.
                conn:
                    quote: Proceed with the 3. You have the conn to approve gate and open pr
                    source: Captain message in this Codex conversation on 2026-11-03
              application:
                target-stage: ideation
                state: pending
---
Paseo can show gate history from the binary-owned schema without parsing Markdown itself.

## Problem

Issue https://github.com/spacedock-dev/spacedock/issues/814 reports that status --read --json flattens nested gates into an empty string.
The consumer needs every attempt, briefing binding, resolution reason, and application state to distinguish spent approvals from the live approval.
Top-level help also advertises gate validate although the dispatcher rejects it.

## Proposed approach

Expose the stored gate tree through the smallest suitable read-only CLI surface.
Ideation must select either a compatible status JSON extension or a focused gate reader, with a concrete schema and command examples.
Resolve the advertised gate-validate mismatch within that command boundary.

## Risk evidence

Reproduce the issue with a fixture containing three attempts and distinct superseded and pending applications.
Compare the emitted tree with the independent fixture, including exact Unicode reason bytes and array order.
Follow the debugging skill before fixes.

## Out of scope

Changing gate decision authority, application semantics, stored schema, the Paseo plugin, or unrelated status projections.

## Acceptance criteria

**AC-1**: A machine-readable command exposes complete stored gate records, attempts, briefing data, resolutions, reasons, and application states.
Verified by: exact independent JSON fixtures through the real command. Dropping a nested field or confusing pending and superseded must fail.

**AC-2**: The reader preserves order and exact scalar values, including captain reasons, and handles empty gate history.
Verified by: command fixtures with repeated attempts, Unicode, and empty history. Reordering or rewriting a value must fail.

**AC-3**: Reading history changes no entity or workflow bytes, and existing status output remains stable outside the explicitly selected JSON surface.
Verified by: before/after file and Git checks plus existing output fixtures.

**AC-4**: Gate help names supported commands and documents the chosen reader with usable examples.
Verified by: command behavior and help fixtures, plus review of the documentation diff.

## Test plan

Add focused failing tests before implementation. Reuse the existing status reader and gate CLI proof owners.
Run go test ./..., go test ./... -race, and gofmt -w ./cmd ./internal.
Apply the workflow's audit trigger to the actual diff. No new standing enforcement process is commissioned.
Declare expected files, net LOC, tolerance, semantic changes, and concrete documentation diff during ideation.

### Feedback Cycles
