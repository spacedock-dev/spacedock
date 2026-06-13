---
id: rzp0ndsqn8bak8v1p7k3bj0w
title: Simplify the dev-workflow README + commission templates under the crew "what awesome looks like" principles
status: backlog
source: "captain (2026-06-13, this session) — j9's P1 added the \"what awesome looks like\" ethos to the FO contract; propagate the same principles to the workflow process doc (docs/dev/README.md) and the commission templates (development.md / experiment.md / refinement.md) so new workflows inherit the slim shape, not the bloat. Surfaced by the dogfood friction log (FR-2/FR-3/FR-5) + the README-simplification discussion."
started:
completed:
verdict:
score:
worktree:
issue:
sprint:
---

Align the dev-workflow process doc AND the commission templates to the crew "what awesome looks like" principles — the same ethos j9's P1 just added to the FO contract:

> - Begin with the end, be clear about the value.
> - Do the hardest things first, de-risk when it is cheap.
> - Communicate and act concisely, choose the simplest approach, JFDI.

This is the process-doc / template companion to T3 (which audits the FO *contract* refs). Distinct artifacts, same family — together they land the ethos consistently across all FO/workflow prose.

## Problem

`docs/dev/README.md` (~315 lines) carries the same content patterns that bloat the FO contract: the proof-policy doctrine is restated ~4× (ideation / validation / detached-audit / done), a ~100-line Runtime-Live-CI section is test-suite *reference* (not workflow process) embedded in the process doc, and the value/DoD is buried under stage mechanics. The commission templates new workflows are generated from (`development.md`, `experiment.md`, `refinement.md`) carry the same shape — so every commissioned workflow inherits the bloat. A one-off README fix does not propagate; the templates must change too.

## Proposed approach (seed — ideation fleshes the concrete before/after wording)

The three ethos moves, applied to the README and then mirrored into the templates:

1. **State the proof policy ONCE** *(simplest approach / do the hard thing once).* Hoist a single `## Proof policy` section (external-oracle-only · no prose-grep · spike-the-riskiest-first · code-gate-over-prose); stages link to it instead of restating it ~4×. Biggest single cut.
2. **Relocate the Runtime-Live-CI reference out** *(be clear about value / concise).* The ~100-line test-suite reference moves to a testing-guide doc (e.g. `internal/ensigncycle/README.md` or `docs/dev/runtime-live-ci.md`) with a 3-line pointer left behind. Same bulk the FO loads-but-never-needs at boot (friction FR-3).
3. **Lead with the end + DoD** *(begin with the end).* Open with what *done* looks like before the stage mechanics; compress each stage's Good/Bad. Fix the stale "No PR merge flow, mods, or lifecycle hooks are in scope" clause (friction FR-5 — mods ARE registered).

Target: README ~315 → ~130 lines, then the same structural principles applied to the three commission templates so future workflows are born slim.

## Out of scope

- The FO *contract* refs (shared-core / runtime adapters) — that is T3.
- The proof-policy *content* port to scaffolding — that is `ey` (which also touches `development.md`; coordinate, do not collide).
- Any behavioral change to the workflow stages themselves — this is prose/structure only.

## Scaffolding guardrail (dispatch shape)

`docs/dev/README.md` is FO-owned process doc (FO-editable directly). The commission templates under `skills/commission/references/templates/` are SHIPPED SCAFFOLDING — they must be changed by a dispatched worker in a worktree under test, never an FO-direct edit. So this is a dispatched task; bundle the README + template edits in one change for consistency.

## Relation to other tasks

- **T3** (`fo-contract-prose-audit`): sibling cleanup — the FO contract half. Coordinate so the ethos lands consistently.
- **`ey`** (`proof-policy-shipped-scaffolding`): touches the same `development.md` template (proof-policy propagation). Sequence to avoid stepping on each other.

## Acceptance criteria (seed — ideation defines the external proofs)

Ideation must define ACs with EXTERNAL proof, not a prose-grep over the changed docs. Candidate proof shapes to design (flagging the prose-proof challenge a "simplification" raises): a structural check that each template/README carries ONE proof-policy section (not N), a line-count floor, and — for behavior — a commission smoke test asserting a workflow generated from the simplified `development.md` produces the slim structure and still drives correctly. Specific before/after wording per the workflow's "template/skill text changes" rule.
