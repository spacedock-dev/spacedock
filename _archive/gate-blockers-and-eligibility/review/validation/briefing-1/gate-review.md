# Validation gate: one-use gate applications and eligibility

Recommendation: **APPROVE** exact candidate `c7612661cce857e90ae2073ac861f5b8b32b72c0` for landing.

## Capability

The candidate extends the canonical gate recorder with a one-use application layer. A durable approval can remain `advance/pending`, be consumed exactly once by atomically co-writing workflow status and `application.state: consumed`, become `superseded` if the reviewed input drifts, and surface a pure fail-closed eligibility condition. The implementation uses the recorder's sole `gates:` write surface and keeps dispatch retry/recovery separate from authorization consumption.

## Validation evidence

- Fresh validation passed all 10 checklist items and AC-A1 through AC-B1 at exact clean head `c7612661`.
- Fresh-process restart and consumption drives proved durable `approved-pending`, atomic status/application co-write, exactly one authorization consumption, byte-identical second-pass state, and no hidden dispatch receipt/recovery machinery.
- Staleness, same-/cross-attempt supersession, second-close refusal, closure-shape, and full fail-closed eligibility tables passed.
- The mandatory job-563 attack reproduced the original same-gate test's false-green mutant, then killed that mutant with a new direct-state control and the existing cross-gate control.
- Legacy/pilot non-strict application shapes and terminal/single-stage approvals without successors were explicitly refused with byte-identical zero mutation; supported successor gates remained green.
- Workflow ownership, split-root behavior, requested-only eligibility I/O, independent literal canonical-byte replay, and unrelated frontmatter/body preservation all passed, including their negative mutants.
- Detached adversarial audit, focused tests, `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, `git diff --check`, and exact-head cleanliness all passed.

## Scope and deferred semantics

The final surface is 14 files, 1,146 additions and 52 deletions: 642 production additions, 480 test additions, and 24 documentation additions, within the approved 2x ceiling. It adds no blocker evaluator, hold-authoring command, Subspace coupling, second gates writer, compatibility layer, successorless terminal semantics, dispatch receipt, or crash-recovery state machine.

- Blocker/hold authoring promotes only when a real cross-session, machine-enforced predicate consumer appears.
- Strict-v1 compatibility promotes only if a future captain promises compatibility to a released external consumer.
- Successorless approvals promote only when a real supported commissioned workflow gates a stage with no successor.

## Decision

- `approve`: authorize landing and terminal completion for the exact candidate.
- `revise`: return a concrete material finding to implementation.
- `hold`: retain validation for a named prerequisite.

