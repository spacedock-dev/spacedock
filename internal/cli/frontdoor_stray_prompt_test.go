// ABOUTME: Stray-positional-after-`--` advisory guard tests — the warn-not-change
// ABOUTME: oracle (AC-1 positive, AC-2 negatives) and the classifier unit table (AC-3).
package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// strayWarnMarker is the stable substring both the positive and negative oracles
// key on: it survives benign rewording of the rest of the warning while still
// failing if the advisory is dropped. The production warning text must contain it.
const strayWarnMarker = "after `--`"

// TestClaudeStrayPromptAfterDashWarns (AC-1 positive): a bare positional after
// `--` with no task before `--` is named in a stderr warning AND the assembled
// inner argv is byte-identical to the pre-guard argv — proving warn-not-change.
func TestClaudeStrayPromptAfterDashWarns(t *testing.T) {
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--", "--model", "gpt-x", "@/tmp/handoff.md"}, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	warn := stderr.String()
	if !strings.Contains(warn, strayWarnMarker) {
		t.Fatalf("stderr missing the stray-positional warning marker %q: %q", strayWarnMarker, warn)
	}
	if !strings.Contains(warn, "@/tmp/handoff.md") {
		t.Fatalf("warning does not name the stray positional: %q", warn)
	}
	if !strings.Contains(warn, "BEFORE") {
		t.Fatalf("warning does not name the corrected form (prompt BEFORE `--`): %q", warn)
	}
	want := []string{"claude", "--agent", "spacedock:first-officer", "--model", "gpt-x", "@/tmp/handoff.md", wantBootstrapPrompt}
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v (warn must not change the argv)", fake.launchedArg, want)
	}
}

// TestClaudeStrayPromptSession12Shape (AC-1, the captain's actual session-12
// case): `--plugin-dir "$(pwd)" -- --model gpt-x '@/tmp/handoff.md'`. The
// spacedock-injected `--plugin-dir <dir>` prefix (re-prepended to fd.passthrough
// by parseFrontDoorArgs) must NOT shadow the operator's real stray prompt — the
// warning names `@/tmp/handoff.md`, never the injected dir, and the inner argv is
// unchanged.
func TestClaudeStrayPromptSession12Shape(t *testing.T) {
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runClaude(context.Background(), []string{"--plugin-dir", "/co", "--", "--model", "gpt-x", "@/tmp/handoff.md"}, t.TempDir(), fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	warn := stderr.String()
	if !strings.Contains(warn, strayWarnMarker) {
		t.Fatalf("stderr missing the stray-positional warning marker %q: %q", strayWarnMarker, warn)
	}
	if !strings.Contains(warn, "@/tmp/handoff.md") {
		t.Fatalf("warning does not name the operator's stray positional: %q", warn)
	}
	if strings.Contains(warn, "/co") {
		t.Fatalf("warning names the spacedock-injected --plugin-dir value (shadows the real prompt): %q", warn)
	}
	want := []string{"claude", "--agent", "spacedock:first-officer", "--plugin-dir", "/co", "--model", "gpt-x", "@/tmp/handoff.md", wantBootstrapPrompt}
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v (warn must not change the argv)", fake.launchedArg, want)
	}
}

// TestCodexStrayPromptAfterDashWarns (AC-1 positive, codex): the subcommand-less
// passthrough leading with a value-taking flag still surfaces the stray prompt,
// and the codex inner argv is unchanged.
func TestCodexStrayPromptAfterDashWarns(t *testing.T) {
	dir := safehouseFixtureDir(t)
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runCodex(context.Background(), []string{"--", "-m", "gpt-x", "@/tmp/handoff.md"}, dir, fake, lookFound, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	warn := stderr.String()
	if !strings.Contains(warn, strayWarnMarker) {
		t.Fatalf("stderr missing the stray-positional warning marker %q: %q", strayWarnMarker, warn)
	}
	if !strings.Contains(warn, "@/tmp/handoff.md") {
		t.Fatalf("warning does not name the stray positional: %q", warn)
	}
	want := []string{"safehouse", "--trust-workdir-config", "--",
		"codex", "--dangerously-bypass-approvals-and-sandbox", "-m", "gpt-x", "@/tmp/handoff.md", wantCodexBootstrapPrompt}
	if !equalArgv(fake.launchedArg, want) {
		t.Fatalf("launch argv = %v, want %v (warn must not change the argv)", fake.launchedArg, want)
	}
}

// TestStrayPromptGuardNegatives (AC-2): a legitimate host positional after `--`
// does NOT trip the guard, and the argv is unchanged in each case. The three
// negatives are the value of a value-taking flag (`-p <prompt>`), the argument of
// a known leading subcommand (`exec <prompt>`), and the hasTask short-circuit.
func TestStrayPromptGuardNegatives(t *testing.T) {
	cases := []struct {
		name string
		run  func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int
		args []string
		want []string
	}{
		{
			name: "claude -p value-of-value-taking-flag",
			run: func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
				var stdout bytes.Buffer
				return runClaude(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
			},
			args: []string{"--", "-p", "do the thing"},
			want: []string{"claude", "--agent", "spacedock:first-officer", "-p", "do the thing", wantBootstrapPrompt},
		},
		{
			name: "codex exec subcommand-arg",
			run: func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
				var stdout bytes.Buffer
				return runCodex(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
			},
			args: []string{"--", "exec", "do the thing"},
			want: []string{"codex", "exec", "do the thing", wantCodexBootstrapPrompt},
		},
		{
			name: "claude hasTask short-circuit",
			run: func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
				var stdout bytes.Buffer
				return runClaude(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
			},
			args: []string{"task before", "--", "@/tmp/handoff.md"},
			want: []string{"claude", "--agent", "spacedock:first-officer", "@/tmp/handoff.md", wantBootstrapPrompt + " task before"},
		},
		{
			// An unrecognized `-`-prefixed flag's value must NOT get the prescriptive
			// "put X before --" advice — it is likely that flag's value, so the guard
			// stays silent rather than mis-advise. The argv is still unchanged.
			name: "claude value of unrecognized flag — no wrong advice",
			run: func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
				var stdout bytes.Buffer
				return runClaude(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
			},
			args: []string{"--", "--some-new-flag", "the-value"},
			want: []string{"claude", "--agent", "spacedock:first-officer", "--some-new-flag", "the-value", wantBootstrapPrompt},
		},
		{
			// The spacedock-injected `--plugin-dir <dir>` prefix lands the `exec`
			// subcommand at index 2; the leading-subcommand exemption must see it
			// THROUGH the prefix and stay silent. This case REDS if skipInjectedPrefix
			// is removed (the bare index-0 check then names `exec` as stray), so it
			// pins the structural skip as the load-bearing mechanism.
			name: "codex --plugin-dir then exec subcommand behind injected prefix",
			run: func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
				var stdout bytes.Buffer
				return runCodex(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
			},
			args: []string{"--plugin-dir", "/co", "--", "exec", "do the thing"},
			want: []string{"codex", "--plugin-dir", "/co", "exec", "do the thing", wantCodexBootstrapPrompt},
		},
		{
			// Same structural skip for the codex `resume` subcommand behind the
			// injected prefix — no stray warning, argv unchanged.
			name: "codex --plugin-dir then resume subcommand behind injected prefix",
			run: func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
				var stdout bytes.Buffer
				return runCodex(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
			},
			args: []string{"--plugin-dir", "/co", "--", "resume", "abc123"},
			want: []string{"codex", "--plugin-dir", "/co", "resume", "abc123", wantCodexBootstrapPrompt},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeHost{manifest: compatibleManifest(t)}
			var stderr bytes.Buffer
			code := tc.run(tc.args, t.TempDir(), fake, &stderr)
			if code != 0 {
				t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
			}
			if strings.Contains(stderr.String(), strayWarnMarker) {
				t.Fatalf("guard fired on a legitimate host positional: %q", stderr.String())
			}
			if !equalArgv(fake.launchedArg, tc.want) {
				t.Fatalf("launch argv = %v, want %v", fake.launchedArg, tc.want)
			}
		})
	}
}

