---
id: sxcjbvk7tvefs8rswcseg3fd
title: Restore model-only Codex bootstrap with exact resume-token heuristic
status: validation
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

## Stage Report: implementation

- DONE: Replaced the broad no-task post-fence suppression rule with direct exact
  membership of `resume` in the forwarded Codex slice, after pre-fence
  `--plugin-dir` consumption. The implementation uses
  `slices.Contains(fd.passthrough, "resume")`; it adds no Codex parser, option
  table, argv reconstruction, flag, or subprocess.
- DONE: Added fake-host argv coverage proving `-- --model gpt-5.6-sol` retains
  unsandboxed `--ask-for-approval on-request`, the Codex bootstrap prompt, and
  the launch banner, while `-- --model gpt-5.6-sol resume <id>` remains
  prompt-free. A `--resume=<id>` case pins full-token matching.
- DONE: Preserved and extended the pre-fence `--plugin-dir` seam: it is consumed
  before the heuristic, never reaches Codex, and composes with both the
  model-only fresh launch and exact-resume launch.
- DONE: Updated stale launch, banner, safehouse, and permission tests plus the
  command reference to describe the exact-token heuristic and the separate
  explicit-approval override.
- DONE: Requested independent review. The reviewer approved the diff with no
  Critical or Important findings; its minor documentation and exact-token test
  suggestions were incorporated.
- DONE: Committed and pushed the code branch:
  `d0082e3d fix: restore Codex model-only bootstrap`
  on `spacedock-ensign/codex-resume-token-bootstrap-heuristic`.
- DONE: Verification completed successfully:
  `go test ./internal/cli -run '^TestCodexPostFenceUsesExactResumeToken$' -count=1 -v`;
  `go test ./internal/cli -run 'Codex|OperatorPermissionFlag|ResumeUnsandboxed|LaunchBannerSuppressedOnResume' -count=1 -v`;
  `go test ./internal/contractlint -count=1`;
  `go test ./... -count=1`;
  `go test ./... -race -count=1`;
  `gofmt -w ./cmd ./internal`; and `git diff --check`.

### Summary

Codex now bootstraps model-only launches and suppresses its bootstrap posture
only when its forwarded argv contains the exact token `resume`. The
model-plus-resume path remains prompt-free, and the pre-fence plugin-dir path
remains correctly consumed.
