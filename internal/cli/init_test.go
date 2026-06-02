// ABOUTME: AC-4 init seam tests — claude install issues host plugin commands
// ABOUTME: (not a file copy); codex emits the documented add command pair as prose.
package cli

import (
	"bytes"
	"context"
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

// TestInitMarketplaceSourceIsMigratedRepo guards the migration cleanup: the
// marketplace-add target is `spacedock-dev/spacedock`, not the pre-migration
// `clkao/spacedock`. Without this, a silent revert of the marketplaceSource
// constant would not fail `go test` — both the claude install seam and the codex
// add-prose carry the source, so both paths are asserted.
func TestInitMarketplaceSourceIsMigratedRepo(t *testing.T) {
	const wantSource = "spacedock-dev/spacedock"

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

	t.Run("codex-add-prose", func(t *testing.T) {
		fake := &fakeHost{}
		var stdout, stderr bytes.Buffer

		code := runInit(context.Background(), []string{"--host", "codex"}, fake, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		out := stdout.String()
		if !strings.Contains(out, "codex plugin marketplace add "+wantSource) {
			t.Fatalf("codex add-prose marketplace source not %q:\n%s", wantSource, out)
		}
		if strings.Contains(out, "clkao/spacedock") {
			t.Fatalf("codex add-prose still names the pre-migration clkao/spacedock:\n%s", out)
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
	t.Run("compatible-installed", func(t *testing.T) {
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer

		code := runInit(context.Background(), []string{"--host", "codex"}, fake, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		out := stdout.String()
		if !strings.Contains(out, "OK: binary contract") {
			t.Fatalf("codex init should report compatible installed plugin first; stdout = %q", out)
		}
		for _, banned := range []string{"codex plugin marketplace add", "codex plugin add spacedock@spacedock"} {
			if strings.Contains(out, banned) {
				t.Errorf("compatible codex init must not print manual add command %q:\n%s", banned, out)
			}
		}
	})

	t.Run("not-installed", func(t *testing.T) {
		fake := &fakeHost{}
		var stdout, stderr bytes.Buffer

		code := runInit(context.Background(), []string{"--host", "codex"}, fake, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		out := stdout.String()
		for _, want := range []string{
			"codex plugin marketplace add",
			"codex plugin add spacedock@spacedock",
		} {
			if !strings.Contains(out, want) {
				t.Errorf("codex init prose missing %q:\n%s", want, out)
			}
		}
		if strings.Contains(out, "codex plugin install") {
			t.Errorf("codex init prose must not use 'codex plugin install' (verb is add):\n%s", out)
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
		if !strings.Contains(errOut, "too-old-binary") {
			t.Fatalf("codex init should surface doctor mismatch; stderr = %q", errOut)
		}
		for _, banned := range []string{"codex plugin marketplace add", "codex plugin add spacedock@spacedock"} {
			if strings.Contains(stdout.String(), banned) || strings.Contains(errOut, banned) {
				t.Errorf("incompatible codex init must not imply absent plugin via %q\nstdout=%s\nstderr=%s", banned, stdout.String(), errOut)
			}
		}
	})
}
