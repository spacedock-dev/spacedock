# Pi Runtime Frontdoor Retry Validation

State validation commit: `d7348a12`
Validated implementation commit: `255bb0da`
Worktree: `.worktrees/spacedock-ensign-pi-runtime-support`

Recommendation: REJECTED

## Commands

- PASS: `gofmt -w ./cmd ./internal` (no tracked product diffs after formatting).
- PASS: `go test ./... -count=1`.
- PASS: `go test ./... -race -count=1`.
- FAIL: `go test -tags live -run 'TestLivePi(SubagentEnsignSmoke|FrontDoorSmoke)' ./internal/ensigncycle -v -count=1`.
  - `TestLivePiFrontDoorSmoke` passed through `spacedock pi`.
  - `TestLivePiSubagentEnsignSmoke` failed because the live worker wrote `## implementation report` instead of the asserted `## Stage Report: implementation`.

## AC Summary

AC-1, AC-2, AC-3, and AC-5 passed by code inspection and tests. AC-4 was not accepted because the required combined live validation command failed, even though the front-door subtest itself passed. No Pi credential or package path was missing.

Residual risk: live Pi output remains stochastic/brittle around exact stage-report heading compliance.
