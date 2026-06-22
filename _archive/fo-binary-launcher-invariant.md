---
id: v9pvmzhxvcmvps9tnz73vs4v
title: FO binary launcher invariant — prevent SPACEDOCK_BIN/PATH drift during boot and helper calls
status: done
source: captain request after Pi FO used PATH spacedock for later helper calls despite SPACEDOCK_BIN being set
started: 2026-06-22T00:31:00Z
completed: 2026-06-22T05:33:52Z
verdict: passed
score: 0.35
worktree: .worktrees/spacedock-ensign-fo-binary-launcher-invariant
issue:
sprint: 0230-stable-finalization
mod-block:
pr: pr-merge:433
archived: 2026-06-22T05:33:52Z
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
  `first-officer-shared-core.md` line 7 invariant rewritten + out-of-range bare-shorthand parenthetical dropped; the blanket "may use bare spacedock as shorthand" allowance replaced with a carve-out matching the AC-2 lint allowlist exactly (naming a command / `→` binding lines / fallback probe / diagnostic-install — never a post-gate invocation), so contract and gate agree. shared-core lands at 28583/28586. Commits 2a1011b3, 4d933f6e.
- DONE: Add the code-gate RED-first — a STRUCTURAL lint flagging bare `spacedock` helper examples outside allowlisted fallback/diagnostic context (AC-2), plus a BEHAVIORAL Go test proving SPACEDOCK_BIN-vs-PATH no drift (AC-3).
  AC-2: `internal/contractlint/launcher_invariant_test.go` (lint + discriminator control), RED on 7 doc lines → GREEN after converting them to `${SPACEDOCK_BIN:-spacedock}`. AC-3: `internal/dispatch/launcher_invariant_drift_test.go`, two real binaries, RED-proven via drift control. Commits 2a1011b3, 8ce7049b.
- DONE: Folded-in AC-4 — harden assertClaudeFilingViaNew / assertCodexFilingViaNew to accept the `B=${SPACEDOCK_BIN:-spacedock}; $B new <slug>` var-capture filing idiom, RED first then loosen; do NOT change FO behavior.
  AC-4 added to task body above; `internal/ensigncycle/shared_filing_test.go` loosened (captured-var launcher), RED-first in `shared_filing_negative_test.go` + a guard that an unrelated `$X new` still fails. Commit 11e1c927.

### Byte budget (wc -c, now / v0.22.0 baseline)

first-officer-shared-core 28583/28586 (3)  fo-dispatch-core 17464/17488 (24)  fo-merge-core 7454/8059 (605)
claude-first-officer-runtime 4387/4575  codex-first-officer-runtime 5991/6004  pi-first-officer-runtime 3754/3754 (0)
ensign-shared-core 8829/8829 (0)  codex-ensign 1522/2390  pi-ensign 1333/1768  claude-ensign 2556/2556 (0)
claude-fo-dispatch 22979 (was 22943, +36) — NOT in the value-gate baseline table; the only file that grew.

### Summary

All four ACs implemented TDD/RED-first and green; `go test ./...` passes (exit 0, no FAIL/panic). AC-1 tightens the launcher invariant to pin one resolved launcher for the whole session and replaces the blanket bare-`spacedock` shorthand allowance (which taught the drift habit) with a carve-out that matches the AC-2 lint's allowlist exactly, so the contract and the gate agree. AC-2 is a structural doc-authoring lint in `internal/contractlint` (the only place tests may read instruction files), scoped to runnable bare `spacedock <helper> --<flag>` invocation examples — exempting `→` capability-binding lines, `${SPACEDOCK_BIN:-spacedock}` forms, and fallback/diagnostic/install contexts — with a paired discriminator control; the 7 flagged invocation examples were converted to the resolved launcher. AC-3 is behavioral: two distinct real binaries (one on SPACEDOCK_BIN, a different `spacedock` on PATH) and an assertion that `LauncherCommand()` resolves the SPACEDOCK_BIN binary for status/state/dispatch/merge, proven non-vacuous by a drift-shape RED control. AC-4 widens the filing assertions to the contract-blessed var-capture idiom, scoped to the captured launcher var. Notable: a pre-existing `go vet` warning in `internal/cli/pi_frontdoor_test.go:701` (untouched by me, out of scope) remains.

