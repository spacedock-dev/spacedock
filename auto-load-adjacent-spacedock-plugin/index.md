---
title: Auto-load the adjacent Spacedock plugin
status: validation
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