// TestStrayPromptAfterDashClassifier (AC-3): the classifier over the after-`--`
// token grammar, independent of launch. Pins the boundary logic directly so a
// regression in flag-set membership or subcommand recognition fails a fast unit
// test, not only the launch-seam tests.
func TestStrayPromptAfterDashClassifier(t *testing.T) {
	cases := []struct {
		name           string
		passthrough    []string
		hasTask        bool
		host           string
		wantPositional string
		wantOK         bool
	}{
		{
			name:           "stray after value-flag pair",
			passthrough:    []string{"--model", "gpt-x", "@/tmp/handoff.md"},
			host:           "claude",
			wantPositional: "@/tmp/handoff.md",
			wantOK:         true,
		},
		{
			name:        "value of value-taking flag is not stray",
			passthrough: []string{"-p", "do the thing"},
			host:        "claude",
		},
		{
			name:        "leading subcommand arg is not stray",
			passthrough: []string{"exec", "do the thing"},
			host:        "codex",
		},
		{
			name:           "equals-form successor positional is stray",
			passthrough:    []string{"--model=gpt-x", "@/tmp/handoff.md"},
			host:           "claude",
			wantPositional: "@/tmp/handoff.md",
			wantOK:         true,
		},
		{
			name:        "hasTask short-circuit suppresses",
			passthrough: []string{"@/tmp/handoff.md"},
			hasTask:     true,
			host:        "claude",
		},
		{
			// skipInjectedPrefix strips the leading `--plugin-dir <dir>` pair so the
			// dir is NOT named and the operator's real prompt after it IS.
			name:           "spacedock-injected --plugin-dir prefix does not shadow the real prompt",
			passthrough:    []string{"--plugin-dir", "/co", "--model", "gpt-x", "@/tmp/handoff.md"},
			host:           "claude",
			wantPositional: "@/tmp/handoff.md",
			wantOK:         true,
		},
		{
			// A leading subcommand BEHIND the injected `--plugin-dir <dir>` prefix must
			// still be recognized as a subcommand (its args legitimate). This reds if
			// skipInjectedPrefix is removed — the bare index-0 check then names `exec`.
			name:        "subcommand behind injected --plugin-dir prefix is not stray",
			passthrough: []string{"--plugin-dir", "/co", "exec", "do the thing"},
			host:        "codex",
		},
		{
			// Multiple injected `--plugin-dir <dir>` pairs are all skipped.
			name:        "multiple injected --plugin-dir pairs all skipped before subcommand",
			passthrough: []string{"--plugin-dir", "/a", "--plugin-dir", "/b", "exec", "p"},
			host:        "codex",
		},
		{
			// A value of an UNRECOGNIZED `-`-prefixed flag is ambiguous; the
			// conservative fallback suppresses rather than prescribe wrong advice.
			name:        "value of unrecognized flag suppresses (no wrong advice)",
			passthrough: []string{"--some-new-flag", "the-value"},
			host:        "claude",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fd := frontDoorArgs{passthrough: tc.passthrough, hasTask: tc.hasTask}
			pos, ok := strayPromptAfterDash(fd, tc.host)
			if ok != tc.wantOK || pos != tc.wantPositional {
				t.Fatalf("strayPromptAfterDash(%v, %q) = (%q, %v), want (%q, %v)",
					tc.passthrough, tc.host, pos, ok, tc.wantPositional, tc.wantOK)
			}
		})
	}
}
