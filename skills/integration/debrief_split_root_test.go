// ABOUTME: Deterministic regression guard for the debrief skill's split-root routing —
// ABOUTME: transcribes skills/debrief/SKILL.md's resolve+write+commit over real-git fixtures (AC-1..AC-4).
package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// splitRootDebriefReadme declares a state checkout, so the binary classifies the
// workflow as split-root and resolves entity_dir to the checkout.
const splitRootDebriefReadme = `---
commissioned-by: spacedock@1
id-style: slug
state: .spacedock-state
stages:
  states:
    - name: ideation
      initial: true
    - name: done
      terminal: true
---

# Split Workflow
`

// singleRootDebriefReadme omits state:, so the binary classifies the workflow as
// single-root and resolves entity_dir to the definition dir.
const singleRootDebriefReadme = `---
commissioned-by: spacedock@1
id-style: slug
stages:
  states:
    - name: ideation
      initial: true
    - name: done
      terminal: true
---

# Single Workflow
`

// debriefBoot holds the three fields the debrief skill reads from the boot record
// (Step 1b) to resolve the debrief home. The binary emits them as JSON strings.
type debriefBoot struct {
	StateBackend     string `json:"state_backend"`
	EntityDir        string `json:"entity_dir"`
	EntityDirPresent string `json:"entity_dir_present"`
}

// debriefBootRecord runs the real `spacedock status --boot --json` for defDir and
// returns the resolution fields — the discriminator the skill depends on.
func debriefBootRecord(t *testing.T, defDir string) debriefBoot {
	t.Helper()
	cmd := exec.Command(spacedockBinary(t), "status", "--boot", "--json", "--workflow-dir", defDir)
	cmd.Dir = defDir
	cmd.Env = append(os.Environ(), "HOME="+t.TempDir())
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("spacedock status --boot --json: %v\nstderr=%s", err, stderr.String())
	}
	var rec debriefBoot
	if err := json.Unmarshal(stdout.Bytes(), &rec); err != nil {
		t.Fatalf("parse boot json: %v\n%s", err, stdout.String())
	}
	return rec
}

// debriefResult is the outcome of one transcribed debrief flow.
type debriefResult struct {
	backend string // state_backend the boot record reported
	root    string // resolved {debrief_root} (entity_dir)
	halted  bool   // true if the absent-checkout halt-gate fired
	written string // absolute path of the new debrief file, "" if halted
}

// runDebriefFlow transcribes skills/debrief/SKILL.md's resolve+write+commit
// sequence (Phase 1 Step 1b + Phase 4) as a deterministic regression guard: it
// drives the real boot record to resolve the debrief home, applies the prescribed
// split-root vs single-root routing, and — unless the absent-checkout halt-gate
// fires — writes and commits the debrief. date is fixed for determinism. This
// transcribes the skill's commands, so it is a regression guard / live-drive
// seed, NOT a substitute for the live skill drive.
func runDebriefFlow(t *testing.T, defDir, date string) debriefResult {
	t.Helper()

	// Step 1b — resolve the debrief home from the boot record.
	boot := debriefBootRecord(t, defDir)
	res := debriefResult{backend: boot.StateBackend, root: boot.EntityDir}

	// Halt-gate: a declared-but-absent split-root checkout writes/commits NOTHING.
	if boot.StateBackend == "split-root" && boot.EntityDirPresent != "true" {
		res.halted = true
		return res
	}

	debriefsDir := filepath.Join(res.root, "_debriefs")

	// Phase 4 Step 1 — sequence number from {debrief_root}/_debriefs only.
	seq := nextDebriefSeq(t, debriefsDir, date)
	name := fmt.Sprintf("%s-%02d.md", date, seq)
	rel := filepath.Join("_debriefs", name)
	abs := filepath.Join(debriefsDir, name)

	// Phase 4 Step 3 — write (writeFile creates the _debriefs parent).
	writeFile(t, abs, fmt.Sprintf("---\nsequence: %d\n---\n# Session Debrief — %s #%d\n", seq, date, seq))
	res.written = abs

	// Phase 4 Step 4 — commit.
	msg := fmt.Sprintf("debrief: session %s #%d — guard", date, seq)
	if boot.StateBackend == "split-root" {
		// Path-scoped commit in the state checkout on its own branch, then push.
		git(t, res.root, "add", "--", rel)
		git(t, res.root, "commit", "-m", msg, "--", rel)
		if checkoutHasOrigin(t, res.root) {
			branch := strings.TrimSpace(git(t, res.root, "rev-parse", "--abbrev-ref", "HEAD"))
			git(t, res.root, "push", "origin", branch)
		}
	} else {
		// Single-root: bare commit on the current branch.
		git(t, res.root, "add", "--", rel)
		git(t, res.root, "commit", "-m", msg)
	}
	return res
}

