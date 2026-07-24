# VALIDATION GATE — Workflow-declared section reads assemble stage context by pointer (`0c`)

Chosen direction: land exact candidate `396c92a302e1ac220d552d6f16600e9f08ccb622`.

Recommendation: **APPROVE** because all six acceptance criteria have exact-tip behavioral evidence and no material finding remains.

## Capability

Workflow stages may name ordered README sections as dispatch context. The builder emits only pointers; each host-neutral bootstrap performs one live, bounded read and returns the selected canonical bytes exactly. Invalid, missing, ambiguous, or overlapping selections stop before a worker artifact can be consumed, while undeclared workflows retain their prior bytes.

## Verification

- Validation checklist: 15 DONE, 0 SKIPPED, 0 FAILED.
- Exact surface: 10 files, 560 changed lines.
- Focused, full, race, formatting, and strict documentation checks passed.
- Roborev job 1990 found no blocking defect at the exact candidate.
- GitHub run 30087598672 passed offline, Pi live, Codex live, Claude Sonnet, and Claude Opus at the exact candidate.
- Independent probes reproduced exact mixed-newline/UTF-8 bytes, tri-state ordering, refusal paths, live-read updates, legacy byte identity, and one-fetch host envelopes.

## Findings

No material finding remains. Four deferred adversarial risks and one cleanup note have concrete promotion conditions in the validation report; none affects a supported workflow.

## Decision

Approve to authorize landing exact candidate `396c92a302e1ac220d552d6f16600e9f08ccb622`; revise only for a new material defect, or hold for a named prerequisite.
