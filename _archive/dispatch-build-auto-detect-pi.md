---
title: spacedock dispatch build auto-detects host pi when running under Pi agent harness
status: done
score: 0.85
id: 769mybp649pj160n17x13r8g
worktree: .worktrees/spacedock-ensign-dispatch-build-auto-detect-pi
mod-block:
pr: pr-merge:550
verdict: passed
completed: 2026-07-22T14:38:16Z
archived: 2026-07-22T14:38:16Z
---

## Problem

When running `spacedock dispatch build` inside a Pi coding agent session without passing `--host pi` explicitly, the helper currently exits 1 with:

```text
error: missing host source: pass --host, set JSON host, or run under CODEX_THREAD_ID or CLAUDECODE
```

This breaks host-autodetect parity with Claude Code (`CLAUDECODE`) and Codex (`CODEX_THREAD_ID`). Pi sessions expose Pi-specific environment markers, so First Officer dispatches launched from Pi should not need a bespoke explicit `--host pi` flag.

## Research and Spike

- Live Pi harness marker observed in this session: `PI_CODING_AGENT=true`; `PI_CODING_AGENT_DIR` and `PI_CODING_AGENT_SESSION_DIR` are empty in this particular child shell. Repository/runtime docs and Pi runtime code also use `PI_CODING_AGENT_DIR` and `PI_CODING_AGENT_SESSION_DIR` as supported Pi agent/session locations.
- Current host resolver in `internal/dispatch/build.go` only checks `CODEX_THREAD_ID` and `CLAUDECODE`, then emits the missing-host error. Its explicit-source precedence is already flag host -> JSON host -> environment.
- Spike result: with `PI_CODING_AGENT=true` and no Claude/Codex markers, `${SPACEDOCK_BIN:-spacedock} dispatch build ...` exits 1 with the missing-host error. The same command with `--host pi` exits 0 and emits a Pi-shaped dispatch JSON, so the unverified risk is isolated to environment host resolution, not the Pi adapter itself.

## Proposed Approach

1. **Extend host auto-detection in `resolveBuildHost`:** After explicit `--host` / JSON host handling, derive `pi := getenv("PI_CODING_AGENT") != "" || getenv("PI_CODING_AGENT_DIR") != ""`. Include Pi in the same ambiguity check as Codex and Claude: if more than one runtime marker family is set, return an explicit ambiguous-runtime error naming the set markers and telling the user to pass `--host claude`, `--host codex`, or `--host pi`.
   - Value AC served: AC-1, AC-2, AC-4.
   - Simplest alternative considered: detect only `PI_CODING_AGENT_DIR`. Insufficient because the live Pi child shell for this task has `PI_CODING_AGENT=true` while `PI_CODING_AGENT_DIR` is empty, so it would still fail in the target runtime.
   - Simplest alternative considered: silently prefer Pi over Claude/Codex when multiple markers are set. Insufficient because current behavior rejects ambiguous Claude+Codex sources; preserving that safety avoids wrong-host dispatch bodies in nested or inherited environments.

2. **Update host-resolution tests in `internal/dispatch/build_json_ergonomics_test.go`:** Add table rows for `PI_CODING_AGENT=true`, `PI_CODING_AGENT_DIR=<tmp>`, explicit `--host` overriding Pi markers, JSON host overriding Pi markers, and Pi+Claude/Codex ambiguity.
   - Value AC served: AC-1 through AC-4.
   - Simplest alternative considered: test only `resolveBuildHost` directly. Insufficient because the user-facing failure occurs through `dispatch build`; existing tests already exercise the real command path and inspect the emitted dispatch shape.

3. **Keep documentation changes minimal:** This is user-visible CLI behavior. Update the `dispatch build` help/environment text (or the closest existing help sentence for host derivation) so the missing-host/remediation wording includes Pi markers.
   - Proposed wording change:

```diff
- error: missing host source: pass --host, set JSON host, or run under CODEX_THREAD_ID or CLAUDECODE
+ error: missing host source: pass --host, set JSON host, or run under CODEX_THREAD_ID, CLAUDECODE, PI_CODING_AGENT, or PI_CODING_AGENT_DIR
```

