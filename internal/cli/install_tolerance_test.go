// ABOUTME: AC-2/AC-3 install-tolerance — drives execHost.Install against a per-PATH
// ABOUTME: stub claude that exits 1 on a chosen subcommand, asserts tolerance asymmetry.
package cli

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestInstallToleratesMarketplaceUpdateStepFailure asserts AC-2 on the
// already-current path. The stub makes `claude plugin marketplace update
// spacedock-edge` exit 1, and every other subcommand exits 0. execHost.Install must
// return a nil error. The combined output must contain the stub marker of all 6
// steps. Without the tolerance flag on each step, this test fails at the update
// step.
func TestInstallToleratesMarketplaceUpdateStepFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	dir := writeClaudeStub(t, "plugin marketplace update")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	out, err := execHost{}.Install("claude", "spacedock-dev/marketplace", "next")
	if err != nil {
		t.Fatalf("Install returned error on tolerated update failure: %v\nout=%q", err, out)
	}
	for _, want := range []string{
		"stub:plugin marketplace update spacedock-edge:exit=1",
		"stub:plugin marketplace add spacedock-dev/marketplace:exit=0",
		"stub:plugin uninstall spacedock@spacedock-edge:exit=0",
		"stub:plugin uninstall spacedock@spacedock:exit=0",
		"stub:plugin uninstall spacedock-edge@spacedock:exit=0",
		"stub:plugin install spacedock@spacedock-edge:exit=0",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("combined output missing %q\nout=%q", want, out)
		}
	}
}

// TestInstallToleratesUninstallStepFailure asserts the fresh-box uninstall half of
// AC-2. The stub makes `claude plugin uninstall <channel-id>` exit 1, and every
// other subcommand exits 0. Claude 2.1.160 gives the shape "Plugin not found in
// installed plugins" on a box where the plugin is not installed. execHost.Install
// must return a nil error. The combined output must contain the stub marker of all
// 6 steps. Without the tolerance flag on the uninstall step, this test fails at
// step 0.
func TestInstallToleratesUninstallStepFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	dir := writeClaudeStub(t, "plugin uninstall spacedock@spacedock-edge")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	out, err := execHost{}.Install("claude", "spacedock-dev/marketplace", "next")
	if err != nil {
		t.Fatalf("Install returned error on tolerated uninstall failure: %v\nout=%q", err, out)
	}
	for _, want := range []string{
		"stub:plugin uninstall spacedock@spacedock-edge:exit=1",
		"stub:plugin uninstall spacedock@spacedock:exit=0",
		"stub:plugin uninstall spacedock-edge@spacedock:exit=0",
		"stub:plugin marketplace add spacedock-dev/marketplace:exit=0",
		"stub:plugin marketplace update spacedock-edge:exit=0",
		"stub:plugin install spacedock@spacedock-edge:exit=0",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("combined output missing %q\nout=%q", want, out)
		}
	}
}

// TestInstallToleratesRouteAMigrationStepFailure asserts the tolerance of the
// round-1 migration step. The stub makes `claude plugin uninstall
// spacedock-edge@spacedock`, the retired route-A id, exit 1. Every other subcommand
// exits 0. execHost.Install must return a nil error. The combined output must
// contain the stub marker of all 6 steps. Without the tolerance flag on this step,
// every install fails for a captain who never used the retired route.
func TestInstallToleratesRouteAMigrationStepFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	dir := writeClaudeStub(t, "plugin uninstall spacedock-edge@spacedock")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	out, err := execHost{}.Install("claude", "spacedock-dev/marketplace", "next")
	if err != nil {
		t.Fatalf("Install returned error on tolerated route-A migration failure: %v\nout=%q", err, out)
	}
	for _, want := range []string{
		"stub:plugin uninstall spacedock@spacedock-edge:exit=0",
		"stub:plugin uninstall spacedock@spacedock:exit=0",
		"stub:plugin uninstall spacedock-edge@spacedock:exit=1",
		"stub:plugin marketplace add spacedock-dev/marketplace:exit=0",
		"stub:plugin marketplace update spacedock-edge:exit=0",
		"stub:plugin install spacedock@spacedock-edge:exit=0",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("combined output missing %q\nout=%q", want, out)
		}
	}
}

// TestInstallFailsFastOnAddStep locks AC-3: tolerance is asymmetric. The
// marketplace add step is NOT tolerated, so a stub exiting 1 on add must surface
// a non-nil error wrapping the add subcommand. Guards against silent regression
// toward tolerate-every-step.
func TestInstallFailsFastOnAddStep(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	dir := writeClaudeStub(t, "plugin marketplace add")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	out, err := execHost{}.Install("claude", "spacedock-dev/marketplace", "next")
	if err == nil {
		t.Fatalf("Install returned nil error; want add-step failure\nout=%q", out)
	}
	if !strings.Contains(err.Error(), "plugin marketplace add spacedock-dev/marketplace") {
		t.Errorf("error %q does not wrap the add subcommand argv", err)
	}
	if !strings.Contains(out, "stub:plugin marketplace add spacedock-dev/marketplace:exit=1") {
		t.Errorf("combined output missing add-step stub stderr; out=%q", out)
	}
}

// writeClaudeStub writes a `claude` host stub (see writeHostStub).
func writeClaudeStub(t *testing.T, failOn string) string {
	t.Helper()
	return writeHostStub(t, "claude", failOn)
}

// writeHostStub writes a host CLI shell script named binName under a temp dir and
// returns the dir (suitable for PATH prepend). The stub echoes
// `stub:<joined args>:exit=<code>` to stdout, then exits with code 1 if the
// argv joined with single spaces starts with failOn, else exit 0. failOn is
// quoted into the case pattern so it can contain spaces (e.g. "plugin
// marketplace remove"). An empty failOn never matches, so every step exits 0.
// This lets a single helper drive both "fail on <step>" and all-zero tests for
// either host binary.
func writeHostStub(t *testing.T, binName, failOn string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, binName)
	matchArm := ""
	if failOn != "" {
		matchArm = `"` + failOn + `"*)
  echo "stub:$args:exit=1"
  exit 1
  ;;
`
	}
	body := `#!/bin/sh
args="$*"
case "$args" in
` + matchArm + `*)
  echo "stub:$args:exit=0"
  exit 0
  ;;
esac
`
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatalf("write stub: %v", err)
	}
	return dir
}
