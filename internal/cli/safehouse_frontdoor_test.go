// ABOUTME: AC-1..AC-5 oracles for `spacedock claude` safehouse interposition:
// ABOUTME: argv shape with/without .safehouse, --resume suppression, gate ordering.
package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// bootstrapPrompt is the fixed launch-and-go FO prompt the launcher appends as
// the last inner-argv token. Pinned here so the oracles fail loudly if the
// production constant drifts.
const wantBootstrapPrompt = "You totally got this. Take your time. I love you. And tell all subagents and team members you love them too."

// lookFound resolves any binary (safehouse Available → ok).
func lookFound(string) (string, error) { return "/usr/bin/safehouse", nil }

// lookMissing fails to resolve (safehouse Available → not ok).
func lookMissing(string) (string, error) { return "", errors.New("not found") }

// safehouseFixtureDir returns a temp dir containing a `.safehouse` profile file.
func safehouseFixtureDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".safehouse"), []byte("profile"), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

// AC-1: .safehouse present → canonical safehouse-wrapped argv with the bootstrap
// prompt appended LAST, after the operator passthrough. The wrap carries
// `--env-pass SPACEDOCK_BIN` among the safehouse flags (before `--`) so safehouse
// forwards the launcher binary through its env sanitization; the inner program is
// `claude` directly after `--` (NO /usr/bin/env, NO SPACEDOCK_BIN= argv token).
func TestClaudeSafehousePresentWrapsArgv(t *testing.T) {
	dir := safehouseFixtureDir(t)
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--", "--foo"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := []string{"safehouse", "--trust-workdir-config", "--env-pass", spacedockBinEnv, "--",
		"claude", "--dangerously-skip-permissions", "--agent", "spacedock:first-officer",
		"--foo", wantBootstrapPrompt}
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
	}
}

// AC-4: when resolvedLauncherBin() cannot resolve a binary (executablePath
// errors), the wrap carries NO `--env-pass SPACEDOCK_BIN` flag — never a stale
// pass-through — exactly as launchEnv omits the env entry. The wrap is otherwise
// the canonical shape.
func TestClaudeSafehouseWrapsWithoutEnvPassWhenBinUnresolved(t *testing.T) {
	dir := safehouseFixtureDir(t)
	withExecutablePath(t, "", errors.New("boom"))
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--", "--foo"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := []string{"safehouse", "--trust-workdir-config", "--",
		"claude", "--dangerously-skip-permissions", "--agent", "spacedock:first-officer",
		"--foo", wantBootstrapPrompt}
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
	}
	for _, tok := range fake.launchedArg {
		if tok == "--env-pass" {
			t.Fatalf("wrapped argv carried --env-pass despite unresolved bin: %v", fake.launchedArg)
		}
	}
}

// AC-6 (regression guard for the profile-detection ship-blocker): the wrapped
// inner argv's FIRST token after the safehouse `--` is the host program name
// (`claude`), NOT a wrapper like `/usr/bin/env` — so safehouse's program-keyed
// profile auto-detection sees the real program. This is the oracle that would
// have caught the argv-prefix mechanism breaking the default claude profile.
func TestClaudeSafehouseWrapInnerProgramIsHost(t *testing.T) {
	dir := safehouseFixtureDir(t)
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	assertInnerProgramIsHost(t, fake.launchedArg, "claude")
}

// assertInnerProgramIsHost asserts the token immediately after the safehouse `--`
// separator is the host program name — never a wrapper (`/usr/bin/env`, `env`)
// that would mask the program from safehouse's profile auto-detection.
func assertInnerProgramIsHost(t *testing.T, argv []string, host string) {
	t.Helper()
	for i, tok := range argv {
		if tok == "--" {
			if i+1 >= len(argv) {
				t.Fatalf("no inner program after the safehouse `--`: %v", argv)
			}
			if got := argv[i+1]; got != host {
				t.Fatalf("inner program after `--` = %q, want host %q (a wrapper here breaks safehouse profile auto-detection): %v", got, host, argv)
			}
			return
		}
	}
	t.Fatalf("no safehouse `--` separator in wrapped argv: %v", argv)
}

