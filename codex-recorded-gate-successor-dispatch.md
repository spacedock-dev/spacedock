---
title: Repair Codex recorded-gate successor dispatch proof
status: backlog
score: 0.9
milestone: 0.27.0
source: "z5 validation cycle 7, exact head 64f2f8773 on base 48a7ea0d9"
id: 2g7akd6m5s3shwsbkrwpbhdn
---

## Problem

A Codex First Officer can consume a recorded gate and build a successor dispatch. The durable oracle can still reject the required commit ordering and ancestry.

This failure blocks the required Codex live lane for z5. The z5 launcher change does not own the dispatch or state-commit integration. Keep this repair in a separate task.

## Acceptance criteria

**AC-1 (VALUE)** The recorded-gate lifecycle proves one durable successor dispatch after consume. Verified by: `TestLiveCodexSharedScenarios/recorded-gate-lifecycle` passes on the current main tip, with JSONL and commit ancestry evidence.

**AC-2 (VALUE)** The oracle rejects missing, late, duplicate, or unrelated successor dispatch evidence. Verified by: existing negative fixtures pass without weakening the predicate.

**AC-3 (VALUE)** The repair identifies the owner of the failure. Verified by: the report classifies runtime/state behavior versus oracle behavior and retains the exact failure artifacts.

## Boundary

Do not change the z5 launcher candidate. Do not weaken the durable oracle. Do not add a new provider lane.

## Test TODO

- TODO: Restore green `TestLiveCodexSharedScenarios/recorded-gate-lifecycle` after the runtime or state-commit repair.

## Evidence

The full Codex lane and one retained-artifact rerun failed the same predicate at z5 cycle 7. Artifacts: `/tmp/spacedock-codex-final-live-artifacts/codex-shared-scenarios/recorded-gate-lifecycle`.
