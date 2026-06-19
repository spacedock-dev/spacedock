// ABOUTME: Real-git e2e — `state sweep` is read-only (HEAD unchanged, AC-7), and the
// ABOUTME: state verbs' usage/routing exit codes (AC-8).
package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// TestStateSweepIsReadOnly pins AC-7's read-only property at the verb level: `state
// sweep` makes no commit, push, or state mutation — the state checkout's HEAD is
// unchanged across the verb. Entities carry no `pr:`, so the un-advanced-PR scan is
// offline (no gh call) and deterministic.
func TestStateSweepIsReadOnly(t *testing.T) {
	_, workflowA, _, _ := twoHostStateWorkflow(t)
	checkoutA := filepath.Join(workflowA, ".spacedock-state")
	hostA := filepath.Dir(filepath.Dir(workflowA))

	headBefore := strings.TrimSpace(git(t, checkoutA, "rev-parse", "HEAD"))

	var out, errBuf strings.Builder
	code := run(context.Background(), []string{"state", "sweep", "--workflow-dir", workflowA, "--json"},
		os.Environ(), hostA, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("state sweep should exit 0; got exit=%d stderr=%q", code, errBuf.String())
	}
	if !strings.Contains(out.String(), `"command": "state sweep"`) {
		t.Fatalf("sweep --json should carry the command envelope; json:\n%s", out.String())
	}
	headAfter := strings.TrimSpace(git(t, checkoutA, "rev-parse", "HEAD"))
	if headAfter != headBefore {
		t.Fatalf("state sweep must be read-only; HEAD changed %s -> %s", headBefore, headAfter)
	}
	// No mutation: a clean checkout afterward.
	if porcelain := git(t, checkoutA, "status", "--porcelain"); strings.TrimSpace(porcelain) != "" {
		t.Fatalf("state sweep must leave a clean checkout; porcelain:\n%s", porcelain)
	}
}

// TestStateVerbUsageExitCodes pins AC-8: each verb exits 2 on a missing/unknown
// required argument with a diagnostic naming the subcommand, and `state <unknown>`
// exits 2 enumerating init|new|ready|sweep|commit.
func TestStateVerbUsageExitCodes(t *testing.T) {
	cases := []struct {
		name      string
		args      []string
		wantCode  int
		wantInErr string
	}{
		{"commit missing slug", []string{"state", "commit"}, 2, "state commit"},
		{"commit unknown flag", []string{"state", "commit", "slug", "--bogus"}, 2, "state commit"},
		{"commit extra positional", []string{"state", "commit", "one", "two"}, 2, "state commit"},
		{"ready unknown arg", []string{"state", "ready", "--bogus"}, 2, "state ready"},
		{"sweep unknown arg", []string{"state", "sweep", "--bogus"}, 2, "state sweep"},
		{"unknown subcommand", []string{"state", "frobnicate"}, 2, "init|new|ready|sweep|commit"},
		{"missing subcommand", []string{"state"}, 2, "init|new|ready|sweep|commit"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var out, errBuf strings.Builder
			code := run(context.Background(), tc.args, os.Environ(), t.TempDir(), nil, &out, &errBuf, &status.NativeRunner{}, nil)
			if code != tc.wantCode {
				t.Fatalf("args %v: want exit %d, got %d (stderr=%q)", tc.args, tc.wantCode, code, errBuf.String())
			}
			if !strings.Contains(errBuf.String(), tc.wantInErr) {
				t.Fatalf("args %v: stderr should contain %q; got:\n%s", tc.args, tc.wantInErr, errBuf.String())
			}
		})
	}
}
