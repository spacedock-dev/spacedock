---
title: Pi dev-override skill loading — runPi must pass the Spacedock extension + skills when --plugin-dir is set
status: validation
source: "Captain (2026-06-20): eq (#406) retired the --skill cfg.firstOfficerDir() / --skill cfg.ensignDir() flags (D4) and moved skill discovery to the .pi/extensions/spacedock.ts extension's resources_discover. But runPi never passes the Spacedock extension to pi — it passes only --extension <pi-subagents> + --skill <pi-subagents>. The installed path works by accident (pi auto-loads registered extensions from settings.json packages); the dev path (--plugin-dir .) is BROKEN — the Spacedock skills are not loaded. pi does not auto-discover .pi/extensions/ from cwd (verified empirically). This is a regression from eq, and the pi parity gap with z2t's --plugin-dir caveat (z2t AC-6 documents that --plugin-dir bypasses installed-plugin resolution; pi's --plugin-dir bypasses skill loading entirely)."
score:
started: 2026-06-21T04:22:56Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-pi-dev-override-skill-loading
issue:
sprint:
sprint-readiness: ready
id: ev8gzecy33zm84spxrj239md
mod-block: merge:pr-merge
pr: "#421"
---

# Pi dev-override skill loading

## End value

`./spacedock pi --plugin-dir .` loads the Spacedock first-officer + ensign skills from the local checkout — the FO boots with its contract, not skill-less. The dev-override path is self-sufficient: it doesn't depend on the package being separately `pi install`'d into `settings.json packages`. Parity with how Codex's `--plugin-dir` loads the local plugin checkout.

## Problem — root cause already determined

