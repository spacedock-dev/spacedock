---
id: fxv06gk9nzevpbmzj6a2q3nx
title: Repeat gate consume after status-advance reports condition=ineligible instead of condition=consumed
status: backlog
source: "Found while landing collapse-gate-approval-ceremony (PR #600), 2026-08-02: PR #599 (Simplify v1 gate state schema, f566f821b) replaced the old stable-ID lookup (doc.Current.Gate) with recordForStage(doc, status), matching a gate record by the entity's live status against record.Stage. Consequence: once status advances past a gate whose approval was already consumed, a repeat `gate consume`/`gate eligibility` call no longer finds that record at all (no record exists for the new status) and reports condition=ineligible, where the old mechanism found the same original record and correctly reported condition=consumed. Caught only because collapse-gate-approval-ceremony's own AC-2 recovery tests (internal/cli/gate_consume_sync_test.go: TestGateConsumeSyncFailedRecoversWithStateCommit, TestGateConsumeRepeatAfterSupersedeRunsNoSync) assert the repeat-call diagnostic string and broke on merge with main."
started:
completed:
verdict:
score: 0.3
worktree:
issue:
---

A repeat `gate consume` call issued after the entity's status has already advanced past the consumed gate now reports `condition=ineligible` rather than the more accurate `condition=consumed`. The underlying guarantee is unaffected: the call is still read-only, byte-clean, and correctly refuses to re-apply -- only the diagnostic label changed, and verified (see below) that nothing in the codebase parses this specific string as a decision input.

## Why this is not blocking anything

Grepped the tree: no production code, skill, or doc anywhere matches `condition=consumed`/`Condition == "consumed"` as a decision point. The two CLI call sites (`internal/cli/gate_ceremony.go:76`, `internal/cli/cli.go:360`) only print the condition string for a human/diagnostic reader. The fo-gate-lifecycle skill's own Resume prose reads entity/boot state directly ("gate validate/gate eligibility are optional diagnostics, never positive-path requirements") rather than parsing this output. Classified Deferred-risk, not Material: correct today, real but narrow diagnostic-clarity loss if a human re-runs `gate consume` by hand after a gate has already been consumed and status has moved on -- they now see "ineligible" (reads as "something's wrong") instead of "consumed" (reads as "already done, safe").

## Proposed approach

In `EvaluateEligibility` (internal/gates/application.go), when `recordForStage(doc, status)` finds no record for the entity's current status, before defaulting to "ineligible" check whether any record in `doc.Records` for a *prior* stage has `Application.State == "consumed"` and, if so, report `condition=consumed` (read-only; this is a diagnostic improvement, not a behavior change to Consume's actual write path). Needs its own spike/ideation to confirm this doesn't paper over a genuinely-ineligible case (e.g., a gate for a stage the entity never reached).

## Out of scope

Any change to Consume's actual write/mutation behavior -- this is a diagnostic-string-only concern.
