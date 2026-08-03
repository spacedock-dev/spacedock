---
title: Warn on legacy gate application fields
status: backlog
source: "Captain directive 2026-08-03: legacy application.action and application.blockers must warn, not fail the state read."
score: "1.0"
sprint: durable-decisions
id: jympnaf11wg4qmd4z85a3ayv
gates:
    version: 1
    records:
        - id: gate:jympnaf11wg4qmd4z85a3ayv:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:jympnaf11wg4qmd4z85a3ayv-backlog-1
              briefing:
                id: briefing:jympnaf11wg4qmd4z85a3ayv:backlog:attempt-1:revision-1
                digest: sha256:899a41fbeb26148c4e30be8a8fbe36c8c18fde57e87823e09a3f772a843a0a68
                request-digest: sha256:612e30cdf011456f2204d8c05c092a8213f074722a427806282324bb35215dd8
                room-ref: ./warn-on-legacy-gate-application-fields/review/backlog/briefing-1
---

Legacy state written before the v1 application-schema cut can contain `application.action` and `application.blockers`. The canonical gate reader currently rejects that state before it can expose the valid `target-stage` and `state` fields.

## Problem

`internal/gates/io.go` uses strict YAML known-field decoding for the complete gates tree. The removed application keys therefore produce a fatal `field action not found in type gates.Application` error and block status, eligibility, and gate recovery. These two legacy keys must be advisory compatibility findings; unrelated unknown keys and invalid canonical values must remain fatal.

## Acceptance criteria

**AC-1 — Legacy application.action and application.blockers do not prevent a valid canonical gate read.**
Verified by fixture-backed reads and status/eligibility commands containing each legacy key; both retain canonical `target-stage` and `state` behavior and emit a warning.

**AC-2 — Unrelated unknown gate/application keys and invalid canonical values remain fatal.**
Verified by negative fixtures for an unrelated key, missing target-stage, and invalid state; each exits nonzero without mutation.

**AC-3 — Warning behavior is stable and visible at the operator boundary.**
Verified by CLI golden output that names the entity and legacy field, while successful reads preserve the original bytes until an explicit state write.

## Test plan

Add the smallest reader/validation fixtures and CLI coverage, update the old fail-closed rows only for `action` and `blockers`, retain strict negatives for every other removed shape, run focused gates/status tests, then `go test ./...`, `go test ./... -race`, and gofmt.
