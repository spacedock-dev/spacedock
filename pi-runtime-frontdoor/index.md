---
title: Pi runtime front door and install UX
status: implementation
score: "0.44"
source: captain (2026-06-03) — after AC-2 proved Pi subagent dispatch, expose the working runtime contract as user-facing Spacedock commands
id: s9piruntimefrontdoorux01
started: 2026-06-03T00:00:00Z
worktree: .worktrees/spacedock-ensign-pi-runtime-support
---
# Pi runtime front door and install UX

Expose the proven Pi runtime path as Spacedock UX. The previous `pi-runtime-support` task proved that a Pi parent can dispatch a worker through `pi-subagents` from an isolated Pi home with copied OAuth credentials. This task turns that mechanism into command behavior: `spacedock pi`, `spacedock install --host pi`, and `spacedock doctor --host pi`.

## Problem

Pi runtime support currently works only through a live test harness that shells `pi` directly with explicit extension and skill paths. Users cannot discover or launch it from the Spacedock command surface, and `spacedock install --host pi` currently fails with `unknown host "pi"`.

## Scope

Implement the smallest compatibility-first UX over the proven mechanism:

- Add a `spacedock pi` launch command that starts Pi as the Spacedock first officer using Pi-native resources.
- Add `spacedock install --host pi` behavior that verifies or wires the required Pi substrate without pretending Pi is a Claude/Codex plugin marketplace.
- Add `spacedock doctor --host pi` checks for the Pi CLI, auth file, `pi-subagents` package/extension, and local Spacedock skill paths.
- Preserve current Claude/Codex behavior and output.
- Keep PR/mod behavior out of scope.

## Proven mechanism to reuse

Use `docs/runtime-support.md` as the implementation guide. The known-good Pi bringup shape is:

- host CLI: `pi`
- substrate package: `~/.pi/agent/npm/node_modules/pi-subagents`
- extension: `pi-subagents/src/extension/index.ts`
- skill: `pi-subagents/skills/pi-subagents`
- Spacedock skills: `<checkout>/skills/first-officer` and `<checkout>/skills/ensign`
- launch mode: Pi parent uses `subagent(...)`; no Claude `Agent`, `SendMessage`, `TeamCreate`, or `TeamDelete`

## Acceptance criteria

**AC-1 - `spacedock pi` is a registered launch command that uses Pi-native resources.**
Verified by: focused CLI tests that assert help lists `pi`, the launch argv begins with `pi`, includes the local or installed `pi-subagents` extension/skill paths plus Spacedock first-officer/ensign skill paths, and excludes Claude/Codex-only flags/tool names.

**AC-2 - `spacedock install --host pi` is accepted and idempotent.**
Verified by: focused CLI tests for `install --host pi` that do not call Claude/Codex plugin commands and either confirm the existing `pi-subagents` package path or print exact next-step instructions when it is missing.

**AC-3 - `spacedock doctor --host pi` reports actionable Pi health.**
Verified by: fixture/injected-host tests covering missing `pi`, missing `auth.json`, missing `pi-subagents`, and healthy local resources with stable output and exit codes.

**AC-4 - The live Pi smoke can launch through `spacedock pi` or an equivalent Spacedock-owned wrapper.**
Verified by: updating or adding a live-gated test, ideally `go test -tags live -run TestLivePiFrontDoorSmoke ./internal/ensigncycle -v -count=1`, that keeps isolated `PI_CODING_AGENT_DIR`/session dirs and asserts durable split-root entity/git-log outcomes.

**AC-5 - Existing Claude and Codex front-door/install/doctor behavior remains stable.**
Verified by: `go test ./... -count=1` and focused existing CLI tests for Claude/Codex front doors/install/doctor.

## Test plan

- Add TDD-first CLI tests for command registration/help and Pi launch argv shape.
- Add TDD-first install/doctor tests with injected host/filesystem seams before implementation.
- Reuse the previous live smoke fixture where possible, but route through the new Spacedock-owned Pi command or wrapper.
- Run `gofmt -w ./cmd ./internal`, `go test ./... -count=1`, `go test ./... -race -count=1`, and the live Pi smoke when credentials are available.

## Stage Report: implementation

- TODO: Implement the Pi front-door/install/doctor UX and tests.

### Summary

Implementation has been dispatched through the Pi-native `subagent(...)` runtime contract. The worker should update this report when code and verification evidence are ready.
