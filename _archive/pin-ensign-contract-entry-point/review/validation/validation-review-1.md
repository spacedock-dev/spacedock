# Validation gate review — pin-ensign-contract-entry-point (mxaa)

- Stage: validation (attempt 1), candidate `694d21db7`
- Reviewer: agent:first-officer (Pi host)
- Verdict: **REJECTED** — on AC-1 gate evidence (three material evidence defects in the committed live harness); revise routes to implementation.

## Verdict support

- **AC-2/AC-3/AC-4 PASS, independently reproduced.** All pinned tests green on candidate; red-first replayed on throwaway `af42c77d5` + new tests applied (AC-2 nil fields; AC-3 injected-under-marker; AC-4 missing binding sentence; wrapper round-trip red as build error). Full `go test ./...` and `-race` green with markers scrubbed; gofmt clean.
- **AC-1 FAIL as graded harness (evidence defect, 3 findings, Material/owned-by-task):**
  1. `pi_live_runner_test.go:181` runs `spacedock build` — CLI exits 2 `unknown command: build`.
  2. `cmd.CombinedOutput()` merges the bare-mode stderr advisory into the JSON envelope parse (`invalid character 'W'…`).
  3. Hermetic `piHome` lacks `settings.json` package registration → basename `ensign` unresolvable in the nested child (live-observed `skills: []`).
- **End-state already proven on a 3-patch throwaway:** the same live leg then PASSED in 124.7s; grader artifact `pi-ensign-boot-grade.json = {ensign_skill_read_rank: 1, first_officer_reads: 0, spawn_agent: "worker", spawn_skills: ["ensign"], verdict: "pass"}` — the product behavior is correct; the graded harness cannot attest it as committed.
- **Surface necessity upheld** (591+/24−): ~200 lines = ruling-attributed smoke/grader; ~25 closes a real hermeticity hole this session itself relies on; ~40 collapsible fixture-prelude polish, non-blocking.
- **Deferred risk (accepted residual):** `piLiveEnv` doesn't scrub ambient `PI_SUBAGENT_*` markers — bites only when the lane runs nested under pi-subagents; CI lane env is clean; promote condition stated.
- **Pre-existing, unrelated:** ambient-`PI_CODING_AGENT` red in `TestVersionAmbiguousMarkersExitZero` (reproduced on base too; my session hit the same mechanism this morning). Recommended as a separate small filing if the captain wants it.

## FO disposition (per workflow Review-finding policy)

Findings 1–3: Material × evidence-kind, owned by this task's diff (all in code the implementation added) → eligible for one FO-authorized narrow fix cycle, assignment: apply the three validator-proven patches (argv `dispatch`; stdout/stderr split, parse stdout only; seed `piHome/settings.json` with the repo as a path-package), then re-run the tagged live leg **on the candidate worktree** to produce the gate-grade boot artifact. Optional-if-adjacent: extend `piLiveEnv` with the `PI_SUBAGENT_*` scrub (same hole class as the landed harnessEnv).

## Ask

REJECTED → revise, routed to a fresh implementation worker with the assignment above and zero candidate-design changes.
