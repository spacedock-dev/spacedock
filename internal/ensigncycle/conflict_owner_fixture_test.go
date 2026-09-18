package ensigncycle

import (
	"encoding/json"
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
		Entity:       status.EntitySlug(entity),
		Stage:        fields["status"],
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

func TestConflictOwnerStampedIdentity(t *testing.T) {
	root := t.TempDir()
	entity := filepath.Join(root, "conflict-owner.md")
	writeFile(t, filepath.Join(root, "README.md"), conflictOwnerWorkflow())
	writeFile(t, entity, conflictOwnerEntity())
	gitInit(t, root)
	owner := stampConflictOwner(t, buildRecordedGateBinary(t), root, entity)

	for _, branches := range [][]string{{"main", owner.Branch}, {owner.Branch, "main"}} {
		if !exactOwnerBranches(branches, owner.Branch) {
			t.Fatalf("valid branch inventory rejected: %q", branches)
		}
	}
	for _, branches := range [][]string{{"main"}, {owner.Branch}, {"main", "wrong"}, {"main", owner.Branch, "extra"}} {
		if exactOwnerBranches(branches, owner.Branch) {
			t.Fatalf("invalid branch inventory accepted: %q", branches)
		}
	}
	if owner.Entity != "conflict-owner" || owner.Stage != "implementation" || owner.WorkerName != "conflict-owner-implementation" || owner.Branch != "conflict-owner" {
		t.Fatalf("stamped owner identity = %#v", owner)
	}
}

func exactOwnerBranches(branches []string, owner string) bool {
	return len(branches) == 2 && ((branches[0] == "main" && branches[1] == owner) || (branches[1] == "main" && branches[0] == owner))
}
