package main

import (
	"os"
	"path/filepath"
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
