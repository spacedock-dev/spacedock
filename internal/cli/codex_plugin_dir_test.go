// ABOUTME: `spacedock codex --plugin-dir` / `install --host codex --plugin-dir`
// ABOUTME: tests — no-flag-passthrough regression guard, edge-channel resolve, advisory.
package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// codexPluginDirHost inspects the local marketplace's plugin symlink at Install
// time — the moment the real codex host reads it — to confirm the checkout is wired
// into the marketplace the install consumes. The channel-NAME property (that an edge
// build names the marketplace `spacedock-edge` so it resolves) is proven
// behaviorally against real codex in AC-2, not by re-reading the JSON here.
type codexPluginDirHost struct {
	fakeHost
	installedSymlinkDest string
	inspectErr           error
}

func (h *codexPluginDirHost) Install(host, source, branch string) (string, error) {
	if dest, err := os.Readlink(filepath.Join(source, "plugins", "spacedock")); err == nil {
		h.installedSymlinkDest = dest
	} else {
		h.inspectErr = err
	}
	return h.fakeHost.Install(host, source, branch)
}

// TestRunCodexPluginDirInstallsThenLaunchesWithoutTheFlag is AC-1: `spacedock codex
// --plugin-dir <checkout>` installs the local checkout in one command (no
// operator-authored marketplace) and forwards NO --plugin-dir token into the real
// codex argv. The no-flag assertion is the direct regression guard against Spike A's
// reproduced baseline (today's argv DOES carry it and the real codex rejects it): a
// re-introduction of the passthrough bug flips this assertion.
func TestRunCodexPluginDirInstallsThenLaunchesWithoutTheFlag(t *testing.T) {
	t.Setenv("CODEX_HOME", t.TempDir()) // isolate the persistent local marketplace
	checkout := t.TempDir()
	host := &codexPluginDirHost{fakeHost: fakeHost{manifest: compatibleManifest(t)}}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), []string{"--plugin-dir", checkout}, t.TempDir(), host, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if host.inspectErr != nil {
		t.Fatalf("marketplace inspection at Install time failed: %v", host.inspectErr)
	}
	// (a) Install called exactly once, for codex, with a plugins/spacedock symlink
	// resolving to the checkout. (The marketplace-name-is-the-channel property is
	// AC-2's behavioral proof against real codex.)
	if len(host.installCmds) != 3 || host.installCmds[0] != "codex" || host.installCmds[2] != devBranch {
		t.Fatalf("install seam = %v, want exactly one {codex, <marketplace>, %q} call", host.installCmds, devBranch)
	}
	if host.installedSymlinkDest != checkout {
		t.Fatalf("plugins/spacedock symlink = %q, want the checkout %q", host.installedSymlinkDest, checkout)
	}
	// (b) No --plugin-dir anywhere in the launched codex argv.
	if host.launchedArg == nil {
		t.Fatalf("launch seam not reached after --plugin-dir install")
	}
	for _, a := range host.launchedArg {
		if a == "--plugin-dir" || a == checkout {
			t.Fatalf("launch argv forwards a --plugin-dir token: %v", host.launchedArg)
		}
	}
}

