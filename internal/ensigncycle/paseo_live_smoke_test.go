//go:build live

// ABOUTME: Isolated-daemon live E2E smoke proving a Paseo-hosted FO dispatches a
// ABOUTME: fresh Paseo worker via --host paseo, verifies Git commits, and archives it.
package ensigncycle

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spacedock-dev/spacedock/internal/status"
)

const paseoLiveSmokeMarker = "PASEO-LIVE-SMOKE-OK"

func TestLivePaseoDispatchSmoke(t *testing.T) {
	d := startIsolatedPaseoDaemon(t)
	fx := setupPaseoSplitRootFixture(t, d.CanonicalRoot)
	d.stageNonCodexSkillsIfNeeded(t, fx.RepoRoot)

	var foWs struct {
		WorkspaceID string `json:"workspaceId"`
	}
	d.runJSON(t, &foWs, "workspace", "create", "--isolation", "local", "--path", fx.RepoRoot, "--title", "paseo-smoke-fo", "--json")
	if foWs.WorkspaceID == "" {
		t.Fatal("FO workspace create returned empty workspaceId")
	}

	thinkingNote := ""
	if d.Thinking != "" {
		thinkingNote = fmt.Sprintf(" with --thinking %s", d.Thinking)
	}
	foPrompt := fmt.Sprintf("Load the installed Spacedock first-officer skill by name (do not hardcode checkout paths) and run single-entity implementation dispatch for entity %q in workflow %s (entity file: %s).\n\n"+
		"1. Follow your loaded first-officer contract and Paseo dispatch binding: run `${SPACEDOCK_BIN:-spacedock} dispatch build --host paseo --stamp` for stage `implementation` with a checklist asking the ensign to:\n"+
		"   - Write `receipt.txt` containing `%s` in the assigned code worktree and commit it cleanly on the worktree branch.\n"+
		"   - Append `## Stage Report: implementation` (with `- DONE:` items and `### Summary`) to %s without modifying YAML frontmatter, and commit ONLY that file path-scoped in %s.\n"+
		"2. Adopt the stamped code worktree via Paseo find-or-create workspace and spawn ONE fresh Paseo worker using `--provider %s`%s. Forward the emitted `prompt` from `dispatch build --host paseo` verbatim (do NOT add raw checkout skill paths).\n"+
		"3. Wait until the worker becomes idle, verify the durable commits in both the code worktree and state checkout, then archive ONLY the worker (never archive yourself).\n"+
		"4. Once verified and archived, send one concise final summary and stop.",
		fx.Slug, fx.WorkflowDir, fx.EntityPath, paseoLiveSmokeMarker, fx.EntityPath, fx.StateRoot, d.WorkerProvider, thinkingNote)

	var foRun struct {
		AgentID string `json:"agentId"`
	}
	d.runJSON(t, &foRun, paseoRunArgs(foWs.WorkspaceID, d.FOProvider, d.Thinking, "paseo-smoke-fo", []string{"role=first-officer"}, []string{"SPACEDOCK_BIN=" + d.SpacedockBin}, foPrompt)...)
	if foRun.AgentID == "" {
		t.Fatal("paseo run for FO returned empty agentId")
	}
	waitForPaseoAgentIdle(t, d, foRun.AgentID, d.Timeout)

	fm := status.ParseFrontmatter(fx.EntityPath)
	wtRel := strings.TrimSpace(fm["worktree"])
	if fm["status"] != "implementation" || wtRel == "" {
		t.Fatalf("frontmatter status=%q worktree=%q", fm["status"], wtRel)
	}
	codeWorktree := status.PyJoin(fx.RepoRoot, wtRel)
	if info, err := os.Stat(codeWorktree); err != nil || !info.IsDir() {
		t.Fatalf("stamped worktree %s missing: %v", codeWorktree, err)
	}
	if branch := strings.TrimSpace(git(t, codeWorktree, "rev-parse", "--abbrev-ref", "HEAD")); branch == "" || branch == "main" || branch == "HEAD" {
		t.Fatalf("code worktree branch = %q, want non-main branch", branch)
	}
	if receipt := readFile(t, filepath.Join(codeWorktree, "receipt.txt")); !strings.Contains(receipt, paseoLiveSmokeMarker) {
		t.Fatalf("receipt.txt missing %q:\n%s", paseoLiveSmokeMarker, receipt)
	}
	if changed := strings.TrimSpace(git(t, codeWorktree, "log", "-1", "--name-only", "--pretty=format:")); !strings.Contains(changed, "receipt.txt") {
		t.Fatalf("latest code commit did not touch receipt.txt:\n%s", changed)
	}
	if dirty := strings.TrimSpace(git(t, codeWorktree, "status", "--porcelain")); dirty != "" {
		t.Fatalf("code worktree not clean:\n%s", dirty)
	}

	entityText := readFile(t, fx.EntityPath)
	for _, want := range []string{"## Stage Report: implementation", "- DONE:", "### Summary"} {
		if !strings.Contains(entityText, want) {
			t.Fatalf("entity missing %q:\n%s", want, entityText)
		}
	}
	if stateChanged := strings.TrimSpace(git(t, fx.StateRoot, "log", "-1", "--name-only", "--pretty=format:")); stateChanged != filepath.ToSlash(filepath.Join(fx.Slug, "index.md")) {
		t.Fatalf("state commit touched %q, want only %s/index.md", stateChanged, fx.Slug)
	}
	if dirty := strings.TrimSpace(git(t, fx.StateRoot, "status", "--porcelain")); dirty != "" {
		t.Fatalf("state checkout not clean:\n%s", dirty)
	}

	var workers []struct {
		ID     string            `json:"id"`
		Labels map[string]string `json:"labels"`
	}
	d.runJSON(t, &workers, "ls", "--all", "--global", "--label", "spacedock.worker=paseo-smoke-implementation", "--json")
	if len(workers) != 1 || workers[0].Labels["spacedock.epoch"] != "paseo-smoke-implementation-c1-a1" {
		t.Fatalf("unexpected workers: %+v", workers)
	}
	var wi struct {
		ParentAgentID string  `json:"ParentAgentId"`
		Archived      bool    `json:"Archived"`
		ArchivedAt    *string `json:"ArchivedAt"`
	}
	d.runJSON(t, &wi, "inspect", workers[0].ID, "--json")
	if wi.ParentAgentID != foRun.AgentID || !wi.Archived || wi.ArchivedAt == nil || *wi.ArchivedAt == "" {
		t.Fatalf("worker inspect mismatch: %+v (want parent %s, archived)", wi, foRun.AgentID)
	}
	var fi struct {
		Archived bool `json:"Archived"`
	}
	d.runJSON(t, &fi, "inspect", foRun.AgentID, "--json")
	if fi.Archived {
		t.Fatalf("FO agent %s must remain unarchived", foRun.AgentID)
	}
}

func waitForPaseoAgentIdle(t *testing.T, d *paseoDaemon, agentID string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			t.Fatalf("Paseo agent %s did not reach idle within %s", agentID, timeout)
		}
		waitSecs := min(60, max(1, int(remaining.Seconds())))
		var res struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}
		d.runJSON(t, &res, "wait", agentID, "--timeout", fmt.Sprintf("%d", waitSecs), "--json")
		if res.Status == "idle" {
			return
		}
		if res.Status != "timeout" {
			t.Fatalf("Paseo agent %s stopped with status=%q message=%q", agentID, res.Status, res.Message)
		}
	}
}
