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
// uninstall; tolerated, fresh-box exit 1 with "Plugin not found in installed
// plugins"), then `plugin marketplace remove spacedock` (tolerated, fresh-box
// exit 1 with "not found"), then `plugin marketplace add` (with the @ref pin),
// then `plugin install spacedock@spacedock`. The marketplace-remove step is
// what defeats the "already on disk" no-op in marketplace add when a stale pin
// is declared. The asymmetry: BOTH cleanup steps (uninstall + remove) are
// tolerated; BOTH pinning steps (add + install) are fail-fast. With an empty
// branch the marketplace add arg is the bare source.
func TestInstallArgvSequence(t *testing.T) {
	wantWithBranch := []installStep{
		{argv: []string{"plugin", "uninstall", "spacedock@spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "marketplace", "remove", "spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "marketplace", "add", "spacedock-dev/spacedock@next"}},
		{argv: []string{"plugin", "install", "spacedock@spacedock"}},
	}
	if got := installArgvSequence("spacedock-dev/spacedock", "next"); !reflect.DeepEqual(got, wantWithBranch) {
		t.Errorf("installArgvSequence(branch=next) = %v, want %v", got, wantWithBranch)
	}

	wantBareSource := []installStep{
		{argv: []string{"plugin", "uninstall", "spacedock@spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "marketplace", "remove", "spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "marketplace", "add", "spacedock-dev/spacedock"}},
		{argv: []string{"plugin", "install", "spacedock@spacedock"}},
	}
	if got := installArgvSequence("spacedock-dev/spacedock", ""); !reflect.DeepEqual(got, wantBareSource) {
		t.Errorf("installArgvSequence(no branch) = %v, want %v", got, wantBareSource)
	}

	// Lock the tolerance asymmetry explicitly: the two cleanup steps (uninstall,
	// marketplace remove) are tolerated; the two pinning steps (marketplace add,
	// plugin install) are fail-fast. Any future drift toward tolerate-every-step
	// (or shifting tolerance onto a pinning step) fails here.
	seq := installArgvSequence("spacedock-dev/spacedock", "next")
	for i, step := range seq {
		isCleanup := isUninstallStep(step.argv) || isMarketplaceRemoveStep(step.argv)
		if isCleanup && !step.tolerateExit {
			t.Errorf("step %d (%v) tolerateExit = false, want true (cleanup step)", i, step.argv)
		}
		if !isCleanup && step.tolerateExit {
			t.Errorf("step %d (%v) tolerateExit = true, want false (pinning step must fail-fast)", i, step.argv)
		}
	}
}

func isUninstallStep(argv []string) bool {
	return len(argv) >= 2 && argv[0] == "plugin" && argv[1] == "uninstall"
}

func isMarketplaceRemoveStep(argv []string) bool {
	return len(argv) >= 3 && argv[0] == "plugin" && argv[1] == "marketplace" && argv[2] == "remove"
}
