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
	want := []string{"claude", "--agent", "spacedock:first-officer", "--permission-mode", "auto", "--model", "gpt-x", "@/tmp/handoff.md", wantBootstrapPrompt}
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
	want := []string{"claude", "--agent", "spacedock:first-officer", "--permission-mode", "auto", "--plugin-dir", "/co", "--model", "gpt-x", "@/tmp/handoff.md", wantBootstrapPrompt}
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
			want: []string{"claude", "--agent", "spacedock:first-officer", "--permission-mode", "auto", "-p", "do the thing", wantBootstrapPrompt},
		},
		{
			name: "claude hasTask short-circuit",
			run: func(args []string, dir string, fake *fakeHost, stderr *bytes.Buffer) int {
				var stdout bytes.Buffer
				return runClaude(context.Background(), args, dir, fake, lookFound, &stdout, stderr)
			},
			args: []string{"task before", "--", "@/tmp/handoff.md"},
			want: []string{"claude", "--agent", "spacedock:first-officer", "--permission-mode", "auto", "@/tmp/handoff.md", wantBootstrapPrompt + " task before"},
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
			want: []string{"claude", "--agent", "spacedock:first-officer", "--permission-mode", "auto", "--some-new-flag", "the-value", wantBootstrapPrompt},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("CODEX_HOME", t.TempDir()) // isolate any codex --plugin-dir local marketplace
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

// TestStrayPromptAfterDashClassifier (AC-3) pins Claude's advisory-only after-`--`
// token grammar directly, so a value-flag regression fails a fast unit test rather
// than only a launch-seam test.
func TestStrayPromptAfterDashClassifier(t *testing.T) {
	cases := []struct {
		name           string
		passthrough    []string
		hasTask        bool
		wantPositional string
		wantOK         bool
	}{
		{
			name:           "stray after value-flag pair",
			passthrough:    []string{"--model", "gpt-x", "@/tmp/handoff.md"},
			wantPositional: "@/tmp/handoff.md",
			wantOK:         true,
		},
		{
			name:        "value of value-taking flag is not stray",
			passthrough: []string{"-p", "do the thing"},
		},
		{
			name:           "equals-form successor positional is stray",
			passthrough:    []string{"--model=gpt-x", "@/tmp/handoff.md"},
			wantPositional: "@/tmp/handoff.md",
			wantOK:         true,
		},
		{
			name:        "hasTask short-circuit suppresses",
			passthrough: []string{"@/tmp/handoff.md"},
			hasTask:     true,
		},
		{
			// skipInjectedPrefix strips the leading `--plugin-dir <dir>` pair so the
			// dir is NOT named and the operator's real prompt after it IS.
			name:           "spacedock-injected --plugin-dir prefix does not shadow the real prompt",
			passthrough:    []string{"--plugin-dir", "/co", "--model", "gpt-x", "@/tmp/handoff.md"},
			wantPositional: "@/tmp/handoff.md",
			wantOK:         true,
		},
		{
			// A value of an UNRECOGNIZED `-`-prefixed flag is ambiguous; the
			// conservative fallback suppresses rather than prescribe wrong advice.
			name:        "value of unrecognized flag suppresses (no wrong advice)",
			passthrough: []string{"--some-new-flag", "the-value"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fd := frontDoorArgs{passthrough: tc.passthrough, hasTask: tc.hasTask}
			pos, ok := strayPromptAfterDash(fd)
			if ok != tc.wantOK || pos != tc.wantPositional {
				t.Fatalf("strayPromptAfterDash(%v) = (%q, %v), want (%q, %v)",
					tc.passthrough, pos, ok, tc.wantPositional, tc.wantOK)
			}
		})
	}
}
