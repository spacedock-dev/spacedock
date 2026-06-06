---
id: rbjkna5jem4vgj3vtv072gzq
title: spacedock claude auto-installs the plugin when absent (--no-install opt-out)
status: ideation
source: "captain (2026-06-05) — friction F8: `spacedock claude` refuses to launch with no installed plugin, forcing the user to pass --skip-contract-check (a launch blocker). Captain direction: 'we need something simpler [than task 44]. maybe just install the plugin unless --no-install is specified in spacedock claude.' Interim relief; task 44 (bundle-into-binary) is the eventual structural fix but is deferred."
score: "0.36"
started: 2026-06-06T03:06:22Z
completed:
verdict:
worktree:
issue:
---

`spacedock claude` runs a fail-fast contract gate before launch (`internal/cli/frontdoor.go:167-171`, `gateHost(ops, "claude", stderr)`) that denies when no plugin is installed AND no `--plugin-dir` is passed. With no installed plugin, `ResolveManifest` returns empty and the gate prints "no installed claude plugin found. Run `spacedock install --host claude` (or --skip-contract-check to bootstrap)" and exits 1. No-plugin is the most common first-run state, yet the only advertised escape is a niche bootstrap flag or a separate install round-trip — so a fresh user is blocked from the one command they tried.

## Direction (for ideation)

- On `spacedock claude` when no plugin is resolvable (the NoPluginFound / empty-manifest case ONLY — NOT a version mismatch), AUTO-RUN the install (`spacedock install --host claude`, which already exists programmatically at `internal/cli/host_exec.go` installArgvSequence) and then proceed to launch — so the single command the user typed yields a working FO session.
- Gate it with a `--no-install` opt-out for users who want the old refuse-and-instruct behavior.
- This is verdict-scoped: a real version mismatch (TooOldBinary / TooOldPlugin / MalformedRange) must STILL fail fast — do not auto-install over an incompatibility. Have gateHost return the verdict (not just a bool) so runClaude can distinguish no-plugin (auto-install) from incompatible (hard-fail).
- Distinct from option A (silent launch with NO plugin -> broken session, rejected): here we INSTALL the plugin so the launch actually works.

## Out of scope

