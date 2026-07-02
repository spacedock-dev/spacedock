package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// fakeRunLister returns canned `gh run list` JSON without shelling out, so the
// e2e-gate subcommand's exit-code + step-summary behavior is exercised against
// fixtures rather than a live gh.
func fakeRunLister(json string) runListFunc {
	return func(commit string) ([]byte, error) {
		return []byte(json), nil
	}
}

const greenForCommit = `[{"databaseId": 42, "headSha": "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef", "conclusion": "success", "status": "completed"}]`
const parkedForCommit = `[{"databaseId": 43, "headSha": "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef", "conclusion": "", "status": "waiting"}]`

const gateCommit = "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"

// TestE2EGateCommandPassesOnGreenRun — the subcommand exits 0 and records the
// matched run in $GITHUB_STEP_SUMMARY when a green run matches the commit.
func TestE2EGateCommandPassesOnGreenRun(t *testing.T) {
	summary := newSummaryFile(t)
	t.Setenv("GITHUB_STEP_SUMMARY", summary)
	t.Setenv("SPACEDOCK_E2E_GATE_WAIVER", "")

	code := runE2EGate([]string{gateCommit}, fakeRunLister(greenForCommit))
	if code != 0 {
		t.Fatalf("e2e-gate exit = %d, want 0 on a green run", code)
	}
	if got := readFile(t, summary); !strings.Contains(got, gateCommit) {
		t.Errorf("step summary did not record the matched commit:\n%s", got)
	}
}

// TestE2EGateCommandBlocksOnParkedRun — the subcommand exits 1 on a parked run
// (the cut is blocked because the live lanes did not run green).
func TestE2EGateCommandBlocksOnParkedRun(t *testing.T) {
	summary := newSummaryFile(t)
	t.Setenv("GITHUB_STEP_SUMMARY", summary)
	t.Setenv("SPACEDOCK_E2E_GATE_WAIVER", "")

	code := runE2EGate([]string{gateCommit}, fakeRunLister(parkedForCommit))
	if code == 0 {
		t.Fatalf("e2e-gate exit = 0 on a parked run; want non-zero (cut blocked)")
	}
}

// TestE2EGateCommandPassesWhenWaived — with SPACEDOCK_E2E_GATE_WAIVER set the
// subcommand exits 0 even on a parked run, and the waiver reason is written to
// the step summary for the audit trail.
func TestE2EGateCommandPassesWhenWaived(t *testing.T) {
	summary := newSummaryFile(t)
	t.Setenv("GITHUB_STEP_SUMMARY", summary)
	const reason = "emergency cut approved by captain clkao"
	t.Setenv("SPACEDOCK_E2E_GATE_WAIVER", reason)

	code := runE2EGate([]string{gateCommit}, fakeRunLister(parkedForCommit))
	if code != 0 {
		t.Fatalf("e2e-gate exit = %d on a waived gate; want 0", code)
	}
	got := readFile(t, summary)
	if !strings.Contains(got, reason) {
		t.Errorf("step summary did not record the waiver reason for the audit trail:\n%s", got)
	}
	if !strings.Contains(strings.ToUpper(got), "WAIV") {
		t.Errorf("step summary did not mark the cut as waived:\n%s", got)
	}
}

// TestE2EGateCommandRejectsMissingCommit — the subcommand needs the release
// commit argument; with none it exits with a usage error and does not pass.
func TestE2EGateCommandRejectsMissingCommit(t *testing.T) {
	t.Setenv("SPACEDOCK_E2E_GATE_WAIVER", "")
	code := runE2EGate(nil, fakeRunLister(greenForCommit))
	if code == 0 {
		t.Fatalf("e2e-gate exit = 0 with no release commit argument; want non-zero")
	}
}

