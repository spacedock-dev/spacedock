package cli

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

// TestTopLevelHelpFormsAreIdentical pins AC-1's invariant that bare `spacedock`,
// `-h`, and the `help` subcommand render byte-identical output to `--help`.
func TestTopLevelHelpFormsAreIdentical(t *testing.T) {
	ref := helpStdout(t, "--help")
	wantGateRow := "  gate        prepare | withdraw | record | validate | consume\n                                      Record, inspect, or consume durable gate resolutions\n"
	if strings.Count(ref, wantGateRow) != 1 {
		t.Fatalf("top-level help gate row differs from the published contract:\n%s", ref)
	}
	for _, form := range [][]string{nil, {"-h"}, {"help"}} {
		got := helpStdout(t, form...)
		if got != ref {
			t.Errorf("help form %v output differs from --help:\n--- form ---\n%s\n--- --help ---\n%s", form, got, ref)
		}
	}
}

func helpStdout(t *testing.T, args ...string) string {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := Run(args, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("Run(%v) = %d, want 0 (stderr=%q)", args, code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("Run(%v) stderr = %q, want empty", args, stderr.String())
	}
	return stdout.String()
}

func TestVersion(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"--version"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("Run returned %d, want 0", code)
	}
	// AC-4: the FIRST line is the load-bearing, FO-parsed version line and carries
	// the version token and NOTHING after it — the frozen contract token moved
	// below (asserted in version_session_test.go, which also asserts the session
	// lines). Leaving the token on line 1 turns the regex red.
	want := "spacedock " + displayVersion()
	got := strings.SplitN(stdout.String(), "\n", 2)[0]
	if got != want {
		t.Fatalf("version first line = %q, want %q", got, want)
	}
	if !regexp.MustCompile(`^spacedock \S+$`).MatchString(got) {
		t.Fatalf("version first line = %q, want to match ^spacedock \\S+$ (the shape the FO version gate parses)", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestUnknownCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"bogus"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("Run returned %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "unknown command: bogus") {
		t.Fatalf("stderr missing unknown-command message: %q", stderr.String())
	}
}

// TestLeadingUnknownFlagIsUsageError pins AC-6's leading-flag half: a stray
// leading unknown flag that resolves to no valid subcommand (before any command
// token) must print the usage block to stderr and exit NON-ZERO — never the
// silent exit-0+help the UnknownFlags whitelist would otherwise yield. The
// space-form `--foo install` (where `--foo` would consume `install` as its value)
// must NOT silently swallow the valid command token either. Distinct from the
// bare invocation, which still prints help + exit 0 (TestLeadingFlagBareStillHelp).
func TestLeadingUnknownFlagIsUsageError(t *testing.T) {
	for _, args := range [][]string{
		{"--bogus"},
		{"--boot"},
		{"--foo", "install"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Run(args, &stdout, &stderr)

			if code == 0 {
				t.Fatalf("Run(%v) exited 0 on a leading unknown flag; want non-zero usage error (stdout=%q)", args, stdout.String())
			}
			if !strings.Contains(stderr.String(), tagline) {
				t.Fatalf("Run(%v) stderr missing the usage block (tagline): %q", args, stderr.String())
			}
		})
	}
}

// TestLeadingFlagBareStillHelp guards the regression boundary: bare `spacedock`
// (no args) still renders help to stdout and exits 0 — the leading-unknown-flag
// fix must not turn the legitimate bare invocation into a usage error.
func TestLeadingFlagBareStillHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run(nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("bare spacedock exited %d, want 0 (stderr=%q)", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), tagline) {
		t.Fatalf("bare spacedock stdout missing help: %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("bare spacedock stderr = %q, want empty", stderr.String())
	}
}

// TestUnknownCommandWithFlag pins AC-6: an unknown subcommand carrying a trailing
// flag must STILL print `unknown command: <name>` + the grouped help to stderr and
// exit 2 — never a silent exit 2. The removed `spacedock init --host claude` is the
// captain's live case: with flag parsing enabled the root flagset errored on the
// unknown --host before RunE ran, silencing all output. The bogus-with-flag case is
// the generalized form.
func TestUnknownCommandWithFlag(t *testing.T) {
	for _, args := range [][]string{
		{"init", "--host", "claude"},
		{"bogus", "--someflag"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Run(args, &stdout, &stderr)

			if code != 2 {
				t.Fatalf("Run(%v) = %d, want 2", args, code)
			}
			out := stderr.String()
			if !strings.Contains(out, "unknown command: "+args[0]) {
				t.Fatalf("Run(%v) stderr missing %q: %q", args, "unknown command: "+args[0], out)
			}
			// The grouped help block follows the diagnostic — the tagline pins it.
			if !strings.Contains(out, tagline) {
				t.Fatalf("Run(%v) stderr missing the usage block (tagline): %q", args, out)
			}
		})
	}
}
