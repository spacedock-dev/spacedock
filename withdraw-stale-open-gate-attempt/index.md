---
title: Withdraw a stale open gate attempt without fabricating a decision
status: backlog
source: "Observed by the Subspace Shaping FO on 2026-07-26 after a legitimate sprint re-scope left a frozen request-backed attempt open with no truthful exit."
started:
completed:
verdict:
score: 1.0
worktree:
issue:
sprint: durable-decisions
id: 0m6vtrw4qh9w4x6bn06x5hen
gates:
    version: 1
    current:
        gate: gate:0m6vtrw4qh9w4x6bn06x5hen:backlog
    records:
        - id: gate:0m6vtrw4qh9w4x6bn06x5hen:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:0m6vtrw4qh9w4x6bn06x5hen-backlog-1
              briefing:
                id: briefing:0m6vtrw4qh9w4x6bn06x5hen:backlog:attempt-1:revision-1
                digest: sha256:9bfedeb38906e04bae528cedfdb96f101efaa1d63c819b44922a8ee6e5db60f6
                digest-domain: canonical-bytes
                request-digest: sha256:0526bd61039ea579f7595d23e5ccfd8bd3d3f18ee7ce5b64b211456df37f8524
                room-ref: ./review/backlog/briefing-1
---

A prepared request-backed gate attempt is correctly frozen, but a legitimate re-scope currently has no truthful operation that retires the open attempt. The operator must either record a `hold` or other Resolution against a Briefing it already knows is stale, or hand-edit/revert gate state. Both paths corrupt the meaning of the durable record.

## Boundary to shape

Preserve the immutable Briefing and room as historical evidence while giving an authorized operator one mechanical way to withdraw the still-open attempt with a reason, leave no decision or application, and prepare a new current attempt. Do not reinterpret withdrawal as approval, revise, hold, rejection, or application supersession. Reuse the existing recorder authority, locking, atomic replacement, state commit, and boot projection; do not add a compatibility path or a second writer.

The public command grammar, minimum durable withdrawal shape, actor authority, and relationship to re-prepare belong to ideation. Prefer the smallest semantic addition that makes the observed re-scope truthful.

## Acceptance criteria

**AC-1 (VALUE)** After a prepared open attempt becomes stale before decision, an agent can retire it and prepare a replacement without recording a false captain decision, hand-editing `gates:`, or deleting the frozen room.

**AC-2 (INTEGRITY)** The withdrawn attempt remains byte-verifiable historical evidence, carries no Resolution or application, cannot be recorded or consumed, and no withdrawal failure changes entity or room bytes.

**AC-3 (RECOVERY)** Boot/status and the First Officer lifecycle distinguish the withdrawn history from the new current open attempt and provide one unambiguous next action after restart.

## Test plan

Exercise a real prepared folder-form room, withdraw it with a reason, cold-boot, prepare attempt N+1, and record its actual decision. Mutate actor, reason, attempt currency, room bytes, and repeat-withdraw inputs; each refusal must be byte-clean. Verify the original room remains frozen and `state commit` durably captures the entity plus both attempt rooms without sibling dirt.
