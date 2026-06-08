// ABOUTME: AC-3/AC-4 front-door + init seam tests — version-gate fail-fast,
// ABOUTME: launch-seam argv on compatible, install-seam host commands, codex prose.
package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/contract"
)

// fakeHost records every seam interaction and returns canned results so the
// front door / init paths run with no real host CLI, no exec, no network.
type fakeHost struct {
	// manifest is the path returned by ResolveManifest; "" means no plugin found.
	manifest    string
	resolveErr  error
	launchedArg []string // argv captured by Launch
	launchedEnv []string // env captured by Launch
	launchErr   error
	installCmds []string // host commands captured by Install
	installOut  string
}

func (f *fakeHost) ResolveManifest(host string) (string, error) {
	return f.manifest, f.resolveErr
}

func (f *fakeHost) Launch(argv []string, env []string) error {
	f.launchedArg = argv
	f.launchedEnv = env
	return f.launchErr
}

func (f *fakeHost) Install(host, source, branch string) (string, error) {
	f.installCmds = append(f.installCmds, host, source, branch)
	return f.installOut, nil
}

// compatibleManifest returns a fixture path whose requires-contract brackets
// CONTRACT_VERSION (the testdata/compatible.json fixture is >=1,<2).
func compatibleManifest(t *testing.T) string {
	t.Helper()
	p, err := filepath.Abs(filepath.Join("..", "contract", "testdata", "compatible.json"))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func tooOldBinaryManifest(t *testing.T) string {
	t.Helper()
	p, err := filepath.Abs(filepath.Join("..", "contract", "testdata", "too-old-binary.json"))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func withExecutablePath(t *testing.T, path string, err error) {
	t.Helper()
	orig := executablePath
	executablePath = func() (string, error) { return path, err }
	t.Cleanup(func() { executablePath = orig })
}

func executableFixture(t *testing.T) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "spacedock")
	if err := os.WriteFile(p, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	real, err := filepath.EvalSymlinks(p)
	if err != nil {
		t.Fatal(err)
	}
	return real
}

func envValue(env []string, key string) (string, bool) {
	prefix := key + "="
	for _, entry := range env {
		if strings.HasPrefix(entry, prefix) {
			return strings.TrimPrefix(entry, prefix), true
		}
	}
	return "", false
}

// TestClaudeFrontDoorLaunchesOnCompatible: on a compatible contract the front
// door invokes the launch seam with argv beginning `claude --agent
// spacedock:first-officer` and passes through the operator's trailing args.
func TestClaudeFrontDoorLaunchesOnCompatible(t *testing.T) {
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--", "-p", "do the thing"}, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := []string{"claude", "--agent", "spacedock:first-officer", "-p", "do the thing", wantBootstrapPrompt}
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
	}
}

// TestClaudeFrontDoorFailFastOnMismatch: on a mismatch verdict the launch seam is
// NOT invoked and the process exits non-zero with the pinned remedy on stderr.
func TestClaudeFrontDoorInjectsResolvedLauncherBin(t *testing.T) {
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	t.Setenv(spacedockBinEnv, "/old/spacedock")
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	got, ok := envValue(fake.launchedEnv, spacedockBinEnv)
	if !ok || got != bin {
		t.Fatalf("%s in launch env = %q, %v; want %q, true (env=%v)", spacedockBinEnv, got, ok, bin, fake.launchedEnv)
	}
}

func TestClaudeFrontDoorOmitsStaleLauncherBinWhenResolutionFails(t *testing.T) {
	withExecutablePath(t, "", errors.New("boom"))
	t.Setenv(spacedockBinEnv, "/old/spacedock")
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if got, ok := envValue(fake.launchedEnv, spacedockBinEnv); ok {
		t.Fatalf("%s in launch env = %q, want omitted", spacedockBinEnv, got)
	}
}