`eq` (#406, commit `a3f7efa8`) shipped the install-managed skill placement:
- **D1:** root `package.json` with `pi.extensions: ["./.pi/extensions/spacedock.ts"]` + `pi.skills: ["./skills"]`.
- **D2:** `.pi/extensions/spacedock.ts` — a pi extension implementing `resources_discover` returning `{ skillPaths: [<repo>/skills] }`.
- **D3:** `spacedock install --host pi` runs `pi install git:github.com/spacedock-dev/spacedock` → registers the package in `settings.json packages` → pi auto-loads the extension at startup → `resources_discover` fires → skills announced.
- **D4:** retired the `--skill cfg.firstOfficerDir()` / `--skill cfg.ensignDir()` flags from `runPi`.

**The gap:** D4 retired the dev-path skill loading (the `--skill` flags) before confirming the extension-based replacement covers the dev path. `runPi` currently passes:
```go
argv := []string{
    "pi",
    "--extension", cfg.extensionPath,        // pi-subagents extension (NOT Spacedock)
    "--skill", cfg.subagentsSkill,            // pi-subagents skill (NOT Spacedock)
}
```
The Spacedock extension (`.pi/extensions/spacedock.ts`) is **NOT passed**. pi does NOT auto-discover `.pi/extensions/` from cwd (verified empirically — a test extension in `/tmp/pi-ext-test/.pi/extensions/test-ext.ts` did NOT fire when `pi` was launched from that cwd). Extensions are only loaded if (a) registered via `pi install` in `settings.json packages`, OR (b) passed explicitly via `--extension <path>`.

**The installed path works by accident:** the Spacedock package is registered in `settings.json packages` (from a prior `pi install`), so pi auto-loads the extension → skills announced. But this is independent of `--plugin-dir` — it works because the package was separately installed, not because the dev-override does anything for skill loading.

**The dev path (`--plugin-dir .`) is broken:** no `pi install` in the dev path → the Spacedock extension isn't registered → not loaded → `resources_discover` never fires → **the Spacedock FO/ensign skills are not loaded**. The FO boots without its contract.

## Parity with z2t

`z2t` (`spacedock-marketplace-source-env`) is the Codex equivalent: it ensures `spacedock install --host codex` works programmatically, the resolver is channel-aware, and `--plugin-dir` is documented as "bypasses installed-plugin resolution." The pi parity gap: `z2t` AC-6 documents that Codex's `--plugin-dir` bypasses installed-plugin resolution; **pi's `--plugin-dir` bypasses skill loading entirely** — it's not just a resolution bypass, it's a skill-loading regression.

## Approach

`runPi` must pass the Spacedock extension + skills explicitly when the dev-override is active (`--plugin-dir` or `SPACEDOCK_REPO_ROOT` set, i.e. `cfg.repoRoot != ""`). This is the dev-override equivalent of what `pi install` registers for the installed path:

```go
// After the existing argv (pi, --extension <pi-subagents>, --skill <pi-subagents>),
// add the Spacedock extension + skills from the dev-override checkout.
if cfg.repoRoot != "" {
    spacedockExt := filepath.Join(cfg.repoRoot, ".pi", "extensions", "spacedock.ts")
    spacedockSkills := filepath.Join(cfg.repoRoot, "skills")
    if _, err := os.Stat(spacedockExt); err == nil {
        argv = append(argv, "--extension", spacedockExt, "--skill", spacedockSkills)
    }
}
```

This makes the dev path self-sufficient: `pi` loads the Spacedock extension via `--extension`, the extension fires `resources_discover`, and the skills are announced to the parent session. pi-subagents children discover the same skills via `collectSettingsPackageSkillPaths` (if the package is also installed) OR via the parent's `--skill` forwarding (pi-subagents inherits `--skill` flags? — verify; if not, the dev-path child discovery needs the same `--skill` on the subagent spawn, which is the Q3/cwd concern from 0223).

### Installed path (no --plugin-dir)

When `cfg.repoRoot == ""` (installed path, no dev override), the Spacedock extension is loaded by pi's registered-package mechanism (`settings.json packages` → pi auto-loads). `runPi` does NOT need to pass the extension explicitly — the installed path works. But verify this: does `runPi` need to pass `--extension <installed-package-path>/.pi/extensions/spacedock.ts` even in the installed case, or does pi auto-load registered extensions at startup? If pi auto-loads, the installed path is fine as-is. If not, `runPi` must resolve the installed package's extension path from `settings.json packages` and pass it.

## Scope

In scope:
- `runPi` passes `--extension <repo>/.pi/extensions/spacedock.ts` + `--skill <repo>/skills` when `cfg.repoRoot != ""` (dev-override path).
- Verify the installed path (no `--plugin-dir`) — does pi auto-load registered extensions, or does `runPi` need to pass the extension there too?
- A Go test asserting the dev-override argv includes the Spacedock extension + skills paths.
- A live/manual smoke: `./spacedock pi --plugin-dir . -- --version` (or a full FO boot) confirms the Spacedock skills load (the FO greets with the workflow state summary, proving the FO contract loaded).

Out of scope:
- The installed-path `pi install` mechanism (D3, already shipped by `eq`).
- The `z2t` Codex parity (separate task).
- The `t0g` binding-block restructure (separate task).
- The fnm-multishell race (already shipped #416).
- The prompt-warmth stopgap (already merged #417's sibling).

## Acceptance criteria (provisional — ideation finalizes; proof = behavior)

**AC-1 — `./spacedock pi --plugin-dir .` loads the Spacedock FO/ensign skills.**
Verified by: a live/manual smoke — launch `spacedock pi --plugin-dir .` from the repo root and confirm the FO boots with its contract (greets with the workflow state summary, not a skill-less prompt). The FO contract loading is the behavioral proof.

**AC-2 — `runPi` passes `--extension <repo>/.pi/extensions/spacedock.ts` + `--skill <repo>/skills` when `cfg.repoRoot != ""`.**
Verified by: a Go test with a fake `ops` capturing the argv; assert the dev-override argv includes the Spacedock extension + skills paths.

**AC-3 — The installed path (no `--plugin-dir`) still works.**
Verified by: a Go test or live smoke confirming the installed path (no `cfg.repoRoot`) does NOT add the dev-override extension/skill flags (the installed mechanism handles it), OR does add them from the resolved package path (if pi doesn't auto-load registered extensions — verify + record).

## Test plan

- Go unit test (AC-2): fake `ops` capturing argv; dev-override argv includes Spacedock extension + skills.
- Live/manual smoke (AC-1): `spacedock pi --plugin-dir .` boots the FO with the contract.
- `pi-live` lane (touches `internal/cli/pi.go` — pi-only surface).

## Related

- `eq` `pi-install-managed-skill-placement` (#406, merged) — retired the `--skill` flags (D4) that this task restores for the dev path.
- `z2t` `spacedock-marketplace-source-env` — the Codex parity task; AC-6 documents the `--plugin-dir` caveat this task closes for pi.
- `internal/cli/pi.go` (`runPi`, `piRuntimeConfigFromEnv`) — the source of truth.
- `.pi/extensions/spacedock.ts` — the extension that fires `resources_discover` (shipped by `eq`).
- `docs/runtime-support.md` — the runtime guide (may need updating: the `--plugin-dir` dev path now loads the extension explicitly).

## Stage Report: implementation

- DONE: Fixed the eq regression in `internal/cli/pi.go` `runPi`. When the dev-override is active (`cfg.repoRoot != ""`, i.e. `--plugin-dir` / `SPACEDOCK_REPO_ROOT`), runPi now appends `--extension <repo>/.pi/extensions/spacedock.ts` and `--skill <repo>/skills` after the existing pi-subagents extension/skill argv, guarded by `os.Stat(spacedockExt)` so an absent extension falls back gracefully (no crash, no flags). This makes the dev path self-sufficient: pi loads the Spacedock extension via `--extension`, the extension fires `resources_discover`, and the Spacedock FO/ensign skills are announced to the parent session. pi does not auto-discover `.pi/extensions/` from cwd (verified empirically in the entity body), so without these flags the dev path boots skill-less.
- DONE: Verified the installed path (no `--plugin-dir`, `cfg.repoRoot == ""`). Determination: the installed path does NOT need the same treatment. pi auto-loads registered extensions from `~/.pi/agent/settings.json` `packages` at startup — the install-managed contract shipped by eq (D3: `spacedock install --host pi` runs `pi install` → registers the package → pi auto-loads `.pi/extensions/spacedock.ts` → `resources_discover` fires → skills announced). Confirmed against the live `~/.pi/agent/settings.json` which lists the Spacedock package (`git:github.com/spacedock-dev/spacedock` and a local path entry) under `packages`. runPi therefore does NOT add the Spacedock extension/skill flags when `cfg.repoRoot == ""`; the installed mechanism handles it. A Go test (`TestRunPi_InstalledPathDoesNotPassSpacedockExtension`) pins this: the installed-path argv has exactly one `--skill` (pi-subagents) and no Spacedock extension/skill tokens.
- DONE: Added three focused Go tests in `internal/cli/pi_frontdoor_test.go`: (1) `TestRunPi_DevOverridePassesSpacedockExtensionAndSkills` — dev-override argv includes `--extension <repo>/.pi/extensions/spacedock.ts` + `--skill <repo>/skills`, exactly two `--skill` flags, pi-subagents extension precedes the Spacedock extension; (2) `TestRunPi_DevOverrideWithoutExtensionFallsBackGracefully` — dev-override with the extension file absent adds no Spacedock flags, exactly one `--skill` (pi-subagents), no crash; (3) `TestRunPi_InstalledPathDoesNotPassSpacedockExtension` — installed path (no dev override) does not add the Spacedock extension/skill flags. The pre-existing `TestPiFrontDoorLaunchesWithNativeResourcePaths` still passes unchanged (its fixture repo has no `.pi/extensions/spacedock.ts`, so the os.Stat guard skips — the graceful-fallback path).
- DONE: Required gates green. `gofmt -w ./cmd ./internal` clean. `go test ./...` all packages PASS. `go test ./internal/cli/ -race` PASS.
- DONE: Manual smoke (load-bearing). Built the worktree binary and ran from the repo root: `./spacedock pi --plugin-dir . -- --version` → launches clean, prints `0.79.8`, no MODULE_NOT_FOUND, no skill-load failure, exit 0. Full FO boot: `./spacedock pi --plugin-dir . "report the current workflow state summary in one line, then stop"` → the FO greets with the workflow state summary (`Workflow docs/dev is discovered; boot shows sd-b32 IDs (next_id=...), no orphans, no PRs, no dispatchables, no team present, and split-root state is not initialized because .spacedock-state is missing`), exit 0. The FO contract loaded (sd-b32 IDs, next_id, orphans/PRs/dispatchables, split-root state) — the behavioral proof for AC-1. Worktree commit: `d11f7cc9d7b646df7f0415e11bb4939b36df6608` on `spacedock-ensign/pi-dev-override-skill-loading`.

### Summary

Closed the eq (#406) dev-override skill-loading regression. `runPi` now passes `--extension <repo>/.pi/extensions/spacedock.ts` + `--skill <repo>/skills` when `cfg.repoRoot != ""` (guarded by `os.Stat`), so `./spacedock pi --plugin-dir .` loads the Spacedock FO/ensign skills from the local checkout and the FO boots with its contract. The installed path (`cfg.repoRoot == ""`) is unchanged and verified correct: pi auto-loads the registered extension from `settings.json` `packages` at startup, so runPi does not need to pass the extension there. Three Go tests pin the dev-override add, the graceful fallback, and the installed-path no-op; all gates green; live smoke confirms a clean `--version` launch and a full FO boot that greets with the workflow state summary.

## Stage Report: validation

- **AC-1 — `./spacedock pi --plugin-dir .` loads the Spacedock FO/ensign skills.** Re-confirmed in the validation env (pi present, `~/.pi/agent/settings.json` present). Built the worktree binary and ran from the repo root: `spacedock pi --plugin-dir . -- --version` → launches clean, prints `0.79.8`, no MODULE_NOT_FOUND, no skill-load failure, exit 0. Full FO boot: `spacedock pi --plugin-dir . "report the current workflow state summary in one line, then stop"` → FO greets with the workflow state summary (`Workflow docs/dev: split-root state checkout is not initialized (.spacedock-state missing), no dispatchable entities, no PRs, team absent, ID style sd-b32 (next z0jbe5mne3fks4b99019223x).`), exit 0. The sd-b32 IDs / next_id / split-root / dispatchables / PRs / team fields prove the FO contract loaded — the dev path is no longer skill-less. PASSED.
- **AC-2 — `runPi` passes `--extension <repo>/.pi/extensions/spacedock.ts` + `--skill <repo>/skills` when `cfg.repoRoot != ""`.** Verified `TestRunPi_DevOverridePassesSpacedockExtensionAndSkills` (pi_frontdoor_test.go): (a) asserts `--extension <repo>/.pi/extensions/spacedock.ts` present, (b) asserts `--skill <repo>/skills` present, (c) asserts exactly two `--skill` flags (pi-subagents + spacedock), (d) asserts pi-subagents extension precedes the Spacedock extension via index comparison. The implementation in `internal/cli/pi.go` matches: appends `--extension`, `spacedockExt`, `--skill`, `spacedockSkills` after the pi-subagents argv, guarded by `os.Stat(spacedockExt)`. PASSED.
- **AC-3 — installed path (no `--plugin-dir`) does NOT add the Spacedock flags.** Verified `TestRunPi_InstalledPathDoesNotPassSpacedockExtension`: launches with no `--plugin-dir`/`SPACEDOCK_REPO_ROOT` (cfg.repoRoot == ""), asserts the argv contains no `.pi/extensions/spacedock.ts` token, no `--skill <repo>/skills`, and exactly one `--skill` (pi-subagents). The implementation's `if cfg.repoRoot != ""` guard excludes the installed path. PASSED.

### Adversarial audit

- **(a) Removed the `os.Stat` guard** (always append the Spacedock flags when `cfg.repoRoot != ""`): `TestRunPi_DevOverrideWithoutExtensionFallsBackGracefully` went RED — "dev-override argv must not reference absent extension". Confirms the guard is load-bearing for the graceful-fallback case. ✅ mutation caught.
- **(b) Inverted the condition to `cfg.repoRoot == ""`**: `TestRunPi_DevOverridePassesSpacedockExtensionAndSkills` went RED — "dev-override argv missing --extension". Confirms the dev-override guard is correctly keyed on the non-empty repoRoot. ✅ mutation caught.
- **(c) Removed `--skill <repo>/skills`** (kept only `--extension`): `TestRunPi_DevOverridePassesSpacedockExtensionAndSkills` went RED — "dev-override argv missing --skill" (the exactly-two-`--skill` assertion). Confirms both flags are asserted, not just the extension. ✅ mutation caught.

### Gate results

- `go test ./...` — PASS (all packages).
- `go test ./... -race` — PASS (all packages).
- `gofmt -l ./cmd ./internal` — clean (no output).

### Summary

VERDICT: **PASSED**. The fix at commit `d11f7cc9` correctly closes the eq (#406) dev-override skill-loading regression. `runPi` appends `--extension <repo>/.pi/extensions/spacedock.ts` + `--skill <repo>/skills` when `cfg.repoRoot != ""` (after the existing pi-subagents extension/skill), guarded by `os.Stat` so an absent extension falls back gracefully. The installed path (`cfg.repoRoot == ""`) is unchanged and correctly omits the flags (pi auto-loads the registered extension from `settings.json` `packages`). All three ACs are met, all three adversarial mutations turned the right tests RED, all required gates are green, and the live FO-boot smoke re-confirms AC-1 in the validation env.
