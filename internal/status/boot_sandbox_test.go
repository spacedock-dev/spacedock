// ABOUTME: AC-1 — status --boot emits a SANDBOX: section (and --boot --json a
// ABOUTME: sandbox field) reporting whether THIS process is sandboxed.
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

// sandboxBootFixture builds a git-rooted workflow and returns the workflow dir
// plus a request env pinning the two independent inputs the SANDBOX state turns
// on: `inside` sets the APP_SANDBOX_CONTAINER_ID signal proving this process is
// already running inside the safehouse sandbox, and `available` pins whether the
// safehouse binary resolves on the request PATH.
//
// It ALSO writes a .safehouse profile at the repo root in every case — deliberately.
// The profile is a launch fact, and boot reports a session; writing it
// unconditionally means a re-collapse of SessionState back into the launch
// renderer would change these rendered strings and turn the assertions red.
func sandboxBootFixture(t *testing.T, inside, available bool) (root string, env []string) {
	t.Helper()
	root = t.TempDir()
	gitC(t, root, "init")
	writeFile(t, filepath.Join(root, "README.md"), sandboxBootReadme)
	writeFile(t, filepath.Join(root, ".safehouse"), "profile")
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
	if inside {
		env = append(env, "APP_SANDBOX_CONTAINER_ID=agent-safehouse")
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

// TestBootSandboxSection (AC-1) drives --boot over the inside-a-sandbox signal
// and the binary found/not-found, and asserts the exact SANDBOX: line for each
// combination. The expected state strings are independent test-supplied values.
//
// The REGRESSION row is `inside-safehouse-absent-from-PATH`: that is the live
// configuration this fix exists for, and it is the exact state the old renderer
// called `unavailable (safehouse not on PATH)` — safehouse being off PATH
// precisely BECAUSE the wrap already happened.
func TestBootSandboxSection(t *testing.T) {
	cases := []struct {
		name      string
		inside    bool
		available bool
		want      string
	}{
		{"inside-safehouse-absent-from-PATH", true, false, "SANDBOX: inside (agent-safehouse)"},
		{"inside-safehouse-on-PATH", true, true, "SANDBOX: inside (agent-safehouse)"},
		{"not-sandboxed-available", false, true, "SANDBOX: not sandboxed (safehouse available)"},
		{"not-sandboxed-not-installed", false, false, "SANDBOX: not sandboxed (safehouse not installed)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, env := sandboxBootFixture(t, tc.inside, tc.available)
			out, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--boot")
			if code != 0 {
				t.Fatalf("--boot exit=%d stderr=%q", code, errOut)
			}
			got := sandboxLineOf(t, out)
			if got != tc.want {
				t.Fatalf("SANDBOX line = %q, want %q", got, tc.want)
			}
			// AC-2b: a .safehouse profile is present in every fixture, and boot
			// must never mention it — boot reports a session, not a launch.
			if strings.Contains(got, ".safehouse") {
				t.Fatalf("SANDBOX line = %q, which mentions a .safehouse profile — a launch fact on a session surface", got)
			}
		})
	}
}

// TestBootSandboxJSONField (AC-1) asserts the --boot --json form's `sandbox`
// string field carries the same session state. This is the field that lands in
// the First Officer's DURABLE boot evidence, which is why the inside-and-absent
// corner mattered most: it wrote a false posture into machine-read state.
func TestBootSandboxJSONField(t *testing.T) {
	cases := []struct {
		name      string
		inside    bool
		available bool
		want      string
	}{
		{"inside-safehouse-absent-from-PATH", true, false, "inside (agent-safehouse)"},
		{"not-sandboxed-not-installed", false, false, "not sandboxed (safehouse not installed)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, env := sandboxBootFixture(t, tc.inside, tc.available)
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