// nextDebriefSeq is one more than the highest {date}-NN.md sequence in debriefsDir,
// or 1 if the dir is absent or empty.
func nextDebriefSeq(t *testing.T, debriefsDir, date string) int {
	t.Helper()
	entries, err := os.ReadDir(debriefsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return 1
		}
		t.Fatal(err)
	}
	max := 0
	prefix := date + "-"
	for _, e := range entries {
		n := e.Name()
		if !strings.HasPrefix(n, prefix) || !strings.HasSuffix(n, ".md") {
			continue
		}
		num := strings.TrimSuffix(strings.TrimPrefix(n, prefix), ".md")
		if v, err := strconv.Atoi(num); err == nil && v > max {
			max = v
		}
	}
	return max + 1
}

// checkoutHasOrigin reports whether the checkout has an origin remote (git remote
// exits 0 with empty output when there is none).
func checkoutHasOrigin(t *testing.T, checkout string) bool {
	t.Helper()
	for _, r := range strings.Fields(git(t, checkout, "remote")) {
		if r == "origin" {
			return true
		}
	}
	return false
}

// stageSplitRootDebrief builds a split-root workflow in a fresh git repo: README
// in the definition dir declaring state: .spacedock-state, the state path
// gitignored on the code branch, and an orphan state branch checked out at that
// path as a linked worktree (the commission mechanics). withOrigin clones from a
// bare origin so the orphan branch is pushed and the flow exercises push; without
// it the repo is local-only (the no-origin carve-out). Returns the definition
// dir, the state checkout, and the code branch name.
func stageSplitRootDebrief(t *testing.T, withOrigin bool) (defDir, stateCheckout, codeBranch string) {
	t.Helper()
	var host string
	if withOrigin {
		bare := filepath.Join(t.TempDir(), "origin.git")
		git(t, t.TempDir(), "init", "-q", "--bare", bare)
		host = filepath.Join(t.TempDir(), "host")
		git(t, t.TempDir(), "clone", "-q", bare, host)
	} else {
		host = t.TempDir()
		git(t, host, "init", "-q")
	}

	defDir = filepath.Join(host, "docs", "dev")
	writeFile(t, filepath.Join(defDir, "README.md"), splitRootDebriefReadme)
	writeFile(t, filepath.Join(host, ".gitignore"), "docs/dev/.spacedock-state/\n.worktrees/\n")
	git(t, host, "add", "docs/dev/README.md", ".gitignore")
	git(t, host, "commit", "-q", "-m", "commission split workflow")
	codeBranch = strings.TrimSpace(git(t, host, "rev-parse", "--abbrev-ref", "HEAD"))

	stateBranch := "spacedock-state/dev"
	stateCheckout = filepath.Join(defDir, ".spacedock-state")

	// Birth the orphan state branch in a temp detached worktree, clearing the
	// inherited tree so the state branch carries no code-branch files.
	tmpWT := filepath.Join(t.TempDir(), "orphan-birth")
	git(t, host, "worktree", "add", "--detach", tmpWT)
	git(t, tmpWT, "checkout", "--orphan", stateBranch)
	git(t, tmpWT, "rm", "-rf", "--cached", ".")
	entries, err := os.ReadDir(tmpWT)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Name() == ".git" {
			continue
		}
		if err := os.RemoveAll(filepath.Join(tmpWT, e.Name())); err != nil {
			t.Fatal(err)
		}
	}
	writeFile(t, filepath.Join(tmpWT, "first-task.md"), "---\nstatus: ideation\n---\n# First Task\n")
	git(t, tmpWT, "add", "-A")
	git(t, tmpWT, "commit", "-q", "-m", "seed state")
	if withOrigin {
		git(t, tmpWT, "push", "origin", stateBranch)
	}
	git(t, host, "worktree", "remove", "--force", tmpWT)

	// Check out the orphan branch at the gitignored state path as a linked worktree.
	git(t, host, "worktree", "add", stateCheckout, stateBranch)
	return defDir, stateCheckout, codeBranch
}

