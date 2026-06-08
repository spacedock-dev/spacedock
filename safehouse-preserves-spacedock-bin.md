---
id: thdf22bca6gn4xr8tra9r4z4
title: safehouse-wrapped `spacedock claude`/`codex` loses launcher-injected SPACEDOCK_BIN (runs PATH binary, not the launched one)
status: backlog
source: "FO OWED (HIGH), carried from the 2026-06-08-01 debrief, unfiled until now (captain-nodded 2026-06-08). A safehouse-wrapped launch does not propagate the launcher-injected SPACEDOCK_BIN, so helper calls run the PATH binary instead of the one that was launched. Workaround: export SPACEDOCK_BIN=<dev binary>."
started:
completed:
verdict:
score:
worktree:
issue:
---

A safehouse-wrapped `spacedock claude` / `spacedock codex` should keep using the binary that launched it. Today the launcher-injected `SPACEDOCK_BIN` is lost across the safehouse boundary, so every in-session helper call silently falls back to the `spacedock` on `$PATH`.

## Problem

The front door launches the host CLI with `SPACEDOCK_BIN` set to the launching binary (the launcher-command invariant: helpers prefer `$SPACEDOCK_BIN` over PATH). When the launch is wrapped by **safehouse**, the wrapper does not propagate `SPACEDOCK_BIN` into the sandboxed environment. Inside the session the variable is unset, so the FO/ensign helper calls (`status`, `dispatch build`, …) resolve `spacedock` from `$PATH` — which may be a different (released/stale) binary than the one under test.

Impact (HIGH for dev/test): a safehouse-wrapped dev launch appears to use the dev binary but actually drives the PATH binary. Behavior diverges silently — exactly the failure mode the launcher invariant exists to prevent. The current workaround is to `export SPACEDOCK_BIN=<dev binary>` before launching, which the operator must remember.

## Proposed approach

Ideation fills this in. Likely: have safehouse preserve/allowlist `SPACEDOCK_BIN` (and any sibling launcher env the front door injects) across the sandbox boundary, OR have the front door re-assert `SPACEDOCK_BIN` inside the wrapped environment. Determine which side owns the fix (safehouse config vs. front-door launch) during ideation; if it is a safehouse-side concern, this task may reduce to documenting the boundary + the required allowlist entry rather than a code change here.

## Out of scope

- safehouse's broader sandbox model.
- The general PATH-fallback behavior when no launcher is involved (that fallback is correct).

## Acceptance criteria

Ideation/implementation fills in. Sketch:

- A safehouse-wrapped `spacedock claude`/`codex` launch resolves helper calls to the launching binary, not the PATH binary (verified by launching a dev binary under safehouse and observing the in-session helper resolves to it — e.g. a probe that reports the resolved binary path/version, compared against the PATH binary's).
- If the resolution is a documentation/allowlist fix rather than code, the required safehouse env allowlist entry is documented and the doc names the failure it prevents.

## Test plan

Ideation/implementation fills in. The riskiest unknown is which side propagates env — exercise it first: launch a dev binary under safehouse with a deliberately-different PATH binary and observe which one the in-session helper resolves to (before/after the fix).
