// ABOUTME: AC-1 fixture harness — drives the compiled command surface through
// ABOUTME: the measured 10-command ceremony and the collapsed 2-command ceremony.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gates"
	"github.com/spacedock-dev/spacedock/internal/status"
	"github.com/spacedock-dev/spacedock/internal/testgit"
)

// gateCeremonyReadme declares the measured shape: a gate:true ideation stage
// advancing to a worktree:true implementation stage, split-root state.
const gateCeremonyReadme = `---
commissioned-by: spacedock@1
id-style: slug
state: .spacedock-state
stages:
  defaults: {worktree: false, concurrency: 1}
  states:
    - name: ideation
      initial: true
      gate: true
    - name: implementation
      worktree: true
    - name: done
      terminal: true
---
# AC-1 Ceremony Fixture

### ideation

Ideate.

- **Outputs:** a plan.

### implementation

Implement.

- **Outputs:** code.

### done

Done.
`

// gateCeremonyFixture births a real split-root workflow — a main repo and an
// independent `.spacedock-state` repo on branch spacedock-state/dev
// (StateBranch's default for a workflow dir named "dev"), matching
// gatePrepareCLIFixture's proven shape (gate_test.go) — then runs a REAL `gate
// prepare` (mustSpacedock) to bind a genuine room + Briefing, so gate
// record/consume's digest and reviewed-input eligibility checks succeed exactly
// as they would in production; hand-crafted gates: YAML cannot satisfy those
// checks without reproducing the recorder's own digest math. Prepare's own
// bind commit is fixture setup, not part of the measured ceremony (the entity
// itself scopes the ceremony's start at the open gate, at `gate record`).
func gateCeremonyFixture(t *testing.T) (mainRoot, workflowDir, entityPath, checklistFile string) {
	t.Helper()
	root := t.TempDir()
	mainRoot = filepath.Join(root, "main")
	workflowDir = filepath.Join(mainRoot, "docs", "dev")
	if err := os.MkdirAll(workflowDir, 0o755); err != nil {
		t.Fatal(err)
	}
	testgit.InitRepo(t, mainRoot, "-q")
	writeFile(t, filepath.Join(workflowDir, "README.md"), gateCeremonyReadme)
	artifact := filepath.Join(mainRoot, "gate-review.md")
	writeFile(t, artifact, "# Review\n\nReady to proceed to implementation.\n")
	git(t, mainRoot, "add", ".")
	git(t, mainRoot, "commit", "-q", "-m", "main fixture")

	statePath := filepath.Join(workflowDir, ".spacedock-state")
	if err := os.MkdirAll(statePath, 0o755); err != nil {
		t.Fatal(err)
	}
	testgit.InitRepo(t, statePath, "-q")
	grandfatherFlatRooms(t, statePath, "task")
	writeFile(t, filepath.Join(statePath, "task.md"), "---\nid: task\nstatus: ideation\ntitle: Task\nstarted:\nworktree:\n---\n# Task\n")
	git(t, statePath, "add", ".")
	git(t, statePath, "commit", "-q", "-m", "state fixture")
	git(t, statePath, "branch", "-M", "spacedock-state/dev")

	mustSpacedock(t, mainRoot, "gate", "prepare", "task",
		"--question", "Advance to implementation?",
		"--artifact", artifact,
		"--summary", "Ready to proceed to implementation.",
		"--workflow-dir", workflowDir)
	mustSpacedock(t, mainRoot, "state", "commit", "task", "--workflow-dir", workflowDir)

	entityPath = filepath.Join(statePath, "task.md")
	checklistFile = filepath.Join(t.TempDir(), "impl.checklist")
	writeFile(t, checklistFile, "- a\n- b\n")
	return mainRoot, workflowDir, entityPath, checklistFile
}

// mustSpacedock runs one spacedock command through the compiled command
// surface, asserting exit 0 — the AC-1 harness's (a) predicate, checked at
// every one of the ten (or two) commands, not just the last.
func mustSpacedock(t *testing.T, dir string, args ...string) string {
	t.Helper()
	var out, errOut bytes.Buffer
	code := run(context.Background(), args, os.Environ(), dir, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("spacedock %v exit=%d stdout=%q stderr=%q", args, code, out.String(), errOut.String())
	}
	return out.String()
}

const gateCeremonyWorkerKey = "spacedock-ensign"

