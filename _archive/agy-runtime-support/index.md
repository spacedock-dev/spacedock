---
title: Add support for the agy (Antigravity) runtime host
status: implementation
source: captain (2026-06-05)
score: "0.50"
started: 2026-06-05T14:10:26Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-agy-runtime-support
issue:
id: 27ax7y7ecbrqyztffmhka9zd
archived: 2026-11-21T18:01:42Z
---
This task implements first-class support for the `agy` (Antigravity) runtime host in Spacedock.

## Problem
Currently, Spacedock supports Claude Code, Codex, and Pi. To use Spacedock within `agy` (Antigravity) as a First Officer, we need custom runtime adapters and command-line support in the Go binary for launching and installing plugins.

## Proposed Approach
1. Add `agy` to `validBuildHost` and update prompt blocks (first action, completion signal, dispatch pointer) in `internal/dispatch/build.go`.
2. Implement `runAgy` in `internal/cli/frontdoor.go` to handle `spacedock agy` invocation.
3. Wire `agy` in `internal/cli/cli.go` subcommands and `internal/cli/init.go` (`install` and `doctor` commands).
4. Create the first-officer and ensign runtime adapter documents under `skills/first-officer/references/agy-first-officer-runtime.md` and `skills/ensign/references/agy-ensign-runtime.md`.

## Acceptance Criteria
- **AC-1 - Agy host is supported by dispatch build:** `spacedock dispatch build --host agy` succeeds, producing the correct `firstActionBlock`, `completionSignalBlock`, and `dispatchPointerPrompt` for `agy` using `send_message`.
- **AC-2 - Agy command launches the agy CLI:** `spacedock agy` runs the `agy` CLI with the bootstrap prompt.
- **AC-3 - Install and Doctor commands support agy:** `spacedock install --host agy` and `spacedock doctor --host agy` are fully wired.
- **AC-4 - Runtime adapters exist:** Both `skills/first-officer/references/agy-first-officer-runtime.md` and `skills/ensign/references/agy-ensign-runtime.md` are created and contain the correct instructions.

## Test Plan
- Unit tests in `internal/dispatch` verifying `agy` host prompt generation.
- CLI tests verifying `spacedock agy`, `install --host agy`, and `doctor --host agy` commands.

## Superseded (2026-11-21)

Superseded by the Paseo-level runtime backend ([`_roadmap/paseo-backend/index.md`](../_roadmap/paseo-backend/index.md)). Running Antigravity (`agy`) through Paseo provides `agy` as a First Officer via the `paseo` CLI binding with zero `agy`-specific Go code (live-verified 2026-11-21: `agy` FO -> `paseo` CLI -> `safehouse-codex` ensign with automatic `ParentAgentId` linkage and durable Git commits), and covers `antigravity` ensigns via a host-neutral `--host paseo` dispatch mode once Paseo's per-session launch-wrapper hook unblocks `paseo-safehouse` for plugin-supplied providers. Direct `spacedock agy` frontdoor, `--host agy`, and native `send_message` adapters are no longer needed.
