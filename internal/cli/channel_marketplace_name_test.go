// ABOUTME: AC-2 channel-in-marketplace-name — the entry stays `spacedock` (== manifest
// ABOUTME: name, so the host name-match passes); the channel lives in the marketplace NAME.
package cli

import "testing"

// TestChannelEntryIsAlwaysManifestName locks the name-match-safe shape: the
// marketplace ENTRY name is always `spacedock`, equal to the plugin's own
// plugin.json `name`, on BOTH channels. Codex (and claude) reject an entry whose
// name differs from the manifest name, so encoding the channel in the entry name
// (the old `spacedock-edge` entry) made `codex plugin add` fail. The entry must be
// channel-invariant.
func TestChannelEntryIsAlwaysManifestName(t *testing.T) {
	for _, devBranch := range []string{"main", "next", "anything-else"} {
		if got := channelEntry(devBranch); got != "spacedock" {
			t.Errorf("channelEntry(%q) = %q, want spacedock (entry name must equal manifest name on every channel)", devBranch, got)
		}
	}
}

// TestChannelMarketplaceCarriesTheChannel locks the inverted model: the channel is
// the marketplace NAME. A stable binary (devBranch=main) resolves the `spacedock`
// marketplace; an edge binary (any other devBranch) resolves `spacedock-edge`.
func TestChannelMarketplaceCarriesTheChannel(t *testing.T) {
	cases := []struct {
		channel         string
		devBranch       string
		wantMarketplace string
	}{
		{channel: "stable", devBranch: "main", wantMarketplace: "spacedock"},
		{channel: "edge", devBranch: "next", wantMarketplace: "spacedock-edge"},
	}
	for _, tc := range cases {
		t.Run(tc.channel, func(t *testing.T) {
			if got := channelMarketplace(tc.devBranch); got != tc.wantMarketplace {
				t.Errorf("channelMarketplace(%q) = %q, want %q (the %s channel marketplace name)", tc.devBranch, got, tc.wantMarketplace, tc.channel)
			}
		})
	}
}

// TestChannelPluginIDIsEntryAtMarketplace locks the resolved id under the inverted
// model: `<entry>@<marketplace>` = `spacedock@spacedock` (stable) and
// `spacedock@spacedock-edge` (edge). The OLD broken shape `spacedock-edge@spacedock`
// (channel in the entry) must NOT appear — it is the id that fails the host
// name-match.
func TestChannelPluginIDIsEntryAtMarketplace(t *testing.T) {
	cases := []struct {
		channel   string
		devBranch string
		wantID    string
	}{
		{channel: "stable", devBranch: "main", wantID: "spacedock@spacedock"},
		{channel: "edge", devBranch: "next", wantID: "spacedock@spacedock-edge"},
	}
	for _, tc := range cases {
		t.Run(tc.channel, func(t *testing.T) {
			got := channelPluginID(tc.devBranch)
			if got != tc.wantID {
				t.Errorf("channelPluginID(%q) = %q, want %q", tc.devBranch, got, tc.wantID)
			}
			if got == "spacedock-edge@spacedock" {
				t.Errorf("channelPluginID(%q) = %q — the channel must be the marketplace name, not the entry name (this id fails the host name-match)", tc.devBranch, got)
			}
		})
	}
}