// TestCodexPluginDirAdvisoryPresenceAndAbsence is AC-3: every --plugin-dir codex
// install prints the version-masquerade advisory; a plain (non---plugin-dir) launch
// prints none. The pair (not a presence-only check) means the test cannot pass by
// printing the advisory unconditionally. The present subtest also guards that the
// advisory carries its meaning-bearing clause (not necessarily its current HEAD) and
// leaks no internal branch identifier.
func TestCodexPluginDirAdvisoryPresenceAndAbsence(t *testing.T) {
	const advisory = "version-masquerade advisory"

	t.Run("present on --plugin-dir", func(t *testing.T) {
		t.Setenv("CODEX_HOME", t.TempDir()) // isolate the persistent local marketplace
		checkout := t.TempDir()
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer
		code := runCodex(context.Background(), []string{"--plugin-dir", checkout}, t.TempDir(), fake, lookFound, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		if !strings.Contains(stderr.String(), advisory) {
			t.Fatalf("stderr missing the version-masquerade advisory: %q", stderr.String())
		}
		if !strings.Contains(stderr.String(), "not necessarily its current HEAD") {
			t.Fatalf("advisory lost its meaning-bearing clause: %q", stderr.String())
		}
		if strings.Contains(stderr.String(), "next-post-release-preversion-bump") {
			t.Fatalf("advisory leaks the internal branch identifier: %q", stderr.String())
		}
	})

	t.Run("absent on plain launch", func(t *testing.T) {
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer
		code := runCodex(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		if strings.Contains(stderr.String(), advisory) {
			t.Fatalf("plain codex launch printed the --plugin-dir advisory: %q", stderr.String())
		}
	})
}

// TestInstallCodexPluginDirInstallsViaSharedHelper covers the persistent primitive:
// `spacedock install --host codex --plugin-dir <checkout>` routes through the same
// shared helper (install seam called once for codex + the advisory printed), rather
// than the pre-change rejection.
func TestInstallCodexPluginDirInstallsViaSharedHelper(t *testing.T) {
	t.Setenv("CODEX_HOME", t.TempDir()) // isolate the persistent local marketplace
	checkout := t.TempDir()
	host := &codexPluginDirHost{fakeHost: fakeHost{manifest: compatibleManifest(t)}}
	var stdout, stderr bytes.Buffer

	code := runInitWithPi(context.Background(), []string{"--host", "codex", "--plugin-dir", checkout}, host, &fakePiRuntimeOps{}, nil, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if len(host.installCmds) != 3 || host.installCmds[0] != "codex" {
		t.Fatalf("install seam = %v, want exactly one codex install", host.installCmds)
	}
	if host.installedSymlinkDest != checkout {
		t.Fatalf("plugins/spacedock symlink = %q, want the checkout %q", host.installedSymlinkDest, checkout)
	}
	if !strings.Contains(stderr.String(), "version-masquerade advisory") {
		t.Fatalf("install --host codex --plugin-dir missing the advisory: %q", stderr.String())
	}
}

// TestInstallCodexLocalPluginDirResolvesOnEdgeChannel is AC-2: a --plugin-dir codex
// install names its marketplace via channelMarketplace(devBranch), so an
// edge-devBranch build's install resolves through the real ResolveManifest — where
// the OLD hardcoded `spacedock` name silently did not (Spike E baseline: empty
// resolve). Pins devBranch="next" (edge), points CODEX_HOME at a fresh temp dir,
// installs a throwaway checkout fixture through the real helper, then asserts
// ResolveManifest("codex") is non-empty. Skips when codex is absent (no auth
// required — Spike C).
func TestInstallCodexLocalPluginDirResolvesOnEdgeChannel(t *testing.T) {
	if _, err := exec.LookPath("codex"); err != nil {
		t.Skip("codex not on PATH; edge-channel resolve test requires the host CLI")
	}
	saved := devBranch
	devBranch = "next" // edge channel → marketplace must be named spacedock-edge
	defer func() { devBranch = saved }()

	tmp := t.TempDir()
	checkout := buildCodexPluginCheckout(t, filepath.Join(tmp, "checkout"), "0.0.0")
	codexHomeDir := filepath.Join(tmp, "codexhome")
	mustMkdir(t, codexHomeDir)
	t.Setenv("CODEX_HOME", codexHomeDir)

	if err := installCodexLocalPluginDir(execHost{}, checkout, io.Discard); err != nil {
		t.Fatalf("installCodexLocalPluginDir(edge) failed: %v", err)
	}

	manifest, err := execHost{}.ResolveManifest("codex")
	if err != nil {
		t.Fatalf("ResolveManifest(codex) after edge install: %v", err)
	}
	if manifest == "" {
		t.Fatalf("ResolveManifest returned empty on the edge channel — the channel-name footgun regressed to the Spike E baseline")
	}
	if _, err := os.Stat(manifest); err != nil {
		t.Fatalf("resolved edge manifest %s does not exist: %v", manifest, err)
	}
}

// buildCodexPluginCheckout writes a minimal valid codex plugin checkout under root
// (a .codex-plugin/plugin.json bracketing CONTRACT_VERSION + one skill), the shape a
// --plugin-dir target must have for codex's `plugin add` to copy it into its cache.
func buildCodexPluginCheckout(t *testing.T, root, version string) string {
	t.Helper()
	mustMkdir(t, filepath.Join(root, ".codex-plugin"))
	mustMkdir(t, filepath.Join(root, "skills", "demo"))
	mustWrite(t, filepath.Join(root, ".codex-plugin", "plugin.json"),
		`{ "name": "spacedock", "version": "`+version+`", "requires-contract": ">=2,<3", "skills": "./skills/" }
`)
	mustWrite(t, filepath.Join(root, "skills", "demo", "SKILL.md"), "---\nname: demo\ndescription: demo skill\n---\ndemo\n")
	return root
}
