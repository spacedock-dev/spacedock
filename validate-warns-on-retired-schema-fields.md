---
id: wyzg6knr8whtmg79mxkb78jg
title: status --validate warns on retired gate-schema fields and exits zero
status: ideation
source: "email-triage field report 2026-08-26: historical batches carry retired gate-schema fields (current, digest-domain — the latter absent from today's source), so --validate exits non-zero forever on any workflow with history; captain ruling: no migration, just warn"
started:
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:wyzg6knr8whtmg79mxkb78jg:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:wyzg6knr8whtmg79mxkb78jg-backlog-1
              briefing:
                id: briefing:wyzg6knr8whtmg79mxkb78jg:backlog:attempt-1:revision-1
                digest: sha256:ae90a8befc3588f499f77645840d99414c911362ce1184fe1bb180afe01575a1
                room-ref: ./validate-warns-on-retired-schema-fields/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:wyzg6knr8whtmg79mxkb78jg:backlog:1
                briefing: briefing:wyzg6knr8whtmg79mxkb78jg:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-26T21:59:45.808528Z"
                decision: approve
              application:
                target-stage: ideation
                state: consumed
---

A workflow with history cannot pass `status --validate` because entities written under an old gate schema carry fields the current schema retired. A check that is always red cannot fail meaningfully. Captain ruling: do not build a migration. Downgrade retired-field findings to warnings and exit zero when they are the only findings.

## Problem

{Ideation fills this in. Seeded: strict schema parsing treats a retired field (for example digest-domain) as an unknown field, which is an error; old entities are immutable history, so the error never clears. Real corruption in current-schema entities must still exit non-zero.}

## Proposed approach

{Ideation fills this in. Seeded: a known-retired-fields list; findings on those fields print as warnings; the exit code reflects only current-schema errors. No rewrite of historical entities — captain ruling 2026-08-26: "no migration, just warn".}

## Risk evidence

{Backlog: the email-triage run's permanently non-zero --validate over clean fresh batches decides design should start.}

## Out of scope

Any migration or rewrite of historical entities. The flat-room conversion warnings (already warnings).

## Expected surface and tolerance

Estimate: production +25 across 2 files; proof +40 across 1 file. {Backlog seed; ideation refines.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - A workflow whose only findings are retired-field warnings exits zero from --validate, and a current-schema corruption still exits non-zero.**
Verified by: {ideation refines — seed: two fixtures — one with a retired-field entity only (exit 0, warning printed), one adding a real current-schema error (exit non-zero). Falsifying edit: drop the retired-field list — the first fixture reds.}

## Test plan

{Ideation fills this in.}

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}
