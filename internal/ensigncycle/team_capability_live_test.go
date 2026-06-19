//go:build live

// ABOUTME: Live wiring of the TeamCreate-capability gate — probes `claude --version` and
// ABOUTME: SKIPs the legacy interactive pty lane on a merged host (native team tools gone).
package ensigncycle

import (
	"os/exec"
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
