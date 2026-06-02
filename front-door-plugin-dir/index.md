---
id: jgc29m3pjb80efvmrc5bkc2n
title: spacedock claude --plugin-dir before `--` (captain dev workflow) + restore the live-e2e net
status: backlog
source: captain (2026-06-02) — CI-E2E live-runtime e2e fails: cobra `spacedock claude` rejects --plugin-dir before `--`; the captain's primary dev workflow and the e2e net both depend on it. 0.19.3 keystone.
started:
completed:
verdict:
score: "0.42"
worktree:
issue:
---

The approval-gated live-runtime e2e (`TestLiveEnsignCycle`, `internal/ensigncycle/live_test.go:91-100`) has been silently broken since the cobra migration (#241): it invokes `spacedock claude --plugin-dir <repo> --skip-contract-check -p … --model … -- <task>` with the host flags BEFORE `--`, but the cobra-migrated front door rejects them (`unknown flag: --plugin-dir`) — host flags must follow `--`. The break stayed hidden because the env-approval gate kept every prior live run "waiting"; it surfaced the first time the job actually ran (this session, after CI approval). The live-e2e net has therefore been dead through the entire 0.19.2 sprint.

This also blocks the captain's PRIMARY dev workflow: nearly every real launch is `claude --plugin-dir /…/spacedock …`. For `spacedock claude` to truly REPLACE the hand-typed launcher (sprint goal #1), `spacedock claude --plugin-dir <dir>` (before `--`) must work — load the local plugin checkout and relax the contract gate — not require the `-- --plugin-dir` form.

## Scope (captain: "test + binary ergonomics")

- **Binary:** `spacedock claude` accepts `--plugin-dir <dir>` / `--plugin-dir=<dir>` BEFORE `--`, forwarding it into the host passthrough and relaxing the contract gate (the existing `hasPluginDir` gate-relax in `frontdoor.go` must see it). Multiple `--plugin-dir` allowed. Ideation decides whether `spacedock codex` gets the symmetric treatment and how broadly other host flags before `--` are forwarded (vs the documented `--` convention).
- **Test:** `TestLiveEnsignCycle` passes green against the fixed binary, restoring the live-e2e net.
- Lands on `next` as its own PR (whose own CI-E2E validates the fix end-to-end), then greens the in-flight #258 (5p) and #259 (#251) once they rebase onto the fixed `next`.

## Acceptance criteria (provisional — harden at ideation)

**AC-1 — `spacedock claude --plugin-dir <dir>` (before `--`) launches with the local checkout and a relaxed gate.** The flag is accepted (no "unknown flag"), forwarded to the host, and `hasPluginDir` relaxes the contract gate; `--plugin-dir=<dir>` and repeated `--plugin-dir` both work.
Verified by: a frontdoor unit test asserting the parsed passthrough carries `--plugin-dir <dir>` and the gate is relaxed, for both before-`--` and after-`--` placements.

**AC-2 — the live-e2e net is restored.** `go test -tags live -run TestLiveEnsignCycle ./internal/ensigncycle/` passes against the fixed binary (the same invocation CI runs).
Verified by: the live job green on this entity's own PR (CI-E2E), plus a local `-tags live` run recorded in the stage report.

**AC-3 — after-`--` host passthrough is unchanged.** Flags placed after `--` still forward verbatim; spacedock's own flags (`--safehouse-*`, `--skip-contract-check`) still parse before `--`.
Verified by: existing frontdoor passthrough tests stay green; `go test ./...` green.
