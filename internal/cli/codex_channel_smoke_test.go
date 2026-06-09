// ABOUTME: AC-a/AC-d — per-channel devBranch drives the claude/codex no-plugin
// ABOUTME: auto-install to @main/@next (claude) and --ref main/next (codex) in argv.
package cli

import (
	"bytes"
	"context"
	"testing"
)

// TestClaudeNoPluginAutoInstallChannelRef is AC-a's hermetic half: the claude
// front door, with devBranch set per channel, drives the no-plugin auto-install
// to the channel-correct marketplace `@ref`. It mirrors the codex test but in
// claude's `source@branch` shorthand: a stable binary (devBranch=main) resolves
// `spacedock-dev/spacedock@main`; an edge binary (devBranch=next) resolves
// `…@next`. The branch is observed off the recorded install seam, and
// marketplaceAddArg is confirmed to compose `source@branch` from it — so the
// observed value IS the production marketplace-add argv, never a constant grep.
// (The live built-binary smoke proves the BUILD stamp end-to-end; this test pins
// the channel→argv contract hermetically.)
func TestClaudeNoPluginAutoInstallChannelRef(t *testing.T) {
	saved := devBranch
	defer func() { devBranch = saved }()

	cases := []struct {
		channel string
		branch  string
	}{
		{channel: "stable", branch: "main"},
		{channel: "edge", branch: "next"},
	}
	for _, tc := range cases {
		t.Run(tc.channel, func(t *testing.T) {
			devBranch = tc.branch

			fake := &fakeHost{manifest: ""} // fresh HOME: no claude plugin installed
			var stdout, stderr bytes.Buffer
			code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("exit = %d, want 0 (no plugin → auto-install + launch) (stderr=%q)", code, stderr.String())
			}

			if len(fake.installCmds) < 3 {
				t.Fatalf("install seam recorded %v, want {host, source, branch}", fake.installCmds)
			}
			if got := fake.installCmds[0]; got != "claude" {
				t.Fatalf("install host = %q, want claude", got)
			}
			if got := fake.installCmds[2]; got != tc.branch {
				t.Fatalf("%s channel install branch = %q, want %q", tc.channel, got, tc.branch)
			}

			wantRef := marketplaceSource + "@" + tc.branch
			if got := marketplaceAddArg(marketplaceSource, fake.installCmds[2]); got != wantRef {
				t.Fatalf("%s channel marketplace add arg = %q, want %q", tc.channel, got, wantRef)
			}
		})
	}
}

// TestCodexNoPluginAutoInstallChannelRef is AC-d's hermetic proof: the codex
// front door is the SOLE-knob host (no .codex-plugin/marketplace.json calendar
// key acts as a secondary channel signal — claude has one, codex does not), so
// the `devBranch` ldflag is the only determinant of the codex `--ref`. The test
// drives the real runCodex no-plugin auto-install with devBranch set per channel
// and OBSERVES the branch the install seam records — exactly the value
// codexInstallArgvSequence threads into codex's `--ref <branch>` flag (asserted
// below). A stable binary (devBranch=main) must resolve `--ref main`; an edge
// binary (devBranch=next) must resolve `--ref next`. The branch is read off the
// recorded seam interaction, never grepped from the constant.
func TestCodexNoPluginAutoInstallChannelRef(t *testing.T) {
	saved := devBranch
	defer func() { devBranch = saved }()

	cases := []struct {
		channel string
		branch  string
	}{
		{channel: "stable", branch: "main"},
		{channel: "edge", branch: "next"},
	}
	for _, tc := range cases {
		t.Run(tc.channel, func(t *testing.T) {
			devBranch = tc.branch // the per-channel ldflag the released binary carries

			fake := &fakeHost{manifest: ""} // fresh HOME: no codex plugin installed
			var stdout, stderr bytes.Buffer
			code := runCodex(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("exit = %d, want 0 (no plugin → auto-install + launch) (stderr=%q)", code, stderr.String())
			}

			// The seam records {host, source, branch}; branch is the channel ref.
			if len(fake.installCmds) < 3 {
				t.Fatalf("install seam recorded %v, want {host, source, branch}", fake.installCmds)
			}
			if got := fake.installCmds[0]; got != "codex" {
				t.Fatalf("install host = %q, want codex", got)
			}
			if got := fake.installCmds[2]; got != tc.branch {
				t.Fatalf("%s channel install branch = %q, want %q (devBranch is the sole codex channel knob)", tc.channel, got, tc.branch)
			}

			// The recorded branch is exactly what codex's `--ref` carries: confirm the
			// argv composition threads it through `--ref <branch>` on the marketplace
			// add step, so the observed seam value IS the production install argv.
			seq := codexInstallArgvSequence(marketplaceSource, fake.installCmds[2])
			if !codexMarketplaceAddHasRef(seq, tc.branch) {
				t.Fatalf("%s channel codexInstallArgvSequence has no `marketplace add … --ref %s`; steps=%v", tc.channel, tc.branch, seq)
			}
		})
	}
}

// codexMarketplaceAddHasRef reports whether the install sequence's marketplace
// add step carries `--ref <branch>` as adjacent argv tokens — the codex
// branch-pinning form (a separate flag, not the claude `source@branch`
// shorthand). It reads the argv structurally so a reordering or a dropped `--ref`
// reds AC-d rather than passing on a substring coincidence.
func codexMarketplaceAddHasRef(steps []installStep, branch string) bool {
	for _, step := range steps {
		argv := step.argv
		if len(argv) < 4 || argv[0] != "plugin" || argv[1] != "marketplace" || argv[2] != "add" {
			continue
		}
		for i := 3; i+1 < len(argv); i++ {
			if argv[i] == "--ref" && argv[i+1] == branch {
				return true
			}
		}
	}
	return false
}