```diff
- ambiguous runtime host sources: CODEX_THREAD_ID and CLAUDECODE are both set; pass --host claude, codex, or pi
+ ambiguous runtime host sources: multiple runtime markers are set; pass --host claude, codex, or pi
```

## Expected Surface

- `internal/dispatch/build.go`: small host resolver change and error wording, about 10-25 LOC changed.
- `internal/dispatch/build_json_ergonomics_test.go`: host-resolution table/rows, about 35-70 LOC changed.
- Optional if help text names env-derived hosts: `internal/dispatch/dispatch.go`, about 1-5 LOC changed.
- Tolerance: total implementation delta should stay under ~120 LOC and should not touch skill text, Pi launch/install flows, status behavior, or runtime-auth code.

## Acceptance Criteria

- **AC-1 (Pi marker auto-detects Pi dispatch):** With `PI_CODING_AGENT=true` set and no `--host`, JSON host, `CODEX_THREAD_ID`, or `CLAUDECODE`, `spacedock dispatch build` exits 0 and emits a Pi dispatch body/envelope rather than the missing-host error. *Verified by:* a command-level Go test using `runNativePreservingHostEnv` that asserts exit 0 and Pi-shaped output; it would fail if `PI_CODING_AGENT` were ignored.
- **AC-2 (Pi directory marker also auto-detects Pi):** With `PI_CODING_AGENT_DIR` set and no other host source, `spacedock dispatch build` exits 0 and emits a Pi dispatch body/envelope. *Verified by:* a command-level Go test; it would fail if detection only handled the boolean marker.
- **AC-3 (Explicit host sources override Pi env):** When `--host claude`, `--host codex`, or JSON `host` is provided, that explicit source wins even if `PI_CODING_AGENT=true` is set. *Verified by:* command-level Go tests asserting the emitted prompt/body matches the explicit host; they would fail if env detection ran before explicit-source precedence.
- **AC-4 (Ambiguous runtime markers fail safely):** When Pi and Claude/Codex runtime markers are both set and no explicit host source is provided, `dispatch build` exits non-zero with an ambiguity error naming the conflicting marker families and the explicit-host remedy. *Verified by:* command-level Go tests; they would fail if the resolver silently selected the wrong host.
- **AC-5 (Operator-facing remediation names Pi):** The missing-host error/help text includes Pi markers so a Pi FO can understand why host inference did or did not happen. *Verified by:* focused test assertion on the error text or help text; it would fail if remediation still mentions only Codex/Claude.

## Test Plan

- Add focused subtests to `TestBuildHostResolutionFromFlagJSONAndEnv` for `derived-pi-from-PI_CODING_AGENT`, `derived-pi-from-PI_CODING_AGENT_DIR`, `explicit-overrides-pi-runtime`, and `ambiguous-pi-runtime`.
- Keep tests command-level through `Run(...)`/`runNativePreservingHostEnv` so they verify output shape and exit behavior, not only a private helper return value.
- Run `go test ./internal/dispatch -run TestBuildHostResolutionFromFlagJSONAndEnv -count=1` for the focused behavior, then `go test ./...` as the stage baseline. Run `go test ./... -race` and `gofmt -w ./cmd ./internal` before completion because Go code will be touched in implementation.

## Directives

- Perform ideation dispatch for this entity.

## Stage Report: ideation

- DONE: Research environment variables set by Pi coding agent harness (PI_CODING_AGENT, PI_CODING_AGENT_DIR).
  Observed live `PI_CODING_AGENT=true`; repo/runtime docs and Pi code reference `PI_CODING_AGENT_DIR`/session dir as Pi agent paths.
- DONE: Design host auto-detection for dispatch build when running under Pi.
  Proposed extending `resolveBuildHost` after explicit sources, with Pi markers and safe multi-runtime ambiguity handling.
- DONE: Author proposed design, acceptance criteria, test plan, and stage report.
  Entity now includes problem, research/spike, expected surface, ACs with verification, test plan, and this report.

### Summary

