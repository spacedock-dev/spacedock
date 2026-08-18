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

// TestHighestKnownEdgeVersionCommandRestoresNextNotchAgainstRealV0251Collision
// replays this repository's REAL v0.25.1 cut (2026-07-20) — the case Finding 1
// of the validation-cycle-1 review named by tag creation date: v0.26.0-pre0 was
// auto-cut 2026-07-15 (three minutes after v0.25.0, the latest-line stable that
// triggered it) and stayed the newest 0.26-line tag until v0.26.0 itself landed
// the next day after v0.25.1. Candidates below are every real tag that existed
// at v0.25.1's cut, taken verbatim from `git tag --list --format='%(refname:short)
// %(creatordate:iso-strict)'`. Before the notch fix, highest-known-edge-version
// printed the raw scan's max (0.26.0-pre0) and edge-advance-decision then said
// "advance" for v0.25.1 — because DevPreVersion(0.25.1) is ALWAYS "0.26.0-pre1"
// (a patch's own line doesn't affect the next-line target), which outranks a
// bare "-pre0" by construction; the auto-cut step would then die re-cutting the
// already-existing v0.26.0-pre0 tag. With the notch (dev-preversion of the
// highest bare stable, 0.25.0 -> "0.26.0-pre1"), the known version is pulled up
// to equal the patch's own target, and the strict `>` boundary skips — exactly
// what the retired next-manifest-based mechanism did (next was stamped to
// 0.26.0-pre1 by the same 0.25.0 cut that auto-cut the pre0 tag).
func TestHighestKnownEdgeVersionCommandRestoresNextNotchAgainstRealV0251Collision(t *testing.T) {
	candidates := []string{
		"v0.24.0-pre1", "v0.24.0-pre2", "v0.24.0",
		"v0.25.0-pre1", "v0.25.0-pre2", "v0.25.0",
		"v0.26.0-pre0",
	}
	const patchTag = "v0.25.1"

	knownOut, code := captureStdout(t, func() int { return highestKnownEdgeVersion(candidates) })
	if code != 0 {
		t.Fatalf("highest-known-edge-version exit = %d, want 0", code)
	}
	known := strings.TrimSpace(knownOut)
	if known != "0.26.0-pre1" {
		t.Fatalf("highest-known-edge-version(%v) = %q, want 0.26.0-pre1 (the restored notch — dev-preversion of the highest bare stable, 0.25.0); the raw scan alone would wrongly stop at 0.26.0-pre0", candidates, known)
	}

	plugin := writeNextPlugin(t, known)
	verdictOut, code := captureStdout(t, func() int { return edgeAdvanceDecision([]string{patchTag, plugin}) })
	if code != 0 {
		t.Fatalf("edge-advance-decision exit = %d, want 0", code)
	}
	if verdict := strings.TrimSpace(verdictOut); verdict != "skip" {
		t.Fatalf("edge-advance-decision(%q, known=%q) = %q, want skip — an old-line patch must not out-rank a newer line that already has a pre0 tag, or the auto-cut step dies re-cutting v0.26.0-pre0", patchTag, known, verdict)
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
