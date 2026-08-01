# Validation gate review (attempt 2) — pin-ensign-contract-entry-point (mxaa)

- Stage: validation, attempt 2; candidate `d9f38ca1e` (implementation `694d21db7` + test-only correction round)
- Reviewer: agent:first-officer (Pi host)
- Verdict: **PASS** — validation reproduced every AC by exercise; recommend approve → terminal delivery ceremony.

## Evidence (all this-run, all exercised)

1. **AC-1 (value):** tagged live leg re-run on the candidate worktree at `d9f38ca1e` on this host (real pi + pi-subagents child): PASS 172.5s. Graded artifact `pi-ensign-boot-grade.json`: `ensign_skill_read_rank: 1`, `first_officer_reads: 0`, `spawn_agent: "worker"`, `spawn_skills: ["ensign"]`, `verdict: "pass"`. Failure modes are falsifiable: a stale binary, an ensign-contract-missing child, or any FO bootstrap in the child each fail this test. Baseline was 0/8; the gate-grade run is 1/1.
2. **Correction discipline:** `694d21db7..d9f38ca1e` = exactly `57 5` in `internal/ensigncycle/pi_live_runner_test.go` — the three findings fixed (argv `dispatch build`, stdout/stderr split with stdout-only parse, `piHome/settings.json` seeded), plus `TestPiLiveEnvScrubsAmbientPiSubagentMarkers` covering the accepted deferred risk (promote condition hit this session).
3. **AC-2/3/4:** every pin PASS with failure modes named; claude goldens byte-identical; full `go test ./...` and `-race` both 19 ok, zero FAIL; gofmt clean. Red-first replays (attempt-1 report) stand on base `af42c77d5`.

## Semantic adversarial pass

Binary identity (graded run exercises the candidate, not the stale ambient SPACEDOCK_BIN), authority (zero first-officer reads binds the extension exemption end-to-end under a real child spawn), cardinality (`worker_transcripts_graded: 1` — no vacuous pass).

## Carried notes

- Deferred risk: ambient-`PI_CODING_AGENT` version-test hermeticity, red on base too, outside this task's surface (candidate small filing if wanted).
- Operational hazard recorded: replays must unset ambient `SPACEDOCK_BIN` or point to a fresh candidate (mirrors CI convention).

## Fan-out accountability (session)

ideation 1 + implementation 1 + correction 1 + validators 2 = 5 workers against the declared 2 ±1 — the correction round and re-validation were the gate-required tolerance class; no unplanned spawns.
