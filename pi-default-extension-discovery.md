---
title: Pi default extension discovery for an isolated home
status: ideation
source: "Pi-UX carve, 2026-08-13: runtime-support first-contact friction"
score: 0.8
sprint:
sprint-readiness: defer
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
            - id: gate-attempt:3w1ncf1thj12aryvkf5gj1rd-ideation-2
              briefing:
                id: briefing:3w1ncf1thj12aryvkf5gj1rd:ideation:attempt-2:revision-1
                digest: sha256:83c51381a27be6ed4dc266d622644f72f322488f0f78eb66c28dd7591741097f
                room-ref: '@review/ideation/briefing-2'
              resolution:
                type: Resolution
                id: resolution:spacedock:3w1ncf1thj12aryvkf5gj1rd:ideation:2
                briefing: briefing:3w1ncf1thj12aryvkf5gj1rd:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-10-03T17:04:58.810792Z"
                decision: revise
                reason: 'Captain: "no, send it back to reword it so we know that was a bad gate attempt". The bound ideation attempt-2 summary was a bad gate attempt: it described the change in invented vocabulary - ''substrate packages'' for the required pi-subagents and pi-intercom extensions, and ''hand-wired'' for three nameable manual steps - so the captain could not tell what the change does or why it matters. Reword the artifact in plain operator language and record in the body that attempt-2 was rejected on wording, so the rejection is visible on the record. Scope, surface, tolerance and criteria are not reopened by this rejection.'
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

**M1 proposed refinements for the captain's gate (not yet approved):** retain
AC-1–AC-4 below as the earlier baseline. Proposed AC-1 requires BOTH
`pi-subagents` and `pi-intercom` to load through Pi package discovery with BOTH
`PI_SUBAGENTS_PACKAGE_ROOT` and `PI_INTERCOM_PACKAGE_ROOT` absent, no
harness-supplied substrate extension paths, and the Spacedock package/ensign
skill still loaded. A successful explicit fallback is not AC-1 evidence.
Proposed AC-2 covers independently located packages and the actual agentDir
settings, including supported absolute/relative local package entries, not
just the conventional npm directory. Proposed AC-3 preserves independent
explicit overrides for both roots and the retained fallback as a separate
compatibility check. Proposed AC-4 adds ordinary (non-live-tagged) helper tests
and a no-model package-loader check to the existing authorized front-door
smoke; the latter still owes durable dispatch/report/commit evidence. See
“M1 staff-review fold” below for falsifiers and the corrected setup contract.
These refinements supersede the earlier proof's symlink-only and single-env
interpretation only if approved at the captain's gate.

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

## M1 staff-review fold — proposed for the captain's gate

This is a correction to the retained ideation, not a new design. All AC,
surface, and harness-semantic changes in this section are **proposed for the
captain's gate**. The earlier body/report remain as history; this proposal
replaces their symlink-only setup, live-tagged helper placement, unconditional
sibling-root expectation, and single-variable/fallback-based proof. No
implementation or live journey is claimed by this ideation fold.

### One isolated-home setup contract (proposed)

- Capture `realHome` from the parent HOME before constructing the child env.
  Capture the real agent directory from the parent's `PI_CODING_AGENT_DIR`
  when explicitly set, otherwise `realHome/.pi/agent`. Read discovery settings
  and installed roots from that directory; never read them from clean HOME.
  Treat real-home packages/settings as read-only; do not copy unrelated packages.
- Allocate `cleanHome` once; define `piHome = cleanHome/.pi/agent` and pass that
  SAME path as `PI_CODING_AGENT_DIR`. This deliberately aligns Pi's agentDir
  loader with the launcher's existing HOME-based package-root probes, without
  changing production resolution. Sessions remain separately isolated.
- Resolve subagents and intercom independently from explicit overrides, then
  real-agent settings/installed-package probes already proposed. For supported
  local sources, recognize absolute and agentDir-relative entries and match
  `package.json` names. Do not mistake a non-sibling intercom install for an
  error, or infer its root from a subagents override. The earlier `file:`
  example is not supported by the exercised Pi 1.0.0 loader; do not emit it.
- `seedPiDefaultExtensions(t, piHome, roots)` links BOTH packages under
  `piHome/npm/node_modules/{pi-subagents,pi-intercom}`. Register BOTH
  `npm:pi-subagents` and `npm:pi-intercom` in `piHome/settings.json`'s `packages`
  array. Preserve the Spacedock package as one absolute checkout-path entry
  (`repo`, not `"file:"+repo`), and preserve any other intentionally seeded
  settings. Root symlinks alone do not register resources with Pi.
