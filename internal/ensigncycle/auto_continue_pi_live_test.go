//go:build live

// ABOUTME: AC-5 Pi live regression — the origin runtime where the bug bit. A real
// ABOUTME: Pi FO, given a completed implementation, must advance + dispatch a validator.
package ensigncycle

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestLivePiAutoContinueAfterImplementation reproduces the auto-continue lifecycle
// invariant on Pi's `pi-subagents` substrate — the ORIGIN runtime where the bug
// concretely bit (the subagent result returns to the parent and the parent must
// itself drive the next lifecycle step; there is no separate event-loop turn to
// resume into). A real Pi first officer is pointed at a split-root dev-shaped
// workflow whose one entity is parked at status: implementation with a filed
// implementation stage report. Under a NEUTRAL runbook (no "advance" / "dispatch
// the validator" coaching — that is the behavior under test), the FO must continue
// the shared `## Completion and Gates` lifecycle in the same turn: advance the
// entity to validation and dispatch a FRESH validator subagent, which leaves a
// durable `## Stage Report: validation`.
//
// The grade is the shared, state-oriented assertAutoContinue over the entity's
// before→after body — the SAME assertion the Claude leg and the offline negative
// use. A run that verifies the implementation report and STOPS (the bug) leaves
// the entity at status: implementation with no validation report and REDS the grade.
func TestLivePiAutoContinueAfterImplementation(t *testing.T) {
	piBin := piBinaryOrSkip(t)
	repo := repoRoot(t)
	piSubagentsRoot := piSubagentsPackageRoot(t)
	binary := piSpacedockBinary(t, repo)

	piHome := t.TempDir()
	sessionDir := t.TempDir()
	cleanHome := t.TempDir()
	seedPiLocalAuth(t, piHome, os.Getenv("HOME"))
	workflowRoot, stateRoot, entityPath := writePiAutoContinueWorkflow(t)
	artifactDir := filepath.Join(piLiveArtifactDir(t, "pi-auto-continue"), "run")
	if err := os.MkdirAll(filepath.Join(artifactDir, "sessions"), 0o755); err != nil {
		t.Fatal(err)
	}
	env := piLiveEnv(piHome, sessionDir, cleanHome, filepath.Dir(binary), piSubagentsRoot)

	before := readFile(t, entityPath)

	prompt := piAutoContinuePrompt(repo, workflowRoot, stateRoot, entityPath)
	runPiLiveCommand(t, artifactDir, workflowRoot, env, piBin,
		"--print",
		"--session-dir", filepath.Join(artifactDir, "sessions"),
		"--extension", filepath.Join(piSubagentsRoot, "src", "extension", "index.ts"),
		"--skill", filepath.Join(piSubagentsRoot, "skills", "pi-subagents"),
		"--skill", filepath.Join(repo, "skills", "first-officer"),
		"--skill", filepath.Join(repo, "skills", "ensign"),
		prompt,
	)

	after := resolvePiAutoContinueEndState(stateRoot, entityPath)
	observed := piLiveObservedOutput(t, artifactDir)
	if err := assertAutoContinue(before, after, observed); err != nil {
		t.Fatalf("Pi auto-continue scenario graded FAIL: %v; artifacts in %s\n--- entity after ---\n%s", err, artifactDir, after)
	}
}

// resolvePiAutoContinueEndState returns the split-root entity's durable end-state
// body. The Pi FO may archive a terminalized entity, moving it out of its original
// path in the state checkout; this reads the archived copy when the original is
// gone, otherwise the original (entity held at the gate, or still at validation).
// A genuinely absent entity yields an empty body and the state-oriented
// assertAutoContinue reds — it never fabricates state.
func resolvePiAutoContinueEndState(stateRoot, entityPath string) string {
	if data, err := os.ReadFile(entityPath); err == nil {
		return string(data)
	}
	for _, p := range []string{
		filepath.Join(stateRoot, "_archive", "auto-continue-task", "index.md"),
		filepath.Join(stateRoot, "_archive", "auto-continue-task.md"),
	} {
		if data, err := os.ReadFile(p); err == nil {
			return string(data)
		}
	}
	return ""
}

// piBinaryOrSkip returns the pi binary path or skips when Pi is not installed —
// the same skip the existing Pi smoke uses, so the offline lane is unaffected.
func piBinaryOrSkip(t *testing.T) string {
	t.Helper()
	p, err := exec.LookPath("pi")
	if err != nil {
		t.Skip("pi not on PATH; install Pi CLI before running the live Pi auto-continue drive")
	}
	return p
}

