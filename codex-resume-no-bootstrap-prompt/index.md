---
title: Codex resume passthrough must not append the default Spacedock prompt
status: ideation
score: 0.7
source: "Captain report 2026-07-10: `spacedock codex -- resume` should not invoke Codex with the default Spacedock prompt."
id: xvcz44jbmye15bpz1ekzxvkc
started: 2026-07-11T04:11:33Z
---

## Problem

The captain reported that `spacedock codex -- resume` appended the default first-officer bootstrap prompt. Investigation did not reproduce that behavior: the installed `0.24.0-pre2` binary, the `v0.24.0` source, and current `main` all preserve the resume invocation without a prompt. The reported symptom therefore needs a failing argv capture before product code changes.

## Proposed approach

Close this entity as already satisfied by the archived launch-parity work (`sandbox-flag-passthrough`, implementation commit `31ad7c8e`). Do not add a second parser or broaden Codex grammar speculatively.

If the symptom recurs, capture the resolved `spacedock --version` and the exact launched argv. Reopen only if that capture contains the bootstrap prompt after a leading post-fence `resume`; a resumed transcript displaying its original bootstrap message is not evidence that the launcher injected a new one.

## Out of scope

- Changing Codex's own resume semantics or session storage.
- Removing the first-officer prompt from fresh launches.
- Broad passthrough-command redesign beyond the minimum subcommand classification.

## Acceptance criteria

- **AC-1:** `spacedock codex -- resume` reaches Codex as `codex resume` with no appended Spacedock bootstrap prompt.
- **AC-2:** `spacedock codex -- resume <session>` preserves `<session>` and all later operator arguments in order, with no injected prompt.
- **AC-3:** A fresh `spacedock codex` launch still receives the normal first-officer bootstrap prompt.
- **AC-4:** Safehouse, no-safehouse, and `--plugin-dir` argv fixtures prove the same resume behavior without leaking launcher flags past `--`.
- **AC-5:** No production change lands while the supported release source and installed binary already satisfy AC-1 through AC-4.

## Test plan

- Keep the existing recorded-launch oracles as the durable proof: `TestCodexResumeSubcommandSuppressesPrompt` covers resume plus session ID and the fresh-launch negative control; `TestCodexFrontDoorInjectsLauncherBinThroughSafehouseResume` covers the safehouse envelope; and `TestFrontDoorSubcommandPassthrough` covers `--plugin-dir` consumption followed by resume.
- For any recurrence, use an argv-recording `codex` stub against the resolved installed binary before changing code. Record both the resume invocation and a fresh-launch control.
- No full/race implementation gate is needed for a no-op close; focused current-source tests and the installed-binary argv capture are the evidence gate.

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
