---
title: Make keep-moving surface in-flight re-shapes across hosts
status: backlog
source: "PR #573 Runtime Live E2E run 30487549723, Sonnet job 90697129966"
started:
completed:
verdict:
score: 0.85
worktree:
issue:
milestone: 0.27.0-pre1
id: z10nqkxg3j2fbm80942kxcv1
---

A First Officer must surface a captain-corrected ticket's re-shape before
stopping at another ticket's gate. The behavior must hold across supported
hosts without grading model wording.

## Problem

PR #573's Sonnet lane failed
`TestLiveClaudeSharedScenarios/keep-moving-posture`. The First Officer stopped
at one valid gate with:

> I'll pause here — the gate is bound and presented, waiting on your decision.

It did not surface the separate `questioned` ticket's corrected re-shape or
name that rework as in flight. This differs from the existing Codex
durable-evidence attribution false-red: the Sonnet final message omitted a real
captain-facing obligation.

The s4 branch changed no keep-moving code. Offline, Pi, Codex, and Opus passed;
the failure therefore belongs to keep-moving reliability rather than gate-room
preparation.

Exact evidence: [PR #573](https://github.com/spacedock-dev/spacedock/pull/573),
[run 30487549723 / Sonnet job 90697129966](https://github.com/spacedock-dev/spacedock/actions/runs/30487549723/job/90697129966).

## Boundary

Preserve the distinction between a true behavioral miss and transcript-shape
variation. Do not add required prose tokens, a general transcript parser, or a
gate-room workaround. Start from the retained run and the existing
keep-moving fixture; determine whether dispatch state, the FO contract, or the
scenario handoff fails to keep the questioned ticket visible at the gate stop.

## Acceptance direction

- A minimized replay or focused live exercise reproduces the Sonnet omission.
- The First Officer surfaces both the current gate and the corrected ticket's
  true state before stopping.
- Negative controls still reject silent absorption and silent waiting.
- Focused Sonnet and Codex keep-moving journeys pass at one exact candidate.
