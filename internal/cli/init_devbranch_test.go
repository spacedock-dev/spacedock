// ABOUTME: AC-3a — `spacedock install --host claude` targets @next when devBranch is
// ABOUTME: pinned to next, and the composed marketplace-add argv carries @next.
package cli

import (
	"bytes"
	"context"
	"reflect"
	"testing"
)

// TestInitTargetsNextWhenDevBranchPinned locks AC-3a: with devBranch pinned to
// `next` (the released binary's default, until `next` is the default branch),
// `spacedock install --host claude` drives the install seam with branch=next, so the
// issued `marketplace add` resolves `spacedock-dev/spacedock@next`. The package
// var is saved/restored so the assertion does not leak into sibling tests.
func TestInitTargetsNextWhenDevBranchPinned(t *testing.T) {
	saved := devBranch
	devBranch = "next"
	defer func() { devBranch = saved }()

	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runInit(context.Background(), []string{"--host", "claude"}, fake, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	// Install records {host, source, branch}; branch carries the @ref pin.
	if len(fake.installCmds) < 3 {
		t.Fatalf("install seam recorded %v, want {host, source, branch}", fake.installCmds)
	}
	if got := fake.installCmds[2]; got != "next" {
		t.Fatalf("install branch = %q, want next (init must target @next)", got)
	}
}

// TestMarketplaceAddArgvCarriesRef locks the argv composition AC-3a asserts: the
// `claude plugin marketplace add` argv pins `source@branch` when a branch is set,
// and is the bare source when it is not. This is the exact 2-command argv shape
// owned today; task 38 changes Install to a 3-command shape (add/uninstall/
// install) and this assertion is updated in lockstep then.
func TestMarketplaceAddArgvCarriesRef(t *testing.T) {
	if got := marketplaceAddArg("spacedock-dev/spacedock", "next"); got != "spacedock-dev/spacedock@next" {
		t.Errorf("marketplaceAddArg with branch = %q, want spacedock-dev/spacedock@next", got)
	}
	if got := marketplaceAddArg("spacedock-dev/spacedock", ""); got != "spacedock-dev/spacedock" {
		t.Errorf("marketplaceAddArg without branch = %q, want bare source", got)
	}
}

// TestInstallArgvSequence locks AC-1 and AC-3's tolerance asymmetry:
// execHost.Install issues the 4-command upgrade shape — `plugin uninstall
// spacedock@spacedock` first (claude tracks an installed plugin via its
// marketplace record, so the marketplace remove later would orphan a live
// uninstall), then `plugin marketplace remove spacedock` (tolerated, fresh-box
// exit 1), then `plugin marketplace add` (with the @ref pin), then `plugin
// install spacedock@spacedock`. The remove step is what defeats the "already
// on disk" no-op in marketplace add when a stale pin is declared; remove is
// tolerated because it exits 1 on a fresh box (claude reports "not found").
// Every other step is fail-fast. With an empty branch the marketplace add arg
// is the bare source.
func TestInstallArgvSequence(t *testing.T) {
	wantWithBranch := []installStep{
		{argv: []string{"plugin", "uninstall", "spacedock@spacedock"}},
		{argv: []string{"plugin", "marketplace", "remove", "spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "marketplace", "add", "spacedock-dev/spacedock@next"}},
		{argv: []string{"plugin", "install", "spacedock@spacedock"}},
	}
	if got := installArgvSequence("spacedock-dev/spacedock", "next"); !reflect.DeepEqual(got, wantWithBranch) {
		t.Errorf("installArgvSequence(branch=next) = %v, want %v", got, wantWithBranch)
	}

	wantBareSource := []installStep{
		{argv: []string{"plugin", "uninstall", "spacedock@spacedock"}},
		{argv: []string{"plugin", "marketplace", "remove", "spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "marketplace", "add", "spacedock-dev/spacedock"}},
		{argv: []string{"plugin", "install", "spacedock@spacedock"}},
	}
	if got := installArgvSequence("spacedock-dev/spacedock", ""); !reflect.DeepEqual(got, wantBareSource) {
		t.Errorf("installArgvSequence(no branch) = %v, want %v", got, wantBareSource)
	}

	// Lock the tolerance asymmetry explicitly: ONLY the marketplace-remove step
	// is tolerated. Any future drift toward tolerate-every-step (or moving the
	// tolerance to a different step) fails here.
	seq := installArgvSequence("spacedock-dev/spacedock", "next")
	for i, step := range seq {
		isRemove := len(step.argv) >= 3 && step.argv[1] == "marketplace" && step.argv[2] == "remove"
		if isRemove && !step.tolerateExit {
			t.Errorf("step %d (remove) tolerateExit = false, want true", i)
		}
		if !isRemove && step.tolerateExit {
			t.Errorf("step %d (%v) tolerateExit = true, want false (only marketplace-remove is tolerated)", i, step.argv)
		}
	}
}
