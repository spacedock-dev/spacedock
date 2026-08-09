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

type conflictOwnerFixture struct {
	root, entity, worktree, marker, before string
	owner                                  conflictOwnerTuple
}

//spacedock:live-fixture id=conflict-owner/stamped-checkout
func writeConflictOwnerFixture(t *testing.T) conflictOwnerFixture {
	t.Helper()
	root := t.TempDir()
	entity := filepath.Join(root, "conflict-owner.md")

	writeFile(t, filepath.Join(root, "README.md"), conflictOwnerWorkflow())
	writeFile(t, filepath.Join(root, "conflict.txt"), "base\n")
	writeFile(t, entity, conflictOwnerEntity())
	gitInit(t, root)
	owner := stampConflictOwner(t, spacedockBinary(t), root, entity)
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
	return conflictOwnerFixture{root: root, entity: entity, worktree: worktree, marker: marker, before: readFile(t, entity), owner: owner}
}

func runConflictOwnerHandoffJourney(t *testing.T, driver liveDriver, scenario sharedRuntimeScenario, build func(*testing.T) conflictOwnerFixture, assert func(*testing.T, conflictOwnerFixture, liveResult)) {
	t.Helper()
	fixture := build(t)
	owner := fixture.owner
	prompt := fmt.Sprintf("Use $spacedock:first-officer for this whole run. The fixture is paused at the runtime boundary of an owned code-worktree rebase conflict for %s in %s. ", owner.Entity, owner.Stage) +
		fmt.Sprintf("Reconstitute the owner recorded by the initial stamped dispatch %s by spawning exactly one worker named %s with identity %s/%s/%s/%s; its first assignment only acknowledges readiness and it must remain addressable. ", owner.DispatchFile, owner.WorkerName, owner.Entity, owner.Stage, owner.Branch, owner.Worktree) +
		"Then handle the existing conflict exactly through the first-officer conflict-owner contract. The owner handoff's next action is to write owner-handoff.marker containing runtime-worker-owner, commit it on the registered branch using the shared Git identity user.name=Captain and user.email=captain@example.test, and leave conflict.txt unresolved beyond the required rebase abort. " +
		"Do not edit the entity. Stop after the owner completes the follow-up."
	result := driver.run(t, scenario, fixture.root, prompt)
	assert(t, fixture, result)
	driver.emitMetrics(t, scenario, result)
}

func assertConflictOwnerHandoff(t *testing.T, fixture conflictOwnerFixture, _ liveResult) {
	t.Helper()
	if after := readFile(t, fixture.entity); after != fixture.before {
		t.Fatalf("owner handoff changed authority bytes\nbefore:\n%s\nafter:\n%s", fixture.before, after)
	}
	if got := strings.TrimSpace(readFile(t, fixture.marker)); got != "runtime-worker-owner" {
		t.Fatalf("worker marker = %q, want runtime-worker-owner", got)
	}
	if got := strings.TrimSpace(git(t, fixture.worktree, "branch", "--show-current")); got != fixture.owner.Branch {
		t.Fatalf("marker branch = %q, want stamped owner branch %q", got, fixture.owner.Branch)
	}
	if got := strings.TrimSpace(git(t, fixture.worktree, "log", "-1", "--format=%an <%ae>")); got != "Captain <captain@example.test>" {
		t.Fatalf("marker Git author = %q, want shared Captain credential", got)
	}
	if got := strings.TrimSpace(git(t, fixture.worktree, "show", "HEAD:owner-handoff.marker")); got != "runtime-worker-owner" {
		t.Fatalf("committed owner marker = %q, want runtime-worker-owner", got)
	}
	if got := strings.TrimSpace(git(t, fixture.worktree, "status", "--porcelain")); got != "" {
		t.Fatalf("registered owner worktree is not clean after committed handoff:\n%s", got)
	}
	rebaseDir := strings.TrimSpace(git(t, fixture.worktree, "rev-parse", "--git-path", "rebase-merge"))
	if _, err := os.Stat(rebaseDir); !os.IsNotExist(err) {
		t.Fatalf("rebase was not cleanly aborted: %v", err)
	}
	worktrees := strings.Split(strings.TrimSpace(git(t, fixture.root, "worktree", "list", "--porcelain")), "\n")
	var worktreePaths []string
	for _, line := range worktrees {
		if path, ok := strings.CutPrefix(line, "worktree "); ok {
			worktreePaths = append(worktreePaths, path)
		}
	}
	if len(worktreePaths) != 2 || canonicalPath(t, worktreePaths[0]) != canonicalPath(t, fixture.root) || canonicalPath(t, worktreePaths[1]) != canonicalPath(t, fixture.worktree) {
		t.Fatalf("worktree inventory = %q, want only root and stamped owner worktree", worktreePaths)
	}
	branches := strings.Fields(git(t, fixture.root, "branch", "--format=%(refname:short)"))
	if len(branches) != 2 || branches[0] != "main" || branches[1] != fixture.owner.Branch {
		t.Fatalf("branch inventory = %q, want only main and stamped owner branch", branches)
	}
	assertConflictOwnerFreshEnvelope(t, spacedockBinary(t), fixture.root, fixture.entity, fixture.owner)
}

func canonicalPath(t *testing.T, path string) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatalf("resolve worktree path %q: %v", path, err)
	}
	return resolved
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
