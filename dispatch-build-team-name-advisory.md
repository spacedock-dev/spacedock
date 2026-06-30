---
title: "dispatch build — warn + document --team-name on claude (auto-team is default; legacy shape is silent)"
status: ideation
score: 0.3
source: "v0.23.0 cut FO session, 2026-06-30. FO reflexively passed --team-name (from the boot team_state hint) to dispatch build on Claude, which silently emitted the LEGACY team-registry shape (team_name present, run_in_background absent) instead of the auto-team merged-mode shape. Self-corrected by rebuilding without --team-name; no bad dispatch shipped."
id: 0qt2r4n577gtwq707abgcfed
started: 2026-06-30T16:54:21Z
---

`spacedock dispatch build` IS Claude-team-mode-aware (`internal/dispatch/build.go:292`: `mergedMode := !bareMode && host == "claude" && teamName == ""`): omitting `--team-name` yields the auto-team shape (`run_in_background:true`, no `team_name`); passing `--team-name` selects the legacy TeamCreate-registry shape. The foot-gun is ergonomic, not correctness.

## Problem
1. `--team-name` / `--bare-mode` are NOT documented in `dispatch build --help` (it only shows the stdin-JSON interface), so the legacy opt-in is invisible at the CLI surface.
2. Nothing warns when `--team-name` is passed on `host=claude`, where auto-team is the default and the legacy path is sunsetting — so a stray `--team-name` silently produces a shape (team_name present, run_in_background absent) that is verbatim-unsafe for auto-team mode. The binary cannot self-detect legacy-vs-auto-team (that is FO-probed at runtime via TeamCreate/SendMessage availability, invisible to the CLI), so `--team-name` must stay the explicit signal — but a guard would catch the slip.

## Proposed approach
Do BOTH halves — they address the two distinct problems above and are each cheap. Neither changes the emitted dispatch shape.

