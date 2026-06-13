---
id: zrcmxzx5c7arew7r8afmxfbw
title: Non-sandboxed launch uses Claude auto-mode (and Codex equivalent)
status: validation
source: captain (2026-06-12)
started: 2026-06-13T04:31:09Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-non-sandboxed-launch-auto-mode
issue:
sprint: 0201-post-flip-release-model
group: ux-cleanup
sprint-readiness: ready
---

When the launch is not sandboxed, `spacedock claude` should launch with Claude Code's auto-mode, and `spacedock codex` should use the equivalent Codex setting. Captain ask, filed verbatim; ideation should confirm the intended rationale and the exact mode mapping with the captain.

Couples to `startup-sandbox-status` (gja5htstcgjxydcz5h2051wc): both need the same sandbox-state detection; this task consumes it to choose the launch mode, that task surfaces it to the operator.

## Problem

When a launch is **not** sandboxed, the launch posture is conservative-by-default and pays an interaction tax that adds nothing to safety in the cases the captain cares about:

- **`spacedock claude` (unsandboxed)** issues `claude --agent spacedock:first-officer …` with **no `--permission-mode` flag** (`frontdoor.go:299-307`). Claude falls back to its `default` permission mode, prompting the operator for each non-trivial tool action. The captain wants the unsandboxed launch to start in Claude Code's **auto-mode** instead, so the first officer can drive without per-action approval friction.
- **`spacedock codex` (unsandboxed)** issues plain `codex …` (`frontdoor.go:457-464`), keeping codex's own default approval/sandbox posture. The captain asks for "the equivalent Codex setting."

This is the unwrapped arm only. The **sandboxed** arms already set the maximally-permissive posture (`claude --dangerously-skip-permissions`, `codex --dangerously-bypass-approvals-and-sandbox`) because safehouse provides the isolation — those are unchanged. The gap is purely the unsandboxed path's permission posture.

This task **consumes** `startup-sandbox-status` (gja5htstcgjxydcz5h2051wc): the wrap decision the launcher already computes (`wrap := safehouse.Present(dir) || fd.forceSafehouse || len(fd.safehouseFlags) > 0`, `frontdoor.go:294`) IS the sandbox-enabled signal. This task does not re-derive sandbox detection; it reads the existing `wrap` boolean and chooses the unwrapped-arm flag from it.

## Captain-confirmation items (resolve at the ideation gate before implementation)

The checklist requires confirming rationale and the exact mode mapping with the captain. Three points need a decision; the recommended answer is given for each:

1. **Rationale.** Confirmed framing: in an unsandboxed launch the operator has accepted running without isolation, so per-action permission prompting is friction without a matching safety gain — auto-mode lets the FO proceed while still stopping short of the sandboxed arm's full skip-permissions. **Recommendation: accept this framing.**

