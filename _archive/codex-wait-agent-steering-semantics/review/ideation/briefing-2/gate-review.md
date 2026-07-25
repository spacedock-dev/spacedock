# Ideation gate review

Capability: Codex First Officers describe captain input as active-loop resumption while unresolved workers continue unchanged, then resume asynchronous monitoring only when they become idle again.

Evidence: the ideation spike observed the same worker running after captain input, 20 useful active-loop calls, zero cancellation or redispatch, and a later monitoring call.

Scope: five Codex-only contract, probe, evidence, and behavioral-test files; no shared-core, Claude, Pi, runtime API, scheduler, or production Go change.

Validation target: planted controls must reject cancellation, redispatch, premature monitoring, wrong or stale completion, wait-return completion, and repeated captain-facing interruption disclaimers.

Authority correction: the prior attempt was consumed without Captain approval and is historical only. This successor attempt is open and undecided.

Recommendation: approve implementation within the declared five-file, 190–285 changed-LOC surface and ±60 LOC tolerance.
