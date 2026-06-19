//go:build live

// ABOUTME: Live wiring of the TeamCreate-capability gate — probes `claude --version` and
// ABOUTME: SKIPs the legacy interactive pty lane on a merged host (native team tools gone).
package ensigncycle

import (
	"os/exec"
	"strings"
	"testing"
)

// skipUnlessTeamCreateCapable SKIPs the calling legacy pty test when the live
// claude is a merged host (≥2.1.178, native TeamCreate/TeamDelete gone). It probes
// the SAME `claude` binary the FO will launch — `claude --version`, the cheap
// model-free harness-observable signal — and applies teamCreateCapable. When CI
// unpins past the merged floor this makes the legacy lane SKIP cleanly rather than
// RED (the FO has no native team tools to drive). An unresolvable/unprobeable
// claude also SKIPs (teamCreateCapable returns false on empty/garbage), never
// fatals — matching the skip-not-fatal discipline.
func skipUnlessTeamCreateCapable(t *testing.T) {
	t.Helper()
	out, err := exec.Command("claude", "--version").Output()
	version := string(out)
	if err != nil {
		// A claude that won't report its version cannot run the legacy lane either;
		// SKIP (the merged-lane test owns the version-independent merged coverage).
		t.Skipf("legacy interactive team-mode lane SKIPPED: `claude --version` failed (%v) — cannot confirm native TeamCreate is present", err)
	}
	if !teamCreateCapable(version) {
		t.Skip(teamCapabilitySkipReason(version))
	}
}

// skipUnlessMergedHost is the symmetric gate for the merged lane: it SKIPs the
// caller when the live claude is a LEGACY host (below the merged floor, native
// TeamCreate still present). On a legacy host the FO finds TeamCreate via ToolSearch
// and drives the native team registry, NOT the in-process named-background shape the
// merged lane asserts — so the merged test would mis-fire there. The two gates are
// mirror images: the legacy lane runs only below the floor, the merged lane only at/
// above it, so whichever claude CI pins to, exactly one team lane runs and the other
// SKIPs cleanly. An unprobeable claude is treated as merged (teamCreateCapable false)
// — the merged lane runs (skip-not-fatal stays with the auth gate downstream).
func skipUnlessMergedHost(t *testing.T) {
	t.Helper()
	out, err := exec.Command("claude", "--version").Output()
	if err != nil {
		// Cannot read the version → assume merged (the unpinned default is merged);
		// let the lane run and fail loudly if the host turns out to be legacy.
		t.Logf("merged team-mode lane: `claude --version` failed (%v) — proceeding (assuming a merged host)", err)
		return
	}
	if teamCreateCapable(string(out)) {
		t.Skipf("merged team-mode lane SKIPPED: the live claude (%q) is BELOW the merged floor (2.1.%d) — native TeamCreate is present, so the FO drives the legacy team registry, not the in-process merged shape this lane asserts. The legacy interactive pty lane covers this host.",
			strings.TrimSpace(string(out)), mergedFloorMinor)
	}
}
