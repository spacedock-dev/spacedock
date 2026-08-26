# Gate decision: reduce first-officer contract reread churn

## Recommendation

Approve ideation. The measured load cost is large enough to justify designing a clear resident-contract rule, while preserving mandatory reloads after compaction or contract replacement.

## Why this matters for 0.27.x

- One audited first-officer day spent 59 reads, about 34% of tool calls, on skill files that did not change.
- The current wording permits two incompatible readings: load once per context, or reread at every gate and mutation trigger.
- The waste compounds exactly where the first officer should be spending context on evidence, routing, and captain decisions.

## Approved scope

- Define when a loaded contract remains resident and satisfies later triggers.
- Keep post-compaction re-satisfaction mandatory.
- Specify how contract replacement or version drift invalidates residency.
- Produce a falsifiable before/after measurement against the 59-read baseline.

## Exclusions

- Do not weaken gate, write-authority, or merge preconditions.
- Do not redesign host-specific context management.
- Do not prescribe the final mechanism before ideation compares the viable contract shapes.

## Proof owed at the next gate

Ideation must identify the smallest contract change that removes within-context rereads, explain invalidation boundaries, estimate the affected surface, and define a replay or live-day measurement that can fail.

