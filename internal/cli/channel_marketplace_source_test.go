// ABOUTME: AC-2 channel-specific marketplace source — edge resolves the source
// ABOUTME: spacedock-dev/marketplace@edge, stable the bare repo; override wins verbatim.
package cli

import (
	"bytes"
	"context"
	"testing"
)

// TestChannelMarketplaceSourceFromDevBranch locks the source-resolution rule the
// marketplace add carries: a stable binary (devBranch=main) installs from the bare
// marketplace repo `spacedock-dev/marketplace` (whose root marketplace.json is named
// `spacedock`), while an edge binary (any other devBranch, e.g. next) installs from
// `spacedock-dev/marketplace@edge` — the `edge` branch whose root marketplace.json is
// named `spacedock-edge`. The @edge ref is what makes codex register a marketplace
// NAMED `spacedock-edge`, so the channel id `spacedock@spacedock-edge` then resolves.
func TestChannelMarketplaceSourceFromDevBranch(t *testing.T) {
	saved := marketplaceSource
	defer func() { marketplaceSource = saved }()
	marketplaceSource = "spacedock-dev/marketplace"

	cases := []struct {
		channel    string
		devBranch  string
		wantSource string
	}{
		{channel: "stable", devBranch: "main", wantSource: "spacedock-dev/marketplace"},
		{channel: "edge", devBranch: "next", wantSource: "spacedock-dev/marketplace@edge"},
	}
	for _, tc := range cases {
		t.Run(tc.channel, func(t *testing.T) {
			if got := channelMarketplaceSource(tc.devBranch); got != tc.wantSource {
				t.Errorf("channelMarketplaceSource(%q) = %q, want %q (the %s channel source)", tc.devBranch, got, tc.wantSource, tc.channel)
			}
		})
	}
}

// TestChannelMarketplaceSourceOverrideWinsVerbatim locks AC-2's override half: an
// explicit SPACEDOCK_MARKETPLACE_SOURCE (dogfooding a local/alternate marketplace)
// replaces the resolved source verbatim on BOTH channels — no `@edge` ref is appended
// to the operator's chosen source, because that source is already the complete
// marketplace they want (appending @edge to a local checkout path is meaningless).
func TestChannelMarketplaceSourceOverrideWinsVerbatim(t *testing.T) {
	saved := marketplaceSource
	defer func() { marketplaceSource = saved }()
	marketplaceSource = "/tmp/local-marketplace"

	for _, devBranch := range []string{"main", "next"} {
		t.Run(devBranch, func(t *testing.T) {
			if got := channelMarketplaceSource(devBranch); got != "/tmp/local-marketplace" {
				t.Errorf("channelMarketplaceSource(%q) = %q, want the overridden /tmp/local-marketplace verbatim on every channel", devBranch, got)
			}
		})
	}
}

// TestEdgeInstallSeamCarriesEdgeSource is AC-2's install-seam half: `spacedock
// install --host codex` on an edge binary passes `spacedock-dev/marketplace@edge` to
// the install seam (observed at installCmds[1]), so the production `marketplace add`
// targets the edge branch — the fix for the masked bug where the bare source
// registered a marketplace named `spacedock` and `plugin add spacedock@spacedock-edge`
// then failed.
func TestEdgeInstallSeamCarriesEdgeSource(t *testing.T) {
	savedSource := marketplaceSource
	savedBranch := devBranch
	marketplaceSource = "spacedock-dev/marketplace"
	devBranch = "next"
	defer func() {
		marketplaceSource = savedSource
		devBranch = savedBranch
	}()

	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer
	code := runInit(context.Background(), []string{"--host", "codex"}, fake, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if len(fake.installCmds) < 2 {
		t.Fatalf("install seam recorded %v, want at least {host, source}", fake.installCmds)
	}
	if got := fake.installCmds[1]; got != "spacedock-dev/marketplace@edge" {
		t.Fatalf("edge install source = %q, want spacedock-dev/marketplace@edge", got)
	}
}
