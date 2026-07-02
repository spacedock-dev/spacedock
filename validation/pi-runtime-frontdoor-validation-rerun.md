# Pi Runtime Frontdoor Validation Rerun

State commit: `46261d5b10bb895802539f5f64c70dece691f8bf`
Product fix commit inspected: `3bcaf487b8b57e810439ac9a39e7f345c8a08c34`
State fix report commit inspected: `f9643dc3`
Worktree validated: `.worktrees/spacedock-ensign-pi-runtime-frontdoor-fixback`

Recommendation: PASSED

## AC Results

- AC-1: PASSED — `spacedock pi` remains registered and uses Pi-native extension/skill paths without Claude/Codex-only runtime tokens.
- AC-2: PASSED — `install --host pi` remains accepted/idempotent and does not call Claude/Codex plugin install commands.
- AC-3: PASSED — `doctor --host pi` reports actionable missing and healthy Pi runtime state.
- AC-4: PASSED — product fix requires the exact `## Stage Report: implementation` heading, and the combined live Pi command passed for both raw subagent and front-door smoke tests.
- AC-5: PASSED — baseline and race suites passed, preserving Claude/Codex behavior.

## Commands

- PASS: `gofmt -w ./cmd ./internal` (no tracked product diffs after formatting).
- PASS: `go test ./... -count=1`.
- PASS: `go test ./... -race -count=1`.
- PASS: `go test -tags live -run 'TestLivePi(SubagentEnsignSmoke|FrontDoorSmoke)' ./internal/ensigncycle -v -count=1`.

## Residual Risks

- Future live validation still depends on local Pi credentials and an available `pi-subagents` package path.
