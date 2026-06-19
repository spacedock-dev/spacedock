// ABOUTME: TeamCreate-capability discriminator — parses `claude --version` so the legacy
// ABOUTME: interactive pty lane SKIPs on a merged host (native team tools gone, ≥ the merged floor).
package ensigncycle

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// mergedFloorMinor is the claude-code minor version at which the native team
// tools (TeamCreate/TeamDelete) were dropped — the "merged floor". 2.1.178
// retired the native team registry in favour of the in-process named-background
// teammate shape (anthropics/claude-code#68721), so 2.1.<minor≥178> is a MERGED
// host where the legacy interactive pty lane cannot run TeamCreate. A host below
// the floor (the pinned 2.1.177) still exposes the native tools — the legacy lane
// runs there.
const mergedFloorMinor = 178

// claudeVersionPattern extracts the leading semver from `claude --version` output
// (e.g. "2.1.181 (Claude Code)" → 2,1,181). The trailing " (Claude Code)" label
// and any build suffix are ignored; only the dotted numeric head is read.
var claudeVersionPattern = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)`)

// mergedFloor is the (major, minor, patch) at which the native team tools were
// dropped. A version BELOW this floor still exposes TeamCreate/TeamDelete; at or
// above it, the host is merged. Kept as the one floor literal so teamCreateCapable
// is a plain lexicographic-by-component compare against it.
var mergedFloor = [3]int{2, 1, mergedFloorMinor}

// teamCreateCapable reports whether the claude version in versionOutput still
// exposes the native TeamCreate/TeamDelete tools — i.e. it is strictly BELOW the
// merged floor (2.1.178). This is the harness-observable analog of the contract's
// `ToolSearch(select:TeamCreate)`-empty discriminator: a Go test cannot run the
// ToolSearch hop, but it can read `claude --version`, and the team-tool presence
// is a deterministic function of the version (present < 2.1.178, gone ≥ 2.1.178).
//
// An unparseable version conservatively reports NOT capable (false): the legacy
// lane SKIPs rather than FAILs when the version is unknown, matching the
// skip-not-fatal discipline (a merged host or an unreadable probe both skip).
func teamCreateCapable(versionOutput string) bool {
	major, minor, patch, ok := parseClaudeVersion(versionOutput)
	if !ok {
		return false
	}
	got := [3]int{major, minor, patch}
	for i := range got {
		if got[i] != mergedFloor[i] {
			return got[i] < mergedFloor[i]
		}
	}
	// Exactly the floor (2.1.178) is the first merged version — not capable.
	return false
}

// parseClaudeVersion pulls the (major, minor, patch) ints out of a `claude
// --version` line. ok is false when no leading dotted-numeric version is present.
func parseClaudeVersion(versionOutput string) (major, minor, patch int, ok bool) {
	m := claudeVersionPattern.FindStringSubmatch(strings.TrimSpace(versionOutput))
	if m == nil {
		return 0, 0, 0, false
	}
	major, err1 := strconv.Atoi(m[1])
	minor, err2 := strconv.Atoi(m[2])
	patch, err3 := strconv.Atoi(m[3])
	if err1 != nil || err2 != nil || err3 != nil {
		return 0, 0, 0, false
	}
	return major, minor, patch, true
}

// teamCapabilitySkipReason returns the SKIP message the legacy pty lane uses when
// the host is merged (no native TeamCreate). It names the observed version so a CI
// log makes the skip self-explaining (the version pin moved past the floor) rather
// than an opaque skip.
func teamCapabilitySkipReason(versionOutput string) string {
	return fmt.Sprintf(
		"legacy interactive team-mode lane SKIPPED: the live claude (%q) is at/above the merged floor (2.1.%d) where native TeamCreate/TeamDelete are gone — this lane needs the native team tools. On a merged host the FO dispatches in-process named background teammates instead (covered by the merged lane); the legacy lane retires when stable Claude catches up.",
		strings.TrimSpace(versionOutput), mergedFloorMinor)
}

// TestTeamCreateCapable is the offline proof of the legacy-lane skip discriminator
// (no model spend, no live claude): teamCreateCapable reads `claude --version` and
// reports whether the native team tools are present, so the legacy interactive pty
// tests SKIP cleanly when CI unpins to a merged host (≥2.1.178) rather than RED.
func TestTeamCreateCapable(t *testing.T) {
	cases := []struct {
		name    string
		version string
		capable bool
	}{
		{"pinned legacy 2.1.177 is team-capable", "2.1.177 (Claude Code)", true},
		{"merged floor 2.1.178 is NOT capable", "2.1.178 (Claude Code)", false},
		{"current local 2.1.181 is NOT capable", "2.1.181 (Claude Code)", false},
		{"older 2.1.161 is team-capable", "2.1.161 (Claude Code)", true},
		{"older minor 2.0.250 is team-capable", "2.0.250 (Claude Code)", true},
		{"future minor 2.2.0 is NOT capable", "2.2.0 (Claude Code)", false},
		{"future major 3.0.0 is NOT capable", "3.0.0 (Claude Code)", false},
		{"older major 1.9.9 is team-capable", "1.9.9 (Claude Code)", true},
		{"bare version with no label still parses", "2.1.177", true},
		{"leading whitespace tolerated", "  2.1.177 (Claude Code)\n", true},
		{"unparseable version is conservatively NOT capable", "garbage", false},
		{"empty output is conservatively NOT capable", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := teamCreateCapable(tc.version); got != tc.capable {
				t.Errorf("teamCreateCapable(%q) = %v, want %v", tc.version, got, tc.capable)
			}
		})
	}

	// The skip reason names the observed version so the CI log is self-explaining.
	if r := teamCapabilitySkipReason("2.1.181 (Claude Code)"); !strings.Contains(r, "2.1.181") || !strings.Contains(r, "SKIPPED") {
		t.Errorf("teamCapabilitySkipReason should name the version and SKIP; got %q", r)
	}
}
