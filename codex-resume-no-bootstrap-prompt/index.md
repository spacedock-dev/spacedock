---
title: Codex resume passthrough must not append the default Spacedock prompt
status: validation
score: 0.7
source: "Captain report 2026-07-10: `spacedock codex -- resume` should not invoke Codex with the default Spacedock prompt."
id: xvcz44jbmye15bpz1ekzxvkc
started: 2026-07-11T04:11:33Z
worktree: .worktrees/spacedock-ensign-codex-resume-no-bootstrap-prompt
---

## Problem

The earlier investigation correctly proved the leading form, `spacedock codex -- resume <session>`, but that proof does not cover Codex global options before its subcommand. The captain's exact form, `spacedock codex -- --model <model> resume <session>`, is valid: local Codex `0.144.1` accepts `codex --model probe-model resume --help` and reports `Usage: codex resume [OPTIONS] [SESSION_ID] [PROMPT]`.

This is a real launcher gap. Safe argv captures against both a binary built from current `main` (`f22360de`) and installed `/opt/homebrew/bin/spacedock` (`0.25.0-pre1`) forwarded `--model probe-model resume xv-session` but appended both launcher defaults: `--ask-for-approval on-request` and the Spacedock bootstrap prompt. Plain leading `resume xv-session` captured exactly those two tokens with neither injection. `codexResume` currently tests only `passthrough[0] == "resume"`, so it classifies the option-before-subcommand form as a fresh launch.

## Proposed approach

Make resume classification find `resume` as the first Codex command token after known top-level Codex options, beginning with the captain's `--model <value>` form. The classifier must be value-aware, preserve every forwarded token byte-for-byte, and stop at the first non-option command token rather than treating any later word `resume` as a resume request. It must cover documented equivalent model spellings (`--model value`, `--model=value`, and `-m value`) and avoid a second independent host-argv parser where existing option knowledge can be shared.

Keep `--plugin-dir` consumption before classification and preserve safehouse wrapping: both paths should receive the same classified inner Codex argv. The fix suppresses the launch banner, default approval mode, and bootstrap prompt only for the recognized resume form; a model-selected fresh launch remains a fresh launch.

Documentation diff proposed for `docs/site/reference/command-reference.md` line 27:

```diff
-Anything after `--` forwards verbatim to the host (`--model`, `--resume`, and the like).
+Anything after `--` forwards verbatim to the host (`--model`, `--resume`, and the like); for Codex this includes `spacedock codex -- --model <model> resume <session>`.
```

## Out of scope

- Changing Codex's own resume semantics or session storage.
- Removing the first-officer prompt from fresh launches.
- Broad passthrough-command redesign or an arbitrary word search for `resume`.

## Acceptance criteria

- **AC-1:** `spacedock codex -- --model <model> resume <session>` reaches Codex as the same model option followed by `resume <session>`, with neither the Spacedock bootstrap prompt nor injected `--ask-for-approval on-request`; an argv fixture checks exact order.
- **AC-2:** Equivalent `--model=value` and `-m <model>` forms receive the same resume treatment, while the first non-option non-`resume` command is not misclassified; table-driven argv fixtures check both outcomes.
- **AC-3:** Leading `spacedock codex -- resume <session>` remains prompt- and approval-free, and `spacedock codex -- --model <model>` remains a fresh launch with the normal bootstrap prompt; focused fixtures provide the two controls.
- **AC-4:** Safehouse and local `--plugin-dir` paths preserve the option-before-resume behavior: their inner Codex argv has no injected bootstrap/default approval, and `--plugin-dir` remains consumed rather than forwarded; wrapper fixtures check exact inner argv.
- **AC-5:** The command reference explicitly gives the supported option-before-resume example; the documentation test/review checks the stated command surface.

## Test plan

- Extend the recorded-launch oracles around `TestCodexResumeSubcommandSuppressesPrompt`, `TestResumeUnsandboxedSuppressesInjection`, and `TestCodexFrontDoorInjectsLauncherBinThroughSafehouseResume` with the three model spellings, a non-resume command control, and exact direct/safehouse argv assertions. Add the local-marketplace (`--plugin-dir`) inner-argv fixture without relying on an installed user plugin.
- Run focused `go test ./internal/cli` cases, then the repository baseline and race gates. Validate the documentation diff in the normal docs review.
- Re-run a temporary argv-recording `codex` stub against a freshly built main binary and the resolved installed launcher. The completed spike already proves the failing pre-fix behavior: both launchers appended the bootstrap/default approval after `--model probe-model resume xv-session`; leading resume did not.

## Implementation plan