- Compose settings once in setup (or merge additions); neither fixture may
  overwrite the substrate registrations with the old repo-only settings write.
  `newPiLiveSmokeFixture` and `newPiSharedLiveDriver` use the same contract.
  `seedPiLiveAuth` still owns auth/models copying into this same piHome (pnc);
  no credential/model policy change is proposed here.
- Capture overrides separately from discovered roots. In default mode scrub
  BOTH package-root variables and do not re-add them, including as empty
  assignments. In explicit mode each nonempty operator override independently
  wins root selection and is forwarded unchanged; the other package still
  discovers normally. Preserve foreign-runtime and `PI_SUBAGENT_*` scrubbing.
  No harness substrate `--extension` paths or additional-extension SDK paths
  are allowed in default-discovery proof. The existing launcher's unregistered
  explicit fallback remains valid only as a separate compatibility path.

This serves proposed AC-1/AC-2. Alternatives: symlinks alone leave the loader
unaware of packages; registration in an unrelated piHome leaves HOME-based
preflight probes empty; root env injection masks the default path. Aligning
piHome with clean HOME is smaller than changing the launcher. Explicit
package-root precedence serves proposed AC-3 without retaining the sibling
assumption for a different package.

### Non-live seam, old test contract, and bounded surface (proposed)

- Add `internal/ensigncycle/pi_default_extensions_test.go` WITHOUT a `live`
  build constraint. Put `piDefaultExtensionRoots`, `seedPiDefaultExtensions`,
  their small shared setup representation if needed, and deterministic tests
  there. A definition in `pi_live_runner_test.go` is invisible to ordinary
  Go tests; ordinary tests must not depend on a live-tagged definition.
- Update `pi_live_controls_test.go`: distinguish explicit overrides from
  discovered roots in `piLiveEnv`/`piLiveEnvForAuth`, preserve env scrubbing,
  and revise `piIntercomPackageRoot` to independent discovery. Replace
  `TestPiIntercomPackageRootDefaultsBesideSubagents` with
  `TestPiIntercomPackageRootDiscoversIndependently`: settings point intercom
  at a non-sibling directory while a plausible sibling exists; require the
  settings root. Keep the independent intercom override case. This is a
  proposed CHANGE to the old test contract, not a promise it passes unchanged.
  `TestPiLiveEnvDropsForeignRuntimeMarkers` currently expects the explicit
  `/target/package` argument, not `/parent/package`; retain that precedence
  and the explicit `/parent/intercom` override in its compatibility case.
- Update `pi_live_runner_test.go` and `pi_shared_live_runner_test.go`: align
  cleanHome/piHome, compose package registration, remove automatic env root
  injection, and use the non-live helper. `piSubagentsPackageRoot` must no
  longer assert a source-layout `.ts` entry: it delegates root discovery only.
- **Estimate net LOC change: +220, across 4 test files** (approximately +260
  insertions / -40 deletions; proposed tolerance net +160..+280, no extra
  files without reapproval). This replaces the earlier three-file estimate.
  The non-live seam serves AC-4; putting it in production or behind `live`
  adds unnecessary scope or defeats ordinary tests.
- Observable semantics allowed to change: isolated harness directory layout,
  package registrations, absence of default root env injections, independent
  intercom discovery, and replacement of the unsupported harness `file:`
  prefix with an absolute path. No command grammar, stored entity format,
  authority, shipped extension/skill, production runtime, or CI-lane change.
  No launcher/docs changes are proposed; no user-facing documentation diff
  is needed for this harness-only correction.
- `mc` (`pi-doctor-probes-stale-subagents-layout`) remains the owner of
  manifest extension/bridge entry resolution and stale launcher probes.
  This task finds and registers package ROOTS, letting Pi load manifest
  entries; it must not implement another entry resolver or restore `.ts`
  assumptions. A front-door preflight failure on current compiled packages
  is an mc dependency to report, not permission to change `internal/cli/pi.go`.
  Any broader launcher change requires separate scope approval.

### Proposed proof refinements

Existing primary proof owner remains `TestLivePiFrontDoorSmoke` for end-to-end
value, with durable entity report/commit and boot-contract grading. Its
successful explicit fallback does NOT prove default discovery. Proposed checks:

