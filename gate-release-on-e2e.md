---
id: nzb7wbwgj93m25ayf9b226xn
title: Gate the release/flip on the live e2e suite (today it's PR-only behind manual env approval)
status: ideation
source: "FO investigation (2026-06-08, during the 0.19.7 cut) — release.yml triggers on tag push and runs goreleaser only; runtime-live-e2e.yml triggers on pull_request to next + manual workflow_dispatch, and its live lanes require per-environment approval (CI-E2E*). So the actual release/tag has NO e2e gate, and PR-time e2e sits 'waiting' for manual approval (it did not gate the 0.19.7 merges)."
score: "0.3"
started: 2026-06-08T15:29:12Z
completed:
verdict:
worktree:
issue:
sprint: 0198-pre-flip-hardening
group: release-gating
sprint-readiness: ready
---

The release/tag path runs no end-to-end check. For the 0.20.0 marketplace flip especially, the release should not publish without the live e2e suite having passed.

## Problem

- `release.yml` (`on: push: tags: v*`) runs goreleaser + the manifest stamp — no e2e.
- `runtime-live-e2e.yml` (`on: pull_request: [next]`, `workflow_dispatch`) runs the shared runtime scenarios, but the live lanes need per-environment approval, so in practice they sit `waiting` and gate neither merges nor releases.
- Net: a release (and the flip) can ship without any live e2e having run green.

## Proposed approach (ideation firms)

Decide the gate shape:
- a pre-release e2e gate (run the live suite on the release tag / a release candidate, block goreleaser on green), or
- require the PR-time e2e (auto-approve the env for trusted branches, or make offline+live a required check before merge to next), or
- a manual `workflow_dispatch` e2e the captain runs before cutting the flip, recorded as the gate.

Scope decision: at minimum the **0.20.0 flip** must not publish without a green live e2e.

## Acceptance criteria (sketch)

- A release (or at least the flip) cannot publish without the live e2e suite having run green — verified by the workflow definition plus an observed/dry-run gate, not just prose.

## Notes

Provenance: surfaced during the 0.19.7 cut. Relevant to the 0200-flip sprint (flip safety). Candidate for a 0.19.8 / flip-prep sprint.
