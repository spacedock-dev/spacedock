package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeNextPlugin writes a plugin.json carrying version to a temp file and
// returns its path — the `<next-plugin.json>` the decision subcommand reads (the
// workflow pipes `git show origin/next:.claude-plugin/plugin.json` into it).
func writeNextPlugin(t *testing.T, version string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "plugin.json")
	body := "{\n  \"name\": \"spacedock\",\n  \"version\": \"" + version + "\"\n}\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write next plugin.json: %v", err)
	}
	return path
}

// TestEdgeAdvanceDecisionCommandPrintsVerdict — the subcommand's contract IS its
// stdout: release.yml's decision step captures `advance`/`skip` and writes
// advance=true|false to $GITHUB_OUTPUT. A forward stable cut advances; an
// old-line patch whose target edge version is not strictly greater than next's
// manifest skips. Exit is 0 for BOTH verdicts (a skip is a normal outcome).
func TestEdgeAdvanceDecisionCommandPrintsVerdict(t *testing.T) {
	cases := []struct {
		name string
		tag  string
		next string
		want string
	}{
		{"forward stable advances", "v0.26.0", "0.26.0-pre1", "advance"},
		{"old-line patch skips", "v0.25.1", "0.27.0-pre1", "skip"},
		{"patch-equality skips", "v0.25.1", "0.26.0-pre1", "skip"},
		{"pre0 vs next pre1 skips", "v0.26.0-pre0", "0.26.0-pre1", "skip"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			plugin := writeNextPlugin(t, c.next)
			out, code := captureStdout(t, func() int { return edgeAdvanceDecision([]string{c.tag, plugin}) })
			if code != 0 {
				t.Fatalf("edge-advance-decision exit = %d, want 0", code)
			}
			if got := strings.TrimSpace(out); got != c.want {
				t.Fatalf("edge-advance-decision(%q, next=%q) stdout = %q, want %q", c.tag, c.next, got, c.want)
			}
		})
	}
}

// TestEdgeAdvanceDecisionCommandRejectsBadArgs — a missing manifest path or a
// wrong arg count is a usage/IO error (non-zero), so a release.yml miswiring
// fails loud rather than silently defaulting to skip (which would leave the edge
// channel broken with no signal).
func TestEdgeAdvanceDecisionCommandRejectsBadArgs(t *testing.T) {
	if code := edgeAdvanceDecision([]string{"v0.26.0"}); code == 0 {
		t.Fatalf("edge-advance-decision exit = 0 with one argument; want non-zero")
	}
	if code := edgeAdvanceDecision([]string{"v0.26.0", filepath.Join(t.TempDir(), "missing.json")}); code == 0 {
		t.Fatalf("edge-advance-decision exit = 0 on a missing manifest; want non-zero")
	}
}

// TestHighestKnownEdgeVersionCommandPrintsHighest — the subcommand release.yml's
// decision step shells out to: given a candidate tag list, prints the greatest
// (prereleases included, per release.HighestKnownEdgeVersion's doc), always
// exit 0.
func TestHighestKnownEdgeVersionCommandPrintsHighest(t *testing.T) {
	out, code := captureStdout(t, func() int {
		return highestKnownEdgeVersion([]string{"v0.26.0", "v0.27.0-pre7", "v0.25.1"})
	})
	if code != 0 {
		t.Fatalf("highest-known-edge-version exit = %d, want 0", code)
	}
	if got := strings.TrimSpace(out); got != "0.27.0-pre7" {
		t.Fatalf("highest-known-edge-version stdout = %q, want 0.27.0-pre7", got)
	}
}

// TestHighestKnownEdgeVersionCommandPrintsNothingWhenNoneParses covers the
// round-3 fail-closed contract: an empty argument list, and a list of only
// malformed entries, must both print NOTHING (empty stdout) with exit 0 — not
// an error — so release.yml's shell `[ -z "$HIGHEST_KNOWN" ]` check can route
// to "skip the auto-pre0 cut" without the job itself failing.
func TestHighestKnownEdgeVersionCommandPrintsNothingWhenNoneParses(t *testing.T) {
	for _, args := range [][]string{nil, {}, {"not-a-version", "garbage"}} {
		out, code := captureStdout(t, func() int { return highestKnownEdgeVersion(args) })
		if code != 0 {
			t.Fatalf("highest-known-edge-version(%v) exit = %d, want 0", args, code)
		}
		if strings.TrimSpace(out) != "" {
			t.Fatalf("highest-known-edge-version(%v) stdout = %q, want empty", args, out)
		}
	}
}

// TestEdgePre0VersionCommandPrintsComputedVersion — the subcommand prints exactly
// X.(Y+1).0-pre0, so release.yml's always-cut-pre0 step forms `vX.(Y+1).0-pre0`
// as the annotated auto-tag, whose minor matches the required minor next's skills
// were stamped to.
func TestEdgePre0VersionCommandPrintsComputedVersion(t *testing.T) {
	out, code := captureStdout(t, func() int { return edgePre0Version([]string{"0.25.0"}) })
	if code != 0 {
		t.Fatalf("edge-pre0-version exit = %d, want 0 on a bare stable semver", code)
	}
	if got := strings.TrimSpace(out); got != "0.26.0-pre0" {
		t.Fatalf("edge-pre0-version stdout = %q, want 0.26.0-pre0", got)
	}
	if code := edgePre0Version([]string{"0.25.0-pre1"}); code == 0 {
		t.Fatalf("edge-pre0-version exit = 0 on a hyphenated input; want non-zero")
	}
}