// AC-1 companion: a `.safehouse` DIRECTORY is detected identically to a file.
func TestClaudeSafehouseDirDetectedLikeFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, ".safehouse"), 0o755); err != nil {
		t.Fatal(err)
	}
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if len(fake.launchedArg) == 0 || fake.launchedArg[0] != "safehouse" {
		t.Fatalf("launch argv = %v, want safehouse-wrapped for a .safehouse directory", fake.launchedArg)
	}
}

// AC-2: no .safehouse → plain claude, NO --dangerously-skip-permissions, the
// token `safehouse` appears nowhere; bootstrap prompt still appended.
func TestClaudeNoSafehouseLaunchesPlain(t *testing.T) {
	dir := t.TempDir() // no .safehouse
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--", "--foo"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := []string{"claude", "--agent", "spacedock:first-officer", "--permission-mode", "auto", "--foo", wantBootstrapPrompt}
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
	}
	for _, tok := range fake.launchedArg {
		if tok == "safehouse" {
			t.Fatalf("unsandboxed launch named safehouse: %v", fake.launchedArg)
		}
		if tok == "--dangerously-skip-permissions" {
			t.Fatalf("unsandboxed launch carried --dangerously-skip-permissions: %v", fake.launchedArg)
		}
		// AC-3: `--env-pass` is wrap-path-only; an unwrapped launch carries no
		// `--env-pass` flag (SPACEDOCK_BIN still reaches the host through launchEnv,
		// asserted by TestClaudeFrontDoorInjectsResolvedLauncherBin).
		if tok == "--env-pass" {
			t.Fatalf("unsandboxed launch carried --env-pass: %v", fake.launchedArg)
		}
	}
}

// AC-3: a refused plugin gate SHORT-CIRCUITS before any safehouse logic. With
// .safehouse present AND safehouse binary absent, --no-install (which refuses the
// no-plugin case rather than auto-installing) fails at the plugin gate first:
// plugin remedy on stderr, NO safehouse hint, Launch never called. (Without
// --no-install the no-plugin case auto-installs and proceeds, so the gate is no
// longer a fail-fast there — the short-circuit-before-safehouse invariant lives
// on the refuse path.)
func TestClaudePluginGateShortCircuitsBeforeSafehouse(t *testing.T) {
	dir := safehouseFixtureDir(t)
	fake := &fakeHost{manifest: ""} // no plugin
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--no-install"}, dir, fake, lookMissing, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("exit = 0, want non-zero when plugin gate fails")
	}
	if len(fake.installCmds) != 0 {
		t.Fatalf("install invoked despite --no-install: %v", fake.installCmds)
	}
	if fake.launchedArg != nil {
		t.Fatalf("Launch invoked despite plugin-gate failure: %v", fake.launchedArg)
	}
	if !strings.Contains(stderr.String(), "spacedock install") && !strings.Contains(stderr.String(), "spacedock doctor") {
		t.Fatalf("stderr missing plugin-gate remedy: %q", stderr.String())
	}
	if strings.Contains(stderr.String(), ".safehouse profile") {
		t.Fatalf("stderr carried the safehouse hint before the plugin gate: %q", stderr.String())
	}
}

// AC-4: .safehouse present, plugin OK, but safehouse binary absent → pinned
// install hint, rc≠0, Launch never called.
func TestClaudeSafehousePresentButBinaryMissing(t *testing.T) {
	dir := safehouseFixtureDir(t)
	fake := &fakeHost{manifest: compatibleManifest(t)} // gate passes
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), nil, dir, fake, lookMissing, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("exit = 0, want non-zero when safehouse binary is missing")
	}
	if fake.launchedArg != nil {
		t.Fatalf("Launch invoked with safehouse binary absent: %v", fake.launchedArg)
	}
	if !strings.Contains(stderr.String(), ".safehouse profile") {
		t.Fatalf("stderr missing pinned safehouse install hint: %q", stderr.String())
	}
}

