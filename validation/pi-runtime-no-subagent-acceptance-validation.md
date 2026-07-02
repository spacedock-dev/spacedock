# Validation: Pi runtime no subagent acceptance contract

Recommendation: PASSED

Validated isolated implementation worktree `.worktrees/spacedock-ensign-pi-runtime-no-subagent-acceptance` at product commit `7119e74d83b8d7077a8824d03399fa3928171b5e`.

Findings:
- Pi first-officer runtime still permits dispatch through `subagent(...)`.
- Spacedock stage dispatches through `pi-subagents` now explicitly forbid `subagent(... acceptance: ...)`.
- The rationale is present: acceptance requirements belong in task prompt/dispatch content; Spacedock gates via entity stage reports, product/state commits, and independent validation, not same-agent finalization.
- `TestPiFirstOfficerRuntimeForbidsSubagentAcceptanceForStages` scopes to `## Dispatch` and checks the required invariant/rationale phrases.

Commands:
- PASS: `go test ./skills/integration -count=1`
- PASS: `go test ./... -count=1`
- PASS: `gofmt -l skills/integration/skill_surface_test.go` produced no output.

Residual risks:
- Instruction/test-surface validation only; no live Pi smokes were run per task constraints.
