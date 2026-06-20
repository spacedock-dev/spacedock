---
title: Pi launch fnm-multishell race — resolve pi to the stable node-installation bin, not the racy per-shell multishell symlink
status: implementation
source: "Captain (2026-06-20): `./spacedock pi --plugin-dir . -- --model z-ai/glm-5.2 --thinking xhigh` failed with `Error: Cannot find module '/Users/clkao/.local/state/fnm_multishells/65968_1781888446149/bin/pi'` (Node MODULE_NOT_FOUND). Root cause: spacedock's `exec.LookPath(\"pi\")` resolves to an fnm per-shell multishell symlink; between LookPath (Go, at launch) and Node's Module._resolveFilename (milliseconds later), fnm tore down the 65968_… multishell dir (that shell exited), so the absolute path Node was handed is dead. Reproduced: `node …/65968_…/bin/pi --version` works after the fact (fnm recreated/cleaned it) — the failure window is tiny but real. Not a Spacedock bug per se, but Spacedock is the victim and can be resilient."
score:
started: 2026-06-20T05:07:34Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-pi-launch-fnm-multishell-race
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

## Approach — ideation decision (finalized 2026-06-19)

**Mechanism chosen: (a) — resolve the fnm-multishell-symlinked `pi` through to the stable node-installation bin, bypassing the racy per-shell layer.** (b) bare-name exec is rejected: `execHost.Launch` already execs by resolving `argv[0]` via `exec.LookPath` at exec time, and the child inherits the same stale `PATH` that holds the tearing-down multishell dir — so bare-name still hits the dead path; (b) does not narrow the race. (c) catch `MODULE_NOT_FOUND` + retry is recorded as a DEFENSIVE BACKSTOP only (the child's Node error is already a process exit by the time Spacedock sees it; retry would need a fresh resolve-and-relaunch loop), DEFERRED past the (a) primary fix — not in this task's scope.

### Code-validated seam (correction to the dispatch's framing)

The dispatch's framing locates the bug at `execPiRuntimeOps.LookPath` (returns `exec.LookPath("pi")`). Reading the actual launch path corrects this: `runPi` builds `argv := []string{"pi", "--extension", …}` (argv[0] is the literal bare name `"pi"`) and calls `ops.Launch(argv)`; `execPiRuntimeOps.Launch` delegates to `execHost{}.Launch(argv, os.Environ())`; `execHost.Launch` does `bin, err := exec.LookPath(argv[0])` (the **stdlib**, not `piRuntimeOps.LookPath`) then `syscall.Exec(bin, argv, env)`. So the actual exec'd binary is chosen by **`execHost.Launch`'s stdlib `exec.LookPath(argv[0]="pi")`**, which BYPASSES `execPiRuntimeOps.LookPath` entirely. Therefore routing the resolution through `execPiRuntimeOps.LookPath` alone would NOT fix the launch (the banner/check would show the stable path, but `execHost.Launch` would re-resolve `"pi"` to the racy multishell path). The resolution must change **what `execHost.Launch` resolves as the exec'd `bin`**, which means setting `argv[0]` to the stable **absolute** path (so stdlib `exec.LookPath(<absolute path>)` returns it unchanged) OR supplying the stable bin to a pi-specific launch. The launch path in `runPi` is the correct seam; `execPiRuntimeOps.LookPath` is left as the raw `exec.LookPath` seam (unchanged) — it is the input the resolution consumes, not the fix site.

### shebang-argv spike (settles the argv[0] question — recorded 2026-06-19)

`syscall.Exec(bin, argv, env)` where `bin` is a `#!/usr/bin/env node` shebang script: the kernel transforms to `execve(node, [node, bin, argv[0], argv[1], …], env)`, so our `argv[0]` becomes the child Node's `process.argv[2]` (the first user arg), and `bin` becomes `process.argv[1]` (the script Node resolves via `Module._resolveFilename`). So setting `argv[0]` to the stable absolute path changes pi's `process.argv[2]` from `"pi"` to the stable path string — a potential pi parser regression. Spiked against the real stable bin `~/.local/share/fnm/node-versions/v24.13.1/installation/bin/pi`:
- `node <stable> pi --version` (today's shape, argv[0]="pi") → `0.79.8`, exit 0.
- `node <stable> <stable> --version` (fixed shape, argv[0]=stable path) → `0.79.8`, exit 0.
- `node <stable> <stable> --extension /no/such --skill /no/such -- --model z-ai/glm-5.2 --thinking xhigh` → arg-parse proceeds identically to the `"pi"` shape; the only errors are `Unknown option: --` and `Failed to load extension`, which are **shape-independent** (identical in the argv[0]="pi" shape) and are the SEPARATE `--` passthrough concern, out of scope here.

**Conclusion:** pi strips/ignores the leading `process.argv[2]` positional whatever it is ("pi" or a path string); setting `argv[0]` to the stable absolute path is SAFE for pi's parser. So the chosen design sets `argv[0]` to the stable path when the fnm-multishell pattern is detected (matches the dispatch's AC-1 framing literally), and leaves `argv[0]="pi"` for all non-multishell cases (no-regression guarantee).

### Finalized design

- New **pure helper** `resolveFnmMultishellPi(lookedUp string) (stable string, ok bool)` in `internal/cli/pi.go`:
  - if `!strings.Contains(lookedUp, "/fnm_multishells/")` → return `"", false` (non-fnm / direct installs unaffected).
  - `binDir := filepath.Dir(lookedUp)`; `stableDir, err := filepath.EvalSymlinks(binDir)`; on err → return `"", false` (fall back, never block the launch).
  - `stable := filepath.Join(stableDir, "pi")`; `os.Lstat(stable)` must succeed; on err → return `"", false`.
  - if `stable == lookedUp` → return `"", false` (already stable, no change).
  - return `stable, true`.
- `execPiRuntimeOps.LookPath`: **UNCHANGED** (`exec.LookPath`) — the raw seam the resolution consumes; routing it through the helper would let `execHost.Launch`'s stdlib re-resolve "pi" to the racy path, so it stays raw.
- `execHost.Launch`: **UNCHANGED** (shared with claude/codex via `hostOps`).
- `runPi`: immediately before `ops.Launch(argv)`, insert:
  ```go
  if lp, err := ops.LookPath("pi"); err == nil {
      if stable, ok := resolveFnmMultishellPi(lp); ok {
          argv[0] = stable
      }
  }
  ```
  On LookPath error or not-ok, `argv[0]` stays `"pi"` (current behavior, no regression). On ok, `argv[0]` becomes the stable absolute path; `execHost.Launch`'s `exec.LookPath(<absolute>)` returns it unchanged; `syscall.Exec(stable, argv, env)` → the kernel's shebang hands Node `process.argv[1]=stable` → `Module._resolveFilename` stats a path under the real, never-torn-down installation dir → the race is closed. `process.argv[2]` becomes the stable path string (spike-proven safe; see above).
- **Banner note (accepted, cosmetic):** `launchBanner` calls `ops.LookPath("pi")` and prints the resolved path; under fnm it will still show the multishell path while the launch uses the stable path. This is an informational mismatch only and is accepted to keep the change single-seamed (one resolution, in `runPi`). The banner does not affect launch correctness. (If a future pass wants banner/launch parity, route `launchBanner`'s `lookPath` through the helper — but that is a cosmetic nicety, NOT in this task's scope.)

### Generalization — nvm/volta (recorded; deferred — no race observed)

- **nvm:** PATH carries `~/.nvm/versions/node/v<ver>/bin` — a REAL directory (nvm only rewrites PATH on `nvm use`; the versions dir is not torn down per shell). No per-shell teardown race. **No action.**
- **volta:** bins under `~/.volta/bin` (real, stable) + shims; node version pinned via project `package.json`/volta config, not per-shell symlinks. No per-shell teardown race. **No action.**
- **fnm is the unique host** with a per-shell symlink-farm teardown (`fnm_multishells/<pid>_<ts>` is unlinked on shell exit). **First cut: fnm-only.** If a future shim with equivalent per-shell teardown appears, extend the pattern set in the helper (rename `resolveFnmMultishellPi` → `resolvePerShellShimPi`, add a pattern table) — recorded as the extension seam, not implemented now. The `"/fnm_multishells/"` substring check is the v1 detection; the helper's pure-function shape keeps the extension localized.

## Scope (finalized)

In scope:
- The `resolveFnmMultishellPi` pure helper (`internal/cli/pi.go`) + the `argv[0]` resolution block in `runPi` (the launch path in `runPi` — NOT `execPiRuntimeOps.LookPath`, per the code-validated seam correction above).
- Detect the fnm multishell pattern (`"/fnm_multishells/"` substring); resolve through to the stable `…/node-versions/<ver>/installation/bin/pi` via `filepath.EvalSymlinks` of the bin dir.
- Fallback to current behavior (argv[0]="pi") for non-matching paths or resolution failure (no regression on non-fnm setups or direct installs).
- A new `internal/cli/pi_launch_test.go` with a synthetic fnm symlink tree asserting the stable path is chosen (AC-1), the direct path is unchanged (AC-2), and the torn-down-parent fallback holds.
- nvm/volta: recorded as no-action (no per-shell teardown race observed; see Generalization) — deferred, not implemented.
- (c) catch `MODULE_NOT_FOUND` + retry: recorded as a deferred defensive backstop, NOT implemented in this task.

Out of scope:
- Changing fnm's teardown behavior (upstream, out of repo control).
- Routing `execPiRuntimeOps.LookPath` through the helper (would not fix the launch; see the seam correction) and `launchBanner` parity (cosmetic; banner note accepted).
- The (c) retry backstop.
- The safehouse-flag-parity work (`pi-safehouse-flag-parity`, sibling) — disjoint `runPi` region.
- Cross-host (claude/codex) — they don't go through this `pi` launch path.
- The `--` passthrough "Unknown option" pi-side error (spike-noted; separate concern).

## Acceptance criteria (finalized — entity-level; proof = behavior, never prose-grep)

**AC-1 — `spacedock pi` launches a fnm-multishell-symlinked `pi` via its stable node-installation bin.**
Property: when `ops.LookPath("pi")` returns a path under `*/fnm_multishells/*/bin/pi`, `runPi` passes the resolved stable install bin (`…/node-versions/<ver>/installation/bin/pi`, a path NOT under `fnm_multishells/`) as the launched `argv[0]` to `ops.Launch`, not the multishell path.
Verified by: a Go unit test (`TestResolveFnmMultishellPi_StableChosen` + a `runPi`-level `TestRunPi_LaunchArgv0_StableForMultishell`) that builds a synthetic fnm tree in `t.TempDir()` — `fnm_multishells/<pid>_<ts>` as a symlink to a synthetic `node-versions/<ver>/installation`, `…/bin/pi` as a relative symlink to the synthetic `lib/node_modules/.../dist/cli.js`; the fake `LookPath` returns the multishell path; the test asserts (a) `resolveFnmMultishellPi` returns the synthetic stable `…/installation/bin/pi` with `ok=true`, and (b) the fake `Launch`'s recorded `argv[0]` equals that stable path and does NOT contain `/fnm_multishells/`.

**AC-2 — Non-multishell `pi` resolution is unchanged (no regression).**
Property: when `ops.LookPath("pi")` returns a non-multishell path (direct install / `/usr/local/bin/pi`-style), the launched `argv[0]` is the literal `"pi"` exactly as today.
Verified by: a Go unit test (`TestResolveFnmMultishellPi_NonMultishell_Unchanged` + `TestRunPi_LaunchArgv0_UnchangedForDirect`) where the fake `LookPath` returns a path with no `/fnm_multishells/` substring; the test asserts `resolveFnmMultishellPi` returns `"", false` and the fake `Launch`'s `argv[0]` equals `"pi"`. The existing `internal/cli/pi_frontdoor_test.go` suite stays green (`go test ./internal/cli/`), confirming no front-door / doctor / check-only regression.

**AC-3 — The race window is closed end-to-end (determination: structural — CI race simulation impractical).**
Property: the path Node `Module._resolveFilename`s at child startup lives under a directory fnm never tears down, so a sibling shell exiting between Go's `LookPath` and the child's resolve cannot ENOENT the launched script.
Determination recorded: a TRUE teardown-during-launch race simulation is impractical in CI — it requires a sibling fnm shell to exit and fnm to `unlink` the `fnm_multishells/<pid>_<ts>` symlink inside the millisecond-wide window between Go's `execHost.Launch` `LookPath` and the child Node's `Module._resolveFilename`, non-deterministically, against a live fnm. AC-3 is therefore satisfied by **AC-1 + AC-2 + the structural no-teardown reasoning**: (i) the stable install bin's parent `~/.local/share/fnm/node-versions/<ver>/installation/bin` is a REAL directory (verified: not a symlink, not per-shell), and fnm's teardown only ever unlinks the per-shell `fnm_multishells/<pid>_<ts>` symlink (a sibling-shell artifact), never the shared `node-versions/` installation tree; (ii) with `argv[0]=stable`, `execHost.Launch`'s `exec.LookPath(<absolute stable>)` returns the stable path unchanged and `syscall.Exec(stable, …)` makes the kernel hand Node `process.argv[1]=stable` as the script; `Module._resolveFilename(stable)` stats a path in the never-torn-down installation dir → cannot ENOENT from a sibling-shell teardown. The dispatched fallback ("if a true race simulation is impractical in CI, AC-1+AC-2 + a reasoning argument that the stable path is never torn down suffices — record the determination") is invoked here. A `pi-live` **manual** smoke is recommended as a best-effort human verification (tear down the `fnm_multishells/<pid>_<ts>` symlink after Go's resolve but before re-invoking, confirm `spacedock pi --version` still launches via the stable path) but is NOT a CI gate.

