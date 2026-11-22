---
id: 936b86qpa8xy4fb6v0n5cafs
title: Opt-in Paseo dispatch backend (--host paseo, CLI-first FO binding, and live smoke)
status: ideation
source: captain (2026-11-21)
started: 2026-11-22T17:09:22Z
completed:
verdict:
score: 0.90
sprint: paseo-backend
worktree:
issue:
pr:
mod-block:
sprint-readiness: ready
gates:
    version: 1
    records:
        - id: gate:936b86qpa8xy4fb6v0n5cafs:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:936b86qpa8xy4fb6v0n5cafs-backlog-1
              briefing:
                id: briefing:936b86qpa8xy4fb6v0n5cafs:backlog:attempt-1:revision-1
                digest: sha256:e05a1a0037fb0cc9666082e61950d5175302ac815fcfe5c22f1aca2ca1904608
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:936b86qpa8xy4fb6v0n5cafs:backlog:1
                briefing: briefing:936b86qpa8xy4fb6v0n5cafs:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-11-22T17:09:05.573077Z"
                decision: approve
                reason: Captain approved backlog->ideation dispatch on 6769 with agy neon high
              application:
                target-stage: ideation
                state: consumed
---

Let any Paseo-hosted first officer (Claude, Codex, Pi, or Antigravity) dispatch fresh ensigns through Paseo using one host-neutral `--host paseo` contract and 5 `paseo` CLI verbs, backed by a real isolated-daemon live E2E smoke test and an optional (non-required) CI lane on `codex/gpt-6-luna`.

## Problem

Today Spacedock only dispatches workers through host-specific native subagent tools (Claude `Agent`, Codex `spawn_agent`, Pi `pi-subagents`). Running under Paseo — especially from hosts like Antigravity (`agy`) or across mixed FO/worker runtimes — requires manual prompt patching, hand-relayed reports, or building per-host Go frontdoors and adapters (`27a`, now superseded). Meanwhile, Paseo moves quickly, so any prose-only dispatch binding without a real-daemon E2E test will silently drift when Paseo flags or JSON shapes change.

## Proposed approach

1. **Host-neutral worker contract (`--host paseo`):** Add `paseo` as a supported `--host` in `internal/dispatch/build.go` (host-neutral `ensign` skill load, split-root state commit guidance, concise `Done:` final turn with no host mailbox tool calls) plus a concise `skills/ensign/references/paseo-ensign-runtime.md` wired from `skills/ensign/SKILL.md`. One envelope covers Codex, Antigravity, Claude, and Pi workers launched by Paseo.
2. **CLI-first FO dispatch binding:** Add `skills/first-officer/references/paseo-fo-dispatch.md` using the 5 core `paseo` CLI verbs (`workspace create`, `run --no-wait`, `wait`, `inspect`/`ls`, `archive`), adopting the `--stamp` worktree under the FO's Paseo project ID, labelling with `spacedock.*`, verifying completion from Git state with `notifyOnFinish=false`, and archiving verified workers. Wire `PASEO_AGENT_ID` / explicit Paseo selection in `skills/first-officer/SKILL.md`.
3. **Real-daemon live E2E smoke + optional CI lane:** Add `TestLivePaseoDispatchSmoke` in `internal/ensigncycle/` (`-tags live`). It launches an isolated `paseo` daemon (`SPACEDOCK_PASEO_BIN` or `paseo` on `PATH`, canonical `filepath.EvalSymlinks(t.TempDir())` home, free port, `relay`/`dictation`/`voiceMode` off, env scrubbed, stopped in `t.Cleanup`), runs a real FO and real worker (`SPACEDOCK_PASEO_LIVE_PROVIDER`, defaulting to `codex/gpt-6-luna`, overridable locally to e.g. `antigravity/gemini-4-neon`), and asserts durable Git commits (code worktree + split-root state report) and worker archival. Add a non-required CI workflow lane (mirroring `pi-live`) using `codex/gpt-6-luna` on the existing Codex subscription (`CODEX_AUTH_JSON`).