- AC-1/AC-2, deterministic Go setup tests (seconds): fake real HOME and custom
  real agentDir, two packages at non-sibling local roots, and separate clean
  HOME. Inspect the constructed child env and follow both symlinks; parse
  settings to require both npm registrations plus exactly one Spacedock path.
  Distinct falsifiers: point PI_CODING_AGENT_DIR elsewhere; remove one
  registration; replace settings with repo-only content; inject either root
  env var; select the sibling decoy. Settings-source cases use non-default
  local directories so a hardcoded npm fallback cannot satisfy the test.
- AC-1/AC-4, no-model real Pi loader check (seconds, Pi install required): feed
  the seeded settings to Pi's `DefaultResourceLoader` without additional
  extension paths. Require BOTH loaded package entries and BOTH `subagent`
  and `intercom` tools, plus Spacedock/ensign. Negative control retains root
  symlinks but removes both substrate registrations and loses both tools.
  This proves supported loader behavior, not just Go's JSON writer. Do not
  add a CI lane or make ordinary Go tests require a local Pi install.
- AC-3, deterministic override tests (seconds): each override alone and both
  together win their respective discovery conflicts; the non-overridden root
  still comes from settings. Preserve marker-scrub tests. Falsifiers: override
  ignored, intercom derived from subagents, or ambient child marker retained.
  The retained explicit fallback is exercised separately, never counted as
  the default-discovery run.
- AC-4/front-door, authorized live run (existing live cost/budget): unset BOTH
  package-root variables in the test process and resulting child env; capture
  launch argv proving no substrate extension path was supplied, including by
  the launcher's fallback. Require both loaded substrate tools in addition to
  the existing report/commit evidence. Run gofmt, ordinary helper tests,
  `go vet -tags live ./internal/ensigncycle`, and
  `go build -tags live ./internal/ensigncycle` before the live proof. A live
  smoke remains pending authorization and mc's entry-resolution readiness.

### Exercised risk evidence and routed findings

A no-model SDK loader spike on **Pi 1.0.0**, Node **v24.13.1**, installed
`pi-subagents` **0.75.0** / `pi-intercom` **0.16.0** exercised real package
loading with clean HOME, aligned agentDir, both root variables deleted, no
additional extension paths, an empty cwd, and parent `PI_SUBAGENT_*` markers
scrubbed. It did not launch a model, dispatch a worker, or invoke intercom.

| Settings `packages` entries | Observed result |
| --- | --- |
| `["file:/Users/clkao/git/spacedock-research/spacedock-v1"]` | 0 extensions, 0 skills, 0 tools; no loader error |
| `["/Users/clkao/git/spacedock-research/spacedock-v1"]` | Spacedock extension, 11 skills including ensign; no substrate tools despite both symlinks |
| `["repo"]` and separately `["./repo"]`, with agentDir/repo linked to that checkout | Same Spacedock extension and 11 skills; both relative forms work |
| `["/Users/clkao/git/spacedock-research/spacedock-v1","npm:pi-subagents","npm:pi-intercom"]` | 3 extensions (Spacedock, subagents/index.js, intercom/index.ts), 14 skills, tools subagent/bg_wait/subagents_enable/intercom, no loader errors |

The initial spike asserted ensign availability for `file:` and failed; the
corrected absolute-path retry succeeded. Another preliminary assertion found
that inherited `PI_SUBAGENT_CHILD=1` suppresses the subagent tool even when
its extension loads; the final spike preserved the existing harness scrub
contract and passed every assertion. This does not justify changing that
contract or count as live front-door evidence.

The `file:` PREFIX is unsupported in the exercised loader; local paths are
supported. FO approved recording the harness correction, subject to the
captain's gate. Search of package-entry writers, front-door/install code,
docs and skills found only two affected writers:
`pi_live_runner_test.go:52` and `pi_shared_live_runner_test.go:28`.
`runInitWithPi` (`internal/cli/pi.go:436-440`) passes raw pluginDir or the
existing git package source to `pi install`, not a constructed `file:` entry;
no user-facing registration example emitting this prefix was found.

Related reader mismatch, routed to FO/launcher owner without a fix here:
`internal/cli/pi.go:944-977` claims/supports stripping `file:` but rejects bare
relative names such as `repo`, which the actual Pi loader accepts. This is
not mc's manifest-entry-resolution work and does not expand this entity.
FO reports real operator settings use `../../git/spacedock-research/spacedock-v1`
and a second `git:github.com/spacedock-dev/spacedock` registration. The duplicate
Spacedock warning is a separate two-entry consequence, not this prefix failure;
this fixture seeds one Spacedock registration and does not edit operator state.

