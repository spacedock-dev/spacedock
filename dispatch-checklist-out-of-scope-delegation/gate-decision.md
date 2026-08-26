# Gate decision: checklist authority boundary

## Recommendation

Approve this design for implementation.

The defect was created when a First Officer placed an FO-owned status transition into two worker checklists. The design fixes that decision point with one host-neutral authority rule.

## Proposed change

Add one paragraph to the shared First Officer dispatch contract:

- Worker checklist items describe worker-owned body and report outcomes only.
- Workers never receive status, frontmatter, or stage-advance obligations.
- After a completion signal and accepted report, the First Officer performs the applicable `spacedock` transition.

The dispatch builder continues transporting free-form checklist bytes unchanged. No binary lint, new schema, worker mutation command, or fixture rewrite is added.

## Surface and semantics

- Expected change: +1 net line in one file, tolerance 0 to +3 lines and exactly one file.
- Command grammar, stored formats, checklist transport, fixtures, and graders remain unchanged.
- Authority is clarified, not transferred: workers own deliverables and reports; the First Officer owns entity-state transitions.

## Required proof

- Re-run the unchanged live Claude journey whose baseline delegated `done` in 2 of 2 worker checklists.
- Candidate target: 0 of 2 delegated transitions.
- Pair every durable status delta with a successful root First Officer `spacedock` call after worker completion.
- Confirm zero descendant worker state commands or frontmatter edits using `parent_tool_use_id` attribution.
- Run full and race Go suites; they are regression checks, not substitutes for the live proof.

## Captain action

Approve to implement the one-paragraph rule, or hold if live model behavior should not be governed through the shared First Officer contract.
