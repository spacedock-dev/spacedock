// ABOUTME: resolveTrunk unit (AC-2) + dispatch trunk command (AC-5) — the single
// ABOUTME: trunk-config resolver and its byte-exact stdout command surface.
package dispatch

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
)

// readmeWithTrunk renders a minimal workflow README frontmatter declaring a
// top-level trunk: key (sibling of state:), plus a stages block so the file
// resembles a real workflow README.
func readmeWithTrunk(trunk string) string {
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("entity-type: task\n")
	b.WriteString("state: .spacedock-state\n")
	if trunk != "" {
		b.WriteString("trunk: ")
		b.WriteString(trunk)
		b.WriteString("\n")
	}
	b.WriteString("stages:\n  defaults:\n    worktree: false\n---\n")
	return b.String()
}

// TestResolveTrunk (AC-2) asserts resolveTrunk returns the declared top-level
// trunk for a README with trunk: set, and the main fallback for a README with
// no trunk key. Expected values come from the fixture READMEs, not the
// function's own source.
func TestResolveTrunk(t *testing.T) {
	cases := []struct {
		name   string
		trunk  string // "" → omit the key entirely
		expect string
	}{
		{"declared-sentinel", "ftrunk", "ftrunk"},
		{"declared-main", "main", "main"},
		{"absent-key-falls-back-to-main", "", "main"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			writeFile(t, filepath.Join(dir, "README.md"), readmeWithTrunk(tc.trunk))
			got := resolveTrunk(dir)
			if got != tc.expect {
				t.Errorf("resolveTrunk(%q trunk=%q) = %q, want %q", dir, tc.trunk, got, tc.expect)
			}
		})
	}
}

// TestResolveTrunkMissingReadme asserts resolveTrunk degrades to the main
// fallback when the README is unreadable/absent — the resolver never panics
// and never returns next.
func TestResolveTrunkMissingReadme(t *testing.T) {
	dir := t.TempDir() // no README written
	if got := resolveTrunk(dir); got != "main" {
		t.Errorf("resolveTrunk(no README) = %q, want main", got)
	}
}

// TestDispatchTrunkCommand (AC-5) drives the `dispatch trunk` subcommand and
// asserts its stdout is BYTE-EXACT — the sole load-bearing contract for the
// $(...)-capturing prose consumers. A declared trunk: ftrunk yields exactly
// "ftrunk\n"; an absent key yields exactly "main\n". Any stray stdout line or a
// missing/extra newline poisons BASE=$(spacedock dispatch trunk ...).
func TestDispatchTrunkCommand(t *testing.T) {
	cases := []struct {
		name       string
		trunk      string
		wantStdout string
	}{
		{"declared-sentinel", "ftrunk", "ftrunk\n"},
		{"absent-key", "", "main\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			writeFile(t, filepath.Join(dir, "README.md"), readmeWithTrunk(tc.trunk))

			var out, errBuf bytes.Buffer
			code := Run(claudeteam.Probe, []string{"trunk", "--workflow-dir", dir},
				strings.NewReader(""), &out, &errBuf)
			if code != 0 {
				t.Fatalf("dispatch trunk exit=%d stderr=%q", code, errBuf.String())
			}
			if out.String() != tc.wantStdout {
				t.Errorf("dispatch trunk stdout = %q, want byte-exact %q", out.String(), tc.wantStdout)
			}
			// Any diagnostics go to stderr, never stdout — a stray stdout line
			// poisons the $(...) capture. stderr is allowed to be empty here.
			if errBuf.Len() != 0 {
				t.Errorf("dispatch trunk wrote to stderr on success: %q (must be silent)", errBuf.String())
			}
		})
	}
}

// TestDispatchTrunkCommandMissingFlag asserts the command surfaces a usage
// error (exit 2) when --workflow-dir is absent, matching the other dispatch
// subcommands' required-flag behavior.
func TestDispatchTrunkCommandMissingFlag(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := Run(claudeteam.Probe, []string{"trunk"}, strings.NewReader(""), &out, &errBuf)
	if code != 2 {
		t.Errorf("dispatch trunk with no flags: exit=%d, want 2", code)
	}
	if out.Len() != 0 {
		t.Errorf("usage error wrote to stdout: %q (must be stderr-only)", out.String())
	}
	if !strings.Contains(errBuf.String(), "trunk") {
		t.Errorf("usage error does not name the subcommand: %q", errBuf.String())
	}
}