## Test plan (finalized)

- **File:** new `internal/cli/pi_launch_test.go` (keeps the launch-resolution tests cohesive and separate from the front-door/doctor/`checkPiRuntime` tests in `pi_frontdoor_test.go`; reuses the existing `fakePiRuntimeOps` seam — no new fake).
- **AC-1 unit tests (fixture, on-disk synthetic tree — Go unit, ~30 LOC each, low cost):**
  - `TestResolveFnmMultishellPi_StableChosen`: in `t.TempDir()`, build `tmp/fnm_multishells/123_456` as a symlink to `tmp/stable-versions/v1.2.3/installation`; `…/installation/bin/pi` as a symlink to `../lib/node_modules/pkg/dist/cli.js` and create that target file; call `resolveFnmMultishellPi(tmp/fnm_multishells/123_456/bin/pi)` → assert returns `tmp/stable-versions/v1.2.3/installation/bin/pi`, `ok=true`. (Pure-function test; the helper takes only the looked-up string and uses `filepath.EvalSymlinks`/`os.Lstat` on the real FS against the synthetic tree.)
  - `TestRunPi_LaunchArgv0_StableForMultishell`: construct the synthetic tree as above; `fakePiRuntimeOps.LookPath` returns the multishell path; drive `runPi` with a minimal passthrough (`--version`) and a `fakeHost`/env that satisfies `piRuntimeLaunchReady`; assert the recorded `launchArgv[0]` equals the stable synthetic `…/installation/bin/pi` and does NOT contain `/fnm_multishells/`.
