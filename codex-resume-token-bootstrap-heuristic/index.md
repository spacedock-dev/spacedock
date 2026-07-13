---
id: sxcjbvk7tvefs8rswcseg3fd
title: Restore model-only Codex bootstrap with exact resume-token heuristic
status: done
source: Captain clarification, 2026-07-13
started: 2026-07-13T00:37:33Z
completed: 2026-07-13T01:38:57Z
verdict: PASSED
score:
worktree: .worktrees/spacedock-ensign-codex-resume-token-bootstrap-heuristic
issue:
mod-block: merge:pr-merge
pr: "#502"
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

## Stage Report: validation

**Verdict: PASSED**

- DONE: Independently validated product commit `d0082e3d` in the assigned
  clean worktree, without changing product code. The focused fake-host proof
  passed:
  `go test ./internal/cli -run '^(TestCodexPostFenceUsesExactResumeToken|TestRunCodexPluginDirPostFenceResumeStaysPromptFree|TestRunCodexPluginDirModelOnlyRetainsBootstrapPosture)$' -count=1 -v`.

- DONE: AC-1 — model-only Codex launch retains fresh bootstrap posture. The
  passing `model_only_retains_bootstrap_posture` subtest verifies
  `-- --model gpt-5.6-sol` preserves the model pair, adds
  `--ask-for-approval on-request`, appends the established first-officer
  bootstrap prompt, and emits the launch banner.

- DONE: AC-2 — exact-token resume stays prompt-free. The passing
  `model_plus_resume_stays_prompt_free` subtest verifies
  `-- --model gpt-5.6-sol resume abc-123` forwards those tokens without the
  automatic approval flag, bootstrap prompt, or banner. The paired
  `resume-like_option_stays_a_fresh_launch` case also proves that
  `--resume=abc-123` is not treated as the token `resume`.

- DONE: AC-3 — source and diff review confirms the classification is the
  direct `slices.Contains(fd.passthrough, "resume")` check after the
  pre-fence plugin-dir seam is consumed. The `d0082e3d^..d0082e3d`
  `frontdoor.go` diff adds no Codex option table, argv parser/reconstruction,
  host-command classifier, subprocess, or front-door flag. The focused
  plugin-dir model-only and exact-resume tests both passed, proving the owned
  pre-fence pair is consumed while the remaining forwarded slice receives the
  same exact-token decision.

- DONE: AC-4 — reviewed
  `docs/site/reference/command-reference.md` at `d0082e3d`. It states the
  exact-token heuristic, gives the model-only fresh-launch and
  model-plus-resume prompt-free examples, and explicitly says Codex option
  grammar is not parsed.

- DONE: Performed the required lean detached audit from a clean throwaway
  worktree at `d0082e3d`. I changed the resume assignment back to broad
  nonempty suppression (and removed only the now-unused `slices` import so
  Go could compile), then ran
  `go test ./internal/cli -run '^TestCodexPostFenceUsesExactResumeToken/model_only_retains_bootstrap_posture$' -count=1 -v`.
  It failed as required: actual argv was
  `[codex --model gpt-5.6-sol]`, missing both the automatic approval default
  and bootstrap prompt expected by the model-only oracle. I restored the
  audit edit, confirmed the checkout clean, and removed the detached worktree
  without `--force`.

- DONE: Validation deliberately reran only the focused proof above; it did not
  run Contractlint, modify product code, create a PR, or merge.

### Summary

PASSED. The model-only launch resumes normal first-officer bootstrap/default
posture, while only a literal forwarded `resume` token suppresses it. The
pre-fence plugin-dir seam composes with both paths, and the isolated broad-rule
mutation demonstrably fails the model-only regression oracle.
