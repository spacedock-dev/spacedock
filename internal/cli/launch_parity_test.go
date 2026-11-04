// ABOUTME: Recorded-Launch oracles for the sandbox knobs (AC-1..AC-8) and the
// ABOUTME: launch-parity gaps (LP-AC-1..3): fence task, codex resume, plugin-dir.
package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"
)

// AC-1: --safehouse-enable=ssh,docker comma-splits into repeated --enable=KEY in
// the pre-`--` extra slot, after --trust-workdir-config.
func TestSafehouseEnableForwardsRepeatedFlags(t *testing.T) {
	dir := safehouseFixtureDir(t)
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--safehouse-enable=ssh,docker"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := []string{"safehouse", "--trust-workdir-config", "--env-pass", spacedockBinEnv, "--enable=ssh", "--enable=docker", "--",
		"claude", "--dangerously-skip-permissions", "--agent", "spacedock:first-officer", wantBootstrapPrompt}
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
	}
}

// AC-2: --safehouse-add-dirs / --safehouse-add-dirs-ro forward path grants into
// the pre-`--` extra slot, in operator order.
func TestSafehouseAddDirsForwardsPathGrants(t *testing.T) {
	dir := safehouseFixtureDir(t)
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--safehouse-add-dirs=/a", "--safehouse-add-dirs-ro=/b"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := []string{"safehouse", "--trust-workdir-config", "--env-pass", spacedockBinEnv, "--add-dirs=/a", "--add-dirs-ro=/b", "--",
		"claude", "--dangerously-skip-permissions", "--agent", "spacedock:first-officer", wantBootstrapPrompt}
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
	}
}

// AC-3: the --safehouse- prefix is stripped by the front-door dispatcher before
// reaching the translator — no --safehouse-* token survives into the safehouse
// extra slot or the inner argv.
func TestSafehousePrefixStrippedByDispatcher(t *testing.T) {
	dir := safehouseFixtureDir(t)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--safehouse-enable=docker"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	for _, tok := range fake.launchedArg {
		if tok == "--safehouse-enable=docker" || tok == "--safehouse" {
			t.Fatalf("a --safehouse* token survived into the argv: %v", fake.launchedArg)
		}
	}
}

// AC-4: explicit --safehouse forces the wrap in a no-profile dir (claude); the
// bare token is consumed and never reaches claude.
func TestClaudeForceSafehouseWrapsNoProfile(t *testing.T) {
	dir := t.TempDir() // no .safehouse
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--safehouse"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := []string{"safehouse", "--trust-workdir-config", "--env-pass", spacedockBinEnv, "--",
		"claude", "--dangerously-skip-permissions", "--agent", "spacedock:first-officer", wantBootstrapPrompt}
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
	}
}

// AC-5: explicit --safehouse forces the wrap in a no-profile dir (codex); the
// bypass flag appears only inside the forced wrap; the bare token never reaches
// codex.
func TestCodexForceSafehouseWrapsNoProfile(t *testing.T) {
	dir := t.TempDir() // no .safehouse
	bin := executableFixture(t)
	withExecutablePath(t, bin, nil)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), []string{"--safehouse"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := append([]string{"safehouse", "--trust-workdir-config", "--env-pass", spacedockBinEnv, "--"},
		wantCodexArgv("--dangerously-bypass-approvals-and-sandbox", wantCodexBootstrapPrompt)...)
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
	}
}

