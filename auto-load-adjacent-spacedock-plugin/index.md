---
title: Auto-load the adjacent Spacedock plugin
status: implementation
source: captain discussion 2026-07-17
started: 2026-07-17T13:00:22Z
completed:
verdict:
score: 0.9
worktree: .worktrees/spacedock-ensign-auto-load-adjacent-spacedock-plugin
issue:
id: zbx5d7c4xgre4aq96d6q3xdy
---

A locally built Spacedock launcher should load the matching plugin checkout beside its own executable without requiring `--plugin-dir`. This separates launcher identity from additional host plugins and prevents unrelated directories from replacing or bypassing Spacedock accidentally.

## Problem

The front doors currently overload `--plugin-dir` as both a local Spacedock development override and, on Claude, a native additional-plugin flag. Any Claude occurrence bypasses the installed Spacedock compatibility gate. Codex has no native session flag, so the pre-fence form is converted into a persistent local marketplace that hard-codes the supplied directory as `spacedock` and removes existing Spacedock channels before validating the replacement. Post-fence Codex forwarding is documented and tested even though Codex rejects it.

The common development case already has a stronger identity source: the resolved launcher recorded in `SPACEDOCK_BIN`. When that executable sits at the root of a valid Spacedock plugin checkout, the launcher can select the adjacent plugin automatically.

## Proposed approach

Resolve the launcher through the existing executable-resolution path and inspect its parent directory. Treat that directory as the local Spacedock plugin root only when the host-specific manifest exists, declares `name: spacedock`, and its configured `skills/` directory exists. A bare adjacent `skills/` directory is insufficient.

For Claude, inject the validated adjacent root as the session-local Spacedock plugin. Keep pre-fence `--plugin-dir` as the compatibility-preserving explicit Spacedock override; forwarded post-fence plugin directories remain native Claude additions and do not themselves bypass the Spacedock gate.

For Codex, feed the validated adjacent root through the existing local-marketplace adapter. Validate the checkout completely before removing or installing any channel. Keep additional Codex plugins under Codex's persistent marketplace/plugin commands; do not build a general plugin manager into Spacedock. Reject or document unsupported post-fence `--plugin-dir` behavior accurately instead of promising forwarding that Codex cannot accept.

When the executable has no valid adjacent Spacedock plugin, preserve the installed-plugin resolution, compatibility gate, auto-install, and release/Homebrew behavior.

The simplest alternative is to continue requiring `--plugin-dir` for local development. It cannot deliver the value because it retains the replacement-versus-additional ambiguity that caused this failure. Automatic discovery reuses the already-established launcher identity and adds no new protocol or persistent registry.

## Out of scope

- Installing arbitrary additional Codex plugins through Spacedock.
- Changing Codex's persistent plugin model.
- Packaging `spacedock-subspace`; it needs its own Claude and Codex manifests separately.
- Changing workflow, status, dispatch, PR, or mod semantics.

## Acceptance criteria

**AC-1 (VALUE) - A locally built launcher in a valid Spacedock checkout starts Claude and Codex with that adjacent checkout selected without an operator-supplied `--plugin-dir`.**
Verified by: front-door behavior tests that place the resolved executable beside valid host manifests and distinct skill markers, then observe Claude argv and Codex's installed provider/source resolving to that exact checkout.

**AC-2 - A launcher without a valid adjacent Spacedock manifest preserves installed-plugin gating and installation behavior.**
Verified by: table tests for missing manifest, wrong manifest name, missing skills directory, and non-adjacent release-style binary layouts; each must exercise the existing resolver/gate seam and must not select a local checkout.

**AC-3 - Claude additional plugin directories do not impersonate the Spacedock override or suppress its compatibility gate.**
Verified by: a behavior test with a compatible installed Spacedock plugin plus a distinct valid post-fence plugin directory; the gate is observed, the additional directory reaches Claude unchanged, and the selected Spacedock provider remains the installed or adjacent one.

**AC-4 - Codex rejects an invalid adjacent or explicit Spacedock checkout before any channel cleanup or plugin installation.**
Verified by: an invalid-checkout test that asserts a non-zero result, an actionable manifest diagnostic, zero install-seam calls, and unchanged isolated Codex configuration state.

**AC-5 - Existing explicit local Spacedock development launches remain compatible.**
Verified by: the existing pre-fence `--plugin-dir` Claude and Codex tests, updated only where provider selection now shares the adjacent-root resolver, plus focused regression tests for explicit override precedence.

**AC-6 - Help and development documentation describe host-specific behavior truthfully.**
Verified by: command help snapshots and rendered documentation review paired with executable behavior tests: Claude supports forwarded additional session plugins; Codex requires persistent plugin installation and does not accept a forwarded `--plugin-dir`.

## Test plan

