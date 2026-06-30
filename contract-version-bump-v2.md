---
title: Bump skill↔binary contract version to 2 — 0.23.0 skills cleanly reject a contract-1 (v0.22.0) binary
status: validation
sprint: 0230-stable-finalization
score: 0.8
id: 7h3qb1atbvp8x2pd597kw802
worktree: .worktrees/spacedock-ensign-contract-version-bump-v2
started: 2026-06-22T13:41:50Z
---

The 0.23.0 skills hard-require binary verbs absent from a contract-1 (v0.22.0) binary — `spacedock state ready/commit/sweep` and `spacedock merge guard` — but the binary contract version stayed `1` and both plugin manifests declare `requires-contract: >=1,<2`. A user who upgrades skills to 0.23.0 without upgrading the binary therefore PASSES the boot contract-version gate (1 ∈ [1,2)) and then breaks cryptically (`unknown subcommand (want: init or new)` / `unknown command: merge`) instead of getting the gate's clean "binary too old — upgrade" abort. Confirmed against the published v0.22.0 binary.

## Problem

The skill↔binary contract version is the designed mechanism to reject an incompatible binary at boot (first-officer-shared-core Startup step 1: the binary's `contract <N>` must satisfy the skill's range; below the lower bound = binary too old). But v0.23.0 added hard skill dependencies on new binary verbs (state commit/ready/sweep #399, merge guard #400/#415) WITHOUT bumping the contract version from 1. A contract-1 binary (v0.22.0, and pre.1–pre.4) passes the gate and then fails on the missing verbs. Stable is the first broad-audience exposure of these skills, so the skew becomes real for existing v0.22.0 users who `plugin update` without `brew upgrade`.

## Proposed approach (captain-directed)

Bump the contract version to 2 so a contract-1 binary is cleanly rejected by the 0.23.0 skills, and (symmetrically) a contract-2 binary is rejected by old `>=1,<2` skills — both skew directions produce the designed abort.

1. Bump the binary's contract-version constant `1` → `2` (internal/cli; the value `--version` reports as `contract <N>`).
2. Update BOTH plugin manifests' `requires-contract`: `>=1,<2` → `>=2,<3` (.claude-plugin/plugin.json, .codex-plugin/plugin.json).
3. Update the contract-range text to `>=2,<3` everywhere it appears in the FO/ensign contract (first-officer-shared-core Startup step 1; any runtime adapter stating the range). The string is the same length, so this is byte-neutral and the value gate is unaffected.
4. Update any test that pins `contract 1` / the `>=1,<2` range (version tests, contractlint, the migration/range-bracketing check) to the new values; the plugin-range-brackets-binary-contract check must hold (2 ∈ [2,3)).

Do NOT change any verb behavior — this is purely the compatibility-version bump.

## Acceptance criteria

- **AC-1 (VALUE)** — a contract-1 binary is REJECTED by the 0.23.0 skills with the clean out-of-range / "binary too old" abort, where today it passes. Verified by: exercise the version-gate predicate (or the boot check) against a contract-1 input vs the v0.23.0 (contract-2) binary — contract 1 → abort/out-of-range, contract 2 → pass. Independent baseline that moves the wrong way: the published v0.22.0 binary reports `contract 1` and today PASSES `>=1,<2`; after the bump it must FAIL `>=2,<3`.
- **AC-2** — the v0.23.0 binary reports `contract 2`; both plugin manifests declare `requires-contract: >=2,<3`; the plugin-range-brackets-binary-contract test passes with the new pair (2 ∈ [2,3)) and would fail if either side were left at `1` / `>=1,<2`.
- **AC-3** — the contract-range text in the FO/ensign contract reads `>=2,<3` consistently with the plugin manifests (no stale `>=1,<2` remains in the contract docs); every value-gate-tracked FO/ensign reference file stays at or under its v0.22.0 byte baseline (the edit is byte-neutral).

## Test plan

- Go test: the binary contract constant == 2; the version-gate predicate rejects contract 1 and accepts contract 2; the plugin-range-bracketing test with the new pair (divergeable — fails if binary=1 or range=`>=1,<2`).
- A fixture or the real published v0.22.0 binary (`contract 1`) confirms AC-1's rejection.
- `go test ./...` green; re-measure the value-gate byte table.
- High-stakes contract/release surface: detached adversarial audit at validation (throwaway checkout).

## Stage Report: implementation

- DONE: Bump the contract version to 2 consistently (constant, both manifests, contract-range TEXT, tests).
  contract.go `CONTRACT_VERSION 1->2` (only non-test source change); `.claude-plugin`/`.codex-plugin` `requires-contract >=2,<3`; first-officer-shared-core Startup step 1 range text `>=2,<3` (byte-neutral, stays 28583); doctor/bracketing fixtures retargeted (compatible.json + testdata/plugin.json `>=2,<3`, too-old-binary.json `>=3,<4`); TestCompatibleUpgradeHint + gate + install/decoupling/codex fixtures tracked. Commit c113a368.
- DONE: AC-1 RED-first: prove a contract-1 binary is now REJECTED where it passed before.
  New `skills/integration/contract_skew_test.go` failed RED (contract-1 verdict=compatible) pre-bump, GREEN post-bump. Confirmed with the REAL published v0.22.0 binary (`spacedock_0.22.0_darwin_arm64`, reports `contract 1`): `doctor` against the bumped `>=2,<3` manifest -> "Upgrade the binary to continue." exit 1; against the pre-bump `>=1,<2` -> "OK" exit 0; the post-bump contract-2 binary -> "OK" exit 0. Bracketing test holds (2 ∈ [2,3)) and goes RED if either side reverts to 1 / `>=1,<2` (both divergence directions exercised).
- DONE: Self-check before handoff: `go test ./...` green; NO verb-behavior change; re-measure the value-gate byte table.
  `go test ./...` exit 0, 15 packages ok. Non-test source diff is the single `CONTRACT_VERSION` line. shared-core 28583 (HEAD=working, ≤ v0.22.0 baseline 28586); no other FO/ensign reference file touched, all stay at their HEAD bytes.

### Summary

Bumped the skill↔binary contract version 1->2 so 0.23.0 skills cleanly reject a contract-1 (v0.22.0) binary at the boot gate instead of breaking on missing verbs. Each `>=1,<2` occurrence was triaged: live-contract declarations and `CONTRACT_VERSION`-paired fixtures moved to `>=2,<3` (or `>=3,<4` for the too-old-binary fixture); abstract Compare-math literals and the StampVersion field-preservation fixture were left unchanged. The range-text edit is byte-neutral. AC-1 proven end-to-end against the real published v0.22.0 binary, and the bracketing guard is divergeable in both directions.

### Value-gate byte table (FO/ensign reference files)

| file | v0.22.0 baseline | this branch | edit |
|------|------------------|-------------|------|
| first-officer-shared-core.md | 28586 | 28583 | byte-neutral range swap, unchanged from HEAD |
| (all other FO/ensign reference files) | — | unchanged from HEAD | not touched |
