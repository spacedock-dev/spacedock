---
title: Bump the live-E2E harness node runtime off node 20 (setup-node node-version) + guard it
status: backlog
score: 0.3
group: cleanup
issue:
id: 3g8h3rysfedynccgt3g09nx4
---

A node-20 deprecation distinct from #359 (`pre-cut-audit-cleanups-0199`). #359 bumped the GitHub ACTIONS off node20 to their node24-runtime majors (checkout@v5, setup-go@v6, goreleaser@v7, setup-node@v6, setup-python@v6, upload-artifact node24), met the 2026-06-16 deadline, and added `internal/release/node24_actions_guard_test.go` to pin the action majors. But that guard checks the action `@vN` majors only — NOT the `node-version:` field.

## Problem

`.github/workflows/runtime-live-e2e.yml` still pins `actions/setup-node@v6` with `node-version: "20"` in THREE places (the claude-live, codex-live, and pi-live jobs — currently lines ~133 / ~341 / ~557). This installs node 20 as the JS RUNTIME the Claude Code / Codex / Pi CLIs execute under during the live E2E. Node 20 reached maintenance/EOL (~mid-2026), so setup-node (or the CLI installers) emit a deprecation. This is the JS runtime the CLIs run on, not the action runtime #359 fixed.

## Why it is NOT a blind bump

Changing the harness node-version changes the runtime the Claude/Codex/Pi CLIs execute on. The CLIs may behave differently (or require a minimum) on node 22/24 — exactly the CLI-compat risk the live lane exists to catch. So the bump must be re-validated by a green live Runtime Live E2E, not merged on a static diff. (Deferred past v0.23.0 deliberately: node 20 still works, the cut should not wait on a deprecation warning, and the compat risk wants its own validation cycle.)

## Proposed approach

1. Bump `node-version: "20"` → `"22"` (current LTS; or `"24"` if the CLIs prefer it) in all three runtime-live-e2e.yml jobs.
2. Extend the node24 guard (or add a sibling test) to assert the `setup-node` `node-version:` pin is ≥ a recorded minimum, so it cannot silently regress to 20 — the independent oracle written in the test.
3. Re-validate with a green live Runtime Live E2E on all lanes (claude opus+sonnet, codex, pi) to confirm the CLIs work on the new node.

## Acceptance criteria

- **AC-1** — all three `setup-node node-version` pins in runtime-live-e2e.yml are off 20 (≥22), and a green live Runtime Live E2E proves the Claude/Codex/Pi CLIs still pass every shared scenario on the new node (the CLI-compat proof — a static diff is NOT sufficient).
- **AC-2** — a guard test asserts the harness `node-version` pin is ≥ the recorded minimum (independent oracle; reds on a regression to 20), alongside the existing node24 action-major guard.
- **AC-3 (no-regression)** — `go test ./internal/release/...` + `go build ./...` green; no other workflow regresses its action majors.

## Test plan

Static: the guard test (red-on-regression). Live: a full Runtime Live E2E dispatch on the bumped branch — the load-bearing proof, since the change alters the CLI runtime. The workflow-file push needs the `workflow` OAuth scope.
