---
title: Own the Pi rejection-worker topology live fault, and give its XFAIL binding a live owner
status: backlog
score: 0.7
source: "pi-xfail-binding-reconciliation ideation, 2026-10-03: the rejection-flow Pi XFAIL names owner p17swb3375rt525fn7f8xt7e, which is archived PASSED, while the topology XFAIL stays registered."
id: 6h3teccccn3qh71yqcmjbjx4
---

The rejection-flow Pi XFAIL names an owner that no longer exists. The fault it
registers still does.

## Problem

The Pi rejection-flow journey registers a live XFAIL that names
`p17swb3375rt525fn7f8xt7e` (`finish-pi-rejection-flow`). That task is archived.
Its PASSED verdict covers the timeout repair; its own validation records that the
topology XFAIL stays registered as a deferred, separate semantic.

No active task owns the Pi rejection-worker topology. The only active topology
entity is `persist-codex-rollout-for-rejection-topology`, which is Codex-only.

So the binding violates the owner standard: a known failure must have an active
product task owner that carries a real approach. Removal is not available either,
because the durable semantic `rejection-worker-topology` is not satisfied on Pi.

## Out of scope

The shared XFAIL policy, the shared assert, and the fixture, which
`repair-pi-recorded-gate-lifecycle` also excludes. The Codex topology work. The
Claude and Codex rejection paths.

## Expected surface and tolerance

Estimate net LOC change: +40, across 2 files (the Pi driver or its extractor, and
the Pi-side test). Insertions ~+55, deletions ~-15. Tolerance: +/-20 net LOC,
+/-1 file.

Declared semantic changes: whether the Pi rejection-flow target satisfies
`rejection-worker-topology`. This task must NOT change the shared XFAIL policy,
the shared assert, the fixture, or any Claude or Codex behavior.

## Acceptance criteria

**AC-1 (VALUE) - A Pi rejection-flow run satisfies the rejection-worker-topology semantic.**
Verified by: the focused live Pi rejection-flow target at the candidate tip
completes and the durable topology semantic passes, rather than grading the
XFAIL. Independent baseline that can move the wrong way: today the same target
does not satisfy that semantic.
Falsifying edit: restore the pre-fix Pi topology handling; the target must fail again.

**AC-2 - The binding names an active owner until the fault is gone.**
Verified by: while AC-1 is unmet, the registered Pi XFAIL names this task; and it
is removed only on exact XPASS plus normal PASS evidence.
Falsifying edit: remove the binding while the fault remains; the check must fail.

**AC-3 (no-regression) -** the shared XFAIL policy, the shared assert, and the
fixture are unchanged, and `go test ./internal/ensigncycle/...` is green.

## Test plan

Primary proof owner: the focused live Pi rejection-flow target
(`TestLiveCommonRejectionFlow` with `SPACEDOCK_LIVE_RUNTIME=pi`). The durable
assert is host-neutral; the stream-based checks read the Pi session JSONL.
Deterministic tests cover the extractor path; one live Pi run proves the semantic.