// stageSingleRootDebrief builds a single-root workflow (README with no state:) in
// a fresh git repo. Returns the definition dir and the code branch name.
func stageSingleRootDebrief(t *testing.T) (defDir, codeBranch string) {
	t.Helper()
	host := t.TempDir()
	git(t, host, "init", "-q")
	defDir = filepath.Join(host, "docs", "dev")
	writeFile(t, filepath.Join(defDir, "README.md"), singleRootDebriefReadme)
	git(t, host, "add", "-A")
	git(t, host, "commit", "-q", "-m", "commission single workflow")
	codeBranch = strings.TrimSpace(git(t, host, "rev-parse", "--abbrev-ref", "HEAD"))
	return defDir, codeBranch
}

// stageAbsentCheckoutDebrief builds a split-root workflow whose state checkout is
// declared but NOT materialized (no linked worktree at the path — a fresh clone /
// removed worktree). Boot reports entity_dir_present:false. Returns the
// definition dir and the code branch name.
func stageAbsentCheckoutDebrief(t *testing.T) (defDir, codeBranch string) {
	t.Helper()
	host := t.TempDir()
	git(t, host, "init", "-q")
	defDir = filepath.Join(host, "docs", "dev")
	writeFile(t, filepath.Join(defDir, "README.md"), splitRootDebriefReadme)
	writeFile(t, filepath.Join(host, ".gitignore"), "docs/dev/.spacedock-state/\n")
	git(t, host, "add", "-A")
	git(t, host, "commit", "-q", "-m", "commission split workflow (state not initialized)")
	codeBranch = strings.TrimSpace(git(t, host, "rev-parse", "--abbrev-ref", "HEAD"))
	return defDir, codeBranch
}

// TestDebriefSplitRootRoutesToCheckout locks AC-1: for a split-root workflow the
// debrief flow writes and commits under the state checkout's _debriefs/ on the
// state branch — not the definition dir — and the trunk does NOT move.
func TestDebriefSplitRootRoutesToCheckout(t *testing.T) {
	const date = "2026-06-30"
	defDir, checkout, codeBranch := stageSplitRootDebrief(t, true)
	trunkBefore := strings.TrimSpace(git(t, checkout, "rev-parse", codeBranch))

	res := runDebriefFlow(t, defDir, date)

	if res.halted {
		t.Fatalf("split-root flow halted unexpectedly")
	}
	if res.backend != "split-root" {
		t.Fatalf("backend = %q, want split-root", res.backend)
	}
	// New file under the checkout, not the definition dir.
	wantPath := filepath.Join(checkout, "_debriefs", date+"-01.md")
	if res.written != wantPath {
		t.Fatalf("written to %q, want %q", res.written, wantPath)
	}
	if _, err := os.Stat(wantPath); err != nil {
		t.Fatalf("debrief not in checkout: %v", err)
	}
	if _, err := os.Stat(filepath.Join(defDir, "_debriefs", date+"-01.md")); err == nil {
		t.Fatalf("debrief leaked into the definition dir")
	}
	// Commit landed on the state branch, naming the debrief file.
	committed := git(t, checkout, "show", "--name-only", "--format=", "HEAD")
	if !strings.Contains(committed, "_debriefs/"+date+"-01.md") {
		t.Fatalf("state-branch HEAD does not contain the debrief:\n%s", committed)
	}
	if branch := strings.TrimSpace(git(t, checkout, "rev-parse", "--abbrev-ref", "HEAD")); branch != "spacedock-state/dev" {
		t.Fatalf("checkout HEAD branch = %q, want spacedock-state/dev", branch)
	}
	// Trunk HEAD byte-identical before and after.
	if trunkAfter := strings.TrimSpace(git(t, checkout, "rev-parse", codeBranch)); trunkBefore != trunkAfter {
		t.Fatalf("trunk moved: before=%s after=%s", trunkBefore, trunkAfter)
	}
}

