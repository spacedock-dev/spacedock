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

// codexPluginDirHost inspects the local marketplace's staged plugin dir at Install
// time — the moment the real codex host reads it — to confirm the checkout's plugin
// surface was staged (copied, not symlinked — see WriteCodexLocalMarketplace) into
// the marketplace the install consumes. The channel-NAME property (that an edge
// build names the marketplace `spacedock-edge` so it resolves) is proven
// behaviorally against real codex in AC-2, not by re-reading the JSON here.
type codexPluginDirHost struct {
	fakeHost
	installedIsRealDir  bool
	installedManifest   []byte
	installedPluginPath string
	inspectErr          error
}

// Install is overridden too (some callers still exercise the channel path
// through a codexPluginDirHost), but the --plugin-dir install seam is
// InstallCodexLocalPluginDir below — that is what installCodexLocalPluginDir
// actually calls since the spacedock-local marketplace name change.
func (h *codexPluginDirHost) Install(host, source, branch string) (string, error) {
	h.inspect(source)
	return h.fakeHost.Install(host, source, branch)
}

func (h *codexPluginDirHost) InstallCodexLocalPluginDir(source string) (string, error) {
	h.inspect(source)
	return h.fakeHost.InstallCodexLocalPluginDir(source)
}

// inspect reads the local marketplace's staged plugin dir at Install time — the
// moment the real codex host reads it — recording whether it is a real
// directory (not a symlink) and its manifest content.
func (h *codexPluginDirHost) inspect(source string) {
	pluginPath := filepath.Join(source, "plugins", "spacedock")
	h.installedPluginPath = pluginPath
	if info, err := os.Lstat(pluginPath); err == nil {
		h.installedIsRealDir = info.IsDir() && info.Mode()&os.ModeSymlink == 0
	} else {
		h.inspectErr = err
	}
	if data, err := os.ReadFile(filepath.Join(pluginPath, ".codex-plugin", "plugin.json")); err == nil {
		h.installedManifest = data
	} else if h.inspectErr == nil {
		h.inspectErr = err
	}
}

// TestRunCodexPluginDirInstallsThenLaunchesWithoutTheFlag is AC-1: `spacedock codex
// --plugin-dir <checkout>` installs the local checkout in one command (no
// operator-authored marketplace) and forwards NO --plugin-dir token into the real
// codex argv. The no-flag assertion is the direct regression guard against Spike A's
// reproduced baseline (today's argv DOES carry it and the real codex rejects it): a
// re-introduction of the passthrough bug flips this assertion.
func TestRunCodexPluginDirInstallsThenLaunchesWithoutTheFlag(t *testing.T) {
	t.Setenv("CODEX_HOME", t.TempDir()) // isolate the persistent local marketplace
	checkout, _ := localPluginCheckout(t, "codex")
	host := &codexPluginDirHost{fakeHost: fakeHost{manifest: compatibleManifest(t)}}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), []string{"--plugin-dir", checkout}, t.TempDir(), host, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if host.inspectErr != nil {
		t.Fatalf("marketplace inspection at Install time failed: %v", host.inspectErr)
	}
	// (a) InstallCodexLocalPluginDir called exactly once, for codex, with a
	// plugins/spacedock staged copy carrying the checkout's manifest. (The
	// dedicated-marketplace-name property is AC-2's behavioral proof against
	// real codex.)
	if len(host.installCmds) != 2 || host.installCmds[0] != "codex" {
		t.Fatalf("install seam = %v, want exactly one {codex, <marketplace>} call", host.installCmds)
	}
	if !host.installedIsRealDir {
		t.Fatalf("plugins/spacedock must be a real staged directory, not a symlink to the checkout")
	}
	wantManifest, err := os.ReadFile(filepath.Join(checkout, ".codex-plugin", "plugin.json"))
	if err != nil {
		t.Fatalf("read checkout manifest: %v", err)
	}
	if !bytes.Equal(host.installedManifest, wantManifest) {
		t.Fatalf("staged plugin manifest = %q, want the checkout's manifest %q", host.installedManifest, wantManifest)
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

// A pre-fence plugin-dir is consumed before the exact-resume check, so neither
// the checkout flag nor its value reaches the real Codex argv.
func TestRunCodexPluginDirPostFenceResumeStaysPromptFree(t *testing.T) {
	t.Setenv("CODEX_HOME", t.TempDir()) // isolate the persistent local marketplace
	checkout, _ := localPluginCheckout(t, "codex")
	host := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), []string{"--plugin-dir", checkout, "--", "--model", "gpt-5.6-sol", "resume", "abc123"}, t.TempDir(), host, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := wantCodexArgv("--model", "gpt-5.6-sol", "resume", "abc123")
	if !equalArgv(host.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", host.launchedArg, want)
	}
	if strings.Contains(stderr.String(), "· launching codex") {
		t.Fatalf("exact resume produced a fresh-launch banner: %q", stderr.String())
	}
}

