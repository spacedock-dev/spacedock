---
id: zvk9cnew2ggpaqb3wty24xtf
title: Remove redundant lint mirrors and the version-gate shell harness
status: backlog
source: "0.27 cut audit (2026-08-14), adversarially verified; captain directed filing"
started:
completed:
verdict:
score:
worktree:
issue:
---

Three verified-redundant test mirrors.

1. The wantGaps XFAIL map (internal/contractlint/live_registry_reconciliation_test.go:51-66). A hand-copied oracle of the liveXFail() calls the same test already parses. The registry doc disclaims exactly this ("does not use a copied gap oracle"). Six lockstep two-file commits in range prove the churn. Keep parseLiveGap shape validation, duplicate-target rejection, and the TODO-owner join.
2. Three install-command literals in version_gate_smoke_test.go:37-39. TestInstallHintNoDrift already anchors the same tokens to their producer, docs/site/get-started/install.md. Keep the uname, go-build, and unsupported-OS tokens (lines 40-42) - they have no other pin.
3. skills/integration/testdata/version_gate_flow.sh plus its 351-line driver version_gate_fixture_test.go. The harness verifies a hand-written re-derivation of the FO prose, not shipped bytes. Its brew command already drifted from the lint-enforced prose and everything stayed green. Nothing binds mirror to prose.

Coordination: this supersedes scope items 2-3 of remove-startup-capability-probe (dav9), which rewrite the files this entity deletes. dav9 items 1 and 4 (shared-core prose, install.md) are unaffected. Whichever lands second reports the overlap as already done.

## Problem

{Ideation fills this in.}

## Proposed approach

{Ideation fills this in.}

## Out of scope

TestDispatchAckMachineryIsAbsent (sole verifier of the hooks.json chain). The three-check failfast bundle (each check carries unique assertions). TestInstallHintNoDrift.

## Expected surface and tolerance

Estimate net LOC change: -400 or more, across ~4 files. No observable semantics change: test-only deletion.

## Acceptance criteria

**AC-1 - The change removes more lines than it adds.**
Verified by: cumulative line delta against origin/main is negative.

**AC-2 - An XFAIL retire is a one-file edit.**
Verified by: wantGaps absent from the tree, parseLiveGap validation still green.

**AC-3 - The suite stays green.**
Verified by: go test ./... and go test ./... -race pass.

## Test plan

Deletion plus the existing suite. The surviving drift and shape lints are the regression floor.