**(a) stderr advisory (Problem #2).** Add `claudeteam.LegacyTeamNameAdvisory(stderr)` alongside the existing `BareModeAdvisory` in `internal/claudeteam/claudeteam.go`. The text names Claude-only concepts (auto-team, the claude default), so it lives in the Claude seam, matching `BareModeAdvisory`. Call it from `runBuildFields` immediately after the `mergedMode` line (`internal/dispatch/build.go:292`), next to the bare-mode advisory block, gated on:

    !bareMode && host == "claude" && teamName != ""

This is the exact `teamName != ""` complement of `mergedMode` within the non-bare claude branch — the SAME three signals `mergedMode` already derives, no new detection. Unlike the bare-mode advisory it needs NO probe: the trigger is wholly in the CLI args (`--team-name` was explicitly passed), so there is no `~/.claude` read. It writes only to `stderr`; it does not touch `out` (`buildOutput`) or any field feeding stdout, so the merged-vs-legacy envelope selection (build.go:679-691) is untouched.

Exact advisory text — one logical line, mirroring `BareModeAdvisory`'s WARN / what / why / corrective / silence shape:

    WARN: --team-name selects the legacy TeamCreate-registry dispatch shape (team_name present, run_in_background absent). On host=claude, auto-team is the default — omit --team-name to emit the auto-team shape (name + run_in_background, no team_name). If you mean the legacy team-registry path, this warning can be ignored.

**(b) `--help` docs (Problem #1).** Extend `printBuildUsage` (`internal/dispatch/dispatch.go:282`) so the Flags block documents the shape-selecting flags. `docs/site/reference/command-reference.md` already defers to `--help` as the source of truth for flags, so the `--help` text IS the doc surface — no separate prose doc to edit.

Out of scope (recorded for the roadmap, not this task): once legacy team mode is fully sunset, rejecting `--team-name` on claude outright — that is a behavior change, not this ergonomic guard.

### Doc diff — `printBuildUsage` Flags block (dispatch.go ~:288)

Before:

    Flags:
      --workflow-dir DIR   Workflow definition directory containing README.md.

After:

    Flags:
      --workflow-dir DIR   Workflow definition directory containing README.md.
      --host HOST          Override the runtime host (claude|codex|pi). Defaults to the detected runtime.
      --team-name NAME     Select the legacy TeamCreate-registry dispatch shape. On host=claude, auto-team is the default — omit this unless you mean legacy team mode.
      --bare-mode          Emit the bare sequential shape (no name, no team_name, no run_in_background).

The advisory text and these descriptions stay aligned: both name the legacy-vs-auto-team distinction in the same words.

## Acceptance criteria
- **AC-1 (end-value, behavioral)** — On `host=claude`, a `dispatch build` invocation that passes `--team-name` emits exactly one stderr advisory line naming the legacy TeamCreate-registry shape and the auto-team default; the SAME invocation WITHOUT `--team-name` (the merged shape) emits NO such advisory. The advisory's presence/absence is the foot-gun signal that today is always absent — the baseline that can move the wrong way: a false negative reverts to today's silence; a false positive fires the warning on every merged build. Verified by: a Go test that runs both arms over one fixture and asserts the advisory marker present on the `--team-name` arm and absent on the merged arm, by inspecting captured stderr — never a prose-grep over source.
- **AC-2 (help docs)** — `dispatch build --help` stdout documents `--team-name` and `--bare-mode` (and `--host`), each with a one-line description, with `--team-name`'s line naming the legacy shape and the auto-team default. Verified by: a `--help` golden test that invokes `dispatch build --help`, captures stdout, and asserts it contains `--team-name` and `--bare-mode` with their descriptions.
- **AC-3 (negative control — envelope unchanged)** — Adding the advisory does NOT change the emitted dispatch envelope. The legacy `--team-name` claude build's stdout JSON keeps its pre-advisory shape: `name` present, `team_name` present, `run_in_background` ABSENT (the advisory does not flip the legacy shape to merged, drop `team_name`, or add `run_in_background`). Because the advisory's trigger (`--team-name` on claude) is the same signal that selects the legacy shape, a same-inputs "advisory off, legacy envelope" arm is impossible by construction; the negative control is therefore (i) the existing `TestBuildLegacyModeUnchanged` (build_merged_mode_test.go:103) staying green — its envelope-shape assertions are the frozen baseline — and (ii) the AC-1 advisory test re-asserting the legacy envelope shape on the `--team-name` arm, so a regression that altered the envelope while adding the advisory fails.

## Test plan
All tests are Go unit tests in `internal/dispatch` (advisory text in `internal/claudeteam`), driven through the in-process `Run(...)` surface via the existing `runNative` / `runNativeWithDefaultClaudeHost` harness — no live workflow, no new fixtures beyond the existing `writeGood`/`writeFlatEntity` helpers. Cost: minutes; no new infra.

1. **AC-1 advisory gate (new)** — `TestBuildLegacyTeamNameAdvisory`, modeled on `TestBuildBareAdvisoryProbeGate` (build_advisory_probe_test.go). One claude fixture, two arms over the same inputs:
   - arm A — `team_name` present → assert stderr CONTAINS a stable advisory marker (fragment `"legacy TeamCreate-registry dispatch shape"`); decode stdout → assert `team_name` present, `run_in_background` absent (legacy shape; this is AC-3 in the same arm).
   - arm B — `team_name` absent (merged) → assert stderr does NOT contain the marker; decode stdout → assert `run_in_background` present, `team_name` absent.
   The marker is a fragment, not the full sentence, so a future wording tweak does not break the gate (matches the `advisoryMarker` convention at build_advisory_probe_test.go:53).
2. **AC-2 help golden (new)** — `TestBuildHelpDocumentsShapeFlags`: invoke `Run(probe, []string{"build", "--help"}, ...)`, capture stdout, assert it contains `--team-name` and `--bare-mode` and the legacy/auto-team phrasing. Behavioral (runs the command, reads stdout), not a source grep.
3. **AC-3 negative control (existing, must stay green)** — `TestBuildLegacyModeUnchanged` already asserts the legacy envelope shape (name + team_name present, run_in_background absent). It is the frozen pre-advisory envelope baseline; a stderr-only advisory must leave it green. No change needed beyond confirming it still passes after the implementation.

## Spike
No spike needed — every mechanism the design rests on is already shipped and proven by existing tests:
- The `!bareMode && host == "claude" && teamName != ""` signal: `mergedMode` derives the exact same three signals at build.go:292; `--team-name`/`--bare-mode`/`--host` are already parsed in `parseBuildOptions` (dispatch.go:126).
- The stderr-only advisory pattern: `claudeteam.BareModeAdvisory` (claudeteam.go:73) writes to a stderr writer separate from stdout in `Run`, and `TestBuildBareAdvisoryProbeGate` proves the envelope is byte-identical across an advisory-present/absent toggle.
- The legacy envelope shape is frozen by `TestBuildLegacyModeUnchanged`; the `--help` dispatch by `wantsHelp` → `printBuildUsage` (dispatch.go:28,282).

## Stage Report: ideation

- DONE: Decide the fix shape — stderr advisory + `--help` docs — with the EXACT advisory wording and exact `--help` additions; confirm the advisory fires off the existing mergedMode detection WITHOUT altering the emitted dispatch shape.
  Proposed approach (a)+(b): exact `WARN:` advisory text and the exact Flags-block before/after recorded in body. Advisory gate is `!bareMode && host == "claude" && teamName != ""` — the exact `teamName != ""` complement of `mergedMode` (build.go:292), stderr-only, no probe, `out`/stdout untouched.
- DONE: An end-value AC that measures the usability fix BEHAVIORALLY — advisory fires on a claude-host `--team-name` build and is silent when omitted, plus a `--help` golden including the flags, verified by running the command and observing stdout/stderr.
  AC-1 (two-arm advisory presence/absence over one fixture, modeled on TestBuildBareAdvisoryProbeGate) + AC-2 (`--help` golden asserting `--team-name`/`--bare-mode`). Test plan names both new tests and the harness.
- DONE: A negative control that the emitted dispatch ENVELOPE bytes are unchanged with/without the advisory.
  AC-3: existing TestBuildLegacyModeUnchanged (build_merged_mode_test.go:103) stays green as the frozen envelope baseline + AC-1 re-asserts the legacy shape on the advisory arm. Body records why a same-inputs "advisory off, legacy envelope" arm is impossible by construction (the trigger IS the shape selector).

### Summary

Refined the ideation body to a complete, grounded design: do BOTH halves (stderr advisory + `--help` docs), with exact wording for each. The advisory text lives in the Claude seam beside `BareModeAdvisory`, fires from build.go right after the existing `mergedMode` line on the `teamName != ""` complement (no new detection, no probe, stderr-only), and the `--help` Flags block gains `--host`/`--team-name`/`--bare-mode`. Three ACs: AC-1 behavioral advisory gate (end-value), AC-2 `--help` golden, AC-3 envelope-unchanged negative control. Recorded "no spike needed" — `mergedMode`, the `BareModeAdvisory` stderr pattern, the frozen legacy-envelope test, and `printBuildUsage` are all already shipped and test-proven; deferred outright-rejection of `--team-name` to the roadmap.
