---
title: Select the actionable Codex default-headless task
status: backlog
score: "0.90"
source: Repeated Codex wrong/queued target entry; DVD run 31432758302 artifact 9080028678; filing run 31434160297 artifact 9080564383, 2026-08-10
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-product
id: 272j6s25f9mry6nxbf4yjxvt
gates:
    version: 1
    records:
        - id: gate:272j6s25f9mry6nxbf4yjxvt:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:272j6s25f9mry6nxbf4yjxvt-backlog-1
              briefing:
                id: briefing:272j6s25f9mry6nxbf4yjxvt:backlog:attempt-1:revision-1
                digest: sha256:d276a71ce0905bc7454b3df2bf92eef7fbe8e938fc575908afa97d493ba37c21
                request-digest: sha256:38daadfeab1726b731607f66dd5ce060bd32fd2eb3bb1ef583818480aa604e69
                room-ref: ./select-actionable-codex-default-headless-task/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:272j6s25f9mry6nxbf4yjxvt:backlog:1
                briefing: briefing:272j6s25f9mry6nxbf4yjxvt:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-11T05:24:15.819247Z"
                decision: approve
                reason: The repaired seed isolates one product-only Codex target-selection outcome, excludes prohibited mechanisms and other hosts, and defines the required bound-to-unbound proof ladder.
              application:
                target-stage: ideation
                state: pending
---
## Problem

Codex default-headless can ignore durable boot and next state, select an unrelated queued task, and call gate prepare before the target is actionable.

## Value

Codex headless uses durable boot and next state to select the prepared actionable default-headless task, never an unrelated queued task. It then completes dispatch acknowledgment, commits the open gate, stops without a successor, and removes its target binding.

## Scope

- Use the retained n28 and kky wrong-entity or queued-entity artifacts as the exact baseline.
- Own only the Codex default-headless binding after kky transfers it.
- Preserve the shared gate-commit product change and n28 acknowledgment mechanism.
- Do not add Pi work or change unrelated host bindings.

## Acceptance criteria

- AC-1: Exact local Codex default-headless selects the actionable target named by durable boot and next state.
- AC-2: The selected target completes pending → armed → consumed implementation acknowledgment.
- AC-3: Codex commits the prepared clean validation gate, stops open, and dispatches no successor.
- AC-4: The Codex binding is removed only after bound XPASS-green and unbound normal PASS.

## Stage Report: backlog

- DONE: Seed end value
  Codex selects the actionable implementation task from durable boot and next state. It reaches a clean, open validation gate.
- DONE: Included scope
  The scope permits the smallest declarative Codex First Officer or public binary behavior correction. It uses existing public-behavior tests. The exact local subscription moves from bound XFAIL to XPASS, binding-only removal, and unbound PASS.
- DONE: Excluded scope
  The scope excludes global or host hooks, observer references, temporary state, transcript-driven product mechanisms, product instrumentation, and new standing CI or lint. It excludes target-only CI and work for Pi, Opus, or Sonnet.
- DONE: Proof plan
  Record the current bound local Codex XFAIL baseline. Run one focused, full, and race ladder, then record the exact XPASS. Remove only the Codex binding and reconciliation. Record the exact unbound PASS, independent validation, and required PR CI.

### Summary

The Captain's current directive supersedes the seed's historical n28 acknowledgment wording. This report defines the current boundary without changes to the seed body.