// runBeforeCeremonyCommands drives the exact 10 binary/git commands transcribed
// from the `## Problem` transcript baseline (the incidental retitle `--set` and
// the five host-side calls are excluded — see the entity's AC-1).
func runBeforeCeremonyCommands(t *testing.T, hostClone, workflowDir, entityPath, checklistFile string) string {
	t.Helper()
	slug := status.EntitySlug(entityPath)
	worktreeRel := ".worktrees/" + gateCeremonyWorkerKey + "-" + slug

	mustSpacedock(t, hostClone, "gate", "record", slug, "--decision", "approve", "--actor", "person:captain", "--workflow-dir", workflowDir)         // 1. gate record
	mustSpacedock(t, hostClone, "state", "commit", slug, "--workflow-dir", workflowDir)                                                              // 2. state commit
	mustSpacedock(t, hostClone, "gate", "consume", slug, "--workflow-dir", workflowDir)                                                              // 3. gate consume
	mustSpacedock(t, hostClone, "state", "commit", slug, "--workflow-dir", workflowDir)                                                              // 4. state commit
	mustSpacedock(t, hostClone, "status", "--workflow-dir", workflowDir, "--read", entityPath, "--json")                                             // 5. status --read (verify)
	mustSpacedock(t, hostClone, "status", "--workflow-dir", workflowDir, "--set", slug, "status=implementation", "worktree="+worktreeRel, "started") // 6. status --set (worktree stamp)
	mustSpacedock(t, hostClone, "state", "commit", slug, "--workflow-dir", workflowDir)                                                              // 7. state commit
	git(t, filepath.Join(workflowDir, ".spacedock-state"), "status", "--porcelain")                                                                  // 8. git status
	git(t, hostClone, "worktree", "add", "-b", gateCeremonyWorkerKey+"/"+slug, filepath.Join(hostClone, worktreeRel))                                // 9. git worktree add
	return mustSpacedock(t, hostClone, "dispatch", "build", "--host", "claude", "--team-name", "fixture-team",                                       // 10. dispatch build
		"--workflow-dir", workflowDir, "--entity-path", entityPath, "--stage", "implementation", "--checklist-file", checklistFile)
}

// runAfterCeremonyCommands drives the collapsed 2-command ceremony pinned by
// AC-1: `gate record --consume` then `dispatch build --stamp`, nothing else.
func runAfterCeremonyCommands(t *testing.T, hostClone, workflowDir, entityPath, checklistFile string) string {
	t.Helper()
	slug := status.EntitySlug(entityPath)

	mustSpacedock(t, hostClone, "gate", "record", slug, "--decision", "approve", "--actor", "person:captain", "--consume", "--workflow-dir", workflowDir) // 1. gate record --consume
	return mustSpacedock(t, hostClone, "dispatch", "build", "--stamp", "--host", "claude", "--team-name", "fixture-team",                                 // 2. dispatch build --stamp
		"--workflow-dir", workflowDir, "--entity-path", entityPath, "--stage", "implementation", "--checklist-file", checklistFile)
}

// assertGateCeremonyEndState is the AC-1 harness's (b) predicate: the
// falsifiable end-state parity both the before-list and after-list must reach —
// Resolution closed and application consumed, status advanced to the successor,
// started/worktree stamped, state-checkout log/tree clean, worktree present on
// its branch, spawn envelope emitted. If a design regression drops a required
// step from the two-command after-list, one of these predicates fails.
func assertGateCeremonyEndState(t *testing.T, label, hostClone, workflowDir, entityPath, envelope string) {
	t.Helper()
	slug := status.EntitySlug(entityPath)
	statePath := filepath.Join(workflowDir, ".spacedock-state")

	fields := status.ParseFrontmatter(entityPath)
	if fields["status"] != "implementation" {
		t.Errorf("%s: status = %q, want implementation", label, fields["status"])
	}
	if fields["started"] == "" {
		t.Errorf("%s: started was not stamped", label)
	}
	wantWorktree := ".worktrees/" + gateCeremonyWorkerKey + "-" + slug
	if fields["worktree"] != wantWorktree {
		t.Errorf("%s: worktree = %q, want %q", label, fields["worktree"], wantWorktree)
	}

	doc, _, err := gates.Read(entityPath)
	if err != nil {
		t.Fatalf("%s: read gate document: %v", label, err)
	}
	if len(doc.Records) == 0 {
		t.Fatalf("%s: no gate records", label)
	}
	attempts := doc.Records[0].Attempts
	attempt := attempts[len(attempts)-1]
	if attempt.Resolution == nil || attempt.Resolution.Decision != "approve" {
		t.Errorf("%s: Resolution not closed approve: %+v", label, attempt.Resolution)
	}
	if attempt.Application == nil || attempt.Application.State != "consumed" {
		t.Errorf("%s: application state = %+v, want consumed", label, attempt.Application)
	}

	if porcelain := git(t, statePath, "status", "--porcelain"); strings.TrimSpace(porcelain) != "" {
		t.Errorf("%s: state checkout dirty after the ceremony: %s", label, porcelain)
	}
	log := git(t, statePath, "log", "--oneline")
	if strings.Count(log, "\n") < 2 {
		t.Errorf("%s: state checkout log missing the ceremony commits:\n%s", label, log)
	}

	worktreePath := filepath.Join(hostClone, wantWorktree)
	if info, err := os.Stat(worktreePath); err != nil || !info.IsDir() {
		t.Fatalf("%s: worktree missing at %s: %v", label, worktreePath, err)
	}
	if branch := strings.TrimSpace(git(t, worktreePath, "rev-parse", "--abbrev-ref", "HEAD")); branch != gateCeremonyWorkerKey+"/"+slug {
		t.Errorf("%s: worktree branch = %q, want %s", label, branch, gateCeremonyWorkerKey+"/"+slug)
	}

	if !strings.Contains(envelope, `"schema_version"`) {
		t.Errorf("%s: dispatch build emitted no spawn envelope:\n%s", label, envelope)
	}
}

