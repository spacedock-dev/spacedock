// ABOUTME: Unit tests for the shared safehouse seam: Present detection,
// ABOUTME: Available lookPath gate + pinned hint, and inner-argv-agnostic Wrap.
package safehouse

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPresentDetectsFile: a `.safehouse` regular file in workdir is Present.
func TestPresentDetectsFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".safehouse"), []byte("profile"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !Present(dir) {
		t.Fatalf("Present(%q) = false, want true for a .safehouse file", dir)
	}
}

// TestPresentDetectsDir: a `.safehouse` directory counts identically to a file
// (os.Stat truthiness), pinning the file-vs-dir case from AC-1.
func TestPresentDetectsDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, ".safehouse"), 0o755); err != nil {
		t.Fatal(err)
	}
	if !Present(dir) {
		t.Fatalf("Present(%q) = false, want true for a .safehouse directory", dir)
	}
}

// TestPresentAbsent: no `.safehouse` entry → not Present.
func TestPresentAbsent(t *testing.T) {
	dir := t.TempDir()
	if Present(dir) {
		t.Fatalf("Present(%q) = true, want false with no .safehouse", dir)
	}
}

// TestAvailableResolvable: a lookPath that resolves the binary → ok, no hint.
func TestAvailableResolvable(t *testing.T) {
	look := func(string) (string, error) { return "/usr/bin/safehouse", nil }
	ok, hint := Available(look)
	if !ok {
		t.Fatalf("Available = (false, %q), want ok=true when lookPath resolves", hint)
	}
	if hint != "" {
		t.Fatalf("Available hint = %q, want empty when resolvable", hint)
	}
}

// TestAvailableNotFound: a lookPath that fails → not ok, with a pinned install
// hint that cites safehouse's real third-party-tap brew install command and its
// canonical site, and never the outdated github.com/anthropics/safehouse URL.
func TestAvailableNotFound(t *testing.T) {
	look := func(string) (string, error) { return "", errors.New("not found") }
	ok, hint := Available(look)
	if ok {
		t.Fatalf("Available ok=true, want false when lookPath fails")
	}
	if !strings.Contains(hint, "brew install eugene1g/safehouse/agent-safehouse") {
		t.Errorf("Available hint %q does not name the brew install command", hint)
	}
	if !strings.Contains(hint, "https://agent-safehouse.dev") {
		t.Errorf("Available hint %q does not name the canonical safehouse site", hint)
	}
	if strings.Contains(hint, "github.com/anthropics/safehouse") {
		t.Errorf("Available hint %q still cites the outdated github.com/anthropics/safehouse URL", hint)
	}
}

// TestWrapComposesPrefix: Wrap prepends the safehouse prefix and the `--`
// separator, then the inner argv, with no extra args.
func TestWrapComposesPrefix(t *testing.T) {
	clearTerminalTargetingEnv(t)
	inner := []string{"claude", "--agent", "spacedock:first-officer", "--foo"}
	got := Wrap(inner, nil)
	want := []string{"safehouse", "--trust-workdir-config", "--",
		"claude", "--agent", "spacedock:first-officer", "--foo"}
	if !equalArgv(got, want) {
		t.Fatalf("Wrap = %v, want %v", got, want)
	}
}

// TestWrapWithExtra: extra safehouse args land between --trust-workdir-config
// and the `--` separator, never mixed into the inner command.
func TestWrapWithExtra(t *testing.T) {
	clearTerminalTargetingEnv(t)
	got := Wrap([]string{"claude"}, []string{"--profile", "x"})
	want := []string{"safehouse", "--trust-workdir-config", "--profile", "x", "--", "claude"}
	if !equalArgv(got, want) {
		t.Fatalf("Wrap with extra = %v, want %v", got, want)
	}
}

// TestWrapAddsTerminalTargetingEnvArgument reshapes the allowance test to the
// tmux pair: bare ZELLIJ is gone from the gate, and the built-in allowance now
// covers every host the consumer probes, not just Zellij.
func TestWrapAddsTerminalTargetingEnvArgument(t *testing.T) {
	clearTerminalTargetingEnv(t)
	t.Setenv("TMUX", "/tmp/tmux-501/default,12345,0")
	t.Setenv("TMUX_PANE", "%3")

	got := Wrap([]string{"claude"}, []string{"--env-pass", "SPACEDOCK_BIN"})
	want := []string{"safehouse", "--trust-workdir-config", "--env-pass=TMUX,TMUX_PANE", "--env-pass", "SPACEDOCK_BIN", "--", "claude"}
	if !equalArgv(got, want) {
		t.Fatalf("Wrap = %v, want %v", got, want)
	}
}

