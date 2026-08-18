// ABOUTME: AC-1 codex install-tolerance — drives execHost.Install against a per-PATH
// ABOUTME: stub codex, asserting the codex install sequence's order + tolerance asymmetry.
package cli

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

// TestCodexInstallIssuesSequenceInOrder locks AC-1a: a stub codex that exits 0 on
// every step -> Install("codex", "spacedock-dev/marketplace", "next") returns nil
// and the combined output carries all six step markers. The sibling channel's
// plugin is removed, then the selected channel's own, then any prior
// --plugin-dir dev install (spacedock@spacedock-local — the round-3 round-trip
// fix), and no step ever spells `marketplace remove` — codex cannot keep a
// stale sibling `spacedock:*` provider enabled beside the selected install, but
// the marketplace record it lives under is never touched. The stub's echoed
// argv is the independent source of truth.
func TestCodexInstallIssuesSequenceInOrder(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	dir := writeHostStub(t, "codex", "")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	out, err := execHost{}.Install("codex", "spacedock-dev/marketplace", "next")
	if err != nil {
		t.Fatalf("Install returned error on all-zero codex stub: %v\nout=%q", err, out)
	}
	wantOrder := []string{
		"stub:plugin remove spacedock@spacedock:exit=0",
		"stub:plugin remove spacedock@spacedock-edge:exit=0",
		"stub:plugin remove spacedock@spacedock-local:exit=0",
		"stub:plugin marketplace add spacedock-dev/marketplace:exit=0",
		"stub:plugin marketplace upgrade spacedock-edge:exit=0",
		"stub:plugin add spacedock@spacedock-edge:exit=0",
	}
	last := -1
	for _, want := range wantOrder {
		idx := strings.Index(out, want)
		if idx < 0 {
			t.Errorf("combined output missing %q\nout=%q", want, out)
			continue
		}
		if idx <= last {
			t.Errorf("step %q out of order (idx=%d, prev=%d)\nout=%q", want, idx, last, out)
		}
		last = idx
	}
}

// TestCodexInstallToleratesMarketplaceUpgradeFailure locks AC-1b: the
// local-source path where `codex plugin marketplace upgrade <marketplace>` prints
// an error but the stub is made to exit 1 here to prove tolerance, while every
// other step exits 0. Install MUST return a nil error and the combined output
// MUST carry both channel-remove markers before the add/upgrade/add markers —
// marketplace upgrade is a tolerated refresh step, and no step spells
// `marketplace remove`.
func TestCodexInstallToleratesMarketplaceUpgradeFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	dir := writeHostStub(t, "codex", "plugin marketplace upgrade")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	out, err := execHost{}.Install("codex", "spacedock-dev/marketplace", "next")
	if err != nil {
		t.Fatalf("Install returned error on tolerated marketplace-upgrade failure: %v\nout=%q", err, out)
	}
	for _, want := range []string{
		"stub:plugin remove spacedock@spacedock:exit=0",
		"stub:plugin remove spacedock@spacedock-edge:exit=0",
		"stub:plugin marketplace add spacedock-dev/marketplace:exit=0",
		"stub:plugin marketplace upgrade spacedock-edge:exit=1",
		"stub:plugin add spacedock@spacedock-edge:exit=0",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("combined output missing %q\nout=%q", want, out)
		}
	}
}

// TestCodexInstallToleratesPluginRemoveFailure locks the cleanup half of AC-1b:
// when `codex plugin remove <spacedock channel id>` exits 1 and every other step
// exits 0, Install MUST still return nil — plugin remove is a tolerated cleanup
// step (idempotent on a fresh box per the spike, exit 0; tolerated here so a
// future codex that exits 1 on a fresh box does not break the install).
func TestCodexInstallToleratesPluginRemoveFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	dir := writeHostStub(t, "codex", "plugin remove")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	out, err := execHost{}.Install("codex", "spacedock-dev/marketplace", "next")
	if err != nil {
		t.Fatalf("Install returned error on tolerated plugin-remove failure: %v\nout=%q", err, out)
	}
	for _, want := range []string{
		"stub:plugin remove spacedock@spacedock:exit=1",
		"stub:plugin remove spacedock@spacedock-edge:exit=1",
		"stub:plugin marketplace add spacedock-dev/marketplace:exit=0",
		"stub:plugin add spacedock@spacedock-edge:exit=0",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("combined output missing %q\nout=%q", want, out)
		}
	}
}

