# IDEATION GATE — Pi default extension discovery (`3w1`)

Recommendation: **APPROVE and dispatch implementation.**

## Selected approach

Test-harness-only discovery that mirrors what a normal Pi home does, so an
isolated home resolves `pi-subagents`/`pi-intercom` without
`PI_SUBAGENTS_PACKAGE_ROOT`.

1. **`piDefaultExtensionRoots(t, realHome)` helper** — reads the operator's
   `~/.pi/agent/settings.json` `packages` array, resolving `npm:<name>` and
   `file:<path>` entries via the same logic as Pi's
   `resolveSettingsPackageRoot` (`getManagedNpmInstallPath` →
   `join(agentDir, "npm", "node_modules", <name>)`), with npm-install-dir and
   extensions-dir probes as fallbacks.
2. **`seedPiDefaultExtensions(t, cleanHome, realHome)`** — symlinks the
   discovered roots into `cleanHome/.pi/agent/npm/node_modules/` so
   `piRuntimeConfigFromEnv`'s fallback resolves without the env var.
3. **`newPiLiveSmokeFixture` calls `seedPiDefaultExtensions` after
   `seedPiLiveAuth`.** Different files/seam from pnc's `models.json`/`auth.json`
   copy — no collision.

The explicit `PI_SUBAGENTS_PACKAGE_ROOT` / `PI_INTERCOM_PACKAGE_ROOT` override
still wins.

## Risk evidence

- **No spike needed — the mechanism is proven by tracing.** The worker traced
  Pi's discovery chain end-to-end: `settings.json` packages →
  `getManagedNpmInstallPath` (`package-manager.js:1670`) →
  `~/.pi/agent/npm/node_modules/<name>` → `package.json` `pi.extensions`. The
  existing hard-coded fallback in `piSubagentsPackageRoot`
  (`pi_live_runner_test.go:258`) mirrors exactly this path; the helper
  generalizes it to read the real `settings.json` instead of hard-coding.
  Verified `~/.pi/agent/npm/node_modules/pi-subagents/src/extension/index.ts`
  exists and `pi-subagents/package.json` declares `pi.extensions: ["./index.ts"]`.
- **Symlink vs production-code alternative.** The proposal seeds
  `cleanHome/.pi/agent/npm/node_modules/` with symlinks. The alternative —
  modifying `piRuntimeConfigFromEnv` to read `PI_CODING_AGENT_DIR` for the
  fallback — is a **production code change**, correctly scoped out (this task
  is test-harness-only). The symlink approach keeps the change test-only.
- **Duplicate `## Test plan` heading** — the body now has two (the original
  brief one at line 97 and the expanded per-AC proof plan at line 225). Cosmetic;
  the implementation worker should consolidate to one when it edits the body.

## Expected surface and tolerance

- `internal/ensigncycle/pi_live_runner_test.go` — `piDefaultExtensionRoots`
  (~40-60 lines) + `seedPiDefaultExtensions` (~15 lines) + `newPiLiveSmokeFixture`
  call site (~1 line).
- **Estimate: ~+60 lines net, tolerance ±15.** No production code changes; no
  semantic changes (test-harness-only discovery + symlink seeding).

## Semantic changes

None. Test-harness only. No CLI output, stored format, authority, or runtime
behavior change. The explicit env-var override is preserved.

## Proposed proof per acceptance criterion

- **AC-1 (VALUE — isolated-home run finds the extension with no env var):**
  `TestLivePiFrontDoorSmoke` passes with no `PI_SUBAGENTS_PACKAGE_ROOT`
  exported; the smoke resolves the subagents extension from the symlinked
  `cleanHome`. Baseline: today the smoke fatal-errors when the env var is unset.
- **AC-2 (reads real installed location, not hard-coded):** the helper reads
  `settings.json` packages and resolves via `getManagedNpmInstallPath`-equivalent
  logic; a machine with extensions installed elsewhere is discovered correctly.
- **AC-3 (explicit override still wins):** when `PI_SUBAGENTS_PACKAGE_ROOT` IS
  exported, it takes precedence over discovery (existing behavior preserved).
- **AC-4 (offline + smoke pass):** `gofmt`, `go vet -tags live`,
  `go build -tags live`, and `TestLivePiFrontDoorSmoke` pass with no env var
  exported.

## Decision ask

Approve to dispatch implementation (worktree) for this test-harness-only
extension discovery, or revise/hold with a concrete boundary.