// AC-6: a single --safehouse-* knob implies sandbox-on in a no-profile dir for
// BOTH claude and codex (the reversal of the old fail-fast: a knob never lands
// on the plain path).
func TestKnobImpliesSandboxOnNoProfile(t *testing.T) {
	t.Run("claude", func(t *testing.T) {
		dir := t.TempDir()
		bin := executableFixture(t)
		withExecutablePath(t, bin, nil)
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer
		code := runClaude(context.Background(), []string{"--safehouse-enable=docker"}, dir, fake, lookFound, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		want := []string{"safehouse", "--trust-workdir-config", "--env-pass", spacedockBinEnv, "--enable=docker", "--",
			"claude", "--dangerously-skip-permissions", "--agent", "spacedock:first-officer", wantBootstrapPrompt}
		if !equalArgv(fake.launchedArg, want) {
			t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
		}
	})
	t.Run("codex", func(t *testing.T) {
		dir := t.TempDir()
		bin := executableFixture(t)
		withExecutablePath(t, bin, nil)
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer
		code := runCodex(context.Background(), []string{"--safehouse-enable=docker"}, dir, fake, lookFound, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		want := append([]string{"safehouse", "--trust-workdir-config", "--env-pass", spacedockBinEnv, "--enable=docker", "--"},
			wantCodexArgv("--dangerously-bypass-approvals-and-sandbox", wantCodexBootstrapPrompt)...)
		if !equalArgv(fake.launchedArg, want) {
			t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
		}
	})
}

// AC-7 codex analog: the plain (unwrapped) launch happens only when none of
// {profile, --safehouse, knob} is present — plain codex with no bypass flag.
// (The claude analog is TestClaudeNoSafehouseLaunchesPlain, unchanged.)
func TestCodexPlainWhenNoTrigger(t *testing.T) {
	dir := t.TempDir() // no .safehouse, no --safehouse, no knob
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), nil, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := wantCodexArgv("--ask-for-approval", "on-request", wantCodexBootstrapPrompt)
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
	}
}

// AC-8: an unknown --safehouse-<key> is a hard error (rc≠0, no Launch).
func TestUnknownSafehouseKeyErrors(t *testing.T) {
	dir := safehouseFixtureDir(t)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--safehouse-bogus=x"}, dir, fake, lookFound, &stdout, &stderr)

	if code == 0 {
		t.Fatalf("exit = 0, want non-zero for an unknown --safehouse-* key")
	}
	if fake.launchedArg != nil {
		t.Fatalf("Launch invoked on an unknown --safehouse-* key: %v", fake.launchedArg)
	}
}

