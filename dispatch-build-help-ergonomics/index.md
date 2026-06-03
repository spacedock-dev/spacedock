---
id: bxav7rd90n4mcw8fhp9myv19
title: Make dispatch subcommand help reachable without required operational flags
status: backlog
source: "codex FO dogfood (2026-06-03) - `spacedock dispatch build --help` returned a required `--workflow-dir` error while trying to inspect the command contract"
score: "0.22"
worktree:
started:
completed:
verdict:
issue:
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

## Proposed approach

1. Teach `spacedock dispatch` to recognize `--help`, `-h`, or `help` for its subcommands before checking required execution flags.
2. Add focused usage text for `dispatch build` and `dispatch show-stage-def` that names required flags, stdin shape, and examples.
3. Preserve current loud errors for real invocations that omit required flags.

## Acceptance criteria

**AC-1 - Dispatch subcommand help is reachable without execution flags.**
Verified by: Go tests over the CLI runner assert `spacedock dispatch build --help`, `spacedock dispatch build -h`, and `spacedock dispatch show-stage-def --help` exit 0, write usage to stdout, and do not emit required-flag errors.

**AC-2 - Missing required flags still fail for real execution.**
Verified by: Go tests assert `spacedock dispatch build` still exits non-zero with `requires --workflow-dir`, and `spacedock dispatch show-stage-def` still exits non-zero with `requires --workflow-dir and --stage`.

**AC-3 - Help text explains the JSON stdin contract without forcing operators to infer it from tests.**
Verified by: a golden or substring test over `dispatch build --help` asserts the output mentions `--workflow-dir`, stdin JSON, `schema_version`, `entity_path`, `stage`, and `checklist`.

## Test plan

- Add CLI or dispatch package tests for the help-before-required-flags cases.
- Add regression tests for the existing missing-flag errors.
- Run `go test ./internal/dispatch ./internal/cli` and the normal repo gate before merge.
