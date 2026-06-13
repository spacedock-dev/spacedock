---
id: tes9th8ncq1p01am9qk7eex4
title: install refresh leaves a stale plugin + no upgrade path is surfaced (the 0.19.8 thing)
status: validation
source: "Captain field report 2026-06-09 — first real 0.20.0 install. `spacedock install --host codex` on a tag-fresh 0.20.0 binary returned `OK: spacedock binary 0.20.0 and plugin 0.19.8 are compatible.` The plugin stayed at 0.19.8 (older than BOTH main HEAD 0.20.0 and next HEAD 0.19.9), and nothing told the user a newer plugin exists, how to get it, or whether the front door upgrades for them."
started: 2026-06-13T04:08:47Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-install-refresh-and-upgrade-hint
issue:
sprint: 0201-post-flip-release-model
group: release-model
sprint-readiness: ready
---

A user with a tag-fresh **0.20.0 binary** ran `spacedock install --host codex` to upgrade, and ended up still on plugin **0.19.8** with a `compatible` message and no path forward. Two separable facets, both real, surfaced by the first 0.20.0 install. **Ideation firmed the root cause and the desired end-state per host; both fixes are small wiring changes, not seam rewrites.**

## Problem

**Facet 1 — `spacedock install --host codex` never re-installs a plugin that is already present (correctness).** Root-caused, not assumed. The `Install` seam itself works — `installArgvSequence`/`codexInstallArgvSequence` (`internal/cli/host_exec.go`) run a cleanup-then-pin sequence (`plugin uninstall`/`remove` + `marketplace remove`, both tolerated, then `marketplace add <source> [--ref <branch>]`, then `plugin install`/`add`), and `TestUpgradeFromStaleMovesToGreen` + `TestCodexPluginInstallIsHostNative` already prove it lands the source HEAD against the real host CLIs. The bug is one level up, in `runInit` (`internal/cli/init.go`): the **claude** arm calls `ops.Install(...)` unconditionally (refresh-on-present), but the **codex** arm does NOT. The codex arm resolves the existing manifest, sees `resolved != ""`, runs `RunDoctor` (which prints "compatible", exit 0), and `return`s at the `if code != 0 || resolved != "" || check` guard — **it never reaches an install call**. So `spacedock install --host codex` on an already-present 0.19.8 plugin is a doctor-only no-op: it reports compatible and leaves the stale plugin in place. That is exactly the captain's field report (`--host codex`). The codified expression of this no-op is the current spec test `TestInitCodexInstallReadiness/compatible-installed`, which asserts the present-and-compatible codex arm prints the OK line and bans any install prose.

**Facet 2 — no upgrade is ever hinted, and the front door is silent on Compatible (UX).** The `Compatible` verdict (`internal/contract/contract.go`) compares CONTRACT versions only (binary contract `1` vs plugin `requires-contract: ">=1,<2"`), not display semvers — so "binary 0.20.0, plugin 0.19.8 are compatible" is correct by design. But it prints two different version numbers side by side and says "you're fine," with no nudge that a newer plugin (matching the binary) exists or how to get it. Worse, the front door's gate (`gateHost`) is **silent** on `Compatible` — `res.Message` is only emitted by doctor/install, never by `spacedock claude`/`codex` — so a launching user with a behind plugin sees nothing at all. The binary already holds both display semvers locally (its own `cli.Version` and the plugin manifest `version`), so the skew is detectable with no network fetch. A contract-compatible-but-behind plugin should surface an honest, opt-in upgrade hint in both doctor and the front-door gate.

## Spike — riskiest mechanism (DONE, throwaway)

The riskiest unverified path was the **codex present→newer refresh**: claude's already-present→newer transition is proven by `TestUpgradeFromStaleMovesToGreen`, but codex's was not (the existing codex behavioral test only proves a fresh install lands). A throwaway test (`spike_codex_refresh_test.go`, since deleted) installed a 0.0.1 codex plugin into an isolated `CODEX_HOME`, then ran `execHost.Install("codex", <0.0.2-marketplace>, "")` with the plugin already present, and observed the resolved cache manifest:

- after seed install: `…/cache/spacedock/spacedock/0.0.1/.codex-plugin/plugin.json` version=**0.0.1**
- after refresh install: `…/cache/spacedock/spacedock/0.0.2/.codex-plugin/plugin.json` version=**0.0.2**

