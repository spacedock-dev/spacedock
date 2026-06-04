---
id: bxav7rd90n4mcw8fhp9myv19
title: Make dispatch subcommand help reachable without required operational flags
status: done
source: "codex FO dogfood (2026-06-03) - `spacedock dispatch build --help` returned a required `--workflow-dir` error while trying to inspect the command contract"
score: "0.22"
worktree:
started: 2026-06-03T21:34:29Z
completed: 2026-06-04T01:14:25Z
verdict: PASSED
issue:
mod-block:
pr: "#286"
---

The dispatch command's subcommand help is blocked by the same required-flag checks used for real execution. During first-officer dogfooding, trying to inspect `dispatch build` before constructing the JSON envelope failed with the operational error instead of showing help.

## Problem

The current routing requires execution flags before recognizing help intent:

```text
$ go run ./cmd/spacedock dispatch build --help
error: dispatch build requires --workflow-dir
exit status 2

$ go run ./cmd/spacedock dispatch show-stage-def --help
error: dispatch show-stage-def requires --workflow-dir and --stage
exit status 2
```

That is backwards for an operator or first officer trying to learn the command. Help should be discoverable before the command has enough data to run. The execution path should still reject missing required flags when no help flag is present.

This is distinct from `dispatch-build-json-ergonomics`: that entity addresses the hand-built JSON payload. This one is only about help/usage routing.

## Spike

The riskiest unknown is whether this is a cobra-front-door issue or the dispatch router's own validation order. It is the dispatch router.

Reproduced on current HEAD:

```text
$ go run ./cmd/spacedock dispatch build --help
error: dispatch build requires --workflow-dir
exit status 2

$ go run ./cmd/spacedock dispatch show-stage-def --help
error: dispatch show-stage-def requires --workflow-dir and --stage
exit status 2
```

Code read: `internal/cli/cli.go` disables cobra flag parsing for `dispatch` and forwards the post-`dispatch` argv verbatim to `dispatch.Run`. In `internal/dispatch/dispatch.go`, the `build` branch calls `requireFlag` before it checks for help, and the `show-stage-def` branch calls `requireStageFlags` before it checks for help. The implementation should therefore live in the dispatch package: detect `-h`/`--help` for the recognized subcommand before required execution-flag validation, print subcommand help to stdout, and return 0.

## Proposed approach

1. Teach the `build` and `show-stage-def` branches in `internal/dispatch.Run` to recognize `--help` and `-h` before checking required execution flags.
2. Add focused stdout help for `dispatch build` and `dispatch show-stage-def`.
3. Keep the existing execution path unchanged when no help flag is present, including the current required-flag errors and exit code.
4. Keep top-level dispatch routing and the other dispatch subcommands out of scope unless the implementation needs a tiny shared helper for help detection.

## Usage expectations

`dispatch build --help` should show:

- `Usage: spacedock dispatch build --workflow-dir DIR` with stdin JSON called out explicitly.
- The `--workflow-dir` flag and what it points at.
- The stdin JSON contract fields used by the current command: `schema_version`, `entity_path`, `workflow_dir`, `stage`, and `checklist`.
- A minimal JSON example or field list that makes clear `checklist` is an array of checklist strings.

`dispatch show-stage-def --help` should show:

- `Usage: spacedock dispatch show-stage-def --workflow-dir DIR --stage STAGE`.
- `--workflow-dir` and `--stage` descriptions.

## Out of scope

- New input forms for `dispatch build`.
- JSON schema printing or validation-only modes.
- Host-default behavior.
- FO runtime-document examples.
- Any PR, mod, or broader dispatch ergonomics work. Coordinate those with `dispatch-build-json-ergonomics` instead.

## Acceptance criteria

**AC-1 - Dispatch subcommand help is reachable without execution flags.**
Verified by: Go tests assert `spacedock dispatch build --help`, `spacedock dispatch build -h`, `spacedock dispatch show-stage-def --help`, and `spacedock dispatch show-stage-def -h` exit 0, write usage to stdout, leave stderr empty, and do not emit required-flag errors.

**AC-2 - Missing required flags still fail for real execution.**
Verified by: Go tests assert `spacedock dispatch build` still exits 2 with `requires --workflow-dir`, and `spacedock dispatch show-stage-def` still exits 2 with `requires --workflow-dir and --stage` when no help flag is present.

**AC-3 - Help text explains the JSON stdin contract without forcing operators to infer it from tests.**
Verified by: a golden or substring test over `dispatch build --help` asserts the output mentions `--workflow-dir`, stdin JSON, `schema_version`, `entity_path`, `stage`, and `checklist`.

**AC-4 - The scope stays help-only.**
Verified by: the implementation changes no `dispatch build` JSON ingest semantics and no `show-stage-def` stage extraction behavior; existing dispatch golden/error tests continue to pass unchanged except for newly added help fixtures/tests.

## Test plan

