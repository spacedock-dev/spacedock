---
title: "spacedock pi dedupes the pi-subagents extension when it is installed as a package"
status: ideation
source: "GitHub issue spacedock-dev/spacedock#746 — spacedock pi fails at startup when pi-subagents is installed as a package: duplicate extension load ('Tool \"subagent\" conflicts with ...')"
started: 2026-08-21T03:45:47Z
completed:
verdict:
score: 0.8
worktree:
issue: spacedock-dev/spacedock#746
id: 5xwwj9c2w50921t16s840p49
gates:
    version: 1
    records:
        - id: gate:5xwwj9c2w50921t16s840p49:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:5xwwj9c2w50921t16s840p49-backlog-1
              briefing:
                id: briefing:5xwwj9c2w50921t16s840p49:backlog:attempt-1:revision-1
                digest: sha256:4a06c5d1093c6c580c7f47bce16a0e356f14ca29418e5c33e6f83e2229f7690d
                request-digest: sha256:dad3caab27a15e3b0e6ab5466355fa1543ee42dca9882c46459255e76ecfad7f
                room-ref: ./pi-subagents-duplicate-extension-load/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:5xwwj9c2w50921t16s840p49:backlog:1
                briefing: briefing:5xwwj9c2w50921t16s840p49:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-21T03:44:29.841226825Z"
                decision: approve
                reason: 'Captain approved backlog gate: seed is clear, issue #746 scope and direction sound, proceed to ideation for design+ACs+test plan.'
              application:
                target-stage: ideation
                state: consumed
---

`spacedock pi` loads the `pi-subagents` extension twice under two different specifiers when the package is also registered in `~/.pi/agent/settings.json` `packages` — the exact setup `spacedock pi --check` recommends. Package discovery loads `<pkg>/index.ts` (re-exporting `./src/extension/index.ts`), and `spacedock pi` additionally passes `--extension <pkg>/src/extension/index.ts --skill <pkg>/skills/pi-subagents` (internal/cli/pi.go, argv built at the `--extension`/`--skill` block; extensionPath is `filepath.Join(pkg, "src", "extension", "index.ts")`). Pi keys extension identity by resolved specifier, not module identity, so the second registration of `subagent`/`subagent_wait` collides and startup fails with `Tool "subagent" conflicts with ...`. Plain `pi` works; `spacedock pi -- --ne` works only by silencing all discovered extensions (lossy — also drops pi-intercom and host extensions).

