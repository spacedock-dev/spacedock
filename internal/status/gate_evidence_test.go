// ABOUTME: Gate-evidence bundle command behavior and fail-closed Git identity controls.
// ABOUTME: The fixture proves one read replaces discovery without touching either repository.
package status

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/testgit"
)

type gateEvidenceFixture struct{ definition, state, entity string }

func newGateEvidenceFixture(t *testing.T) gateEvidenceFixture {
	t.Helper()
	definition := t.TempDir()
	state := filepath.Join(definition, ".spacedock-state")
	writeFile(t, filepath.Join(definition, "README.md"), "---\nid-style: slug\nstate: .spacedock-state\nstages:\n  states:\n    - name: implementation\n      initial: true\n    - name: validation\n      gate: true\n---\n# Workflow\n\n### validation\n\nValidate retained evidence.\n")
	writeFile(t, filepath.Join(definition, "recorder-contract.md"), "# Recorder contract\n\nUse exact objects.\n")
	entity := filepath.Join(state, "recorded-gate-task", "index.md")
	writeFile(t, entity, "---\nid: recorded-gate-task\nstatus: validation\n---\n# Task\n\n## Stage Report: validation\n\n- DONE: retained\n  Evidence.\n")
	writeFile(t, filepath.Join(filepath.Dir(entity), "selected", "gate-review.md"), "# Review\n\nApprove.\n")
	writeFile(t, filepath.Join(filepath.Dir(entity), "selected", "entity-snapshot.md"), "# Snapshot\n\nExact.\n")
	writeFile(t, filepath.Join(state, "dirty-sibling.md"), "untracked sibling\n")
	testgit.InitRepo(t, definition)
	gateEvidenceTestGit(t, definition, "add", "README.md", "recorder-contract.md")
	gateEvidenceTestGit(t, definition, "commit", "-m", "definition evidence")
	testgit.InitRepo(t, state)
	gateEvidenceTestGit(t, state, "add", "recorded-gate-task")
	gateEvidenceTestGit(t, state, "commit", "-m", "entity evidence")
	return gateEvidenceFixture{definition, state, entity}
}

func gateEvidenceTestGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func runGateEvidenceCommand(t *testing.T, f gateEvidenceFixture, ref string) (string, string, int) {
	t.Helper()
	return runNative(t, f.definition, pinnedEnv(t), "--workflow-dir", f.definition, "--read", ref, "--gate-evidence", "--json")
}

func TestGateEvidenceReturnsCanonicalCommittedBundleWithoutMutation(t *testing.T) {
	f := newGateEvidenceFixture(t)
	beforeMain := gateEvidenceTestGit(t, f.definition, "status", "--porcelain=v1")
	beforeState := gateEvidenceTestGit(t, f.state, "status", "--porcelain=v1")
	out, errOut, code := runGateEvidenceCommand(t, f, "recorded-gate-task")
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, errOut)
	}
	var got struct {
		Command    string                       `json:"command"`
		Stage      struct{ Name, Bytes string } `json:"stage"`
		Entity     map[string]string            `json:"entity"`
		Candidates []map[string]string          `json:"candidates"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("parse: %v\n%s", err, out)
	}
	if got.Command != "gate-evidence" || got.Stage.Name != "validation" || !strings.Contains(got.Stage.Bytes, "Validate retained evidence.") {
		t.Fatalf("wrong command/stage: %+v", got)
	}
	want := map[string]string{"main:recorder-contract.md": "", "state:recorded-gate-task/selected/entity-snapshot.md": "", "state:recorded-gate-task/selected/gate-review.md": ""}
	for _, candidate := range got.Candidates {
		key := candidate["root"] + ":" + candidate["path"]
		if _, ok := want[key]; !ok || candidate["input_path"] == "" || candidate["object_id"] == "" || !strings.HasPrefix(candidate["digest"], "sha256:") || candidate["bytes"] == "" {
			t.Fatalf("unexpected or incomplete candidate: %v", candidate)
		}
		delete(want, key)
	}
	if len(want) != 0 || got.Entity["root"] != "state" || got.Entity["path"] != "recorded-gate-task/index.md" || strings.Contains(out, "dirty-sibling") {
		t.Fatalf("bundle mismatch missing=%v entity=%v\n%s", want, got.Entity, out)
	}
	if gateEvidenceTestGit(t, f.definition, "status", "--porcelain=v1") != beforeMain || gateEvidenceTestGit(t, f.state, "status", "--porcelain=v1") != beforeState {
		t.Fatal("gate evidence mutated a repository")
	}
}

func TestGateEvidenceFailsClosed(t *testing.T) {
	cases := []struct{ name, ref, want string }{
		{"dirty selected source identity mismatch", "recorded-gate-task", "differs from its committed Git object"},
		{"traversal", "../README.md", "task reference"}, {"non gate", "recorded-gate-task", "not a gate stage"},
		{"missing candidates", "recorded-gate-task", "no committed Markdown candidates"},
		{"duplicate candidates", "recorded-gate-task", "duplicate candidate name"}, {"swapped roots", "recorded-gate-task", "independent Git roots"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newGateEvidenceFixture(t)
			switch tc.name {
			case "dirty selected source identity mismatch":
				writeFile(t, filepath.Join(filepath.Dir(f.entity), "selected", "gate-review.md"), "dirty\n")
			case "non gate":
				writeFile(t, f.entity, strings.Replace(string(mustRead(t, f.entity)), "status: validation", "status: implementation", 1))
				commitGateEvidence(t, f.state, "recorded-gate-task/index.md")
			case "missing candidates":
				if err := os.Remove(filepath.Join(f.definition, "recorder-contract.md")); err != nil {
					t.Fatal(err)
				}
				commitGateEvidence(t, f.definition, "-u")
			case "duplicate candidates":
				writeFile(t, filepath.Join(filepath.Dir(f.entity), "other", "gate-review.md"), "# Other\n")
				commitGateEvidence(t, f.state, "recorded-gate-task/other/gate-review.md")
			case "swapped roots":
				if err := os.RemoveAll(filepath.Join(f.state, ".git")); err != nil {
					t.Fatal(err)
				}
				commitGateEvidence(t, f.definition, ".spacedock-state/recorded-gate-task")
			}
			_, errOut, code := runGateEvidenceCommand(t, f, tc.ref)
			if code == 0 || !strings.Contains(errOut, tc.want) {
				t.Fatalf("exit=%d stderr=%q", code, errOut)
			}
		})
	}
}

func commitGateEvidence(t *testing.T, root, path string) {
	t.Helper()
	gateEvidenceTestGit(t, root, "add", path)
	gateEvidenceTestGit(t, root, "commit", "-m", "fixture mutation")
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
