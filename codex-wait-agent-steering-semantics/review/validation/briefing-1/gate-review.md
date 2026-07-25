# Validation gate — Codex captain steering semantics

Recommendation: approve candidate `9e3a5ca0b28723850427735611110a720b8bd0b6` for merge and terminalization.

## Capability delivered

- Captain input resumes the First Officer's active loop while unresolved workers continue unchanged.
- Monitoring resumes only after captain-authorized active work is exhausted; later work invalidates any stale empty-scope observation.
- A wait return never completes a worker. Completion requires one matching final-status signal plus a parsed durable stage report.
- The harness label remains tool output and is not repeated as captain-facing interruption language.
- The change is Codex-only: no Claude, Pi, shared-core, scheduler, or runtime-API behavior changes.

## Validation evidence

Cycle-2 validation reports 4 DONE, 0 SKIPPED, and 0 FAILED. AC-1 through AC-4 all have independent evidence.

A real Codex drive kept one sentinel unresolved across captain steering, performed useful active work, resumed monitoring on the same task/epoch, observed one matching final status, and parsed the durable implementation report. Only the reduced facts were retained.

The corrected three-state matrix proves the original trace passes, active work after an empty observation fails without a renewed empty event, and monitoring becomes valid after a fresh empty event. All 36 planted mutants, focused/containment checks, live-tag compile, `go test ./...`, and `go test ./... -race` pass.

## Review disposition

The material stale-empty ordering false pass is closed by commit `9e3a5ca0`. A proposal to require a new mandatory evidence-line schema under every DONE bullet was declined because the accepted completion contract requires checklist results in a parsed durable report, not a new paragraph format.

## Decision requested

Approve to merge and terminalize this exact candidate, or revise only for a remaining material steering or completion-authority defect.