Bundling the plugin into the binary + injecting --plugin-dir (that is task 44 `bundle-asset-distribution`, deferred — it eventually makes the no-plugin case impossible and supersedes this auto-install round-trip). The binary-ABSENT journey (that is qa `spacedock-binary-missing-install-journey` — different root cause: here the binary is present, the plugin is missing). Codex (no --plugin-dir equivalent; this task is Claude's no-plugin launch).

## Proposed approach

The single seam is `internal/cli/frontdoor.go`. Two changes, plus a flag:

1. **`gateHost` returns the verdict, not a `bool`.** Change the signature to `gateHost(ops, host, stderr) contract.Verdict`. Map each existing branch to a verdict: `Compatible` on the compatible path; `NoPluginFound` for BOTH no-plugin states (the empty-`manifestPath` branch AND the resolved-but-missing-file `ManifestVerdict==NoPluginFound` branch); the resolve-error branch maps to a non-`NoPluginFound`, non-`Compatible` verdict (proposed: `MalformedRange`) so it stays a hard fail (a broken host CLI is NOT a missing plugin — auto-installing over it would just fail again); the mismatch path returns `res.Verdict` verbatim (`TooOldBinary`/`TooOldPlugin`/`MalformedRange`/`PluginPredatesContract`). The stderr messages are unchanged — only the return type changes, so the existing remedy text and `TestGateRemedyNamesLiveInstallCommand` keep passing (that test asserts on stderr, and its `if ok :=` call site updates to `if v := …; v == contract.Compatible`).

2. **`runClaude` branches on the verdict.** Replace `if !gateHost(...) { return 1 }` with a switch: `Compatible` → proceed; `NoPluginFound` → if `fd.noInstall` print the instruct message and `return 1`, else `ops.Install("claude", marketplaceSource, devBranch)` (the same call `runInit` makes) and on success proceed to launch (on install error, `return 1`); `default` (any mismatch / resolve-error) → `return 1` (the gate already printed the message). `runCodex` keeps the old all-or-nothing semantics: `if gateHost(...) != contract.Compatible { return 1 }` — codex is out of scope and `Install` rejects non-claude hosts anyway.

3. **`--no-install` flag.** Add a `noInstall *bool` to `frontDoorFlags`, bind it in `bindFrontDoorFlags` (`fs.Bool("no-install", false, …)`), read it into a new `frontDoorArgs.noInstall` in `parseFrontDoorArgs`. It rides the existing pflag pre-`--` grammar like `--skip-contract-check`, so no grammar work is needed.

**Open design decisions (resolve during implementation, surface at gate):**
- **Re-gate after install vs. proceed directly.** The spike proceeds straight to launch after a successful `Install` (matching the AC wording "then proceeds to launch"). Alternative: re-run `gateHost` after install so a silently-broken install (Install returns nil but the plugin still does not resolve) does not launch into a broken FO session. The `Install` fail-fast steps already surface real failures via a non-nil error, so proceed-directly is defensible; re-gate is one extra `ResolveManifest`. Recommend proceed-directly for simplicity, revisit if a broken-install case is observed.
- **Resolved-but-missing-manifest treatment.** This is the `NoPluginFound` verdict from a phantom `installPath`, so under the design it auto-installs too. Defensible (the plugin genuinely is not on disk), but it changes `TestClaudeFrontDoorNonEmptyMissingManifestFailsFast` from a hard fail to auto-install-by-default. Recorded as a deliberate behavior change, not an accident.

## Spike (mechanism validation)

Ran the riskiest unknown first — confirm `gateHost` can return the verdict and `runClaude` can branch `NoPluginFound`-vs-mismatch using the existing `frontdoor_test.go` `fakeHost` stubs. Made the scratch edits above (verdict return + switch + `--no-install`) and a throwaway `spike_autoinstall_test.go` driving `runClaude` against `fakeHost`, then reverted both (spike is throwaway; implementation rebuilds it test-first).

Result — **mechanism confirmed, all three spike tests passed:**
- `fakeHost{manifest: ""}`, no flags → `gateHost` returns `NoPluginFound`; `runClaude` called `ops.Install` (observed via `fake.installCmds` non-empty) AND reached `ops.Launch` (`fake.launchedArg` non-nil); exit 0. (AC-1a)
- `fakeHost{manifest: ""}`, `--no-install` → no `Install` call, no `Launch`, exit non-zero. (AC-1b)
- `fakeHost{manifest: tooOldBinaryManifest()}`, no flags → mismatch hit the `default` branch: no `Install`, no `Launch`, exit non-zero. (AC-2)

**Finding the spike surfaced (drives the test plan):** the existing `TestClaudeFrontDoorUnresolvableManifestFailsFast` and `TestClaudeFrontDoorNonEmptyMissingManifestFailsFast` FAILED under the scratch behavior (both now exit 0). They encode the OLD contract where no-plugin is a fail-fast. They are not stale — they must be REWRITTEN to the new default-auto-installs contract (with `--no-install` as the new fail-fast path). The `fakeHost` harness already records install-invocation, launch-reached, and exit code — every AC observable is an already-present stub field, so no harness work is needed.

## Acceptance criteria

**AC-1 — `spacedock claude` with no installed plugin and no flags auto-installs the plugin then launches a working session; `--no-install` preserves the refuse-and-instruct behavior.**
Verified by: a frontdoor test in `internal/cli/frontdoor_test.go` driving `runClaude` against `fakeHost{manifest: ""}`. With no flags: `fake.installCmds` is non-empty (install invoked) AND `fake.launchedArg != nil` (launch reached) AND exit 0. With `--no-install`: `fake.installCmds` is empty (no install) AND `fake.launchedArg == nil` (no launch) AND exit non-zero with the instruct message on stderr. Independent source: install-invocation, launch-reached, and exit code are observed behaviors recorded by the stub — none is a string match over any file the model reads. Spike-confirmed both arms pass against the existing harness.

**AC-2 — a real version mismatch still fails fast even without --no-install (auto-install does not paper over incompatibility).**
Verified by: a frontdoor test driving `runClaude` against `fakeHost{manifest: tooOldBinaryManifest(t)}` with no flags, asserting exit non-zero AND `fake.installCmds` empty (no auto-install) AND `fake.launchedArg == nil` (no launch). The verdict reaches `runClaude`'s `default` (mismatch) branch, not the `NoPluginFound` branch. Spike-confirmed.

## Test plan

Go frontdoor unit tests in `internal/cli/frontdoor_test.go` (the existing `fakeHost` harness already stubs `ResolveManifest`/`Launch`/`Install` and records `installCmds`/`launchedArg` — every AC observable is an already-present field, so no harness work). Cost: low, pure unit, no network/exec, sub-second.

New/changed tests:
- New: AC-1a (no plugin + no flags → install + launch), AC-1b (no plugin + `--no-install` → no install + no launch + instruct), AC-2 (mismatch → hard fail, no install, no launch).
- **Rewrite (the spike's finding):** `TestClaudeFrontDoorUnresolvableManifestFailsFast` and `TestClaudeFrontDoorNonEmptyMissingManifestFailsFast` currently assert no-plugin/missing-manifest fail fast — the OLD contract. Under default-auto-installs both now exit 0. Migrate them to assert the NEW default (auto-install) and add the `--no-install` fail-fast arm. Do NOT delete; the no-plugin journey still needs a fail-fast test, just gated on `--no-install`.
- Update the `gateHost` call site in `TestGateRemedyNamesLiveInstallCommand` from `if ok :=` to `if v := …; v == contract.Compatible` (the remedy-string assertions on stderr are unchanged).
- `runCodex`: confirm `TestCodexFrontDoorFailFastOnMismatch` and `TestCodexFrontDoorLaunchesOnCompatible` still pass unchanged (codex semantics preserved — no auto-install).

TDD order: write AC-1a/AC-1b/AC-2 (red), make `gateHost` return the verdict + `runClaude` switch + `--no-install` flag (green), then migrate the two FailFast tests. Front-door launcher is a high-stakes surface -> detached adversarial audit before merge.

## Stage Report: ideation

- DONE: Design the verdict-scoped auto-install: gateHost RETURNS the verdict (not a bool) so runClaude distinguishes no-plugin (auto-install via host_exec.go installArgvSequence, the same `ops.Install("claude", marketplaceSource, devBranch)` runInit makes) from incompatible (hard-fail); gated by --no-install opt-out preserving refuse-and-instruct.
  Proposed approach + open decisions written into the entity body; single seam is `internal/cli/frontdoor.go` (gateHost + runClaude + frontDoorFlags/parseFrontDoorArgs for `--no-install`); runCodex preserved via `!= Compatible`.
- DONE: Acceptance criteria are proven by RUNNING the front door, not string-match: AC-1 (no-plugin+no-flags -> install invoked + launch reached; --no-install -> no install, instruct-exit) and AC-2 (real version mismatch hard-fails without --no-install).
  ACs rewritten to assert observed `fake.installCmds`/`fake.launchedArg`/exit code via the existing TestClaudeFrontDoor* `fakeHost` harness — no string-in-file proof; test plan names the new tests, the two FailFast tests to migrate, and the gateHost call-site update.
- DONE: Spike-first the riskiest unknown (gateHost returns verdict + runClaude branches NoPluginFound-vs-mismatch on existing frontdoor_test stubs) and record the result; high-stakes front door -> detached adversarial audit before merge.
  Ran scratch edits + throwaway `spike_autoinstall_test.go`: all 3 spike tests passed (AC-1a/AC-1b/AC-2 mechanism confirmed); spike reverted (`git status` clean, post-revert frontdoor tests 12/12 green). Spike surfaced the FailFast-tests-must-be-rewritten finding. Adversarial-audit-before-merge recorded in the test plan.

### Summary

Fleshed out the no-plugin auto-install design on `internal/cli/frontdoor.go`: `gateHost` returns `contract.Verdict`, `runClaude` switches `NoPluginFound`→(`--no-install`? instruct-exit : install-then-launch) and any mismatch→hard-fail, plus a `--no-install` pflag. Spiked the riskiest unknown first against the real `fakeHost` stubs — all three ACs proven (install-invoked/launch-reached/exit observed, not string-matched) — then reverted the throwaway. The spike surfaced one concrete finding now on the record: the two existing `*FailFastOnMismatch`-style no-plugin tests encode the old fail-fast contract and must be migrated to default-auto-installs with `--no-install` as the new fail-fast arm. Two open decisions (re-gate-after-install vs proceed; phantom-installPath treatment) are recorded with recommendations for the implementer to settle at the gate.
