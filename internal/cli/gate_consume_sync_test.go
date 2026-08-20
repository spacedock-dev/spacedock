// ABOUTME: AC-2 negative/recovery tests — the `--consume` usage-error and
// ABOUTME: room-skip boundaries, and the per-phase sync-failed/HALT/recovery procedure.
package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gates"
	"github.com/spacedock-dev/spacedock/internal/status"
	"github.com/spacedock-dev/spacedock/internal/testgit"
)

// TestGateConsumeFlagRejectsNonApproveDecisionBeforeWrite pins mechanism 2's
// chat-source usage error: `--consume` with `--decision revise|hold` is a
// usage error (exit 2) before any write — the flag never softens a non-approve
// decision.
func TestGateConsumeFlagRejectsNonApproveDecisionBeforeWrite(t *testing.T) {
	for _, decision := range []string{"revise", "hold"} {
		t.Run(decision, func(t *testing.T) {
			root, entity := semanticDecisionFixture(t)
			before, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			var out, errOut bytes.Buffer
			code := run(context.Background(), []string{"gate", "record", "task", "--workflow-dir", root,
				"--decision", decision, "--actor", "person:captain", "--reason", "evidence", "--consume"},
				nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil)
			if code != 2 {
				t.Fatalf("exit=%d stdout=%q stderr=%q, want exit 2", code, out.String(), errOut.String())
			}
			if out.Len() != 0 {
				t.Fatalf("usage-error rejection wrote to stdout: %q", out.String())
			}
			after, err := os.ReadFile(entity)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(after, before) {
				t.Fatal("--consume usage-error rejection changed the entity")
			}
			if _, err := os.Stat(entity + ".gates.lock"); !os.IsNotExist(err) {
				t.Fatalf("usage-error rejection left lock residue: %v", err)
			}
		})
	}
}

// TestGateConsumeFlagSkipsConsumeOnRoomSourceReviseHoldClose pins mechanism 2's
// room-source path: the decision lives inside the room, unresolved until the
// close succeeds, so a revise/hold room close reports the close and skips
// consume (`consume=skipped`, exit 0) rather than erroring.
func TestGateConsumeOnTerminalRouteEmitsNoSyncLine(t *testing.T) {
	hostClone, workflowDir := gatedTerminalSplitRootFixture(t)
	mustSpacedock(t, hostClone, "gate", "record", "task", "--decision", "approve", "--actor", "person:captain", "--consume", "--workflow-dir", workflowDir)

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "consume", "task", "--workflow-dir", workflowDir},
		nil, hostClone, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("terminal consume exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	if !strings.Contains(out.String(), "route=approved-awaiting-merge") {
		t.Fatalf("terminal consume missing route=approved-awaiting-merge: %s", out.String())
	}
	if strings.Contains(out.String(), "sync=") {
		t.Fatalf("terminal (unspent) consume must emit no sync line: %s", out.String())
	}
}

// gatedTerminalSplitRootFixture births a split-root workflow whose gated stage
// advances directly to a terminal stage (validation -> done), for the
// terminal-route no-sync-line assertion.
func gatedTerminalSplitRootFixture(t *testing.T) (hostClone, workflowDir string) {
	t.Helper()
	root := t.TempDir()
	hostClone = filepath.Join(root, "main")
	workflowDir = filepath.Join(hostClone, "docs", "dev")
	if err := os.MkdirAll(workflowDir, 0o755); err != nil {
		t.Fatal(err)
	}
	testgit.InitRepo(t, hostClone, "-q")
	writeFile(t, filepath.Join(workflowDir, "README.md"), "---\nid-style: slug\nstate: .spacedock-state\nstages:\n  states:\n    - name: validation\n      initial: true\n      gate: true\n    - name: done\n      terminal: true\n---\n# Terminal Fixture\n")
	git(t, hostClone, "add", ".")
	git(t, hostClone, "commit", "-q", "-m", "main fixture")

	statePath := filepath.Join(workflowDir, ".spacedock-state")
	if err := os.MkdirAll(statePath, 0o755); err != nil {
		t.Fatal(err)
	}
	testgit.InitRepo(t, statePath, "-q")
	grandfatherFlatRooms(t, statePath, "task")
	writeFile(t, filepath.Join(statePath, "task.md"), "---\nid: task\nstatus: validation\ntitle: Task\n---\n# Task\n")
	git(t, statePath, "add", ".")
	git(t, statePath, "commit", "-q", "-m", "state fixture")
	git(t, statePath, "branch", "-M", "spacedock-state/dev")

	artifact := filepath.Join(hostClone, "gate-review.md")
	writeFile(t, artifact, "# Review\n\nReady.\n")
	git(t, hostClone, "add", ".")
	git(t, hostClone, "commit", "-q", "-m", "artifact")
	mustSpacedock(t, hostClone, "gate", "prepare", "task", "--question", "Ship?", "--artifact", artifact, "--summary", "Ready.", "--workflow-dir", workflowDir)
	mustSpacedock(t, hostClone, "state", "commit", "task", "--workflow-dir", workflowDir)
	return hostClone, workflowDir
}

