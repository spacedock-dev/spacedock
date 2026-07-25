# Ideation gate: Codex captain steering resumes the active loop

## Capability

Codex First Officers treat `wait_agent` as asynchronous monitoring. Captain input
resumes the active loop while unresolved workers continue unchanged; when active
work is exhausted, the First Officer resumes monitoring the same worker. Neither
captain input nor a wait return becomes completion evidence.

## Evidence

The retained real trace observed the same task path still running after captain
input, 20 useful active-loop calls, zero cancellation or redispatch, and a later
monitoring call. The proposed ordered fixture extends that trace through matching
final status and a durable report, with planted cancellation, redispatch,
premature-monitoring, wrong-handle, stale-epoch, wait-return-only, and repeated
disclaimer controls.

## Intended boundary

The semantic change is Codex-only: the Codex runtime adapter, the idle-notification
probe/evidence, and behavioral proof. It does not change shared-core, Claude, Pi,
the `wait_agent` API, cancellation/restart state, scheduler behavior, or durable
completion authority.

The likely five files and roughly 190–285 changed lines are advisory planning
evidence. Implementation must declare and reconcile the actual surface, but file or
line arithmetic cannot block the stage or force a design reset. Only semantic
expansion beyond the boundary above requires a new decision.

## Delta from attempt 2

Attempt 2 correctly described the steering behavior but incorrectly made its
five-file estimate, LOC range, and tolerance binding gate authority. Cycle 2 removes
that arithmetic authority without changing the behavior, acceptance criteria, test
plan, negative controls, or cross-runtime exclusions.

## Recommendation

Approve the corrected Codex-only steering design for implementation.

## Decision

Approve to enter implementation; revise to change the semantic boundary or evidence
model; or hold before implementation.