// codexBootstrapPrompt is the fixed codex launch-and-go prompt the launcher
// appends as the last inner-argv token. Pinned here so the codex oracles fail
// loudly if the production constant drifts. The load-bearing invariant is the
// literal `spacedock:first-officer` skill-name token (codex has no --agent).
const wantCodexBootstrapPrompt = "You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Assume $spacedock:first-officer for the entire session."

// codex AC-2: .safehouse present → canonical safehouse-wrapped codex argv with
// codex's own sandbox bypassed, forwarded post-fence tokens preserved, and the
// bootstrap prompt appended unless an exact `resume` token was forwarded.
// Mirrors runClaude's dir+lookPath threading. The wrap
// carries `--env-pass SPACEDOCK_BIN` among the safehouse flags (before `--`); the
// inner program is `codex` directly after `--` (NO /usr/bin/env, NO SPACEDOCK_BIN=
// argv token).
func TestCodexSafehousePresentWrapsArgv(t *testing.T) {
	dir := safehouseFixtureDir(t)
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), []string{"--", "--foo"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := append([]string{"safehouse", "--trust-workdir-config", "--env-pass", spacedockBinEnv, "--"},
		wantCodexArgv("--dangerously-bypass-approvals-and-sandbox", "--foo", wantCodexBootstrapPrompt)...)
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
	}
}

// codex AC-6 (regression guard): the wrapped inner argv's FIRST token after the
// safehouse `--` is `codex`, NOT a wrapper — so safehouse's program-keyed profile
// auto-detection sees the real program.
func TestCodexSafehouseWrapInnerProgramIsHost(t *testing.T) {
	dir := safehouseFixtureDir(t)
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), nil, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	assertInnerProgramIsHost(t, fake.launchedArg, "codex")
}

// codex AC-2: the emitted codex argv selects the first officer via the literal
// `spacedock:first-officer` skill token in the appended prompt (no --agent flag
// exists in codex).
func TestCodexSafehousePromptNamesFirstOfficerSkill(t *testing.T) {
	dir := safehouseFixtureDir(t)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), nil, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	last := fake.launchedArg[len(fake.launchedArg)-1]
	if !strings.Contains(last, "spacedock:first-officer") {
		t.Fatalf("emitted codex prompt does not name the first-officer skill: %q", last)
	}
	for _, tok := range fake.launchedArg {
		if tok == "--agent" {
			t.Fatalf("codex launch carried --agent (no such flag in codex): %v", fake.launchedArg)
		}
	}
}

// codex no-`.safehouse`: a post-fence argv without an exact `resume` token is a
// plain Codex launch with its default approval posture and bootstrap prompt, but
// no safehouse bypass flag.
func TestCodexNoSafehouseLaunchesPlainNoBypass(t *testing.T) {
	dir := t.TempDir() // no .safehouse
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), []string{"--", "--foo"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := wantCodexArgv("--ask-for-approval", "on-request", "--foo", wantCodexBootstrapPrompt)
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
	}
	for _, tok := range fake.launchedArg {
		if tok == "safehouse" {
			t.Fatalf("unsandboxed codex launch named safehouse: %v", fake.launchedArg)
		}
		if tok == "--dangerously-bypass-approvals-and-sandbox" {
			t.Fatalf("unsandboxed codex launch carried the bypass flag: %v", fake.launchedArg)
		}
		// AC-3: `--env-pass` is wrap-path-only; an unwrapped launch carries no
		// `--env-pass` flag (SPACEDOCK_BIN still reaches the host through launchEnv).
		if tok == "--env-pass" {
			t.Fatalf("unsandboxed codex launch carried --env-pass: %v", fake.launchedArg)
		}
	}
}

