// ABOUTME: AC-3 channel-resolution seam — devBranch selects the marketplace ENTRY
// ABOUTME: NAME (spacedock vs spacedock-edge) from the marketplace repo, no @ref shorthand.
package cli

import (
	"bytes"
	"context"
	"reflect"
	"strings"
	"testing"
)

// TestChannelEntryFromDevBranch locks the channel-selection rule under Model B:
// the binary's devBranch stamp selects the marketplace ENTRY NAME, not a git ref
// pinned into the plugin repo. A stable binary (devBranch=main) resolves the
// `spacedock` entry; an edge binary (devBranch=next) resolves `spacedock-edge`.
// The tag pin lives in the marketplace manifest, so the channel is the entry, not
// an @branch shorthand on the install command.
func TestChannelEntryFromDevBranch(t *testing.T) {
	cases := []struct {
		channel   string
		devBranch string
		wantEntry string
	}{
		{channel: "stable", devBranch: "main", wantEntry: "spacedock"},
		{channel: "edge", devBranch: "next", wantEntry: "spacedock-edge"},
	}
	for _, tc := range cases {
		t.Run(tc.channel, func(t *testing.T) {
			if got := channelEntry(tc.devBranch); got != tc.wantEntry {
				t.Errorf("channelEntry(%q) = %q, want %q (the %s channel)", tc.devBranch, got, tc.wantEntry, tc.channel)
			}
		})
	}
}

// TestClaudeChannelInstallArgvSequence is AC-3's claude half: with marketplaceSource
// repointed to the marketplace repo and devBranch set per channel, the issued
// claude install argv installs the channel-correct ENTRY id (`spacedock@spacedock`
// stable / `spacedock-edge@spacedock` edge) and the marketplace add carries the
// BARE marketplace-repo source — no `@<branch>` shorthand. The plugin id (`@spacedock`
// suffix) is the marketplace NAME; the entry before the `@` is the channel.
func TestClaudeChannelInstallArgvSequence(t *testing.T) {
	cases := []struct {
		channel   string
		devBranch string
		wantID    string
	}{
		{channel: "stable", devBranch: "main", wantID: "spacedock@spacedock"},
		{channel: "edge", devBranch: "next", wantID: "spacedock-edge@spacedock"},
	}
	for _, tc := range cases {
		t.Run(tc.channel, func(t *testing.T) {
			want := []installStep{
				{argv: []string{"plugin", "uninstall", tc.wantID}, tolerateExit: true},
				{argv: []string{"plugin", "marketplace", "remove", "spacedock"}, tolerateExit: true},
				{argv: []string{"plugin", "marketplace", "add", "spacedock-dev/marketplace"}},
				{argv: []string{"plugin", "install", tc.wantID}},
			}
			got := installArgvSequence("spacedock-dev/marketplace", tc.devBranch)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("installArgvSequence(%q) =\n%v\nwant\n%v", tc.devBranch, got, want)
			}
			// No @<branch> shorthand leaked onto the marketplace add.
			for _, step := range got {
				if len(step.argv) >= 3 && step.argv[1] == "marketplace" && step.argv[2] == "add" {
					if strings.Contains(step.argv[3], "@") {
						t.Errorf("%s channel marketplace add %q carries an @ shorthand; the tag pin lives in the manifest", tc.channel, step.argv[3])
					}
				}
			}
		})
	}
}

// TestCodexChannelInstallArgvSequence is AC-3's codex half: the codex install argv
// adds the BARE marketplace-repo source (no `--ref`, since the channel is the entry
// name, not a branch ref) and adds the channel-correct entry id.
func TestCodexChannelInstallArgvSequence(t *testing.T) {
	cases := []struct {
		channel   string
		devBranch string
		wantID    string
	}{
		{channel: "stable", devBranch: "main", wantID: "spacedock@spacedock"},
		{channel: "edge", devBranch: "next", wantID: "spacedock-edge@spacedock"},
	}
	for _, tc := range cases {
		t.Run(tc.channel, func(t *testing.T) {
			want := []installStep{
				{argv: []string{"plugin", "remove", tc.wantID}, tolerateExit: true},
				{argv: []string{"plugin", "marketplace", "remove", "spacedock"}, tolerateExit: true},
				{argv: []string{"plugin", "marketplace", "add", "spacedock-dev/marketplace"}},
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
		{channel: "edge", devBranch: "next", wantID: "spacedock-edge@spacedock"},
	}
	for _, tc := range cases {
		t.Run(tc.channel, func(t *testing.T) {
			devBranch = tc.devBranch

			fake := &fakeHost{manifest: ""} // fresh HOME: no claude plugin installed
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
			if got := fake.installCmds[1]; got != marketplaceSource {
				t.Fatalf("install source = %q, want %q (the marketplace repo)", got, marketplaceSource)
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
// `spacedock-edge@spacedock` edge) must be the `plugin add` target. The values are
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
		{channel: "edge", devBranch: "next", wantID: "spacedock-edge@spacedock"},
	}
	for _, tc := range cases {
		t.Run(tc.channel, func(t *testing.T) {
			devBranch = tc.devBranch

			fake := &fakeHost{manifest: ""} // fresh HOME: no codex plugin installed
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
			if got := fake.installCmds[1]; got != marketplaceSource {
				t.Fatalf("install source = %q, want %q (the marketplace repo)", got, marketplaceSource)
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
