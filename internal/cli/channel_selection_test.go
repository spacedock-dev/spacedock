// ABOUTME: AC-3 channel-resolution seam — devBranch selects the marketplace NAME
// ABOUTME: (spacedock vs spacedock-edge); the entry stays spacedock, no @ref shorthand.
package cli

import (
	"bytes"
	"context"
	"reflect"
	"strings"
	"testing"
)

// TestChannelMarketplaceFromDevBranch locks the channel-selection rule: the
// binary's devBranch stamp selects the marketplace NAME, not a git ref pinned into
// the plugin repo. A stable binary (devBranch=main) resolves the `spacedock`
// marketplace; an edge binary (devBranch=next) resolves `spacedock-edge`. The tag
// pin lives in the marketplace manifest, so the channel is the marketplace name,
// not an @branch shorthand on the install command. The entry name stays `spacedock`
// on both channels (== manifest name, so the host name-match passes).
func TestChannelMarketplaceFromDevBranch(t *testing.T) {
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
				t.Errorf("channelMarketplace(%q) = %q, want %q (the %s channel)", tc.devBranch, got, tc.wantMarketplace, tc.channel)
			}
			if got := channelEntry(tc.devBranch); got != "spacedock" {
				t.Errorf("channelEntry(%q) = %q, want spacedock (entry stays the manifest name on every channel)", tc.devBranch, got)
			}
		})
	}
}

// TestClaudeChannelInstallArgvSequence is AC-3's claude half: with marketplaceSource
// repointed to the marketplace repo and devBranch set per channel, the issued
// claude install argv installs the channel-correct id (`spacedock@spacedock` stable
// / `spacedock@spacedock-edge` edge) and the marketplace add carries the BARE
// marketplace-repo source — no `@<branch>` shorthand. The plugin id suffix is the
// marketplace NAME (the channel); the entry before the `@` is always `spacedock`.
// The marketplace-update refresh targets the channel's marketplace name; no step
// spells `marketplace remove`.
func TestClaudeChannelInstallArgvSequence(t *testing.T) {
	cases := []struct {
		channel         string
		devBranch       string
		wantID          string
		wantMarketplace string
	}{
		{channel: "stable", devBranch: "main", wantID: "spacedock@spacedock", wantMarketplace: "spacedock"},
		{channel: "edge", devBranch: "next", wantID: "spacedock@spacedock-edge", wantMarketplace: "spacedock-edge"},
	}
	for _, tc := range cases {
		t.Run(tc.channel, func(t *testing.T) {
			want := []installStep{
				{argv: []string{"plugin", "uninstall", tc.wantID}, tolerateExit: true},
				{argv: []string{"plugin", "uninstall", "spacedock-edge@spacedock"}, tolerateExit: true},
				{argv: []string{"plugin", "marketplace", "add", "spacedock-dev/marketplace"}},
				{argv: []string{"plugin", "marketplace", "update", tc.wantMarketplace}, tolerateExit: true},
				{argv: []string{"plugin", "install", tc.wantID}},
			}
			got := installArgvSequence("spacedock-dev/marketplace", tc.devBranch)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("installArgvSequence(%q) =\n%v\nwant\n%v", tc.devBranch, got, want)
			}
			// No @<branch> shorthand leaked onto the marketplace add.
			for _, step := range got {
				if len(step.argv) >= 4 && step.argv[1] == "marketplace" && step.argv[2] == "add" {
					if strings.Contains(step.argv[3], "@") {
						t.Errorf("%s channel marketplace add %q carries an @ shorthand; the tag pin lives in the manifest", tc.channel, step.argv[3])
					}
				}
			}
		})
	}
}