## Risk evidence

Two live probes on 2026-11-21 proved the underlying path end-to-end:
- `agy` FO -> `paseo` CLI -> `safehouse-codex` ensign completed in 49 s with automatic `ParentAgentId` linkage, code commit (`a9fb598`), split-root state commit (`cf64887`), `paseo wait` idle, and clean archive.
- Sandboxed `antigravity` (`gemini-4-neon`) ensign on Paseo `launch-wrap-hook` (`9f9f39ce1`) + `paseo-safehouse` (`0ff8d55`, `route: wrap`) completed with host-neutral prompt, verified `ps` under `safehouse --enable=keychain`, code commit (`c24eda1`), and split-root state commit (`9b1476c`).
- Remaining risk to verify in implementation: non-`/tmp` worktree sandbox access to the sibling `.spacedock-state` checkout and CI installation/execution of `paseo` CLI with `CODEX_AUTH_JSON`.

## Out of scope

- FO restart adoption and failure/timeout recovery matrix (tracked separately in `paseo-worker-adoption-and-recovery`).
- Remote Paseo workers (`serverId`) and child-to-parent mid-run questions.
- Fake ACP stub servers or Paseo version-pinning machinery.
- Worker reuse across stages.

## Expected surface and tolerance

Estimate net LOC change: +350 to +550, across 8 to 12 files (Go dispatch build + unit tests, FO/ensign references + SKILL.md wiring, live smoke test, optional CI workflow).

## Acceptance criteria

**AC-1 (VALUE) - Real Paseo FO-to-worker dispatch produces verified durable Git state and archives the worker.**
Verified by: `go test -tags live -count=1 -run '^TestLivePaseoDispatchSmoke$' ./internal/ensigncycle -v` passes against an isolated `paseo` daemon (`codex/gpt-6-luna` default; `SPACEDOCK_PASEO_BIN` and `SPACEDOCK_PASEO_LIVE_PROVIDER` overridable), asserting the task branch commit, path-scoped split-root stage report commit, archived worker status, and clean daemon teardown. Breaking any `paseo` CLI flag in `paseo-fo-dispatch.md` or dropping archive cleanup fails the test.

**AC-2 - `spacedock dispatch build --host paseo` emits a host-neutral worker envelope.**
Verified by: Go unit tests in `internal/dispatch/` assert `--host paseo` succeeds with `--stamp`, emits neutral `ensign` skill-load and `Done:` completion instructions, and contains zero Claude/Codex/Pi-only tool or mailbox directives (`SendMessage`, `$spacedock:ensign`, `/skill:ensign`, `pi-subagents`). Reintroducing any host-specific token or rejecting `--host paseo` fails the test.

**AC-3 - First-officer and ensign skills wire the Paseo runtime adapters cleanly.**
Verified by: `go test ./...` (including `internal/contractlint` and `skills/integration`) passes with `paseo-fo-dispatch.md` and `paseo-ensign-runtime.md` reachable from `SKILL.md` under `PASEO_AGENT_ID` / explicit Paseo selection. Breaking reference closure or contractlint rules fails the suite.

**AC-4 - Optional CI lane runs `TestLivePaseoDispatchSmoke` on `codex/gpt-6-luna` without blocking merges.**
Verified by: structural YAML checks in Go tests confirm the non-required workflow invokes `TestLivePaseoDispatchSmoke` with `codex/gpt-6-luna` and `CODEX_AUTH_JSON` (matching `pi-live` optionality), and skips cleanly locally when `paseo` or auth is absent. Making the lane required or miswiring the selector/model fails the check.

## Test plan

- `go test ./...` and `go test ./... -race` for `internal/dispatch`, `internal/contractlint`, `skills/integration`, and workflow YAML structural checks.
- Targeted local live run: `go test -tags live -count=1 -timeout 15m -run '^TestLivePaseoDispatchSmoke$' ./internal/ensigncycle -v` with `codex/gpt-6-luna` (and optional second run with `antigravity/gemini-4-neon` against `launch-wrap-hook`).

### Feedback Cycles
