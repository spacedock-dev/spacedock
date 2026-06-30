---
title: "dispatch build — warn + document --team-name on claude (auto-team is default; legacy shape is silent)"
status: backlog
score: 0.3
source: "v0.23.0 cut FO session, 2026-06-30. FO reflexively passed --team-name (from the boot team_state hint) to dispatch build on Claude, which silently emitted the LEGACY team-registry shape (team_name present, run_in_background absent) instead of the auto-team merged-mode shape. Self-corrected by rebuilding without --team-name; no bad dispatch shipped."
id: 0qt2r4n577gtwq707abgcfed
---

`spacedock dispatch build` IS Claude-team-mode-aware (`internal/dispatch/build.go:292`: `mergedMode := !bareMode && host == "claude" && teamName == ""`): omitting `--team-name` yields the auto-team shape (`run_in_background:true`, no `team_name`); passing `--team-name` selects the legacy TeamCreate-registry shape. The foot-gun is ergonomic, not correctness.

## Problem
1. `--team-name` / `--bare-mode` are NOT documented in `dispatch build --help` (it only shows the stdin-JSON interface), so the legacy opt-in is invisible at the CLI surface.
2. Nothing warns when `--team-name` is passed on `host=claude`, where auto-team is the default and the legacy path is sunsetting — so a stray `--team-name` silently produces a shape (team_name present, run_in_background absent) that is verbatim-unsafe for auto-team mode. The binary cannot self-detect legacy-vs-auto-team (that is FO-probed at runtime via TeamCreate/SendMessage availability, invisible to the CLI), so `--team-name` must stay the explicit signal — but a guard would catch the slip.

## Proposed approach
Either/both: (a) emit a stderr advisory when `--team-name` is passed on `host=claude` ("legacy team-name path; auto-team is the default on claude — omit --team-name unless you mean legacy mode"); (b) document `--team-name`/`--bare-mode` in `dispatch build --help`. Once legacy team mode is sunset, consider rejecting `--team-name` on claude outright.

## Acceptance criteria
- **AC-1** — passing `--team-name` to `dispatch build` on `host=claude` emits a one-line stderr advisory naming the legacy path and the auto-team default, AND/OR `dispatch build --help` documents `--team-name`/`--bare-mode`. Verified by: a test asserting the stderr advisory on a claude-host build invoked with `--team-name`, or a `--help` golden that includes the flags.