- **AC-2 unit tests (~15 LOC each, low cost):**
  - `TestResolveFnmMultishellPi_NonMultishell_Unchanged`: `resolveFnmMultishellPi("/usr/local/bin/pi")` → `"", false`; also `resolveFnmMultishellPi("/opt/bin/pi")` → `"", false`.
  - `TestRunPi_LaunchArgv0_UnchangedForDirect`: `fakePiRuntimeOps.LookPath` returns `/usr/local/bin/pi`; drive `runPi`; assert `launchArgv[0] == "pi"` (unchanged).
- **AC-2 regression gate:** `go test ./internal/cli/` (the existing `pi_frontdoor_test.go` suite) stays green — no `LookPath`/`Stat`/`checkPiRuntime` semantic change for non-multishell, so no rework of existing fixtures.
- **Edge test (~10 LOC):** `TestResolveFnmMultishellPi_TornDownParent_Fallback` — the multishell parent symlink is removed (simulating post-teardown) before the helper runs → `filepath.EvalSymlinks` errors → returns `"", false` (fall back, never block). This pins the "on resolution failure, leave argv[0]='pi'" guarantee.
- **AC-3:** structural — no CI test beyond AC-1/AC-2 (the determination above records why a live race sim is impractical). The `pi-live` **lane** is required for the manual smoke (this task touches the `internal/cli/pi.go` high-stakes front door), but the lane produces a human-verified smoke note, not an automated CI assertion. Validator gate accepts AC-1 + AC-2 green plus the structural reasoning recorded above (per the dispatch's fallback).
- **Lanes required:** `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal` (the AGENTS.md expected commands); `pi-live` for the manual smoke (best-effort, non-gating).
- **Cost/complexity:** low. The change is one pure helper (~20 LOC) + one 5-line block in `runPi` + ~120 LOC of unit tests. No new dependency, no interface change, no fake change. The synthetic-tree tests use only stdlib `os`/`path/filepath`/`testing`.

## Region ownership & merge coordination

- **This task owns:** the `resolveFnmMultishellPi` pure helper + the `argv[0]` resolution block in `runPi` (`internal/cli/pi.go`). `execPiRuntimeOps.LookPath` is **not** changed (left as the raw seam). `execHost.Launch` is **not** changed (shared with claude/codex).
- **Sibling `pi-safehouse-flag-parity` (`qn`)** owns `parsePiFrontDoorArgs` / the `runPi` safehouse-wrap arm / `setPiHelp` — disjoint region. Both touch `runPi`; the safehouse-wrap arm is earlier in `runPi` (banner → argv build → `safehouse.Wrap`), the resolution block here is immediately before `ops.Launch(argv)` at the end. The two edits are non-overlapping line ranges; coordinate if both land in the same window (the merge is a trivial append of the two blocks, but the FO should sequence the two implementations to avoid a `runPi` rebase churn).

## Spike notes — durable evidence (recorded 2026-06-19, against the live fnm install)

- **Stable install bin exists and is not torn down:** `~/.local/share/fnm/node-versions/v24.13.1/installation/bin/pi` → symlink (`lrwxr-xr-x`, 63 chars) to `../lib/node_modules/@earendil-works/pi-coding-agent/dist/cli.js`. Parent `…/installation/bin` is a REAL directory (not a symlink). Verified `os.Lstat` → mode `Lrwxr-xr-x`.
- **Multishell chain:** `~/.local/state/fnm_multishells/<pid>_<ts>/bin/pi` → the same 63-char relative target. The per-shell `<pid>_<ts>` dir is a **symlink** to `~/.local/share/fnm/aliases/default` → `~/.local/share/fnm/node-versions/v24.13.1/installation`. fnm's teardown unlinks the `<pid>_<ts>` symlink (sibling-shell exit), killing `…/bin/pi`.
- **EvalSymlinks proof (Go, the mechanism's basis):** `filepath.EvalSymlinks(multishellBinDir)` → `~/.local/share/fnm/node-versions/v24.13.1/installation/bin` (the real stable bin dir); `filepath.Join(that, "pi")` → the stable install bin symlink. `filepath.EvalSymlinks` of both the multishell and the stable `pi` resolves to the SAME underlying `…/dist/cli.js`, confirming they are the same script via two routes — only the multishell route's parent is torn down.
- **shebang-argv proof (spike above):** `argv[0]` becomes Node's `process.argv[2]`; pi tolerates both `"pi"` and the stable path string there (both shapes → `0.79.8`, exit 0). The `--` "Unknown option" error is shape-independent (separate passthrough concern, out of scope).
- **No re-spike needed:** the race was reproduced by the captain (the `65968_…` MODULE_NOT_FOUND invocation) and post-hoc (`node …/65968_…/bin/pi --version` works after the fact). The stable-path mechanism is proven by the EvalSymlinks + Lstat evidence above and the shebang-argv spike. This matches the dispatch's "No spike needed" determination, with the argv[0] question added and settled.

## Q15 — cold-boot package workaround

Until this task lands, the 0223 Commander cold-boot package carries Q15: `export PATH="<stable node bin>:$PATH"` (e.g. `export PATH="$HOME/.local/share/fnm/node-versions/v24.13.1/installation/bin:$PATH"`) **ahead** of `spacedock pi`, so `exec.LookPath("pi")` resolves to the stable bin first and bypasses the multishell layer entirely. This workaround expires once `resolveFnmMultishellPi` ships in `spacedock`.

## Related

- `internal/cli/pi.go` (`execPiRuntimeOps.LookPath`, `runPi` launch, `execHost.Launch` via `internal/cli/host_exec.go`) — the source of truth.
- `internal/cli/host_exec.go` — `execHost.Launch`'s `exec.LookPath(argv[0])` + `syscall.Exec` (the actual exec'd-bin selection site; unchanged by this task).
- `internal/cli/pi_frontdoor_test.go` — `fakePiRuntimeOps` seam (reused by the new `pi_launch_test.go`); unchanged.
- `pi-safehouse-flag-parity` (`qn`, sibling) — disjoint region in `runPi`; coordinate the merge window.
- The captain's failing invocation: `./spacedock pi --plugin-dir . -- --model z-ai/glm-5.2 --thinking xhigh` → `MODULE_NOT_FOUND` on `65968_…/bin/pi`.
- The 0223 Commander cold-boot package — Q15 workaround above.

## Stage Report: ideation (2026-06-19)

- DONE: Picked mechanism **(a)** — resolve the fnm-multishell-symlinked `pi` to the stable node-installation bin. Rejected (b) bare-name exec (`execHost.Launch` already LookPaths `argv[0]` at exec time and the child inherits the stale PATH → same dead path; (b) does not narrow the race). Recorded (c) catch-MODULE_NOT_FOUND-retry as a deferred defensive backstop, NOT in scope.
- DONE: Code-validated the seam and **corrected the dispatch's framing**. The dispatch locates the bug at `execPiRuntimeOps.LookPath`; reading `runPi` → `execPiRuntimeOps.Launch` → `execHost{}.Launch` shows the actual exec'd binary is chosen by **`execHost.Launch`'s stdlib `exec.LookPath(argv[0]="pi")`**, which bypasses `execPiRuntimeOps.LookPath`. So routing the resolution through `execPiRuntimeOps.LookPath` alone would NOT fix the launch. The correct seam is the launch path in `runPi`: set `argv[0]` to the stable absolute path so stdlib `exec.LookPath(<absolute>)` returns it unchanged. `execPiRuntimeOps.LookPath` stays the raw seam (unchanged); `execHost.Launch` stays shared (unchanged).
- DONE: Spiked and settled the argv[0] question. `syscall.Exec` of a `#!/usr/bin/env node` shebang script makes our `argv[0]` the child Node's `process.argv[2]` (first user arg) and `bin` the `process.argv[1]` script. Live spike against the real stable bin: `node <stable> pi --version` → 0.79.8 (today's shape); `node <stable> <stable> --version` → 0.79.8 (fixed shape). pi tolerates a path string as the leading `process.argv[2]`; the `--` "Unknown option" error is shape-independent (separate passthrough concern, out of scope). Conclusion: setting `argv[0]` to the stable absolute path is SAFE for pi's parser; no-regression guarantee holds by leaving `argv[0]="pi"` for all non-multishell cases.
- DONE: Finalized the design — a pure helper `resolveFnmMultishellPi(lookedUp string) (stable string, ok bool)` (detect `"/fnm_multishells/"`; `filepath.EvalSymlinks` of the bin dir → `Join(…, "pi")`; `os.Lstat` guard; fall back to `"", false` on any miss or resolution failure) + a 5-line `runPi` block before `ops.Launch(argv)` that sets `argv[0]=stable` only when `ok`. No interface change, no fake change, no new dependency.
- DONE: Finalized 3 behavior-bound ACs. AC-1 (multishell → argv[0]=stable install bin, NOT under fnm_multishells), AC-2 (non-multishell → argv[0]="pi" unchanged + existing pi_frontdoor suite green), AC-3 (race closed: determination = structural; CI race-sim impractical, so AC-1+AC-2 + the no-teardown reasoning; the stable install bin's parent is a REAL dir fnm never tears down). Each AC states its test method.
- DONE: Finalized the test plan — new `internal/cli/pi_launch_test.go` reusing `fakePiRuntimeOps`; synthetic fnm symlink tree in `t.TempDir()` for AC-1; direct-path for AC-2; torn-down-parent fallback edge test; `go test ./... -race` + `gofmt` lanes; `pi-live` manual smoke (best-effort, non-gating). Cost/complexity: low (~20-LOC helper + 5-line runPi block + ~120 LOC tests).
- DONE: Generalization recorded — nvm (`~/.nvm/versions/node/v<ver>/bin`, real dir, not per-shell) and volta (`~/.volta/bin`, stable) exhibit NO per-shell teardown race; **no action**, deferred. fnm is the unique host with per-shell symlink-farm teardown; first cut is fnm-only, with the rename-to-`resolvePerShellShimPi` extension seam recorded.
- DONE: Region ownership recorded — owns the helper + the `runPi` argv[0] block; does NOT touch `execPiRuntimeOps.LookPath` or `execHost.Launch`. Sibling `pi-safehouse-flag-parity` (qn) is disjoint (owns the earlier safehouse-wrap arm in `runPi`); merge coordination note recorded.
- DONE: Spike notes committed as durable evidence — stable install bin `Lrwxr-xr-x` (real parent dir, not torn down); multishell `<pid>_<ts>` is a symlink to `aliases/default` → `node-versions/…/installation`; Go `filepath.EvalSymlinks(multishellBinDir)` → the real stable bin dir; both routes resolve to the same `…/dist/cli.js`. Matches the dispatch's "no spike needed" determination with the argv[0] question added and settled.
- DONE: Q15 cold-boot workaround recorded — `export PATH="<stable node bin>:$PATH"` ahead of `spacedock pi` until this lands.
- SKIPPED: Editing any product file (no `pi.go` edits) — ideation = design only; implementation is a separate stage, per the dispatch.
- SKIPPED: Routing `execPiRuntimeOps.LookPath`/`launchBanner` through the helper — would not fix the launch (banner/check only) and is cosmetic; accepted as a single-seam tradeoff.

### Summary

Ideation finalized. The entity body now carries the code-validated mechanism decision (a, with the dispatch's seam framing corrected to the `runPi` argv[0] block — `execHost.Launch`'s stdlib LookPath is the real exec'd-bin selector, not `execPiRuntimeOps.LookPath`), the shebang-argv spike proving setting `argv[0]` to the stable path is pi-parser-safe, 3 behavior-bound ACs with test methods, a finalized test plan (new `pi_launch_test.go`, synthetic fnm tree, low cost), the nvm/volta no-action determination, region ownership + merge coordination with the sibling `qn` task, spike evidence, and the Q15 workaround. No product files were edited (ideation = design only). Ready for the ideation gate.

### Open risks / questions for the gate

1. **argv[0] becomes pi's `process.argv[2]` as the stable path string.** Spiked safe (`pi --version` exits 0 in both shapes), but the spike used `--version`; the FO/validator should confirm the manual `pi-live` smoke (a real model invocation) launches cleanly with `argv[0]=stable` before declaring AC-1 end-to-end. If a real invocation reveals pi DOES choke on a path as `process.argv[2]`, the fallback design (keep `argv[0]="pi"`, put the resolution in `execPiRuntimeOps.Launch` with a separate exec'd-bin override) is recorded in the seam-correction section and is a small pivot — flag for implementation, not a gate blocker.
2. **Banner/launch path mismatch under fnm** is accepted (cosmetic); recorded for a future cosmetic pass.
3. **AC-3 is structural** (CI race-sim impractical) per the dispatch's explicit fallback; the validator gate accepts AC-1+AC-2 green + the no-teardown reasoning. A `pi-live` manual smoke is recommended but non-gating.