Result: the codex `Install` seam DOES refresh an already-present plugin to source HEAD. The mechanism is sound; the only missing piece is the `runInit` codex arm calling it. This seeds the implementation's first test (assert `runInit --host codex` on a present plugin records an `ops.Install` call, then a live smoke that the resolved version advances). No further spike needed — Facet 2's skew detection composes already-proven behavior (`cli.Version` + manifest `version`, both already threaded into `Compare`).

## Proposed direction (firmed — single task, two facets)

- **Facet 1:** make `runInit`'s codex arm refresh-on-present like the claude arm — call `ops.Install("codex", marketplaceSource, devBranch)` before running doctor when not `--check`, instead of short-circuiting to doctor on `resolved != ""`. The `--check` path keeps its no-install report. The claude arm is unchanged (already correct).
- **Facet 2:** when binary and plugin display semvers are both valid semver and the binary is strictly newer while contract-compatible, append an honest opt-in hint line to the `Compatible` message: a newer plugin is available — run `spacedock install --host {host}` to refresh. Keep the contract-based `compatible` verdict and its exit-0 status. Surface the hint in BOTH doctor output AND the front-door `gateHost` Compatible arm (which is currently silent). The hint NEVER fires when versions are equal or `cli.Version` is non-semver (`"dev"`) — defensive: no false "you must upgrade".

## Out of scope

- Changing the contract-compatibility semantics themselves (contract `1` plugin with a contract `1` binary IS compatible — that stays).
- The branch/release-model decision (HEAD-vs-tag serving, trunk direction) — tracked separately in the post-flip roadmap shaping; this task is the install/upgrade correctness + hint regardless of model.

## Acceptance criteria

(Each verified by command behavior / on-disk plugin manifest state — never a prose-grep of the skill or contract. The expected value comes from an independent source — the recorded `ops.Install` calls on the seam, the resolved on-disk manifest version, or the version inputs fed to `Compare` — not from the file under test.)

**AC-1 — `spacedock install --host codex` re-installs an already-present plugin.** End-state: running `runInit` for codex without `--check` invokes the install seam exactly as the claude arm does, regardless of whether a plugin is already resolved. Verified by: a `fakeHost`-seam test (mirroring `TestInitClaudeIssuesHostPluginCommands`) where the fake reports an already-present compatible codex manifest; after `runInit(--host codex)`, `fake.installCmds` records an `Install` call with host=`codex` and source=`spacedock-dev/spacedock`. The expectation (an install WAS issued) comes from the seam's recorded calls, an independent value the production wiring can fail to produce. This replaces the current `TestInitCodexInstallReadiness/compatible-installed` assertion, which codifies the no-op bug (it asserts NO install on a present-compatible plugin). The `--check` arm keeps `len(fake.installCmds) == 0`.

**AC-2 — codex install advances a behind plugin on disk (live).** End-state: against the real codex CLI, installing over an already-present older plugin leaves the resolved cache manifest at the newer source version. Verified by: a host-CLI smoke (skip when `codex` absent, kin to `TestCodexPluginInstallIsHostNative`) that seeds an older version, runs the production `runInit`/`execHost.Install` codex path against a newer local marketplace, and asserts the resolved manifest version equals the newer marketplace version — read from the resulting on-disk manifest, not the command's claim. (The spike above already proved the seam does this; this AC pins it as a kept regression test on the wired-up `runInit` path.)

**AC-3 — a contract-compatible-but-behind plugin surfaces an opt-in upgrade hint in doctor.** End-state: when `cli.Version` and the plugin manifest `version` are both valid semver and the binary is strictly newer while contract-compatible, the rendered `Compatible` message contains an extra line naming that a newer plugin is available and the `spacedock install --host {host}` command; the verdict stays `Compatible` and `RunDoctor` exits 0. Verified by: a `contract` unit test feeding `Compare(CONTRACT_VERSION, ">=1,<2", host, branch, pluginVersion="0.19.8", binaryVersion="0.20.0")` and asserting the message contains the hint AND the install command for the host; a paired negative case with `pluginVersion == binaryVersion` (and one with `binaryVersion="dev"`) asserting NO hint line. The trigger (semver skew) comes from the version inputs passed in, not from the message under test.

