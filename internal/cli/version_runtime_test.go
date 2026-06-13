// ABOUTME: AC-3/AC-5 oracles for --version — the load-bearing first line, the
// ABOUTME: Sandbox line, and the per-runtime install/enablement block from stubs.
package cli

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/contract"
)

// fakeRuntimeProbe pins each host's install/enablement outcome so the per-runtime
// block renders from injected state, never a live host CLI.
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

// TestVersionFirstLineUnchanged (AC-3 part 1) pins the load-bearing invariant: the
// FIRST line still matches `^spacedock .* \(contract \d+\)$`, the token the FO and
// ensign skills parse. Appending the sandbox + per-runtime block must not disturb
// line 1.
func TestVersionFirstLineUnchanged(t *testing.T) {
	probe := fakeRuntimeProbe{status: map[string]runtimeStatus{
		"claude": {installed: true, enablement: enablementEnabled},
		"codex":  {installed: true, enablement: enablementNotEnabled},
		"pi":     {installed: false},
	}}
	out := renderVersion(probe)
	firstLine := strings.SplitN(out, "\n", 2)[0]
	re := regexp.MustCompile(`^spacedock .* \(contract \d+\)$`)
	if !re.MatchString(firstLine) {
		t.Fatalf("version first line %q does not match the FO-parsed pattern", firstLine)
	}
	want := "spacedock " + Version + " (contract " + strconv.Itoa(contract.CONTRACT_VERSION) + ")"
	if firstLine != want {
		t.Fatalf("version first line = %q, want %q", firstLine, want)
	}
}

// TestVersionPerRuntimeBlock (AC-3 part 2) drives each install/enablement outcome
// through the injected probe and asserts the exact per-runtime line, plus a Sandbox
// line. The expected strings are independent test-supplied values.
func TestVersionPerRuntimeBlock(t *testing.T) {
	probe := fakeRuntimeProbe{status: map[string]runtimeStatus{
		"claude": {installed: true, enablement: enablementEnabled},
		"codex":  {installed: true, enablement: enablementNotEnabled},
		"pi":     {installed: false},
	}}
	out := renderVersion(probe)

	for _, want := range []string{
		"Sandbox: unavailable (safehouse not on PATH)",
		"claude: installed, spacedock enabled",
		"codex: installed",
		"pi: not installed",
	} {
		if !lineEquals(out, want) {
			t.Fatalf("version output missing whole line %q:\n%s", want, out)
		}
	}
}

// TestVersionEnablementUnknown (AC-5) injects a probe whose enablement read errored
// (binary resolves, but the probe could not determine enablement — the sandboxed
// `codex plugin list` "Operation not permitted" mode). The line must read
// `installed, enablement unknown`, never silently `not installed`. A separate
// binary-absent case asserts `not installed`.
func TestVersionEnablementUnknown(t *testing.T) {
	probe := fakeRuntimeProbe{status: map[string]runtimeStatus{
		"claude": {installed: true, enablement: enablementUnknown},
		"codex":  {installed: false}, // binary absent → not installed
		"pi":     {installed: true, enablement: enablementUnknown},
	}}
	out := renderVersion(probe)

	for _, want := range []string{
		"claude: installed, enablement unknown",
		"codex: not installed",
		"pi: installed, enablement unknown",
	} {
		if !lineEquals(out, want) {
			t.Fatalf("version output missing whole line %q:\n%s", want, out)
		}
	}
	// AC-5 guard: an enablement-unknown host must NOT be rendered as not installed.
	if lineEquals(out, "claude: not installed") {
		t.Fatalf("enablement-unknown claude silently rendered as not installed:\n%s", out)
	}
}