// The pre-fence plugin-dir consumption seam composes with model-only Codex
// launches: the host receives the model pair and normal bootstrap posture, but
// never the Spacedock-owned checkout tokens.
func TestRunCodexPluginDirModelOnlyRetainsBootstrapPosture(t *testing.T) {
	t.Setenv("CODEX_HOME", t.TempDir()) // isolate the persistent local marketplace
	checkout, _ := localPluginCheckout(t, "codex")
	host := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	args := []string{"--plugin-dir", checkout, "--", "--model", "gpt-5.6-sol"}
	code := runCodex(context.Background(), args, t.TempDir(), host, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := wantCodexArgv("--ask-for-approval", "on-request", "--model", "gpt-5.6-sol", wantCodexBootstrapPrompt)
	if !equalArgv(host.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", host.launchedArg, want)
	}
	if !strings.Contains(stderr.String(), "· launching codex") {
		t.Fatalf("model-only launch omitted the Codex banner: %q", stderr.String())
	}
}

func TestRunCodexPostFencePluginDirIsRejectedTruthfully(t *testing.T) {
	t.Setenv("CODEX_HOME", t.TempDir())
	postFenceCheckout := t.TempDir()
	host := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	args := []string{"--", "--plugin-dir", postFenceCheckout}
	code := runCodex(context.Background(), args, t.TempDir(), host, lookFound, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("exit = 0, want unsupported forwarded flag failure (stderr=%q)", stderr.String())
	}
	if len(host.installCmds) != 0 {
		t.Fatalf("post-fence --plugin-dir invoked the local-install seam: %v", host.installCmds)
	}
	if host.launchedArg != nil {
		t.Fatalf("post-fence --plugin-dir reached Codex argv: %v", host.launchedArg)
	}
	if !strings.Contains(stderr.String(), "does not accept forwarded --plugin-dir") || !strings.Contains(stderr.String(), "codex plugin add") {
		t.Fatalf("post-fence diagnostic is not actionable: %q", stderr.String())
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
		checkout, _ := localPluginCheckout(t, "codex")
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer
		code := runCodex(context.Background(), []string{"--plugin-dir", checkout}, t.TempDir(), fake, lookFound, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		if !strings.Contains(stderr.String(), advisory) {
			t.Fatalf("stderr missing the version-masquerade advisory: %q", stderr.String())
		}
		if !strings.Contains(stderr.String(), "as "+codexLocalPluginID) {
			t.Fatalf("advisory must name the selected codex plugin id: %q", stderr.String())
		}
		if !strings.Contains(stderr.String(), "Removed other Spacedock Codex channels") {
			t.Fatalf("advisory must name the sibling-provider cleanup: %q", stderr.String())
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
	checkout, _ := localPluginCheckout(t, "codex")
	host := &codexPluginDirHost{fakeHost: fakeHost{manifest: compatibleManifest(t)}}
	var stdout, stderr bytes.Buffer

	code := runInitWithPi(context.Background(), []string{"--host", "codex", "--plugin-dir", checkout}, host, &fakePiRuntimeOps{}, nil, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if len(host.installCmds) != 2 || host.installCmds[0] != "codex" {
		t.Fatalf("install seam = %v, want exactly one codex install", host.installCmds)
	}
	if !host.installedIsRealDir {
		t.Fatalf("plugins/spacedock must be a real staged directory, not a symlink to the checkout")
	}
	wantManifest, err := os.ReadFile(filepath.Join(checkout, ".codex-plugin", "plugin.json"))
	if err != nil {
		t.Fatalf("read checkout manifest: %v", err)
	}
	if !bytes.Equal(host.installedManifest, wantManifest) {
		t.Fatalf("staged plugin manifest = %q, want the checkout's manifest %q", host.installedManifest, wantManifest)
	}
	if !strings.Contains(stderr.String(), "version-masquerade advisory") {
		t.Fatalf("install --host codex --plugin-dir missing the advisory: %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "as "+codexLocalPluginID) {
		t.Fatalf("install --host codex --plugin-dir must name the selected plugin id: %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Removed other Spacedock Codex channels") {
		t.Fatalf("install --host codex --plugin-dir must name the sibling-provider cleanup: %q", stderr.String())
	}
}

// TestInstallCodexLocalPluginDirLeavesOnePromptInputProvider is AC-1's hermetic
// host smoke. It seeds a fresh CODEX_HOME with BOTH Spacedock channels using raw
// codex commands, then runs the production --plugin-dir install helper for the
// selected edge checkout. `codex debug prompt-input` is Codex's own model-visible
// skill list renderer, so the assertion observes provider resolution without an
// LLM call: exactly one `spacedock:first-officer` entry remains, and it is the
// selected checkout's provider.
func TestInstallCodexLocalPluginDirLeavesOnePromptInputProvider(t *testing.T) {
	codexBin, err := exec.LookPath("codex")
	if err != nil {
		t.Skip("codex not on PATH; prompt-input provider smoke requires the host CLI")
	}
	saved := devBranch
	devBranch = "next"
	defer func() { devBranch = saved }()

	tmp := t.TempDir()
	codexHomeDir := filepath.Join(tmp, "codexhome")
	mustMkdir(t, codexHomeDir)
	t.Setenv("CODEX_HOME", codexHomeDir)

	stableCheckout := buildCodexFirstOfficerCheckout(t, filepath.Join(tmp, "stable-checkout"), "1.0.0", "STABLE_PROVIDER_SPIKE")
	oldEdgeCheckout := buildCodexFirstOfficerCheckout(t, filepath.Join(tmp, "old-edge-checkout"), "2.0.0", "OLD_EDGE_PROVIDER_SPIKE")
	selectedCheckout := buildCodexFirstOfficerCheckout(t, filepath.Join(tmp, "selected-checkout"), "2.0.1", "SELECTED_PROVIDER_SPIKE")

	stableInstall, err := WriteCodexLocalMarketplace(filepath.Join(tmp, "stable-marketplace"), stableCheckout, "spacedock")
	if err != nil {
		t.Fatalf("build stable local marketplace: %v", err)
	}
	oldEdgeInstall, err := WriteCodexLocalMarketplace(filepath.Join(tmp, "old-edge-marketplace"), oldEdgeCheckout, "spacedock-edge")
	if err != nil {
		t.Fatalf("build old edge local marketplace: %v", err)
	}

	runHost(t, codexBin, os.Environ(), "plugin", "marketplace", "add", stableInstall.MarketplaceRoot)
	runHost(t, codexBin, os.Environ(), "plugin", "add", "spacedock@spacedock")
	runHost(t, codexBin, os.Environ(), "plugin", "marketplace", "add", oldEdgeInstall.MarketplaceRoot)
	runHost(t, codexBin, os.Environ(), "plugin", "add", "spacedock@spacedock-edge")

	if err := installCodexLocalPluginDir(execHost{}, selectedCheckout, io.Discard); err != nil {
		t.Fatalf("install selected --plugin-dir checkout: %v", err)
	}

	cmd := exec.Command(codexBin, "debug", "prompt-input", "probe")
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("codex debug prompt-input failed: %v\n%s", err, out)
	}
	rendered := string(out)
	if got := strings.Count(rendered, "spacedock:first-officer:"); got != 1 {
		t.Fatalf("prompt input has %d spacedock:first-officer providers, want exactly 1\n%s", got, rendered)
	}
	if strings.Contains(rendered, "STABLE_PROVIDER_SPIKE") || strings.Contains(rendered, "OLD_EDGE_PROVIDER_SPIKE") {
		t.Fatalf("prompt input still exposes stale Spacedock providers:\n%s", rendered)
	}
	if !strings.Contains(rendered, "SELECTED_PROVIDER_SPIKE") {
		t.Fatalf("prompt input missing the selected checkout provider:\n%s", rendered)
	}
}

// TestInstallCodexLocalPluginDirResolvesOnEdgeChannel is AC-2, round-3 shape: a
// --plugin-dir codex install names its marketplace via the FIXED
// codexLocalMarketplaceName ("spacedock-local"), devBranch-INDEPENDENT — so an
// edge-devBranch build's install still resolves through the real
// ResolveManifest, which must check the local id (codexLocalPluginID)
// regardless of devBranch (team-lead's "teach the gate the second id" point).
// Without that check, ResolveManifest would only look for
// channelPluginID(devBranch) and never find the local install — the same
// empty-resolve failure mode the original (pre-round-3) Spike E baseline hit
// for a different reason (a hardcoded `spacedock` name). Pins devBranch="next"
// (edge), points CODEX_HOME at a fresh temp dir, installs a throwaway checkout
// fixture through the real helper, then asserts ResolveManifest("codex") is
// non-empty. Skips when codex is absent (no auth required — Spike C).
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
// (a .codex-plugin/plugin.json with a version + one skill), the shape a
// --plugin-dir target must have for codex's `plugin add` to copy it into its cache.
func buildCodexPluginCheckout(t *testing.T, root, version string) string {
	t.Helper()
	mustMkdir(t, filepath.Join(root, ".codex-plugin"))
	mustMkdir(t, filepath.Join(root, "skills", "demo"))
	mustWrite(t, filepath.Join(root, ".codex-plugin", "plugin.json"),
		`{ "name": "spacedock", "version": "`+version+`", "skills": "./skills/" }
`)
	mustWrite(t, filepath.Join(root, "skills", "demo", "SKILL.md"), "---\nname: demo\ndescription: demo skill\n---\ndemo\n")
	return root
}

func buildCodexFirstOfficerCheckout(t *testing.T, root, version, description string) string {
	t.Helper()
	mustMkdir(t, filepath.Join(root, ".codex-plugin"))
	mustMkdir(t, filepath.Join(root, "skills", "first-officer"))
	mustWrite(t, filepath.Join(root, ".codex-plugin", "plugin.json"),
		`{ "name": "spacedock", "version": "`+version+`", "skills": "./skills/" }
`)
	mustWrite(t, filepath.Join(root, "skills", "first-officer", "SKILL.md"), "---\nname: first-officer\ndescription: "+description+"\n---\n"+description+"\n")
	return root
}
