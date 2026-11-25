//go:build live

// ABOUTME: Isolated-daemon live E2E smoke proving a restarted Paseo FO adopts an
// ABOUTME: in-flight worker purely from Paseo labels and Git state with zero lease files.
package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/spacedock-dev/spacedock/internal/status"
)

const paseoAdoptionMarker = "PASEO-ADOPTION-SMOKE-OK"

func TestLivePaseoAdoptionSmoke(t *testing.T) {
	d := startIsolatedPaseoDaemon(t, func(cfg *paseoDaemonConfig) { cfg.artifactSubdir = "paseo-adoption-smoke" })
	fx := setupPaseoSplitRootFixture(t, d.CanonicalRoot)
	releaseFile := filepath.Join(d.CanonicalRoot, ".adoption-release")
	mustWriteFile(t, filepath.Join(fx.WorkflowDir, "README.md"), fmt.Sprintf("---\ncommissioned-by: spacedock@test\nentity-type: task\nid-style: slug\nstate: .spacedock-state\nstages:\n  defaults:\n    worktree: false\n    concurrency: 1\n  states:\n    - name: backlog\n      initial: true\n    - name: implementation\n      worktree: true\n    - name: done\n      terminal: true\n---\n# Paseo Adoption Workflow\n\n### backlog\n\nSeed.\n\n### implementation\n\n1. Write `receipt.txt` containing `%s` on the task branch and commit cleanly.\n2. Poll up to 180s (1s interval) until `%s` exists.\n3. Append `## Stage Report: implementation` (`- DONE:` items + `### Summary`) without editing frontmatter and commit path-scoped in `.spacedock-state`.\n\n### done\n\nTerminal.\n", paseoAdoptionMarker, releaseFile))
	git(t, fx.RepoRoot, "add", "docs/smoke/README.md")
	git(t, fx.RepoRoot, "commit", "-q", "-m", "configure adoption smoke workflow")
	d.stageNonCodexSkillsIfNeeded(t, fx.RepoRoot)
	leaseDir := filepath.Join(strings.TrimSpace(git(t, fx.WorkflowDir, "rev-parse", "--path-format=absolute", "--git-common-dir")), "spacedock-paseo-leases")

	var foWs struct{ WorkspaceID string }
	d.runJSON(t, &foWs, "workspace", "create", "--isolation", "local", "--path", fx.RepoRoot, "--title", "paseo-adoption-fo", "--json")
	extraFlags := ""
	if d.Thinking != "" {
		extraFlags = " --thinking " + d.Thinking
	}
	if strings.HasPrefix(d.WorkerProvider, "codex/") {
		extraFlags += " --mode full-access"
	}
	workerBase := fx.Slug + "-implementation"
	wantEpoch := workerBase + "-c1-a1"
	fo1Prompt := fmt.Sprintf("Use $spacedock:first-officer (load the installed Spacedock first-officer skill by name; no raw checkout paths) and dispatch entity %q in workflow %s (entity: %s) for stage `implementation` via Paseo.\nUse worker flags `--provider %s%s`, pass the stage checklist (write `receipt.txt` containing `%s` on the task branch and commit cleanly; poll every 1s up to 180s until `%s` exists; append `## Stage Report: implementation` and commit path-scoped in `.spacedock-state`), and wait on the worker.",
		fx.Slug, fx.WorkflowDir, fx.EntityPath, d.WorkerProvider, extraFlags, paseoAdoptionMarker, releaseFile)
	var fo1, fo2 struct{ AgentID string }
	d.runJSON(t, &fo1, paseoRunArgs(foWs.WorkspaceID, d.FOProvider, d.Thinking, "paseo-adoption-fo1", []string{"role=first-officer"}, []string{"SPACEDOCK_BIN=" + d.SpacedockBin}, fo1Prompt)...)
	if foWs.WorkspaceID == "" || fo1.AgentID == "" {
		t.Fatalf("empty workspace=%q or FO-1=%q", foWs.WorkspaceID, fo1.AgentID)
	}
	t.Cleanup(func() { _, _ = runPaseoJSON(d.Bin, d.Home, d.Env, nil, "archive", fo1.AgentID, "--force", "--json") })

	type paseoInspect struct {
		Status, ParentAgentID string
		Archived              bool
		ArchivedAt            *string
	}
	inspect := func(id string) (out paseoInspect) { d.runJSON(t, &out, "inspect", id, "--json"); return }

	w1ID, stampedWorktree := waitForRunningPaseoWorker(t, d, fx, leaseDir, fo1.AgentID, workerBase, wantEpoch, d.Timeout/2)
	dispatchRev := strings.TrimSpace(git(t, fx.StateRoot, "rev-parse", "HEAD"))
	assertNoLeaseDirOrEntityPollution(t, leaseDir, fx.EntityPath, "after FO-1 dispatch")
	d.runJSON(t, nil, "stop", fo1.AgentID, "--json")
	if w1 := inspect(w1ID); w1.Status != "running" || w1.ParentAgentID != fo1.AgentID || w1.Archived {
		t.Fatalf("after FO-1 stop, W1 = %+v, want running parent=%s unarchived", w1, fo1.AgentID)
	}

	fo2Prompt := fmt.Sprintf("Use $spacedock:first-officer (load the installed Spacedock first-officer skill by name; no raw checkout paths) and resume workflow %s for entity %q (entity: %s) via Paseo.\nContinue from durable workflow state through stage completion to `done` (`verdict: PASSED`) and stop. If worker launch flags are ever needed, use `--provider %s%s`.",
		fx.WorkflowDir, fx.Slug, fx.EntityPath, d.WorkerProvider, extraFlags)
	d.runJSON(t, &fo2, paseoRunArgs(foWs.WorkspaceID, d.FOProvider, d.Thinking, "paseo-adoption-fo2", []string{"role=first-officer"}, []string{"SPACEDOCK_BIN=" + d.SpacedockBin}, fo2Prompt)...)
	if fo2.AgentID == "" {
		t.Fatal("FO-2 launch returned empty agentId")
	}
	waitForPaseoAdoptionCall(t, d, fo2.AgentID, w1ID, 90*time.Second)
	if s := inspect(w1ID).Status; s != "running" {
		t.Fatalf("W1 status = %q before .adoption-release, want running", s)
	}
	mustWriteFile(t, releaseFile, "release\n")
	waitForPaseoAgentIdle(t, d, fo2.AgentID, d.Timeout)

	var labeled []struct {
		ID          string            `json:"id"`
		WorkspaceID string            `json:"workspaceId"`
		Labels      map[string]string `json:"labels"`
	}
	d.runJSON(t, &labeled, "ls", "--all", "--global", "--label", "spacedock.worker="+workerBase, "--json")
	if len(labeled) != 1 || labeled[0].ID != w1ID || labeled[0].Labels["spacedock.epoch"] != wantEpoch || labeled[0].WorkspaceID == "" {
		t.Fatalf("want exactly 1 worker %s epoch=%s with workspaceId, got %+v", w1ID, wantEpoch, labeled)
	}
	var allAgents []struct{ ID, Name string }
	d.runJSON(t, &allAgents, "ls", "--all", "--global", "--json")
	for _, a := range allAgents {
		if a.ID != fo1.AgentID && a.ID != fo2.AgentID && a.ID != w1ID {
			t.Fatalf("unexpected extra agent %s (%q); want only FO-1/FO-2/%s", a.ID, a.Name, w1ID)
		}
	}
	w1, fo2Inspect := inspect(w1ID), inspect(fo2.AgentID)
	if !w1.Archived || w1.ArchivedAt == nil || *w1.ArchivedAt == "" || fo2Inspect.Archived {
		t.Fatalf("archive mismatch: W1=%+v FO-2.Archived=%v", w1, fo2Inspect.Archived)
	}

	entityPath := fx.EntityPath
	for _, p := range []string{filepath.Join(fx.StateRoot, "_archive", fx.Slug, "index.md"), filepath.Join(fx.StateRoot, "_archive", fx.Slug+".md")} {
		if _, err := os.Stat(p); err == nil {
			entityPath = p
		}
	}
	assertNoLeaseDirOrEntityPollution(t, leaseDir, entityPath, "after completion")
	fm := status.ParseFrontmatter(entityPath)
	if fm["status"] != "done" || fm["verdict"] != "PASSED" {
		t.Fatalf("entity %s frontmatter status=%q verdict=%q, want done/PASSED", entityPath, fm["status"], fm["verdict"])
	}
	wtRel := strings.TrimSpace(fm["worktree"])
	if wtRel == "" {
		wtRel = stampedWorktree
	}
	codeWorktree := status.PyJoin(fx.RepoRoot, wtRel)
	if canon, err := filepath.EvalSymlinks(codeWorktree); err == nil {
		codeWorktree = canon
	}
	codeDir, merged := codeWorktree, false
	if _, err := os.Stat(codeWorktree); os.IsNotExist(err) {
		codeDir, merged = fx.RepoRoot, true
	} else if err != nil {
		t.Fatalf("stat code worktree %s: %v", codeWorktree, err)
	}
	if b := strings.TrimSpace(git(t, codeDir, "rev-parse", "--abbrev-ref", "HEAD")); (!merged && (b == "" || b == "main" || b == "HEAD")) || (merged && b != "main") {
		t.Fatalf("code branch = %q (merged=%v)", b, merged)
	}
	if (!merged && !slices.Contains(strings.Split(strings.TrimSpace(git(t, codeDir, "log", "-1", "--name-only", "--pretty=format:")), "\n"), "receipt.txt")) ||
		(merged && strings.TrimSpace(git(t, codeDir, "log", "--format=%P", "--", "receipt.txt")) == "") || !strings.Contains(git(t, codeDir, "show", "HEAD:receipt.txt"), paseoAdoptionMarker) {
		t.Fatalf("code checkout %s (merged=%v) missing committed receipt.txt with %q", codeDir, merged, paseoAdoptionMarker)
	}
	entityRel, _ := filepath.Rel(fx.StateRoot, entityPath)
	relSlash, baseRel, foundReport := filepath.ToSlash(entityRel), fx.Slug+"/index.md", false
	for _, sha := range strings.Fields(git(t, fx.StateRoot, "log", "--follow", "--format=%H", dispatchRev+"..HEAD", "--", relSlash)) {
		files := strings.Split(strings.TrimSpace(git(t, fx.StateRoot, "show", "--format=", "--name-only", sha)), "\n")
		if len(files) != 1 || (files[0] != relSlash && files[0] != baseRel) {
			continue
		}
		var hasHdr, hasDone, hasSum bool
		for _, ln := range strings.Split(git(t, fx.StateRoot, "show", "--format=", "-U0", sha, "--", files[0]), "\n") {
			if strings.HasPrefix(ln, "+") && !strings.HasPrefix(ln, "+++ ") {
				s := strings.TrimSpace(ln[1:])
				hasHdr, hasDone, hasSum = hasHdr || s == "## Stage Report: implementation", hasDone || strings.HasPrefix(s, "- DONE:"), hasSum || s == "### Summary"
			}
		}
		foundReport = foundReport || (hasHdr && hasDone && hasSum)
	}
	if !foundReport {
		t.Fatalf("no entity-only post-%s commit in %s added complete ## Stage Report: implementation", dispatchRev, entityPath)
	}
	dirty := strings.TrimSpace(git(t, codeDir, "status", "--porcelain")) + strings.TrimSpace(git(t, fx.StateRoot, "status", "--porcelain"))
	if !merged {
		dirty += strings.TrimSpace(git(t, fx.RepoRoot, "status", "--porcelain"))
	}
	if dirty != "" {
		t.Fatalf("unclean checkout(s): %q", dirty)
	}
}

