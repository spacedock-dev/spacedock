# Roborev evidence — JC validation

- Job: `768`
- Scope: branch `spacedock-ensign/simplify-gate-state-v1-schema` against `origin/main`, candidate `f566f821b`.
- Panel: `branch_final` (`codex` correctness and `claude-code` product, Codex synthesis).
- Result: Changes requested.

## High

1. `internal/gates/model.go:16`: removing `gates.current` and `briefing.digest-domain` in place makes existing v1 documents unreadable after upgrade, blocking status, withdraw, record, and consume. Roborev requests atomic migration or a new schema version with an explicit migration path.
2. `internal/gates/delivery.go:180`: terminal delivery selects the target status, while the pending terminal application belongs to the source gate stage. Roborev requests returning the unique pending terminal application and validating its target against terminal status.

## Medium

1. `internal/gates/model.go:248`: blank record stages pass validation; reject empty stages before uniqueness tracking.
2. `internal/gates/model.go:374-395`, `internal/status/discover.go:218`, `internal/gates/io.go:69-73`: post-gate summaries can become blank when no status-matched record exists; define and test the projection rule.
3. `internal/gates/prepare_test.go:386-408`: the stale-selection test no longer asserts the selected gate or unrelated-record immutability.
4. `internal/ensigncycle/shared_round_recording_test.go:141`: rejection-flow coverage checks count and stage but not gate identity.

## Low

1. `internal/gates/model.go:374-395`, `internal/gates/gates_test.go:30`, `internal/gates/prepare_test.go:179`: make `CurrentSummary` stage semantics explicit.
2. `internal/gates/model.go:175`, `internal/gates/prepare.go:389`: replace error-text branching with a sentinel error.
3. `TestDigestDomainFieldFailsClosed`: assert the error names `digest-domain`.
4. `TestPrototypeAndUnknownGateShapesFailClosed`: rename the stale global-pointer case.
5. `recorded_gate_lifecycle_test.go:109`: anchor the stage-count assertion to record indentation.
6. `internal/cli/gate_test.go:1153`: remove the dead string-concatenation fragment.