1. In `internal/cli/launch_parity_test.go`, first add direct argv fixtures for `--model <value>`, `--model=<value>`, and `-m <value>` followed by `resume`, plus a first-command control that must remain fresh. Run the focused test and record the expected pre-change failure.
2. In `internal/cli/frontdoor_permission_mode_test.go` and `internal/cli/frontdoor_test.go`, first add unsandboxed, safehouse, and local `--plugin-dir` argv fixtures that pin no bootstrap/default approval and exact forwarding. Run them red before changing production code.
3. Make the smallest change in `internal/cli/frontdoor.go`: classify only the first Codex command token after known value-taking global options, without rewriting the passthrough slice.
4. Update `docs/site/reference/command-reference.md` with the approved `--model <model> resume <session>` command reference example.
5. Run focused `go test ./internal/cli`, then `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`; record each red/green result in the implementation report.

## Stage Report: ideation

### Trace and root-cause result

The claimed append site is `runCodex` in `internal/cli/frontdoor.go`, but it is already guarded:

1. `parseFrontDoorArgs` places every token after the first `--` in `fd.passthrough`, so the captain's exact command yields `[]string{"resume"}`.
2. Codex-only `--plugin-dir` handling consumes its injected prefix before classification; the remaining passthrough still begins with `resume`.
3. `codexResume(fd.passthrough)` returns true only for that leading subcommand.
4. The same boolean suppresses the launch banner, unsandboxed `--ask-for-approval on-request`, and `launchPrompt(codexBootstrapPrompt, fd)`. Safehouse wrapping then adds only its outer envelope and sandbox-bypass flag.

This is not a new HEAD behavior. `git blame` traces the prompt guard and `codexResume` helper to `31ad7c8e` (2026-05-30), and `v0.24.0` contains the same logic and `TestCodexResumeSubcommandSuppressesPrompt`. The archived `sandbox-flag-passthrough` entity records the original design decision and LP-AC-2 oracle. New parsing would duplicate that precedent.

### Evidence

- Focused current-source tests passed: `TestCodexResumeSubcommandSuppressesPrompt`, `TestCodexFrontDoorInjectsLauncherBinThroughSafehouseResume`, and `TestResumeUnsandboxedSuppressesInjection` (7 subtests total).
- A PTY argv capture against the installed `/opt/homebrew/bin/spacedock` (`0.24.0-pre2`) produced, for `spacedock codex --skip-compat-check -- resume session-42 --last`, exactly `codex`, `resume`, `session-42`, `--last`; exit 0; no prompt.
- Its fresh-launch control produced `codex`, `--ask-for-approval`, `on-request`, then the bootstrap prompt; exit 0. This proves the stub and prompt detector were capable of observing the injection when appropriate.
- Existing current-source coverage additionally pins safehouse resume as `safehouse ... -- codex --dangerously-bypass-approvals-and-sandbox resume abc123`, and `--plugin-dir /co -- resume abc123` as `codex resume abc123`, both without the prompt.

### Recommendation

Mark ideation complete and close the entity as already satisfied/no-op. If the captain saw the old bootstrap text after resume, the leading hypothesis is transcript history from the resumed session, not a fresh argv injection. A different binary on `PATH` is the other plausible explanation, but the locally resolved installed binary also passed; either hypothesis needs a captured failing argv before reopening.

### Summary

Traced the exact command through parsing, plugin-dir consumption, resume classification, prompt suppression, and safehouse wrapping; compared it with the archived launch-parity design and `v0.24.0`; and reproduced both resume and fresh-launch behavior against the installed binary. The reported bug is already fixed and durably tested, so the smallest correct change is no product change and a no-op close.

## Stage Report: ideation (cycle 2)

- DONE: Reproduce the captain's exact option-before-resume argv shape against the current main/installed launcher.
  Safe stubs captured the bootstrap prompt and `--ask-for-approval on-request` after `--model probe-model resume xv-session` from main `f22360de` and installed `0.25.0-pre1`; both exited 0.
- DONE: Determine whether prompt suppression recognizes a resume subcommand after Codex options.
  It does not: `codexResume` sees only a leading `resume`; leading-resume control captured `codex resume xv-session` with neither injection, while safehouse and isolated `--plugin-dir` reproduced the modeled-resume injection.
- DONE: Update the task's diagnosis and proof plan; make no product edits.
  The body now records the confirmed gap, a narrow value-aware classification direction, exact acceptance tests, and the required command-reference diff; no files outside this entity changed.

### Summary

The captain's `--model … resume` form is a real, reproducible regression, not transcript history or a stale binary. This supersedes the prior no-op recommendation, which applied only to leading `resume`; implementation must change resume recognition while preserving forwarded argv and fresh-launch behavior.

## Stage Report: implementation

### Deliverable