// TestCodexInstallFailsFastOnMarketplaceAdd locks AC-1c: tolerance is asymmetric.
// The marketplace add is NOT tolerated, so a stub exiting 1 on add must surface a
// non-nil error wrapping that subcommand argv (the real-failure backstop).
func TestCodexInstallFailsFastOnMarketplaceAdd(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	dir := writeHostStub(t, "codex", "plugin marketplace add")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	out, err := execHost{}.Install("codex", "spacedock-dev/marketplace", "next")
	if err == nil {
		t.Fatalf("Install returned nil error; want marketplace-add fail-fast\nout=%q", out)
	}
	if !strings.Contains(err.Error(), "plugin marketplace add spacedock-dev/marketplace") {
		t.Errorf("error %q does not wrap the codex add subcommand argv", err)
	}
	if !strings.Contains(out, "stub:plugin marketplace add spacedock-dev/marketplace:exit=1") {
		t.Errorf("combined output missing add-step stub marker; out=%q", out)
	}
}

// TestCodexInstallFailsFastOnPluginAdd locks the other fail-fast arm of AC-1c: the
// final `plugin add` step is not tolerated either. A stub exiting 1 on it must
// surface a non-nil error wrapping that subcommand argv.
func TestCodexInstallFailsFastOnPluginAdd(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	dir := writeHostStub(t, "codex", "plugin add")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	out, err := execHost{}.Install("codex", "spacedock-dev/marketplace", "next")
	if err == nil {
		t.Fatalf("Install returned nil error; want plugin-add fail-fast\nout=%q", out)
	}
	if !strings.Contains(err.Error(), "plugin add spacedock@spacedock-edge") {
		t.Errorf("error %q does not wrap the codex plugin add subcommand argv", err)
	}
	if !strings.Contains(out, "stub:plugin add spacedock@spacedock-edge:exit=1") {
		t.Errorf("combined output missing plugin-add stub marker; out=%q", out)
	}
}

// TestCodexInstallStableEntryOmitsRef locks the stable-channel arm of AC-1 (and
// the no-`--ref` invariant): the marketplace add carries the bare marketplace-repo
// source and NO `--ref` token (the channel is the entry name, not a branch ref),
// and the stable binary (devBranch=main) installs the `spacedock` entry. A stub
// recording the exact argv is the source of truth — a leaked `--ref` would wrongly
// pin a non-existent ref, and a wrong entry id would cross channels.
func TestCodexInstallStableEntryOmitsRef(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	dir := writeHostStub(t, "codex", "")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	out, err := execHost{}.Install("codex", "spacedock-dev/marketplace", "main")
	if err != nil {
		t.Fatalf("Install returned error on stable-channel codex stub: %v\nout=%q", err, out)
	}
	if !strings.Contains(out, "stub:plugin marketplace add spacedock-dev/marketplace:exit=0") {
		t.Errorf("combined output missing bare-source add marker; out=%q", out)
	}
	if !strings.Contains(out, "stub:plugin remove spacedock@spacedock-edge:exit=0") {
		t.Errorf("stable install must still remove the edge sibling; out=%q", out)
	}
	if !strings.Contains(out, "stub:plugin add spacedock@spacedock:exit=0") {
		t.Errorf("combined output missing stable-channel plugin add marker; out=%q", out)
	}
	if strings.Contains(out, "--ref") {
		t.Errorf("combined output carries a --ref token; the channel is the entry name, not a branch ref; out=%q", out)
	}
}
