---
id: w8rrjv6jsmgahc269arx9k5v
title: Restore Codex handoff bootstrap with post-fence options
status: validation
source: Captain-reported regression, 2026-07-13
started: 2026-07-12T23:16:40Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-restore-codex-handoff-bootstrap
issue:
mod-block: merge:pr-merge
pr: "#501"
---

Restore a fresh Codex launch when a Spacedock handoff is supplied before the
front-door fence and Codex options follow it after the fence.

## Problem

`spacedock codex @/tmp/handoff-file.md -- --model gpt-5.6-sol` forwards the
Codex model option but drops the Spacedock first-officer bootstrap prompt and
the handoff text. The resume repair made every nonempty post-fence argv suppress
the bootstrap, which is broader than the resume case.

## Behavior rule

Keep post-fence argv opaque by content and order. After Spacedock has consumed
only its pre-fence `--plugin-dir` prefix, the already-parsed `hasTask` bit is the
sole fresh-launch discriminator:

- With a pre-fence task, retain the normal fresh Codex posture: launch banner,
  unsandboxed default approval mode when applicable, and the final
  `codexBootstrapPrompt + " " + task` prompt. Append that prompt after the
  unchanged post-fence tokens.
- With no pre-fence task and nonempty post-fence argv, retain the current opaque,
  prompt-free launch. This includes `resume` without identifying it as a Codex
  command.

The implementation is one injection-gate correction (`!fd.hasTask &&
len(fd.passthrough) > 0` for the opaque no-task case), not a Codex argv parser,
grammar table, or reconstruction layer. The post-fence slice remains untouched
and in operator order.

## Out of scope

No Codex argv parser, new command surface, approval-policy change for no-task
post-fence launches, or changes to resume invocations without a pre-fence task.

## Acceptance criteria

**AC-1 (VALUE) - A fresh handoff launch retains both the handoff and the first-officer bootstrap when Codex options are supplied after `--`.**
Verified by: a focused `runCodex` fake-host argv assertion for
`@/tmp/handoff-file.md -- --model gpt-5.6-sol` that expects
`codex --ask-for-approval on-request --model gpt-5.6-sol
"<codex bootstrap> @/tmp/handoff-file.md"`. The `--model` pair must remain
unchanged and ordered before the final prompt.

**AC-2 - Resume-family post-fence argv without a pre-fence task remains prompt-free.**
Verified by: the existing
`TestCodexPostFenceSuppressesPrompt/post-fence-argv-no-prompt` exact argv
assertion still yields `codex resume abc-123` with no bootstrap token.

**AC-3 - The repair leaves post-fence Codex tokens in operator order without a Codex grammar table.**
Verified by: the focused argv assertion and the existing opaque-passthrough test run.

**AC-4 - The command reference distinguishes an explicit pre-fence task from a no-task opaque Codex launch.**
Verified by: review against the focused argv behavior and the following exact
`docs/site/reference/command-reference.md` correction:

> For Codex, a task before `--` remains a fresh first-officer launch even when
> host options follow the fence: post-fence tokens forward verbatim, and
> Spacedock retains its normal launch banner, default approval posture, and
> bootstrap prompt. Only a no-task command with nonempty post-fence argv is an
> opaque host-directed launch with none of those additions.

The following sentence in the next launch-posture paragraph changes from
`Codex suppresses them for any nonempty post-fence passthrough.` to
`Codex suppresses them only for a no-task nonempty post-fence passthrough.`

## Test plan

Add one `codex-task-before-fenced-model` subtest at the existing
`TestFenceTaskPromptOverride` fake-host seam. It supplies the recorded argv and
asserts the exact launch vector in AC-1; it does not inspect Codex option
spelling. Keep the existing no-task resume subtest as the control and retain the
opaque-passthrough test unchanged. Run the focused `internal/cli` tests, then
`go test ./...`. The behavior is pure argv construction, so no live host or
timeout test is justified.

## Spike disposition

No spike needed: the recorded reproduction identifies the bad argv, and the
current parser already separates `fd.hasTask` from `fd.passthrough` while the
existing fake-host seam records the final argv without a real Codex executable.
The no-task resume control was exercised with
`go test ./internal/cli -run '^TestCodexPostFenceSuppressesPrompt$' -count=1`
(pass). No Codex option classification is needed or permitted to prove this
branch.

## Stage Report: ideation

