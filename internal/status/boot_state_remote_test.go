// ABOUTME: AC-1/AC-4 boot coverage — STATE_BACKEND surfaces state-remote
// ABOUTME: availability (origin vs none) for split-root, and omits it for single-root.
package status

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/testgit"
)

// initStateRepoWithOrigin turns an existing state checkout dir into a real git
// repo with a named `origin` remote pointing at a throwaway bare upstream, so
// stateHasOrigin reports true. The remote is never contacted (the probe is
// `remote get-url origin`, network-free), so the bare upstream needs no content.
func initStateRepoWithOrigin(t *testing.T, stateDir string) {
	t.Helper()
	testgit.InitRepo(t, stateDir, "-q")
	upstream := filepath.Join(t.TempDir(), "upstream.git")
	gitC(t, t.TempDir(), "init", "-q", "--bare", upstream)
	gitC(t, stateDir, "remote", "add", "origin", upstream)
}

// TestBootTextStateRemoteNone (AC-1, AC-4) asserts a split-root state checkout
// with NO origin remote names the local-only mode on the text STATE_BACKEND line.
func TestBootTextStateRemoteNone(t *testing.T) {
	def, _ := buildSplitRoot(t, splitRootReadme, map[string]string{
		"add-login.md": "---\nstatus: ideation\n---\n",
	})
	env := pinnedEnv(t)

	out, stderr, code := runNative(t, def, env, "--workflow-dir", def, "--boot")
	if code != 0 {
		t.Fatalf("--boot exit=%d stderr=%q", code, stderr)
	}
	line := stateBackendLineOf(t, out)
	if !strings.Contains(line, "remote: none") {
		t.Fatalf("STATE_BACKEND line missing `remote: none`: %q", line)
	}
	if !strings.Contains(line, "state not remotely synced") {
		t.Fatalf("STATE_BACKEND line missing not-remotely-synced phrase: %q", line)
	}
}

// TestBootTextStateRemoteOrigin (AC-1) asserts a split-root state checkout WITH
// an origin remote reports `remote: origin` and does NOT carry the not-synced
// phrase.
func TestBootTextStateRemoteOrigin(t *testing.T) {
	def, state := buildSplitRoot(t, splitRootReadme, map[string]string{
		"add-login.md": "---\nstatus: ideation\n---\n",
	})
	initStateRepoWithOrigin(t, state)
	env := pinnedEnv(t)

	out, stderr, code := runNative(t, def, env, "--workflow-dir", def, "--boot")
	if code != 0 {
		t.Fatalf("--boot exit=%d stderr=%q", code, stderr)
	}
	line := stateBackendLineOf(t, out)
	if !strings.Contains(line, "remote: origin") {
		t.Fatalf("STATE_BACKEND line missing `remote: origin`: %q", line)
	}
	if strings.Contains(line, "state not remotely synced") {
		t.Fatalf("STATE_BACKEND line leaks not-remotely-synced phrase under origin: %q", line)
	}
}

// TestBootTextSingleRootNoRemoteClause (AC-1 negative) asserts a single-root
// workflow's STATE_BACKEND line carries NO remote clause at all — the degrade is
// split-root-only.
func TestBootTextSingleRootNoRemoteClause(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "sdb32-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)

	out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--boot")
	if code != 0 {
		t.Fatalf("--boot exit=%d stderr=%q", code, stderr)
	}
	line := stateBackendLineOf(t, out)
	if strings.Contains(line, "remote:") {
		t.Fatalf("single-root STATE_BACKEND line leaks a remote clause: %q", line)
	}
}

// TestBootJSONStateRemoteNone (AC-1, AC-4) asserts the JSON envelope carries
// state_remote "none" for a no-origin split-root checkout, positioned AFTER
// entity_dir_present so the FO's key-order parse is preserved.
func TestBootJSONStateRemoteNone(t *testing.T) {
	def, _ := buildSplitRoot(t, splitRootReadme, map[string]string{
		"add-login.md": "---\nstatus: ideation\n---\n",
	})
	env := pinnedEnv(t)

	out, stderr, code := runNative(t, def, env, "--workflow-dir", def, "--boot", "--json")
	if code != 0 {
		t.Fatalf("--boot --json exit=%d stderr=%q", code, stderr)
	}

	var b struct {
		StateRemote string `json:"state_remote"`
	}
	if err := json.Unmarshal([]byte(out), &b); err != nil {
		t.Fatalf("parse --boot --json: %v\n%s", err, out)
	}
	if b.StateRemote != "none" {
		t.Fatalf("state_remote = %q, want none", b.StateRemote)
	}

	// Key order: state_remote must appear AFTER entity_dir_present.
	presentIdx := strings.Index(out, `"entity_dir_present"`)
	remoteIdx := strings.Index(out, `"state_remote"`)
	if presentIdx < 0 || remoteIdx < 0 {
		t.Fatalf("--boot --json missing entity_dir_present or state_remote\n%s", out)
	}
	if remoteIdx < presentIdx {
		t.Fatalf("state_remote must follow entity_dir_present\n%s", out)
	}
}

// TestBootJSONStateRemoteOrigin (AC-1) asserts the JSON envelope carries
// state_remote "origin" for an origin-backed split-root checkout.
func TestBootJSONStateRemoteOrigin(t *testing.T) {
	def, state := buildSplitRoot(t, splitRootReadme, map[string]string{
		"add-login.md": "---\nstatus: ideation\n---\n",
	})
	initStateRepoWithOrigin(t, state)
	env := pinnedEnv(t)

	out, stderr, code := runNative(t, def, env, "--workflow-dir", def, "--boot", "--json")
	if code != 0 {
		t.Fatalf("--boot --json exit=%d stderr=%q", code, stderr)
	}
	var b struct {
		StateRemote string `json:"state_remote"`
	}
	if err := json.Unmarshal([]byte(out), &b); err != nil {
		t.Fatalf("parse --boot --json: %v\n%s", err, out)
	}
	if b.StateRemote != "origin" {
		t.Fatalf("state_remote = %q, want origin", b.StateRemote)
	}
}

// TestBootJSONSingleRootNoStateRemote (AC-1 negative) asserts a single-root
// workflow's JSON envelope OMITS state_remote entirely.
func TestBootJSONSingleRootNoStateRemote(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "sdb32-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)

	out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--boot", "--json")
	if code != 0 {
		t.Fatalf("--boot --json exit=%d stderr=%q", code, stderr)
	}
	if strings.Contains(out, `"state_remote"`) {
		t.Fatalf("single-root --boot --json leaks state_remote key\n%s", out)
	}
}

// stateBackendLineOf extracts the STATE_BACKEND line from --boot text output.
func stateBackendLineOf(t *testing.T, out string) string {
	t.Helper()
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "STATE_BACKEND:") {
			return line
		}
	}
	t.Fatalf("--boot output missing STATE_BACKEND line\n%s", out)
	return ""
}