Reproduction inputs are declared, not silent machine dependencies: Node and
Pi SDK installed globally via npm, both named packages installed in the real
home's managed npm directory (adjust the source roots if installed elsewhere),
and this checkout as REPO_ROOT. Save the following as a temporary `.mjs` file
outside `.pi/extensions/`, then run:

```bash
PI_SDK="$(npm root -g)/@earendil-works/pi-coding-agent/dist/index.js" \
REAL_HOME="$HOME" REPO_ROOT="$PWD" node /path/to/spike.mjs
```

```javascript
import fs from 'node:fs'; import os from 'node:os'; import path from 'node:path'; import assert from 'node:assert/strict'; import {pathToFileURL} from 'node:url';
const base = fs.mkdtempSync(path.join(os.tmpdir(), 'pi-discovery-M1-'));
const home = path.join(base, 'home'), agentDir = path.join(home, '.pi/agent'), cwd = path.join(base, 'cwd');
fs.mkdirSync(path.join(agentDir, 'npm/node_modules'), {recursive: true}); fs.mkdirSync(cwd);
for (const name of ['pi-subagents','pi-intercom']) fs.symlinkSync(path.join(process.env.REAL_HOME, '.pi/agent/npm/node_modules', name), path.join(agentDir, 'npm/node_modules', name));
fs.symlinkSync(process.env.REPO_ROOT, path.join(agentDir, 'repo'));
process.env.HOME = home; process.env.PI_CODING_AGENT_DIR = agentDir; process.env.PI_OFFLINE = '1';
delete process.env.PI_SUBAGENTS_PACKAGE_ROOT; delete process.env.PI_INTERCOM_PACKAGE_ROOT;
for (const key of Object.keys(process.env)) if (key.startsWith('PI_SUBAGENT_')) delete process.env[key];
const {DefaultResourceLoader, VERSION} = await import(pathToFileURL(process.env.PI_SDK).href);
try {
  for (const [label, source, registered, expected] of [['file-prefix','file:'+process.env.REPO_ROOT,false,false], ['absolute',process.env.REPO_ROOT,false,true], ['bare-relative','repo',false,true], ['dot-relative','./repo',false,true], ['registered',process.env.REPO_ROOT,true,true]]) {
    const packages = [source, ...(registered ? ['npm:pi-subagents','npm:pi-intercom'] : [])];
    fs.writeFileSync(path.join(agentDir, 'settings.json'), JSON.stringify({packages}));
    const loader = new DefaultResourceLoader({cwd, agentDir}); await loader.reload();
    const result = loader.getExtensions(), skills = loader.getSkills().skills;
    const tools = result.extensions.flatMap(e => [...e.tools.keys()]);
    console.log(JSON.stringify({version: VERSION, label, packages, extensions: result.extensions.map(e => e.path), tools, errors: result.errors, skills: skills.map(s=>s.name)}));
    assert.equal(result.errors.length, 0); assert.equal(tools.includes('subagent'), registered); assert.equal(tools.includes('intercom'), registered);
    assert.equal(skills.some(s=>s.name==='ensign'), expected);
  }
} finally { fs.rmSync(base, {recursive: true, force: true}); }
```

## Stage Report: ideation (cycle 2)

- DONE: Define one isolated-home setup contract covering the real home, the agent directory, the clean HOME, and package registration, with the explicit-override behavior named.
  M1 contract aligns piHome with cleanHome/.pi/agent, captures the real source agentDir before isolation, independently preserves explicit root overrides, and composes settings once with both substrates and one Spacedock entry; all changed AC/surface semantics are proposed for the captain's gate.
- DONE: Register and load BOTH the pi-subagents and pi-intercom packages through Pi's supported discovery path while preserving the Spacedock package entry; the proof must hold with both package-root variables absent and no substrate extension path supplied by the harness, and it must not count the retained explicit fallback.
  Pi 1.0.0 no-model loader spike with both root variables absent loaded all 3 extensions, subagent/intercom tools and ensign; removing npm registrations lost both substrate tools despite symlinks. Unsupported file: baseline loaded 0 extensions/skills; absolute and both relative forms loaded Spacedock. Reproducer and observed results are in M1 risk evidence, not a claim that harness implementation/live smoke shipped.