// codex AC-3 analog: a refused plugin gate SHORT-CIRCUITS before any safehouse
// logic. With .safehouse present AND the safehouse binary absent, --no-install
// (which refuses the no-plugin case rather than auto-installing) fails at the
// plugin gate first: plugin remedy on stderr, NO safehouse hint, Launch never
// called. (Without --no-install the no-plugin case auto-installs and proceeds, so
// the gate is no longer a fail-fast there — the short-circuit-before-safehouse
// invariant lives on the refuse path, mirroring runClaude.)
func TestCodexPluginGateShortCircuitsBeforeSafehouse(t *testing.T) {
	dir := safehouseFixtureDir(t)
	fake := &fakeHost{manifest: ""} // no plugin
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), []string{"--no-install"}, dir, fake, lookMissing, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("exit = 0, want non-zero when plugin gate fails")
	}
	if len(fake.installCmds) != 0 {
		t.Fatalf("install invoked despite --no-install: %v", fake.installCmds)
	}
	if fake.launchedArg != nil {
		t.Fatalf("Launch invoked despite plugin-gate failure: %v", fake.launchedArg)
	}
	if !strings.Contains(stderr.String(), "spacedock install") && !strings.Contains(stderr.String(), "spacedock doctor") {
		t.Fatalf("stderr missing plugin-gate remedy: %q", stderr.String())
	}
	if strings.Contains(stderr.String(), ".safehouse profile") {
		t.Fatalf("stderr carried the safehouse hint before the plugin gate: %q", stderr.String())
	}
}

// codex AC-4 analog: .safehouse present, plugin OK, but the safehouse binary is
// absent → pinned install hint, rc≠0, Launch never called.
func TestCodexSafehousePresentButBinaryMissing(t *testing.T) {
	dir := safehouseFixtureDir(t)
	fake := &fakeHost{manifest: compatibleManifest(t)} // gate passes
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), nil, dir, fake, lookMissing, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("exit = 0, want non-zero when safehouse binary is missing")
	}
	if fake.launchedArg != nil {
		t.Fatalf("Launch invoked with safehouse binary absent: %v", fake.launchedArg)
	}
	if !strings.Contains(stderr.String(), ".safehouse profile") {
		t.Fatalf("stderr missing pinned safehouse install hint: %q", stderr.String())
	}
}

// AC-5: any resume-family arg suppresses the bootstrap prompt; the operator arg
// forwards verbatim. The family is claude's full set of session-resume forms:
// --resume, --resume=<id>, -r, --continue, -c.
func TestClaudeResumeFamilySuppressesBootstrapPrompt(t *testing.T) {
	cases := []struct {
		name string
		arg  string
	}{
		{"resume", "--resume"},
		{"resume-equals", "--resume=abc123"},
		{"resume-short", "-r"},
		{"continue", "--continue"},
		{"continue-short", "-c"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := safehouseFixtureDir(t)
			bin := executableFixture(t)
			withExecutablePath(t, bin, nil)
			fake := &fakeHost{manifest: compatibleManifest(t)}
			var stdout, stderr bytes.Buffer

			code := runClaude(context.Background(), []string{"--", tc.arg}, dir, fake, lookFound, &stdout, &stderr)

			if code != 0 {
				t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
			}
			want := []string{"safehouse", "--trust-workdir-config", "--env-pass", spacedockBinEnv, "--",
				"claude", "--dangerously-skip-permissions", "--agent", "spacedock:first-officer", tc.arg}
			if !equalArgv(fake.launchedArg, want) {
				t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
			}
			for _, tok := range fake.launchedArg {
				if tok == wantBootstrapPrompt {
					t.Fatalf("%s launch carried the bootstrap prompt: %v", tc.arg, fake.launchedArg)
				}
			}
		})
	}
}

// AC-3: the launch flourish is gone. Now that `engage` is a real captain-invoked
// verb, no bootstrap prompt carries the throwaway "Engage." that would collide
// with it. All three (claude, codex, pi) are covered so a flourish regression on
// any is caught here, not only by git-diff emptiness (pi de-risks the 7v parity).
func TestBootstrapPromptsDropEngageFlourish(t *testing.T) {
	for _, tc := range []struct {
		name   string
		prompt string
	}{
		{"claude", bootstrapPrompt},
		{"codex", codexBootstrapPrompt},
		{"pi", piBootstrapPrompt},
	} {
		if strings.Contains(tc.prompt, "Engage.") {
			t.Errorf("%s bootstrap prompt still carries the Engage. flourish: %q", tc.name, tc.prompt)
		}
	}
}
