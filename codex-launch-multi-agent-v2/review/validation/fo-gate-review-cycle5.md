# Codex launcher guarantees multi-agent v2 — validation cycle 5

## Chosen direction

Keep launcher-owned collaboration configuration as the sole source of truth and remove the obsolete forwarded `--enable multi_agent_v2` pair from the shared live runner.

## Evidence

- Validation cycle 5 covers AC-1 through AC-4 at exact commit `60a75ba658`.
- The isolated-home lifecycle passed in 39.91s; the shared live lane passed seven enabled journeys in 792.62s, with three existing TODO journeys explicitly skipped.
- `go test ./...`, `go test ./... -race`, formatting, diff, and detached adversarial checks passed.
- The candidate is exactly 9 files, +380/−42 against `main`; cycle 5 changes only the runner test and runtime-live documentation.

## Recommendation

Present this fresh validation Briefing for Captain approval to advance the candidate to done and the merge boundary. The three TODO skips are not treated as capability evidence.
