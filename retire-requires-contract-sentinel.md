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