2. **Claude mode = `auto`.** Claude Code's `--permission-mode` accepts `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk`, `plan` (verified live, `claude --help`). The captain's words map directly to `--permission-mode auto`. **Recommendation: `auto`.** (Note: `auto` is distinct from `bypassPermissions` — the latter is what the sandboxed arm's `--dangerously-skip-permissions` already covers, so reusing it unsandboxed would erase the sandboxed/unsandboxed distinction.)

3. **Codex equivalent — needs an explicit captain choice.** Codex in this version has **no single "auto-mode" flag** (no `--full-auto`/`--yolo`; spike-verified). Its closest interactive control is `--ask-for-approval on-request` ("the model decides when to ask the user for approval"), which is codex's nearest analog to Claude's `auto`. The two defensible options:
   - **(A) `codex --ask-for-approval on-request`** — explicit parallel to Claude's auto-mode; the model self-decides when to escalate. **Recommended.**
   - **(B) No codex change** — leave codex's own default posture, treating "equivalent" as "already adequate." Lower-risk but does not deliver a behavior change for codex.

   Implementation pins whichever the captain selects. The AC below is written for (A) and notes how it narrows under (B).

4. **Operator override.** The operator can already pass `--permission-mode <mode>` (claude) or `--ask-for-approval <policy>` / `--sandbox <mode>` (codex) after `--`; these are in `valueTakingHostFlags` (`frontdoor.go:501-522`). **Recommendation: the spacedock-injected flag must NOT override an operator-supplied one** — if the operator already passed the host's permission/approval flag in the passthrough, spacedock injects nothing (last-wins host parsing would otherwise let the operator override, but injecting a duplicate is noise and some hosts reject repeats). AC-3 pins this precedence.

## Proposed approach

In `runClaude` / `runCodex`, the unwrapped arm gains the permission flag, gated on `!wrap` and on the operator not having already supplied it:

- **Claude** (`frontdoor.go`, the `if wrap { … } else` shape around line 299-307): when `!wrap` and the passthrough carries no `--permission-mode`, append `--permission-mode auto` to `inner` (before the passthrough, mirroring how `--agent` is positioned). When `wrap`, unchanged (`--dangerously-skip-permissions`).
- **Codex** (option A): when `!wrap` and the passthrough carries no `--ask-for-approval`/`-a`, append `--ask-for-approval on-request`. When `wrap`, unchanged (`--dangerously-bypass-approvals-and-sandbox`).

A small `passthroughHasFlag(passthrough, names…)` helper (mirroring the existing `hasPluginDir`, `frontdoor.go:362-369`) reads whether the operator already supplied the flag in either `--flag value` or `--flag=value` form, so injection is suppressed when the operator set it.

The mechanism is already proven (see spike): the launcher composes the inner argv as a `[]string` and the host accepts the flag. No new seam — the existing `fake.launchedArg` capture in the front-door tests is the AC oracle.

If the codex mapping lands on option A, add `--ask-for-approval` (and its short form `-a`) to the codex entry of `valueTakingHostFlags` (`frontdoor.go:513-521`) so the stray-prompt classifier treats an operator-supplied value as the flag's value, not a stray positional — symmetric with claude's `--permission-mode` already being in that set.

## Implementation co-edit — blast radius (REQUIRED, not optional)

Changing the unsandboxed-arm inner argv changes the value asserted by every existing exact-argv (`equalArgv`) front-door oracle for an unwrapped launch. These tests pin the current unsandboxed argv (`claude --agent spacedock:first-officer …` / `codex …` with NO permission flag), so the implementer **must update each `want` slice in the same change** — an implementer who follows only the ACs above will change the production argv and ship a RED package believing it passed. The complete set of co-edit sites (verified during ideation to assert an unwrapped-launch argv):

- `internal/cli/frontdoor_test.go:108` — claude `-p` headless launch `want`.
- `internal/cli/safehouse_frontdoor_test.go:159` — `TestClaudeNoSafehouseLaunchesPlain` `want` (this is also AC-1's edited baseline).
- `internal/cli/safehouse_frontdoor_test.go:317` — `TestCodexNoSafehouseLaunchesPlainNoBypass` `want` (AC-2's edited baseline under option A).
- `internal/cli/frontdoor_stray_prompt_test.go:39` and `:70` — claude stray-prompt `want` argvs (`:70` also carries `--plugin-dir`).
- `internal/cli/plugin_dir_frontdoor_test.go:64` — claude `--plugin-dir` dev-lane `want`.

**`--plugin-dir` is still an unsandboxed (`!wrap`) launch.** A `--plugin-dir` launch relaxes the contract gate but does NOT enable safehouse, so it stays on the `!wrap` arm and DOES receive the injected `--permission-mode auto`. The two `--plugin-dir`-carrying oracles (`frontdoor_stray_prompt_test.go:70`, `plugin_dir_frontdoor_test.go:64`) must therefore gain the injected flag too, placed consistently with the chosen insertion point relative to `--agent` and the injected `--plugin-dir <dir>` prefix. Decide the token position once (recommendation: immediately after `--agent spacedock:first-officer`, before the passthrough) and apply it uniformly across all co-edit sites.

Any other `equalArgv` oracle over an unwrapped claude/codex launch added before implementation joins this list — the success condition (below) is the package-wide green, which catches a missed site even if this enumeration drifts.

## Riskiest-mechanism spike (run first) — DONE

The design's soundness rests on one unverified mechanism: that the launcher can set the target permission mode per host such that the host actually honors it end-to-end. Exercised live on the dev machine during ideation:

- **Claude flag exists and is the named mode.** `claude --help` → `--permission-mode <mode>` choices include `auto` (alongside `acceptEdits`, `bypassPermissions`, `default`, `dontAsk`, `plan`). The captain's "auto-mode" is `--permission-mode auto` verbatim.
- **Claude honors it end-to-end.** `claude --permission-mode auto -p "say OK" --model claude-haiku-4-5` → printed `OK`, exit 0. The flag is accepted and the session runs; the launcher's job is only to compose this token into the inner argv, which it already does for `--agent`.
- **Codex has no single auto-mode flag.** `codex --help` shows no `--full-auto`/`--yolo`; the nearest is `--ask-for-approval on-request` (and the `--sandbox` policy is a separate axis). This is what makes item 3 a captain decision rather than a mechanical mapping — recorded so the gate review sees it.

These observations seed the first implementation tests: the AC-1/AC-2 expected argv strings come from the proven flag tokens, not from prose.

## Out of scope

- **No change to the sandboxed (wrapped) arms.** `--dangerously-skip-permissions` (claude) and `--dangerously-bypass-approvals-and-sandbox` (codex) stay exactly as they are; safehouse owns that posture.
- **No new sandbox detection.** The `wrap` boolean is the existing signal from the safehouse probes; this task reads it, it does not compute it (that is `startup-sandbox-status`'s concern).
- **No resume-path change.** A `--resume`/`resume` launch already suppresses the bootstrap prompt and carries its own session intent; the injected permission flag rides the same non-resume gate as the bootstrap prompt, so a resumed session is not forced into auto-mode.
- **No `--permission-mode` value beyond `auto` for claude.** Other modes remain operator-supplied via passthrough.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 — An unsandboxed `spacedock claude` launch issues `claude --permission-mode auto …`, while a sandboxed launch is unchanged (`--dangerously-skip-permissions`, no `--permission-mode`).**
Verified by: a Go test in `internal/cli` driving `runClaude` with (a) no `.safehouse` profile + `lookFound` and asserting `fake.launchedArg` contains `--permission-mode auto` and NOT `--dangerously-skip-permissions`, and (b) a `.safehouse` profile present and asserting the wrapped argv still carries `--dangerously-skip-permissions` and NOT `--permission-mode`. Same `fake.launchedArg` oracle shape as `safehouse_frontdoor_test.go:159` (`TestClaudeNoSafehouseLaunchesPlain`). The expected argv tokens are test-supplied independent values from the proven flag, not read from any instruction file.

**AC-2 — An unsandboxed `spacedock codex` launch issues the captain-chosen codex equivalent (option A: `codex --ask-for-approval on-request …`), while a sandboxed launch is unchanged (`--dangerously-bypass-approvals-and-sandbox`).**
Verified by: a Go test in `internal/cli` driving `runCodex` with (a) no `.safehouse` + `lookFound` asserting `fake.launchedArg` contains `--ask-for-approval on-request` and NOT the bypass flag, and (b) `.safehouse` present asserting the wrapped argv still carries `--dangerously-bypass-approvals-and-sandbox`. Mirrors `safehouse_frontdoor_test.go:317` (`TestCodexNoSafehouseLaunchesPlainNoBypass`). Under captain option (B) this AC narrows to: an unsandboxed codex launch carries no spacedock-injected approval flag (the existing `TestCodexNoSafehouseLaunchesPlainNoBypass` already covers it, and AC-2 is dropped).

**AC-3 — An operator-supplied permission/approval flag in the passthrough suppresses the spacedock-injected one (no duplicate, operator wins).**
Verified by: a Go test driving `runClaude` with `--permission-mode plan` after `--` (no `.safehouse`) and asserting `fake.launchedArg` contains exactly one `--permission-mode` and its value is `plan` (not `auto`); and the codex analog with `--ask-for-approval untrusted`. Expected values (`plan`, `untrusted`, single-occurrence count) are test-supplied and independent of the implementation file.

**AC-4 — The injected flag rides the non-resume gate: a resumed unsandboxed launch is NOT forced into the auto/approval mode.**
Verified by: a Go test driving `runClaude` with `--resume` (no `.safehouse`) and asserting `fake.launchedArg` contains no spacedock-injected `--permission-mode`, mirroring the resume-suppression oracle at `safehouse_frontdoor_test.go:390` (`TestClaudeResumeFamilySuppressesBootstrapPrompt`); codex analog with the leading `resume` subcommand.

## Test plan

All AC proofs are Go unit tests in `internal/cli` driving the real `runClaude`/`runCodex` over the existing `fakeHost`/`lookPath` seam, asserting the issued `fake.launchedArg` per `(wrap, operator-flag, resume)` combination — the same harness `safehouse_frontdoor_test.go` already uses, no new seam, no live host CLI in the test path (the live host probes were the ideation spike). Estimated complexity: **low** — four small table-style tests reusing existing helpers (`safehouseFixtureDir`, `lookFound`, `lookMissing`, `equalArgv`, `executableFixture`).

**Success condition: `go test ./internal/cli/` is green over the WHOLE package after the change — not merely that the new tests pass.** The new ACs add coverage, but the change also rewrites the `want` argv in the existing oracles enumerated under *Implementation co-edit — blast radius*; the package is RED until those are folded into the same change. Whole-package green is what proves both the new behavior AND that no existing unwrapped-launch oracle was left asserting the stale argv. A run that asserts only the new tests in isolation would pass while the package is broken — exactly the RED-package-believed-green failure this fold guards against.

- New tests live alongside the existing front-door oracles in `internal/cli/safehouse_frontdoor_test.go` (or a sibling `frontdoor_permission_mode_test.go`).
- A one-off manual `claude --permission-mode auto -p "say OK"` was the mechanism spike (already run, exit 0) — a sanity check, not an AC proof.
- No live-workflow smoke test needed: the claim is the issued launch argv, fully observable through the `Launch` seam without a running host session.
- Cost note: the change adds at most two argv tokens on the unwrapped arm; it does not touch the gate, the wrap path, or the resume path beyond reading the existing booleans.

## Documentation diff

The `spacedock claude` / `spacedock codex` launch posture is user-visible launch behavior, so per the ideation doc-diff rule the concrete before/after is recorded here for implementation to apply. The exact wording depends on the captain's codex choice (option A vs B); the claude half is fixed.

**`README.md`** — the launch-command lines (the `spacedock claude` / `spacedock codex` front-door entries) currently describe the launch without naming the unsandboxed permission posture. Implementation adds a sentence to the launch section:

> Before:
> `spacedock claude     # launch Claude Code as your first officer`
>
> After:
> `spacedock claude     # launch Claude Code as your first officer (unsandboxed launches start in auto permission-mode; sandboxed launches skip permissions)`

(and the codex line gains the parallel `--ask-for-approval on-request` note under option A, or no change under option B).

**`docs/install-journey.md`** and any host-launch doc that enumerates the launch flags: if a section documents the inner host argv the launcher composes, add the unsandboxed-arm flag (`--permission-mode auto` / `--ask-for-approval on-request`) to that enumeration. Implementation greps `docs/` for `--dangerously-skip-permissions` (the documented sandboxed-arm flag) and adds the unsandboxed-arm counterpart wherever the sandboxed one is described, so the two arms are documented symmetrically. If no doc enumerates the inner argv, the README line above is the only doc edit.

Note: exact line numbers are deliberately not pinned here because the README launch-command block wording is being actively reconciled in the post-flip sprint; implementation locates the `spacedock claude` launch line by its command token, not a line number.

## Stage Report: ideation

- DONE: Confirm with the captain the intended rationale and exact mode mapping: which Claude permission mode equals "auto-mode" and the Codex equivalent, and whether the operator can override. Build on gj's banked sandbox-detection design (do NOT re-derive detection; this consumes it).
  Pinned claude mapping = `--permission-mode auto` (verified live in `claude --help` choices); raised codex as a captain-decision item (no single auto-mode flag exists; recommended `--ask-for-approval on-request`, option B = no change); operator-override precedence pinned (operator-supplied flag suppresses injection, AC-3). Consumes the existing `wrap` boolean (= gj's sandbox-enabled signal) rather than re-deriving detection — recorded in Out of scope.
- DONE: Run the riskiest-mechanism check FIRST: that the launcher can actually set the target permission mode per host (the flag/setting exists and is honored end-to-end), or record "no spike needed" naming the proven launch-arg mechanism.
  Spike section records: `claude --permission-mode auto -p "say OK" --model claude-haiku-4-5` → printed OK, exit 0 (flag exists AND honored end-to-end); `codex --help` shows no single auto-mode flag (forces the captain-decision on the codex side). Launcher composes the token into the inner `[]string` argv exactly as it already does for `--agent`, so no launcher-mechanism unknown remains.
- DONE: Propose the doc diff for the launch-behavior change (per the ideation doc-diff rule) and a test plan whose AC proof is the issued launch argv / observed mode, not prose. Note the dependency on gj's detection.
  Documentation diff section gives README before/after for the `spacedock claude`/`codex` launch lines plus a symmetric `docs/` enumeration rule; all four ACs proven by Go tests asserting `fake.launchedArg` per `(wrap, operator-flag, resume)` over the existing `safehouse_frontdoor_test.go` harness — no live host CLI in the test path, no prose/substring proof. Dependency on `startup-sandbox-status` (gja5htstcgjxydcz5h2051wc) noted in Problem and Out of scope.

### Summary
Defined the unsandboxed-launch auto-mode change as: on the unwrapped arm only, inject `claude --permission-mode auto` (and the captain-chosen codex equivalent, recommended `--ask-for-approval on-request`), suppressed when the operator already supplied the flag or when resuming. The riskiest mechanism — that the host accepts and honors the mode flag end-to-end — was exercised live (`claude --permission-mode auto` ran to exit 0); the codex side surfaced as a genuine captain decision because codex has no single auto-mode flag (no `--full-auto`/`--yolo`), recorded as a gate item with a recommendation. The design consumes the existing `wrap` boolean as gj's sandbox-enabled signal and does not re-derive detection. Four ACs are all proven by the issued launch argv through the existing `fakeHost` seam, never by prose. Open gate decision: confirm rationale, the codex mapping (option A vs B), and the operator-override precedence.

## Stage Report: ideation (cycle 2 — M1 fold)

- DONE: Add an "Implementation co-edit" / blast-radius line naming the five oracle sites as required updates when the unsandboxed argv changes.
  Added the "Implementation co-edit — blast radius (REQUIRED, not optional)" section naming all five sites verbatim (frontdoor_test.go:108; safehouse_frontdoor_test.go:159 and :317; frontdoor_stray_prompt_test.go:39 and :70; plugin_dir_frontdoor_test.go:64), each verified during the fold to assert an unwrapped-launch `want` argv. Also pinned the `--plugin-dir`-is-still-unsandboxed nuance (its two oracles must gain the injected flag too) and a single-insertion-point rule.
- DONE: Broaden the test-plan success condition from "mirror the equalArgv shape" to whole-package green.
  Test plan now states the success condition is `go test ./internal/cli/` green over the WHOLE package, explicitly distinguishing it from new-tests-pass-in-isolation, and names this as the catch for the RED-package-believed-green failure class.
- DONE: (Polish) add `--ask-for-approval` to valueTakingHostFlags under option A; check the gj-codex-enablement prose.
  Added the `valueTakingHostFlags` (`frontdoor.go:513-521`) co-edit note for `--ask-for-approval`/`-a` under option A in Proposed approach. No gj-codex-enablement overstatement to trim: the body references gj only as "consumes the existing `wrap` boolean / no new sandbox detection," which is accurate and does not claim gj's per-runtime enablement probe.

### Summary
Folded staff-review M1: the ideation now carries an explicit, REQUIRED blast-radius section naming the five existing `equalArgv` oracle sites that the unsandboxed-arm argv change breaks, plus the nuance that a `--plugin-dir` launch is still `!wrap` (so its two oracles also gain the injected flag). The test-plan success condition is broadened to whole-package green (`go test ./internal/cli/`), not new-tests-in-isolation, which is the catch against shipping a RED package believed green. Added the `--ask-for-approval` `valueTakingHostFlags` polish under option A; confirmed no gj-enablement prose overstatement exists to trim. No code touched — this is an ideation-body amendment only.

## Stage Report: implementation

- DONE: Unwrapped (`!wrap`) arm ONLY: claude injects `--permission-mode auto`; codex injects `--ask-for-approval on-request` (captain OPTION A) and `--ask-for-approval`/`-a` added to codex `valueTakingHostFlags`. Injection suppressed when operator already supplied the flag OR on resume; sandboxed arms unchanged (AC-1..AC-4).
  `runClaude`/`runCodex` in `internal/cli/frontdoor.go` gain a `!wrap && !resume && !passthroughHasFlag(...)`-gated append after `--agent spacedock:first-officer` (claude) / after `codex` (codex); new `passthroughHasFlag` helper mirrors `hasPluginDir`; codex `valueTakingHostFlags` gains `--ask-for-approval` + `-a`. Commit 2289e998.
- DONE: ALL co-edit oracle sites updated in the SAME change, single insertion point immediately after `--agent spacedock:first-officer`: frontdoor_test.go:108; safehouse_frontdoor_test.go:159 & :317; frontdoor_stray_prompt_test.go:39 & :70; plugin_dir_frontdoor_test.go:64 — including the two `--plugin-dir`-is-still-`!wrap` oracles.
  Folded those plus the spec-unlisted unwrapped oracles caught by the whole-package guard: frontdoor_test.go:488 (skip-check), launch_parity_test.go:167/224/286/306/318/330, and three `TestStrayPromptGuardNegatives` cases. Resume/wrapped oracles correctly left unchanged. New AC coverage in `internal/cli/frontdoor_permission_mode_test.go`.
- DONE: `go test ./internal/cli/` green over the WHOLE package (NOT new-tests-in-isolation). ZERO .md edits — README launch-line diff delegated to yw. `git show --stat` no *.md.
  `go test ./internal/cli/` → 292 passed (was 278; +14 new subtests); full repo `go test ./...` → 1263 passed in 16 packages; `go vet` + `go build` clean. `git show --stat HEAD` = 7 files, all `.go`, zero `*.md`.

### Summary
Injected the unsandboxed-arm permission posture on the `!wrap` path only: claude gets `--permission-mode auto`, codex gets `--ask-for-approval on-request` (captain OPTION A). A new `passthroughHasFlag` helper suppresses injection when the operator already passed the flag (space, equals, or codex `-a` short form), and the injection rides the same non-resume gate as the bootstrap prompt; the sandboxed arms are byte-for-byte unchanged. `--ask-for-approval`/`-a` were added to codex's `valueTakingHostFlags` so an operator value is not misclassified as a stray positional. Every exact-argv oracle the change touched was folded in the same commit — including several unwrapped-launch oracles in `launch_parity_test.go` and `frontdoor_stray_prompt_test.go` that the ideation enumeration did not list but the whole-package-green success condition surfaced. Whole package is green (292) and the full repo passes (1263); the diff is code+tests only, zero `.md` (the README launch-line doc-diff is yw's per the sprint anti-collision design).