func TestClaudeFrontDoorLaunchEnvResolvesSymlink(t *testing.T) {
	real := executableFixture(t)
	link := filepath.Join(t.TempDir(), "spacedock-link")
	if err := os.Symlink(real, link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	withExecutablePath(t, link, nil)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if got, ok := envValue(fake.launchedEnv, spacedockBinEnv); !ok || got != real {
		t.Fatalf("%s = %q, %v; want symlink target %q, true", spacedockBinEnv, got, ok, real)
	}
}

// TestClaudeFrontDoorFailFastOnMismatch (AC-2): a real version mismatch still
// fails fast even without --no-install — auto-install must NOT paper over an
// incompatibility. The verdict reaches runClaude's mismatch branch, not the
// no-plugin branch, so no install is invoked and no launch is reached.
func TestClaudeFrontDoorFailFastOnMismatch(t *testing.T) {
	fake := &fakeHost{manifest: tooOldBinaryManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("exit = 0, want non-zero on mismatch")
	}
	if len(fake.installCmds) != 0 {
		t.Fatalf("auto-install invoked over a mismatch: %v", fake.installCmds)
	}
	if fake.launchedArg != nil {
		t.Fatalf("launch seam invoked on mismatch: %v", fake.launchedArg)
	}
	assertGateMismatchMessage(t, stderr.String())
}

// assertGateMismatchMessage is the gate-parity oracle: the live launch gate emits
// the same version-bearing mismatch message the doctor does — it names both the
// plugin display version (0.13.0, from the too-old-binary fixture) and the binary
// display version (Version), and carries no `contract N` token in the user-facing
// line.
func assertGateMismatchMessage(t *testing.T, out string) {
	t.Helper()
	if !strings.Contains(out, "plugin 0.13.0") {
		t.Fatalf("gate mismatch message missing plugin version 0.13.0: %q", out)
	}
	if !strings.Contains(out, "binary "+Version) {
		t.Fatalf("gate mismatch message missing binary version %q: %q", Version, out)
	}
	if regexp.MustCompile(`contract \d`).MatchString(out) {
		t.Fatalf("gate mismatch message must not carry a contract-N token: %q", out)
	}
}

// TestClaudeFrontDoorNoPluginAutoInstalls (AC-1a): with no installed plugin and
// no flags the front door auto-installs the plugin then launches a working
// session — the single command the user typed yields a working FO session
// rather than refusing. Install-invocation and launch-reached are observed
// behaviors recorded by the stub, not a string match.
func TestClaudeFrontDoorNoPluginAutoInstalls(t *testing.T) {
	fake := &fakeHost{manifest: ""} // no plugin found
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 when no plugin → auto-install + launch (stderr=%q)", code, stderr.String())
	}
	if len(fake.installCmds) == 0 {
		t.Fatalf("install seam not invoked: auto-install did not run")
	}
	if fake.launchedArg == nil {
		t.Fatalf("launch seam not invoked after auto-install")
	}
}

// TestClaudeFrontDoorNoPluginNoInstallRefuses (AC-1b): --no-install preserves
// the old refuse-and-instruct behavior — no install, no launch, non-zero exit,
// with the actionable instruct message on stderr.
func TestClaudeFrontDoorNoPluginNoInstallRefuses(t *testing.T) {
	fake := &fakeHost{manifest: ""} // no plugin found
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--no-install"}, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("exit = 0, want non-zero with --no-install and no plugin")
	}
	if len(fake.installCmds) != 0 {
		t.Fatalf("install seam invoked despite --no-install: %v", fake.installCmds)
	}
	if fake.launchedArg != nil {
		t.Fatalf("launch seam invoked despite --no-install: %v", fake.launchedArg)
	}
	if !strings.Contains(stderr.String(), "spacedock install") {
		t.Fatalf("stderr missing actionable instruct message: %q", stderr.String())
	}
}