// TestCodexChannelInstallArgvSequence is AC-3's codex half: the codex install argv
// removes the sibling channel's plugin, then its own, before it adds the BARE
// marketplace-repo source (no `--ref`, since the channel is the marketplace name,
// not a branch ref) and adds the channel-correct id. Removing both channels' PLUGIN
// content (never their marketplace records — no step spells `marketplace remove`)
// keeps Codex's global `spacedock:*` skill namespace authoritative for the selected
// install.
func TestCodexChannelInstallArgvSequence(t *testing.T) {
	cases := []struct {
		channel         string
		devBranch       string
		wantID          string
		wantOtherID     string
		wantMarketplace string
	}{
		{channel: "stable", devBranch: "main", wantID: "spacedock@spacedock", wantOtherID: "spacedock@spacedock-edge", wantMarketplace: "spacedock"},
		{channel: "edge", devBranch: "next", wantID: "spacedock@spacedock-edge", wantOtherID: "spacedock@spacedock", wantMarketplace: "spacedock-edge"},
	}
	for _, tc := range cases {
		t.Run(tc.channel, func(t *testing.T) {
			want := []installStep{
				{argv: []string{"plugin", "remove", tc.wantOtherID}, tolerateExit: true},
				{argv: []string{"plugin", "remove", tc.wantID}, tolerateExit: true},
				{argv: []string{"plugin", "marketplace", "add", "spacedock-dev/marketplace"}},
				{argv: []string{"plugin", "marketplace", "upgrade", tc.wantMarketplace}, tolerateExit: true},
				{argv: []string{"plugin", "add", tc.wantID}},
			}
			got := codexInstallArgvSequence("spacedock-dev/marketplace", tc.devBranch)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("codexInstallArgvSequence(%q) =\n%v\nwant\n%v", tc.devBranch, got, want)
			}
			for _, step := range got {
				for _, a := range step.argv {
					if a == "--ref" {
						t.Errorf("%s channel codex sequence carries a --ref token; the channel is the entry name, not a branch ref", tc.channel)
					}
				}
			}
		})
	}
}

// TestInstallSequencesNeverSpellMarketplaceRemove is the regression lock this
// entity exists to add: `plugin marketplace remove` measurably CASCADE-UNINSTALLS
// every plugin installed from the removed marketplace on both hosts (probe 1) —
// on stable that marketplace also hosts unrelated co-installed plugins (subspace,
// cargento). Neither installArgvSequence (claude) nor codexInstallArgvSequence
// (codex) may ever issue a step whose argv is `plugin marketplace remove ...`, on
// either channel. A later edit that reintroduces the remove step (e.g. as a
// "simpler" re-pin) fails here before it fails live.
func TestInstallSequencesNeverSpellMarketplaceRemove(t *testing.T) {
	for _, devBranch := range []string{"main", "next"} {
		for _, step := range installArgvSequence("spacedock-dev/marketplace", devBranch) {
			if isMarketplaceRemoveStep(step.argv) {
				t.Errorf("installArgvSequence(devBranch=%q) spells marketplace remove: %v", devBranch, step.argv)
			}
		}
		for _, step := range codexInstallArgvSequence("spacedock-dev/marketplace", devBranch) {
			if isMarketplaceRemoveStep(step.argv) {
				t.Errorf("codexInstallArgvSequence(devBranch=%q) spells marketplace remove: %v", devBranch, step.argv)
			}
		}
	}
}

// isMarketplaceRemoveStep reports whether argv is exactly `plugin marketplace
// remove <name>` — the cascade-uninstall step neither install sequence may ever
// issue again (see TestInstallSequencesNeverSpellMarketplaceRemove).
func isMarketplaceRemoveStep(argv []string) bool {
	return len(argv) >= 3 && argv[0] == "plugin" && argv[1] == "marketplace" && argv[2] == "remove"
}

