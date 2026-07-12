---
title: Codex resume passthrough must not append the default Spacedock prompt
status: ideation
score: 0.7
source: "Captain report 2026-07-10: `spacedock codex -- resume` should not invoke Codex with the default Spacedock prompt."
id: xvcz44jbmye15bpz1ekzxvkc
started: 2026-07-11T04:11:33Z
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
