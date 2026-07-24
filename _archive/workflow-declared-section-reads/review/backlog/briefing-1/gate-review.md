# BACKLOG GATE — Workflow-declared section reads (`0c`)

Recommendation: **APPROVE entry to ideation under the reduced design below. Do not dispatch yet.**

## Value

A dispatched ensign receives only the workflow README sections declared for its stage. The workflow author selects the context; the First Officer neither parses headings nor copies section text.

## Approved direction

- Workflow frontmatter may declare ordered `context-sections` in stage defaults and individual stages.
- `dispatch build` resolves the current stage, validates every selector before spawn, and emits the existing stage-definition fetch command.
- The ensign reads the generated dispatch file and runs that command.
- `dispatch show-stage-def` derives the selectors from `--workflow-dir` and `--stage`, then returns the selected sections with the stage definition.
- Default sections precede stage sections. Exact duplicates are removed without changing order.
- Missing or ambiguous headings fail with the workflow path, stage, and selector.
- Existing workflows and dispatch packets remain byte-identical when they declare no sections.

## Deliberate cuts

The First Officer does not construct section content or pass `--include-section`. The design adds no `status --read` projection, dispatch-time content snapshot, ensign discovery step, bootstrap multi-read protocol, general include language, or product-specific heading semantics.

The current stage-definition fetch already reads the workflow after dispatch. This task does not add revision pinning solely for selected sections. A consistent snapshot for all dispatch inputs would be a separate design.

## Ideation obligation

Ideation must revise the seed specification and acceptance criteria to match this boundary, choose the smallest frontmatter representation, and augment the existing stage-definition and dispatch fixtures. It must reuse the current fence-safe heading reader and avoid a second parser.

## Decision

Approve to open the ideation application with this direction bound to the Briefing. Leave the application pending; do not consume it or dispatch an ensign.
