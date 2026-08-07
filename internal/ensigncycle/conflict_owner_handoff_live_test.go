//go:build live

package ensigncycle

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestLiveCodexOwnedConflictReturnsToRegisteredWorker(t *testing.T) {
	runner := newCodexLiveRunner(t)
	root := t.TempDir()
	worktreeRel := filepath.Join(".worktrees", "spacedock-ensign-conflict-owner")
	worktree := filepath.Join(root, worktreeRel)
	entity := filepath.Join(root, "conflict-owner.md")
	marker := filepath.Join(worktree, "owner-handoff.marker")

	writeFile(t, filepath.Join(root, "README.md"), conflictOwnerWorkflow())
	writeFile(t, filepath.Join(root, "conflict.txt"), "base\n")
	writeFile(t, entity, conflictOwnerEntity(worktreeRel))
	gitInit(t, root)
	git(t, root, "worktree", "add", "-q", "-b", "spacedock-ensign/conflict-owner", worktree)
	writeFile(t, filepath.Join(worktree, "conflict.txt"), "worker branch\n")
	gitAsCaptain(t, worktree, "add", "conflict.txt")
	gitAsCaptain(t, worktree, "commit", "-q", "-m", "worker: change conflict target")
	writeFile(t, filepath.Join(root, "conflict.txt"), "moved main\n")
	gitAsCaptain(t, root, "add", "conflict.txt")
	gitAsCaptain(t, root, "commit", "-q", "-m", "captain: move target")

	cmd := exec.Command("git", "-C", worktree, "rebase", "main")
	if out, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("fixture rebase unexpectedly succeeded: %s", out)
	}
	before := readFile(t, entity)

	scenario := sharedRuntimeScenario{name: "owned-conflict-owner-handoff"}
	prompt := "Use $spacedock:first-officer for this whole run. The fixture is paused at the runtime boundary of an owned code-worktree rebase conflict for conflict-owner in implementation. " +
		"Reconstitute the recorded live owner by spawning exactly one implementation ensign with identity conflict-owner/implementation/spacedock-ensign/conflict-owner/" + worktreeRel + "; its first assignment only acknowledges readiness and it must remain addressable. " +
		"Then handle the existing conflict exactly through the first-officer conflict-owner contract. The owner handoff's next action is to write owner-handoff.marker containing runtime-worker-owner, commit it on the registered branch using the shared Git identity user.name=Captain and user.email=captain@example.test, and leave conflict.txt unresolved beyond the required rebase abort. " +
		"Do not edit the entity. Stop after the owner completes the follow-up."
	result, err := runner.run(t, scenario, root, prompt)
	if err != nil {
		t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
	}

	if after := readFile(t, entity); after != before {
		t.Fatalf("owner handoff changed authority bytes\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if got := strings.TrimSpace(readFile(t, marker)); got != "runtime-worker-owner" {
		t.Fatalf("worker marker = %q, want runtime-worker-owner", got)
	}
	if got := strings.TrimSpace(git(t, worktree, "branch", "--show-current")); got != "spacedock-ensign/conflict-owner" {
		t.Fatalf("marker branch = %q, want registered branch", got)
	}
	if got := strings.TrimSpace(git(t, worktree, "log", "-1", "--format=%an <%ae>")); got != "Captain <captain@example.test>" {
		t.Fatalf("marker Git author = %q, want shared Captain credential", got)
	}
	rebaseDir := strings.TrimSpace(git(t, worktree, "rev-parse", "--git-path", "rebase-merge"))
	if _, err := os.Stat(rebaseDir); !os.IsNotExist(err) {
		t.Fatalf("rebase was not cleanly aborted: %v", err)
	}
	for _, event := range []string{"spawn_agent", "followup_task"} {
		if !strings.Contains(result.jsonl, event) {
			t.Errorf("Codex stream lacks actual %s runtime call; artifacts: %s", event, result.artifactDir)
		}
	}
	assertConflictOwnerFreshEnvelope(t, runner.binary, root, entity, worktreeRel)
}

func assertConflictOwnerFreshEnvelope(t *testing.T, binary, root, entity, worktree string) {
	t.Helper()
	checklist := filepath.Join(root, "fresh-owner.checklist")
	scope := filepath.Join(root, "fresh-owner.scope.md")
	writeFile(t, checklist, "Write the owner marker on the registered branch.\n")
	writeFile(t, scope, "entity: conflict-owner\nstage: implementation\npr: 632\nbranch: spacedock-ensign/conflict-owner\nworktree: "+worktree+"\nold base/head: base / worker\nmoved base: main\nconflict paths: conflict.txt\nnext owner action: reconcile and write marker\n")
	cmd := exec.Command(binary, "dispatch", "build", "--host", "codex", "--workflow-dir", root, "--entity-path", entity, "--stage", "implementation", "--checklist-file", checklist, "--scope-notes-file", scope)
	cmd.Env = append(os.Environ(), "HOME="+t.TempDir())
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("fresh owner build: %v\n%s", err, out)
	}
	var envelope struct {
		Name         string `json:"name"`
		DispatchFile string `json:"dispatch_file_path"`
	}
	if err := json.Unmarshal(out, &envelope); err != nil {
		t.Fatalf("decode fresh owner envelope: %v\n%s", err, out)
	}
	if envelope.Name == "" || envelope.DispatchFile == "" {
		t.Fatalf("fresh owner envelope lacks spawn identity: %s", out)
	}
	body := readFile(t, envelope.DispatchFile)
	for _, want := range []string{"implementation", "spacedock-ensign/conflict-owner", worktree, "conflict paths: conflict.txt"} {
		if !strings.Contains(body, want) {
			t.Errorf("fresh owner dispatch body missing %q\n%s", want, body)
		}
	}
}

func gitAsCaptain(t *testing.T, dir string, args ...string) string {
	t.Helper()
	full := append([]string{"-C", dir, "-c", "user.name=Captain", "-c", "user.email=captain@example.test"}, args...)
	out, err := exec.Command("git", full...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func conflictOwnerWorkflow() string {
	return "---\n" +
		"commissioned-by: spacedock@1\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: true\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: backlog\n" +
		"      initial: true\n" +
		"    - name: implementation\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Conflict owner fixture\n\n" +
		"### implementation\n\nReconcile only through the registered owner.\n\n- **Outputs:** owner marker on the registered branch.\n"
}

func conflictOwnerEntity(worktree string) string {
	return "---\n" +
		"title: Conflict owner\n" +
		"status: implementation\n" +
		"started: 2026-08-07T00:00:00Z\n" +
		"pr: 632\n" +
		"mod-block: preserved\n" +
		"worktree: " + worktree + "\n" +
		"gates:\n" +
		"  version: 1\n" +
		"  records:\n" +
		"    - id: gate:fixture:backlog\n" +
		"      stage: backlog\n" +
		"      application:\n" +
		"        target-stage: implementation\n" +
		"        state: consumed\n" +
		"---\n\nOwned moving-target conflict fixture.\n"
}
