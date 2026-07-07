---
title: "status --next orders dispatchables by score — a code-gated value-first event loop"
status: backlog
source: "Principle-derived contract review, 0250 Commander session 2026-07-07 (captain-requested filing). Finding: the event loop («dispatch.next-action» step 2) says 'dispatch any newly ready entity' — arrival order; the frontmatter score: field exists but no contract line or binary path ever orders by it. Ranked highest combined efficacy x token impact of the review's residuals: ensign lifecycles are the dominant token line-item (observed 130k-230k resident per worker this sprint; a full member pipeline runs to high hundreds of k), so one misallocated dispatch in a multi-dispatchable autonomous drive outweighs every prose optimization shipped in 0250."
started:
completed:
verdict:
score: 0.6
worktree:
issue:
id: 21150nfmgckg4ykh80d4vrqw
---

The FO's scheduler is value-blind: `status --next` emits dispatchables unordered by declared priority, and the contract's loop dispatches in arrival order, so autonomous multi-entity drives spend their freshest, most expensive worker lifecycles on whatever unblocked first. Direction (zm doctrine — code gate, zero resident bytes): sort `--next`'s dispatchable output in the binary by `score` descending (tie-break: oldest `started`/filed first; document the null-score position), plus one contract half-line pointing at the guarantee. Acceptance sketch: value — on a fixture with >=3 dispatchables of distinct scores, the dispatch order matches score order (behavior test over `--next --json`; baseline: current arrival order moves the wrong way); mechanism — the sort ships with ordering tests incl. ties and null scores. Read-path projection change in internal/status; low blast radius.
