---
title: Pi launch fnm-multishell race — resolve pi to the stable node-installation bin, not the racy per-shell multishell symlink
status: ideation
source: "Captain (2026-06-20): `./spacedock pi --plugin-dir . -- --model z-ai/glm-5.2 --thinking xhigh` failed with `Error: Cannot find module '/Users/clkao/.local/state/fnm_multishells/65968_1781888446149/bin/pi'` (Node MODULE_NOT_FOUND). Root cause: spacedock's `exec.LookPath(\"pi\")` resolves to an fnm per-shell multishell symlink; between LookPath (Go, at launch) and Node's Module._resolveFilename (milliseconds later), fnm tore down the 65968_… multishell dir (that shell exited), so the absolute path Node was handed is dead. Reproduced: `node …/65968_…/bin/pi --version` works after the fact (fnm recreated/cleaned it) — the failure window is tiny but real. Not a Spacedock bug per se, but Spacedock is the victim and can be resilient."
score:
started: 2026-06-20T05:07:34Z
completed:
verdict:
worktree:
issue:
sprint:
sprint-readiness:
id: j7nhrmghyy0kmwtphd0fmq32
---

# Pi launch fnm-multishell race

## End value

`spacedock pi` launches reliably under fnm (and other per-shell Node-version shims) regardless of whether sibling shells are exiting. The `pi` binary is resolved to its **stable node-installation bin** (e.g. `~/.local/share/fnm/node-versions/<ver>/installation/bin/pi`), bypassing the racy per-shell multishell symlink layer, OR launched by bare name with a retry on `MODULE_NOT_FOUND`. No more launch failures from a stale/tearing-down multishell path.

## Problem — root cause already determined

`spacedock pi` (internal/cli/pi.go) launches via `execPiRuntimeOps.LookPath("pi")` → `exec.LookPath("pi")`, then `ops.Launch(argv)` where `argv[0]` is the resolved absolute path. Under fnm:

- `PATH` includes `~/.local/state/fnm_multishells/<pid>_<timestamp>/bin` — a per-shell symlink farm fnm creates on shell startup and tears down on shell exit. The `pi` entry there is a symlink to `../lib/node_modules/@earendil-works/pi-coding-agent/dist/cli.js` (relative), resolving through the node-installation.
- `exec.LookPath("pi")` returns that multishell absolute path (first `PATH` hit).
- Spacedock execs `node <multishell-path>`. Between `LookPath` (Go) and Node's `Module._resolveFilename` (child, milliseconds later), fnm can tear down the multishell dir (a sibling shell exited). The absolute path Node was handed is now dead → `Cannot find module '…/65968_…/bin/pi'` (`MODULE_NOT_FOUND`, `requireStack: []`).
- By the time anyone re-checks (`ls`, `node …/65968_…/bin/pi --version`), fnm has recreated or left the path, so it works — the failure window is tiny but real, and it strikes at launch.

Verified: `node /Users/clkao/.local/state/fnm_multishells/65968_1781888446149/bin/pi --version` returns `0.79.8` after the fact (the path resolved again). The stable install bin `~/.local/share/fnm/node-versions/v24.13.1/installation/bin/pi` exists and is NOT subject to teardown.

## Approach (candidate fixes — ideation confirms and picks)

- **(a) Resolve `pi` to the stable node-installation bin, bypassing the multishell layer.** When `LookPath("pi")` returns a path under `*/fnm_multishells/*/bin/pi`, resolve the symlink chain to the stable `~/.local/share/fnm/node-versions/<ver>/installation/bin/pi` (or the equivalent for nvm/volta/etc.) and launch that. The stable path is not torn down. Recommend — addresses the root cause.
- **(b) Exec `pi` by bare name (`exec.Command("pi", …)`) without pre-resolving.** Let the child's own `PATH` resolution happen at exec time, narrowing the race window. But Go's `exec.Command` may still `LookPath` internally; and if the child inherits the same stale `PATH`, it hits the same dead path. Weaker than (a).
- **(c) Catch `MODULE_NOT_FOUND` from the child and re-resolve + retry once.** Defensive backstop on top of (a) or (b); not a primary fix (the child's Node error is already a process exit).

Ideation picks one (recommend (a) — root-cause), records the decision, and generalizes: detect the multishell pattern (`*/fnm_multishells/*`, and the nvm/volta equivalents if straightforward) and resolve through to the stable install bin. Fallback to `LookPath` if the pattern doesn't match (non-fnm environments unaffected).

## Scope

In scope:
- `execPiRuntimeOps.LookPath` (or the launch path in `runPi`) resolves a multishell-symlinked `pi` to its stable node-installation bin.
- Detect the fnm multishell pattern; optionally nvm/volta if low-cost.
- Fallback to current `LookPath` behavior for non-matching paths (no regression on non-fnm setups or direct installs).
- A Go unit test with a fake multishell symlink tree asserting the stable path is chosen.

Out of scope:
- Changing fnm's teardown behavior (upstream, out of repo control).
- The safehouse-flag-parity work (`pi-safehouse-flag-parity`, sibling) — different pi front-door concern.
- Cross-host (claude/codex) — they don't go through this `pi` launch path.

## Acceptance criteria (provisional — ideation finalizes; proof = behavior, never prose-grep)

**AC-1 — `spacedock pi` resolves a multishell-symlinked `pi` to the stable node-installation bin.**
Verified by: a Go test with a fake `LookPath` returning a path under a synthetic `fnm_multishells/<pid>_<ts>/bin/pi` symlink that resolves to a stable install bin; assert the launched `argv[0]` is the stable install path, not the multishell path.

**AC-2 — Non-multishell `pi` resolution is unchanged (no regression).**
Verified by: a Go test where `LookPath` returns a direct (non-multishell) path; assert the launched `argv[0]` is that path unchanged. Existing frontdoor/pi tests stay green.

**AC-3 — The race window is closed end-to-end.**
Verified by: a `pi-live` or harness run that simulates a tearing-down multishell dir (remove the multishell symlink after LookPath but before Launch) and confirms the launch still succeeds via the stable path. (If a true race-simulation is impractical in CI, AC-1+AC-2 + a reasoning argument that the stable path is never torn down suffices — record the determination.)

## Test plan

- Go unit tests (AC-1, AC-2) in `internal/cli/pi_frontdoor_test.go` (or a new `pi_launch_test.go`): synthetic multishell symlink tree → stable path chosen; direct path → unchanged.
- AC-3: harness or `pi-live` race simulation if feasible; otherwise the structural tests + the no-teardown reasoning.
- `pi-live` lane required (touches `internal/cli/pi.go` launch path — high-stakes front door).

## Related

- `internal/cli/pi.go` (`execPiRuntimeOps.LookPath`, `runPi` launch) — the source of truth.
- `pi-safehouse-flag-parity` (`qn`, sibling) — also a pi front-door robustness task; disjoint region (`qn` owns `parsePiFrontDoorArgs`/`runPi` wrapping; this task owns `LookPath`/launch resolution). Coordinate if both land in the same window.
- The captain's failing invocation: `./spacedock pi --plugin-dir . -- --model z-ai/glm-5.2 --thinking xhigh` → `MODULE_NOT_FOUND` on `65968_…/bin/pi`.
- The 0223 Commander cold-boot package — add a Q15 quirk (workaround: `export PATH="<stable node bin>:$PATH"` ahead of `spacedock pi`) until this lands.