func waitForRunningPaseoWorker(t *testing.T, d *paseoDaemon, fx paseoSplitRootFixture, leaseDir, fo1ID, workerBase, wantEpoch string, timeout time.Duration) (string, string) {
	t.Helper()
	for deadline := time.Now().Add(timeout); time.Now().Before(deadline); time.Sleep(time.Second) {
		var workers []struct {
			ID     string            `json:"id"`
			Status string            `json:"status"`
			Labels map[string]string `json:"labels"`
		}
		d.runJSON(t, &workers, "ls", "--all", "--global", "--label", "spacedock.worker="+workerBase, "--json")
		assertNoLeaseDirOrEntityPollution(t, leaseDir, fx.EntityPath, "during FO-1 dispatch")
		fm := status.ParseFrontmatter(fx.EntityPath)
		wtRel := strings.TrimSpace(fm["worktree"])
		if len(workers) == 1 && workers[0].Status == "running" && workers[0].Labels["spacedock.epoch"] == wantEpoch &&
			fm["status"] == "implementation" && strings.TrimSpace(fm["started"]) != "" && wtRel != "" &&
			strings.TrimSpace(git(t, fx.StateRoot, "status", "--porcelain")) == "" {
			return workers[0].ID, wtRel
		}
		var fo1 struct{ Status string }
		d.runJSON(t, &fo1, "inspect", fo1ID, "--json")
		if fo1.Status == "error" || fo1.Status == "closed" {
			t.Fatalf("FO-1 exited (%q) before spawning running worker %s (workers=%+v)", fo1.Status, wantEpoch, workers)
		}
	}
	t.Fatalf("timed out waiting for FO-1 to spawn running worker %s (%s)", workerBase, wantEpoch)
	return "", ""
}

