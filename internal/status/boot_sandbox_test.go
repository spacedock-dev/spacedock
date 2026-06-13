// ABOUTME: AC-2 — status --boot emits a SANDBOX: section (and --boot --json a
// ABOUTME: sandbox field) whose state matches the profile-present/binary inputs.
package status

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sandboxBootReadme = `---
commissioned-by: spacedock@1
id-style: slug
stages:
  states:
    - name: build
      initial: true
---

# Sandbox Boot Workflow
`

// sandboxBootFixture builds a git-rooted workflow, optionally writing a .safehouse
// profile at the repo root, and returns the workflow dir plus a PATH whose
// safehouse-binary resolution is pinned by `available`. The .safehouse profile and
// the safehouse binary are the two independent inputs the SANDBOX state turns on.
func sandboxBootFixture(t *testing.T, profilePresent, available bool) (root string, env []string) {
	t.Helper()
	root = t.TempDir()
	gitC(t, root, "init")
	writeFile(t, filepath.Join(root, "README.md"), sandboxBootReadme)
	if profilePresent {
		writeFile(t, filepath.Join(root, ".safehouse"), "profile")
	}
	pathDir := t.TempDir()
	if available {
		bin := filepath.Join(pathDir, "safehouse")
		if err := os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	env = append(pinnedEnv(t), "")
	// Replace PATH so safehouse resolution is controlled, not read from the machine.
	for i, kv := range env {
		if strings.HasPrefix(kv, "PATH=") {
			env[i] = "PATH=" + pathDir
		}
	}
	if env[len(env)-1] == "" {
		env = env[:len(env)-1]
	}
	return root, env
}

// sandboxLineOf returns the whole SANDBOX: line from a --boot text body.
func sandboxLineOf(t *testing.T, boot string) string {
	t.Helper()
	for _, line := range strings.Split(boot, "\n") {
		if strings.HasPrefix(line, "SANDBOX:") {
			return line
		}
	}
	t.Fatalf("--boot output has no SANDBOX: line:\n%s", boot)
	return ""
}

// TestBootSandboxSection (AC-2) drives --boot over a fixture with the .safehouse
// profile present/absent and the binary found/not-found, and asserts the exact
// SANDBOX: line for each combination. The expected state strings are independent
// test-supplied values.
func TestBootSandboxSection(t *testing.T) {
	cases := []struct {
		name           string
		profilePresent bool
		available      bool
		want           string
	}{
		{"enabled", true, true, "SANDBOX: enabled (safehouse)"},
		{"available-not-enabled", false, true, "SANDBOX: available, not enabled (no .safehouse profile)"},
		{"unavailable-with-profile", true, false, "SANDBOX: unavailable (safehouse not on PATH)"},
		{"unavailable-no-profile", false, false, "SANDBOX: unavailable (safehouse not on PATH)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, env := sandboxBootFixture(t, tc.profilePresent, tc.available)
			out, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--boot")
			if code != 0 {
				t.Fatalf("--boot exit=%d stderr=%q", code, errOut)
			}
			if got := sandboxLineOf(t, out); got != tc.want {
				t.Fatalf("SANDBOX line = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestBootSandboxJSONField (AC-2) asserts the --boot --json form gains a `sandbox`
// string field carrying the same three-way state, for the present+available and
// absent+unavailable corners.
func TestBootSandboxJSONField(t *testing.T) {
	cases := []struct {
		name           string
		profilePresent bool
		available      bool
		want           string
	}{
		{"enabled", true, true, "enabled (safehouse)"},
		{"unavailable", false, false, "unavailable (safehouse not on PATH)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, env := sandboxBootFixture(t, tc.profilePresent, tc.available)
			out, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--boot", "--json")
			if code != 0 {
				t.Fatalf("--boot --json exit=%d stderr=%q", code, errOut)
			}
			var boot struct {
				Sandbox string `json:"sandbox"`
			}
			if err := json.Unmarshal([]byte(out), &boot); err != nil {
				t.Fatalf("parse --boot --json: %v\n%s", err, out)
			}
			if boot.Sandbox != tc.want {
				t.Fatalf("sandbox field = %q, want %q", boot.Sandbox, tc.want)
			}
		})
	}
}