// TestClaudeNoPluginAutoInstallSelectsChannelEntry is AC-3's end-to-end seam: the
// claude front door, with devBranch set per channel, drives the no-plugin
// auto-install through the real runClaude path and the recorded install seam
// carries the marketplace-repo source + the devBranch the channel entry derives
// from. The entry id installed is confirmed via installArgvSequence on the
// observed seam values — so the observed values ARE the production install argv,
// never a constant grep.
func TestClaudeNoPluginAutoInstallSelectsChannelEntry(t *testing.T) {
	saved := devBranch
	defer func() { devBranch = saved }()

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
			devBranch = tc.devBranch

			fake := &fakeHost{manifest: "", manifestAfterInstall: compatibleManifest(t)} // fresh HOME: no claude plugin installed; install lands a compatible one
			var stdout, stderr bytes.Buffer
			code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("exit = %d, want 0 (no plugin → auto-install + launch) (stderr=%q)", code, stderr.String())
			}

			if len(fake.installCmds) < 3 {
				t.Fatalf("install seam recorded %v, want {host, source, devBranch}", fake.installCmds)
			}
			if got := fake.installCmds[0]; got != "claude" {
				t.Fatalf("install host = %q, want claude", got)
			}
			wantSource := channelMarketplaceSource(tc.devBranch)
			if got := fake.installCmds[1]; got != wantSource {
				t.Fatalf("install source = %q, want %q (the channel-resolved marketplace repo)", got, wantSource)
			}
			if got := fake.installCmds[2]; got != tc.devBranch {
				t.Fatalf("%s channel install devBranch = %q, want %q", tc.channel, got, tc.devBranch)
			}

			// The observed seam values reconstruct the production install argv: the
			// channel-correct entry id must be the install target.
			seq := installArgvSequence(fake.installCmds[1], fake.installCmds[2])
			if !sequenceInstallsID(seq, tc.wantID) {
				t.Fatalf("%s channel install sequence does not install %q; steps=%v", tc.channel, tc.wantID, seq)
			}
		})
	}
}

// sequenceInstallsID reports whether the claude install sequence's final pinning
// step installs the given plugin id (`plugin install <id>`).
func sequenceInstallsID(steps []installStep, id string) bool {
	for _, step := range steps {
		if len(step.argv) == 3 && step.argv[0] == "plugin" && step.argv[1] == "install" && step.argv[2] == id {
			return true
		}
	}
	return false
}

// TestCodexNoPluginAutoInstallSelectsChannelEntry is AC-3's codex front-door half:
// codex is the SOLE-knob host (no marketplace calendar key acts as a secondary
// channel signal), so devBranch alone determines the channel entry. The test
// drives the real runCodex no-plugin auto-install with devBranch set per channel
// and OBSERVES the seam values, then reconstructs the production codex install
// argv from them — the channel-correct entry id (`spacedock@spacedock` stable /
// `spacedock@spacedock-edge` edge) must be the `plugin add` target. The values are
// read off the recorded seam, never grepped from a constant.
func TestCodexNoPluginAutoInstallSelectsChannelEntry(t *testing.T) {
	saved := devBranch
	defer func() { devBranch = saved }()

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
			devBranch = tc.devBranch

			fake := &fakeHost{manifest: "", manifestAfterInstall: compatibleManifest(t)} // fresh HOME: no codex plugin installed; install lands a compatible one
			var stdout, stderr bytes.Buffer
			code := runCodex(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("exit = %d, want 0 (no plugin → auto-install + launch) (stderr=%q)", code, stderr.String())
			}

			if len(fake.installCmds) < 3 {
				t.Fatalf("install seam recorded %v, want {host, source, devBranch}", fake.installCmds)
			}
			if got := fake.installCmds[0]; got != "codex" {
				t.Fatalf("install host = %q, want codex", got)
			}
			wantSource := channelMarketplaceSource(tc.devBranch)
			if got := fake.installCmds[1]; got != wantSource {
				t.Fatalf("install source = %q, want %q (the channel-resolved marketplace repo)", got, wantSource)
			}
			if got := fake.installCmds[2]; got != tc.devBranch {
				t.Fatalf("%s channel install devBranch = %q, want %q (devBranch is the sole codex channel knob)", tc.channel, got, tc.devBranch)
			}

			seq := codexInstallArgvSequence(fake.installCmds[1], fake.installCmds[2])
			if !codexSequenceAddsID(seq, tc.wantID) {
				t.Fatalf("%s channel codex sequence does not `plugin add %q`; steps=%v", tc.channel, tc.wantID, seq)
			}
		})
	}
}

// codexSequenceAddsID reports whether the codex install sequence's final pinning
// step adds the given plugin id (`plugin add <id>`).
func codexSequenceAddsID(steps []installStep, id string) bool {
	for _, step := range steps {
		if len(step.argv) == 3 && step.argv[0] == "plugin" && step.argv[1] == "add" && step.argv[2] == id {
			return true
		}
	}
	return false
}
