---
id: sxcjbvk7tvefs8rswcseg3fd
title: Restore model-only Codex bootstrap with exact resume-token heuristic
status: implementation
source: Captain clarification, 2026-07-13
started: 2026-07-13T00:37:33Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-codex-resume-token-bootstrap-heuristic
issue:
---

Restore the default first-officer prompt when Codex receives host options after
the fence, while preserving prompt-free session resumes without reconstructing
Codex argv grammar.

## Behavior rule

After Spacedock consumes its own pre-fence flags, inspect forwarded Codex tokens
only for an exact token equal to `resume`. That exact token suppresses the
bootstrap/default launch posture. Otherwise preserve every forwarded token in
order and append the normal Codex bootstrap prompt. This is a deliberate
heuristic accepted by the captain: an exact `resume` value for a known
value-taking option is sufficiently unlikely to prefer no duplicated option
arity table.

## Out of scope

No Codex flag parser, option grammar table, generic host-command classifier, or
new front-door flag.

## Acceptance criteria

**AC-1 (VALUE) - A no-task model-only Codex launch starts the first officer.**
Verified by: focused fake-host argv test for `-- --model gpt-5.6-sol` expecting the unchanged model pair plus the existing bootstrap prompt.

**AC-2 - A model-plus-resume launch stays prompt-free when `resume` is an exact forwarded token.**
Verified by: focused fake-host argv test for `-- --model gpt-5.6-sol resume <id>`.

**AC-3 - The decision uses only exact-token membership, without a Codex option table or argv reconstruction.**
Verified by: focused source/diff review and behavior tests for both argv shapes.

**AC-4 - Launch documentation describes the exact-token resume heuristic.**
Verified by: updated command reference aligned to the focused tests.

## Test plan

Update the existing fake-host launch-parity tests with model-only and
model-plus-resume cases, run focused `internal/cli` tests, then `go test ./...`.
No live host test is required for deterministic argv assembly.