Committed [`243ca50`](https://github.com/spacedock-dev/spacedock/commit/243ca50753b9d46e0b9dec0a5151c5f615da7c61) on `spacedock-ensign/codex-resume-no-bootstrap-prompt`.

- `internal/cli/frontdoor.go` now finds the first Codex command token after known global options without mutating forwarded argv. It handles space, equals, and compact short value forms (`--model value`, `--model=value`, `-m value`, and `-mvalue`), stops at the first non-option command, and shares the same classifier with the stray-prompt advisory.
- Direct, safehouse, and local `--plugin-dir` launch fixtures pin exact inner argv, no default approval injection, no bootstrap prompt, and no launch banner for option-before-resume forms.
- `docs/site/reference/command-reference.md` documents `spacedock codex -- --model <model> resume <session>`.

### TDD evidence

- RED: before the classifier change, focused option-before-resume fixtures had 8 expected failures and 2 controls passing; every resume form carried injected approval and the bootstrap prompt.
- RED: after the initial classifier, the new stray-advisory unit fixture failed because `warnStrayPromptAfterDash` still treated option-before-resume as a bare positional. The shared command-token classifier fixed that independent warning path.
- RED: review identified compact `-mvalue`; its focused fixture failed with injected approval/bootstrap. Installed Codex `0.144.1` confirmed `codex -mprobe-model resume --help` enters the `resume` command. The compact-short recognition was then added.
- GREEN: focused CLI regression suite: 29 passed.

### Acceptance evidence

- DONE AC-1: exact direct `--model <model> resume <session>` argv is preserved with no bootstrap, approval injection, or banner.
- DONE AC-2: exact `--model=<model>`, `-m <model>`, and compact `-m<model>` argv are preserved; `--model <model> exec resume …` remains a fresh launch.
- DONE AC-3: existing leading-resume and model-only fresh-launch controls remain green.
- DONE AC-4: safehouse and local `--plugin-dir` fixtures preserve the classified inner argv; the local plugin flag is still consumed.
- DONE AC-5: command reference contains the approved option-before-resume example.

### Verification

- `gofmt -w ./cmd ./internal` run; the pre-existing unrelated formatting drift in `internal/release/journeydelta.go` was restored and excluded.
- `go test ./...` — 2129 passed in 17 packages.
- `go test ./... -race` — 2129 passed in 17 packages.
- Independent code review found and drove the compact-short regression; its follow-up found no remaining Critical or Important issue.

### Dispatch checklist

- DONE: Wrote and committed the task-local implementation plan before production edits, then recorded the focused red/green TDD sequence.
- DONE: Added the shared option-aware Codex resume classifier, direct/safehouse/local-plugin argv fixtures, and the approved command-reference example.
- DONE: Ran focused CLI, repository-wide, and race verification; the final full and race runs each passed 2129 tests in 17 packages.

### Summary

Codex resume classification now recognizes options before `resume` while preserving operator argv and fresh-launch behavior. The committed deliverable is ready for independent validation.

## Stage Report: validation

- DONE: **AC-1:** Preserve `--model <model> resume <session>` with no bootstrap/default approval/banner.
  A freshly built `243ca507` binary and argv-recording `codex` emitted only `--model`, `probe-model`, `resume`, `xv-session`.
- DONE: **AC-2:** Cover equals, short, compact-short, and first-command controls.
  Direct recorder probes preserved `--model=probe-model`, `-m probe-model`, and `-mprobe-model`; `--model probe-model exec resume …` remained a fresh launch. Installed Codex 0.144.1 accepted both long and compact forms before `resume --help`.
- DONE: **AC-3:** Retain leading-resume suppression and model-only fresh launch behavior.
  Leading `resume xv-session` had no extra argv/banner; model-only recorded approval plus bootstrap and the normal banner.
- DONE: **AC-4:** Preserve safehouse/local-plugin inner argv and consume `--plugin-dir`.
  The executable recorder observed safehouse's exact inner `codex --dangerously-bypass-approvals-and-sandbox --model probe-model resume xv-session`; combined local `--plugin-dir` output contained no forwarded plugin flag, approval, or bootstrap.
- DONE: **AC-5:** Document the supported command surface.
  Reviewed `docs/site/reference/command-reference.md`: it explicitly shows `spacedock codex -- --model <model> resume <session>`.
- DONE: Independently reproduce option-before-resume argv behavior across direct, safehouse, and plugin-dir paths.
  Reproduced direct, safehouse, local-plugin, and combined safehouse+local-plugin launches through an executable argv recorder; all preserved the intended inner argv.
- DONE: Verify leading resume and fresh-launch controls remain correct, including compact `-mVALUE`.
  Focused CLI suite and the independent direct probes covered leading `resume`, model-only, later-`resume` after `exec`, and `-mprobe-model`.
- DONE: Validate AC evidence from actual commands/tests and report PASSED or REJECTED without changing product code.
  Focused CLI and focused `-race` suites passed; serialized `go test -p 1 ./...` passed all 17 packages. Code worktree remains clean at `243ca507`.
- SKIPPED: Repository-wide `go test -p 1 ./... -race` completion.
  It was attempted but the host volume exhausted temporary build space (`no space left on device`); failures were environmental before XV assertions. The focused XV race suite passed.

### Summary

Independent executable argv evidence covers every acceptance criterion, including the captain's option-before-resume form and the safehouse/plugin-dir paths. The implementation is behaviorally correct and clean at `243ca507`; recommend PASSED, with the full-repository race gate noted as an environment-capacity follow-up rather than a product defect.