### Feedback Cycles

**Cycle 1 — validation (detached adversarial audit): one material finding, REJECTED.**

A live counterexample to the invariant survives in a shipped FO reference file, and the AC-2 lint has a hole that let it through:

- **`claude-fo-dispatch.md:132`** instructs the FO to "Run `spacedock dispatch context-budget --name {ensign-name}`" — a bare, imperative, post-gate Spacedock helper invocation. AC-1's tightened line-7 prose forbids exactly this ("never a post-gate helper invocation" for bare `spacedock`), and AC-2 says the lint must flag "bare `spacedock` helper examples outside an allowlisted fallback/diagnostic context." The implementer converted two sibling helper calls in this same file (`spawn-standing-all`, `reconcile`) but missed this third one.
- **AC-2 lint gap:** `lineHasBareLauncherHelperCall` only fires when the invocation carries a flag from its enumerated allowlist (`--workflow-dir|boot|discover|set|next-id|validate|json|resolve|where|next|archived`). `--name` is not in that set, so the lint does not catch line 132. Verified on a throwaway checkout: `lineHasBareLauncherHelperCall(line132) == false`.
- **This is not a vacuity defect** — every AC's discriminator/RED control is genuinely non-vacuous (proven below). It is an *incomplete fix + an under-scoped gate*: the invariant is asserted but a real drift-teaching example ships, and the lint that should catch its class can't.

**Bounce-back (both halves, both proven closable on the throwaway checkout):**
1. Convert `claude-fo-dispatch.md:132` to `${SPACEDOCK_BIN:-spacedock} dispatch context-budget --name {ensign-name}`.
2. Add `name` to the AC-2 lint's invocation-flag allowlist (`launcher_invariant_test.go` `launcherHelperInvocation` regex) so this class is caught going forward; a control proving the `→`-binding twin at `fo-dispatch-core.md:98` stays exempt. With `name` added, the lint flags the bare line-132 form and passes the resolved-launcher fix (audit-confirmed).

Note: `fo-dispatch-core.md:98` carries the same `spacedock dispatch context-budget --name` text but as a `- → **Claude:** PRESENT` capability-binding line — correctly exempt (names the shipped surface, not an FO-emitted call). Only line 132 is a defect.

## Stage Report: validation

- DONE: Verify the value-gate byte budget — every tracked FO + ensign reference file <= v0.22.0 baseline; assess the +36 B growth of untracked `claude-fo-dispatch.md`.
  All 10 tracked files PASS (re-measured `wc -c` vs the index.md gate table, baselines independently re-derived from `git show v0.22.0:`). Per-file table + the +36 B verdict below.
- DONE: Reproduce AC-1..AC-4 from independent evidence with non-vacuous controls; `go test ./...` green; confirm NO FO behavior change.
  AC-2: lint REDs on a planted bare `spacedock status --discover` in a FO ref, GREENs on restore; discriminator control passes. AC-3: two real distinct binaries — REDs (all 5 helper subtests resolve `path-bin`) when `LauncherCommand()` is mutated to the bare-PATH drift shape, GREENs as shipped. AC-4: accepts `B=${SPACEDOCK_BIN:-spacedock}; $B new` and rejects an unrelated `$C new` (var-scoping probed independently). `go test ./...` exit 0, all packages ok.
- FAILED: DETACHED adversarial audit on a throwaway checkout — clean audit OR material findings.
  Audit ran on `/tmp/launcher-invariant-audit-checkout` (a fresh clone, never the worktree). All three named claim-breaking edits REDed correctly and restored green. BUT the audit's "lint passes a planted bare-spacedock helper" probe surfaced a REAL escape the deliverable ships: `claude-fo-dispatch.md:132`. Recorded as `### Feedback Cycles` Cycle 1. Contract/lint agreement otherwise confirmed.

### Summary

