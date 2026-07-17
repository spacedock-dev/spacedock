---
title: Document dispatch entity path base
status: validation
source: captain discussion 2026-07-17
started: 2026-07-17T13:36:06Z
completed:
verdict:
score: 0.7
worktree: .worktrees/spacedock-ensign-document-dispatch-entity-path-base
issue:
id: 93y35a0yrwxnxa5qzfghtj4h
---

`spacedock dispatch build` accepts a relative `entity_path`, but resolves it against the process working directory rather than `workflow_dir`. The built-in help says only “Path to the entity file,” the JSON schema supplies no description, and the FO contract uses an undefined `{entity_file_path}` placeholder. In a split-root workflow this ambiguity makes `.spacedock-state/{slug}/index.md` look canonical even when the first officer is running from the repository root and must pass `docs/dev/.spacedock-state/{slug}/index.md`.

## Proposed approach

Document one rule at the command boundary: `entity_path` may be absolute or relative to the caller's current working directory; it is never relative to `workflow_dir`. State that it must identify the project-root/state-checkout entity, never a code-worktree copy. Add a split-root example whose workflow is `docs/dev` and whose entity is `docs/dev/.spacedock-state/example/index.md`.

Keep this slice on the built-in command surface and generated JSON schema. Do not edit `docs/site/reference/command-reference.md` while the adjacent-plugin branch owns that file.

## Acceptance criteria

**AC-1 (VALUE) - An operator can determine the correct `entity_path` base from `spacedock dispatch build --help` without reading source code.**
Verified by: a help golden asserting the cwd-relative rule, the explicit non-`workflow_dir` distinction, and the split-root example.

**AC-2 - `spacedock dispatch build --print-schema` carries the same path semantics.**
Verified by: decoding the emitted schema and asserting the `entity_path` description names the current working directory and excludes workflow-relative interpretation.

**AC-3 - Existing absolute and cwd-relative behavior remains unchanged.**
Verified by: the existing relative-path behavioral test plus focused help/schema tests; no path-resolution implementation change.

## Test plan

Add focused tests before changing the help/schema text. Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.

## Stage Report: implementation

- DONE: Add focused red-then-green help and schema tests that independently assert entity_path is absolute or relative to the caller current working directory, never relative to workflow_dir.
  Red run failed on both missing contracts; green focused run passed help, decoded-schema, and existing relative-path behavior tests.
- DONE: Document the project-root/state-checkout boundary and include the docs/dev/.spacedock-state/example/index.md split-root example without changing path resolution behavior or the overlapping public command-reference page.
  Commit `d2929c6d` changes only built-in dispatch help, emitted schema metadata, and their focused tests; no resolver or public reference changes.
- DONE: Run gofmt -w ./cmd ./internal, go test ./..., and go test ./... -race; append exact evidence and a recommendation to the entity stage report.
  `gofmt -w ./cmd ./internal` completed; `go test ./...` and `go test ./... -race` both exited 0 across all packages.

### Summary

`dispatch build --help` and `--print-schema` now state that `entity_path` is absolute or caller-cwd-relative, never `workflow_dir`-relative, and points to shared project/state rather than a code worktree. The split-root example uses `docs/dev/.spacedock-state/example/index.md`; existing resolution behavior remains unchanged. Recommendation: advance commit `d2929c6d` to fresh validation.
