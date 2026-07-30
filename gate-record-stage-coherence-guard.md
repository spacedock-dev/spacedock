---
title: "Reject gate recording when the Briefing and current workflow stage disagree"
status: backlog
source: "se0 live CI recovery, 2026-07-28: TestLiveDefaultHeadlessStopsAtGate reproduced a validation Briefing durably bound as an implementation gate; captain separated the feature defect from CI restoration."
score: 0.85
sprint: durable-decisions
group: recorder
sprint-readiness: ready
issue:
id: q3vpb8hes1b3k3f1jps1kvpk
gates:
    version: 1
    current:
        gate: gate:q3vpb8hes1b3k3f1jps1kvpk:backlog
    records:
        - id: gate:q3vpb8hes1b3k3f1jps1kvpk:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:q3vpb8hes1b3k3f1jps1kvpk-backlog-1
              briefing:
                id: briefing:q3vpb8hes1b3k3f1jps1kvpk:backlog:attempt-1:revision-1
                digest: sha256:d84af4e055a57f211a97dbaa3faa7ec5c8dbc2dbcb85b0134c3096bcbc215652
                digest-domain: canonical-bytes
                request-digest: sha256:67af66f54068626fea34e776301eb17b8f5e59bee743606828c6f3713ca394af
                room-ref: ./gate-record-stage-coherence-guard/review/backlog/briefing-1
---

## Problem

`gate record` can bind a stage-qualified validation Briefing while the task is still in `implementation`. `RecordSemantic` holds the entity lock, but `recordBriefingLocked` has no workflow-stage authority check. `initialIDs` recognizes a structured Briefing ID only when its stage matches current status; on mismatch it silently falls back to `gate:<entity>:<current-stage>`. The result is an implementation attempt whose retained Briefing is qualified for validation.

The retained `se0` Sonnet journey demonstrated the supported trigger: a task resumed at `implementation` while a later validation Briefing remained on disk. The First Officer read the current-stage failure, selected the retained validation evidence anyway, and `gate record` accepted it. The final durable state remained `status: implementation`, `gate:...:implementation`, with a validation-qualified Briefing. This is a genuine recorder authority defect, not a test-oracle defect.

`TestLiveDefaultHeadlessStopsAtGate` is temporarily TODO under this task. It is one test definition executed by the Sonnet and Opus Claude CI matrix legs; no Codex, Pi, or shared gate/rejection/keep-moving scenario is quarantined. Re-enable it when this task lands.

## Recovered semantic work

An uncommitted spike established the smallest credible enforcement seam:

- Pass the already-resolved workflow directory from `RecordSemantic` into `recordBriefingLocked`.
- Parse the existing stage-qualified Briefing ID convention once and share it with `initialIDs`.
- Before any read-transition/write mutation, compare the qualified Briefing stage with authoritative current status.
- Resolve current stage metadata from the workflow definition and require an actionable `gate: true`, nonterminal stage for ordinary Briefing recording.
- Keep historical `gate record --round` behavior unchanged; review rounds intentionally need not equal current workflow status.
- Reject byte-clean and remove the entity lock on every failure.

The proposed stable diagnostics were:

```text
Briefing stage validation does not match current workflow stage implementation
current workflow stage implementation is not an actionable gate:true stage
```

A prose-only gate-entry precondition was also exercised live. Sonnet loaded it and ignored it, then reproduced the same malformed binding. The shipped contract should still explain recovery—resume the authoritative current stage after refusal—but prose is not the authority mechanism.

The spike identified likely changes in `internal/gates/operation.go`, existing gates/application/CLI tests, and `docs/specs/gate-resolution-frontmatter-contract.md`. Its first test draft duplicated too much fixture scaffolding; ideation should preserve the semantic findings while consolidating refusal cases into existing tables/helpers. V1 is unreleased: add no compatibility path or migration fallback. Existing malformed records may remain readable; new malformed writes must fail.

## Acceptance criteria

**AC-1 (VALUE) - A captain cannot receive or spend a decision whose Briefing belongs to a different workflow stage than the recorded gate.**
Verified by: a command-level fixture beginning in `implementation` with a validation-qualified Briefing exits nonzero, reports the specific mismatch, leaves entity bytes unchanged, and leaves no lock residue.

**AC-2 - Ordinary Briefing recording is allowed only at the authoritative current `gate: true`, nonterminal stage.**
Verified by: package and CLI cases for non-gated refusal plus an existing genuine gated-stage positive.

**AC-3 - Valid recorder behavior is preserved without compatibility machinery.**
Verified by: matching/unqualified Briefings on a genuine gated stage still record; `gate record --round` remains green and unchanged.

**AC-4 - The quarantined default-headless live journey is restored in both Claude matrix legs.**
Verified by: remove the TODO, then run focused Sonnet and Opus `TestLiveDefaultHeadlessStopsAtGate` against the exact candidate tip, followed by the applicable live CI lane.

## Test plan

Add the smallest failing package and CLI cases before production changes. Reuse existing positive fixtures and table helpers; do not create a parallel gate harness. Run focused gates/CLI tests, `go test ./...`, `go test ./... -race`, focused Sonnet and Opus default-headless journeys, then the affected Claude live CI lane once local proof is green.
