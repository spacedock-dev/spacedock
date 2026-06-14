// ABOUTME: AC-4 — the integration-trunk axis and the marketplace channel stamp
// ABOUTME: stay un-conflated: trunk resolution never reads or mutates devBranch.
package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
	"github.com/spacedock-dev/spacedock/internal/dispatch"
)

// TestTrunkNotConflatedWithChannel (AC-4) is the un-conflation invariant over the
// two already-separate axes. The integration trunk (here declared trunk: main in
// the workflow README) and the marketplace channel stamp (cli.devBranch, default
// "next") are different axes; next is CORRECT for the edge channel. This test
// drives the real `dispatch trunk` resolution path on a trunk: main workflow and
// asserts:
//   - the command resolves the trunk to "main" (from the README), AND
//   - cli.devBranch is unchanged at its "next" default afterward.
//
// The two expected values come from two independent declared sources (the README's
// trunk: main and devBranch's default), not from one another. A regression that
// routed the trunk through devBranch (or vice-versa) would couple them and trip
// this invariant.
func TestTrunkNotConflatedWithChannel(t *testing.T) {
	// devBranch is the package channel stamp; capture it as the independent axis.
	saved := devBranch
	defer func() { devBranch = saved }()
	if devBranch != "next" {
		t.Fatalf("precondition: devBranch=%q, want default \"next\" (the edge-channel stamp)", devBranch)
	}

	dir := t.TempDir()
	writeReadme(t, dir, "main")

	var out, errBuf bytes.Buffer
	code := dispatch.Run(claudeteam.Probe, []string{"trunk", "--workflow-dir", dir},
		strings.NewReader(""), &out, &errBuf)
	if code != 0 {
		t.Fatalf("dispatch trunk exit=%d stderr=%q", code, errBuf.String())
	}

	if got := out.String(); got != "main\n" {
		t.Errorf("trunk resolution = %q, want \"main\\n\" (from the README's trunk: main)", got)
	}
	// The load-bearing invariant: resolving the trunk left the channel stamp
	// untouched — the two axes are not conflated.
	if devBranch != "next" {
		t.Errorf("trunk resolution mutated devBranch to %q; the channel stamp must stay \"next\" (un-conflated)", devBranch)
	}
}

// writeReadme writes a minimal workflow README declaring the given top-level
// trunk: key under dir.
func writeReadme(t *testing.T, dir, trunk string) {
	t.Helper()
	body := "---\nentity-type: task\nstate: .spacedock-state\ntrunk: " + trunk + "\nstages:\n  defaults:\n    worktree: false\n---\n"
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