// TestTerminalEnvPassArgsComposer is the AC-2 oracle: the argv composer names
// a variable only when it is present in the (stubbed) parent, in probe order,
// and an empty parent yields no allowance at all — never a fixed list, never
// gated on a different variable than the one being passed.
func TestTerminalEnvPassArgsComposer(t *testing.T) {
	cases := []struct {
		name   string
		lookup func(string) (string, bool)
		want   []string
	}{
		{
			name:   "empty parent yields no allowance",
			lookup: mapLookup(nil),
			want:   nil,
		},
		{
			name: "tmux pair only",
			lookup: mapLookup(map[string]string{
				"TMUX":      "/tmp/tmux-501/default,12345,0",
				"TMUX_PANE": "%3",
			}),
			want: []string{"--env-pass=TMUX,TMUX_PANE"},
		},
		{
			name: "all nine set, emitted in probe order",
			lookup: mapLookup(map[string]string{
				"ZELLIJ_SESSION_NAME": "excellent-pheasant",
				"ZELLIJ_PANE_ID":      "51",
				"TMUX":                "/tmp/tmux-501/default,12345,0",
				"TMUX_PANE":           "%3",
				"HERDR_ENV":           "1",
				"HERDR_PANE_ID":       "7",
				"CMUX_WORKSPACE_ID":   "ws-1",
				"CMUX_SURFACE_ID":     "sf-1",
				"TERM_PROGRAM":        "ghostty",
			}),
			want: []string{"--env-pass=ZELLIJ_SESSION_NAME,ZELLIJ_PANE_ID,TMUX,TMUX_PANE,HERDR_ENV,HERDR_PANE_ID,CMUX_WORKSPACE_ID,CMUX_SURFACE_ID,TERM_PROGRAM"},
		},
		{
			name: "set-but-empty name is still named",
			lookup: mapLookup(map[string]string{
				"TMUX": "",
			}),
			want: []string{"--env-pass=TMUX"},
		},
		{
			name: "a name outside the list is never named",
			lookup: mapLookup(map[string]string{
				"NOT_A_TERMINAL_SIGNAL": "1",
			}),
			want: nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := terminalEnvPassArgs(tc.lookup)
			if !equalArgv(got, tc.want) {
				t.Fatalf("terminalEnvPassArgs() = %v, want %v", got, tc.want)
			}
		})
	}
}

// mapLookup adapts a plain map to the (string) (string, bool) lookup shape
// terminalEnvPassArgs takes, so the composer's table test is deterministic
// under any developer terminal — it never touches real process env.
func mapLookup(m map[string]string) func(string) (string, bool) {
	return func(name string) (string, bool) {
		v, ok := m[name]
		return v, ok
	}
}

// clearTerminalTargetingEnv unsets all nine signals terminalHostEnvVars
// probes, so tests that exercise Wrap's built-in allowance start from a clean
// baseline regardless of the developer's own terminal host.
func clearTerminalTargetingEnv(t *testing.T) {
	t.Helper()
	for _, key := range terminalHostEnvVars {
		value, present := os.LookupEnv(key)
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if present {
				_ = os.Setenv(key, value)
				return
			}
			_ = os.Unsetenv(key)
		})
	}
}

// TestTranslateFlags pins the de-prefixed-knob → safehouse-extra translation:
// `enable=ssh,docker` comma-splits into repeated `--enable=KEY`; `add-dirs=P` and
// `add-dirs-ro=P` map to `--add-dirs=P` / `--add-dirs-ro=P`; an unknown key is a
// hard error. The translator holds NO `--safehouse-` namespace knowledge — the
// dispatcher in internal/cli strips the prefix before calling it (AC-3).
func TestTranslateFlags(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want []string
	}{
		{"empty", nil, nil},
		{"enable-single", []string{"enable=docker"}, []string{"--enable=docker"}},
		{"enable-comma-split", []string{"enable=ssh,docker"}, []string{"--enable=ssh", "--enable=docker"}},
		{"add-dirs", []string{"add-dirs=/a"}, []string{"--add-dirs=/a"}},
		{"add-dirs-ro", []string{"add-dirs-ro=/b"}, []string{"--add-dirs-ro=/b"}},
		{"order-preserved", []string{"enable=ssh", "add-dirs=/a", "add-dirs-ro=/b"},
			[]string{"--enable=ssh", "--add-dirs=/a", "--add-dirs-ro=/b"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := TranslateFlags(tc.in)
			if err != nil {
				t.Fatalf("TranslateFlags(%v) err = %v, want nil", tc.in, err)
			}
			if !equalArgv(got, tc.want) {
				t.Fatalf("TranslateFlags(%v) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

// TestTranslateFlagsUnknownKeyErrors: an unrecognized de-prefixed key is a hard
// error so a typo can never silently fall through to the host (AC-8).
func TestTranslateFlagsUnknownKeyErrors(t *testing.T) {
	if _, err := TranslateFlags([]string{"bogus=x"}); err == nil {
		t.Fatalf("TranslateFlags([bogus=x]) err = nil, want a hard error for an unknown key")
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
