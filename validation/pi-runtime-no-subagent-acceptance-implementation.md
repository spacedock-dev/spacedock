# Pi runtime no-subagent-acceptance implementation

Product commit: `7119e74d83b8d7077a8824d03399fa3928171b5e`
State commit: `c82079eb72e2f98711c122f3e379ef03f386d139`

Changed files:
- `skills/first-officer/references/pi-first-officer-runtime.md`
- `skills/integration/skill_surface_test.go`
- `docs/dev/.spacedock-state/pi-runtime-frontdoor/index.md` (state report)

Tests run from isolated worktree:
- PASS: `go test ./skills/integration -count=1`
- PASS: `go test ./... -count=1`
- Not needed: `gofmt -w ./cmd ./internal` because no `cmd/` or `internal/` Go files changed. `gofmt` was run on `skills/integration/skill_surface_test.go`.

Residual risks:
- This is an instruction-surface invariant; independent validation should still review the contract text for intended semantics.