// TestGateCeremonyCollapseAC1 is the AC-1 fixture harness: it drives the
// compiled command surface twice over the same fixture shape (a gated stage ->
// worktree-declaring successor) — once through the 10-command before-list, once
// through the after-list pinned to the two documented FO commands — and asserts
// both runs converge to equivalent end-state predicates. The 10 -> 2 delta is a
// property of the two committed command lists this test executes verbatim, not
// a separate self-referential length assertion.
func TestGateCeremonyCollapseAC1(t *testing.T) {
	t.Run("before", func(t *testing.T) {
		hostClone, workflowDir, entityPath, checklistFile := gateCeremonyFixture(t)
		envelope := runBeforeCeremonyCommands(t, hostClone, workflowDir, entityPath, checklistFile)
		assertGateCeremonyEndState(t, "before", hostClone, workflowDir, entityPath, envelope)
	})
	t.Run("after", func(t *testing.T) {
		hostClone, workflowDir, entityPath, checklistFile := gateCeremonyFixture(t)
		envelope := runAfterCeremonyCommands(t, hostClone, workflowDir, entityPath, checklistFile)
		assertGateCeremonyEndState(t, "after", hostClone, workflowDir, entityPath, envelope)
	})
}

// TestGateCeremonyDispatchCommitContainsEnteredStageEvidence pins AC-2 to the
// exact dispatch commit. A later commit that adds started makes this test fail.
func TestGateCeremonyDispatchCommitContainsEnteredStageEvidence(t *testing.T) {
	hostClone, workflowDir, entityPath, checklistFile := gateCeremonyFixture(t)
	runAfterCeremonyCommands(t, hostClone, workflowDir, entityPath, checklistFile)

	statePath := filepath.Join(workflowDir, ".spacedock-state")
	wantSubject := "dispatch: task entering implementation"
	commit := strings.TrimSpace(git(t, statePath, "log", "-1", "--format=%H", "--grep=^"+wantSubject+"$"))
	if commit == "" {
		t.Fatalf("no exact dispatch commit with subject %q", wantSubject)
	}
	blob := git(t, statePath, "show", commit+":task.md")
	fields := status.ParseFrontmatterData([]byte(blob))
	if fields["status"] != "implementation" || strings.TrimSpace(fields["started"]) == "" {
		t.Fatalf("exact dispatch commit lacks complete entered-stage evidence: status=%q started=%q", fields["status"], fields["started"])
	}
}

func TestGateJourneyUsesStatusAndActingCommandsWithoutEligibilityPreflight(t *testing.T) {
	hostClone, workflowDir, entityPath, checklistFile := gateCeremonyFixture(t)
	commands := []string{}
	run := func(args ...string) string {
		commands = append(commands, strings.Join(args, " "))
		return mustSpacedock(t, hostClone, args...)
	}

	statusOut := run("status", "--boot", "--identify", "--json", "--workflow-dir", workflowDir)
	var boot struct {
		ReadyGates []map[string]string `json:"ready_gates"`
	}
	if err := json.Unmarshal([]byte(statusOut), &boot); err != nil || len(boot.ReadyGates) != 1 || boot.ReadyGates[0]["readiness"] != "awaiting-captain" {
		t.Fatalf("status did not project the gate's next action: %s", statusOut)
	}
	run("gate", "record", "task", "--decision", "approve", "--actor", "person:captain", "--consume", "--workflow-dir", workflowDir)
	envelope := run("dispatch", "build", "--stamp", "--host", "claude", "--team-name", "fixture-team",
		"--workflow-dir", workflowDir, "--entity-path", entityPath, "--stage", "implementation", "--checklist-file", checklistFile)
	assertGateCeremonyEndState(t, "status-to-acting-command", hostClone, workflowDir, entityPath, envelope)

	for _, command := range commands {
		if strings.Contains(command, "eligibility") {
			t.Fatalf("operator journey invoked removed eligibility preflight: %v", commands)
		}
	}
}
