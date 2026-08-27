# Ideation gate: make First Officer contract loads resident within one context

## Recommendation

Approve implementation. Add one host-neutral residency/invalidation rule to the shared First Officer core so an unchanged loaded contract body satisfies later same-context triggers while compaction and direct replacement evidence still invalidate it.

## Value and proof

- The reproducible historical slice contains 28 deferred-body reads among 151 root tool calls: 18.5%.
- Removing 23 redundant same-context reads leaves 5 required reads among 128 calls: 3.9%.
- AC-1 caps the implementation replay at 5 reads and 4.0%.
- A separate safety suffix proves lazy exactly-once reload after compaction and replacement evidence.
- Gated mutations must preserve `gate-load → write-load → mutation`.
- Gated terminal mutations must preserve `gate-load → write-load → merge-load → transition`.
- Every adjacent swap and omitted required load must fail the trace oracle.

## Surface and boundaries

| Surface | Estimate | Tolerance |
| --- | ---: | ---: |
| Net LOC | +2 | +1 to +4 |
| Files | 1 | exactly 1 |

The only expected file is `skills/first-officer/references/first-officer-shared-core.md`. Projected size is 23,170 bytes, 330 bytes below the component cap. Command grammar, output, stored formats, authority, gate decisions, merge decisions, runtime adapters, and host-specific context management remain unchanged.

## Independent review

The reviewer reproduced 28/151 and 5/128 from the bound session slice, verified both prerequisite-order sequences and their falsifiers, and returned PASS with no remaining findings.

## Decision effect

Approval dispatches implementation in an isolated worktree. Rejection returns the design for correction without changing product files.