- DONE: Name the non-live helper seam, since a live-tagged definition is invisible to ordinary Go tests, and name the changed old sibling-root test contract; keep mc as entry-resolution owner and do not propose a broader launcher change without separate scope approval.
  Proposed pi_default_extensions_test.go has no live constraint; replaces TestPiIntercomPackageRootDefaultsBesideSubagents with independent settings-root/decoy coverage. mc retains manifest-entry ownership; the separate launcher reader mismatch was routed to FO without a fix or scope expansion.
- DONE: Validate the existing focused offline baseline and document bounded validation limits.
  `go test ./internal/ensigncycle -run 'PiLiveEnv|PiIntercom|TestPiLive' -count=1` passed; these existing checks detect env-scrub/old sibling-root regressions, not the proposed discovery contract. `gofmt -w ./cmd ./internal` ran; its unrelated pre-existing formatting delta was undone. `go test ./...` timed out at 120s; own race run was stopped after >3m with cli/ensigncycle/status still running, so neither full suite is claimed green.
- SKIPPED: Implement proposed harness changes and run authorized live front-door smoke.
  This is a staff-review ideation fold; captain approval, mc readiness, and live authorization remain implementation/validation prerequisites. No code, launcher, docs, operator settings, or YAML frontmatter change is delivered.

- DONE: AC-1 proof plan cited: M1 “Proposed proof refinements” requires both packages/tools and ensign through registered discovery with both root variables absent and no substrate extension paths; the recorded no-model loader spike supports registration, while authorized `TestLivePiFrontDoorSmoke` remains pending and explicit fallback does not count.
- DONE: AC-2 proof plan cited: M1 deterministic setup tests use a custom real agentDir and independent non-sibling local roots; selecting the sibling decoy or a hard-coded npm root falsifies discovery from actual settings (supported absolute/relative sources).
- DONE: AC-3 proof plan cited: M1 deterministic override tests exercise each root override alone and both together against conflicting settings, preserve marker scrubbing, and check the retained explicit fallback separately; ignoring an override or deriving intercom from subagents falsifies precedence.
- DONE: AC-4 proof plan cited: M1 requires ordinary non-live helper tests, gofmt, live-tagged vet/build, the no-model loader check, and authorized front-door smoke with durable report/commit evidence; cycle 2 records only the focused offline baseline and loader spike as passed, not full-suite or live completion. All M1 refinements remain proposed for the captain's gate.

### Summary

Folded M1 into the retained design with one composable isolated-home contract, a non-live helper seam, explicit old-test replacement, and falsifiable default-discovery proof for both packages. Exercised Pi's actual loader (not a model session) to show registration is necessary and that the harness's file: prefix must become a supported absolute checkout entry; all AC/surface refinements await the captain's gate, and mc remains the entry-resolution owner.

### Report repair validation

`spacedock status --read docs/dev/.spacedock-state/pi-default-extension-discovery.md --ac-scan --json --workflow-dir docs/dev`

```json
{"command":"read","stage":"ideation","acs":[{"id":"AC-1","line":"100","unevidenced":"false","citations":[{"line":"531","text":"- DONE: AC-1 proof plan cited: M1 “Proposed proof refinements” requires both packages/tools and ensign through registered discovery with both root variables absent and no substrate extension paths; the recorded no-model loader spike supports registration, while authorized `TestLivePiFrontDoorSmoke` remains pending and explicit fallback does not count."}]},{"id":"AC-2","line":"108","unevidenced":"false","citations":[{"line":"532","text":"- DONE: AC-2 proof plan cited: M1 deterministic setup tests use a custom real agentDir and independent non-sibling local roots; selecting the sibling decoy or a hard-coded npm root falsifies discovery from actual settings (supported absolute/relative sources)."}]},{"id":"AC-3","line":"115","unevidenced":"false","citations":[{"line":"533","text":"- DONE: AC-3 proof plan cited: M1 deterministic override tests exercise each root override alone and both together against conflicting settings, preserve marker scrubbing, and check the retained explicit fallback separately; ignoring an override or deriving intercom from subagents falsifies precedence."}]},{"id":"AC-4","line":"121","unevidenced":"false","citations":[{"line":"534","text":"- DONE: AC-4 proof plan cited: M1 requires ordinary non-live helper tests, gofmt, live-tagged vet/build, the no-model loader check, and authorized front-door smoke with durable report/commit evidence; cycle 2 records only the focused offline baseline and loader spike as passed, not full-suite or live completion. All M1 refinements remain proposed for the captain's gate."}]}]}
```
