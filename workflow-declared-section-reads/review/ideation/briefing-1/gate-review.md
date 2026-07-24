# IDEATION GATE — Workflow-declared section reads (`0c`)

Recommendation: **APPROVE the cycle-3 design and advance to implementation.**

## End value

A dispatched worker receives exactly the workflow sections its current stage declares, through the existing one-file bootstrap fetch. Workflows without declarations remain byte-identical.

## Chosen design

- Generic ordered `context-sections` metadata supports defaults, stage replacement, and explicit empty clearing.
- `dispatch build` and `show-stage-def` use one shared in-memory resolution operation over one README buffer per invocation.
- The existing stage extractor remains the authority for legacy stage bytes; the fence-safe scanner remains the authority for selected headings. Raw half-open spans provide one overlap coordinate without introducing a third parser.
- Missing, ambiguous, malformed, repeated, or intersecting selections fail before stage work.
- The pointer has honest live-read semantics: build preflights one version; bootstrap validates and adopts the then-current valid version.
- Declared context has an exact canonical rendering for LF, CRLF, lone CR, Unicode separators, UTF-8, internal blanks, trailing blank lines, component separators, and final termination.

## Evidence

- The tri-state YAML spike proved inheritance, explicit clear, replacement order, and wrong-kind rejection.
- The real-parser mixed-buffer spike exercised UTF-8, CRLF/lone-CR, an exotic separator, parent/child/stage/disjoint spans, and exposed then corrected a CRLF double-advance panic.
- Independent staff review ran three rounds. The final re-review reports **APPROVE** with no remaining material finding and confirms no snapshot, status projection, FO assembly, bootstrap multi-read, third parser, or second harness.

## Surface

Implementation expects 8–10 files and 380–560 changed lines, with a hard tolerance of 12 files or 700 changed lines. It reuses existing status, dispatch, host-envelope, golden, and approval-gated live lanes.

## Decision

Approve this bounded design. Revise if implementation changes undeclared bytes, adds another read/projection/parser/harness, or cannot preserve the raw-span and mixed-newline contracts demonstrated by the spikes.