// writeGhStatusLagShim writes a `gh` shim that reproduces the v0.23.0 gate-time
// filter-consistency lag: when argv carries a bare `--status` flag it prints
// `[]` (the observed behavior of `--status success` right after a re-run flips
// a run to green), and otherwise it prints fixtureJSON (the unfiltered query,
// which already saw the success). This drives the REAL ghRunListForCommit
// end-to-end — the fakeRunLister seam above bypasses the args entirely and
// cannot see the bug.
//
// It also pins the query's ONLY binding to the live matrix: `--workflow
// "Runtime Live E2E"` is not among the fetched `--json` fields, so nothing else
// in the predicate distinguishes a Runtime Live E2E run from an unrelated green
// run on the release commit (e.g. ordinary push CI). The shim prints `[]`
// unless argv carries that exact `--workflow "Runtime Live E2E"` pair, so
// dropping the token from ghRunListForCommit fails the gate closed here — a
// regression that dropped it would otherwise fail OPEN (any green run on the
// commit would satisfy the gate).
func writeGhStatusLagShim(t *testing.T, fixtureJSON string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("gh stub shim is a POSIX shell script")
	}
	dir := t.TempDir()
	// echo, not cat: the stub PATH contains only this shim, so an external `cat`
	// would itself fail to resolve. echo is a shell builtin.
	script := `#!/bin/sh
have_workflow=0
prev=""
for arg in "$@"; do
  if [ "$prev" = "--workflow" ] && [ "$arg" = "Runtime Live E2E" ]; then
    have_workflow=1
  fi
  if [ "$arg" = "--status" ]; then
    echo '[]'
    exit 0
  fi
  prev="$arg"
done
if [ "$have_workflow" != "1" ]; then
  echo '[]'
  exit 0
fi
echo '` + fixtureJSON + `'
`
	path := filepath.Join(dir, "gh")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

// rerunToGreenJSON pins the 2026-06-30 gate-time observation: the only success
// entry for the release commit is a run that was re-run to green (attempt 2).
const rerunToGreenJSON = `[
  {"databaseId": 28429490220, "headSha": "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef", "conclusion": "success", "status": "completed", "attempt": 2}
]`

// parkedOnlyJSON has no success entry at all — only a parked/waiting run.
const parkedOnlyJSON = `[
  {"databaseId": 28400000000, "headSha": "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef", "conclusion": "", "status": "waiting"}
]`

// TestE2EGateCommandPassesOnRerunToGreenLiveRun is T1 (AC-1): drives the REAL
// ghRunListForCommit against the gh shim. RED under today's query shape (argv
// carries `--status success`, so the shim replays the filter-consistency lag
// and returns `[]`, blocking a genuinely green cut); GREEN once
// ghRunListForCommit drops `--status success` (the shim then sees no bare
// `--status` arg and returns the fixture, so the gate matches the run).
func TestE2EGateCommandPassesOnRerunToGreenLiveRun(t *testing.T) {
	stubDir := writeGhStatusLagShim(t, rerunToGreenJSON)
	t.Setenv("PATH", stubDir)

	summary := newSummaryFile(t)
	t.Setenv("GITHUB_STEP_SUMMARY", summary)
	t.Setenv("SPACEDOCK_E2E_GATE_WAIVER", "")

	code := runE2EGate([]string{gateCommit}, ghRunListForCommit)
	if code != 0 {
		t.Fatalf("e2e-gate exit = %d, want 0 for a re-run-to-green Runtime Live E2E run (AC-1)", code)
	}
	if got := readFile(t, summary); !strings.Contains(got, "28429490220") {
		t.Errorf("step summary did not cite the matched re-run-to-green run:\n%s", got)
	}
}

// TestE2EGateCommandBlocksParkedOnlyRunViaLiveQuery is T2 (AC-2): a
// parked-only fixture still blocks the cut under the new (unfiltered) query
// shape. The parked-run defense lives in EvaluateE2EGate's predicate (empty
// conclusion never matches), not in the dropped `--status` flag, so removing
// the flag must not regress this.
func TestE2EGateCommandBlocksParkedOnlyRunViaLiveQuery(t *testing.T) {
	stubDir := writeGhStatusLagShim(t, parkedOnlyJSON)
	t.Setenv("PATH", stubDir)

	summary := newSummaryFile(t)
	t.Setenv("GITHUB_STEP_SUMMARY", summary)
	t.Setenv("SPACEDOCK_E2E_GATE_WAIVER", "")

	code := runE2EGate([]string{gateCommit}, ghRunListForCommit)
	if code == 0 {
		t.Fatalf("e2e-gate exit = 0 on a parked-only run list; want non-zero (cut blocked, AC-2)")
	}
}

func newSummaryFile(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "step_summary")
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
