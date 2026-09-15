---
title: Use semantic branch and worker names
status: ideation
source: Captain request 2026-09-15; GitHub issue 624
started: 2026-09-15T18:27:49Z
completed:
verdict:
score: 0.8
worktree:
issue: spacedock-dev/spacedock#624
pr:
mod-block:
id: 6es505tn1zz2597hetvnqn7y
---

Use short, readable task names for public branches and dispatched workers. Keep agent implementation details out of these names.

## Problem

Issue #624 requests semantic public branch names instead of spacedock-ensign/<slug>. The captain also requested better subagent namespaces and shorter slugs. Current capWorkerName replaces long readable slugs with an opaque ID when the worker name exceeds 64 characters. Branch creation derives worker_key/slug in build.go and stamp.go. Reconciliation requires the spacedock-ensign- prefix. A dispatch-only rename breaks lifecycle routing.

## Proposed approach

Use short semantic slugs at task creation and put descriptive detail in titles. Example: title 'Schedule live journeys with committed duration hints', slug and branch 'ci-duration-hints', worker 'ci-duration-hints-ideation'. Codex maps unsupported separators to its supported underscore spelling.

Derive branches from the slug and workers from slug-stage. Preserve stable identity and readable text when names require shortening. Preserve active branch and worker names; do not rename existing public PR heads. Decide the smallest compatibility surface during ideation, including whether a workflow naming option is necessary. Avoid general template engines, new identity ledgers, and unnecessary aliases.

This task owns naming and lifecycle compatibility. Use a new worktree for implementation, based on current main. Do not change other active tasks, their branches, or remote metadata.

## Risk evidence

At origin/main 438053493838dc70c9478b3d991309d566783e85, internal/dispatch/build.go derives names and branches; stamp.go creates and verifies worktree branches; reconcile.go decompose requires spacedock-ensign-. Existing build_namecap_test.go and build_stamp_test.go cover name limits and branch registration. Spike new and legacy names through actual dispatch and reconciliation before proposing implementation.

## Out of scope

Bulk renaming active tasks or branches, changing stable entity IDs, provider changes, scheduler changes, automatic PR creation, and unrelated contract cleanup.

## Expected surface and tolerance

Ideation must record a concrete net-line and file estimate with tolerance before implementation. Likely owners: internal/dispatch naming, branch creation, reconciliation, their existing tests, and the minimal filing guidance or workflow configuration needed for short semantic slugs.

## Acceptance criteria

**AC-1 — New public branches and worker names retain readable task meaning without the agent prefix.**
Verified by: real dispatch build/stamp in a temporary Git workflow using an independently specified slug and stage. Verify resulting branch, worktree registration and worker envelope; a retained agent prefix or unexpected opaque replacement fails.

**AC-2 — Naming remains safe under host limits and ambiguous task names.**
Verified by: focused long-name, normalization, collision and stage/cycle tests against existing naming proof owners. Distinct active tasks must remain distinguishable. Include sequential, sd-b32 and slug identities where supported.

**AC-3 — Existing active workers and branches remain usable through lifecycle routing.**
Verified by: actual legacy and new dispatch/reconcile fixtures covering reuse, retries, completion and cleanup. Existing recorded worktrees keep their branches. Rejecting a valid legacy worker or targeting the wrong task fails verification.

**AC-4 — Supported runtime dispatch preserves the intended name and lifecycle.**
Verified by: a targeted local live journey on the affected runtime path, observing native worker names and completed durable stage evidence. Ideation must name the journey and concrete expected outcome. Offline envelope checks alone do not prove host behavior.

## Test plan

Reuse existing naming, stamp and reconciliation tests. Before implementation, spike the smallest real CLI path for legacy and proposed names in isolated temporary workflows. Compare supported host naming rules and reject invalid names before mutation. Record the exact runtime proof and cost. Run required normal/race suites, formatting and the applicable detached audit after implementation. No remote push or expensive CI until separately authorized and individually verified.

### Feedback Cycles