- Add a focused `internal/dispatch` test for help-before-required-flags: `build --help`, `build -h`, `show-stage-def --help`, and `show-stage-def -h` return 0, stdout contains `Usage:`, stderr is empty, and the required-flag diagnostics are absent. Cost: low.
- Add or extend a routing regression test for real missing-flag execution: `build` and `show-stage-def` without flags still return exit 2 and the existing diagnostics. Cost: low.
- Add a golden or substring assertion for `build --help` covering `--workflow-dir`, stdin JSON, `schema_version`, `entity_path`, `stage`, and `checklist`; include `show-stage-def` help assertions for `--workflow-dir` and `--stage`. Cost: low.
- Run `go test ./internal/dispatch ./internal/cli` first, then the repo gate required by `AGENTS.md` before implementation completion.

## Stage Report: ideation

- DONE: Help-before-required-flags behavior is scoped separately from real missing-flag failures.
  Reproduced the current failure and located the validation-order issue in `internal/dispatch.Run`; the proposed fix handles help before required flags only for recognized dispatch subcommands, while AC-2 preserves current real-invocation failures.
- DONE: Usage expectations cover `build`, `show-stage-def`, and the JSON stdin contract fields.
  Added explicit usage expectations for both subcommands and the `build` stdin JSON fields: `schema_version`, `entity_path`, `workflow_dir`, `stage`, and `checklist`.
- DONE: Acceptance criteria and test plan are focused, failable, and low-risk to implement first.
  ACs now assert observable exit codes, stdout/stderr placement, required diagnostic absence/presence, and a help-only scope guard. The test plan starts with focused dispatch tests before the broader repo gate.

### Summary

Completed ideation for a narrow help-routing fix. The current failure is a dispatch-router validation-order problem, not a cobra parsing problem: `build` and `show-stage-def` validate required execution flags before recognizing help. The next stage should add subcommand help detection in `internal/dispatch.Run`, print focused stdout help for `build` and `show-stage-def`, and preserve the existing loud missing-flag errors for real invocations.

## Stage Report: implementation

- DONE: `build` and `show-stage-def` help flags exit 0 with stdout usage before required-flag validation.
  Commit `9f9a7ef0` adds help-first routing; `go run ./cmd/spacedock dispatch build --help` and `show-stage-def --help` printed stdout usage.
- DONE: Real missing-flag executions still exit 2 with the existing required-flag diagnostics.
  `TestRequiredFlagGuard` now pins exit 2, empty stdout, and the existing required-flag messages for real invocations.
- DONE: Focused tests cover the build JSON-field help text and the help-only/no-behavior-drift scope.
  `internal/dispatch/help_test.go` covers `-h`/`--help`, JSON field substrings, and required-diagnostic absence; `go test ./internal/dispatch ./internal/cli` passed.

### Summary

Implemented the help-only dispatch router change in `internal/dispatch`: `build` and `show-stage-def` now detect `-h`/`--help` before required operational flag validation and print focused stdout usage. Added focused help tests and tightened the missing-flag guard; `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` all completed successfully.

## Stage Report: validation

- DONE: AC-1 and AC-3 help output behavior is verified with command or test evidence.
  Recommendation: PASSED. `go test ./internal/dispatch ./internal/cli` passed 324 tests across 2 packages, including the new `TestDispatchBuildHelpBeforeRequiredFlags` and `TestDispatchShowStageDefHelpBeforeRequiredFlags` coverage for `--help` and `-h`. Direct command checks also passed: `go run ./cmd/spacedock dispatch build --help`, `go run ./cmd/spacedock dispatch build -h`, `go run ./cmd/spacedock dispatch show-stage-def --help`, and `go run ./cmd/spacedock dispatch show-stage-def -h` all exited 0 and printed usage. A temporary built binary confirmed `dispatch build --help` and `dispatch show-stage-def --help` each exited 0, wrote usage to stdout, and wrote 0 stderr bytes. The `build --help` output includes `--workflow-dir`, stdin JSON, `schema_version`, `entity_path`, `workflow_dir`, `stage`, and `checklist`.
- DONE: AC-2 missing required flags still fail with the existing diagnostics.
  Recommendation: PASSED. The temporary built binary verified exact real-command failures: `dispatch build` exited 2 with empty stdout and `error: dispatch build requires --workflow-dir`; `dispatch show-stage-def` exited 2 with empty stdout and `error: dispatch show-stage-def requires --workflow-dir and --stage`. `TestRequiredFlagGuard` also passed in the focused package run.
- DONE: AC-4 help-only scope is checked against changed files and existing dispatch behavior tests.
  Recommendation: PASSED. `git diff --name-status origin/next...HEAD` shows only `internal/dispatch/dispatch.go`, `internal/dispatch/guard_test.go`, and `internal/dispatch/help_test.go` changed. The production delta is help detection and help text in the dispatch router; `internal/dispatch/build.go` and `internal/dispatch/showstagedef.go` are unchanged, so JSON ingest semantics and stage extraction behavior are not modified. Existing dispatch and CLI coverage passed via `go test ./internal/dispatch ./internal/cli`; repo-wide gates also passed via `go test ./...` and `go test ./... -race`, each reporting 885 tests across 12 packages. `gofmt -l ./cmd ./internal` produced no output.

### Summary

Validation PASSED for code branch HEAD `9f9a7ef0`. All acceptance criteria are satisfied by focused tests, direct command evidence, exact missing-flag exit checks, changed-file scope review, and full repo gates.
