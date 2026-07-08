// ABOUTME: AC-4 init seam tests — claude install issues host plugin commands
// ABOUTME: (not a file copy); codex emits the documented add command pair as prose.
package cli

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
)

// TestInitClaudeIssuesHostPluginCommands: `spacedock install --host claude` drives
// the install seam (host plugin marketplace add + install) rather than a
// filesystem copy of skill files.
func TestInitClaudeIssuesHostPluginCommands(t *testing.T) {
	fake := &fakeHost{
		manifest:   compatibleManifest(t), // doctor after install sees a compatible plugin
		installOut: "installed spacedock@spacedock",
	}
	var stdout, stderr bytes.Buffer

	code := runInit(context.Background(), []string{"--host", "claude"}, fake, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if len(fake.installCmds) == 0 {
		t.Fatalf("install seam not invoked — init must use the host plugin mechanism, not a file copy")
	}
	// The seam was called with host=claude.
	if fake.installCmds[0] != "claude" {
		t.Fatalf("install seam host = %q, want claude", fake.installCmds[0])
	}
	// After install, init runs doctor — a compatible report on stdout.
	if !strings.Contains(stdout.String(), "OK") {
		t.Fatalf("init should run doctor after install; stdout = %q", stdout.String())
	}
}

// TestInitMarketplaceSourceIsMarketplaceRepo guards the Model B decouple: the
// marketplace-add target is the standalone marketplace repo
// `spacedock-dev/marketplace` (channel-resolved by channelMarketplaceSource — the
// edge binary adds the `@edge` branch ref), NOT the plugin repo
// `spacedock-dev/spacedock` (the manifest moved out of the plugin branch). Without
// this, a silent revert of the marketplaceSource back to the plugin repo would not
// fail `go test` — both hosts' install seams carry the source, so both paths are
// asserted.
func TestInitMarketplaceSourceIsMarketplaceRepo(t *testing.T) {
	wantSource := channelMarketplaceSource(devBranch)

	t.Run("claude-install-seam", func(t *testing.T) {
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer

		code := runInit(context.Background(), []string{"--host", "claude"}, fake, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		// Install records {host, source, branch}; the source is the marketplace target.
		if len(fake.installCmds) < 2 {
			t.Fatalf("install seam recorded %v, want at least {host, source}", fake.installCmds)
		}
		if got := fake.installCmds[1]; got != wantSource {
			t.Fatalf("claude marketplace source = %q, want %q (pre-migration clkao/spacedock must not return)", got, wantSource)
		}
	})

	t.Run("codex-install-seam", func(t *testing.T) {
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer

		code := runInit(context.Background(), []string{"--host", "codex"}, fake, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		if len(fake.installCmds) < 2 {
			t.Fatalf("codex install seam recorded %v, want at least {host, source}", fake.installCmds)
		}
		if got := fake.installCmds[1]; got != wantSource {
			t.Fatalf("codex marketplace source = %q, want %q (pre-migration plugin repo must not return)", got, wantSource)
		}
	})
}

// TestInitCheckRunsDoctorWithoutInstalling: `--check` runs the compatibility
// report without invoking the install seam.
func TestInitCheckRunsDoctorWithoutInstalling(t *testing.T) {
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runInit(context.Background(), []string{"--host", "claude", "--check"}, fake, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if len(fake.installCmds) != 0 {
		t.Fatalf("--check must not install: %v", fake.installCmds)
	}
	if !strings.Contains(stdout.String(), "OK") {
		t.Fatalf("--check should print the doctor report; stdout = %q", stdout.String())
	}
}

func TestInitCodexInstallReadiness(t *testing.T) {
	// compatible-installed: `spacedock install --host codex` re-installs an
	// already-present compatible plugin instead of short-circuiting to a
	// doctor-only no-op. The codex arm must drive the install seam exactly like
	// the claude arm — the recorded {host, source, branch} call is the
	// independent source of truth that the refresh actually fired, not the
	// no-op the prior assertion codified.
	t.Run("compatible-installed", func(t *testing.T) {
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer

		code := runInit(context.Background(), []string{"--host", "codex"}, fake, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		wantInstall := []string{"codex", channelMarketplaceSource(devBranch), devBranch}
		if !equalArgv(fake.installCmds, wantInstall) {
			t.Fatalf("install seam = %v, want %v — codex init on a present plugin must refresh, not no-op", fake.installCmds, wantInstall)
		}
		// After install, init runs doctor — a compatible report on stdout.
		major, minor := binaryMinor(t)
		out := stdout.String()
		wantOK := fmt.Sprintf("OK: spacedock binary %s and plugin %d.%d.8", displayVersion(), major, minor)
		if !strings.Contains(out, wantOK) {
			t.Fatalf("codex init should run doctor after install and report compatible; stdout = %q, want substring %q", out, wantOK)
		}
	})

	// compatible-installed-check: `--check` keeps the no-install report — it runs
	// doctor without driving the install seam.
	t.Run("compatible-installed-check", func(t *testing.T) {
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer

		code := runInit(context.Background(), []string{"--host", "codex", "--check"}, fake, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		if len(fake.installCmds) != 0 {
			t.Fatalf("--check must not install: %v", fake.installCmds)
		}
		if !strings.Contains(stdout.String(), "OK") {
			t.Fatalf("--check should print the doctor report; stdout = %q", stdout.String())
		}
	})

	// not-installed: AC-3 — on a fresh box (no plugin resolves) the codex arm RUNS
	// the install seam (marketplace add + plugin add, re-pinning the source), like
	// the claude arm, instead of printing manual prose. The recorded {host, source,
	// branch} call is the source of truth; no "Run these in your shell" prose.
	t.Run("not-installed", func(t *testing.T) {
		saved := devBranch
		devBranch = "main"
		defer func() { devBranch = saved }()

		fake := &fakeHost{}
		var stdout, stderr bytes.Buffer

		code := runInit(context.Background(), []string{"--host", "codex"}, fake, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		wantInstall := []string{"codex", marketplaceSource, "main"}
		if !equalArgv(fake.installCmds, wantInstall) {
			t.Fatalf("install seam = %v, want %v — fresh-box codex install must run the seam, not print prose", fake.installCmds, wantInstall)
		}
		if out := stdout.String(); strings.Contains(out, "Run these in your shell") {
			t.Errorf("fresh-box codex install must not fall back to manual prose:\n%s", out)
		}
	})

	// not-installed-edge: an edge binary (devBranch=next) drives the install seam
	// with devBranch=next, so the channel selection reaches the codex install — the
	// seam then targets the `spacedock@spacedock-edge` id (the edge marketplace).
	t.Run("not-installed-edge", func(t *testing.T) {
		saved := devBranch
		devBranch = "next"
		defer func() { devBranch = saved }()

		fake := &fakeHost{}
		var stdout, stderr bytes.Buffer

		code := runInit(context.Background(), []string{"--host", "codex"}, fake, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		wantInstall := []string{"codex", channelMarketplaceSource("next"), "next"}
		if !equalArgv(fake.installCmds, wantInstall) {
			t.Fatalf("edge codex install seam = %v, want %v", fake.installCmds, wantInstall)
		}
		// The observed seam values reconstruct the production codex install argv: the
		// edge channel id must be the `plugin add` target.
		seq := codexInstallArgvSequence(fake.installCmds[1], fake.installCmds[2])
		if !codexSequenceAddsID(seq, "spacedock@spacedock-edge") {
			t.Fatalf("edge codex install sequence does not `plugin add spacedock@spacedock-edge`; steps=%v", seq)
		}
	})

	t.Run("incompatible-installed", func(t *testing.T) {
		fake := &fakeHost{manifest: tooOldBinaryManifest(t)}
		var stdout, stderr bytes.Buffer

		code := runInit(context.Background(), []string{"--host", "codex"}, fake, &stdout, &stderr)

		if code == 0 {
			t.Fatalf("exit = 0, want non-zero for incompatible installed plugin")
		}
		errOut := stderr.String()
		major, minor := binaryMinor(t)
		wantMismatch := fmt.Sprintf("Spacedock version mismatch: binary %s, plugin %d.%d.0", displayVersion(), major, minor+1)
		if !strings.Contains(errOut, wantMismatch) {
			t.Fatalf("codex init should surface doctor mismatch; stderr = %q, want substring %q", errOut, wantMismatch)
		}
		for _, banned := range []string{"codex plugin marketplace add", "codex plugin add spacedock@spacedock"} {
			if strings.Contains(stdout.String(), banned) || strings.Contains(errOut, banned) {
				t.Errorf("incompatible codex init must not imply absent plugin via %q\nstdout=%s\nstderr=%s", banned, stdout.String(), errOut)
			}
		}
	})
}