// TestClaudeFrontDoorNonEmptyMissingManifestAutoInstalls: a host that reports a
// non-empty installPath to a directory LACKING the plugin manifest is the
// NoPluginFound verdict from a phantom installPath — the plugin genuinely is not
// on disk, so the default auto-installs (a deliberate behavior change recorded in
// ideation). --no-install preserves the refuse-and-instruct arm.
func TestClaudeFrontDoorNonEmptyMissingManifestAutoInstalls(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "no-such-dir", ".claude-plugin", "plugin.json")

	t.Run("default auto-installs", func(t *testing.T) {
		fake := &fakeHost{manifest: missing} // non-empty path, but the file is absent
		var stdout, stderr bytes.Buffer

		code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

		if code != 0 {
			t.Fatalf("exit = %d, want 0 when phantom manifest → auto-install + launch (stderr=%q)", code, stderr.String())
		}
		if len(fake.installCmds) == 0 {
			t.Fatalf("install seam not invoked for a phantom (missing) manifest")
		}
		if fake.launchedArg == nil {
			t.Fatalf("launch seam not invoked after auto-install for a phantom manifest")
		}
	})

	t.Run("--no-install refuses", func(t *testing.T) {
		fake := &fakeHost{manifest: missing}
		var stdout, stderr bytes.Buffer

		code := runClaude(context.Background(), []string{"--no-install"}, t.TempDir(), fake, lookFound, &stdout, &stderr)

		if code == 0 {
			t.Fatalf("exit = 0, want non-zero with --no-install and a phantom manifest")
		}
		if len(fake.installCmds) != 0 {
			t.Fatalf("install seam invoked despite --no-install: %v", fake.installCmds)
		}
		if fake.launchedArg != nil {
			t.Fatalf("launch seam invoked despite --no-install: %v", fake.launchedArg)
		}
	})
}

// TestGateRemedyNamesLiveInstallCommand: every gateHost remedy must point at a
// command the binary actually recognizes. After the init->install rename a user
// who hits the gate and is told to "run spacedock init --host …" runs a command
// that now exits 2 (unknown command). Drive each remedy branch (resolve error,
// no plugin, missing manifest) and assert the printed remedy names `spacedock
// install` and never `spacedock init`; then prove the named command resolves by
// feeding it through cli.Run and asserting it is not the unknown-command exit 2.
func TestGateRemedyNamesLiveInstallCommand(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "no-such-dir", ".claude-plugin", "plugin.json")
	cases := []struct {
		name string
		fake *fakeHost
	}{
		{"resolve error", &fakeHost{resolveErr: errors.New("host CLI failed")}},
		{"no plugin", &fakeHost{manifest: ""}},
		{"missing manifest", &fakeHost{manifest: missing}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var stderr bytes.Buffer
			if v := gateHost(tc.fake, "claude", &stderr); v == contract.Compatible {
				t.Fatalf("gateHost = Compatible, want denied for %s", tc.name)
			}
			remedy := stderr.String()
			if !strings.Contains(remedy, "spacedock install") {
				t.Fatalf("remedy does not name the live install command: %q", remedy)
			}
			if strings.Contains(remedy, "spacedock init") {
				t.Fatalf("remedy names the removed init command (exits 2): %q", remedy)
			}
		})
	}

	// The command the remedy names must resolve in the live command tree. cobra's
	// Find returns the matched command for a registered name and falls back to the
	// root for an unknown one — so `install` must resolve to a non-root command
	// while the removed `init` must fall back to root (the unknown-command path
	// that exits 2). Resolution is deterministic and touches no host CLI.
	var stdout, stderr bytes.Buffer
	root := newRootCommand(context.Background(), nil, nil, t.TempDir(), nil, &stdout, &stderr, &fakeRunner{}, nil)
	if cmd, _, err := root.Find([]string{"install"}); err != nil || cmd == root {
		t.Fatalf("`install` did not resolve to a registered command (cmd=%v, err=%v)", cmd.Name(), err)
	}
	if cmd, _, _ := root.Find([]string{"init"}); cmd != root {
		t.Fatalf("`init` resolved to a command (%v) — the removed verb must fall back to the unknown-command path", cmd.Name())
	}
}

