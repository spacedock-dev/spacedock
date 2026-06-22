---
id: v9pvmzhxvcmvps9tnz73vs4v
title: FO binary launcher invariant — prevent SPACEDOCK_BIN/PATH drift during boot and helper calls
status: implementation
source: captain request after Pi FO used PATH spacedock for later helper calls despite SPACEDOCK_BIN being set
started: 2026-06-22T00:31:00Z
completed:
verdict:
score: 0.35
worktree: .worktrees/spacedock-ensign-fo-binary-launcher-invariant
issue:
sprint: 0230-stable-finalization
---

Tighten the first-officer contract and tests so a FO cannot silently switch from the resolved `SPACEDOCK_BIN` launcher to a different `spacedock` binary on PATH after startup.

## Problem

The FO startup contract says to prefer `${SPACEDOCK_BIN:-spacedock}`, but it also permits bare `spacedock` shorthand in examples. In this session, the initial version gate used the repo-local `SPACEDOCK_BIN`, while later status/state helper calls used bare `spacedock` from `/opt/homebrew/bin`, causing a false subcommand-missing result for `state ready` / `state sweep`.

This matters because launcher drift changes the command surface mid-session and can make the FO reason from the wrong binary's capabilities.

## Proposed approach

Clarify the first-officer contract so startup resolves one launcher variable, reports the resolved binary identity, and requires every Spacedock helper call after the version gate to use that resolved launcher. Remove or narrow the bare-command shorthand allowance so examples do not teach the wrong habit.

Prefer a code gate over prose: add a lint or test that fails on accidental bare `spacedock` helper invocations in FO contract examples, while allowing the explicit PATH fallback/diagnostic probes.

## Out of scope

- Changing launcher behavior for Claude/Codex/Pi hosts beyond the FO contract and tests.
- Reworking all non-FO documentation examples.
- Adding new `spacedock` subcommands.

## Acceptance criteria

**AC-1 - FO startup contract pins a single resolved launcher for the session.**
Verified by: the first-officer contract text requires resolving a launcher variable from `SPACEDOCK_BIN`/PATH, reporting its path/version, and using it for all later Spacedock helper calls.

**AC-2 - Bare `spacedock` examples are not allowed in FO helper-call guidance except explicit fallback diagnostics.**
Verified by: an automated test or lint fails when FO contract/reference docs contain bare `spacedock` helper examples outside an allowlisted fallback/diagnostic context.

**AC-3 - The regression case is covered.**
Verified by: a test fixture or live/fixture-backed check demonstrates that when `SPACEDOCK_BIN` and `command -v spacedock` point to different binaries, FO boot guidance keeps using the `SPACEDOCK_BIN` path for `status`, `state ready`, and other helper calls.

**AC-4 - Filing assertions accept the contract-blessed var-capture launcher idiom.**
Verified by: `assertClaudeFilingViaNew` / `assertCodexFilingViaNew` recognize the `B=${SPACEDOCK_BIN:-spacedock}; $B new <slug>` filing form (where the `$B new` call segment carries no literal `spacedock`/`SPACEDOCK_BIN` token) as atomic filing, scoped to the captured launcher var so an unrelated `$X new` still fails. FO filing behavior is unchanged — the var-capture form is the correct launcher idiom.

## Test plan

Add a focused contract-doc lint or Go test over `skills/first-officer/references/*.md` (and any generated shipped scaffolding if applicable). Include a fixture for the allowed initial fallback probe and disallowed post-boot bare helper examples. Run `go test ./...` as the baseline gate.

## Stage Report: implementation

- DONE: Tighten the FO launcher-invariant contract so startup resolves ONE launcher, reports path+version, and every helper call after the version gate uses it (AC-1), keeping every FO/ensign reference file at/under its v0.22.0 baseline.
  `first-officer-shared-core.md` line 7 invariant rewritten + out-of-range bare-shorthand parenthetical dropped; shared-core lands at 28586/28586. Commit 2a1011b3.
- DONE: Add the code-gate RED-first — a STRUCTURAL lint flagging bare `spacedock` helper examples outside allowlisted fallback/diagnostic context (AC-2), plus a BEHAVIORAL Go test proving SPACEDOCK_BIN-vs-PATH no drift (AC-3).
  AC-2: `internal/contractlint/launcher_invariant_test.go` (lint + discriminator control), RED on 7 doc lines → GREEN after converting them to `${SPACEDOCK_BIN:-spacedock}`. AC-3: `internal/dispatch/launcher_invariant_drift_test.go`, two real binaries, RED-proven via drift control. Commits 2a1011b3, 8ce7049b.
- DONE: Folded-in AC-4 — harden assertClaudeFilingViaNew / assertCodexFilingViaNew to accept the `B=${SPACEDOCK_BIN:-spacedock}; $B new <slug>` var-capture filing idiom, RED first then loosen; do NOT change FO behavior.
  AC-4 added to task body above; `internal/ensigncycle/shared_filing_test.go` loosened (captured-var launcher), RED-first in `shared_filing_negative_test.go` + a guard that an unrelated `$X new` still fails. Commit 11e1c927.

### Byte budget (wc -c, now / v0.22.0 baseline)

first-officer-shared-core 28586/28586 (0)  fo-dispatch-core 17464/17488 (24)  fo-merge-core 7454/8059 (605)
claude-first-officer-runtime 4387/4575  codex-first-officer-runtime 5991/6004  pi-first-officer-runtime 3754/3754 (0)
ensign-shared-core 8829/8829 (0)  codex-ensign 1522/2390  pi-ensign 1333/1768  claude-ensign 2556/2556 (0)
claude-fo-dispatch 22979 (was 22943, +36) — NOT in the value-gate baseline table; the only file that grew.

### Summary

All four ACs implemented TDD/RED-first and green; `go test ./...` passes (exit 0, no FAIL/panic). AC-1 tightens the launcher invariant to pin one resolved launcher for the whole session and removes the bare-`spacedock` shorthand allowance that taught the drift habit. AC-2 is a structural doc-authoring lint in `internal/contractlint` (the only place tests may read instruction files), scoped to runnable bare `spacedock <helper> --<flag>` invocation examples — exempting `→` capability-binding lines, `${SPACEDOCK_BIN:-spacedock}` forms, and fallback/diagnostic/install contexts — with a paired discriminator control; the 7 flagged invocation examples were converted to the resolved launcher. AC-3 is behavioral: two distinct real binaries (one on SPACEDOCK_BIN, a different `spacedock` on PATH) and an assertion that `LauncherCommand()` resolves the SPACEDOCK_BIN binary for status/state/dispatch/merge, proven non-vacuous by a drift-shape RED control. AC-4 widens the filing assertions to the contract-blessed var-capture idiom, scoped to the captured launcher var. Notable: a pre-existing `go vet` warning in `internal/cli/pi_frontdoor_test.go:701` (untouched by me, out of scope) remains.
