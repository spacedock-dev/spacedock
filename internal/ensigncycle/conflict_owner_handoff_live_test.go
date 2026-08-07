//go:build live

package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/dispatch"
	"github.com/spacedock-dev/spacedock/internal/status"
)

type conflictOwnerTuple struct {
	Entity       string
	Stage        string
	WorkerName   string
	Branch       string
	Worktree     string
	DispatchFile string
}

func TestLiveCodexOwnedConflictReturnsToRegisteredWorker(t *testing.T) {
	runner := newCodexLiveRunner(t)
	root := t.TempDir()
	entity := filepath.Join(root, "conflict-owner.md")

	writeFile(t, filepath.Join(root, "README.md"), conflictOwnerWorkflow())
	writeFile(t, filepath.Join(root, "conflict.txt"), "base\n")
	writeFile(t, entity, conflictOwnerEntity())
	gitInit(t, root)
	owner := stampConflictOwner(t, runner.binary, root, entity)
	worktree := filepath.Join(root, owner.Worktree)
	marker := filepath.Join(worktree, "owner-handoff.marker")
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
	prompt := fmt.Sprintf("Use $spacedock:first-officer for this whole run. The fixture is paused at the runtime boundary of an owned code-worktree rebase conflict for %s in %s. ", owner.Entity, owner.Stage) +
		fmt.Sprintf("Reconstitute the owner recorded by the initial stamped dispatch %s by spawning exactly one worker named %s with identity %s/%s/%s/%s; its first assignment only acknowledges readiness and it must remain addressable. ", owner.DispatchFile, owner.WorkerName, owner.Entity, owner.Stage, owner.Branch, owner.Worktree) +
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
	if got := strings.TrimSpace(git(t, worktree, "branch", "--show-current")); got != owner.Branch {
		t.Fatalf("marker branch = %q, want stamped owner branch %q", got, owner.Branch)
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
	assertConflictOwnerFreshEnvelope(t, runner.binary, root, entity, owner)
}

func stampConflictOwner(t *testing.T, binary, root, entity string) conflictOwnerTuple {
	t.Helper()
	checklist := filepath.Join(root, "initial-owner.checklist")
	writeFile(t, checklist, "Acknowledge readiness and remain addressable; do not change entity or code before the same-stage follow-up.\n")
	cmd := exec.Command(binary, "dispatch", "build", "--stamp", "--host", "codex", "--workflow-dir", root, "--entity-path", entity, "--stage", "implementation", "--checklist-file", checklist)
	cmd.Env = append(os.Environ(), "HOME="+t.TempDir())
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("initial stamped owner build: %v\n%s", err, out)
	}
	spawn, err := dispatch.CodexMultiAgentV2SpawnInput(out)
	if err != nil {
		t.Fatalf("record initial stamped owner: %v\n%s", err, out)
	}
	fields := status.ParseFrontmatter(entity)
	worktree := fields["worktree"]
	if fields["started"] == "" || worktree == "" {
		t.Fatalf("initial dispatch did not stamp owner checkout: %#v", fields)
	}
	if subject := strings.TrimSpace(git(t, root, "log", "-1", "--format=%s")); subject != "dispatch: conflict-owner entering implementation" {
		t.Fatalf("initial stamped dispatch commit = %q, want owner-entry commit", subject)
	}
	worktreePath := filepath.Join(root, worktree)
	branch := strings.TrimSpace(git(t, worktreePath, "branch", "--show-current"))
	owner := conflictOwnerTuple{
		Entity:       spawn.Identity.Slug,
		Stage:        spawn.Identity.Stage,
		WorkerName:   spawn.Identity.Name,
		Branch:       branch,
		Worktree:     worktree,
		DispatchFile: readInitialDispatchPath(t, out),
	}
	if owner.Entity == "" || owner.Stage == "" || owner.WorkerName == "" || owner.Branch == "" {
		t.Fatalf("initial stamped dispatch produced incomplete owner tuple: %#v", owner)
	}
	return owner
}

func readInitialDispatchPath(t *testing.T, raw []byte) string {
	t.Helper()
	var envelope struct {
		DispatchFile string `json:"dispatch_file_path"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		t.Fatalf("decode initial stamped owner envelope: %v\n%s", err, raw)
	}
	if envelope.DispatchFile == "" {
		t.Fatalf("initial stamped owner envelope lacks dispatch file: %s", raw)
	}
	return envelope.DispatchFile
}

func assertConflictOwnerFreshEnvelope(t *testing.T, binary, root, entity string, owner conflictOwnerTuple) {
	t.Helper()
	checklist := filepath.Join(root, "fresh-owner.checklist")
	scope := filepath.Join(root, "fresh-owner.scope.md")
	writeFile(t, checklist, "Write the owner marker on the registered branch.\n")
	writeFile(t, scope, "entity: "+owner.Entity+"\nstage: "+owner.Stage+"\npr: 632\nbranch: "+owner.Branch+"\nworktree: "+owner.Worktree+"\nold base/head: base / worker\nmoved base: main\nconflict paths: conflict.txt\nnext owner action: reconcile and write marker\n")
	cmd := exec.Command(binary, "dispatch", "build", "--host", "codex", "--workflow-dir", root, "--entity-path", entity, "--stage", owner.Stage, "--checklist-file", checklist, "--scope-notes-file", scope)
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
	if envelope.Name != owner.WorkerName || envelope.DispatchFile == "" {
		t.Fatalf("fresh owner envelope identity = %q, want stamped owner %q: %s", envelope.Name, owner.WorkerName, out)
	}
	body := readFile(t, envelope.DispatchFile)
	for _, want := range []string{owner.Entity, owner.Stage, owner.Branch, owner.Worktree, "conflict paths: conflict.txt"} {
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

func conflictOwnerEntity() string {
	return "---\n" +
		"title: Conflict owner\n" +
		"status: implementation\n" +
		"started:\n" +
		"pr: 632\n" +
		"mod-block: preserved\n" +
		"worktree:\n" +
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
