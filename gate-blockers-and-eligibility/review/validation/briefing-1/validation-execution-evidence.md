# Validation execution evidence

## Value and lifecycle drives

- Fresh-process approval restart preserved exact `approve`, `advance/pending`, `approved-pending`, and eligibility output without status or worktree mutation.
- Fresh-process consumption atomically co-wrote workflow status and `application.state: consumed`; a second process returned `condition=consumed consumed=false` with byte-identical state.
- Stale, same-gate, cross-gate, second-close, closure-shape, and fail-closed eligibility controls passed.

## Adversarial controls

- The job-563 mutant left the old same-gate replacement test false-green, but the new direct-state control and existing cross-gate control each failed it with exit 1.
- Forced unrequested eligibility I/O, wrong split-root ownership, and pending-left-live mutants each failed their independent controls.
- Non-strict pilot application records and successorless terminal/single-stage approvals were refused with byte-identical zero mutation.
- Eight supported application shapes replayed against independent literal expected bytes; unrelated frontmatter and body bytes survived the atomic co-write.

## Repository checks

- Detached adversarial audit: no material finding.
- Focused gates/status/CLI tests: passed.
- `gofmt -w ./cmd ./internal`: no diff.
- `go test ./...`: passed.
- `go test ./... -race`: passed.
- `git diff --check`: passed.
- Candidate worktree: clean at exact head `c7612661cce857e90ae2073ac861f5b8b32b72c0`.

The three deferred semantics and their captain-owned promotion conditions are preserved in `gate-review.md` and the exact validator report.