// blockPush installs a pre-push hook that always fails in checkout, so a
// publish attempt fails without touching origin — the broken-origin shape
// spike round 3 used to verify the per-phase recovery procedure. Returns a
// restore func that removes the hook.
func blockPush(t *testing.T, checkout string) (restore func()) {
	t.Helper()
	hooks := t.TempDir()
	if err := os.WriteFile(filepath.Join(hooks, "pre-push"), []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, checkout, "config", "core.hooksPath", hooks)
	return func() { git(t, checkout, "config", "--unset", "core.hooksPath") }
}

// TestGateRecordSyncFailedRecoversWithStateCommit pins the verified phase=record
// recovery procedure (spike round 3): a close is durable+locally committed when
// its publish fails; the composite/standalone record is not re-entrant (frozen
// closed); `state commit <slug>` publishes the existing close; a standalone
// consume then proceeds normally.
func TestGateRecordSyncFailedRecoversWithStateCommit(t *testing.T) {
	hostClone := singleHostOrigin(t)
	workflowDir, _ := commissionGatedSplitWorkflow(t, hostClone)
	checkout := filepath.Join(workflowDir, ".spacedock-state")
	restore := blockPush(t, checkout)

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "record", "task", "--decision", "approve", "--actor", "person:captain", "--workflow-dir", workflowDir},
		nil, hostClone, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 1 {
		t.Fatalf("phase=record sync-failed exit=%d stdout=%q stderr=%q, want 1", code, out.String(), errOut.String())
	}
	if !strings.Contains(out.String(), "sync=failed phase=record") {
		t.Fatalf("stdout missing sync=failed phase=record: %s", out.String())
	}
	headBefore := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD"))

	// Not re-entrant: repeat record refuses the frozen-closed attempt, byte-clean.
	out.Reset()
	errOut.Reset()
	code = run(context.Background(), []string{"gate", "record", "task", "--decision", "approve", "--actor", "person:captain", "--workflow-dir", workflowDir},
		nil, hostClone, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 1 || !strings.Contains(errOut.String(), "frozen closed") {
		t.Fatalf("repeat record exit=%d stdout=%q stderr=%q, want a frozen-closed refusal", code, out.String(), errOut.String())
	}
	if strings.Contains(out.String(), "sync=") {
		t.Fatalf("repeat record refusal must run no sync: %s", out.String())
	}
	if head := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD")); head != headBefore {
		t.Fatalf("repeat record refusal created a new commit: before=%s after=%s", headBefore, head)
	}

	// Recovery: unblock, then `state commit <slug>` publishes the existing close.
	restore()
	out.Reset()
	errOut.Reset()
	code = run(context.Background(), []string{"state", "commit", "task", "--workflow-dir", workflowDir},
		nil, hostClone, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 || !strings.Contains(out.String(), "Published previously committed state") {
		t.Fatalf("recovery state commit exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	if head := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD")); head != headBefore {
		t.Fatalf("recovery state commit created a NEW commit rather than publishing the existing one: before=%s after=%s", headBefore, head)
	}

	// Resume from the durable position: standalone consume proceeds normally.
	out.Reset()
	errOut.Reset()
	code = run(context.Background(), []string{"gate", "consume", "task", "--workflow-dir", workflowDir},
		nil, hostClone, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 || !strings.Contains(out.String(), "consumed=true") || !strings.Contains(out.String(), "sync=pushed phase=consume") {
		t.Fatalf("recovered consume exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
}

// TestGateConsumeSyncFailedRecoversWithStateCommit pins the verified
// phase=consume recovery procedure (spike round 3): an advance is durable when
// its publish fails; repeat consume refuses byte-clean (already consumed);
// `state commit <slug>` publishes the advance.
func TestGateConsumeSyncFailedRecoversWithStateCommit(t *testing.T) {
	hostClone := singleHostOrigin(t)
	workflowDir, _ := commissionGatedSplitWorkflow(t, hostClone)
	checkout := filepath.Join(workflowDir, ".spacedock-state")

	mustSpacedock(t, hostClone, "gate", "record", "task", "--decision", "approve", "--actor", "person:captain", "--workflow-dir", workflowDir)

	restore := blockPush(t, checkout)
	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "consume", "task", "--workflow-dir", workflowDir},
		nil, hostClone, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 1 {
		t.Fatalf("phase=consume sync-failed exit=%d stdout=%q stderr=%q, want 1", code, out.String(), errOut.String())
	}
	if !strings.Contains(out.String(), "sync=failed phase=consume") {
		t.Fatalf("stdout missing sync=failed phase=consume: %s", out.String())
	}
	headBefore := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD"))

	// Not re-entrant: repeat consume is a byte-clean refusal, no new commit.
	// Since PR #599 (Simplify v1 gate state schema, 9ab72ea9f) dropped the
	// stable gates.current gate-ID pointer, the current gate is now resolved
	// by matching the entity's live status against a record's stage
	// (recordForStage) rather than by a stable ID — so once the first call's
	// local (unpushed) advance has moved status past ideation, a repeat call
	// finds no record for the new status at all and reports "ineligible"
	// rather than the old "consumed". The guarantee this test exists to pin
	// is unaffected: no double-write, no new commit, no sync.
	out.Reset()
	errOut.Reset()
	code = run(context.Background(), []string{"gate", "consume", "task", "--workflow-dir", workflowDir},
		nil, hostClone, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 1 || !strings.Contains(out.String(), "condition=ineligible") || !strings.Contains(out.String(), "consumed=false") {
		t.Fatalf("repeat consume exit=%d stdout=%q stderr=%q, want an ineligible refusal", code, out.String(), errOut.String())
	}
	if strings.Contains(out.String(), "sync=") {
		t.Fatalf("repeat consume refusal must run no sync: %s", out.String())
	}
	if head := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD")); head != headBefore {
		t.Fatalf("repeat consume refusal created a new commit: before=%s after=%s", headBefore, head)
	}

	// Recovery: unblock, then `state commit <slug>` publishes the existing advance.
	restore()
	out.Reset()
	errOut.Reset()
	code = run(context.Background(), []string{"state", "commit", "task", "--workflow-dir", workflowDir},
		nil, hostClone, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 || !strings.Contains(out.String(), "Published previously committed state") {
		t.Fatalf("recovery state commit exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	if head := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD")); head != headBefore {
		t.Fatalf("recovery state commit created a NEW commit rather than publishing the existing one: before=%s after=%s", headBefore, head)
	}
	if fields := status.ParseFrontmatter(filepath.Join(checkout, "task.md")); fields["status"] != "implementation" {
		t.Fatalf("recovered advance status = %q, want implementation", fields["status"])
	}
}

// TestGateConsumeRepeatAfterSupersedeRunsNoSync pins finding 5 (roborev,
// branch_final): EvaluateEligibility copies the attempt's CURRENT application
// state into ApplicationState on every read, including a pure-refusal repeat
// against an application ALREADY superseded from a prior call — so a naive
// `ApplicationState == "superseded"` write-detector wrongly treats that repeat
// as a fresh write and re-runs sync, sweeping up whatever else is dirty in the
// entity's companion path and even risking a spurious HALT. The fix
// (gates.ConsumeResult.Wrote, set only on a real mutation) must keep the
// repeat byte-clean and side-effect-free: no sync line, no commit, and any
// unrelated dirt left in the entity's companion directory between the two
// calls stays untouched.
func TestGateConsumeRepeatAfterSupersedeRunsNoSync(t *testing.T) {
	workflowDir, checkout, briefing := staleableGatedSplitRootFixture(t)
	hostDir := filepath.Dir(checkout)

	mustSpacedock(t, hostDir, "gate", "record", "task", "--decision", "approve", "--actor", "person:captain", "--workflow-dir", workflowDir)

	// Drift the bound briefing (no request-digest binding in this hand-crafted
	// fixture, so this is the softer "stale reviewed input" path, not the
	// harder tamper-detection error a request-digest-bound Briefing would
	// raise) so reviewed input reads as definitively stale, not merely
	// unreadable/"unknown" (which ConsumeAt treats as ineligible instead).
	driftedOnce := strings.Replace(mustReadFile(t, briefing), "review", "revised review", 1)
	if driftedOnce == mustReadFile(t, briefing) {
		t.Fatal("drift fixture did not change the briefing content")
	}
	writeFile(t, briefing, driftedOnce)
	git(t, checkout, "add", "-A")
	git(t, checkout, "commit", "-q", "-m", "drift briefing")

	// The stale-supersede path is a documented nonzero-exit outcome (consume
	// spends nothing but does write pending->superseded), so this call is not
	// run through mustSpacedock, which requires exit 0.
	var firstOut, firstErr bytes.Buffer
	firstCode := run(context.Background(), []string{"gate", "consume", "task", "--workflow-dir", workflowDir},
		nil, hostDir, nil, &firstOut, &firstErr, &status.NativeRunner{}, nil)
	if firstCode != 1 || !strings.Contains(firstOut.String(), "condition=stale") || !strings.Contains(firstOut.String(), "sync=") {
		t.Fatalf("first (real) stale-supersede consume exit=%d stdout=%q stderr=%q, want exit 1, condition=stale, and a sync line", firstCode, firstOut.String(), firstErr.String())
	}
	if porcelain := git(t, checkout, "status", "--porcelain"); strings.TrimSpace(porcelain) != "" {
		t.Fatalf("first supersede left the checkout dirty: %q", porcelain)
	}

	// Fresh, deliberately UNCOMMITTED dirt in the flat entity's own companion
	// directory (checkout/task/, per flatEntityCommitPaths) — under the bug, a
	// repeat consume's wrongly-triggered sync would path-scope-commit this too,
	// even though nothing about the gate itself changed.
	strayArtifact := filepath.Join(checkout, "task", "stray.txt")
	if err := os.MkdirAll(filepath.Dir(strayArtifact), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, strayArtifact, "unrelated dirt\n")
	headBefore := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD"))

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "consume", "task", "--workflow-dir", workflowDir},
		nil, hostDir, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 1 || !strings.Contains(out.String(), "condition=superseded") {
		t.Fatalf("repeat consume exit=%d stdout=%q stderr=%q, want a condition=superseded refusal", code, out.String(), errOut.String())
	}
	if strings.Contains(out.String(), "sync=") {
		t.Fatalf("repeat consume after supersede wrongly ran a sync: %q", out.String())
	}
	if head := strings.TrimSpace(git(t, checkout, "rev-parse", "HEAD")); head != headBefore {
		t.Fatalf("repeat consume after supersede created a commit: before=%s after=%s", headBefore, head)
	}
	if porcelain := git(t, checkout, "status", "--porcelain", "--untracked-files=all"); !strings.Contains(porcelain, "stray.txt") {
		t.Fatalf("repeat consume swept up the unrelated stray artifact; porcelain=%q", porcelain)
	}
}

// staleableGatedSplitRootFixture hand-crafts (no real `gate prepare`) a
// split-root workflow with a flat "task" entity holding an OPEN ideation gate
// attempt bound WITHOUT a request-digest — the shape whose reviewed-input
// check (inspectReviewedInput) reads the room-ref path directly as the
// canonical Briefing file, so drifting that file's bytes produces the softer
// "stale" condition rather than gate prepare's own request-digest-bound hard
// tamper-detection error. Returns the workflow dir, the state checkout, and
// the bound briefing's path.
func staleableGatedSplitRootFixture(t *testing.T) (workflowDir, checkout, briefing string) {
	t.Helper()
	root := t.TempDir()
	mainRoot := filepath.Join(root, "main")
	workflowDir = filepath.Join(mainRoot, "docs", "dev")
	if err := os.MkdirAll(workflowDir, 0o755); err != nil {
		t.Fatal(err)
	}
	testgit.InitRepo(t, mainRoot, "-q")
	writeFile(t, filepath.Join(workflowDir, "README.md"), gateCeremonyReadme)
	git(t, mainRoot, "add", ".")
	git(t, mainRoot, "commit", "-q", "-m", "main fixture")

	checkout = filepath.Join(workflowDir, ".spacedock-state")
	if err := os.MkdirAll(checkout, 0o755); err != nil {
		t.Fatal(err)
	}
	testgit.InitRepo(t, checkout, "-q")

	briefing = filepath.Join(checkout, "review", "ideation", "briefing-1", "briefing.json")
	if err := os.MkdirAll(filepath.Dir(briefing), 0o755); err != nil {
		t.Fatal(err)
	}
	briefingBytes := []byte(`{"type":"Briefing","version":"1","id":"briefing:task:ideation:attempt-1:revision-1","question":"Advance?","artifacts":[{"id":"artifact:review","uri":"review.md","rev":"sha256:` + strings.Repeat("2", 64) + `"}],"context":[]}`)
	writeFile(t, briefing, string(briefingBytes))
	digest, err := gates.CanonicalDigest(briefingBytes)
	if err != nil {
		t.Fatal(err)
	}
	entity := filepath.Join(checkout, "task.md")
	writeFile(t, entity, "---\nid: task\nstatus: ideation\ntitle: Task\nstarted:\nworktree:\ngates:\n"+
		"  version: 1\n"+
		"  records:\n"+
		"    - id: gate:task:ideation\n"+
		"      stage: ideation\n"+
		"      attempts:\n"+
		"        - id: gate-attempt:task-ideation-1\n"+
		"          briefing: {id: 'briefing:task:ideation:attempt-1:revision-1', digest: '"+digest+"', room-ref: ./review/ideation/briefing-1/briefing.json}\n"+
		"---\n# Task\n")
	git(t, checkout, "add", "-A")
	git(t, checkout, "commit", "-q", "-m", "state fixture")
	git(t, checkout, "branch", "-M", "spacedock-state/dev")
	return workflowDir, checkout, briefing
}

func mustReadFile(t *testing.T, path string) string {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

// singleHostOrigin returns a fresh clone of a throwaway bare origin, so
// commissionGatedSplitWorkflow's push calls have a real (reachable) origin to
// push to before blockPush breaks it.
func singleHostOrigin(t *testing.T) string {
	t.Helper()
	bare := filepath.Join(t.TempDir(), "origin.git")
	git(t, t.TempDir(), "init", "-q", "--bare", bare)
	hostClone := filepath.Join(t.TempDir(), "host")
	git(t, t.TempDir(), "clone", "-q", bare, hostClone)
	git(t, hostClone, "config", "user.email", "t@t")
	git(t, hostClone, "config", "user.name", "t")
	return hostClone
}

// commissionGatedSplitWorkflow births a gate-bearing split-root workflow inside
// hostClone (already cloned from a bare origin): a real `gate prepare` binds a
// genuine room + Briefing for entity "task" at stage ideation, then the bind
// commit is pushed — the shape gate record/consume's sync-failure, HALT, and
// recovery tests need. Mirrors commissionSplitWorkflow (state_init_test.go)
// with a gate-bearing README and a prepared (not bare) entity.
func commissionGatedSplitWorkflow(t *testing.T, hostClone string) (workflowDir, stateBranch string) {
	t.Helper()
	workflowDir = filepath.Join(hostClone, "docs", "dev")
	if err := os.MkdirAll(workflowDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(workflowDir, "README.md"), gateCeremonyReadme)
	artifact := filepath.Join(hostClone, "gate-review.md")
	writeFile(t, artifact, "# Review\n\nReady to proceed to implementation.\n")
	writeFile(t, filepath.Join(hostClone, ".gitignore"), "docs/dev/.spacedock-state/\n.worktrees/\n")
	git(t, hostClone, "add", "docs/dev/README.md", "gate-review.md", ".gitignore")
	git(t, hostClone, "commit", "-q", "-m", "commission gated split workflow")

	statePath := filepath.Join(workflowDir, ".spacedock-state")
	stateBranch = "spacedock-state/dev"

	tmpWT := filepath.Join(t.TempDir(), "orphan-birth")
	git(t, hostClone, "worktree", "add", "--detach", tmpWT)
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
	grandfatherFlatRooms(t, tmpWT, "task")
	writeFile(t, filepath.Join(tmpWT, "task.md"), "---\nid: task\nstatus: ideation\ntitle: Task\nstarted:\nworktree:\n---\n# Task\n")
	git(t, tmpWT, "add", "-A")
	git(t, tmpWT, "commit", "-q", "-m", "seed state")
	git(t, tmpWT, "push", "origin", stateBranch)
	git(t, hostClone, "worktree", "remove", "--force", tmpWT)
	git(t, hostClone, "worktree", "add", statePath, stateBranch)

	mustSpacedock(t, hostClone, "gate", "prepare", "task",
		"--question", "Advance to implementation?", "--artifact", artifact,
		"--summary", "Ready to proceed to implementation.", "--workflow-dir", workflowDir)
	mustSpacedock(t, hostClone, "state", "commit", "task", "--workflow-dir", workflowDir)
	return workflowDir, stateBranch
}

// twoHostGatedWorkflow mirrors twoHostStateWorkflow (state_commit_test.go) but
// commissions a gate-bearing README + a genuinely prepared "task" gate, so both
// hosts resume the SAME prepared-but-open gate before diverging — the shape the
// two-writer HALT tests need.
func twoHostGatedWorkflow(t *testing.T) (bare, workflowA, workflowB, stateBranch string) {
	t.Helper()
	bare = filepath.Join(t.TempDir(), "origin.git")
	git(t, t.TempDir(), "init", "-q", "--bare", bare)

	hostA := filepath.Join(t.TempDir(), "hostA")
	git(t, t.TempDir(), "clone", "-q", bare, hostA)
	git(t, hostA, "config", "user.email", "a@t")
	git(t, hostA, "config", "user.name", "hostA")
	workflowA, stateBranch = commissionGatedSplitWorkflow(t, hostA)
	git(t, hostA, "push", "-q", "origin", "HEAD")

	hostB := filepath.Join(t.TempDir(), "hostB")
	git(t, t.TempDir(), "clone", "-q", bare, hostB)
	git(t, hostB, "config", "user.email", "b@t")
	git(t, hostB, "config", "user.name", "hostB")
	workflowB = filepath.Join(hostB, "docs", "dev")
	stateInit(t, hostB, workflowB)
	return bare, workflowA, workflowB, stateBranch
}

// TestGateRecordHaltsOnSameEntityConflict pins the genuinely new risk mechanism
// 1 introduces: two writers independently close the SAME open gate; A's close
// syncs and pushes first; B's close is durable locally but its OWN sync hits a
// same-entity rebase conflict and must HALT (exit 3, phase=record), leaving the
// checkout clean (rebase aborted) and A's edit — not B's — on origin.
func TestGateRecordHaltsOnSameEntityConflict(t *testing.T) {
	bare, workflowA, workflowB, stateBranch := twoHostGatedWorkflow(t)
	hostA := filepath.Dir(filepath.Dir(workflowA))
	hostB := filepath.Dir(filepath.Dir(workflowB))
	checkoutB := filepath.Join(workflowB, ".spacedock-state")

	mustSpacedock(t, hostA, "gate", "record", "task", "--decision", "approve", "--actor", "person:captain", "--workflow-dir", workflowA)

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "record", "task", "--decision", "approve", "--actor", "agent:first-officer", "--reason", "B's independent close", "--conn-quote", "you have the conn toward the sprint goal", "--conn-source", "launch runbook for this headless session", "--workflow-dir", workflowB},
		nil, hostB, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 3 {
		t.Fatalf("same-entity conflict must HALT with exit 3; got exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	if !strings.Contains(out.String(), "sync=halted phase=record") {
		t.Fatalf("stdout missing sync=halted phase=record: %s", out.String())
	}
	if !strings.Contains(errOut.String(), "HALT") || !strings.Contains(errOut.String(), "task.md") {
		t.Fatalf("HALT stderr should name the conflicting entity path; got:\n%s", errOut.String())
	}
	if porcelain := git(t, checkoutB, "status", "--porcelain"); strings.TrimSpace(porcelain) != "" {
		t.Fatalf("HALT must leave a clean checkout; porcelain:\n%s", porcelain)
	}
	if _, ok := gitOK(t, checkoutB, "push", "origin", stateBranch); ok {
		t.Fatalf("a plain push after HALT must stay rejected; B must not have force-pushed")
	}
	originTask := showOriginFile(t, bare, stateBranch, "task.md")
	if !strings.Contains(originTask, "by: person:captain") {
		t.Fatalf("origin task should carry A's Resolution; got:\n%s", originTask)
	}
}

// TestGateConsumeHaltsOnSameEntityConflict mirrors the record-phase HALT test
// one step later: both hosts resume an ALREADY-CLOSED gate; A consumes first
// and pushes; B's independent consume is durable locally (its own write) but
// its sync hits a same-entity rebase conflict and HALTs (exit 3, phase=consume).
func TestGateConsumeHaltsOnSameEntityConflict(t *testing.T) {
	bare, workflowA, workflowB, stateBranch := twoHostGatedWorkflow(t)
	hostA := filepath.Dir(filepath.Dir(workflowA))
	hostB := filepath.Dir(filepath.Dir(workflowB))
	checkoutB := filepath.Join(workflowB, ".spacedock-state")

	// Close the gate on A and let B resume that closed-but-pending state before
	// either side consumes, so both start from the identical closed gate.
	mustSpacedock(t, hostA, "gate", "record", "task", "--decision", "approve", "--actor", "person:captain", "--workflow-dir", workflowA)
	git(t, checkoutB, "pull", "-q", "--rebase", "origin", stateBranch)

	mustSpacedock(t, hostA, "gate", "consume", "task", "--workflow-dir", workflowA)

	// Consume's own write is otherwise fully deterministic (no timestamp), so an
	// UNTOUCHED B would produce byte-identical output to A's already-pushed
	// commit and git's rebase would recognize the redundant patch as already
	// applied rather than conflicting. Touch an unrelated field on B's local copy
	// first so B's consume commit genuinely diverges from A's pushed commit and
	// the rebase hits a real same-entity conflict.
	entityB := filepath.Join(checkoutB, "task.md")
	bodyB, err := os.ReadFile(entityB)
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, entityB, strings.Replace(string(bodyB), "title: Task", "title: Task (B)", 1))

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"gate", "consume", "task", "--workflow-dir", workflowB},
		nil, hostB, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 3 {
		t.Fatalf("same-entity conflict must HALT with exit 3; got exit=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	if !strings.Contains(out.String(), "sync=halted phase=consume") {
		t.Fatalf("stdout missing sync=halted phase=consume: %s", out.String())
	}
	if !strings.Contains(errOut.String(), "HALT") || !strings.Contains(errOut.String(), "task.md") {
		t.Fatalf("HALT stderr should name the conflicting entity path; got:\n%s", errOut.String())
	}
	if porcelain := git(t, checkoutB, "status", "--porcelain"); strings.TrimSpace(porcelain) != "" {
		t.Fatalf("HALT must leave a clean checkout; porcelain:\n%s", porcelain)
	}
	originTask := showOriginFile(t, bare, stateBranch, "task.md")
	if !strings.Contains(originTask, "status: implementation") {
		t.Fatalf("origin task should carry A's advance; got:\n%s", originTask)
	}

	// The HALT's whole guarantee is that B's own divergent write survives
	// LOCALLY (aborted rebase restores B's original branch tip; nothing is
	// force-pushed or discarded) even though it never reached origin — a test
	// that only inspects the clean checkout and origin's content would still
	// pass if B's local HEAD had been silently reset to match origin instead
	// of genuinely halting with B's edit intact.
	localTask := showOriginFile(t, checkoutB, "HEAD", "task.md")
	if !strings.Contains(localTask, "title: Task (B)") {
		t.Fatalf("B's local HEAD lost its own divergent edit after the HALT; got:\n%s", localTask)
	}
	if !strings.Contains(localTask, "state: consumed") || !strings.Contains(localTask, "status: implementation") {
		t.Fatalf("B's local HEAD lost its own consumed write after the HALT; got:\n%s", localTask)
	}
	if _, ok := gitOK(t, checkoutB, "push", "origin", stateBranch); ok {
		t.Fatal("a plain push after HALT must stay rejected; B's branch must still diverge from origin")
	}
}