// TestClaudeFrontDoorSkipContractCheckBootstrap: the --skip-contract-check
// override launches without resolving the manifest (bootstrap case where the
// plugin is being installed for the first time).
func TestClaudeFrontDoorSkipContractCheckBootstrap(t *testing.T) {
	fake := &fakeHost{manifest: tooOldBinaryManifest(t)} // would mismatch if checked
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--skip-contract-check"}, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 with --skip-contract-check (stderr=%q)", code, stderr.String())
	}
	want := []string{"claude", "--agent", "spacedock:first-officer", wantBootstrapPrompt}
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v (skip-check must not pass the flag through)", fake.launchedArg, want)
	}
}

// TestCodexFrontDoorLaunchesOnCompatible: on a compatible contract the codex
// front door invokes the launch seam with argv beginning `codex
// --dangerously-bypass-approvals-and-sandbox` (under .safehouse) and passes
// through the operator's trailing args before the FO-skill prompt.
func TestCodexFrontDoorLaunchesOnCompatible(t *testing.T) {
	dir := safehouseFixtureDir(t)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), []string{"--", "-m", "gpt-x"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := []string{"safehouse", "--trust-workdir-config", "--",
		"codex", "--dangerously-bypass-approvals-and-sandbox", "-m", "gpt-x", wantCodexBootstrapPrompt}
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
	}
}

// TestCodexFrontDoorFailFastOnMismatch: codex fails fast on a mismatch verdict
// with the pinned remedy and does NOT launch.
func TestCodexFrontDoorInjectsLauncherBinThroughSafehouseResume(t *testing.T) {
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	t.Setenv(spacedockBinEnv, "/old/spacedock")
	dir := safehouseFixtureDir(t)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), []string{"--", "resume", "abc123"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	wantArgv := []string{"safehouse", "--trust-workdir-config", "--",
		"codex", "--dangerously-bypass-approvals-and-sandbox", "resume", "abc123"}
	if !equalArgv(fake.launchedArg, wantArgv) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, wantArgv)
	}
	got, ok := envValue(fake.launchedEnv, spacedockBinEnv)
	if !ok || got != bin {
		t.Fatalf("%s in launch env = %q, %v; want %q, true (env=%v)", spacedockBinEnv, got, ok, bin, fake.launchedEnv)
	}
}

func TestCodexFrontDoorFailFastOnMismatch(t *testing.T) {
	fake := &fakeHost{manifest: tooOldBinaryManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("exit = 0, want non-zero on mismatch")
	}
	if len(fake.installCmds) != 0 {
		t.Fatalf("install seam invoked on codex mismatch: %v", fake.installCmds)
	}
	if fake.launchedArg != nil {
		t.Fatalf("launch seam invoked on mismatch: %v", fake.launchedArg)
	}
	assertGateMismatchMessage(t, stderr.String())
}

func TestCodexFrontDoorNoPluginFailsFastWithoutInstalling(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "no-such-dir", ".codex-plugin", "plugin.json")
	cases := []struct {
		name     string
		manifest string
	}{
		{name: "empty manifest", manifest: ""},
		{name: "missing manifest", manifest: missing},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeHost{manifest: tc.manifest}
			var stdout, stderr bytes.Buffer

			code := runCodex(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

			if code == 0 {
				t.Fatalf("exit = 0, want non-zero when codex has no plugin")
			}
			if len(fake.installCmds) != 0 {
				t.Fatalf("install seam invoked on codex no-plugin path: %v", fake.installCmds)
			}
			if fake.launchedArg != nil {
				t.Fatalf("launch seam invoked on codex no-plugin path: %v", fake.launchedArg)
			}
			if !strings.Contains(stderr.String(), "codex plugin") {
				t.Fatalf("stderr missing codex no-plugin remedy: %q", stderr.String())
			}
		})
	}
}

func equalArgv(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
