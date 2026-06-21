// ABOUTME: AC-5 SPACEDOCK_MARKETPLACE_SOURCE override — the env var repoints the
// ABOUTME: install source on both hosts' install paths; unset keeps the default.
package cli

import (
	"bytes"
	"context"
	"testing"
)

// TestMarketplaceSourceOverrideDefault locks the unset case: with no
// SPACEDOCK_MARKETPLACE_SOURCE in the env, applyMarketplaceSourceOverride leaves the
// default `spacedock-dev/marketplace` in place — the released binary keeps its
// production source.
func TestMarketplaceSourceOverrideDefault(t *testing.T) {
	saved := marketplaceSource
	defer func() { marketplaceSource = saved }()

	applyMarketplaceSourceOverride([]string{"PATH=/usr/bin", "HOME=/home/x"})
	if marketplaceSource != "spacedock-dev/marketplace" {
		t.Fatalf("marketplaceSource = %q, want default spacedock-dev/marketplace when env unset", marketplaceSource)
	}
}

// TestMarketplaceSourceOverrideApplied locks the override case: an explicit
// SPACEDOCK_MARKETPLACE_SOURCE wins, repointing the source (e.g. at a local
// marketplace checkout for dogfooding).
func TestMarketplaceSourceOverrideApplied(t *testing.T) {
	saved := marketplaceSource
	defer func() { marketplaceSource = saved }()

	applyMarketplaceSourceOverride([]string{"SPACEDOCK_MARKETPLACE_SOURCE=/tmp/local-marketplace"})
	if marketplaceSource != "/tmp/local-marketplace" {
		t.Fatalf("marketplaceSource = %q, want /tmp/local-marketplace (override must win)", marketplaceSource)
	}
}

// TestInstallHonorsMarketplaceSourceOverride is the install-seam half of AC-5:
// after the override is applied, `spacedock install --host codex` passes the
// overridden source to the install seam — observed at fake.installCmds[1], the
// actual source the production install would `marketplace add`.
func TestInstallHonorsMarketplaceSourceOverride(t *testing.T) {
	savedSource := marketplaceSource
	savedBranch := devBranch
	devBranch = "main"
	defer func() {
		marketplaceSource = savedSource
		devBranch = savedBranch
	}()

	applyMarketplaceSourceOverride([]string{"SPACEDOCK_MARKETPLACE_SOURCE=/tmp/local-marketplace"})

	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer
	code := runInit(context.Background(), []string{"--host", "codex"}, fake, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if len(fake.installCmds) < 2 {
		t.Fatalf("install seam recorded %v, want at least {host, source}", fake.installCmds)
	}
	if got := fake.installCmds[1]; got != "/tmp/local-marketplace" {
		t.Fatalf("install source = %q, want the overridden /tmp/local-marketplace", got)
	}
}

// TestFrontDoorAutoInstallHonorsMarketplaceSourceOverride is the front-door half of
// AC-5: the codex no-plugin auto-install passes the overridden source to the
// install seam, so dogfooding a local marketplace works through the single-command
// front door, not only `install`.
func TestFrontDoorAutoInstallHonorsMarketplaceSourceOverride(t *testing.T) {
	savedSource := marketplaceSource
	savedBranch := devBranch
	devBranch = "main"
	defer func() {
		marketplaceSource = savedSource
		devBranch = savedBranch
	}()

	applyMarketplaceSourceOverride([]string{"SPACEDOCK_MARKETPLACE_SOURCE=/tmp/local-marketplace"})

	fake := &fakeHost{manifest: ""} // fresh HOME: no codex plugin installed
	var stdout, stderr bytes.Buffer
	code := runCodex(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if len(fake.installCmds) < 2 {
		t.Fatalf("install seam recorded %v, want at least {host, source}", fake.installCmds)
	}
	if got := fake.installCmds[1]; got != "/tmp/local-marketplace" {
		t.Fatalf("front-door auto-install source = %q, want the overridden /tmp/local-marketplace", got)
	}
}
