// ABOUTME: AC-3a — `spacedock install --host claude` drives the install seam with the
// ABOUTME: binary's devBranch, which selects the marketplace channel entry to install.
package cli

import (
	"bytes"
	"context"
	"reflect"
	"testing"
)

// TestInitTargetsNextWhenDevBranchPinned locks AC-3a: with devBranch pinned to
// `next` (the released edge binary's default, until `next` is the default
// branch), `spacedock install --host claude` drives the install seam with
// devBranch=next, so the issued sequence installs the `spacedock-edge` channel
// entry. The package var is saved/restored so the assertion does not leak into
// sibling tests.
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
	// Install records {host, source, devBranch}; devBranch selects the channel entry.
	if len(fake.installCmds) < 3 {
		t.Fatalf("install seam recorded %v, want {host, source, devBranch}", fake.installCmds)
	}
	if got := fake.installCmds[2]; got != "next" {
		t.Fatalf("install devBranch = %q, want next (edge binary installs the edge channel)", got)
	}
}

// TestInstallArgvSequence asserts the tolerance asymmetry of AC-1 and AC-3, and the
// non-destructive shape. execHost.Install runs 6 commands in this order: `plugin
// uninstall <id>`, `plugin uninstall <sibling channel id>`, `plugin uninstall
// spacedock-edge@spacedock`, `plugin marketplace add <source>`, `plugin marketplace
// update <marketplace>`, and `plugin install <id>`.
//
// The first uninstall is tolerated. Claude tracks an installed plugin through its
// marketplace record, and it exits 1 on a fresh box with "Plugin not found in
// installed plugins". The second uninstall keeps one channel on the host. The
// sibling ids below are independent literals for each channel. They never come from
// otherChannelMarketplace, so this test and the production sequence can diverge.
// The third uninstall is the round-1 route-A migration, and devBranch does not gate
// it.
//
// The builder takes any source. channelMarketplaceSource resolves the channel
// source before this call: the bare repo for stable and `<repo>@edge` for edge. The
// id is the channel that devBranch selects, with the channel in the marketplace
// name. For the edge binary the id is `spacedock@spacedock-edge`, and the
// marketplace update targets `spacedock-edge`. For stable the id is
// `spacedock@spacedock`.
//
// No step uses `plugin marketplace remove`. Probe 1 measured that this command
// uninstalls every co-hosted plugin. The four cleanup steps (three uninstalls and
// the marketplace update) are tolerated. The two pinning steps (marketplace add and
// plugin install) are fail-fast.
func TestInstallArgvSequence(t *testing.T) {
	wantEdge := []installStep{
		{argv: []string{"plugin", "uninstall", "spacedock@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "uninstall", "spacedock@spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "uninstall", "spacedock-edge@spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "marketplace", "add", "spacedock-dev/marketplace"}},
		{argv: []string{"plugin", "marketplace", "update", "spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "install", "spacedock@spacedock-edge"}},
	}
	if got := installArgvSequence("spacedock-dev/marketplace", "next"); !reflect.DeepEqual(got, wantEdge) {
		t.Errorf("installArgvSequence(devBranch=next) = %v, want %v", got, wantEdge)
	}

	wantStable := []installStep{
		{argv: []string{"plugin", "uninstall", "spacedock@spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "uninstall", "spacedock@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "uninstall", "spacedock-edge@spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "marketplace", "add", "spacedock-dev/marketplace"}},
		{argv: []string{"plugin", "marketplace", "update", "spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "install", "spacedock@spacedock"}},
	}
	if got := installArgvSequence("spacedock-dev/marketplace", "main"); !reflect.DeepEqual(got, wantStable) {
		t.Errorf("installArgvSequence(devBranch=main) = %v, want %v", got, wantStable)
	}

	// This loop asserts the tolerance asymmetry directly. The four cleanup steps
	// (three uninstalls and the marketplace update) are tolerated. The two pinning
	// steps (marketplace add and plugin install) are fail-fast. A change toward
	// tolerance on every step, or tolerance on a pinning step, fails here.
	seq := installArgvSequence("spacedock-dev/marketplace", "next")
	for i, step := range seq {
		isCleanup := isUninstallStep(step.argv) || isMarketplaceUpdateStep(step.argv)
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

func isMarketplaceUpdateStep(argv []string) bool {
	return len(argv) >= 3 && argv[0] == "plugin" && argv[1] == "marketplace" && argv[2] == "update"
}