- DONE: Define the smallest behavior rule from the reproduced argv.
  `fd.hasTask`, not any Codex token, selects fresh bootstrap posture; the full rule is recorded above.
- DONE: Specify one focused argv regression test and resume control.
  `TestFenceTaskPromptOverride/codex-task-before-fenced-model` is the new exact-vector oracle; the existing no-task resume control passed.
- DONE: Name the exact command-reference correction.
  The two replacement sentences for `docs/site/reference/command-reference.md` are quoted verbatim in AC-4.

### Summary

The repair is deliberately one boolean boundary: an explicit pre-fence handoff
means fresh launch even when opaque Codex options follow `--`; no-task post-fence
launches stay prompt-free, including resume. The implementation retains operator
tokens verbatim and adds neither a Codex grammar nor a framework.

## Stage Report: implementation

- DONE: Add the exact handoff-plus-model argv regression test first.
  `TestFenceTaskPromptOverride/codex-task-before-fenced-model` failed before the fix and passes in product commit `7b71385a`.
- DONE: Restore fresh-launch posture only when a pre-fence task exists.
  `opaquePassthrough` is now `!fd.hasTask && len(fd.passthrough) > 0`; the existing no-task resume and opaque-passthrough controls pass.
- DONE: Update the launch reference and run focused plus full Go tests.
  The approved command-reference wording landed in `7b71385a`; focused `internal/cli`, `go test ./...`, and `go test ./... -race` passed.

### Summary

The Codex handoff-plus-model launch now preserves the operator's model tokens,
the unsandboxed approval default, and the bootstrap prompt as the final token.
No-task post-fence invocations remain opaque and prompt-free; the product change
is committed and pushed as `7b71385a`.

## Stage Report: validation

- DONE: AC-1 independently exercised the focused fake-host handoff regression.
  `go test ./internal/cli -run '^TestFenceTaskPromptOverride$/^codex-task-before-fenced-model$' -count=1 -v` passed at `7b71385a`. Its exact-vector oracle requires `codex --ask-for-approval on-request --model gpt-5.6-sol "<bootstrap> @/tmp/handoff-file.md"`, preserving the model pair before the final bootstrap-plus-handoff prompt.
- DONE: AC-2 independently exercised the existing no-task resume control.
  `go test ./internal/cli -run '^TestCodexPostFenceSuppressesPrompt$/^post-fence-argv-no-prompt$' -count=1 -v` passed; its exact argv remains `codex resume abc-123`, with no bootstrap token.
- DONE: AC-3 independently exercised opaque no-task forwarding and reviewed the mechanism diff.
  `go test ./internal/cli -run '^TestCodexPostFencePassthroughIsOpaque$' -count=1 -v` passed both the future-token opaque case and bare-Codex bootstrap control. The implementation diff changes the injection gate only to `!fd.hasTask && len(fd.passthrough) > 0`; no Codex argv parser, grammar table, or reconstruction layer appears in the five-file diff. `git diff --check 7b71385a^ 7b71385a` was clean.
- DONE: AC-4 matches the shipped command-reference correction against the exercised behavior.
  The scoped document diff says a task before `--` remains a fresh first-officer launch while only a no-task nonempty passthrough is opaque, and changes default suppression to the no-task case. This agrees with the independently passing fresh-handoff and no-task opaque argv controls.
- DONE: Detached adversarial audit proved the AC-1 oracle can fail.
  Created detached `/tmp/spacedock-codex-handoff-audit.z55qoS` at `7b71385a`, changed only `opaquePassthrough := !fd.hasTask && len(fd.passthrough) > 0` to `opaquePassthrough := len(fd.passthrough) > 0`, then ran `go test ./internal/cli -run '^TestFenceTaskPromptOverride$/^codex-task-before-fenced-model$' -count=1 -v`. It failed with actual argv `[codex --model gpt-5.6-sol]` versus the expected approval/model/final-bootstrap vector (exit 1). The audit worktree was removed with `git worktree remove --force` and confirmed unregistered.

### Recommendation

PASSED. All four acceptance criteria have direct argv or shipped-artifact evidence, and the focused regression fails under the smallest claim-breaking gate edit.

### Summary

Validation reproduced the new fresh-handoff behavior, preserved no-task resume and opaque controls, confirmed the narrow no-parser implementation boundary, and verified that the new oracle rejects the prior broad opaque condition in a disposable checkout.