// LP-AC-1 (Option-2 grammar): a task positional (BEFORE any `--`) becomes
// base + " " + task as the LAST inner token; bare → base exactly; a host flag
// AFTER the `--` still forwards verbatim, with the task riding before the fence.
func TestFenceTaskPromptOverride(t *testing.T) {
	t.Run("claude-task-positional", func(t *testing.T) {
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer
		code := runClaude(context.Background(), []string{"do the thing"}, t.TempDir(), fake, lookFound, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		last := fake.launchedArg[len(fake.launchedArg)-1]
		if last != wantBootstrapPrompt+" do the thing" {
			t.Fatalf("last token = %q, want base+space+task", last)
		}
	})
	t.Run("claude-bare-base-only", func(t *testing.T) {
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer
		code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		last := fake.launchedArg[len(fake.launchedArg)-1]
		if last != wantBootstrapPrompt {
			t.Fatalf("last token = %q, want bare base prompt (no trailing space)", last)
		}
	})
	t.Run("claude-task-before-fenced-host-flag", func(t *testing.T) {
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer
		code := runClaude(context.Background(), []string{"do the thing", "--", "--model", "gpt-x"}, t.TempDir(), fake, lookFound, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		want := []string{"claude", "--agent", "spacedock:first-officer", "--permission-mode", "auto", "--model", "gpt-x", wantBootstrapPrompt + " do the thing"}
		if !equalArgv(fake.launchedArg, want) {
			t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
		}
	})
	t.Run("codex-task-positional", func(t *testing.T) {
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer
		code := runCodex(context.Background(), []string{"do the thing"}, t.TempDir(), fake, lookFound, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		last := fake.launchedArg[len(fake.launchedArg)-1]
		if last != wantCodexBootstrapPrompt+" do the thing" {
			t.Fatalf("last token = %q, want codexBase+space+task", last)
		}
	})
	t.Run("codex-task-before-fenced-model", func(t *testing.T) {
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer
		code := runCodex(context.Background(), []string{"@/tmp/handoff-file.md", "--", "--model", "gpt-5.6-sol"}, t.TempDir(), fake, lookFound, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		want := wantCodexArgv("--ask-for-approval", "on-request", "--model", "gpt-5.6-sol", wantCodexBootstrapPrompt+" @/tmp/handoff-file.md")
		if !equalArgv(fake.launchedArg, want) {
			t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
		}
	})
}

// LP-AC-2: an exact post-fence Codex `resume` token suppresses the prompt; bare
// Codex gets the bootstrap prompt.
func TestCodexResumeSuppressesPrompt(t *testing.T) {
	t.Run("exact-resume-token-no-prompt", func(t *testing.T) {
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer
		code := runCodex(context.Background(), []string{"--", "resume", "abc-123"}, t.TempDir(), fake, lookFound, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		want := wantCodexArgv("resume", "abc-123")
		if !equalArgv(fake.launchedArg, want) {
			t.Fatalf("launch argv = %v, want %v (exact resume forwards without a prompt)", fake.launchedArg, want)
		}
		for _, tok := range fake.launchedArg {
			if tok == wantCodexBootstrapPrompt {
				t.Fatalf("codex post-fence argv carried the bootstrap prompt: %v", fake.launchedArg)
			}
		}
	})
	t.Run("bare-codex-base-prompt", func(t *testing.T) {
		fake := &fakeHost{manifest: compatibleManifest(t)}
		var stdout, stderr bytes.Buffer
		code := runCodex(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
		}
		last := fake.launchedArg[len(fake.launchedArg)-1]
		if last != wantCodexBootstrapPrompt {
			t.Fatalf("bare codex last token = %q, want codex base prompt", last)
		}
	})
}

// Codex treats only an exact `resume` token in the forwarded argv as a resume.
// Other post-fence argv stays byte-for-byte intact but retains the normal
// first-officer launch posture.
func TestCodexPostFenceUsesExactResumeToken(t *testing.T) {
	tests := []struct {
		name        string
		passthrough []string
		want        []string
		wantBanner  bool
	}{
		{
			name:        "model only retains bootstrap posture",
			passthrough: []string{"--model", "gpt-5.6-sol"},
			want:        wantCodexArgv("--ask-for-approval", "on-request", "--model", "gpt-5.6-sol", wantCodexBootstrapPrompt),
			wantBanner:  true,
		},
		{
			name:        "model plus resume stays prompt free",
			passthrough: []string{"--model", "gpt-5.6-sol", "resume", "abc-123"},
			want:        wantCodexArgv("--model", "gpt-5.6-sol", "resume", "abc-123"),
		},
		{
			name:        "resume-like option stays a fresh launch",
			passthrough: []string{"--model", "gpt-5.6-sol", "--resume=abc-123"},
			want:        wantCodexArgv("--ask-for-approval", "on-request", "--model", "gpt-5.6-sol", "--resume=abc-123", wantCodexBootstrapPrompt),
			wantBanner:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &fakeHost{manifest: compatibleManifest(t)}
			var stdout, stderr bytes.Buffer
			args := append([]string{"--"}, tt.passthrough...)
			code := runCodex(context.Background(), args, t.TempDir(), fake, lookFound, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
			}
			if !equalArgv(fake.launchedArg, tt.want) {
				t.Fatalf("launch argv = %v, want %v", fake.launchedArg, tt.want)
			}
			if tt.wantBanner && !strings.Contains(stderr.String(), "\u00b7 launching codex") {
				t.Fatalf("fresh post-fence launch omitted the Codex banner: %q", stderr.String())
			}
			if !tt.wantBanner && stderr.Len() != 0 {
				t.Fatalf("resume launch produced Spacedock output: %q", stderr.String())
			}
		})
	}
}

// A post-fence Codex argv keeps its original token order. The launcher does not
// parse its grammar; absent an exact `resume` token, it adds normal fresh-launch
// posture around the forwarded tokens.
func TestCodexPostFenceWithoutResumeBootstrapsFirstOfficer(t *testing.T) {
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer
	args := []string{"--", "--future-codex-flag=handoff", "opaque-argument"}
	code := runCodex(context.Background(), args, t.TempDir(), fake, lookFound, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	want := wantCodexArgv("--ask-for-approval", "on-request", "--future-codex-flag=handoff", "opaque-argument", wantCodexBootstrapPrompt)
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
	}
	if !strings.Contains(stderr.String(), "\u00b7 launching codex") {
		t.Fatalf("fresh post-fence launch omitted the Codex banner: %q", stderr.String())
	}
}

// LP-AC-3: for claude, only a declared pre-fence --plugin-dir passes through AND
// relaxes the gate; a post-fence directory is a native addition and the gate
// still runs. Codex is the exception: its CLI has no --plugin-dir flag, so the
// flag is consumed into a real local-marketplace install and the gate then runs
// against that install rather than being relaxed (the codex sub-case below).
func TestPluginDirRelaxesGate(t *testing.T) {
	t.Run("claude-post-fence-does-not-relax-failing-manifest", func(t *testing.T) {
		fake := &fakeHost{manifest: tooOldBinaryManifest(t)} // gate would FAIL
		var stdout, stderr bytes.Buffer
		code := runClaude(context.Background(), []string{"--", "--plugin-dir", "/a", "--plugin-dir", "/b"}, t.TempDir(), fake, lookFound, &stdout, &stderr)
		if code == 0 {
			t.Fatalf("exit = 0, want installed-gate failure for post-fence additions; stderr=%q", stderr.String())
		}
		if fake.launchedArg != nil {
			t.Fatalf("post-fence additions bypassed gate and launched: %v", fake.launchedArg)
		}
	})
	t.Run("codex-plugin-dir-installs-then-gates-not-relaxed", func(t *testing.T) {
		// Unlike claude's ephemeral --plugin-dir bypass, codex's --plugin-dir is a real
		// install: the flag is consumed (never forwarded to codex) and the gate then
		// runs against the installed plugin. With a fake still reporting a mismatched
		// manifest, the gate FAILS — proving codex --plugin-dir installs then gates
		// rather than relaxing.
		t.Setenv("CODEX_HOME", t.TempDir()) // isolate the persistent local marketplace
		checkout, _ := localPluginCheckout(t, "codex")
		fake := &fakeHost{manifest: tooOldBinaryManifest(t)}
		var stdout, stderr bytes.Buffer
		code := runCodex(context.Background(), []string{"--plugin-dir", checkout}, t.TempDir(), fake, lookFound, &stdout, &stderr)
		if code == 0 {
			t.Fatalf("exit = 0, want non-zero (codex --plugin-dir gate-checks the install, not relaxes); stderr=%q", stderr.String())
		}
		if len(fake.installCmds) == 0 {
			t.Fatalf("codex --plugin-dir did not install before gating: installCmds empty")
		}
		if fake.launchedArg != nil {
			t.Fatalf("launch reached despite the post-install gate mismatch: %v", fake.launchedArg)
		}
	})
	t.Run("claude-before-dash-forwards-and-relaxes", func(t *testing.T) {
		fake := &fakeHost{manifest: tooOldBinaryManifest(t)} // gate would FAIL
		var stdout, stderr bytes.Buffer
		code := runClaude(context.Background(), []string{"--plugin-dir", "/a", "--plugin-dir=/b"}, t.TempDir(), fake, lookFound, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit = %d, want 0 (before-`--` --plugin-dir relaxes the gate); stderr=%q", code, stderr.String())
		}
		want := []string{"claude", "--agent", "spacedock:first-officer", "--permission-mode", "auto", "--plugin-dir", "/a", "--plugin-dir", "/b", wantBootstrapPrompt}
		if !equalArgv(fake.launchedArg, want) {
			t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
		}
	})
	t.Run("claude-before-dash-plugin-dir-with-task", func(t *testing.T) {
		fake := &fakeHost{manifest: tooOldBinaryManifest(t)} // gate would FAIL
		var stdout, stderr bytes.Buffer
		code := runClaude(context.Background(), []string{"--plugin-dir", "/a", "review the PRs"}, t.TempDir(), fake, lookFound, &stdout, &stderr)
		if code != 0 {
			t.Fatalf("exit = %d, want 0 (captain no-`--` form relaxes the gate); stderr=%q", code, stderr.String())
		}
		want := []string{"claude", "--agent", "spacedock:first-officer", "--permission-mode", "auto", "--plugin-dir", "/a", wantBootstrapPrompt + " review the PRs"}
		if !equalArgv(fake.launchedArg, want) {
			t.Fatalf("launch argv = %v, want %v", fake.launchedArg, want)
		}
	})
	t.Run("no-plugin-dir-still-fails-fast", func(t *testing.T) {
		fake := &fakeHost{manifest: tooOldBinaryManifest(t)}
		var stdout, stderr bytes.Buffer
		code := runClaude(context.Background(), nil, t.TempDir(), fake, lookFound, &stdout, &stderr)
		if code == 0 {
			t.Fatalf("exit = 0, want non-zero with a failing manifest and no --plugin-dir")
		}
		if fake.launchedArg != nil {
			t.Fatalf("Launch invoked despite failing gate and no --plugin-dir: %v", fake.launchedArg)
		}
	})
}

func recordSandboxLaunch(t *testing.T, host string, args []string, dir string, available bool) (argv, env []string, code int, diagnostic string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	look := lookFound
	if !available {
		look = lookMissing
	}
	if host == "pi" {
		repo, pkg := t.TempDir(), t.TempDir()
		writePiSkillFixtures(t, repo)
		writePiSubagentsFixtures(t, pkg)
		ops := piSafehouseReadyOps(repo, pkg)
		if !available {
			delete(ops.lookPath, "safehouse")
		}
		code = runPi(context.Background(), args, dir, append(os.Environ(), piTestEnv(pkg, t.TempDir())...), ops, &stdout, &stderr)
		argv, env = ops.launched, ops.launchedEnv
	} else {
		ops := &fakeHost{manifest: compatibleManifest(t)}
		if host == "claude" {
			code = runClaude(context.Background(), args, dir, ops, look, &stdout, &stderr)
		} else {
			code = runCodex(context.Background(), args, dir, ops, look, &stdout, &stderr)
		}
		argv, env = ops.launchedArg, ops.launchedEnv
	}
	return argv, env, code, stderr.String()
}

func TestSandboxEnvironmentPrecedence(t *testing.T) {
	withExecutablePath(t, executableFixture(t), nil)
	for _, host := range []string{"claude", "codex", "pi"} {
		for _, provider := range []bool{false, true} {
			for _, tc := range []struct {
				name, force, enable          string
				flags, extra                 []string
				profile, wrap, missing, fail bool
			}{
				{name: "unset"}, {name: "empty"}, {name: "force", force: "1", wrap: true},
				{name: "false", force: "0"}, {name: "parse bool", force: "TRUE", wrap: true},
				{name: "invalid", force: "bad", fail: true}, {name: "flag false", force: "bad", flags: []string{"--safehouse=false"}},
				{name: "flag true", force: "bad", flags: []string{"--safehouse=true"}, wrap: true},
				{name: "capabilities", enable: "ssh,docker", wrap: true, extra: []string{"--enable=ssh", "--enable=docker"}},
				{name: "replace", enable: "network", flags: []string{"--safehouse-enable=docker", "--safehouse-enable=ssh,docker"}, wrap: true, extra: []string{"--enable=docker", "--enable=ssh", "--enable=docker"}},
				{name: "explicit empty", enable: "ssh", flags: []string{"--safehouse-enable="}, wrap: true, extra: []string{"--enable="}},
				{name: "profile", force: "0", profile: true, wrap: true},
				{name: "independent grant", force: "1", flags: []string{"--safehouse=false", "--safehouse-add-dirs=/chosen"}, wrap: true, extra: []string{"--add-dirs=/chosen"}},
				{name: "missing sandbox", force: "1", missing: true, fail: true},
			} {
				t.Run(fmt.Sprintf("%s/provider=%t/%s", host, provider, tc.name), func(t *testing.T) {
					t.Setenv("SPACEDOCK_SAFEHOUSE", tc.force)
					t.Setenv("SPACEDOCK_SAFEHOUSE_ENABLE", tc.enable)
					if tc.name == "unset" {
						os.Unsetenv("SPACEDOCK_SAFEHOUSE")
						os.Unsetenv("SPACEDOCK_SAFEHOUSE_ENABLE")
					}
					dir, args := t.TempDir(), slices.Clone(tc.flags)
					if tc.profile {
						dir = safehouseFixtureDir(t)
					}
					if provider {
						args = append(args, "--passthrough")
					}
					argv, _, code, diagnostic := recordSandboxLaunch(t, host, args, dir, !tc.missing)
					if tc.fail {
						if code == 0 || argv != nil || diagnostic == "" {
							t.Fatalf("exit=%d argv=%q stderr=%q", code, argv, diagnostic)
						}
						return
					}
					if code != 0 || len(argv) == 0 || (argv[0] == "safehouse") != tc.wrap {
						t.Fatalf("exit=%d argv=%q stderr=%q", code, argv, diagnostic)
					}
					var actual []string
					for _, arg := range piSafehouseExtra(argv) {
						if strings.HasPrefix(arg, "--enable=") || strings.HasPrefix(arg, "--add-dirs=") {
							actual = append(actual, arg)
						}
					}
					if !slices.Equal(actual, tc.extra) {
						t.Fatalf("wrapper flags=%q, want %q", actual, tc.extra)
					}
				})
			}
		}
	}
}

func TestProviderCanonicalSettings(t *testing.T) {
	withExecutablePath(t, executableFixture(t), nil)
	names := []string{"CLAUDE_CONFIG_DIR", "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS", "CODEX_HOME", "PI_CODING_AGENT_DIR"}
	for _, host := range []string{"claude", "codex", "pi"} {
		for _, mode := range []string{"wrapped provider", "plain provider", "ordinary"} {
			for _, value := range []string{"absent", "", "0", "settings with space=value"} {
				t.Run(host+"/"+mode+"/"+value, func(t *testing.T) {
					for _, name := range names {
						t.Setenv(name, value)
						if value == "absent" {
							os.Unsetenv(name)
						}
					}
					args := []string{"--safehouse"}
					if mode != "ordinary" {
						args = append(args, "--passthrough")
					}
					if mode == "plain provider" {
						args = []string{"--passthrough"}
					}
					argv, env, code, diagnostic := recordSandboxLaunch(t, host, args, t.TempDir(), true)
					if code != 0 {
						t.Fatalf("exit=%d stderr=%q", code, diagnostic)
					}
					for _, name := range names {
						selected := host == "claude" && (name == "CLAUDE_CONFIG_DIR" || name == "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS") || host == "codex" && name == "CODEX_HOME" || host == "pi" && name == "PI_CODING_AGENT_DIR"
						wantAllowance := mode == "wrapped provider" && selected && value != "absent"
						if slices.Contains(piSafehouseExtra(argv), name) != wantAllowance {
							t.Fatalf("allowance %s: argv=%q want=%t", name, argv, wantAllowance)
						}
						if mode != "ordinary" && (hasEnv(env, name) != (value != "absent") || value != "absent" && envGetenv(env)(name) != value) {
							t.Fatalf("value %s=%q present=%t, want %q", name, envGetenv(env)(name), hasEnv(env, name), value)
						}
					}
				})
			}
		}
	}
}
