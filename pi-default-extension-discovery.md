---
title: Pi default extension discovery for an isolated home
status: ideation
source: "Pi-UX carve, 2026-08-13: runtime-support first-contact friction"
score: 0.8
sprint: live-evidence-followups
sprint-readiness: ready
group: tooling
id: 3w1ncf1thj12aryvkf5gj1rd
gates:
    version: 1
    records:
        - id: gate:3w1ncf1thj12aryvkf5gj1rd:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:3w1ncf1thj12aryvkf5gj1rd-backlog-1
              briefing:
                id: briefing:3w1ncf1thj12aryvkf5gj1rd:backlog:attempt-1:revision-1
                digest: sha256:ea7fae5fe635ebedf7504254e21a6e10ac5b467da3049466a2ae2ffaf9f47856
                request-digest: sha256:5bd2c1bebf60fb5a6261e2d6ec9e2f2b54564577d606af9f5d87079b59d884fd
                room-ref: ./pi-default-extension-discovery/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:3w1ncf1thj12aryvkf5gj1rd:backlog:1
                briefing: briefing:3w1ncf1thj12aryvkf5gj1rd:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-14T06:43:06.589777Z"
                decision: approve
                reason: Captain approved backlog gate; advance to ideation for extension discovery.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:3w1ncf1thj12aryvkf5gj1rd:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:3w1ncf1thj12aryvkf5gj1rd-ideation-1
              briefing:
                id: briefing:3w1ncf1thj12aryvkf5gj1rd:ideation:attempt-1:revision-1
                digest: sha256:c2df2a3c3d2f448a2e15a6f34eace51298496bc45cfff53ab2e56acca1f9524f
                request-digest: sha256:a2c10784371baba39305a9926ec9246bf9e14e33c4faa7dc736e3727308c909e
                room-ref: ./pi-default-extension-discovery/review/ideation/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-14T14:30:13.698305Z"
                reason: Captain dropped Phase 0; not relevant right now. Entity moved to live-evidence-followups; ideation work retained on the body.
started: 2026-08-14T06:44:31Z
---

## Problem

An isolated Pi home does not auto-discover the `pi-subagents` / `pi-intercom`
extensions. The live harness works around this by hand-wiring paths: `piSubagentsPackageRoot`
requires `PI_SUBAGENTS_PACKAGE_ROOT` or falls back to `~/.pi/agent/npm/node_modules/pi-subagents`,
and the smoke launches pi with an explicit `--extension .../src/extension/index.ts`.
The operator's normal Pi install knows where these extensions live; an isolated
home forgets and has to be told.