**AC-4 — the front-door gate surfaces the same hint on Compatible.** End-state: `gateHost` (and thus `spacedock claude`/`codex`) prints the upgrade hint to stderr when the resolved plugin is contract-compatible but behind, then proceeds to launch (the hint never blocks or changes the exit/launch path). Verified by: a `frontdoor`-seam test (kin to the `gate_test.go` / `frontdoor_test.go` fakes) with a fake reporting a behind-but-compatible manifest; assert the launch still happens (the fake's `Launch` is invoked) AND stderr contains the hint line. A paired equal-version case asserts the gate stays silent. Proven by invoking the gate and observing stderr + the recorded launch, never by grepping the source.

## Documentation changes (doc diff for implementation to apply)

The behavior-visible doc is `docs/install-journey.md`. Its "Keep things in sync" section currently tells the user doctor will report an out-of-date plugin — but a contract-compatible-but-behind plugin is reported `compatible`, so today the user gets no signal. Update it to match the new hint behavior:

```diff
 ## Keep things in sync

-`spacedock doctor` is the compatibility check. If it reports your installed
-plugin is out of date, refresh it:
+`spacedock doctor` is the compatibility check. When your installed plugin still
+works with the binary but a newer plugin is available, doctor and the
+`spacedock claude`/`codex` launch both print an upgrade hint. To refresh:

 ```bash
 spacedock install --host claude
 ```

 If the `spacedock` command itself is missing, install the launcher with Homebrew
 first, then run `spacedock install --host claude`.
```

No other doc names the install/launch output for this path; the front-door grammar and Codex-install prose sections are unchanged (the codex `install` command surface is unchanged — only its already-present behavior changes).

## Test plan

All tests live next to existing kin in `internal/cli` and `internal/contract`; no new fixtures beyond a versioned local marketplace builder (the spike's `buildCodexMarketplaceAtVersion` can be revived for AC-2). Estimated cost: low — four tests, three of them fast in-process seam/unit tests, one host-CLI smoke gated on `codex` presence.

- **AC-1 (fast, seam):** `fakeHost` test in `init_test.go`; replaces the `compatible-installed` no-op assertion. Asserts `fake.installCmds` records the codex install. ~20 lines.
- **AC-2 (live, gated):** host-CLI smoke in `install_behavior_codex_test.go`; seeds 0.0.1, runs the wired `runInit` codex path against a 0.0.2 marketplace, asserts the resolved manifest advances. Skips when `codex` not on PATH.
- **AC-3 (fast, unit):** `contract` test feeding `Compare` a semver-skewed-but-contract-compatible pair; asserts the hint line + host install command, plus equal-version and `dev`-binary negative cases. ~25 lines.
- **AC-4 (fast, seam):** `frontdoor` fake-ops test; asserts stderr hint + launch-still-happens on a behind-but-compatible manifest, silent on equal versions.

Risk note: AC-3's hint is rendered in `contract` where `Compare` already receives both display versions; the semver compare is a small helper in `contract` (the existing `compareVersion` lives in `cli` and must not be reached from `contract` — duplicate the tiny dotted-int compare or lift it to a shared spot, implementer's call). The hint must be additive to the existing `Compatible` message so the `compatible` token and exit-0 behavior other tests assert stay intact.

## Notes

First real-world 0.20.0 install feedback. Lands in the 0.20.x cleanup/UX band. Related: the post-flip branch-model decision (separate concern). The "binary follows tags, plugin follows branch HEAD" asymmetry is the backdrop — a tag-fresh binary against a stale branch-HEAD plugin is exactly the skew a user hits.

**Sequencing:** this task (`te`) lands the install-refresh + hint behavior independent of release model, and should land BEFORE `gp`'s marketplace-migration step — the migration changes where the newer plugin comes from, but the refresh-on-present wiring and the skew hint are correct regardless of source. Landing `te` first means the migration step inherits a working refresh path rather than the codex no-op.

## Stage Report: ideation

- DONE: Pin the desired install/refresh behavior on a stale plugin vs a newer binary; define the end-state for each host.
  Per-host end-state firmed in "Proposed direction": claude `install` unchanged (already refresh-on-present); codex `install` arm must call `ops.Install` instead of short-circuiting to doctor; front door surfaces the hint on Compatible but does NOT silently auto-refresh.
- DONE: Run the riskiest-mechanism check FIRST (how install detects/refreshes a stale plugin and where the newer plugin comes from).
  Throwaway spike (now deleted) proved codex `execHost.Install` advances an already-present plugin 0.0.1→0.0.2 on disk (resolved manifest version observed); claude path already proven by `TestUpgradeFromStaleMovesToGreen`. Root cause is `runInit` codex arm never calling the proven seam, not the seam — recorded in "Spike" section.
- DONE: Propose the doc diff for the install-output / upgrade-path change and a test plan with command-output/exit-code AC proof; note the te-before-gp sequencing.
  Unified diff against `docs/install-journey.md` "Keep things in sync" recorded; four ACs each bound to an independent source (recorded seam calls / on-disk manifest / `Compare` inputs); sequencing note added.

### Summary

Root-caused Facet 1: `spacedock install --host codex` never re-installs an already-present plugin because `runInit`'s codex arm short-circuits to doctor on `resolved != ""`, while the claude arm calls `Install` unconditionally — the `Install` seam itself is proven-good (spike confirmed codex present→newer refresh works on disk). Facet 2 is a silent-on-Compatible gap: the binary holds both display semvers locally, so an additive opt-in upgrade hint can fire in doctor and the front-door gate with no network fetch and no change to the contract verdict. Firmed four behavioral ACs (seam-recorded install call, live on-disk version advance, contract hint render with negative cases, front-door gate hint), a `docs/install-journey.md` doc diff, and the te-before-gp sequencing. Both fixes are small wiring changes, not seam rewrites.

## Stage Report: implementation

- DONE: Facet 1 — make `runInit`'s codex arm call `ops.Install("codex", marketplaceSource, devBranch)` before doctor when NOT `--check`, mirroring the claude arm, instead of short-circuiting to doctor on `resolved != ""`. The `--check` arm keeps its no-install report.
  `internal/cli/init.go` codex arm: `--check` → doctor only; present plugin → Install-then-doctor; absent → unchanged prose. Commit fc77cd10.
- DONE: REPLACE `TestInitCodexInstallReadiness/compatible-installed` (it codified the no-op bug). (AC-1)
  Rewrote the subtest to assert `fake.installCmds == [codex spacedock-dev/spacedock next]` on a present-compatible plugin; added a `--check` subtest asserting zero installCmds. `internal/cli/init_test.go`. Fails before the init.go fix, passes after.
- DONE: AC-2 = live codex smoke that the on-disk resolved manifest advances, skip when codex absent.
  `TestCodexInitRefreshAdvancesBehindPlugin` in `install_behavior_codex_test.go` seeds 0.0.1 then refreshes to 0.0.2 via the production `execHost.Install` seam, asserts `resolveCodexManifest`'s on-disk version == 0.0.2. RAN (not skipped) — codex on PATH — PASS in 0.30s.
- DONE: Facet 2 (UX hint) — additive opt-in upgrade hint when binary + plugin display semvers are BOTH valid semver AND binary strictly newer AND contract-compatible, in BOTH doctor (contract Compare) AND front-door `gateHost`. Verdict stays Compatible, exit 0, launch still happens. NEVER on equal versions or `cli.Version=="dev"`.
  `internal/contract/contract.go`: `upgradeHint`/`semverCompare`/`parseDottedInts` (contract-local, NOT `cli.compareVersion`) append the hint to the Compatible `Message` and populate a new `Result.Hint`. `gateHost` prints `res.Hint` alone on Compatible (silent on the bare OK line), then launches. Commits b3cf4ec8 + 5d527f9b.
- DONE: AC-3 contract unit + negative equal/dev cases.
  `TestCompatibleUpgradeHint`: behind-plugin (0.19.8 vs 0.20.0) hints for both hosts + preserves the OK line; equal-version, `dev`-binary, and older-binary cases assert NO hint. `internal/contract/contract_test.go`.
- DONE: AC-4 frontdoor seam — stderr hint + launch-still-happens, silent on equal.
  `TestFrontDoorUpgradeHintOnBehindPlugin` (claude+codex): with `Version` stamped 0.20.0 over the 0.12.1 fixture, asserts launch invoked + stderr carries the hint + host install command; equal-version (0.12.1) asserts silent. New `withVersion` test helper. `internal/cli/frontdoor_test.go`.
- DONE: `go test ./internal/cli/ ./internal/contract/` green WHOLE-package; ZERO `.md` edits; hint is ADDITIVE.
  `322 passed in 2 packages` (whole-package); full repo `1261 passed in 16 packages`; `go vet` clean. `git diff --name-only origin/main...HEAD` → 7 Go files, ZERO `*.md`.
- SKIPPED: the `docs/install-journey.md` doc diff in the spec's "Documentation changes" section.
  Delegated to yw per the dispatch anti-collision rule (yw is deleting install-journey.md and moving the doctor-refresh line to `reference/command-reference.md ## Setup`); my branch makes ZERO `.md` edits.

### Summary

Both facets landed as small wiring changes on `origin/main` HEAD 41bd47ae, three commits (fc77cd10, b3cf4ec8, 5d527f9b). Facet 1 wires the codex arm to refresh-on-present like the claude arm and replaces the no-op-codifying test with a seam-recorded install assertion plus a live on-disk-advance smoke (ran against the real codex CLI). Facet 2 adds a contract-local semver-skew gate that appends an opt-in upgrade hint to the Compatible message in doctor and surfaces it (via a new `Result.Hint`) at the previously-silent front-door gate, additive and never firing on equal/dev. Whole-package and full-repo tests green; zero `.md` edits (doc-diff delegated to yw). Note per dispatch: a detached adversarial audit of the front-door changes is owed at validation before merge.

## Stage Report: validation

- DONE: Reproduce `go test ./internal/cli/ ./internal/contract/` green WHOLE-package (implementer: 322).
  `322 passed in 2 packages`; full-repo `go test ./...` = `1261 passed in 16 packages`; `go vet` clean. Matches implementer.
- DONE: AC-1 (init_test — codex arm records `Install [codex spacedock-dev/spacedock next]` on present-compatible; `--check` records zero installCmds; REPLACED `compatible-installed` no longer codifies the no-op).
  `TestInitCodexInstallReadiness/{compatible-installed,compatible-installed-check}` PASS; rewritten subtest asserts the recorded install argv, the `--check` subtest asserts `len(installCmds)==0`.
- DONE: AC-2 (live codex smoke `TestCodexInitRefreshAdvancesBehindPlugin` — on-disk resolved manifest advances 0.0.1→0.0.2).
  codex on PATH (`/opt/homebrew/bin/codex`) — RAN (not skipped), PASS in 0.29s; seeds + refreshes via production `execHost.Install`, asserts `resolveCodexManifest` on-disk version, not the command's claim.
- DONE: AC-3 (contract hint on behind 0.19.8-vs-0.20.0 for both hosts + OK line preserved; NO hint on equal/dev/older).
  `TestCompatibleUpgradeHint` all four subtests PASS (behind-plugin-hints, equal/dev/older-binary no-hint); behind case asserts the OK line is preserved alongside the hint.
- DONE: AC-4 (frontdoor — stderr hint + host install command + launch-still-happens on behind plugin; SILENT on equal).
  `TestFrontDoorUpgradeHintOnBehindPlugin` 4 subtests PASS (claude+codex behind hints+launches, equal silent); asserts recorded launch + stderr, never a source grep.
- DONE: ADVERSARIALLY confirm the guards discriminate (init.go revert; additive hint; equal/dev guard; contract-local compare).
  Reverted init.go codex arm → AC-1 reds with `install seam = [], want [codex spacedock-dev/spacedock next]`. Forced `upgradeHint` always-on → equal/dev/older contract cases AND both front-door equal-version cases red (launch banner confirms launch still happens). Additive: 322 whole-package + 1261 full-repo green with the hint present, no `compatible`/exit-0 test broke. Layering: `internal/contract` does not import `internal/cli`; `cli.compareVersion` (in host_exec.go) is never referenced from contract — `semverCompare`/`parseDottedInts` are contract-local. All edits reverted; tree clean against HEAD.
- DONE: Confirm zero `.md` edits (7 Go files); recommend PASSED/REJECTED.
  `git diff --name-only 41bd47ae...HEAD` → exactly 7 Go files, 0 `.md` (0 non-`.go`). Recommendation: PASSED.

### Summary

PASSED. All four ACs verified by running the behavior and observing test results, exit codes, and on-disk manifest state — not prose-grep. The whole-package (322) and full-repo (1261) suites reproduce the implementer's counts; `go vet` clean. Four adversarial discrimination checks confirm the guards are non-vacuous: reverting the init.go fix reds AC-1's recorded-install assertion; forcing the hint always-on reds every equal/dev/older negative case in both contract and the front door; the hint is additive (verdict stays Compatible, exit 0, launch proceeds, no pre-existing test broke); and the semver compare is contract-local (no `internal/cli` import, `cli.compareVersion` unreferenced from contract). te is the front-door = high-stakes: a detached read-only adversarial audit on a throwaway checkout is owed before merge — that is the FO's step, dispatched after this PASSED.