Ideation scoped the task to a small dispatch host-resolution change plus focused command-level tests. The live spike confirms Pi dispatch works with explicit `--host pi` and the current failure is specifically missing Pi env auto-detection.

## Stage Report: implementation

- DONE: Implement PI_CODING_AGENT / PI_CODING_AGENT_DIR auto-detection for host: pi in resolveBuildHost in internal/cli/dispatch_build.go.
  Implemented in `internal/dispatch/build.go` (actual package path): Pi markers now resolve to `pi`, explicit sources still win, and mixed runtime families fail with an ambiguity error; code commit `2fed964b`.
- DONE: Add unit test coverage in internal/cli/dispatch_build_test.go covering pi auto-resolution and explicit host override.
  Added command-level coverage in `internal/dispatch/build_json_ergonomics_test.go` (actual package path); tests fail if `PI_CODING_AGENT`/`PI_CODING_AGENT_DIR` are ignored or explicit flag/JSON host stops overriding Pi env.
- DONE: Run go test ./... and go test ./... -race to verify all tests pass cleanly.
  `go test ./...` and `go test ./... -race` both passed after `gofmt -w ./cmd ./internal`.

### Summary

`dispatch build` now auto-detects Pi when either `PI_CODING_AGENT` or `PI_CODING_AGENT_DIR` is present, while retaining safe ambiguity behavior across distinct runtime marker families. The implementation also clears Pi markers in legacy default-Claude test harness paths so the suite remains deterministic when run from a live Pi shell.

## Stage Report: validation

- DONE: Run go test ./... and go test ./... -race in worktree .worktrees/spacedock-ensign-dispatch-build-auto-detect-pi.
  `go test ./...` passed; `go test ./... -race` passed; `go test ./internal/dispatch -run TestBuildHostResolutionFromFlagJSONAndEnv -count=1` passed and would fail if dispatch host resolution emitted the wrong host shape or exit code.
- DONE: Verify dispatch_build_test.go covers both automatic host resolution and explicit overrides.
  Coverage lives in `internal/dispatch/build_json_ergonomics_test.go` (actual dispatch build test file): subtests cover `PI_CODING_AGENT`, `PI_CODING_AGENT_DIR`, flag override, JSON override, missing-source text, and Pi+Codex ambiguity.
- DONE: Record validation evidence in stage report.
  This report records focused/full/race test evidence plus AC-specific semantic checks and the validation recommendation.

### Acceptance Criteria Validation

- DONE: AC-1 (Pi marker auto-detects Pi dispatch).
  `derived-pi-from-PI_CODING_AGENT` runs command-level `dispatch build` with only `PI_CODING_AGENT=true`, asserts exit 0 and Pi read-dispatch body; it would fail if the boolean marker were ignored.
- DONE: AC-2 (Pi directory marker also auto-detects Pi).
  `derived-pi-from-PI_CODING_AGENT_DIR` runs command-level `dispatch build` with only a temp `PI_CODING_AGENT_DIR`, asserts exit 0 and Pi body; it would fail if only `PI_CODING_AGENT` were recognized.
- DONE: AC-3 (Explicit host sources override Pi env).
  `host-flag-overrides-pi-runtime` and `json-host-overrides-pi-runtime` assert Claude/Codex output under `PI_CODING_AGENT=true`; they would fail if environment host detection ran before explicit sources.
- DONE: AC-4 (Ambiguous runtime markers fail safely).
  `ambiguous-pi-runtime` asserts non-zero command-level failure naming `CODEX_THREAD_ID`, `PI_CODING_AGENT`, and `--host claude, codex, or pi`; it would fail if the resolver silently selected a host.
- DONE: AC-5 (Operator-facing remediation names Pi).
  `missing-source` asserts the missing-host remediation includes `PI_CODING_AGENT` and `PI_CODING_AGENT_DIR`; it would fail if the old Codex/Claude-only message returned.

### Summary

Validation PASSED. The focused, full, and race suites pass, and the AC evidence is command-level output/exit behavior rather than static prose; semantic review found explicit-source precedence is preserved before environment detection and same-host Pi markers do not create false ambiguity.