// TestDebriefSingleRootUnchanged locks AC-2 (the discriminator / negative
// control): a workflow with no state: routes the read, write, and commit through
// the definition dir on the trunk — the routing branch is taken only when state:
// resolves to a real checkout.
func TestDebriefSingleRootUnchanged(t *testing.T) {
	const date = "2026-06-30"
	defDir, codeBranch := stageSingleRootDebrief(t)

	res := runDebriefFlow(t, defDir, date)

	if res.halted {
		t.Fatalf("single-root flow halted unexpectedly")
	}
	if res.backend != "single-root" {
		t.Fatalf("backend = %q, want single-root", res.backend)
	}
	if res.root != defDir {
		t.Fatalf("debrief_root = %q, want definition dir %q", res.root, defDir)
	}
	wantPath := filepath.Join(defDir, "_debriefs", date+"-01.md")
	if res.written != wantPath {
		t.Fatalf("written to %q, want %q", res.written, wantPath)
	}
	// Commit landed on the trunk, naming the debrief.
	committed := git(t, defDir, "show", "--name-only", "--format=", "HEAD")
	if !strings.Contains(committed, "_debriefs/"+date+"-01.md") {
		t.Fatalf("trunk HEAD does not contain the debrief:\n%s", committed)
	}
	if branch := strings.TrimSpace(git(t, defDir, "rev-parse", "--abbrev-ref", "HEAD")); branch != codeBranch {
		t.Fatalf("commit on branch %q, want trunk %q", branch, codeBranch)
	}
}

// TestDebriefSplitRootNumberingIgnoresOrphan locks AC-3: for split-root the
// sequence is computed from the state-checkout debriefs only, so a same-named
// orphan in the definition dir neither perturbs the count nor is overwritten.
func TestDebriefSplitRootNumberingIgnoresOrphan(t *testing.T) {
	const date = "2026-06-30"
	defDir, checkout, _ := stageSplitRootDebrief(t, false)

	// Seed the checkout's _debriefs with -01/-02/-03 on the state branch.
	for _, n := range []string{"01", "02", "03"} {
		rel := filepath.Join("_debriefs", date+"-"+n+".md")
		writeFile(t, filepath.Join(checkout, rel), "seed "+n+"\n")
		git(t, checkout, "add", "--", rel)
	}
	git(t, checkout, "commit", "-m", "seed prior debriefs", "--", "_debriefs")

	// Seed a colliding -01 orphan (a DIFFERENT session) in the definition dir.
	orphan := filepath.Join(defDir, "_debriefs", date+"-01.md")
	const orphanBody = "DIFFERENT SESSION — must not be touched\n"
	writeFile(t, orphan, orphanBody)
	git(t, defDir, "add", "--", filepath.Join("_debriefs", date+"-01.md"))
	git(t, defDir, "commit", "-m", "orphan debrief in definition dir")

	res := runDebriefFlow(t, defDir, date)

	// New file is -04 in the checkout, continuing the checkout's history.
	want := filepath.Join(checkout, "_debriefs", date+"-04.md")
	if res.written != want {
		t.Fatalf("written %q, want %q (numbering must continue checkout history, not collide with the orphan)", res.written, want)
	}
	// The orphan is untouched, byte-identical.
	got, err := os.ReadFile(orphan)
	if err != nil {
		t.Fatalf("orphan disappeared: %v", err)
	}
	if string(got) != orphanBody {
		t.Fatalf("orphan was modified: %q", got)
	}
}

// TestDebriefAbsentCheckoutHalts locks AC-4: a split-root workflow whose checkout
// is declared but not materialized HALTS — it writes/commits nothing and falls
// back to neither the definition dir nor an opaque git -C {missing} error.
func TestDebriefAbsentCheckoutHalts(t *testing.T) {
	const date = "2026-06-30"
	defDir, codeBranch := stageAbsentCheckoutDebrief(t)
	trunkBefore := strings.TrimSpace(git(t, defDir, "rev-parse", codeBranch))

	res := runDebriefFlow(t, defDir, date)

	if !res.halted {
		t.Fatalf("absent-checkout flow did not halt (backend=%q, written=%q)", res.backend, res.written)
	}
	if res.backend != "split-root" {
		t.Fatalf("backend = %q, want split-root", res.backend)
	}
	if res.written != "" {
		t.Fatalf("halt must write nothing; wrote %q", res.written)
	}
	// No debrief and no _debriefs dir in the definition dir; no checkout created.
	if _, err := os.Stat(filepath.Join(defDir, "_debriefs")); err == nil {
		t.Fatalf("halt created a _debriefs dir in the definition dir")
	}
	if _, err := os.Stat(filepath.Join(defDir, ".spacedock-state")); err == nil {
		t.Fatalf("absent checkout unexpectedly exists")
	}
	// The trunk did not advance.
	if trunkAfter := strings.TrimSpace(git(t, defDir, "rev-parse", codeBranch)); trunkBefore != trunkAfter {
		t.Fatalf("trunk advanced during halt: %s -> %s", trunkBefore, trunkAfter)
	}
}
