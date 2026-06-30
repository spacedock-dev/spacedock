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
	launchCode  int      // host exit code Launch returns (default 0)
	launchErr   error
	installCmds []string // host commands captured by Install
	installOut  string
}

func (f *fakeHost) ResolveManifest(host string) (string, error) {
	return f.manifest, f.resolveErr
}

func (f *fakeHost) Launch(argv []string, env []string) (int, error) {
	f.launchedArg = argv
	f.launchedEnv = env
	return f.launchCode, f.launchErr
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

// withVersion stamps the package Version (the binary display semver the gate
// feeds to the upgrade-hint compare), restoring it after the test. The package
// default is the `dev` sentinel, which suppresses the hint, so a test that
// exercises the behind-plugin hint must stamp a real semver.
func withVersion(t *testing.T, v string) {
	t.Helper()
	orig := Version
	Version = v
	t.Cleanup(func() { Version = orig })
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
	want := []string{"claude", "--agent", "spacedock:first-officer", "--permission-mode", "auto", "-p", "do the thing", wantBootstrapPrompt}
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
	}
}

// TestFrontDoorUpgradeHintOnBehindPlugin is AC-4: the front-door gate prints the
// opt-in upgrade hint to stderr when the resolved plugin is contract-compatible
// but behind the binary, then proceeds to launch (the hint never blocks). The
// compatible.json fixture is plugin 0.12.1; stamping the binary Version to a
// strictly-newer semver makes it behind-but-compatible. A paired equal-version
// case (binary == plugin) asserts the gate stays silent. Both arms observe the
// recorded launch + stderr, never a source grep.
func TestFrontDoorUpgradeHintOnBehindPlugin(t *testing.T) {
	cases := []struct {
		name string
		host string
		run  func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int
	}{
		{"claude", "claude", func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runClaude(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
		}},
		{"codex", "codex", func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runCodex(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name+" behind-plugin hints + launches", func(t *testing.T) {
			withVersion(t, "0.20.0") // strictly newer than the fixture plugin 0.12.1
			fake := &fakeHost{manifest: compatibleManifest(t)}
			var stderr bytes.Buffer

			code := tc.run(nil, t.TempDir(), fake, &stderr)

			if code != 0 {
				t.Fatalf("exit = %d, want 0 (the hint must not block launch) (stderr=%q)", code, stderr.String())
			}
			if fake.launchedArg == nil {
				t.Fatalf("launch seam not invoked — the hint must not change the launch path")
			}
			out := stderr.String()
			if !strings.Contains(out, "newer plugin") {
				t.Fatalf("stderr missing the behind-plugin upgrade hint: %q", out)
			}
			if !strings.Contains(out, "spacedock install --host "+tc.host) {
				t.Fatalf("stderr hint missing host install command for %s: %q", tc.host, out)
			}
		})

		t.Run(tc.name+" equal-version stays silent", func(t *testing.T) {
			withVersion(t, "0.12.1") // exactly the fixture plugin version
			fake := &fakeHost{manifest: compatibleManifest(t)}
			var stderr bytes.Buffer

			code := tc.run(nil, t.TempDir(), fake, &stderr)

			if code != 0 {
				t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
			}
			if fake.launchedArg == nil {
				t.Fatalf("launch seam not invoked on an equal-version compatible plugin")
			}
			if strings.Contains(stderr.String(), "newer plugin") || strings.Contains(stderr.String(), "spacedock install") {
				t.Fatalf("equal-version gate must stay silent (no upgrade hint): %q", stderr.String())
			}
		})
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

// TestClaudeFrontDoorEnablesAgentTeamsWhenParentUnset (AC-1): with no parent
// value for CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS, the launched claude child env
// carries it set to 1 — a best-effort export, not the authoritative enabler (the
// authoritative enabler is ~/.claude/settings.json; see agentTeamsEnv).
func TestClaudeFrontDoorEnablesAgentTeamsWhenParentUnset(t *testing.T) {
	t.Setenv(agentTeamsEnv, "") // register restore-on-cleanup, then unset for real
	os.Unsetenv(agentTeamsEnv)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	got, ok := envValue(fake.launchedEnv, agentTeamsEnv)
	if !ok || got != "1" {
		t.Fatalf("%s in launch env = %q, %v; want %q, true (env=%v)", agentTeamsEnv, got, ok, "1", fake.launchedEnv)
	}
}

// TestClaudeFrontDoorPreservesExplicitAgentTeams (AC-2): an explicit parent value
// for CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS is preserved, not overridden — an
// operator who set =0 keeps =0. The launcher respects the override either way.
func TestClaudeFrontDoorPreservesExplicitAgentTeams(t *testing.T) {
	t.Setenv(agentTeamsEnv, "0")
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	got, ok := envValue(fake.launchedEnv, agentTeamsEnv)
	if !ok || got != "0" {
		t.Fatalf("%s in launch env = %q, %v; want %q, true (env=%v)", agentTeamsEnv, got, ok, "0", fake.launchedEnv)
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

// TestGateRemedyNamesLiveInstallCommand: every remedy must point at a command
// the binary actually recognizes. After the init->install rename a user who hits
// the gate and is told to "run spacedock init --host …" runs a command that now
// exits 2 (unknown command). gateHost owns only the always-fail-fast remedies
// (resolve error → MalformedRange); the NoPluginFound message is the caller's
// (it auto-installs by default, refuses under --no-install), so the no-plugin /
// missing-manifest remedies are asserted on the launcher's --no-install output.
// Each remedy must name `spacedock install` and never `spacedock init`.
func TestGateRemedyNamesLiveInstallCommand(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "no-such-dir", ".claude-plugin", "plugin.json")

	// gateHost owns the resolve-error remedy (a host-CLI failure is a hard fail,
	// not a missing plugin). It MUST print here.
	t.Run("resolve error (gateHost-owned)", func(t *testing.T) {
		var stderr bytes.Buffer
		if v := gateHost(&fakeHost{resolveErr: errors.New("host CLI failed")}, "claude", &stderr); v == contract.Compatible {
			t.Fatalf("gateHost = Compatible, want denied")
		}
		assertRemedyNamesInstall(t, stderr.String())
	})

	// The NoPluginFound remedies (no plugin, phantom manifest) are caller-owned:
	// gateHost no longer prints for them, so assert the remedy on the launcher's
	// --no-install refuse output.
	noPluginCases := []struct {
		name     string
		manifest string
	}{
		{"no plugin", ""},
		{"missing manifest", missing},
	}
	for _, tc := range noPluginCases {
		t.Run(tc.name+" (caller-owned via --no-install)", func(t *testing.T) {
			fake := &fakeHost{manifest: tc.manifest}
			var stdout, stderr bytes.Buffer
			if code := runClaude(context.Background(), []string{"--no-install"}, t.TempDir(), fake, lookFound, &stdout, &stderr); code == 0 {
				t.Fatalf("exit = 0, want non-zero with --no-install and no plugin")
			}
			assertRemedyNamesInstall(t, stderr.String())
		})
		t.Run(tc.name+" (gateHost stays silent)", func(t *testing.T) {
			var stderr bytes.Buffer
			if v := gateHost(&fakeHost{manifest: tc.manifest}, "claude", &stderr); v != contract.NoPluginFound {
				t.Fatalf("gateHost verdict = %v, want NoPluginFound", v)
			}
			if stderr.Len() != 0 {
				t.Fatalf("gateHost printed for NoPluginFound (caller owns the message): %q", stderr.String())
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

// TestNoPluginAutoInstallAnnouncesHostCorrectly (AC-A, auto-install arm): with
// no installed plugin the launcher announces `Installing the {host} plugin…` on
// stderr before installing, and a codex run never names `claude` (the old
// gateHost remedy hardcoded a `spacedock claude` hint that was wrong in a codex
// run). Install + launch are observed via the recorded seams, not a string
// match.
func TestNoPluginAutoInstallAnnouncesHostCorrectly(t *testing.T) {
	cases := []struct {
		name string
		host string
		run  func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int
	}{
		{"claude", "claude", func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runClaude(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
		}},
		{"codex", "codex", func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runCodex(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeHost{manifest: ""} // no plugin found
			var stderr bytes.Buffer

			code := tc.run(nil, t.TempDir(), fake, &stderr)

			if code != 0 {
				t.Fatalf("exit = %d, want 0 (auto-install + launch) (stderr=%q)", code, stderr.String())
			}
			want := "Installing the " + tc.host + " plugin"
			if !strings.Contains(stderr.String(), want) {
				t.Fatalf("stderr missing %q announcement: %q", want, stderr.String())
			}
			if tc.host == "codex" && strings.Contains(stderr.String(), "spacedock claude") {
				t.Fatalf("codex auto-install stderr names claude: %q", stderr.String())
			}
			if len(fake.installCmds) == 0 {
				t.Fatalf("install seam not invoked: auto-install did not run")
			}
			if fake.launchedArg == nil {
				t.Fatalf("launch seam not invoked after auto-install")
			}
		})
	}
}

// TestNoPluginNoInstallRemedyIsHostCorrect (AC-A, refuse arm): with --no-install
// and no plugin the launcher prints the manual remedy naming the host-correct
// install/bootstrap commands (`spacedock install --host {host}`, `spacedock
// {host} --skip-contract-check`) and NEVER a `claude` hint in a codex run, then
// exits non-zero without installing or launching. This is the message the caller
// now owns (gateHost stopped printing for NoPluginFound).
func TestNoPluginNoInstallRemedyIsHostCorrect(t *testing.T) {
	cases := []struct {
		name string
		host string
		run  func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int
	}{
		{"claude", "claude", func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runClaude(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
		}},
		{"codex", "codex", func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
			var stdout bytes.Buffer
			return runCodex(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeHost{manifest: ""} // no plugin found
			var stderr bytes.Buffer

			code := tc.run([]string{"--no-install"}, t.TempDir(), fake, &stderr)

			if code == 0 {
				t.Fatalf("exit = 0, want non-zero with --no-install and no plugin")
			}
			out := stderr.String()
			if !strings.Contains(out, "spacedock install --host "+tc.host) {
				t.Fatalf("remedy missing host-correct install command for %s: %q", tc.host, out)
			}
			if !strings.Contains(out, "spacedock "+tc.host+" --skip-contract-check") {
				t.Fatalf("remedy missing host-correct bootstrap command for %s: %q", tc.host, out)
			}
			if tc.host == "codex" && strings.Contains(out, "claude") {
				t.Fatalf("codex --no-install remedy names claude: %q", out)
			}
			if len(fake.installCmds) != 0 {
				t.Fatalf("install seam invoked despite --no-install: %v", fake.installCmds)
			}
			if fake.launchedArg != nil {
				t.Fatalf("launch seam invoked despite --no-install: %v", fake.launchedArg)
			}
		})
	}
}

// assertRemedyNamesInstall: every gate / refuse remedy must name the live
// `spacedock install` command and never the removed `spacedock init` (which now
// exits 2).
func assertRemedyNamesInstall(t *testing.T, remedy string) {
	t.Helper()
	if !strings.Contains(remedy, "spacedock install") {
		t.Fatalf("remedy does not name the live install command: %q", remedy)
	}
	if strings.Contains(remedy, "spacedock init") {
		t.Fatalf("remedy names the removed init command (exits 2): %q", remedy)
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
	want := []string{"claude", "--agent", "spacedock:first-officer", "--permission-mode", "auto", wantBootstrapPrompt}
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
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), []string{"--", "-m", "gpt-x"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := []string{"safehouse", "--trust-workdir-config", "--env-pass", spacedockBinEnv, "--",
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
	wantArgv := []string{"safehouse", "--trust-workdir-config", "--env-pass", spacedockBinEnv, "--",
		"codex", "--dangerously-bypass-approvals-and-sandbox", "resume", "abc123"}
	if !equalArgv(fake.launchedArg, wantArgv) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, wantArgv)
	}
	got, ok := envValue(fake.launchedEnv, spacedockBinEnv)
	if !ok || got != bin {
		t.Fatalf("%s in launch env = %q, %v; want %q, true (env=%v)", spacedockBinEnv, got, ok, bin, fake.launchedEnv)
	}
}

// TestCodexFrontDoorDoesNotEnableAgentTeams: the agent-teams flag is a
// claude-only concern; the codex launch path must NOT inject it (scope guard).
func TestCodexFrontDoorDoesNotEnableAgentTeams(t *testing.T) {
	t.Setenv(agentTeamsEnv, "") // register restore-on-cleanup, then unset for real
	os.Unsetenv(agentTeamsEnv)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if got, ok := envValue(fake.launchedEnv, agentTeamsEnv); ok {
		t.Fatalf("%s in codex launch env = %q, want omitted", agentTeamsEnv, got)
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

// TestCodexFrontDoorNoPluginAutoInstalls (AC-2): with no installed codex plugin
// and no flags the front door auto-installs the plugin then launches a working
// session — the single command the user typed yields a working FO session rather
// than refusing. Both no-plugin verdicts (an empty manifest AND a phantom
// installPath to an absent manifest) auto-install, matching the claude branch.
// The recorded install seam ({host, source, branch}) and the captured launch argv
// are the independent sources of truth.
func TestCodexFrontDoorNoPluginAutoInstalls(t *testing.T) {
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

			if code != 0 {
				t.Fatalf("exit = %d, want 0 when codex has no plugin → auto-install + launch (stderr=%q)", code, stderr.String())
			}
			wantInstall := []string{"codex", channelMarketplaceSource(devBranch), devBranch}
			if !equalArgv(fake.installCmds, wantInstall) {
				t.Fatalf("install seam = %v, want %v (the {host, source, branch} seam)", fake.installCmds, wantInstall)
			}
			if fake.launchedArg == nil {
				t.Fatalf("launch seam not invoked after codex auto-install")
			}
		})
	}
}

// TestCodexFrontDoorNoPluginNoInstallRefuses (AC-3): --no-install preserves the
// refuse-and-instruct behavior — no install, no launch, non-zero exit, with the
// codex no-plugin remedy on stderr. This is the migrated assertion from the old
// fail-fast test, now gated behind --no-install. Both no-plugin verdicts refuse.
func TestCodexFrontDoorNoPluginNoInstallRefuses(t *testing.T) {
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

			code := runCodex(context.Background(), []string{"--no-install"}, t.TempDir(), fake, lookFound, &stdout, &stderr)

			if code == 0 {
				t.Fatalf("exit = 0, want non-zero with --no-install and no codex plugin")
			}
			if len(fake.installCmds) != 0 {
				t.Fatalf("install seam invoked despite --no-install: %v", fake.installCmds)
			}
			if fake.launchedArg != nil {
				t.Fatalf("launch seam invoked despite --no-install: %v", fake.launchedArg)
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
