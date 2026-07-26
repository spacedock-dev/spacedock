---
id: ppfr0qqafmaqmhy3kyj418dw
title: "Claim mechanism hardening — the edges the walking skeleton holds fixed"
status: backlog
source: "Skeleton-first re-carve of multi-fo-coordination, 2026-07-26: entity-session-claim-lease was cut to the minimum walking skeleton (acquire, disclose, refuse at one door). These are the edges it defers, each with the condition its journey holds fixed."
started:
completed:
verdict:
score: 0.6
worktree:
issue:
---

Harden the entity-session claim beyond the walking skeleton: concurrent acquisition, lease expiry, ownership transfer, duplicate fail-closed reads, multi-host identity, and the remaining mutation doors.

## Problem

`entity-session-claim-lease` is the walking skeleton for multi-FO coordination. It proves the journey — session A claims and drives an entity, session B cold-boots, sees the claim, attempts a mutation and is refused — with everything else held fixed. A skeleton is legitimate only because its journey never triggers the hazards it defers. The moment the journey widens, each held condition unfixes and the corresponding hazard becomes live.

This task owns those edges. Each was deferred against a named condition, and each is designed rather than unknown: the skeleton's spike already exercised the two-writer race and found the split-brain hazard, so these are informed deferrals with recorded evidence, not gaps.

## The deferred edges and the condition each was held under

| edge | condition the skeleton holds fixed | unfixes when |
| --- | --- | --- |
| Concurrent acquisition and lost-race rollback | one acquirer at a time; the journey never races | two sessions can acquire independently |
| Duplicate claim fields fail closed (the spike's split-brain) | no concurrent write, so no divergent placement | concurrent acquisition lands |
| Lease expiry, EXPIRED marking, owner renewal | the journey is short and the owner stays alive | a drive outlives a session, or a session dies without releasing |
| Ownership transfer (`--take --from`) and release | the journey ends at the refusal; B never obtains the entity | a handoff is needed, which the channel member introduces |
| Multi-host identity resolution and ambiguity refusal | both sessions run on one host | sessions run on different runtimes |
| The remaining mutation doors | the journey exercises one door | any other door can mutate a claimed entity |
| Exact refusal rendering grammar | the refusal names the owner; format is unspecified | operators depend on the message shape |

## Recorded evidence to build against

The skeleton's `## Spike: two-writer claim CAS over real Git` section is the design input for the first two rows and must be read before designing them. It established that a real two-clone race resolves to exactly one owner through ordinary git compare-and-swap, AND that divergently placed claim writes automerge into a pushed split-brain carrying two claim lines — so textual conflict is not a git guarantee, and canonical fixed placement plus duplicate-fail-closed reads are load-bearing rather than defensive.

## Out of scope

- The claim primitive itself, boot disclosure, and the guard at the skeleton's door — those are the skeleton's.
- Session scope semantics and the cross-session channel — separate members of the same sprint.
- Any second writer of claim state; hardening extends the existing verb and publisher.

## Acceptance criteria

Ideation fills these in. Each edge needs an end-state property with a falsifier, and each must name the condition whose unfixing makes it necessary — an edge whose triggering condition is still held fixed elsewhere does not need to ship here.

## Test plan

Ideation fills this in. The skeleton's real-Git two-clone harness is the substrate; do not build a second one.
