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
