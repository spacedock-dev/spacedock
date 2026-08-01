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

	"github.com/spacedock-dev/spacedock/internal/gates"
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
    - name: validation
      gate: true
    - name: implementation
    - name: done
      terminal: true
---

# Identify Ready Gates Workflow
`

func openGateEntity(slug, status, score string) string {
	return "---\nid: " + slug + "\nstatus: " + status + "\nscore: " + score + "\ngates:\n" +
		"  version: 1\n  current: {gate: 'gate:" + slug + ":" + status + "'}\n  records:\n" +
		"    - id: gate:" + slug + ":" + status + "\n      stage: " + status + "\n      attempts:\n" +
		"        - id: gate-attempt:" + slug + "-" + status + "-1\n" +
		"          briefing: {id: 'briefing:" + slug + ":" + status + ":attempt-1', digest: 'sha256:" + strings.Repeat("1", 64) + "', digest-domain: canonical-bytes, room-ref: ./review/" + status + "/briefing-1}\n" +
		"---\n# " + slug + "\n"
}

func approvedGateEntity(slug, status, target, score string) string {
	body := strings.TrimSuffix(openGateEntity(slug, status, score), "---\n# "+slug+"\n")
	return body +
		"          resolution: {type: Resolution, id: 'resolution:" + slug + ":1', briefing: 'briefing:" + slug + ":" + status + ":attempt-1', by: 'person:captain', at: '2026-07-23T00:00:00Z', decision: approve}\n" +
		"          application: {action: advance, target-stage: " + target + ", state: pending, blockers: []}\n" +
		"---\n# " + slug + "\n"
}

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
	ordinary, ordinaryErr, ordinaryCode := runNative(t, def, env, "--workflow-dir", def, "--boot", "--json")
	if ordinaryCode != 0 || strings.Contains(ordinary, `"ready_gates"`) {
		t.Fatalf("ordinary boot gained identify-only ready_gates: exit=%d stderr=%q output=%s", ordinaryCode, ordinaryErr, ordinary)
	}
}

// TestBootIdentifyReadyGates is AC-1/AC-6's 3-of-5 native counterexample: all
// five entities share the validation stage, but only the three with selected
// durable attempts are scheduled. It also pins ordering and dispatch separation.
func TestBootIdentifyReadyGates(t *testing.T) {
	def, state := buildSplitRoot(t, identifyReadyGatesReadme, map[string]string{
		"sp.md":          "---\nid: sp\nstatus: validation\nscore: 100\n---\n# Still validating\n",
		"mf.md":          openGateEntity("mf", "validation", "90"),
		"r4.md":          strings.Replace(openGateEntity("r4", "validation", "80"), "id: r4\n", "", 1),
		"2n.md":          approvedGateEntity("2n", "validation", "done", "70"),
		"qc.md":          "---\nid: qc\nstatus: validation\nscore: 60\n---\n# Still validating\n",
		"dispatch-me.md": "---\nstatus: draft\nscore: 1000\n---\n",
		"done.md":        "---\nstatus: done\n---\n",
		"unknown.md":     "---\nstatus: vanished\n---\n",
	})
	archiveDir := filepath.Join(state, "_archive")
	if err := os.MkdirAll(archiveDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(archiveDir, "archived-gate.md"), openGateEntity("archived-gate", "validation", "9999"))

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

	wantReady := `[{"id":"mf","slug":"mf","current":"validation","readiness":"awaiting-captain"},{"id":"r4","slug":"r4","current":"validation","readiness":"awaiting-captain"},{"id":"2n","slug":"2n","current":"validation","readiness":"approved-awaiting-merge"}]`
	if got := string(identify.ReadyGates); got != wantReady {
		t.Fatalf("ready_gates = %s\nwant        = %s", got, wantReady)
	}
	wantDispatchable := `[{"id":"dispatch-me","slug":"dispatch-me","current":"draft","next":"validation","worktree":"no"}]`
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
	mf := filepath.Join(state, "mf.md")
	mfBytes, err := os.ReadFile(mf)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(mf, append(mfBytes, []byte("\n## Stage Report: validation\n\n- DONE: prose only\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	if afterBodyEdit := identifyReadyRows(t, def); afterBodyEdit != wantReady {
		t.Fatalf("body-only report edit changed ready_gates: %s", afterBodyEdit)
	}
}

func TestBootReadyGatesRequiresCurrentStageSelection(t *testing.T) {
	def, state := buildSplitRoot(t, identifyReadyGatesReadme, nil)
	room := filepath.Join(state, "review", "validation", "briefing-1")
	writeFile(t, filepath.Join(room, "briefing.json"),
		`{"type":"Briefing","version":"1","id":"briefing:mf:validation:attempt-1","question":"ready?","artifacts":[{"id":"artifact:1","uri":"artifact.md","rev":"sha256:`+strings.Repeat("a", 64)+`"}]}`)
	briefing := filepath.Join(room, "briefing.json")
	data, err := os.ReadFile(briefing)
	if err != nil {
		t.Fatal(err)
	}
	digest, err := gates.CanonicalDigest(data)
	if err != nil {
		t.Fatal(err)
	}
	entity := filepath.Join(state, "mf.md")
	writeFile(t, entity, "---\nid: mf\nstatus: validation\ngates:\n"+
		"  version: 1\n  current: {gate: 'gate:mf:ideation'}\n  records:\n"+
		"    - id: gate:mf:ideation\n      stage: ideation\n      attempts:\n"+
		"        - id: gate-attempt:mf-ideation-1\n          briefing: {id: 'briefing:mf:ideation:attempt-1', digest: 'sha256:"+strings.Repeat("2", 64)+"', digest-domain: canonical-bytes, room-ref: ./review/ideation/briefing-1}\n"+
		"    - id: gate:mf:validation\n      stage: validation\n      attempts:\n"+
		"        - id: gate-attempt:mf-validation-1\n          briefing: {id: 'briefing:mf:validation:attempt-1', digest: '"+digest+"', digest-domain: canonical-bytes, room-ref: ./review/validation/briefing-1}\n"+
		"---\n# MF\n")

	before := identifyReadyRows(t, def)
	if before != "[]" {
		t.Fatalf("stale old-stage selection scheduled mf: %s", before)
	}
	body, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	selected := strings.Replace(string(body), "current: {gate: 'gate:mf:ideation'}", "current: {gate: 'gate:mf:validation'}", 1)
	if selected == string(body) {
		t.Fatal("current-stage selection fixture was not updated")
	}
	if err := os.WriteFile(entity, []byte(selected), 0o644); err != nil {
		t.Fatal(err)
	}
	after := identifyReadyRows(t, def)
	want := `[{"id":"mf","slug":"mf","current":"validation","readiness":"awaiting-captain"}]`
	if after != want {
		t.Fatalf("ready rows after same-Briefing selection repair = %s, want %s", after, want)
	}
}

func identifyReadyRows(t *testing.T, def string) string {
	t.Helper()
	out, errOut, code := runNative(t, def, pinnedEnv(t), "--workflow-dir", def, "--boot", "--identify", "--json")
	if code != 0 {
		t.Fatalf("identify exit=%d stderr=%q", code, errOut)
	}
	var rec struct {
		Ready json.RawMessage `json:"ready_gates"`
	}
	if err := json.Unmarshal([]byte(out), &rec); err != nil {
		t.Fatal(err)
	}
	return string(rec.Ready)
}

func TestBootReadyGateTerminalApprovalPersistsAtConsumeUntilDeliveryEnvelope(t *testing.T) {
	readme := strings.Replace(identifyReadyGatesReadme, "    - name: implementation\n", "", 1)
	def, state := buildSplitRoot(t, readme, map[string]string{
		"2n.md": "---\nid: 2n\nstatus: validation\n---\n# 2n\n",
	})
	room := filepath.Join(state, "review", "validation", "briefing-1")
	briefing := filepath.Join(room, "briefing.json")
	writeFile(t, briefing,
		`{"type":"Briefing","version":"1","id":"briefing:2n:validation:attempt-1:revision-1","question":"ship?","artifacts":[{"id":"artifact:1","uri":"artifact.md","rev":"sha256:`+strings.Repeat("a", 64)+`"}]}`)
	entity := filepath.Join(state, "2n.md")
	briefingBytes, err := os.ReadFile(briefing)
	if err != nil {
		t.Fatal(err)
	}
	digest, err := gates.CanonicalDigest(briefingBytes)
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, entity, "---\nid: 2n\nstatus: validation\ngates:\n"+
		"  version: 1\n  current: {gate: 'gate:2n:validation'}\n  records:\n"+
		"    - id: gate:2n:validation\n      stage: validation\n      attempts:\n"+
		"        - id: gate-attempt:2n-validation-1\n"+
		"          briefing: {id: 'briefing:2n:validation:attempt-1:revision-1', digest: '"+digest+"', digest-domain: canonical-bytes, room-ref: ./review/validation/briefing-1/briefing.json}\n"+
		"---\n# 2n\n")
	if _, _, err := gates.Read(entity); err != nil {
		t.Fatal(err)
	}
	if err := gates.RecordSemantic(entity, gates.RecordInput{
		Decision: "approve", Actor: "person:captain", WorkflowDir: def,
	}); err != nil {
		t.Fatal(err)
	}
	want := `[{"id":"2n","slug":"2n","current":"validation","readiness":"approved-awaiting-merge"}]`
	if got := identifyReadyRows(t, def); got != want {
		t.Fatalf("pending terminal approval = %s, want %s", got, want)
	}
	eligibility, err := gates.EligibilityFileAt(entity, def)
	if err != nil || !eligibility.Eligible || eligibility.TargetStage != "done" {
		t.Fatalf("terminal approval eligibility = %#v, err=%v", eligibility, err)
	}
	// Consume routes a terminal-target approval WITHOUT spending it: the
	// application stays pending (merge guard's delivery envelope is the only
	// terminal consumer), so the ready row persists after consume — the routed
	// entity must keep surfacing through the existing approved-awaiting-merge
	// display until delivery proof lands.
	result, err := gates.ConsumeAt(entity, def)
	if err != nil || result.Consumed || !gates.ApprovedAwaitingMergeRoute(entity, def) {
		t.Fatalf("consume = %#v routed=%t, err=%v", result, gates.ApprovedAwaitingMergeRoute(entity, def), err)
	}
	if got := identifyReadyRows(t, def); got != want {
		t.Fatalf("routed terminal approval lost its ready row: %s, want %s", got, want)
	}
	doc, _, err := gates.Read(entity)
	if err != nil || doc.Records[0].Attempts[0].Application.State != "pending" {
		t.Fatalf("routed approval must stay pending: %#v, err=%v", doc, err)
	}
}

func TestBootReadyGatesFailClosedLifecycleControls(t *testing.T) {
	blocked := strings.Replace(approvedGateEntity("blocked", "validation", "done", "90"),
		"blockers: []", "blockers: [{id: blocker:x, state: unsatisfied}]", 1)
	held := strings.Replace(approvedGateEntity("held", "validation", "done", "80"),
		"blockers: []}", "blockers: [], execution-hold: {state: active}}", 1)
	feedback := strings.Replace(approvedGateEntity("feedback", "validation", "done", "70"),
		"decision: approve}", "decision: revise, reason: revise}", 1)
	feedback = strings.Replace(feedback, "action: advance", "action: feedback", 1)
	consumed := strings.Replace(approvedGateEntity("consumed", "validation", "done", "60"),
		"state: pending", "state: consumed", 1)
	superseded := strings.Replace(approvedGateEntity("superseded", "validation", "done", "50"),
		"state: pending", "state: superseded", 1)
	def, _ := buildSplitRoot(t, identifyReadyGatesReadme, map[string]string{
		"validating.md": "---\nstatus: validation\n---\n",
		"blocked.md":    blocked,
		"held.md":       held,
		"feedback.md":   feedback,
		"consumed.md":   consumed,
		"superseded.md": superseded,
		"malformed.md":  "---\nstatus: validation\ngates:\n  version: 1\n  current: {gate: missing}\n  records: []\n---\n",
		"terminal.md":   openGateEntity("terminal", "done", "100"),
		"ordinary.md":   openGateEntity("ordinary", "implementation", "100"),
	})
	if got := identifyReadyRows(t, def); got != "[]" {
		t.Fatalf("fail-closed lifecycle controls scheduled rows: %s", got)
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
