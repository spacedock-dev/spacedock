# Delivery gate: name the failing completion sub-check

## Recommendation

Approve delivery. Validation found the candidate satisfies AC-1 through AC-4 with behavioral and on-disk-state evidence, and the detached audit found no candidate defect.

## Delivered outcome

- Four core completion failures now emit four distinct, actionable remedies instead of one generic message.
- An annotated near-miss such as `DONE (annotation):` refuses visibly with its line and repair instead of silently advancing.
- Consumed nonterminal journeys can complete directly or through merge guard without `ineligible`, while pending terminal approval remains merge-guard-only.
- Local-HEAD durability, dirty-path boundaries, latest-cycle selection, and force behavior remain intact.

## Evidence that matters

- Real CLI tables assert exit 1, empty stdout, exact stderr, and unchanged entity bytes for all four failure classes.
- Parser and CLI controls cover canonical and near-miss checklist forms.
- Focused status and terminal tests, `go test ./...`, and `go test ./... -race` passed.
- A detached throwaway-checkout audit found no outcome, evidence, material, deferred-risk, or polish finding in the candidate.
- Commit `d24e90a8e` is +175 net lines across six files, inside the approved +145 +/-40 and six +/-1 bounds.

## Delivery boundary

Approval authorizes terminal delivery through the workflow merge guard. It does not expand the diagnostic order, Git durability rules, or pending-gate authority model.