Recommendation: **REJECTED** — one material finding, small and bounded. The value gate is met (all 10 tracked files at/under v0.22.0; `first-officer-shared-core.md` is even 28583, 3 B under the 28586 ceiling — the report's "28586/28586" predates the final commit `4d933f6e`). The +36 B in untracked `claude-fo-dispatch.md` (22943→22979) is two `spacedock`→`${SPACEDOCK_BIN:-spacedock}` conversions on the only two runnable helper-call examples (`spawn-standing-all`, `reconcile`) — **necessary, not trimmable**: +18 B/call is the irreducible cost of the resolved-launcher idiom, and leaving them bare would reintroduce the exact drift bug. AC-1/AC-2/AC-3/AC-4 are all reproduced from independent evidence with genuinely non-vacuous controls, `go test ./...` is green, and there is NO shipped FO behavior change (the 4 doc diffs are pure launcher-token swaps; `shared_filing_test.go` is `_test.go`; `build.go`/`launcher_command.go` untouched this task). The blocker: the detached audit caught that `claude-fo-dispatch.md:132` ("Run `spacedock dispatch context-budget --name {ensign-name}`") is a bare post-gate helper invocation the invariant forbids — the implementer converted two sibling calls in that file but missed this third — and the AC-2 lint can't catch it because `--name` is absent from its invocation-flag allowlist. Both fix halves (convert line 132; add `name` to the lint allowlist) are proven to work on the throwaway checkout. This is a live counterexample to the very invariant the task ships, so it must close before PASSED.

## Stage Report: implementation (cycle 2)

- DONE: Convert `claude-fo-dispatch.md:132` `spacedock dispatch context-budget --name {ensign-name}` → `${SPACEDOCK_BIN:-spacedock} ...` — the third runnable helper invocation in that file (sibling to the already-converted `spawn-standing-all` + `reconcile`).
  Converted; landed in untracked `claude-fo-dispatch.md` (+18 B → 22997). Commit 029a4b14.
- DONE: Close the AC-2 lint hole — add `name` to `launcherHelperInvocation`'s invocation-flag allowlist, RED-first.
  Added `name` to the regex; with it added the lint REDs on the still-bare line 132 (and ONLY line 132 — full re-scan flagged no other line), then GREENs after the conversion. Re-confirmed adversarially: lint reds on a line-132 regression, greens on restore. Commit 029a4b14.
- DONE: Keep the `→`-binding twin at `fo-dispatch-core.md:98` (same `dispatch context-budget --name` text) EXEMPT, with a control.
  Added a discriminator control asserting the `--name` bare form IS flagged and the line-98 `- → **Claude:** PRESENT` binding twin is NOT; both pass. The existing `- → ` prefix exemption already covers line 98.
- DONE: Re-scan ALL FO/ensign refs with the broadened lint; convert any new catch within byte budget; re-run `go test ./...` green; re-measure.
  Comprehensive sweep (any bare `spacedock <helper> --<anyflag>` outside `${SPACEDOCK_BIN}`/`→`/diagnostic) surfaces only lines 77 & 185 of shared-core — both command-NAME/signature mentions ("`spacedock new <slug> --id-seed …` mints it", `[--folder]` usage brackets), NOT imperative invocations, so correctly NOT in the defect class and NOT flagged (consistent with the validator scoping the bounce-back to `--name` only). `go test ./...` exit 0, no FAIL/panic.

### Byte budget (cycle 2, wc -c, now / v0.22.0 baseline)

first-officer-shared-core 28583/28586 (3)  fo-dispatch-core 17464/17488 (24)  fo-merge-core 7454/8059 (605)
claude-first-officer-runtime 4387/4575  codex-first-officer-runtime 5991/6004  pi-first-officer-runtime 3754/3754 (0)
ensign-shared-core 8829/8829 (0)  codex-ensign 1522/2390  pi-ensign 1333/1768  claude-ensign 2556/2556 (0)
claude-fo-dispatch 22997 (was 22943 base, +54: 3 launcher conversions) — NOT in the value-gate baseline table.

### Summary (cycle 2)

The Cycle-1 material finding is closed. Line 132's bare context-budget invocation now resolves the pinned launcher, and the AC-2 lint's invocation-flag allowlist gained `name` so this class is caught going forward — driven RED-first (lint flags the still-bare line 132 before the conversion) with a paired control locking the line-132-flagged / line-98-exempt distinction. A full re-scan of every FO/ensign reference surfaces no other escape (the only other bare-flag occurrences are descriptive command-name mentions, outside the invariant's post-gate-invocation scope). `go test ./...` is green, and every value-gate-tracked file remains at/under its v0.22.0 baseline; the one converted line lands in untracked `claude-fo-dispatch.md`, so no net-trim was needed. New tip: 029a4b14.

## Stage Report: validation (cycle 2)

Re-review of the Cycle-1 fix at worktree tip `029a4b14`. Verified closure independently — not by trusting the fix diff.

- DONE: Confirm `claude-fo-dispatch.md:132` converted + full independent re-scan for any other escape.
  Line 132 now reads `${SPACEDOCK_BIN:-spacedock} dispatch context-budget --name {ensign-name}`. Re-scanned every FO/ensign ref: the 5 remaining bare `spacedock <verb> --<flag>` occurrences are all command-NAMING (shared-core 77/185 `new <slug>` signatures), `→`-binding lines (fo-dispatch-core 98, shared-core 154), or `--help`/diagnostic (claude-fo-runtime 35) — none is an imperative call site. Also checked the flagless `spacedock state commit/ready/sweep` refs: all live inside `## «state.X»` prose-function definition blocks (`- **effect:**` / `- → shipped:`), invoked at call sites via the `«state.X»` indirection, not typed bare — correctly out of the post-gate-invocation defect class.
- DONE: Verify the lint `--name` fix is RED-first and non-vacuous; `→`-binding twin stays exempt.
  AC-2 lint + discriminator GREEN at this tip. Proved the two new controls are non-vacuous by mutation on a throwaway copy: removing `name` from the allowlist FAILS the `leakName` control (line-132 class not flagged); breaking the `- → ` exemption FAILS the `bindingTwin` control (line-98 wrongly flagged). Both restore green. The "line-132-flagged / line-98-exempt" lock is real.
- DONE: Re-run the DETACHED adversarial audit on a fresh throwaway clone of `029a4b14`; `go test ./...` green; re-confirm the value-gate byte table.
  Audit on `/tmp/launcher-invariant-audit-cycle2` (fresh clone, never the worktree, removed after). Edit 1: regressing line 132 back to bare now correctly REDs the lint at `claude-fo-dispatch.md:132` — the escape no longer survives. Edit 2 (generic planted bare helper) and Edit 3 (drift `LauncherCommand()` to bare PATH) both RED and restore green. `go test ./...` exit 0, all packages ok. Byte table: all 10 tracked files unchanged from cycle 1, each ≤ v0.22.0 baseline (shared-core 28583/28586). Untracked `claude-fo-dispatch.md` is 22997 (+18 from the line-132 conversion); not in the gate table, and the +18 is the necessary cost of the fix.

### Summary (validation cycle 2)

Recommendation: **PASSED**. The Cycle-1 material finding is fully closed, verified independently at tip `029a4b14`. Both halves landed exactly as prescribed: line 132 now resolves the pinned launcher, and `name` joined the AC-2 lint's invocation-flag allowlist with two paired discriminator controls — both proven non-vacuous by mutation (each fails on the corresponding regression and restores green). My re-run detached audit on a fresh `029a4b14` clone confirms the original escape no longer survives: regressing line 132 to bare now REDs the lint, which it did not before the fix. The full re-scan finds no other genuine post-gate helper invocation — the remaining bare `spacedock` references are command-naming, `→`-bindings, `--help`/diagnostic hints, or `«state.X»` prose-function definitions, all correctly outside the invariant's scope. `go test ./...` is green, the value gate remains met (all 10 tracked files at/under v0.22.0; the converted line lands in untracked `claude-fo-dispatch.md`, +18 B, justified), the contract's line-7 carve-out still agrees with the lint's exemptions (the fix only tightens that agreement by catching the `--name` invocation class), and there is no shipped FO behavior change. AC-1..AC-4 stand as reproduced in Cycle 1.