Direction (from issue #746): the explicit flags exist for the `--plugin-dir`/`SPACEDOCK_REPO_ROOT` dev-override case (cfg.repoRoot != "", where the extension is NOT registered), per the code comments. Gate them on the package NOT already being registered — skip the explicit `--extension`/`--skill` when pi-subagents is in `~/.pi/agent/settings.json` `packages` (or when extensionPath resolves inside `~/.pi/agent/npm/node_modules`) — preserving the dev-override purpose. Alternative if the explicit flag must stay unconditional: point it at the package's declared entry (`<pkg>/index.ts`) so both load paths produce one specifier (more fragile, depends on pi-subagents' internal layout).

Acceptance sketch: value — with `npm:pi-subagents` in `settings.json` `packages`, `spacedock pi` starts and the subagent tools register exactly once (no `Tool "subagent" conflicts`), and the dev-override path still loads the extension explicitly; mechanism — a behavior test driving the flag-selection path against a registered-vs-unregistered package config asserts one extension load in the registered case and the explicit flag in the dev-override case. Expected surface: `internal/cli/pi.go` plus a test; small diff.

## Proposed approach

Gate the explicit `--extension`/`--skill` (pi-subagents) on the package NOT being registered, so pi's own package discovery is the sole load path when the package is in `settings.json` `packages`:

1. **Detection: extend `piSpacedockPackageStatus` to also report pi-subagents registration.** The function already iterates `settings.json` `packages`, resolves each entry to a root, and reads the package `name` from `package.json`. While iterating, also set `subagentsRegistered: true` on the returned `piPackageStatus` when a package named `"pi-subagents"` is found. This reuses the existing settings.json read and the existing `ops.SpacedockPackageStatus` test seam — no new I/O, no new interface method. The iteration must not early-return on the first match: it finds "spacedock" (for the existing gate) and continues scanning for "pi-subagents" so both flags are set in one pass.
2. **Gate in `runPi`.** Replace the unconditional `argv := []string{"pi", "--extension", cfg.extensionPath, "--skill", cfg.subagentsSkill}` with:
   ```go
   argv := []string{"pi"}
   if !check.packageStatus.subagentsRegistered {
       argv = append(argv, "--extension", cfg.extensionPath, "--skill", cfg.subagentsSkill)
   }
   ```
   When pi-subagents IS registered, pi auto-loads `<pkg>/index.ts` (which re-exports `./src/extension/index.ts`) — one specifier, no collision. When NOT registered (custom `PI_SUBAGENTS_PACKAGE_ROOT` outside the npm store, or a fresh pi-home without `pi install npm:pi-subagents`), the explicit flags are the only load path — current behavior preserved.
3. **Dev-override interaction.** The existing dev-override fallback in `checkPiRuntime` overwrites `check.packageStatus` when `!spacedockPackageOK && cfg.repoRoot != ""`. It must preserve `subagentsRegistered` from the original status (pi-subagents registration is independent of the Spacedock package registration). The dev-override Spacedock extension/skill block (`cfg.repoRoot != ""`) is unchanged — it always appends the Spacedock extension when present at `<repo>/.pi/extensions/spacedock.ts`.

### New mechanism — justification

- **Settings.json pi-subagents registration check** → serves AC-1 (value: no duplicate load). Simplest alternative: test whether `cfg.extensionPath` resolves inside `~/.pi/agent/npm/node_modules` (path prefix check, no I/O). Insufficient: it is a proxy — a package installed in the default npm location but NOT registered in `settings.json` would false-positive, silently dropping the explicit flags and breaking a custom-install user who has the files but no registration. The settings.json read is already performed for the Spacedock gate; checking for pi-subagents in the same pass is zero additional I/O.
- **Alternative approach (not selected): point the flag at `<pkg>/index.ts`.** Instead of gating, change `extensionPath` from `filepath.Join(pkg, "src", "extension", "index.ts")` to `filepath.Join(pkg, "index.ts")` so both the explicit flag and pi's discovery load the same specifier — pi deduplicates. More fragile: it depends on pi-subagents' internal layout (that `index.ts` re-exports `./src/extension/index.ts`), so a future pi-subagents refactor that moves or removes the re-export silently breaks the load. The gating approach depends only on the stable settings.json contract.

## Spike (riskiest unverified mechanism)

Spike target: does `settings.json` `packages` reliably list pi-subagents when installed via `pi install npm:pi-subagents`, and does `<pkg>/index.ts` re-export the extension entry that causes the duplicate load? Exercised live on this machine:

- `~/.pi/agent/settings.json` `packages` contains `"npm:pi-subagents"` — confirmed (settings.json reliably lists the package after `pi install`).
- `~/.pi/agent/npm/node_modules/pi-subagents/index.ts` contains `export { default } from "./src/extension/index.ts"` — confirmed (the re-export that makes pi's discovery load `<pkg>/index.ts` collide with the explicit `<pkg>/src/extension/index.ts`).
- `~/.pi/agent/npm/node_modules/pi-subagents/src/extension/index.ts` exists — confirmed (the path `spacedock pi` currently passes as `--extension`).

Result: both detection mechanisms are proven viable. The settings.json approach is selected (more precise). No further spike needed: the `piSpacedockPackageStatus` settings.json read is already a proven mechanism (shipped, tested), and the collision mechanism (pi keys by resolved specifier, two different specifiers → `Tool "subagent" conflicts`) is confirmed by the issue report and the spike evidence.

## Acceptance criteria

**AC-1 (value-measuring) — With `npm:pi-subagents` in `settings.json` `packages`, `spacedock pi` produces a launch argv with zero explicit `--extension`/`--skill` flags for pi-subagents, so the subagent tools register exactly once with no `Tool "subagent" conflicts` error.**

This measures the end-value the entity exists for against a baseline that can move the wrong way: the current unconditional-flag behavior always passes `--extension <pkg>/src/extension/index.ts`, which collides with pi's discovery of `<pkg>/index.ts` (two specifiers for the same extension) — the baseline fails. After the fix, the registered case drops the explicit flags (one specifier, no collision); the unregistered case (custom `PI_SUBAGENTS_PACKAGE_ROOT`, no registration) keeps the explicit flags (the only load path). Verified by: a behavior test asserting the launch argv has 0 `--extension` and 0 `--skill` flags for pi-subagents when `subagentsRegistered` is true, and has exactly 1 `--extension` + 1 `--skill` when false.

**AC-2 (mechanism) — The flag-selection path is exercised by a behavior test: registered config → no explicit pi-subagents flags; unregistered config (dev-override or custom path) → explicit pi-subagents flags; the dev-override Spacedock extension block is unaffected.**

Verified by: a `runPi` behavior test in `internal/cli/pi_frontdoor_test.go` driving the flag-selection path against (a) a registered config (`subagentsRegistered: true` → 0 pi-subagents `--extension`/`--skill`) and (b) an unregistered config (`subagentsRegistered: false` → 1 pi-subagents `--extension`/`--skill`), with the dev-override Spacedock extension (cfg.repoRoot) still appended in both cases. The test that would make it fail: flipping `subagentsRegistered` to true and seeing the explicit flags appear (or vice versa).

## Expected surface and tolerance

- **Files:** `internal/cli/pi.go` (the gate in `runPi`, the `subagentsRegistered` field on `piPackageStatus`, the extended iteration in `piSpacedockPackageStatus`, the dev-override `subagentsRegistered` preservation) + `internal/cli/pi_frontdoor_test.go` (the new behavior test).
- **Net LOC:** +25 to +35 across 2 files. Insertions: ~35 (field + detection + gate + test). Deletions: ~5 (the unconditional argv literal becomes conditional). Net: +30.
- **Tolerance:** ±10 LOC net.
- **Observable semantics:** NO change to command grammar, dev-override path, or runtime behavior beyond dropping the duplicate extension load when the package is registered. The `spacedock pi` command surface, flags, and help text are unchanged. The dev-override (`--plugin-dir`/`SPACEDOCK_REPO_ROOT`) Spacedock extension block is unchanged. The only observable change is: the launch argv for the registered case no longer carries `--extension <pkg>/src/extension/index.ts --skill <pkg>/skills/pi-subagents` — pi's discovery handles both.

## Test plan

- **Flag-selection behavior test (cheap, primary):** new `TestRunPi_RegisteredSubagentsDedupesExtensionLoad` in `pi_frontdoor_test.go` — drives `runPi` with `fakePiRuntimeOps{packageStatus: healthyPiPackageStatus() + subagentsRegistered: true}` and asserts the launch argv has 0 `--extension` and 0 `--skill` flags (pi-subagents fully delegated to discovery). A paired sub-test with `subagentsRegistered: false` asserts the explicit flags are present (current behavior). The dev-override case (cfg.repoRoot != "") with `subagentsRegistered: true` asserts the Spacedock extension is still appended but pi-subagents flags are absent. Fixture/CLI test — no live pi needed.
- **Existing tests (regression guard):** all existing `pi_frontdoor_test.go` tests use `healthyPiPackageStatus()` with `subagentsRegistered: false` (zero value) and assert the explicit flags ARE present — these continue to pass unchanged. `TestRunPi_DevOverridePassesSpacedockExtensionAndSkills` still gets 2 `--skill` flags (pi-subagents + Spacedock) when `subagentsRegistered: false`.
- **`piSpacedockPackageStatus` unit test (cheap):** extend or add a test that writes a `settings.json` with `npm:pi-subagents` and asserts the returned `piPackageStatus.subagentsRegistered` is true; a settings.json without it returns false.
- **Live smoke (optional, not blocking):** `spacedock pi` on this machine (where `npm:pi-subagents` is registered) starts without `Tool "subagent" conflicts` — the regression the entity exists to fix.

## Out of scope

- Changing pi-subagents' internal layout (index.ts re-export).
- The `spacedock install --host pi` install flow (already registers the package correctly).
- The dev-override Spacedock extension/skill block (unchanged).
- The safehouse wrap path (inner argv construction is the same; only the pi-subagents flags are gated).
- No documentation change needed: no user-visible command/output surface changes (the launch argv is an implementation detail; the help text and doctor output are unchanged).

## Stage Report: ideation

- DONE: Concrete approach specifying the registered-vs-dev-override detection (read ~/.pi/agent/settings.json packages, or test whether cfg.extensionPath resolves inside ~/.pi/agent/npm/node_modules) and gate the explicit --extension/--skill on the package NOT being registered — preserving the dev-override (cfg.repoRoot != "") purpose.
  Proposed approach section: extend `piSpacedockPackageStatus` to also check for "pi-subagents" while iterating settings.json packages (zero additional I/O), gate the explicit flags in `runPi` on `!check.packageStatus.subagentsRegistered`, preserve `subagentsRegistered` in the dev-override fallback. The value-AC this serves is AC-1 (no duplicate extension load). Simplest alternative: point the flag at `<pkg>/index.ts` so both specifiers coincide — more fragile because it depends on pi-subagents' internal layout (index.ts re-exporting src/extension/index.ts); a future refactor that moves or removes the re-export silently breaks the load.
- DONE: At least one value-measuring AC: with npm:pi-subagents in settings.json packages, spacedock pi registers the subagent tools exactly once and starts with no 'Tool "subagent" conflicts' error, measured against a registered-package baseline that can move the wrong way (the current unconditional-flag behavior fails).
  AC-1: the registered case drops the explicit `--extension`/`--skill` (one specifier via pi discovery, no collision); the unregistered case keeps them (the only load path). Baseline that can move wrong: the current unconditional-flag behavior always passes both specifiers → collision. Verified by a behavior test asserting 0 pi-subagents `--extension`/`--skill` flags when `subagentsRegistered: true`, 1 of each when false.
- DONE: Pair it with a mechanism AC for the flag-selection path (registered -> no explicit flag; dev-override -> explicit flag) exercised by a behavior test.
  AC-2: `TestRunPi_RegisteredSubagentsDedupesExtensionLoad` drives `runPi` against registered (`subagentsRegistered: true` → 0 flags) and unregistered (`false` → 1 flag) configs, with the dev-override Spacedock extension block unaffected in both cases. The falsifying change: flipping `subagentsRegistered` and seeing the wrong flag count.
- DONE: Expected surface and tolerance: net LOC change and files (internal/cli/pi.go + test), with observable-semantics declaration (no change to command grammar, dev-override path, or runtime behavior beyond dropping the duplicate load).
  Expected surface: `internal/cli/pi.go` + `internal/cli/pi_frontdoor_test.go`, +25 to +35 net LOC, tolerance ±10. Observable semantics: no change to command grammar, dev-override path, or runtime behavior beyond dropping the duplicate extension load when registered.
- DONE: Record the riskiest-mechanism spike result — whether settings.json packages reliably lists pi-subagents and whether extensionPath resolves inside ~/.pi/agent/npm/node_modules — exercised first, or "no spike needed: {proven mechanisms}".
  Spike exercised live: `~/.pi/agent/settings.json` packages contains `"npm:pi-subagents"` (confirmed); `<pkg>/index.ts` contains `export { default } from "./src/extension/index.ts"` (confirmed re-export); `<pkg>/src/extension/index.ts` exists (confirmed). Both detection mechanisms proven viable; settings.json approach selected (more precise). No further spike needed: `piSpacedockPackageStatus` settings.json read is a proven shipped mechanism.

### Summary

Fleshed out the ideation into a gate-ready design: extend the existing `piSpacedockPackageStatus` settings.json scan to also detect pi-subagents registration (zero additional I/O), gate the explicit `--extension`/`--skill` in `runPi` on `!subagentsRegistered`, and preserve the field through the dev-override fallback. The spike confirmed both detection paths (settings.json packages list `npm:pi-subagents`; extensionPath resolves inside the npm store) and the collision mechanism (index.ts re-exporting src/extension/index.ts → two specifiers). Two ACs: value (AC-1, no duplicate load when registered) and mechanism (AC-2, flag-selection path behavior test). Expected surface is internal/cli/pi.go + test, ~+30 net LOC, no observable semantic change beyond dropping the duplicate load. The simplest alternative (point the flag at `<pkg>/index.ts`) is named and rejected as more fragile (depends on pi-subagents' internal layout).
