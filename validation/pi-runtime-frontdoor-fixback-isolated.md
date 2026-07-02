# Pi runtime frontdoor fixback (isolated)

Implemented a narrow validation fix for the AC-4 live heading failure on isolated branch `spacedock-ensign/pi-runtime-frontdoor-fixback`.

## Commits

- Product: `3bcaf487b8b57e810439ac9a39e7f345c8a08c34` (`test: require exact Pi live smoke report heading`)
- State: `f9643dc3` (`implementation: pi runtime frontdoor validation fix`)

## Changed files

- `.worktrees/spacedock-ensign-pi-runtime-frontdoor-fixback/internal/ensigncycle/pi_live_runner_test.go`
- `docs/dev/.spacedock-state/pi-runtime-frontdoor/index.md`

## Commands run

- FAIL before fix: `go test -tags live -run TestPiLiveSmokePromptRequiresExactStageReportHeading ./internal/ensigncycle -count=1` failed because the prompt did not contain `exact heading '## Stage Report: implementation'`.
- PASS: `gofmt -w ./cmd ./internal`
- PASS: `go test -tags live -run TestPiLiveSmokePromptRequiresExactStageReportHeading ./internal/ensigncycle -count=1`
- PASS: `go test ./... -count=1`
- PASS: `go test ./... -race -count=1`
- PASS: `go test -tags live -run 'TestLivePi(SubagentEnsignSmoke|FrontDoorSmoke)' ./internal/ensigncycle -v -count=1`; both live tests passed.

## Diff summary

Added focused live-build prompt contract coverage and changed the raw Pi live smoke prompt to require the exact `## Stage Report: implementation` heading already asserted by `assertPiLiveSmokeResult`.

## Residual risks

- Independent validation rerun still depends on Pi credentials and `pi-subagents` availability in the validation environment.
- Final validation was not self-approved.
