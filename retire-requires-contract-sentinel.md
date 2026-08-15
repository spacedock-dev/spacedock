---
id: 6qhgsezz7v4g4h76t0jf98b0
title: Retire the extinct requires-contract manifest sentinel
status: backlog
source: "0.27 audit (2026-08-14) Priority 2; pre-ship cleanup companion to remove-startup-capability-probe"
started:
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:6qhgsezz7v4g4h76t0jf98b0:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:6qhgsezz7v4g4h76t0jf98b0-backlog-1
              briefing:
                id: briefing:6qhgsezz7v4g4h76t0jf98b0:backlog:attempt-1:revision-1
                digest: sha256:886472682b4b65555d37be5b39151f04f478b618bebafce013e3c41dea56d935
                request-digest: sha256:e3dfe8c776eebd1adfb7660b57d36d1886b392d9e26c2d6f3cc392d374eab563
                room-ref: ./retire-requires-contract-sentinel/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:6qhgsezz7v4g4h76t0jf98b0:backlog:1
                briefing: briefing:6qhgsezz7v4g4h76t0jf98b0:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T02:53:19.566854Z"
                decision: approve
                reason: 'Captain ruling 2026-08-14 (dispatch them): approved into ideation'
              application:
                target-stage: ideation
                state: pending
---

Retire the pre-0.19 `requires-contract` sentinel; its audience is extinct and production code explicitly ignores it. Pure deletion, no behavior change.

The audit scoped this as "testdata fixtures only", but a live grep (2026-08-14) finds 8 source/test files referencing it: internal/contract/doctor.go (a comment stating the field "is not read"), internal/release/release_test.go, internal/cli/codex_plugin_dir_test.go, internal/cli/codex_name_match_test.go, internal/cli/upgrade_from_stale_test.go, internal/cli/decoupling_behavior_test.go, internal/cli/install_behavior_codex_test.go, skills/integration/plugin_manifest_test.go. Historical mentions in docs/roadmap and archived state entities stay — they are records, not live surface.

## Problem

{Ideation fills this in.}

## Proposed approach

{Ideation fills this in.}

## Out of scope

Priority-3 audit items (`gate validate` demotion, `provider-evidence` schema strip, `hold` decision) — each needs its own captain call.

## Expected surface and tolerance

Estimate net LOC change: -NN, across ~9 files.

## Acceptance criteria

**AC-1 - No live source, test, or fixture file references requires-contract.**
Verified by: `grep -rl "requires-contract" internal skills cmd` returns no matches.

**AC-2 - The suite stays green after the deletion.**
Verified by: `go test ./...` and `go test ./... -race` pass.

## Test plan

Deletion-only: the existing suite must stay green; no new tests needed.