`docs/runtime-support.md:147` names this as first-contact friction that is
supposed to be harness work ("an extension not auto-discovered in a temp home...
is harness work"), and the "assume it works" operating prompt expects auth and
package paths to be ironed out without a real blocker. Today the ironing is
manual and per-harness; a better default probe would make the isolated home
discover the operator's installed extensions the same way a normal home does.

## Visible value

A Pi runner in an isolated home resolves `pi-subagents` and `pi-intercom`
without the operator exporting `PI_SUBAGENTS_PACKAGE_ROOT` or the harness
hard-coding the `--extension` path. Measured against baseline: before, an
isolated-home Pi run with no `PI_SUBAGENTS_PACKAGE_ROOT` exported fails to find
the subagents extension; after, the same run discovers it from the operator's
installed package location and proceeds.

## Out of scope

- Changing where Pi itself stores or loads extensions.
- The `models.json` / `auth.json` copy (owned by `repair-pi-live-harness-parallelism-and-custom-model`, pnc).
- The intercom supervisor-talkback capability (archived spike `pi-intercom-runtime-capability-probe`).
- A new runtime, fixture, result format, or CI lane.

## Acceptance criteria

**AC-1 (VALUE) — An isolated-home Pi run finds the subagents extension with no env var exported.**

Verified by: an isolated-home Pi run that does NOT export `PI_SUBAGENTS_PACKAGE_ROOT`
resolves the `pi-subagents` extension (and `pi-intercom`, sibling package) from
the operator's installed package location, instead of erroring that the package
extension was not found. The baseline is the current `piSubagentsPackageRoot` fatal
path when the env var is unset and the fallback path is absent.

**AC-2 — The discovery is read from a real installed location, not a hard-coded fallback.**

Verified by: the probe reads the operator's actual installed extension/package
location (e.g. the Pi home's npm node_modules or package manifest), not a
hard-coded absolute path; a machine with extensions installed elsewhere is
discovered correctly.

**AC-3 — The explicit env var override still wins.**

Verified by: when `PI_SUBAGENTS_PACKAGE_ROOT` IS exported, it takes precedence
over discovery (existing harness behavior preserved); the `--extension` wiring
in the smoke stays valid.

**AC-4 — Offline and front-door smoke pass.**

Verified by: `gofmt`, `go vet -tags live ./internal/ensigncycle`,
`go build -tags live ./internal/ensigncycle`, and `TestLivePiFrontDoorSmoke`
pass with no `PI_SUBAGENTS_PACKAGE_ROOT` exported.

## Test plan

Use the offline `PiLiveEnv|PiIntercom|TestPiLive` unit tests first. Then one
`TestLivePiFrontDoorSmoke` run with the env var unset only when Pi work is
authorized. Preserve the explicit-override path.

## Notes

- General first-contact-friction reduction, not a journey repair.
- Coordinate with `repair-pi-live-harness-parallelism-and-custom-model` (pnc),
  which owns the `models.json`/`auth.json` copy into the isolated home; both
  touch `seedPiLiveAuth` / isolated-home setup.

## Proposed approach

### Discovery path a normal Pi home uses

Pi's package manager (`@earendil-works/pi-coding-agent/dist/core/package-manager.js`)
discovers extensions through two mechanisms:

1. **Settings.json `packages` array** — entries like `npm:pi-subagents` are
   resolved by `getManagedNpmInstallPath(source, "user")` to
   `join(agentDir, "npm", "node_modules", <name>)` (i.e.
   `~/.pi/agent/npm/node_modules/pi-subagents`). Each package's `package.json`
   `pi.extensions` field lists the extension entry points (e.g. `./index.ts`).
   `file:<path>` entries resolve as local packages.
2. **Auto-discovered extensions** — `addAutoDiscoveredResources` scans
   `~/.pi/agent/extensions/` for loose `.ts`/`.js` files.

The existing test-harness fallback
`filepath.Join(home, ".pi", "agent", "npm", "node_modules", "pi-subagents")`
(in `piSubagentsPackageRoot` at `pi_live_runner_test.go:258`) mirrors exactly
`getManagedNpmInstallPath` for user scope — the path Pi's package manager
computes for `npm:pi-subagents` in `settings.json`.

The spacedock `pi` command (`internal/cli/pi.go:piRuntimeConfigFromEnv`, line
~505) reads `PI_SUBAGENTS_PACKAGE_ROOT` from the env, or falls back to the same
`join(home, ".pi", "agent", "npm", "node_modules", "pi-subagents")` path, then
passes `--extension <pkg>/src/extension/index.ts` to the pi binary. In the test,
`HOME` is set to `cleanHome` (a temp dir), so the fallback resolves to the temp
dir — which has no npm install. The test currently works around this by setting
`PI_SUBAGENTS_PACKAGE_ROOT` in `piLiveEnv` (`pi_live_controls_test.go:64`).

### `piDefaultExtensionRoots` helper

A new test helper that reads the operator's **real** Pi home (the test-time
`os.Getenv("HOME")`, not the isolated `cleanHome`) to discover where
`pi-subagents` and `pi-intercom` are actually installed — dynamically, not via a
hard-coded absolute path. The resolution mirrors Pi's own package-manager logic
(`resolveSettingsPackageRoot` in `pi.go:721`):

1. Read the operator's real `~/.pi/agent/settings.json` `packages` array.
2. For `npm:pi-subagents` / `npm:pi-intercom` entries, resolve to
   `~/.pi/agent/npm/node_modules/<name>`.
3. For `file:<path>` entries, read the package's `package.json` `name` field and
   match `pi-subagents` / `pi-intercom`.
4. If not found in settings.json, probe `~/.pi/agent/npm/node_modules/<name>/`
   (the managed npm install path) and `~/.pi/agent/extensions/subagent/`
   (the `install.mjs` clone location) as fallbacks.
5. Return the resolved roots (subagents root, intercom root).

The explicit env-var override (`PI_SUBAGENTS_PACKAGE_ROOT` /
`PI_INTERCOM_PACKAGE_ROOT`) still wins: the helper returns the env-var value
immediately when set, without probing.

### Seeding the isolated home

The test fixture (`newPiLiveSmokeFixture` / `newPiSharedLiveDriver`) calls the
helper, then seeds `cleanHome/.pi/agent/npm/node_modules/` with symlinks to the
discovered package roots. This makes `piRuntimeConfigFromEnv`'s fallback path
(`join(home, ".pi", "agent", "npm", "node_modules", "pi-subagents")`) resolve
through the symlink — without setting `PI_SUBAGENTS_PACKAGE_ROOT` in the env.

A new `seedPiDefaultExtensions(t, cleanHome, realHome)` function performs the
symlinking, called from the same setup area as `seedPiLiveAuth` but as a
separate function (different files: auth copy vs. extension symlinks — no
collision with the pnc's `models.json`/`auth.json` copy).

`piLiveEnv` drops the hard-coded `PI_SUBAGENTS_PACKAGE_ROOT` / `PI_INTERCOM_PACKAGE_ROOT`
env-var assignments (or makes them conditional on the env var being set at test
time, preserving the explicit-override path).

## Expected surface and tolerance

**Files (test-harness only, no production code changes):**

- `internal/ensigncycle/pi_live_runner_test.go`:
  - New `piDefaultExtensionRoots(t, realHome) (subagentsRoot, intercomRoot string)`
    helper (~40-60 lines) — reads `settings.json` packages, resolves via
    `resolveSettingsPackageRoot`-equivalent logic, probes npm install +
    extensions dir as fallbacks.
  - New `seedPiDefaultExtensions(t, cleanHome, realHome)` (~15 lines) —
    symlinks discovered roots into `cleanHome/.pi/agent/npm/node_modules/`.
  - `newPiLiveSmokeFixture` calls `seedPiDefaultExtensions` after `seedPiLiveAuth`
    (~2 lines added).
  - `piSubagentsPackageRoot` updated to delegate to `piDefaultExtensionRoots`
    instead of the hard-coded fallback.

- `internal/ensigncycle/pi_live_controls_test.go`:
  - `piLiveEnv` signature: `piSubagentsRoot` parameter becomes optional (empty =
    rely on seeded discovery); the `PI_SUBAGENTS_PACKAGE_ROOT` /
    `PI_INTERCOM_PACKAGE_ROOT` env-var lines become conditional (set only when
    the root is non-empty).
  - `piIntercomPackageRoot` updated to use `piDefaultExtensionRoots` when the
    env var is unset, instead of deriving from the subagents root.
  - `TestPiLiveEnvDropsForeignRuntimeMarkers` updated: when
    `PI_SUBAGENTS_PACKAGE_ROOT` is not in the source env, it should not appear in
    the target env (the assertion at line 73 for `/parent/package` stays valid
    because that test sets the env var explicitly).

- `internal/ensigncycle/pi_shared_live_runner_test.go`:
  - `newPiSharedLiveDriver` calls `seedPiDefaultExtensions` and drops the
    `piSubagentsPackageRoot(t)` argument to `piLiveEnv` (or passes empty).

**Tolerance:**
- No changes to `internal/cli/pi.go` or any non-test file.
- The existing `--extension` wiring in `runPi` stays valid (AC-3).
- The explicit env-var override must still win when set.
- Offline unit tests (`TestPiLiveEnvDropsForeignRuntimeMarkers`,
  `TestPiLiveEnvScrubsAmbientPiSubagentMarkers`,
  `TestPiIntercomPackageRootDefaultsBesideSubagents`,
  `TestPiLiveSmokePromptRequiresExactStageReportHeading`) must still pass.
- `piLiveEnv`'s env-scrub contract (dropping foreign `PI_SUBAGENT_*` markers)
  is preserved — the scrub list stays, only the *additive* env-var lines change.

**Semantic changes:** none expected — test harness only. No production behavior
changes; the spacedock `pi` command's resolution logic is untouched.

## Test plan

Use the offline `PiLiveEnv|PiIntercom|TestPiLive` unit tests first. Then one
`TestLivePiFrontDoorSmoke` run with the env var unset only when Pi work is
authorized. Preserve the explicit-override path.

### Per-AC proof plan

- **AC-1 (isolated-home Pi run finds subagents with no env var):**
  `TestLivePiFrontDoorSmoke` run with `PI_SUBAGENTS_PACKAGE_ROOT` **unset** in
  the test process env. The seeded symlinks make `piRuntimeConfigFromEnv`'s
  fallback resolve; the smoke passes (subagent tool available, stage report
  written, boot contract graded). Falsifiable: remove the symlink seeding and
  the smoke fails with "pi-subagents package extension not found".

- **AC-2 (discovery reads real installed location, not hard-coded path):**
  New offline unit test `TestPiDefaultExtensionRootsReadsSettings` creates a
  fake Pi home with `settings.json` containing `"npm:pi-subagents"` and a mock
  `node_modules/pi-subagents/` dir; asserts the helper resolves the root from
  the settings entry, not from a hard-coded path. A second case uses
  `file:<path>` with a `package.json` name `pi-subagents` and asserts
  resolution from the `file:` entry. Falsifiable: point settings.json at a
  different location and the helper returns that location, not the default.

- **AC-3 (explicit env var override still wins):**
  `TestPiDefaultExtensionRootsEnvOverride` sets `PI_SUBAGENTS_PACKAGE_ROOT` to
  a sentinel path and asserts the helper returns it immediately without
  probing settings.json or the npm dir. The existing
  `TestPiLiveEnvDropsForeignRuntimeMarkers` assertion for
  `PI_SUBAGENTS_PACKAGE_ROOT: /parent/package` stays valid (that test sets the
  env var explicitly). Falsifiable: set the env var and the helper ignores
  settings.json entirely.

- **AC-4 (offline and front-door smoke pass):**
  `gofmt -l ./internal/ensigncycle`, `go vet -tags live ./internal/ensigncycle`,
  `go build -tags live ./internal/ensigncycle`, and
  `TestLivePiFrontDoorSmoke` (live, env var unset) all pass. Offline unit tests
  pass with no live gate. Falsifiable: any of these commands exits non-zero.

## Stage Report: ideation

- DONE: Investigate how a normal (non-isolated) Pi home discovers installed extensions
  Read `piSubagentsPackageRoot` (`pi_live_runner_test.go:251-262`) and `piIntercomPackageRoot` (`pi_live_controls_test.go:121-125`); traced the fallback to Pi's `getManagedNpmInstallPath` (`package-manager.js:1670`) which resolves `npm:<name>` settings entries to `join(agentDir, "npm", "node_modules", <name>)`; confirmed the fallback `~/.pi/agent/npm/node_modules/pi-subagents` mirrors exactly this path.
- DONE: Determine what real installed location the fallback mirrors
  The fallback mirrors Pi's managed npm install path for user-scope packages (`getManagedNpmInstallPath` → `join(agentDir, "npm", "node_modules", source.name)`); verified `~/.pi/agent/npm/node_modules/pi-subagents/src/extension/index.ts` exists on this machine and `pi-subagents/package.json` declares `pi.extensions: ["./index.ts"]`.
- DONE: Propose a `piDefaultExtensionRoots` helper that reads the operator's real installed location
  Proposed helper reads `~/.pi/agent/settings.json` packages (resolving `npm:` and `file:` entries via the same logic as `pi.go:resolveSettingsPackageRoot`), with npm-install-dir and extensions-dir probes as fallbacks; explicit `PI_SUBAGENTS_PACKAGE_ROOT` / `PI_INTERCOM_PACKAGE_ROOT` override still wins.
- DONE: Record proposed approach, expected surface, semantic changes, and per-AC proof plan into entity body
  Added `## Proposed approach`, `## Expected surface and tolerance`, and expanded `## Test plan` with per-AC proof plan; ACs and other existing sections left unedited.
- DONE: Coordinate with the merged pnc `seedPiLiveAuth` seam
  The extension-discovery seeding is a separate `seedPiDefaultExtensions` function called from the same setup area (`newPiLiveSmokeFixture`) after `seedPiLiveAuth`; different files (symlinks vs. auth/models copy) — no collision with pnc's `models.json`/`auth.json` copy.

### Summary

Investigated Pi's extension discovery chain end-to-end: `settings.json` packages → `getManagedNpmInstallPath` → `~/.pi/agent/npm/node_modules/<name>` → `package.json` `pi.extensions`. The test-harness fallback mirrors this path but is hard-coded and requires `PI_SUBAGENTS_PACKAGE_ROOT` to be set for isolated homes. Proposed `piDefaultExtensionRoots` — a test helper that reads the operator's real `settings.json` packages to dynamically discover the installed package roots, plus `seedPiDefaultExtensions` to symlink them into the isolated `cleanHome` so `piRuntimeConfigFromEnv`'s fallback resolves without the env var. No production code changes; explicit override preserved; per-AC proof plan recorded.
