---
id: v9pvmzhxvcmvps9tnz73vs4v
title: FO binary launcher invariant — prevent SPACEDOCK_BIN/PATH drift during boot and helper calls
status: backlog
source: captain request after Pi FO used PATH spacedock for later helper calls despite SPACEDOCK_BIN being set
started:
completed:
verdict:
score: 0.35
worktree:
issue:
sprint: 0230-stable-finalization
---

Tighten the first-officer contract and tests so a FO cannot silently switch from the resolved `SPACEDOCK_BIN` launcher to a different `spacedock` binary on PATH after startup.

## Problem

The FO startup contract says to prefer `${SPACEDOCK_BIN:-spacedock}`, but it also permits bare `spacedock` shorthand in examples. In this session, the initial version gate used the repo-local `SPACEDOCK_BIN`, while later status/state helper calls used bare `spacedock` from `/opt/homebrew/bin`, causing a false subcommand-missing result for `state ready` / `state sweep`.

This matters because launcher drift changes the command surface mid-session and can make the FO reason from the wrong binary's capabilities.

## Proposed approach

Clarify the first-officer contract so startup resolves one launcher variable, reports the resolved binary identity, and requires every Spacedock helper call after the version gate to use that resolved launcher. Remove or narrow the bare-command shorthand allowance so examples do not teach the wrong habit.

Prefer a code gate over prose: add a lint or test that fails on accidental bare `spacedock` helper invocations in FO contract examples, while allowing the explicit PATH fallback/diagnostic probes.

## Out of scope

- Changing launcher behavior for Claude/Codex/Pi hosts beyond the FO contract and tests.
- Reworking all non-FO documentation examples.
- Adding new `spacedock` subcommands.

## Acceptance criteria

**AC-1 - FO startup contract pins a single resolved launcher for the session.**
Verified by: the first-officer contract text requires resolving a launcher variable from `SPACEDOCK_BIN`/PATH, reporting its path/version, and using it for all later Spacedock helper calls.

**AC-2 - Bare `spacedock` examples are not allowed in FO helper-call guidance except explicit fallback diagnostics.**
Verified by: an automated test or lint fails when FO contract/reference docs contain bare `spacedock` helper examples outside an allowlisted fallback/diagnostic context.

**AC-3 - The regression case is covered.**
Verified by: a test fixture or live/fixture-backed check demonstrates that when `SPACEDOCK_BIN` and `command -v spacedock` point to different binaries, FO boot guidance keeps using the `SPACEDOCK_BIN` path for `status`, `state ready`, and other helper calls.

## Test plan

Add a focused contract-doc lint or Go test over `skills/first-officer/references/*.md` (and any generated shipped scaffolding if applicable). Include a fixture for the allowed initial fallback probe and disallowed post-boot bare helper examples. Run `go test ./...` as the baseline gate.
