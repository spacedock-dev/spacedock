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
// every step → Install("codex", "spacedock-dev/spacedock", "next") returns nil and
// the combined output carries all four step markers, including the codex-specific
// `plugin marketplace add spacedock-dev/spacedock --ref next` and `plugin add
// spacedock@spacedock`. The stub's echoed argv is the independent source of truth.
func TestCodexInstallIssuesSequenceInOrder(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	dir := writeHostStub(t, "codex", "")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	out, err := execHost{}.Install("codex", "spacedock-dev/spacedock", "next")
	if err != nil {
		t.Fatalf("Install returned error on all-zero codex stub: %v\nout=%q", err, out)
	}
	wantOrder := []string{
		"stub:plugin remove spacedock@spacedock:exit=0",
		"stub:plugin marketplace remove spacedock:exit=0",
		"stub:plugin marketplace add spacedock-dev/spacedock --ref next:exit=0",
		"stub:plugin add spacedock@spacedock:exit=0",
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

// TestCodexInstallToleratesMarketplaceRemoveFailure locks AC-1b: the fresh-box
// path where `codex plugin marketplace remove spacedock` exits 1 ("is not
// configured or installed") and every other step exits 0. Install MUST return a
// nil error and the combined output MUST carry all four step markers — the
// marketplace remove is a tolerated cleanup step.
func TestCodexInstallToleratesMarketplaceRemoveFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	dir := writeHostStub(t, "codex", "plugin marketplace remove")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	out, err := execHost{}.Install("codex", "spacedock-dev/spacedock", "next")
	if err != nil {
		t.Fatalf("Install returned error on tolerated marketplace-remove failure: %v\nout=%q", err, out)
	}
	for _, want := range []string{
		"stub:plugin remove spacedock@spacedock:exit=0",
		"stub:plugin marketplace remove spacedock:exit=1",
		"stub:plugin marketplace add spacedock-dev/spacedock --ref next:exit=0",
		"stub:plugin add spacedock@spacedock:exit=0",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("combined output missing %q\nout=%q", want, out)
		}
	}
}

// TestCodexInstallToleratesPluginRemoveFailure locks the cleanup half of AC-1b:
// when `codex plugin remove spacedock@spacedock` exits 1 and every other step
// exits 0, Install MUST still return nil — the plugin remove is a tolerated
// cleanup step (idempotent on a fresh box per the spike, exit 0; tolerated here
// so a future codex that exits 1 on a fresh box does not break the install).
func TestCodexInstallToleratesPluginRemoveFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	dir := writeHostStub(t, "codex", "plugin remove")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	out, err := execHost{}.Install("codex", "spacedock-dev/spacedock", "next")
	if err != nil {
		t.Fatalf("Install returned error on tolerated plugin-remove failure: %v\nout=%q", err, out)
	}
	for _, want := range []string{
		"stub:plugin remove spacedock@spacedock:exit=1",
		"stub:plugin marketplace add spacedock-dev/spacedock --ref next:exit=0",
		"stub:plugin add spacedock@spacedock:exit=0",
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

	out, err := execHost{}.Install("codex", "spacedock-dev/spacedock", "next")
	if err == nil {
		t.Fatalf("Install returned nil error; want marketplace-add fail-fast\nout=%q", out)
	}
	if !strings.Contains(err.Error(), "plugin marketplace add spacedock-dev/spacedock --ref next") {
		t.Errorf("error %q does not wrap the codex add subcommand argv", err)
	}
	if !strings.Contains(out, "stub:plugin marketplace add spacedock-dev/spacedock --ref next:exit=1") {
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

	out, err := execHost{}.Install("codex", "spacedock-dev/spacedock", "next")
	if err == nil {
		t.Fatalf("Install returned nil error; want plugin-add fail-fast\nout=%q", out)
	}
	if !strings.Contains(err.Error(), "plugin add spacedock@spacedock") {
		t.Errorf("error %q does not wrap the codex plugin add subcommand argv", err)
	}
	if !strings.Contains(out, "stub:plugin add spacedock@spacedock:exit=1") {
		t.Errorf("combined output missing plugin-add stub marker; out=%q", out)
	}
}

// TestCodexInstallOmitsRefWhenBranchEmpty locks the empty-branch arm of AC-1: with
// branch=="" the marketplace add carries the bare source and NO `--ref` token. A
// stub recording the exact argv is the source of truth — a leaked `--ref` would
// pin against a non-existent ref on the default-branch install path.
func TestCodexInstallOmitsRefWhenBranchEmpty(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	dir := writeHostStub(t, "codex", "")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	out, err := execHost{}.Install("codex", "spacedock-dev/spacedock", "")
	if err != nil {
		t.Fatalf("Install returned error on empty-branch codex stub: %v\nout=%q", err, out)
	}
	if !strings.Contains(out, "stub:plugin marketplace add spacedock-dev/spacedock:exit=0") {
		t.Errorf("combined output missing bare-source add marker; out=%q", out)
	}
	if strings.Contains(out, "--ref") {
		t.Errorf("combined output carries a --ref token with empty branch; out=%q", out)
	}
}