// piLiveObservedOutput concatenates the Pi run's stdout + stderr artifacts as the
// observed final output for the grade. assertAutoContinue grades primarily on the
// durable entity state; the observed string is carried for parity with the
// Claude/offline legs (which also pass observed) but the state checks are the
// load-bearing ones, so a transcript-only narration cannot satisfy the grade.
func piLiveObservedOutput(t *testing.T, artifactDir string) string {
	t.Helper()
	var b strings.Builder
	for _, name := range []string{"pi-stdout.txt", "pi-stderr.txt"} {
		p := filepath.Join(artifactDir, name)
		if data, err := os.ReadFile(p); err == nil {
			b.Write(data)
			b.WriteString("\n")
		}
	}
	return b.String()
}

// writePiAutoContinueWorkflow stages a split-root dev-shaped workflow
// (backlog → implementation → validation(fresh,gate) → done) with the entity
// parked at an implementation-ready state in the STATE checkout. It mirrors
// writePiSplitRootSmokeWorkflow's split-root layout (separate code root + state
// checkout, both git-initialized) so the FO drives the real split-root path. The
// validation stage is non-worktree (like the proven Pi split-root smoke); the
// behavior under test is the parent advancing + dispatching the fresh validator,
// not worktree mechanics.
func writePiAutoContinueWorkflow(t *testing.T) (workflowRoot, stateRoot, entityPath string) {
	t.Helper()
	workflowRoot = t.TempDir()
	stateRoot = filepath.Join(workflowRoot, ".spacedock-state")
	writeFile(t, filepath.Join(workflowRoot, "README.md"), piAutoContinueReadme())
	entityPath = filepath.Join(stateRoot, "auto-continue-task", "index.md")
	writeFile(t, entityPath, autoContinueEntity())
	gitInit(t, workflowRoot)
	gitInit(t, stateRoot)
	return workflowRoot, stateRoot, entityPath
}

// piAutoContinueReadme is the split-root variant of autoContinueReadme(): same
// dev-shaped stage graph, plus `state: .spacedock-state` so the entity body and
// stage reports live in the state checkout (the split-root contract the bug's
// origin entity used).
func piAutoContinueReadme() string {
	return "---\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"state: .spacedock-state\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: backlog\n" +
		"      initial: true\n" +
		"    - name: implementation\n" +
		"    - name: validation\n" +
		"      fresh: true\n" +
		"      feedback-to: implementation\n" +
		"      gate: true\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Pi Auto-Continue Fixture\n\n" +
		"### backlog\n\nSeed the task.\n\n- **Outputs:** A seed task.\n\n" +
		"### implementation\n\nProduce the deliverable.\n\n- **Outputs:** The deliverable plus an implementation stage report.\n\n" +
		"### validation\n\n" +
		"Verify the implementation against the acceptance criteria. Append a `## Stage Report: validation` " +
		"section to the entity with one `- DONE:` item and a PASSED or REJECTED recommendation.\n\n" +
		"- **Outputs:** A PASSED or REJECTED validation stage report.\n\n" +
		"### done\n\nTerminal state.\n"
}

// piAutoContinuePrompt is the NEUTRAL Pi runbook. It points the FO at the
// split-root workflow and the one entity whose implementation worker has just
// completed, names the substrate facts the Pi FO needs (use pi-subagents
// subagent(...) for dispatch; this is split-root), and asks it to PROCEED per its
// contract — it deliberately does NOT instruct the FO to advance, dispatch a
// validator, or drive to done. Whether the FO continues the lifecycle is exactly
// the behavior under test.
func piAutoContinuePrompt(repo, workflowRoot, stateRoot, entityPath string) string {
	return fmt.Sprintf(`You are the Spacedock first officer running on Pi for a live auto-continue regression.

Use the local Spacedock first-officer skill and the Pi first-officer runtime adapter. Dispatch any workers with the pi-subagents subagent(...) tool. Do not use or mention Claude Agent, SendMessage, TeamCreate, or TeamDelete tools.

This is a split-root Spacedock workflow.
Repo root: %[1]s
Workflow directory: %[2]s
State checkout: %[3]s
Entity file: %[4]s

The entity `+"`auto-continue-task`"+` is parked at status: implementation, and its implementation worker has just completed and filed its `+"`## Stage Report: implementation`"+`. Proceed with the workflow exactly as the first-officer contract directs, then give your final response.

When you dispatch a worker, give it: load and follow the local Spacedock ensign skill at %[1]s/skills/ensign/SKILL.md and the Pi ensign adapter at %[1]s/skills/ensign/references/pi-ensign-runtime.md; the workflow directory, state checkout, and entity file above; and the target stage. The worker must not edit YAML frontmatter, must append its stage report with the exact heading for its stage, and must path-scoped commit only the entity path in the state checkout.`, repo, workflowRoot, stateRoot, entityPath)
}
