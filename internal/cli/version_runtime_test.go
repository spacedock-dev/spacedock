// ABOUTME: AC oracles for --version — the load-bearing first line, the Sandbox
// ABOUTME: line, and the version-forward per-runtime block rendered from stubs.
package cli

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

// fakeRuntimeProbe pins each host's install/version/marker outcome so the
// per-runtime block renders from injected state, never a live host CLI.
type fakeRuntimeProbe struct {
	status map[string]runtimeStatus
}

func (f fakeRuntimeProbe) ProbeRuntime(host string) runtimeStatus {
	return f.status[host]
}

// renderVersion captures printVersion over an injected probe and a pinned lookPath
// (so the Sandbox line is a stable `unavailable`, asserted on its own elsewhere)
// and returns stdout.
func renderVersion(probe runtimeProbe) string {
	var buf bytes.Buffer
	printVersion(&buf, t1TempDir, probe, lookMissing)
	return buf.String()
}

// t1TempDir is a directory with no .safehouse profile, so the version Sandbox line
// renders deterministically as unavailable under lookMissing. A literal keeps the
// helper allocation-free; the path need not exist for safehouse.Present.
const t1TempDir = "/nonexistent-version-dir"

// TestVersionFirstLineUnchanged pins the load-bearing invariant: the FIRST line
// still matches `^spacedock .* \(contract \d+\)$`, the token the FO and ensign
// skills parse. Appending the sandbox + per-runtime block must not disturb line 1.
func TestVersionFirstLineUnchanged(t *testing.T) {
	probe := fakeRuntimeProbe{status: map[string]runtimeStatus{
		"claude": {installed: true, version: "0.20.0", marker: markerEnabled},
		"codex":  {installed: true, version: "0.20.0", marker: markerDisabled},
		"pi":     {installed: false},
	}}
	out := renderVersion(probe)
	firstLine := strings.SplitN(out, "\n", 2)[0]
	re := regexp.MustCompile(`^spacedock .* \(contract \d+\)$`)
	if !re.MatchString(firstLine) {
		t.Fatalf("version first line %q does not match the FO-parsed pattern", firstLine)
	}
	want := "spacedock " + displayVersion() + " " + frozenContractToken
	if firstLine != want {
		t.Fatalf("version first line = %q, want %q", firstLine, want)
	}
}

// TestVersionPerRuntimeBlock drives each install/version outcome through the
// injected probe and asserts the exact version-forward per-runtime line, plus a
// Sandbox line. The expected strings are independent test-supplied values.
func TestVersionPerRuntimeBlock(t *testing.T) {
	probe := fakeRuntimeProbe{status: map[string]runtimeStatus{
		"claude": {installed: true, version: "0.20.0", marker: markerEnabled},
		"codex":  {installed: true, version: "0.20.0", marker: markerDisabled},
		"pi":     {installed: true, ready: true},
	}}
	out := renderVersion(probe)

	for _, want := range []string{
		"Sandbox: unavailable (safehouse not on PATH)",
		"claude: spacedock 0.20.0",
		"codex: spacedock 0.20.0 (disabled)",
		"pi: spacedock ready",
	} {
		if !lineEquals(out, want) {
			t.Fatalf("version output missing whole line %q:\n%s", want, out)
		}
	}
}

// TestVersionHostAbsentAndNoPlugin asserts the two not-installed shapes: a host
// whose binary is absent reads `<host>: not installed`, while a host that resolves
// but carries no plugin reads `<host>: spacedock not installed`.
func TestVersionHostAbsentAndNoPlugin(t *testing.T) {
	probe := fakeRuntimeProbe{status: map[string]runtimeStatus{
		"claude": {installed: false},              // binary absent
		"codex":  {installed: true, version: ""},  // host present, no plugin
		"pi":     {installed: true, ready: false}, // host present, not launch-ready
	}}
	out := renderVersion(probe)

	for _, want := range []string{
		"claude: not installed",
		"codex: spacedock not installed",
		"pi: spacedock not installed",
	} {
		if !lineEquals(out, want) {
			t.Fatalf("version output missing whole line %q:\n%s", want, out)
		}
	}
}

// TestVersionMarkerUnknownRendersBareVersion is the AC-2 oracle: when the
// enabled/disabled probe could not read state (markerUnknown) but the manifest
// resolved a version, the line shows the bare `spacedock <version>` with NO marker
// — the version still renders from the manifest even though the probe failed. It
// must never invent an "unknown" word or silently read as not installed.
func TestVersionMarkerUnknownRendersBareVersion(t *testing.T) {
	probe := fakeRuntimeProbe{status: map[string]runtimeStatus{
		"claude": {installed: true, version: "0.20.0", marker: markerUnknown},
		"codex":  {installed: false},
		"pi":     {installed: false},
	}}
	out := renderVersion(probe)

	if !lineEquals(out, "claude: spacedock 0.20.0") {
		t.Fatalf("marker-unknown claude did not render the bare version line:\n%s", out)
	}
	// Guard: a probe-failed-but-version-resolved host must NOT read as not installed.
	if lineEquals(out, "claude: not installed") || lineEquals(out, "claude: spacedock not installed") {
		t.Fatalf("marker-unknown claude with a resolved version rendered as not installed:\n%s", out)
	}
}

// TestVersionVocabularyHasNoEnablementJargon is the AC-1 oracle on the rendered
// output: the noun "enablement" and the retired phrasings must appear nowhere in
// the per-runtime block across the install/version/marker outcomes.
func TestVersionVocabularyHasNoEnablementJargon(t *testing.T) {
	probe := fakeRuntimeProbe{status: map[string]runtimeStatus{
		"claude": {installed: true, version: "0.20.0", marker: markerUnknown},
		"codex":  {installed: true, version: "0.20.0", marker: markerDisabled},
		"pi":     {installed: true, ready: true},
	}}
	out := renderVersion(probe)
	for _, banned := range []string{
		"enablement",
		"spacedock enabled",
		"enablement unknown",
	} {
		if strings.Contains(out, banned) {
			t.Fatalf("version output contains retired jargon %q:\n%s", banned, out)
		}
	}
}
