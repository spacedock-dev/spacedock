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

// TestInstallArgvSequence asserts the argv sequence and the tolerance asymmetry of
// AC-1 and AC-3. The sibling ids below are independent literals for each channel,
// never derived from otherChannelMarketplace.
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

	// A change toward tolerance on every step, or tolerance on a pinning step,
	// fails here.
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
