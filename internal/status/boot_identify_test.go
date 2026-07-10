// ABOUTME: --boot --identify (opt-in local identify) — folds discovery + taxonomy +
// ABOUTME: a local pr: mirror into the record, provably side-effect-free, uniform zero/one/many.
package status

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// identifyPRReadme declares a split-root slug workflow whose implementation stage
// is non-terminal, so a pr-bearing entity there surfaces in the PR mirror.
const identifyPRReadme = `---
commissioned-by: spacedock@1
id-style: slug
state: .spacedock-state
stages:
  states:
    - name: ideation
      initial: true
    - name: implementation
    - name: review
      terminal: true
---

# Identify Boot Workflow
`

const identifyReadyGatesReadme = `---
commissioned-by: spacedock@1
id-style: slug
state: .spacedock-state
stages:
  states:
    - name: draft
      initial: true
    - name: review
      gate: true
    - name: review-b
      gate: true
    - name: terminal-gate
      gate: true
      terminal: true
    - name: done
      terminal: true
---

# Identify Ready Gates Workflow
`

// writeRecordingGh writes a `gh` shim that appends to sentinelPath whenever it is
// invoked, so a test can prove `gh` was NEVER run by asserting the sentinel is
// absent afterward. It still prints a merge state, so a boot that DID shell out to
// gh would both leave the sentinel AND resolve a non-local pr_state.
func writeRecordingGh(t *testing.T, sentinelPath string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("gh recording shim is a POSIX shell script")
	}
	dir := t.TempDir()
	script := "#!/bin/sh\n" +
		"echo called >> " + sentinelPath + "\n" +
		"echo MERGED\n"
	if err := os.WriteFile(filepath.Join(dir, "gh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

// TestBootIdentifyFoldsDiscoveryTaxonomyLocalPR (AC-2) asserts identify mode on a
// healthy split-root workflow exits 0 and emits one record carrying every existing
// boot section, the folded discovery result and stages taxonomy (appended AFTER the
// existing key set), and a LOCAL pr: mirror (pr-pending entity by number, status
// "local"), with no gh call.
func TestBootIdentifyFoldsDiscoveryTaxonomyLocalPR(t *testing.T) {
	def, _ := buildSplitRoot(t, identifyPRReadme, map[string]string{
		"add-login.md": "---\nstatus: implementation\npr: \"#42\"\n---\n",
	})
	env := pinnedEnv(t)

	out, errOut, code := runNative(t, def, env, "--workflow-dir", def, "--boot", "--identify", "--json")
	if code != 0 {
		t.Fatalf("--boot --identify --json exit=%d stderr=%q", code, errOut)
	}

	// The existing key set, then the folded discovery + stages, all in order.
	orderedKeys := []string{
		"command", "mods", "id_style", "next_id",
		"orphans", "pr_state", "dispatchable", "team_state",
		"state_backend", "definition_dir", "entity_dir", "entity_dir_present",
		"sandbox", "discovery", "stages", "ready_gates",
	}
	last := -1
	for _, key := range orderedKeys {
		idx := strings.Index(out, `"`+key+`"`)
		if idx < 0 {
			t.Fatalf("identify record missing key %q\n%s", key, out)
		}
		if idx < last {
			t.Fatalf("identify record key %q out of order (discovery/stages must append after the existing set)\n%s", key, out)
		}
		last = idx
	}

	var rec struct {
		Discovery []string `json:"discovery"`
		Stages    []struct {
			Name string `json:"name"`
		} `json:"stages"`
		PRState struct {
			Status  string              `json:"status"`
			Entries []map[string]string `json:"entries"`
		} `json:"pr_state"`
		ReadyGates json.RawMessage `json:"ready_gates"`
	}
	if err := json.Unmarshal([]byte(out), &rec); err != nil {
		t.Fatalf("parse identify record: %v\n%s", err, out)
	}
	if len(rec.Discovery) != 1 {
		t.Fatalf("discovery = %v, want the one discovered workflow", rec.Discovery)
	}
	if len(rec.Stages) != 3 || rec.Stages[0].Name != "ideation" {
		t.Fatalf("stages taxonomy not folded in: %+v", rec.Stages)
	}
	if rec.PRState.Status != "local" {
		t.Fatalf("pr_state.status = %q, want \"local\" (identify renders the local pr: mirror, no gh)", rec.PRState.Status)
	}
	if len(rec.PRState.Entries) != 1 || rec.PRState.Entries[0]["pr"] != "#42" {
		t.Fatalf("pr_state entries = %+v, want the pr-pending entity by its stored #42", rec.PRState.Entries)
	}
	if rec.PRState.Entries[0]["state"] != "local" {
		t.Fatalf("pr_state entry state = %q, want \"local\" (not-gh-checked)", rec.PRState.Entries[0]["state"])
	}
	if got := string(rec.ReadyGates); got != "[]" {
		t.Fatalf("ready_gates = %s, want [] for a workflow with no current gates", got)
	}
}

// TestBootIdentifyReadyGates is AC-1/AC-2's native-path proof: identify exposes
// every active non-terminal current gate in status order, leaves dispatchable
// byte-stable with --next, and excludes every non-ready class.
func TestBootIdentifyReadyGates(t *testing.T) {
	def, state := buildSplitRoot(t, identifyReadyGatesReadme, map[string]string{
		"alpha-gate.md":    "---\nstatus: review\nscore: 5\n---\n",
		"dispatch-me.md":   "---\nstatus: draft\nscore: 100\n---\n",
		"done.md":          "---\nstatus: done\n---\n",
		"gate-check.md":    "---\nstatus: review\nscore: 5\n---\n",
		"later-gate.md":    "---\nstatus: review-b\nscore: 1000\n---\n",
		"ordinary.md":      "---\nstatus: draft\nworktree: busy\n---\n",
		"terminal-gate.md": "---\nstatus: terminal-gate\n---\n",
		"unknown.md":       "---\nstatus: vanished\n---\n",
	})
	archiveDir := filepath.Join(state, "_archive")
	if err := os.MkdirAll(archiveDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(archiveDir, "archived-gate.md"), "---\nstatus: review\nscore: 9999\n---\n")

	identifyOut, errOut, code := runNative(t, def, pinnedEnv(t), "--workflow-dir", def, "--boot", "--identify", "--json")
	if code != 0 {
		t.Fatalf("--boot --identify --json exit=%d stderr=%q", code, errOut)
	}
	var identify struct {
		Dispatchable json.RawMessage `json:"dispatchable"`
		ReadyGates   json.RawMessage `json:"ready_gates"`
	}
	if err := json.Unmarshal([]byte(identifyOut), &identify); err != nil {
		t.Fatalf("parse identify boot: %v\n%s", err, identifyOut)
	}

	wantReady := `[{"id":"alpha-gate","slug":"alpha-gate","current":"review"},{"id":"gate-check","slug":"gate-check","current":"review"},{"id":"later-gate","slug":"later-gate","current":"review-b"}]`
	if got := string(identify.ReadyGates); got != wantReady {
		t.Fatalf("ready_gates = %s\nwant        = %s", got, wantReady)
	}
	wantDispatchable := `[{"id":"dispatch-me","slug":"dispatch-me","current":"draft","next":"review","worktree":"no"}]`
	if got := string(identify.Dispatchable); got != wantDispatchable {
		t.Fatalf("identify dispatchable = %s\nwant                  = %s", got, wantDispatchable)
	}

	nextOut, nextErr, nextCode := runNative(t, def, pinnedEnv(t), "--workflow-dir", def, "--next", "--json")
	if nextCode != 0 {
		t.Fatalf("--next --json exit=%d stderr=%q", nextCode, nextErr)
	}
	var next struct {
		Dispatchable json.RawMessage `json:"dispatchable"`
	}
	if err := json.Unmarshal([]byte(nextOut), &next); err != nil {
		t.Fatalf("parse next: %v\n%s", err, nextOut)
	}
	if string(next.Dispatchable) != string(identify.Dispatchable) {
		t.Fatalf("identify dispatchable = %s, --next dispatchable = %s", identify.Dispatchable, next.Dispatchable)
	}
}

// TestBootIdentifyIsSideEffectFree (AC-3) is the core boundary guarantee: identify
// mode makes NO gh call (a recording shim on PATH is never invoked, pr_state stays
// local) and NO mutation (the git-backed state checkout's HEAD and working tree are
// byte-identical before and after).
func TestBootIdentifyIsSideEffectFree(t *testing.T) {
	def, state := buildSplitRoot(t, identifyPRReadme, map[string]string{
		"add-login.md": "---\nstatus: implementation\npr: \"#42\"\n---\n",
	})
	// The state checkout is a real git repo so HEAD/tree can be diffed.
	gitC(t, state, "init", "-q")
	gitC(t, state, "config", "user.email", "t@t")
	gitC(t, state, "config", "user.name", "t")
	gitC(t, state, "add", "-A")
	gitC(t, state, "commit", "-q", "-m", "seed")

	headBefore := gitOut(t, state, "rev-parse", "HEAD")
	treeBefore := gitOut(t, state, "status", "--porcelain")

	sentinel := filepath.Join(t.TempDir(), "gh-was-called")
	shimDir := writeRecordingGh(t, sentinel)
	env := pinnedEnv(t)
	// Prepend the shim dir so a boot that shells out to gh WOULD resolve + record it.
	for i, kv := range env {
		if strings.HasPrefix(kv, "PATH=") {
			env[i] = "PATH=" + shimDir + string(os.PathListSeparator) + strings.TrimPrefix(kv, "PATH=")
		}
	}

	out, errOut, code := runNative(t, def, env, "--workflow-dir", def, "--boot", "--identify", "--json")
	if code != 0 {
		t.Fatalf("--boot --identify --json exit=%d stderr=%q", code, errOut)
	}

	if _, err := os.Stat(sentinel); err == nil {
		t.Fatal("identify boot invoked `gh` — the recording shim fired; the greet must make no network call")
	}
	if !strings.Contains(out, `"status":"local"`) {
		t.Fatalf("pr_state is not the local mirror — a gh path was taken\n%s", out)
	}
	if got := gitOut(t, state, "rev-parse", "HEAD"); got != headBefore {
		t.Fatalf("state checkout HEAD moved: %q -> %q — identify boot mutated the state repo", headBefore, got)
	}
	if got := gitOut(t, state, "status", "--porcelain"); got != treeBefore {
		t.Fatalf("state checkout working tree changed: %q -> %q — identify boot wrote to the state checkout", treeBefore, got)
	}
}

// TestBootIdentifyUniformZeroOneMany (AC-5) asserts the uniform discovery contract
// with no N==1 eager convergence: an empty root reports no workflow and does NOT
// broad-search; a one-workflow root lists 1; a two-workflow root lists 2 — and in
// no case does any convergence run (identify never calls state ready/sweep, so the
// state checkouts are untouched — the same side-effect-free guarantee as AC-3).
func TestBootIdentifyUniformZeroOneMany(t *testing.T) {
	// Zero: an empty git-less root discovers nothing → report-and-stop, no sweep.
	empty := t.TempDir()
	_, errOut, code := runNative(t, empty, pinnedEnv(t), "--boot", "--identify", "--json")
	if code == 0 {
		t.Fatalf("empty-root identify boot exited 0, want a no-workflow halt; stderr=%q", errOut)
	}
	if !strings.Contains(errOut, "do NOT search the filesystem") {
		t.Fatalf("zero-discovery halt missing the no-broad-search directive: %q", errOut)
	}

	// One: a root holding exactly one workflow → discovery list of 1.
	oneRoot := t.TempDir()
	buildWorkflowUnder(t, oneRoot, "wf-a")
	if got := identifyDiscovery(t, oneRoot); len(got) != 1 {
		t.Fatalf("one-workflow root: discovery = %v, want length 1", got)
	}

	// Many: a root holding two workflows → discovery list of 2, no convergence.
	twoRoot := t.TempDir()
	buildWorkflowUnder(t, twoRoot, "wf-a")
	buildWorkflowUnder(t, twoRoot, "wf-b")
	if got := identifyDiscovery(t, twoRoot); len(got) != 2 {
		t.Fatalf("two-workflow root: discovery = %v, want length 2", got)
	}
}

// buildWorkflowUnder materializes a commissioned split-root workflow named `name`
// under root, so discovery finds it.
func buildWorkflowUnder(t *testing.T, root, name string) {
	t.Helper()
	def := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Join(def, ".spacedock-state"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(def, "README.md"), identifyPRReadme)
}

// identifyDiscovery runs `--boot --identify --json` at root (no --workflow-dir) and
// returns the discovery list, failing on a non-zero exit.
func identifyDiscovery(t *testing.T, root string) []string {
	t.Helper()
	out, errOut, code := runNative(t, root, pinnedEnv(t), "--boot", "--identify", "--json")
	if code != 0 {
		t.Fatalf("identify boot at %s exit=%d stderr=%q", root, code, errOut)
	}
	var rec struct {
		Discovery []string `json:"discovery"`
	}
	if err := json.Unmarshal([]byte(out), &rec); err != nil {
		t.Fatalf("parse identify discovery: %v\n%s", err, out)
	}
	return rec.Discovery
}

// gitOut runs a read-only git subcommand in dir and returns trimmed stdout.
func gitOut(t *testing.T, dir string, args ...string) string {
	t.Helper()
	full := append([]string{"-C", dir}, args...)
	out, err := exec.Command("git", full...).Output()
	if err != nil {
		t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out))
}