Add focused unit tests in `internal/cli` before implementation for adjacent-root qualification, explicit-override precedence, Claude gate behavior, and Codex preflight-before-mutation. Use isolated temporary launcher roots and `CODEX_HOME` directories. Exercise current Claude and Codex CLI schemas with fixture-backed command behavior; use live host commands only where the claim depends on provider resolution rather than argv construction.

Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`. Because the change touches the high-stakes front-door launcher and both host lanes, validation must also run the required Claude and Codex live lanes and a detached adversarial audit that removes or corrupts an adjacent manifest and confirms the tests catch the fallback or fail-closed boundary.

## Stage Report: implementation

- DONE: Local binaries automatically select an adjacent valid Spacedock checkout for both Claude and Codex.
  Commit `a69032c8` adds host-manifest-qualified discovery; isolated Claude 2.1.212 and Codex 0.144.4 smokes launched the adjacent checkout, and Codex listed it as the sole enabled provider.
- DONE: Missing or invalid adjacent checkouts preserve installed fallback or fail before any Codex mutation.
  Fixture tests cover missing/wrong-name/missing-skills/release layouts for both hosts; invalid explicit Codex input left isolated `CODEX_HOME` unchanged with zero install calls.
- DONE: Explicit overrides and host-specific additional-plugin behavior stay compatible, truthful, tested, and documented.
  Pre-fence precedence tests pass for both hosts; Claude post-fence additions retain installed gating, Codex rejects unsupported forwarding, strict MkDocs rendering succeeds, and help snapshots pin both behaviors.

### Summary

Implemented the smallest boundary around the resolved launcher directory, reusing Claude's session flag and Codex's existing local-marketplace adapter only after complete validation. `go test ./...`, `go test ./... -race`, real isolated Codex provider tests, host CLI smokes, strict documentation rendering, and a detached adversarial audit all passed; the audit found no blocking defects.

## Stage Report: validation

- DONE: Reproduce the live 0.25.0+dev adjacent launcher versus installed Claude 0.26.0-pre1 mismatch and prove both Claude and Codex select the adjacent checkout without an explicit --plugin-dir.
  Claude 2.1.212 debug output loaded `spacedock:first-officer` and skills from the branch worktree; isolated Codex 0.144.4 returned `ADJACENT_CODEX_OK`, and its installed provider symlink resolved to the same worktree.
- DONE: Reproduce every acceptance criterion with focused behavior evidence, then run gofmt, go test ./..., and go test ./... -race; reject false-green or self-referential test evidence.
  AC-1 through AC-6 focused tests passed; both complete Go suites passed, changed Go files are gofmt-clean, and repository-pinned MkDocs 1.6.1 built strictly.
- DONE: Run the required detached adversarial audit on a throwaway checkout by removing or corrupting adjacent manifests and confirm the deliverable tests catch the fallback and fail-closed boundaries.
  Corrupting both manifest names failed both repository manifest guards; automatic Claude fell back to the exact 0.25.0+dev/0.26.0-pre1 mismatch, Codex fell back to installed resolution, and explicit invalid Codex exited before writing isolated config.
- DONE: AC-1 adjacent local selection for Claude and Codex.
  Live host launches used no explicit `--plugin-dir`; Claude's debug log named the exact adjacent skills root, while Codex's provider symlink and successful agent output named and exercised the exact root.
- DONE: AC-2 invalid adjacent layouts preserve installed resolution and healing.
  Eight host/layout table cases passed for missing manifest, wrong name, missing skills, and release-style `bin/`; both host auto-install fallback cases re-resolved, gated, and launched.
- DONE: AC-3 Claude additional plugins remain additions and cannot suppress the installed gate.
  The behavior test observed installed-manifest resolution, preserved the additional directory in exact Claude argv, and reached the launch seam only after a compatible verdict.
- DONE: AC-4 invalid Codex input fails before mutation.
  The focused test asserted non-zero exit, manifest-name diagnostic, zero install calls, and an empty isolated `CODEX_HOME`; the detached real-CLI probe reproduced all four boundaries.
- DONE: AC-5 explicit development override compatibility and precedence.
  Existing explicit Claude/Codex tests and both adjacent-versus-explicit precedence tests passed.
- DONE: AC-6 truthful host-specific help and development documentation.
  Help behavior tests passed, Codex post-fence rejection was executable rather than prose-only, and `mkdocs build --strict` completed with the pinned requirements.
- DONE: Semantic adversarial pass and false-green audit.
  Disabling discovery failed both host selection tests; disabling manifest-name validation failed both host fallback lanes and the Codex no-mutation test, proving the new evidence is mutation-sensitive rather than self-confirming.
- DONE: Validation recommendation: PASSED.
  No material outcome or evidence defects remain; no deferred risks were identified for the promised workflows.

### Summary

Validation independently reproduced the release-autobump mismatch and proved that the adjacent checkout now wins for both live hosts without an explicit override. All acceptance criteria, complete test gates, strict docs build, and detached adversarial mutations passed; recommendation is PASSED with no blocking findings.

### Feedback Cycles

- Cycle 1 — evidence defect in AC-3: the Claude additional-plugin test proves the installed compatibility gate runs and the extra directory reaches argv, but it does not observe which `spacedock:first-officer` provider Claude selects. Because the additional fixture is itself named `spacedock`, the test can pass while that directory impersonates the selected provider. Route to implementation: add provider-identity observation plus an adversarial mutation that fails if the additional directory wins while the gate and argv assertions still pass; then rerun focused, full, and race tests before fresh validation.

- Cycle 2 — RoboRev branch-final rejection: `internal/cli/frontdoor.go` gates the installed Claude manifest whenever the development plugin appears after `--`, even when that post-fence manifest is itself named `spacedock` and becomes Claude's selected provider. Preserve differently named post-fence additions, but inspect post-fence manifests and retain the supported install-free Spacedock override. The new provider-identity test must also leave the baseline suite deterministic: put live-Claude execution behind an explicit integration opt-in or build tag, check the supported host version, and bound external commands with a timeout while retaining fixture-backed baseline coverage.

## Stage Report: implementation (cycle 2)

- DONE: Add a Claude behavior test that observes the selected spacedock:first-officer provider when a distinct additional plugin directory is forwarded.
  `TestClaudeAdditionalPluginKeepsInstalledSpacedockProvider` uses Claude's pre-API `agent_listing_delta` to observe the installed provider's unique description while `additional-tools` remains a distinct forwarded plugin.
- DONE: Add an adversarial mutation proving the test fails if the additional directory wins while compatibility-gate and argv assertions still pass.
  The same test mutates only the additional manifest name to `spacedock`, observes the provider description flip to `MUTANT_ADDITIONAL_PROVIDER_IDENTITY`, and proves installed-manifest resolution plus normalized argv remain unchanged.
- DONE: Run focused tests, gofmt -w ./cmd ./internal, go test ./..., and go test ./... -race; append a replacement implementation report with exact evidence.
  Focused real-Claude provider evidence, the complete suite, and the complete race suite passed; commit `26200066` contains the cycle-2 test only.

### Summary

Closed only the AC-3 evidence gap using Claude's own durable agent-listing event before the intentionally invalid isolated API key exits. The positive and adversarial arms now distinguish installed-provider selection from additional-plugin impersonation without changing the accepted adjacent-discovery implementation.

## Stage Report: validation (cycle 2)

- DONE: Independently reproduce the positive and mutant real-Claude provider-identity observations while confirming compatibility-gate and argv evidence remain unchanged.
  Claude 2.1.212 reported `spacedock:first-officer: INSTALLED_PROVIDER_IDENTITY` for `additional-tools` and `spacedock:first-officer: MUTANT_ADDITIONAL_PROVIDER_IDENTITY` when only the additional manifest name changed; both arms retained installed resolution and byte-equivalent normalized argv.
- DONE: Re-check AC-3 and the full AC-1 through AC-6 set, rejecting tautological evidence and recording any material or deferred finding.
  AC-1 through AC-6 focused behavior tests passed; the new AC-3 evidence observes Claude's final agent registry rather than fixture prose, and no material or deferred finding remains.
- DONE: Run focused tests, gofmt -w ./cmd ./internal, go test ./..., and go test ./... -race; append a cycle-2 PASSED or REJECTED validation report.
  The real-Claude test passed three consecutive isolated runs plus one instrumented observation run; focused AC tests, the complete suite, and the complete race suite passed, and changed Go files are gofmt-clean.
- DONE: Reproduce the cycle-2 evidence boundary adversarially on a detached checkout.
  Naming the positive additional fixture `spacedock` failed on the host-reported mutant marker; restoring the old post-fence gate bypass failed both the original gate assertion and the new real-Claude resolver assertion.
- DONE: AC-3 provider identity is independent of the gate and argv oracles.
  Candidate descriptions differ before launch, `agent_listing_delta` supplies the selected qualified-agent description after Claude resolves both plugins, and fresh config roots prevent stale-session evidence.
- DONE: AC-1, AC-2, AC-4, AC-5, and AC-6 remain satisfied.
  Cycle 2 changes only `internal/cli/claude_provider_identity_test.go`; all prior adjacent selection, fallback, no-mutation, precedence, help, and documentation evidence remained green.
- DONE: Validation recommendation: PASSED.
  The cycle-1 AC-3 evidence defect is closed; there are no material findings, deferred risks, or polish findings in cycle 2.

### Summary

Fresh validation observed the installed and impersonating Claude providers through Claude's own session event and proved the oracle turns red under both provider-identity and gate-bypass mutations. The full normal and race suites pass, so cycle 2 recommends PASSED with no remaining findings.
