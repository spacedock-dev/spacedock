# Validation gate: Git-only keep-moving attribution

## Capability

Codex keep-moving journeys are graded from each ticket’s durable Git history instead of provider transcript dialects. Supported ordinary, atomic, complete-batch, and corroborated delayed-persistence histories pass; reordered, stale, partial, foreign, premature-terminal, reopened-questioned, and cross-attributed histories fail.

## Reviewed snapshot

- Candidate: `fba615fdf`
- Surface: 10 files, `+875/-1,449` (net `-574`)
- Provider observers removed: four transcript/parser files

## Evidence

- All four retained real live roots grade `3/3`.
- The deterministic matrix covers per-ticket identity, room binding, ordering, terminal state, archive state, batch boundaries, and delayed persistence.
- `go test ./...`, `go test ./... -race`, formatting, and diff checks passed on unchanged Go bytes; the final docs-only correction passed focused and full checks.
- Independent validation reproduced AC-1 through AC-4 and recommends PASSED.

## Findings

- `verdict: questioned`, split terminal fields, and missing-current-archive after external removal are declined outside supported behavior, each with a promotion condition.
- No Material supported-workflow finding remains.

## Recommendation

Approve exact `fba615fdf` for PR and required CI. This should land before rebasing PR #584 so Codex grades commissioned completion from durable state.
