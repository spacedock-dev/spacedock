// ABOUTME: Locks the env-aware launcher command token generated into dispatch prompts.
// ABOUTME: Ensures ensign fetch commands prefer executable SPACEDOCK_BIN with spacedock PATH fallback.
package dispatch

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestLauncherCommandUsesExecutableSpacedockBinWhenAvailable(t *testing.T) {
	bin := writeExecutable(t, t.TempDir(), "spacedock-bin", "#!/bin/sh\necho env-bin:$1\n")
	out := runLauncherCommand(t, []string{"SPACEDOCK_BIN=" + bin}, nil, "status")

	if got, want := strings.TrimSpace(out), "env-bin:status"; got != want {
		t.Fatalf("launcher command output = %q, want %q", got, want)
	}
}

func TestLauncherCommandFallsBackToPathWhenSpacedockBinUnsetEmptyOrUnusable(t *testing.T) {
	pathDir := t.TempDir()
	writeExecutable(t, pathDir, "spacedock", "#!/bin/sh\necho path-bin:$1\n")
	nonExecutable := filepath.Join(t.TempDir(), "not-executable")
	if err := os.WriteFile(nonExecutable, []byte("not executable"), 0o644); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(t.TempDir(), "missing-spacedock")

	for _, tc := range []struct {
		name string
		env  []string
	}{
		{name: "unset"},
		{name: "empty", env: []string{"SPACEDOCK_BIN="}},
		{name: "non-executable", env: []string{"SPACEDOCK_BIN=" + nonExecutable}},
		{name: "missing", env: []string{"SPACEDOCK_BIN=" + missing}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out := runLauncherCommand(t, tc.env, []string{pathDir}, "dispatch")
			if got, want := strings.TrimSpace(out), "path-bin:dispatch"; got != want {
				t.Fatalf("launcher command output = %q, want %q", got, want)
			}
		})
	}
}

func runLauncherCommand(t *testing.T, env []string, pathDirs []string, arg string) string {
	t.Helper()
	cmd := exec.Command("sh", "-c", launcherCommand()+" "+arg)
	cmd.Env = append(os.Environ(), env...)
	if pathDirs != nil {
		cmd.Env = append(cmd.Env, "PATH="+strings.Join(pathDirs, string(os.PathListSeparator)))
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("launcher command failed: %v\n%s", err, out)
	}
	return string(out)
}

func writeExecutable(t *testing.T, dir, name, body string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}
