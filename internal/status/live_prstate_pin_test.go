// ABOUTME: Shallow-boot accuracy pin — pr_state.entries[].state in `status --boot
// ABOUTME: --json` reflects LIVE gh merge state (gh present) or an explicit unknown (gh absent).
package status

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// writeStubGh writes a `gh` shim into a fresh temp dir that prints the given
// merge state for `gh pr view ... --json state --jq .state` and returns the dir.
// The shim lets the offline pin drive checkPRStates' live `gh pr view` shell-out
// deterministically — no network, no real PR.
func writeStubGh(t *testing.T, state string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("gh stub shim is a POSIX shell script")
	}
	dir := t.TempDir()
	script := "#!/bin/sh\n" +
		"# Stub gh: emit a fixed PR state for `gh pr view ... --json state --jq .state`.\n" +
		"echo " + state + "\n"
	path := filepath.Join(dir, "gh")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

// bootWithPATH runs `status --boot --json` over the workflow root with PATH set to
// pathValue in BOTH the process env (so checkPRStates' bare exec.Command("gh")
// resolves the shim) and the Request env (so lookupExecutable finds it). It returns
// the parsed pr_state section.
func bootPRState(t *testing.T, root, pathValue string) (status string, entries []map[string]string) {
	t.Helper()
	t.Setenv("PATH", pathValue) // checkPRStates runs exec.Command("gh") against the process PATH
	env := []string{
		"PYTHONUTF8=1",
		"LANG=C.UTF-8",
		"LC_ALL=C.UTF-8",
		"USER=pinned-actor",
		"HOME=" + t.TempDir(),
		"PATH=" + pathValue,
	}
	out, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--boot", "--json")
	if code != 0 {
		t.Fatalf("--boot --json exit=%d stderr=%q", code, errOut)
	}
	var boot struct {
		PRState struct {
			Status  string              `json:"status"`
			Entries []map[string]string `json:"entries"`
		} `json:"pr_state"`
	}
	if err := json.Unmarshal([]byte(out), &boot); err != nil {
		t.Fatalf("parse --boot --json: %v\n%s", err, out)
	}
	return boot.PRState.Status, boot.PRState.Entries
}

// TestBootPRStateCarriesLiveMergeState is the shallow-boot accuracy pin: with `gh`
// on PATH, a PR-bearing non-terminal entity's pr_state entry reflects the LIVE
// merge state (`gh pr view` → MERGED), not just the stored `pr:` field. The
// shallow-boot greet's accurate local state summary rests on this; a regression
// that dropped the live `gh pr view` and echoed only the stored field would not
// report a freshly-merged PR accurately. The stub `gh` makes the live state
// deterministic and offline.
func TestBootPRStateCarriesLiveMergeState(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "live-prstate-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	stubDir := writeStubGh(t, "MERGED")

	status, entries := bootPRState(t, root, stubDir)
	if status != "ok" {
		t.Fatalf("pr_state.status = %q, want ok (gh present)", status)
	}
	if len(entries) != 1 {
		t.Fatalf("pr_state has %d entries, want 1 (the PR-pending entity)", len(entries))
	}
	e := entries[0]
	if e["pr"] != "#42" {
		t.Fatalf("pr_state entry pr = %q, want #42 (the stored field)", e["pr"])
	}
	// The load-bearing pin: the serialized state is the LIVE gh state, not the
	// stored field. The stub reports MERGED for a stored pr of #42, so a MERGED
	// state proves the boot ran `gh pr view` live and serialized its result.
	if e["state"] != "MERGED" {
		t.Fatalf("pr_state entry state = %q, want MERGED (the LIVE gh state) — the boot must run `gh pr view` live, not echo the stored pr field", e["state"])
	}
}

// TestBootPRStateGhAbsentReportsUnknown pins the M6 degraded branch: when `gh` is
// absent from PATH, checkPRStates returns status "gh not available" with NO
// entries. The shallow-boot greet keys off this to state PR merge status is
// UNKNOWN rather than asserting an unknowable state. A PATH stripped of `gh`
// reproduces the branch.
func TestBootPRStateGhAbsentReportsUnknown(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "live-prstate-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	// An isolated empty bin dir as PATH: no `gh`, but also no `git` etc. — boot only
	// needs the workflow files for this section, and the PR-state probe is the part
	// under test. If boot needs git, it is found via the stored fixture, not PATH.
	emptyDir := t.TempDir()
	// Guard the fixture: there must be no `gh` anywhere on this PATH.
	if _, statErr := os.Stat(filepath.Join(emptyDir, "gh")); statErr == nil {
		t.Fatal("empty PATH dir unexpectedly contains a gh shim")
	}
	status, entries := bootPRState(t, root, emptyDir)
	if status != "gh not available" {
		t.Fatalf("pr_state.status = %q, want \"gh not available\" (gh absent)", status)
	}
	if len(entries) != 0 {
		t.Fatalf("pr_state has %d entries with gh absent, want 0 — the greet must report UNKNOWN, not a stale state", len(entries))
	}
}