func assertNoLeaseDirOrEntityPollution(t *testing.T, leaseDir, entityPath, phase string) {
	t.Helper()
	if _, err := os.Stat(leaseDir); !os.IsNotExist(err) {
		t.Fatalf("%s: forbidden lease directory %s exists (stat err=%v)", phase, leaseDir, err)
	}
	doc := readFile(t, entityPath)
	if strings.Contains(doc, "## Paseo Worker") || strings.Contains(doc, "### Dispatch Retries") {
		t.Fatalf("%s: entity %s contains forbidden ## Paseo Worker or ### Dispatch Retries heading", phase, entityPath)
	}
}

func waitForPaseoAdoptionCall(t *testing.T, d *paseoDaemon, fo2ID, w1ID string, timeout time.Duration) {
	t.Helper()
	var lastTools string
	for deadline := time.Now().Add(timeout); time.Now().Before(deadline); time.Sleep(time.Second) {
		var logs struct {
			Items []struct {
				Type   string
				Detail struct {
					Command       string
					Input, Output json.RawMessage
				}
			}
		}
		d.runJSON(t, &logs, "logs", fo2ID, "--filter", "tools", "--json")
		var summary []string
		for _, it := range logs.Items {
			cmd := it.Detail.Command
			if cmd == "" {
				cmd = string(it.Detail.Input)
			}
			adoptCmd := strings.Contains(cmd, "paseo inspect") || strings.Contains(cmd, "paseo wait") ||
				strings.Contains(cmd, "paseo agent inspect") || strings.Contains(cmd, "paseo agent wait") ||
				(strings.Contains(cmd, `"paseo"`) && (strings.Contains(cmd, `"inspect"`) || strings.Contains(cmd, `"wait"`)))
			if it.Type == "tool_call" && adoptCmd && (strings.Contains(cmd, w1ID) || strings.Contains(string(it.Detail.Output), w1ID)) {
				return
			}
			snippet := strings.Join(strings.Fields(cmd), " ")
			if r := []rune(snippet); len(r) > 120 {
				snippet = string(r[:120]) + "..."
			}
			summary = append(summary, fmt.Sprintf("%s:%q", it.Type, snippet))
		}
		if len(summary) > 16 {
			summary = append(append(summary[:8:8], fmt.Sprintf("...+%d omitted", len(summary)-16)), summary[len(summary)-8:]...)
		}
		lastTools = fmt.Sprintf("%d items [%s]", len(logs.Items), strings.Join(summary, "; "))
		var fo2, w1 struct{ Status string }
		d.runJSON(t, &fo2, "inspect", fo2ID, "--json")
		d.runJSON(t, &w1, "inspect", w1ID, "--json")
		if fo2.Status == "idle" || fo2.Status == "error" || fo2.Status == "closed" || w1.Status != "running" {
			t.Fatalf("before FO-2 adoption call: FO-2=%q W1=%q (want running W1 before release); FO-2 tools: %s", fo2.Status, w1.Status, lastTools)
		}
	}
	t.Fatalf("timed out after %s waiting for FO-2 %s to inspect/wait on running worker %s; FO-2 tools: %s", timeout, fo2ID, w1ID, lastTools)
}
